package main

import "testing"

func Test(t *testing.T) {
	//in{x>3:R,x<2:R,A}
	// slp{x<583:A,A}
	// in{x<4000:R,m<4000:R,a<4000:R,s<4000:R,A}
	input := make(map[string]Workflow)
	input["in"] = Workflow{"in", nil,
		[]Condition{
			{&Rule{"x", "<", 4000}, "R"},
			{&Rule{"m", "<", 4000}, "R"},
			{&Rule{"a", "<", 4000}, "R"},
			{&Rule{"s", "<", 4000}, "R"},
			{nil, "A"}}}
	res := part2(input)
	if res != 1 {
		t.FailNow()
	}
}
