package stream

import (
	"context"
	"github.com/jpfourny/papaya/pkg/stream/mapper"

	"github.com/jpfourny/papaya/pkg/pair"
)

// CollectSlice returns a slice containing all elements from the stream.
// The stream is fully consumed.
//
// Example usage:
//
//	s := stream.CollectSlice(stream.Of(1, 2, 3)) // []int{1, 2, 3}
func CollectSlice[E any](s Stream[E]) []E {
	return Aggregate(
		s,
		[]E(nil), // Initialize with nil slice.
		func(a []E, e E) []E { return append(a, e) }, // Accumulate: Append element to slice.
		mapper.Identity[[]E](),                       // Finish: Return the slice as is.
	)
}

// CollectMap returns a map containing all key-value pair elements from the stream.
// The stream is fully consumed.
//
// Example usage:
//
//	s := stream.CollectMap(stream.Of(pair.Of(1, "foo"), pair.Of(2, "bar"))) // map[int]string{1: "foo", 2: "bar"}
func CollectMap[K comparable, V any](s Stream[pair.Pair[K, V]]) map[K]V {
	return Aggregate(
		s,
		make(map[K]V), // Initialize with empty map.
		func(a map[K]V, e pair.Pair[K, V]) map[K]V { // Accumulate: Add key-value pair to map.
			a[e.First()] = e.Second()
			return a
		},
		mapper.Identity[map[K]V](), // Finish: Return the map as is.
	)
}

// CollectChannel sends all elements from the stream to the given channel.
// The channel must be buffered or have a receiver ready to receive the elements in another goroutine.
// The method returns when the stream is exhausted and all elements have been sent to the channel.
//
// Example usage:
//
//	ch := make(chan int, 3)
//	go func() {
//	  for e := range ch {
//	    fmt.Println(e)
//	  }
//	}()
//	stream.CollectChannel(stream.Of(1, 2, 3), ch)
//	close(ch)
//
// Output:
//
//	1
//	2
//	3
func CollectChannel[E any](s Stream[E], ch chan<- E) {
	s(func(e E) bool {
		ch <- e
		return true
	})
}

// CollectChannelCtx behaves like CollectChannel, but it returns immediately when the context is done.
//
// Example usage:
//
//	ch := make(chan int, 3)
//	go func() {
//	  for e := range ch {
//	    fmt.Println(e)
//	  }
//	}()
//	ctx := context.Background()
//	stream.CollectChannelCtx(ctx, stream.Of(1, 2, 3), ch)
//	close(ch)
//
// Output:
//
//	1
//	2
//	3
func CollectChannelCtx[E any](ctx context.Context, s Stream[E], ch chan<- E) {
	s(func(e E) bool {
		select {
		case <-ctx.Done():
			return false
		case ch <- e:
			return true
		}
	})
}

// CollectChannelAsync sends all elements from the stream to a new channel returned by the method.
// A goroutine is started to send the elements to the channel and the method returns immediately.
// The caller must ensure the channel is fully consumed, otherwise the goroutine will be blocked forever.
// The channel is buffered with the given buffer size, and is closed when the stream is exhausted.
//
// Example usage:
//
//	ch := stream.CollectChannelAsync(stream.Of(1, 2, 3), 3)
//	for e := range ch {
//	  fmt.Println(e)
//	}
//
// Output:
//
//	1
//	2
//	3
func CollectChannelAsync[E any](s Stream[E], buf int) <-chan E {
	ch := make(chan E, buf)
	go func() {
		defer close(ch)
		s(func(e E) bool {
			ch <- e
			return true
		})
	}()
	return ch
}

// CollectChannelAsyncCtx behaves like CollectChannelAsync, but consumption of the stream stops when the context is done.
// The caller can assume the channel will eventually close when the context is done.
//
// Example usage:
//
//	ctx := context.Background()
//	ch := stream.CollectChannelAsyncCtx(ctx, stream.Of(1, 2, 3), 3)
//	for e := range ch {
//	  fmt.Println(e)
//	}
//
// Output:
//
//	1
//	2
//	3
func CollectChannelAsyncCtx[E any](ctx context.Context, s Stream[E], buf int) <-chan E {
	ch := make(chan E, buf)
	go func() {
		defer close(ch)
		s(func(e E) bool {
			select {
			case <-ctx.Done():
				return false
			case ch <- e:
				return true
			}
		})
	}()
	return ch
}
