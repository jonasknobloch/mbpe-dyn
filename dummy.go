package main

type Dummy struct {
	template Segmenter
}

func NewDummySegmenter(template Segmenter) *Dummy {
	return &Dummy{
		template: template,
	}
}

func (d *Dummy) Segment(text string) ([]string, float64) {
	template, alpha := d.template.Segment(text)

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
