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

	base := "out-m100"

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

	// fmt.Println(runner.RunAll(1<<17, 1<<16, 1<<15, 1<<14, 1<<13))
	fmt.Println(runner.RunAll(1 << 15))

	// fmt.Print(runner.RunAll(1 << 17))
	// fmt.Println()
	// fmt.Print(runner.RunAll(1 << 16))
	// fmt.Println()
	// fmt.Print(runner.RunAll(1 << 15))
	// fmt.Println()
	// fmt.Print(runner.RunAll(1 << 14))
	// fmt.Println()
	// fmt.Print(runner.RunAll(1 << 13))
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

	morfessor := func(alpha float64) Segmenter {
		m := NewMorfessor(alpha)

		if err := m.LoadModel("data/morfessor/semisup_model.proto"); err != nil {
			log.Fatal(err)
		}

		// return NewDummySegmenter(m)

		return m
	}

	m000 := morfessor(0.0)
	m010 := morfessor(0.1)
	m020 := morfessor(0.2)
	m030 := morfessor(0.3)
	m040 := morfessor(0.4)
	m050 := morfessor(0.5)
	m060 := morfessor(0.6)
	m070 := morfessor(0.7)
	m080 := morfessor(0.8)
	m090 := morfessor(0.9)
	m100 := morfessor(1.0)

	newTrainer := func(segmenter Segmenter) *MBPETrainer {
		return NewMBPETrainer(NewByteLevel(true), segmenter, NewMBPE(), 1<<17)
	}

	trainers := []struct {
		*MBPETrainer
		string
	}{
		{newTrainer(m000), "en-m000"},
		{newTrainer(m010), "en-m010"},
		{newTrainer(m020), "en-m020"},
		{newTrainer(m030), "en-m030"},
		{newTrainer(m040), "en-m040"},
		{newTrainer(m050), "en-m050"},
		{newTrainer(m060), "en-m060"},
		{newTrainer(m070), "en-m070"},
		{newTrainer(m080), "en-m080"},
		{newTrainer(m090), "en-m090"},
		{newTrainer(m100), "en-m100"},
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
