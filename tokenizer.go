package main

import "strconv"

type Tokenizer struct {
	vocab  []string
	atoi   map[string]int
	itoa   map[int]string
	merges [][2]string
}

func NewTokenizer(n int) *Tokenizer {
	return &Tokenizer{
		vocab:  make([]string, 0, n),
		atoi:   make(map[string]int, n),
		itoa:   make(map[int]string, n),
		merges: make([][2]string, 0, n-256),
	}
}

func (t *Tokenizer) Size() int {
	return cap(t.vocab)
}

func (t *Tokenizer) Init() {
	for i := range 256 {
		t.AddToken(strconv.Itoa(i))
	}
}

func (t *Tokenizer) AddToken(token string) {
	idx := len(t.vocab)

	t.vocab = append(t.vocab, token)

	t.atoi[token] = idx
	t.itoa[idx] = token
}

func (t *Tokenizer) AddMerge(left, right string) {
	t.merges = append(t.merges, [2]string{left, right})
}
