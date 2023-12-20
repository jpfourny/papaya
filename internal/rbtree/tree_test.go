package rbtree

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/jpfourny/papaya/pkg/cmp"
)

func TestTree_Put(t *testing.T) {
	tree := New[int, string](cmp.Natural[int]())
	if !tree.Put(1, "foo") {
		t.Errorf("Expected no replacement")
	}
	if !tree.Put(2, "bar") {
		t.Errorf("Expected no replacement")
	}
	if !tree.Put(3, "baz") {
		t.Errorf("Expected no replacement")
	}
	if tree.Put(1, "qux") {
		t.Errorf("Expected replacement")
	}
}

func TestTree_Contains(t *testing.T) {
	tree := New[int, string](cmp.Natural[int]())
	tree.Put(1, "foo")
	tree.Put(2, "bar")
	tree.Put(3, "baz")

	if !tree.Contains(1) {
		t.Errorf("Expected to find key 1")
	}
	if !tree.Contains(2) {
		t.Errorf("Expected to find key 2")
	}
	if !tree.Contains(3) {
		t.Errorf("Expected to find key 3")
	}
	if tree.Contains(4) {
		t.Errorf("Expected not to find key 4")
	}
}

func TestTree_Get(t *testing.T) {
	tree := New[int, string](cmp.Natural[int]())
	tree.Put(1, "foo")
	tree.Put(2, "bar")
	tree.Put(3, "baz")

	want := "foo"
	got, ok := tree.Get(1)
	if !ok {
		t.Errorf("Expected to find key 1")
	}
	if got != want {
		t.Errorf("Expected value %s, got %s", want, got)
	}

	want = "bar"
	got, ok = tree.Get(2)
	if !ok {
		t.Errorf("Expected to find key 2")
	}
	if got != want {
		t.Errorf("Expected value %s, got %s", want, got)
	}

	want = "baz"
	got, ok = tree.Get(3)
	if !ok {
		t.Errorf("Expected to find key 3")
	}
	if got != want {
		t.Errorf("Expected value %s, got %s", want, got)
	}

	got, ok = tree.Get(4)
	if ok {
		t.Errorf("Expected not to find key 4")
	}
	if got != "" {
		t.Errorf("Expected empty value, got %s", got)
	}
}

func TestTree_Size(t *testing.T) {
	tree := New[int, string](cmp.Natural[int]())
	want := 0
	got := tree.Size()
	if got != want {
		t.Errorf("Expected size %d, got %d", want, got)
	}

	tree.Put(1, "foo")
	want = 1
	got = tree.Size()
	if got != want {
		t.Errorf("Expected size %d, got %d", want, got)
	}

	tree.Put(2, "bar")
	want = 2
	got = tree.Size()
	if got != want {
		t.Errorf("Expected size %d, got %d", want, got)
	}

	tree.Put(3, "baz")
	want = 3
	got = tree.Size()
	if got != want {
		t.Errorf("Expected size %d, got %d", want, got)
	}
}

func TestTree_Delete(t *testing.T) {
	tree := New[int, string](cmp.Natural[int]())
	tree.Put(1, "foo")
	tree.Put(2, "bar")
	tree.Put(3, "baz")

	if !tree.Delete(1) {
		t.Errorf("Expected to remove key 1")
	}
	if tree.Delete(1) {
		t.Errorf("Expected not to remove key 1")
	}
	if !tree.Delete(2) {
		t.Errorf("Expected to remove key 2")
	}
	if tree.Delete(2) {
		t.Errorf("Expected not to remove key 2")
	}
	if !tree.Delete(3) {
		t.Errorf("Expected to remove key 3")
	}
	if tree.Delete(3) {
		t.Errorf("Expected not to remove key 3")
	}
}

func TestTree_Clear(t *testing.T) {
	tree := New[int, string](cmp.Natural[int]())
	tree.Put(1, "foo")
	tree.Put(2, "bar")
	tree.Put(3, "baz")

	tree.Clear()
	if tree.Size() != 0 {
		t.Errorf("Expected size 0, got %d", tree.Size())
	}
}

func TestTree_ForEach(t *testing.T) {
	tree := New[int, string](cmp.Natural[int]())
	tree.Put(1, "foo")
	tree.Put(2, "bar")
	tree.Put(3, "baz")

	t.Run("all", func(t *testing.T) {
		var gotKeys []int
		var gotValues []string
		tree.ForEach(func(k int, v string) bool {
			gotKeys = append(gotKeys, k)
			gotValues = append(gotValues, v)
			return true
		})

		wantKeys := []int{1, 2, 3}
		wantValues := []string{"foo", "bar", "baz"}
		if !reflect.DeepEqual(gotKeys, wantKeys) {
			t.Errorf("Expected keys %v, got %v", wantKeys, gotKeys)
		}
		if !reflect.DeepEqual(gotValues, wantValues) {
			t.Errorf("Expected values %v, got %v", wantValues, gotValues)
		}
	})

	t.Run("limited", func(t *testing.T) {
		var gotKeys []int
		var gotValues []string
		tree.ForEach(func(k int, v string) bool {
			gotKeys = append(gotKeys, k)
			gotValues = append(gotValues, v)
			return len(gotKeys) < 2
		})

		wantKeys := []int{1, 2}
		wantValues := []string{"foo", "bar"}
		if !reflect.DeepEqual(gotKeys, wantKeys) {
			t.Errorf("Expected keys %v, got %v", wantKeys, gotKeys)
		}
		if !reflect.DeepEqual(gotValues, wantValues) {
			t.Errorf("Expected values %v, got %v", wantValues, gotValues)
		}
	})
}

func TestTree_mixed(t *testing.T) {
	tree := New[int, string](cmp.Natural[int]())

	// Generate key-value pairs for insertion.
	var keys []int
	var values []string
	for i := 0; i < 100000; i++ {
		keys = append(keys, i)
		values = append(values, fmt.Sprintf("value-%d", i))
	}

	// Shuffle the order of pairs to insert.
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
		values[i], values[j] = values[j], values[i]
	})

	// Insert the pairs.
	for i := 0; i < len(keys); i++ {
		if !tree.Put(keys[i], values[i]) {
			t.Fatalf("Expected no replacement")
		}
	}

	// Check that all pairs are present.
	if tree.Size() != len(keys) {
		t.Fatalf("Expected size %d, got %d", len(keys), tree.Size())
	}
	for i := 0; i < len(keys); i++ {
		if !tree.Contains(keys[i]) {
			t.Fatalf("Expected to find key %d", keys[i])
		}
		if got, ok := tree.Get(keys[i]); !ok {
			t.Fatalf("Expected to find key %d", keys[i])
		} else if got != values[i] {
			t.Fatalf("Expected value %s, got %s", values[i], got)
		}
	}

	// Shuffle the order of pairs to delete.
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
		values[i], values[j] = values[j], values[i]
	})

	// Delete the pairs.
	for i := 0; i < len(keys); i++ {
		if !tree.Delete(keys[i]) {
			t.Fatalf("Expected to remove key %d", keys[i])
		}
		if tree.Size() != (len(keys) - (i + 1)) {
			t.Fatalf("Expected size %d, got %d", len(keys)-(i+1), tree.Size())
		}
	}
}
