package main

import (
	"fmt"
	"testing"
)

func TestChunk_Pairs(t *testing.T) {
	c := NewChunk("Ġthth", 2, nil, 1)

	pairs, weights := c.Pairs()

	fmt.Println(pairs, weights)

	// TODO implement
}

func TestChunk_MergePairIdx(t *testing.T) {
	c := NewChunk("hello", 2, nil, 1)

	c.MergePairIdx(0)
	c.MergePairIdx(0)
	c.MergePairIdx(0)
	c.MergePairIdx(0)

	fmt.Println(c.bounds)

	// TODO implement
}

func TestChunk_MergePair(t *testing.T) {
	c := NewChunk("hello", 2, nil, 1)

	c.MergePair("h", "e")
	c.MergePair("he", "l")
	c.MergePair("hel", "l")
	c.MergePair("hell", "o")

	fmt.Println(c.bounds)

	// TODO implement
}

func TestChunk_TrackedMerge(t *testing.T) {
	c := NewChunk("Ġthither", 1, nil, 0)

	changes := c.TrackedMerge(Merge{
		pair:      Pair{"Ġ", "t"},
		idx:       [2]int{0, 1},
		weight:    0,
		positions: nil,
	})

	if len(changes) != 3 {
		t.Errorf("expected %d changes but got %d\n", 3, len(changes))
	}

	if delta, ok := changes[[2]string{"Ġ", "t"}]; !ok || delta != -1 {
		if !ok {
			t.Errorf("expected change not found\n")
		} else {
			t.Errorf("expected delta %d but got %d\n", -1, int(delta))
		}
	}

	if delta, ok := changes[[2]string{"Ġ", "t"}]; !ok || delta != -1 {
		if !ok {
			t.Errorf("expected change not found\n")
		} else {
			t.Errorf("expected delta %d but got %d\n", -1, int(delta))
		}
	}

	if delta, ok := changes[[2]string{"t", "h"}]; !ok || delta != -1 {
		if !ok {
			t.Errorf("expected change not found\n")
		} else {
			t.Errorf("expected delta %d but got %d\n", -1, int(delta))
		}
	}

	if delta, ok := changes[[2]string{"Ġt", "h"}]; !ok || delta != 1 {
		if !ok {
			t.Errorf("expected change not found\n")
		} else {
			t.Errorf("expected delta %d but got %d\n", 1, int(delta))
		}
	}
}
