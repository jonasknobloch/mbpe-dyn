package mbpe

import (
	"bytes"
	"encoding/gob"
	"os"
)

type Serializable struct {
	Vocab  []string
	Merges [][2]string
}

func SerializeModel(model *MBPE, name string) error {
	var serializable Serializable

	serializable.Vocab = model.vocab
	serializable.Merges = model.merges

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

func DeserializeModel(blob []byte) (*MBPE, error) {
	buf := bytes.NewBuffer(blob)

	decoder := gob.NewDecoder(buf)

	var serializable Serializable

	if err := decoder.Decode(&serializable); err != nil {
		return nil, err
	}

	var model MBPE

	model.vocab = serializable.Vocab
	model.merges = serializable.Merges

	atoi := make(map[string]int, len(model.vocab))
	itoa := make(map[int]string, len(model.vocab))

	for i, s := range serializable.Vocab {
		atoi[s] = i
		itoa[i] = s
	}

	model.atoi = atoi
	model.itoa = itoa

	ranks := make(map[[2]string]int, len(model.merges))

	for i, merge := range model.merges {
		ranks[merge] = i
	}

	model.ranks = ranks

	return &model, nil
}
