package hf

import (
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
	Model         json.RawMessage `json:"model"`
	Vocab         Vocab           `json:"vocab"`
	Merges        Merges          `json:"merges"`
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		Version: "1.0",
	}
}

func (t *Tokenizer) Encode(name string) error {
	file, err := os.Create(name)

	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewEncoder(file).Encode(t)
}

func (t *Tokenizer) Decode(name string) error {
	file, err := os.Open(name)

	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewDecoder(file).Decode(t)
}
