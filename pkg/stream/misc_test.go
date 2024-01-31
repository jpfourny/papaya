package stream

import (
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
