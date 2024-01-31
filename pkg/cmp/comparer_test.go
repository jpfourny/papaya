package cmp

import (
	"strings"
	"testing"
	"time"

	"github.com/jpfourny/papaya/pkg/pair"
	"github.com/jpfourny/papaya/pkg/ptr"
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

func TestSelf(t *testing.T) {
	c := Self[Person]()
	got := c(Person{"John", "Doe"}, Person{"John", "Doe"})
	want := 0
	if got != want {
		t.Errorf("Self[Person]()(Person{FirstName:\"John\", LastName:\"Doe\"}, Person{FirstName:\"Jane\", LastName:\"Doe\"}): expected %d, got %d", want, got)
	}

	got = c(Person{"John", "Doe"}, Person{"John", "Smith"})
	want = -1
	if got != want {
		t.Errorf("Self[Person]()(Person{FirstName:\"John\", LastName:\"Doe\"}, Person{FirstName:\"John\", LastName:\"Smith\"}): expected %d, got %d", want, got)
	}

	got = c(Person{"John", "Smith"}, Person{"John", "Doe"})
	want = 1
	if got != want {
		t.Errorf("Self[Person]()(Person{FirstName:\"John\", LastName:\"Smith\"}, Person{FirstName:\"John\", LastName:\"Doe\"}): expected %d, got %d", want, got)
	}
}

func TestBool(t *testing.T) {
	c := Bool()
	if c(true, false) != 1 {
		t.Errorf("Bool()(true, false): expected 1, got %d", c(true, false))
	}
	if c(false, true) != -1 {
		t.Errorf("Bool()(false, true): expected -1, got %d", c(false, true))
	}
	if c(true, true) != 0 {
		t.Errorf("Bool()(true, true): expected 0, got %d", c(true, true))
	}
}

func TestComplex64(t *testing.T) {
	c := Complex64()
	if c(1+2i, 3+4i) != -1 {
		t.Errorf("Complex64()(1+2i, 3+4i): expected -1, got %d", c(1+2i, 3+4i))
	}
	if c(1+1i, 1+2i) != -1 {
		t.Errorf("Complex64()(1+1i, 1+2i): expected -1, got %d", c(1+1i, 1+2i))
	}
	if c(3+4i, 1+2i) != 1 {
		t.Errorf("Complex64()(3+4i, 1+2i): expected 1, got %d", c(3+4i, 1+2i))
	}
	if c(1+2i, 1+1i) != 1 {
		t.Errorf("Complex64()(1+2i, 1+1i): expected 1, got %d", c(1+2i, 1+1i))
	}
	if c(1+2i, 1+2i) != 0 {
		t.Errorf("Complex64()(1+2i, 1+2i): expected 0, got %d", c(1+2i, 1+2i))
	}
}

func TestComplex128(t *testing.T) {
	c := Complex128()
	if c(1+2i, 3+4i) != -1 {
		t.Errorf("Complex128()(1+2i, 3+4i): expected -1, got %d", c(1+2i, 3+4i))
	}
	if c(1+1i, 1+2i) != -1 {
		t.Errorf("Complex128()(1+1i, 1+2i): expected -1, got %d", c(1+1i, 1+2i))
	}
	if c(3+4i, 1+2i) != 1 {
		t.Errorf("Complex128()(3+4i, 1+2i): expected 1, got %d", c(3+4i, 1+2i))
	}
	if c(1+2i, 1+1i) != 1 {
		t.Errorf("Complex128()(1+2i, 1+1i): expected 1, got %d", c(1+2i, 1+1i))
	}
	if c(1+2i, 1+2i) != 0 {
		t.Errorf("Complex128()(1+2i, 1+2i): expected 0, got %d", c(1+2i, 1+2i))
	}
}

func TestTime(t *testing.T) {
	c := Time()
	if c(time.Unix(0, 0), time.Unix(0, 1)) != -1 {
		t.Errorf("Time()(time.Unix(0, 0), time.Unix(0, 1)): expected -1, got %d", c(time.Unix(0, 0), time.Unix(0, 1)))
	}
	if c(time.Unix(0, 1), time.Unix(0, 0)) != 1 {
		t.Errorf("Time()(time.Unix(0, 1), time.Unix(0, 0)): expected 1, got %d", c(time.Unix(0, 1), time.Unix(0, 0)))
	}
	if c(time.Unix(0, 0), time.Unix(0, 0)) != 0 {
		t.Errorf("Time()(time.Unix(0, 0), time.Unix(0, 0)): expected 0, got %d", c(time.Unix(0, 0), time.Unix(0, 0)))
	}
}

func TestPair(t *testing.T) {
	c := Pair(Natural[int](), Natural[int]())

	got := c(pair.Of(1, 2), pair.Of(1, 2))
	want := 0
	if got != want {
		t.Errorf("Pair(Natural[int](), Natural[int]())(pair.Of(1, 2), pair.Of(1, 2)): expected %d, got %d", want, got)
	}

	got = c(pair.Of(1, 2), pair.Of(2, 2))
	want = -1
	if got != want {
		t.Errorf("Pair(Natural[int](), Natural[int]())(pair.Of(1, 2), pair.Of(2, 2)): expected %d, got %d", want, got)
	}

	got = c(pair.Of(1, 2), pair.Of(1, 3))
	want = -1
	if got != want {
		t.Errorf("Pair(Natural[int](), Natural[int]())(pair.Of(1, 2), pair.Of(1, 3)): expected %d, got %d", want, got)
	}

	got = c(pair.Of(2, 2), pair.Of(1, 2))
	want = 1
	if got != want {
		t.Errorf("Pair(Natural[int](), Natural[int]())(pair.Of(2, 2), pair.Of(1, 2)): expected %d, got %d", want, got)
	}

	got = c(pair.Of(2, 2), pair.Of(2, 1))
	want = 1
	if got != want {
		t.Errorf("Pair(Natural[int](), Natural[int]())(pair.Of(2, 2), pair.Of(2, 1)): expected %d, got %d", want, got)
	}
}

func TestSlice(t *testing.T) {
	c := Slice(Natural[int]())
	got := c([]int{3, 1, 2}, []int{1, 2, 3})
	want := 1
	if got != want {
		t.Errorf("Slice(Natural[int]())([]int{3, 1, 2}, []int{1, 2, 3}): expected %d, got %d", want, got)
	}

	got = c([]int{1, 2, 3}, []int{3, 1, 2})
	want = -1
	if got != want {
		t.Errorf("Slice(Natural[int]())([]int{1, 2, 3}, []int{3, 1, 2}): expected %d, got %d", want, got)
	}

	got = c([]int{1, 2, 3}, []int{1, 2, 3})
	want = 0
	if got != want {
		t.Errorf("Slice(Natural[int]())([]int{1, 2, 3}, []int{1, 2, 3}): expected %d, got %d", want, got)
	}

	got = c([]int{1, 2, 3}, []int{1, 2, 3, 4})
	want = -1
	if got != want {
		t.Errorf("Slice(Natural[int]())([]int{1, 2, 3}, []int{1, 2, 3, 4}): expected %d, got %d", want, got)
	}

	got = c([]int{1, 2, 3, 4}, []int{1, 2, 3})
	want = 1
	if got != want {
		t.Errorf("Slice(Natural[int]())([]int{1, 2, 3, 4}, []int{1, 2, 3}): expected %d, got %d", want, got)
	}
}

func TestDerefNilFirst(t *testing.T) {
	c := DerefNilFirst(Natural[int]())
	got := c(nil, nil)
	want := 0
	if got != want {
		t.Errorf("DerefNilFirst()(nil, nil): expected %d, got %d", want, got)
	}

	got = c(nil, ptr.Ref(1))
	want = -1
	if got != want {
		t.Errorf("DerefNilFirst()(nil, 1): expected %d, got %d", want, got)
	}

	got = c(ptr.Ref(1), nil)
	want = 1
	if got != want {
		t.Errorf("DerefNilFirst()(1, nil): expected %d, got %d", want, got)
	}

	got = c(ptr.Ref(1), ptr.Ref(2))
	want = -1
	if got != want {
		t.Errorf("DerefNilFirst()(1, 2): expected %d, got %d", want, got)
	}

	got = c(ptr.Ref(2), ptr.Ref(1))
	want = 1
	if got != want {
		t.Errorf("DerefNilFirst()(2, 1): expected %d, got %d", want, got)
	}

	got = c(ptr.Ref(1), ptr.Ref(1))
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

	got = c(nil, ptr.Ref(1))
	want = 1
	if got != want {
		t.Errorf("DerefNilLast()(nil, 1): expected %d, got %d", want, got)
	}

	got = c(ptr.Ref(1), nil)
	want = -1
	if got != want {
		t.Errorf("DerefNilLast()(1, nil): expected %d, got %d", want, got)
	}

	got = c(ptr.Ref(1), ptr.Ref(2))
	want = -1
	if got != want {
		t.Errorf("DerefNilLast()(1, 2): expected %d, got %d", want, got)
	}

	got = c(ptr.Ref(2), ptr.Ref(1))
	want = 1
	if got != want {
		t.Errorf("DerefNilLast()(2, 1): expected %d, got %d", want, got)
	}

	got = c(ptr.Ref(1), ptr.Ref(1))
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

func TestComparer_Equal(t *testing.T) {
	c := Natural[int]()
	if !c.Equal(1, 1) {
		t.Errorf("Natural().Equal(1, 1): expected true, got false")
	}
	if c.Equal(1, 2) {
		t.Errorf("Natural().Equal(1, 2): expected false, got true")
	}
}

func TestComparer_NotEqual(t *testing.T) {
	c := Natural[int]()
	if c.NotEqual(1, 1) {
		t.Errorf("Natural().NotEqual(1, 1): expected false, got true")
	}
	if !c.NotEqual(1, 2) {
		t.Errorf("Natural().NotEqual(1, 2): expected true, got false")
	}
}

func TestComparer_LessThan(t *testing.T) {
	c := Natural[int]()
	if !c.LessThan(1, 2) {
		t.Errorf("Natural().LessThan(1, 2): expected true, got false")
	}
	if c.LessThan(2, 1) {
		t.Errorf("Natural().LessThan(2, 1): expected false, got true")
	}
}

func TestComparer_LessThanOrEqual(t *testing.T) {
	c := Natural[int]()
	if !c.LessThanOrEqual(1, 2) {
		t.Errorf("Natural().LessThanOrEqual(1, 2): expected true, got false")
	}
	if c.LessThanOrEqual(2, 1) {
		t.Errorf("Natural().LessThanOrEqual(2, 1): expected false, got true")
	}
	if !c.LessThanOrEqual(1, 1) {
		t.Errorf("Natural().LessThanOrEqual(1, 1): expected true, got false")
	}
}

func TestComparer_GreaterThan(t *testing.T) {
	c := Natural[int]()
	if c.GreaterThan(1, 2) {
		t.Errorf("Natural().GreaterThan(1, 2): expected false, got true")
	}
	if !c.GreaterThan(2, 1) {
		t.Errorf("Natural().GreaterThan(2, 1): expected true, got false")
	}
}

func TestComparer_GreaterThanOrEqual(t *testing.T) {
	c := Natural[int]()
	if c.GreaterThanOrEqual(1, 2) {
		t.Errorf("Natural().GreaterThanOrEqual(1, 2): expected false, got true")
	}
	if !c.GreaterThanOrEqual(2, 1) {
		t.Errorf("Natural().GreaterThanOrEqual(2, 1): expected true, got false")
	}
	if !c.GreaterThanOrEqual(1, 1) {
		t.Errorf("Natural().GreaterThanOrEqual(1, 1): expected true, got false")
	}
}

func TestComparer_Min(t *testing.T) {
	c := Natural[int]()
	if got := c.Min(1, 2); got != 1 {
		t.Errorf("Natural().Min(1, 2): expected 1, got %d", got)
	}
	if got := c.Min(2, 1); got != 1 {
		t.Errorf("Natural().Min(2, 1): expected 1, got %d", got)
	}
}

func TestComparer_Max(t *testing.T) {
	c := Natural[int]()
	if got := c.Max(1, 2); got != 2 {
		t.Errorf("Natural().Max(1, 2): expected 2, got %d", got)
	}
	if got := c.Max(2, 1); got != 2 {
		t.Errorf("Natural().Max(2, 1): expected 2, got %d", got)
	}
}
