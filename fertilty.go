package main

import "errors"

type FertilityEvaluator struct {
	dict *Dict
}

func NewFertilityEvaluator() *FertilityEvaluator {
	return &FertilityEvaluator{
		dict: NewDict(),
	}
}

func (f *FertilityEvaluator) InitDict(names ...string) error {
	f.dict.ProcessFiles(names...)

	return nil
}

func (f *FertilityEvaluator) LoadDict(name string) error {
	return f.dict.Load(name)
}

func (f *FertilityEvaluator) Eval(tokenizer *Tokenizer, maxRank int) ([]float64, error) {
	numTokens := 0
	numChunks := 0

	m, ok := tokenizer.model.(*MBPE)

	if !ok {
		return nil, errors.New("unexpected model type")
	}

	for _, chunk := range f.dict.Items() {
		numTokens += len(m.tokenize(chunk.src, nil, maxRank)) * chunk.n
		numChunks += chunk.n
	}

	fertility := float64(numTokens) / float64(numChunks)

	return []float64{fertility}, nil
}
