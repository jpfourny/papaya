package mapper

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
