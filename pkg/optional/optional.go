package optional

import "fmt"

// Optional is a generic type that takes one type parameter V and represents a value that may or may not be present.
// It is similar to Java's Optional type.
type Optional[V any] interface {
	fmt.Stringer

	// Present returns true if the Optional contains a value, false otherwise.
	Present() bool

	// Get returns the value contained in the Optional, or the zero value of type V if the Optional is empty.
	Get() V

	// IfPresent calls the provided function with the value contained in the Optional if the Optional is not empty.
	// Returns true if the function was called, false otherwise.
	IfPresent(func(V)) bool

	// IfPresentElse calls the first function with the value contained in the Optional if the Optional is not empty, or the second function otherwise.
	// Returns true if the first function was called, false otherwise.
	IfPresentElse(func(V), func()) bool

	// OrElse returns the value contained in the Optional if the Optional is not empty, or the provided value otherwise.
	OrElse(V) V

	// OrElseGet returns the value contained in the Optional if the Optional is not empty, or the result of the provided function otherwise.
	OrElseGet(func() V) V
}

// None is an Optional that represents the absence of a value.
type None[V any] struct {
}

func (n None[V]) Present() bool {
	return false
}

func (n None[V]) Get() V {
	var zero V
	return zero
}

func (n None[V]) IfPresent(_ func(V)) bool {
	return false
}

func (n None[V]) IfPresentElse(_ func(V), f func()) bool {
	f()
	return false
}

func (n None[V]) OrElse(v V) V {
	return v
}

func (n None[V]) OrElseGet(f func() V) V {
	return f()
}

func (n None[V]) String() string {
	return "None"
}

// Some is an Optional that represents the presence of a value.
type Some[V any] struct {
	Value V
}

func (s Some[V]) Present() bool {
	return true
}

func (s Some[V]) IfPresent(f func(V)) bool {
	f(s.Value)
	return true
}

func (s Some[V]) IfPresentElse(f func(V), _ func()) bool {
	f(s.Value)
	return true
}

func (s Some[V]) OrElse(V) V {
	return s.Value
}

func (s Some[V]) OrElseGet(_ func() V) V {
	return s.Value
}

func (s Some[V]) Get() V {
	return s.Value
}

func (s Some[V]) String() string {
	return fmt.Sprintf("Some(%#v)", s.Value)
}

// Of returns an Optional containing the provided value.
func Of[V any](value V) Optional[V] {
	return Some[V]{Value: value}
}

// Maybe returns an Optional containing the provided value if the provided boolean is true, or an empty Optional otherwise.
func Maybe[V any](value V, ok bool) Optional[V] {
	if ok {
		return Some[V]{Value: value}
	}
	return None[V]{}
}

// Empty returns an empty Optional.
func Empty[V any]() Optional[V] {
	return None[V]{}
}
