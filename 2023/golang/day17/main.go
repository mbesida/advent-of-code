package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
	"github.com/rdleal/go-priorityq/kpq"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Key struct {
	point common.Point
	path  string
}

func NewKey(p common.Point) Key {
	return Key{p, fmt.Sprintf("p%d,%d;", p.I, p.J)}
}

func (k Key) NextKey(nextPoint common.Point) Key {
	return Key{nextPoint, k.path + fmt.Sprintf("p%d,%d;", nextPoint.I, nextPoint.J)}
}

func (k Key) buildPath() []common.Point {
	var res []common.Point
	points := strings.Split(k.path, ";")
	for _, p := range points {
		if p != "" {
			ij := strings.Split(strings.TrimLeft(p, "p"), ",")
			i, _ := strconv.Atoi(ij[0])
			j, _ := strconv.Atoi(ij[1])
			res = append(res, common.Point{i, j})
		}
	}
	return res
}

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
	pq := kpq.NewKeyedPriorityQueue[Key, int](func(a, b int) bool {
		return a < b
	})

	distances := make([][]int, n)
	visited := make([][]bool, n)
	for i := 0; i < n; i++ {
		visited[i] = make([]bool, m)
		distances[i] = make([]int, m)
		for j := 0; j < m; j++ {
			distances[i][j] = math.MaxInt
		}
	}

	pq.Push(NewKey(common.Point{startI, startJ}), 0)
	distances[startI][startJ] = 0

	// var path []common.Point
	for !pq.IsEmpty() {
		key, dist, _ := pq.Pop()
		current := key.point
		visited[current.I][current.J] = true
		if visited[endI][endJ] {
			// path = key.buildPath()
			break
		}
		for _, delta := range deltas {
			ni, nj := current.I+delta[0], current.J+delta[1]
			next := common.Point{I: ni, J: nj}

			if inBounds(ni, nj, n, m) && !visited[ni][nj] {
				newDist := dist + input[ni][nj]
				if !hasMoreThan3consecutive(next, key) {
					pq.Push(key.NextKey(next), newDist)
				}
				if newDist <= distances[ni][nj] {
					distances[ni][nj] = newDist
				}
			}
		}
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

	return distances[endI][endJ]
}

func hasMoreThan3consecutive(next common.Point, key Key) bool {
	points := key.buildPath()
	l := len(points)
	if l > 3 {
		current := points[l-1]
		p1 := points[l-2]
		p2 := points[l-3]
		p3 := points[l-4]
		pn := direction(current, next)
		p1p := direction(p1, current)
		p2p1 := direction(p2, p1)
		p3p2 := direction(p3, p2)
		compactedDirections := slices.Compact([]Direction{pn, p1p, p2p1, p3p2})
		return len(compactedDirections) == 1
	}

	return false
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
