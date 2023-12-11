package optional

import "fmt"

// Optional is a generic type that takes one type parameter V and represents a value that may or may not be present.
// It is similar to Java's Optional type.
type Optional[V any] interface {
	fmt.Stringer

	// IsPresent returns true if the Optional contains a value, false otherwise.
	IsPresent() bool

	// Get returns the value contained in the Optional, or the zero value of type V if the Optional is empty.
	Get() V

	// IfPresent calls the provided function with the value contained in the Optional if the Optional is not empty.
	IfPresent(func(V))

	// OrElse returns the value contained in the Optional if the Optional is not empty, or the provided value otherwise.
	OrElse(V) V

	// OrElseGet returns the value contained in the Optional if the Optional is not empty, or the result of the provided function otherwise.
	OrElseGet(func() V) V
}

// None is an Optional that represents the absence of a value.
type None[V any] struct {
}

func (n None[V]) IsPresent() bool {
	return false
}

func (n None[V]) Get() V {
	var zero V
	return zero
}

func (n None[V]) IfPresent(func(V)) {
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

func (s Some[V]) IsPresent() bool {
	return true
}

func (s Some[V]) IfPresent(f func(V)) {
	f(s.Value)
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

// OfNillable returns an Optional containing the provided value, or None if the value is nil.
func OfNillable[V any](value *V) Optional[V] {
	if value == nil {
		return None[V]{}
	}
	return Some[V]{Value: *value}
}

// Empty returns an empty Optional.
func Empty[V any]() Optional[V] {
	return None[V]{}
}
