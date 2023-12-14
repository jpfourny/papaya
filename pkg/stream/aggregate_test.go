package stream

import (
	"testing"

	"github.com/jpfourny/papaya/pkg/optional"
)

func TestReduce(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := Reduce(Empty[int](), func(a, b int) int {
			return a + b
		})
		want := optional.Empty[int]()
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := Reduce(Of(1, 2, 3), func(a, b int) int {
			return a + b
		})
		want := optional.Of(6)
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
		got := Sum(Empty[int]())
		want := 0
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := Sum(Of(1, 2, 3))
		want := 6
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

}

func TestSumInteger(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := SumInteger(Empty[int]())
		want := int64(0)
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := SumInteger(Of(1, 2, 3))
		want := int64(6)
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}

func TestSumUnsignedInteger(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := SumUnsignedInteger(Empty[uint]())
		want := uint64(0)
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := SumUnsignedInteger(Of[uint](1, 2, 3))
		want := uint64(6)
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}

func TestSumFloat(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := SumFloat(Empty[float64]())
		want := 0.0
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := SumFloat(Of(1.0, 2.0, 3.0))
		want := 6.0
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
		want := optional.Empty[int]()
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := Min(Of(3, 1, 2))
		want := optional.Of(1)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestMax(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := Max(Empty[int]())
		want := optional.Empty[int]()
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := Max(Of(1, 3, 2))
		want := optional.Of(3)
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
