package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"

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

type Node struct {
	p         common.Point
	distnceTo int
}
type PreviousValue struct {
	point     common.Point
	direction Direction
}

var input [][]int

var graph = make(map[common.Point][]Node)

func main() {
	f := common.InputFileHandle("day17")
	defer f.Close()

	data, _ := io.ReadAll(f)
	parseInput(data)
	n, m := buildGraph()
	res := dijkstra(common.Point{0, 0}, common.Point{n - 1, m - 1})

	fmt.Println(res)

}

func buildGraph() (int, int) {
	n := len(input)
	m := len(input[0])
	for i, row := range input {
		for j := range row {
			var neighbors []Node
			for _, points := range [][2]int{{i - 1, j}, {i + 1, j}, {i, j - 1}, {i, j + 1}} {
				a := points[0]
				b := points[1]
				if 0 <= a && a < n && 0 <= b && b < m {
					neighbors = append(neighbors, Node{common.Point{a, b}, input[a][b]})
				}
			}
			graph[common.Point{i, j}] = neighbors
		}
	}
	return n, m
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

func dijkstra(start, end common.Point) int {
	distances := make(map[common.Point]int)
	previous := make(map[common.Point]*PreviousValue)
	for point := range graph {
		distances[point] = math.MaxInt
	}
	distances[start] = 0

	pq := kpq.NewKeyedPriorityQueue[common.Point, int](func(a, b int) bool {
		return a < b
	})
	pq.Push(start, 0)

	for !pq.IsEmpty() {
		currentPoint, currentDistance, _ := pq.Pop()
		if currentDistance > distances[currentPoint] {
			continue
		}
		for _, neighbor := range graph[currentPoint] {
			distanceToNeighbor := currentDistance + neighbor.distnceTo
			if distanceToNeighbor < distances[neighbor.p] {
				if hasMoreThan3consecutive(currentPoint, neighbor.p, previous) {
					fmt.Println(neighbor.p)
					continue
				}
				distances[neighbor.p] = distanceToNeighbor
				dir := direction(currentPoint, neighbor.p)
				previous[neighbor.p] = &PreviousValue{currentPoint, dir}
				pq.Push(neighbor.p, distanceToNeighbor)
			}
		}
	}

	var path []common.Point
	current := end
	for previous[current] != nil {
		path = append([]common.Point{current}, path...)
		previousValue := *previous[current]
		current = previousValue.point
	}

	sum := 0
	for _, p := range path {
		sum += input[p.I][p.J]
		// fmt.Println(p)
	}
	for i, row := range input {
		for j := range row {
			if slices.Contains(path, common.Point{i, j}) {
				input[i][j] = 0
			}
		}
	}
	for _, row := range input {
		fmt.Println(row)

	}

	return sum

}

func hasMoreThan3consecutive(currentPoint, nextPoint common.Point, previous map[common.Point]*PreviousValue) bool {
	nextDirection := direction(currentPoint, nextPoint)
	p1 := previous[currentPoint]
	if p1 != nil {
		p2 := previous[p1.point]
		if p2 != nil {
			currentDirection := direction(p1.point, currentPoint)
			compactedDirections := slices.Compact([]Direction{nextDirection, currentDirection, p1.direction, p2.direction})
			return len(compactedDirections) == 1
		}

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
