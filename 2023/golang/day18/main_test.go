package main

import (
	"fmt"
	"testing"

	"github.com/mbesida/advent-of-code-2023/common"
)

func TestIsInPolygon(t *testing.T) {
	data := [][]int{
		{1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0},
		{1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 0, 0, 0},
		{1, 1, 1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0},
		{0, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0},
	}

	terrain = data
	lines = []Line{
		NewLine(common.Point{0, 0}, common.Point{0, 3}),
		NewLine(common.Point{0, 3}, common.Point{3, 3}),
		NewLine(common.Point{3, 3}, common.Point{3, 9}),
		NewLine(common.Point{3, 9}, common.Point{5, 9}),
		NewLine(common.Point{5, 9}, common.Point{5, 10}),
		NewLine(common.Point{5, 10}, common.Point{8, 10}),
		NewLine(common.Point{8, 6}, common.Point{8, 10}),
		NewLine(common.Point{6, 6}, common.Point{8, 6}),
		NewLine(common.Point{6, 4}, common.Point{6, 6}),
		NewLine(common.Point{6, 4}, common.Point{8, 4}),
		NewLine(common.Point{8, 2}, common.Point{8, 4}),
		NewLine(common.Point{7, 2}, common.Point{8, 2}),
		NewLine(common.Point{7, 0}, common.Point{7, 2}),
		NewLine(common.Point{0, 0}, common.Point{7, 2}),
	}
	fillTerrain(9, 14)
	for _, v := range terrain {
		fmt.Println(v)
	}
}
