package main

import (
	"bufio"
	"fmt"
	"log"
)

func Fertility(tokenizer *Tokenizer) {
	// dict := NewDict()

	// dict.ProcessFiles("data/culturax/fi_part_00000.txt")
	//
	// if err := dict.Load("out/en-base/dict.txt"); err != nil {
	// 	log.Fatal(err)
	// }
	//
	// // TODO use dict load
	//
	// numTokens := 0
	//
	// for _, chunk := range dict.Items() {
	// 	numTokens += len(tokenizer.Tokenize(chunk.src))
	// }
	//
	// f := float64(numTokens) / float64(len(dict.Items()))
	//
	// fmt.Println(f, "fertility")

	nTokens := 0
	nChunks := 0

	if err := fromFile("data/shakespeare.txt", func(scanner *bufio.Scanner) error {
		for scanner.Scan() {
			line := scanner.Text()

			if err := scanner.Err(); err != nil {
				return err
			}

			chunks := tokenizer.preTokenizer.PreTokenize(line)

			nChunks += len(chunks)

			for _, chunk := range chunks {
				tokens := tokenizer.model.Tokenize(chunk)

				nTokens += len(tokens)
			}
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	f := float64(nTokens) / float64(nChunks)

	fmt.Println(f, "fertility")
}
