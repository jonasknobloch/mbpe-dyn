package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strings"
)

type CELEX struct {
	dict map[string][]string
}

func NewCELEX() *CELEX {
	return &CELEX{
		dict: make(map[string][]string),
	}
}

func (c *CELEX) Init(name string) {
	file, err := os.Open(name)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	bufferedReader := bufio.NewReader(file)
	reader := csv.NewReader(bufferedReader)
	reader.Comma = '\t'

	for {
		record, err := reader.Read()

		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			log.Fatal(err)
		}

		if len(record) != 2 {
			panic("Invalid record")
		}

		c.dict[record[0]] = strings.Split(record[1], " ")
	}
}

func (c *CELEX) Split(text string) []string {
	prefixSpace := strings.HasPrefix(text, "Ġ")

	if prefixSpace {
		text = text[len("Ġ"):]
	}

	substrings, ok := c.dict[text]

	if !ok {
		substrings = []string{text}
	}

	if prefixSpace {
		substrings[0] = "Ġ" + substrings[0]
	}

	return substrings
}
