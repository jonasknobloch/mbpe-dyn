package main

import "fmt"

func Fertility(tokenizer *Tokenizer) {
	dict := NewDict()

	dict.ProcessFiles("data/shakespeare.txt")

	numTokens := 0

	for _, chunk := range dict.Items() {
		numTokens += len(tokenizer.Tokenize(chunk.src))
	}

	f := float64(numTokens) / float64(len(dict.Items()))

	fmt.Println(f, "fertility")
}
