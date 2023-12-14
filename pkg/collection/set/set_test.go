package set

import (
	"slices"
	"testing"

	"github.com/jpfourny/papaya/pkg/stream"
)

func TestOf(t *testing.T) {
	s := Of(1, 2)
	if !s.Contains(1) {
		t.Errorf("expected Contains(42) to be true")
	}
	if s.Contains(3) {
		t.Errorf("expected Contains(44) to be false")
	}
}

func TestMake(t *testing.T) {
	s := Make[int]()
	if !s.Empty() {
		t.Errorf("expected Empty() to be true")
	}
}

func TestMakeWithCapacity(t *testing.T) {
	s := MakeWithCapacity[int](10)
	if !s.Empty() {
		t.Errorf("expected Empty() to be true")
	}
}

func TestFromMap(t *testing.T) {
	m := map[int]struct{}{1: {}, 2: {}}
	s := FromMap(m)
	if !s.Contains(1) {
		t.Errorf("expected Contains(1) to be true")
	}
	if !s.Contains(2) {
		t.Errorf("expected Contains(2) to be true")
	}
	if s.Contains(3) {
		t.Errorf("expected Contains(3) to be false")
	}
}

func TestSet_Add(t *testing.T) {
	s := Make[int]()
	s.Add(42)
	if s.Size() != 1 {
		t.Errorf("expected Size() to be 1")
	}
	if !s.Contains(42) {
		t.Errorf("expected Contains(42) to be true")
	}
}

func TestSet_AddAll(t *testing.T) {
	s := Make[int]()
	s.AddAll(42, 43)
	if s.Size() != 2 {
		t.Errorf("expected Size() to be 2")
	}
	if !s.Contains(42) {
		t.Errorf("expected Contains(42) to be true")
	}
	if !s.Contains(43) {
		t.Errorf("expected Contains(43) to be true")
	}
}

func TestSet_Remove(t *testing.T) {
	s := Of(1, 2, 3)
	s.Remove(1)
	if s.Size() != 2 {
		t.Errorf("expected Size() to be 2")
	}
	if s.Contains(1) {
		t.Errorf("expected Contains(1) to be false")
	}
}

func TestSet_RemoveAll(t *testing.T) {
	s := Of(1, 2, 3)
	s.RemoveAll(1, 2)
	if s.Size() != 1 {
		t.Errorf("expected Size() to be 1")
	}
	if s.Contains(1) {
		t.Errorf("expected Contains(1) to be false")
	}
	if s.Contains(2) {
		t.Errorf("expected Contains(2) to be false")
	}
}

func TestSet_Contains(t *testing.T) {
	s := Of(1, 2, 3)
	if !s.Contains(1) {
		t.Errorf("expected Contains(1) to be true")
	}
	if !s.Contains(2) {
		t.Errorf("expected Contains(2) to be true")
	}
	if !s.Contains(3) {
		t.Errorf("expected Contains(3) to be true")
	}
	if s.Contains(4) {
		t.Errorf("expected Contains(4) to be false")
	}
}

func TestSet_ContainsAll(t *testing.T) {
	s := Of(1, 2, 3)
	if !s.ContainsAll(1, 2) {
		t.Errorf("expected ContainsAll(1, 2) to be true")
	}
	if !s.ContainsAll(1, 2, 3) {
		t.Errorf("expected ContainsAll(1, 2, 3) to be true")
	}
	if s.ContainsAll(1, 2, 3, 4) {
		t.Errorf("expected ContainsAll(1, 2, 3, 4) to be false")
	}
}

func TestSet_ContainsAny(t *testing.T) {
	s := Of(1, 2, 3)
	if !s.ContainsAny(1, 2) {
		t.Errorf("expected ContainsAny(1, 2) to be true")
	}
	if !s.ContainsAny(1, 2, 3) {
		t.Errorf("expected ContainsAny(1, 2, 3) to be true")
	}
	if !s.ContainsAny(1, 2, 3, 4) {
		t.Errorf("expected ContainsAny(1, 2, 3, 4) to be true")
	}
	if s.ContainsAny(4, 5) {
		t.Errorf("expected ContainsAny(4, 5) to be false")
	}
}

func TestSet_ContainsNone(t *testing.T) {
	s := Of(1, 2, 3)
	if !s.ContainsNone(4, 5) {
		t.Errorf("expected ContainsNone(4, 5) to be true")
	}
	if s.ContainsNone(1, 2) {
		t.Errorf("expected ContainsNone(1, 2) to be false")
	}
	if s.ContainsNone(1, 2, 3) {
		t.Errorf("expected ContainsNone(1, 2, 3) to be false")
	}
	if s.ContainsNone(1, 2, 3, 4) {
		t.Errorf("expected ContainsNone(1, 2, 3, 4) to be false")
	}
}

func TestSet_Equal(t *testing.T) {
	s := Of(1, 2, 3)
	other := Of(3, 1, 2)
	if !s.Equal(other) {
		t.Errorf("expected Equal() to be true")
	}

	s = Of(1, 2, 3)
	other = Of(4, 5, 6)
	if s.Equal(other) {
		t.Errorf("expected Equal() to be true")
	}

	s = Of(1, 2, 3)
	other = Of(1, 2)
	if s.Equal(other) {
		t.Errorf("expected Equal() to be false")
	}

	s = Of(1, 2)
	other = Of(1, 2, 3)
	if s.Equal(other) {
		t.Errorf("expected Equal() to be false")
	}
}

func TestSet_Union(t *testing.T) {
	s := Of(1, 2, 3)
	other := Of(2, 3, 4)
	result := s.Union(other)
	if result.Size() != 4 {
		t.Errorf("expected Size() to be 4")
	}
	if !result.Contains(1) {
		t.Errorf("expected Contains(1) to be true")
	}
	if !result.Contains(2) {
		t.Errorf("expected Contains(2) to be true")
	}
	if !result.Contains(3) {
		t.Errorf("expected Contains(3) to be true")
	}
	if !result.Contains(4) {
		t.Errorf("expected Contains(4) to be true")
	}
}

func TestSet_Intersection(t *testing.T) {
	s := Of(1, 2, 3)
	other := Of(2, 3)
	result := s.Intersection(other)
	if result.Size() != 2 {
		t.Errorf("expected Size() to be 2")
	}
	if !result.Contains(2) {
		t.Errorf("expected Contains(2) to be true")
	}
	if !result.Contains(3) {
		t.Errorf("expected Contains(3) to be true")
	}

	s = Of(2, 3)
	other = Of(1, 2, 3)
	result = s.Intersection(other)
	if result.Size() != 2 {
		t.Errorf("expected Size() to be 2")
	}
	if !result.Contains(2) {
		t.Errorf("expected Contains(2) to be true")
	}
	if !result.Contains(3) {
		t.Errorf("expected Contains(3) to be true")
	}
}

func TestSet_Difference(t *testing.T) {
	s := Of(1, 2, 3)
	other := Of(2, 3)
	result := s.Difference(other)
	if result.Size() != 1 {
		t.Errorf("expected Size() to be 1")
	}
	if !result.Contains(1) {
		t.Errorf("expected Contains(1) to be true")
	}

	s = Of(2, 3)
	other = Of(1, 2, 3)
	result = s.Difference(other)
	if result.Size() != 0 {
		t.Errorf("expected Size() to be 0")
	}
}

func TestSet_Clear(t *testing.T) {
	s := Of(1, 2, 3)
	s.Clear()
	if !s.Empty() {
		t.Errorf("expected Empty() to be true")
	}
}

func TestSet_Stream(t *testing.T) {
	s := Of(1, 2, 3)
	st := s.Stream()
	if stream.Count(st) != 3 {
		t.Errorf("expected Count() to be 3")
	}
	if !stream.ContainsAll(st, 1, 2, 3) {
		t.Errorf("expected all elements to be between 1 and 3")
	}
}

func TestSet_String(t *testing.T) {
	s := Make[int]()
	if s.String() != "[]" {
		t.Errorf("expected String() to be []")
	}

	s = Of(1, 2, 3)
	if s.String() != "[1, 2, 3]" {
		t.Errorf("expected String() to be [1, 2, 3]")
	}
}

func TestSet_ToSlice(t *testing.T) {
	s := Make[int]()
	if len(s.ToSlice()) != 0 {
		t.Errorf("expected len(ToSlice()) to be 0")
	}

	s = Of(1, 2, 3)
	sl := s.ToSlice()
	if len(sl) != 3 {
		t.Errorf("expected len(ToSlice()) to be 3")
	}
	if !slices.Contains(sl, 1) {
		t.Errorf("expected Contains(1) to be true")
	}
	if !slices.Contains(sl, 2) {
		t.Errorf("expected Contains(2) to be true")
	}
	if !slices.Contains(sl, 3) {
		t.Errorf("expected Contains(3) to be true")
	}
}
