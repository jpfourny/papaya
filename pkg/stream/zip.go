package stream

import "github.com/jpfourny/papaya/pkg/pair"

// Zip returns a stream that pairs each element in the first stream with the corresponding element in the second stream.
// The resulting stream will have the same number of elements as the shorter of the two input streams.
//
// Example usage:
//
//	s := stream.Zip(stream.Of(1, 2, 3), stream.Of("foo", "bar"))
//	out := stream.DebugString(s) // "<(1, foo), (2, bar)>"
func Zip[E, F any](s1 Stream[E], s2 Stream[F]) Stream[pair.Pair[E, F]] {
	return func(yield Consumer[pair.Pair[E, F]]) bool {
		done := make(chan struct{})
		defer close(done)

		ch1 := make(chan E)
		go func() {
			defer close(ch1)
			s1(func(e E) bool {
				select {
				case <-done:
					return false
				case ch1 <- e:
					return true
				}
			})
		}()

		ch2 := make(chan F)
		go func() {
			defer close(ch2)
			s2(func(f F) bool {
				select {
				case <-done:
					return false
				case ch2 <- f:
					return true
				}
			})
		}()

		for {
			e, ok1 := <-ch1
			f, ok2 := <-ch2
			if !ok1 || !ok2 {
				return true
			}
			if !yield(pair.Of(e, f)) {
				return false
			}
		}
	}
}

// ZipWithIndexInt returns a stream that pairs each element in the input stream with its index, starting at the given offset.
// The index is of type int.
//
// Example usage:
//
//	s := stream.ZipWithIndexInt(stream.Of("foo", "bar"), 1)
//	out := stream.DebugString(s) // "<(foo, 1), (bar, 2)>"
func ZipWithIndexInt[E any](s Stream[E], offset int) Stream[pair.Pair[E, int]] {
	return func(yield Consumer[pair.Pair[E, int]]) bool {
		i := offset - 1
		return s(func(e E) bool {
			i++
			return yield(pair.Of(e, i))
		})
	}
}

// ZipWithIndexInt64 returns a stream that pairs each element in the input stream with its index, starting at the given offset.
// The index is of type int64.
//
// Example usage:
//
//	s := stream.ZipWithIndexInt64(stream.Of("foo", "bar"), 1)
//	out := stream.DebugString(s) // "<(foo, 1), (bar, 2)>"
func ZipWithIndexInt64[E any](s Stream[E], offset int64) Stream[pair.Pair[E, int64]] {
	return func(yield Consumer[pair.Pair[E, int64]]) bool {
		i := offset - 1
		return s(func(e E) bool {
			i++
			return yield(pair.Of(e, i))
		})
	}
}

// ZipWithKey returns a stream that pairs each element in the input stream with the key extracted from the element using the given key extractor.
// The resulting stream will have the same number of elements as the input stream.
//
// Example usage:
//
//	s := stream.ZipWithKey(stream.Of("foo", "bar"), func(s string) string {
//	    return strings.ToUpper(s)
//	})
//	out := stream.DebugString(s) // "<(FOO, foo), (BAR, bar)>"
func ZipWithKey[E any, K any](s Stream[E], ke KeyExtractor[E, K]) Stream[pair.Pair[K, E]] {
	return func(yield Consumer[pair.Pair[K, E]]) bool {
		return s(func(e E) bool {
			return yield(pair.Of(ke(e), e))
		})
	}
}
