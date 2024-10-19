package mapper

import (
	"github.com/jpfourny/papaya/v2/pkg/opt"
	"github.com/jpfourny/papaya/v2/pkg/ptr"
)

// PtrRef returns a function that accepts a value of type V and returns a pointer to a copy of the value (on the heap).
func PtrRef[V any]() func(V) *V {
	return func(v V) *V {
		return ptr.Ref(v)
	}
}

// PtrDerefOptional returns a function that accepts a pointer and returns the dereferenced value as an opt.Optional, or an empty opt.Optional if the pointer is nil.
func PtrDerefOptional[V any]() func(*V) opt.Optional[V] {
	return func(v *V) opt.Optional[V] {
		return ptr.DerefOptional(v)
	}
}

// PtrDerefOrZero returns a function that accepts a pointer and returns the dereferenced value, or the zero value of type V if the pointer is nil.
func PtrDerefOrZero[V any]() func(*V) V {
	return func(v *V) V {
		return ptr.DerefOrZero(v)
	}
}

// PtrDerefOrDefault returns a function that accepts a pointer and returns the dereferenced value, or the provided default value if the pointer is nil.
func PtrDerefOrDefault[V any](or V) func(*V) V {
	return func(v *V) V {
		return ptr.DerefOrDefault(v, or)
	}
}
