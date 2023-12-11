package cmp

import (
	"github.com/jpfourny/papaya/pkg/pointer"
	"strings"
	"testing"
)

type Person struct {
	FirstName string
	LastName  string
}

func (p Person) Compare(other Person) int {
	if r := strings.Compare(p.LastName, other.LastName); r != 0 {
		return r
	}
	return strings.Compare(p.FirstName, other.FirstName)
}

func TestNatural(t *testing.T) {
	c := Natural[int]()
	if c(1, 2) != -1 {
		t.Errorf("Natural()(1, 2): expected -1, got %d", c(1, 2))

	}
	if c(2, 1) != 1 {
		t.Errorf("Natural()(2, 1): expected 1, got %d", c(2, 1))
	}
	if c(1, 1) != 0 {
		t.Errorf("Natural()(1, 1): expected 0, got %d", c(1, 1))
	}
}

func TestReverse(t *testing.T) {
	c := Reverse[int]()
	if c(1, 2) != 1 {
		t.Errorf("Reverse()(1, 2): expected 1, got %d", c(1, 2))
	}
	if c(2, 1) != -1 {
		t.Errorf("Reverse()(2, 1): expected -1, got %d", c(2, 1))
	}
	if c(1, 1) != 0 {
		t.Errorf("Reverse()(1, 1): expected 0, got %d", c(1, 1))
	}
}

func TestDerefNilFirst(t *testing.T) {
	c := DerefNilFirst(Natural[int]())
	got := c(nil, nil)
	want := 0
	if got != want {
		t.Errorf("DerefNilFirst()(nil, nil): expected %d, got %d", want, got)
	}

	got = c(nil, pointer.Ref(1))
	want = -1
	if got != want {
		t.Errorf("DerefNilFirst()(nil, 1): expected %d, got %d", want, got)
	}

	got = c(pointer.Ref(1), nil)
	want = 1
	if got != want {
		t.Errorf("DerefNilFirst()(1, nil): expected %d, got %d", want, got)
	}

	got = c(pointer.Ref(1), pointer.Ref(2))
	want = -1
	if got != want {
		t.Errorf("DerefNilFirst()(1, 2): expected %d, got %d", want, got)
	}

	got = c(pointer.Ref(2), pointer.Ref(1))
	want = 1
	if got != want {
		t.Errorf("DerefNilFirst()(2, 1): expected %d, got %d", want, got)
	}

	got = c(pointer.Ref(1), pointer.Ref(1))
	want = 0
	if got != want {
		t.Errorf("DerefNilFirst()(1, 1): expected %d, got %d", want, got)
	}
}

func TestDerefNilLast(t *testing.T) {
	c := DerefNilLast(Natural[int]())
	got := c(nil, nil)
	want := 0
	if got != want {
		t.Errorf("DerefNilLast()(nil, nil): expected %d, got %d", want, got)
	}

	got = c(nil, pointer.Ref(1))
	want = 1
	if got != want {
		t.Errorf("DerefNilLast()(nil, 1): expected %d, got %d", want, got)
	}

	got = c(pointer.Ref(1), nil)
	want = -1
	if got != want {
		t.Errorf("DerefNilLast()(1, nil): expected %d, got %d", want, got)
	}

	got = c(pointer.Ref(1), pointer.Ref(2))
	want = -1
	if got != want {
		t.Errorf("DerefNilLast()(1, 2): expected %d, got %d", want, got)
	}

	got = c(pointer.Ref(2), pointer.Ref(1))
	want = 1
	if got != want {
		t.Errorf("DerefNilLast()(2, 1): expected %d, got %d", want, got)
	}

	got = c(pointer.Ref(1), pointer.Ref(1))
	want = 0
	if got != want {
		t.Errorf("DerefNilLast()(1, 1): expected %d, got %d", want, got)
	}
}

func TestComparing(t *testing.T) {
	c := Comparing(func(p Person) string { return p.LastName })
	got := c(Person{"John", "Doe"}, Person{"Jane", "Doe"})
	want := 0
	if got != want {
		t.Errorf("Comparing(func(p Person) string { return p.LastName })(Person{FirstName:\"John\", LastName:\"Doe\"}, Person{FirstName:\"Jane\", LastName:\"Doe\"}): expected %d, got %d", want, got)
	}

	got = c(Person{"John", "Doe"}, Person{"John", "Smith"})
	want = -1
	if got != want {
		t.Errorf("Comparing(func(p Person) string { return p.LastName })(Person{FirstName:\"John\", LastName:\"Doe\"}, Person{FirstName:\"John\", LastName:\"Smith\"}): expected %d, got %d", want, got)
	}

	got = c(Person{"John", "Smith"}, Person{"John", "Doe"})
	want = 1
	if got != want {
		t.Errorf("Comparing(func(p Person) string { return p.LastName })(Person{FirstName:\"John\", LastName:\"Smith\"}, Person{FirstName:\"John\", LastName:\"Doe\"}): expected %d, got %d", want, got)
	}
}

func TestComparingBy(t *testing.T) {
	c := ComparingBy(func(p Person) string { return p.LastName }, Natural[string]())
	got := c(Person{"John", "Doe"}, Person{"Jane", "Doe"})
	want := 0
	if got != want {
		t.Errorf("ComparingBy(Person{}, func(p Person) string { return p.LastName })(Person{FirstName:\"John\", LastName:\"Doe\"}, Person{FirstName:\"Jane\", LastName:\"Doe\"}): expected %d, got %d", want, got)
	}

	got = c(Person{"John", "Doe"}, Person{"John", "Smith"})
	want = -1
	if got != want {
		t.Errorf("ComparingBy(Person{}, func(p Person) string { return p.LastName })(Person{FirstName:\"John\", LastName:\"Doe\"}, Person{FirstName:\"John\", LastName:\"Smith\"}): expected %d, got %d", want, got)
	}

	got = c(Person{"John", "Smith"}, Person{"John", "Doe"})
	want = 1
	if got != want {
		t.Errorf("ComparingBy(Person{}, func(p Person) string { return p.LastName })(Person{FirstName:\"John\", LastName:\"Smith\"}, Person{FirstName:\"John\", LastName:\"Doe\"}): expected %d, got %d", want, got)
	}
}

func TestComparer_Reverse(t *testing.T) {
	c := Natural[int]().Reverse()
	if c(1, 2) != 1 {
		t.Errorf("Natural().Reverse(1, 2): expected 1, got %d", c(1, 2))
	}
	if c(2, 1) != -1 {
		t.Errorf("Natural().Reverse(2, 1): expected -1, got %d", c(2, 1))
	}
	if c(1, 1) != 0 {
		t.Errorf("Natural().Reverse(1, 1): expected 0, got %d", c(1, 1))
	}
}

func TestComparer_Then(t *testing.T) {
	c := Comparing(func(p Person) string { return p.LastName }).
		Then(Comparing(func(p Person) string { return p.FirstName }))

	got := c(Person{"John", "Doe"}, Person{"Jane", "Doe"})
	want := 1
	if got != want {
		t.Errorf("Comparing(func(p Person) string { return p.LastName }).Then(Comparing(func(p Person) string { return p.FirstName }))(Person{FirstName:\"John\", LastName:\"Doe\"}, Person{FirstName:\"Jane\", LastName:\"Doe\"}): expected %d, got %d", want, got)
	}

	got = c(Person{"John", "Doe"}, Person{"John", "Smith"})
	want = -1
	if got != want {
		t.Errorf("Comparing(func(p Person) string { return p.LastName }).Then(Comparing(func(p Person) string { return p.FirstName }))(Person{FirstName:\"John\", LastName:\"Doe\"}, Person{FirstName:\"John\", LastName:\"Smith\"}): expected %d, got %d", want, got)
	}

	got = c(Person{"John", "Smith"}, Person{"John", "Doe"})
	want = 1
	if got != want {
		t.Errorf("Comparing(func(p Person) string { return p.LastName }).Then(Comparing(func(p Person) string { return p.FirstName }))(Person{FirstName:\"John\", LastName:\"Smith\"}, Person{FirstName:\"John\", LastName:\"Doe\"}): expected %d, got %d", want, got)
	}

	got = c(Person{"John", "Doe"}, Person{"John", "Doe"})
	want = 0
	if got != want {
		t.Errorf("Comparing(func(p Person) string { return p.LastName }).Then(Comparing(func(p Person) string { return p.FirstName }))(Person{FirstName:\"John\", LastName:\"Doe\"}, Person{FirstName:\"John\", LastName:\"Doe\"}): expected %d, got %d", want, got)
	}
}
