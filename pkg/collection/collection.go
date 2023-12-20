package collection

import (
	"github.com/jpfourny/papaya/pkg/optional"
	"github.com/jpfourny/papaya/pkg/pair"
	"github.com/jpfourny/papaya/pkg/stream"
)

type Collection[E any] interface {
	Size() int
	IsEmpty() bool
	Contains(elem E) bool
	ContainsAll(elems Collection[E]) bool
	ForEach(func(elem E) bool) bool
	Stream() stream.Stream[E]
	ToSlice() []E
	Clone() Collection[E]
	String() string
}

type MutableCollection[E any] interface {
	Collection[E]

	Add(elem E) bool
	AddAll(elems Collection[E]) bool
	Remove(elem E)
	RemoveAll(elems Collection[E])
	RetainAll(elems Collection[E])
	Clear()
}

type List[E any] interface {
	Collection[E]

	At(index int) optional.Optional[E]
	IndexOf(elem E) int
	LastIndexOf(elem E) int
	SubList(fromIndex int, toIndex int) List[E]
	ForEachIndex(func(index int, elem E) bool) bool
}

type MutableList[E any] interface {
	List[E]
	MutableCollection[E]

	Set(index int, elem E) optional.Optional[E]
	Insert(index int, elem E)
	InsertAll(index int, elems Collection[E])
	RemoveAt(index int) optional.Optional[E]
	RemoveRange(fromIndex int, toIndex int)
}

type Set[E any] interface {
	Collection[E]

	Union(other Set[E]) Set[E]
	Intersection(other Set[E]) Set[E]
	Difference(other Set[E]) Set[E]
	SymmetricDifference(other Set[E]) Set[E]
	IsSubsetOf(other Set[E]) bool
	IsSupersetOf(other Set[E]) bool
	IsProperSubsetOf(other Set[E]) bool
	IsProperSupersetOf(other Set[E]) bool
}

type MutableSet[E any] interface {
	Set[E]
	MutableCollection[E]
}

type Map[K any, V any] interface {
	Collection[pair.Pair[K, V]]

	ContainsKeyValue(key K, value V) bool
	ContainsKey(key K) bool
	ContainerValue(value V) bool
	Get(key K) optional.Optional[V]

	Keys() Set[K]
	Values() Collection[V]
}

type MutableMap[K any, V any] interface {
	Map[K, V]
	MutableCollection[pair.Pair[K, V]]

	Put(key K, value V)
	PutAll(m Map[K, V])
	RemoveKey(key K)
	RemoveKeyValue(key K, value V)
}
