package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	trainer := NewTrainer(5000)

	if err := trainer.Train("data/shakespeare.txt"); err != nil {
		log.Fatal(err)
	}

	tokenizer := trainer.model.tokenizer

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
