package main

import "testing"

func TestArrangment(t *testing.T) {
	n1 := arrangements(RowData{[]int{1, 1, 3}, "???.###."})
	n2 := arrangements(RowData{[]int{1, 1, 3}, ".??..??...?##."})
	n3 := arrangements(RowData{[]int{1, 3, 1, 6}, "?#?#?#?#?#?#?#?"})
	n4 := arrangements(RowData{[]int{3, 2, 1}, "?###????????"})
	if n1 != 1 {
		t.Errorf("%d", n1)
	}
	if n2 != 4 {
		t.Errorf("%d", n2)
	}
	if n3 != 1 {
		t.Errorf("%d", n3)
	}
	if n4 != 10 {
		t.Errorf("%d", n4)
	}
}
