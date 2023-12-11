package pair

import "testing"

func TestZero(t *testing.T) {
	p := Zero[int, string]()
	if p.First() != 0 {
		t.Errorf("Zero().First() = %#v; want 0", p.First())
	}
	if p.Second() != "" {
		t.Errorf(`Zero().Second() = %#v; want ""`, p.Second())
	}
}

func TestOf(t *testing.T) {
	p := Of[int, string](42, "foo")
	if p.First() != 42 {
		t.Errorf(`Of(42, "foo").First() = %#v; want 42`, p.First())
	}
	if p.Second() != "foo" {
		t.Errorf(`Of(42, "foo").Second() = %#v; want "foo"`, p.Second())
	}
}

func TestPair_Explode(t *testing.T) {
	p := Of[int, string](42, "foo")
	if first, second := p.Explode(); first != 42 || second != "foo" {
		t.Errorf(`Of(42, "foo").Explode() = (%#v, %#v); want (42, "foo")`, first, second)
	}
}

func TestPair_Reverse(t *testing.T) {
	p := Of[int, string](42, "foo")
	if first, second := p.Reverse().Explode(); first != "foo" || second != 42 {
		t.Errorf(`Of(42, "foo").Reverse().Explode() = (%#v, %#v); want ("foo", 42)`, first, second)
	}
}

func TestPair_String(t *testing.T) {
	p := Of[int, string](42, "foo")
	if got := p.String(); got != `(42, "foo")` {
		t.Errorf(`Of(42, "foo").String() = %v; want (42, "foo")`, got)
	}
}
