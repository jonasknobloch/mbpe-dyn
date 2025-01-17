package main

import (
	"fmt"
	"strings"
)

type Evaluator interface {
	Eval(tokenizer *Tokenizer) ([]float64, error)
}

type EvalRunner struct {
	tokenizers     []*Tokenizer
	evaluators     []Evaluator
	tokenizerNames []string
	evaluatorNames []string
}

func NewRunner() *EvalRunner {
	return &EvalRunner{
		tokenizers:     make([]*Tokenizer, 0),
		evaluators:     make([]Evaluator, 0),
		tokenizerNames: make([]string, 0),
		evaluatorNames: make([]string, 0),
	}
}

func (r *EvalRunner) AddTokenizer(tokenizer *Tokenizer, name string) {
	r.tokenizers = append(r.tokenizers, tokenizer)
	r.tokenizerNames = append(r.tokenizerNames, name)
}

func (r *EvalRunner) AddEvaluator(evaluator Evaluator, name string) {
	r.evaluators = append(r.evaluators, evaluator)
	r.evaluatorNames = append(r.evaluatorNames, name)
}

func (r *EvalRunner) RunAll() string {
	results := make([][]string, len(r.tokenizers)+1)
	columns := make([]string, len(r.evaluators)+1)

	widths := make([]int, len(columns))

	for _, name := range r.tokenizerNames {
		if len(name) > widths[0] {
			widths[0] = len(name)
		}
	}

	for i, name := range r.evaluatorNames {
		columns[i+1] = name
		widths[i+1] = len(name)
	}

	results[0] = append([]string{"#"}, r.evaluatorNames...)

	for i, tokenizer := range r.tokenizers {
		row := make([]string, len(columns))

		row[0] = r.tokenizerNames[i]

		for j, evaluator := range r.evaluators {
			result, err := evaluator.Eval(tokenizer)

			if err != nil {
				row[j+1] = "x"

				continue
			}

			s := make([]string, len(result))

			for k, v := range result {
				s[k] = fmt.Sprintf("%.4f", v)
			}

			row[j+1] = strings.Join(s, ", ")

			if len(row[j+1]) > widths[j+1] {
				widths[j+1] = len(row[j+1])
			}
		}

		results[i+1] = row
	}

	return markdownTable(results, widths)
}

func markdownTable(table [][]string, widths []int) string {
	var sb strings.Builder

	divider := make([]string, len(widths))

	for i, width := range widths {
		divider[i] = strings.Repeat("-", width+2)
	}

	for i, row := range table {
		cells := make([]string, len(row))

		for j, cell := range row {
			cells[j] = fmt.Sprintf("%-*s", widths[j], cell)
		}

		sb.WriteString(fmt.Sprintf("| %s |\n", strings.Join(cells, " | ")))

		if i == 0 {
			sb.WriteString(fmt.Sprintf("|%s|\n", strings.Join(divider, "|")))
		}
	}

	return sb.String()
}
