package main

import (
	"fmt"
	stok "github.com/sugarme/tokenizer"
	sbpe "github.com/sugarme/tokenizer/model/bpe"
	spre "github.com/sugarme/tokenizer/pretokenizer"
)

func CompareStateToReference() error {
	modA := NewMBPE()

	if err := modA.Load("vocab.json", "merges.txt"); err != nil {
		return err
	}

	modB := NewMBPE()

	if err := modB.Load("reference/vocab.json", "reference/merges.txt"); err != nil {
		return err
	}

	fmt.Printf("\nvocab overlap: %f\n", VocabOverlap(modA.atoi, modB.atoi))
	fmt.Printf("\nmerge overlap: %f\n", MergeOverlap(modA.merges, modB.merges))

	return nil
}

func VocabOverlap(a, b map[string]int) float64 {
	if len(a) != len(b) {
		// panic("vocabularies have different sizes")
	}

	n := 0

	for k := range a {
		if _, ok := b[k]; ok {
			n++
		}
	}

	fmt.Println("\nmissed tokens")

	for k := range b {
		if _, ok := a[k]; !ok {
			fmt.Println(k)
		}
	}

	fmt.Println("\nextra tokens")

	for k := range a {
		if _, ok := b[k]; !ok {
			fmt.Println(k)
		}
	}

	return float64(n) / float64(len(a))
}

func MergeOverlap(a, b [][2]string) float64 {
	if len(a) != len(b) {
		// panic("merge lists have different sizes")
	}

	n := 0

	for _, ma := range a {
		for _, mb := range b {
			if ma == mb {
				n++
				break
			}
		}
	}

	fmt.Println("\nmissed merges")

	for _, mb := range b {
		found := false

		for _, ma := range a {
			if ma == mb {
				found = true
				break
			}
		}

		if !found {
			fmt.Println(mb)
		}
	}

	fmt.Println("\nextra merges")

	for _, ma := range a {
		found := false

		for _, mb := range b {
			if ma == mb {
				found = true
				break
			}
		}

		if !found {
			fmt.Println(ma)
		}
	}

	return float64(n) / float64(len(a))
}

func CompareTokenizerToReference() error {
	a, err := Tokenize("To infinity and beyond!")

	if err != nil {
		return err
	}

	b, err := TokenizeReference("To infinity and beyond!")

	if err != nil {
		return err
	}

	fmt.Printf("mbpe: %v\n", a)
	fmt.Printf("reference: %v\n", b)

	return nil
}

func Tokenize(sequence string) ([]string, error) {
	model := NewMBPE()

	err := model.Load("vocab.json", "merges.txt")

	if err != nil {
		return nil, err
	}

	preTokenizer := NewByteLevel(true)

	tokenizer := NewTokenizer(model)

	tokenizer.SetPreTokenizer(preTokenizer)

	tokens := tokenizer.Tokenize(sequence)

	return model.ToString(tokens), nil
}

func TokenizeReference(sequence string) ([]string, error) {
	model, err := sbpe.NewBpeFromFiles("vocab.json", "merges.txt")

	if err != nil {
		return nil, err
	}

	tokenizer := stok.NewTokenizer(model)

	// tokenizer.WithPreTokenizer(spre.NewByteLevel())

	pretokenizer := spre.NewByteLevel()

	pretokenized, err := pretokenizer.PreTokenize(stok.NewPreTokenizedString(sequence))

	if err != nil {
		return nil, err
	}

	splits := pretokenized.GetSplits(0, stok.Byte)

	input := make([]string, len(splits))

	for i, s := range splits {
		input[i] = s.Value
	}

	encoding, err := tokenizer.Encode(stok.NewSingleEncodeInput(stok.NewInputSequence(input)), false)

	if err != nil {
		return nil, err
	}

	return encoding.Tokens, nil
}
