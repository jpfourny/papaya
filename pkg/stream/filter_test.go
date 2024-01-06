package stream

import (
	"github.com/jpfourny/papaya/pkg/cmp"
	"testing"

	"github.com/jpfourny/papaya/internal/assert"
)

func TestFilter(t *testing.T) {
	s := Filter(Of(1, 2, 3, 4, 5), func(e int) bool {
		return e%2 == 0
	})
	got := CollectSlice(s)
	want := []int{2, 4}
	assert.ElementsMatch(t, got, want)
}

func TestLimit(t *testing.T) {
	t.Run("limit-0", func(t *testing.T) {
		s := Limit(Of(1, 2, 3, 4, 5), 0)
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limit-negative", func(t *testing.T) {
		s := Limit(Of(1, 2, 3, 4, 5), -1) // Will be treated as 0.
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limit-positive", func(t *testing.T) {
		s := Limit(Of(1, 2, 3, 4, 5), 3)
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		assert.ElementsMatch(t, got, want)
	})
}

func TestSkip(t *testing.T) {
	t.Run("skip-0", func(t *testing.T) {
		s := Skip(Of(1, 2, 3, 4, 5), 0)
		got := CollectSlice(s)
		want := []int{1, 2, 3, 4, 5}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("skip-negative", func(t *testing.T) {
		s := Skip(Of(1, 2, 3, 4, 5), -1) // Will be treated as 0.
		got := CollectSlice(s)
		want := []int{1, 2, 3, 4, 5}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("skip-positive", func(t *testing.T) {
		s := Skip(Of(1, 2, 3, 4, 5), 3)
		got := CollectSlice(s)
		want := []int{4, 5}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("skip-all", func(t *testing.T) {
		s := Skip(Of(1, 2, 3, 4, 5), 6)
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})
}

func TestDistinct(t *testing.T) {
	s := Distinct(Of(1, 2, 3, 2, 1))
	got := CollectSlice(s)
	want := []int{1, 2, 3}
	assert.ElementsMatch(t, got, want)
}

func TestDistinctBy(t *testing.T) {
	s := DistinctBy(Of(1, 2, 3, 2, 1), cmp.Natural[int]())
	got := CollectSlice(s)
	want := []int{1, 2, 3}
	assert.ElementsMatch(t, got, want)
}
