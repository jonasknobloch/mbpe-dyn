package main

import "testing"

func TestCountLines(t *testing.T) {
	n, err := countLines("data/shakespeare.txt")

	if err != nil {
		t.Fatal(err)
	}

	expected := 2469

	if n != expected {
		t.Errorf("expected %d lines but got %d", expected, n)
	}
}
