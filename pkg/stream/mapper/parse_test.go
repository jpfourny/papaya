package mapper

import (
	"testing"
	"time"

	"github.com/jpfourny/papaya/pkg/opt"
)

func TestTryParseBool(t *testing.T) {
	m := TryParseBool[string]()
	got := m("true")
	want := opt.Of(true)
	if got != want {
		t.Errorf("TryParseBool()(true) = %#v; want %#v", got, want)
	}

	m = TryParseBool[string]()
	got = m("false")
	want = opt.Of(false)
	if got != want {
		t.Errorf("TryParseBool()(false) = %#v; want %#v", got, want)
	}

	m = TryParseBool[string]()
	got = m("foo")
	want = opt.Empty[bool]()
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
	want := opt.Of[int](42)
	if got != want {
		t.Errorf("TryParseInt(10, 32)(`42`) = %#v; want %#v", got, want)
	}

	m = TryParseInt[string, int](10, 32)
	got = m("foo")
	want = opt.Empty[int]()
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
	want := opt.Of[uint](42)
	if got != want {
		t.Errorf("TryParseUint(10, 32)(`42`) = %#v; want %#v", got, want)
	}

	m = TryParseUint[string, uint](10, 32)
	got = m("foo")
	want = opt.Empty[uint]()
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
	want32 := opt.Of[float32](42.42)
	if got32 != want32 {
		t.Errorf("TryParseFloat(32)(`42.42`) = %#v; want %#v", got32, want32)
	}

	m32 = TryParseFloat[string, float32](32)
	got32 = m32("foo")
	want32 = opt.Empty[float32]()
	if got32 != want32 {
		t.Errorf("TryParseFloat(32)(`foo`) = %#v; want %#v", got32, want32)
	}

	m64 := TryParseFloat[string, float64](64)
	got64 := m64("42.42")
	want64 := opt.Of[float64](42.42)
	if got64 != want64 {
		t.Errorf("TryParseFloat(64)(`42.42`) = %#v; want %#v", got64, want64)
	}

	m64 = TryParseFloat[string, float64](64)
	got64 = m64("foo")
	want64 = opt.Empty[float64]()
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
	want64 := 42.42
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
	want64 := opt.Of[complex64](42.42 + 42.42i)
	if got64 != want64 {
		t.Errorf("TryParseComplex(64)(`(42.42+42.42i)`) = %#v; want %#v", got64, want64)
	}

	m64 = TryParseComplex[string, complex64](64)
	got64 = m64("foo")
	want64 = opt.Empty[complex64]()
	if got64 != want64 {
		t.Errorf("TryParseComplex(64)(`foo`) = %#v; want %#v", got64, want64)
	}

	m128 := TryParseComplex[string, complex128](128)
	got128 := m128("(42.42+42.42i)")
	want128 := opt.Of[complex128](42.42 + 42.42i)
	if got128 != want128 {
		t.Errorf("TryParseComplex(128)(`(42.42+42.42i)`) = %#v; want %#v", got128, want128)
	}

	m128 = TryParseComplex[string, complex128](128)
	got128 = m128("foo")
	want128 = opt.Empty[complex128]()
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
	want128 := 42.42 + 42.42i
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

func TestTryParseDuration(t *testing.T) {
	m := TryParseDuration[string]()
	got := m("1s")
	want := opt.Of(time.Second)
	if got != want {
		t.Errorf("TryParseDuration()(`1s`) = %#v; want %#v", got, want)
	}

	m = TryParseDuration[string]()
	got = m("foo")
	want = opt.Empty[time.Duration]()
	if got != want {
		t.Errorf("TryParseDuration()(`foo`) = %#v; want %#v", got, want)
	}

	m = TryParseDuration[string]()
	got = m("")
	want = opt.Empty[time.Duration]()
	if got != want {
		t.Errorf("TryParseDuration()(``) = %#v; want %#v", got, want)
	}
}

func TestParseDurationOr(t *testing.T) {
	m := ParseDurationOr[string](time.Second)
	got := m("1s")
	want := time.Second
	if got != want {
		t.Errorf("ParseDurationOr(time.Second)(`1s`) = %#v; want %#v", got, want)
	}

	m = ParseDurationOr[string](time.Second)
	got = m("foo")
	want = time.Second
	if got != want {
		t.Errorf("ParseDurationOr(time.Second)(`foo`) = %#v; want %#v", got, want)
	}

	m = ParseDurationOr[string](time.Second)
	got = m("")
	want = time.Second
	if got != want {
		t.Errorf("ParseDurationOr(time.Second)(``) = %#v; want %#v", got, want)
	}
}

func TestTryParseTime(t *testing.T) {
	m := TryParseTime[string](time.RFC3339)
	got := m("2006-01-02T15:04:05Z")
	want := opt.Of(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC))
	if got != want {
		t.Errorf("TryParseTime(time.RFC3339)(`2006-01-02T15:04:05Z`) = %#v; want %#v", got, want)
	}

	m = TryParseTime[string](time.RFC3339)
	got = m("foo")
	want = opt.Empty[time.Time]()
	if got != want {
		t.Errorf("TryParseTime(time.RFC3339)(`foo`) = %#v; want %#v", got, want)
	}

	m = TryParseTime[string](time.RFC3339)
	got = m("")
	want = opt.Empty[time.Time]()
	if got != want {
		t.Errorf("TryParseTime(time.RFC3339)(``) = %#v; want %#v", got, want)
	}
}

func TestParseTimeOr(t *testing.T) {
	m := ParseTimeOr[string](time.RFC3339, time.Time{})
	got := m("2006-01-02T15:04:05Z")
	want := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
	if got != want {
		t.Errorf("ParseTimeOr(time.RFC3339, time.Time{})(`2006-01-02T15:04:05Z`) = %#v; want %#v", got, want)
	}

	m = ParseTimeOr[string](time.RFC3339, time.Time{})
	got = m("foo")
	want = time.Time{}
	if got != want {
		t.Errorf("ParseTimeOr(time.RFC3339, time.Time{})(`foo`) = %#v; want %#v", got, want)
	}

	m = ParseTimeOr[string](time.RFC3339, time.Time{})
	got = m("")
	want = time.Time{}
	if got != want {
		t.Errorf("ParseTimeOr(time.RFC3339, time.Time{})(``) = %#v; want %#v", got, want)
	}
}

func TestTryParseTimeInLocation(t *testing.T) {
	loc := time.UTC
	m := TryParseTimeInLocation[string](time.RFC3339, loc)
	got := m("2006-01-02T15:04:05Z")
	want := opt.Of(time.Date(2006, 1, 2, 15, 4, 5, 0, loc))
	if got != want {
		t.Errorf("TryParseTimeInLocation(time.RFC3339, loc)(`2006-01-02T15:04:05Z`) = %#v; want %#v", got, want)
	}

	m = TryParseTimeInLocation[string](time.RFC3339, loc)
	got = m("foo")
	want = opt.Empty[time.Time]()
	if got != want {
		t.Errorf("TryParseTimeInLocation(time.RFC3339, loc)(`foo`) = %#v; want %#v", got, want)
	}

	m = TryParseTimeInLocation[string](time.RFC3339, loc)
	got = m("")
	want = opt.Empty[time.Time]()
	if got != want {
		t.Errorf("TryParseTimeInLocation(time.RFC3339, loc)(``) = %#v; want %#v", got, want)
	}
}

func TestParseTimeInLocationOr(t *testing.T) {
	loc := time.UTC
	m := ParseTimeInLocationOr[string](time.RFC3339, loc, time.Time{})
	got := m("2006-01-02T15:04:05Z")
	want := time.Date(2006, 1, 2, 15, 4, 5, 0, loc)
	if got != want {
		t.Errorf("ParseTimeInLocationOr(time.RFC3339, loc, time.Time{})(`2006-01-02T15:04:05Z`) = %#v; want %#v", got, want)
	}

	m = ParseTimeInLocationOr[string](time.RFC3339, loc, time.Time{})
	got = m("foo")
	want = time.Time{}
	if got != want {
		t.Errorf("ParseTimeInLocationOr(time.RFC3339, loc, time.Time{})(`foo`) = %#v; want %#v", got, want)
	}

	m = ParseTimeInLocationOr[string](time.RFC3339, loc, time.Time{})
	got = m("")
	want = time.Time{}
	if got != want {
		t.Errorf("ParseTimeInLocationOr(time.RFC3339, loc, time.Time{})(``) = %#v; want %#v", got, want)
	}
}
