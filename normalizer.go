package main

type DefaultNormalizer struct {
	//
}

type Normalizer interface {
	normalize(string) string
}

func NewDefaultNormalizer() *DefaultNormalizer {
	return &DefaultNormalizer{}
}

func (n *DefaultNormalizer) normalize(phrase string) string {
	return phrase
}
