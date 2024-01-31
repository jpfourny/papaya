package stream

import (
	"github.com/jpfourny/papaya/pkg/opt"
	"github.com/jpfourny/papaya/pkg/stream/mapper"
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
	return func(yield Consumer[E]) bool {
		return s(func(e E) bool {
			peek(e)
			return yield(e)
		})
	}
}

// DebugString returns a string representation of up to the first 100 elements in the stream.
// The string is formatted like `<e1, e2, e3>` where e1, e2, e3 are the string representations of the elements.
// If the stream has more than 100 elements, the string will end with `...>` to indicate that the stream was truncated.
// Useful for debugging.
func DebugString[E any](s Stream[E]) string {
	return "<" + StringJoin(Truncate(Map(s, mapper.Sprintf[E]("%#v")), 100, "..."), ", ") + ">"
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

// Cache returns a stream that caches all elements from the given stream.
// The first call to the returned stream will cache all elements from the given stream (and exhaust it).
// Subsequent calls to the returned stream will replay the cache.
func Cache[E any](s Stream[E]) Stream[E] {
	var cache []E
	return func(yield Consumer[E]) bool {
		if cache == nil {
			// First call to stream caches all elements.
			cache = []E{}
			ok := true
			s(func(e E) bool {
				cache = append(cache, e)
				if ok { // Yield until the consumer returns false.
					ok = yield(e)
				}
				return true // Always exhaust the stream for caching.
			})
			return ok
		} else {
			// Replay from the cache.
			for _, e := range cache {
				if !yield(e) {
					return false
				}
			}
		}
		return true
	}
}
