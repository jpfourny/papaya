package stream

import (
	"github.com/jpfourny/papaya/internal/assert"
	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/optional"
	"testing"
)

func TestSortedKeyStore_get(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ks := sortedKeyStoreFactory[int, string](cmp.Natural[int]())()
		got := ks.get(0)
		want := optional.Empty[string]()
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		ks := sortedKeyStoreFactory[int, string](cmp.Natural[int]())()
		ks.put(1, "one")
		ks.put(2, "two")
		ks.put(1, "uno")
		ks.put(3, "three")
		ks.put(2, "dos")
		got := ks.get(1)
		want := optional.Of("uno")
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		got = ks.get(2)
		want = optional.Of("dos")
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		got = ks.get(3)
		want = optional.Of("three")
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		got = ks.get(4)
		want = optional.Empty[string]()
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})
}

func TestSortedKeyStore_forEach(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ks := sortedKeyStoreFactory[int, string](cmp.Natural[int]())()
		var got []int
		ks.forEach(func(key int, value string) bool {
			got = append(got, key)
			return true
		})
		if len(got) != 0 {
			t.Fatalf("got %#v, want %#v", got, []int{})
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		ks := sortedKeyStoreFactory[int, string](cmp.Natural[int]())()
		ks.put(1, "one")
		ks.put(2, "two")
		ks.put(1, "uno")
		ks.put(3, "three")
		ks.put(2, "dos")
		var got []int
		ks.forEach(func(key int, value string) bool {
			got = append(got, key)
			return true
		})
		want := []int{1, 2, 3}
		assert.ElementsMatch(t, got, want)
	})
}

func TestMapKeyStore_get(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ks := mapKeyStoreFactory[int, string]()()
		got := ks.get(0)
		want := optional.Empty[string]()
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		ks := mapKeyStoreFactory[int, string]()()
		ks.put(1, "one")
		ks.put(2, "two")
		ks.put(1, "uno")
		ks.put(3, "three")
		ks.put(2, "dos")
		got := ks.get(1)
		want := optional.Of("uno")
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		got = ks.get(2)
		want = optional.Of("dos")
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		got = ks.get(3)
		want = optional.Of("three")
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		got = ks.get(4)
		want = optional.Empty[string]()
		if got != want {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})
}

func TestMapKeyStore_forEach(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ks := mapKeyStoreFactory[int, string]()()
		var got []int
		ks.forEach(func(key int, value string) bool {
			got = append(got, key)
			return true
		})
		if len(got) != 0 {
			t.Fatalf("got %#v, want %#v", got, []int{})
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		ks := mapKeyStoreFactory[int, string]()()
		ks.put(1, "one")
		ks.put(2, "two")
		ks.put(1, "uno")
		ks.put(3, "three")
		ks.put(2, "dos")
		var got []int
		ks.forEach(func(key int, value string) bool {
			got = append(got, key)
			return true
		})
		want := []int{1, 2, 3}
		assert.ElementsMatchAnyOrder(t, got, want)
	})
}
