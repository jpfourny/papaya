package cmp

import (
	stdcmp "cmp"

	"github.com/jpfourny/papaya/pkg/constraint"
)

// Comparer is a function that compares two values of the same type E and returns an integer.
// It returns a negative integer if a < b, zero if a == b, or a positive integer if a > b.
type Comparer[E any] func(a, b E) int

// Reverse returns a Comparer that compares two values of the same type E in reverse order.
func (c Comparer[E]) Reverse() Comparer[E] {
	return func(a, b E) int {
		return c(b, a)
	}
}

// Then returns a Comparer that compares two values of the same type E using the provided 'other' Comparer if the receiver Comparer returns zero.
// If the receiver Comparer returns a non-zero value, the provided 'other' Comparer is not used.
//
// Example:
//
//	type Person struct {
//		FirstName string
//		LastName  string
//	}
//
//	people := []Person{
//		{"John", "Doe"},
//		{"Jane", "Doe"},
//		{"John", "Smith"},
//	}
//
//	// Sort by LastName, then by FirstName.
//	s := stream.SortBy(
//		stream.FromSlice(people),
//		Comparing(func(p Person) string { return p.LastName }).
//			Then(Comparing(func(p Person) string { return p.FirstName })),
//	)
//	out := stream.DebugString(s) // <Person{FirstName:"Jane", LastName:"Doe"}, Person{FirstName:"John", LastName:"Doe"}, Person{FirstName:"John", LastName:"Smith"}>
func (c Comparer[E]) Then(other Comparer[E]) Comparer[E] {
	return func(a, b E) int {
		if r := c(a, b); r != 0 {
			return r
		}
		return other(a, b)
	}
}

// KeyExtractor is a function that extracts a sort key of type K from a value of type E.
type KeyExtractor[E, K any] func(E) K

// Natural returns a Comparer that compares two values of the same type E using the standard library's Compare function for type E.
// The type parameter E must implement the Ordered constraint.
func Natural[E constraint.Ordered]() Comparer[E] {
	return stdcmp.Compare[E]
}

// Reverse returns a Comparer that compares two values of the same type E in reverse order using the standard library's Compare function for type E.
// The type parameter E must implement the Ordered constraint.
func Reverse[E constraint.Ordered]() Comparer[E] {
	return Natural[E]().Reverse()
}

// DerefNilFirst returns a Comparer that compares two values of the same type *E by dereferencing them and comparing the resulting values using the provided Comparer.
// If the first value is nil, it is considered less than the second value.
// If the second value is nil, it is considered greater than the first value.
// If both values are nil, they are considered equal.
// Otherwise, the provided Comparer is used to compare the dereferenced values.
func DerefNilFirst[E any](cmp Comparer[E]) Comparer[*E] {
	return func(a, b *E) int {
		if a == nil {
			if b == nil {
				return 0
			}
			return -1
		}
		if b == nil {
			return 1
		}
		return cmp(*a, *b)
	}
}

// DerefNilLast returns a Comparer that compares two values of the same type *E by dereferencing them and comparing the resulting values using the provided Comparer.
// If the first value is nil, it is considered greater than the second value.
// If the second value is nil, it is considered less than the first value.
// If both values are nil, they are considered equal.
// Otherwise, the provided Comparer is used to compare the dereferenced values.
func DerefNilLast[E any](cmp Comparer[E]) Comparer[*E] {
	return func(a, b *E) int {
		if a == nil {
			if b == nil {
				return 0
			}
			return 1
		}
		if b == nil {
			return -1
		}
		return cmp(*a, *b)
	}
}

// Comparing returns a Comparer that compares two values of the same type E by extracting a sort key of type K from each value using the provided KeyExtractor,
// then comparing the resulting keys using the standard library's Compare function for type K.
// The type parameter K must implement the Ordered constraint.
//
// Example:
//
//	type Person struct {
//		FirstName string
//		LastName  string
//	}
//
//	people := []Person{
//		{"John", "Doe"},
//		{"Jane", "Doe"},
//		{"John", "Smith"},
//	}
//
//	// Sort by LastName, then by FirstName.
//	s := stream.SortBy(
//		stream.FromSlice(people),
//		Comparing(func(p Person) string { return p.LastName }),
//
//	)
//	out := stream.DebugString(s) // <Person{FirstName:"Jane", LastName:"Doe"}, Person{FirstName:"John", LastName:"Doe"}, Person{FirstName:"John", LastName:"Smith"}>
func Comparing[E any, K constraint.Ordered](ke KeyExtractor[E, K]) Comparer[E] {
	return func(a, b E) int {
		return stdcmp.Compare[K](ke(a), ke(b))
	}
}

// ComparingBy returns a Comparer that compares two values of the same type E by extracting a sort key of type K from each value using the provided KeyExtractor,
// then comparing the resulting keys using the provided Comparer.
//
// Example:
//
//	type Person struct {
//		FirstName string
//		LastName  string
//	}
//
//	people := []Person{
//		{"John", "Doe"},
//		{"Jane", "Doe"},
//		{"John", "Smith"},
//	}
//
//	// Sort by LastName, then by FirstName.
//	s := stream.SortBy(
//		stream.FromSlice(people),
//		ComparingBy(func(p Person) string { return p.LastName }, Natural[string]()),
//	)
//	out := stream.DebugString(s) // <Person{FirstName:"Jane", LastName:"Doe"}, Person{FirstName:"John", LastName:"Doe"}, Person{FirstName:"John", LastName:"Smith"}>
func ComparingBy[E any, K any](ke KeyExtractor[E, K], cmp Comparer[K]) Comparer[E] {
	return func(a, b E) int {
		return cmp(ke(a), ke(b))
	}
}
