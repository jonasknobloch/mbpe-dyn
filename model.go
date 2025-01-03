package main

import (
	"bufio"
	"fmt"
)

type Model interface {
	Tokenize(string) []int
	Save(string, string) error
	Load(string, string) error
}

type MBPE struct {
	vocab  []string
	atoi   map[string]int
	itoa   map[int]string
	merges [][2]string
}

func NewMBPE() *MBPE {
	return &MBPE{}
}

func (m *MBPE) InitVocab(n int) {
	m.vocab = make([]string, 0, n)

	m.atoi = make(map[string]int, n)
	m.itoa = make(map[int]string, n)
}

func (m *MBPE) InitMerges(n int) {
	m.merges = make([][2]string, 0, n)
}

func (m *MBPE) Len() int {
	return len(m.vocab)
}

func (m *MBPE) Cap() int {
	return cap(m.vocab)
}

func (m *MBPE) AddToken(token string) {
	idx := len(m.vocab)

	m.vocab = append(m.vocab, token)

	m.atoi[token] = idx
	m.itoa[idx] = token
}

func (m *MBPE) AddMerge(left, right string) {
	m.merges = append(m.merges, [2]string{left, right})
}

func (m *MBPE) Tokenize(phrase string) []int {
	// pairs := make([]string, 0, len(chunk)-1)
	//
	// for i := 0; i < len(chunk)-1; i++ {
	// 	pairs = append(pairs, string(chunk[i]) + string(chunk[i+1]))
	// }

	c := NewChunk(phrase, 1, 0)

	var tokenize func()

	tokenize = func() {
		pairs := make([][2]string, len(c.bounds)-2)

		for i := 0; i < len(c.bounds)-2; i++ {
			pairs[i] = [2]string{
				c.src[c.bounds[i]:c.bounds[i+1]],
				c.src[c.bounds[i+1]:c.bounds[i+2]],
			}
		}

		if len(pairs) == 0 {
			return
		}

		for _, merge := range m.merges {
			for _, pair := range pairs {
				if merge == pair {
					c.MergePair(pair[0], pair[1])

					tokenize()
				}
			}
		}

		return
	}

	tokenize()

	r := make([]int, len(c.bounds)-1)

	for i, token := range c.Tokens() {
		idx, ok := m.atoi[token]

		if !ok {
			panic("unknown token")
		}

		r[i] = idx
	}

	return r
}

func (m *MBPE) ToString(tokens []int) []string {
	r := make([]string, len(tokens))

	for i, token := range tokens {
		s, ok := m.itoa[token]

		if !ok {
			panic("unknown token")
		}

		r[i] = s
	}

	return r
}

func (m *MBPE) Save(vocab, merges string) error {
	if err := toJSON(vocab, m.atoi); err != nil {
		return err
	}

	if err := toFile(merges, func(writer *bufio.Writer) error {
		for _, merge := range m.merges {
			if _, err := writer.WriteString(merge[0] + " " + merge[1] + "\n"); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (m *MBPE) Load(vocab, merges string) error {
	atoi := make(map[string]int)

	if err := fromJSON(vocab, &atoi); err != nil {
		return err
	}

	itoa := make(map[int]string, len(atoi))

	for token, idx := range atoi {
		itoa[idx] = token
	}

	vs := make([]string, len(itoa))

	for i := range len(itoa) {
		vs[i] = itoa[i]
	}

	ms := make([][2]string, 0) // unknown number of ms

	if err := fromFile(merges, func(scanner *bufio.Scanner) error {
		for scanner.Scan() {
			line := scanner.Text()

			if line[0] == '#' {
				continue
			}

			if err := scanner.Err(); err != nil {
			}

			var merge [2]string

			if _, err := fmt.Sscanf(line, "%s %s", &merge[0], &merge[1]); err != nil {
				return err
			}

			ms = append(ms, merge)
		}

		return nil
	}); err != nil {
		return err
	}

	m.vocab = vs
	m.atoi = atoi
	m.itoa = itoa
	m.merges = ms

	return nil
}
