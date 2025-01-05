package main

import (
	"github.com/dlclark/regexp2"
	"testing"
)

func TestFSA_FindAll(t *testing.T) {
	fsm := NewFSA()

	expected := []string{
		"Theirs",
		" for",
		" their",
		" style",
		" I",
		"'ll",
		" read",
		",",
		" his",
		" for",
		" his",
		" love",
		"'.",
		"\n",
	}

	for i, s := range fsm.FindAll("Theirs for their style I'll read, his for his love'.\n") {
		if s != expected[i] {
			t.Errorf("expected [%s] but got [%s]", expected[i], s)
		}
	}
}

func BenchmarkFSA_FindAll(b *testing.B) {
	fsm := NewFSA()

	for i := 0; i < b.N; i++ {
		fsm.FindAll("Theirs for their style I'll read, his for his love'.\n")
	}
}

func BenchmarkRegexp2_FindAll(b *testing.B) {
	regexp2FindAllString := func(re *regexp2.Regexp, s string) []string {
		var matches []string

		m, _ := re.FindStringMatch(s)

		for m != nil {
			matches = append(matches, m.String())

			m, _ = re.FindNextMatch(m)
		}

		return matches
	}

	re := regexp2.MustCompile(`'s|'t|'re|'ve|'m|'ll|'d| ?\pL+| ?\pN+| ?[^\s\pL\pN]+|\s+(?!\S)|\s+`, 0)

	for i := 0; i < b.N; i++ {
		regexp2FindAllString(re, "Theirs for their style I'll read, his for his love'.\n")
	}
}
