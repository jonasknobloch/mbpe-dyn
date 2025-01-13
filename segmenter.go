package main

import "strings"

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
