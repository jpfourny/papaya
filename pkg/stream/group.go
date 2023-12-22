package stream

import (
	"slices"

	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/constraint"
	"github.com/jpfourny/papaya/pkg/optional"
	"github.com/jpfourny/papaya/pkg/pair"
)

// GroupByKey returns a stream that groups key-value pairs by key.
// The resulting stream contains key-value pairs where the key is the same, and the value is a slice of all the groups that had that key.
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
	return groupByKey(s, newMapGrouper[K, []V]())
}

// GroupBySortedKey returns a stream that groups key-value pairs by key using the given cmp.Comparer to compare keys.
// The resulting stream contains key-value pairs where the key is the same, and the value is a slice of all the groups that had that key.
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
	return groupByKey(s, newSortedGrouper[K, []V](keyCompare))
}

func groupByKey[K any, V any](s Stream[pair.Pair[K, V]], ng newGrouper[K, []V]) Stream[pair.Pair[K, []V]] {
	return func(yield Consumer[pair.Pair[K, []V]]) bool {
		grpr := ng()
		s(func(p pair.Pair[K, V]) bool {
			g := grpr.get(p.First()).OrElse(nil)
			g = append(g, p.Second())
			grpr.put(p.First(), g)
			return true
		})
		return grpr.forEach(func(k K, vs []V) bool {
			return yield(pair.Of(k, vs))
		})
	}
}

// ReduceByKey returns a stream that reduces key-value pairs by key using the given Reducer to reduce values.
// The resulting stream contains key-value pairs where the key is the same, and the value is the result of reducing all the groups that had that key.
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
	return reduceByKey(s, newMapGrouper[K, V](), reduce)
}

// ReduceBySortedKey returns a stream that reduces key-value pairs by key using the given cmp.Comparer to compare keys and the given Reducer to reduce values.
// The resulting stream contains key-value pairs where the key is the same, and the value is the result of reducing all the groups that had that key.
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
	return reduceByKey(s, newSortedGrouper[K, V](keyCompare), reduce)
}

func reduceByKey[K any, V any](s Stream[pair.Pair[K, V]], ng newGrouper[K, V], reduce Reducer[V]) Stream[pair.Pair[K, V]] {
	return func(yield Consumer[pair.Pair[K, V]]) bool {
		grpr := ng()
		s(func(p pair.Pair[K, V]) bool {
			grpr.get(p.First()).IfPresentElse(
				func(v V) { // If present
					grpr.put(p.First(), reduce(v, p.Second()))
				},
				func() { // Else
					grpr.put(p.First(), p.Second())
				},
			)
			return true
		})
		return grpr.forEach(func(k K, v V) bool {
			return yield(pair.Of(k, v))
		})
	}
}

// AggregateByKey returns a stream that aggregates key-value pairs by key.
// The resulting stream contains key-value pairs where the key is the same, and the value is the result of aggregating all the groups that had that key.
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
	return aggregateByKey(s, newMapGrouper[K, A](), identity, accumulate, finish)
}

// AggregateBySortedKey returns a stream that aggregates key-value pairs by key using the given cmp.Comparer to compare keys.
// The resulting stream contains key-value pairs where the key is the same, and the value is the result of aggregating all the groups that had that key.
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
	return aggregateByKey(s, newSortedGrouper[K, A](keyCompare), identity, accumulate, finish)
}

func aggregateByKey[K any, V, A, F any](s Stream[pair.Pair[K, V]], ng newGrouper[K, A], identity A, accumulate Accumulator[A, V], finish Finisher[A, F]) Stream[pair.Pair[K, F]] {
	return func(yield Consumer[pair.Pair[K, F]]) bool {
		grpr := ng()
		s(func(p pair.Pair[K, V]) bool {
			grpr.get(p.First()).IfPresentElse(
				func(a A) { // If present
					grpr.put(p.First(), accumulate(a, p.Second()))
				},
				func() { // Else
					grpr.put(p.First(), accumulate(identity, p.Second()))
				},
			)
			return true
		})
		grpr.forEach(func(k K, a A) bool {
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

// grouper provides map-like interface to index groups by key.
// A group could be anything; for example, a slice of elements or an accumulator for an aggregation.
// This interface and all implementations are for internal use only, so they are not exported.
type grouper[K, G any] interface {
	get(key K) optional.Optional[G]
	put(key K, group G)
	forEach(func(key K, group G) bool) bool
}

// newGrouper represents a factory function that creates a grouper implementation.
type newGrouper[K, G any] func() grouper[K, G]

// newMapGrouper returns a newGrouper factory function for creating a map-based grouper with O(1) access time.
// The key type K must be comparable.
func newMapGrouper[K comparable, G any]() newGrouper[K, G] {
	return func() grouper[K, G] {
		return make(mapGrouper[K, G])
	}
}

// newSortedGrouper returns a newGrouper factory function for creating a sorted grouper with O(log n) access time.
// This is used when the key type is not comparable, ruling out the use of a map-based grouper.
// The keys are sorted using the given cmp.Comparer.
func newSortedGrouper[K any, G any](compare cmp.Comparer[K]) newGrouper[K, G] {
	return func() grouper[K, G] {
		return &sortedGrouper[K, G]{
			compare: compare,
		}
	}
}

// mapGrouper provides an implementation of grouper using a map.
// The key type K must be comparable.
type mapGrouper[K comparable, G any] map[K]G

func (mg mapGrouper[K, G]) get(key K) optional.Optional[G] {
	if g, ok := mg[key]; ok {
		return optional.Of(g)
	}
	return optional.Empty[G]()
}

func (mg mapGrouper[K, G]) put(key K, group G) {
	mg[key] = group
}

func (mg mapGrouper[K, G]) forEach(yield func(key K, group G) bool) bool {
	for k, v := range mg {
		if !yield(k, v) {
			return false
		}
	}
	return true
}

// sortedGrouper provides an implementation of grouper using sorted slices and binary search.
type sortedGrouper[K any, G any] struct {
	compare cmp.Comparer[K]
	keys    []K
	groups  []G
}

func (sg *sortedGrouper[K, G]) get(key K) optional.Optional[G] {
	if i, ok := sg.indexOf(key); ok {
		return optional.Of(sg.groups[i])
	}
	return optional.Empty[G]()
}

func (sg *sortedGrouper[K, G]) put(key K, group G) {
	i, ok := sg.indexOf(key)
	if ok {
		sg.groups[i] = group
	} else {
		sg.keys = append(sg.keys, key)
		sg.groups = append(sg.groups, group)
		copy(sg.keys[i+1:], sg.keys[i:])     // Shift keys.
		copy(sg.groups[i+1:], sg.groups[i:]) // Shift groups.
		sg.keys[i] = key                     // Insert key.
		sg.groups[i] = group                 // Insert group.
	}
}

func (sg *sortedGrouper[K, G]) indexOf(key K) (int, bool) {
	return slices.BinarySearchFunc(sg.keys, key, sg.compare)
}

func (sg *sortedGrouper[K, G]) forEach(f func(key K, group G) bool) bool {
	for i, k := range sg.keys {
		if !f(k, sg.groups[i]) {
			return false
		}
	}
	return true
}
