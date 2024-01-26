package main

import (
	"bytes"
	"fmt"
	"io"
	"slices"

	"github.com/mbesida/advent-of-code-2023/common"
	"golang.org/x/exp/maps"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

var grid [][]rune
var start, end common.Point

type Edge struct {
	to     common.Point
	weight int
}
type Graph struct {
	adjList map[common.Point][]Edge
}

func (g *Graph) AddEdge(from common.Point, edge Edge) {
	if g.adjList == nil {
		g.adjList = make(map[common.Point][]Edge)
	}

	g.adjList[from] = append(g.adjList[from], edge)

	common.ExecuteActions(func() {}, func() {
		// in task 2 we need to take into account reverse paths
		// because '>' '<' '^' 'v' doesn't exist in the grid
		g.adjList[edge.to] = append(g.adjList[edge.to], Edge{from, edge.weight})
	})
}
func (g *Graph) ContainsEdge(from common.Point, e Edge) bool {
	existing1, ok1 := g.adjList[from]
	if ok1 {
		if slices.Contains(existing1, e) {
			return true
		}
	}
	existing2, ok2 := g.adjList[e.to]
	if ok2 {
		if slices.Contains(existing2, Edge{from, e.weight}) {
			return true
		}
	}
	return false
}
func (g *Graph) VerticesToIndex() map[common.Point]int {
	vertixToIndex := make(map[common.Point]int)
	vertices := maps.Keys(g.adjList)
	for i, v := range vertices {
		vertixToIndex[v] = i
	}
	if _, ok := vertixToIndex[end]; !ok {
		vertixToIndex[end] = len(vertixToIndex)
	}
	return vertixToIndex
}
func (g *Graph) Grid(lookupTable map[common.Point]int) [][]int {
	n := len(lookupTable)
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, n)
	}
	for p, neigbours := range g.adjList {
		for _, n := range neigbours {
			from := lookupTable[p]
			to := lookupTable[n.to]
			grid[from][to] = n.weight
		}
	}
	return grid
}

func main() {
	f := common.InputFileHandle("day23")
	defer f.Close()
	parseData(f)
	graph := &Graph{}
	buildGraph(start, nil, graph)
	lookupTable := graph.VerticesToIndex()
	newGrid := graph.Grid(lookupTable)

	res := longestWalk(newGrid, lookupTable[start], lookupTable[end])
	fmt.Println(res)
}

func parseData(reader io.Reader) {
	data, _ := io.ReadAll(reader)
	splitted := bytes.Split(data, []byte("\n"))
	grid = make([][]rune, len(splitted))
	for i, line := range splitted {
		str := string(line)
		grid[i] = make([]rune, len(str))
		for j, r := range str {
			grid[i][j] = r
		}
	}
	firstLine := grid[0]
	lastLine := grid[len(grid)-1]
	for j, r := range firstLine {
		if r == '.' {
			start = common.Point{0, j}
			break
		}
	}
	for j, r := range lastLine {
		if r == '.' {
			end = common.Point{len(grid) - 1, j}
			break
		}
	}
}

func neigbours(p common.Point) []common.Point {
	var next []common.Point
	for _, np := range []common.Point{{p.I + 1, p.J}, {p.I - 1, p.J}, {p.I, p.J + 1}, {p.I, p.J - 1}} {
		if np.I >= 0 && np.I < len(grid) && np.J >= 0 && np.J < len(grid[0]) {
			if grid[np.I][np.J] != '#' {
				next = append(next, np)
			}
		}
	}
	return next
}

func validNeigbourPoints(p common.Point) []common.Point {
	allNeighbours := neigbours(p)
	var validNeigbours []common.Point
	for _, np := range allNeighbours {
		r := grid[np.I][np.J]
		dir := direction(p, np)
		if isValidNeigbour(r, dir) {
			validNeigbours = append(validNeigbours, np)
		}
	}
	return validNeigbours
}

func direction(current, next common.Point) Direction {
	var direction Direction
	switch {
	case current.I == next.I && current.J > next.J:
		direction = Left
	case current.I == next.I && current.J < next.J:
		direction = Right
	case current.J == next.J && current.I > next.I:
		direction = Up
	case current.J == next.J && current.I < next.I:
		direction = Down
	}
	return direction
}

func isValidNeigbour(r rune, dir Direction) bool {
	res := common.HandleValue(
		r == '.' || (r == '<' && dir != Right) || (r == '>' && dir != Left) ||
			(r == '^' && dir != Down) || (r == 'v' && dir != Up),
		r != '#',
	)
	return res
}

func longestWalk(matrix [][]int, start, end int) int {
	visited := make([]bool, len(matrix))

	var dfs func(int, int) int
	dfs = func(currentVertex int, length int) int {
		if currentVertex == end {
			return length
		}
		visited[currentVertex] = true
		var currentMax int
		for j, v := range matrix[currentVertex] {
			if v != 0 && !visited[j] {
				newLength := length + matrix[currentVertex][j]
				currentMax = max(currentMax, dfs(j, newLength))
			}
		}
		visited[currentVertex] = false
		return currentMax
	}

	return dfs(start, 0)
}

func buildGraph(p common.Point, next *common.Point, g *Graph) {
	var edge Edge
	var prev common.Point
	if next == nil {
		edge, prev = buildEdge(p, common.Point{p.I + 1, p.J})
	} else {
		edge, prev = buildEdge(p, *next)
	}
	if !g.ContainsEdge(p, edge) {
		g.AddEdge(p, edge)
		if edge.to == end {
			return
		}
		for _, next := range neigboursExceptPrevious(edge.to, prev) {
			buildGraph(edge.to, &next, g)
		}
	}
}

func buildEdge(current, next common.Point) (Edge, common.Point) {
	weight := 1
	neigbours := neigboursExceptPrevious(next, current)
	for len(neigbours) == 1 {
		weight++
		current = next
		next = neigbours[0]
		neigbours = neigboursExceptPrevious(next, current)
	}
	return Edge{next, weight}, current
}

func neigboursExceptPrevious(current common.Point, previous common.Point) []common.Point {
	neigbours := validNeigbourPoints(current)
	return slices.DeleteFunc(neigbours, func(p common.Point) bool { return p == previous })
}
