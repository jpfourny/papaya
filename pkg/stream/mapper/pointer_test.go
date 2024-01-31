package mapper

import (
	"testing"

	"github.com/jpfourny/papaya/pkg/opt"
	"github.com/jpfourny/papaya/pkg/ptr"
)

func TestPtrRef(t *testing.T) {
	m := PtrRef[int]()
	got := m(42)
	want := ptr.Ref(42)
	if *got != *want {
		t.Errorf("*PtrRef()(42) = %#v; want %#v", *got, *want)
	}
}

func TestPtrDerefOptional(t *testing.T) {
	m := PtrDerefOptional[int]()
	got := m(ptr.Ref(42))
	want := opt.Of(42)
	if got != want {
		t.Errorf("PtrDerefOptional()(&42) = %#v; want %#v", got, want)
	}

	m = PtrDerefOptional[int]()
	got = m(nil)
	want = opt.Empty[int]()
	if got != want {
		t.Errorf("PtrDerefOptional()(nil) = %#v; want %#v", got, want)
	}
}

func TestPtrDerefOrDefault(t *testing.T) {
	m := PtrDerefOrDefault[int](-1)
	got := m(ptr.Ref(42))
	want := 42
	if got != want {
		t.Errorf("PtrDerefOrDefault(-1)(ptr.Ref(42)) = %#v; want %#v", got, want)
	}

	m = PtrDerefOrDefault[int](-1)
	got = m(nil)
	want = -1
	if got != want {
		t.Errorf("PtrDerefOrDefault(-1)(nil) = %#v; want %#v", got, want)
	}
}

func TestPtrDerefOrZero(t *testing.T) {
	m := PtrDerefOrZero[int]()
	got := m(ptr.Ref(42))
	want := 42
	if got != want {
		t.Errorf("PtrDerefOrZero()(ptr.Ref(42)) = %#v; want %#v", got, want)
	}

	m = PtrDerefOrZero[int]()
	got = m(nil)
	want = 0
	if got != want {
		t.Errorf("PtrDerefOrZero()(nil) = %#v; want %#v", got, want)
	}
}
