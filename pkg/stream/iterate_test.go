package stream

import (
	"github.com/jpfourny/papaya/pkg/stream/mapper"
	"github.com/jpfourny/papaya/pkg/stream/pred"
	"testing"

	"github.com/jpfourny/papaya/internal/assert"
	"github.com/jpfourny/papaya/pkg/opt"
)

func TestIterate(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		iter := func() (int, bool) {
			return 0, false
		}
		s := Iterate(iter)
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		i := 0
		iter := func() (int, bool) {
			if i < 3 {
				i++
				return i, true
			}
			return 0, false
		}
		s := Iterate(iter)
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		assert.ElementsMatch(t, got, want)
	})
}

func TestIterateOptional(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		iter := func() opt.Optional[int] {
			return opt.Empty[int]()
		}
		s := IterateOptional(iter)
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		i := 0
		iter := func() opt.Optional[int] {
			if i < 3 {
				i++
				return opt.Of(i)
			}
			return opt.Empty[int]()
		}
		s := IterateOptional(iter)
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		i := 0
		iter := func() opt.Optional[int] {
			if i < 3 {
				i++
				return opt.Of(i)
			}
			return opt.Empty[int]()
		}
		s := Limit(IterateOptional(iter), 2)
		got := CollectSlice(s)
		want := []int{1, 2}
		assert.ElementsMatch(t, got, want)
	})
}

func TestWalk(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := Walk(0, pred.LessThan(0), mapper.Increment(1))
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := Walk(1, pred.LessThanOrEqual(5), mapper.Increment(2))
		got := CollectSlice(s)
		want := []int{1, 3, 5}
		assert.ElementsMatch(t, got, want)
	})
}

func TestInterval(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := Interval(0, 0, 1)
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty-increasing", func(t *testing.T) {
		s := Interval(1, 5, 2)
		got := CollectSlice(s)
		want := []int{1, 3}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty-decreasing", func(t *testing.T) {
		s := Interval(5, 1, -2)
		got := CollectSlice(s)
		want := []int{5, 3}
		assert.ElementsMatch(t, got, want)
	})
}
