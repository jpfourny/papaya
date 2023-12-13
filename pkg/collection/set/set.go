package set

import (
	"github.com/jpfourny/papaya/pkg/mapper"
	"github.com/jpfourny/papaya/pkg/stream"
)

// Set represents map-based set of elements of type E.
type Set[E comparable] map[E]struct{}

// Of creates a new Set with the provided elements.
func Of[E comparable](elems ...E) Set[E] {
	s := make(Set[E], len(elems))
	for _, e := range elems {
		s.Add(e)
	}
	return s
}

// Make creates a new empty Set of element type E with a small initial capacity.
func Make[E comparable]() Set[E] {
	return make(Set[E])
}

// MakeWithCapacity creates a new empty Set of element type E with the specified initial capacity.
func MakeWithCapacity[E comparable](initialCapacity int) Set[E] {
	return make(Set[E], initialCapacity)
}

// FromMap converts the provided map to a Set.
func FromMap[E comparable](m map[E]struct{}) Set[E] {
	return m
}

// Add adds the provided element to the Set.
func (s Set[E]) Add(e E) {
	s[e] = struct{}{}
}

// AddAll adds the provided elements to the Set.
func (s Set[E]) AddAll(elems ...E) {
	for _, e := range elems {
		s.Add(e)
	}
}

// Remove removes the provided element from the Set.
func (s Set[E]) Remove(e E) {
	delete(s, e)
}

// RemoveAll removes all provided elements from the Set.
func (s Set[E]) RemoveAll(elems ...E) {
	for _, e := range elems {
		delete(s, e)
	}
}

// Union returns a new Set containing the union of the `s` Set and the `other` Set.
func (s Set[E]) Union(other Set[E]) Set[E] {
	result := s.Clone()
	for e := range other {
		result.Add(e)
	}
	return result
}

// Intersection returns a new Set containing the intersection of the `s` Set and the `other` Set.
func (s Set[E]) Intersection(other Set[E]) Set[E] {
	sSize := s.Size()
	otherSize := other.Size()
	result := MakeWithCapacity[E](min(sSize, otherSize))

	// Range over smaller set for best performance.
	if sSize > otherSize {
		for e := range other {
			if s.Contains(e) {
				result.Add(e)
			}
		}
	} else {
		for e := range s {
			if other.Contains(e) {
				result.Add(e)
			}
		}
	}

	return result
}

// Difference returns a new Set containing the difference of the `s` Set and the `other` Set.
func (s Set[E]) Difference(other Set[E]) Set[E] {
	result := s.Clone()
	for e := range other {
		result.Remove(e)
	}
	return result
}

// Contains returns true if the Set contains the provided element; false otherwise.
func (s Set[E]) Contains(e E) bool {
	_, ok := s[e]
	return ok
}

// Size returns the number of elements in the Set.
func (s Set[E]) Size() int {
	return len(s)
}

// Empty returns true if the Set is empty; false otherwise.
func (s Set[E]) Empty() bool {
	return s.Size() == 0
}

// Clear removes all elements from the Set.
func (s Set[E]) Clear() {
	clear(s)
}

// Clone returns a shallow copy of the Set.
func (s Set[E]) Clone() Set[E] {
	result := MakeWithCapacity[E](s.Size())
	for e := range s {
		result.Add(e)
	}
	return result
}

// Stream returns a stream of the elements in the Set.
func (s Set[E]) Stream() stream.Stream[E] {
	return stream.FromMapKeys(s)
}

// ToSlice returns a slice of the elements in the Set.
func (s Set[E]) ToSlice() []E {
	return stream.CollectSlice(s.Stream())
}

// String returns a string representation of the Set, formatted as "[<elem1>, <elem2>, ...]".
// The elements are sorted in ascending order after being converted to strings using fmt.Sprint.
func (s Set[E]) String() string {
	return "[" + stream.StringJoin(stream.SortAsc(stream.Map(s.Stream(), mapper.Sprint[E]())), ", ") + "]"
}
