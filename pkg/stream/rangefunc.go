package stream

import "github.com/jpfourny/papaya/pkg/pair"

// FromRangeFunc adapts a range function of parameter type E to a Stream of the same type.
// This exists for compatibility with the range-over-func support in Go (1.22+).
func FromRangeFunc[E any](f func(func(E) bool) bool) Stream[E] {
	return func(yield Consumer[E]) bool {
		return f(yield)
	}
}

// FromRangeBiFunc converts a range bi-function of parameter types K and V into a Stream of `pair.Pair` with the same types.
// This exists for compatibility with the range-over-func support in Go (1.22+).
func FromRangeBiFunc[K, V any](f func(func(K, V) bool) bool) Stream[pair.Pair[K, V]] {
	return func(yield Consumer[pair.Pair[K, V]]) bool {
		return f(func(k K, v V) bool {
			return yield(pair.Of(k, v))
		})
	}
}

// ToRangeFunc converts the given Stream of type E into a range function of the same parameter type.
// This exists for compatibility with the range-over-func support in Go (1.22+).
func ToRangeFunc[E any](s Stream[E]) func(func(E) bool) bool {
	return func(yield func(E) bool) bool {
		return s(yield)
	}
}

// ToRangeBiFunc converts the given Stream of `pair.Pair` with types K and V into a range bi-function with the same parameter types.
// This exists for compatibility with the range-over-func support in Go (1.22+).
func ToRangeBiFunc[K, V any](s Stream[pair.Pair[K, V]]) func(func(K, V) bool) bool {
	return func(yield func(K, V) bool) bool {
		return s(func(p pair.Pair[K, V]) bool {
			return yield(p.Explode())
		})
	}
}
