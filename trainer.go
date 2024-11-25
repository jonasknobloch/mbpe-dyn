package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Trainer struct {
	model *MBPE
	dict  map[string]*Chunk
}

func NewTrainer(n int) *Trainer {
	return &Trainer{
		model: NewMBPE(n),
		dict:  make(map[string]*Chunk),
	}
}

func (t *Trainer) Init(name string) error {
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

		compounds := t.model.preTokenizer.preTokenize(line)

		for _, compound := range compounds {
			if _, ok := t.dict[compound]; !ok {
				t.dict[compound] = NewChunk(compound, 0, 0.1) // TODO move alpha into trainer ?!
			}

			t.dict[compound].n += 1
		}
	}

	return nil
}

func (t *Trainer) Pairs(k int) [][2]string {
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
		return pairsList[i].weight > pairsList[j].weight
	})

	k = min(k, len(pairsList))

	result := make([][2]string, k)

	for i := range k {
		result[i] = pairsList[i].pair
	}

	return result
}

func (t *Trainer) AddToken(left, right string) {
	t.model.tokenizer.AddToken(left + right)
	t.model.tokenizer.AddMerge(left, right)
}

func (t *Trainer) Merge(left, right string) {
	// fmt.Printf("merging %s and %s\n", left, right)

	for _, compound := range t.dict {
		pairs := compound.Pairs() // TODO don't recompute again

		if _, ok := pairs[[2]string{left, right}]; ok {
			compound.MergePair(left, right)
		}
	}
}

func (t *Trainer) Train(name string) error {
	// k number of merges
	// base vocab is 256

	k := t.model.tokenizer.Size()

	if err := t.Init(name); err != nil {
		return err
	}

	for i := 0; i < k-256; i++ {
		pairs := t.Pairs(1)

		if len(pairs) == 0 {
			return nil
		}

		left, right := pairs[0][0], pairs[0][1]

		t.AddToken(left, right)

		t.Merge(left, right)

		fmt.Printf("%d\n", int(float64(i)/float64(k-256)*100))
	}

	return nil
}
