package main

type MBPE struct {
	normalizer   Normalizer
	preTokenizer PreTokenizer
	tokenizer    *Tokenizer
}

func NewMBPE() *MBPE {
	return &MBPE{
		normalizer:   NewDefaultNormalizer(),
		preTokenizer: NewPreTokenizer(),
		tokenizer:    NewTokenizer(),
	}
}

func (m *MBPE) Tokenize(phrase string) []int {
	phrase = m.normalizer.normalize(phrase)
	chunks := m.preTokenizer.preTokenize(phrase)

	r := make([]int, 0)

	for _, chunk := range chunks {
		r = append(r, m.tokenizer.Tokenize(chunk)...)
	}

	return r
}
