package stream

import (
	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/constraint"
	"github.com/jpfourny/papaya/pkg/pair"
)

// GroupByKey returns a stream that values key-value pairs by key.
// The resulting stream contains key-value pairs where the key is the same, and the value is a slice of all the values that had that key.
// The key type K must be comparable.
// The order of the key-value pairs is not guaranteed.
//
// Example usage:
//
//	s := stream.GroupByKey(
//	  stream.Of(
//	    pair.Of("foo", 1),
//	    pair.Of("bar", 2),
//	    pair.Of("foo", 3),
//	  ),
//	)
//	out := stream.DebugString(s) // "<(foo, [1, 3]), (bar, [2])>"
func GroupByKey[K comparable, V any](s Stream[pair.Pair[K, V]]) Stream[pair.Pair[K, []V]] {
	return groupByKey(s, mapKeyStoreFactory[K, []V]())
}

// GroupBySortedKey returns a stream that values key-value pairs by key using the given cmp.Comparer to compare keys.
// The resulting stream contains key-value pairs where the key is the same, and the value is a slice of all the values that had that key.
// The order of the key-value pairs is determined by the given cmp.Comparer.
//
// Example usage:
//
//	s := stream.GroupBySortedKey(
//	  stream.Of(
//	    pair.Of("foo", 1),
//	    pair.Of("bar", 2),
//	    pair.Of("foo", 3),
//	  ),
//	  cmp.Natural[string](), // Compare keys naturally
//	)
//	out := stream.DebugString(s) // "<(bar, [2]), (foo, [1, 3])>"
func GroupBySortedKey[K any, V any](s Stream[pair.Pair[K, V]], keyCompare cmp.Comparer[K]) Stream[pair.Pair[K, []V]] {
	return groupByKey(s, sortedKeyStoreFactory[K, []V](keyCompare))
}

func groupByKey[K any, V any](s Stream[pair.Pair[K, V]], ksf keyStoreFactory[K, []V]) Stream[pair.Pair[K, []V]] {
	return func(yield Consumer[pair.Pair[K, []V]]) bool {
		groups := ksf()
		s(func(p pair.Pair[K, V]) bool {
			g := groups.get(p.First()).OrElse(nil)
			g = append(g, p.Second())
			groups.put(p.First(), g)
			return true
		})
		return groups.forEach(func(k K, vs []V) bool {
			return yield(pair.Of(k, vs))
		})
	}
}

// ReduceByKey returns a stream that reduces key-value pairs by key using the given Reducer to reduce values.
// The resulting stream contains key-value pairs where the key is the same, and the value is the result of reducing all the values that had that key.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.ReduceByKey(
//	  stream.Of(
//	    pair.Of("foo", 1),
//	    pair.Of("bar", 2),
//	    pair.Of("foo", 3),
//	  ),
//	  func(a, b int) int { // Reduce values with addition
//	    return a + b
//	  },
//	)
//	out := stream.DebugString(s) // "<("foo", 4), ("bar", 2)>"
func ReduceByKey[K comparable, V any](s Stream[pair.Pair[K, V]], reduce Reducer[V]) Stream[pair.Pair[K, V]] {
	return reduceByKey(s, mapKeyStoreFactory[K, V](), reduce)
}

// ReduceBySortedKey returns a stream that reduces key-value pairs by key using the given cmp.Comparer to compare keys and the given Reducer to reduce values.
// The resulting stream contains key-value pairs where the key is the same, and the value is the result of reducing all the values that had that key.
// The order of the elements is determined by the given cmp.Comparer.
//
// Example usage:
//
//	s := stream.ReduceBySortedKey(
//	  stream.Of(
//	    pair.Of("foo", 1),
//	    pair.Of("bar", 2),
//	    pair.Of("foo", 3),
//	  ),
//	  cmp.Natural[string](), // Compare keys naturally
//	  func(a, b int) int { // Reduce values with addition
//	    return a + b
//	  },
//	)
//	out := stream.DebugString(s) // "<("bar", 2), ("foo", 4)>"
func ReduceBySortedKey[K any, V any](s Stream[pair.Pair[K, V]], keyCompare cmp.Comparer[K], reduce Reducer[V]) Stream[pair.Pair[K, V]] {
	return reduceByKey(s, sortedKeyStoreFactory[K, V](keyCompare), reduce)
}

func reduceByKey[K any, V any](s Stream[pair.Pair[K, V]], ksf keyStoreFactory[K, V], reduce Reducer[V]) Stream[pair.Pair[K, V]] {
	return func(yield Consumer[pair.Pair[K, V]]) bool {
		groups := ksf()
		s(func(p pair.Pair[K, V]) bool {
			groups.get(p.First()).IfPresentElse(
				func(v V) { // If present
					groups.put(p.First(), reduce(v, p.Second()))
				},
				func() { // Else
					groups.put(p.First(), p.Second())
				},
			)
			return true
		})
		return groups.forEach(func(k K, v V) bool {
			return yield(pair.Of(k, v))
		})
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
//	s := stream.AggregateByKey(
//	  stream.Of(
//	    pair.Of("foo", 1),
//	    pair.Of("bar", 2),
//	    pair.Of("foo", 3),
//	  ),
//	  0, // Initial value
//	  func(a int, b int) int { // Accumulate values with addition
//	    return a + b
//	  },
//	  func(a int) string { // Finish values with string conversion
//	    return fmt.Sprintf("%d", a)
//	  },
//	)
//	out := stream.DebugString(s) // "<("foo", "4"), ("bar", "2")>"
func AggregateByKey[K comparable, V, A, F any](s Stream[pair.Pair[K, V]], identity A, accumulate Accumulator[A, V], finish Finisher[A, F]) Stream[pair.Pair[K, F]] {
	return aggregateByKey(s, mapKeyStoreFactory[K, A](), identity, accumulate, finish)
}

// AggregateBySortedKey returns a stream that aggregates key-value pairs by key using the given cmp.Comparer to compare keys.
// The resulting stream contains key-value pairs where the key is the same, and the value is the result of aggregating all the values that had that key.
// This is a generalization of ReduceBySortedKey that allows an intermediate accumulated value to be of a different type than both the input and the final result.
// The accumulated value is initialized with the given identity value, and then each element from the input stream is combined with the accumulated value using the given `accumulate` function.
// Once all elements have been accumulated, the accumulated value is transformed into the final result using the given `finish` function.
// The order of the elements is determined by the given cmp.Comparer.
//
// Example usage:
//
//	s := stream.AggregateBySortedKey(
//	  stream.Of(
//	    pair.Of("foo", 1),
//	    pair.Of("bar", 2),
//	    pair.Of("foo", 3),
//	  ),
//	  cmp.Natural[string](), // Compare keys naturally
//	  0, // Initial value
//	  func(a int, b int) int { // Accumulate values with addition
//	    return a + b
//	  },
//	  func(a int) string { // Finish values with string conversion
//	    return fmt.Sprintf("%d", a)
//	  },
//	)
//	out := stream.DebugString(s) // "<("bar", "2"), ("foo", "4")>"
func AggregateBySortedKey[K any, V, A, F any](s Stream[pair.Pair[K, V]], keyCompare cmp.Comparer[K], identity A, accumulate Accumulator[A, V], finish Finisher[A, F]) Stream[pair.Pair[K, F]] {
	return aggregateByKey(s, sortedKeyStoreFactory[K, A](keyCompare), identity, accumulate, finish)
}

func aggregateByKey[K any, V, A, F any](s Stream[pair.Pair[K, V]], ksf keyStoreFactory[K, A], identity A, accumulate Accumulator[A, V], finish Finisher[A, F]) Stream[pair.Pair[K, F]] {
	return func(yield Consumer[pair.Pair[K, F]]) bool {
		groups := ksf()
		s(func(p pair.Pair[K, V]) bool {
			groups.get(p.First()).IfPresentElse(
				func(a A) { // If present
					groups.put(p.First(), accumulate(a, p.Second()))
				},
				func() { // Else
					groups.put(p.First(), accumulate(identity, p.Second()))
				},
			)
			return true
		})
		groups.forEach(func(k K, a A) bool {
			return yield(pair.Of(k, finish(a)))
		})
		return true
	}
}

// CountByKey returns a stream that counts the number of elements for each key.
// The resulting stream contains key-value pairs where the key is the same, and the value is the number of elements that had that key.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.CountByKey(
//	  stream.Of(
//	    pair.Of("foo", 1),
//	    pair.Of("bar", 2),
//	    pair.Of("foo", 3),
//	  ),
//	)
//	out := stream.DebugString(s) // "<("foo", 2), ("bar", 1)>"
func CountByKey[K comparable, V any](s Stream[pair.Pair[K, V]]) Stream[pair.Pair[K, int64]] {
	return AggregateByKey(
		s,
		int64(0),
		func(a int64, _ V) int64 { return a + 1 },
		func(a int64) int64 { return a },
	)
}

// CountBySortedKey returns a stream that counts the number of elements for each key using the given cmp.Comparer to compare keys.
// The resulting stream contains key-value pairs where the key is the same, and the value is the number of elements that had that key.
// The order of the elements is determined by the given cmp.Comparer.
//
// Example usage:
//
//	s := stream.CountBySortedKey(
//	  stream.Of(
//	    pair.Of("foo", 1),
//	    pair.Of("bar", 2),
//	    pair.Of("foo", 3),
//	  ),
//	  cmp.Natural[string](), // Compare keys naturally
//	)
//	out := stream.DebugString(s) // "<("bar", 1), ("foo", 2)>"
func CountBySortedKey[K any, V any](s Stream[pair.Pair[K, V]], keyCompare cmp.Comparer[K]) Stream[pair.Pair[K, int64]] {
	return AggregateBySortedKey(
		s,
		keyCompare,
		int64(0),
		func(a int64, _ V) int64 { return a + 1 },
		func(a int64) int64 { return a },
	)
}

// MinByKey returns a stream that finds the minimum value for each key.
// The resulting stream contains key-value pairs where the key is the same, and the value is the minimum value that had that key.
// The key type K must be comparable.
// The value type V must be ordered.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.MinByKey(
//	  stream.Of(
//	    pair.Of("foo", 1),
//	    pair.Of("bar", 2),
//	    pair.Of("foo", 3),
//	  ),
//	)
//	out := stream.DebugString(s) // "<("foo", 1), ("bar", 2)>"
func MinByKey[K comparable, V constraint.Ordered](s Stream[pair.Pair[K, V]]) Stream[pair.Pair[K, V]] {
	return ReduceByKey(
		s,
		func(a, b V) V { return min(a, b) },
	)
}

// MinBySortedKey returns a stream that finds the minimum value for each key using the given cmp.Comparer to compare keys.
// The resulting stream contains key-value pairs where the key is the same, and the value is the minimum value that had that key.
// The value type V must be ordered.
// The order of the elements is determined by the given cmp.Comparer.
//
// Example usage:
//
//	s := stream.MinBySortedKey(
//	  stream.Of(
//	    pair.Of("foo", 1),
//	    pair.Of("bar", 2),
//	    pair.Of("foo", 3),
//	  ),
//	  cmp.Natural[string](), // Compare keys naturally
//	)
//	out := stream.DebugString(s) // "<("bar", 2), ("foo", 1)>"
func MinBySortedKey[K any, V constraint.Ordered](s Stream[pair.Pair[K, V]], keyCompare cmp.Comparer[K]) Stream[pair.Pair[K, V]] {
	return ReduceBySortedKey(
		s,
		keyCompare,
		func(a, b V) V { return min(a, b) },
	)
}

// MaxByKey returns a stream that finds the maximum value for each key.
// The resulting stream contains key-value pairs where the key is the same, and the value is the maximum value that had that key.
// The key type K must be comparable.
// The value type V must be ordered.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//		s := stream.MaxByKey(
//	 	stream.Of(
//		    pair.Of("foo", 1),
//		    pair.Of("bar", 2),
//		    pair.Of("foo", 3),
//		  ),
//		)
//		out := stream.DebugString(s) // "<("foo", 3), ("bar", 2)>"
func MaxByKey[K comparable, V constraint.Ordered](s Stream[pair.Pair[K, V]]) Stream[pair.Pair[K, V]] {
	return ReduceByKey(
		s,
		func(a, b V) V { return max(a, b) },
	)
}

// MaxBySortedKey returns a stream that finds the maximum value for each key using the given cmp.Comparer to compare keys.
// The resulting stream contains key-value pairs where the key is the same, and the value is the maximum value that had that key.
// The value type V must be ordered.
// The order of the elements is determined by the given cmp.Comparer.
//
// Example usage:
//
//	s := stream.MaxBySortedKey(
//	  stream.Of(
//	    pair.Of("foo", 1),
//	    pair.Of("bar", 2),
//	    pair.Of("foo", 3),
//	  ),
//	  cmp.Natural[string](), // Compare keys naturally
//	)
//	out := stream.DebugString(s) // "<("bar", 2), ("foo", 3)>"
func MaxBySortedKey[K any, V constraint.Ordered](s Stream[pair.Pair[K, V]], keyCompare cmp.Comparer[K]) Stream[pair.Pair[K, V]] {
	return ReduceBySortedKey(
		s,
		keyCompare,
		func(a, b V) V { return max(a, b) },
	)
}

// SumByKey returns a stream that sums the values for each key.
// The resulting stream contains key-value pairs where the key is the same, and the value is the sum of all the values that had that key.
// The key type K must be comparable.
// The value type V must be numeric.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//		s := stream.SumByKey(
//	 	stream.Of(
//		    pair.Of("foo", 1),
//		    pair.Of("bar", 2),
//		    pair.Of("foo", 3),
//		  ),
//		)
//		out := stream.DebugString(s) // "<("foo", 4), ("bar", 2)>"
func SumByKey[K comparable, V constraint.Numeric](s Stream[pair.Pair[K, V]]) Stream[pair.Pair[K, V]] {
	return ReduceByKey(
		s,
		func(a, b V) V { return a + b },
	)
}

// SumBySortedKey returns a stream that sums the values for each key using the given cmp.Comparer to compare keys.
// The resulting stream contains key-value pairs where the key is the same, and the value is the sum of all the values that had that key.
// The value type V must be numeric.
// The order of the elements is determined by the given cmp.Comparer.
//
// Example usage:
//
//	s := stream.SumBySortedKey(
//	  stream.Of(
//	    pair.Of("foo", 1),
//	    pair.Of("bar", 2),
//	    pair.Of("foo", 3),
//	  ),
//	  cmp.Natural[string](), // Compare keys naturally
//	)
//	out := stream.DebugString(s) // "<("bar", 2), ("foo", 4)>"
func SumBySortedKey[K any, V constraint.Numeric](s Stream[pair.Pair[K, V]], keyCompare cmp.Comparer[K]) Stream[pair.Pair[K, V]] {
	return ReduceBySortedKey(
		s,
		keyCompare,
		func(a, b V) V { return a + b },
	)
}
