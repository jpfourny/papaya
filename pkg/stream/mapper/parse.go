package mapper

import (
	"strconv"
	"time"

	"github.com/jpfourny/papaya/pkg/constraint"
	"github.com/jpfourny/papaya/pkg/opt"
)

// TryParseBool returns a function that accepts a value of any string type and returns the result of calling strconv.ParseBool on it.
// If the string cannot be parsed as a boolean, then an empty opt is returned.
// See the documentation for strconv.ParseBool for details.
func TryParseBool[E constraint.String]() func(E) opt.Optional[bool] {
	return func(e E) opt.Optional[bool] {
		b, err := strconv.ParseBool(string(e))
		return opt.Maybe(b, err == nil)
	}
}

// ParseBoolOr returns a function that accepts a value of any string type and returns the result of calling strconv.ParseBool on it.
// If the string cannot be parsed as a boolean, the provided `or` default value is returned.
// See the documentation for strconv.ParseBool for details.
func ParseBoolOr[E constraint.String](or bool) func(E) bool {
	return func(e E) bool {
		return TryParseBool[E]()(e).GetOrDefault(or)
	}
}

// TryParseInt returns a function that accepts a value of any string type and returns the result of calling strconv.ParseInt on it.
// If the string cannot be parsed as an integer, then an empty opt is returned.
// The `base` and `bitSize` parameters are passed as parameters to strconv.ParseInt.
// See the documentation for strconv.ParseInt for details.
func TryParseInt[E constraint.String, F constraint.SignedInteger](base int, bitSize int) func(E) opt.Optional[F] {
	return func(e E) opt.Optional[F] {
		i, err := strconv.ParseInt(string(e), base, bitSize)
		return opt.Maybe(F(i), err == nil)
	}
}

// ParseIntOr returns a function that accepts a value of any string type and returns the result of calling strconv.ParseInt on it.
// If the string cannot be parsed as an integer, the provided `or` default value is returned.
// The `base` and `bitSize` parameters are passed as parameters to strconv.ParseInt.
// See the documentation for strconv.ParseInt for details.
func ParseIntOr[E constraint.String, F constraint.SignedInteger](base int, bitSize int, or F) func(E) F {
	return func(e E) F {
		return TryParseInt[E, F](base, bitSize)(e).GetOrDefault(or)
	}
}

// TryParseUint returns a function that accepts a value of any string type and returns the result of calling strconv.ParseUint on it.
// If the string cannot be parsed as an unsigned integer, then an empty opt is returned.
// The `base` and `bitSize` parameters are passed as parameters to strconv.ParseUint.
// See the documentation for strconv.ParseUint for details.
func TryParseUint[E constraint.String, F constraint.UnsignedInteger](base int, bitSize int) func(E) opt.Optional[F] {
	return func(e E) opt.Optional[F] {
		i, err := strconv.ParseUint(string(e), base, bitSize)
		return opt.Maybe(F(i), err == nil)
	}
}

// ParseUintOr returns a function that accepts a value of any string type and returns the result of calling strconv.ParseUint on it.
// If the string cannot be parsed as an unsigned integer, the provided `or` default value is returned.
// The `base` and `bitSize` parameters are passed as parameters to strconv.ParseUint.
// See the documentation for strconv.ParseUint for details.
func ParseUintOr[E constraint.String, F constraint.UnsignedInteger](base int, bitSize int, or F) func(E) F {
	return func(e E) F {
		return TryParseUint[E, F](base, bitSize)(e).GetOrDefault(or)
	}
}

// TryParseFloat returns a function that accepts a value of any string type and returns the result of calling strconv.ParseFloat on it.
// If the string cannot be parsed as a floating-point number, the provided `or` default value is returned.
// The `bitSize` parameter is passed as a parameter to strconv.ParseFloat.
// See the documentation for strconv.ParseFloat for details.
func TryParseFloat[E constraint.String, F constraint.Float](bitSize int) func(E) opt.Optional[F] {
	return func(e E) opt.Optional[F] {
		f, err := strconv.ParseFloat(string(e), bitSize)
		return opt.Maybe(F(f), err == nil)
	}
}

// ParseFloatOr returns a function that accepts a value of any string type and returns the result of calling strconv.ParseFloat on it.
// If the string cannot be parsed as a floating-point number, the provided `or` default value is returned.
// The `bitSize` parameter is passed as a parameter to strconv.ParseFloat.
// See the documentation for strconv.ParseFloat for details.
func ParseFloatOr[E constraint.String, F constraint.Float](bitSize int, or F) func(E) F {
	return func(e E) F {
		return TryParseFloat[E, F](bitSize)(e).GetOrDefault(or)
	}
}

// TryParseComplex returns a function that accepts a value of any string type and returns the result of calling strconv.ParseComplex on it.
// If the string cannot be parsed as a complex number, then an empty opt is returned.
// The `bitSize` parameter is passed as a parameter to strconv.ParseComplex.
// See the documentation for strconv.ParseComplex for details.
func TryParseComplex[E constraint.String, F constraint.Complex](bitSize int) func(E) opt.Optional[F] {
	return func(e E) opt.Optional[F] {
		c, err := strconv.ParseComplex(string(e), bitSize)
		return opt.Maybe(F(c), err == nil)
	}
}

// ParseComplexOr returns a function that accepts a value of any string type and returns the result of calling strconv.ParseComplex on it.
// If the string cannot be parsed as a complex number, the provided `or` default value is returned.
// The `bitSize` parameter is passed as a parameter to strconv.ParseComplex.
// See the documentation for strconv.ParseComplex for details.
func ParseComplexOr[E constraint.String, F constraint.Complex](bitSize int, or F) func(E) F {
	return func(e E) F {
		return TryParseComplex[E, F](bitSize)(e).GetOrDefault(or)
	}
}

// TryParseDuration returns a function that accepts a value of any string type and returns the result of calling time.ParseDuration on it.
// If the string cannot be parsed as a duration, then an empty opt is returned.
// See the documentation for time.ParseDuration for details.
func TryParseDuration[E constraint.String]() func(E) opt.Optional[time.Duration] {
	return func(e E) opt.Optional[time.Duration] {
		d, err := time.ParseDuration(string(e))
		return opt.Maybe(d, err == nil)
	}
}

// ParseDurationOr returns a function that accepts a value of any string type and returns the result of calling time.ParseDuration on it.
// If the string cannot be parsed as a duration, the provided `or` default value is returned.
// See the documentation for time.ParseDuration for details.
func ParseDurationOr[E constraint.String](or time.Duration) func(E) time.Duration {
	return func(e E) time.Duration {
		return TryParseDuration[E]()(e).GetOrDefault(or)
	}
}

// TryParseTime returns a function that accepts a value of any string type and returns the result of calling time.Parse on it.
// If the string cannot be parsed as a time, then an empty opt is returned.
// The `layout` parameter is passed as a parameter to time.Parse.
// See the documentation for time.Parse for details.
func TryParseTime[E constraint.String](layout string) func(E) opt.Optional[time.Time] {
	return func(e E) opt.Optional[time.Time] {
		t, err := time.Parse(layout, string(e))
		return opt.Maybe(t, err == nil)
	}
}

// ParseTimeOr returns a function that accepts a value of any string type and returns the result of calling time.Parse on it.
// If the string cannot be parsed as a time, the provided `or` default value is returned.
// The `layout` parameter is passed as a parameter to time.Parse.
// See the documentation for time.Parse for details.
func ParseTimeOr[E constraint.String](layout string, or time.Time) func(E) time.Time {
	return func(e E) time.Time {
		return TryParseTime[E](layout)(e).GetOrDefault(or)
	}
}

// TryParseTimeInLocation returns a function that accepts a value of any string type and returns the result of calling time.ParseInLocation on it.
// If the string cannot be parsed as a time, then an empty opt is returned.
// The `layout` and `loc` parameters are passed as parameters to time.ParseInLocation.
// See the documentation for time.ParseInLocation for details.
func TryParseTimeInLocation[E constraint.String](layout string, loc *time.Location) func(E) opt.Optional[time.Time] {
	return func(e E) opt.Optional[time.Time] {
		t, err := time.ParseInLocation(layout, string(e), loc)
		return opt.Maybe(t, err == nil)
	}
}

// ParseTimeInLocationOr returns a function that accepts a value of any string type and returns the result of calling time.ParseInLocation on it.
// If the string cannot be parsed as a time, the provided `or` default value is returned.
// The `layout` and `loc` parameters are passed as parameters to time.ParseInLocation.
// See the documentation for time.ParseInLocation for details.
func ParseTimeInLocationOr[E constraint.String](layout string, loc *time.Location, or time.Time) func(E) time.Time {
	return func(e E) time.Time {
		return TryParseTimeInLocation[E](layout, loc)(e).GetOrDefault(or)
	}
}
