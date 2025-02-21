package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"syscall/js"
)

type WToken struct {
	Token string
	ID    int
}

type WChunk struct {
	Segmentations [][]WToken
}

//go:embed en-m000.gob
var m000 []byte

//go:embed en-m010.gob
var m010 []byte

//go:embed en-m020.gob
var m020 []byte

//go:embed en-m030.gob
var m030 []byte

//go:embed en-m040.gob
var m040 []byte

//go:embed en-m050.gob
var m050 []byte

//go:embed en-m060.gob
var m060 []byte

//go:embed en-m070.gob
var m070 []byte

//go:embed en-m080.gob
var m080 []byte

//go:embed en-m090.gob
var m090 []byte

//go:embed en-m100.gob
var m100 []byte

func wrapTokenizeWeb() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 3 {
			log.Fatal("not enough arguments")
		}

		input := args[0].String()
		modelChoice := args[1].String()
		vocabSize := args[2].Int()

		return WTokenize(input, modelChoice, vocabSize)
	})
}

func WTokenize(input string, modelChoice string, vocabSize int) js.Value {
	modelMapping := map[string][]byte{
		"m000": m000,
		"m010": m010,
		"m020": m020,
		"m030": m030,
		"m040": m040,
		"m050": m050,
		"m060": m060,
		"m070": m070,
		"m080": m080,
		"m090": m090,
		"m100": m100,
	}

	serialized, ok := modelMapping[modelChoice]

	if !ok {
		log.Fatal("invalid model choice")
	}

	result, err := WTokenizeWithSerialized(input, serialized, vocabSize)

	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(result)

	if err != nil {
		log.Fatal(err)
	}

	return js.ValueOf(string(jsonData))
}

func WTokenizeWithSerialized(input string, serialized []byte, vocabSize int) ([]string, error) {
	var model *MBPE

	if m, err := DeserializeModel(serialized); err != nil {
		return nil, err
	} else {
		model = m
	}

	tokenizer := NewTokenizer(model)

	byteLevel := NewByteLevel(true)

	tokenizer.SetPreTokenizer(byteLevel)
	tokenizer.SetDecoder(byteLevel)

	chunks := tokenizer.preTokenizer.PreTokenize(input)

	result := make([]string, len(chunks))

	maxRank := -1

	if vocabSize > -1 {
		alphabet := model.Alphabet()

		if !(vocabSize < len(alphabet)) && !(vocabSize > len(model.vocab)) {
			maxRank = vocabSize - len(alphabet)
		}
	}

	for i, chunk := range chunks {
		result[i] = WTokenizeToMermaid(model, chunk, maxRank)
	}

	return result, nil
}

func WTokenizeChunk(model *MBPE, chunk string, maxRank int) WChunk {
	layers := model.TokenizeLayered(chunk, maxRank)

	segmentations := make([][]WToken, len(layers))

	for i, layer := range layers {
		segmentations[i] = make([]WToken, len(layer))

		tokens := model.ToString(layer)

		for j, id := range layer {
			segmentations[i][j] = WToken{
				Token: tokens[j],
				ID:    id,
			}
		}
	}

	return WChunk{
		Segmentations: segmentations,
	}
}

func WTokenizeToMermaid(model *MBPE, chunk string, maxRank int) string {
	layers := model.TokenizeLayered(chunk, maxRank)

	diff := func(a, b []int) []int {
		seen := make(map[int]struct{})

		for _, v := range b {
			seen[v] = struct{}{}
		}

		var diff []int

		for _, v := range a {
			if _, ok := seen[v]; !ok {
				diff = append(diff, v)
			}
		}

		return diff
	}

	seen := make(map[int][2]int)

	nodes := make([][2]int, 0)
	edges := make([][2][2]int, 0)

	for i, layer := range layers {
		for j, id := range layer {
			if _, ok := seen[id]; ok {
				continue
			}

			node := [2]int{i, j}

			seen[id] = node
			nodes = append(nodes, node)
		}

		if i == 0 {
			continue
		}

		added := diff(layer, layers[i-1])
		removed := diff(layers[i-1], layer)

		if len(added) != 1 {
			panic("expected 1 but got ??")
		}

		if len(removed) != 2 {
			panic(fmt.Sprintf("expected 2 but got %v", removed))
		}

		edges = append(edges, [2][2]int{seen[removed[1]], seen[added[0]]})
		edges = append(edges, [2][2]int{seen[removed[0]], seen[added[0]]})
	}

	var sb strings.Builder

	sb.WriteString("graph TD\n")

	for i := len(nodes) - 1; i >= 0; i-- {
		node := nodes[i]

		sb.WriteString(fmt.Sprintf("%d%d[%s]\n", node[0], node[1], model.ToString([]int{layers[node[0]][node[1]]})[0]))
	}

	for i := len(edges) - 1; i >= 0; i-- {
		edge := edges[i]

		sb.WriteString(fmt.Sprintf("%d%d-%s->%d%d\n", edge[1][0], edge[1][1], strings.Repeat("-", edge[1][0]-edge[0][0]), edge[0][0], edge[0][1]))
	}

	return sb.String()
}
