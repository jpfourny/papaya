package stream

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"sync"
	"testing"

	"github.com/jpfourny/papaya/pkg/mapper"
	"github.com/jpfourny/papaya/pkg/optional"
	"github.com/jpfourny/papaya/pkg/pair"
)

func TestEmpty(t *testing.T) {
	s := Empty[int]()
	got := CollectSlice(s)
	var want []int
	assertElementsMatch(t, got, want)
}

func TestOf(t *testing.T) {
	s := Of(1, 2, 3)
	got := CollectSlice(s)
	want := []int{1, 2, 3}
	assertElementsMatch(t, got, want)
}

type emptyIter struct{}

func (_ emptyIter) Next() (int, bool) {
	return 0, false
}

type sliceIter struct {
	slice []int
	index int
}

func (s *sliceIter) Next() (e int, ok bool) {
	if s.index >= len(s.slice) {
		return
	}
	e, ok = s.slice[s.index], true
	s.index++
	return
}

func TestIterator(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := FromIterator[int](emptyIter{})
		got := CollectSlice(s)
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := FromIterator[int](&sliceIter{slice: []int{1, 2, 3}})
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		assertElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := FromIterator[int](&sliceIter{slice: []int{1, 2, 3}})
		got := CollectSlice(Limit(s, 1)) // Stops stream after 1 element.
		want := []int{1}
		assertElementsMatch(t, got, want)
	})
}

func TestFromSlice(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := FromSlice([]int{})
		got := CollectSlice(s)
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := FromSlice([]int{1, 2, 3})
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		assertElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := FromSlice([]int{1, 2, 3})
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []int{1, 2}
		assertElementsMatch(t, got, want)
	})
}

func TestFromSliceWithIndex(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := FromSliceWithIndex([]int{})
		got := CollectSlice(s)
		var want []pair.Pair[int, int]
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := FromSliceWithIndex([]int{1, 2, 3})
		got := CollectSlice(s)
		want := []pair.Pair[int, int]{
			pair.Of(0, 1),
			pair.Of(1, 2),
			pair.Of(2, 3),
		}
		dd := DebugString(s)
		fmt.Println(dd)
		assertElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := FromSliceWithIndex([]int{1, 2, 3})
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []pair.Pair[int, int]{
			pair.Of(0, 1),
			pair.Of(1, 2),
		}
		assertElementsMatch(t, got, want)
	})
}

func TestFromMap(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := FromMap(map[int]string{})
		got := CollectSlice(s)
		var want []pair.Pair[int, string]
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := FromMap(map[int]string{1: "one", 2: "two", 3: "three"})
		got := CollectSlice(s)
		want := []pair.Pair[int, string]{
			pair.Of(1, "one"),
			pair.Of(2, "two"),
			pair.Of(3, "three"),
		}
		// All elements in map iteration order.
		assertElementsMatchAnyOrder(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := FromMap(map[int]string{1: "one", 2: "two", 3: "three"})
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []pair.Pair[int, string]{
			pair.Of(1, "one"),
			pair.Of(2, "two"),
			pair.Of(3, "three"),
		}
		// 2 elements returned are unpredictable due to map iteration order.
		assertSomeElementsMatchAnyOrder(t, got, want, 2)
	})
}

func TestFromMapKeys(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := FromMapKeys(map[int]string{})
		got := CollectSlice(s)
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := FromMapKeys(map[int]string{1: "one", 2: "two", 3: "three"})
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		// All elements in map iteration order.
		assertElementsMatchAnyOrder(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := FromMapKeys(map[int]string{1: "one", 2: "two", 3: "three"})
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []int{1, 2, 3}
		// 2 elements returned are unpredictable due to map iteration order.
		assertSomeElementsMatchAnyOrder(t, got, want, 2)
	})
}

func TestFromMapValues(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := FromMapValues(map[int]string{})
		got := CollectSlice(s)
		var want []string
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := FromMapValues(map[int]string{1: "one", 2: "two", 3: "three"})
		got := CollectSlice(s)
		want := []string{"one", "two", "three"}
		// All elements in map iteration order.
		assertElementsMatchAnyOrder(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := FromMapValues(map[int]string{1: "one", 2: "two", 3: "three"})
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []string{"one", "two", "three"}
		// 2 elements returned are unpredictable due to map iteration order.
		assertSomeElementsMatchAnyOrder(t, got, want, 2)
	})
}

func TestFromChannel(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ch := make(chan int)
		close(ch)
		s := FromChannel(ch)
		got := CollectSlice(s)
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Run("unbuffered", func(t *testing.T) {
			ch := make(chan int)
			go func() {
				ch <- 1
				ch <- 2
				ch <- 3
				close(ch)
			}()
			s := FromChannel(ch)
			got := CollectSlice(s)
			want := []int{1, 2, 3}
			assertElementsMatch(t, got, want)
		})

		t.Run("buffered", func(t *testing.T) {
			ch := make(chan int, 3)
			go func() {
				ch <- 1
				ch <- 2
				ch <- 3
				close(ch)
			}()
			s := FromChannel(ch)
			got := CollectSlice(s)
			want := []int{1, 2, 3}
			assertElementsMatch(t, got, want)
		})
	})

	t.Run("limited", func(t *testing.T) {
		ch := make(chan int, 3)
		go func() {
			ch <- 1
			ch <- 2
			ch <- 3
			close(ch)
		}()
		s := FromChannel(ch)
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []int{1, 2}
		assertElementsMatch(t, got, want)
	})
}

func TestFromChannelCtx(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ch := make(chan int)
		go close(ch)
		s := FromChannelCtx(context.Background(), ch)
		got := CollectSlice(s)
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		ch := make(chan int, 3)
		go cancel()
		s := FromChannelCtx(ctx, ch)
		got := CollectSlice(s)
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Run("unbuffered", func(t *testing.T) {
			ch := make(chan int)
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				ch <- 1
				ch <- 2
				ch <- 3
				close(ch)
			}()
			s := FromChannelCtx(context.Background(), ch)
			got := CollectSlice(s)
			want := []int{1, 2, 3}
			assertElementsMatch(t, got, want)
		})
		t.Run("buffered", func(t *testing.T) {
			ch := make(chan int, 3)
			go func() {
				ch <- 1
				ch <- 2
				ch <- 3
				close(ch)
			}()
			s := FromChannelCtx(context.Background(), ch)
			got := CollectSlice(s)
			want := []int{1, 2, 3}
			assertElementsMatch(t, got, want)
		})
	})

	t.Run("limited", func(t *testing.T) {
		ch := make(chan int, 3)
		go func() {
			ch <- 1
			ch <- 2
			ch <- 3
			close(ch)
		}()
		s := FromChannelCtx(context.Background(), ch)
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []int{1, 2}
		assertElementsMatch(t, got, want)
	})
}

func TestRangeInteger(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := RangeInteger(0, 0)
		got := CollectSlice(s)
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := RangeInteger(1, 4)
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		assertElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := RangeInteger(1, 4)
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []int{1, 2}
		assertElementsMatch(t, got, want)
	})
}

func TestRangeIntegerStep(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := RangeIntegerStep(0, 0, 1)
		got := CollectSlice(s)
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := RangeIntegerStep(1, 4, 2)
		got := CollectSlice(s)
		want := []int{1, 3}
		assertElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := RangeIntegerStep(1, 4, 2)
		got := CollectSlice(Limit(s, 1)) // Stops stream after 1 element.
		want := []int{1}
		assertElementsMatch(t, got, want)
	})
}

func TestRangeFloatStep(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := RangeFloatStep(0.0, 0.0, 1.0)
		got := CollectSlice(s)
		var want []float64
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := RangeFloatStep(1.0, 4.0, 2.0)
		got := CollectSlice(s)
		want := []float64{1.0, 3.0}
		assertElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := RangeFloatStep(1.0, 4.0, 2.0)
		got := CollectSlice(Limit(s, 1)) // Stops stream after 1 element.
		want := []float64{1.0}
		assertElementsMatch(t, got, want)
	})
}

func TestUnion(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := Union[int]()
		got := CollectSlice(s)
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := Union(Of(1, 2, 3), Of(4, 5, 6))
		got := CollectSlice(s)
		want := []int{1, 2, 3, 4, 5, 6}
		assertElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := Union(Of(1, 2, 3), Of(4, 5, 6))
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []int{1, 2}
		assertElementsMatch(t, got, want)
	})
}

func TestIntersection(t *testing.T) {
	t.Run("zero-streams", func(t *testing.T) {
		s := Intersection[int]()
		got := CollectSlice(s)
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("one-stream", func(t *testing.T) {
		s := Intersection(Of(1, 2, 3))
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		assertElementsMatch(t, got, want)
	})

	t.Run("two-streams", func(t *testing.T) {
		s := Intersection(Of(1, 2, 3), Of(2, 3, 4))
		got := CollectSlice(s)
		want := []int{2, 3}
		assertElementsMatch(t, got, want)
	})

	t.Run("three-streams", func(t *testing.T) {
		s := Intersection(Of(1, 2, 3), Of(2, 3, 4), Of(3, 4, 5))
		got := CollectSlice(s)
		want := []int{3}
		assertElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := Intersection(Of(1, 2, 3), Of(2, 3, 4))
		got := CollectSlice(Limit(s, 1)) // Stops stream after 1 element.
		want := []int{2}
		assertElementsMatch(t, got, want)
	})
}

func TestSortAsc(t *testing.T) {
	s := SortAsc(Of(3, 1, 2))
	got := CollectSlice(s)
	want := []int{1, 2, 3}
	assertElementsMatch(t, got, want)
}

func TestSortDesc(t *testing.T) {
	s := SortDesc(Of(3, 1, 2))
	got := CollectSlice(s)
	want := []int{3, 2, 1}
	assertElementsMatch(t, got, want)
}

func TestSortBy(t *testing.T) {
	s := SortBy(Of(3, 1, 2), func(a, b int) int {
		return a - b
	})
	got := CollectSlice(s)
	want := []int{1, 2, 3}
	assertElementsMatch(t, got, want)
}

func TestDistinct(t *testing.T) {
	s := Distinct(Of(1, 2, 3, 2, 1))
	got := CollectSlice(s)
	want := []int{1, 2, 3}
	assertElementsMatch(t, got, want)
}

func TestDistinctBy(t *testing.T) {
	s := DistinctBy(Of(1, 2, 3, 2, 1), func(e int) int {
		return e % 2
	})
	got := CollectSlice(s)
	want := []int{1, 2}
	assertElementsMatch(t, got, want)
}

func TestForEach(t *testing.T) {
	var got []int
	ForEach(Of(1, 2, 3), func(e int) {
		got = append(got, e)
	})
	want := []int{1, 2, 3}
	assertElementsMatch(t, got, want)
}

func TestPeek(t *testing.T) {
	var got []int
	s := Peek(Of(1, 2, 3), func(e int) {
		got = append(got, e)
	})
	Count(s) // Force evaluation so peek is called.
	want := []int{1, 2, 3}
	assertElementsMatch(t, got, want)
	assertElementsMatch(t, CollectSlice(s), want)
}

func TestMap(t *testing.T) {
	s := Map(Of(1, 2, 3), mapper.Sprintf[int]("%d"))
	got := CollectSlice(s)
	want := []string{"1", "2", "3"}
	assertElementsMatch(t, got, want)
}

func TestFlatMap(t *testing.T) {
	s := FlatMap(Of(1, 2, 3), func(e int) Stream[string] {
		return Of(fmt.Sprintf("%dA", e), fmt.Sprintf("%dB", e))
	})
	got := CollectSlice(s)
	want := []string{"1A", "1B", "2A", "2B", "3A", "3B"}
	assertElementsMatch(t, got, want)
}

func TestLimit(t *testing.T) {
	t.Run("limit-0", func(t *testing.T) {
		s := Limit(Of(1, 2, 3, 4, 5), 0)
		got := CollectSlice(s)
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("limit-negative", func(t *testing.T) {
		s := Limit(Of(1, 2, 3, 4, 5), -1) // Will be treated as 0.
		got := CollectSlice(s)
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("limit-positive", func(t *testing.T) {
		s := Limit(Of(1, 2, 3, 4, 5), 3)
		got := CollectSlice(s)
		want := []int{1, 2, 3}
		assertElementsMatch(t, got, want)
	})
}

func TestSkip(t *testing.T) {
	t.Run("skip-0", func(t *testing.T) {
		s := Skip(Of(1, 2, 3, 4, 5), 0)
		got := CollectSlice(s)
		want := []int{1, 2, 3, 4, 5}
		assertElementsMatch(t, got, want)
	})

	t.Run("skip-negative", func(t *testing.T) {
		s := Skip(Of(1, 2, 3, 4, 5), -1) // Will be treated as 0.
		got := CollectSlice(s)
		want := []int{1, 2, 3, 4, 5}
		assertElementsMatch(t, got, want)
	})

	t.Run("skip-positive", func(t *testing.T) {
		s := Skip(Of(1, 2, 3, 4, 5), 3)
		got := CollectSlice(s)
		want := []int{4, 5}
		assertElementsMatch(t, got, want)
	})

	t.Run("skip-all", func(t *testing.T) {
		s := Skip(Of(1, 2, 3, 4, 5), 6)
		got := CollectSlice(s)
		var want []int
		assertElementsMatch(t, got, want)
	})
}

func TestFilter(t *testing.T) {
	s := Filter(Of(1, 2, 3, 4, 5), func(e int) bool {
		return e%2 == 0
	})
	got := CollectSlice(s)
	want := []int{2, 4}
	assertElementsMatch(t, got, want)
}

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

func TestCollectSlice(t *testing.T) {
	got := CollectSlice(Of(1, 2, 3))
	want := []int{1, 2, 3}
	assertElementsMatch(t, got, want)
}

func TestCollectMap(t *testing.T) {
	got := CollectMap(Of(
		pair.Of(1, "one"),
		pair.Of(2, "two"),
		pair.Of(3, "three"),
	))
	want := map[int]string{1: "one", 2: "two", 3: "three"}
	if len(got) != len(want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("got[%d] %#v, want %#v", k, got, want)
		}
	}
}

func TestCollectChannel(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ch := make(chan int)
		go func() {
			CollectChannel(Empty[int](), ch)
			close(ch)
		}()
		got := CollectSlice(FromChannel(ch))
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		ch := make(chan int)
		go func() {
			CollectChannel(Of(1, 2, 3), ch)
			close(ch)
		}()
		got := CollectSlice(FromChannel(ch))
		want := []int{1, 2, 3}
		assertElementsMatch(t, got, want)
	})
}

func TestCollectChannelCtx(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ch := make(chan int)
		go func() {
			CollectChannelCtx(context.Background(), Empty[int](), ch)
			close(ch)
		}()
		got := CollectSlice(FromChannel(ch))
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		ch := make(chan int)
		go func() {
			CollectChannelCtx(context.Background(), Of(1, 2, 3), ch)
			close(ch)
		}()
		got := CollectSlice(FromChannel(ch))
		want := []int{1, 2, 3}
		assertElementsMatch(t, got, want)
	})

	t.Run("cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		ch := make(chan int)
		cancel()
		go func() {
			CollectChannelCtx(ctx, Of(1, 2, 3), ch)
			close(ch)
		}()
		// Due to race condition, we may get 0 or 1 element before cancel is seen.
		got := CollectSlice(FromChannel(ch))
		if len(got) > 1 {
			t.Errorf("expected no more than 1 element, got %#v", got)
		}
	})
}

func TestCollectChannelAsync(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ch := CollectChannelAsync(Empty[int](), 0) // ch closed at end of stream.
		got := CollectSlice(FromChannel(ch))
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Run("unbuffered", func(t *testing.T) {
			ch := CollectChannelAsync(Of(1, 2, 3), 0) // ch closed at end of stream.
			got := CollectSlice(FromChannel(ch))
			want := []int{1, 2, 3}
			assertElementsMatch(t, got, want)
		})
		t.Run("buffered", func(t *testing.T) {
			ch := CollectChannelAsync(Of(1, 2, 3), 3) // ch closed at end of stream.
			got := CollectSlice(FromChannel(ch))
			want := []int{1, 2, 3}
			assertElementsMatch(t, got, want)
		})
	})
}

func TestCollectChannelAsyncCtx(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ch := CollectChannelAsyncCtx(context.Background(), Empty[int](), 0) // ch closed at end of stream.
		got := CollectSlice(FromChannel(ch))
		var want []int
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Run("unbuffered", func(t *testing.T) {
			ch := CollectChannelAsyncCtx(context.Background(), Of(1, 2, 3), 0) // ch closed at end of stream.
			got := CollectSlice(FromChannel(ch))
			want := []int{1, 2, 3}
			assertElementsMatch(t, got, want)
		})
		t.Run("buffered", func(t *testing.T) {
			ch := CollectChannelAsyncCtx(context.Background(), Of(1, 2, 3), 3) // ch closed at end of stream.
			got := CollectSlice(FromChannel(ch))
			want := []int{1, 2, 3}
			assertElementsMatch(t, got, want)
		})
	})

	t.Run("cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		ch := CollectChannelAsyncCtx(ctx, Of(1, 2, 3), 0) // ch closed at end of stream.
		// Due to race condition, we may get 0 or 1 element before cancel is seen.
		got := CollectSlice(FromChannel(ch))
		if len(got) > 1 {
			t.Errorf("expected no more than 1 element, got %#v", got)
		}
	})
}

func TestDebugString(t *testing.T) {
	got := DebugString(Of(1, 2, 3))
	want := "<1, 2, 3>"
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
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

func TestFirst(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got, ok := First(Empty[int]())
		want := 0
		if got != want || ok {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got, ok := First(Of(1, 2, 3))
		want := 1
		if got != want || !ok {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}

func TestLast(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got, ok := Last(Empty[int]())
		want := 0
		if got != want || ok {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got, ok := Last(Of(1, 2, 3))
		want := 3
		if got != want || !ok {
			t.Errorf("got %#v, want %#v", got, want)
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

func TestZip(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := Zip(Empty[int](), Empty[int]())
		got := CollectSlice(s)
		var want []pair.Pair[int, int]
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := Zip(Of(1, 2, 3), Of(4, 5, 6))
		got := CollectSlice(s)
		want := []pair.Pair[int, int]{
			pair.Of(1, 4),
			pair.Of(2, 5),
			pair.Of(3, 6),
		}
		assertElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := Zip(Of(1, 2, 3), Of(4, 5, 6))
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []pair.Pair[int, int]{
			pair.Of(1, 4),
			pair.Of(2, 5),
		}
		assertElementsMatch(t, got, want)
	})

	t.Run("different-length", func(t *testing.T) {
		s := Zip(Of(1, 2, 3), Of(4, 5))
		got := CollectSlice(s)
		want := []pair.Pair[int, int]{
			pair.Of(1, 4),
			pair.Of(2, 5),
		}
		assertElementsMatch(t, got, want)
	})
}

func TestZipWithIndexInt(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := ZipWithIndexInt(Empty[int](), 0)
		got := CollectSlice(s)
		var want []pair.Pair[int, int]
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := ZipWithIndexInt(Of(1, 2, 3), 0)
		got := CollectSlice(s)
		want := []pair.Pair[int, int]{
			pair.Of(1, 0),
			pair.Of(2, 1),
			pair.Of(3, 2),
		}
		assertElementsMatch(t, got, want)
	})

	t.Run("non-zero-offset", func(t *testing.T) {
		s := ZipWithIndexInt(Of(1, 2, 3), -1)
		got := CollectSlice(s)
		want := []pair.Pair[int, int]{
			pair.Of(1, -1),
			pair.Of(2, 0),
			pair.Of(3, 1),
		}
		assertElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := ZipWithIndexInt(Of(1, 2, 3), 0)
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []pair.Pair[int, int]{
			pair.Of(1, 0),
			pair.Of(2, 1),
		}
		assertElementsMatch(t, got, want)
	})
}

func TestZipWithIndexInt64(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := ZipWithIndexInt64(Empty[int](), 0)
		got := CollectSlice(s)
		var want []pair.Pair[int, int64]
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := ZipWithIndexInt64(Of(1, 2, 3), 0)
		got := CollectSlice(s)
		want := []pair.Pair[int, int64]{
			pair.Of(1, int64(0)),
			pair.Of(2, int64(1)),
			pair.Of(3, int64(2)),
		}
		assertElementsMatch(t, got, want)
	})

	t.Run("non-zero-offset", func(t *testing.T) {
		s := ZipWithIndexInt64(Of(1, 2, 3), -1)
		got := CollectSlice(s)
		want := []pair.Pair[int, int64]{
			pair.Of(1, int64(-1)),
			pair.Of(2, int64(0)),
			pair.Of(3, int64(1)),
		}
		assertElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := ZipWithIndexInt64(Of(1, 2, 3), 0)
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []pair.Pair[int, int64]{
			pair.Of(1, int64(0)),
			pair.Of(2, int64(1)),
		}
		assertElementsMatch(t, got, want)
	})
}

func TestZipWithKey(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := ZipWithKey(Empty[int](), func(e int) int {
			return e * 10
		})
		got := CollectSlice(s)
		var want []pair.Pair[int, int]
		assertElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		s := ZipWithKey(Of(1, 2, 3), func(e int) int {
			return e * 10
		})
		got := CollectSlice(s)
		want := []pair.Pair[int, int]{
			pair.Of(10, 1),
			pair.Of(20, 2),
			pair.Of(30, 3),
		}
		assertElementsMatch(t, got, want)
	})

	t.Run("limited", func(t *testing.T) {
		s := ZipWithKey(Of(1, 2, 3), func(e int) int {
			return e * 10
		})
		got := CollectSlice(Limit(s, 2)) // Stops stream after 2 elements.
		want := []pair.Pair[int, int]{
			pair.Of(10, 1),
			pair.Of(20, 2),
		}
		assertElementsMatch(t, got, want)
	})
}

func assertElementsMatch[E comparable](t *testing.T, got, want []E) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %#v, want %#v exactly", got, want)
	}
	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("got %#v, want %#v exactly", got, want)
		}
	}
}

func assertElementsMatchAnyOrder[E comparable](t *testing.T, got, want []E) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %3v, want all elements from %#v in any order", got, want)
	}
	for _, e := range got {
		if !slices.Contains(want, e) {
			t.Fatalf("got %#v, want all elements from %#v in any order", got, want)
		}
	}
}

func assertSomeElementsMatchAnyOrder[E comparable](t *testing.T, got, want []E, n int) {
	t.Helper()
	if len(got) != n {
		t.Fatalf("got %#v, want %d elements from %#v in any order", got, n, want)
	}
	for _, e := range got {
		if !slices.Contains(want, e) {
			t.Fatalf("got %#v, want %d elements from %#v in any order", got, n, want)
		}
	}
}
