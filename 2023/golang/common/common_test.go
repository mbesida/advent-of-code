package common

import (
	"slices"
	"testing"
)

func TestTranspose(t *testing.T) {
	/*
	  0,1,2    0,1,2,3
	  1,2,3 -- 1,2,3,4
	  2,3,4 -- 2,3,4,5
	  3,4,5
	*/
	data := [][]int{
		{0, 1, 2},
		{1, 2, 3},
		{2, 3, 4},
		{3, 4, 5},
	}

	trasnposed := TransposeMatrix(data)
	if len(trasnposed) != 3 {
		t.Errorf("length of transposed matrix is %d", len(trasnposed))
	}
	if slices.Compare(trasnposed[0], []int{0, 1, 2, 3}) != 0 {
		t.Error("not valid transposition")
	}
	if slices.Compare(trasnposed[1], []int{1, 2, 3, 4}) != 0 {
		t.Error("not valid transposition")
	}
	if slices.Compare(trasnposed[2], []int{2, 3, 4, 5}) != 0 {
		t.Error("not valid transposition")
	}
}
