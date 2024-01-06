package stream

import (
	"github.com/jpfourny/papaya/pkg/cmp"
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
//	    return e % 2 == 0
//	})
//	out := stream.DebugString(s) // "<2>"
func Filter[E any](s Stream[E], p Predicate[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		return s(func(e E) bool {
			if p(e) {
				return yield(e)
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
	return func(yield Consumer[E]) bool {
		n := n // Shadow with a copy.
		if n <= 0 {
			return true
		}
		return s(func(e E) bool {
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
	return func(yield Consumer[E]) bool {
		n := n // Shadow with a copy.
		return s(func(e E) bool {
			if n > 0 {
				n--
				return true
			}
			return yield(e)
		})
	}
}

// Distinct returns a stream that only contains distinct elements of some comparable type E.
//
// Example usage:
//
//	s := stream.Distinct(stream.Of(1, 2, 2, 3))
//	out := stream.DebugString(s) // "<1, 2, 3>"
func Distinct[E comparable](s Stream[E]) Stream[E] {
	return distinct(s, mapKeyStoreFactory[E, struct{}]())
}

// DistinctBy returns a stream that only contains distinct elements using the given comparer to compare elements.
//
// Example usage:
//
//	s := stream.DistinctBy(stream.Of(1, 2, 2, 3), cmp.Natural[int]())
//	out := stream.DebugString(s) // "<1, 2, 3>"
func DistinctBy[E any](s Stream[E], compare cmp.Comparer[E]) Stream[E] {
	return distinct(s, sortedKeyStoreFactory[E, struct{}](compare))
}

func distinct[E any](s Stream[E], ksf keyStoreFactory[E, struct{}]) Stream[E] {
	return func(yield Consumer[E]) bool {
		seen := ksf()
		return s(func(e E) bool {
			if seen.get(e).Present() {
				return true // Skip.
			}
			seen.put(e, struct{}{})
			return yield(e)
		})
	}
}
