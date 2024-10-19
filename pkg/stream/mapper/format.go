package mapper

import (
	"fmt"
	"strconv"

	"github.com/jpfourny/papaya/v2/pkg/constraint"
)

// Sprint returns a function that accepts a value of any type and returns the result of calling fmt.Sprint on it.
// See the documentation for fmt.Sprint for details.
func Sprint[E any]() func(E) string {
	return func(e E) string {
		return fmt.Sprint(e)
	}
}

// Sprintf returns a function that accepts a value of any type and returns the result of calling fmt.Sprintf on it.
// The `format` string is passed as a parameter to fmt.Sprintf.
// See the documentation for fmt.Sprintf for details.
func Sprintf[E any](format string) func(E) string {
	return func(e E) string {
		return fmt.Sprintf(format, e)
	}
}

// FormatBool returns a function that accepts a value of any boolean type and returns the result of calling strconv.FormatBool on it.
// See the documentation for strconv.FormatBool for details.
func FormatBool[E constraint.Boolean]() func(E) string {
	return func(e E) string {
		return strconv.FormatBool(bool(e))
	}
}

// FormatInt returns a function that accepts a value of any integer type and returns the result of calling strconv.FormatInt on it.
// The `base` parameter is passed as a parameter to strconv.FormatInt.
// See the documentation for strconv.FormatInt for details.
func FormatInt[E constraint.Integer](base int) func(E) string {
	return func(e E) string {
		return strconv.FormatInt(int64(e), base)
	}
}

// FormatUint returns a function that accepts a value of any unsigned integer type and returns the result of calling strconv.FormatUint on it.
// The `base` parameter is passed as a parameter to strconv.FormatUint.
// See the documentation for strconv.FormatUint for details.
func FormatUint[E constraint.UnsignedInteger](base int) func(E) string {
	return func(e E) string {
		return strconv.FormatUint(uint64(e), base)
	}
}

// FormatFloat returns a function that accepts a value of any floating-point type and returns the result of calling strconv.FormatFloat on it.
// The `format`, `prec` and `bitSize` parameters are passed as parameters to strconv.FormatFloat.
// See the documentation for strconv.FormatFloat for details.
func FormatFloat[E constraint.Float](format byte, prec int, bitSize int) func(E) string {
	return func(e E) string {
		return strconv.FormatFloat(float64(e), format, prec, bitSize)
	}
}

// FormatComplex returns a function that accepts a value of any complex numeric type and returns the result of calling strconv.FormatComplex on it.
// The `format` and `prec` parameters are passed as parameters to strconv.FormatComplex.
// See the documentation for strconv.FormatComplex for details.
func FormatComplex[E constraint.Complex](format byte, prec int) func(E) string {
	return func(e E) string {
		return strconv.FormatComplex(complex128(e), format, prec, 128)
	}
}
