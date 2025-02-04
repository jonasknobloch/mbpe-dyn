package main

import (
	"errors"
	"fmt"
	"mbpe-dyn/bpr"
	"strings"
)

type MergeLayerEvaluator struct {
	gold [][]string
}

func NewMergeLayerEvaluator() *MergeLayerEvaluator {
	return &MergeLayerEvaluator{}
}

func (ml *MergeLayerEvaluator) LoadSegmentations(name string) error {
	return readTsv(name, func(record []string) error {
		if len(record) != 2 {
			return errors.New("unexpected number of fields")
		}

		ml.gold = append(ml.gold, append([]string{record[0]}, strings.Split(record[1], " ")...))

		return nil
	})
}

func (ml *MergeLayerEvaluator) Eval(tokenizer Tokenizer, maxRank int) ([]float64, error) {
	return ml.evalMorph(tokenizer, maxRank)
}

func (ml *MergeLayerEvaluator) evalMorph(tokenizer Tokenizer, maxRank int) ([]float64, error) {
	model := tokenizer.model.(*MBPE)

	sum := 0.0
	total := 0

	for _, split := range ml.gold {
		word := "Ġ" + split[0]
		gold := split[1:]

		layers := func() [][]string {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("unknown token in", word)
				}
			}()

			layers := model.TokenizeLayered(word, maxRank)

			result := make([][]string, len(layers))

			for i, layer := range layers {
				result[i] = model.ToString(layer)
			}

			return result
		}()

		if layers == nil {
			continue // skip words containing unknown tokens
		}

		duplicate := -1

		for i, tokens := range layers {
			if tokens[0] == "Ġ" {
				layers[i] = tokens[1:]
			} else if len(tokens[0]) > 1 && tokens[0][:len("Ġ")] == "Ġ" {
				layers[i][0] = tokens[0][len("Ġ"):]

				if duplicate == -1 {
					duplicate = i - 1
				}
			}
		}

		if duplicate != -1 {
			layers = append(layers[:duplicate], layers[duplicate+1:]...)
		}

		gV := bpr.BoundaryVector(gold)
		pV := make([][]int, len(layers))

		for i, layer := range layers {
			pV[i] = bpr.BoundaryVector(layer)
		}

		s := 0
		n := 0

		for i, v := range gV {
			if v == 0 {
				continue
			}

			for j, p := range pV {
				if p[i] == 0 {
					s += j

					break
				}

				if j == len(pV)-1 {
					s += len(pV) // bound not destroyed
				}
			}

			n++
		}

		if n == 0 {
			sum += 1
			total++

			continue
		}

		average := float64(s) / float64(n)

		normalized := average / float64(len(pV))

		if normalized > 1 {
			panic("invalid layer")
		}

		sum += normalized

		total++
	}

	return []float64{sum / float64(total)}, nil
}
