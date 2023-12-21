package stream

import "math/rand"

// Generator represents a function that produces an infinite sequence of elements of type E.
// Used by the Generate function.
// Example: Random numbers.
type Generator[E any] func() E

// Generate returns an infinite stream that produces elements by invoking the given generator function.
//
// Example usage:
//
//	s := stream.Generate(func() int { return 1	})
//	out := stream.DebugString(s) // "<1, 1, 1, 1, 1, 1, 1, 1, 1, 1, ...>"
func Generate[E any](next Generator[E]) Stream[E] {
	return func(yield Consumer[E]) bool {
		for {
			if !yield(next()) {
				return false
			}
		}
	}
}

// RandomInt returns a stream that produces pseudo-random, non-negative int values using the given rand.Source.
func RandomInt(source rand.Source) Stream[int] {
	rnd := rand.New(source)
	return Generate(rnd.Int)
}

// RandomIntn returns a stream that produces pseudo-random int values in the range [0, n) using the given rand.Source.
func RandomIntn(source rand.Source, n int) Stream[int] {
	rnd := rand.New(source)
	return Generate(func() int { return rnd.Intn(n) })
}

// RandomUint32 returns a stream that produces pseudo-random uint32 values using the given rand.Source.
func RandomUint32(source rand.Source) Stream[uint32] {
	rnd := rand.New(source)
	return Generate(rnd.Uint32)
}

// RandomUint64 returns a stream that produces pseudo-random uint64 values using the given rand.Source.
func RandomUint64(source rand.Source) Stream[uint64] {
	rnd := rand.New(source)
	return Generate(rnd.Uint64)
}

// RandomFloat32 returns a stream that produces pseudo-random float32 values in the range [0.0, 1.0) using the given rand.Source.
func RandomFloat32(source rand.Source) Stream[float32] {
	rnd := rand.New(source)
	return Generate(rnd.Float32)
}

// RandomFloat64 returns a stream that produces pseudo-random float64 values in the range [0.0, 1.0) using the given rand.Source.
func RandomFloat64(source rand.Source) Stream[float64] {
	rnd := rand.New(source)
	return Generate(rnd.Float64)
}

// RandomNormFloat64 returns a stream that produces normally-distributed float64 values in the range [-math.MaxFloat64, +math.MaxFloat64] using the given rand.Source.
// The mean is 0.0 and the standard deviation is 1.0.
func RandomNormFloat64(source rand.Source) Stream[float64] {
	rnd := rand.New(source)
	return Generate(rnd.NormFloat64)
}

// RandomExpFloat64 returns a stream that produces exponentially-distributed float64 values in the range [0.0, +math.MaxFloat64] using the given rand.Source.
// The rate (lambda) parameter is 1.0.
func RandomExpFloat64(source rand.Source) Stream[float64] {
	rnd := rand.New(source)
	return Generate(rnd.ExpFloat64)
}

// RandomBool returns a stream that produces pseudo-random boolean values using the given rand.Source.
func RandomBool(source rand.Source) Stream[bool] {
	rnd := rand.New(source)
	return Generate(func() bool { return rnd.Intn(2) == 0 })
}

// RandomBytes returns a stream that produces pseudo-random byte values using the given rand.Source.
func RandomBytes(source rand.Source) Stream[byte] {
	rnd := rand.New(source)
	return Generate(func() byte { return byte(rnd.Intn(256)) })
}
