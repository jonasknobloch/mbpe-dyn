package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	stok "github.com/sugarme/tokenizer"
	sbpe "github.com/sugarme/tokenizer/model/bpe"
	spre "github.com/sugarme/tokenizer/pretokenizer"
	"log"
	"os"
)

func main() {
	if _, err := os.Stat("temp.gob"); err != nil {
		if os.IsNotExist(err) {
			train()
		} else {
			log.Fatal(err)
		}
	}

	tokenizer, err := DeserializeTokenizer("temp.gob")

	if err != nil {
		log.Fatal(err)
	}

	model := NewMBPE()

	model.tokenizer = tokenizer

	tokens := model.Tokenize("To infinity and beyond!")

	fmt.Println(tokens)
	fmt.Println(tokenizer.ToString(tokens))
}

func train() {
	trainer := NewTrainer(5000)

	if err := trainer.Train("data/shakespeare.txt"); err != nil {
		log.Fatal(err)
	}

	tokenizer := trainer.model.tokenizer

	if err := SerializeTokenizer(tokenizer, "temp.gob"); err != nil {
		log.Fatal(err)
	}

	if err := toJSON("vocab.json", tokenizer.atoi); err != nil {
		log.Fatal(err)
	}

	// if err := toFile("vocab.txt", func(writer *bufio.Writer) error {
	// 	for i, token := range tokenizer.vocab {
	// 		if _, err := writer.WriteString(strconv.Itoa(i) + " " + token + "\n"); err != nil {
	// 			return err
	// 		}
	// 	}
	//
	// 	return nil
	// }); err != nil {
	// 	log.Fatal(err)
	// }

	if err := toFile("merges.txt", func(writer *bufio.Writer) error {
		for _, merge := range tokenizer.merges {
			if _, err := writer.WriteString(merge[0] + " " + merge[1] + "\n"); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func fromFile(name string, callback func(scanner *bufio.Scanner) error) error {
	file, err := os.Open(name)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	if err := callback(scanner); err != nil {
		return err
	}

	return nil
}

func toFile(name string, callback func(writer *bufio.Writer) error) error {
	file, err := os.Create(name)

	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	if err := callback(writer); err != nil {
		return err
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	return nil
}

func fromJSON(name string, data interface{}) error {
	file, err := os.Open(name)

	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(data); err != nil {
		return err
	}

	return nil
}

func toJSON(name string, data interface{}) error {
	file, err := os.Create(name)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
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

	if err := toFile("reference/vocab.txt", func(writer *bufio.Writer) error {
		for token := range result.GetVocab() {
			if _, err := writer.WriteString(token + "\n"); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func fromReference(nameVocab, nameMerges string) (*Tokenizer, error) {
	atoi := make(map[string]int)

	if err := fromJSON(nameVocab, &atoi); err != nil {
		return nil, err
	}

	itoa := make(map[int]string, len(atoi))

	for token, idx := range atoi {
		itoa[idx] = token
	}

	vocab := make([]string, len(itoa))

	for i := range len(itoa) {
		vocab[i] = itoa[i]
	}

	merges := make([][2]string, 0) // unknown number of merges

	if err := fromFile(nameMerges, func(scanner *bufio.Scanner) error {
		for scanner.Scan() {
			line := scanner.Text()
			var merge [2]string

			if _, err := fmt.Sscanf(line, "%s %s", &merge[0], &merge[1]); err != nil {
				return err
			}

			merges = append(merges, merge)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &Tokenizer{
		vocab:  vocab,
		atoi:   atoi,
		itoa:   itoa,
		merges: merges,
	}, nil
}
