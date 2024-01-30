package optional

import "fmt"

// Empty returns an empty Optional (None).
func Empty[V any]() Optional[V] {
	return None[V]{}
}

// Of returns a non-empty Optional (Some) wrapping the provided value.
func Of[V any](value V) Optional[V] {
	return Some[V]{Value: value}
}

// Maybe returns a non-empty Optional wrapping the provided value if ok is true; an empty Optional, otherwise.
func Maybe[V any](value V, ok bool) Optional[V] {
	if ok {
		return Some[V]{Value: value}
	}
	return None[V]{}
}

// Map returns an Optional containing the result of applying the provided function to the value contained in the provided Optional.
// If the provided Optional is empty, an empty Optional is returned.
//
// Example usage:
//
//	o := optional.Map(
//	  optional.Of(1),
//	  func(i int) string { return fmt.Sprintf("%d", i) },
//	) // optional.Some("1")
//
//	o = optional.Map(
//	  optional.Empty[int](),
//	  func(i int) string { return fmt.Sprintf("%d", i) },
//	) // optional.None()
func Map[V, U any](o Optional[V], mapper func(V) U) Optional[U] {
	if o.Present() {
		return Of(mapper(o.Get()))
	}
	return Empty[U]()
}

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

	// OrElseZero returns the value contained in the Optional if the Optional is not empty, or the zero value of type V otherwise.
	OrElseZero() V

	// OrElseGet returns the value contained in the Optional if the Optional is not empty, or the result of the provided function otherwise.
	OrElseGet(func() V) V

	// Explode returns the value contained in the Optional and an indicator of whether the Optional is empty.
	// If the Optional is empty, the value returned is the zero value of type V.
	Explode() (V, bool)

	// Filter returns an Optional containing the value contained in the Optional if the provided predicate returns true for that value.
	// If the Optional is empty, an empty Optional is returned.
	Filter(func(V) bool) Optional[V]
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

func (n None[V]) Explode() (V, bool) {
	var zero V
	return zero, false
}

func (n None[V]) Filter(_ func(V) bool) Optional[V] {
	return n
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

func (n None[V]) OrElseZero() V {
	var zero V
	return zero
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

func (s Some[V]) Get() V {
	return s.Value
}

func (s Some[V]) Explode() (V, bool) {
	return s.Value, true
}

func (s Some[V]) Filter(f func(V) bool) Optional[V] {
	if f(s.Value) {
		return s
	}
	return Empty[V]()
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

func (s Some[V]) OrElseZero() V {
	return s.Value
}

func (s Some[V]) OrElseGet(_ func() V) V {
	return s.Value
}

func (s Some[V]) String() string {
	return fmt.Sprintf("Some(%#v)", s.Value)
}
