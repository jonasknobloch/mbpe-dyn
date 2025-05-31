package mbpe

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
	remove bool
}

var InvertWeightFunction = false

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

func (c *Chunk) WeightedPairs() ([]Pair, []float64, float64) {
	return c.weightedPairs(InvertWeightFunction)
}

func (c *Chunk) weightedPairs(inverse bool) ([]Pair, []float64, float64) {
	pairs := c.Pairs()

	if len(pairs) == 0 {
		return pairs, []float64{}, 0.0
	}

	clashes := make([]bool, len(pairs))

	for i := 0; i < len(c.bounds)-2; i++ {
		lower := c.bounds[i]
		upper := c.bounds[i+2]

		for _, b := range c.morphs {
			if b > lower && b < upper {
				clashes[i] = true

				break
			}
		}
	}

	weights, epsilon := c.pairWeights(pairs, clashes, inverse)

	for i := range weights {
		weights[i] *= float64(c.n)
	}

	epsilon *= float64(c.n)

	return pairs, weights, epsilon
}

func (c *Chunk) pairWeights(pairs []Pair, clashes []bool, inverse bool) ([]float64, float64) {
	weights := make([]float64, len(pairs))

	n := float64(len(weights))
	k := 0.0

	for _, v := range clashes {
		if v != inverse {
			k++
		}
	}

	for i := range pairs {
		var w float64

		if clashes[i] != inverse {
			w = (1 - c.alpha) + (c.alpha * (k - 1) / n)
		} else {
			w = 1 + (c.alpha * k / n)
		}

		weights[i] = w
	}

	epsilon := c.alpha * k / n // no merge

	return weights, epsilon
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

func (c *Chunk) TrackedMerge(merge Merge) (map[Pair]Change, float64) {
	changes := make(map[Pair]Change)

	pairsBefore, weightsBefore, epsilonBefore := c.WeightedPairs()

	c.MergePair(merge.pair[0], merge.pair[1])

	pairsAfter, weightsAfter, epsilonAfter := c.WeightedPairs()

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
				remove: false,
			}
		} else {
			changes[pair] = Change{
				delta:  -weightBefore,
				update: false,
				remove: true,
			}
		}
	}

	for pair, weightAfter := range after {
		if _, ok := before[pair]; !ok {
			changes[pair] = Change{
				delta:  weightAfter,
				update: false,
				remove: false,
			}
		}
	}

	return changes, epsilonAfter - epsilonBefore
}

func (c *Chunk) Tokens() []string {
	r := make([]string, 0, len(c.bounds)-1)

	for i := 0; i < len(c.bounds)-1; i++ {
		r = append(r, c.src[c.bounds[i]:c.bounds[i+1]])
	}

	return r
}

// Deprecated: Use Inverter segmenter instead. Disable MergePrefixWhiteSpace
// to maintain behavior of Invert during training.
func (c *Chunk) Invert() {
	n := len(c.bounds) - len(c.morphs) - 2

	r := make([]int, 0, n)

	for _, b := range c.bounds[1 : len(c.bounds)-1] {
		found := false

		for _, m := range c.morphs {
			if b == m {
				found = true

				break
			}
		}

		if !found {
			r = append(r, b)
		}
	}

	if len(r) != n {
		panic("unexpected number of morphs")
	}

	c.morphs = r
}
