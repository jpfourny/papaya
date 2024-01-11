package mapper

// SliceFrom returns a function that accepts a slice of any type E and returns a slice of the same type containing all elements from the provided `start` index to the end of the slice.
func SliceFrom[E any](start int) func([]E) []E {
	return func(s []E) []E {
		return s[start:]
	}
}

// SliceTo returns a function that accepts a slice of any type E and returns a slice of the same type containing all elements from the start of the slice up to, but excluding, the provided `end` index.
func SliceTo[E any](end int) func([]E) []E {
	return func(s []E) []E {
		return s[:end]
	}
}

// SliceFromTo returns a function that accepts a slice of any type E and returns a slice of the same type containing all elements from the provided `start` index up to, but excluding, the provided `end` index.
func SliceFromTo[E any](start int, end int) func([]E) []E {
	return func(s []E) []E {
		return s[start:end]
	}
}
