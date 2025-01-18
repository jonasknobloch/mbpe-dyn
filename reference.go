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

func (r *ReferenceEvaluator) Eval(tokenizer Tokenizer, maxRank int) ([]float64, error) {
	var model *MBPE

	if m, ok := tokenizer.model.(*MBPE); !ok {
		return nil, errors.New("unexpected model type")
	} else {
		model = m
	}

	vocab := r.vocabOverlap(model, maxRank)
	merges := r.mergeOverlap(model, maxRank)

	return []float64{vocab, merges}, nil
}

func (r *ReferenceEvaluator) vocabOverlap(target *MBPE, maxRank int) float64 {
	sizeR := len(r.model.vocab)
	sizeT := len(target.vocab)

	if maxRank > -1 {
		sizeR = len(r.model.Alphabet()) + maxRank
		sizeT = len(target.Alphabet()) + maxRank
	}

	n := 0

	for i, token := range target.vocab {
		if i > sizeT-1 {
			break
		}

		j, ok := r.model.atoi[token]

		if !ok || j > sizeR-1 {
			continue
		}

		n++
	}

	return float64(n) / float64(sizeR)
}

func (r *ReferenceEvaluator) mergeOverlap(target *MBPE, maxRank int) float64 {
	n := 0

	for _, merge := range target.merges {
		if maxRank > -1 && target.ranks[merge] > maxRank-1 {
			continue
		}

		rank, ok := r.model.ranks[merge]

		if !ok || (maxRank > -1 && rank > maxRank-1) {
			continue
		}

		n++
	}

	if maxRank == -1 {
		return float64(n) / float64(len(r.model.merges))
	}

	return float64(n) / float64(maxRank)
}
