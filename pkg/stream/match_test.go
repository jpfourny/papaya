package stream

import (
	"github.com/jpfourny/papaya/v2/pkg/cmp"
	"testing"
)

func TestContains(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := Contains(Empty[int](), 1)
		if got != false {
			t.Errorf("got %#v, want %#v", got, false)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := Contains(Of(1, 2, 3), 2)
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}
	})
}

func TestContainsBy(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := ContainsBy(Empty[int](), cmp.Natural[int](), 1)
		if got != false {
			t.Errorf("got %#v, want %#v", got, false)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := ContainsBy(Of(1, 2, 3), cmp.Natural[int](), 2)
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}
	})
}

func TestContainsAny(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := ContainsAny(Empty[int](), 1, 2)
		if got != false {
			t.Errorf("got %#v, want %#v", got, false)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := ContainsAny(Of(1, 2, 3), 2, 4)
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}

		got = ContainsAny(Of(1, 2, 3), 5)
		if got != false {
			t.Errorf("got %#v, want %#v", got, true)
		}
	})
}

func TestContainsAnyBy(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := ContainsAnyBy(Empty[int](), cmp.Natural[int](), 1, 2)
		if got != false {
			t.Errorf("got %#v, want %#v", got, false)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := ContainsAnyBy(Of(1, 2, 3), cmp.Natural[int](), 2, 4)
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}

		got = ContainsAnyBy(Of(1, 2, 3), cmp.Natural[int](), 5)
		if got != false {
			t.Errorf("got %#v, want %#v", got, true)
		}
	})
}

func TestContainsNone(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := ContainsNone(Empty[int](), 1, 2)
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := ContainsNone(Of(1, 2, 3), 2, 4)
		if got != false {
			t.Errorf("got %#v, want %#v", got, false)
		}

		got = ContainsNone(Of(1, 2, 3), 5)
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}
	})
}

func TestContainsNoneBy(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := ContainsNoneBy(Empty[int](), cmp.Natural[int](), 1, 2)
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := ContainsNoneBy(Of(1, 2, 3), cmp.Natural[int](), 2, 4)
		if got != false {
			t.Errorf("got %#v, want %#v", got, false)
		}

		got = ContainsNoneBy(Of(1, 2, 3), cmp.Natural[int](), 5)
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}
	})
}

func TestContainsAll(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := ContainsAll(Empty[int](), 1, 2)
		if got != false {
			t.Errorf("got %#v, want %#v", got, false)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := ContainsAll(Of(1, 2, 3), 2, 1)
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}

		got = ContainsAll(Of(1, 2, 3), 2, 1, 1) // duplicate should be ignored.
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}

		got = ContainsAll(Of(1, 2, 3), 3, 5)
		if got != false {
			t.Errorf("got %#v, want %#v", got, false)
		}
	})
}

func TestContainsAllBy(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := ContainsAllBy(Empty[int](), cmp.Natural[int](), 1, 2)
		if got != false {
			t.Errorf("got %#v, want %#v", got, false)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := ContainsAllBy(Of(1, 2, 3), cmp.Natural[int](), 2, 1)
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}

		got = ContainsAllBy(Of(1, 2, 3), cmp.Natural[int](), 2, 1, 1) // duplicate should be ignored.
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}

		got = ContainsAllBy(Of(1, 2, 3), cmp.Natural[int](), 3, 5)
		if got != false {
			t.Errorf("got %#v, want %#v", got, false)
		}
	})
}

func TestAnyMatch(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := AnyMatch(Empty[int](), func(e int) bool {
			return true
		})
		if got != false {
			t.Errorf("got %#v, want %#v", got, false)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Run("match", func(t *testing.T) {
			got := AnyMatch(Of(1, 2, 3), func(e int) bool {
				return e%2 == 0
			})
			if got != true {
				t.Errorf("got %#v, want %#v", got, true)
			}
		})

		t.Run("no-match", func(t *testing.T) {
			got := AnyMatch(Of(1, 2, 3), func(e int) bool {
				return e == 4
			})
			if got != false {
				t.Errorf("got %#v, want %#v", got, false)
			}
		})
	})
}

func TestAllMatch(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := AllMatch(Empty[int](), func(e int) bool {
			return true
		})
		if got != false {
			t.Errorf("got %#v, want %#v", got, false)
		}
	})

	got := AllMatch(Of(1, 2, 3), func(e int) bool {
		return e%2 == 0
	})
	if got != false {
		t.Errorf("got %#v, want %#v", got, false)
	}
	got = AllMatch(Of(1, 2, 3), func(e int) bool {
		return e <= 3
	})
	if got != true {
		t.Errorf("got %#v, want %#v", got, true)
	}
}

func TestNoneMatch(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := NoneMatch(Empty[int](), func(e int) bool {
			return true
		})
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Run("match", func(t *testing.T) {
			got := NoneMatch(Of(1, 2, 3), func(e int) bool {
				return e%2 == 0
			})
			if got != false {
				t.Errorf("got %#v, want %#v", got, false)
			}
		})

		t.Run("no-match", func(t *testing.T) {
			got := NoneMatch(Of(1, 2, 3), func(e int) bool {
				return e == 4
			})
			if got != true {
				t.Errorf("got %#v, want %#v", got, true)
			}
		})
	})
}

func TestExactlySame(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := ExactlySame(Empty[int](), Empty[int]())
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Run("same", func(t *testing.T) {
			got := ExactlySame(Of(1, 2, 3), Of(1, 2, 3))
			if got != true {
				t.Errorf("got %#v, want %#v", got, true)
			}
		})

		t.Run("different", func(t *testing.T) {
			got := ExactlySame(Of(1, 2, 3), Of(1, 2, 4))
			if got != false {
				t.Errorf("got %#v, want %#v", got, false)
			}
		})
	})
}

func TestExactlySameBy(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := ExactlySameBy(Empty[int](), Empty[int](), cmp.Natural[int]())
		if got != true {
			t.Errorf("got %#v, want %#v", got, true)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Run("same", func(t *testing.T) {
			got := ExactlySameBy(Of(1, 2, 3), Of(1, 2, 3), cmp.Natural[int]())
			if got != true {
				t.Errorf("got %#v, want %#v", got, true)
			}
		})

		t.Run("different", func(t *testing.T) {
			got := ExactlySameBy(Of(1, 2, 3), Of(1, 2, 4), cmp.Natural[int]())
			if got != false {
				t.Errorf("got %#v, want %#v", got, false)
			}
		})
	})
}

func TestIndexOfMatches(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := IndexOfMatches(Empty[int](), func(e int) bool {
			return true
		})
		if !ExactlySame(got, Empty[int64]()) {
			t.Errorf("expected empty stream")
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Run("match", func(t *testing.T) {
			got := IndexOfMatches(Of(1, 2, 3, 4), func(e int) bool {
				return e%2 == 0
			})
			if !ExactlySame(got, Of[int64](1, 3)) {
				t.Errorf("expected stream <1, 3>")
			}
		})

		t.Run("no-match", func(t *testing.T) {
			got := IndexOfMatches(Of(1, 2, 3), func(e int) bool {
				return e == 4
			})
			if !ExactlySame(got, Empty[int64]()) {
				t.Errorf("expected empty stream")
			}
		})
	})
}

func TestIndexOfFirstMatch(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := IndexOfFirstMatch(Empty[int](), func(e int) bool {
			return true
		})
		if got.Present() { // Expected empty
			t.Errorf("expected empty, got %#v", got)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Run("match", func(t *testing.T) {
			got := IndexOfFirstMatch(Of(1, 2, 3, 4), func(e int) bool {
				return e%2 == 0
			})
			if got.GetOrZero() != 1 {
				t.Errorf("got %#v, want %#v", got.GetOrZero(), 2)
			}
		})

		t.Run("no-match", func(t *testing.T) {
			got := IndexOfFirstMatch(Of(1, 2, 3), func(e int) bool {
				return e == 4
			})
			if got.Present() { // Expected empty
				t.Errorf("expected empty, got %#v", got)
			}
		})
	})
}

func TestIndexOfLastMatch(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := IndexOfLastMatch(Empty[int](), func(e int) bool {
			return true
		})
		if got.Present() { // Expected empty
			t.Errorf("expected empty, got %#v", got)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Run("match", func(t *testing.T) {
			got := IndexOfLastMatch(Of(1, 2, 3, 4), func(e int) bool {
				return e%2 == 0
			})
			if got.GetOrZero() != 3 {
				t.Errorf("got %#v, want %#v", got.GetOrZero(), 3)
			}
		})

		t.Run("no-match", func(t *testing.T) {
			got := IndexOfLastMatch(Of(1, 2, 3), func(e int) bool {
				return e == 4
			})
			if got.Present() { // Expected empty
				t.Errorf("expected empty, got %#v", got)
			}
		})
	})
}
