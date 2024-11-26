package main

type MBPE struct {
	normalizer   Normalizer
	preTokenizer PreTokenizer
	tokenizer    *Tokenizer
}

func NewMBPE(n int) *MBPE {
	tokenizer := NewTokenizer(n)

	tokenizer.Init()

	return &MBPE{
		normalizer:   NewDefaultNormalizer(),
		preTokenizer: NewPreTokenizer(),
		tokenizer:    tokenizer,
	}
}

func (m *MBPE) Tokenize(phrase string) []string {
	phrase = m.normalizer.normalize(phrase)
	chunks := m.preTokenizer.preTokenize(phrase)

	r := make([]string, 0)

	for _, chunk := range chunks {
		r = append(r, m.tokenizer.Tokenize(chunk)...)
	}

	return r
}
