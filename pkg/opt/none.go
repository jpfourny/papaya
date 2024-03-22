package opt

// None represents the absence of a value.
type None[V any] struct {
}

// Assert that None[V] implements Optional[V]
var _ Optional[any] = None[any]{}

func (n None[V]) Present() bool {
	return false
}

func (n None[V]) Get() (V, bool) {
	var zero V
	return zero, false
}

func (n None[V]) GetOrZero() V {
	var zero V
	return zero
}

func (n None[V]) GetOrDefault(defaultValue V) V {
	return defaultValue
}

func (n None[V]) GetOrFunc(f func() V) V {
	return f()
}

func (n None[V]) Filter(_ func(V) bool) Optional[V] {
	return n
}

func (n None[V]) IfPresent(_ func(V)) bool {
	return false
}

func (n None[V]) IfPresentElse(_ func(V), f func()) bool {
	f()
	return false
}

func (n None[V]) String() string {
	return "None"
}
