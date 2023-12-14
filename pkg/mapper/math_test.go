package mapper

import "testing"

func TestIncrement(t *testing.T) {
	m := Increment[int](2)
	got := m(42)
	want := 44
	if got != want {
		t.Errorf("Increment(2)(42) = %#v; want %#v", got, want)
	}
}

func TestDecrement(t *testing.T) {
	m := Decrement[int](2)
	got := m(42)
	want := 40
	if got != want {
		t.Errorf("Decrement(2)(42) = %#v; want %#v", got, want)
	}
}
