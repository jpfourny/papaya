package mapper

import "testing"

func TestSprint(t *testing.T) {
	m := Sprint[int]()
	got := m(42)
	want := "42"
	if got != want {
		t.Errorf("Sprint()(testStringer(42)) = %#v; want %#v", got, want)
	}
}

func TestSprintf(t *testing.T) {
	m := Sprintf[int]("%d!")
	got := m(42)
	want := "42!"
	if got != want {
		t.Errorf("Sprintf(`%%d!`)(42) = %#v; want %#v", got, want)
	}
}

func TestFormatBool(t *testing.T) {
	m := FormatBool[bool]()
	got := m(true)
	want := "true"
	if got != want {
		t.Errorf("FormatBool()(true) = %#v; want %#v", got, want)
	}
}

func TestFormatInt(t *testing.T) {
	m := FormatInt[int](10)
	got := m(42)
	want := "42"
	if got != want {
		t.Errorf("FormatInt(10)(42) = %#v; want %#v", got, want)
	}

	m = FormatInt[int](2)
	got = m(42)
	want = "101010"
	if got != want {
		t.Errorf("FormatInt(2)(42) = %#v; want %#v", got, want)
	}
}

func TestFormatUint(t *testing.T) {
	m := FormatUint[uint](10)
	got := m(42)
	want := "42"
	if got != want {
		t.Errorf("FormatUint(10)(42) = %#v; want %#v", got, want)
	}

	m = FormatUint[uint](2)
	got = m(42)
	want = "101010"
	if got != want {
		t.Errorf("FormatUint(2)(42) = %#v; want %#v", got, want)
	}
}

func TestFormatFloat(t *testing.T) {
	m := FormatFloat[float64]('E', 3, 64)
	got := m(42.42)
	want := "4.242E+01"
	if got != want {
		t.Errorf("FormatFloat('E', 3, 64)(42.42) = %#v; want %#v", got, want)
	}

	m = FormatFloat[float64]('f', 2, 64)
	got = m(42.42)
	want = "42.42"
	if got != want {
		t.Errorf("FormatFloat('g', 2, 64)(42.42) = %#v; want %#v", got, want)
	}
}

func TestFormatComplex(t *testing.T) {
	m64 := FormatComplex[complex64]('E', 3)
	got64 := m64(42.42 + 42.42i)
	want64 := "(4.242E+01+4.242E+01i)"
	if got64 != want64 {
		t.Errorf("FormatComplex('E', 3)(42.42+42.42i) = %#v; want %#v", got64, want64)
	}

	m128 := FormatComplex[complex128]('f', 2)
	got128 := m128(42.42 + 42.42i)
	want128 := "(42.42+42.42i)"
	if got128 != want128 {
		t.Errorf("FormatComplex('g', 2)(42.42+42.42i) = %#v; want %#v", got128, want128)
	}
}
