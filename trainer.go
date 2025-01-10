package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"mbpe-dyn/morfessor"
	"os"
	"sort"
	"strings"
	"time"
)

type MBPETrainer struct {
	preTokenizer PreTokenizer
	model        *MBPE
	celex        *CELEX
	morfessor    *morfessor.Model
	alpha        float64
	vocabSize    int
	dict         map[string]*Chunk
}

func NewMBPETrainer(preTokenizer PreTokenizer, model *MBPE, celex *CELEX, morfessor *morfessor.Model, alpha float64, vocabSize int) *MBPETrainer {
	if alpha < 0 || alpha > 1 {
		panic("alpha must be in [0, 1]")
	}

	return &MBPETrainer{
		preTokenizer: preTokenizer,
		model:        model,
		celex:        celex,
		morfessor:    morfessor,
		alpha:        alpha,
		vocabSize:    vocabSize,
		dict:         make(map[string]*Chunk),
	}
}

func (t *MBPETrainer) InitDict(name string) error {
	file, err := os.Open(name)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := strings.IndexAny(string(data), "\r\n"); i >= 0 {
			if i+1 < len(data) && data[i] == '\r' && data[i+1] == '\n' {
				return i + 2, data[0 : i+2], nil
			}

			return i + 1, data[0 : i+1], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	})

	for scanner.Scan() {
		line := scanner.Text()

		if err := scanner.Err(); err != nil {
			return err
		}

		compounds := t.preTokenizer.PreTokenize(line)

		for _, compound := range compounds {
			if _, ok := t.dict[compound]; !ok {
				split, ok := t.celex.Split(compound)

				if !ok {
					split, _ = t.morfessor.Segment(compound)
				}

				t.dict[compound] = NewChunk(compound, 0, split, t.alpha)
			}

			t.dict[compound].n += 1
		}
	}

	return nil
}

func (t *MBPETrainer) AddToken(left, right string) {
	t.model.AddToken(left + right)
	t.model.AddMerge(left, right)
}

func (t *MBPETrainer) Train(name string) error {
	pb0 := NewProgressBar("Pre-process files", 40, 1, time.Now())

	if err := t.InitDict(name); err != nil {
		return err
	}

	pb0.Update()
	pb0.Finish()

	if err := t.SaveDict("dict.txt"); err != nil {
		return err
	}

	t.model.InitVocab(t.vocabSize)

	t.InitVocab()

	k := t.model.Cap() - t.model.Len()

	t.model.InitMerges(k)

	var chunks = make([]Chunk, 0, len(t.dict))

	for _, chunk := range t.dict {
		chunks = append(chunks, *chunk)
	}

	var mergeWeights = make(map[Pair]float64) // pair_counts
	var pairPositions = make(map[Pair][]int)  // where_to_update

	pb1 := NewProgressBar("Initialize pairs", 40, len(chunks), time.Now())

	for i, chunk := range chunks {
		pairs, weights := chunk.Pairs()

		for j, pair := range pairs {
			if _, ok := mergeWeights[pair]; !ok {
				mergeWeights[pair] = 0
			}

			if _, ok := pairPositions[pair]; !ok {
				pairPositions[pair] = make([]int, 0)
			}

			mergeWeights[pair] += weights[j]
			pairPositions[pair] = append(pairPositions[pair], i)
		}

		pb1.Update()
	}

	pb1.Finish()

	queue := NewQueue(make([]Merge, 0))

	for pair, positions := range pairPositions {
		idx := [2]int{
			t.model.atoi[pair[0]],
			t.model.atoi[pair[1]],
		}

		heap.Push(queue, Merge{
			pair:      pair,
			idx:       idx,
			weight:    mergeWeights[pair],
			positions: positions,
		})

		delete(pairPositions, pair)
	}

	pb2 := NewProgressBar("Compute merges", 40, t.vocabSize-t.model.Len(), time.Now())

	for t.model.Len() < t.vocabSize && queue.Len() != 0 {
		top := heap.Pop(queue).(Merge)

		if top.weight != mergeWeights[top.pair] {
			top.weight = mergeWeights[top.pair]

			queue.Push(top)

			continue
		}

		// dumpQueue(*queue, top, mergeWeights)

		if top.weight < 1 {
			break
		}

		t.AddToken(top.pair[0], top.pair[1])

		for _, position := range top.positions {
			chunk := &chunks[position]

			for pair, change := range chunk.TrackedMerge(top) {
				mergeWeights[pair] += change.delta

				if change.delta <= 0 || change.update {
					continue // don't queue removals and positive weight updates
				}

				if _, ok := pairPositions[pair]; !ok {
					pairPositions[pair] = make([]int, 0)
				}

				pairPositions[pair] = append(pairPositions[pair], position)
			}
		}

		for pair := range pairPositions {
			idx := [2]int{
				t.model.atoi[pair[0]],
				t.model.atoi[pair[1]],
			}

			merge := Merge{
				pair:      pair,
				idx:       idx,
				weight:    mergeWeights[pair],
				positions: pairPositions[pair],
			}

			heap.Push(queue, merge)

			delete(pairPositions, pair)
		}

		pb2.Update()
	}

	pb2.Finish()

	return nil
}

func (t *MBPETrainer) InitVocab() {
	tokens := make(map[string]int)

	// TODO side effects of byte replacements on strings.Split(..., "")

	for _, chunk := range t.dict {
		for _, token := range chunk.Tokens() {
			if _, ok := tokens[token]; !ok {
				tokens[token] = 0

				// fmt.Printf("discovered new token %s %02x \n", token, []byte(token))
			}

			tokens[token]++
		}
	}

	alphabet := make([]string, 0, len(tokens))

	for token := range tokens {
		alphabet = append(alphabet, token)
	}

	sort.Strings(alphabet)

	for _, token := range alphabet {
		t.model.AddToken(token)
	}
}

func (t *MBPETrainer) SaveDict(name string) error {
	order := make([]string, 0, len(t.dict))

	for token := range t.dict {
		order = append(order, token)
	}

	sort.Slice(order, func(i, j int) bool {
		if t.dict[order[i]].n != t.dict[order[j]].n {
			return t.dict[order[i]].n > t.dict[order[j]].n
		}

		return t.dict[order[i]].src < t.dict[order[j]].src
	})

	if err := toFile(name, func(writer *bufio.Writer) error {
		for _, key := range order {
			if _, err := writer.WriteString(fmt.Sprintf("%s %d\n", t.dict[key].src, t.dict[key].n)); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func saveMergeFrame(name string, frame []Merge) error {
	if _, err := os.Stat("merge_frames"); os.IsNotExist(err) {
		_ = os.Mkdir("merge_frames", os.ModePerm)
	}

	if err := toFile(name, func(writer *bufio.Writer) error {
		for _, merge := range frame {
			if _, err := writer.WriteString(fmt.Sprintf("%s (%d) %s (%d) %f\n", merge.pair[0], merge.idx[0], merge.pair[1], merge.idx[1], merge.weight)); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func dumpQueue(queue Queue, top Merge, mergeWeights map[Pair]float64) {
	list := make([]Merge, queue.Len())

	copy(list, queue)

	list = append(list, top) // TODO just use heap helpers instead of sorting the underlying slice

	for _, merge := range list {
		if merge.weight != mergeWeights[merge.pair] {
			merge.weight = mergeWeights[merge.pair]
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[j].Less(list[i])
	})

	fmt.Println()

	// _ = saveMergeFrame(fmt.Sprintf("merge_frames/frame_%d.txt", time.Now().UnixNano()), list[0:10])

	for i, merge := range list {
		if i == 10 {
			break
		}

		fmt.Printf("%s (%d) %s (%d) %f\n", merge.pair[0], merge.idx[0], merge.pair[1], merge.idx[1], merge.weight)
	}
}
