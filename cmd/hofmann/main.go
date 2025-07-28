package main

import (
	"fmt"
	"log"
	"math"
	"slices"
)

func main() {
	// paper()
	babyLM()
}

func paper() {
	fmt.Println("Table 4:")
	table4()
	fmt.Println("\nFigure 5a:")
	figure5a()
	fmt.Println("\nFigure 5b:")
	figure5b()
}

func babyLM() {
	paths, stubs := walkResultsStatic("data/wug_results/out/gpt2_%d_%s%s_babylm_v2_ity_ness_nonce.json")

	ratios, _, _, _ := surveyResponses("data/wug_results/survey_responses.json")

	columns := getKeys(toSet(ratios))

	slices.Reverse(columns)

	fmt.Printf("vocab,prefix,alpha,")

	for _, c := range columns {
		fmt.Printf("%.2f,", c)
	}

	fmt.Println("average")

	for i, path := range paths {
		fmt.Printf("%s,%s,%s", stubs[i][0], stubs[i][1], stubs[i][2])

		results, _ := againstGold2(path, columns, 1) // set group size 12 to average across prompts per nonce adjective

		for _, v := range results {
			fmt.Printf(",%.2f", v)
		}

		// for _, v := range deviations {
		// 	fmt.Printf(",%.3f", v)
		// }

		fmt.Printf(",%.2f\n", average(results))
	}
}

func againstGold(name string, groupSize int) ([]float64, []float64) {
	ratiosGold, binaryGold, keys, err := surveyResponses("data/wug_results/survey_responses.json")

	if err != nil {
		log.Fatal(err)
	}

	predictions, err := processPredictions(name)

	if err != nil {
		log.Fatal(err)
	}

	ratios, binary := evalPredictions(predictions, groupSize)

	groups := len(binary) / len(binaryGold)

	r := make([]float64, 0)
	e := make([]float64, 0)

	for _, adj := range [][]string{able[:], ish[:], ive[:], ous[:]} {
		p := 0
		n := 0

		totalError := 0.0

		allowed := toSet(adj)

		for i, key := range keys {
			if key != nonce[i] {
				panic("unexpected nonce adjective: " + key)
			}

			if _, ok := allowed[key]; !ok {
				continue
			}

			for j := 0; j < groups; j++ {
				if binary[(i*groups)+j] == binaryGold[i] {
					p += 1
				} else {
					n += 1
				}

				totalError += math.Abs(ratios[(i*groups)+j] - ratiosGold[i])
			}
		}

		r = append(r, float64(p)/float64(p+n))
		e = append(e, totalError/float64(p+n))
	}

	return r, e
}

func againstGold2(name string, ratios []float64, groupSize int) ([]float64, []float64) {
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

	r := make([]float64, 0)
	e := make([]float64, 0)

	for _, ratio := range ratios {
		p := 0
		n := 0

		totalError := 0.0

		for i, key := range keys {
			if key != nonce[i] {
				panic("unexpected nonce adjective: " + key)
			}

			if ratiosGold[i] != ratio {
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

		r = append(r, float64(p)/float64(p+n))
		e = append(e, totalError/float64(p+n))
	}

	return r, e
}

func table4() {
	acc, _ := againstGold("data/wug_results/gptj_predictions_nonce.json", 1)

	for _, v := range acc {
		fmt.Printf("%.3f\n", v)
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
