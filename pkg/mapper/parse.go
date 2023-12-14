package mapper

import (
	"strconv"

	"github.com/jpfourny/papaya/pkg/constraint"
)

// ParseBoolOr returns a function that accepts a value of any string type and returns the result of calling strconv.ParseBool on it.
// If the string cannot be parsed as a boolean, the provided `or` default value is returned.
// See the documentation for strconv.ParseBool for details.
func ParseBoolOr[E constraint.String](or bool) func(E) bool {
	return func(e E) bool {
		b, err := strconv.ParseBool(string(e))
		if err != nil {
			return or
		}
		return b
	}
}

// ParseIntOr returns a function that accepts a value of any string type and returns the result of calling strconv.ParseInt on it.
// If the string cannot be parsed as an integer, the provided `or` default value is returned.
// The `base` and `bitSize` parameters are passed as parameters to strconv.ParseInt.
// See the documentation for strconv.ParseInt for details.
func ParseIntOr[E constraint.String, F constraint.SignedInteger](base int, bitSize int, or F) func(E) F {
	return func(e E) F {
		i, err := strconv.ParseInt(string(e), base, bitSize)
		if err != nil {
			return or
		}
		return F(i)
	}
}

// ParseUintOr returns a function that accepts a value of any string type and returns the result of calling strconv.ParseUint on it.
// If the string cannot be parsed as an unsigned integer, the provided `or` default value is returned.
// The `base` and `bitSize` parameters are passed as parameters to strconv.ParseUint.
// See the documentation for strconv.ParseUint for details.
func ParseUintOr[E constraint.String, F constraint.UnsignedInteger](base int, bitSize int, or F) func(E) F {
	return func(e E) F {
		i, err := strconv.ParseUint(string(e), base, bitSize)
		if err != nil {
			return or
		}
		return F(i)
	}
}

// ParseFloatOr returns a function that accepts a value of any string type and returns the result of calling strconv.ParseFloat on it.
// If the string cannot be parsed as a floating-point number, the provided `or` default value is returned.
// The `bitSize` parameter is passed as a parameter to strconv.ParseFloat.
// See the documentation for strconv.ParseFloat for details.
func ParseFloatOr[E constraint.String, F constraint.Float](bitSize int, or F) func(E) F {
	return func(e E) F {
		f, err := strconv.ParseFloat(string(e), bitSize)
		if err != nil {
			return or
		}
		return F(f)
	}
}

// ParseComplexOr returns a function that accepts a value of any string type and returns the result of calling strconv.ParseComplex on it.
// If the string cannot be parsed as a complex number, the provided `or` default value is returned.
// The `bitSize` parameter is passed as a parameter to strconv.ParseComplex.
// See the documentation for strconv.ParseComplex for details.
func ParseComplexOr[E constraint.String, F constraint.Complex](bitSize int, or F) func(E) F {
	return func(e E) F {
		c, err := strconv.ParseComplex(string(e), bitSize)
		if err != nil {
			return or
		}
		return F(c)
	}
}
