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

// FromMapValues takes a map and returns a stream of its groups.
// The order of the groups is not guaranteed.
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
