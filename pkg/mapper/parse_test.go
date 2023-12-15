package mapper

import (
	"testing"

	"github.com/jpfourny/papaya/pkg/optional"
)

func TestTryParseBool(t *testing.T) {
	m := TryParseBool[string]()
	got := m("true")
	want := optional.Of(true)
	if got != want {
		t.Errorf("TryParseBool()(true) = %#v; want %#v", got, want)
	}

	m = TryParseBool[string]()
	got = m("false")
	want = optional.Of(false)
	if got != want {
		t.Errorf("TryParseBool()(false) = %#v; want %#v", got, want)
	}

	m = TryParseBool[string]()
	got = m("foo")
	want = optional.Empty[bool]()
	if got != want {
		t.Errorf("TryParseBool()(foo) = %#v; want %#v", got, want)
	}
}

func TestParseBoolOr(t *testing.T) {
	m := ParseBoolOr[string](false)
	got := m("true")
	if got != true {
		t.Errorf("ParseBoolOr(false)(`true`) = %#v; want %#v", got, true)
	}

	m = ParseBoolOr[string](true)
	got = m("false")
	if got != false {
		t.Errorf("ParseBoolOr(true)(`false`) = %#v; want %#v", got, false)
	}

	m = ParseBoolOr[string](true)
	got = m("foo")
	if got != true {
		t.Errorf("ParseBoolOr(true)(`foo`) = %#v; want %#v", got, true)
	}
}

func TestTryParseInt(t *testing.T) {
	m := TryParseInt[string, int](10, 32)
	got := m("42")
	want := optional.Of[int](42)
	if got != want {
		t.Errorf("TryParseInt(10, 32)(`42`) = %#v; want %#v", got, want)
	}

	m = TryParseInt[string, int](10, 32)
	got = m("foo")
	want = optional.Empty[int]()
	if got != want {
		t.Errorf("TryParseInt(10, 32)(`foo`) = %#v; want %#v", got, want)
	}
}

func TestParseIntOr(t *testing.T) {
	m := ParseIntOr[string](10, 32, 0)
	got := m("42")
	want := 42
	if got != want {
		t.Errorf("ParseIntOr(10, 32, 0)(`42`) = %#v; want %#v", got, want)
	}

	m = ParseIntOr[string](10, 32, 0)
	got = m("foo")
	want = 0
	if got != want {
		t.Errorf("ParseIntOr(10, 32, 0)(`foo`) = %#v; want %#v", got, want)
	}
}

func TestTryParseUint(t *testing.T) {
	m := TryParseUint[string, uint](10, 32)
	got := m("42")
	want := optional.Of[uint](42)
	if got != want {
		t.Errorf("TryParseUint(10, 32)(`42`) = %#v; want %#v", got, want)
	}

	m = TryParseUint[string, uint](10, 32)
	got = m("foo")
	want = optional.Empty[uint]()
	if got != want {
		t.Errorf("TryParseUint(10, 32)(`foo`) = %#v; want %#v", got, want)
	}
}

func TestParseUintOr(t *testing.T) {
	m := ParseUintOr[string](10, 32, uint32(0))
	got := m("42")
	want := uint32(42)
	if got != want {
		t.Errorf("ParseUintOr(10, 32, 0)(`42`) = %#v; want %#v", got, want)
	}

	m = ParseUintOr[string](10, 32, uint32(0))
	got = m("foo")
	want = 0
	if got != want {
		t.Errorf("ParseUintOr(10, 32, 0)(`foo`) = %#v; want %#v", got, want)
	}
}

func TestTryParseFloat(t *testing.T) {
	m32 := TryParseFloat[string, float32](32)
	got32 := m32("42.42")
	want32 := optional.Of[float32](42.42)
	if got32 != want32 {
		t.Errorf("TryParseFloat(32)(`42.42`) = %#v; want %#v", got32, want32)
	}

	m32 = TryParseFloat[string, float32](32)
	got32 = m32("foo")
	want32 = optional.Empty[float32]()
	if got32 != want32 {
		t.Errorf("TryParseFloat(32)(`foo`) = %#v; want %#v", got32, want32)
	}

	m64 := TryParseFloat[string, float64](64)
	got64 := m64("42.42")
	want64 := optional.Of[float64](42.42)
	if got64 != want64 {
		t.Errorf("TryParseFloat(64)(`42.42`) = %#v; want %#v", got64, want64)
	}

	m64 = TryParseFloat[string, float64](64)
	got64 = m64("foo")
	want64 = optional.Empty[float64]()
	if got64 != want64 {
		t.Errorf("TryParseFloat(64)(`foo`) = %#v; want %#v", got64, want64)
	}
}

func TestParseFloatOr(t *testing.T) {
	m32 := ParseFloatOr[string, float32](32, 0)
	got32 := m32("42.42")
	want32 := float32(42.42)
	if got32 != want32 {
		t.Errorf("ParseFloatOr(32, 0)(`42.42`) = %#v; want %#v", got32, want32)
	}

	m32 = ParseFloatOr[string, float32](32, -1)
	got32 = m32("foo")
	want32 = float32(-1)
	if got32 != want32 {
		t.Errorf("ParseFloatOr(32, -1)(`foo`) = %#v; want %#v", got32, want32)
	}

	m64 := ParseFloatOr[string, float64](64, 0)
	got64 := m64("42.42")
	want64 := float64(42.42)
	if got64 != want64 {
		t.Errorf("ParseFloatOr(64, 0)(`42.42`) = %#v; want %#v", got64, want64)
	}

	m64 = ParseFloatOr[string, float64](64, 0)
	got64 = m64("foo")
	want64 = float64(0)
	if got64 != want64 {
		t.Errorf("ParseFloatOr(64, 0)(`foo`) = %#v; want %#v", got64, want64)
	}
}

func TestTryParseComplex(t *testing.T) {
	m64 := TryParseComplex[string, complex64](64)
	got64 := m64("(42.42+42.42i)")
	want64 := optional.Of[complex64](42.42 + 42.42i)
	if got64 != want64 {
		t.Errorf("TryParseComplex(64)(`(42.42+42.42i)`) = %#v; want %#v", got64, want64)
	}

	m64 = TryParseComplex[string, complex64](64)
	got64 = m64("foo")
	want64 = optional.Empty[complex64]()
	if got64 != want64 {
		t.Errorf("TryParseComplex(64)(`foo`) = %#v; want %#v", got64, want64)
	}

	m128 := TryParseComplex[string, complex128](128)
	got128 := m128("(42.42+42.42i)")
	want128 := optional.Of[complex128](42.42 + 42.42i)
	if got128 != want128 {
		t.Errorf("TryParseComplex(128)(`(42.42+42.42i)`) = %#v; want %#v", got128, want128)
	}

	m128 = TryParseComplex[string, complex128](128)
	got128 = m128("foo")
	want128 = optional.Empty[complex128]()
	if got128 != want128 {
		t.Errorf("TryParseComplex(128)(`foo`) = %#v; want %#v", got128, want128)
	}
}

func TestParseComplexOr(t *testing.T) {
	m64 := ParseComplexOr[string, complex64](64, 0)
	got64 := m64("(42.42+42.42i)")
	want64 := complex64(42.42 + 42.42i)
	if got64 != want64 {
		t.Errorf("ParseComplexOr(64, 0)(`(42.42+42.42i)`) = %#v; want %#v", got64, want64)
	}

	m64 = ParseComplexOr[string, complex64](64, -1)
	got64 = m64("foo")
	want64 = complex64(-1)
	if got64 != want64 {
		t.Errorf("ParseComplexOr(64, -1)(`foo`) = %#v; want %#v", got64, want64)
	}

	m128 := ParseComplexOr[string, complex128](128, 0)
	got128 := m128("(42.42+42.42i)")
	want128 := complex128(42.42 + 42.42i)
	if got128 != want128 {
		t.Errorf("ParseComplexOr(128, 0)(`(42.42+42.42i)`) = %#v; want %#v", got128, want128)
	}

	m128 = ParseComplexOr[string, complex128](128, -1)
	got128 = m128("foo")
	want128 = complex128(-1)
	if got128 != want128 {
		t.Errorf("ParseComplexOr(128, 0)(`foo`) = %#v; want %#v", got128, want128)
	}
}
