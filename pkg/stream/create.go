package stream

import (
	"context"

	"github.com/jpfourny/papaya/pkg/optional"
	"github.com/jpfourny/papaya/pkg/pair"
)

// Stream represents a function that produces a sequence of elements of type E and sends them to the given Consumer.
//
// Streams are lazy, meaning they only produce elements when the consumer is invoked.
// Furthermore, streams are idempotent, meaning they can be invoked multiple times with the same result.
// However, the order of the elements is not guaranteed to be the same across multiple invocations.
//
// If the Consumer returns false, the stream must stop producing elements and return false immediately.
// If the stream is exhausted, it must return true.
type Stream[E any] func(c Consumer[E]) bool

// Consumer represents a function that accepts a yielded element of type E and returns a boolean value.
// The boolean value indicates whether the consumer wishes to continue accepting elements.
// If the consumer returns false, the caller must stop yielding elements.
type Consumer[E any] func(yield E) (cont bool)

// Empty returns a stream that does not contain any elements.
// It always returns true when invoked with a consumer.
//
// Example usage:
//
//	s := stream.Empty[int]()
//	out := stream.DebugString(s) // "<>"
func Empty[E any]() Stream[E] {
	return func(_ Consumer[E]) bool {
		return true
	}
}

// Of creates a stream from the given elements.
//
// Example usage:
//
//	s := stream.Of(1, 2, 3)
//	out := stream.DebugString(s) // "<1, 2, 3>"
func Of[E any](e ...E) Stream[E] {
	return FromSlice(e)
}

// FromSlice creates a stream that iterates over the elements of the given slice.
// The order of the elements is guaranteed to be the same as the order in the slice.
//
// Example usage:
//
//	s := stream.FromSlice([]int{1, 2, 3})
//	out := stream.DebugString(s) // "<1, 2, 3>"
func FromSlice[E any](s []E) Stream[E] {
	return func(yield Consumer[E]) bool {
		for _, e := range s {
			if !yield(e) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}

// FromSliceWithIndex returns a stream that iterates over the elements of the input slice along with their indices.
// The stream returns a `pair.Pair` with both the index and the value of each element.
// The order of the elements is guaranteed to be the same as the order in the slice.
//
// Example usage:
//
//	s := stream.FromSliceWithIndex([]int{1, 2, 3})
//	out := stream.DebugString(s) // "<(0, 1), (1, 2), (2, 3)>"
func FromSliceWithIndex[E any](s []E) Stream[pair.Pair[int, E]] {
	return func(yield Consumer[pair.Pair[int, E]]) bool {
		for i, e := range s {
			if !yield(pair.Of(i, e)) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}

// FromMap returns a stream that iterates over the key-value pairs in the given map.
// The key-value pairs are encapsulated in `pair.Pair` objects.
// The order of the key-value pairs is not guaranteed.
//
// Example usage:
//
//	s := stream.FromMap(map[int]string{1: "foo", 2: "bar"})
//	out := stream.DebugString(s) // "<(1, foo), (2, bar)>"
func FromMap[K comparable, V any](m map[K]V) Stream[pair.Pair[K, V]] {
	return func(yield Consumer[pair.Pair[K, V]]) bool {
		for k, v := range m {
			if !yield(pair.Of(k, v)) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}

// FromMapKeys takes a map and returns a stream of its keys.
// The order of the keys is not guaranteed.
//
// Example usage:
//
//	s := stream.FromMapKeys(map[int]string{1: "foo", 2: "bar"})
//	out := stream.DebugString(s) // "<1, 2>" // Order not guaranteed.
func FromMapKeys[K comparable, V any](m map[K]V) Stream[K] {
	return func(yield Consumer[K]) bool {
		for k := range m {
			if !yield(k) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}

// FromMapValues takes a map and returns a stream of its values.
// The order of the values is not guaranteed.
//
// Example usage:
//
//	s := stream.FromMapValues(map[int]string{1: "foo", 2: "bar"})
//	out := stream.DebugString(s) // "<foo, bar>" // Order not guaranteed.
func FromMapValues[K comparable, V any](m map[K]V) Stream[V] {
	return func(yield Consumer[V]) bool {
		for _, v := range m {
			if !yield(v) {
				return false
			}
		}
		return true
	}
}

// FromChannel returns a stream that reads elements from the given channel until it is closed.
//
//	Note: If the channel is not closed, the stream will block forever.
//	Note: If the consumer returns false, channel writes may block forever if the channel is unbuffered.
//
// Example usage:
//
//	ch := make(chan int)
//	go func() {
//	    ch <- 1
//	    ch <- 2
//	    close(ch)
//	}()
//	s := stream.FromChannel(ch)
//	out := stream.DebugString(s) // "<1, 2>"
func FromChannel[E any](ch <-chan E) Stream[E] {
	return func(yield Consumer[E]) bool {
		for e := range ch {
			if !yield(e) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}

// FromChannelCtx behaves like FromChannel, but it returns immediately when the context is done.
//
//	Note: If the context is done, channel writes may block forever if the channel is unbuffered.
//
// Example usage:
//
//	ch := make(chan int)
//	go func() {
//	    ch <- 1
//	    ch <- 2
//	    close(ch)
//	}()
//	s := stream.FromChannelCtx(ctx, ch)
//	out := stream.DebugString(s) // "<1, 2>"
func FromChannelCtx[E any](ctx context.Context, ch <-chan E) Stream[E] {
	return func(yield Consumer[E]) bool {
		for {
			select {
			case e, ok := <-ch:
				if !ok {
					return true
				}
				if !yield(e) {
					return false
				}
			case <-ctx.Done():
				return true
			}
		}
	}
}

// Iterator represents a type that provides a way to iterate over a sequence of elements.
// The Next() method is used to obtain the next element in the iteration, if any.
// If there are no more elements, the method returns an empty optional.Optional..
type Iterator[E any] interface {
	Next() optional.Optional[E]
}

// FromIterator returns a stream that produces elements by calling the provided Iterator function.
func FromIterator[E any](it Iterator[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		for e := it.Next(); e.Present(); e = it.Next() {
			if !yield(e.Get()) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}

// Range returns a stream that produces elements over a range beginning at `start`, advanced by the `next` function, and ending when `cond` predicate returns false.
//
// Example usage:
//
//	s := stream.Range(1, pred.LessThanOrEqual(5), mapper.Increment(2))
//	out := stream.DebugString(s) // "<1, 3, 5>"
func Range[E any](start E, cond Predicate[E], next Mapper[E, E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		for e := start; cond(e); e = next(e) {
			if !yield(e) {
				return false
			}
		}
		return true
	}
}
