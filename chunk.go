package main

type Chunk struct {
	src    string
	n      int
	bounds []int
	morphs []int
	alpha  float64
}

func NewChunk(src string, n int, alpha float64) *Chunk {
	bounds := []int{0}

	for _, r := range src {
		j := bounds[len(bounds)-1] + len(string(r))

		bounds = append(bounds, j)
	}

	var morphs []int

	// suffixes := []string{"ing", "s", "ed"}
	//
	// for _, suffix := range suffixes {
	// 	if strings.HasSuffix(src, suffix) {
	// 		morphs = append(morphs, len(src)-len(suffix))
	// 	}
	// }

	return &Chunk{
		src:    src,
		n:      n,
		bounds: bounds,
		morphs: morphs,
		alpha:  alpha,
	}
}

func (c *Chunk) Pairs() ([]Pair, []float64) {
	pairs := make([]Pair, len(c.bounds)-2)
	clashes := make([]bool, len(c.bounds)-2)

	for i := 0; i < len(c.bounds)-2; i++ {
		lower := c.bounds[i]
		upper := c.bounds[i+2]

		for _, b := range c.morphs {
			if b > lower && b < upper {
				clashes[i] = true

				break
			}
		}

		pairs[i] = Pair{
			c.src[c.bounds[i]:c.bounds[i+1]],
			c.src[c.bounds[i+1]:c.bounds[i+2]],
		}
	}

	weights := make([]float64, len(pairs))

	nclashes := func() float64 {
		r := 0

		for _, b := range clashes {
			if b {
				r++
			}
		}

		return float64(r)
	}()

	for i := range pairs {
		sum := float64(1)

		if clashes[i] {
			sum -= c.alpha
		}

		sum += nclashes * c.alpha / float64(len(pairs))

		weights[i] = sum
	}

	for i := range weights {
		weights[i] *= float64(c.n)
	}

	return pairs, weights
}

func (c *Chunk) MergePairIdx(i int) {
	if i > len(c.bounds)-2 {
		panic("merge out of bounds")
	}

	c.bounds = append(c.bounds[:i+1], c.bounds[i+2:]...)
}

func (c *Chunk) MergePair(left, right string) {
	for i := 0; i < len(c.bounds)-2; i++ {
		l := c.src[c.bounds[i]:c.bounds[i+1]]
		r := c.src[c.bounds[i+1]:c.bounds[i+2]]

		if l == left && r == right {
			c.MergePairIdx(i)
			c.MergePair(left, right)

			return
		}
	}
}

func (c *Chunk) TrackedMerge(merge Merge) map[Pair]float64 {
	changes := make(map[Pair]float64)

	pairsBefore, weightsBefore := c.Pairs()

	c.MergePair(merge.pair[0], merge.pair[1])

	pairsAfter, weightsAfter := c.Pairs()

	before := make(map[Pair]float64)
	after := make(map[Pair]float64)

	for i, pair := range pairsBefore {
		before[pair] += weightsBefore[i]
	}

	for i, pair := range pairsAfter {
		after[pair] += weightsAfter[i]
	}

	for pair, weightBefore := range before {
		if weightAfter, ok := after[pair]; ok {
			if weightBefore == weightAfter {
				continue
			}

			changes[pair] = weightAfter - weightBefore // changed weight
		} else {
			changes[pair] = -weightBefore // removed pair
		}
	}

	for pair, weightAfter := range after {
		if _, ok := before[pair]; !ok {
			changes[pair] = weightAfter // new pair
		}
	}

	return changes
}

func (c *Chunk) Tokens() []string {
	r := make([]string, 0, len(c.bounds)-1)

	for i := 0; i < len(c.bounds)-1; i++ {
		r = append(r, c.src[c.bounds[i]:c.bounds[i+1]])
	}

	return r
}
