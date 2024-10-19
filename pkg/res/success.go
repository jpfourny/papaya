package res

import (
	"fmt"
	"github.com/jpfourny/papaya/v2/pkg/opt"
)

// Success represents a successful result.
// A successful result has a value and no error.
type Success[T any] struct {
	Val T
}

// Assert that Success[T] implements Result[T]
var _ Result[any] = Success[any]{}

func (r Success[T]) Succeeded() bool {
	return true
}

func (r Success[T]) PartiallySucceeded() bool {
	return false
}

func (r Success[T]) Failed() bool {
	return false
}

func (r Success[T]) HasError() bool {
	return false
}

func (r Success[T]) HasValue() bool {
	return true
}

func (r Success[T]) Value() opt.Optional[T] {
	return opt.Of(r.Val)
}

func (r Success[T]) Error() opt.Optional[error] {
	return opt.Empty[error]()
}

func (r Success[T]) String() string {
	return fmt.Sprintf("Success(%#v)", r.Val)
}
