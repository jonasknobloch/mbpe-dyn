package main

import (
	"fmt"
	"log"
	"math"
	"slices"
	"strconv"
)

func main() {
	paper()
	// babyLM()
	// babyLM2()
}

func paper() {
	fmt.Println("Table 4:")
	table4()
	fmt.Println("\nFigure 5a:")
	figure5a()
	fmt.Println("\nFigure 5b:")
	figure5b()
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

func babyLM() {
	paths, stubs := walkResultsStatic("data/wug_results/out/gpt2_%d_%s%s_babylm_v2_ity_ness_nonce.json")

	fmt.Printf("vocab,prefix,alpha,able,ish,ive,ous,able_err,ish_err,ive_err,ous_err\n")

	for i, path := range paths {
		fmt.Printf("%s,%s,%s", stubs[i][0], stubs[i][1], stubs[i][2])

		results, deviations := againstGold(path, [][]string{able[:], ish[:], ive[:], ous[:]}, []float64{-1}, 1) // set group size 12 to average across prompts per nonce adjective

		for _, v := range results {
			fmt.Printf(",%.3f", v[0])
		}

		for _, v := range deviations {
			fmt.Printf(",%.3f", v[0])
		}

		fmt.Println()
	}
}

func babyLM2() {
	paths, stubs := walkResultsStatic("data/wug_results/out/gpt2_%d_%s%s_babylm_v2_ity_ness_nonce.json")

	ratios, _, _, _ := surveyResponses("data/wug_results/survey_responses.json")

	columns := getKeys(toSet(ratios))

	slices.Reverse(columns)

	columns = append(columns, -1)

	fmt.Printf("vocab,prefix,alpha,")

	for _, c := range columns {
		if c == -1 {
			fmt.Println("average")

			continue
		}

		fmt.Printf("%.2f,", c)
	}

	for i, path := range paths {
		fmt.Printf("%s,%s,%s", stubs[i][0], stubs[i][1], stubs[i][2])

		results, _ := againstGold(path, [][]string{nonce}, columns, 1) // set group size 12 to average across prompts per nonce adjective

		for _, v := range results[len(results)-1] {
			fmt.Printf(",%.2f", v)
		}

		fmt.Println()
	}
}

func againstGold(name string, adjectives [][]string, ratios []float64, groupSize int) ([][]float64, [][]float64) {
	ratiosGold, binaryGold, keys, err := surveyResponses("data/wug_results/survey_responses.json")

	if err != nil {
		log.Fatal(err)
	}

	predictions, err := processPredictions(name)

	if err != nil {
		log.Fatal(err)
	}

	ratiosPred, binaryPred := evalPredictions(predictions, groupSize)

	groups := len(binaryPred) / len(binaryGold)

	r := make([][]float64, 0)
	e := make([][]float64, 0)

	for _, adj := range adjectives {
		acc := make([]float64, 0)
		dev := make([]float64, 0)

		allowed := toSet(adj)

		for _, ratio := range ratios {
			p := 0
			n := 0

			totalError := 0.0

			for i, key := range keys {
				if key != nonce[i] {
					panic("unexpected nonce adjective: " + key)
				}

				if _, ok := allowed[key]; !ok {
					continue
				}

				if ratio != -1 && ratiosGold[i] != ratio {
					continue
				}

				for j := 0; j < groups; j++ {
					if binaryPred[(i*groups)+j] == binaryGold[i] {
						p += 1
					} else {
						n += 1
					}

					totalError += math.Abs(ratiosPred[(i*groups)+j] - ratiosGold[i])
				}
			}

			acc = append(acc, float64(p)/float64(p+n))
			dev = append(dev, totalError/float64(p+n))
		}

		r = append(r, acc)
		e = append(e, dev)
	}

	return r, e
}

func table4() {
	acc, _ := againstGold("data/wug_results/gptj_predictions_nonce.json", [][]string{able[:], ish[:], ive[:], ous[:]}, []float64{-1}, 1)

	for _, v := range acc {
		fmt.Printf("%.3f\n", v[0])
	}

	return
}

func figure5a() {
	predictions, err := processPredictions("data/wug_results/gptj_predictions_nonce.json")

	if err != nil {
		log.Fatal(err)
	}

	_, binary := evalPredictions(predictions, 1)

	groups := len(binary) / len(nonce)

	for _, adj := range [][]string{able[:], ish[:], ive[:], ous[:]} {
		allowed := toSet(adj)

		p := 0
		n := 0

		for i, key := range nonce {
			if _, ok := allowed[key]; !ok {
				continue
			}

			for j := 0; j < groups; j++ {
				if binary[(i*groups)+j] {
					p++
				} else {
					n++
				}
			}
		}

		fmt.Printf("%.3f\n", 1-(float64(p)/float64(p+n)))
	}
}

func figure5b() {
	_, binaryGold, keys, _ := surveyResponses("data/wug_results/survey_responses.json")

	for _, adj := range [][]string{able[:], ish[:], ive[:], ous[:]} {
		allowed := toSet(adj)

		p := 0
		n := 0

		for i, key := range keys {
			if _, ok := allowed[key]; !ok {
				continue
			}

			if binaryGold[i] {
				p++
			} else {
				n++
			}
		}

		fmt.Printf("%.3f\n", 1-(float64(p)/float64(p+n)))
	}
}
