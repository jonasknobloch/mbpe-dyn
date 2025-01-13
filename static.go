package main

import (
	"bufio"
	"encoding/csv"
	"os"
	"strings"
)

type Static struct {
	dict  map[string][]string
	alpha float64
}

func NewStatic(alpha float64) *Static {
	if alpha < 0 || alpha > 1 {
		panic("alpha must be in [0, 1]")
	}

	return &Static{
		dict:  make(map[string][]string),
		alpha: alpha,
	}
}

func (c *Static) LoadDict(name string) error {
	var file *os.File

	if f, err := os.Open(name); err != nil {
		return err
	} else {
		file = f
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

			return err
		}

		if len(record) != 2 {
			panic("Invalid record")
		}

		c.dict[record[0]] = strings.Split(record[1], " ")
	}

	return nil
}

func (c *Static) Segment(text string) ([]string, float64) {
	substrings, ok := c.dict[text]

	if !ok {
		return []string{text}, 0
	}

	return substrings, c.alpha
}
