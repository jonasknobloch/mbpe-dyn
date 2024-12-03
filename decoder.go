package main

type Decoder interface {
	Decode([]string) string
}
