package res

import (
	"errors"
	"fmt"
	"testing"
)

func TestOK(t *testing.T) {
	r := OK(42)
	if !r.Succeeded() {
		t.Errorf("expected Succeeded() to be true")
	}
	if r.PartiallySucceeded() {
		t.Errorf("expected PartiallySucceeded() to be false")
	}
	if r.Failed() {
		t.Errorf("expected Failed() to be false")
	}
	if !r.HasValue() {
		t.Errorf("expected HasValue() to be true")
	}
	if r.HasError() {
		t.Errorf("expected HasError() to be false")
	}
	if r.Value().GetOrZero() != 42 {
		t.Errorf("expected Value() to return 42")
	}
	if r.Error().Present() {
		t.Errorf("expected Error() to be empty")
	}
	if r.String() != "Success(42)" {
		t.Errorf("expected String() to return Success(42)")
	}
}

func TestFail(t *testing.T) {
	r := Fail[any](errors.New("error"))
	if r.Succeeded() {
		t.Errorf("expected Succeeded() to be false")
	}
	if r.PartiallySucceeded() {
		t.Errorf("expected PartiallySucceeded() to be false")
	}
	if !r.Failed() {
		t.Errorf("expected Failed() to be true")
	}
	if r.HasValue() {
		t.Errorf("expected HasValue() to be false")
	}
	if !r.HasError() {
		t.Errorf("expected HasError() to be true")
	}
	if r.Error().GetOrZero().Error() != "error" {
		t.Errorf("expected Error() to return error")
	}
	if r.Value().Present() {
		t.Errorf("expected Value() to be empty")
	}
	if r.String() != "Failure(error)" {
		t.Errorf("expected String() to return Failure(error)")
	}
}

func TestPartial(t *testing.T) {
	r := Partial(42, errors.New("error"))
	if r.Succeeded() {
		t.Errorf("expected Succeeded() to be false")
	}
	if !r.PartiallySucceeded() {
		t.Errorf("expected PartiallySucceeded() to be true")
	}
	if r.Failed() {
		t.Errorf("expected Failed() to be false")
	}
	if !r.HasValue() {
		t.Errorf("expected HasValue() to be true")
	}
	if !r.HasError() {
		t.Errorf("expected HasError() to be true")
	}
	if r.Value().GetOrZero() != 42 {
		t.Errorf("expected Value() to return 42")
	}
	if !r.Error().Present() {
		t.Errorf("expected Error() to be non-empty")
	}
	if r.String() != "PartialSuccess(42, error)" {
		t.Errorf("expected String() to return PartialSuccess(42, error)")
	}
}

func TestMaybe(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r := Maybe(42, nil)
		if !r.Succeeded() {
			t.Errorf("expected Succeeded() to be true")
		}
		if r.PartiallySucceeded() {
			t.Errorf("expected PartiallySucceeded() to be false")
		}
		if r.Failed() {
			t.Errorf("expected Failed() to be false")
		}
		if !r.HasValue() {
			t.Errorf("expected HasValue() to be true")
		}
		if r.HasError() {
			t.Errorf("expected HasError() to be false")
		}
		if r.Value().GetOrZero() != 42 {
			t.Errorf("expected Value() to return 42")
		}
		if r.Error().Present() {
			t.Errorf("expected Error() to be empty")
		}
		if r.String() != "Success(42)" {
			t.Errorf("expected String() to return Success(42)")
		}
	})

	t.Run("Failure", func(t *testing.T) {
		r := Maybe(42, errors.New("error"))
		if r.Succeeded() {
			t.Errorf("expected Succeeded() to be false")
		}
		if r.PartiallySucceeded() {
			t.Errorf("expected PartiallySucceeded() to be false")
		}
		if !r.Failed() {
			t.Errorf("expected Failed() to be true")
		}
		if r.HasValue() {
			t.Errorf("expected HasValue() to be false")
		}
		if !r.HasError() {
			t.Errorf("expected HasError() to be true")
		}
		if r.Error().GetOrZero().Error() != "error" {
			t.Errorf("expected Error() to return error")
		}
		if r.Value().Present() {
			t.Errorf("expected Value() to be empty")
		}
		if r.String() != "Failure(error)" {
			t.Errorf("expected String() to return Failure(error)")
		}
	})
}

func TestMapValue(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r := OK(42)
		mapped := MapValue(r, func(val int) string {
			return fmt.Sprintf("%d", val)
		})
		if !mapped.Succeeded() {
			t.Errorf("expected Succeeded() to be true")
		}
		if mapped.PartiallySucceeded() {
			t.Errorf("expected PartiallySucceeded() to be false")
		}
		if mapped.Failed() {
			t.Errorf("expected Failed() to be false")
		}
		if !mapped.HasValue() {
			t.Errorf("expected HasValue() to be true")
		}
		if mapped.HasError() {
			t.Errorf("expected HasError() to be false")
		}
		if mapped.Value().GetOrZero() != "42" {
			t.Errorf("expected Value() to return 42")
		}
		if mapped.Error().Present() {
			t.Errorf("expected Error() to be empty")
		}
		if mapped.String() != "Success(\"42\")" {
			t.Errorf("expected String() to return Success(\"42\")")
		}
	})

	t.Run("Failure", func(t *testing.T) {
		r := Fail[int](errors.New("error"))
		mapped := MapValue(r, func(val int) string {
			return fmt.Sprintf("%d", val)
		})
		if mapped.Succeeded() {
			t.Errorf("expected Succeeded() to be false")
		}
		if mapped.PartiallySucceeded() {
			t.Errorf("expected PartiallySucceeded() to be false")
		}
		if !mapped.Failed() {
			t.Errorf("expected Failed() to be true")
		}
		if mapped.HasValue() {
			t.Errorf("expected HasValue() to be false")
		}
		if !mapped.HasError() {
			t.Errorf("expected HasError() to be true")
		}
		if mapped.Error().GetOrZero().Error() != "error" {
			t.Errorf("expected Error() to return error")
		}
		if mapped.Value().Present() {
			t.Errorf("expected Value() to be empty")
		}
		if mapped.String() != "Failure(error)" {
			t.Errorf("expected String() to return Failure(error)")
		}
	})

	t.Run("PartialSuccess", func(t *testing.T) {
		r := Partial(42, errors.New("error"))
		mapped := MapValue(r, func(val int) string {
			return fmt.Sprintf("%d", val)
		})
		if mapped.Succeeded() {
			t.Errorf("expected Succeeded() to be false")
		}
		if !mapped.PartiallySucceeded() {
			t.Errorf("expected PartiallySucceeded() to be true")
		}
		if mapped.Failed() {
			t.Errorf("expected Failed() to be false")
		}
		if !mapped.HasValue() {
			t.Errorf("expected HasValue() to be true")
		}
		if !mapped.HasError() {
			t.Errorf("expected HasError() to be true")
		}
		if mapped.Value().GetOrZero() != "42" {
			t.Errorf("expected Value() to return 42")
		}
		if !mapped.Error().Present() {
			t.Errorf("expected Error() to be non-empty")
		}
		if mapped.String() != "PartialSuccess(\"42\", error)" {
			t.Errorf("expected String() to return PartialSuccess(\"42\", error)")
		}
	})
}

func TestMapError(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r := OK(42)
		mapped := MapError[int](r, func(err error) error {
			return errors.New("mapped error")
		})
		if !mapped.Succeeded() {
			t.Errorf("expected Succeeded() to be true")
		}
		if mapped.PartiallySucceeded() {
			t.Errorf("expected PartiallySucceeded() to be false")
		}
		if mapped.Failed() {
			t.Errorf("expected Failed() to be false")
		}
		if !mapped.HasValue() {
			t.Errorf("expected HasValue() to be true")
		}
		if mapped.HasError() {
			t.Errorf("expected HasError() to be false")
		}
		if mapped.Value().GetOrZero() != 42 {
			t.Errorf("expected Value() to return 42")
		}
		if mapped.Error().Present() {
			t.Errorf("expected Error() to be empty")
		}
		if mapped.String() != "Success(42)" {
			t.Errorf("expected String() to return Success(42)")
		}
	})

	t.Run("Failure", func(t *testing.T) {
		r := Fail[int](errors.New("error"))
		mapped := MapError[int](r, func(err error) error {
			return errors.New("mapped error")
		})
		if mapped.Succeeded() {
			t.Errorf("expected Succeeded() to be false")
		}
		if mapped.PartiallySucceeded() {
			t.Errorf("expected PartiallySucceeded() to be false")
		}
		if !mapped.Failed() {
			t.Errorf("expected Failed() to be true")
		}
		if mapped.HasValue() {
			t.Errorf("expected HasValue() to be false")
		}
		if !mapped.HasError() {
			t.Errorf("expected HasError() to be true")
		}
		if mapped.Error().GetOrZero().Error() != "mapped error" {
			t.Errorf("expected Error() to return mapped error")
		}
		if mapped.Value().Present() {
			t.Errorf("expected Value() to be empty")
		}
		if mapped.String() != "Failure(mapped error)" {
			t.Errorf("expected String() to return Failure(mapped error)")
		}
	})

	t.Run("PartialSuccess", func(t *testing.T) {
		r := Partial(42, errors.New("error"))
		mapped := MapError[int](r, func(err error) error {
			return errors.New("mapped error")
		})
		if mapped.Succeeded() {
			t.Errorf("expected Succeeded() to be false")
		}
		if !mapped.PartiallySucceeded() {
			t.Errorf("expected PartiallySucceeded() to be true")
		}
		if mapped.Failed() {
			t.Errorf("expected Failed() to be false")
		}
		if !mapped.HasValue() {
			t.Errorf("expected HasValue() to be true")
		}
		if !mapped.HasError() {
			t.Errorf("expected HasError() to be true")
		}
		if mapped.Value().GetOrZero() != 42 {
			t.Errorf("expected Value() to return 42")
		}
		if mapped.Error().GetOrZero().Error() != "mapped error" {
			t.Errorf("expected Error() to return mapped error")
		}
		if mapped.String() != "PartialSuccess(42, mapped error)" {
			t.Errorf("expected String() to return PartialSuccess(42, mapped error)")
		}
	})
}
