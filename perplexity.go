package mbpe

import (
	"fmt"
	ort "github.com/yalue/onnxruntime_go"
	"math/rand"
	"time"
)

type Perplexity struct {
	alpha     float64
	tokenizer *Tokenizer
}

func NewPerplexity() *Perplexity {
	return &Perplexity{}
}

func (p *Perplexity) SetTokenizer(tokenizer *Tokenizer) {
	p.tokenizer = tokenizer
}

func (p *Perplexity) Segment(text string) ([]string, float64) {
	ids := toInt64Slice(p.tokenizer.Tokenize(text))

	ort.SetSharedLibraryPath("onnxruntime.so")

	if err := ort.InitializeEnvironment(); err != nil {
		panic(err)
	}

	// inputs, outputs, err := ort.GetInputOutputInfo("model_quantized.onnx")
	//
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(inputs)
	// fmt.Println(outputs)

	in, err := ort.NewTensor(ort.NewShape(1, int64(len(ids))), ids)

	if err != nil {
		panic(err)
	}

	out, err := ort.NewEmptyTensor[float32](ort.NewShape(1, 1, 50257))

	if err != nil {
		panic(err)
	}

	// TODO postion_ids
	// TODO attention mask (all 1s?)
	// TODO past values

	att := attentionMask(ids)
	pos := positionIDs(ids)

	names, tensors := pastKeyValues()

	inputNames := []string{
		"input_ids",
		names[0],
		names[1],
		names[2],
		names[3],
		names[4],
		names[5],
		names[6],
		names[7],
		names[8],
		names[9],
		names[10],
		names[11],
		names[12],
		names[13],
		names[14],
		names[15],
		names[16],
		names[17],
		names[18],
		names[19],
		names[20],
		names[21],
		names[22],
		names[23],
		"attention_mask",
		"position_ids",
	}

	inputValues := []ort.Value{
		in,
		tensors[0],
		tensors[1],
		tensors[2],
		tensors[3],
		tensors[4],
		tensors[5],
		tensors[6],
		tensors[7],
		tensors[8],
		tensors[9],
		tensors[10],
		tensors[11],
		tensors[12],
		tensors[13],
		tensors[14],
		tensors[15],
		tensors[16],
		tensors[17],
		tensors[18],
		tensors[19],
		tensors[20],
		tensors[21],
		tensors[22],
		tensors[23],
		att,
		pos,
	}

	session, err := ort.NewAdvancedSession("gpt2/model_quantized.onnx", inputNames, []string{"logits"}, inputValues, []ort.Value{out}, nil)

	if err != nil {
		panic(err)
	}

	if err := session.Run(); err != nil {
		panic(err)
	}

	foo := out.GetData()

	_, id := findMinValue(foo)

	tokens := p.tokenizer.model.(*MBPE).ToString([]int{id})
	decoded := p.tokenizer.Decoder().Decode(tokens)

	fmt.Println(id)
	fmt.Println(tokens)
	fmt.Println(decoded)

	// TODO calc perplexity of all substrings starting at zero index -> track how PPL changes

	return []string{text}, p.alpha
}

func toInt64Slice(s []int) []int64 {
	r := make([]int64, len(s))

	for i, v := range s {
		r[i] = int64(v)
	}

	return r
}

func positionIDs(ids []int64) *ort.Tensor[int64] {
	pos := make([]int64, len(ids))

	for i := range pos {
		pos[i] = int64(i)
	}

	t, err := ort.NewTensor(ort.NewShape(1, int64(len(ids))), pos)

	if err != nil {
		panic(err)
	}

	return t
}

func attentionMask(ids []int64) *ort.Tensor[int64] {
	att := make([]int64, len(ids))

	for i := range att {
		att[i] = int64(1)
	}

	t, err := ort.NewTensor(ort.NewShape(1, int64(len(ids))), att)

	if err != nil {
		panic(err)
	}

	return t
}

func generateFloat32Array(size int) []float32 {
	rand.Seed(time.Now().UnixNano())

	r := make([]float32, size)

	for i := range r {
		r[i] = 0.0
	}

	return r
}

func pastKeyValues() ([]string, []*ort.Tensor[float32]) {
	names := make([]string, 0, 24)
	tensors := make([]*ort.Tensor[float32], 0, 24)
	shape := ort.NewShape(1, 12, 1, 64)

	for l := 0; l < 12; l++ {
		for _, kv := range []string{"key", "value"} {
			name := fmt.Sprintf("past_key_values.%d.%s", l, kv)

			pastData := generateFloat32Array(12 * 64)

			tensor, err := ort.NewTensor(shape, pastData)

			if err != nil {
				panic(err)
			}

			names = append(names, name)
			tensors = append(tensors, tensor)
		}
	}

	return names, tensors
}

func findMinValue(outputData []float32) (float32, int) {
	minValue := float32(0)
	minIndex := -1

	for i, v := range outputData {
		if v < minValue {
			minValue = v
			minIndex = i
		}
	}

	return minValue, minIndex
}
