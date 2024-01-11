package pred

// True returns a function that always returns true.
func True[E any]() func(E) bool {
	return func(E) bool {
		return true
	}
}

// False returns a function that always returns false.
func False[E any]() func(E) bool {
	return func(E) bool {
		return false
	}
}

// Not returns a function that returns the opposite of the provided predicate.
//
// Examples:
//
//	t := pred.True[any]()
//	f := pred.False[any]()
//	pred.Not(t) // evaluates to false
//	pred.Not(f) // evaluates to true
func Not[E any](p func(E) bool) func(E) bool {
	return func(e E) bool {
		return !p(e)
	}
}

// And returns a function that returns true if both the provided predicates return true.
//
// Examples:
//
//	t := pred.True[any]()
//	f := pred.False[any]()
//	pred.And(t, t) // evaluates to true
//	pred.And(t, f) // evaluates to false
//	pred.And(f, f) // evaluates to false
func And[E any](p1, p2 func(E) bool) func(E) bool {
	return func(e E) bool {
		return p1(e) && p2(e)
	}
}

// Or returns a function that returns true if either the provided predicates return true.
//
// Examples:
//
//	t := pred.True[any]()
//	f := pred.False[any]()
//	pred.Or(t, t) // evaluates to true
//	pred.Or(t, f) // evaluates to true
//	pred.Or(f, f) // evaluates to false
func Or[E any](p1, p2 func(E) bool) func(E) bool {
	return func(e E) bool {
		return p1(e) || p2(e)
	}
}

// OneOf returns a function that returns true if exactly one of the provided predicates returns true.
// It returns false when:
//   - none of the provided predicates return true, or
//   - more than one of the provided predicates return true or
//   - no predicates are provided (degenerate case).
//
// Examples:
//
//	t := pred.True[any]()
//	f := pred.False[any]()
//	pred.OneOf() // false (no predicates provided)
//	pred.OneOf(t, t) // false (both predicates match)
//	pred.OneOf(t, f) // true (one predicate matches)
//	pred.OneOf(f, f) // false (no predicates match)
func OneOf[E any](ps ...func(E) bool) func(E) bool {
	return func(e E) bool {
		match := false
		for _, p := range ps {
			if p(e) {
				if match {
					return false // More than one match.
				}
				match = true // First match.
			}
		}
		return match // False if no predicates are provided.
	}
}

// AllOf returns a function that returns true if all the provided predicates return true.
// It returns false when:
//   - any of the provided predicates return false, or
//   - no predicates are provided.
//
// Examples:
//
//	t := pred.True[any]()
//	f := pred.False[any]()
//	pred.AllOf() // false (no predicates provided)
//	pred.AllOf(t, t) // true (both predicates match)
//	pred.AllOf(t, f) // false (one predicate matches)
//	pred.AllOf(f, f) // false (no predicates match)
func AllOf[E any](ps ...func(E) bool) func(E) bool {
	return func(e E) bool {
		for _, p := range ps {
			if !p(e) {
				return false
			}
		}
		return len(ps) > 0 // False if no predicates are provided.
	}
}

// AnyOf returns a function that returns true if any of the provided predicates return true.
// It returns false when:
//   - none of the provided predicates return true, or
//   - no predicates are provided.
//
// Examples:
//
//	t := pred.True[any]()
//	f := pred.False[any]()
//	pred.AnyOf() // false (no predicates provided)
//	pred.AnyOf(t, t) // true (both predicates match)
//	pred.AnyOf(t, f) // true (one predicate matches)
//	pred.AnyOf(f, f) // false (no predicates match)
func AnyOf[E any](ps ...func(E) bool) func(E) bool {
	return func(e E) bool {
		for _, p := range ps {
			if p(e) {
				return true
			}
		}
		return false
	}
}

// NoneOf returns a function that returns true if none of the provided predicates return true.
// It returns false when any of the provided predicates return true.
// If no predicates are provided, it returns true (degenerate case).
//
// Examples:
//
//	t := pred.True[any]()
//	f := pred.False[any]()
//	pred.NoneOf() // true (no predicates provided)
//	pred.NoneOf(t, t) // false (both predicates match)
//	pred.NoneOf(t, f) // false (one predicate matches)
//	pred.NoneOf(f, f) // true (no predicates match)
func NoneOf[E any](ps ...func(E) bool) func(E) bool {
	return func(e E) bool {
		for _, p := range ps {
			if p(e) {
				return false
			}
		}
		return true
	}
}
