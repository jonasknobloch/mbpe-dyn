package main

type Pair [2]string

type Merge struct {
	pair      Pair
	ids       [2]int
	weight    float64
	positions []int
}

func (p *Merge) Less(than Merge) bool {
	if p.weight != than.weight {
		return p.weight < than.weight
	}

	if p.ids[0] != than.ids[0] {
		return p.ids[0] > than.ids[0]
	}

	return p.ids[1] > than.ids[1]
}
