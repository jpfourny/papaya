package stream

import (
	"fmt"
	"testing"

	"github.com/jpfourny/papaya/examples"
	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/stream"
)

func TestSortPersonComparing(t *testing.T) {
	// Sort by LastName, then by FirstName.
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

func TestSortPersonCompare(t *testing.T) {
	// Sort by Person.Compare.
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

func TestSortPersonCompareReverse(t *testing.T) {
	// Sort by reverse of Person.Compare.
	s := stream.SortBy(
		stream.FromSlice(examples.People),
		cmp.Comparer[examples.Person](examples.Person.Compare).Reverse(),
	)
	stream.ForEach(s, func(p examples.Person) {
		fmt.Println(p)
	})
	// Output:
	// Person{FirstName:"John", LastName:"Smith"}
	// Person{FirstName:"John", LastName:"Doe"}
	// Person{FirstName:"Jane", LastName:"Doe"}
}
