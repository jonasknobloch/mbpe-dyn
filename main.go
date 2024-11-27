package main

import (
	"bufio"
	"fmt"
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

	model := NewMBPE(5000)

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

	if err := toFile("vocab.txt", func(writer *bufio.Writer) error {
		for _, token := range tokenizer.vocab {
			if _, err := writer.WriteString(token + "\n"); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

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
