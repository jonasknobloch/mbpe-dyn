package mbpe

import (
	"bufio"
	"container/heap"
	"context"
	"fmt"
	"os"
	"sort"
	"time"
)

type MBPETrainer struct {
	preTokenizer PreTokenizer
	segmenter    Segmenter
	model        *MBPE
	vocabSize    int
	dict         *Dict
}

func NewMBPETrainer(preTokenizer PreTokenizer, segmenter Segmenter, model *MBPE, vocabSize int) *MBPETrainer {
	return &MBPETrainer{
		preTokenizer: preTokenizer,
		segmenter:    segmenter,
		model:        model,
		vocabSize:    vocabSize,
		dict:         NewDict(),
	}
}

func (t *MBPETrainer) Model() *MBPE {
	return t.model
}

func (t *MBPETrainer) Dict() *Dict {
	return t.dict
}

func (t *MBPETrainer) InitDict(names ...string) error {
	var lines int

	if l, err := countLines(names...); err != nil {
		return err
	} else {
		lines = l
	}

	pb := NewProgressBar("Pre-process files", 20, lines, time.Now())

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	done := make(chan struct{})

	go func(ctx context.Context) {
	main:
		for {
			select {
			case <-ctx.Done():
				break main
			default:
				time.Sleep(time.Second * 1)

				l := t.dict.Lines()

				pb.Update(l)
				pb.Print()

				if l >= lines {
					break main
				}
			}
		}

		pb.Finish()

		close(done)
	}(ctx)

	err := t.dict.ProcessFiles(names...)

	<-done

	return err
}

func (t *MBPETrainer) LoadDict(name string) error {
	return t.dict.Load(name)
}

func (t *MBPETrainer) AddToken(left, right string) {
	t.model.AddToken(left + right)
	t.model.AddMerge(left, right)
}

func (t *MBPETrainer) Train() {
	t.model.InitVocab(t.vocabSize)

	t.InitVocab()

	k := t.model.Cap() - t.model.Len()

	t.model.InitMerges(k)

	chunks := t.dict.Items()

	if t.segmenter != nil {
		pbSplit := NewProgressBar("Segment chunks", 20, len(chunks), time.Now())

		for i := range chunks {
			segments, alpha := SegmentWithoutPrefixWhitespace(chunks[i].src, t.segmenter)

			chunks[i].Split(segments)
			chunks[i].Alpha(alpha)

			// chunks[i].Invert()

			pbSplit.Increment()
			pbSplit.Print()
		}

		pbSplit.Finish()
	}

	var mergeWeights = make(map[Pair]float64) // pair_counts
	var pairPositions = make(map[Pair][]int)  // where_to_update

	pbPairs := NewProgressBar("Initialize pairs", 20, len(chunks), time.Now())

	epsilon := float64(0)

	for i, chunk := range chunks {
		pairs, weights, delta := chunk.WeightedPairs()

		epsilon += delta

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

		pbPairs.Increment()
		pbPairs.Print()
	}

	pbPairs.Finish()

	queue := NewQueue(make([]Merge, 0))

	for pair, positions := range pairPositions {
		ids := [2]int{
			t.model.atoi[pair[0]],
			t.model.atoi[pair[1]],
		}

		heap.Push(queue, Merge{
			pair:      pair,
			ids:       ids,
			weight:    mergeWeights[pair],
			positions: positions,
		})

		delete(pairPositions, pair)
	}

	pbMerges := NewProgressBar("Compute merges", 20, t.vocabSize-t.model.Len(), time.Now())

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

			changes, delta := chunk.TrackedMerge(top)

			epsilon += delta

			for pair, change := range changes {
				mergeWeights[pair] += change.delta

				// After updating counts (weights), we need to collect postions for all potential merges.
				// With alpha set to zero (reference), checking if delta > 0 is sufficient, since pair counts
				// can only decrease - meaning only new pairs have a positive count.
				// In fact the Hugging Face implementation checks for count > 0.

				// Alpha > 0 enables positive weight updates, as well as the removal and addition of new pairs
				// with a zero count. Sticking with checking for count > 0 could then lead to queueing duplicated
				// merge operations.

				if change.remove || change.update {
					continue // don't queue removals and weight updates
				}

				if _, ok := pairPositions[pair]; !ok {
					pairPositions[pair] = make([]int, 0)
				}

				pairPositions[pair] = append(pairPositions[pair], position)
			}
		}

		for pair := range pairPositions {
			id := [2]int{
				t.model.atoi[pair[0]],
				t.model.atoi[pair[1]],
			}

			merge := Merge{
				pair:      pair,
				ids:       id,
				weight:    mergeWeights[pair],
				positions: pairPositions[pair],
			}

			heap.Push(queue, merge)

			delete(pairPositions, pair)
		}

		pbMerges.Increment()
		pbMerges.Print()
	}

	pbMerges.Finish()

	fmt.Println("\nepsilon", epsilon)
}

func (t *MBPETrainer) InitVocab() {
	tokens := make(map[string]int)

	// TODO side effects of byte replacements on strings.Split(..., "")

	for _, chunk := range t.dict.items {
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

func saveMergeFrame(name string, frame []Merge) error {
	if _, err := os.Stat("merge_frames"); os.IsNotExist(err) {
		_ = os.Mkdir("merge_frames", os.ModePerm)
	}

	if err := toFile(name, func(writer *bufio.Writer) error {
		for _, merge := range frame {
			if _, err := writer.WriteString(fmt.Sprintf("%s (%d) %s (%d) %f\n", merge.pair[0], merge.ids[0], merge.pair[1], merge.ids[1], merge.weight)); err != nil {
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

		fmt.Printf("%s (%d) %s (%d) %f\n", merge.pair[0], merge.ids[0], merge.pair[1], merge.ids[1], merge.weight)
	}
}
