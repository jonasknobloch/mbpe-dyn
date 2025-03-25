package mbpe

import (
	"testing"
)

func TestDict_Init(t *testing.T) {
	d := NewDict()

	if err := d.ProcessFiles("data/shakespeare.txt"); err != nil {
		panic(err)
	}

	if err := d.Save("dict.txt"); err != nil {
		panic(err)
	}

	// TODO implement
}

func BenchmarkDict_ProcessFiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d := NewDict()

		err := d.ProcessFiles("data/shakespeare.txt")

		if err != nil {
			panic(err)
		}
	}
}
