package main

import (
	"fmt"
	"testing"
)

func TestUTF8_BytesChar(t *testing.T) {
	fmt.Println(BytesChar[[]byte(" ")[0]])

	// TODO implement
}

func TestUTF8_CharBytes(t *testing.T) {
	fmt.Println(string(CharBytes["Ġ"]))

	// TODO implement
}
