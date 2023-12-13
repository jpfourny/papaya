package stream

import (
	"context"
	"slices"
	"strings"

	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/constraint"
	"github.com/jpfourny/papaya/pkg/mapper"
	"github.com/jpfourny/papaya/pkg/optional"
	"github.com/jpfourny/papaya/pkg/pair"
	"github.com/jpfourny/papaya/pkg/pred"
)

// Stream represents a function that produces a sequence of elements of type E and sends them to the given Consumer.
//
// Streams are lazy, meaning they only produce elements when the consumer is invoked.
// Furthermore, streams are idempotent, meaning they can be invoked multiple times with the same result.
// However, the order of the elements is not guaranteed to be the same across multiple invocations.
//
// If the Consumer returns false, the stream must stop producing elements and return false immediately.
// If the stream is exhausted, it must return true.
type Stream[E any] func(c Consumer[E]) bool

// Consumer represents a function that accepts a yielded element of type E and returns a boolean value.
// The boolean value indicates whether the consumer wishes to continue accepting elements.
// If the consumer returns false, the caller must stop yielding elements.
type Consumer[E any] func(yield E) (cont bool)

// Iterator represents a type that provides a way to iterate over a sequence of elements.
// The Next() method is used to obtain the next element in the iteration, if any.
// If there are no more elements, the method returns an empty optional.Optional..
type Iterator[E any] interface {
	Next() optional.Optional[E]
}

// Predicate is a function that accepts a value of type E and returns a boolean.
// It is used to test values for a given property.
// It must be idempotent, free of side effects, and thread-safe.
type Predicate[E any] func(e E) (pass bool)

// Mapper represents a function that transforms an input of type E to an output of type F.
// It is used in the Map operation.
// It must be idempotent, free of side effects, and thread-safe.
type Mapper[E, F any] func(from E) (to F)

// FlatMapper represents a function that takes an input of type E and returns an output stream of type F.
// It is used to map each element of the input stream to a new stream of elements of type F.
// The FlatMapper function is typically used as a parameter in the FlatMap function.
// It must be idempotent, free of side effects, and thread-safe.
type FlatMapper[E, F any] func(from E) (to Stream[F])

// Reducer represents a function that takes two inputs of type E and returns an output of type E.
// The Reducer is commonly used in the `Reduce` function to combine elements of a stream into a single result.
type Reducer[E any] func(e1, e2 E) (result E)

// Accumulator represents a function that takes an accumulated value of type A and an element of type E,
// and returns the updated accumulated value of type A.
// The Accumulator is commonly used in the `Aggregate` function to combine elements of a stream into a single result.
type Accumulator[A, E any] func(a A, e E) (result A)

// Finisher represents a function that takes an accumulated value of type A and returns the finished result of type F.
// The Finisher is commonly used in the `Aggregate` function to compute the final result after all elements have been accumulated.
type Finisher[A, F any] func(a A) (result F)

// KeyExtractor represents a function that extracts a key of type K from a value of type E.
type KeyExtractor[E, K any] func(e E) K

// Empty returns a stream that does not contain any elements.
// It always returns true when invoked with a consumer.
//
// Example usage:
//
//	s := stream.Empty[int]()
//	out := stream.DebugString(s) // "<>"
func Empty[E any]() Stream[E] {
	return func(_ Consumer[E]) bool {
		return true
	}
}

// Of creates a stream from the given elements.
//
// Example usage:
//
//	s := stream.Of(1, 2, 3)
//	out := stream.DebugString(s) // "<1, 2, 3>"
func Of[E any](e ...E) Stream[E] {
	return FromSlice(e)
}

// FromSlice creates a stream that iterates over the elements of the given slice.
// The order of the elements is guaranteed to be the same as the order in the slice.
//
// Example usage:
//
//	s := stream.FromSlice([]int{1, 2, 3})
//	out := stream.DebugString(s) // "<1, 2, 3>"
func FromSlice[E any](s []E) Stream[E] {
	return func(yield Consumer[E]) bool {
		for _, e := range s {
			if !yield(e) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}

// FromSliceWithIndex returns a stream that iterates over the elements of the input slice along with their indices.
// The stream returns a `pair.Pair` with both the index and the value of each element.
// The order of the elements is guaranteed to be the same as the order in the slice.
//
// Example usage:
//
//	s := stream.FromSliceWithIndex([]int{1, 2, 3})
//	out := stream.DebugString(s) // "<(0, 1), (1, 2), (2, 3)>"
func FromSliceWithIndex[E any](s []E) Stream[pair.Pair[int, E]] {
	return func(yield Consumer[pair.Pair[int, E]]) bool {
		for i, e := range s {
			if !yield(pair.Of(i, e)) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}

// FromMap returns a stream that iterates over the key-value pairs in the given map.
// The key-value pairs are encapsulated in `pair.Pair` objects.
// The order of the key-value pairs is not guaranteed.
//
// Example usage:
//
//	s := stream.FromMap(map[int]string{1: "foo", 2: "bar"})
//	out := stream.DebugString(s) // "<(1, foo), (2, bar)>"
func FromMap[K comparable, V any](m map[K]V) Stream[pair.Pair[K, V]] {
	return func(yield Consumer[pair.Pair[K, V]]) bool {
		for k, v := range m {
			if !yield(pair.Of(k, v)) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}

// FromMapKeys takes a map and returns a stream of its keys.
// The order of the keys is not guaranteed.
//
// Example usage:
//
//	s := stream.FromMapKeys(map[int]string{1: "foo", 2: "bar"})
//	out := stream.DebugString(s) // "<1, 2>" // Order not guaranteed.
func FromMapKeys[K comparable, V any](m map[K]V) Stream[K] {
	return func(yield Consumer[K]) bool {
		for k := range m {
			if !yield(k) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}

// FromMapValues takes a map and returns a stream of its values.
// The order of the values is not guaranteed.
//
// Example usage:
//
//	s := stream.FromMapValues(map[int]string{1: "foo", 2: "bar"})
//	out := stream.DebugString(s) // "<foo, bar>" // Order not guaranteed.
func FromMapValues[K comparable, V any](m map[K]V) Stream[V] {
	return func(yield Consumer[V]) bool {
		for _, v := range m {
			if !yield(v) {
				return false
			}
		}
		return true
	}
}

// FromChannel returns a stream that reads elements from the given channel until it is closed.
//
//	Note: If the channel is not closed, the stream will block forever.
//	Note: If the consumer returns false, channel writes may block forever if the channel is unbuffered.
//
// Example usage:
//
//	ch := make(chan int)
//	go func() {
//	    ch <- 1
//	    ch <- 2
//	    close(ch)
//	}()
//	s := stream.FromChannel(ch)
//	out := stream.DebugString(s) // "<1, 2>"
func FromChannel[E any](ch <-chan E) Stream[E] {
	return func(yield Consumer[E]) bool {
		for e := range ch {
			if !yield(e) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}

// FromChannelCtx behaves like FromChannel, but it returns immediately when the context is done.
//
//	Note: If the context is done, channel writes may block forever if the channel is unbuffered.
//
// Example usage:
//
//	ch := make(chan int)
//	go func() {
//	    ch <- 1
//	    ch <- 2
//	    close(ch)
//	}()
//	s := stream.FromChannelCtx(ctx, ch)
//	out := stream.DebugString(s) // "<1, 2>"
func FromChannelCtx[E any](ctx context.Context, ch <-chan E) Stream[E] {
	return func(yield Consumer[E]) bool {
		for {
			select {
			case e, ok := <-ch:
				if !ok {
					return true
				}
				if !yield(e) {
					return false
				}
			case <-ctx.Done():
				return true
			}
		}
	}
}

// FromIterator returns a stream that produces elements by calling the provided Iterator function.
func FromIterator[E any](it Iterator[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		for e := it.Next(); e.Present(); e = it.Next() {
			if !yield(e.Get()) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}

// Range returns a stream that produces elements over a range beginning at `start`, advanced by the `next` function, and ending when `cond` predicate returns false.
//
// Example usage:
//
//	s := stream.Range(1, pred.LessThanOrEqual(5), mapper.Increment(2))
//	out := stream.DebugString(s) // "<1, 3, 5>"
func Range[E any](start E, cond Predicate[E], next Mapper[E, E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		for e := start; cond(e); e = next(e) {
			if !yield(e) {
				return false
			}
		}
		return true
	}
}

// Union combines multiple streams into a single stream (concatenation).
// The length of the resulting stream is the sum of the lengths of the input streams.
// If any of the input streams return false when invoked with the consumer, the concatenation stops.
//
// Example usage:
//
//	s := stream.Union(stream.Of(1, 2, 3, 4), stream.Of(4, 5, 6))
//	out := stream.DebugString(s) // "<1, 2, 3, 4, 4, 5, 6>"
func Union[E any](ss ...Stream[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		for _, s := range ss {
			if !s(yield) {
				return false // Consumer saw enough.
			}
		}
		return true
	}
}

// Intersection returns a stream that contains elements that are in all the given streams.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.Intersection(stream.Of(1, 2, 3, 4, 5), stream.Of(4, 5, 6), stream.Of(4, 5))
//	out := stream.DebugString(s) // "<4, 5>"
func Intersection[E comparable](ss ...Stream[E]) Stream[E] {
	switch {
	case len(ss) == 0:
		return Empty[E]() // No streams.
	case len(ss) == 1:
		return ss[0] // One stream.
	case len(ss) > 2:
		// Recursively intersect the first stream with the intersection of rest.
		return Intersection(ss[0], Intersection(ss[1:]...))
	}

	return func(yield Consumer[E]) bool {
		// Index elements of the first stream into a set.
		seen := make(map[E]struct{})
		ss[0](func(e E) bool {
			seen[e] = struct{}{}
			return true
		})
		// Yield elements of the second stream that are in the set.
		return ss[1](func(e E) bool {
			if _, ok := seen[e]; ok {
				return yield(e)
			}
			return true
		})
	}
}

// Difference returns a stream that contains elements that are in the first stream but not in the second stream.
// The order of the elements is not guaranteed.
//
// Example usage:
//
//	s := stream.Difference(stream.Of(1, 2, 3, 4, 5), stream.Of(4, 5, 6))
//	out := stream.DebugString(s) // "<1, 2, 3>"
func Difference[E comparable](s1, s2 Stream[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		// Index elements of the second stream into a set.
		seen := make(map[E]struct{})
		s2(func(e E) bool {
			seen[e] = struct{}{}
			return true
		})
		// Yield elements of the first stream that are not in the set.
		return s1(func(e E) bool {
			if _, ok := seen[e]; !ok {
				return yield(e)
			}
			return true
		})
	}
}

// Peek decorates the given stream to invoke the given function for each element passing through it.
// This is useful for debugging or logging elements as they pass through the stream.
//
// Example usage:
//
//	s := stream.Peek(stream.Of(1, 2, 3), func(e int) {
//	    fmt.Println(e)
//	})
//	stream.Count(s) // Force the stream to materialize.
//
// Output:
//
//	1
//	2
//	3
func Peek[E any](s Stream[E], peek func(e E)) Stream[E] {
	return func(yield Consumer[E]) bool {
		return s(func(e E) bool {
			peek(e)
			return yield(e)
		})
	}
}

// Map applies a Mapper function to each element in a stream and returns a new stream containing the mapped elements.
//
// Example usage:
//
//	s := stream.Map(stream.Of(1, 2, 3), mapper.Sprint)
//	out := stream.DebugString(s) // "<1, 2, 3>"
func Map[E, F any](s Stream[E], m Mapper[E, F]) Stream[F] {
	return func(yield Consumer[F]) bool {
		return s(func(e E) bool {
			return yield(m(e))
		})
	}
}

// FlatMap applies a FlatMapper function to each element in a stream and returns a new stream containing the mapped elements.
// All mapped streams are flattened into a single stream.
//
// Example usage:
//
//	s := stream.FlatMap(stream.Of(1, 2, 3), func(e int) stream.Stream[string] {
//	    return stream.Map(stream.RangeInteger(0, e), mapper.Sprint)
//	})
//	out := stream.DebugString(s) // "<0, 0, 1, 0, 1, 2>"
func FlatMap[E, F any](s Stream[E], fm FlatMapper[E, F]) Stream[F] {
	return func(yield Consumer[F]) bool {
		return s(func(e E) bool {
			return fm(e)(yield)
		})
	}
}

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

// Limit returns a stream that is limited to the first `n` elements.
// If the input stream has fewer than `n` elements, the returned stream will have the same number of elements.
//
// Example usage:
//
//	s := stream.Limit(stream.Of(1, 2, 3), 2)
//	out := stream.DebugString(s) // "<1, 2>"
func Limit[E any](s Stream[E], n int64) Stream[E] {
	return func(yield Consumer[E]) bool {
		n := n // Shadow with a copy.
		if n <= 0 {
			return true
		}
		return s(func(e E) bool {
			n--
			return yield(e) && n > 0
		})
	}
}

// Skip returns a stream that skips the first `n` elements.
// If the input stream has fewer than `n` elements, the returned stream will be empty.
//
// Example usage:
//
//	s := stream.Skip(stream.Of(1, 2, 3), 2)
//	out := stream.DebugString(s) // "<3>"
func Skip[E any](yield Stream[E], n int64) Stream[E] {
	return func(c Consumer[E]) bool {
		n := n // Shadow with a copy.
		return yield(func(e E) bool {
			if n > 0 {
				n--
				return true
			}
			return c(e)
		})
	}
}

// Filter returns a stream that only contains elements that pass the given Predicate.
//
// Example usage:
//
//	s := stream.Filter(stream.Of(1, 2, 3), func(e int) bool {
//	    return e % 2 == 0
//	})
//	out := stream.DebugString(s) // "<2>"
func Filter[E any](yield Stream[E], p Predicate[E]) Stream[E] {
	return func(c Consumer[E]) bool {
		return yield(func(e E) bool {
			if p(e) {
				return c(e)
			}
			return true
		})
	}
}

// SortAsc returns a stream that sorts the elements in ascending order.
// The elements must implement the Ordered interface.
//
// Example usage:
//
//	s := stream.SortAsc(stream.Of(3, 1, 2))
//	out := stream.DebugString(s) // "<1, 2, 3>"
func SortAsc[E constraint.Ordered](s Stream[E]) Stream[E] {
	return SortBy(s, cmp.Natural[E]())
}

// SortDesc returns a stream that sorts the elements in descending order.
// The elements must implement the Ordered interface.
//
// Example usage:
//
//	s := stream.SortDesc(stream.Of(3, 1, 2))
//	out := stream.DebugString(s) // "<3, 2, 1>"
func SortDesc[E constraint.Ordered](s Stream[E]) Stream[E] {
	return SortBy(s, cmp.Reverse[E]())
}

// SortBy returns a stream that sorts the elements using the given cmp.Comparer.
//
// Example usage:
//
//	s := stream.SortBy(stream.Of(3, 1, 2), cmp.Natural[int]())
//	out := stream.DebugString(s) // "<1, 2, 3>"
func SortBy[E any](s Stream[E], cmp cmp.Comparer[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		sl := CollectSlice(s)
		slices.SortFunc(sl, cmp)
		return FromSlice(sl)(yield)
	}
}

// Distinct returns a stream that only contains distinct elements.
// The elements must implement the comparable interface.
//
// Example usage:
//
//	s := stream.Distinct(stream.Of(1, 2, 2, 3))
//	out := stream.DebugString(s) // "<1, 2, 3>"
func Distinct[E comparable](s Stream[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		seen := make(map[E]struct{})
		return s(func(e E) bool {
			if _, ok := seen[e]; !ok {
				seen[e] = struct{}{}
				return yield(e)
			}
			return true
		})
	}
}

// DistinctBy returns a stream that only contains distinct elements, as determined by the given key extractor.
// The key extractor is used to extract a key from each element, and the keys are compared to determine distinctness.
// The key must implement the comparable interface.
//
// Example usage:
//
//	type Person struct {
//	    FirstName string
//	    LastName  string
//	}
//
//	s := stream.DistinctBy(stream.Of(
//	    Person{"John", "Doe"},
//	    Person{"Jane", "Doe"},
//	    Person{"John", "Smith"},
//	), func(p Person) string {
//	    return p.FirstName
//	})
//	out := stream.DebugString(s) // "<Person{FirstName:"John", LastName:"Doe"}, Person{FirstName:"Jane", LastName:"Doe"}>"
func DistinctBy[E any, K comparable](s Stream[E], ke KeyExtractor[E, K]) Stream[E] {
	return func(yield Consumer[E]) bool {
		seen := make(map[K]struct{})
		return s(func(e E) bool {
			k := ke(e)
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				return yield(e)
			}
			return true
		})
	}
}

// Zip returns a stream that pairs each element in the first stream with the corresponding element in the second stream.
// The resulting stream will have the same number of elements as the shorter of the two input streams.
//
// Example usage:
//
//	s := stream.Zip(stream.Of(1, 2, 3), stream.Of("foo", "bar"))
//	out := stream.DebugString(s) // "<(1, foo), (2, bar)>"
func Zip[E, F any](s1 Stream[E], s2 Stream[F]) Stream[pair.Pair[E, F]] {
	return func(yield Consumer[pair.Pair[E, F]]) bool {
		done := make(chan struct{})
		defer close(done)

		ch1 := make(chan E)
		go func() {
			defer close(ch1)
			s1(func(e E) bool {
				select {
				case <-done:
					return false
				case ch1 <- e:
					return true
				}
			})
		}()

		ch2 := make(chan F)
		go func() {
			defer close(ch2)
			s2(func(f F) bool {
				select {
				case <-done:
					return false
				case ch2 <- f:
					return true
				}
			})
		}()

		for {
			e, ok1 := <-ch1
			f, ok2 := <-ch2
			if !ok1 || !ok2 {
				return true
			}
			if !yield(pair.Of(e, f)) {
				return false
			}
		}
	}
}

// ZipWithIndexInt returns a stream that pairs each element in the input stream with its index, starting at the given offset.
// The index is of type int.
//
// Example usage:
//
//	s := stream.ZipWithIndexInt(stream.Of("foo", "bar"), 1)
//	out := stream.DebugString(s) // "<(foo, 1), (bar, 2)>"
func ZipWithIndexInt[E any](s Stream[E], offset int) Stream[pair.Pair[E, int]] {
	return func(yield Consumer[pair.Pair[E, int]]) bool {
		i := offset - 1
		return s(func(e E) bool {
			i++
			return yield(pair.Of(e, i))
		})
	}
}

// ZipWithIndexInt64 returns a stream that pairs each element in the input stream with its index, starting at the given offset.
// The index is of type int64.
//
// Example usage:
//
//	s := stream.ZipWithIndexInt64(stream.Of("foo", "bar"), 1)
//	out := stream.DebugString(s) // "<(foo, 1), (bar, 2)>"
func ZipWithIndexInt64[E any](s Stream[E], offset int64) Stream[pair.Pair[E, int64]] {
	return func(yield Consumer[pair.Pair[E, int64]]) bool {
		i := offset - 1
		return s(func(e E) bool {
			i++
			return yield(pair.Of(e, i))
		})
	}
}

// ZipWithKey returns a stream that pairs each element in the input stream with the key extracted from the element using the given key extractor.
// The resulting stream will have the same number of elements as the input stream.
//
// Example usage:
//
//	s := stream.ZipWithKey(stream.Of("foo", "bar"), func(s string) string {
//	    return strings.ToUpper(s)
//	})
//	out := stream.DebugString(s) // "<(FOO, foo), (BAR, bar)>"
func ZipWithKey[E any, K any](s Stream[E], ke KeyExtractor[E, K]) Stream[pair.Pair[K, E]] {
	return func(yield Consumer[pair.Pair[K, E]]) bool {
		return s(func(e E) bool {
			return yield(pair.Of(ke(e), e))
		})
	}
}

// ForEach invokes the given consumer for each element in the stream.
//
// Example usage:
//
//	stream.ForEach(stream.Of(1, 2, 3), func(e int) bool {
//	    fmt.Println(e)
//	})
//
// Output:
//
//	1
//	2
//	3
func ForEach[E any](s Stream[E], yield func(E)) {
	s(func(e E) bool {
		yield(e)
		return true
	})
}

// Count returns the number of elements in the stream.
// The stream is fully consumed.
//
// Example usage:
//
//	n := stream.Count(stream.Of(1, 2, 3)) // 3
func Count[E any](s Stream[E]) (count int64) {
	s(func(e E) bool {
		count++
		return true
	})
	return
}

// IsEmpty returns true if the stream is empty; otherwise false.
//
// Example usage:
//
//	empty := stream.IsEmpty(stream.Of(1, 2, 3)) // false
//	empty := stream.IsEmpty(stream.Empty[int]()) // true
func IsEmpty[E any](s Stream[E]) (empty bool) {
	empty = true
	s(func(e E) bool {
		empty = false
		return false
	})
	return
}

// Reduce combines the elements of the stream into a single value using the given reducer function.
// If the stream is empty, then an empty optional.Optional is returned.
// The stream is fully consumed.
//
// Example usage:
//
//	out := stream.Reduce(stream.Of(1, 2, 3), func(a, e int) int {
//	    return a + e
//	}) // Some(6)
//	out = stream.Reduce(stream.Empty[int](), func(a, e int) int {
//	    return a + e
//	} // None()
func Reduce[E any](s Stream[E], reduce Reducer[E]) (result optional.Optional[E]) {
	result = optional.Empty[E]()
	s(func(e E) bool {
		if result.Present() {
			result = optional.Of(reduce(result.Get(), e))
		} else {
			result = optional.Of(e)
		}
		return true
	})
	return
}

// Aggregate combines the elements of the stream into a single value using the given identity value, accumulator function and finisher function.
// The accumulated value is initialized to the identity value.
// The accumulator function is used to combine each element with the accumulated value.
// The finisher function is used to compute the final result after all elements have been accumulated.
// The stream is fully consumed.
//
// Example usage:
//
//	s := stream.Aggregate(
//	  stream.Of(1, 2, 3),
//	  0,                   // Initial value.
//	  func(a, e int) int {
//	      return a + e     // Accumulate with addition.
//	  },
//	  func(a int) int {
//	      return a * 2     // Finish with multiplication by 2.
//	  },
//	) // (1+2+3) * 2 = 12
func Aggregate[E, A, F any](s Stream[E], identity A, accumulate Accumulator[A, E], finish Finisher[A, F]) F {
	a := identity
	s(func(e E) bool {
		a = accumulate(a, e)
		return true
	})
	return finish(a)
}

// Sum computes the sum of all elements in the stream of any number type E and returns the result as type E.
// The result of an empty stream is the zero value of type E.
// The stream is fully consumed.
//
// Example usage:
//
//	n := stream.Sum(stream.Of(1, 2, 3)) // 6 (int)
func Sum[E constraint.RealNumber](s Stream[E]) E {
	return Aggregate(
		s,
		E(0),
		func(a E, e E) E { return a + e },
		func(a E) E { return a },
	)
}

// SumInteger computes the sum of all elements in the stream of any signed-integer type E and returns the result as type int64.
// The result of an empty stream is the zero value of type int64.
// The stream is fully consumed.
//
// Example usage:
//
//	n := stream.SumInteger(stream.Of(1, 2, 3)) // 6 (int64)
func SumInteger[E constraint.SignedInteger](s Stream[E]) int64 {
	return Aggregate(
		s,
		int64(0),
		func(a int64, e E) int64 { return a + int64(e) },
		func(a int64) int64 { return a },
	)
}

// SumUnsignedInteger computes the sum of all elements in the stream of any unsigned-integer type E and returns the result as type uint64.
// The result of an empty stream is the zero value of type uint64.
// The stream is fully consumed.
//
// Example usage:
//
//	n := stream.SumUnsignedInteger(stream.Of[uint](1, 2, 3)) // 6 (uint64)
func SumUnsignedInteger[E constraint.UnsignedInteger](s Stream[E]) uint64 {
	return Aggregate(
		s,
		uint64(0),
		func(a uint64, e E) uint64 { return a + uint64(e) },
		func(a uint64) uint64 { return a },
	)
}

// SumFloat computes the sum of all elements in the stream of any floating-point type E and returns the result as type float64.
// The result of an empty stream is the zero value of type float64.
// The stream is fully consumed.
//
// Example usage:
//
//	n := stream.SumFloat(stream.Of(1.0, 2.0, 3.0)) // 6.0 (float64)
func SumFloat[E constraint.RealNumber](s Stream[E]) float64 {
	return Aggregate(
		s,
		float64(0),
		func(a float64, e E) float64 { return a + float64(e) },
		func(a float64) float64 { return a },
	)
}

// Average computes the average of all elements in the stream of any number type E and returns the result as type float64.
// The result of an empty stream is the zero value of type float64.
// The stream is fully consumed.
//
// Example usage:
//
//	n := stream.Average(stream.Of(1, 2, 3)) // 2.0 (float64)
func Average[E constraint.RealNumber](s Stream[E]) float64 {
	var count uint64
	return Aggregate(
		s,
		float64(0),
		func(a float64, e E) float64 {
			count++
			return a + float64(e)
		},
		func(a float64) float64 {
			if count == 0 {
				return 0
			}
			return a / float64(count)
		},
	)
}

// CollectSlice returns a slice containing all elements from the stream.
// The stream is fully consumed.
//
// Example usage:
//
//	s := stream.CollectSlice(stream.Of(1, 2, 3)) // []int{1, 2, 3}
func CollectSlice[E any](s Stream[E]) []E {
	return Aggregate(
		s,
		[]E(nil),
		func(a []E, e E) []E { return append(a, e) },
		func(a []E) []E { return a },
	)
}

// CollectMap returns a map containing all key-value pair elements from the stream.
// The stream is fully consumed.
//
// Example usage:
//
//	s := stream.CollectMap(stream.Of(pair.Of(1, "foo"), pair.Of(2, "bar"))) // map[int]string{1: "foo", 2: "bar"}
func CollectMap[K comparable, V any](s Stream[pair.Pair[K, V]]) map[K]V {
	return Aggregate(
		s,
		make(map[K]V),
		func(a map[K]V, e pair.Pair[K, V]) map[K]V {
			a[e.First()] = e.Second()
			return a
		},
		func(a map[K]V) map[K]V { return a },
	)
}

// CollectChannel sends all elements from the stream to the given channel.
// The channel must be buffered or have a receiver ready to receive the elements in another goroutine.
// The method returns when the stream is exhausted and all elements have been sent to the channel.
//
// Example usage:
//
//	ch := make(chan int, 3)
//	go func() {
//	  for e := range ch {
//	    fmt.Println(e)
//	  }
//	}()
//	stream.CollectChannel(stream.Of(1, 2, 3), ch)
//	close(ch)
//
// Output:
//
//	1
//	2
//	3
func CollectChannel[E any](s Stream[E], ch chan<- E) {
	s(func(e E) bool {
		ch <- e
		return true
	})
}

// CollectChannelCtx behaves like CollectChannel, but it returns immediately when the context is done.
//
// Example usage:
//
//	ch := make(chan int, 3)
//	go func() {
//	  for e := range ch {
//	    fmt.Println(e)
//	  }
//	}()
//	ctx := context.Background()
//	stream.CollectChannelCtx(ctx, stream.Of(1, 2, 3), ch)
//	close(ch)
//
// Output:
//
//	1
//	2
//	3
func CollectChannelCtx[E any](ctx context.Context, s Stream[E], ch chan<- E) {
	s(func(e E) bool {
		select {
		case <-ctx.Done():
			return false
		case ch <- e:
			return true
		}
	})
}

// CollectChannelAsync sends all elements from the stream to a new channel returned by the method.
// A goroutine is started to send the elements to the channel and the method returns immediately.
// The caller must ensure the channel is fully consumed, otherwise the goroutine will be blocked forever.
// The channel is buffered with the given buffer size, and is closed when the stream is exhausted.
//
// Example usage:
//
//	ch := stream.CollectChannelAsync(stream.Of(1, 2, 3), 3)
//	for e := range ch {
//	    fmt.Println(e)
//	}
//
// Output:
//
//	1
//	2
//	3
func CollectChannelAsync[E any](s Stream[E], buf int) <-chan E {
	ch := make(chan E, buf)
	go func() {
		defer close(ch)
		s(func(e E) bool {
			ch <- e
			return true
		})
	}()
	return ch
}

// CollectChannelAsyncCtx behaves like CollectChannelAsync, but consumption of the stream stops when the context is done.
// The caller can assume the channel will eventually close when the context is done.
//
// Example usage:
//
//	ctx := context.Background()
//	ch := stream.CollectChannelAsyncCtx(ctx, stream.Of(1, 2, 3), 3)
//	for e := range ch {
//	    fmt.Println(e)
//	}
//
// Output:
//
//	1
//	2
//	3
func CollectChannelAsyncCtx[E any](ctx context.Context, s Stream[E], buf int) <-chan E {
	ch := make(chan E, buf)
	go func() {
		defer close(ch)
		s(func(e E) bool {
			select {
			case <-ctx.Done():
				return false
			case ch <- e:
				return true
			}
		})
	}()
	return ch
}

// DebugString returns a string representation of all elements from the stream.
// The string is formatted like `<e1, e2, e3>` where e1, e2, e3 are the string representations of the elements.
// Useful for debugging.
func DebugString[E any](s Stream[E]) string {
	return "<" + StringJoin(Map(s, mapper.Sprintf[E]("%#v")), ", ") + ">"
}

// StringJoin concatenates the elements of the provided stream of strings into a single string, using the specified separator.
//
// Example usage:
//
//	s := stream.Of("foo", "bar", "baz")
//	out := stream.StringJoin(s, ", ") // "foo, bar, baz"
func StringJoin(s Stream[string], sep string) string {
	sb := strings.Builder{}
	s(func(e string) bool {
		if sb.Len() > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(e)
		return true
	})
	return sb.String()
}

// Min returns the minimum element in the stream, or the zero value of the type parameter E if the stream is empty.
// If the stream is empty, the 'ok' return value is false; otherwise it is true.
// Uses the natural ordering of type E to compare elements.
//
// Example usage:
//
//	min := stream.Min(stream.Of(3, 1, 2)) // Some(1)
//	min = stream.Min(stream.Empty[int]()) // None()
func Min[E constraint.Ordered](s Stream[E]) (min optional.Optional[E]) {
	return MinBy(s, cmp.Natural[E]())
}

// MinBy returns the minimum element in the stream.
// Uses the given cmp.Comparer to compare elements.
// If the stream is empty, then an empty optional.Optional is returned.
//
// Example usage:
//
//	min := stream.MinBy(stream.Of(3, 1, 2), cmp.Natural[int]()) // Some(1)
//	min = stream.MinBy(stream.Empty[int](), cmp.Natural[int]()) // None()
func MinBy[E any](s Stream[E], cmp cmp.Comparer[E]) (min optional.Optional[E]) {
	return Reduce(
		s,
		func(a, e E) E {
			if cmp(e, a) < 0 {
				return e
			}
			return a
		},
	)
}

// Max returns the maximum element in the stream.
// Uses the natural ordering of type E to compare elements.
// If the stream is empty, then an empty optional.Optional is returned.
//
// Example usage:
//
//	max := stream.Max(stream.Of(3, 1, 2)) // Some(3)
//	max = stream.Max(stream.Empty[int]()) // None()
func Max[E constraint.Ordered](s Stream[E]) (max optional.Optional[E]) {
	return MaxBy(s, cmp.Natural[E]())
}

// MaxBy returns the maximum element in the stream, or the zero value of the type parameter E if the stream is empty.
// Uses the given cmp.Comparer to compare elements.
// If the stream is empty, then an empty optional.Optional is returned.
//
// Example usage:
//
//	max := stream.MaxBy(stream.Of(3, 1, 2), cmp.Natural[int]()) // Some(3)
//	max = stream.MaxBy(stream.Empty[int](), cmp.Natural[int]()) // None()
func MaxBy[E any](s Stream[E], cmp cmp.Comparer[E]) (max optional.Optional[E]) {
	return Reduce(
		s,
		func(a, e E) E {
			if cmp(e, a) > 0 {
				return e
			}
			return a
		},
	)
}

// First returns the first element in the stream; an empty optional.Optional, if the stream is empty.
//
// Example usage:
//
//	out := stream.First(stream.Of(1, 2, 3)) // Some(1)
//	out = stream.First(stream.Empty[int]()) // None()
func First[E any](s Stream[E]) (first optional.Optional[E]) {
	first = optional.Empty[E]()
	s(func(e E) bool {
		first = optional.Of(e)
		return false
	})
	return
}

// Last returns the last element in the stream; an empty optional.Optional, if the stream is empty.
//
// Example usage:
//
//	out := stream.Last(stream.Of(1, 2, 3)) // Some(3)
//	out = stream.Last(stream.Empty[int]()) // None()
func Last[E any](s Stream[E]) (last optional.Optional[E]) {
	last = optional.Empty[E]()
	s(func(e E) bool {
		last = optional.Of(e)
		return true
	})
	return
}

// Contains returns true if the stream contains the given element; false otherwise.
//
// Example usage:
//
//	out := stream.Contains(stream.Of(1, 2, 3), 2) // true
//	out = stream.Contains(stream.Of(1, 2, 3), 4) // false
func Contains[E comparable](s Stream[E], e E) bool {
	return AnyMatch(s, pred.Equal(e))
}

// ContainsAny returns true if the stream contains any of the given elements; false otherwise.
//
// Example usage:
//
//	out := stream.ContainsAny(stream.Of(1, 2, 3), 2, 4) // true
//	out = stream.ContainsAny(stream.Of(1, 2, 3), 4, 5) // false
func ContainsAny[E comparable](s Stream[E], es ...E) bool {
	return AnyMatch(s, pred.In(es...))
}

// ContainsAll returns true if the stream contains all the given elements; false otherwise.
//
// Example usage:
//
//	out := stream.ContainsAll(stream.Of(1, 2, 3), 2, 3) // true
//	out = stream.ContainsAll(stream.Of(1, 2, 3), 2, 4) // false
func ContainsAll[E comparable](s Stream[E], es ...E) bool {
	return Count(Distinct(Filter(s, pred.In(es...)))) == int64(len(es))
}

// ContainsNone returns true if the stream contains none of the given elements; false otherwise.
//
// Example usage:
//
//	out := stream.ContainsNone(stream.Of(1, 2, 3), 4, 5) // true
//	out = stream.ContainsNone(stream.Of(1, 2, 3), 2, 4) // false
func ContainsNone[E comparable](s Stream[E], es ...E) bool {
	return NoneMatch(s, pred.In(es...))
}

// AnyMatch returns true if any element in the stream matches the given Predicate.
// If the stream is empty, it returns false.
//
// Example usage:
//
//	out := stream.AnyMatch(stream.Of(1, 2, 3), pred.GreaterThan(2)) // true
//	out = stream.AnyMatch(stream.Of(1, 2, 3), pred.GreaterThan(3)) // false
func AnyMatch[E any](s Stream[E], p Predicate[E]) (anyMatch bool) {
	s(func(e E) bool {
		if p(e) {
			anyMatch = true
			return false // Stop the stream.
		}
		return true
	})
	return
}

// AllMatch returns true if all elements in the stream match the given Predicate.
// If the stream is empty, it returns false.
//
// Example usage:
//
//	out := stream.AllMatch(stream.Of(1, 2, 3), pred.GreaterThan(0)) // true
//	out = stream.AllMatch(stream.Of(1, 2, 3), pred.GreaterThan(1)) // false
func AllMatch[E any](s Stream[E], p Predicate[E]) (allMatch bool) {
	allMatch = true
	empty := true
	s(func(e E) bool {
		empty = false
		if !p(e) {
			allMatch = false
			return false // Stop the stream.
		}
		return true
	})
	allMatch = allMatch && !empty
	return
}

// NoneMatch returns true if no elements in the stream match the given Predicate.
// If the stream is empty, it returns true.
//
// Example usage:
//
//	out := stream.NoneMatch(stream.Of(1, 2, 3), pred.GreaterThan(3)) // true
//	out = stream.NoneMatch(stream.Of(1, 2, 3), pred.GreaterThan(2)) // false
func NoneMatch[E any](s Stream[E], p Predicate[E]) bool {
	return !AnyMatch(s, p)
}
