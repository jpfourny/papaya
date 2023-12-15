package stream

import (
	"reflect"
	"testing"

	"github.com/jpfourny/papaya/pkg/pair"
)

func TestGroupByKey(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := GroupByKey(Empty[pair.Pair[int, string]]())
		got := CollectMap(s)
		if len(got) != 0 {
			t.Fatalf("got %#v, want %#v", got, map[int][]string{})
		}
		if !reflect.DeepEqual(got[0], []string(nil)) {
			t.Fatalf("got %#v, want %#v", got[0], []string{})
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		s := GroupByKey(Of(
			pair.Of(1, "one"),
			pair.Of(2, "two"),
			pair.Of(1, "uno"),
			pair.Of(3, "three"),
			pair.Of(2, "dos"),
		))
		got := CollectMap(s)
		want := map[int][]string{
			1: {"one", "uno"},
			2: {"two", "dos"},
			3: {"three"},
		}
		if len(got) != len(want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		for k, v := range want {
			if !reflect.DeepEqual(got[k], v) {
				t.Fatalf("got %#v, want %#v", got[k], want[k])
			}
		}
	})

	t.Run("limited", func(t *testing.T) {
		s := GroupByKey(Of(
			pair.Of(1, "one"),
			pair.Of(2, "two"),
			pair.Of(1, "uno"),
			pair.Of(3, "three"),
			pair.Of(2, "dos"),
		))
		got := CollectMap(Limit(s, 1)) // Stops stream after 1 elements.
		if len(got) != 1 {
			t.Fatal("expected 1 element; got", len(got)) // Actual value is unpredictable due to map iteration order.
		}
	})
}

func TestReduceByKey(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := ReduceByKey(
			Empty[pair.Pair[int, string]](),
			func(a, b string) string { return a + ", " + b },
		)
		got := CollectMap(s)
		if len(got) != 0 {
			t.Fatalf("got %#v, want %#v", got, map[int]string{})
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		s := ReduceByKey(
			Of(
				pair.Of(1, "one"),
				pair.Of(2, "two"),
				pair.Of(1, "uno"),
				pair.Of(3, "three"),
				pair.Of(2, "dos"),
			),
			func(a, b string) string { return a + ", " + b },
		)
		got := CollectMap(s)
		want := map[int]string{
			1: "one, uno",
			2: "two, dos",
			3: "three",
		}
		if len(got) != len(want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		for k, v := range want {
			if got[k] != v {
				t.Fatalf("got %#v, want %#v", got[k], want[k])
			}
		}
	})

	t.Run("limited", func(t *testing.T) {
		s := ReduceByKey(
			Of(
				pair.Of(1, "one"),
				pair.Of(2, "two"),
				pair.Of(1, "uno"),
				pair.Of(3, "three"),
				pair.Of(2, "dos"),
			),
			func(a, b string) string { return a + ", " + b },
		)
		got := CollectMap(Limit(s, 1)) // Stops stream after 1 elements.
		if len(got) != 1 {
			t.Fatal("expected 1 element; got", len(got)) // Actual value is unpredictable due to map iteration order.
		}
	})
}

func TestAggregateByKey(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := AggregateByKey(
			Empty[pair.Pair[int, string]](),
			"",
			func(a string, e string) string {
				if a == "" {
					return e
				}
				return a + ", " + e
			}, func(a string) string {
				return "<" + a + ">"
			})
		got := CollectMap(s)
		if len(got) != 0 {
			t.Fatalf("got %#v, want %#v", got, map[int]string{})
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		s := AggregateByKey(
			Of(
				pair.Of(1, "one"),
				pair.Of(2, "two"),
				pair.Of(1, "uno"),
				pair.Of(3, "three"),
				pair.Of(2, "dos"),
			),
			"",
			func(a string, e string) string {
				if a == "" {
					return e
				}
				return a + ", " + e
			},
			func(a string) string {
				return "<" + a + ">"
			},
		)
		got := CollectMap(s)
		want := map[int]string{
			1: "<one, uno>",
			2: "<two, dos>",
			3: "<three>",
		}
		if len(got) != len(want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		for k, v := range want {
			if got[k] != v {
				t.Fatalf("got %#v, want %#v", got[k], want[k])
			}
		}
	})

	t.Run("limited", func(t *testing.T) {
		s := AggregateByKey(
			Of(
				pair.Of(1, "one"),
				pair.Of(2, "two"),
				pair.Of(1, "uno"),
				pair.Of(3, "three"),
				pair.Of(2, "dos"),
			),
			"",
			func(a string, e string) string {
				if a == "" {
					return e
				}
				return a + ", " + e
			},
			func(a string) string {
				return "<" + a + ">"
			},
		)
		got := CollectMap(Limit(s, 1)) // Stops stream after 1 elements.
		if len(got) != 1 {
			t.Fatal("expected 1 element; got", len(got)) // Actual value is unpredictable due to map iteration order.
		}
	})
}

func TestCountByKey(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := CountByKey(Empty[pair.Pair[int, string]]())
		got := CollectMap(s)
		if len(got) != 0 {
			t.Fatalf("got %#v, want %#v", got, map[int]int{})
		}
		if got[0] != 0 {
			t.Fatalf("got %#v, want %#v", got[0], 0)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		s := CountByKey(Of(
			pair.Of(1, "one"),
			pair.Of(2, "two"),
			pair.Of(1, "uno"),
			pair.Of(3, "three"),
			pair.Of(2, "dos"),
		))
		got := CollectMap(s)
		want := map[int]int64{
			1: 2,
			2: 2,
			3: 1,
		}
		if len(got) != len(want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		for k, v := range want {
			if got[k] != v {
				t.Fatalf("got %#v, want %#v", got[k], want[k])
			}
		}
	})

	t.Run("limited", func(t *testing.T) {
		s := CountByKey(Of(
			pair.Of(1, "one"),
			pair.Of(2, "two"),
			pair.Of(1, "uno"),
			pair.Of(3, "three"),
			pair.Of(2, "dos"),
		))
		got := CollectMap(Limit(s, 1)) // Stops stream after 1 elements.
		if len(got) != 1 {
			t.Fatal("expected 1 element; got", len(got)) // Actual value is unpredictable due to map iteration order.
		}
	})
}