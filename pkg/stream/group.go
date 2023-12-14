package stream

import "github.com/jpfourny/papaya/pkg/pair"

// GroupByKey returns a stream that groups key-value pairs by key.
// The resulting stream contains key-value pairs where the key is the same, and the value is a slice of all the values that had that key.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.GroupByKey(stream.Of(
//	    pair.Of("foo", 1),
//	    pair.Of("bar", 2),
//	    pair.Of("foo", 3),
//	))
//	out := stream.DebugString(s) // "<(foo, [1, 3]), (bar, [2])>"
func GroupByKey[K comparable, V any](s Stream[pair.Pair[K, V]]) Stream[pair.Pair[K, []V]] {
	return func(yield Consumer[pair.Pair[K, []V]]) bool {
		// Group elements by key into map.
		m := make(map[K][]V)
		s(func(p pair.Pair[K, V]) bool {
			m[p.First()] = append(m[p.First()], p.Second())
			return true
		})
		// Yield map entries as pairs.
		for k, vs := range m {
			if !yield(pair.Of(k, vs)) {
				return false
			}
		}
		return true
	}
}

// ReduceByKey returns a stream that reduces key-value pairs by key.
// The resulting stream contains key-value pairs where the key is the same, and the value is the result of reducing all the values that had that key.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.ReduceByKey(stream.Of(
//	    pair.Of("foo", 1),
//	    pair.Of("bar", 2),
//	    pair.Of("foo", 3),
//	), func(a, b int) int {
//	    return a + b
//	})
//	out := stream.DebugString(s) // "<("foo", 4), ("bar", 2)>"
func ReduceByKey[K comparable, V any](s Stream[pair.Pair[K, V]], reduce Reducer[V]) Stream[pair.Pair[K, V]] {
	return func(yield Consumer[pair.Pair[K, V]]) bool {
		// Group reduced elements by key into map.
		m := make(map[K]V)
		s(func(p pair.Pair[K, V]) bool {
			if v, ok := m[p.First()]; ok {
				m[p.First()] = reduce(v, p.Second())
			} else {
				m[p.First()] = p.Second()
			}
			return true
		})
		// Yield map entries as pairs.
		for k, v := range m {
			if !yield(pair.Of(k, v)) {
				return false
			}
		}
		return true
	}
}

// AggregateByKey returns a stream that aggregates key-value pairs by key.
// The resulting stream contains key-value pairs where the key is the same, and the value is the result of aggregating all the values that had that key.
// This is a generalization of ReduceByKey that allows an intermediate accumulated value to be of a different type than both the input and the final result.
// The accumulated value is initialized with the given identity value, and then each element from the input stream is combined with the accumulated value using the given `accumulate` function.
// Once all elements have been accumulated, the accumulated value is transformed into the final result using the given `finish` function.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.AggregateByKey(stream.Of(
//	    pair.Of("foo", 1),
//	    pair.Of("bar", 2),
//	    pair.Of("foo", 3),
//	), 0, func(a int, b int) int {
//	    return a + b
//	}, func(a int) string {
//	    return fmt.Sprintf("%d", a)
//	})
//	out := stream.DebugString(s) // "<("foo", "4"), ("bar", "2")>"
func AggregateByKey[K comparable, V, A, F any](s Stream[pair.Pair[K, V]], identity A, accumulate Accumulator[A, V], finish Finisher[A, F]) Stream[pair.Pair[K, F]] {
	return func(yield Consumer[pair.Pair[K, F]]) bool {
		// Group accumulated elements by key into map.
		m := make(map[K]A)
		s(func(p pair.Pair[K, V]) bool {
			if v, ok := m[p.First()]; ok {
				m[p.First()] = accumulate(v, p.Second())
			} else {
				m[p.First()] = accumulate(identity, p.Second())
			}
			return true
		})
		// Yield map entries as pairs.
		for k, v := range m {
			if !yield(pair.Of(k, finish(v))) {
				return false
			}
		}
		return true
	}
}

// CountByKey returns a stream that counts the number of elements for each key.
// The resulting stream contains key-value pairs where the key is the same, and the value is the number of elements that had that key.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.CountByKey(stream.Of("foo", "bar", "foo"))
//	out := stream.DebugString(s) // "<("foo", 2), ("bar", 1)>"
func CountByKey[K comparable, V any](s Stream[pair.Pair[K, V]]) Stream[pair.Pair[K, int64]] {
	return AggregateByKey(
		s,
		int64(0),
		func(a int64, _ V) int64 { return a + 1 },
		func(a int64) int64 { return a },
	)
}
