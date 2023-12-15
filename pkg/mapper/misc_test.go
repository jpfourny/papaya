package mapper

import (
	"testing"
)

func TestConstant(t *testing.T) {
	m := Constant[int](42)
	got := m(0)
	want := 42
	if got != want {
		t.Errorf("Constant(42)(0) = %#v; want %#v", got, want)
	}
}

func TestIdentity(t *testing.T) {
	m := Identity[int]()
	got := m(42)
	want := 42
	if got != want {
		t.Errorf("Identity()(42) = %#v; want %#v", got, want)
	}
}
