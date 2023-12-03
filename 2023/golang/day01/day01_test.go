package main

import (
	"testing"
)

func TestParsing(t *testing.T) {
	s := "foneight3fsbhdqzr5twojbsdnntwohd9seven"
	res := digits2(s)
	if res != [2]int{1, 7} {
		t.Fatalf("%v", res)
	}
}
