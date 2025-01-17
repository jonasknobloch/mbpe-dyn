package main

import "fmt"

func Fertility(tokenizer *Tokenizer, dict *Dict) {
	numTokens := 0
	numChunks := 0

	for _, chunk := range dict.Items() {
		numTokens += len(tokenizer.Tokenize(chunk.src)) * chunk.n
		numChunks += chunk.n
	}

	f := float64(numTokens) / float64(numChunks)

	fmt.Println(f, "fertility")
}
