package mapper

import (
	"slices"
	"testing"

	"github.com/jpfourny/papaya/pkg/pointer"
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

func TestPointerRef(t *testing.T) {
	m := PointerRef[int]()
	got := m(42)
	want := pointer.Ref(42)
	if *got != *want {
		t.Errorf("*PointerRef()(42) = %#v; want %#v", *got, *want)
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

func TestSliceFrom(t *testing.T) {
	m := SliceFrom[int](1)
	got := m([]int{1, 2, 3})
	want := []int{2, 3}
	if !slices.Equal(got, want) {
		t.Errorf("SliceFrom(1)([1, 2, 3]) = %#v; want %#v", got, want)
	}
}

func TestSliceTo(t *testing.T) {
	m := SliceTo[int](2)
	got := m([]int{1, 2, 3})
	want := []int{1, 2}
	if !slices.Equal(got, want) {
		t.Errorf("SliceTo(2)([1, 2, 3]) = %#v; want %#v", got, want)
	}
}

func TestSliceFromTo(t *testing.T) {
	m := SliceFromTo[int](1, 2)
	got := m([]int{1, 2, 3})
	want := []int{2}
	if !slices.Equal(got, want) {
		t.Errorf("SliceFromTo(1, 2)([1, 2, 3]) = %#v; want %#v", got, want)
	}
}
