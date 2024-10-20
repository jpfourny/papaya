package stream

import (
	"fmt"
	"github.com/jpfourny/papaya/v2/pkg/stream/mapper"
	"testing"

	"github.com/jpfourny/papaya/v2/pkg/pair"
	"github.com/jpfourny/papaya/v2/pkg/stream"
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

func TestMapStringToInt(t *testing.T) {
	// Map stream of string to int; default to -1 when parse fails.
	s := stream.Map(
		stream.Of("1", "2", "3", "foo"),
		mapper.ParseIntOr[string, int](10, 64, -1),
	)
	stream.ForEach(s, func(i int) {
		fmt.Println(i)
	})
	// Output:
	// 1
	// 2
	// 3
	// -1

	// Same as above, but discard the parse failures.
	s = stream.MapOrDiscard(
		stream.Of("1", "2", "3", "foo"),
		mapper.TryParseInt[string, int](10, 64),
	)
	stream.ForEach(s, func(i int) {
		fmt.Println(i)
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

func TestMapNumToBool(t *testing.T) {
	// Given stream of numbers, map each number to whether it is even.
	s := stream.Map(
		stream.Of(0, 2, 0),
		mapper.NumToBool[int, bool](),
	)
	stream.ForEach(s, func(b bool) {
		fmt.Println(b)
	})
	// Output:
	// false
	// true
	// false
}

func TestMapBoolNumber(t *testing.T) {
	// Given stream of bools, map each bool to 0 if false and 1 if true.
	s := stream.Map(
		stream.Of(false, true, false),
		mapper.BoolToNum[bool, int](),
	)
	stream.ForEach(s, func(i int) {
		fmt.Println(i)
	})
	// Output:
	// 0
	// 1
	// 0
}
