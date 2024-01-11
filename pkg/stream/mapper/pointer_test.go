package mapper

import (
	"testing"

	"github.com/jpfourny/papaya/pkg/optional"
	"github.com/jpfourny/papaya/pkg/pointer"
)

func TestPointerRef(t *testing.T) {
	m := PointerRef[int]()
	got := m(42)
	want := pointer.Ref(42)
	if *got != *want {
		t.Errorf("*PointerRef()(42) = %#v; want %#v", *got, *want)
	}
}

func TestPointerDeref(t *testing.T) {
	m := PointerDeref[int]()
	got := m(pointer.Ref(42))
	want := optional.Of(42)
	if got != want {
		t.Errorf("PointerDeref()(&42) = %#v; want %#v", got, want)
	}

	m = PointerDeref[int]()
	got = m(nil)
	want = optional.Empty[int]()
	if got != want {
		t.Errorf("PointerDeref()(nil) = %#v; want %#v", got, want)
	}
}

func TestPointerDerefOr(t *testing.T) {
	m := PointerDerefOr[int](-1)
	got := m(pointer.Ref(42))
	want := 42
	if got != want {
		t.Errorf("PointerDerefOr(-1)(pointer.Ref(42)) = %#v; want %#v", got, want)
	}

	m = PointerDerefOr[int](-1)
	got = m(nil)
	want = -1
	if got != want {
		t.Errorf("PointerDerefOr(-1)(nil) = %#v; want %#v", got, want)
	}
}

func TestPointerDerefOrZero(t *testing.T) {
	m := PointerDerefOrZero[int]()
	got := m(pointer.Ref(42))
	want := 42
	if got != want {
		t.Errorf("PointerDerefOrZero()(pointer.Ref(42)) = %#v; want %#v", got, want)
	}

	m = PointerDerefOrZero[int]()
	got = m(nil)
	want = 0
	if got != want {
		t.Errorf("PointerDerefOrZero()(nil) = %#v; want %#v", got, want)
	}
}
