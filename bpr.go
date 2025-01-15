package main

import (
	"errors"
	"fmt"
	"strings"
)

// type Evaluator interface {
// 	Eval(tokenizer Tokenizer)
// }

const (
	MaxPrecision = iota
	MaxRecall
	MaxAccuracy
	MaxF1
)

type BoundaryPrecisionRecall struct {
	skipSingletonSegmentations  bool
	skipSingletonTokens         bool
	chooseBestTokenizationLayer bool
	maxRank                     int
	gold                        [][]string
}

func NewBoundaryPrecisionRecall(skipSingletonSegmentations, skipSingletonTokens, chooseBestTokenizationLayer bool, maxRank int) *BoundaryPrecisionRecall {
	return &BoundaryPrecisionRecall{
		skipSingletonSegmentations:  skipSingletonSegmentations,
		skipSingletonTokens:         skipSingletonTokens,
		chooseBestTokenizationLayer: chooseBestTokenizationLayer,
		maxRank:                     maxRank,
		gold:                        make([][]string, 0),
	}
}

func (bpr *BoundaryPrecisionRecall) LoadDict(name string) error {
	return readTsv(name, func(record []string) error {
		if len(record) != 2 {
			return errors.New("unexpected number of fields")
		}

		bpr.gold = append(bpr.gold, append([]string{record[0]}, strings.Split(record[1], " ")...))

		return nil
	})
}

func (bpr *BoundaryPrecisionRecall) Eval(tokenizer *Tokenizer) {
	tp := 0
	fp := 0
	tn := 0
	fn := 0

	model := tokenizer.model.(*MBPE)

	for _, split := range bpr.gold {
		if bpr.skipSingletonSegmentations && len(split) == 2 {
			continue
		}

		word := "Ġ" + split[0]
		gold := split[1:]

		layers := func() [][]string {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("unknown token in", word)
				}
			}()

			layers := model.TokenizeLayered(word, bpr.maxRank)

			if !bpr.chooseBestTokenizationLayer {
				layers = layers[len(layers)-1:]
			}

			result := make([][]string, len(layers))

			for i, layer := range layers {
				result[i] = model.ToString(layer)
			}

			return result
		}()

		if layers == nil {
			continue // skip words containing unknown tokens
		}

		for _, tokens := range layers {
			if tokens[0] == "Ġ" {
				tokens = tokens[1:]
			} else if len(tokens[0]) > 1 && tokens[0][:len("Ġ")] == "Ġ" {
				tokens[0] = tokens[0][len("Ġ"):]
			}
		}

		if bpr.skipSingletonTokens && len(layers[len(layers)-1]) == 1 {
			continue
		}

		_, counts := evalSegmentations(layers, gold, word, MaxF1)

		tp += counts[0]
		fp += counts[1]
		tn += counts[2]
		fn += counts[3]
	}

	precision := float64(tp) / float64(tp+fp)
	recall := float64(tp) / float64(tp+fn)
	accuracy := float64(tp+tn) / float64(tp+tn+fp+fn)
	f1 := 2 * precision * recall / (precision + recall)

	fmt.Println(precision, "precision")
	fmt.Println(recall, "recall")
	fmt.Println(accuracy, "accuracy")
	fmt.Println(f1, "f1")
}

func evalSegmentations(segmentations [][]string, gold []string, src string, mode int) ([]string, [4]int) {
	r := make([][4]int, len(segmentations))

	boundsSet := func(tokens []string) map[int]struct{} {
		b := make(map[int]struct{})

		i := 0

		for j, g := range tokens {
			if j == len(tokens)-1 {
				break
			}

			l := len(g)

			b[i+l] = struct{}{}
			i += l
		}

		return b
	}

	b := boundsSet(gold)

	for n, tokens := range segmentations {
		tp := 0
		fp := 0
		tn := 0
		fn := 0

		if tokens[0] == "Ġ" {
			tokens = tokens[1:]
		} else if len(tokens[0]) > 1 && tokens[0][:len("Ġ")] == "Ġ" {
			tokens[0] = tokens[0][len("Ġ"):]
		}

		if len(tokens) == 0 {
			panic("no tokens")
		}

		a := boundsSet(tokens)

		for i := range src {
			if i == 0 {
				continue
			}

			if i == len([]rune(src))-1 {
				break
			}

			_, inA := a[i]
			_, inB := b[i]

			if inA {
				if inB {
					tp++
				} else {
					fp++
				}
			} else {
				if inB {
					fn++
				} else {
					tn++
				}
			}
		}

		r[n] = [4]int{tp, fp, tn, fn}
	}

	best := float64(0)
	idx := -1

	for i, v := range r {
		tp := v[0]
		fp := v[1]
		tn := v[2]
		fn := v[3]

		precision := float64(tp) / float64(tp+fp)
		recall := float64(tp) / float64(tp+fn)
		accuracy := float64(tp+tn) / float64(tp+tn+fp+fn)
		f1 := 2 * precision * recall / (precision + recall)

		switch mode {
		case MaxPrecision:
			if precision >= best {
				best = precision
				idx = i
			}
		case MaxRecall:
			if recall >= best {
				best = recall
				idx = i
			}
		case MaxAccuracy:
			if accuracy >= best {
				best = accuracy
				idx = i
			}
		case MaxF1:
			if f1 >= best {
				best = f1
				idx = i
			}
		}
	}

	if idx == -1 {
		idx = len(r) - 1
	}

	// fmt.Println(src, ":", gold, ":", idx, ":", segmentations)

	return segmentations[idx], r[idx]
}
