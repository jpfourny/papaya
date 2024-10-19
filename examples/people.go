package examples

import (
	"fmt"

	"github.com/jpfourny/papaya/v2/pkg/cmp"
)

type Person struct {
	FirstName string
	LastName  string
}

func (p Person) Compare(other Person) int {
	// Compare by LastName, then by FirstName.
	return cmp.Comparing(func(p Person) string { return p.LastName }).
		Then(cmp.Comparing(func(p Person) string { return p.FirstName }))(p, other)
}

func (p Person) String() string {
	return fmt.Sprintf("Person{FirstName:%q, LastName:%q}", p.FirstName, p.LastName)
}

var People = []Person{
	{"John", "Doe"},
	{"Jane", "Doe"},
	{"John", "Smith"},
}
