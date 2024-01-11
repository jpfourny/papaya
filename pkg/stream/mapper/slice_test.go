package mapper

import (
	"slices"
	"testing"
)

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
