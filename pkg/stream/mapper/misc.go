package mapper

import "github.com/jpfourny/papaya/v2/pkg/constraint"

// Identity returns a function that accepts a value of any type E and returns that value.
func Identity[E any]() func(E) E {
	return func(e E) E {
		return e
	}
}

// Constant returns a function that accepts a value of any type E and returns the provided constant value of type F.
func Constant[E any, F any](c F) func(E) F {
	return func(E) F {
		return c
	}
}

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
