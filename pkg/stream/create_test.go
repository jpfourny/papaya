package stream

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/jpfourny/papaya/internal/assert"
	"github.com/jpfourny/papaya/pkg/mapper"
	"github.com/jpfourny/papaya/pkg/optional"
	"github.com/jpfourny/papaya/pkg/pair"
	"github.com/jpfourny/papaya/pkg/pred"
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

type emptyIter struct{}

func (_ emptyIter) Next() optional.Optional[int] {
	return optional.Empty[int]()
}

type sliceIter struct {
	slice []int
}

func (it *sliceIter) Next() optional.Optional[int] {
	if len(it.slice) == 0 {
		return optional.Empty[int]()
	}
	e := it.slice[0]
	it.slice = it.slice[1:]
	return optional.Of(e)
}

func TestFromIterator(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := FromIterator[int](emptyIter{})
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := FromIterator[int](&sliceIter{slice: []int{1, 2, 3}})
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := FromIterator[int](&sliceIter{slice: []int{1, 2, 3}})
		got := CollectSlice(Limit(s, 1)) // Stops stream after 1 element.
		want := []int{1}
		assert.ElementsMatch(t, got, want)
	})
}

func TestFromSlice(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := FromSlice([]int{})
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := FromSlice([]int{1, 2, 3})
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := FromSlice([]int{1, 2, 3})
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []int{1, 2}
		assert.ElementsMatch(t, got, want)
	})
}

func TestFromSliceWithIndex(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := FromSliceWithIndex([]int{})
		got := CollectSlice(s)
		var want []pair.Pair[int, int]
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := FromSliceWithIndex([]int{1, 2, 3})
		got := CollectSlice(s)
		want := []pair.Pair[int, int]{
			pair.Of(0, 1),
			pair.Of(1, 2),
			pair.Of(2, 3),
		}
		dd := DebugString(s)
		fmt.Println(dd)
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := FromSliceWithIndex([]int{1, 2, 3})
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []pair.Pair[int, int]{
			pair.Of(0, 1),
			pair.Of(1, 2),
		}
		assert.ElementsMatch(t, got, want)
	})
}

func TestFromMap(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := FromMap(map[int]string{})
		got := CollectSlice(s)
		var want []pair.Pair[int, string]
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := FromMap(map[int]string{1: "one", 2: "two", 3: "three"})
		got := CollectSlice(s)
		want := []pair.Pair[int, string]{
			pair.Of(1, "one"),
			pair.Of(2, "two"),
			pair.Of(3, "three"),
		}
		// All elements in map iteration order.
		assert.ElementsMatchAnyOrder(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := FromMap(map[int]string{1: "one", 2: "two", 3: "three"})
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []pair.Pair[int, string]{
			pair.Of(1, "one"),
			pair.Of(2, "two"),
			pair.Of(3, "three"),
		}
		// 2 elements returned are unpredictable due to map iteration order.
		assert.SomeElementsMatchAnyOrder(t, got, want, 2)
	})
}

func TestFromMapKeys(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := FromMapKeys(map[int]string{})
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := FromMapKeys(map[int]string{1: "one", 2: "two", 3: "three"})
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		// All elements in map iteration order.
		assert.ElementsMatchAnyOrder(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := FromMapKeys(map[int]string{1: "one", 2: "two", 3: "three"})
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []int{1, 2, 3}
		// 2 elements returned are unpredictable due to map iteration order.
		assert.SomeElementsMatchAnyOrder(t, got, want, 2)
	})
}

func TestFromMapValues(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := FromMapValues(map[int]string{})
		got := CollectSlice(s)
		var want []string
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := FromMapValues(map[int]string{1: "one", 2: "two", 3: "three"})
		got := CollectSlice(s)
		want := []string{"one", "two", "three"}
		// All elements in map iteration order.
		assert.ElementsMatchAnyOrder(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := FromMapValues(map[int]string{1: "one", 2: "two", 3: "three"})
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []string{"one", "two", "three"}
		// 2 elements returned are unpredictable due to map iteration order.
		assert.SomeElementsMatchAnyOrder(t, got, want, 2)
	})
}

func TestFromChannel(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ch := make(chan int)
		close(ch)
		s := FromChannel(ch)
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Run("unbuffered", func(t *testing.T) {
			ch := make(chan int)
			go func() {
				ch <- 1
				ch <- 2
				ch <- 3
				close(ch)
			}()
			s := FromChannel(ch)
			got := CollectSlice(s)
			want := []int{1, 2, 3}
			assert.ElementsMatch(t, got, want)
		})

		t.Run("buffered", func(t *testing.T) {
			ch := make(chan int, 3)
			go func() {
				ch <- 1
				ch <- 2
				ch <- 3
				close(ch)
			}()
			s := FromChannel(ch)
			got := CollectSlice(s)
			want := []int{1, 2, 3}
			assert.ElementsMatch(t, got, want)
		})
	})

	t.Run("limited", func(t *testing.T) {
		ch := make(chan int, 3)
		go func() {
			ch <- 1
			ch <- 2
			ch <- 3
			close(ch)
		}()
		s := FromChannel(ch)
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []int{1, 2}
		assert.ElementsMatch(t, got, want)
	})
}

func TestFromChannelCtx(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ch := make(chan int)
		go close(ch)
		s := FromChannelCtx(context.Background(), ch)
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		ch := make(chan int, 3)
		go cancel()
		s := FromChannelCtx(ctx, ch)
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Run("unbuffered", func(t *testing.T) {
			ch := make(chan int)
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				ch <- 1
				ch <- 2
				ch <- 3
				close(ch)
			}()
			s := FromChannelCtx(context.Background(), ch)
			got := CollectSlice(s)
			want := []int{1, 2, 3}
			assert.ElementsMatch(t, got, want)
		})
		t.Run("buffered", func(t *testing.T) {
			ch := make(chan int, 3)
			go func() {
				ch <- 1
				ch <- 2
				ch <- 3
				close(ch)
			}()
			s := FromChannelCtx(context.Background(), ch)
			got := CollectSlice(s)
			want := []int{1, 2, 3}
			assert.ElementsMatch(t, got, want)
		})
	})

	t.Run("limited", func(t *testing.T) {
		ch := make(chan int, 3)
		go func() {
			ch <- 1
			ch <- 2
			ch <- 3
			close(ch)
		}()
		s := FromChannelCtx(context.Background(), ch)
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []int{1, 2}
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
