package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Sample struct {
	Left  string
	Right string
	Ratio float64
}

type Results struct {
	WugAdjNominalization struct {
		Predictions []struct {
			ID   string `json:"id"`
			Pred string `json:"pred"`
		} `json:"predictions"`
	} `json:"wug_adj_nominalization"`
}

func main() {
	var samples []Sample
	var ratios []float64

	if s, r, err := loadSamples("data/wug_results/wug_adj_nominalization.jsonl"); err != nil {
		log.Fatal(err)
	} else {
		samples = s
		ratios = r
	}

	paths, stubs := walkResultsStatic("data/wug_results/results/gpt2_%d_%s%s_babylm_v2/main/zero_shot/causal/wug/wug_adj_nominalization/predictions.json")

	out := make([]string, 0, len(paths)+1)

	for i, path := range paths {
		fmt.Printf("\n%s\n\n", path)

		accuracies, total, err := evalPredictions(samples, ratios, path)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("### RATIO ACCURACY\n")

		for j, ratio := range ratios {
			fmt.Printf("%.16f: %.2f\n", ratio, accuracies[j]*100)
		}

		fmt.Printf("\n### AVERAGE ACCURACY\n%.2f\n", total*100)

		if len(out) == 0 {
			row := "vocab,prefix,alpha,"

			for _, r := range ratios {
				row += fmt.Sprintf("%.2f,", r)
			}

			row += "average"

			out = append(out, row)
		}

		row := fmt.Sprintf("%s,%s,%s,", stubs[i][0], stubs[i][1], stubs[i][2])

		for _, accuracy := range accuracies {
			if math.IsNaN(accuracy) {
				accuracy = 0.0
			}

			row += fmt.Sprintf("%.2f,", accuracy)
		}

		row += fmt.Sprintf("%.2f", total)

		out = append(out, row)
	}

	fmt.Println()

	for _, row := range out {
		fmt.Println(row)
	}
}

func loadSamples(name string) ([]Sample, []float64, error) {
	samples := make([]Sample, 0, 200)
	ratios := make(map[float64]struct{})

	if err := fromFile(name, func(scanner *bufio.Scanner) error {
		for scanner.Scan() {
			line := scanner.Text()

			var raw struct {
				Sentences string  `json:"sentences"`
				Ratio     float64 `json:"ratio"`
			}

			if err := json.Unmarshal([]byte(line), &raw); err != nil {
				return fmt.Errorf("error parsing JSON: %w", err)
			}

			parts := strings.SplitN(raw.Sentences, "\t", 2)

			if len(parts) != 2 {
				return fmt.Errorf("invalid sentence format: %q", raw.Sentences)
			}

			sample := Sample{
				Left:  parts[0],
				Right: parts[1],
				Ratio: raw.Ratio,
			}

			samples = append(samples, sample)
			ratios[sample.Ratio] = struct{}{}
		}

		return scanner.Err()
	}); err != nil {
		log.Fatal(err)
	}

	return samples, mapToSortedSlice(ratios), nil
}

func evalPredictions(samples []Sample, ratios []float64, name string) ([]float64, float64, error) {
	// re := regexp.MustCompile(`Create a noun out of the following adjective:\s+(\w+)\.\s+(\w+)`)
	re := regexp.MustCompile(`Create a noun out of the following adjective:(\s+\w+)\.(\s+\w+)`)

	match := func(line string) (string, string, bool) {
		matches := re.FindStringSubmatch(line)

		if len(matches) == 3 {
			return matches[1], matches[2], true
		}

		return "", "", false
	}

	results := Results{}

	if err := fromJSON(name, &results); err != nil {
		return nil, 0, err
	}

	classes := []string{"able", "ish", "ive", "ous"}

	totalPositives := 0
	totalNegatives := 0

	accuracies := make([]float64, len(ratios))

	for i, ratio := range ratios {
		positives := 0
		negatives := 0

		// if ratio < 0.7 && ratio > 0.3 {
		// 	continue
		// }

		for _, class := range classes {
			p := 0
			n := 0

			for i, prediction := range results.WugAdjNominalization.Predictions {
				s := samples[i]

				if s.Ratio != ratio {
					continue
				}

				adj, _, ok := match(prediction.Pred)

				if !ok {
					return nil, 0, errors.New("could not match prediction: " + prediction.Pred)
				}

				if adj[len(adj)-len(class):] != class {
					continue
				}

				// _, left, _ := match(s.Left)
				// _, right, _ := match(s.Right)

				if prediction.Pred == s.Left {
					p++
				} else {
					n++
				}

				// if s.Ratio > 0.5 {
				// 	if prediction.Pred == s.Left {
				// 		p++
				// 	} else {
				// 		n++
				// 	}
				// } else {
				// 	if prediction.Pred == s.Right {
				// 		p++
				// 	} else {
				// 		n++
				// 	}
				// }
			}

			positives += p
			negatives += n
		}

		accuracy := float64(positives) / float64(positives+negatives)

		accuracies[i] = accuracy

		totalPositives += positives
		totalNegatives += negatives
	}

	average := float64(totalPositives) / float64(totalPositives+totalNegatives)

	return accuracies, average, nil
}

func fromFile(name string, callback func(scanner *bufio.Scanner) error) error {
	file, err := os.Open(name)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	if err := callback(scanner); err != nil {
		return err
	}

	return nil
}

func fromJSON(name string, data interface{}) error {
	file, err := os.Open(name)

	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(data); err != nil {
		return err
	}

	return nil
}

func mapToSortedSlice(m map[float64]struct{}) []float64 {
	s := make([]float64, 0, len(m))

	for k := range m {
		s = append(s, k)
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(s)))

	return s
}

func walkResults() ([]string, error) {
	root := "data/wug_results/results"
	match := "wug_adj_nominalization/predictions.json"

	paths := make([]string, 0)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.Contains(path, match) {
			paths = append(paths, path)
		}

		return nil
	})

	return paths, err
}

func walkResultsStatic(format string) ([]string, [][3]string) {
	vocabSizes := []int{8192, 16384, 32768, 50256, 100512}
	prefixes := []string{"m", "mi"}
	alphas := []float64{0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0}

	paths := make([]string, 0, len(vocabSizes)*len(prefixes)*len(alphas))
	stubs := make([][3]string, 0, len(vocabSizes)*len(prefixes)*len(alphas))

	for _, vocab := range vocabSizes {
		for _, prefix := range prefixes {
			for _, alpha := range alphas {
				path := fmt.Sprintf(format, vocab, prefix, fmt.Sprintf("%03d", int(alpha*100)))

				stub := [3]string{strconv.Itoa(vocab), prefix, fmt.Sprintf("%.2f", alpha)}

				paths = append(paths, path)
				stubs = append(stubs, stub)
			}
		}
	}

	return paths, stubs
}
