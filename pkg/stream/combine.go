package stream

import (
	"github.com/jpfourny/papaya/pkg/constraint"
	"github.com/jpfourny/papaya/pkg/mapper"
	"github.com/jpfourny/papaya/pkg/optional"
	"github.com/jpfourny/papaya/pkg/pair"
	"github.com/jpfourny/papaya/pkg/pred"
)

// Combiner represents a function that combines two elements of type E1 and E2 into an element of type F.
// It is used in the Combine operation.
type Combiner[E1, E2, F any] func(E1, E2) F

// Combine combines the elements of two streams into a single stream using the given Combiner function.
// The resulting stream will have the same number of elements as the shorter of the two input streams.
//
// Example usage:
//
//	s := stream.Combine(
//	  stream.Of(1, 2, 3),
//	  stream.Of("foo", "bar"),
//	  func(i int, s string) string {
//	    return fmt.Sprintf("%s%d", s, i)
//	  },
//	)
//	out := stream.DebugString(s) // "<foo1, bar2>"
func Combine[E1, E2, F any](s1 Stream[E1], s2 Stream[E2], combine Combiner[E1, E2, F]) Stream[F] {
	return CombineOrDiscard(s1, s2, func(e1 E1, e2 E2) optional.Optional[F] {
		return optional.Of(combine(e1, e2))
	})
}

// OptionalCombiner represents a function that combines two elements of type E1 and E2 into an optional element of type F.
// If the elements cannot be combined, the function must return an empty optional.
// It is used in the CombineOrDiscard operation.
type OptionalCombiner[E1, E2, F any] func(E1, E2) optional.Optional[F]

// CombineOrDiscard combines the elements of two streams into a single stream using the given OptionalCombiner function or discards them, if the combiner returns an empty optional.
// The resulting stream will have at most the same number of elements as the shorter of the two input streams.
//
// Example usage:
//
//	s := stream.CombineOrDiscard(
//	  stream.Of(1, 2, 3),
//	  stream.Of("foo", "bar"),
//	  func(i int, s string) optional.Optional[string] {
//	    if i == 2 {
//	      return optional.Empty[string]()
//	    }
//	    return optional.Of(fmt.Sprintf("%s%d", s, i))
//	  },
//	)
//	out := stream.DebugString(s) // "<foo1>"
func CombineOrDiscard[E1, E2, F any](s1 Stream[E1], s2 Stream[E2], combine OptionalCombiner[E1, E2, F]) Stream[F] {
	return func(yield Consumer[F]) bool {
		done := make(chan struct{})
		defer close(done)

		ch1 := make(chan E1)
		go func() {
			defer close(ch1)
			s1(func(e E1) bool {
				select {
				case <-done:
					return false
				case ch1 <- e:
					return true
				}
			})
		}()

		ch2 := make(chan E2)
		go func() {
			defer close(ch2)
			s2(func(e E2) bool {
				select {
				case <-done:
					return false
				case ch2 <- e:
					return true
				}
			})
		}()

		for {
			e1, ok1 := <-ch1
			e2, ok2 := <-ch2
			if !ok1 || !ok2 {
				return true
			}

			if o := combine(e1, e2); o.Present() {
				if !yield(o.Get()) {
					return false
				}
			}
		}
	}
}

// Zip returns a stream that pairs each element in the first stream with the corresponding element in the second stream.
// The resulting stream will have the same number of elements as the shorter of the two input streams.
//
// Example usage:
//
//	s := stream.Zip(stream.Of(1, 2, 3), stream.Of("foo", "bar"))
//	out := stream.DebugString(s) // "<(1, foo), (2, bar)>"
func Zip[E, F any](s1 Stream[E], s2 Stream[F]) Stream[pair.Pair[E, F]] {
	return Combine(s1, s2, func(e E, f F) pair.Pair[E, F] {
		return pair.Of(e, f)
	})
}

// ZipWithIndex returns a stream that pairs each element in the input stream with its index, starting at the given offset.
// The index type I must be an integer type.
//
// Example usage:
//
//	s := stream.ZipWithIndex(stream.Of("foo", "bar"), 1)
//	out := stream.DebugString(s) // "<(foo, 1), (bar, 2)>"
func ZipWithIndex[E any, I constraint.Integer](s Stream[E], offset I) Stream[pair.Pair[E, I]] {
	return Combine(
		s,
		Range[I](offset, pred.True[I](), mapper.Increment[I](1)),
		func(e E, i I) pair.Pair[E, I] {
			return pair.Of(e, i)
		},
	)
}

// KeyExtractor represents a function that extracts a key of type K from a value of type E.
type KeyExtractor[E, K any] func(e E) K

// ZipWithKey returns a stream that pairs each element in the input stream with the key extracted from the element using the given key extractor.
// The resulting stream will have the same number of elements as the input stream.
//
// Example usage:
//
//	s := stream.ZipWithKey(
//	  stream.Of("foo", "bar"),
//	  func(s string) string {
//	    return strings.ToUpper(s)
//	  },
//	)
//	out := stream.DebugString(s) // "<(FOO, foo), (BAR, bar)>"
func ZipWithKey[K, E any](s Stream[E], ke KeyExtractor[E, K]) Stream[pair.Pair[K, E]] {
	return Map(s, func(e E) pair.Pair[K, E] {
		return pair.Of(ke(e), e)
	})
}
