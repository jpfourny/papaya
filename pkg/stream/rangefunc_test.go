package stream

import (
	"github.com/jpfourny/papaya/internal/assert"
	"github.com/jpfourny/papaya/pkg/pair"
	"testing"
)

func TestFromRangeFunc(t *testing.T) {
	f := func(yield func(int) bool) bool {
		return yield(1) && yield(2) && yield(3)
	}
	s := FromRangeFunc(f)
	got := CollectSlice(s)
	want := []int{1, 2, 3}
	assert.ElementsMatch(t, got, want)
}

func TestFromRangeBiFunc(t *testing.T) {
	f := func(yield func(int, int) bool) bool {
		return yield(1, 2) && yield(3, 4) && yield(5, 6)
	}
	s := FromRangeBiFunc(f)
	got := CollectSlice(s)
	want := []pair.Pair[int, int]{
		pair.Of(1, 2),
		pair.Of(3, 4),
		pair.Of(5, 6),
	}
	assert.ElementsMatch(t, got, want)
}

func TestToRangeFunc(t *testing.T) {
	s := Of(1, 2, 3)
	f := ToRangeFunc(s)
	got := make([]int, 0, 3)
	want := []int{1, 2, 3}
	f(func(e int) bool {
		got = append(got, e)
		return true
	})
	assert.ElementsMatch(t, got, want)
}

func TestToRangeBiFunc(t *testing.T) {
	s := Of(
		pair.Of(1, 2),
		pair.Of(3, 4),
		pair.Of(5, 6),
	)
	f := ToRangeBiFunc(s)
	got := make([]pair.Pair[int, int], 0, 3)
	want := []pair.Pair[int, int]{
		pair.Of(1, 2),
		pair.Of(3, 4),
		pair.Of(5, 6),
	}
	f(func(k, v int) bool {
		got = append(got, pair.Of(k, v))
		return true
	})
	assert.ElementsMatch(t, got, want)
}
