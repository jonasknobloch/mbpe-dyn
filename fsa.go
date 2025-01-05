package main

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	RuneWhitespace = iota
	RuneLetter
	RuneNumber
	RuneOther
	StateInitial
	StateWhitespace
	StateLetter
	StateNumber
	StateOther
	StateWhitespaceLookAhead
)

type FSA struct {
	state  int
	input  string
	static []string
}

func NewFSA() *FSA {
	return &FSA{
		state:  StateInitial,
		input:  "",
		static: []string{"'s", "'t", "'re", "'m", "'ll", "'d"},
	}
}

func (f *FSA) Reset() {
	f.state = StateInitial
	f.input = ""
}

func (f *FSA) Read(next rune) bool {
	var r int

	if unicode.IsSpace(next) {
		r = RuneWhitespace
	} else if unicode.IsLetter(next) {
		r = RuneLetter
	} else if unicode.IsNumber(next) {
		r = RuneNumber
	} else {
		r = RuneOther
	}

	if f.state == StateInitial {
		switch r {
		case RuneWhitespace:
			f.input += string(next)
			f.state = StateWhitespace

			break
		case RuneLetter:
			f.input += string(next)
			f.state = StateLetter

			break
		case RuneNumber:
			f.input += string(next)
			f.state = StateNumber

			break
		default:
			f.input += string(next)
			f.state = StateOther
		}
	} else if f.state == StateWhitespace {
		switch r {
		case RuneWhitespace:
			f.input += string(next)
			f.state = StateWhitespaceLookAhead

			break
		case RuneLetter:
			f.input += string(next)
			f.state = StateLetter

			break
		case RuneNumber:
			f.input += string(next)
			f.state = StateNumber

			break
		default:
			f.input += string(next)
			f.state = StateOther
		}
	} else if f.state == StateNumber {
		switch r {
		case RuneNumber:
			f.input += string(next)
			f.state = StateNumber

			break
		default:
			return false
		}
	} else if f.state == StateLetter {
		switch r {
		case RuneLetter:
			f.input += string(next)
			f.state = StateLetter

			break
		default:
			return false
		}
	} else if f.state == StateOther {
		switch r {
		case RuneOther:
			f.input += string(next)
			f.state = StateOther

			break
		default:
			return false
		}
	} else if f.state == StateWhitespaceLookAhead {
		switch r {
		case RuneWhitespace:
			f.input += string(next)
			f.state = StateWhitespaceLookAhead

			break
		default:
			return false
		}
	} else {
		panic("invalid state")
	}

	return true
}

func (f *FSA) FindAll(s string) []string {
	var findAll func(runes []rune, matches []string) []string

	findAll = func(runes []rune, matches []string) []string {
		s = string(runes)

		for _, v := range f.static {
			if strings.HasPrefix(s, v) {
				matches = append(matches, v)
				runes = runes[utf8.RuneCountInString(v):]

				if len(runes) == 0 {
					return matches
				}

				break
			}
		}

		for i, r := range runes {
			ok := f.Read(r)

			if !ok {
				if f.state == StateInitial {
					return matches
				}

				if f.state == StateWhitespaceLookAhead {
					matches = append(matches, f.input[:len(f.input)-1])

					f.Reset()

					return findAll(runes[i-1:], matches)
				}

				matches = append(matches, f.input)

				f.Reset()

				return findAll(runes[i:], matches)
			}
		}

		matches = append(matches, f.input)

		return matches
	}

	defer f.Reset()

	return findAll([]rune(s), make([]string, 0))
}
