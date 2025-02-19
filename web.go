package main

import (
	_ "embed"
	"encoding/json"
	"log"
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
		if len(args) < 2 {
			log.Fatal("not enough arguments")
		}

		input := args[0].String()
		modelChoice := args[1].String()

		return WTokenize(input, modelChoice)
	})
}

func WTokenize(input string, modelChoice string) js.Value {
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

	result, err := WTokenizeWithSerialized(input, serialized)

	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(result)

	if err != nil {
		log.Fatal(err)
	}

	return js.ValueOf(string(jsonData))
}

func WTokenizeWithSerialized(input string, serialized []byte) ([]WChunk, error) {
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

	result := make([]WChunk, len(chunks))

	for i, chunk := range chunks {
		result[i] = WTokenizeChunk(model, chunk, -1)
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
