package pred

import (
	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/constraint"
	"math"
	"reflect"
)

// True returns a function that always returns true.
func True[E any]() func(E) bool {
	return func(E) bool {
		return true
	}
}

// False returns a function that always returns false.
func False[E any]() func(E) bool {
	return func(E) bool {
		return false
	}
}

// Not returns a function that returns the opposite of the provided predicate.
//
// Examples:
//
//	t := pred.True[any]()
//	f := pred.False[any]()
//	pred.Not(t) // evaluates to false
//	pred.Not(f) // evaluates to true
func Not[E any](p func(E) bool) func(E) bool {
	return func(e E) bool {
		return !p(e)
	}
}

// And returns a function that returns true if both the provided predicates return true.
//
// Examples:
//
//	t := pred.True[any]()
//	f := pred.False[any]()
//	pred.And(t, t) // evaluates to true
//	pred.And(t, f) // evaluates to false
//	pred.And(f, f) // evaluates to false
func And[E any](p1, p2 func(E) bool) func(E) bool {
	return func(e E) bool {
		return p1(e) && p2(e)
	}
}

// Or returns a function that returns true if either the provided predicates return true.
//
// Examples:
//
//	t := pred.True[any]()
//	f := pred.False[any]()
//	pred.Or(t, t) // evaluates to true
//	pred.Or(t, f) // evaluates to true
//	pred.Or(f, f) // evaluates to false
func Or[E any](p1, p2 func(E) bool) func(E) bool {
	return func(e E) bool {
		return p1(e) || p2(e)
	}
}

// OneOf returns a function that returns true if exactly one of the provided predicates returns true.
// It returns false when:
//   - none of the provided predicates return true, or
//   - more than one of the provided predicates return true or
//   - no predicates are provided (degenerate case).
//
// Examples:
//
//	t := pred.True[any]()
//	f := pred.False[any]()
//	pred.OneOf() // false (no predicates provided)
//	pred.OneOf(t, t) // false (both predicates match)
//	pred.OneOf(t, f) // true (one predicate matches)
//	pred.OneOf(f, f) // false (no predicates match)
func OneOf[E any](ps ...func(E) bool) func(E) bool {
	return func(e E) bool {
		match := false
		for _, p := range ps {
			if p(e) {
				if match {
					return false // More than one match.
				}
				match = true // First match.
			}
		}
		return match // False if no predicates are provided.
	}
}

// AllOf returns a function that returns true if all the provided predicates return true.
// It returns false when:
//   - any of the provided predicates return false, or
//   - no predicates are provided.
//
// Examples:
//
//	t := pred.True[any]()
//	f := pred.False[any]()
//	pred.AllOf() // false (no predicates provided)
//	pred.AllOf(t, t) // true (both predicates match)
//	pred.AllOf(t, f) // false (one predicate matches)
//	pred.AllOf(f, f) // false (no predicates match)
func AllOf[E any](ps ...func(E) bool) func(E) bool {
	return func(e E) bool {
		for _, p := range ps {
			if !p(e) {
				return false
			}
		}
		return len(ps) > 0 // False if no predicates are provided.
	}
}

// AnyOf returns a function that returns true if any of the provided predicates return true.
// It returns false when:
//   - none of the provided predicates return true, or
//   - no predicates are provided.
//
// Examples:
//
//	t := pred.True[any]()
//	f := pred.False[any]()
//	pred.AnyOf() // false (no predicates provided)
//	pred.AnyOf(t, t) // true (both predicates match)
//	pred.AnyOf(t, f) // true (one predicate matches)
//	pred.AnyOf(f, f) // false (no predicates match)
func AnyOf[E any](ps ...func(E) bool) func(E) bool {
	return func(e E) bool {
		for _, p := range ps {
			if p(e) {
				return true
			}
		}
		return false
	}
}

// NoneOf returns a function that returns true if none of the provided predicates return true.
// It returns false when any of the provided predicates return true.
// If no predicates are provided, it returns true (degenerate case).
//
// Examples:
//
//	t := pred.True[any]()
//	f := pred.False[any]()
//	pred.NoneOf() // true (no predicates provided)
//	pred.NoneOf(t, t) // false (both predicates match)
//	pred.NoneOf(t, f) // false (one predicate matches)
//	pred.NoneOf(f, f) // true (no predicates match)
func NoneOf[E any](ps ...func(E) bool) func(E) bool {
	return func(e E) bool {
		for _, p := range ps {
			if p(e) {
				return false
			}
		}
		return true
	}
}

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
func EqualBy[E any](want E, cmp cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		return cmp(got, want) == 0
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
func NotEqualBy[E any](want E, cmp cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		return cmp(got, want) != 0
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
func LessThanBy[E any](want E, cmp cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		return cmp(got, want) < 0
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
func LessThanOrEqualBy[E any](want E, cmp cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		return cmp(got, want) <= 0
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
func GreaterThanBy[E any](want E, cmp cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		return cmp(got, want) > 0
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
func GreaterThanOrEqualBy[E any](want E, cmp cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		return cmp(got, want) >= 0
	}
}

// RoughlyEqual returns a function specialized for floats that returns true if the provided value is roughly equal to the provided want value relative to the provided epsilon.
//
// Note: this function is symmetric, but less meaningful for values near or smaller than epsilon.
func RoughlyEqual[E constraint.Float](want, epsilon E) func(E) bool {
	return func(got E) bool {
		return roughlyEqual(float64(got), float64(want), float64(epsilon))
	}
}

// NotRoughlyEqual returns a function specialized for floats that returns true if the provided value is not roughly equal to the provided want value relative to the provided epsilon.
//
// Note: this function is symmetric, but less meaningful for values near or smaller than epsilon.
func NotRoughlyEqual[E constraint.Float](want, epsilon E) func(E) bool {
	return Not(RoughlyEqual(want, epsilon))
}

// roughlyEqual returns true if the provided values are roughly equal relative to the provided epsilon.
//
// Note: this function is symmetric, but less meaningful for values near or smaller than epsilon.
func roughlyEqual(a, b, epsilon float64) bool {
	if a == b {
		// shortcut, handles infinities
		return true
	} else { // use absolute error
		return math.Abs(a-b) < epsilon
	}
}

// Nil returns a function that returns true if the provided pointer is nil.
//
// Examples:
//
//	p := pred.Nil[int]()
//	p(nil) // true
//	p(0) // false
func Nil[E any]() func(*E) bool {
	return func(e *E) bool {
		return e == nil
	}
}

// NotNil returns a function that returns true if the provided pointer is not nil.
//
// Examples:
//
//	p := pred.NotNil[int]()
//	p(nil) // false
//	p(0) // true
func NotNil[E any]() func(*E) bool {
	return func(e *E) bool {
		return e != nil
	}
}

// Zero returns a function that returns true if the provided value is the zero value of the type parameter E.
//
// Examples:
//
//	p := pred.Zero[int]()
//	p(0) // true
//	p(1) // false
func Zero[E comparable]() func(E) bool {
	var zero E
	return func(e E) bool {
		return e == zero
	}
}

// NotZero returns a function that returns true if the provided value is not the zero value of the type parameter E.
//
// Examples:
//
//	p := pred.NotZero[int]()
//	p(0) // false
//	p(1) // true
func NotZero[E comparable]() func(E) bool {
	var zero E
	return func(e E) bool {
		return e != zero
	}
}

// In returns a function that returns true if the provided value is equal to any of the provided want values.
// It uses the == operator to compare values.
//
// Examples:
//
//	p := pred.In(0, 1, 2)
//	p(0) // true
//	p(1) // true
//	p(2) // true
//	p(3) // false
func In[E comparable](want ...E) func(E) bool {
	return InSlice(want)
}

// NotIn returns a function that returns true if the provided value is not equal to any of the provided want values.
// It uses the == operator to compare values.
//
// Examples:
//
//	p := pred.NotIn(0, 1, 2)
//	p(0) // false
//	p(1) // false
//	p(2) // false
//	p(3) // true
func NotIn[E comparable](e ...E) func(E) bool {
	return NotInSlice(e)
}

// InSlice returns a function that returns true if the provided value is equal to any element in the provided want slice.
// It uses the == operator to compare values.
//
// Examples:
//
//	p := pred.InSlice([]int{0, 1, 2})
//	p(0) // true
//	p(1) // true
//	p(2) // true
//	p(3) // false
func InSlice[E comparable](want []E) func(E) bool {
	return func(got E) bool {
		for _, w := range want {
			if got == w {
				return true
			}
		}
		return false
	}
}

// NotInSlice returns a function that returns true if the provided value is not equal to any element in the provided want slice.
// It uses the == operator to compare values.
//
// Examples:
//
//	p := pred.NotInSlice([]int{0, 1, 2})
//	p(0) // false
//	p(1) // false
//	p(2) // false
//	p(3) // true
func NotInSlice[E comparable](want []E) func(E) bool {
	return func(got E) bool {
		for _, w := range want {
			if got == w {
				return false
			}
		}
		return true
	}
}

// InBy returns a function that returns true if the provided value is equal to any of the provided want values.
// It uses the provided cmp.Comparer to compare values.
//
// Examples:
//
//	p := pred.IntBy(cmp.Natural[int], 0, 1, 2)
//	p(0) // true
//	p(1) // true
//	p(2) // true
//	p(3) // false
func InBy[E any](cmp cmp.Comparer[E], want ...E) func(E) bool {
	return InSliceBy(want, cmp)
}

// NotInBy returns a function that returns true if the provided value is not equal to any of the provided want values.
// It uses the provided cmp.Comparer to compare values.
//
// Examples:
//
//	p := pred.NotInBy(cmp.Natural[int], 0, 1, 2)
//	p(0) // false
//	p(1) // false
//	p(2) // false
//	p(3) // true
func NotInBy[E any](cmp cmp.Comparer[E], want ...E) func(E) bool {
	return NotInSliceBy[E](want, cmp)
}

// InSliceBy returns a function that returns true if the provided value is equal to any element in the provided want slice.
// It uses the provided cmp.Comparer to compare values.
//
// Examples:
//
//	p := pred.InSliceBy([]int{0, 1, 2}, cmp.Natural[int])
//	p(0) // true
//	p(1) // true
//	p(2) // true
//	p(3) // false
func InSliceBy[E any](want []E, cmp cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		for _, w := range want {
			if cmp(got, w) == 0 {
				return true
			}
		}
		return false
	}
}

// NotInSliceBy returns a function that returns true if the provided value is not equal to any element in the provided want slice.
// It uses the provided cmp.Comparer to compare values.
//
// Examples:
//
//	p := pred.NotInSliceBy([]int{0, 1, 2}, cmp.Natural[int])
//	p(0) // false
//	p(1) // false
//	p(2) // false
//	p(3) // true
func NotInSliceBy[E any](want []E, cmp cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		for _, w := range want {
			if cmp(got, w) == 0 {
				return false
			}
		}
		return true
	}
}

// InSet returns a function that returns true if the provided value is found in the provided want set.
// The set is represented with map keys; the map values are ignored.
//
// Examples:
//
//	set := map[int]struct{}{0: {}, 1: {}, 2: {}}
//	p := pred.InSet(set)
//	p(0) // true
//	p(1) // true
//	p(2) // true
//	p(3) // false
func InSet[E comparable, F any](want map[E]F) func(E) bool {
	return func(got E) bool {
		_, ok := want[got]
		return ok
	}
}

// NotInSet returns a function that returns true if the provided value is not found in the provided want set.
// The set is represented with map keys; the map values are ignored.
//
// Examples:
//
//	set := map[int]struct{}{0: {}, 1: {}, 2: {}}
//	p := pred.NotInSet(set)
//	p(0) // false
//	p(1) // false
//	p(2) // false
//	p(3) // true
func NotInSet[E comparable, F any](want map[E]F) func(E) bool {
	return func(got E) bool {
		_, ok := want[got]
		return !ok
	}
}
