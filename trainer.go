package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type MBPETrainer struct {
	preTokenizer PreTokenizer
	model        *MBPE
	vocabSize    int
	dict         map[string]*Chunk
}

func NewMBPETrainer(preTokenizer PreTokenizer, model *MBPE, vocabSize int) *MBPETrainer {
	return &MBPETrainer{
		preTokenizer: preTokenizer,
		model:        model,
		vocabSize:    vocabSize,
		dict:         make(map[string]*Chunk),
	}
}

func (t *MBPETrainer) InitDict(name string) error {
	file, err := os.Open(name)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if err := scanner.Err(); err != nil {
			return err
		}

		compounds := t.preTokenizer.PreTokenize(line + "\n")

		for _, compound := range compounds {
			if _, ok := t.dict[compound]; !ok {
				t.dict[compound] = NewChunk(compound, 0, 0.1) // TODO move alpha into trainer ?!
			}

			t.dict[compound].n += 1
		}
	}

	return nil
}

func (t *MBPETrainer) Pairs(k int) [][2]string {
	pairs := make(map[[2]string]float64)

	for _, compound := range t.dict {
		for pair, weight := range compound.Pairs() {
			if _, ok := pairs[pair]; !ok {
				pairs[pair] = 0
			}

			pairs[pair] += weight
		}
	}

	pairsList := make([]struct {
		pair   [2]string
		weight float64
	}, 0, len(pairs))

	for pair, weight := range pairs {
		pairsList = append(pairsList, struct {
			pair   [2]string
			weight float64
		}{pair: pair, weight: weight})
	}

	sort.Slice(pairsList, func(i, j int) bool {
		if pairsList[i].weight != pairsList[j].weight {
			return pairsList[i].weight > pairsList[j].weight
		}

		if pairsList[i].pair[0] != pairsList[j].pair[0] {
			return t.model.atoi[pairsList[i].pair[0]] < t.model.atoi[pairsList[j].pair[0]]
		}

		return t.model.atoi[pairsList[i].pair[1]] < t.model.atoi[pairsList[j].pair[1]]
	})

	// for _, pair := range pairsList {
	// 	fmt.Printf("%s (%d) %s (%d) %f\n", pair.pair[0], t.model.atoi[pair.pair[0]], pair.pair[1], t.model.atoi[pair.pair[1]], pair.weight)
	// }

	k = min(k, len(pairsList))

	result := make([][2]string, k)

	for i := range k {
		result[i] = pairsList[i].pair

		// fmt.Printf("merging %s %s\n", result[i][0], result[i][1])
	}

	return result
}

func (t *MBPETrainer) AddToken(left, right string) {
	t.model.AddToken(left + right)
	t.model.AddMerge(left, right)
}

func (t *MBPETrainer) Merge(left, right string) {
	// fmt.Printf("merging %s and %s\n", left, right)

	for _, compound := range t.dict {
		pairs := compound.Pairs() // TODO don't recompute again

		if _, ok := pairs[[2]string{left, right}]; ok {
			compound.MergePair(left, right)
		}
	}
}

func (t *MBPETrainer) Train(name string) error {
	if err := t.InitDict(name); err != nil {
		return err
	}

	if err := t.SaveDict("dict.txt"); err != nil {
		return err
	}

	t.model.InitVocab(t.vocabSize)

	t.InitVocab()

	k := t.model.Cap() - t.model.Len()

	t.model.InitMerges(k)

	pb := NewProgressBar(60, k)

	fmt.Print(pb.String())

	for i := 0; i < k; i++ {
		pairs := t.Pairs(1)

		if len(pairs) == 0 {
			return nil
		}

		left, right := pairs[0][0], pairs[0][1]

		t.AddToken(left, right)

		t.Merge(left, right)

		pb.Increment()

		fmt.Print(pb.String())
	}

	return nil
}

func (t *MBPETrainer) InitVocab() {
	tokens := make(map[string]int)

	// TODO side effects of byte replacements on strings.Split(..., "")

	for _, chunk := range t.dict {
		for _, token := range chunk.Tokens() {
			if _, ok := tokens[token]; !ok {
				tokens[token] = 0

				// fmt.Printf("discovered new token %s %02x \n", token, []byte(token))
			}

			tokens[token]++
		}
	}

	alphabet := make([]string, 0, len(tokens))

	for token := range tokens {
		alphabet = append(alphabet, token)
	}

	sort.Strings(alphabet)

	for _, token := range alphabet {
		t.model.AddToken(token)
	}
}

func (t *MBPETrainer) SaveDict(name string) error {
	if err := toFile(name, func(writer *bufio.Writer) error {
		for _, chunk := range t.dict {
			if _, err := writer.WriteString(fmt.Sprintf("%s %d\n", chunk.src, chunk.n)); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
