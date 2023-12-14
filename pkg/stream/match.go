package stream

import "github.com/jpfourny/papaya/pkg/pred"

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
//
// Example usage:
//
//	out := stream.Contains(stream.Of(1, 2, 3), 2) // true
//	out = stream.Contains(stream.Of(1, 2, 3), 4) // false
func Contains[E comparable](s Stream[E], e E) bool {
	return AnyMatch(s, pred.Equal(e))
}

// ContainsAny returns true if the stream contains any of the given elements; false otherwise.
//
// Example usage:
//
//	out := stream.ContainsAny(stream.Of(1, 2, 3), 2, 4) // true
//	out = stream.ContainsAny(stream.Of(1, 2, 3), 4, 5) // false
func ContainsAny[E comparable](s Stream[E], es ...E) bool {
	return AnyMatch(s, pred.In(es...))
}

// ContainsAll returns true if the stream contains all the given elements; false otherwise.
//
// Example usage:
//
//	out := stream.ContainsAll(stream.Of(1, 2, 3), 2, 3) // true
//	out = stream.ContainsAll(stream.Of(1, 2, 3), 2, 4) // false
func ContainsAll[E comparable](s Stream[E], es ...E) bool {
	return Count(Distinct(Filter(s, pred.In(es...)))) == int64(len(es))
}

// ContainsNone returns true if the stream contains none of the given elements; false otherwise.
//
// Example usage:
//
//	out := stream.ContainsNone(stream.Of(1, 2, 3), 4, 5) // true
//	out = stream.ContainsNone(stream.Of(1, 2, 3), 2, 4) // false
func ContainsNone[E comparable](s Stream[E], es ...E) bool {
	return NoneMatch(s, pred.In(es...))
}
