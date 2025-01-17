package main

import (
	"fmt"
	"log"
)

func main() {
	// tokenize()
	// eval()
	// train()
}

func eval() {
	initTokenizer := func(vocab, merges string) *Tokenizer {
		model := NewMBPE()

		if err := model.Load(vocab, merges); err != nil {
			log.Fatal(err)
		}

		tokenizer := NewTokenizer(model)

		byteLevel := NewByteLevel(true)

		tokenizer.SetPreTokenizer(byteLevel)
		tokenizer.SetDecoder(byteLevel)

		return tokenizer
	}

	runner := NewRunner()

	runner.AddTokenizer(initTokenizer("out/en-base/vocab.json", "out/en-base/merges.txt"), "en-base")
	runner.AddTokenizer(initTokenizer("out/en-c050/vocab.json", "out/en-c050/merges.txt"), "en-c050")
	runner.AddTokenizer(initTokenizer("out/en-c100/vocab.json", "out/en-c100/merges.txt"), "en-c100")
	runner.AddTokenizer(initTokenizer("out/en-c100-m050/vocab.json", "out/en-c100-m050/merges.txt"), "en-c100-m050")
	runner.AddTokenizer(initTokenizer("out/en-c100-m100/vocab.json", "out/en-c100-m100/merges.txt"), "en-c100-m100")

	runner.AddEvaluator(func() Evaluator {
		bprEval := NewBPREvaluator()

		if err := bprEval.LoadSegmentations("data/mbpe/goldstd_trainset.segmentation.eng.tsv"); err != nil {
			log.Fatal(err)
		}

		return bprEval
	}(), "BPR")

	runner.AddEvaluator(func() Evaluator {
		fertilityEval := NewFertilityEvaluator()

		if err := fertilityEval.InitDict("data/shakespeare.txt"); err != nil {
			log.Fatal(err)
		}

		return fertilityEval
	}(), "Fertility")

	runner.AddEvaluator(func() Evaluator {
		refEval := NewReferenceEvaluator()

		if err := refEval.LoadModel("out/en-base/vocab.json", "out/en-base/merges.txt"); err != nil {
			log.Fatal(err)
		}

		return refEval
	}(), "Overlap")

	fmt.Print(runner.RunAll())
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
