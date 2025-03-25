package mbpe

type Sequence []Segmenter

func NewSequence(segmenters ...Segmenter) *Sequence {
	return (*Sequence)(&segmenters)
}

func (s Sequence) Segment(compound string) ([]string, float64) {
	if len(s) == 0 {
		return []string{compound}, 0
	}

	var segments []string
	var alpha float64

	for _, segmenter := range s {
		segments, alpha = segmenter.Segment(compound)

		if alpha > 0 {
			return segments, alpha
		}
	}

	return []string{compound}, 0
}
