package main

import (
	"bytes"
	"fmt"
	"io"
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

type Point struct {
	x, y int
}

func (p Point) validNeigbourPoints(n, m int) []Point {
	var neigbours []Point
	for _, delta := range [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		nx := p.x + delta[0]
		ny := p.y + delta[1]
		if nx >= 0 && nx < n && ny >= 0 && ny < m {
			neigbours = append(neigbours, Point{nx, ny})
		}
	}
	return neigbours
}

type Node struct {
	point     Point
	direction Direction
	counter   int
}

func (node Node) neighbours(n, m int, min, max int) []Node {
	var neigbours []Node
	for _, p := range node.point.validNeigbourPoints(n, m) {
		newDirection := direction(node.point, p)
		if newDirection == oppositeDirection(node.direction) {
			continue
		} else if node.direction != newDirection && node.counter >= min {
			neigbours = append(neigbours, Node{p, newDirection, 1})
		} else if node.direction == newDirection && node.counter < max {
			neigbours = append(neigbours, Node{p, newDirection, node.counter + 1})
		}
	}
	return neigbours
}

var input [][]int

func main() {
	f := common.InputFileHandle("day17")
	defer f.Close()

	data, _ := io.ReadAll(f)
	parseInput(data)
	min := common.HandleValue(1, 4)
	max := common.HandleValue(3, 10)
	n, m := len(input), len(input[0])
	res := dijkstra(Point{0, 0}, Point{n - 1, m - 1}, min, max)
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

func dijkstra(start, end Point, min, max int) int {
	n, m := len(input), len(input[0])
	pq := kpq.NewKeyedPriorityQueue[Node, int](func(a, b int) bool {
		return a < b
	})

	distances := make(map[Node]int)
	distances[Node{start, Right, 0}] = 0
	distances[Node{start, Down, 0}] = 0
	pq.Push(Node{start, Right, 0}, 0)
	pq.Push(Node{start, Down, 0}, 0)

	var cost int
	for !pq.IsEmpty() {
		node, dist, _ := pq.Pop()

		if node.point == end {
			cost = dist
			break
		}

		for _, neighbour := range node.neighbours(n, m, min, max) {
			newDist := dist + input[neighbour.point.x][neighbour.point.y]
			gridDistance, ok := distances[neighbour]
			if !ok || (ok && newDist < gridDistance) {
				pq.Push(neighbour, newDist)
				distances[neighbour] = newDist
			}
		}
	}

	return cost
}

func direction(current, next Point) Direction {
	var res Direction
	switch {
	case current.x == next.x && current.y+1 == next.y:
		res = Right
	case current.x == next.x && current.y-1 == next.y:
		res = Left
	case current.x+1 == next.x && current.y == next.y:
		res = Down
	case current.x-1 == next.x && current.y == next.y:
		res = Up
	}

	return res
}

func oppositeDirection(d Direction) Direction {
	var opposite Direction
	switch d {
	case Up:
		opposite = Down
	case Down:
		opposite = Up
	case Right:
		opposite = Left
	case Left:
		opposite = Right
	}
	return opposite
}

func tracePath(path []Point, n, m int) {
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if slices.Contains(path, Point{i, j}) {
				fmt.Print(" " + strconv.Itoa(input[i][j]) + "*")
			} else {
				fmt.Print(" " + strconv.Itoa(input[i][j]) + " ")
			}
		}
		fmt.Println()
	}
}
