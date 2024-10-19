package stream

import (
	"fmt"
	"testing"

	"github.com/jpfourny/papaya/v2/internal/assert"
	"github.com/jpfourny/papaya/v2/pkg/opt"
	"github.com/jpfourny/papaya/v2/pkg/pair"
)

func TestCombine(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := Combine(Empty[int](), Empty[string](), func(i int, s string) string {
			return fmt.Sprintf("%s%d", s, i)
		})
		got := CollectSlice(s)
		var want []string
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := Combine(Of(1, 2, 3), Of("foo", "bar"), func(i int, s string) string {
			return fmt.Sprintf("%s%d", s, i)
		})
		got := CollectSlice(s)
		want := []string{"foo1", "bar2"}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := Combine(Of(1, 2, 3), Of("foo", "bar"), func(i int, s string) string {
			return fmt.Sprintf("%s%d", s, i)
		})
		got := CollectSlice(Limit(s, 1)) // Stops stream after 1 element.
		want := []string{"foo1"}
		assert.ElementsMatch(t, got, want)
	})
}

func TestCombineOrDiscard(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := CombineOrDiscard(Empty[int](), Empty[string](), func(i int, s string) opt.Optional[string] {
			return opt.Of(fmt.Sprintf("%s%d", s, i))
		})
		got := CollectSlice(s)
		var want []string
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := CombineOrDiscard(Of(1, 2, 3), Of("foo", "bar"), func(i int, s string) opt.Optional[string] {
			if i == 2 {
				return opt.Empty[string]()
			}
			return opt.Of(fmt.Sprintf("%s%d", s, i))
		})
		got := CollectSlice(s)
		want := []string{"foo1"}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := CombineOrDiscard(Of(1, 2, 3), Of("foo", "bar"), func(i int, s string) opt.Optional[string] {
			if i == 2 {
				return opt.Empty[string]()
			}
			return opt.Of(fmt.Sprintf("%s%d", s, i))
		})
		got := CollectSlice(Limit(s, 1)) // Stops stream after 1 element.
		want := []string{"foo1"}
		assert.ElementsMatch(t, got, want)
	})
}

func TestZip(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := Zip(Empty[int](), Empty[int]())
		got := CollectSlice(s)
		var want []pair.Pair[int, int]
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := Zip(Of(1, 2, 3), Of(4, 5, 6))
		got := CollectSlice(s)
		want := []pair.Pair[int, int]{
			pair.Of(1, 4),
			pair.Of(2, 5),
			pair.Of(3, 6),
		}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := Zip(Of(1, 2, 3), Of(4, 5, 6))
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []pair.Pair[int, int]{
			pair.Of(1, 4),
			pair.Of(2, 5),
		}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("different-length", func(t *testing.T) {
		s := Zip(Of(1, 2, 3), Of(4, 5))
		got := CollectSlice(s)
		want := []pair.Pair[int, int]{
			pair.Of(1, 4),
			pair.Of(2, 5),
		}
		assert.ElementsMatch(t, got, want)
	})
}

func TestZipWithIndex(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := ZipWithIndex(Empty[int](), 0)
		got := CollectSlice(s)
		var want []pair.Pair[int, int]
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := ZipWithIndex(Of(1, 2, 3), 0)
		got := CollectSlice(s)
		want := []pair.Pair[int, int]{
			pair.Of(1, 0),
			pair.Of(2, 1),
			pair.Of(3, 2),
		}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-zero-offset", func(t *testing.T) {
		s := ZipWithIndex(Of(1, 2, 3), -1)
		got := CollectSlice(s)
		want := []pair.Pair[int, int]{
			pair.Of(1, -1),
			pair.Of(2, 0),
			pair.Of(3, 1),
		}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := ZipWithIndex(Of(1, 2, 3), 0)
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []pair.Pair[int, int]{
			pair.Of(1, 0),
			pair.Of(2, 1),
		}
		assert.ElementsMatch(t, got, want)
	})
}

func TestUnzipFirst(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := UnzipFirst(Empty[pair.Pair[int, string]]())
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := UnzipFirst(Of(
			pair.Of(1, "foo"),
			pair.Of(2, "bar"),
		))
		got := CollectSlice(s)
		want := []int{1, 2}
		assert.ElementsMatch(t, got, want)
	})
}

func TestUnzipSecond(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := UnzipSecond(Empty[pair.Pair[int, string]]())
		got := CollectSlice(s)
		var want []string
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := UnzipSecond(Of(
			pair.Of(1, "foo"),
			pair.Of(2, "bar"),
		))
		got := CollectSlice(s)
		want := []string{"foo", "bar"}
		assert.ElementsMatch(t, got, want)
	})
}
