package stream

// Stream represents a function that produces a sequence of elements of type E and sends them to the given Consumer.
//
// Streams are lazy, meaning they only produce elements when the consumer is invoked.
// Furthermore, streams are idempotent, meaning they can be invoked multiple times with the same result.
// However, the order of the elements is not guaranteed to be the same across multiple invocations.
//
// If the Consumer returns false, the stream must stop producing elements and return false immediately.
// If the stream is exhausted, it must return true.
type Stream[E any] func(c Consumer[E])

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
	return func(_ Consumer[E]) {
		return
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
