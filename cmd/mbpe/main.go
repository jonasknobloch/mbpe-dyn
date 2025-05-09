package main

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"golang.org/x/image/colornames"
	"log"
	mbpe "mbpe-dyn"
	"mbpe-dyn/hf"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// tokenize()
	// eval()
	// train()
	// server()
	// serialize()
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

	initTokenizer := func(vocab, merges string) *mbpe.Tokenizer {
		model := mbpe.NewMBPE()

		if err := model.Load(vocab, merges); err != nil {
			log.Fatal(err)
		}

		tokenizer := mbpe.NewTokenizer(model)

		byteLevel := mbpe.NewByteLevel(true)

		tokenizer.SetPreTokenizer(byteLevel)
		tokenizer.SetDecoder(byteLevel)

		return tokenizer
	}

	runner := mbpe.NewRunner()

	for _, name := range tokenizers {
		runner.AddTokenizer(*initTokenizer(filepath.Join(name, "vocab.json"), filepath.Join(name, "merges.txt")), filepath.Base(name))
	}

	bpr := func() mbpe.Evaluator {
		bprEval := mbpe.NewBPREvaluator()

		if err := bprEval.LoadSegmentations("data/mbpe/goldstd_trainset.segmentation.eng.tsv"); err != nil {
			log.Fatal(err)
		}

		return bprEval
	}()

	ml := func() mbpe.Evaluator {
		mlEval := mbpe.NewMergeLayerEvaluator()

		if err := mlEval.LoadSegmentations("data/mbpe/goldstd_trainset.segmentation.eng.tsv"); err != nil {
			log.Fatal(err)
		}

		return mlEval
	}()

	fert := func() mbpe.Evaluator {
		fertilityEval := mbpe.NewFertilityEvaluator()

		if err := fertilityEval.InitDict("data/culturax/en_part_00001-10k.txt"); err != nil {
			log.Fatal(err)
		}

		return fertilityEval
	}()

	ref := func() mbpe.Evaluator {
		refEval := mbpe.NewReferenceEvaluator()

		if err := refEval.LoadModel(filepath.Join(tokenizers[0], "vocab.json"), filepath.Join(tokenizers[0], "merges.txt")); err != nil {
			log.Fatal(err)
		}

		return refEval
	}()

	runner.AddEvaluator(bpr, "Boundary Precision Recall")
	runner.AddEvaluator(ml, "Merge Layer")
	runner.AddEvaluator(fert, "Fertility")
	runner.AddEvaluator(ref, "Reference Overlap")

	baseline := mbpe.NewRunner()

	baseline.AddTokenizer(*initTokenizer(filepath.Join(tokenizers[0], "vocab.json"), filepath.Join(tokenizers[0], "merges.txt")), filepath.Base(tokenizers[0]))

	baseline.AddEvaluator(ml, "Merge Layer")
	baseline.AddEvaluator(fert, "Fertility")

	md00, raw00 := runner.RunAll(1 << 16)
	md01, raw01 := runner.RunAll(1 << 15)
	md02, raw02 := runner.RunAll(1 << 14)

	_, rawBase := baseline.RunAll(100000, 90000, 80000, 70000, 60000, 50000, 40000, 30000, 20000, 10000, 5000)

	s00 := mbpe.NewPlotData(mbpe.SelectColumn(raw00[2], 0), mbpe.SelectColumn(raw00[1], 0), true, false, "2^16", colornames.Red)
	s01 := mbpe.NewPlotData(mbpe.SelectColumn(raw01[2], 0), mbpe.SelectColumn(raw01[1], 0), true, false, "2^15", colornames.Green)
	s02 := mbpe.NewPlotData(mbpe.SelectColumn(raw02[2], 0), mbpe.SelectColumn(raw02[1], 0), true, false, "2^14", colornames.Blue)

	sBase := mbpe.NewPlotData(mbpe.SelectColumn(rawBase[1], 0), mbpe.SelectColumn(rawBase[0], 0), false, true, "baseline", colornames.Black)

	data := []mbpe.PlotData{s00, s01, s02, sBase}

	mbpe.Plot(data, [2]float64{1.05, 1.32}, [2]float64{0.76, 0.86}, "Fertility", "Merge Layer")

	fmt.Printf(md00)
	fmt.Println()
	fmt.Printf(md01)
	fmt.Println()
	fmt.Printf(md02)
}

func tokenize() {
	model := mbpe.NewMBPE()

	// err := model.Load("out-m100/00-en-m000/vocab.json", "out-m100/00-en-m000/merges.txt")
	// err := model.Load("out-m100/01-en-m010/vocab.json", "out-m100/01-en-m010/merges.txt")
	// err := model.Load("out-m100/02-en-m020/vocab.json", "out-m100/02-en-m020/merges.txt")
	// err := model.Load("out-m100/03-en-m030/vocab.json", "out-m100/03-en-m030/merges.txt")
	// err := model.Load("out-m100/04-en-m040/vocab.json", "out-m100/04-en-m040/merges.txt")
	// err := model.Load("out-m100/05-en-m050/vocab.json", "out-m100/05-en-m050/merges.txt")
	// err := model.Load("out-m100/06-en-m060/vocab.json", "out-m100/06-en-m060/merges.txt")
	// err := model.Load("out-m100/07-en-m070/vocab.json", "out-m100/07-en-m070/merges.txt")
	// err := model.Load("out-m100/08-en-m080/vocab.json", "out-m100/08-en-m080/merges.txt")
	// err := model.Load("out-m100/09-en-m090/vocab.json", "out-m100/09-en-m090/merges.txt")
	// err := model.Load("out-m100/10-en-m100/vocab.json", "out-m100/10-en-m100/merges.txt")

	err := model.Load("out-m100/05-en-m050/vocab.json", "out-m100/05-en-m050/merges.txt")
	// err := model.Load("out-m100/10-en-m100/vocab.json", "out-m100/10-en-m100/merges.txt")

	if err != nil {
		log.Fatal(err)
	}
	//
	// if err := SerializeModel(model, "en-m100.gob"); err != nil {
	// 	log.Fatal(err)
	// }

	// var model *MBPE
	//
	// if m, err := DeserializeModel(m000); err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	model = m
	// }

	tokenizer := mbpe.NewTokenizer(model)

	byteLevel := mbpe.NewByteLevel(true)

	tokenizer.SetPreTokenizer(byteLevel)
	tokenizer.SetDecoder(byteLevel)

	ids := tokenizer.Tokenize(" airsickness")
	tokens := model.ToString(ids)

	fmt.Println(ids)
	fmt.Println(tokens)

	fmt.Println(tokenizer.Decoder().Decode(tokens))
}

func segmentFile(name string, vocabSize int) {
	model := mbpe.NewMBPE()

	tokenizer := mbpe.NewTokenizer(model)

	byteLevel := mbpe.NewByteLevel(true)

	tokenizer.SetPreTokenizer(byteLevel)
	tokenizer.SetDecoder(byteLevel)

	if err := model.Load("out/00-en-base/vocab.json", "out/00-en-base/merges.txt"); err != nil {
		log.Fatal(err)
	}

	compounds := make([]string, 0)

	if err := mbpe.ReadTsv(name, func(record []string) error {
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
		segmentation, ok := mbpe.GetTokenizerSegmentation(*tokenizer, compound, maxRank)

		if !ok {
			continue
		}

		segmentations[i] = segmentation
	}

	if err := mbpe.ToFile("segmentations.txt", func(writer *bufio.Writer) error {
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

	morfessor := func(alpha float64) mbpe.Segmenter {
		m := mbpe.NewMorfessor(alpha)

		if err := m.LoadModel("data/morfessor/semisup_model.proto"); err != nil {
			log.Fatal(err)
		}

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

	newTrainer := func(segmenter mbpe.Segmenter) *mbpe.MBPETrainer {
		return mbpe.NewMBPETrainer(mbpe.NewByteLevel(true), segmenter, mbpe.NewMBPE(), 1<<17)
	}

	trainers := []struct {
		*mbpe.MBPETrainer
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

			if err := t.Dict().Save(dict); err != nil {
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

		if err := t.Model().Save(filepath.Join(dir, "vocab.json"), filepath.Join(dir, "merges.txt")); err != nil {
			log.Fatal(err)
		}
	}
}

func serialize() {
	model := mbpe.NewMBPE()

	if err := model.Load("out/00-en-m000/vocab.json", "out/00-en-m000/merges.txt"); err != nil {
		log.Fatal(err)
	}

	tokenizer := hf.NewTokenizer()

	if err := tokenizer.Decode("data/mbpe/empty.json"); err != nil {
		log.Fatal(err)
	}

	atoi := make(map[string]int)

	last := 0

	for i, token := range model.Vocab()[:1024] {
		atoi[token] = i

		if i > last {
			last = i
		}
	}

	tokenizer.Model.Vocab = atoi
	tokenizer.Model.Merges = model.Merges()[:1024-len(model.Alphabet())]

	eot := hf.AddedToken{
		ID:         last + 1,
		Content:    "<|endoftext|>",
		SingleWord: false,
		Lstrip:     false,
		Rstrip:     false,
		Normalized: true,
		Special:    true,
	}

	tokenizer.AddedTokens = []hf.AddedToken{eot}

	if err := tokenizer.Encode("tokenizer.json"); err != nil {
		log.Fatal(err)
	}
}
