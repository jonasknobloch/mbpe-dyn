package mbpe

import (
	"reflect"
	"testing"
)

type baseUniformizer struct{}

func (s *baseUniformizer) Segment(text string) ([]string, float64) {
	return []string{"foo", "bar"}, 1.0
}

func TestUniformizer_Segment(t *testing.T) {
	u := NewUniformizer(&baseUniformizer{})

	segmentation, _ := u.Segment("foobarbaz")

	expected := []string{"foob", "arbaz"}

	if !reflect.DeepEqual(segmentation, expected) {
		t.Errorf("expected %v but got %v", expected, segmentation)
	}
}
