package stream

import (
	"github.com/jpfourny/papaya/internal/kvstore"
	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/pair"
)

// Predicate is a function that accepts a value of type E and returns a boolean.
// It is used to test values for a given property.
// It must be idempotent, free of side effects, and thread-safe.
type Predicate[E any] func(e E) (pass bool)

// Filter returns a stream that only contains elements that pass the given Predicate.
//
// Example usage:
//
//	s := stream.Filter(stream.Of(1, 2, 3), func(e int) bool {
//	  return e % 2 == 0
//	})
//	out := stream.DebugString(s) // "<2>"
func Filter[E any](s Stream[E], p Predicate[E]) Stream[E] {
	return func(yield Consumer[E]) {
		s(func(e E) bool {
			if p(e) {
				return yield(e)
			}
			return true
		})
	}
}

// FilterIndexed returns a stream that only contains elements and their index that pass the given Predicate.
//
// Example usage:
//
//	s := stream.FilterIndexed(stream.Of(1, 2, 3, 4), func(e int) bool {
//	  return e % 2 == 0
//	})
//	out := stream.DebugString(s) // "<(1, 2), (3, 4)>"
func FilterIndexed[E any](s Stream[E], p func(E) bool) Stream[pair.Pair[int64, E]] {
	return func(yield Consumer[pair.Pair[int64, E]]) {
		var i int64
		s(func(e E) bool {
			i++
			if p(e) {
				return yield(pair.Of(i-1, e))
			}
			return true
		})
	}
}

// Limit returns a stream that is limited to the first `n` elements.
// If the input stream has fewer than `n` elements, the returned stream will have the same number of elements.
//
// Example usage:
//
//	s := stream.Limit(stream.Of(1, 2, 3), 2)
//	out := stream.DebugString(s) // "<1, 2>"
func Limit[E any](s Stream[E], n int64) Stream[E] {
	return func(yield Consumer[E]) {
		n := n // Shadow with a copy.
		if n <= 0 {
			return // Nothing to do.
		}
		s(func(e E) bool {
			n--
			return yield(e) && n > 0
		})
	}
}

// Skip returns a stream that skips the first `n` elements.
// If the input stream has fewer than `n` elements, the returned stream will be empty.
//
// Example usage:
//
//	s := stream.Skip(stream.Of(1, 2, 3), 2)
//	out := stream.DebugString(s) // "<3>"
func Skip[E any](s Stream[E], n int64) Stream[E] {
	return func(yield Consumer[E]) {
		n := n // Shadow with a copy.
		s(func(e E) bool {
			if n > 0 {
				n--
				return true
			}
			return yield(e)
		})
	}
}

// Slice returns a stream that contains elements from the start index (inclusive) to the end index (exclusive).
// If the start index is greater than the end index, the returned stream will be empty.
// If the end index is greater than the number of elements in the input stream, the returned stream will contain all elements from the start index.
//
// Example usage:
//
//	s := stream.Slice(stream.Of(1, 2, 3), 1, 2)
//	out := stream.DebugString(s) // "<2>"
func Slice[E any](s Stream[E], start, end int64) Stream[E] {
	return Limit(Skip(s, start), end-start)
}

// Distinct returns a stream that only contains distinct elements of some comparable type E.
//
// Example usage:
//
//	s := stream.Distinct(stream.Of(1, 2, 2, 3))
//	out := stream.DebugString(s) // "<1, 2, 3>"
func Distinct[E comparable](s Stream[E]) Stream[E] {
	return distinct(s, kvstore.MappedMaker[E, struct{}]())
}

// DistinctBy returns a stream that only contains distinct elements using the given comparer to compare elements.
//
// Example usage:
//
//	s := stream.DistinctBy(stream.Of(1, 2, 2, 3), cmp.Natural[int]())
//	out := stream.DebugString(s) // "<1, 2, 3>"
func DistinctBy[E any](s Stream[E], compare cmp.Comparer[E]) Stream[E] {
	return distinct(s, kvstore.SortedMaker[E, struct{}](compare))
}

func distinct[E any](s Stream[E], kv kvstore.Maker[E, struct{}]) Stream[E] {
	return func(yield Consumer[E]) {
		seen := kv()
		s(func(e E) bool {
			if seen.Get(e).Present() {
				return true // Skip.
			}
			seen.Put(e, struct{}{})
			return yield(e)
		})
	}
}
