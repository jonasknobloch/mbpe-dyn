package main

import (
	"fmt"
	stok "github.com/sugarme/tokenizer"
	sbpe "github.com/sugarme/tokenizer/model/bpe"
	spre "github.com/sugarme/tokenizer/pretokenizer"
	"log"
)

func main() {
	tokenize()
}

func tokenize() {
	model := NewMBPE()

	err := model.Load("vocab.json", "merges.txt")

	if err != nil {
		log.Fatal(err)
	}

	preTokenizer := NewByteLevel()

	tokenizer := NewTokenizer(preTokenizer, model, nil)

	tokens := tokenizer.Tokenize("To infinity and beyond!")

	fmt.Println(tokens)
	fmt.Println(model.ToString(tokens))
}

func train() {
	model := NewMBPE()

	preTokenizer := NewByteLevel()

	trainer := NewMBPETrainer(preTokenizer, model, 5000)

	if err := trainer.Train("data/shakespeare.txt"); err != nil {
		log.Fatal(err)
	}

	if err := model.Save("vocab.json", "merges.txt"); err != nil {
		log.Fatal(err)
	}
}

func trainReference() {
	files := []string{
		"data/shakespeare.txt",
	}

	var vocab = make(map[string]int)
	var merges = make(map[sbpe.Pair]sbpe.PairVal)

	model := sbpe.NewBPE(vocab, merges)

	trainer := sbpe.NewBpeTrainer(0, 5000)

	tokenizer := stok.NewTokenizer(model)

	preTokenizer := spre.NewByteLevel()

	preTokenizer.SetAddPrefixSpace(false)

	tokenizer.WithPreTokenizer(preTokenizer)

	if err := tokenizer.Train(trainer, files); err != nil {
		log.Fatal(err)
	}

	result := tokenizer.GetModel()

	if err := result.Save("reference"); err != nil {
		log.Fatal(err)
	}
}
