package stream

import (
	"github.com/jpfourny/papaya/pkg/stream/mapper"
	"github.com/jpfourny/papaya/pkg/stream/pred"
	"testing"

	"github.com/jpfourny/papaya/internal/assert"
	"github.com/jpfourny/papaya/pkg/opt"
)

func TestForEach(t *testing.T) {
	var got []int
	ForEach(Of(1, 2, 3), func(e int) {
		got = append(got, e)
	})
	want := []int{1, 2, 3}
	assert.ElementsMatch(t, got, want)
}

func TestPeek(t *testing.T) {
	var got []int
	s := Peek(Of(1, 2, 3), func(e int) {
		got = append(got, e)
	})
	Count(s) // Force evaluation so peek is called.
	want := []int{1, 2, 3}
	assert.ElementsMatch(t, got, want)
	assert.ElementsMatch(t, CollectSlice(s), want)
}

func TestDebugString(t *testing.T) {
	got := DebugString(Of(1, 2, 3))
	want := "<1, 2, 3>"
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}

	got = DebugString(Generate(func() int { return 1 })) // Infinite stream will be truncated to 100 elements (+ tailing ...).
	want = "<1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, ...>"
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestIsEmpty(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := IsEmpty(Empty[int]())
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := IsEmpty(Of(1, 2, 3))
		if got != false {
			t.Errorf("got %#v, want %#v", got, false)
		}
	})
}

func TestFirst(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := First(Empty[int]())
		want := opt.Empty[int]()
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := First(Of(1, 2, 3))
		want := opt.Of(1)
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}

func TestLast(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := Last(Empty[int]())
		want := opt.Empty[int]()
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := Last(Of(1, 2, 3))
		want := opt.Of(3)
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}

func TestCache(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := Cache(Empty[int]())
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := Cache(Of(1, 2, 3))
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := Cache(Of(1, 2, 3))
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []int{1, 2}
		assert.ElementsMatch(t, got, want)

		got = CollectSlice(s) // Replay from cache without the limit.
		want = []int{1, 2, 3} // Despite limit from first call to stream, cache has all elements.
		assert.ElementsMatch(t, got, want)

		got = CollectSlice(Limit(s, 1)) // Stops stream after 1 elements.
		want = []int{1}                 // Limit is applied to the cache.
		assert.ElementsMatch(t, got, want)
	})
}

func TestTruncate(t *testing.T) {
	s := Truncate(Of("a", "b", "c"), 4, "...")
	got := CollectSlice(s)
	want := []string{"a", "b", "c"}
	assert.ElementsMatch(t, got, want)

	s = Truncate(Of("a", "b", "c"), 3, "...")
	got = CollectSlice(s)
	want = []string{"a", "b", "c"}
	assert.ElementsMatch(t, got, want)

	s = Truncate(Of("a", "b", "c"), 2, "...")
	got = CollectSlice(s)
	want = []string{"a", "b", "..."}
	assert.ElementsMatch(t, got, want)

	s = Truncate(Of("a", "b", "c"), 1, "...")
	got = CollectSlice(s)
	want = []string{"a", "..."}
	assert.ElementsMatch(t, got, want)

	s = Truncate(Of("a", "b", "c"), 0, "...")
	got = CollectSlice(s)
	want = []string{"..."}
	assert.ElementsMatch(t, got, want)

	s = Truncate(Empty[string](), 1, "...")
	got = CollectSlice(s)
	want = []string{}
	assert.ElementsMatch(t, got, want)
}

func TestPad(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := Pad(Of[int](), 0, 5)
		got := CollectSlice(s)
		want := []int{0, 0, 0, 0, 0}
		assert.ElementsMatch(t, got, want)
	})
	t.Run("short", func(t *testing.T) {
		s := Pad(Of(1, 2, 3), 0, 5)
		got := CollectSlice(s)
		want := []int{1, 2, 3, 0, 0}
		assert.ElementsMatch(t, got, want)
	})
	t.Run("long", func(t *testing.T) {
		s := Pad(Of(1, 2, 3, 4, 5, 6), 0, 5)
		got := CollectSlice(s)
		want := []int{1, 2, 3, 4, 5, 6}
		assert.ElementsMatch(t, got, want)
	})
	t.Run("limited-before-padding", func(t *testing.T) {
		s := Pad(Of(1, 2, 3), 0, 5)
		got := CollectSlice(Limit(s, 2))
		want := []int{1, 2}
		assert.ElementsMatch(t, got, want)
	})
	t.Run("limited-during-padding", func(t *testing.T) {
		s := Pad(Of(1, 2, 3), 0, 5)
		got := CollectSlice(Limit(s, 4))
		want := []int{1, 2, 3, 0}
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
