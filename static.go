package mbpe

import (
	"errors"
	"strings"
)

type Static struct {
	dict  map[string][]string
	alpha float64
}

func NewStatic(alpha float64) *Static {
	if alpha < 0 || alpha > 1 {
		panic("alpha must be in [0, 1]")
	}

	return &Static{
		dict:  make(map[string][]string),
		alpha: alpha,
	}
}

func (c *Static) LoadDict(name string) error {
	return readTsv(name, func(record []string) error {
		if len(record) != 2 {
			return errors.New("unexpected number of fields")
		}

		c.dict[record[0]] = strings.Split(record[1], " ")

		return nil
	})
}

func (c *Static) Segment(text string) ([]string, float64) {
	substrings, ok := c.dict[text]

	if !ok {
		return []string{text}, 0
	}

	return substrings, c.alpha
}
