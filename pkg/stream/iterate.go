package stream

import "github.com/jpfourny/papaya/pkg/optional"

// Iterator represents a function that produces a (typically finite) sequence of elements of type E.
// It is similar to Generator, but it can return an empty optional.Optional to signal the end of the sequence.
// Used by the Iterate function.
type Iterator[E any] func() optional.Optional[E]

// Iterate returns a stream of elements produced by the given Iterator.
// When the Iterator returns an empty optional.Optional, the stream ends.
//
// Example usage:
//
//	i := 0
//	s := stream.Iterate(func() optional.Optional[int] {
//	  if i < 3 {
//	    i++
//	    return optional.Of(i)
//	  }
//	  return optional.Empty[int]()
//	})
//	out := stream.DebugString(s) // "<1, 2, 3>"
func Iterate[E any](next Iterator[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		for e := next(); e.Present(); e = next() {
			if !yield(e.Get()) {
				return false
			}
		}
		return true
	}
}

// Range returns a stream that produces elements over a range beginning at `start`, advanced by the `next` function, and ending when `cond` predicate returns false.
//
// Example usage:
//
//	s := stream.Range(1, pred.LessThanOrEqual(5), mapper.Increment(2))
//	out := stream.DebugString(s) // "<1, 3, 5>"
func Range[E any](start E, cond Predicate[E], next Mapper[E, E]) Stream[E] {
	e := start
	return Iterate(func() (result optional.Optional[E]) {
		if cond(e) {
			result = optional.Of(e)
			e = next(e)
		} else {
			result = optional.Empty[E]()
		}
		return
	})
}
