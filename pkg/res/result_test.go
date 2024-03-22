package res

import (
	"errors"
	"testing"
)

func TestOfSuccess(t *testing.T) {
	r := OfSuccess(42)
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

func TestOfFailure(t *testing.T) {
	r := OfFailure[any](errors.New("error"))
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

func TestOfPartialSuccess(t *testing.T) {
	r := OfPartialSuccess(42, errors.New("error"))
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

func TestOf(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r := Of(42, nil)
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
		r := Of(42, errors.New("error"))
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
