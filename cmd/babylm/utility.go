package main

import "sort"

func mapToSortedSlice(m map[float64]struct{}) []float64 {
	s := make([]float64, 0, len(m))

	for k := range m {
		s = append(s, k)
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(s)))

	return s
}
