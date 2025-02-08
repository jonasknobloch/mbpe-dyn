package bpr

import (
	"reflect"
	"testing"
)

func TestBoundaryVector(t *testing.T) {
	segmentation := []string{
		"adoptio",
		"tyttäre",
		"nsä",
	}

	bounds := BoundaryVector(segmentation)

	expected := []int{
		0, 0, 0, 0, 0, 0, 1,
		0, 0, 0, 0, 0, 0, 0, 1,
		0, 0, 0,
	}

	if !reflect.DeepEqual(bounds, expected) {
		t.Errorf("expecetd %v but got %v", expected, bounds)
	}
}

func TestEval(t *testing.T) {
	gold := [][]string{
		{
			"ablative", "s",
		},
	}

	pred := [][]string{
		{
			"ab", "lative", "s",
		},
	}

	precision, recall, f1 := Eval(gold, pred)

	if precision != 0.5 || recall != 1.0 || f1 != 0.6666666666666666 {
		t.Errorf("expected [0.5, 1.0, 0.666667] but got [%f, %f, %f]", precision, recall, f1)
	}
}
