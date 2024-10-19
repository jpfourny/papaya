package stream

import (
	"github.com/jpfourny/papaya/internal/assert"
	"github.com/jpfourny/papaya/pkg/pair"
	"maps"
	"slices"
	"testing"
)

func TestFromIterSeq(t *testing.T) {
	s := FromIterSeq(slices.Values([]int{1, 2, 3}))
	got := CollectSlice(s)
	want := []int{1, 2, 3}
	assert.ElementsMatch(t, got, want)
}

func TestFromIterSeq2(t *testing.T) {
	s := FromIterSeq2(maps.All(map[int]int{1: 2, 3: 4, 5: 6}))
	got := CollectSlice(s)
	want := []pair.Pair[int, int]{
		pair.Of(1, 2),
		pair.Of(3, 4),
		pair.Of(5, 6),
	}
	assert.ElementsMatch(t, got, want)
}
