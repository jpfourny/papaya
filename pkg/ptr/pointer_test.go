package ptr

import "testing"

func TestNil(t *testing.T) {
	var p = Nil[int]()
	if p != nil {
		t.Errorf("Nil() = %v; want nil", p)
	}
}

func TestRef(t *testing.T) {
	var i = 42
	p := Ref(i)
	if p == nil {
		t.Errorf("Ref(%v) = %v; want non-nil", i, p)
		return
	}
	if *p != i {
		t.Errorf("Ref(%v) = %v; want %v", i, p, i)
	}
}

func TestDerefOptional(t *testing.T) {
	var i = 42
	var nilPtr *int
	if got := DerefOptional(nilPtr); got.Present() {
		t.Errorf("DerefOptional(%v) = %v; want empty opt", nilPtr, got)
	}
	if got := DerefOptional(Ref(i)); !got.Present() || got.GetOrZero() != i {
		t.Errorf("DerefOptional(%v) = %v; want %v", Ref(i), got, i)
	}
}

func TestDerefOrZero(t *testing.T) {
	var i = 42
	var nilPtr *int
	if got := DerefOrZero(nilPtr); got != 0 {
		t.Errorf("DerefOrZero(%v) = %v; want 0", nilPtr, got)
	}
	if got := DerefOrZero(Ref(i)); got != i {
		t.Errorf("DerefOrZero(%v) = %v; want %v", Ref(i), got, i)
	}
}

func TestDerefOrDefault(t *testing.T) {
	var i = 42
	var nilPtr *int
	if got := DerefOrDefault(nilPtr, i); got != i {
		t.Errorf("DerefOr(%v, %v) = %v; want %v", nilPtr, i, got, i)
	}
	if got := DerefOrDefault(Ref(i), i+1); got != i {
		t.Errorf("DerefOr(%v, %v) = %v; want %v", Ref(i), i+1, got, i)
	}
}

func TestDerefOrFunc(t *testing.T) {
	var i = 42
	var nilPtr *int
	if got := DerefOrFunc(nilPtr, func() int { return i }); got != i {
		t.Errorf("DerefOrElse(%v, %v) = %v; want %v", nilPtr, i, got, i)
	}
	if got := DerefOrFunc(Ref(i), func() int { return i + 1 }); got != i {
		t.Errorf("DerefOrElse(%v, %v) = %v; want %v", Ref(i), i+1, got, i)
	}
}
