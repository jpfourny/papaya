package stream

import (
	"context"

	"github.com/jpfourny/papaya/pkg/pair"
)

// FromSlice creates a stream that iterates over the elements of the given slice.
// The order of the elements is guaranteed to be the same as the order in the slice.
//
// Example usage:
//
//	s := stream.FromSlice([]int{1, 2, 3})
//	out := stream.DebugString(s) // "<1, 2, 3>"
func FromSlice[E any](s []E) Stream[E] {
	return func(yield Consumer[E]) {
		for _, e := range s {
			if !yield(e) {
				break // Consumer saw enough.
			}
		}
	}
}

// FromSliceBackwards creates a stream that iterates over the elements of the given slice in reverse order.
// The order of the elements is guaranteed to be the same as the order in the slice (but backwards).
//
// Example usage:
//
//	s := stream.FromSliceBackwards([]int{1, 2, 3})
//	out := stream.DebugString(s) // "<3, 2, 1>"
func FromSliceBackwards[E any](s []E) Stream[E] {
	return func(yield Consumer[E]) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(s[i]) {
				break // Consumer saw enough.
			}
		}
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
	return func(yield Consumer[pair.Pair[int, E]]) {
		for i, e := range s {
			if !yield(pair.Of(i, e)) {
				break // Consumer saw enough.
			}
		}
	}
}

// FromSliceWithIndexBackwards returns a stream that iterates over the elements of the input slice along with their indices in reverse order.
// The stream returns a `pair.Pair` with both the index and the value of each element.
// The order of the elements is guaranteed to be the same as the order in the slice (but backwards).
//
// Example usage:
//
//	s := stream.FromSliceWithIndexBackwards([]int{1, 2, 3})
//	out := stream.DebugString(s) // "<(2, 3), (1, 2), (0, 1)>"
func FromSliceWithIndexBackwards[E any](s []E) Stream[pair.Pair[int, E]] {
	return func(yield Consumer[pair.Pair[int, E]]) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(pair.Of(i, s[i])) {
				break // Consumer saw enough.
			}
		}
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
	return func(yield Consumer[pair.Pair[K, V]]) {
		for k, v := range m {
			if !yield(pair.Of(k, v)) {
				break // Consumer saw enough.
			}
		}
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
	return func(yield Consumer[K]) {
		for k := range m {
			if !yield(k) {
				break // Consumer saw enough.
			}
		}
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
	return func(yield Consumer[V]) {
		for _, v := range m {
			if !yield(v) {
				break // Consumer saw enough.
			}
		}
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
//	  ch <- 1
//	  ch <- 2
//	  close(ch)
//	}()
//	s := stream.FromChannel(ch)
//	out := stream.DebugString(s) // "<1, 2>"
func FromChannel[E any](ch <-chan E) Stream[E] {
	return func(yield Consumer[E]) {
		for e := range ch {
			if !yield(e) {
				break // Consumer saw enough.
			}
		}
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
//	  ch <- 1
//	  ch <- 2
//	  close(ch)
//	}()
//	s := stream.FromChannelCtx(ctx, ch)
//	out := stream.DebugString(s) // "<1, 2>"
func FromChannelCtx[E any](ctx context.Context, ch <-chan E) Stream[E] {
	return func(yield Consumer[E]) {
		for {
			select {
			case e, ok := <-ch:
				if !ok {
					return // Channel closed.
				}
				if !yield(e) {
					return // Consumer saw enough.
				}
			case <-ctx.Done():
				return // Context done.
			}
		}
	}
}
