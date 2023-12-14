package mapper

import (
	"testing"

	"github.com/jpfourny/papaya/pkg/pred"
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

func TestIfElse(t *testing.T) {
	m := IfElse(
		pred.GreaterThan(0),
		Constant[int](1),
		Constant[int](-1),
	)
	got := m(42)
	want := 1
	if got != want {
		t.Errorf("IfElse(GreaterThan(0), Constant(1), Constant(-1))(42) = %#v; want %#v", got, want)
	}

	got = m(-42)
	want = -1
	if got != want {
		t.Errorf("IfElse(GreaterThan(0), Constant(1), Constant(-1))(-42) = %#v; want %#v", got, want)
	}
}
