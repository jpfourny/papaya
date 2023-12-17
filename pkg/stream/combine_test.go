package stream

import (
	"fmt"
	"testing"

	"github.com/jpfourny/papaya/internal/assert"
	"github.com/jpfourny/papaya/pkg/optional"
	"github.com/jpfourny/papaya/pkg/pair"
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
		s := CombineOrDiscard(Empty[int](), Empty[string](), func(i int, s string) optional.Optional[string] {
			return optional.Of(fmt.Sprintf("%s%d", s, i))
		})
		got := CollectSlice(s)
		var want []string
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := CombineOrDiscard(Of(1, 2, 3), Of("foo", "bar"), func(i int, s string) optional.Optional[string] {
			if i == 2 {
				return optional.Empty[string]()
			}
			return optional.Of(fmt.Sprintf("%s%d", s, i))
		})
		got := CollectSlice(s)
		want := []string{"foo1"}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := CombineOrDiscard(Of(1, 2, 3), Of("foo", "bar"), func(i int, s string) optional.Optional[string] {
			if i == 2 {
				return optional.Empty[string]()
			}
			return optional.Of(fmt.Sprintf("%s%d", s, i))
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

func TestZipWithIndexInt(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := ZipWithIndexInt(Empty[int](), 0)
		got := CollectSlice(s)
		var want []pair.Pair[int, int]
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := ZipWithIndexInt(Of(1, 2, 3), 0)
		got := CollectSlice(s)
		want := []pair.Pair[int, int]{
			pair.Of(1, 0),
			pair.Of(2, 1),
			pair.Of(3, 2),
		}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-zero-offset", func(t *testing.T) {
		s := ZipWithIndexInt(Of(1, 2, 3), -1)
		got := CollectSlice(s)
		want := []pair.Pair[int, int]{
			pair.Of(1, -1),
			pair.Of(2, 0),
			pair.Of(3, 1),
		}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := ZipWithIndexInt(Of(1, 2, 3), 0)
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []pair.Pair[int, int]{
			pair.Of(1, 0),
			pair.Of(2, 1),
		}
		assert.ElementsMatch(t, got, want)
	})
}

func TestZipWithIndexInt64(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := ZipWithIndexInt64(Empty[int](), 0)
		got := CollectSlice(s)
		var want []pair.Pair[int, int64]
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := ZipWithIndexInt64(Of(1, 2, 3), 0)
		got := CollectSlice(s)
		want := []pair.Pair[int, int64]{
			pair.Of(1, int64(0)),
			pair.Of(2, int64(1)),
			pair.Of(3, int64(2)),
		}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-zero-offset", func(t *testing.T) {
		s := ZipWithIndexInt64(Of(1, 2, 3), -1)
		got := CollectSlice(s)
		want := []pair.Pair[int, int64]{
			pair.Of(1, int64(-1)),
			pair.Of(2, int64(0)),
			pair.Of(3, int64(1)),
		}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := ZipWithIndexInt64(Of(1, 2, 3), 0)
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []pair.Pair[int, int64]{
			pair.Of(1, int64(0)),
			pair.Of(2, int64(1)),
		}
		assert.ElementsMatch(t, got, want)
	})
}

func TestZipWithKey(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := ZipWithKey(Empty[int](), func(e int) int {
			return e * 10
		})
		got := CollectSlice(s)
		var want []pair.Pair[int, int]
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := ZipWithKey(Of(1, 2, 3), func(e int) int {
			return e * 10
		})
		got := CollectSlice(s)
		want := []pair.Pair[int, int]{
			pair.Of(10, 1),
			pair.Of(20, 2),
			pair.Of(30, 3),
		}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := ZipWithKey(Of(1, 2, 3), func(e int) int {
			return e * 10
		})
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []pair.Pair[int, int]{
			pair.Of(10, 1),
			pair.Of(20, 2),
		}
		assert.ElementsMatch(t, got, want)
	})
}
