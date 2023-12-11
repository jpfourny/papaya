package pair

import "fmt"

// Pair represents a pair of two values.
// It is a generic type that takes two type parameters A and B.
type Pair[A, B any] struct {
	first  A
	second B
}

// First returns the first element of a Pair.
func (p Pair[A, B]) First() A {
	return p.first
}

// Second returns the second element of a Pair.
func (p Pair[A, B]) Second() B {
	return p.second
}

// Explode returns the first and second elements of a Pair.
func (p Pair[A, B]) Explode() (A, B) {
	return p.first, p.second
}

// Reverse returns a new Pair with the first and second elements of the original Pair reversed.
func (p Pair[A, B]) Reverse() Pair[B, A] {
	return Pair[B, A]{first: p.second, second: p.first}
}

// String returns a string representation of the Pair, formatted as "(%#v, %#v)" where the first value is %#v and the second value is %#v.
func (p Pair[A, B]) String() string {
	return fmt.Sprintf("(%#v, %#v)", p.first, p.second)
}

// Zero returns a new Pair with default zero values for the type parameters A and B.
func Zero[A, B any]() Pair[A, B] {
	return Pair[A, B]{}
}

// Of creates a new Pair with the provided values for the first and second elements.
// The type parameters A and B determine the types of the first and second elements, respectively.
func Of[A, B any](first A, second B) Pair[A, B] {
	return Pair[A, B]{first: first, second: second}
}
