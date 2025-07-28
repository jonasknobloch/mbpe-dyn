package main

import (
	"fmt"
	"log"
	"math"
)

func main() {
	fmt.Println(againstGold("data/wug_results/gptj_predictions_nonce.json"))

	// return

	fmt.Println()

	// fmt.Println("Table 4:")
	// table4()
	// fmt.Println("\nFigure 5b:")
	// figure5b()
	// fmt.Println("\nFigure 5a:")
	// figure5a()
	// fmt.Println("\nFigure 5a (raw):")
	// figure5aRaw()
	//
	// fmt.Println()

	paths, stubs := walkResultsStatic("data/wug_results/out/gpt2_%d_%s%s_babylm_v2_ity_ness_nonce.json")

	fmt.Printf("vocab,prefix,alpha,able,ish,ive,ous\n")

	for i, path := range paths {
		fmt.Printf("%s,%s,%s", stubs[i][0], stubs[i][1], stubs[i][2])

		results := againstGold(path)

		for _, v := range results {
			fmt.Printf(",%.3f", v)
		}

		fmt.Println()
	}
}

func againstGold(name string) []float64 {
	_, binaryGold, keys, err := surveyResponses("data/wug_results/survey_responses.json")

	if err != nil {
		log.Fatal(err)
	}

	predictions, err := processPredictions(name)

	if err != nil {
		log.Fatal(err)
	}

	_, binary := evalPredictions(predictions, 1)
	groups := len(binary) / len(binaryGold)

	r := make([]float64, 0)

	for _, adj := range [][]string{able[:], ish[:], ive[:], ous[:]} {
		p := 0
		n := 0

		allowed := toSet(adj)

		for i, key := range keys {
			if key != nonce[i] {
				panic("unexpected nonce adjective: " + key)
			}

			if _, ok := allowed[key]; !ok {
				continue
			}

			// fmt.Println("Processing key:", key, "at index:", i)

			for j := 0; j < groups; j++ {
				// fmt.Println("Processing pair:", j)

				if binary[(i*groups)+j] == binaryGold[i] {
					p += 1
				} else {
					n += 1
				}
			}
		}

		r = append(r, float64(p)/float64(p+n))
	}

	return r
}

func table4() {
	ratioGold, binaryGold, keys, err := surveyResponses("data/wug_results/survey_responses.json")

	if err != nil {
		log.Fatal(err)
	}

	predictions, err := processPredictions("data/wug_results/gptj_predictions_nonce.json")

	if err != nil {
		log.Fatal(err)
	}

	ratio, binary := evalPredictions(predictions, 12)

	for _, adj := range [][]string{able[:], ish[:], ive[:], ous[:]} {
		p := 0
		n := 0

		errSum := 0.0

		allowed := toSet(adj)

		for i, key := range keys {
			if key != nonce[i] {
				panic("unexpected nonce adjective: " + key)
			}

			if _, ok := allowed[key]; !ok {
				continue
			}

			if binary[i] == binaryGold[i] {
				p += 1
			} else {
				n += 1
			}

			errSum += math.Abs(ratio[i] - ratioGold[i])
		}

		fmt.Printf("%.3f (%.3f)\n", float64(p)/float64(len(allowed)), errSum/float64(len(allowed)))
	}
}

func figure5a() {
	predictions, err := processPredictions("data/wug_results/gptj_predictions_nonce.json")

	if err != nil {
		log.Fatal(err)
	}

	_, binary := evalPredictions(predictions, 12)

	for _, adj := range [][]string{able[:], ish[:], ive[:], ous[:]} {
		allowed := toSet(adj)

		p := 0
		n := 0

		for i, key := range nonce {
			if _, ok := allowed[key]; !ok {
				continue
			}

			if binary[i] {
				p += 1
			} else {
				n += 1
			}
		}

		fmt.Printf("%.3f\n", 1-(float64(p)/float64(len(allowed))))
	}
}

func figure5aRaw() {
	predictions, err := processPredictions("data/wug_results/gptj_predictions_nonce.json")

	if err != nil {
		log.Fatal(err)
	}

	s := predictionsToSlice(predictions)

	for _, adj := range [][]string{able[:], ish[:], ive[:], ous[:]} {
		allowed := toSet(adj)

		p := 0
		n := 0

		for i, key := range nonce {
			if _, ok := allowed[key]; !ok {
				continue
			}

			houston := s[i]

			for _, entry := range houston {
				a, b := evalPair(entry)

				if a > b {
					p += 1
				} else {
					n += 1
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
				p += 1
			} else {
				n += 1
			}
		}

		fmt.Printf("%.3f\n", 1-(float64(p)/float64(len(allowed))))
	}
}

func foo() {
	adjectives := 50
	prompts := 12

	for i := 0; i < adjectives*prompts; i++ {
		fmt.Printf("%.4f %d/%d\n", float64(i)/float64(adjectives*prompts), i, adjectives*prompts)
	}
}
