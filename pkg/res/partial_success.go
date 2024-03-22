package res

import (
	"fmt"
	"github.com/jpfourny/papaya/pkg/opt"
)

// PartialSuccess represents a partially-successful result.
// A partially successful result has a value and an error.
// This is useful when a result is expected to have multiple values, and some of them are missing.
type PartialSuccess[T any] struct {
	Val T
	Err error
}

// Assert that PartialSuccess[T] implements Result[T]
var _ Result[any] = PartialSuccess[any]{}

func (r PartialSuccess[T]) Succeeded() bool {
	return false
}

func (r PartialSuccess[T]) PartiallySucceeded() bool {
	return true
}

func (r PartialSuccess[T]) Failed() bool {
	return false
}

func (r PartialSuccess[T]) HasError() bool {
	return true
}

func (r PartialSuccess[T]) HasValue() bool {
	return true
}

func (r PartialSuccess[T]) Value() opt.Optional[T] {
	return opt.Of(r.Val)
}

func (r PartialSuccess[T]) Error() opt.Optional[error] {
	return opt.Of(r.Err)
}

func (r PartialSuccess[T]) String() string {
	return fmt.Sprintf("PartialSuccess(%#v, %v)", r.Val, r.Err)
}
