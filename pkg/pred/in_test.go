package pred

import "testing"

func TestIn(t *testing.T) {
	p := In(0, 1, 2)
	if !p(0) {
		t.Errorf("In(0, 1, 2)(0) = false; want true")
	}
	if !p(1) {
		t.Errorf("In(0, 1, 2)(1) = false; want true")
	}
	if !p(2) {
		t.Errorf("In(0, 1, 2)(2) = false; want true")
	}
	if p(3) {
		t.Errorf("In(0, 1, 2)(3) = true; want false")
	}
}

func TestNotIn(t *testing.T) {
	p := NotIn(0, 1, 2)
	if p(0) {
		t.Errorf("NotIn(0, 1, 2)(0) = true; want false")
	}
	if p(1) {
		t.Errorf("NotIn(0, 1, 2)(1) = true; want false")
	}
	if p(2) {
		t.Errorf("NotIn(0, 1, 2)(2) = true; want false")
	}
	if !p(3) {
		t.Errorf("NotIn(0, 1, 2)(3) = false; want true")
	}
}

func TestInSlice(t *testing.T) {
	p := InSlice[int]([]int{0, 1, 2})
	if !p(0) {
		t.Errorf("InSlice([]int{0, 1, 2})(0) = false; want true")
	}
	if !p(1) {
		t.Errorf("InSlice([]int{0, 1, 2})(1) = false; want true")
	}
	if !p(2) {
		t.Errorf("InSlice([]int{0, 1, 2})(2) = false; want true")
	}
	if p(3) {
		t.Errorf("InSlice([]int{0, 1, 2})(3) = true; want false")
	}
}

func TestNotInSlice(t *testing.T) {
	p := NotInSlice[int]([]int{0, 1, 2})
	if p(0) {
		t.Errorf("NotInSlice([]int{0, 1, 2})(0) = true; want false")
	}
	if p(1) {
		t.Errorf("NotInSlice([]int{0, 1, 2})(1) = true; want false")
	}
	if p(2) {
		t.Errorf("NotInSlice([]int{0, 1, 2})(2) = true; want false")
	}
	if !p(3) {
		t.Errorf("NotInSlice([]int{0, 1, 2})(3) = false; want true")
	}
}

func TestInBy(t *testing.T) {
	p := InBy[int](func(a, b int) int {
		return a - b
	}, 0, 1, 2)
	if !p(0) {
		t.Errorf("InBy(cmp, 0, 1, 2)(0) = false; want true")
	}
	if !p(1) {
		t.Errorf("InBy(cmp, 0, 1, 2)(1) = false; want true")
	}
	if !p(2) {
		t.Errorf("InBy(cmp, 0, 1, 2)(2) = false; want true")
	}
	if p(3) {
		t.Errorf("InBy(cmp, 0, 1, 2)(3) = true; want false")
	}
}

func TestNotInBy(t *testing.T) {
	p := NotInBy[int](func(a, b int) int {
		return a - b
	}, 0, 1, 2)
	if p(0) {
		t.Errorf("NotInBy(cmp, 0, 1, 2)(0) = true; want false")
	}
	if p(1) {
		t.Errorf("NotInBy(cmp, 0, 1, 2)(1) = true; want false")
	}
	if p(2) {
		t.Errorf("NotInBy(cmp, 0, 1, 2)(2) = true; want false")
	}
	if !p(3) {
		t.Errorf("NotInBy(cmp, 0, 1, 2)(3) = false; want true")
	}
}

func TestInSliceBy(t *testing.T) {
	p := InSliceBy[int]([]int{0, 1, 2}, func(a, b int) int {
		return a - b
	})
	if !p(0) {
		t.Errorf("InSliceBy([]int{0, 1, 2}, cmp)(0) = false; want true")
	}
	if !p(1) {
		t.Errorf("InSliceBy([]int{0, 1, 2}, cmp)(1) = false; want true")
	}
	if !p(2) {
		t.Errorf("InSliceBy([]int{0, 1, 2}, cmp)(2) = false; want true")
	}
	if p(3) {
		t.Errorf("InSliceBy([]int{0, 1, 2}, cmp)(3) = true; want false")
	}
}

func TestNotInSliceBy(t *testing.T) {
	p := NotInSliceBy[int]([]int{0, 1, 2}, func(a, b int) int {
		return a - b
	})
	if p(0) {
		t.Errorf("NotInSliceBy([]int{0, 1, 2}, cmp)(0) = true; want false")
	}
	if p(1) {
		t.Errorf("NotInSliceBy([]int{0, 1, 2}, cmp)(1) = true; want false")
	}
	if p(2) {
		t.Errorf("NotInSliceBy([]int{0, 1, 2}, cmp)(2) = true; want false")
	}
	if !p(3) {
		t.Errorf("NotInSliceBy([]int{0, 1, 2}, cmp)(3) = false; want true")
	}
}

func TestInSet(t *testing.T) {
	set := map[int]struct{}{0: {}, 1: {}, 2: {}}
	p := InSet[int](set)
	if !p(0) {
		t.Errorf("InSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(0) = false; want true")
	}
	if !p(1) {
		t.Errorf("InSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(1) = false; want true")
	}
	if !p(2) {
		t.Errorf("InSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(2) = false; want true")
	}
	if p(3) {
		t.Errorf("InSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(3) = true; want false")
	}
}

func TestNotInSet(t *testing.T) {
	set := map[int]struct{}{0: {}, 1: {}, 2: {}}
	p := NotInSet[int](set)
	if p(0) {
		t.Errorf("NotInSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(0) = true; want false")
	}
	if p(1) {
		t.Errorf("NotInSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(1) = true; want false")
	}
	if p(2) {
		t.Errorf("NotInSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(2) = true; want false")
	}
	if !p(3) {
		t.Errorf("NotInSet(map[int]struct{}{0: {}, 1: {}, 2: {}})(3) = false; want true")
	}
}
