package stream

// Predicate is a function that accepts a value of type E and returns a boolean.
// It is used to test values for a given property.
// It must be idempotent, free of side effects, and thread-safe.
type Predicate[E any] func(e E) (pass bool)

// Filter returns a stream that only contains elements that pass the given Predicate.
//
// Example usage:
//
//	s := stream.Filter(stream.Of(1, 2, 3), func(e int) bool {
//	    return e % 2 == 0
//	})
//	out := stream.DebugString(s) // "<2>"
func Filter[E any](s Stream[E], p Predicate[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		return s(func(e E) bool {
			if p(e) {
				return yield(e)
			}
			return true
		})
	}
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
func Skip[E any](s Stream[E], n int64) Stream[E] {
	return func(yield Consumer[E]) bool {
		n := n // Shadow with a copy.
		return s(func(e E) bool {
			if n > 0 {
				n--
				return true
			}
			return yield(e)
		})
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

// KeyExtractor represents a function that extracts a key of type K from a value of type E.
type KeyExtractor[E, K any] func(e E) K

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
