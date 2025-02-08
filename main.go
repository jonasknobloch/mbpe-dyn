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
	// plot()
	// train()
}

func eval() {
	tokenizers := make([]string, 0)

	base := "out"

	if err := filepath.WalkDir(base, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(base, path)

		if err != nil {
			return err
		}

		depth := strings.Count(rel, string(os.PathSeparator))

		if d.IsDir() && rel != "." && depth == 0 {
			tokenizers = append(tokenizers, path)
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

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

	for _, name := range tokenizers {
		runner.AddTokenizer(*initTokenizer(filepath.Join(name, "vocab.json"), filepath.Join(name, "merges.txt")), filepath.Base(name))
	}

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

		if err := refEval.LoadModel(filepath.Join(tokenizers[0], "vocab.json"), filepath.Join(tokenizers[0], "merges.txt")); err != nil {
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

	ids := tokenizer.Tokenize("To infinity and beyond!")
	tokens := model.ToString(ids)

	fmt.Println(ids)
	fmt.Println(tokens)

	fmt.Println(tokenizer.decoder.Decode(tokens))
}

func segmentFile(name string, vocabSize int) {
	model := NewMBPE()

	tokenizer := NewTokenizer(model)

	byteLevel := NewByteLevel(true)

	tokenizer.SetPreTokenizer(byteLevel)
	tokenizer.SetDecoder(byteLevel)

	if err := model.Load("out/00-en-base/vocab.json", "out/00-en-base/merges.txt"); err != nil {
		log.Fatal(err)
	}

	compounds := make([]string, 0)

	if err := readTsv(name, func(record []string) error {
		if len(record) == 0 {
			return errors.New("unexpected number of fields")
		}

		compounds = append(compounds, " "+strings.TrimLeft(record[0], " "))

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	segmentations := make([][]string, len(compounds))

	maxRank := -1

	if vocabSize > -1 {
		maxRank = vocabSize - len(model.Alphabet())
	}

	for i, compound := range compounds {
		segmentation, ok := getTokenizerSegmentation(*tokenizer, compound, maxRank)

		if !ok {
			continue
		}

		segmentations[i] = segmentation
	}

	if err := toFile("segmentations.txt", func(writer *bufio.Writer) error {
		for i, segmentation := range segmentations {
			if _, err := writer.WriteString(fmt.Sprintf("%s\t%s\n", strings.TrimLeft(compounds[i], " "), strings.TrimLeft(strings.Join(segmentation, " "), " "))); err != nil {
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

		dir := filepath.Join(out, fmt.Sprintf("%02d-%s", i, t.string))

		if err := os.Mkdir(dir, 0755); err != nil && !errors.Is(err, os.ErrExist) {
			log.Fatal(err)
		}

		if err := t.model.Save(filepath.Join(dir, "vocab.json"), filepath.Join(dir, "merges.txt")); err != nil {
			log.Fatal(err)
		}
	}
}
