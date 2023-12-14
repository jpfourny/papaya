package stream

// Union combines multiple streams into a single stream (concatenation).
// The length of the resulting stream is the sum of the lengths of the input streams.
// If any of the input streams return false when invoked with the consumer, the concatenation stops.
//
// Example usage:
//
//	s := stream.Union(stream.Of(1, 2, 3, 4), stream.Of(4, 5, 6))
//	out := stream.DebugString(s) // "<1, 2, 3, 4, 4, 5, 6>"
func Union[E any](ss ...Stream[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		for _, s := range ss {
			if !s(yield) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}

// Intersection returns a stream that contains elements that are in all the given streams.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.Intersection(stream.Of(1, 2, 3, 4, 5), stream.Of(4, 5, 6), stream.Of(4, 5))
//	out := stream.DebugString(s) // "<4, 5>"
func Intersection[E comparable](ss ...Stream[E]) Stream[E] {
	switch {
	case len(ss) == 0:
		return Empty[E]() // No streams.
	case len(ss) == 1:
		return ss[0] // One stream.
	case len(ss) > 2:
		// Recursively intersect the first stream with the intersection of rest.
		return Intersection(ss[0], Intersection(ss[1:]...))
	}

	return func(yield Consumer[E]) bool {
		// Index elements of the first stream into a set.
		seen := make(map[E]struct{})
		ss[0](func(e E) bool {
			seen[e] = struct{}{}
			return true
		})
		// Yield elements of the second stream that are in the set.
		return ss[1](func(e E) bool {
			if _, ok := seen[e]; ok {
				return yield(e)
			}
			return true
		})
	}
}

// Difference returns a stream that contains elements that are in the first stream but not in the second stream.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.Difference(stream.Of(1, 2, 3, 4, 5), stream.Of(4, 5, 6))
//	out := stream.DebugString(s) // "<1, 2, 3>"
func Difference[E comparable](s1, s2 Stream[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		// Index elements of the second stream into a set.
		seen := make(map[E]struct{})
		s2(func(e E) bool {
			seen[e] = struct{}{}
			return true
		})
		// Yield elements of the first stream that are not in the set.
		return s1(func(e E) bool {
			if _, ok := seen[e]; !ok {
				return yield(e)
			}
			return true
		})
	}
}
