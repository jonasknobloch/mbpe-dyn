package main

import (
	"fmt"
	"testing"
)

func TestMerge_Less(t *testing.T) {
	a := Merge{
		idx:    [2]int{0, 1},
		weight: 1,
	}

	b := Merge{
		idx:    [2]int{0, 0},
		weight: 1,
	}

	fmt.Println(a.Less(b), true)

	// TODO implement
}
