package main

import (
	"errors"
	"fmt"
	legacy "mbpe-dyn/bpr"
	"strings"
)

const (
	MaxPrecision = iota
	MaxRecall
	MaxAccuracy
	MaxF1
)

type BPREvaluator struct {
	skipSingletonGoldSegmentations bool
	skipSingletonTestSegmentations bool
	chooseBestTokenizationLayer    bool
	gold                           [][]string
}

func NewBPREvaluator() *BPREvaluator {
	return &BPREvaluator{
		skipSingletonGoldSegmentations: true,
		skipSingletonTestSegmentations: false,
		chooseBestTokenizationLayer:    false,
		gold:                           make([][]string, 0),
	}
}

func (bpr *BPREvaluator) LoadSegmentations(name string) error {
	return readTsv(name, func(record []string) error {
		if len(record) != 2 {
			return errors.New("unexpected number of fields")
		}

		bpr.gold = append(bpr.gold, append([]string{record[0]}, strings.Split(record[1], " ")...))

		return nil
	})
}

func (bpr *BPREvaluator) Eval(tokenizer Tokenizer, maxRank int) ([]float64, error) {
	return bpr.evalLegacy(tokenizer, maxRank)
}

func (bpr *BPREvaluator) evalLegacy(tokenizer Tokenizer, maxRank int) ([]float64, error) {
	gold := make([][]string, 0)
	pred := make([][]string, 0)

	model := tokenizer.model.(*MBPE)

	for _, split := range bpr.gold {
		if bpr.skipSingletonGoldSegmentations && len(split) == 2 {
			continue
		}

		word := "Ġ" + split[0]

		tokens := func() []string {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("unknown token in", word)
				}
			}()

			layers := model.TokenizeLayered(word, maxRank)

			return model.ToString(layers[len(layers)-1])
		}()

		if tokens == nil {
			continue // tokens = []string{word}
		}

		if tokens[0] == "Ġ" {
			tokens = tokens[1:]
		} else if len(tokens[0]) > 1 && tokens[0][:len("Ġ")] == "Ġ" {
			tokens[0] = tokens[0][len("Ġ"):]
		}

		if bpr.skipSingletonTestSegmentations && len(tokens) == 1 {
			continue
		}

		gold = append(gold, split[1:])
		pred = append(pred, tokens)
	}

	precision, recall, f1 := legacy.Eval(gold, pred)

	return []float64{precision, recall, f1}, nil
}

func (bpr *BPREvaluator) eval(tokenizer Tokenizer, maxRank int) ([]float64, error) {
	tp := 0
	fp := 0
	tn := 0
	fn := 0

	model := tokenizer.model.(*MBPE)

	for _, split := range bpr.gold {
		if bpr.skipSingletonGoldSegmentations && len(split) == 2 {
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

			layers := model.TokenizeLayered(word, maxRank)

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

		if bpr.skipSingletonTestSegmentations && len(layers[len(layers)-1]) == 1 {
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

	return []float64{precision, recall, accuracy, f1}, nil
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

func (bpr *BPREvaluator) EvalFile(name string) ([]float64, error) {
	gold := make([][]string, len(bpr.gold))
	pred := make([][]string, len(bpr.gold))

	i := 0

	if err := readTsv(name, func(record []string) error {
		if len(record) != 2 {
			return errors.New("unexpected number of fields")
		}

		if record[0] != bpr.gold[i][0] {
			return errors.New("unexpected compound")
		}

		gold[i] = bpr.gold[i][1:]
		pred[i] = strings.Split(record[1], " ")

		i++

		return nil
	}); err != nil {
		return nil, err
	}

	precision, recall, f1 := legacy.Eval(gold, pred)

	return []float64{precision, recall, f1}, nil
}
