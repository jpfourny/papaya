package pred

// Nil returns a function that returns true if the provided pointer is nil.
//
// Examples:
//
//	p := pred.Nil[int]()
//	p(nil) // true
//	p(0) // false
func Nil[E any]() func(*E) bool {
	return func(e *E) bool {
		return e == nil
	}
}

// NotNil returns a function that returns true if the provided pointer is not nil.
//
// Examples:
//
//	p := pred.NotNil[int]()
//	p(nil) // false
//	p(0) // true
func NotNil[E any]() func(*E) bool {
	return func(e *E) bool {
		return e != nil
	}
}

// Zero returns a function that returns true if the provided value is the zero value of the type parameter E.
//
// Examples:
//
//	p := pred.Zero[int]()
//	p(0) // true
//	p(1) // false
func Zero[E comparable]() func(E) bool {
	var zero E
	return func(e E) bool {
		return e == zero
	}
}

// NotZero returns a function that returns true if the provided value is not the zero value of the type parameter E.
//
// Examples:
//
//	p := pred.NotZero[int]()
//	p(0) // false
//	p(1) // true
func NotZero[E comparable]() func(E) bool {
	var zero E
	return func(e E) bool {
		return e != zero
	}
}
