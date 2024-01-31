package stream

import (
	"github.com/jpfourny/papaya/pkg/cmp"
	"testing"

	"github.com/jpfourny/papaya/pkg/opt"
)

func TestReduce(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := Reduce(Empty[int](), func(a, b int) int {
			return a + b
		})
		want := opt.Empty[int]()
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := Reduce(Of(1, 2, 3), func(a, b int) int {
			return a + b
		})
		want := opt.Of(6)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestAggregate(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := Aggregate(
			Empty[int](),
			0, func(a, b int) int { return a + b },
			func(r int) int { return r * 2 },
		)
		want := 0
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := Aggregate(
			Of(1, 2, 3),
			0,
			func(a, b int) int { return a + b },
			func(r int) int { return r * 2 },
		)
		want := 12 // (1+2+3)*2
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}

func TestSum(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := Sum[int](Empty[int]())
		want := 0
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := Sum[int](Of(1, 2, 3))
		want := 6
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}

func TestSumComplex(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := SumComplex[complex128](Empty[complex128]())
		want := complex128(0)
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := SumComplex[complex128](Of(complex(1, 2), complex(3, 4)))
		want := complex(4, 6)
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}

func TestAverage(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := Average(Empty[int]())
		want := 0.0
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := Average(Of(1, 2, 3))
		want := 2.0
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}

func TestMin(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := Min(Empty[int]())
		want := opt.Empty[int]()
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := Min(Of(3, 1, 2))
		want := opt.Of(1)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestMinBy(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := MinBy(Empty[int](), cmp.Natural[int]())
		want := opt.Empty[int]()
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := MinBy(Of(3, 1, 2), cmp.Natural[int]())
		want := opt.Of(1)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestMax(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := Max(Empty[int]())
		want := opt.Empty[int]()
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := Max(Of(1, 3, 2))
		want := opt.Of(3)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestMaxBy(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := MaxBy(Empty[int](), cmp.Natural[int]())
		want := opt.Empty[int]()
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := MaxBy(Of(1, 3, 2), cmp.Natural[int]())
		want := opt.Of(3)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestCount(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := Count(Empty[int]())
		want := int64(0)
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := Count(Of(1, 2, 3))
		want := int64(3)
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}
