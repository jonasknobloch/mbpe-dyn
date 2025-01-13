package main

type Sequence []Segmenter

func NewSequence(segmenters ...Segmenter) *Sequence {
	return (*Sequence)(&segmenters)
}

func (s Sequence) Segment(compound string) ([]string, bool) {
	if len(s) == 0 {
		return []string{compound}, false
	}

	var segments []string
	var ok bool

	for _, segmenter := range s {
		segments, ok = segmenter.Segment(compound)

		if ok {
			return segments, true
		}
	}

	return []string{compound}, false
}
