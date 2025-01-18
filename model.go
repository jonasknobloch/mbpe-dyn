package main

import (
	"bufio"
	"fmt"
	"strings"
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
	ranks  map[[2]string]int
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

	m.ranks = make(map[[2]string]int, n)
}

func (m *MBPE) Len() int {
	return len(m.vocab)
}

func (m *MBPE) Cap() int {
	return cap(m.vocab)
}

func (m *MBPE) AddToken(token string) {
	if _, ok := m.atoi[token]; ok {
		panic("token already exists")
	}

	idx := len(m.vocab)

	m.vocab = append(m.vocab, token)

	m.atoi[token] = idx
	m.itoa[idx] = token
}

func (m *MBPE) AddMerge(left, right string) {
	if _, ok := m.ranks[[2]string{left, right}]; ok {
		panic("merge already exists")
	}

	idx := len(m.merges)

	m.merges = append(m.merges, [2]string{left, right})

	m.ranks[[2]string{left, right}] = idx
}

func (m *MBPE) TokenizeLayered(phrase string, maxRank int) [][]int {
	merges := make([]int, 0)

	m.tokenize(phrase, &merges, maxRank)

	c := NewChunk(phrase, 1, nil, 0)

	r := make([][]int, len(merges)+1)

	idx := func(tokens []string) []int {
		layer := make([]int, len(tokens))

		for j, token := range tokens {
			idx, ok := m.atoi[token]

			if !ok {
				panic("unknown token")
			}

			layer[j] = idx
		}

		return layer
	}

	r[0] = idx(c.Tokens())

	for i, pos := range merges {
		c.MergePairIdx(pos)

		r[i+1] = idx(c.Tokens())
	}

	return r
}

func (m *MBPE) Tokenize(phrase string) []int {
	return m.tokenize(phrase, nil, -1)
}

func (m *MBPE) tokenize(phrase string, merges *[]int, maxRank int) []int {
	c := NewChunk(phrase, 1, nil, 0)

	for {
		pairs := c.Pairs()

		if len(pairs) == 0 {
			break
		}

		idx := -1
		rank := -1

		for i, pair := range pairs {
			r, ok := m.ranks[pair]

			if !ok {
				continue
			}

			if maxRank > -1 && r > maxRank-1 {
				continue
			}

			if idx == -1 || r < rank {
				idx = i
				rank = r
			}
		}

		if idx == -1 {
			break
		}

		if merges != nil {
			*merges = append(*merges, idx)
		}

		c.MergePairIdx(idx)
	}

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
	if err := toJSON(vocab, Vocab[string, int](m.atoi)); err != nil {
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

	first := true

	if err := fromFile(merges, func(scanner *bufio.Scanner) error {
		for scanner.Scan() {
			line := scanner.Text()

			if first {
				first = false

				if strings.HasPrefix(line, "#version:") {
					continue
				}
			}

			if err := scanner.Err(); err != nil {
				return err
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

	ranks := make(map[[2]string]int, len(ms))

	for i, merge := range ms {
		ranks[merge] = i
	}

	m.vocab = vs
	m.atoi = atoi
	m.itoa = itoa
	m.merges = ms
	m.ranks = ranks

	return nil
}
