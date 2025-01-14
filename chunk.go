package main

type Chunk struct {
	src    string
	n      int
	bounds []int
	morphs []int
	alpha  float64
}

type Change struct {
	delta  float64
	update bool
}

func NewChunk(src string, n int, splits []string, alpha float64) *Chunk {
	bounds := []int{0}

	for _, r := range src {
		j := bounds[len(bounds)-1] + len(string(r))

		bounds = append(bounds, j)
	}

	var morphs []int

	if len(splits) > 1 {
		morphs = make([]int, 0)

		i := 0

		for _, sub := range splits[:len(splits)-1] {
			i += len(sub)

			morphs = append(morphs, i)
		}
	}

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

func (c *Chunk) Split(segments []string) {
	var morphs []int

	if len(segments) > 1 {
		morphs = make([]int, 0)

		i := 0

		for _, sub := range segments[:len(segments)-1] {
			i += len(sub)

			morphs = append(morphs, i)
		}
	}

	c.morphs = morphs
}

func (c *Chunk) Alpha(alpha float64) {
	c.alpha = alpha
}

func (c *Chunk) Pairs() []Pair {
	pairs := make([]Pair, len(c.bounds)-2)

	for i := 0; i < len(c.bounds)-2; i++ {
		pairs[i] = Pair{
			c.src[c.bounds[i]:c.bounds[i+1]],
			c.src[c.bounds[i+1]:c.bounds[i+2]],
		}
	}

	return pairs
}

func (c *Chunk) WeightedPairs() ([]Pair, []float64) {
	pairs := c.Pairs()

	clashes := make([]bool, len(pairs))
	nClashes := 0

	for i := 0; i < len(c.bounds)-2; i++ {
		lower := c.bounds[i]
		upper := c.bounds[i+2]

		for _, b := range c.morphs {
			if b > lower && b < upper {
				clashes[i] = true
				nClashes++

				break
			}
		}
	}

	weights := make([]float64, len(pairs))

	for i := range pairs {
		var w float64

		if len(pairs) == 1 {
			if clashes[i] {
				w = 1 - c.alpha // TODO does this make sense panics otherwise lol
				// probbaly panics because no weight change ist reqsterd otherwise idk -> verify conditions
			} else {
				w = 1
			}

			weights[i] = w
			break
		}

		if clashes[i] {
			w = (1 - c.alpha) + (c.alpha * float64(nClashes-1) / float64(len(pairs)-1))
		} else {
			w = 1 + (c.alpha * float64(nClashes) / float64(len(pairs)-1))
		}

		weights[i] = w
	}

	sum := 0.0

	for _, w := range weights {
		sum += w
	}

	// if sum != float64(len(pairs)) {
	// 	panic("houston")
	// }

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

func (c *Chunk) TrackedMerge(merge Merge) map[Pair]Change {
	changes := make(map[Pair]Change)

	pairsBefore, weightsBefore := c.WeightedPairs()

	c.MergePair(merge.pair[0], merge.pair[1])

	pairsAfter, weightsAfter := c.WeightedPairs()

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

			changes[pair] = Change{
				delta:  weightAfter - weightBefore,
				update: true,
			}
		} else {
			changes[pair] = Change{
				delta:  -weightBefore,
				update: false, // remove
			}
		}
	}

	for pair, weightAfter := range after {
		if _, ok := before[pair]; !ok {
			changes[pair] = Change{
				delta:  weightAfter,
				update: false, // add
			}
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
