package pred

import (
	"github.com/jpfourny/papaya/pkg/pair"
	"github.com/jpfourny/papaya/pkg/pointer"
	"testing"
)

func TestTrue(t *testing.T) {
	p := True[int]()
	if !p(0) {
		t.Errorf("True()(0) = false; want true")
	}
}

func TestFalse(t *testing.T) {
	p := False[int]()
	if p(0) {
		t.Errorf("False()(0) = true; want false")
	}
}

func TestNot(t *testing.T) {
	p := Not(True[int]())
	if p(0) {
		t.Errorf("Not(True())(0) = true; want false")
	}

	p = Not(False[int]())
	if !p(0) {
		t.Errorf("Not(False())(0) = false; want true")
	}
}

func TestAnd(t *testing.T) {
	p := And(True[int](), True[int]())
	if !p(0) {
		t.Errorf("And(True(), True())(0) = false; want true")
	}

	p = And(True[int](), False[int]())
	if p(0) {
		t.Errorf("And(True(), False())(0) = true; want false")
	}

	p = And(False[int](), True[int]())
	if p(0) {
		t.Errorf("And(False(), True())(0) = true; want false")
	}

	p = And(False[int](), False[int]())
	if p(0) {
		t.Errorf("And(False(), False())(0) = true; want false")
	}
}

func TestOr(t *testing.T) {
	p := Or(True[int](), True[int]())
	if !p(0) {
		t.Errorf("Or(True(), True())(0) = false; want true")
	}

	p = Or(True[int](), False[int]())
	if !p(0) {
		t.Errorf("Or(True(), False())(0) = false; want true")
	}

	p = Or(False[int](), True[int]())
	if !p(0) {
		t.Errorf("Or(False(), True())(0) = false; want true")
	}

	p = Or(False[int](), False[int]())
	if p(0) {
		t.Errorf("Or(False(), False())(0) = true; want false")
	}
}

func TestEqual(t *testing.T) {
	p := Equal(0)
	if !p(0) {
		t.Errorf("Equal(0)(0) = false; want true")
	}
	if p(1) {
		t.Errorf("Equal(0)(1) = true; want false")
	}
}

func TestNotEqual(t *testing.T) {
	p := NotEqual(0)
	if p(0) {
		t.Errorf("NotEqual(0)(0) = true; want false")
	}
	if !p(1) {
		t.Errorf("NotEqual(0)(1) = false; want true")
	}
}

func TestEqualBy(t *testing.T) {
	p := EqualBy(0, func(a, b int) int {
		return a - b
	})
	if !p(0) {
		t.Errorf("EqualBy(0)(0) = false; want true")
	}
	if p(1) {
		t.Errorf("EqualBy(0)(1) = true; want false")
	}
}

func TestNotEqualBy(t *testing.T) {
	p := NotEqualBy(0, func(a, b int) int {
		return a - b
	})
	if p(0) {
		t.Errorf("NotEqualBy(0)(0) = true; want false")
	}
	if !p(1) {
		t.Errorf("NotEqualBy(0)(1) = false; want true")
	}
}

func TestDeepEqual(t *testing.T) {
	p := DeepEqual(0)
	if !p(0) {
		t.Errorf("DeepEqual(0)(0) = false; want true")
	}
	if p(1) {
		t.Errorf("DeepEqual(0)(1) = true; want false")
	}

	p2 := DeepEqual(pointer.Ref(0))
	if !p2(pointer.Ref(0)) {
		t.Errorf("DeepEqual(&0)(&0) = false; want true")
	}
	if p2(pointer.Ref(1)) {
		t.Errorf("DeepEqual(&0)(&1) = true; want false")
	}
}

func TestNotDeepEqual(t *testing.T) {
	p := NotDeepEqual(0)
	if p(0) {
		t.Errorf("NotDeepEqual(0)(0) = true; want false")
	}
	if !p(1) {
		t.Errorf("NotDeepEqual(0)(1) = false; want true")
	}

	p2 := NotDeepEqual(pointer.Ref(0))
	if p2(pointer.Ref(0)) {
		t.Errorf("NotDeepEqual(&0)(&0) = true; want false")
	}
	if !p2(pointer.Ref(1)) {
		t.Errorf("NotDeepEqual(&0)(&1) = false; want true")
	}
}

func TestLessThan(t *testing.T) {
	p := LessThan(0)
	if p(0) {
		t.Errorf("LessThan(0)(0) = true; want false")
	}
	if !p(-1) {
		t.Errorf("LessThan(0)(-1) = false; want true")
	}
	if p(1) {
		t.Errorf("LessThan(0)(1) = true; want false")
	}
}

func TestLessThanOrEqual(t *testing.T) {
	p := LessThanOrEqual(0)
	if !p(0) {
		t.Errorf("LessThanOrEqual(0)(0) = false; want true")
	}
	if !p(-1) {
		t.Errorf("LessThanOrEqual(0)(-1) = false; want true")
	}
	if p(1) {
		t.Errorf("LessThanOrEqual(0)(1) = true; want false")
	}
}

func TestGreaterThan(t *testing.T) {
	p := GreaterThan(0)
	if p(0) {
		t.Errorf("GreaterThan(0)(0) = true; want false")
	}
	if p(-1) {
		t.Errorf("GreaterThan(0)(-1) = true; want false")
	}
	if !p(1) {
		t.Errorf("GreaterThan(0)(1) = false; want true")
	}
}

func TestGreaterThanOrEqual(t *testing.T) {
	p := GreaterThanOrEqual(0)
	if !p(0) {
		t.Errorf("GreaterThanOrEqual(0)(0) = false; want true")
	}
	if p(-1) {
		t.Errorf("GreaterThanOrEqual(0)(-1) = true; want false")
	}
	if !p(1) {
		t.Errorf("GreaterThanOrEqual(0)(1) = false; want true")
	}
}

func TestLessThanBy(t *testing.T) {
	p := LessThanBy(0, func(a, b int) int {
		return a - b
	})
	if p(0) {
		t.Errorf("LessThanBy(0)(0) = true; want false")
	}
	if !p(-1) {
		t.Errorf("LessThanBy(0)(-1) = false; want true")
	}
	if p(1) {
		t.Errorf("LessThanBy(0)(1) = true; want false")
	}
}

func TestLessThanOrEqualBy(t *testing.T) {
	p := LessThanOrEqualBy(0, func(a, b int) int {
		return a - b
	})
	if !p(0) {
		t.Errorf("LessThanOrEqualBy(0)(0) = false; want true")
	}
	if !p(-1) {
		t.Errorf("LessThanOrEqualBy(0)(-1) = false; want true")
	}
	if p(1) {
		t.Errorf("LessThanOrEqualBy(0)(1) = true; want false")
	}
}

func TestGreaterThanBy(t *testing.T) {
	p := GreaterThanBy(0, func(a, b int) int {
		return a - b
	})
	if p(0) {
		t.Errorf("GreaterThanBy(0)(0) = true; want false")
	}
	if p(-1) {
		t.Errorf("GreaterThanBy(0)(-1) = true; want false")
	}
	if !p(1) {
		t.Errorf("GreaterThanBy(0)(1) = false; want true")
	}
}

func TestGreaterThanOrEqualBy(t *testing.T) {
	p := GreaterThanOrEqualBy(0, func(a, b int) int {
		return a - b
	})
	if !p(0) {
		t.Errorf("GreaterThanOrEqualBy(0)(0) = false; want true")
	}
	if p(-1) {
		t.Errorf("GreaterThanOrEqualBy(0)(-1) = true; want false")
	}
	if !p(1) {
		t.Errorf("GreaterThanOrEqualBy(0)(1) = false; want true")
	}
}

func TestRoughlyEqual(t *testing.T) {
	t.Run("far-from-epsilon", func(t *testing.T) {
		p := RoughlyEqual(5.0, 0.1)
		if !p(5.01) {
			t.Errorf("RoughlyEqual(5.0, 0.1)(5.01) = false; want true")
		}
		if !p(5.09999) {
			t.Errorf("RoughlyEqual(5.0, 0.09999)(5.09999) = false; want true")
		}
		if !p(5.1) {
			t.Errorf("RoughlyEqual(5.0, 0.1)(5.1) = false; want true")
		}
		if p(5.10001) {
			t.Errorf("RoughlyEqual(5.0, 0.1)(5.10001) = true; want false")
		}
	})

	t.Run("close-to-epsilon", func(t *testing.T) {
		p := RoughlyEqual(0, 0.001)
		if !p(0) {
			t.Errorf("RoughlyEqual(0, 0.001)(0) = false; want true")
		}
		if !p(0.0009) {
			t.Errorf("RoughlyEqual(0, 0.001)(0.0009) = false; want true")
		}
		if !p(0.00099999999) {
			t.Errorf("RoughlyEqual(0, 0.001)(00099999999) = false; want true")
		}
		if p(0.1) {
			t.Errorf("RoughlyEqual(0, 0.001)(0.1) = true; want false")
		}
	})
}

func TestNotRoughlyEqual(t *testing.T) {
	t.Run("far-from-epsilon", func(t *testing.T) {
		p := NotRoughlyEqual(5.0, 0.1)
		if p(5.01) {
			t.Errorf("NotRoughlyEqual(5.0, 0.1)(5.01) = true; want false")
		}
		if p(5.09999) {
			t.Errorf("NotRoughlyEqual(5.0, 0.09999)(5.09999) = true; want false")
		}
		if p(5.1) {
			t.Errorf("NotRoughlyEqual(5.0, 0.1)(5.1) = true; want false")
		}
		if !p(5.10001) {
			t.Errorf("NotRoughlyEqual(5.0, 0.1)(5.10001) = false; want true")
		}
	})

	t.Run("close-to-epsilon", func(t *testing.T) {
		p := NotRoughlyEqual(0, 0.001)
		if p(0) {
			t.Errorf("NotRoughlyEqual(0, 0.001)(0) = true; want false")
		}
		if p(0.0009) {
			t.Errorf("NotRoughlyEqual(0, 0.001)(0.0009) = true; want false")
		}
		if p(0.00099999999) {
			t.Errorf("NotRoughlyEqual(0, 0.001)(00099999999) = true; want false")
		}
		if !p(0.1) {
			t.Errorf("NotRoughlyEqual(0, 0.001)(0.1) = false; want true")
		}
	})
}

func TestNil(t *testing.T) {
	p := Nil[int]()
	if !p(nil) {
		t.Errorf("Nil()(nil) = false; want true")
	}
	if p(pointer.Ref(0)) {
		t.Errorf("Nil()(0) = true; want false")
	}
}

func TestNotNil(t *testing.T) {
	p := NotNil[int]()
	if p(nil) {
		t.Errorf("NotNil()(nil) = true; want false")
	}
	if !p(pointer.Ref(0)) {
		t.Errorf("NotNil()(0) = false; want true")
	}
}

func TestZero(t *testing.T) {
	p := Zero[int]()
	if !p(0) {
		t.Errorf("Zero()(0) = false; want true")
	}
	if p(1) {
		t.Errorf("Zero()(1) = true; want false")
	}

	p2 := Zero[pair.Pair[int, int]]()
	if !p2(pair.Of(0, 0)) {
		t.Errorf("Zero()(Pair{0, 0}) = false; want true")
	}
	if p2(pair.Of(0, 1)) {
		t.Errorf("Zero()(Pair{0, 1}) = true; want false")
	}
}

func TestNotZero(t *testing.T) {
	p := NotZero[int]()
	if p(0) {
		t.Errorf("NotZero()(0) = true; want false")
	}
	if !p(1) {
		t.Errorf("NotZero()(1) = false; want true")
	}

	p2 := NotZero[pair.Pair[int, int]]()
	if p2(pair.Of(0, 0)) {
		t.Errorf("NotZero()(Pair{0, 0}) = true; want false")
	}
	if !p2(pair.Of(0, 1)) {
		t.Errorf("NotZero()(Pair{0, 1}) = false; want true")
	}
}

func TestIn(t *testing.T) {
	p := In(0, 1, 2)
	if !p(0) {
		t.Errorf("In(0, 1, 2)(0) = false; want true")
	}
	if !p(1) {
		t.Errorf("In(0, 1, 2)(1) = false; want true")
	}
	if !p(2) {
		t.Errorf("In(0, 1, 2)(2) = false; want true")
	}
	if p(3) {
		t.Errorf("In(0, 1, 2)(3) = true; want false")
	}
}

func TestNotIn(t *testing.T) {
	p := NotIn(0, 1, 2)
	if p(0) {
		t.Errorf("NotIn(0, 1, 2)(0) = true; want false")
	}
	if p(1) {
		t.Errorf("NotIn(0, 1, 2)(1) = true; want false")
	}
	if p(2) {
		t.Errorf("NotIn(0, 1, 2)(2) = true; want false")
	}
	if !p(3) {
		t.Errorf("NotIn(0, 1, 2)(3) = false; want true")
	}
}

func TestInSlice(t *testing.T) {
	p := InSlice[int]([]int{0, 1, 2})
	if !p(0) {
		t.Errorf("InSlice([]int{0, 1, 2})(0) = false; want true")
	}
	if !p(1) {
		t.Errorf("InSlice([]int{0, 1, 2})(1) = false; want true")
	}
	if !p(2) {
		t.Errorf("InSlice([]int{0, 1, 2})(2) = false; want true")
	}
	if p(3) {
		t.Errorf("InSlice([]int{0, 1, 2})(3) = true; want false")
	}
}

func TestNotInSlice(t *testing.T) {
	p := NotInSlice[int]([]int{0, 1, 2})
	if p(0) {
		t.Errorf("NotInSlice([]int{0, 1, 2})(0) = true; want false")
	}
	if p(1) {
		t.Errorf("NotInSlice([]int{0, 1, 2})(1) = true; want false")
	}
	if p(2) {
		t.Errorf("NotInSlice([]int{0, 1, 2})(2) = true; want false")
	}
	if !p(3) {
		t.Errorf("NotInSlice([]int{0, 1, 2})(3) = false; want true")
	}
}

func TestInBy(t *testing.T) {
	p := InBy[int](func(a, b int) int {
		return a - b
	}, 0, 1, 2)
	if !p(0) {
		t.Errorf("InBy(cmp, 0, 1, 2)(0) = false; want true")
	}
	if !p(1) {
		t.Errorf("InBy(cmp, 0, 1, 2)(1) = false; want true")
	}
	if !p(2) {
		t.Errorf("InBy(cmp, 0, 1, 2)(2) = false; want true")
	}
	if p(3) {
		t.Errorf("InBy(cmp, 0, 1, 2)(3) = true; want false")
	}
}

func TestNotInBy(t *testing.T) {
	p := NotInBy[int](func(a, b int) int {
		return a - b
	}, 0, 1, 2)
	if p(0) {
		t.Errorf("NotInBy(cmp, 0, 1, 2)(0) = true; want false")
	}
	if p(1) {
		t.Errorf("NotInBy(cmp, 0, 1, 2)(1) = true; want false")
	}
	if p(2) {
		t.Errorf("NotInBy(cmp, 0, 1, 2)(2) = true; want false")
	}
	if !p(3) {
		t.Errorf("NotInBy(cmp, 0, 1, 2)(3) = false; want true")
	}
}

func TestInSliceBy(t *testing.T) {
	p := InSliceBy[int]([]int{0, 1, 2}, func(a, b int) int {
		return a - b
	})
	if !p(0) {
		t.Errorf("InSliceBy([]int{0, 1, 2}, cmp)(0) = false; want true")
	}
	if !p(1) {
		t.Errorf("InSliceBy([]int{0, 1, 2}, cmp)(1) = false; want true")
	}
	if !p(2) {
		t.Errorf("InSliceBy([]int{0, 1, 2}, cmp)(2) = false; want true")
	}
	if p(3) {
		t.Errorf("InSliceBy([]int{0, 1, 2}, cmp)(3) = true; want false")
	}
}

func TestNotInSliceBy(t *testing.T) {
	p := NotInSliceBy[int]([]int{0, 1, 2}, func(a, b int) int {
		return a - b
	})
	if p(0) {
		t.Errorf("NotInSliceBy([]int{0, 1, 2}, cmp)(0) = true; want false")
	}
	if p(1) {
		t.Errorf("NotInSliceBy([]int{0, 1, 2}, cmp)(1) = true; want false")
	}
	if p(2) {
		t.Errorf("NotInSliceBy([]int{0, 1, 2}, cmp)(2) = true; want false")
	}
	if !p(3) {
		t.Errorf("NotInSliceBy([]int{0, 1, 2}, cmp)(3) = false; want true")
	}
}

func TestInSet(t *testing.T) {
	set := map[int]struct{}{0: {}, 1: {}, 2: {}}
	p := InSet[int](set)
	if !p(0) {
		t.Errorf("InSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(0) = false; want true")
	}
	if !p(1) {
		t.Errorf("InSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(1) = false; want true")
	}
	if !p(2) {
		t.Errorf("InSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(2) = false; want true")
	}
	if p(3) {
		t.Errorf("InSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(3) = true; want false")
	}
}

func TestNotInSet(t *testing.T) {
	set := map[int]struct{}{0: {}, 1: {}, 2: {}}
	p := NotInSet[int](set)
	if p(0) {
		t.Errorf("NotInSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(0) = true; want false")
	}
	if p(1) {
		t.Errorf("NotInSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(1) = true; want false")
	}
	if p(2) {
		t.Errorf("NotInSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(2) = true; want false")
	}
	if !p(3) {
		t.Errorf("NotInSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(3) = false; want true")
	}
}

func TestAllOf(t *testing.T) {
	p := AllOf(True[int](), True[int]())
	if !p(0) {
		t.Errorf("AllOf(True(), True())(0) = false; want true")
	}

	p = AllOf(True[int](), False[int]())
	if p(0) {
		t.Errorf("AllOf(True(), False())(0) = true; want false")
	}

	p = AllOf(False[int](), True[int]())
	if p(0) {
		t.Errorf("AllOf(False(), True())(0) = true; want false")
	}

	p = AllOf(False[int](), False[int]())
	if p(0) {
		t.Errorf("AllOf(False(), False())(0) = true; want false")
	}
}

func TestAnyOf(t *testing.T) {
	p := AnyOf[int]()
	if p(0) {
		t.Errorf("AnyOf()(0) = true; want false")
	}

	p = AnyOf(True[int]())
	if !p(0) {
		t.Errorf("AnyOf(True())(0) = false; want true")
	}

	p = AnyOf(False[int]())
	if p(0) {
		t.Errorf("AnyOf(False())(0) = true; want false")
	}

	p = AnyOf(True[int](), True[int]())
	if !p(0) {
		t.Errorf("AnyOf(True(), True())(0) = false; want true")
	}

	p = AnyOf(True[int](), False[int]())
	if !p(0) {
		t.Errorf("AnyOf(True(), False())(0) = false; want true")
	}

	p = AnyOf(False[int](), True[int]())
	if !p(0) {
		t.Errorf("AnyOf(False(), True())(0) = false; want true")
	}

	p = AnyOf(False[int](), False[int]())
	if p(0) {
		t.Errorf("AnyOf(False(), False())(0) = true; want false")
	}
}

func TestNoneOf(t *testing.T) {
	p := NoneOf[int]()
	if !p(0) {
		t.Errorf("NoneOf()(0) = false; want true")
	}

	p = NoneOf(True[int]())
	if p(0) {
		t.Errorf("NoneOf(True())(0) = true; want false")
	}

	p = NoneOf(False[int]())
	if !p(0) {
		t.Errorf("NoneOf(False())(0) = false; want true")
	}

	p = NoneOf(True[int](), True[int]())
	if p(0) {
		t.Errorf("NoneOf(True(), True())(0) = true; want false")
	}

	p = NoneOf(True[int](), False[int]())
	if p(0) {
		t.Errorf("NoneOf(True(), False())(0) = true; want false")
	}

	p = NoneOf(False[int](), True[int]())
	if p(0) {
		t.Errorf("NoneOf(False(), True())(0) = true; want false")
	}

	p = NoneOf(False[int](), False[int]())
	if !p(0) {
		t.Errorf("NoneOf(False(), False())(0) = false; want true")
	}
}

func TestOneOf(t *testing.T) {
	p := OneOf[int]()
	if p(0) {
		t.Errorf("OneOf()(0) = true; want false")
	}

	p = OneOf(True[int]())
	if !p(0) {
		t.Errorf("OneOf(True()(0) = false; want true")
	}

	p = OneOf(False[int]())
	if p(0) {
		t.Errorf("OneOf(False()(0) = true; want false")
	}

	p = OneOf(True[int](), True[int]())
	if p(0) {
		t.Errorf("OneOf(True(), True())(0) = true; want false")
	}

	p = OneOf(True[int](), False[int]())
	if !p(0) {
		t.Errorf("OneOf(True(), False())(0) = false; want true")
	}

	p = OneOf(False[int](), True[int]())
	if !p(0) {
		t.Errorf("OneOf(False(), True())(0) = false; want true")
	}

	p = OneOf(False[int](), False[int]())
	if p(0) {
		t.Errorf("OneOf(False(), False())(0) = true; want false")
	}
}
