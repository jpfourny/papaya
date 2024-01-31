package cmp

import (
	stdcmp "cmp"
	"time"

	"github.com/jpfourny/papaya/pkg/constraint"
	"github.com/jpfourny/papaya/pkg/pair"
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
//	// Sort by LastName ascending, then by FirstName descending.
//	sort.Slice(
//		people,
//		cmp.Comparing(func(p Person) string { return p.LastName }).
//			Then(cmp.Comparing(func(p Person) string { return p.FirstName }).Reverse()),
//	) // [Person{FirstName:"John", LastName:"Doe"}, Person{FirstName:"Jane", LastName:"Doe"}, Person{FirstName:"John", LastName:"Smith"}]
func (c Comparer[E]) Then(other Comparer[E]) Comparer[E] {
	return func(a, b E) int {
		if r := c(a, b); r != 0 {
			return r
		}
		return other(a, b)
	}
}

// Min returns the minimum of two values of the same type E using the provided Comparer.
func (c Comparer[E]) Min(a, b E) E {
	if c(a, b) <= 0 {
		return a
	}
	return b
}

// Max returns the maximum of two values of the same type E using the provided Comparer.
func (c Comparer[E]) Max(a, b E) E {
	if c(a, b) >= 0 {
		return a
	}
	return b
}

// Equal returns true if the two values of the same type E are equal using the provided Comparer; otherwise it returns false.
func (c Comparer[E]) Equal(a, b E) bool {
	return c(a, b) == 0
}

// NotEqual returns true if the two values of the same type E are not equal using the provided Comparer; otherwise it returns false.
func (c Comparer[E]) NotEqual(a, b E) bool {
	return c(a, b) != 0
}

// LessThan returns true if the first value of type E is less than the second value of type E using the provided Comparer; otherwise it returns false.
func (c Comparer[E]) LessThan(a, b E) bool {
	return c(a, b) < 0
}

// LessThanOrEqual returns true if the first value of type E is less than or equal to the second value of type E using the provided Comparer; otherwise it returns false.
func (c Comparer[E]) LessThanOrEqual(a, b E) bool {
	return c(a, b) <= 0
}

// GreaterThan returns true if the first value of type E is greater than the second value of type E using the provided Comparer; otherwise it returns false.
func (c Comparer[E]) GreaterThan(a, b E) bool {
	return c(a, b) > 0
}

// GreaterThanOrEqual returns true if the first value of type E is greater than or equal to the second value of type E using the provided Comparer; otherwise it returns false.
func (c Comparer[E]) GreaterThanOrEqual(a, b E) bool {
	return c(a, b) >= 0
}

// KeyExtractor is a function that extracts a sort key of type K from a value of type E.
type KeyExtractor[E, K any] func(E) K

// SelfComparer is an interface that represents a type that can compare itself to another value of the same type.
type SelfComparer[E any] interface {
	Compare(E) int
}

// Natural returns a Comparer that compares two values of the same ordered type E using <, ==, > operators.
//
// Example:
//
//	s := []int{3, 1, 2}
//	sort.Slice(s, cmp.Natural[int]()) // [1, 2, 3]
func Natural[E constraint.Ordered]() Comparer[E] {
	return stdcmp.Compare[E]
}

// Reverse returns a Comparer that compares two values of the same ordered type E in reverse order using <, ==, > operators.
//
// Example:
//
//	s := []int{3, 1, 2}
//	sort.Slice(s, cmp.Reverse[int]()) // [3, 2, 1]
func Reverse[E constraint.Ordered]() Comparer[E] {
	return Natural[E]().Reverse()
}

// Self returns a Comparer that compares two values of the same type E by calling the Compare method on the first value with the second value as the argument.
// The type parameter E must implement the SelfComparer interface.
func Self[E SelfComparer[E]]() Comparer[E] {
	return func(a, b E) int {
		return a.Compare(b)
	}
}

// Bool returns a Comparer that compares two values of type bool.
// The zero value of type bool is considered less than false, which is considered less than true.
func Bool() Comparer[bool] {
	return func(a, b bool) int {
		if a == b {
			return 0
		}
		if a {
			return 1
		}
		return -1
	}
}

// Time returns a Comparer that compares two values of type time.Time.
func Time() Comparer[time.Time] {
	return func(a, b time.Time) int {
		if a.Before(b) {
			return -1
		} else if a.After(b) {
			return 1
		} else {
			return 0
		}
	}
}

// Complex64 returns a Comparer that compares two values of type complex64.
// The real and imaginary parts of the complex numbers are compared in order.
// If the real parts are equal, the imaginary parts are compared.
// If both the real and imaginary parts are equal, the complex numbers are considered equal.
func Complex64() Comparer[complex64] {
	return func(a, b complex64) int {
		aReal, aImag := real(a), imag(a)
		bReal, bImag := real(b), imag(b)

		if aReal < bReal {
			return -1
		} else if aReal > bReal {
			return 1
		} else if aImag < bImag {
			return -1
		} else if aImag > bImag {
			return 1
		} else {
			return 0
		}
	}
}

// Complex128 returns a Comparer that compares two values of type complex128.
// The real and imaginary parts of the complex numbers are compared in order.
// If the real parts are equal, the imaginary parts are compared.
// If both the real and imaginary parts are equal, the complex numbers are considered equal.
func Complex128() Comparer[complex128] {
	return func(a, b complex128) int {
		aReal, aImag := real(a), imag(a)
		bReal, bImag := real(b), imag(b)

		if aReal < bReal {
			return -1
		} else if aReal > bReal {
			return 1
		} else if aImag < bImag {
			return -1
		} else if aImag > bImag {
			return 1
		} else {
			return 0
		}
	}
}

// Pair returns a Comparer that compares two values of type pair.Pair[A, B] by comparing the elements of the pairs using the provided Comparer functions.
// If the first elements are equal, the second elements of the pairs are compared.
// If both the first and second elements are equal, the pairs are considered equal.
//
// Example:
//
//	s := []pair.Pair[int, int]{{3, 1}, {1, 2}, {1, 1}, {1, 3}}
//	sort.Slice(s, cmp.Pair(cmp.Natural[int](), cmp.Natural[int]())) // [[1, 1], [1, 2], [1, 3], [3, 1]]
func Pair[A, B any](compareA Comparer[A], compareB Comparer[B]) Comparer[pair.Pair[A, B]] {
	return func(a, b pair.Pair[A, B]) int {
		if c := compareA(a.First(), b.First()); c != 0 {
			return c
		}
		return compareB(a.Second(), b.Second())
	}
}

// Slice returns a Comparer that compares two slices of type []E by comparing the elements of the slices using the provided Comparer.
// The first N elements of the slices are compared, where N is the length of the shorter slice.
// If all N elements are equal, the length of the slices are compared.
//
// Example:
//
//	s := [][]int{{3, 1, 2}, {1, 2, 3}, {1, 2}, {1, 2, 3, 4}}
//	sort.Slice(s, cmp.Slice(cmp.Natural[int]())) // [[1, 2], [1, 2, 3], [1, 2, 3, 4], [3, 1, 2]]
func Slice[E any](compare Comparer[E]) Comparer[[]E] {
	return func(a, b []E) int {
		// Sort by common elements.
		n := min(len(a), len(b))
		for i := 0; i < n; i++ {
			if c := compare(a[i], b[i]); c != 0 {
				return c
			}
		}
		// Sort by length.
		if len(a) < len(b) {
			return -1
		} else if len(a) > len(b) {
			return 1
		}
		return 0
	}
}

// DerefNilFirst returns a Comparer that compares two values of the same type *E by dereferencing them and comparing the resulting values using the provided Comparer.
// If the first value is nil, it is considered less than the second value.
// If the second value is nil, it is considered greater than the first value.
// If both values are nil, they are considered equal.
// Otherwise, the provided Comparer is used to compare the dereferenced values.
//
// Example:
//
//	s := []*int{ptr.Ref(3), nil, ptr.Ref(1), ptr.Ref(2)}
//	sort.Slice(s, cmp.DerefNilFirst(cmp.Natural[int]())) // [nil, 1, 2, 3]
func DerefNilFirst[E any](compare Comparer[E]) Comparer[*E] {
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
		return compare(*a, *b)
	}
}

// DerefNilLast returns a Comparer that compares two values of the same type *E by dereferencing them and comparing the resulting values using the provided Comparer.
// If the first value is nil, it is considered greater than the second value.
// If the second value is nil, it is considered less than the first value.
// If both values are nil, they are considered equal.
// Otherwise, the provided Comparer is used to compare the dereferenced values.
//
// Example:
//
//	s := []*int{ptr.Ref(3), nil, ptr.Ref(1), ptr.Ref(2)}
//	sort.Slice(s, cmp.DerefNilLast(cmp.Natural[int]())) // [1, 2, 3, nil]
func DerefNilLast[E any](compare Comparer[E]) Comparer[*E] {
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
		return compare(*a, *b)
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
//	// Sort by LastName.
//	sort.Slice(
//		people,
//		cmp.Comparing(func(p Person) string { return p.LastName }),
//	) // [Person{FirstName:"Jane", LastName:"Doe"}, Person{FirstName:"John", LastName:"Doe"}, Person{FirstName:"John", LastName:"Smith"}]
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
//	// Sort by LastName.
//	sort.Slice(
//		people,
//		cmp.ComparingBy(func(p Person) string { return p.LastName }, cmp.Natural[string]()),
//	) // [Person{FirstName:"Jane", LastName:"Doe"}, Person{FirstName:"John", LastName:"Doe"}, Person{FirstName:"John", LastName:"Smith"}]
func ComparingBy[E any, K any](ke KeyExtractor[E, K], compare Comparer[K]) Comparer[E] {
	return func(a, b E) int {
		return compare(ke(a), ke(b))
	}
}
