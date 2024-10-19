package pred

import (
	"testing"

	"github.com/jpfourny/papaya/v2/pkg/pair"
	"github.com/jpfourny/papaya/v2/pkg/ptr"
)

func TestNil(t *testing.T) {
	p := Nil[int]()
	if !p(nil) {
		t.Errorf("Nil()(nil) = false; want true")
	}
	if p(ptr.Ref(0)) {
		t.Errorf("Nil()(0) = true; want false")
	}
}

func TestNotNil(t *testing.T) {
	p := NotNil[int]()
	if p(nil) {
		t.Errorf("NotNil()(nil) = true; want false")
	}
	if !p(ptr.Ref(0)) {
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
