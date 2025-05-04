package hf

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Merges [][2]string

func (m *Merges) MarshalJSON() ([]byte, error) {
	var raw []string

	for _, pair := range *m {
		raw = append(raw, pair[0]+" "+pair[1])
	}

	return json.Marshal(raw)
}

func (m *Merges) UnmarshalJSON(data []byte) error {
	var raw []string

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var result [][2]string

	for _, item := range raw {
		parts := strings.Split(item, " ")

		if len(parts) != 2 {
			return fmt.Errorf("invalid merge: %q", item)
		}

		result = append(result, [2]string{parts[0], parts[1]})
	}

	*m = result

	return nil
}
