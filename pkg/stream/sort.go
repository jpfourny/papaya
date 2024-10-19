package stream

import (
	"github.com/jpfourny/papaya/v2/pkg/cmp"
	"github.com/jpfourny/papaya/v2/pkg/constraint"
	"github.com/jpfourny/papaya/v2/pkg/pair"
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

// SortKeyAsc returns a stream that sorts the elements by key in ascending order.
// The elements must be pairs of key and value.
// The keys must implement the Ordered interface.
//
// Example usage:
//
//	s := stream.SortKeyAsc(
//	  stream.Of(
//	    pair.Of(3, "c"),
//	    pair.Of(1, "a"),
//	    pair.Of(2, "b"),
//	  ),
//	)
//	out := stream.DebugString(s) // "<1:a, 2:b, 3:c>"
func SortKeyAsc[K constraint.Ordered, V any](s Stream[pair.Pair[K, V]]) Stream[pair.Pair[K, V]] {
	return SortKeyBy(s, cmp.Natural[K]())
}

// SortKeyDesc returns a stream that sorts the elements by key in descending order.
// The elements must be pairs of key and value.
// The keys must implement the Ordered interface.
//
// Example usage:
//
//	s := stream.SortKeyDesc(
//	  stream.Of(
//	    pair.Of(3, "c"),
//	    pair.Of(1, "a"),
//	    pair.Of(2, "b"),
//	  ),
//	)
//	out := stream.DebugString(s) // "<3:c, 2:b, 1:a>"
func SortKeyDesc[K constraint.Ordered, V any](s Stream[pair.Pair[K, V]]) Stream[pair.Pair[K, V]] {
	return SortKeyBy(s, cmp.Reverse[K]())
}

// SortKeyBy returns a stream that sorts the elements by key using the given cmp.Comparer.
// The elements must be pairs of key and value.
// The order of the elements is determined by the comparer.
//
// Example usage:
//
//	s := stream.SortKeyBy(
//	  stream.Of(
//	    pair.Of(3, "c"),
//	    pair.Of(1, "a"),
//	    pair.Of(2, "b"),
//	  ),
//	  cmp.Natural[int](),
//	)
//	out := stream.DebugString(s) // "<1:a, 2:b, 3:c>"
func SortKeyBy[K any, V any](s Stream[pair.Pair[K, V]], compare cmp.Comparer[K]) Stream[pair.Pair[K, V]] {
	return SortBy(s, cmp.ComparingBy(pair.Pair[K, V].First, compare))
}

// SortValueAsc returns a stream that sorts the elements by value in ascending order.
// The elements must be pairs of key and value.
// The values must implement the Ordered interface.
//
// Example usage:
//
//	s := stream.SortValueAsc(
//	  stream.Of(
//	    pair.Of("c", 3),
//	    pair.Of("a", 1),
//	    pair.Of("b", 2),
//	  ),
//	)
//	out := stream.DebugString(s) // "a:1, b:2, c:3"
func SortValueAsc[K any, V constraint.Ordered](s Stream[pair.Pair[K, V]]) Stream[pair.Pair[K, V]] {
	return SortValueBy(s, cmp.Natural[V]())
}

// SortValueDesc returns a stream that sorts the elements by value in descending order.
// The elements must be pairs of key and value.
// The values must implement the Ordered interface.
//
// Example usage:
//
//	s := stream.SortValueDesc(
//	  stream.Of(
//	    pair.Of("c", 3),
//	    pair.Of("a", 1),
//	    pair.Of("b", 2),
//	  ),
//	)
//	out := stream.DebugString(s) // "c:3, b:2, a:1"
func SortValueDesc[K any, V constraint.Ordered](s Stream[pair.Pair[K, V]]) Stream[pair.Pair[K, V]] {
	return SortValueBy(s, cmp.Reverse[V]())
}

// SortValueBy returns a stream that sorts the elements by value using the given cmp.Comparer.
// The elements must be pairs of key and value.
// The order of the elements is determined by the comparer.
//
// Example usage:
//
//	s := stream.SortValueBy(
//	  stream.Of(
//	    pair.Of("c", 3),
//	    pair.Of("a", 1),
//	    pair.Of("b", 2),
//	  ),
//	  cmp.Natural[int](),
//	)
//	out := stream.DebugString(s) // "a:1, b:2, c:3"
func SortValueBy[K any, V any](s Stream[pair.Pair[K, V]], compare cmp.Comparer[V]) Stream[pair.Pair[K, V]] {
	return SortBy(s, cmp.ComparingBy(pair.Pair[K, V].Second, compare))
}
