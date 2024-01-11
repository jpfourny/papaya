package mapper

import "testing"

type testBool bool

func TestBoolToBool(t *testing.T) {
	t.Run("bool-to-testBool", func(t *testing.T) {
		m := BoolToBool[bool, testBool]()
		got := m(true)
		want := testBool(true)
		if got != want {
			t.Errorf("BoolToBool()(true) = %#v; want %#v", got, want)
		}
		got = m(false)
		want = false
		if got != want {
			t.Errorf("BoolToBool()(false) = %#v; want %#v", got, want)
		}
	})

	t.Run("testBool-to-bool", func(t *testing.T) {
		m := BoolToBool[testBool, bool]()
		got := m(true)
		if got != true {
			t.Errorf("BoolToBool()(true) = %#v; want %#v", got, true)
		}
		got = m(false)
		if got != false {
			t.Errorf("BoolToBool()(false) = %#v; want %#v", got, false)
		}
	})
}

type testString string

func TestStringToString(t *testing.T) {
	t.Run("string-to-testString", func(t *testing.T) {
		m := StringToString[string, testString]()
		got := m("foo")
		want := testString("foo")
		if got != want {
			t.Errorf("StringToString()(foo) = %#v; want %#v", got, want)
		}
	})

	t.Run("testString-to-string", func(t *testing.T) {
		m := StringToString[testString, string]()
		got := m("foo")
		want := "foo"
		if got != want {
			t.Errorf("StringToString()(foo) = %#v; want %#v", got, want)
		}
	})
}

type testInt int

func TestNumberToNumber(t *testing.T) {
	t.Run("int-to-float64", func(t *testing.T) {
		m := NumberToNumber[int, float64]()
		got := m(42)
		want := float64(42)
		if got != want {
			t.Errorf("NumberToNumber()(42) = %#v; want %#v", got, want)
		}
	})

	t.Run("float64-to-int", func(t *testing.T) {
		m := NumberToNumber[float64, int]()
		got := m(42)
		want := 42
		if got != want {
			t.Errorf("NumberToNumber()(42) = %#v; want %#v", got, want)
		}
	})

	t.Run("int-to-testInt", func(t *testing.T) {
		m := NumberToNumber[int, testInt]()
		got := m(42)
		want := testInt(42)
		if got != want {
			t.Errorf("NumberToNumber()(42) = %#v; want %#v", got, want)
		}
	})

	t.Run("testInt-to-int", func(t *testing.T) {
		m := NumberToNumber[testInt, int]()
		got := m(42)
		want := 42
		if got != want {
			t.Errorf("NumberToNumber()(42) = %#v; want %#v", got, want)
		}
	})
}

func TestComplexToComplex(t *testing.T) {
	t.Run("complex64-to-complex128", func(t *testing.T) {
		m := ComplexToComplex[complex64, complex128]()
		got := m(42 + 42i)
		want := 42 + 42i
		if got != want {
			t.Errorf("ComplexToComplex()(42.42+42.42i) = %#v; want %#v", got, want)
		}
	})

	t.Run("complex128-to-complex64", func(t *testing.T) {
		m := ComplexToComplex[complex128, complex64]()
		got := m(42 + 42i)
		want := complex64(42 + 42i)
		if got != want {
			t.Errorf("ComplexToComplex()(42.42+42.42i) = %#v; want %#v", got, want)
		}
	})
}

func TestBoolToNumber(t *testing.T) {
	t.Run("bool-to-int", func(t *testing.T) {
		m := BoolToNumber[bool, int]()
		got := m(true)
		want := 1
		if got != want {
			t.Errorf("BoolToNumber()(true) = %#v; want %#v", got, want)
		}
		got = m(false)
		want = 0
		if got != want {
			t.Errorf("BoolToNumber()(false) = %#v; want %#v", got, want)
		}
	})

	t.Run("bool-to-testInt", func(t *testing.T) {
		m := BoolToNumber[bool, testInt]()
		got := m(true)
		want := testInt(1)
		if got != want {
			t.Errorf("BoolToNumber()(true) = %#v; want %#v", got, want)
		}
		got = m(false)
		want = 0
		if got != want {
			t.Errorf("BoolToNumber()(false) = %#v; want %#v", got, want)
		}
	})
}

func TestNumberToBool(t *testing.T) {
	t.Run("int-to-bool", func(t *testing.T) {
		m := NumberToBool[int, bool]()
		got := m(42)
		if got != true {
			t.Errorf("NumberToBool()(42) = %#v; want %#v", got, true)
		}
		got = m(0)
		if got != false {
			t.Errorf("NumberToBool()(0) = %#v; want %#v", got, false)
		}
	})

	t.Run("testInt-to-bool", func(t *testing.T) {
		m := NumberToBool[testInt, bool]()
		got := m(42)
		if got != true {
			t.Errorf("NumberToBool()(42) = %#v; want %#v", got, true)
		}
		got = m(0)
		if got != false {
			t.Errorf("NumberToBool()(0) = %#v; want %#v", got, false)
		}
	})
}
