package bpr

import "strings"

func RecallEvalSingle(pre, ref []int) (float64, int) {
	if len(pre) != len(ref) {
		panic("length missmatch")
	}

	sum := func(x []int) int {
		s := 0

		for _, v := range x {
			s += v
		}

		return s
	}

	total := sum(ref)

	if total == 0 {
		return 1.0, 0 // no bounds in reference
	}

	diff := make([]int, len(ref))

	for i, v := range ref {
		diff[i] = v - pre[i]
	}

	missed := make([]int, len(ref))

	for i := range ref {
		missed[i] = max(0, diff[i]) // (abs(diff[i]) + diff[i]) / 2
	}

	r := float64(total-sum(missed)) / float64(total)

	return r, total
}

func BoundaryVector(segmentation []string) []int {
	r := make([]int, 0)

	for i, s := range segmentation {
		for range s[:len(s)-1] {
			r = append(r, 0)
		}

		if i != len(segmentation)-1 {
			r = append(r, 1)
		}
	}

	if len(r) != len([]rune(strings.Join(segmentation, "")))-1 {
		panic("length missmatch")
	}

	return r
}

func Eval(gold, predicted [][]string) (float64, float64, float64) {
	sumPrecision := 0.0
	sumRecall := 0.0

	total := 0.0

	for i := range gold {
		a := BoundaryVector(gold[i])
		b := BoundaryVector(predicted[i])

		p, _ := RecallEvalSingle(a, b)
		r, _ := RecallEvalSingle(b, a)

		sumPrecision += p
		sumRecall += r

		total++
	}

	precision := sumPrecision / total
	recall := sumRecall / total

	f1 := 2.0 / (1.0/precision + 1.0/recall)

	return precision, recall, f1
}
