package main

type MBPE struct {
	normalizer   *Normalizer
	preTokenizer *PreTokenizer
	tokenizer    *Tokenizer
}

func NewMBPE(n int) *MBPE {
	tokenizer := NewTokenizer(n)

	tokenizer.Init()

	return &MBPE{
		normalizer:   NewNormalizer(),
		preTokenizer: NewPreTokenizer(),
		tokenizer:    tokenizer,
	}
}
