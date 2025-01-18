package main

import (
	"fmt"
	"strings"
)

type Evaluator interface {
	Eval(tokenizer *Tokenizer, maxRank int) ([]float64, error)
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

func (r *EvalRunner) RunAll(vocabSizes ...int) string {
	if len(vocabSizes) == 0 {
		vocabSizes = []int{-1}
	}

	results := make([][]string, 0, len(r.tokenizers)*len(vocabSizes)+1)

	columns := make([]string, len(r.evaluators)+2)
	widths := make([]int, len(r.evaluators)+2)

	results = append(results, append([]string{"#", "Vocabulary"}, r.evaluatorNames...))

	for i, name := range results[0] {
		columns[i] = name
		widths[i] = len(name)
	}

	for _, vocabSize := range vocabSizes {
		for i, tokenizer := range r.tokenizers {
			model, ok := tokenizer.model.(*MBPE)

			if !ok {
				panic("unexpected model type")
			}

			row := make([]string, len(columns))

			row[0] = r.tokenizerNames[i]

			if vocabSize == -1 {
				row[1] = fmt.Sprintf("%d", len(model.vocab))
			} else {
				row[1] = fmt.Sprintf("%d", vocabSize)
			}

			maxRank := -1

			if vocabSize > -1 {
				alphabet := model.Alphabet()

				if vocabSize < len(alphabet) {
					panic("vocab size smaller than alphabet")
				}

				if vocabSize > len(model.vocab) {
					panic("vocab size larger than model vocabulary")
				}

				maxRank = vocabSize - len(alphabet)
			}

			for j, evaluator := range r.evaluators {
				result, err := evaluator.Eval(tokenizer, maxRank)

				if err != nil {
					row[j+2] = "error"

					continue
				}

				s := make([]string, len(result))

				for k, v := range result {
					s[k] = fmt.Sprintf("%.4f", v)
				}

				row[j+2] = strings.Join(s, ", ")
			}

			for j, cell := range row {
				if len(cell) > widths[j] {
					widths[j] = len(cell)
				}
			}

			results = append(results, row)
		}
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
