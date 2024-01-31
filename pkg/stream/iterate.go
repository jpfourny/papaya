package stream

import (
	"github.com/jpfourny/papaya/pkg/constraint"
	"github.com/jpfourny/papaya/pkg/opt"
	"github.com/jpfourny/papaya/pkg/stream/mapper"
	"github.com/jpfourny/papaya/pkg/stream/pred"
)

// Iterate returns a stream of elements produced by the given iterator function.
// When the iterator function returns false, the stream ends.
//
// Example usage:
//
//	i := 0
//	s := stream.Iterate(func() (int, bool) {
//	  if i < 3 {
//	    i++
//	    return i, true
//	  }
//	  return 0, false
//	})
//	out := stream.DebugString(s) // "<1, 2, 3>"
func Iterate[E any](next func() (E, bool)) Stream[E] {
	return func(yield Consumer[E]) bool {
		for e, ok := next(); ok; e, ok = next() {
			if !yield(e) {
				return false
			}
		}
		return true
	}
}

// IterateOptional returns a stream of elements produced by the given iterator function.
// When the iterator function returns an empty opt.Optional, the stream ends.
//
// Example usage:
//
//	i := 0
//	s := stream.IterateOptional(func() opt.Optional[int] {
//	  if i < 3 {
//	    i++
//	    return opt.Of(i)
//	  }
//	  return opt.Empty[int]()
//	})
//	out := stream.DebugString(s) // "<1, 2, 3>"
func IterateOptional[E any](next func() opt.Optional[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		for e := next(); e.Present(); e = next() {
			if !yield(e.GetOrZero()) {
				return false
			}
		}
		return true
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
	return Iterate(func() (e E, ok bool) {
		if cond(next) {
			e, ok = next, true
			next = advance(next)
		}
		return
	})
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
