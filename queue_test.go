package main

import (
	"container/heap"
	"testing"
)

func TestQueue_Pop(t *testing.T) {
	ab := Merge{
		pair:   Pair{"a", "b"},
		weight: 2,
	}

	bc := Merge{
		pair:   Pair{"b", "c"},
		weight: 3,
	}

	cd := Merge{
		pair:   Pair{"c", "d"},
		weight: 1,
	}

	pairs := []Merge{ab, bc, cd}

	q := NewQueue(pairs)

	top := heap.Pop(q).(Merge)

	if top.pair != bc.pair {
		t.Errorf("expected %v but got %v", bc, top)
	}

	l := q.Len()

	if l != 2 {
		t.Errorf("expected len %d but got %d", 2, l)
	}
}
