package mapper

import "github.com/jpfourny/papaya/pkg/constraint"

// Increment returns a function that accepts a value of real number type E and returns the result of adding the provided `step` value to it.
func Increment[E constraint.RealNumber](step E) func(E) E {
	return func(e E) E {
		return e + step
	}
}

// Decrement returns a function that accepts a value of real number type E and returns the result of subtracting the provided `step` value from it.
func Decrement[E constraint.RealNumber](step E) func(E) E {
	return func(e E) E {
		return e - step
	}
}
