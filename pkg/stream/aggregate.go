package stream

import (
	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/constraint"
	"github.com/jpfourny/papaya/pkg/opt"
	"github.com/jpfourny/papaya/pkg/stream/mapper"
	"github.com/jpfourny/papaya/pkg/stream/reducer"
)

// Reducer represents a function that takes two inputs of type E and returns an output of type E.
// The Reducer is commonly used in the `Reduce` function to combine elements of a stream into a single result.
type Reducer[E any] func(e1, e2 E) (result E)

// Reduce combines the elements of the stream into a single value using the given reducer function.
// If the stream is empty, then an empty opt.Optional is returned.
// The stream is fully consumed.
//
// Example usage:
//
//	out := stream.Reduce(
//	  stream.Of(1, 2, 3),
//	  func(a, e int) int { // Reduce values by addition.
//	    return a + e
//	  },
//	) // Some(6)
//
//	out = stream.Reduce(
//	  stream.Empty[int](),
//	  func(a, e int) int { // Reduce values by addition.
//	    return a + e
//	  },
//	) // None()
func Reduce[E any](s Stream[E], reduce Reducer[E]) opt.Optional[E] {
	var accum E
	var ok bool
	s(func(e E) bool {
		if ok {
			accum = reduce(accum, e)
		} else {
			accum = e
			ok = true
		}
		return true
	})
	return opt.Maybe(accum, ok)
}

// Sum computes the sum of all elements in the stream of any real-number type E and returns the result as real-number type F.
// The result of an empty stream is the zero value of type F.
// The stream is fully consumed.
//
// Example usage:
//
//	n1 := stream.Sum[int](stream.Of(1, 2, 3)) // 6 (int)
//	n2 := stream.Sum[float64](stream.Of(1, 2, 3)) // 6.0 (float64)
func Sum[R, E constraint.RealNumber](s Stream[E]) R {
	return Reduce(
		Map(s, mapper.NumToNum[E, R]()),
		reducer.Sum[R](),
	).GetOrZero()
}

// SumComplex computes the sum of all elements in the stream of any complex-number type E and returns the result as complex-number type F.
// The result of an empty stream is the zero value of type F.
// The stream is fully consumed.
//
// Example usage:
//
//	n := stream.SumComplex[complex128](stream.Of(1+i, 2+i, 3+i)) // 6+3i (complex128)
func SumComplex[R, E constraint.Complex](s Stream[E]) R {
	return Reduce(
		Map(s, mapper.ComplexToComplex[E, R]()),
		reducer.Sum[R](),
	).GetOrZero()
}

// Min returns the minimum element in the stream, or the zero value of the type parameter E if the stream is empty.
// If the stream is empty, the 'ok' return value is false; otherwise it is true.
// Uses the natural ordering of type E to compare elements.
//
// Example usage:
//
//	min := stream.Min(stream.Of(3, 1, 2)) // Some(1)
//	min = stream.Min(stream.Empty[int]()) // None()
func Min[E constraint.Ordered](s Stream[E]) (min opt.Optional[E]) {
	return Reduce(s, reducer.Min[E]())
}

// MinBy returns the minimum element in the stream.
// Uses the given cmp.Comparer to compare elements.
// If the stream is empty, then an empty opt.Optional is returned.
//
// Example usage:
//
//	min := stream.MinBy(stream.Of(3, 1, 2), cmp.Natural[int]()) // Some(1)
//	min = stream.MinBy(stream.Empty[int](), cmp.Natural[int]()) // None()
func MinBy[E any](s Stream[E], compare cmp.Comparer[E]) (min opt.Optional[E]) {
	return Reduce(s, reducer.MinBy(compare))
}

// Max returns the maximum element in the stream.
// Uses the natural ordering of type E to compare elements.
// If the stream is empty, then an empty opt.Optional is returned.
//
// Example usage:
//
//	max := stream.Max(stream.Of(3, 1, 2)) // Some(3)
//	max = stream.Max(stream.Empty[int]()) // None()
func Max[E constraint.Ordered](s Stream[E]) (max opt.Optional[E]) {
	return Reduce(s, reducer.Max[E]())
}

// MaxBy returns the maximum element in the stream, or the zero value of the type parameter E if the stream is empty.
// Uses the given cmp.Comparer to compare elements.
// If the stream is empty, then an empty opt.Optional is returned.
//
// Example usage:
//
//	max := stream.MaxBy(stream.Of(3, 1, 2), cmp.Natural[int]()) // Some(3)
//	max = stream.MaxBy(stream.Empty[int](), cmp.Natural[int]()) // None()
func MaxBy[E any](s Stream[E], compare cmp.Comparer[E]) (max opt.Optional[E]) {
	return Reduce(s, reducer.MaxBy(compare))
}

// Accumulator represents a function that takes an accumulated value of type A and an element of type E,
// and returns the updated accumulated value of type A.
// The Accumulator is commonly used in the `Aggregate` function to combine elements of a stream into a single result.
type Accumulator[A, E any] func(a A, e E) (result A)

// Finisher represents a function that takes an accumulated value of type A and returns the finished result of type F.
// The Finisher is commonly used in the `Aggregate` function to compute the final result after all elements have been accumulated.
type Finisher[A, F any] func(a A) (result F)

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
//	  0,                 // Initial value
//	  func(a, e int) int {
//	    return a + e     // Accumulate with addition
//	  },
//	  func(a int) int {
//	    return a * 2     // Finish with multiplication by 2
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
		float64(0), // Initialize with zero.
		func(a float64, e E) float64 { // Accumulate with addition; increment count.
			count++
			return a + float64(e)
		},
		func(a float64) float64 { // Finish with division by count, or zero if count is zero.
			if count == 0 {
				return 0
			}
			return a / float64(count)
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
	s(func(e E) bool {
		count++
		return true
	})
	return
}
