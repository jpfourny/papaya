package stream

import (
	"github.com/jpfourny/papaya/internal/assert"
	"testing"
)

func TestSortAsc(t *testing.T) {
	s := SortAsc(Of(3, 1, 2))
	got := CollectSlice(s)
	want := []int{1, 2, 3}
	assert.ElementsMatch(t, got, want)
}

func TestSortDesc(t *testing.T) {
	s := SortDesc(Of(3, 1, 2))
	got := CollectSlice(s)
	want := []int{3, 2, 1}
	assert.ElementsMatch(t, got, want)
}

func TestSortBy(t *testing.T) {
	s := SortBy(Of(3, 1, 2), func(a, b int) int {
		return a - b
	})
	got := CollectSlice(s)
	want := []int{1, 2, 3}
	assert.ElementsMatch(t, got, want)
}
