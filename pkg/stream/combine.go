package stream

import (
	"context"
	"github.com/jpfourny/papaya/v2/pkg/constraint"
	"github.com/jpfourny/papaya/v2/pkg/opt"
	"github.com/jpfourny/papaya/v2/pkg/pair"
	"github.com/jpfourny/papaya/v2/pkg/stream/mapper"
	"github.com/jpfourny/papaya/v2/pkg/stream/pred"
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
	return CombineOrDiscard(s1, s2, func(e1 E1, e2 E2) opt.Optional[F] {
		return opt.Of(combine(e1, e2))
	})
}

// OptionalCombiner represents a function that combines two elements of type E1 and E2 into an opt element of type F.
// If the elements cannot be combined, the function must return an empty opt.
// It is used in the CombineOrDiscard operation.
type OptionalCombiner[E1, E2, F any] func(E1, E2) opt.Optional[F]

// CombineOrDiscard combines the elements of two streams into a single stream using the given OptionalCombiner function or discards them, if the combiner returns an empty opt.
// The resulting stream will have at most the same number of elements as the shorter of the two input streams.
//
// Example usage:
//
//	s := stream.CombineOrDiscard(
//	  stream.Of(1, 2, 3),
//	  stream.Of("foo", "bar"),
//	  func(i int, s string) opt.Optional[string] {
//	    if i == 2 {
//	      return opt.Empty[string]()
//	    }
//	    return opt.Of(fmt.Sprintf("%s%d", s, i))
//	  },
//	)
//	out := stream.DebugString(s) // "<foo1>"
func CombineOrDiscard[E1, E2, F any](s1 Stream[E1], s2 Stream[E2], combine OptionalCombiner[E1, E2, F]) Stream[F] {
	return func(yield Consumer[F]) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ch1 := CollectChannelAsyncCtx(ctx, s1, 0)
		ch2 := CollectChannelAsyncCtx(ctx, s2, 0)

		for {
			e1, ok1 := <-ch1
			e2, ok2 := <-ch2
			if !ok1 || !ok2 {
				return
			}

			if o := combine(e1, e2); o.Present() {
				if !yield(o.GetOrZero()) {
					return
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
	return Combine(s1, s2, pair.Of[E, F])
}

// ZipWithIndex returns a stream that pairs each element in the input stream with its index, starting at the given offset.
// The index type I must be an integer type.
//
// Example usage:
//
//	s := stream.ZipWithIndex(stream.Of("foo", "bar"), 1)
//	out := stream.DebugString(s) // "<(foo, 1), (bar, 2)>"
func ZipWithIndex[E any, I constraint.Integer](s Stream[E], offset I) Stream[pair.Pair[E, I]] {
	return Zip(s, Walk(offset, pred.True[I](), mapper.Increment[I](1)))
}

// UnzipFirst returns a stream that contains the first elements of each pair in the input stream.
//
// Example usage:
//
//	s := stream.UnzipFirst(
//	  stream.Of(
//	    pair.Of(1, "foo"),
//	    pair.Of(2, "bar"),
//	  ),
//	)
//	out := stream.DebugString(s) // "<1, 2>"
func UnzipFirst[E, F any](s Stream[pair.Pair[E, F]]) Stream[E] {
	return Map(s, pair.Pair[E, F].First)
}

// UnzipSecond returns a stream that contains the second elements of each pair in the input stream.
//
// Example usage:
//
//	s := stream.UnzipSecond(
//	  stream.Of(
//	    pair.Of(1, "foo"),
//	    pair.Of(2, "bar"),
//	  ),
//	)
//	out := stream.DebugString(s) // "<foo, bar>"
func UnzipSecond[E, F any](s Stream[pair.Pair[E, F]]) Stream[F] {
	return Map(s, pair.Pair[E, F].Second)
}
