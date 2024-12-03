package main

type Tokenizer struct {
	preTokenizer PreTokenizer
	model        Model
	decoder      Decoder
}

func NewTokenizer(preTokenizer PreTokenizer, model Model, decoder Decoder) *Tokenizer {
	return &Tokenizer{
		preTokenizer: preTokenizer,
		model:        model,
		decoder:      decoder,
	}
}

func (t *Tokenizer) Tokenize(phrase string) []int {
	chunks := t.preTokenizer.PreTokenize(phrase)

	r := make([]int, 0)

	for _, chunk := range chunks {
		r = append(r, t.model.Tokenize(chunk)...)
	}

	return r
}
