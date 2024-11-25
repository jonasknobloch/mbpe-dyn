package main

import (
	"fmt"
	"testing"
)

func TestChunk_Pairs(t *testing.T) {
	c := NewChunk("hello", 2, 1)

	pairs := c.Pairs()

	fmt.Println(pairs)

	// TODO implement
}

func TestChunk_MergePairIdx(t *testing.T) {
	c := NewChunk("hello", 2, 1)

	c.MergePairIdx(0)
	c.MergePairIdx(0)
	c.MergePairIdx(0)
	c.MergePairIdx(0)

	fmt.Println(c.bounds)

	// TODO implement
}

func TestChunk_MergePair(t *testing.T) {
	c := NewChunk("hello", 2, 1)

	c.MergePair("h", "e")
	c.MergePair("he", "l")
	c.MergePair("hel", "l")
	c.MergePair("hell", "o")

	fmt.Println(c.bounds)

	// TODO implement
}
