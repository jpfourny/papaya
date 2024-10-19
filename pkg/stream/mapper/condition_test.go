package mapper

import (
	"github.com/jpfourny/papaya/v2/pkg/stream/pred"
	"testing"

	"github.com/jpfourny/papaya/v2/pkg/opt"
)

func TestIf(t *testing.T) {
	m := If(pred.GreaterThan(0), Constant[int](1))
	got := m(42)
	want := opt.Of(1)
	if got != want {
		t.Errorf("If(GreaterThan(0), Constant(1))(42) = %#v; want %#v", got, want)
	}

	got = m(-42)
	want = opt.Empty[int]()
	if got != want {
		t.Errorf("If(GreaterThan(0), Constant(1))(-42) = %#v; want %#v", got, want)
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

func TestSwitch(t *testing.T) {
	m := Switch[int, string](
		[]Case[int, string]{
			{Cond: pred.GreaterThan(0), Mapper: Constant[int]("positive")},
			{Cond: pred.LessThan(0), Mapper: Constant[int]("negative")},
		},
	)
	got := m(-1)
	want := opt.Of("negative")
	if got != want {
		t.Errorf("Switch(..)(-1) = %#v; want %#v", got, want)
	}

	got = m(0)
	want = opt.Empty[string]()
	if got != want {
		t.Errorf("Switch(..)(0) = %#v; want %#v", got, want)
	}

	got = m(1)
	want = opt.Of("positive")
	if got != want {
		t.Errorf("Switch(..)(1) = %#v; want %#v", got, want)
	}
}

func TestSwitchWithDefault(t *testing.T) {
	m := SwitchWithDefault[int, string](
		[]Case[int, string]{
			{Cond: pred.GreaterThan(0), Mapper: Constant[int]("positive")},
			{Cond: pred.LessThan(0), Mapper: Constant[int]("negative")},
		},
		Constant[int]("neutral"), // Default case.
	)
	if m(-1) != "negative" {
		t.Errorf("Switch(..)(-1) = %#v; want %#v", m(-1), "negative")
	}
	if m(0) != "neutral" {
		t.Errorf("Switch(..)(0) = %#v; want %#v", m(0), "neutral")
	}
	if m(1) != "positive" {
		t.Errorf("Switch(..)(1) = %#v; want %#v", m(1), "positive")
	}
}
