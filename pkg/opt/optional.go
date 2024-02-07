package opt

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

// Any returns the first non-empty Optional from the provided options.
// If all options are empty, an empty Optional is returned.
func Any[V any](options ...Optional[V]) Optional[V] {
	for _, o := range options {
		if o.Present() {
			return o
		}
	}
	return None[V]{}
}

// Map returns an Optional containing the result of applying the provided mapper function to the value contained in the provided Optional.
// If the provided Optional is empty, an empty Optional is returned; otherwise, a non-empty Optional is returned.
//
// Example usage:
//
//	o := opt.Map(
//	  opt.Of(1),
//	  func(i int) string { return fmt.Sprintf("%d", i) },
//	) // opt.Some("1")
//
//	o = opt.Map(
//	  opt.Empty[int](),
//	  func(i int) string { return fmt.Sprintf("%d", i) },
//	) // opt.None()
func Map[V, U any](o Optional[V], mapper func(V) U) Optional[U] {
	if value, ok := o.Get(); ok {
		return Of[U](mapper(value))
	}
	return Empty[U]()
}

// OptionalMap returns an Optional containing the result of applying the provided mapper function to the value contained in the provided Optional.
// If the provided Optional is empty, of if the mapper returns an empty Optional, an empty Optional is returned; otherwise, a non-empty Optional is returned.
//
// Example usage:
//
//	o := opt.OptionalMap(
//	  opt.Of(1),
//	  func(i int) opt.Optional[string] { return opt.Of(fmt.Sprintf("%d", i)) },
//	) // opt.Some("1")
//
//	o = opt.OptionalMap(
//	  opt.Of(1),
//	  func(i int) opt.Optional[string] { return opt.Empty[string]() },
//	) // opt.None()
//
//	o = opt.OptionalMap(
//	  opt.Empty[int](),
//	  func(i int) opt.Optional[string] { return opt.Of(fmt.Sprintf("%d", i)) },
//	) // opt.None()
func OptionalMap[V, U any](o Optional[V], mapper func(V) Optional[U]) Optional[U] {
	if value, ok := o.Get(); ok {
		return mapper(value)
	}
	return Empty[U]()
}

// MaybeMap returns an Optional containing the result of applying the provided mapper function to the value contained in the provided Optional.
// If the provided Optional is empty, or if the mapper returns false, an empty Optional is returned; otherwise, a non-empty Optional is returned.
//
// Example usage:
//
//	o := opt.MaybeMap(
//	  opt.Of(1),
//	  func(i int) (string, bool) { return fmt.Sprintf("%d", i), true },
//	) // opt.Some("1")
//
//	o = opt.MaybeMap(
//	  opt.Of(1),
//	  func(i int) (string, bool) { return fmt.Sprintf("%d", i), false },
//	) // opt.None()
//
//	o = opt.MaybeMap(
//	  opt.Empty[int](),
//	  func(i int) (string, bool) { return fmt.Sprintf("%d", i), true },
//	) // opt.None()
func MaybeMap[V, U any](o Optional[V], mapper func(V) (U, bool)) Optional[U] {
	if value, ok := o.Get(); ok {
		return Maybe(mapper(value))
	}
	return Empty[U]()
}

// Optional is a generic type that takes one type parameter V and represents a value that may or may not be present.
// It is similar to Java's Optional type.
type Optional[V any] interface {
	fmt.Stringer

	// Present returns true if the Optional contains a value, false otherwise.
	Present() bool

	// Get returns the value contained in the Optional and an indicator of whether the Optional is empty.
	// If the Optional is empty, the value returned is the zero value of type V.
	Get() (V, bool)

	// GetOrZero returns the value contained in the Optional, or the zero value of type V if the Optional is empty.
	GetOrZero() V

	// GetOrDefault returns the value contained in the Optional, or the provided default value if the Optional is empty.
	GetOrDefault(defaultValue V) V

	// GetOrFunc returns the value contained in the Optional, or the result of calling the provided function if the Optional is empty.
	GetOrFunc(func() V) V

	// IfPresent calls the provided function with the value contained in the Optional if the Optional is not empty.
	// Returns true if the function was called, false otherwise.
	IfPresent(func(V)) bool

	// IfPresentElse calls the first function with the value contained in the Optional if the Optional is not empty, or the second function otherwise.
	// Returns true if the first function was called, false otherwise.
	IfPresentElse(func(V), func()) bool

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

func (n None[V]) Get() (V, bool) {
	var zero V
	return zero, false
}

func (n None[V]) GetOrZero() V {
	var zero V
	return zero
}

func (n None[V]) GetOrDefault(defaultValue V) V {
	return defaultValue
}

func (n None[V]) GetOrFunc(f func() V) V {
	return f()
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

func (s Some[V]) Get() (V, bool) {
	return s.Value, true
}

func (s Some[V]) GetOrZero() V {
	return s.Value
}

func (s Some[V]) GetOrDefault(_ V) V {
	return s.Value
}

func (s Some[V]) GetOrFunc(_ func() V) V {
	return s.Value
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

func (s Some[V]) String() string {
	return fmt.Sprintf("Some(%#v)", s.Value)
}
