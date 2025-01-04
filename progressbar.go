package main

import (
	"fmt"
	"strings"
)

type ProgressBar struct {
	width int
	steps int
	n     int
}

func NewProgressBar(width, steps int) *ProgressBar {
	return &ProgressBar{width - 15, steps, 0}
}

func (p *ProgressBar) Increment() {
	p.n++
}

func (p *ProgressBar) String() string {
	progress := int(float64(p.n) / float64(p.steps) * float64(p.width))
	percentage := int(float64(p.n) / float64(p.steps) * 100)

	return fmt.Sprintf("\r%d%% [%s%s] %d/%d", percentage, strings.Repeat("=", progress), strings.Repeat("-", p.width-progress), p.n, p.steps)
}
