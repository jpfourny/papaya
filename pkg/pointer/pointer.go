package pointer

import "github.com/jpfourny/papaya/pkg/optional"

// Nil returns a nil pointer to the type parameter E.
func Nil[E any]() *E {
	return nil
}

// Ref returns a pointer to a copy of the provided value.
func Ref[E any](e E) *E {
	return &e
}

// Deref returns the value after dereferencing the provided pointer as an optional.Optional, or an empty optional if the pointer is nil.
func Deref[E any](e *E) optional.Optional[E] {
	if e == nil {
		return optional.Empty[E]()
	}
	return optional.Of[E](*e)
}

// DerefOr returns the value after dereferencing the provided pointer, or the provided `or` default value if the pointer is nil.
func DerefOr[E any](e *E, or E) E {
	if e == nil {
		return or
	}
	return *e
}

// DerefOrZero returns the value after dereferencing the provided pointer, or the zero value of the type parameter E if the pointer is nil.
func DerefOrZero[E any](e *E) E {
	if e == nil {
		var zero E
		return zero
	}
	return *e
}

// DerefOrFunc returns the value after dereferencing the provided pointer, or the result of calling the provided function if the pointer is nil.
func DerefOrFunc[E any](e *E, f func() E) E {
	if e == nil {
		return f()
	}
	return *e
}
