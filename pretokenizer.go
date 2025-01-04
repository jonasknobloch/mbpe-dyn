package main

import (
	"github.com/dlclark/regexp2"
)

type ByteLevel struct {
	addPrefixSpace bool
	re             *regexp2.Regexp
}

type PreTokenizer interface {
	PreTokenize(string) []string
}

func NewByteLevel(addPrefixSpace bool) *ByteLevel {
	re := regexp2.MustCompile(`'s|'t|'re|'ve|'m|'ll|'d| ?\pL+| ?\pN+| ?[^\s\pL\pN]+|\s+`, 0)

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

	regexp2FindAllString := func(re *regexp2.Regexp, s string) []string {
		var matches []string

		m, _ := re.FindStringMatch(s)

		for m != nil {
			matches = append(matches, m.String())

			m, _ = re.FindNextMatch(m)
		}

		return matches
	}

	compounds := regexp2FindAllString(p.re, phrase)

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
