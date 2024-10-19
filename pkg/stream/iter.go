package stream

import (
	"github.com/jpfourny/papaya/v2/pkg/pair"
	"iter"
)

// FromIterSeq creates a Stream from the given iter.Seq.
func FromIterSeq[E any](seq iter.Seq[E]) Stream[E] {
	return func(yield Consumer[E]) {
		seq(yield)
	}
}

// FromIterSeq2 creates a Stream of pairs from the given iter.Seq2.
func FromIterSeq2[K, V any](seq iter.Seq2[K, V]) Stream[pair.Pair[K, V]] {
	return func(yield Consumer[pair.Pair[K, V]]) {
		seq(func(k K, v V) bool {
			return yield(pair.Of(k, v))
		})
	}
}

// ToIterSeq converts a Stream to an iter.Seq.
func ToIterSeq[E any](s Stream[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		s(yield)
	}
}

// ToIterSeq2 converts a Stream of pairs to an iter.Seq2.
func ToIterSeq2[K, V any](s Stream[pair.Pair[K, V]]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		s(func(p pair.Pair[K, V]) bool {
			return yield(p.Explode())
		})
	}
}
