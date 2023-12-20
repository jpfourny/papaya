package collection

import (
	"github.com/jpfourny/papaya/pkg/cmp"
	"github.com/jpfourny/papaya/pkg/optional"
	"github.com/jpfourny/papaya/pkg/stream"
)

type SliceList[E any] struct {
	compare cmp.Comparer[E]
	elems   []E
}

var _ MutableList[int] = (*SliceList[int])(nil)

func (s *SliceList[E]) Size() int {
	return len(s.elems)
}

func (s *SliceList[E]) IsEmpty() bool {
	return len(s.elems) == 0
}

func (s *SliceList[E]) Contains(elem E) bool {
	return stream.ContainsBy(stream.FromSlice(s.elems), s.compare, elem)
}

func (s SliceList[E]) ContainsAll(elems Collection[E]) bool {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) ForEach(f func(elem E) bool) bool {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) Stream() stream.Stream[E] {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) ToSlice() []E {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) Clone() Collection[E] {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) String() string {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) At(index int) optional.Optional[E] {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) IndexOf(elem E) int {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) LastIndexOf(elem E) int {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) SubList(fromIndex int, toIndex int) List[E] {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) ForEachIndex(f func(index int, elem E) bool) bool {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) Add(elem E) bool {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) AddAll(elems Collection[E]) bool {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) Remove(elem E) {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) RemoveAll(elems Collection[E]) {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) RetainAll(elems Collection[E]) {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) Clear() {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) Set(index int, elem E) optional.Optional[E] {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) Insert(index int, elem E) {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) InsertAll(index int, elems Collection[E]) {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) RemoveAt(index int) optional.Optional[E] {
	//TODO implement me
	panic("implement me")
}

func (s SliceList[E]) RemoveRange(fromIndex int, toIndex int) {
	//TODO implement me
	panic("implement me")
}
