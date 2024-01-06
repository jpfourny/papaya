package stream

import (
	"slices"

	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/constraint"
	"github.com/jpfourny/papaya/pkg/optional"
)

// Mapper represents a function that transforms an input of type E to an output of type F.
// It is used in the Map operation.
// It must be idempotent, free of side effects, and thread-safe.
type Mapper[E, F any] func(from E) (to F)

// OptionalMapper represents a function that transforms an input of type E to an optional output of type F.
// If the input cannot be transformed, the function must return an empty optional.
// It is used in the MapOrDiscard operation.
// It must be idempotent, free of side effects, and thread-safe.
type OptionalMapper[E, F any] func(from E) optional.Optional[F]

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

// MapOrDiscard applies an OptionalMapper function to each element in a stream and returns a new stream containing the mapped elements.
// If the OptionalMapper returns an empty optional, the element is discarded from the stream.
//
// Example usage:
//
//	s := stream.MapOrDiscard(stream.Of("1", "foo", "3"), mapper.TryParseInt[int](10, 64))
//	out := stream.DebugString(s) // "<1, 3>"
func MapOrDiscard[E, F any](s Stream[E], m OptionalMapper[E, F]) Stream[F] {
	return func(yield Consumer[F]) bool {
		return s(func(e E) (ok bool) {
			ok = true
			m(e).IfPresent(func(f F) {
				ok = yield(f)
			})
			return
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
//	s := stream.FlatMap(
//	  stream.Of(1, 2, 3),
//	  func(e int) stream.Stream[string] { // e -> <"e", "e">
//	    return stream.Of(mapper.Sprint(e), mapper.Sprint(e))
//	  },
//	)
//	out := stream.DebugString(s) // "<1, 1, 2, 2, 3, 3>"
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
func SortBy[E any](s Stream[E], compare cmp.Comparer[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		sl := CollectSlice(s)
		slices.SortFunc(sl, compare)
		return FromSlice(sl)(yield)
	}
}

// Truncate returns a stream that limits the given stream to the desired length and appends the given 'tail' value, if the stream is longer than the desired length.
// The tail value is appended only once, even if the stream is longer than the desired.
// If the stream is already shorter than the desired length, then the stream is returned as-is.
//
// Example usage:
//
//	s := stream.Truncate(stream.Of("a", "b", "c""), 2, "...")
//	out := stream.DebugString(s) // "<a, b, ...>"
//
//	s = stream.Truncate(stream.Of("a", "b", "c""), 3, "...")
//	out = stream.DebugString(s) // "<a, b, c>"
func Truncate[E any](s Stream[E], length int, tail E) Stream[E] {
	return func(yield Consumer[E]) bool {
		i := 0
		return s(func(e E) bool {
			i++
			if i <= length {
				return yield(e)
			}
			yield(tail)
			return false // Stop after the tail.
		})
	}
}

// PadTail returns a stream that pads the tail of the given stream with the given 'pad' value until the stream reaches the given length.
// If the stream is already longer than the given length, then the stream is returned as-is.
//
// Example usage:
//
//	s := stream.Pad(stream.Of(1, 2, 3), 0, 5)
//	out := stream.DebugString(s) // "<1, 2, 3, 0, 0>"
func PadTail[E any](s Stream[E], pad E, length int) Stream[E] {
	return func(yield Consumer[E]) bool {
		i := 0
		if !s(func(e E) bool {
			i++
			return yield(e)
		}) {
			return false // Consumer saw enough.
		}
		for ; i < length; i++ {
			if !yield(pad) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}
