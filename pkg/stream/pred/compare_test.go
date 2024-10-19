package pred

import (
	"testing"

	"github.com/jpfourny/papaya/v2/pkg/ptr"
)

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

	p2 := DeepEqual(ptr.Ref(0))
	if !p2(ptr.Ref(0)) {
		t.Errorf("DeepEqual(&0)(&0) = false; want true")
	}
	if p2(ptr.Ref(1)) {
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

	p2 := NotDeepEqual(ptr.Ref(0))
	if p2(ptr.Ref(0)) {
		t.Errorf("NotDeepEqual(&0)(&0) = true; want false")
	}
	if !p2(ptr.Ref(1)) {
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
