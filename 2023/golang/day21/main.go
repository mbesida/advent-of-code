package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

func main() {
	f := common.InputFileHandle("day21")
	defer f.Close()

	bytes, _ := io.ReadAll(f)
	start, garden := parseGarden(string(bytes))
	res := countReachability(garden, map[common.Point]struct{}{start: {}}, 0)
	fmt.Println(res)
}

func parseGarden(data string) (common.Point, [][]int) {
	splitted := strings.Split(data, "\n")
	var start common.Point
	matrix := make([][]int, len(splitted))
	for i, line := range splitted {
		matrix[i] = make([]int, len(line))
		for j, r := range line {
			if r == '.' {
				matrix[i][j] = 0
			} else if r == '#' {
				matrix[i][j] = 1
			} else if r == 'S' {
				start = common.Point{i, j}
				matrix[i][j] = 0
			}
		}
	}
	return start, matrix
}

func countReachability(garden [][]int, starts map[common.Point]struct{}, steps int) int {
	if steps == 64 {
		return len(starts)
	}
	nextStarts := make(map[common.Point]struct{})
	for p := range starts {
		if p.J+1 < len(garden[p.I]) && garden[p.I][p.J+1] == 0 {
			nextStarts[common.Point{p.I, p.J + 1}] = struct{}{}
		}
		if p.J-1 >= 0 && garden[p.I][p.J-1] == 0 {
			nextStarts[common.Point{p.I, p.J - 1}] = struct{}{}
		}
		if p.I+1 < len(garden) && garden[p.I+1][p.J] == 0 {
			nextStarts[common.Point{p.I + 1, p.J}] = struct{}{}
		}
		if p.I-1 >= 0 && garden[p.I-1][p.J] == 0 {
			nextStarts[common.Point{p.I - 1, p.J}] = struct{}{}
		}
	}
	return countReachability(garden, nextStarts, steps+1)
}
