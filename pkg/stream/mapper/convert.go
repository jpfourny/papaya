package mapper

import "github.com/jpfourny/papaya/pkg/constraint"

// BoolToBool returns a function that accepts a value of boolean type E and returns a value of boolean type F.
func BoolToBool[E constraint.Boolean, F constraint.Boolean]() func(E) F {
	return func(e E) F {
		return F(e)
	}
}

// BoolToNum returns a function that accepts a value of boolean type E and returns a value of real-number type F.
// A value of true is converted to 1, and a value of false is converted to 0.
func BoolToNum[E constraint.Boolean, F constraint.RealNumber]() func(E) F {
	return func(e E) F {
		if e {
			return 1
		}
		return 0
	}
}

// NumToBool returns a function that accepts a value of real-number type E and returns a value of boolean type F.
// A value of 0 is converted to false, and any other value is converted to true.
func NumToBool[E constraint.RealNumber, F constraint.Boolean]() func(E) F {
	return func(e E) F {
		return F(e != 0)
	}
}

// NumToNum returns a function that accepts a value of real-number type E and returns a value of real-number type F.
// The value is converted using standard type conversion rules.
func NumToNum[E constraint.RealNumber, F constraint.RealNumber]() func(E) F {
	return func(e E) F {
		return F(e)
	}
}

// ComplexToComplex returns a function that accepts a value of complex-number type E and returns a value of complex-number type F.
func ComplexToComplex[E constraint.Complex, F constraint.Complex]() func(E) F {
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
