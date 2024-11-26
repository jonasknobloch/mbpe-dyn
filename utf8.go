package main

import "github.com/sugarme/tokenizer/pretokenizer"

var BytesChar = pretokenizer.BytesChar
var CharBytes = pretokenizer.CharBytes

func Alphabet() []string {
	alphabet := make([]string, 256)

	for i, c := range BytesChar {
		alphabet[i] = c
	}

	return alphabet
}
