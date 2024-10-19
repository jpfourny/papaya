package stream

import (
	"github.com/jpfourny/papaya/v2/pkg/constraint"
	"github.com/jpfourny/papaya/v2/pkg/opt"
	"github.com/jpfourny/papaya/v2/pkg/stream/mapper"
	"github.com/jpfourny/papaya/v2/pkg/stream/pred"
	"strings"
)

// ForEach invokes the given consumer for each element in the stream.
//
// Example usage:
//
//	stream.ForEach(stream.Of(1, 2, 3), func(e int) bool {
//	    fmt.Println(e)
//	})
//
// Output:
//
//	1
//	2
//	3
func ForEach[E any](s Stream[E], yield func(E)) {
	s(func(e E) bool {
		yield(e)
		return true
	})
}

// IsEmpty returns true if the stream is empty; otherwise false.
//
// Example usage:
//
//	empty := stream.IsEmpty(stream.Of(1, 2, 3)) // false
//	empty := stream.IsEmpty(stream.Empty[int]()) // true
func IsEmpty[E any](s Stream[E]) (empty bool) {
	empty = true
	s(func(e E) bool {
		empty = false
		return false
	})
	return
}

// First returns the first element in the stream; an empty opt.Optional, if the stream is empty.
//
// Example usage:
//
//	out := stream.First(stream.Of(1, 2, 3)) // Some(1)
//	out = stream.First(stream.Empty[int]()) // None()
func First[E any](s Stream[E]) (first opt.Optional[E]) {
	first = opt.Empty[E]()
	s(func(e E) bool {
		first = opt.Of(e)
		return false
	})
	return
}

// Last returns the last element in the stream; an empty opt.Optional, if the stream is empty.
//
// Example usage:
//
//	out := stream.Last(stream.Of(1, 2, 3)) // Some(3)
//	out = stream.Last(stream.Empty[int]()) // None()
func Last[E any](s Stream[E]) (last opt.Optional[E]) {
	last = opt.Empty[E]()
	s(func(e E) bool {
		last = opt.Of(e)
		return true
	})
	return
}

// Peek decorates the given stream to invoke the given function for each element passing through it.
// This is useful for debugging or logging elements as they pass through the stream.
//
// Example usage:
//
//	s := stream.Peek(stream.Of(1, 2, 3), func(e int) {
//	  fmt.Println(e)
//	})
//	stream.Count(s) // Force the stream to materialize.
//
// Output:
//
//	1
//	2
//	3
func Peek[E any](s Stream[E], peek func(e E)) Stream[E] {
	return func(yield Consumer[E]) {
		s(func(e E) bool {
			peek(e)
			return yield(e)
		})
	}
}

// Walk returns a stream that walks elements beginning at `start`, advanced by the `advance` function, and ending when `cond` predicate returns false.
//
// Example usage:
//
//	s := stream.Walk(1, pred.LessThanOrEqual(5), mapper.Increment(2))
//	out := stream.DebugString(s) // "<1, 3, 5>"
func Walk[E any](start E, cond Predicate[E], advance Mapper[E, E]) Stream[E] {
	next := start
	return func(yield Consumer[E]) {
		for cond(next) {
			if !yield(next) {
				return
			}
			next = advance(next)
		}
		return
	}
}

// Interval returns a stream of real-number type N from the half-open interval `[start, end)` using the given step size.
// If the step size is negative, then the stream will be decreasing; otherwise, it will be increasing.
//
// Example usage:
//
//	s := stream.Interval(0, 5, 1)
//	out := stream.DebugString(s) // "<0, 1, 2, 3, 4>"
//	s := stream.Interval(0, 5, 2)
//	out := stream.DebugString(s) // "<0, 2, 4>"
//	s = stream.Interval(5, 0, -2)
//	out = stream.DebugString(s) // "<5, 3, 1>"
func Interval[N constraint.RealNumber](start, end, step N) Stream[N] {
	if step < 0 {
		return Walk[N](start, pred.GreaterThan(end), mapper.Increment(step))
	}
	return Walk[N](start, pred.LessThan(end), mapper.Increment(step))
}

// StringJoin concatenates the elements of the provided stream of strings into a single string, using the specified separator.
//
// Example usage:
//
//	s := stream.Of("foo", "bar", "baz")
//	out := stream.StringJoin(s, ", ") // "foo, bar, baz"
func StringJoin(s Stream[string], sep string) string {
	return Aggregate(
		s,
		&strings.Builder{},
		func(a *strings.Builder, e string) *strings.Builder {
			if a.Len() > 0 {
				a.WriteString(sep)
			}
			a.WriteString(e)
			return a
		},

		func(a *strings.Builder) string { return a.String() },
	)
}

// Pad returns a stream that pads the tail of the given stream with the given 'pad' value until the stream reaches the given length.
// If the stream is already longer than the given length, then the stream is returned as-is.
//
// Example usage:
//
//	s := stream.Pad(stream.Of(1, 2, 3), 0, 5)
//	out := stream.DebugString(s) // "<1, 2, 3, 0, 0>"
func Pad[E any](s Stream[E], pad E, length int) Stream[E] {
	return func(yield Consumer[E]) {
		yield2, stopped := stopSensingConsumer(yield)

		i := 0
		s(func(e E) bool {
			i++
			return yield2(e)
		})
		if *stopped {
			return // Consumer saw enough.
		}
		for ; i < length; i++ {
			if !yield(pad) {
				return // Consumer saw enough.
			}
		}
	}
}

// Truncate returns a stream that limits the given stream to the desired length and appends the given 'tail' value, if the stream is longer than the desired length.
// The tail value is appended only once, even if the stream is longer than the desired.
// If the stream is already shorter than the desired length, then the stream is returned as-is.
//
// Example usage:
//
//	s := stream.Truncate(stream.Of("a", "b", "c""), 2, "...")
//	out := stream.DebugString(s) // "<a, b, ...>"
//
//	s = stream.Truncate(stream.Of("a", "b", "c""), 3, "...")
//	out = stream.DebugString(s) // "<a, b, c>"
func Truncate[E any](s Stream[E], length int, tail E) Stream[E] {
	return func(yield Consumer[E]) {
		i := 0
		s(func(e E) bool {
			i++
			if i <= length {
				return yield(e)
			}
			yield(tail)
			return false // Stop after the tail.
		})
	}
}

// DebugString returns a string representation of up to the first 100 elements in the stream.
// The string is formatted like `<e1, e2, e3>` where e1, e2, e3 are the string representations of the elements.
// If the stream has more than 100 elements, the string will end with `...>` to indicate that the stream was truncated.
// Useful for debugging.
func DebugString[E any](s Stream[E]) string {
	return "<" +
		StringJoin(
			Truncate(
				Map(s, mapper.Sprintf[E]("%#v")),
				100,
				"...",
			),
			", ",
		) + ">"
}

// Cache returns a stream that caches all elements from the given stream.
// The first call to the returned stream will cache all elements from the given stream (and exhaust it).
// Subsequent calls to the returned stream will replay the cache.
func Cache[E any](s Stream[E]) Stream[E] {
	var cache []E
	return func(yield Consumer[E]) {
		if cache == nil {
			// First call to stream caches all elements.
			cache = []E{}
			cont := true
			s(func(e E) bool {
				cache = append(cache, e)
				if cont { // Yield until the consumer requests to stop.
					cont = yield(e)
				}
				return true // Always exhaust the stream for caching.
			})
		} else {
			// Replay from the cache.
			FromSlice(cache)(yield)
		}
	}
}
