package hf

import (
	"testing"
)

func TestTokenizer_Decode(t *testing.T) {
	var tokenizer Tokenizer

	if err := tokenizer.Decode("../data/mbpe/empty.json"); err != nil {
		t.Fatal(err)
	}

	if len(tokenizer.Vocab) == 0 {
		t.Error("empty vocab")
	}

	if len(tokenizer.Merges) == 0 {
		t.Error("empty merges")
	}
}
