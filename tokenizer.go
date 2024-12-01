package main

type Tokenizer struct {
	vocab  []string
	atoi   map[string]int
	itoa   map[int]string
	merges [][2]string
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{}
}

func (t *Tokenizer) InitVocab(n int) {
	t.vocab = make([]string, 0, n)

	t.atoi = make(map[string]int, n)
	t.itoa = make(map[int]string, n)
}

func (t *Tokenizer) InitMerges(n int) {
	t.merges = make([][2]string, 0, n)
}

func (t *Tokenizer) Len() int {
	return len(t.vocab)
}

func (t *Tokenizer) Cap() int {
	return cap(t.vocab)
}

// func (t *Tokenizer) Init() {
// 	for _, c := range Alphabet() {
// 		t.AddToken(c)
// 	}
// }

func (t *Tokenizer) AddToken(token string) {
	idx := len(t.vocab)

	t.vocab = append(t.vocab, token)

	t.atoi[token] = idx
	t.itoa[idx] = token
}

func (t *Tokenizer) AddMerge(left, right string) {
	t.merges = append(t.merges, [2]string{left, right})
}

func (t *Tokenizer) Tokenize(phrase string) []int {
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

	r := make([]int, len(c.bounds)-1)

	for i, token := range c.Tokens() {
		idx, ok := t.atoi[token]

		if !ok {
			panic("unknown token")
		}

		r[i] = idx
	}

	return r
}

func (t *Tokenizer) ToString(tokens []int) []string {
	r := make([]string, len(tokens))

	for i, token := range tokens {
		s, ok := t.itoa[token]

		if !ok {
			panic("unknown token")
		}

		r[i] = s
	}

	return r
}
