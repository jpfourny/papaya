package stream

import (
	"github.com/jpfourny/papaya/v2/internal/kvstore"
	"github.com/jpfourny/papaya/v2/pkg/cmp"
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
	return func(yield Consumer[E]) {
		yield2, stopped := stopSensingConsumer(yield)

		for _, s := range ss {
			s(yield2)
			if *stopped {
				return // Consumer saw enough.
			}
		}
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
		return Intersection[E](ss[0], IntersectionAll[E](ss[1:]...))
	}

	// Exactly 2 streams; intersect ss[0] and ss[1].
	return Intersection[E](ss[0], ss[1])
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
		return IntersectionBy[E](ss[0], IntersectionAllBy[E](compare, ss[1:]...), compare)
	}

	// Exactly 2 streams; intersect ss[0] and ss[1].
	return IntersectionBy[E](ss[0], ss[1], compare)
}

func intersection[E any](s1, s2 Stream[E], kv kvstore.Maker[E, struct{}]) Stream[E] {
	return func(yield Consumer[E]) {
		set2 := toSet(s2, kv)
		// Yield elements of the first stream that are in the set.
		s1(func(e E) bool {
			if set2.Get(e).Present() {
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
	return func(yield Consumer[E]) {
		set2 := toSet(s2, kv)
		// Yield elements of the first stream that are not in the set.
		s1(func(e E) bool {
			if !set2.Get(e).Present() {
				return yield(e)
			}
			return true
		})
	}
}

// SymmetricDifference returns a stream that contains elements that are in either of the given streams, but not in both.
// The element type E must be comparable.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.SymmetricDifference(stream.Of(1, 2, 3, 4, 5), stream.Of(4, 5, 6))
//	out := stream.DebugString(s) // "<1, 2, 3, 6>"
func SymmetricDifference[E comparable](s1, s2 Stream[E]) Stream[E] {
	return Union(Difference(s1, s2), Difference(s2, s1))
}

// SymmetricDifferenceBy returns a stream that contains elements that are in either of the given streams, but not in both, compared by the given cmp.Comparer.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.SymmetricDifferenceBy(stream.Of(1, 2, 3, 4, 5), stream.Of(4, 5, 6), cmp.Natural[int]())
//	out := stream.DebugString(s) // "<1, 2, 3, 6>"
func SymmetricDifferenceBy[E any](s1, s2 Stream[E], compare cmp.Comparer[E]) Stream[E] {
	return Union(DifferenceBy(s1, s2, compare), DifferenceBy(s2, s1, compare))
}

// Subset returns true if all elements of the first stream are in the second stream.
// The element type E must be comparable.
//
// Example usage:
//
//	ok := stream.Subset(stream.Of(1, 2), stream.Of(1, 2, 3, 4))
//	fmt.Println(ok) // "true"
func Subset[E comparable](s1, s2 Stream[E]) bool {
	return subset(s1, s2, kvstore.MappedMaker[E, struct{}]())
}

// SubsetBy returns true if all elements of the first stream are in the second stream, compared by the given cmp.Comparer.
//
// Example usage:
//
//	ok := stream.SubsetBy(stream.Of(1, 2), stream.Of(1, 2, 3, 4), cmp.Natural[int]())
//	fmt.Println(ok) // "true"
func SubsetBy[E any](s1, s2 Stream[E], compare cmp.Comparer[E]) bool {
	return subset(s1, s2, kvstore.SortedMaker[E, struct{}](compare))
}

func subset[E any](s1, s2 Stream[E], kv kvstore.Maker[E, struct{}]) bool {
	// Index elements of the second stream into a set.
	set2 := toSet(s2, kv)
	// Check if all elements of the first stream are in the set.
	// If an element is not in the set, the result is false
	allInSet := true
	s1(func(e E) bool {
		if !set2.Get(e).Present() {
			allInSet = false
			return false
		}
		return true
	})
	return allInSet
}

// Superset returns true if all elements of the second stream are in the first stream.
// The element type E must be comparable.
//
// Example usage:
//
//	ok := stream.Superset(stream.Of(1, 2, 3, 4), stream.Of(1, 2))
//	fmt.Println(ok) // "true"
func Superset[E comparable](s1, s2 Stream[E]) bool {
	return Subset(s2, s1)
}

// SupersetBy returns true if all elements of the second stream are in the first stream, compared by the given cmp.Comparer.
//
// Example usage:
//
//	ok := stream.SupersetBy(stream.Of(1, 2, 3, 4), stream.Of(1, 2), cmp.Natural[int]())
//	fmt.Println(ok) // "true"
func SupersetBy[E any](s1, s2 Stream[E], compare cmp.Comparer[E]) bool {
	return SubsetBy(s2, s1, compare)
}

// SetEqual returns true if the two streams contain the same elements, ignoring order and duplicates (ie: set equality).
// The element type E must be comparable.
//
// Example usage:
//
//	ok := stream.SetEqual(stream.Of(1, 2, 3), stream.Of(3, 2, 1))
//	fmt.Println(ok) // "true"
func SetEqual[E comparable](s1, s2 Stream[E]) bool {
	return setEqual(s1, s2, kvstore.MappedMaker[E, struct{}]())
}

// SetEqualBy returns true if the two streams contain the same elements (in any order), compared by the given cmp.Comparer.
//
// Example usage:
//
//	ok := stream.SetEqualBy(stream.Of(1, 2, 3), stream.Of(3, 2, 1), cmp.Natural[int]())
//	fmt.Println(ok) // "true"
func SetEqualBy[E any](s1, s2 Stream[E], compare cmp.Comparer[E]) bool {
	return setEqual(s1, s2, kvstore.SortedMaker[E, struct{}](compare))
}

func setEqual[E any](s1, s2 Stream[E], kv kvstore.Maker[E, struct{}]) bool {
	set1 := toSet(s1, kv)
	set2 := toSet(s2, kv)
	if set1.Size() != set2.Size() {
		return false
	}
	allFound := true
	set1.ForEachKey(func(e E) bool {
		if !set2.Get(e).Present() {
			allFound = false
			return false
		}
		return true
	})
	return allFound
}

func toSet[E any](s Stream[E], kv kvstore.Maker[E, struct{}]) kvstore.Store[E, struct{}] {
	set := kv()
	s(func(e E) bool {
		set.Put(e, struct{}{})
		return true
	})
	return set
}
