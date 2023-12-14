package mapper

import "github.com/jpfourny/papaya/pkg/pointer"

// PointerRef returns a function that accepts any value and returns a pointer to a copy of that value.
// The pointer is created using the pointer.Ref function.
func PointerRef[E any]() func(E) *E {
	return func(e E) *E {
		return pointer.Ref(e)
	}
}

// PointerDerefOr returns a function that accepts a pointer to a value and returns the value of that pointer, or the provided `or` default value if the pointer is nil.
// The value is retrieved using the pointer.DerefOr function.
func PointerDerefOr[E any](or E) func(*E) E {
	return func(e *E) E {
		return pointer.DerefOr(e, or)
	}
}

// PointerDerefOrZero returns a function that accepts a pointer to a value and returns the value of that pointer, or the zero value of that type if the pointer is nil.
// The value is retrieved using the pointer.DerefOrZero function.
func PointerDerefOrZero[E any]() func(*E) E {
	return func(e *E) E {
		return pointer.DerefOrZero(e)
	}
}
