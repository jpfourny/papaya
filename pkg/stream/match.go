package stream

import (
	"context"
	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/opt"
	"github.com/jpfourny/papaya/pkg/stream/pred"
)

// AnyMatch returns true if any element in the stream matches the given Predicate.
// If the stream is empty, it returns false.
//
// Example usage:
//
//	out := stream.AnyMatch(stream.Of(1, 2, 3), pred.GreaterThan(2)) // true
//	out = stream.AnyMatch(stream.Of(1, 2, 3), pred.GreaterThan(3)) // false
func AnyMatch[E any](s Stream[E], p Predicate[E]) (anyMatch bool) {
	s(func(e E) bool {
		if p(e) {
			anyMatch = true
			return false // Stop the stream.
		}
		return true
	})
	return
}

// AllMatch returns true if all elements in the stream match the given Predicate.
// If the stream is empty, it returns false.
//
// Example usage:
//
//	out := stream.AllMatch(stream.Of(1, 2, 3), pred.GreaterThan(0)) // true
//	out = stream.AllMatch(stream.Of(1, 2, 3), pred.GreaterThan(1)) // false
func AllMatch[E any](s Stream[E], p Predicate[E]) (allMatch bool) {
	allMatch = true
	empty := true
	s(func(e E) bool {
		empty = false
		if !p(e) {
			allMatch = false
			return false // Stop the stream.
		}
		return true
	})
	allMatch = allMatch && !empty
	return
}

// NoneMatch returns true if no elements in the stream match the given Predicate.
// If the stream is empty, it returns true.
//
// Example usage:
//
//	out := stream.NoneMatch(stream.Of(1, 2, 3), pred.GreaterThan(3)) // true
//	out = stream.NoneMatch(stream.Of(1, 2, 3), pred.GreaterThan(2)) // false
func NoneMatch[E any](s Stream[E], p Predicate[E]) bool {
	return !AnyMatch(s, p)
}

// Contains returns true if the stream contains the given element; false otherwise.
// The element type E must be comparable.
//
// Example usage:
//
//	out := stream.Contains(stream.Of(1, 2, 3), 2) // true
//	out = stream.Contains(stream.Of(1, 2, 3), 4) // false
func Contains[E comparable](s Stream[E], e E) bool {
	return AnyMatch(s, pred.Equal(e))
}

// ContainsBy returns true if the stream contains the given element; false otherwise.
// The elements are compared using the given cmp.Comparer.
//
// Example usage:
//
//	out := stream.ContainsBy(stream.Of(1, 2, 3), 2, cmp.Natural[int]()) // true
//	out = stream.ContainsBy(stream.Of(1, 2, 3), 4, cmp.Natural[int]()) // false
func ContainsBy[E any](s Stream[E], compare cmp.Comparer[E], e E) bool {
	return AnyMatch(s, pred.EqualBy(e, compare))
}

// ContainsAny returns true if the stream contains any of the given elements; false otherwise.
// The element type E must be comparable.
//
// Example usage:
//
//	out := stream.ContainsAny(stream.Of(1, 2, 3), 2, 4) // true
//	out = stream.ContainsAny(stream.Of(1, 2, 3), 4, 5) // false
func ContainsAny[E comparable](s Stream[E], es ...E) bool {
	return AnyMatch(s, pred.In(es...))
}

// ContainsAnyBy returns true if the stream contains any of the given elements; false otherwise.
// The elements are compared using the given cmp.Comparer.
//
// Example usage:
//
//	out := stream.ContainsAnyBy(stream.Of(1, 2, 3), cmp.Natural[int](), 2, 4) // true
//	out = stream.ContainsAnyBy(stream.Of(1, 2, 3), cmp.Natural[int](), 4, 5) // false
func ContainsAnyBy[E any](s Stream[E], compare cmp.Comparer[E], es ...E) bool {
	return AnyMatch(s, pred.InBy(compare, es...))
}

// ContainsAll returns true if the stream contains all the given elements; false otherwise.
// The element type E must be comparable.
//
// Example usage:
//
//	out := stream.ContainsAll(stream.Of(1, 2, 3), 2, 3) // true
//	out = stream.ContainsAll(stream.Of(1, 2, 3), 2, 4) // false
func ContainsAll[E comparable](s Stream[E], es ...E) bool {
	return IsEmpty(
		Difference(
			Distinct(FromSlice(es)),
			Distinct(s),
		),
	)
}

// ContainsAllBy returns true if the stream contains all the given elements; false otherwise.
// The elements are compared using the given cmp.Comparer.
//
// Example usage:
//
//	out := stream.ContainsAllBy(stream.Of(1, 2, 3), cmp.Natural[int](), 2, 3) // true
//	out = stream.ContainsAllBy(stream.Of(1, 2, 3), cmp.Natural[int](), 2, 4) // false
func ContainsAllBy[E any](s Stream[E], compare cmp.Comparer[E], es ...E) bool {
	return IsEmpty(
		DifferenceBy(
			DistinctBy(FromSlice(es), compare),
			DistinctBy(s, compare),
			compare,
		),
	)
}

// ContainsNone returns true if the stream contains none of the given elements; false otherwise.
// The element type E must be comparable.
//
// Example usage:
//
//	out := stream.ContainsNone(stream.Of(1, 2, 3), 4, 5) // true
//	out = stream.ContainsNone(stream.Of(1, 2, 3), 2, 4) // false
func ContainsNone[E comparable](s Stream[E], es ...E) bool {
	return NoneMatch(s, pred.In(es...))
}

// ContainsNoneBy returns true if the stream contains none of the given elements; false otherwise.
// The elements are compared using the given cmp.Comparer.
//
// Example usage:
//
//	out := stream.ContainsNoneBy(stream.Of(1, 2, 3), cmp.Natural[int](), 4, 5) // true
//	out = stream.ContainsNoneBy(stream.Of(1, 2, 3), cmp.Natural[int](), 2, 4) // false
func ContainsNoneBy[E any](s Stream[E], compare cmp.Comparer[E], es ...E) bool {
	return NoneMatch(s, pred.InBy(compare, es...))
}

// ExactlySame returns true if the two streams are exactly the same; false otherwise.
// The element type E must be comparable.
//
// Example usage:
//
//	out := stream.ExactlySame(stream.Of(1, 2, 3), stream.Of(1, 2, 3)) // true
//	out = stream.ExactlySame(stream.Of(1, 2, 3), stream.Of(1, 2, 4)) // false
func ExactlySame[E comparable](s1, s2 Stream[E]) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch1 := CollectChannelAsyncCtx(ctx, s1, 0)
	ch2 := CollectChannelAsyncCtx(ctx, s2, 0)

	for {
		e1, ok1 := <-ch1
		e2, ok2 := <-ch2

		if !ok1 && !ok2 {
			break
		} else if !ok1 || !ok2 || e1 != e2 {
			return false
		}
	}

	return true
}

// ExactlySameBy returns true if the two streams are exactly the same; false otherwise.
// The elements are compared using the given cmp.Comparer.
//
// Example usage:
//
//	out := stream.ExactlySameBy(stream.Of(1, 2, 3), stream.Of(1, 2, 3), cmp.Natural[int]()) // true
//	out = stream.ExactlySameBy(stream.Of(1, 2, 3), stream.Of(1, 2, 4), cmp.Natural[int]()) // false
func ExactlySameBy[E any](s1, s2 Stream[E], compare cmp.Comparer[E]) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch1 := CollectChannelAsyncCtx(ctx, s1, 0)
	ch2 := CollectChannelAsyncCtx(ctx, s2, 0)

	for {
		e1, ok1 := <-ch1
		e2, ok2 := <-ch2
		if !ok1 && !ok2 {
			break
		} else if !ok1 || !ok2 || !compare.Equal(e1, e2) {
			return false
		}
	}

	return true
}

// IndexOfMatches returns a stream of indices of elements that pass the given Predicate.
//
// Example usage:
//
//	out := stream.IndexOfMatches(stream.Of(1, 2, 3), pred.GreaterThan(2)) // "<2>"
//	out = stream.IndexOfMatches(stream.Of(1, 2, 3), pred.GreaterThan(3)) // "<>"
func IndexOfMatches[E any](s Stream[E], p Predicate[E]) Stream[int64] {
	return UnzipFirst(FilterIndexed(s, p))
}

// IndexOfFirstMatch returns the index of the first element that passes the given Predicate.
// If no element passes the Predicate, an empty Optional is returned.
//
// Example usage:
//
//	out := stream.IndexOfFirstMatch(stream.Of(1, 2, 3, 4), pred.GreaterThan(2)) // opt.Of(2)
//	out = stream.IndexOfFirstMatch(stream.Of(1, 2, 3), pred.GreaterThan(3)) // opt.Empty[int64]()
func IndexOfFirstMatch[E any](s Stream[E], p Predicate[E]) opt.Optional[int64] {
	return First(IndexOfMatches(s, p))
}

// IndexOfLastMatch returns the index of the last element that passes the given Predicate.
// If no element passes the Predicate, an empty Optional is returned.
//
// Example usage:
//
//	out := stream.IndexOfLastMatch(stream.Of(1, 2, 3, 4), pred.GreaterThan(2)) // opt.Of(4)
//	out = stream.IndexOfLastMatch(stream.Of(1, 2, 3), pred.GreaterThan(3)) // opt.Empty[int64]()
func IndexOfLastMatch[E any](s Stream[E], p Predicate[E]) opt.Optional[int64] {
	return Last(IndexOfMatches(s, p))
}
