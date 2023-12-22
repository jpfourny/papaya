package stream

import (
	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/optional"
	"slices"
)

// keyStore represents a store of keys and their associated values.
// Used internally for various stream operations, such as distinct, groupBy, aggregateBy, etc.
type keyStore[K, G any] interface {
	get(key K) optional.Optional[G]
	put(key K, value G)
	forEach(func(key K, value G) bool) bool
}

// keyStoreFactory represents a factory function for creating a keyStore.
type keyStoreFactory[K, V any] func() keyStore[K, V]

// mapKeyStoreFactory returns a keyStoreFactory for creating a map-based keyStore with O(1) access time.
// The key type K must be comparable.
func mapKeyStoreFactory[K comparable, V any]() keyStoreFactory[K, V] {
	return func() keyStore[K, V] {
		return make(mapKeyStore[K, V])
	}
}

// sortedKeyStoreFactory returns a keyStoreFactory for creating a sorted keyStore with O(log n) access time.
// This is used when the key type is not comparable, ruling out the use of a map-based keyStore.
// The keys are sorted using the given cmp.Comparer.
func sortedKeyStoreFactory[K any, V any](compare cmp.Comparer[K]) keyStoreFactory[K, V] {
	return func() keyStore[K, V] {
		return &sortedKeyStore[K, V]{
			compare: compare,
		}
	}
}

// mapKeyStore provides an implementation of keyStore using a map.
// The key type K must be comparable.
type mapKeyStore[K comparable, V any] map[K]V

func (ks mapKeyStore[K, V]) get(key K) optional.Optional[V] {
	if v, ok := ks[key]; ok {
		return optional.Of(v)
	}
	return optional.Empty[V]()
}

func (ks mapKeyStore[K, V]) put(key K, value V) {
	ks[key] = value
}

func (ks mapKeyStore[K, V]) forEach(yield func(K, V) bool) bool {
	for k, v := range ks {
		if !yield(k, v) {
			return false
		}
	}
	return true
}

// sortedKeyStore provides an implementation of keyStore using sorted slices and binary search.
// The keys are sorted using the given cmp.Comparer.
type sortedKeyStore[K any, V any] struct {
	compare cmp.Comparer[K]
	keys    []K
	values  []V
}

func (ks *sortedKeyStore[K, V]) get(key K) optional.Optional[V] {
	if i, ok := ks.indexOf(key); ok {
		return optional.Of(ks.values[i])
	}
	return optional.Empty[V]()
}

func (ks *sortedKeyStore[K, V]) put(key K, value V) {
	i, ok := ks.indexOf(key)
	if ok {
		ks.values[i] = value
	} else {
		ks.keys = append(ks.keys, key)
		ks.values = append(ks.values, value)
		copy(ks.keys[i+1:], ks.keys[i:])     // Shift keys.
		copy(ks.values[i+1:], ks.values[i:]) // Shift values.
		ks.keys[i] = key                     // Insert key.
		ks.values[i] = value                 // Insert value.
	}
}

func (ks *sortedKeyStore[K, V]) indexOf(key K) (int, bool) {
	return slices.BinarySearchFunc(ks.keys, key, ks.compare)
}

func (ks *sortedKeyStore[K, V]) forEach(f func(K, V) bool) bool {
	for i, k := range ks.keys {
		if !f(k, ks.values[i]) {
			return false
		}
	}
	return true
}
