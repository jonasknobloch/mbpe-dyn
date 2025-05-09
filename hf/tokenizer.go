package hf

import (
	"bytes"
	"encoding/json"
	"os"
)

type Tokenizer struct {
	Version       string          `json:"version"`
	Truncation    json.RawMessage `json:"truncation"`
	Padding       json.RawMessage `json:"padding"`
	AddedTokens   []AddedToken    `json:"added_tokens"`
	Normalizer    json.RawMessage `json:"normalizer"`
	PreTokenizer  json.RawMessage `json:"pre_tokenizer"`
	PostProcessor json.RawMessage `json:"post_processor"`
	Decoder       json.RawMessage `json:"decoder"`
	Model         BPE             `json:"model"` // TODO model interface
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		Version: "1.0",
	}
}

func (t *Tokenizer) Encode(name string) error {
	var file *os.File

	if f, err := os.Create(name); err != nil {
		return err
	} else {
		file = f
	}

	defer file.Close()

	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)

	enc.SetEscapeHTML(false)

	if err := enc.Encode(t); err != nil {
		return err
	}

	data := bytes.TrimRight(buf.Bytes(), "\n")

	_, err := file.Write(data)

	return err
}

func (t *Tokenizer) Decode(name string) error {
	file, err := os.Open(name)

	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewDecoder(file).Decode(t)
}
