package stream

// Consumer represents a function that accepts a yielded element of type E and returns a boolean value.
// The boolean value indicates whether the consumer wishes to continue accepting elements.
// If the consumer returns false, the caller must stop yielding elements.
type Consumer[E any] func(yield E) (cont bool)

func stopSensingConsumer[E any](c Consumer[E]) (Consumer[E], *bool) {
	stopped := new(bool)
	c2 := func(e E) bool {
		if c(e) {
			return true
		}
		*stopped = true
		return false
	}
	return c2, stopped
}
