package main

type Pair [2]string

type Merge struct {
	pair      Pair
	idx       [2]int
	weight    float64
	positions []int
}

func (p *Merge) Less(than Merge) bool {
	if p.weight != than.weight {
		return p.weight < than.weight
	}

	if p.idx[0] != than.idx[0] {
		return p.idx[0] > than.idx[0]
	}

	return p.idx[1] > than.idx[1]
}
