package main

type MBPE struct {
	normalizer   *Normalizer
	preTokenizer *PreTokenizer
}

func NewMBPE() *MBPE {
	return &MBPE{
		normalizer:   NewNormalizer(),
		preTokenizer: NewPreTokenizer(),
	}
}
