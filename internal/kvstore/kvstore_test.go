package kvstore

import (
	"github.com/jpfourny/papaya/v2/internal/assert"
	"github.com/jpfourny/papaya/v2/pkg/cmp"
	"github.com/jpfourny/papaya/v2/pkg/opt"
	"testing"
)

func TestMappedMaker(t *testing.T) {
	m := MappedMaker[int, string]()
	if m == nil {
		t.Fatalf("got %#v, want non-nil", m)
	}
	ks := m()
	if ks == nil {
		t.Fatalf("got %#v, want non-nil", ks)
	}
}

func TestOrderedMaker(t *testing.T) {
	m := SortedMaker[int, string](cmp.Natural[int]())
	if m == nil {
		t.Fatalf("got %#v, want non-nil", m)
	}
	ks := m()
	if ks == nil {
		t.Fatalf("got %#v, want non-nil", ks)
	}
}

func TestOrderedStore_Get(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ks := NewSorted[int, string](cmp.Natural[int]())
		got := ks.Get(0)
		want := opt.Empty[string]()
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		ks := NewSorted[int, string](cmp.Natural[int]())
		ks.Put(1, "one")
		ks.Put(2, "two")
		ks.Put(1, "uno")
		ks.Put(3, "three")
		ks.Put(2, "dos")
		got := ks.Get(1)
		want := opt.Of("uno")
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		got = ks.Get(2)
		want = opt.Of("dos")
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		got = ks.Get(3)
		want = opt.Of("three")
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		got = ks.Get(4)
		want = opt.Empty[string]()
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})
}

func TestOrderedStore_ForEach(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ks := NewSorted[int, string](cmp.Natural[int]())
		var got []int
		ks.ForEach(func(key int, value string) bool {
			got = append(got, key)
			return true
		})
		if len(got) != 0 {
			t.Fatalf("got %#v, want %#v", got, []int{})
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		ks := NewSorted[int, string](cmp.Natural[int]())
		ks.Put(1, "one")
		ks.Put(2, "two")
		ks.Put(1, "uno")
		ks.Put(3, "three")
		ks.Put(2, "dos")
		var got []int
		ks.ForEach(func(key int, value string) bool {
			got = append(got, key)
			return true
		})
		want := []int{1, 2, 3}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("stop", func(t *testing.T) {
		ks := NewSorted[int, string](cmp.Natural[int]())
		ks.Put(1, "one")
		ks.Put(2, "two")
		ks.Put(1, "uno")
		ks.Put(3, "three")
		ks.Put(2, "dos")
		var got []int
		ks.ForEach(func(key int, value string) bool {
			got = append(got, key)
			return false
		})
		want := []int{1}
		assert.ElementsMatch(t, got, want)
	})
}

func TestMappedStore_Get(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ks := NewMapped[int, string]()
		got := ks.Get(0)
		want := opt.Empty[string]()
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		ks := NewMapped[int, string]()
		ks.Put(1, "one")
		ks.Put(2, "two")
		ks.Put(1, "uno")
		ks.Put(3, "three")
		ks.Put(2, "dos")
		got := ks.Get(1)
		want := opt.Of("uno")
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		got = ks.Get(2)
		want = opt.Of("dos")
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		got = ks.Get(3)
		want = opt.Of("three")
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		got = ks.Get(4)
		want = opt.Empty[string]()
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})
}

func TestMappedStore_ForEach(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ks := NewMapped[int, string]()
		var got []int
		ks.ForEach(func(key int, value string) bool {
			got = append(got, key)
			return true
		})
		if len(got) != 0 {
			t.Fatalf("got %#v, want %#v", got, []int{})
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		ks := NewMapped[int, string]()
		ks.Put(1, "one")
		ks.Put(2, "two")
		ks.Put(1, "uno")
		ks.Put(3, "three")
		ks.Put(2, "dos")
		var got []int
		ks.ForEach(func(key int, value string) bool {
			got = append(got, key)
			return true
		})
		want := []int{1, 2, 3}
		assert.ElementsMatchAnyOrder(t, got, want)
	})

	t.Run("stop", func(t *testing.T) {
		ks := NewMapped[int, string]()
		ks.Put(1, "one")
		ks.Put(2, "two")
		ks.Put(1, "uno")
		ks.Put(3, "three")
		ks.Put(2, "dos")
		var got []int
		ks.ForEach(func(key int, value string) bool {
			got = append(got, key)
			return false
		})
		if len(got) != 1 {
			t.Fatalf("expected exactly 1 element, got %d", len(got))
		}
	})
}
