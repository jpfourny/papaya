package stream

import (
	"strings"

	"github.com/jpfourny/papaya/pkg/mapper"
	"github.com/jpfourny/papaya/pkg/optional"
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

// Peek decorates the given stream to invoke the given function for each element passing through it.
// This is useful for debugging or logging elements as they pass through the stream.
//
// Example usage:
//
//	s := stream.Peek(stream.Of(1, 2, 3), func(e int) {
//	    fmt.Println(e)
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

// First returns the first element in the stream; an empty optional.Optional, if the stream is empty.
//
// Example usage:
//
//	out := stream.First(stream.Of(1, 2, 3)) // Some(1)
//	out = stream.First(stream.Empty[int]()) // None()
func First[E any](s Stream[E]) (first optional.Optional[E]) {
	first = optional.Empty[E]()
	s(func(e E) bool {
		first = optional.Of(e)
		return false
	})
	return
}

// Last returns the last element in the stream; an empty optional.Optional, if the stream is empty.
//
// Example usage:
//
//	out := stream.Last(stream.Of(1, 2, 3)) // Some(3)
//	out = stream.Last(stream.Empty[int]()) // None()
func Last[E any](s Stream[E]) (last optional.Optional[E]) {
	last = optional.Empty[E]()
	s(func(e E) bool {
		last = optional.Of(e)
		return true
	})
	return
}
