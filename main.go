package main

import (
	"fmt"
	stok "github.com/sugarme/tokenizer"
	sbpe "github.com/sugarme/tokenizer/model/bpe"
	spre "github.com/sugarme/tokenizer/pretokenizer"
	"log"
)

func main() {
	// tokenize()
	eval()
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

	bpr := NewBoundaryPrecisionRecall(false, false, true)

	if err := bpr.LoadDict("data/goldstd_trainset.segmentation.eng.tsv"); err != nil {
		log.Fatal(err)
	}

	bpr.Eval(tokenizer)

	Fertility(tokenizer)
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

	static := NewStatic(0.5)

	if err := static.LoadDict("data/en.splits.tsv"); err != nil {
		log.Fatal(err)
	}

	morfessor := NewMorfessor(0.5)

	if err := morfessor.LoadModel("data/morfessor/semisup_model.proto"); err != nil {
		log.Fatal(err)
	}

	segmenter := NewSequence(static, morfessor)

	trainer := NewMBPETrainer(preTokenizer, segmenter, model, 5000)

	if err := trainer.Train("data/shakespeare.txt"); err != nil {
		log.Fatal(err)
	}

	if err := model.Save("vocab.json", "merges.txt"); err != nil {
		log.Fatal(err)
	}

	if err := trainer.dict.Save("dict.txt"); err != nil {
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

	preTokenizer.SetTrimOffsets(false)

	tokenizer.WithPreTokenizer(preTokenizer)

	if err := tokenizer.Train(trainer, files); err != nil {
		log.Fatal(err)
	}

	result := tokenizer.GetModel()

	if err := result.Save("reference"); err != nil {
		log.Fatal(err)
	}
}
