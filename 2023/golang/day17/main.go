package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/mbesida/advent-of-code-2023/common"
	"github.com/rdleal/go-priorityq/kpq"
)

type Direction int

const (
	Undefined Direction = iota
	Up
	Down
	Left
	Right
)

var input [][]int

func main() {
	f := common.InputFileHandle("day17")
	defer f.Close()

	data, _ := io.ReadAll(f)
	parseInput(data)
	n, m := len(input), len(input[0])
	res := dijkstra(0, 0, n-1, m-1)
	fmt.Println(res)
}

func parseInput(data []byte) {
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		row := make([]int, len(line))
		for j, d := range line {
			value, _ := strconv.Atoi(string(d))
			row[j] = value
		}
		input = append(input, row)
	}
}

func dijkstra(startI, startJ, endI, endJ int) int {
	n, m := len(input), len(input[0])
	deltas := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	pq := kpq.NewKeyedPriorityQueue[common.Point, int](func(a, b int) bool {
		return a < b
	})

	distances := make(map[common.Point]int)
	prev := make(map[common.Point]common.Point)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			point := common.Point{i, j}
			distances[point] = math.MaxInt
			pq.Push(point, distances[point])
		}
	}
	start := common.Point{startI, startJ}
	distances[start] = 0
	pq.Update(start, distances[start])

	// var path []common.Point
	for !pq.IsEmpty() {
		point, dist, _ := pq.Pop()

		for _, delta := range deltas {
			ni, nj := point.I+delta[0], point.J+delta[1]
			next := common.Point{ni, nj}

			if inBounds(ni, nj, n, m) {
				newDist := dist + input[ni][nj]
				if newDist < distances[next] {
					distances[next] = newDist
					pq.Update(next, newDist)
					prev[next] = point
				}
			}
		}
	}
	var path []common.Point
	current := common.Point{endI, endJ}
	for {
		p, ok := prev[current]
		if !ok {
			break
		}
		path = append([]common.Point{current}, path...)
		current = p
	}
	for _, p := range path {
		fmt.Println(p)
	}
	// sum := 0
	// for _, p := range path {
	// 	sum += input[p.I][p.J]
	// }
	// fmt.Println("sum is", sum-input[startI][startJ])

	// runes := make([][]string, n)
	// for i := 0; i < n; i++ {
	// 	runes[i] = make([]string, m)
	// 	for j := 0; j < m; j++ {
	// 		if slices.Contains(path, common.Point{i, j}) {
	// 			runes[i][j] = " " + strconv.Itoa(input[i][j]) + "*"
	// 		} else {
	// 			runes[i][j] = " " + strconv.Itoa(input[i][j]) + " "
	// 		}
	// 	}
	// }
	// for _, p := range runes {
	// 	fmt.Println(p)
	// }

	return distances[common.Point{endI, endJ}]
}

func direction(current, next common.Point) Direction {
	var res Direction
	switch {
	case current.I == next.I && current.J+1 == next.J:
		res = Right
	case current.I == next.I && current.J-1 == next.J:
		res = Left
	case current.I+1 == next.I && current.J == next.J:
		res = Down
	case current.I-1 == next.I && current.J == next.J:
		res = Up
	}

	return res
}

func inBounds(i, j int, n, m int) bool {
	return 0 <= i && i < n && 0 <= j && j < m
}
