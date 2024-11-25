package main

type Chunk struct {
	src    string
	n      int
	bounds []int
	morphs []int
	alpha  float64
}

func NewChunk(src string, n int, alpha float64) *Chunk {
	bounds := make([]int, len(src)+1)

	for i := range len(src) + 1 {
		bounds[i] = i
	}

	var morphs []int

	return &Chunk{
		src:    src,
		n:      n,
		bounds: bounds,
		morphs: morphs,
		alpha:  alpha,
	}
}

func (c *Chunk) Pairs() map[[2]string]float64 {
	pairs := make([][2]string, len(c.bounds)-2)
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

		pairs[i] = [2]string{
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

	result := make(map[[2]string]float64)

	for i, pair := range pairs {
		if _, ok := result[pair]; ok {
			result[pair] = 0
		}

		result[pair] += float64(c.n) * weights[i]
	}

	return result
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
