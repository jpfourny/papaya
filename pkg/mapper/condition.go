package mapper

import "github.com/jpfourny/papaya/pkg/optional"

// If returns a function that accepts a value of any type E and returns the result of calling the `ifTrue` function as an optional.Optional or an empty optional, if the `cond` function returns false.
//
// Example usage:
//
//	m := mapper.If[int, string](
//	  pred.GreaterThan(0),
//	  mapper.Constant[int]("positive"),
//	)
//	out := m(-1) // optional.None
//	out = m(0)   // optional.None
//	out = m(1)   // optional.Some("positive")
func If[E, F any](cond func(E) bool, ifTrue func(E) F) func(E) optional.Optional[F] {
	return func(e E) optional.Optional[F] {
		if cond(e) {
			return optional.Of[F](ifTrue(e))
		}
		return optional.Empty[F]()
	}
}

// IfElse returns a function that accepts a value of any type E and returns the result of calling either the `ifTrue` or `ifFalse` function, which return a value of type F.
// If the given `cond` function returns true, the `ifTrue` function is used; otherwise, the `ifFalse` function is used.
//
// Example usage:
//
//	m := mapper.IfElse[int, string](
//	  pred.GreaterThan(0),
//	  mapper.Constant[int]("positive"),
//	  mapper.Constant[int]("negative"),
//	)
//	out := m(-1) // "negative"
//	out = m(0)   // "negative"
//	out = m(1)   // "positive"
func IfElse[E, F any](cond func(E) bool, ifTrue func(E) F, ifFalse func(E) F) func(E) F {
	return func(e E) F {
		if cond(e) {
			return ifTrue(e)
		}
		return ifFalse(e)
	}
}

// Case represents a conditional mapper function.
// It is used in the Switch operation.
type Case[E, F any] struct {
	Cond   func(E) bool
	Mapper func(E) F
}

// Switch returns a function that accepts a value of any type E and returns the result of applying the `Mapper` from the first case whose `Cond` returns true as an optional.Optional.
// If no case's `Cond` returns true, an empty optional is returned.
//
// Example usage:
//
//	m := mapper.Switch[int, string](
//	  []mapper.Case[int, string] {
//	    { Cond: pred.GreaterThan(0), Mapper: mapper.Constant[int]("positive") },
//	    { Cond: pred.LessThan(0), Mapper: mapper.Constant[int]("negative") },
//	  },
//	)
//	out := m(-1) // optional.Some("negative")
//	out = m(0)   // optional.None
//	out = m(1)   // optional.Some("positive")
func Switch[E, F any](cases []Case[E, F]) func(E) optional.Optional[F] {
	return func(e E) optional.Optional[F] {
		for _, c := range cases {
			if c.Cond(e) {
				return optional.Of[F](c.Mapper(e))
			}
		}
		return optional.Empty[F]()
	}
}

// SwitchWithDefault returns a function that accepts a value of any type E and returns the result of applying the `Mapper` from the first case whose `Cond` returns true.
// If no case's `Cond` returns true, the `defaultMapper` function is used.
//
// Example usage:
//
//	m := mapper.SwitchWithDefault[int, string](
//	  []mapper.Case[int, string] {
//	    { Cond: pred.GreaterThan(0), Mapper: mapper.Constant[int]("positive") },
//	    { Cond: pred.LessThan(0), Mapper: mapper.Constant[int]("negative") },
//	  },
//	  mapper.Constant[int]("neutral"), // Default case.
//	)
//	out := m(-1) // "negative"
//	out = m(0)   // "neutral"
//	out = m(1)   // "positive"
func SwitchWithDefault[E, F any](cases []Case[E, F], defaultMapper func(E) F) func(E) F {
	return func(e E) F {
		for _, c := range cases {
			if c.Cond(e) {
				return c.Mapper(e)
			}
		}
		return defaultMapper(e)
	}
}
