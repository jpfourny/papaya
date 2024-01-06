package stream

import (
	"testing"

	"github.com/jpfourny/papaya/internal/assert"
	"github.com/jpfourny/papaya/pkg/mapper"
	"github.com/jpfourny/papaya/pkg/optional"
	"github.com/jpfourny/papaya/pkg/pred"
)

func TestIterate(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		iter := func() optional.Optional[int] {
			return optional.Empty[int]()
		}
		s := Iterate(iter)
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		i := 0
		iter := func() optional.Optional[int] {
			if i < 3 {
				i++
				return optional.Of(i)
			}
			return optional.Empty[int]()
		}
		s := Iterate(iter)
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		assert.ElementsMatch(t, got, want)
	})
}

func TestRange(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := Range(0, pred.LessThan(0), mapper.Increment(1))
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := Range(1, pred.LessThanOrEqual(5), mapper.Increment(2))
		got := CollectSlice(s)
		want := []int{1, 3, 5}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := Range(1, pred.LessThanOrEqual(5), mapper.Increment(2))
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []int{1, 3}
		assert.ElementsMatch(t, got, want)
	})
}
