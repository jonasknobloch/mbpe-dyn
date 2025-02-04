package main

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

type Segmenter interface {
	Segment(string) ([]string, float64)
}

func SegmentWithoutPrefixWhitespace(compound string, segmenter Segmenter) ([]string, float64) {
	removedPrefixSpace := stripPrefixWhitespace(&compound)

	substrings, alpha := segmenter.Segment(compound)

	if removedPrefixSpace {
		addPrefixWhitespace(&substrings)
	}

	return substrings, alpha
}

func stripPrefixWhitespace(compound *string) bool {
	hasPrefixSpace := strings.HasPrefix(*compound, "Ġ")

	if hasPrefixSpace {
		*compound = (*compound)[len("Ġ"):]
	}

	return hasPrefixSpace
}

func addPrefixWhitespace(segments *[]string) {
	if len(*segments) == 0 {
		return
	}

	(*segments)[0] = "Ġ" + (*segments)[0]
}

func SegmentFile(name string, segmenter Segmenter) error {
	gold := make([]string, 0)

	if err := readTsv(name, func(record []string) error {
		if len(record) == 0 {
			return errors.New("unexpected number of fields")
		}

		gold = append(gold, record[0])

		return nil
	}); err != nil {
		return err
	}

	segmentations := make([][]string, len(gold))

	for i, compound := range gold {
		segmentations[i], _ = segmenter.Segment(compound)
	}

	if err := toFile("segmentations.txt", func(writer *bufio.Writer) error {
		for i, segmentation := range segmentations {
			if _, err := writer.WriteString(fmt.Sprintf("%s\t%s\n", gold[i], strings.Join(segmentation, " "))); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
