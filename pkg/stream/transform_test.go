package stream

import (
	"fmt"
	"testing"

	"github.com/jpfourny/papaya/internal/assert"
	"github.com/jpfourny/papaya/pkg/mapper"
)

func TestMap(t *testing.T) {
	s := Map(Of(1, 2, 3), mapper.Sprintf[int]("%d"))
	got := CollectSlice(s)
	want := []string{"1", "2", "3"}
	assert.ElementsMatch(t, got, want)
}

func TestMapOrDiscard(t *testing.T) {
	s := MapOrDiscard(Of("1", "2", "3", "foo"), mapper.TryParseInt[string, int](10, 64))
	got := CollectSlice(s)
	want := []int{1, 2, 3}
	assert.ElementsMatch(t, got, want)
}

func TestFlatMap(t *testing.T) {
	s := FlatMap(Of(1, 2, 3), func(e int) Stream[string] {
		return Of(fmt.Sprintf("%dA", e), fmt.Sprintf("%dB", e))
	})
	got := CollectSlice(s)
	want := []string{"1A", "1B", "2A", "2B", "3A", "3B"}
	assert.ElementsMatch(t, got, want)
}

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
