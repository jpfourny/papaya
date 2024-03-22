package res

import (
	"fmt"
	"github.com/jpfourny/papaya/pkg/opt"
)

// Failure represents a failed result.
// A failed result has an error and no value.
type Failure[T any] struct {
	Err error
}

// Assert that Failure[T] implements Result[T]
var _ Result[any] = Failure[any]{}

func (r Failure[T]) Succeeded() bool {
	return false
}

func (r Failure[T]) PartiallySucceeded() bool {
	return false
}

func (r Failure[T]) Failed() bool {
	return true
}

func (r Failure[T]) HasError() bool {
	return true
}

func (r Failure[T]) HasValue() bool {
	return false
}

func (r Failure[T]) Value() opt.Optional[T] {
	return opt.Empty[T]()
}

func (r Failure[T]) Error() opt.Optional[error] {
	return opt.Of(r.Err)
}

func (r Failure[T]) String() string {
	return fmt.Sprintf("Failure(%v)", r.Err)
}
