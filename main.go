package main

import (
	"fmt"
	"log"
)

func main() {
	trainer := NewTrainer()

	if err := trainer.Train("data/shakespeare.txt", 500); err != nil {
		log.Fatal(err)
	}

	fmt.Println(trainer.dict)
}
