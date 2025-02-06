package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	runner.AddTokenizer(*initTokenizer("out/en-base/vocab.json", "out/en-base/merges.txt"), "en-base")
	runner.AddTokenizer(*initTokenizer("out/en-s050/vocab.json", "out/en-s050/merges.txt"), "en-s050")
	runner.AddTokenizer(*initTokenizer("out/en-s100/vocab.json", "out/en-s100/merges.txt"), "en-s100")
	runner.AddTokenizer(*initTokenizer("out/en-m050/vocab.json", "out/en-m050/merges.txt"), "en-m050")
	runner.AddTokenizer(*initTokenizer("out/en-m100/vocab.json", "out/en-m100/merges.txt"), "en-m100")
	runner.AddTokenizer(*initTokenizer("out/en-s100-m050/vocab.json", "out/en-s100-m050/merges.txt"), "en-s100-m050")
	runner.AddTokenizer(*initTokenizer("out/en-s100-m100/vocab.json", "out/en-s100-m100/merges.txt"), "en-s100-m100")

	runner.AddEvaluator(func() Evaluator {
		bprEval := NewBPREvaluator()

		if err := bprEval.LoadSegmentations("data/mbpe/goldstd_trainset.segmentation.eng.tsv"); err != nil {
			log.Fatal(err)
		}

		return bprEval
	}(), "Boundary Precision Recall")

	runner.AddEvaluator(func() Evaluator {
		mlEval := NewMergeLayerEvaluator()

		if err := mlEval.LoadSegmentations("data/mbpe/goldstd_trainset.segmentation.eng.tsv"); err != nil {
			log.Fatal(err)
		}

		return mlEval
	}(), "Merge Layer")

	runner.AddEvaluator(func() Evaluator {
		fertilityEval := NewFertilityEvaluator()

		if err := fertilityEval.InitDict("data/culturax/en_part_00001-10k.txt"); err != nil {
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
	}(), "Reference Overlap")

	fmt.Print(runner.RunAll(1 << 15))
}

func tokenize() {
	model := NewMBPE()

	err := model.Load("out/en-base/vocab.json", "out/en-base/merges.txt")

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

func tokenizeFile(name string, vocabSize int) {
	model := NewMBPE()

	if err := model.Load("out/en-base/vocab.json", "out/en-base/merges.txt"); err != nil {
		log.Fatal(err)
	}

	gold := make([]string, 0)

	if err := readTsv(name, func(record []string) error {
		if len(record) == 0 {
			return errors.New("unexpected number of fields")
		}

		gold = append(gold, record[0])

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	segmentations := make([][]string, len(gold))

	maxRank := -1

	if vocabSize > -1 {
		maxRank = vocabSize - len(model.Alphabet())
	}

	for i, compound := range gold {
		tokens := model.ToString(model.tokenize(compound, nil, maxRank))

		if tokens[0] == "Ġ" {
			tokens = tokens[1:]
		} else if len(tokens[0]) > 1 && tokens[0][:len("Ġ")] == "Ġ" {
			tokens[0] = tokens[0][len("Ġ"):]
		}

		segmentations[i] = tokens
	}

	if err := toFile("segmentations.txt", func(writer *bufio.Writer) error {
		for i, segmentation := range segmentations {
			if _, err := writer.WriteString(fmt.Sprintf("%s\t%s\n", gold[i], strings.Join(segmentation, " "))); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func train() {
	out := "out"

	s050 := NewStatic(0.5)

	if err := s050.LoadDict("data/webcelex/en.splits.tsv"); err != nil {
		log.Fatal(err)
	}

	s100 := NewStatic(1)

	if err := s100.LoadDict("data/webcelex/en.splits.tsv"); err != nil {
		log.Fatal(err)
	}

	m050 := NewMorfessor(0.5)

	if err := m050.LoadModel("data/morfessor/semisup_model.proto"); err != nil {
		log.Fatal(err)
	}

	m100 := NewMorfessor(1)

	if err := m100.LoadModel("data/morfessor/semisup_model.proto"); err != nil {
		log.Fatal(err)
	}

	newTrainer := func(segmenter Segmenter) *MBPETrainer {
		return NewMBPETrainer(NewByteLevel(true), segmenter, NewMBPE(), 1<<16)
	}

	trainers := []struct {
		*MBPETrainer
		string
	}{
		{newTrainer(nil), "en-base"},
		{newTrainer(s050), "en-s050"},
		{newTrainer(s100), "en-s100"},
		{newTrainer(m050), "en-m050"},
		{newTrainer(m100), "en-m100"},
		{newTrainer(NewSequence(s100, m050)), "en-s100-m050"},
		{newTrainer(NewSequence(s100, m100)), "en-s100-m100"},
	}

	for i, t := range trainers {
		dict := filepath.Join(out, "dict.txt")

		if err := t.LoadDict(dict); err != nil {
			if err := t.InitDict("data/culturax/en_part_00000.txt"); err != nil {
				log.Fatal(err)
			}

			if err := t.dict.Save(dict); err != nil {
				log.Fatal(err)
			}
		}

		if i > 0 {
			fmt.Println()
		}

		fmt.Printf("%s\n\n", t.string)

		t.Train()

		dir := filepath.Join(out, t.string)

		if err := os.Mkdir(dir, 0755); err != nil && !errors.Is(err, os.ErrExist) {
			log.Fatal(err)
		}

		if err := t.model.Save(filepath.Join(dir, "vocab.json"), filepath.Join(dir, "merges.txt")); err != nil {
			log.Fatal(err)
		}
	}
}
