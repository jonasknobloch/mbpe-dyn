package main

import (
	"regexp"
)

type ByteLevel struct {
	addPrefixSpace bool
	re             *regexp.Regexp
}

type PreTokenizer interface {
	PreTokenize(string) []string
}

func NewByteLevel(addPrefixSpace bool) *ByteLevel {
	// r"'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+"

	// TODO work around for missing negative lookahead

	re := regexp.MustCompile(`'s|'t|'re|'ve|'m|'ll|'d| ?\pL+| ?\pN+| ?[^\s\pL\pN]+|\s+`)

	return &ByteLevel{
		addPrefixSpace: addPrefixSpace,
		re:             re,
	}
}

func (p *ByteLevel) PreTokenize(phrase string) []string {
	if phrase == "" {
		return []string{}
	}

	if p.addPrefixSpace && phrase[0] != ' ' {
		phrase = " " + phrase
	}

	compounds := p.re.FindAllString(phrase, -1)

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

func (p *ByteLevel) Decode(tokens []string) string {
	phrase := ""

	for _, token := range tokens {
		r := make([]byte, 0)

		for _, c := range token {
			r = append(r, CharBytes[string(c)])
		}

		phrase += string(r)
	}

	// never remove prefix space since we have no way of knowing if a prefix space
	// was added during pre-tokenization or if it was part of the original string

	// if p.addPrefixSpace && phrase[0] == ' ' {
	// 	phrase = phrase[1:]
	// }

	return phrase
}
