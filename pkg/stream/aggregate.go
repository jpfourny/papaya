package stream

import (
	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/constraint"
	"github.com/jpfourny/papaya/pkg/optional"
)

// Reducer represents a function that takes two inputs of type E and returns an output of type E.
// The Reducer is commonly used in the `Reduce` function to combine elements of a stream into a single result.
type Reducer[E any] func(e1, e2 E) (result E)

// Accumulator represents a function that takes an accumulated value of type A and an element of type E,
// and returns the updated accumulated value of type A.
// The Accumulator is commonly used in the `Aggregate` function to combine elements of a stream into a single result.
type Accumulator[A, E any] func(a A, e E) (result A)

// Finisher represents a function that takes an accumulated value of type A and returns the finished result of type F.
// The Finisher is commonly used in the `Aggregate` function to compute the final result after all elements have been accumulated.
type Finisher[A, F any] func(a A) (result F)

// Reduce combines the elements of the stream into a single value using the given reducer function.
// If the stream is empty, then an empty optional.Optional is returned.
// The stream is fully consumed.
//
// Example usage:
//
//	out := stream.Reduce(stream.Of(1, 2, 3), func(a, e int) int {
//	    return a + e
//	}) // Some(6)
//	out = stream.Reduce(stream.Empty[int](), func(a, e int) int {
//	    return a + e
//	} // None()
func Reduce[E any](s Stream[E], reduce Reducer[E]) (result optional.Optional[E]) {
	result = optional.Empty[E]()
	s(func(e E) bool {
		if result.Present() {
			result = optional.Of(reduce(result.Get(), e))
		} else {
			result = optional.Of(e)
		}
		return true
	})
	return
}

// Aggregate combines the elements of the stream into a single value using the given identity value, accumulator function and finisher function.
// The accumulated value is initialized to the identity value.
// The accumulator function is used to combine each element with the accumulated value.
// The finisher function is used to compute the final result after all elements have been accumulated.
// The stream is fully consumed.
//
// Example usage:
//
//	s := stream.Aggregate(
//	  stream.Of(1, 2, 3),
//	  0,                   // Initial value.
//	  func(a, e int) int {
//	      return a + e     // Accumulate with addition.
//	  },
//	  func(a int) int {
//	      return a * 2     // Finish with multiplication by 2.
//	  },
//	) // (1+2+3) * 2 = 12
func Aggregate[E, A, F any](s Stream[E], identity A, accumulate Accumulator[A, E], finish Finisher[A, F]) F {
	a := identity
	s(func(e E) bool {
		a = accumulate(a, e)
		return true
	})
	return finish(a)
}

// Sum computes the sum of all elements in the stream of any number type E and returns the result as type E.
// The result of an empty stream is the zero value of type E.
// The stream is fully consumed.
//
// Example usage:
//
//	n := stream.Sum(stream.Of(1, 2, 3)) // 6 (int)
func Sum[E constraint.RealNumber](s Stream[E]) E {
	return Aggregate(
		s,
		E(0),
		func(a E, e E) E { return a + e },
		func(a E) E { return a },
	)
}

// SumInteger computes the sum of all elements in the stream of any signed-integer type E and returns the result as type int64.
// The result of an empty stream is the zero value of type int64.
// The stream is fully consumed.
//
// Example usage:
//
//	n := stream.SumInteger(stream.Of(1, 2, 3)) // 6 (int64)
func SumInteger[E constraint.SignedInteger](s Stream[E]) int64 {
	return Aggregate(
		s,
		int64(0),
		func(a int64, e E) int64 { return a + int64(e) },
		func(a int64) int64 { return a },
	)
}

// SumUnsignedInteger computes the sum of all elements in the stream of any unsigned-integer type E and returns the result as type uint64.
// The result of an empty stream is the zero value of type uint64.
// The stream is fully consumed.
//
// Example usage:
//
//	n := stream.SumUnsignedInteger(stream.Of[uint](1, 2, 3)) // 6 (uint64)
func SumUnsignedInteger[E constraint.UnsignedInteger](s Stream[E]) uint64 {
	return Aggregate(
		s,
		uint64(0),
		func(a uint64, e E) uint64 { return a + uint64(e) },
		func(a uint64) uint64 { return a },
	)
}

// SumFloat computes the sum of all elements in the stream of any floating-point type E and returns the result as type float64.
// The result of an empty stream is the zero value of type float64.
// The stream is fully consumed.
//
// Example usage:
//
//	n := stream.SumFloat(stream.Of(1.0, 2.0, 3.0)) // 6.0 (float64)
func SumFloat[E constraint.RealNumber](s Stream[E]) float64 {
	return Aggregate(
		s,
		float64(0),
		func(a float64, e E) float64 { return a + float64(e) },
		func(a float64) float64 { return a },
	)
}

// Average computes the average of all elements in the stream of any number type E and returns the result as type float64.
// The result of an empty stream is the zero value of type float64.
// The stream is fully consumed.
//
// Example usage:
//
//	n := stream.Average(stream.Of(1, 2, 3)) // 2.0 (float64)
func Average[E constraint.RealNumber](s Stream[E]) float64 {
	var count uint64
	return Aggregate(
		s,
		float64(0),
		func(a float64, e E) float64 {
			count++
			return a + float64(e)
		},
		func(a float64) float64 {
			if count == 0 {
				return 0
			}
			return a / float64(count)
		},
	)
}

// Min returns the minimum element in the stream, or the zero value of the type parameter E if the stream is empty.
// If the stream is empty, the 'ok' return value is false; otherwise it is true.
// Uses the natural ordering of type E to compare elements.
//
// Example usage:
//
//	min := stream.Min(stream.Of(3, 1, 2)) // Some(1)
//	min = stream.Min(stream.Empty[int]()) // None()
func Min[E constraint.Ordered](s Stream[E]) (min optional.Optional[E]) {
	return MinBy(s, cmp.Natural[E]())
}

// MinBy returns the minimum element in the stream.
// Uses the given cmp.Comparer to compare elements.
// If the stream is empty, then an empty optional.Optional is returned.
//
// Example usage:
//
//	min := stream.MinBy(stream.Of(3, 1, 2), cmp.Natural[int]()) // Some(1)
//	min = stream.MinBy(stream.Empty[int](), cmp.Natural[int]()) // None()
func MinBy[E any](s Stream[E], cmp cmp.Comparer[E]) (min optional.Optional[E]) {
	return Reduce(
		s,
		func(a, e E) E {
			if cmp(e, a) < 0 {
				return e
			}
			return a
		},
	)
}

// Max returns the maximum element in the stream.
// Uses the natural ordering of type E to compare elements.
// If the stream is empty, then an empty optional.Optional is returned.
//
// Example usage:
//
//	max := stream.Max(stream.Of(3, 1, 2)) // Some(3)
//	max = stream.Max(stream.Empty[int]()) // None()
func Max[E constraint.Ordered](s Stream[E]) (max optional.Optional[E]) {
	return MaxBy(s, cmp.Natural[E]())
}

// MaxBy returns the maximum element in the stream, or the zero value of the type parameter E if the stream is empty.
// Uses the given cmp.Comparer to compare elements.
// If the stream is empty, then an empty optional.Optional is returned.
//
// Example usage:
//
//	max := stream.MaxBy(stream.Of(3, 1, 2), cmp.Natural[int]()) // Some(3)
//	max = stream.MaxBy(stream.Empty[int](), cmp.Natural[int]()) // None()
func MaxBy[E any](s Stream[E], cmp cmp.Comparer[E]) (max optional.Optional[E]) {
	return Reduce(
		s,
		func(a, e E) E {
			if cmp(e, a) > 0 {
				return e
			}
			return a
		},
	)
}

// Count returns the number of elements in the stream.
// The stream is fully consumed.
//
// Example usage:
//
//	n := stream.Count(stream.Of(1, 2, 3)) // 3
func Count[E any](s Stream[E]) (count int64) {
	// Count is so simple that we don't bother to use a Aggregate.
	s(func(e E) bool {
		count++
		return true
	})
	return
}
