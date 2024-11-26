package main

import (
	"regexp"
)

type RegexpPreTokenizer struct {
	re *regexp.Regexp
}

type PreTokenizer interface {
	preTokenize(string) []string
}

func NewPreTokenizer() *RegexpPreTokenizer {
	// r"'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+"

	// TODO work around for missing negative lookahead

	re := regexp.MustCompile(`'s|'t|'re|'ve|'m|'ll|'d| ?\pL+| ?\pN+| ?[^\s\pL\pN]+|\s+`)

	return &RegexpPreTokenizer{
		re: re,
	}
}

func (p *RegexpPreTokenizer) preTokenize(phrase string) []string {
	compounds := p.re.FindAllString(phrase, -1)

	if phrase == "" {
		return []string{}
	}

	if compounds == nil {
		panic("could not match phrase")
	}

	return compounds
}
