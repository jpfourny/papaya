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

// IfElse returns a function that accepts a value of any type E and returns the result of calling either the `ifTrue` or `ifFalse` function, which return a value of type F.
// If the given `cond` function returns true, the `ifTrue` function is used; otherwise, the `ifFalse` function is used.
func IfElse[E, F any](cond func(E) bool, ifTrue func(E) F, ifFalse func(E) F) func(E) F {
	return func(e E) F {
		if cond(e) {
			return ifTrue(e)
		}
		return ifFalse(e)
	}
}
