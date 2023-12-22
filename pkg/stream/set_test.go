package stream

import (
	"github.com/jpfourny/papaya/pkg/cmp"
	"testing"

	"github.com/jpfourny/papaya/internal/assert"
)

func TestUnion(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := Union[int]()
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := Union(Of(1, 2, 3), Of(4, 5, 6))
		got := CollectSlice(s)
		want := []int{1, 2, 3, 4, 5, 6}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := Union(Of(1, 2, 3), Of(4, 5, 6))
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []int{1, 2}
		assert.ElementsMatch(t, got, want)
	})
}

func TestIntersection(t *testing.T) {
	t.Run("two-streams", func(t *testing.T) {
		s := Intersection(Of(1, 2, 3), Of(2, 3, 4))
		got := CollectSlice(s)
		want := []int{2, 3}
		assert.ElementsMatchAnyOrder(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := Intersection(Of(1, 2, 3), Of(2, 3, 4))
		got := CollectSlice(Limit(s, 1)) // Stops stream after 1 element.
		want := []int{2}
		assert.ElementsMatchAnyOrder(t, got, want)
	})
}

func TestIntersectionBy(t *testing.T) {
	t.Run("two-streams", func(t *testing.T) {
		s := IntersectionBy(Of(1, 2, 3), Of(2, 3, 4), cmp.Natural[int]())
		got := CollectSlice(s)
		want := []int{2, 3}
		assert.ElementsMatchAnyOrder(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := IntersectionBy(Of(1, 2, 3), Of(2, 3, 4), cmp.Natural[int]())
		got := CollectSlice(Limit(s, 1)) // Stops stream after 1 element.
		want := []int{2}
		assert.ElementsMatchAnyOrder(t, got, want)
	})
}

func TestIntersectionAll(t *testing.T) {
	t.Run("zero-streams", func(t *testing.T) {
		s := IntersectionAll[int]()
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("one-stream", func(t *testing.T) {
		s := IntersectionAll(Of(1, 2, 3))
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		assert.ElementsMatchAnyOrder(t, got, want)
	})

	t.Run("two-streams", func(t *testing.T) {
		s := IntersectionAll(Of(1, 2, 3), Of(2, 3, 4))
		got := CollectSlice(s)
		want := []int{2, 3}
		assert.ElementsMatchAnyOrder(t, got, want)
	})

	t.Run("three-streams", func(t *testing.T) {
		s := IntersectionAll(Of(1, 2, 3), Of(2, 3, 4), Of(3, 4, 5))
		got := CollectSlice(s)
		want := []int{3}
		assert.ElementsMatchAnyOrder(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := IntersectionAll(Of(1, 2, 3), Of(2, 3, 4))
		got := CollectSlice(Limit(s, 1)) // Stops stream after 1 element.
		want := []int{2}
		assert.ElementsMatchAnyOrder(t, got, want)
	})
}

func TestIntersectionAllBy(t *testing.T) {
	t.Run("zero-streams", func(t *testing.T) {
		s := IntersectionAllBy[int](cmp.Natural[int]())
		got := CollectSlice(s)
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("one-stream", func(t *testing.T) {
		s := IntersectionAllBy(cmp.Natural[int](), Of(1, 2, 3))
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("two-streams", func(t *testing.T) {
		s := IntersectionAllBy(cmp.Natural[int](), Of(1, 2, 3), Of(2, 3, 4))
		got := CollectSlice(s)
		want := []int{2, 3}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("three-streams", func(t *testing.T) {
		s := IntersectionAllBy(cmp.Natural[int](), Of(1, 2, 3), Of(2, 3, 4), Of(3, 4, 5))
		got := CollectSlice(s)
		want := []int{3}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := IntersectionAllBy(cmp.Natural[int](), Of(1, 2, 3), Of(2, 3, 4))
		got := CollectSlice(Limit(s, 1)) // Stops stream after 1 element.
		want := []int{2}
		assert.ElementsMatch(t, got, want)
	})
}

func TestDifference(t *testing.T) {
	t.Run("two-streams", func(t *testing.T) {
		s := Difference(Of(1, 2, 3), Of(2, 3, 4))
		got := CollectSlice(s)
		want := []int{1}
		assert.ElementsMatch(t, got, want)
	})
}

func TestDifferenceBy(t *testing.T) {
	t.Run("two-streams", func(t *testing.T) {
		s := DifferenceBy(Of(1, 2, 3), Of(2, 3, 4), cmp.Natural[int]())
		got := CollectSlice(s)
		want := []int{1}
		assert.ElementsMatch(t, got, want)
	})
}
