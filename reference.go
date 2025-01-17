package main

import "errors"

type ReferenceEvaluator struct {
	model *MBPE
}

func NewReferenceEvaluator() *ReferenceEvaluator {
	return &ReferenceEvaluator{
		NewMBPE(),
	}
}

func (r *ReferenceEvaluator) LoadModel(vocab, merges string) error {
	return r.model.Load(vocab, merges)
}

func (r *ReferenceEvaluator) Eval(tokenizer *Tokenizer) ([]float64, error) {
	m, ok := tokenizer.model.(*MBPE)

	if !ok {
		return nil, errors.New("unexpected model type")
	}

	vocab := vocabOverlap(m.atoi, r.model.atoi)
	merges := mergeOverlap(m.merges, r.model.merges)

	return []float64{vocab, merges}, nil
}

func vocabOverlap(a, b map[string]int) float64 {
	if len(a) != len(b) {
		// panic("vocabularies have different sizes")
	}

	n := 0

	for k := range a {
		if _, ok := b[k]; ok {
			n++
		}
	}

	// fmt.Println("\nmissed tokens")
	//
	// for k := range b {
	// 	if _, ok := a[k]; !ok {
	// 		fmt.Println(k)
	// 	}
	// }

	// fmt.Println("\nextra tokens")
	//
	// for k := range a {
	// 	if _, ok := b[k]; !ok {
	// 		fmt.Println(k)
	// 	}
	// }

	return float64(n) / float64(len(a))
}

func mergeOverlap(a, b [][2]string) float64 {
	if len(a) != len(b) {
		// panic("merge lists have different sizes")
	}

	n := 0

	for _, ma := range a {
		for _, mb := range b {
			if ma == mb {
				n++
				break
			}
		}
	}

	// fmt.Println("\nmissed merges")
	//
	// for _, mb := range b {
	// 	found := false
	//
	// 	for _, ma := range a {
	// 		if ma == mb {
	// 			found = true
	// 			break
	// 		}
	// 	}
	//
	// 	if !found {
	// 		fmt.Println(mb)
	// 	}
	// }

	// fmt.Println("\nextra merges")
	//
	// for _, ma := range a {
	// 	found := false
	//
	// 	for _, mb := range b {
	// 		if ma == mb {
	// 			found = true
	// 			break
	// 		}
	// 	}
	//
	// 	if !found {
	// 		fmt.Println(ma)
	// 	}
	// }

	return float64(n) / float64(len(a))
}
