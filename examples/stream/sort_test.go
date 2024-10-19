package stream

import (
	"fmt"
	"testing"

	"github.com/jpfourny/papaya/v2/examples"
	"github.com/jpfourny/papaya/v2/pkg/cmp"
	"github.com/jpfourny/papaya/v2/pkg/stream"
)

func TestSort(t *testing.T) {
	// Sort an ordered type in ascending order (natural order).
	s := stream.SortAsc(stream.Of(3, 1, 2))
	stream.ForEach(s, func(i int) {
		fmt.Println(i)
	})
	// Output:
	// 1
	// 2
	// 3

	// Sort an ordered type in descending order (reverse of natural order).
	s = stream.SortDesc(stream.Of(3, 1, 2))
	stream.ForEach(s, func(i int) {
		fmt.Println(i)
	})
	// Output:
	// 3
	// 2
	// 1
}

func TestSortBy(t *testing.T) {
	// Sort by cmp.Natural[E]() order, explicitly.
	s := stream.SortBy(stream.Of(3, 1, 2), cmp.Natural[int]())
	stream.ForEach(s, func(i int) {
		fmt.Println(i)
	})
	// Output:
	// 1
	// 2
	// 3

	// Sort by cmp.Reverse[E]() order, explicitly.
	// Note: cmp.Reverse[E]() is shorthand for cmp.Natural[E]().Reverse().
	s = stream.SortBy(stream.Of(3, 1, 2), cmp.Reverse[int]())
	stream.ForEach(s, func(i int) {
		fmt.Println(i)
	})
	// Output:
	// 3
	// 2
	// 1

	// Sort strings in case-insensitive order.
	//s2 := stream.SortBy(stream.Of("foo", "Bar", "bar", "Baz"), cmp.StringCompareFold("bar"))
}

func TestSortByComparing(t *testing.T) {
	// Sort Person by LastName, then by FirstName.
	// Use cmp.Comparing to create a cmp.Comparer[Person] that extracts sort keys from LastName and FirstName fields.
	s := stream.SortBy(
		stream.FromSlice(examples.People),
		cmp.Comparing(func(p examples.Person) string { return p.LastName }).
			Then(cmp.Comparing(func(p examples.Person) string { return p.FirstName })),
	)
	stream.ForEach(s, func(p examples.Person) {
		fmt.Println(p)
	})
	// Output:
	// Person{FirstName:"Jane", LastName:"Doe"}
	// Person{FirstName:"John", LastName:"Doe"}
	// Person{FirstName:"John", LastName:"Smith"}
}

func TestSortByMethod(t *testing.T) {
	// Sort by existing Person.Compare method, which satisfies cmp.Comparer[Person] interface.
	s := stream.SortBy(
		stream.FromSlice(examples.People),
		examples.Person.Compare, // Implicitly converted to cmp.Comparer[Person]
	)
	stream.ForEach(s, func(p examples.Person) {
		fmt.Println(p)
	})
	// Output:
	// Person{FirstName:"Jane", LastName:"Doe"}
	// Person{FirstName:"John", LastName:"Doe"}
	// Person{FirstName:"John", LastName:"Smith"}
}

func TestSortBySlice(t *testing.T) {
	// Sort a slice of int elements in natural order.
	s := stream.SortBy(
		stream.Of(
			[]int{3, 1, 2},
			[]int{2, 1, 3},
			[]int{2, 1},
		),
		cmp.Slice[int](cmp.Natural[int]()), // []int sorted lexicographically, elements in natural order.
	)
	stream.ForEach(s, func(slice []int) {
		fmt.Println(slice)
	})
	// Output:
	// [2 1]
	// [2 1 3]
	// [3 1 2]
}
