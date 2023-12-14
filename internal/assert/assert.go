package assert

import (
	"reflect"
	"slices"
	"testing"
)

func ElementsMatch[E comparable](t *testing.T, got, want []E) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %#v, want %#v exactly", got, want)
	}
	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("got %#v, want %#v exactly", got, want)
		}
	}
}

func ElementsMatchAnyOrder[E comparable](t *testing.T, got, want []E) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %#v, want all elements from %#v in any order", got, want)
	}
	for _, e := range got {
		if !slices.Contains(want, e) {
			t.Fatalf("got %#v, want all elements from %#v in any order", got, want)
		}
	}
}

func SomeElementsMatchAnyOrder[E comparable](t *testing.T, got, want []E, n int) {
	t.Helper()
	if len(got) != n {
		t.Fatalf("got %#v, want %d elements from %#v in any order", got, n, want)
	}
	for _, e := range got {
		if !slices.Contains(want, e) {
			t.Fatalf("got %#v, want %d elements from %#v in any order", got, n, want)
		}
	}
}
