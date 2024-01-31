package ptr

import "github.com/jpfourny/papaya/pkg/opt"

// Nil returns a nil pointer to the type parameter V.
func Nil[V any]() *V {
	return nil
}

// Ref returns a pointer to a copy of the given value (on the heap).
func Ref[V any](v V) *V {
	return &v
}

// DerefOptional returns the value from de-referencing the given pointer wrapped in an opt.Optional, or an empty opt.Optional if the pointer is nil.
func DerefOptional[V any](v *V) opt.Optional[V] {
	if v == (*V)(nil) {
		return opt.Empty[V]()
	}
	return opt.Of[V](*v)
}

// DerefOrDefault returns the value from de-referencing the given pointer, or the provided defaultValue if the pointer is nil.
func DerefOrDefault[V any](v *V, defaultValue V) V {
	return DerefOptional(v).GetOrDefault(defaultValue)
}

// DerefOrZero returns the value after de-referencing the given pointer, or the zero value of the type parameter V if the pointer is nil.
func DerefOrZero[V any](v *V) V {
	return DerefOptional(v).GetOrZero()
}

// DerefOrFunc returns the value after de-referencing the provided ptr, or the result of calling the provided function if the ptr is nil.
func DerefOrFunc[V any](v *V, f func() V) V {
	return DerefOptional(v).GetOrFunc(f)
}
