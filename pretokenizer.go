package main

import (
	"regexp"
)

type PreTokenizer struct {
	re *regexp.Regexp
}

// type PreTokenizer interface {
// 	preTokenize(string) []string
// }

func NewPreTokenizer() *PreTokenizer {
	// r"'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+"

	// TODO work around for missing negative lookahead

	re := regexp.MustCompile(`'s|'t|'re|'ve|'m|'ll|'d| ?\pL+| ?\pN+| ?[^\s\pL\pN]+|\s+`)

	return &PreTokenizer{
		re: re,
	}
}

func (p *PreTokenizer) preTokenize(phrase string) []string {
	compounds := p.re.FindAllString(phrase, -1)

	if phrase == "" {
		return []string{}
	}

	if compounds == nil {
		panic("could not match phrase")
	}

	return compounds
}
