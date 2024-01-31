package mapper

import (
	"strconv"

	"github.com/jpfourny/papaya/pkg/constraint"
	"github.com/jpfourny/papaya/pkg/opt"
)

// TryParseBool returns a function that accepts a value of any string type and returns the result of calling strconv.ParseBool on it.
// If the string cannot be parsed as a boolean, then an empty opt is returned.
// See the documentation for strconv.ParseBool for details.
func TryParseBool[E constraint.String]() func(E) opt.Optional[bool] {
	return func(e E) opt.Optional[bool] {
		b, err := strconv.ParseBool(string(e))
		if err != nil {
			return opt.Empty[bool]()
		}
		return opt.Of[bool](b)
	}
}

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

// TryParseInt returns a function that accepts a value of any string type and returns the result of calling strconv.ParseInt on it.
// If the string cannot be parsed as an integer, then an empty opt is returned.
// The `base` and `bitSize` parameters are passed as parameters to strconv.ParseInt.
// See the documentation for strconv.ParseInt for details.
func TryParseInt[E constraint.String, F constraint.SignedInteger](base int, bitSize int) func(E) opt.Optional[F] {
	return func(e E) opt.Optional[F] {
		i, err := strconv.ParseInt(string(e), base, bitSize)
		if err != nil {
			return opt.Empty[F]()
		}
		return opt.Of(F(i))
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

// TryParseUint returns a function that accepts a value of any string type and returns the result of calling strconv.ParseUint on it.
// If the string cannot be parsed as an unsigned integer, then an empty opt is returned.
// The `base` and `bitSize` parameters are passed as parameters to strconv.ParseUint.
// See the documentation for strconv.ParseUint for details.
func TryParseUint[E constraint.String, F constraint.UnsignedInteger](base int, bitSize int) func(E) opt.Optional[F] {
	return func(e E) opt.Optional[F] {
		i, err := strconv.ParseUint(string(e), base, bitSize)
		if err != nil {
			return opt.Empty[F]()
		}
		return opt.Of(F(i))
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

// TryParseFloat returns a function that accepts a value of any string type and returns the result of calling strconv.ParseFloat on it.
// If the string cannot be parsed as a floating-point number, the provided `or` default value is returned.
// The `bitSize` parameter is passed as a parameter to strconv.ParseFloat.
// See the documentation for strconv.ParseFloat for details.
func TryParseFloat[E constraint.String, F constraint.Float](bitSize int) func(E) opt.Optional[F] {
	return func(e E) opt.Optional[F] {
		f, err := strconv.ParseFloat(string(e), bitSize)
		if err != nil {
			return opt.Empty[F]()
		}
		return opt.Of(F(f))
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

// TryParseComplex returns a function that accepts a value of any string type and returns the result of calling strconv.ParseComplex on it.
// If the string cannot be parsed as a complex number, then an empty opt is returned.
// The `bitSize` parameter is passed as a parameter to strconv.ParseComplex.
// See the documentation for strconv.ParseComplex for details.
func TryParseComplex[E constraint.String, F constraint.Complex](bitSize int) func(E) opt.Optional[F] {
	return func(e E) opt.Optional[F] {
		c, err := strconv.ParseComplex(string(e), bitSize)
		if err != nil {
			return opt.Empty[F]()
		}
		return opt.Of(F(c))
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
