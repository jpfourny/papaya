package stream

import (
	"github.com/jpfourny/papaya/v2/pkg/opt"
)

// Mapper represents a function that transforms an input of type E to an output of type F.
// It is used in the Map operation.
// It must be idempotent, free of side effects, and thread-safe.
type Mapper[E, F any] func(from E) (to F)

// OptionalMapper represents a function that transforms an input of type E to an opt output of type F.
// If the input cannot be transformed, the function must return an empty opt.
// It is used in the MapOrDiscard operation.
// It must be idempotent, free of side effects, and thread-safe.
type OptionalMapper[E, F any] func(from E) opt.Optional[F]

// Map applies a Mapper function to each element in a stream and returns a new stream containing the mapped elements.
//
// Example usage:
//
//	s := stream.Map(stream.Of(1, 2, 3), mapper.Sprint)
//	out := stream.DebugString(s) // "<1, 2, 3>"
func Map[E, F any](s Stream[E], m Mapper[E, F]) Stream[F] {
	return func(yield Consumer[F]) {
		s(func(e E) bool {
			return yield(m(e))
		})
	}
}

// MapOrDiscard applies an OptionalMapper function to each element in a stream and returns a new stream containing the mapped elements.
// If the OptionalMapper returns an empty opt, the element is discarded from the stream.
//
// Example usage:
//
//	s := stream.MapOrDiscard(stream.Of("1", "foo", "3"), mapper.TryParseInt[int](10, 64))
//	out := stream.DebugString(s) // "<1, 3>"
func MapOrDiscard[E, F any](s Stream[E], m OptionalMapper[E, F]) Stream[F] {
	return func(yield Consumer[F]) {
		s(func(e E) (ok bool) {
			ok = true
			m(e).IfPresent(func(f F) {
				ok = yield(f)
			})
			return
		})
	}
}

// StreamMapper represents a function that takes an input of type E and returns an output stream of type F.
// The StreamMapper function is typically used as a parameter of the FlatMap function.
// It must be idempotent, free of side effects, and thread-safe.
type StreamMapper[E, F any] func(from E) (to Stream[F])

// FlatMap applies a StreamMapper function to each element from the given stream and flattens the returned streams into an output stream.
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
func FlatMap[E, F any](s Stream[E], m StreamMapper[E, F]) Stream[F] {
	return func(yield Consumer[F]) {
		yield2, stopped := stopSensingConsumer(yield)

		s(func(e E) bool {
			m(e)(yield2)
			if *stopped {
				return false // Consumer saw enough.
			}
			return true
		})
	}
}

// SliceMapper represents a function that takes an input of type E and returns an output slice of type F.
// The SliceMapper function is typically used as a parameter of the FlatMapSlice function.
// It must be idempotent, free of side effects, and thread-safe.
type SliceMapper[E, F any] func(from E) (to []F)

// FlatMapSlice applies a SliceMapper function to each element from the given stream and flattens the returned slices into an output stream.
//
// Example usage:
//
//	s := stream.FlatMapSlice(
//	  stream.Of(1, 2, 3),
//	  func(e int) []string { // e -> ["e", "e"]
//	    return []string{mapper.Sprint(e), mapper.Sprint(e)}
//	  },
//	)
//	out := stream.DebugString(s) // "<1, 1, 2, 2, 3, 3>"
func FlatMapSlice[E, F any](s Stream[E], m SliceMapper[E, F]) Stream[F] {
	return func(yield Consumer[F]) {
		s(func(e E) bool {
			yield2, stopped := stopSensingConsumer(yield)

			FromSlice(m(e))(yield2)
			if *stopped {
				return false // Consumer saw enough.
			}
			return true
		})
	}
}
