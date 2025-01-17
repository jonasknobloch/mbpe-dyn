package main

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

func (f *FertilityEvaluator) Eval(tokenizer *Tokenizer) ([]float64, error) {
	numTokens := 0
	numChunks := 0

	for _, chunk := range f.dict.Items() {
		numTokens += len(tokenizer.model.Tokenize(chunk.src)) * chunk.n
		numChunks += chunk.n
	}

	fertility := float64(numTokens) / float64(numChunks)

	return []float64{fertility}, nil
}
