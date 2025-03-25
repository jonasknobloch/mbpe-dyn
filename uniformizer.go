package mbpe

type Uniformizer struct {
	segmenter Segmenter
}

func NewUniformizer(segmenter Segmenter) *Uniformizer {
	return &Uniformizer{
		segmenter: segmenter,
	}
}

func (u *Uniformizer) Segment(text string) ([]string, float64) {
	template, alpha := u.segmenter.Segment(text)

	n := len(template)

	if n == 1 {
		return template, alpha
	}

	segmentation := make([]string, 0, n)

	runes := []rune(text)
	step := len(runes) / n

	prev := 0

	for i := 1; i < n; i++ {
		b := i * step
		segmentation = append(segmentation, string(runes[prev:b]))
		prev = b
	}

	segmentation = append(segmentation, string(runes[prev:]))

	return segmentation, alpha
}
