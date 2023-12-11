package stream

import (
	"fmt"
	"github.com/jpfourny/papaya/pkg/mapper"
	"github.com/jpfourny/papaya/pkg/pair"
	"github.com/jpfourny/papaya/pkg/stream"
	"testing"
)

func TestMapIntToString(t *testing.T) {
	// Map stream of int to string.
	s := stream.Map(
		stream.Of(1, 2, 3),
		mapper.Sprintf[int]("%d"),
	)
	stream.ForEach(s, func(s string) {
		fmt.Println(s)
	})
	// Output:
	// 1
	// 2
	// 3
}

func TestMapPairFirst(t *testing.T) {
	// Given stream of pairs, project the first element of each pair.
	s := stream.Map(
		stream.Of(
			pair.Of(1, "foo"),
			pair.Of(2, "bar"),
		),
		pair.Pair[int, string].First,
	)
	stream.ForEach(s, func(i int) {
		fmt.Println(i)
	})
	// Output:
	// 1
	// 2
}

func TestMapPairSecond(t *testing.T) {
	// Given stream of pairs, project the second element of each pair.
	s := stream.Map(
		stream.Of(
			pair.Of(1, "foo"),
			pair.Of(2, "bar"),
		),
		pair.Pair[int, string].Second,
	)
	stream.ForEach(s, func(s string) {
		fmt.Println(s)
	})
	// Output:
	// foo
	// bar
}

func TestMapPairReverse(t *testing.T) {
	// Given stream of pairs, reverse each pair.
	s := stream.Map(
		stream.Of(
			pair.Of(1, "foo"),
			pair.Of(2, "bar"),
		),
		pair.Pair[int, string].Reverse,
	)
	stream.ForEach(s, func(p pair.Pair[string, int]) {
		fmt.Println(p)
	})
	// Output:
	// ("foo", 1)
	// ("bar", 2)
}
