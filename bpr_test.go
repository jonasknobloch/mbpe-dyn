package mbpe

import (
	"fmt"
	"testing"
)

func TestEvalSegmentations(t *testing.T) {
	pred := [][]string{
		{
			"ab", "lative", "s",
		},
	}

	gold := []string{
		"ablative", "s",
	}

	_, counts := evalSegmentations(pred, gold, "ablatives", MaxF1)

	tp := counts[0] // 1
	fp := counts[1] // 1
	tn := counts[2] // 6
	fn := counts[3] // 0

	precision := float64(tp) / float64(tp+fp)
	recall := float64(tp) / float64(tp+fn)

	f1 := 2 * precision * recall / (precision + recall)

	fmt.Println(tp, fp, tn, fn)

	if precision != 0.5 || recall != 1.0 || f1 != 0.6666666666666666 {
		t.Errorf("expected [0.5, 1.0, 0.666667] but got [%f, %f, %f]", precision, recall, f1)
	}
}
