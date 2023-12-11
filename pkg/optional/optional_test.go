package optional

import (
	"testing"

	"github.com/jpfourny/papaya/pkg/pointer"
)

func TestOf(t *testing.T) {
	o := Of(42)
	if !o.IsPresent() {
		t.Errorf("expected IsPresent() to be true")
	}
	if o.Get() != 42 {
		t.Errorf("expected Get() to return 42")
	}
}

func TestOfNillable(t *testing.T) {
	o := OfNillable(pointer.Ref(42))
	if !o.IsPresent() {
		t.Errorf("expected IsPresent() to be true")
	}
	if o.Get() != 42 {
		t.Errorf("expected Get() to return 42")
	}

	o = OfNillable(pointer.Nil[int]())
	if o.IsPresent() {
		t.Errorf("expected IsPresent() to be false")
	}
}

func TestEmpty(t *testing.T) {
	o := Empty[int]()
	if o.IsPresent() {
		t.Errorf("expected IsPresent() to be false")
	}
	if o.Get() != 0 {
		t.Errorf("expected Get() to return 0")
	}
}

func TestIfPresent(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		var called bool
		o := Of(42)
		o.IfPresent(func(i int) {
			called = true
		})
		if !called {
			t.Errorf("expected callback to be called")
		}
	})

	t.Run("None", func(t *testing.T) {
		var called bool
		o := Empty[int]()
		o.IfPresent(func(i int) {
			called = true
		})
		if called {
			t.Errorf("expected callback to not be called")
		}
	})
}

func TestOrElse(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Of(42)
		if o.OrElse(0) != 42 {
			t.Errorf("expected OrElse to return 42")
		}
	})

	t.Run("None", func(t *testing.T) {
		o := Empty[int]()
		if o.OrElse(0) != 0 {
			t.Errorf("expected OrElse to return 0")
		}
	})
}

func TestOrElseGet(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Of(42)
		if o.OrElseGet(func() int { return 0 }) != 42 {
			t.Errorf("expected OrElseGet to return 42")
		}
	})

	t.Run("None", func(t *testing.T) {
		o := Empty[int]()
		if o.OrElseGet(func() int { return 0 }) != 0 {
			t.Errorf("expected OrElseGet to return 0")
		}
	})
}

func TestString(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Of(42)
		if o.String() != "Some(42)" {
			t.Errorf("expected String() to return %q", "Some(42)")
		}
	})

	t.Run("None", func(t *testing.T) {
		o := Empty[int]()
		if o.String() != "None" {
			t.Errorf("expected String() to return %q", "None")
		}
	})
}