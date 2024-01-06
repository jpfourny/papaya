package stream

import (
	"github.com/jpfourny/papaya/internal/kvstore"
	"github.com/jpfourny/papaya/pkg/cmp"
)

// Union combines multiple streams into a single stream (concatenation).
// The length of the resulting stream is the sum of the lengths of the input streams.
// If any of the input streams return false when invoked with the consumer, the concatenation stops.
//
// Example usage:
//
//	s := stream.Union(stream.Of(1, 2, 3, 4), stream.Of(4, 5, 6))
//	out := stream.DebugString(s) // "<1, 2, 3, 4, 4, 5, 6>"
func Union[E any](ss ...Stream[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		for _, s := range ss {
			if !s(yield) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}

// Intersection returns a stream that contains elements that are in the given streams.
// The element type E must be comparable.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.Intersection(stream.Of(1, 2, 3, 4, 5), stream.Of(4, 5, 6))
//	out := stream.DebugString(s) // "<4, 5>"
func Intersection[E comparable](s1, s2 Stream[E]) Stream[E] {
	return intersection(s1, s2, kvstore.MappedMaker[E, struct{}]())
}

// IntersectionAll returns a stream that contains elements that are in all the given streams.
// The element type E must be comparable.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.IntersectionAll(stream.Of(1, 2, 3, 4, 5), stream.Of(4, 5, 6), stream.Of(4, 5, 7))
//	out := stream.DebugString(s) // "<4, 5>"
func IntersectionAll[E comparable](ss ...Stream[E]) Stream[E] {
	switch {
	case len(ss) == 0: // No streams; result is empty.
		return Empty[E]()
	case len(ss) == 1: // One stream; result is the same stream.
		return ss[0]
	case len(ss) > 2: // More than 2 streams; recursively intersect stream pairs.
		return Intersection(ss[0], IntersectionAll(ss[1:]...))
	}

	// Exactly 2 streams; intersect ss[0] and ss[1].
	return Intersection(ss[0], ss[1])
}

// IntersectionBy returns a stream that contains elements that are in the given streams, compared by the given cmp.Comparer.
// The order of the elements is determined by the comparer.
//
// Example usage:
//
//	s := stream.IntersectionBy(stream.Of(1, 2, 3, 4, 5), stream.Of(4, 5, 6), cmp.Natural[int]())
//	out := stream.DebugString(s) // "<4, 5>"
func IntersectionBy[E any](s1, s2 Stream[E], compare cmp.Comparer[E]) Stream[E] {
	return intersection(s1, s2, kvstore.SortedMaker[E, struct{}](compare))
}

// IntersectionAllBy returns a stream that contains elements that are in all the given streams, compared by the given cmp.Comparer.
// The order of the elements is determined by the comparer.
//
// Example usage:
//
//	s := stream.IntersectionAllBy(cmp.Natural[int](), stream.Of(1, 2, 3, 4, 5), stream.Of(4, 5, 6), stream.Of(4, 5, 7))
//	out := stream.DebugString(s) // "<4, 5>"
func IntersectionAllBy[E any](compare cmp.Comparer[E], ss ...Stream[E]) Stream[E] {
	switch {
	case len(ss) == 0: // No streams; result is empty.
		return Empty[E]()
	case len(ss) == 1: // One stream; result is the same stream.
		return ss[0]
	case len(ss) > 2: // More than 2 streams; recursively intersect stream pairs.
		return IntersectionBy(ss[0], IntersectionAllBy(compare, ss[1:]...), compare)
	}

	// Exactly 2 streams; intersect ss[0] and ss[1].
	return IntersectionBy(ss[0], ss[1], compare)
}

func intersection[E any](s1, s2 Stream[E], kv kvstore.Maker[E, struct{}]) Stream[E] {
	// Exactly 2 streams; intersect ss[0] and ss[1].
	return func(yield Consumer[E]) bool {
		// Index elements of the first stream into a set.
		seen := kv()
		s1(func(e E) bool {
			seen.Put(e, struct{}{})
			return true
		})
		// Yield elements of the second stream that are in the set.
		return s2(func(e E) bool {
			if seen.Get(e).Present() {
				return yield(e)
			}
			return true
		})
	}
}

// Difference returns a stream that contains elements that are in the first stream but not in the second stream.
// The element type E must be comparable.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.Difference(stream.Of(1, 2, 3, 4, 5), stream.Of(4, 5, 6))
//	out := stream.DebugString(s) // "<1, 2, 3>"
func Difference[E comparable](s1, s2 Stream[E]) Stream[E] {
	return difference(s1, s2, kvstore.MappedMaker[E, struct{}]())
}

// DifferenceBy returns a stream that contains elements that are in the first stream but not in the second stream, compared by the given cmp.Comparer.
// The order of the elements is determined by the comparer.
//
// Example usage:
//
//	s := stream.DifferenceBy(stream.Of(1, 2, 3, 4, 5), stream.Of(4, 5, 6), cmp.Natural[int]())
//	out := stream.DebugString(s) // "<1, 2, 3>"
func DifferenceBy[E any](s1, s2 Stream[E], compare cmp.Comparer[E]) Stream[E] {
	return difference(s1, s2, kvstore.SortedMaker[E, struct{}](compare))
}

func difference[E any](s1, s2 Stream[E], kv kvstore.Maker[E, struct{}]) Stream[E] {
	return func(yield Consumer[E]) bool {
		// Index elements of the second stream into a set.
		seen := kv()
		s2(func(e E) bool {
			seen.Put(e, struct{}{})
			return true
		})
		// Yield elements of the first stream that are not in the set.
		return s1(func(e E) bool {
			if !seen.Get(e).Present() {
				return yield(e)
			}
			return true
		})
	}
}
