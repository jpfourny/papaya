package mapper

import (
	"fmt"
	"strconv"

	"github.com/jpfourny/papaya/pkg/constraint"
	"github.com/jpfourny/papaya/pkg/pointer"
)

// Identity returns a function that accepts a value of any type E and returns that value.
func Identity[E any]() func(E) E {
	return func(e E) E {
		return e
	}
}

// Constant returns a function that accepts a value of any type E and returns the provided constant value of type F.
func Constant[E any, F any](c F) func(E) F {
	return func(E) F {
		return c
	}
}

// IfElse returns a function that accepts a value of any type E and returns the result of calling either the `ifTrue` or `ifFalse` function, which return a value of type F.
// If the given `cond` function returns true, the `ifTrue` function is used; otherwise, the `ifFalse` function is used.
func IfElse[E, F any](cond func(E) bool, ifTrue func(E) F, ifFalse func(E) F) func(E) F {
	return func(e E) F {
		if cond(e) {
			return ifTrue(e)
		}
		return ifFalse(e)
	}
}

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

// FormatInt returns a function that accepts a value of any integer type and returns the result of calling strconv.FormatInt on it.
// The `base` parameter is passed as a parameter to strconv.FormatInt.
// See the documentation for strconv.FormatInt for details.
func FormatInt[E constraint.Integer](base int) func(E) string {
	return func(e E) string {
		return strconv.FormatInt(int64(e), base)
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

// FormatUint returns a function that accepts a value of any unsigned integer type and returns the result of calling strconv.FormatUint on it.
// The `base` parameter is passed as a parameter to strconv.FormatUint.
// See the documentation for strconv.FormatUint for details.
func FormatUint[E constraint.UnsignedInteger](base int) func(E) string {
	return func(e E) string {
		return strconv.FormatUint(uint64(e), base)
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

// FormatFloat returns a function that accepts a value of any floating-point type and returns the result of calling strconv.FormatFloat on it.
// The `format`, `prec` and `bitSize` parameters are passed as parameters to strconv.FormatFloat.
// See the documentation for strconv.FormatFloat for details.
func FormatFloat[E constraint.Float](format byte, prec int, bitSize int) func(E) string {
	return func(e E) string {
		return strconv.FormatFloat(float64(e), format, prec, bitSize)
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

// FormatComplex returns a function that accepts a value of any complex numeric type and returns the result of calling strconv.FormatComplex on it.
// The `format` and `prec` parameters are passed as parameters to strconv.FormatComplex.
// See the documentation for strconv.FormatComplex for details.
func FormatComplex[E constraint.Complex](format byte, prec int) func(E) string {
	return func(e E) string {
		return strconv.FormatComplex(complex128(e), format, prec, 128)
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

// PointerRef returns a function that accepts any value and returns a pointer to a copy of that value.
// The pointer is created using the pointer.Ref function.
func PointerRef[E any]() func(E) *E {
	return func(e E) *E {
		return pointer.Ref(e)
	}
}

// PointerDerefOr returns a function that accepts a pointer to a value and returns the value of that pointer, or the provided `or` default value if the pointer is nil.
// The value is retrieved using the pointer.DerefOr function.
func PointerDerefOr[E any](or E) func(*E) E {
	return func(e *E) E {
		return pointer.DerefOr(e, or)
	}
}

// PointerDerefOrZero returns a function that accepts a pointer to a value and returns the value of that pointer, or the zero value of that type if the pointer is nil.
// The value is retrieved using the pointer.DerefOrZero function.
func PointerDerefOrZero[E any]() func(*E) E {
	return func(e *E) E {
		return pointer.DerefOrZero(e)
	}
}

// BoolToBool returns a function that accepts a value of boolean type E and returns a value of boolean type F.
func BoolToBool[E constraint.Boolean, F constraint.Boolean]() func(E) F {
	return func(e E) F {
		return F(e)
	}
}

// BoolToNumber returns a function that accepts a value of boolean type E and returns a value of real number type F.
// A value of true is converted to 1, and a value of false is converted to 0.
func BoolToNumber[E constraint.Boolean, F constraint.RealNumber]() func(E) F {
	return func(e E) F {
		if e {
			return 1
		}
		return 0
	}
}

// NumberToBool returns a function that accepts a value of real number type E and returns a value of boolean type F.
// A value of 0 is converted to false, and any other value is converted to true.
func NumberToBool[E constraint.RealNumber, F constraint.Boolean]() func(E) F {
	return func(e E) F {
		return F(e != 0)
	}
}

// NumberToNumber returns a function that accepts a value of real number type E and returns a value of real number type F.
// The value is converted using standard type conversion rules.
func NumberToNumber[E constraint.RealNumber, F constraint.RealNumber]() func(E) F {
	return func(e E) F {
		return F(e)
	}
}

// StringToString returns a function that accepts a value of string type E and returns a value of string type F.
func StringToString[E constraint.String, F constraint.String]() func(E) F {
	return func(e E) F {
		return F(e)
	}
}

// ComplexToComplex returns a function that accepts a value of complex number type E and returns a value of complex number type F.
func ComplexToComplex[E constraint.Complex, F constraint.Complex]() func(E) F {
	return func(e E) F {
		return F(e)
	}
}
