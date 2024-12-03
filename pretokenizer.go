package main

import (
	"regexp"
)

type ByteLevel struct {
	re *regexp.Regexp
}

type PreTokenizer interface {
	PreTokenize(string) []string
}

func NewByteLevel() *ByteLevel {
	// r"'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+"

	// TODO work around for missing negative lookahead

	re := regexp.MustCompile(`'s|'t|'re|'ve|'m|'ll|'d| ?\pL+| ?\pN+| ?[^\s\pL\pN]+|\s+`)

	return &ByteLevel{
		re: re,
	}
}

func (p *ByteLevel) PreTokenize(phrase string) []string {
	compounds := p.re.FindAllString(phrase, -1)

	if phrase == "" {
		return []string{}
	}

	if compounds == nil {
		panic("could not match phrase")
	}

	for i, compound := range compounds {
		r := ""

		for _, b := range []byte(compound) {
			r += BytesChar[b]
		}

		compounds[i] = r
	}

	return compounds
}
