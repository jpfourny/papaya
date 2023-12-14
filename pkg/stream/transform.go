package stream

import (
	"slices"

	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/constraint"
)

// Peek decorates the given stream to invoke the given function for each element passing through it.
// This is useful for debugging or logging elements as they pass through the stream.
//
// Example usage:
//
//	s := stream.Peek(stream.Of(1, 2, 3), func(e int) {
//	    fmt.Println(e)
//	})
//	stream.Count(s) // Force the stream to materialize.
//
// Output:
//
//	1
//	2
//	3
func Peek[E any](s Stream[E], peek func(e E)) Stream[E] {
	return func(yield Consumer[E]) bool {
		return s(func(e E) bool {
			peek(e)
			return yield(e)
		})
	}
}

// Mapper represents a function that transforms an input of type E to an output of type F.
// It is used in the Map operation.
// It must be idempotent, free of side effects, and thread-safe.
type Mapper[E, F any] func(from E) (to F)

// Map applies a Mapper function to each element in a stream and returns a new stream containing the mapped elements.
//
// Example usage:
//
//	s := stream.Map(stream.Of(1, 2, 3), mapper.Sprint)
//	out := stream.DebugString(s) // "<1, 2, 3>"
func Map[E, F any](s Stream[E], m Mapper[E, F]) Stream[F] {
	return func(yield Consumer[F]) bool {
		return s(func(e E) bool {
			return yield(m(e))
		})
	}
}

// FlatMapper represents a function that takes an input of type E and returns an output stream of type F.
// It is used to map each element of the input stream to a new stream of elements of type F.
// The FlatMapper function is typically used as a parameter in the FlatMap function.
// It must be idempotent, free of side effects, and thread-safe.
type FlatMapper[E, F any] func(from E) (to Stream[F])

// FlatMap applies a FlatMapper function to each element in a stream and returns a new stream containing the mapped elements.
// All mapped streams are flattened into a single stream.
//
// Example usage:
//
//	s := stream.FlatMap(stream.Of(1, 2, 3), func(e int) stream.Stream[string] {
//	    return stream.Map(stream.RangeInteger(0, e), mapper.Sprint)
//	})
//	out := stream.DebugString(s) // "<0, 0, 1, 0, 1, 2>"
func FlatMap[E, F any](s Stream[E], fm FlatMapper[E, F]) Stream[F] {
	return func(yield Consumer[F]) bool {
		return s(func(e E) bool {
			return fm(e)(yield)
		})
	}
}

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
func SortBy[E any](s Stream[E], cmp cmp.Comparer[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		sl := CollectSlice(s)
		slices.SortFunc(sl, cmp)
		return FromSlice(sl)(yield)
	}
}
