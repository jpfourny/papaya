package kvstore

import (
	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/optional"
	"slices"
)

// Store represents a container for key-value pairs.
// Used internally for key-grouping and key-joining operations.
type Store[K, V any] interface {
	Get(key K) optional.Optional[V]
	Put(key K, value V)
	ForEach(func(key K, value V) bool) bool
}

// NewMapped creates a new Store backed by a map.
// The key type K must be comparable.
func NewMapped[K comparable, V any]() Store[K, V] {
	return make(mappedStore[K, V])
}

// NewSorted creates a new Store of sorted keys, ordered by the given cmp.Comparer.
func NewSorted[K any, V any](compare cmp.Comparer[K]) Store[K, V] {
	return &sortedStore[K, V]{
		compare: compare,
	}
}

// Maker is a factory function for creating a Store.
type Maker[K, V any] func() Store[K, V]

// MappedMaker returns a Maker that calls NewMapped.
// The key type K must be comparable.
func MappedMaker[K comparable, V any]() Maker[K, V] {
	return func() Store[K, V] {
		return NewMapped[K, V]()
	}
}

// SortedMaker returns a Maker that calls NewSorted with the given cmp.Comparer.
func SortedMaker[K any, V any](compare cmp.Comparer[K]) Maker[K, V] {
	return func() Store[K, V] {
		return NewSorted[K, V](compare)
	}
}

// mappedStore provides an implementation of Store using the builtin map.
// The key type K must be comparable.
type mappedStore[K comparable, V any] map[K]V

func (s mappedStore[K, V]) Get(key K) optional.Optional[V] {
	if v, ok := s[key]; ok {
		return optional.Of(v)
	}
	return optional.Empty[V]()
}

func (s mappedStore[K, V]) Put(key K, value V) {
	s[key] = value
}

func (s mappedStore[K, V]) ForEach(yield func(K, V) bool) bool {
	for k, v := range s {
		if !yield(k, v) {
			return false
		}
	}
	return true
}

// sortedStore provides an implementation of Store using sorted slices and binary-search.
// The keys are ordered using the given cmp.Comparer.
type sortedStore[K any, V any] struct {
	compare cmp.Comparer[K]
	keys    []K
	values  []V
}

func (s *sortedStore[K, V]) Get(key K) optional.Optional[V] {
	if i, ok := s.indexOf(key); ok {
		return optional.Of(s.values[i])
	}
	return optional.Empty[V]()
}

func (s *sortedStore[K, V]) Put(key K, value V) {
	i, ok := s.indexOf(key)
	if ok {
		s.values[i] = value
	} else {
		s.keys = append(s.keys, key)
		s.values = append(s.values, value)
		copy(s.keys[i+1:], s.keys[i:])     // Shift keys.
		copy(s.values[i+1:], s.values[i:]) // Shift values.
		s.keys[i] = key                    // Insert key.
		s.values[i] = value                // Insert value.
	}
}

func (s *sortedStore[K, V]) indexOf(key K) (int, bool) {
	return slices.BinarySearchFunc(s.keys, key, s.compare)
}

func (s *sortedStore[K, V]) ForEach(f func(K, V) bool) bool {
	for i, k := range s.keys {
		if !f(k, s.values[i]) {
			return false
		}
	}
	return true
}
