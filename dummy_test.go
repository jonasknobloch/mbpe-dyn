package main

import (
	"reflect"
	"testing"
)

type testSegmenter struct{}

func (s *testSegmenter) Segment(text string) ([]string, float64) {
	return []string{"foo", "bar"}, 1.0
}

func TestDummySegmenter_Segment(t *testing.T) {
	d := NewDummySegmenter(&testSegmenter{})

	segmentation, _ := d.Segment("foobarbaz")

	expected := []string{"foob", "arbaz"}

	if !reflect.DeepEqual(segmentation, expected) {
		t.Errorf("expected %v but got %v", expected, segmentation)
	}
}
