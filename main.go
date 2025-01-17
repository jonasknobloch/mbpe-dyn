package main

import (
	"fmt"
	"log"
)

func main() {
	// tokenize()
	// eval()
	// train()
	CompareStateToReference()
}

func eval() {
	mbpe := NewMBPE()

	if err := mbpe.Load("out/00/vocab.json", "out/00/merges.txt"); err != nil {
		log.Fatal(err)
	}

	tokenizer := NewTokenizer(mbpe)

	byteLevel := NewByteLevel(true)

	tokenizer.SetPreTokenizer(byteLevel)
	tokenizer.SetDecoder(byteLevel)

	bpr := NewBoundaryPrecisionRecall(false, false, true, -1)

	if err := bpr.LoadDict("data/goldstd_trainset.segmentation.eng.tsv"); err != nil {
		log.Fatal(err)
	}

	bpr.Eval(tokenizer)

	dict := NewDict()

	dict.ProcessFiles("data/shakespeare.txt")

	Fertility(tokenizer, dict)
}

func tokenize() {
	model := NewMBPE()

	err := model.Load("vocab.json", "merges.txt")

	if err != nil {
		log.Fatal(err)
	}

	tokenizer := NewTokenizer(model)

	byteLevel := NewByteLevel(true)

	tokenizer.SetPreTokenizer(byteLevel)
	tokenizer.SetDecoder(byteLevel)

	tokens := tokenizer.Tokenize("To infinity and beyond!")

	fmt.Println(tokens)
	fmt.Println(model.ToString(tokens))
	fmt.Println(tokenizer.decoder.Decode(model.ToString(tokens)))
}

func train() {
	model := NewMBPE()

	preTokenizer := NewByteLevel(true)

	// static := NewStatic(0.5)
	//
	// if err := static.LoadDict("data/en.splits.tsv"); err != nil {
	// 	log.Fatal(err)
	// }
	//
	// morfessor := NewMorfessor(0.5)
	//
	// if err := morfessor.LoadModel("data/morfessor/semisup_model.proto"); err != nil {
	// 	log.Fatal(err)
	// }

	segmenter := NewSequence()

	trainer := NewMBPETrainer(preTokenizer, segmenter, model, 6000)

	if err := trainer.InitDict("data/shakespeare.txt"); err != nil {
		log.Fatal(err)
	}

	trainer.Train()

	if err := model.Save("vocab.json", "merges.txt"); err != nil {
		log.Fatal(err)
	}

	if err := trainer.dict.Save("dict.txt"); err != nil {
		log.Fatal(err)
	}
}
