package constraint

// Boolean is a constraint that permits any boolean type.
type Boolean interface {
	~bool
}

// String is a constraint that permits any string type.
type String interface {
	~string
}

// SignedInteger is a constraint that permits any signed integer type.
type SignedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// UnsignedInteger is a constraint that permits any unsigned integer type.
type UnsignedInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a constraint that permits any integer type (signed or unsigned).
type Integer interface {
	SignedInteger | UnsignedInteger
}

// Float is a constraint that permits any floating-point type.
type Float interface {
	~float32 | ~float64
}

// Complex is a constraint that permits any complex numeric type.
type Complex interface {
	~complex64 | ~complex128
}

// RealNumber is a constraint that permits any real number type (integer or float).
type RealNumber interface {
	Integer | Float
}

// Numeric is a constraint that permits any numeric type (integer, float, or complex).
type Numeric interface {
	RealNumber | Complex
}

// Ordered is a constraint that permits any ordered type (ie: integer, float, or string).
type Ordered interface {
	Integer | Float | String
}
