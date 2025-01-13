package main

import "strings"

type Segmenter interface {
	Segment(string) ([]string, bool)
}

func SegmentWithoutPrefixWhitespace(compound string, segmenter Segmenter) ([]string, bool) {
	removedPrefixSpace := stripPrefixWhitespace(&compound)

	substrings, ok := segmenter.Segment(compound)

	if removedPrefixSpace {
		addPrefixWhitespace(&substrings)
	}

	return substrings, ok
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
