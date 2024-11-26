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

func (t *Tokenizer) Tokenize(phrase string) []string {
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

		for _, merge := range t.merges {
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

	return c.Tokens()

	// r := make([]int, len(c.bounds)-1)
	//
	// for _, token := range c.Tokens() {
	// 	idx, ok := t.atoi[token]
	//
	// 	if !ok {
	// 		panic("unknown token")
	// 	}
	//
	// 	r = append(r, idx)
	// }
	//
	// return r
}
