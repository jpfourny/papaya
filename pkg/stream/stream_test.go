package stream

import (
	"testing"

	"github.com/jpfourny/papaya/v2/internal/assert"
)

func TestEmpty(t *testing.T) {
	s := Empty[int]()
	got := CollectSlice(s)
	var want []int
	assert.ElementsMatch(t, got, want)
}

func TestOf(t *testing.T) {
	s := Of(1, 2, 3)
	got := CollectSlice(s)
	want := []int{1, 2, 3}
	assert.ElementsMatch(t, got, want)
}
