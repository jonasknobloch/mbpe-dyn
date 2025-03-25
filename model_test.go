package mbpe

import (
	"log"
	"testing"
)

func BenchmarkMBPE_Tokenize(b *testing.B) {
	model := NewMBPE()

	err := model.Load("out/vocab.json", "out/merges.txt")

	if err != nil {
		log.Fatal(err)
	}

	tokenizer := NewTokenizer(model)

	byteLevel := NewByteLevel(true)

	tokenizer.SetPreTokenizer(byteLevel)
	tokenizer.SetDecoder(byteLevel)

	for i := 0; i < b.N; i++ {
		tokenizer.Tokenize("To infinity and beyond!")
	}
}
