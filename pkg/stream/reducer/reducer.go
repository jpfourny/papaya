package reducer

import (
	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/constraint"
)

// Sum returns a function that computes the sum of two given values.
func Sum[E constraint.Numeric]() func(a, b E) E {
	return func(a, b E) E {
		return a + b
	}
}

// Min returns a function that returns the minimum of two given values.
func Min[E constraint.Ordered]() func(a, b E) E {
	return func(a, b E) E {
		if a <= b {
			return a
		}
		return b
	}
}

// MinBy returns a function that returns the minimum of two given values using the provided cmp.Comparer.
func MinBy[E any](compare cmp.Comparer[E]) func(a, b E) E {
	return func(a, b E) E {
		if compare.LessThanOrEqual(a, b) {
			return a
		}
		return b
	}
}

// Max returns a function that returns the maximum of two given values.
func Max[E constraint.Ordered]() func(a, b E) E {
	return func(a, b E) E {
		if a >= b {
			return a
		}
		return b
	}
}

// MaxBy returns a function that returns the maximum of two given values using the provided cmp.Comparer.
func MaxBy[E any](compare cmp.Comparer[E]) func(a, b E) E {
	return func(a, b E) E {
		if compare.GreaterThanOrEqual(a, b) {
			return a
		}
		return b
	}
}
