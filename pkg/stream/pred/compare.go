package pred

import (
	"math"
	"reflect"

	"github.com/jpfourny/papaya/v2/pkg/cmp"
	"github.com/jpfourny/papaya/v2/pkg/constraint"
)

// Equal returns a function that returns true if the provided value is equal to the provided want value.
// It uses the == operator to compare values.
//
// Examples:
//
//	p := pred.Equal(0)
//	p(0) // true
//	p(1) // false
func Equal[E comparable](want E) func(E) bool {
	return func(got E) bool {
		return got == want
	}
}

// NotEqual returns a function that returns true if the provided value is not equal to the provided want value.
// It uses the != operator to compare values.
//
// Examples:
//
//	p := pred.NotEqual(0)
//	p(0) // false
//	p(1) // true
func NotEqual[E comparable](want E) func(E) bool {
	return func(got E) bool {
		return got != want
	}
}

// EqualBy returns a function that returns true if the provided value is equal to the provided want value.
// It uses the provided cmp.Comparer to compare values.
//
// Examples:
//
//	p := pred.EqualBy(0, cmp.Natural[int])
//	p(0) // true
//	p(1) // false
func EqualBy[E any](want E, compare cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		return compare.Equal(got, want)
	}
}

// NotEqualBy returns a function that returns true if the provided value is not equal to the provided want value.
// It uses the provided cmp.Comparer to compare values.
//
// Examples:
//
//	p := pred.NotEqualBy(0, cmp.Natural[int])
//	p(0) // false
//	p(1) // true
func NotEqualBy[E any](want E, compare cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		return compare.NotEqual(got, want)
	}
}

// DeepEqual returns a function that returns true if the provided value is equal to the provided want value.
// It uses the reflect.DeepEqual function to compare values.
//
// Examples:
//
//	p := pred.DeepEqual(0)
//	p(0) // true
//	p(1) // false
func DeepEqual[E any](want E) func(E) bool {
	return func(got E) bool {
		return reflect.DeepEqual(got, want)
	}
}

// NotDeepEqual returns a function that returns true if the provided value is not equal to the provided want value.
// It uses the reflect.DeepEqual function to compare values.
//
// Examples:
//
//	p := pred.NotDeepEqual(0)
//	p(0) // false
//	p(1) // true
func NotDeepEqual[E any](want E) func(E) bool {
	return func(got E) bool {
		return !reflect.DeepEqual(got, want)
	}
}

// LessThan returns a function that returns true if the provided value is less than the provided want value.
// It uses the < operator to compare values.
//
// Examples:
//
//	p := pred.LessThan(0)
//	p(-1) // true
//	p(0) // false
//	p(1) // false
func LessThan[E constraint.Ordered](want E) func(E) bool {
	return func(got E) bool {
		return got < want
	}
}

// LessThanOrEqual returns a function that returns true if the provided value is less than or equal to the provided want value.
// It uses the <= operator to compare values.
//
// Examples:
//
//	p := pred.LessThanOrEqual(0)
//	p(-1) // true
//	p(0) // true
//	p(1) // false
func LessThanOrEqual[E constraint.Ordered](want E) func(E) bool {
	return func(got E) bool {
		return got <= want
	}
}

// GreaterThan returns a function that returns true if the provided value is greater than the provided want value.
// It uses the > operator to compare values.
//
// Examples:
//
//	p := pred.GreaterThan(0)
//	p(-1) // false
//	p(0) // false
//	p(1) // true
func GreaterThan[E constraint.Ordered](want E) func(E) bool {
	return func(got E) bool {
		return got > want
	}
}

// GreaterThanOrEqual returns a function that returns true if the provided value is greater than or equal to the provided want value.
// It uses the >= operator to compare values.
//
// Examples:
//
//	p := pred.GreaterThanOrEqual(0)
//	p(-1) // false
//	p(0) // true
//	p(1) // true
func GreaterThanOrEqual[E constraint.Ordered](want E) func(E) bool {
	return func(got E) bool {
		return got >= want
	}
}

// LessThanBy returns a function that returns true if the provided value is less than the provided want value.
// It uses the provided cmp.Comparer to compare values.
//
// Examples:
//
//	p := pred.LessThanBy(0, cmp.Natural[int])
//	p(-1) // true
//	p(0) // false
//	p(1) // false
func LessThanBy[E any](want E, compare cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		return compare.LessThan(got, want)
	}
}

// LessThanOrEqualBy returns a function that returns true if the provided value is less than or equal to the provided want value.
// It uses the provided cmp.Comparer to compare values.
//
// Examples:
//
//	p := pred.LessThanOrEqualBy(0, cmp.Natural[int])
//	p(-1) // true
//	p(0) // true
//	p(1) // false
func LessThanOrEqualBy[E any](want E, compare cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		return compare.LessThanOrEqual(got, want)
	}
}

// GreaterThanBy returns a function that returns true if the provided value is greater than the provided want value.
// It uses the provided cmp.Comparer to compare values.
//
// Examples:
//
//	p := pred.GreaterThanBy(0, cmp.Natural[int])
//	p(-1) // false
//	p(0) // false
//	p(1) // true
func GreaterThanBy[E any](want E, compare cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		return compare.GreaterThan(got, want)
	}
}

// GreaterThanOrEqualBy returns a function that returns true if the provided value is greater than or equal to the provided want value.
// It uses the provided cmp.Comparer to compare values.
//
// Examples:
//
//	p := pred.GreaterThanOrEqualBy(0, cmp.Natural[int])
//	p(-1) // false
//	p(0) // true
//	p(1) // true
func GreaterThanOrEqualBy[E any](want E, compare cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		return compare.GreaterThanOrEqual(got, want)
	}
}

// RoughlyEqual returns a function specialized for floats that returns true if the provided value is roughly equal to the provided want value relative to the provided epsilon.
//
// Note: this function is symmetric, but less meaningful for values near or smaller than epsilon.
func RoughlyEqual[E constraint.Float](want, epsilon E) func(E) bool {
	return func(got E) bool {
		if got == want { // shortcut, handles infinities
			return true
		} else { // use absolute error
			return math.Abs(float64(got-want)) < float64(epsilon)
		}
	}
}

// NotRoughlyEqual returns a function specialized for floats that returns true if the provided value is not roughly equal to the provided want value relative to the provided epsilon.
//
// Note: this function is symmetric, but less meaningful for values near or smaller than epsilon.
func NotRoughlyEqual[E constraint.Float](want, epsilon E) func(E) bool {
	return Not(RoughlyEqual(want, epsilon))
}
