package pred

import "testing"

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
