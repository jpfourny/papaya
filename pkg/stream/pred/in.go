package pred

import "github.com/jpfourny/papaya/v2/pkg/cmp"

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
func InBy[E any](compare cmp.Comparer[E], want ...E) func(E) bool {
	return InSliceBy(want, compare)
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
func NotInBy[E any](compare cmp.Comparer[E], want ...E) func(E) bool {
	return NotInSliceBy[E](want, compare)
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
func InSliceBy[E any](want []E, compare cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		for _, w := range want {
			if compare.Equal(got, w) {
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
func NotInSliceBy[E any](want []E, compare cmp.Comparer[E]) func(E) bool {
	return func(got E) bool {
		for _, w := range want {
			if compare.Equal(got, w) {
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
