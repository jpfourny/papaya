package res

import (
	"fmt"
	"github.com/jpfourny/papaya/pkg/opt"
)

// Result represents the result of an operation that may have a value and/or an error.
// The result can be successful, failed, or partially successful.
// A successful result has a value and no error.
// A failed result has an error and no value.
// A partially successful result has both a value and an error.
// The result can be queried for its value and/or error, and for its success/failure status.
type Result[T any] interface {
	fmt.Stringer

	// Succeeded returns true for a successful result; false otherwise.
	// A successful result is one that has a value and no error.
	Succeeded() bool

	// PartiallySucceeded returns true for a partially successful result that includes both a value and an error.
	// This is useful when a result is expected to have multiple values, and some of them are missing.
	PartiallySucceeded() bool

	// Failed returns true for a failed result; false otherwise.
	// A failed result is one that has an error and no value.
	Failed() bool

	// HasError returns true if the result has an error; false otherwise.
	// This is true for both failed and partially successful results.
	HasError() bool

	// HasValue returns true if the result has a value; false otherwise.
	// This is true for both successful and partially successful results.
	HasValue() bool

	// Value returns the value of the result as an optional.
	// If the result has a value, the optional is non-empty; otherwise, it is empty.
	Value() opt.Optional[T]

	// Error returns the error of the result as an optional.
	// If the result has an error, the optional is non-empty; otherwise, it is empty.
	Error() opt.Optional[error]
}

// OK returns a successful result with the provided value.
func OK[T any](val T) Success[T] {
	return Success[T]{Val: val}
}

// Partial returns a partially-successful result with the provided value and error.
func Partial[T any](val T, err error) PartialSuccess[T] {
	return PartialSuccess[T]{Val: val, Err: err}
}

// Fail returns a failed result with the provided error.
func Fail[T any](err error) Failure[T] {
	return Failure[T]{Err: err}
}

// Maybe returns a result based on the provided value and error.
// If the error is nil, the result is successful; otherwise, it is failed.
// This is a convenience function that creates a successful or failed result based on the provided error.
func Maybe[T any](val T, err error) Result[T] {
	if err != nil {
		return Fail[T](err)
	}
	return OK[T](val)
}

// MapValue maps the value of the result to a new value using the provided mapper function.
// The error of the result, if any, is unchanged.
func MapValue[T, U any](r Result[T], valueMapper func(T) U) Result[U] {
	errorMapper := func(err error) error { return err }
	return Map[T, U](r, valueMapper, errorMapper)
}

// MapError maps the error of the result to a new error using the provided mapper function.
// The value of the result, if any, is unchanged.
func MapError[T any](r Result[T], errorMapper func(error) error) Result[T] {
	valueMapper := func(val T) T { return val }
	return Map[T, T](r, valueMapper, errorMapper)
}

// Map maps the result to a new result using the provided value and error mapper functions.
// The value mapper function maps the value of the result to a new value.
// The error mapper function maps the error of the result to a new error.
func Map[T, U any](r Result[T], valueMapper func(T) U, errorMapper func(error) error) Result[U] {
	if r.Succeeded() {
		return OK[U](valueMapper(r.Value().GetOrZero()))
	}
	if r.PartiallySucceeded() {
		return Partial[U](valueMapper(r.Value().GetOrZero()), errorMapper(r.Error().GetOrZero()))
	}
	return Fail[U](errorMapper(r.Error().GetOrZero()))
}
