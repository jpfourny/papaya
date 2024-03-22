package opt

import "fmt"

// Some represents the presence of a value.
type Some[V any] struct {
	Value V
}

// Assert that Some[V] implements Optional[V]
var _ Optional[any] = Some[any]{}

func (s Some[V]) Present() bool {
	return true
}

func (s Some[V]) Get() (V, bool) {
	return s.Value, true
}

func (s Some[V]) GetOrZero() V {
	return s.Value
}

func (s Some[V]) GetOrDefault(_ V) V {
	return s.Value
}

func (s Some[V]) GetOrFunc(_ func() V) V {
	return s.Value
}

func (s Some[V]) Filter(f func(V) bool) Optional[V] {
	if f(s.Value) {
		return s
	}
	return Empty[V]()
}

func (s Some[V]) IfPresent(f func(V)) bool {
	f(s.Value)
	return true
}

func (s Some[V]) IfPresentElse(f func(V), _ func()) bool {
	f(s.Value)
	return true
}

func (s Some[V]) String() string {
	return fmt.Sprintf("Some(%#v)", s.Value)
}
