package opt

import (
	"testing"
)

func TestOf(t *testing.T) {
	o := Of(42)
	if !o.Present() {
		t.Errorf("expected Present() to be true")
	}
	if o.GetOrZero() != 42 {
		t.Errorf("expected Get() to return 42")
	}
}

func TestMaybe(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Maybe[int](42, true)
		if !o.Present() {
			t.Errorf("expected Present() to be true")
		}
		if o.GetOrZero() != 42 {
			t.Errorf("expected Get() to return 42")
		}
	})

	t.Run("None", func(t *testing.T) {
		o := Maybe[int](42, false)
		if o.Present() {
			t.Errorf("expected Present() to be false")
		}
		if o.GetOrZero() != 0 {
			t.Errorf("expected Get() to return 0")
		}
	})
}

func TestEmpty(t *testing.T) {
	o := Empty[int]()
	if o.Present() {
		t.Errorf("expected Present() to be false")
	}
	if o.GetOrZero() != 0 {
		t.Errorf("expected Get() to return 0")
	}
}

func TestAny(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Any(Of(42), Empty[int]())
		if !o.Present() {
			t.Errorf("expected Present() to be true")
		}
		if o.GetOrZero() != 42 {
			t.Errorf("expected Get() to return 42")
		}

		o = Any(Empty[int](), Of(42))
		if !o.Present() {
			t.Errorf("expected Present() to be true")
		}
		if o.GetOrZero() != 42 {
			t.Errorf("expected Get() to return 42")
		}
	})

	t.Run("None", func(t *testing.T) {
		o := Any(Empty[int](), Empty[int]())
		if o.Present() {
			t.Errorf("expected Present() to be false")
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Of(42)
		m := Map(o, func(i int) string { return "foo" })
		if !m.Present() {
			t.Errorf("expected Present() to be true")
		}
		if m.GetOrZero() != "foo" {
			t.Errorf("expected Get() to return %q", "foo")
		}
	})

	t.Run("None", func(t *testing.T) {
		o := Empty[int]()
		m := Map(o, func(i int) string { return "foo" })
		if m.Present() {
			t.Errorf("expected Present() to be false")
		}
	})
}

func TestPresent(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Of(42)
		if !o.Present() {
			t.Errorf("expected Present() to be true")
		}
	})

	t.Run("None", func(t *testing.T) {
		o := Empty[int]()
		if o.Present() {
			t.Errorf("expected Present() to be false")
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Of(42)
		if v, ok := o.Get(); !ok || v != 42 {
			t.Errorf("expected Get() to return (42, true)")
		}
	})

	t.Run("None", func(t *testing.T) {
		o := Empty[int]()
		if v, ok := o.Get(); ok || v != 0 {
			t.Errorf("expected Get() to return (0, false)")
		}
	})
}

func TestGetOrDefault(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Of(42)
		if o.GetOrDefault(0) != 42 {
			t.Errorf("expected GetOrDefault to return 42")
		}
	})

	t.Run("None", func(t *testing.T) {
		o := Empty[int]()
		if o.GetOrDefault(0) != 0 {
			t.Errorf("expected GetOrDefault to return 0")
		}
	})
}

func TestGetOrFunc(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Of(42)
		if o.GetOrFunc(func() int { return 0 }) != 42 {
			t.Errorf("expected GetOrFunc to return 42")
		}
	})

	t.Run("None", func(t *testing.T) {
		o := Empty[int]()
		if o.GetOrFunc(func() int { return 0 }) != 0 {
			t.Errorf("expected GetOrFunc to return 0")
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Of(42)

		o2 := o.Filter(func(i int) bool { return i == 42 })
		if !o2.Present() {
			t.Errorf("expected Present() to be true")
		}
		if o2.GetOrZero() != 42 {
			t.Errorf("expected Get() to return 42")
		}

		o3 := o.Filter(func(i int) bool { return i != 42 })
		if o3.Present() {
			t.Errorf("expected Present() to be false")
		}
	})

	t.Run("None", func(t *testing.T) {
		o := Empty[int]()

		o2 := o.Filter(func(i int) bool { return i == 42 })
		if o2.Present() {
			t.Errorf("expected Present() to be false")
		}

		o3 := o.Filter(func(i int) bool { return i != 42 })
		if o3.Present() {
			t.Errorf("expected Present() to be false")
		}
	})
}

func TestIfPresent(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		var called bool
		o := Of(42)
		if !o.IfPresent(func(i int) {
			called = true
		}) {
			t.Errorf("expected IfPresent to return true")
		}
		if !called {
			t.Errorf("expected callback to be called")
		}
	})

	t.Run("None", func(t *testing.T) {
		var called bool
		o := Empty[int]()
		if o.IfPresent(func(i int) {
			called = true
		}) {
			t.Errorf("expected IfPresent to return false")
		}
		if called {
			t.Errorf("expected callback to not be called")
		}
	})
}

func TestIfPresentElse(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		var called bool
		o := Of(42)
		if !o.IfPresentElse(
			func(i int) { called = true },
			func() { t.Errorf("expected Else callback to not be called") },
		) {
			t.Errorf("expected IfPresentElse to return true")
		}
		if !called {
			t.Errorf("expected callback to be called")
		}
	})

	t.Run("None", func(t *testing.T) {
		var called bool
		o := Empty[int]()
		if o.IfPresentElse(
			func(i int) { called = true },
			func() { called = true },
		) {
			t.Errorf("expected IfPresentElse to return false")
		}
		if !called {
			t.Errorf("expected Else callback to be called")
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
