package main

import (
	"bytes"
	"encoding/gob"
	"os"
)

type Serializable struct {
	Vocab  []string
	Merges [][2]string
}

func SerializeTokenizer(tokenizer *Tokenizer, name string) error {
	var serializable Serializable

	serializable.Vocab = tokenizer.vocab
	serializable.Merges = tokenizer.merges

	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(serializable)

	if err != nil {
		return err
	}

	file, err := os.Create(name)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(buf.Bytes())

	if err != nil {
		return err
	}

	return nil
}

func DeserializeTokenizer(name string) (*Tokenizer, error) {
	file, err := os.Open(name)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	stat, err := file.Stat()

	if err != nil {
		return nil, err
	}

	data := make([]byte, stat.Size())

	if _, err := file.Read(data); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)

	decoder := gob.NewDecoder(buf)

	var serializable Serializable

	if err := decoder.Decode(&serializable); err != nil {
		return nil, err
	}

	var tokenizer Tokenizer

	tokenizer.vocab = serializable.Vocab
	tokenizer.merges = serializable.Merges

	atoi := make(map[string]int, len(tokenizer.vocab))
	itoa := make(map[int]string, len(tokenizer.vocab))

	for i, s := range serializable.Vocab {
		atoi[s] = i
		itoa[i] = s
	}

	tokenizer.atoi = atoi
	tokenizer.itoa = itoa

	return &tokenizer, nil
}
