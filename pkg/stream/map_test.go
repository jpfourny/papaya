package stream

import (
	"fmt"
	"github.com/jpfourny/papaya/pkg/stream/mapper"
	"testing"

	"github.com/jpfourny/papaya/internal/assert"
)

func TestMap(t *testing.T) {
	s := Map(Of(1, 2, 3), mapper.Sprintf[int]("%d"))
	got := CollectSlice(s)
	want := []string{"1", "2", "3"}
	assert.ElementsMatch(t, got, want)
}

func TestMapOrDiscard(t *testing.T) {
	s := MapOrDiscard(Of("1", "2", "3", "foo"), mapper.TryParseInt[string, int](10, 64))
	got := CollectSlice(s)
	want := []int{1, 2, 3}
	assert.ElementsMatch(t, got, want)
}

func TestFlatMap(t *testing.T) {
	s := FlatMap(Of(1, 2, 3), func(e int) Stream[string] {
		return Of(fmt.Sprintf("%dA", e), fmt.Sprintf("%dB", e))
	})
	got := CollectSlice(s)
	want := []string{"1A", "1B", "2A", "2B", "3A", "3B"}
	assert.ElementsMatch(t, got, want)
}

func TestFlatMapSlice(t *testing.T) {
	s := FlatMapSlice(Of(1, 2, 3), func(e int) []string {
		return []string{fmt.Sprintf("%dA", e), fmt.Sprintf("%dB", e)}
	})
	got := CollectSlice(s)
	want := []string{"1A", "1B", "2A", "2B", "3A", "3B"}
	assert.ElementsMatch(t, got, want)
}
