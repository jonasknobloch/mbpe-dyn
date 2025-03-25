package mbpe

import (
	"bytes"
	"cmp"
	"encoding/json"
	"sort"
)

type Vocab[K comparable, V cmp.Ordered] map[K]V

func (m Vocab[K, V]) MarshalJSON() ([]byte, error) {
	vtok := make(map[V]K, len(m))

	for key, value := range m {
		vtok[value] = key
	}

	values := make([]V, 0, len(vtok))

	for value := range vtok {
		values = append(values, value)
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})

	var buf bytes.Buffer

	buf.WriteByte('{')

	enc := json.NewEncoder(&buf)

	for i, value := range values {
		if i > 0 {
			buf.WriteByte(',')
		}

		if err := enc.Encode(vtok[value]); err != nil {
			return nil, err
		}

		buf.WriteByte(':')

		if err := enc.Encode(value); err != nil {
			return nil, err
		}
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}
