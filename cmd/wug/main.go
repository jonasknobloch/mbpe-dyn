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
	foo()

	return

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

func foo() {
	paths, stubs := walkResultsStatic("data/wug_results/out/gpt2_%d_%s%s_babylm_v2_ity_ness_nonce.json")

	for _, path := range paths {
		bar(path)
	}

	out := make([]string, 0, len(paths)+1)

	out = append(out, "vocab,prefix,alpha,acc,error")

	for i, path := range paths {
		fmt.Println(path)

		acc, avgErr := bar(path)

		row := fmt.Sprintf("%s,%s,%s,%.2f,%.2f", stubs[i][0], stubs[i][1], stubs[i][2], acc, avgErr)

		out = append(out, row)
	}

	fmt.Println()

	for _, row := range out {
		fmt.Println(row)
	}
}

func bar(name string) (float64, float64) {
	// pred, _ := processPredictions("data/wug_results/gptj_predictions_nonce.json")
	pred, _ := processPredictions(name)

	bases, bin, rat := cumulate(pred)

	// fmt.Println(bases)
	// fmt.Println(bin)
	// fmt.Println(rat)

	goldRat, goldBin, _ := surveyResponses("data/wug_results/survey_responses.json")

	pos := 0
	neg := 0

	totalError := 0.0

	for i, base := range bases {
		if bin[i] != goldBin[base] { // TODO != seems to work better ?!
			pos++
		} else {
			neg++
		}

		e := math.Abs(goldRat[base] - rat[i][0])

		totalError += e
	}

	acc := float64(pos) / float64(pos+neg)
	avgErr := totalError / float64(len(bases))

	return acc, 1 - avgErr
}

type Entry struct {
	text     string
	tokens   []string
	logProbs []float64
}

type Predictions = map[string]map[string][][2]Entry

func processPredictions(name string) (Predictions, error) {
	var data map[string]interface{}

	if err := fromJSON(name, &data); err != nil {
		return nil, err
	}

	types := [3]string{
		"target",
		"base",
		"instruction",
	}

	instructions := [3][4]string{
		{
			"Nominalized adjective:",
			"Noun:",
			"The following is a nominalized adjective:",
			"The following is a noun:",
		},
		{
			"{} ->",
			"{} :",
			"{} -",
			"{}",
		},
		{
			"Adjective: {}\nNominalization:",
			"Form the nominalization of the given adjective.\n\n{} ->",
			"Nominalize the given adjective.\n\n{} ->",
			"Turn the given adjective into a noun.\n\n{} ->",
		},
	}

	results := make(Predictions)

	for i, t := range types {
		results[t] = make(map[string][][2]Entry)

		for _, instruction := range instructions[i] {
			results[t][instruction] = make([][2]Entry, 200)
		}
	}

	for i, t := range types {
		var typesData map[string]interface{}

		if d, ok := data[types[i]].(map[string]interface{}); ok {
			typesData = d
		} else {
			return nil, errors.New("unexpected data format")
		}

		for j, instruction := range instructions[i] {
			var instructionData map[string]interface{}

			if d, ok := typesData[instructions[i][j]].(map[string]interface{}); ok {
				instructionData = d
			} else {
				return nil, errors.New("unexpected data format")
			}

			keys := getKeys(instructionData)

			for k, key := range keys {
				pairData := instructionData[key].(map[string]interface{})

				stims := strings.Split(key, "_")

				left := pairData[stims[0]].([]interface{})
				right := pairData[stims[1]].([]interface{})

				tokensLeft := left[0].([]interface{})
				logProbsLeft := left[1].([]interface{})

				tokensRight := right[0].([]interface{})
				logProbsRight := right[1].([]interface{})

				pair := [2]Entry{
					{
						text:     stims[0],
						tokens:   make([]string, len(tokensLeft)),
						logProbs: make([]float64, len(logProbsLeft)),
					},
					{
						text:     stims[1],
						tokens:   make([]string, len(tokensRight)),
						logProbs: make([]float64, len(logProbsRight)),
					},
				}

				for pos, token := range tokensLeft {
					pair[0].tokens[pos] = token.(string)
					pair[0].logProbs[pos] = logProbsLeft[pos].(float64)
				}

				for pos, token := range tokensRight {
					pair[1].tokens[pos] = token.(string)
					pair[1].logProbs[pos] = logProbsRight[pos].(float64)
				}

				results[t][instruction][k] = pair
			}
		}
	}

	return results, nil
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func evalPair(pair [2]Entry) (float64, float64) {
	sumLeft := 0.0
	sumRight := 0.0

	n := 0
	m := 0

	for i, v := range pair[0].logProbs {
		if (i < len(pair[1].tokens)) && (pair[0].tokens[i] == pair[1].tokens[i]) {
			continue
		}

		sumLeft += v
		n += 1
	}

	for i, v := range pair[1].logProbs {
		if (i < len(pair[0].tokens)) && (pair[1].tokens[i] == pair[0].tokens[i]) {
			continue
		}

		sumRight += v
		m += 1
	}

	// avgLeft := sumLeft / float64(len(pair[0].logProbs))
	avgLeft := sumLeft / float64(n)
	// avgRight := sumRight / float64(len(pair[1].logProbs))
	avgRight := sumRight / float64(m)

	pLeft := math.Exp(avgLeft)
	pRight := math.Exp(avgRight)

	pLeftNormalized := pLeft / (pLeft + pRight)
	pRightNormalized := pRight / (pLeft + pRight)

	return pLeftNormalized, pRightNormalized
}

func cumulate(predictions Predictions) ([]string, []bool, [][2]float64) {
	bases := make([]string, 200)
	binary := make([]bool, 200)
	ratio := make([][2]float64, 200)

	for _, v := range predictions {
		for _, w := range v {
			for i, pair := range w {
				bases[i] = pair[1].text[:len(pair[1].text)-4]

				a, b := evalPair(pair)

				if a > b {
					binary[i] = true // true for "ity"; false for "ness"
				}

				ratio[i][0] += a
				ratio[i][1] += b
			}
		}
	}

	for i := range ratio {
		total := ratio[i][0] + ratio[i][1]

		ratio[i][0] /= total
		ratio[i][1] /= total
	}

	return bases, binary, ratio
}

func surveyResponses(name string) (map[string]float64, map[string]bool, error) {
	var data map[string]interface{}

	if err := fromJSON(name, &data); err != nil {
		return nil, nil, err
	}

	ratios := make(map[string]float64, 200)
	binary := make(map[string]bool, 200)

	keys := getKeys(data)

	for _, key := range keys {
		responsesData := data[key].(map[string]interface{})

		ity := 0
		ness := 0

		for _, v := range responsesData {
			s, ok := v.(string)

			if !ok {
				continue
			}

			if s == "ity" {
				ity += 1

				continue
			}

			ness += 1
		}

		ratios[key] = float64(ity) / float64(ity+ness) // 1.0 for all "ity"
		binary[key] = ratios[key] > 0.5                // true for "ity"; false for "ness"
	}

	return ratios, binary, nil
}
