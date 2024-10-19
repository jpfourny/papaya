package stream

import (
	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/constraint"
	"slices"
)

// SortAsc returns a stream that sorts the elements in ascending order.
// The elements must implement the Ordered interface.
//
// Example usage:
//
//	s := stream.SortAsc(stream.Of(3, 1, 2))
//	out := stream.DebugString(s) // "<1, 2, 3>"
func SortAsc[E constraint.Ordered](s Stream[E]) Stream[E] {
	return SortBy(s, cmp.Natural[E]())
}

// SortDesc returns a stream that sorts the elements in descending order.
// The elements must implement the Ordered interface.
//
// Example usage:
//
//	s := stream.SortDesc(stream.Of(3, 1, 2))
//	out := stream.DebugString(s) // "<3, 2, 1>"
func SortDesc[E constraint.Ordered](s Stream[E]) Stream[E] {
	return SortBy(s, cmp.Reverse[E]())
}

// SortBy returns a stream that sorts the elements using the given cmp.Comparer.
//
// Example usage:
//
//	s := stream.SortBy(stream.Of(3, 1, 2), cmp.Natural[int]())
//	out := stream.DebugString(s) // "<1, 2, 3>"
func SortBy[E any](s Stream[E], compare cmp.Comparer[E]) Stream[E] {
	return func(yield Consumer[E]) {
		sl := CollectSlice(s)
		slices.SortFunc(sl, compare)
		FromSlice(sl)(yield)
	}
}
