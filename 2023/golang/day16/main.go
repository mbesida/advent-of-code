package main

import (
	"bytes"
	"fmt"
	"io"
	"slices"

	"github.com/mbesida/advent-of-code-2023/common"
)

type Direction int

const (
	Left Direction = iota
	Up
	Right
	Down
)

type PathItem struct {
	i, j          int
	energizedWhen Direction
}

var layout [][]rune
var energized [][]byte
var path map[PathItem]struct{} = make(map[PathItem]struct{})

func main() {
	f := common.InputFileHandle("day16")
	defer f.Close()

	parseLayout(f)

	t1 := func() int {
		traceBeam(0, 0, Right)
		return countEnegized()
	}
	t2 := func() int {
		return traceAllEdgeBeams()
	}
	fmt.Println(common.HandleTasks(t1, t2))

}

func parseLayout(reader io.Reader) {
	data, _ := io.ReadAll(reader)
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		newLine := make([]rune, len(line))
		for i, b := range line {
			newLine[i] = rune(b)
		}
		layout = append(layout, newLine)
		energized = append(energized, make([]byte, len(line)))
	}
}

func resetData(n, m int) {
	energized = nil
	for i := 0; i < n; i++ {
		energized = append(energized, make([]byte, m))
	}
	path = make(map[PathItem]struct{})
}

func traceAllEdgeBeams() int {
	var results []int
	n := len(layout)
	m := len(layout[0])

	for i := 0; i < n; i++ {
		traceBeam(i, 0, Right)
		results = append(results, countEnegized())
		resetData(n, m)
		traceBeam(i, m-1, Left)
		results = append(results, countEnegized())
		resetData(n, m)
	}
	for i := 0; i < m; i++ {
		traceBeam(0, i, Down)
		results = append(results, countEnegized())
		resetData(n, m)
		traceBeam(n-1, i, Up)
		results = append(results, countEnegized())
		resetData(n, m)
	}

	return slices.Max(results)
}

func traceBeam(i, j int, current Direction) {
	if i >= 0 && i < len(layout) {
		if j >= 0 && j < len(layout[i]) {
			item := PathItem{i, j, current}
			if _, ok := path[item]; ok {
				return
			}

			energized[i][j] = 1
			path[item] = struct{}{}

			switch layout[i][j] {
			case '.':
				follow(i, j, current)
			case '-':
				if current == Left || current == Right {
					follow(i, j, current)
				} else {
					split(i, j, current)
				}
			case '|':
				if current == Up || current == Down {
					follow(i, j, current)
				} else {
					split(i, j, current)
				}
			case '/':
				bend(i, j, current, true)
			case '\\':
				bend(i, j, current, false)

			}
		}
	}
}

func follow(i, j int, current Direction) {
	switch current {
	case Up:
		traceBeam(i-1, j, current)
	case Down:
		traceBeam(i+1, j, current)
	case Left:
		traceBeam(i, j-1, current)
	case Right:
		traceBeam(i, j+1, current)
	}
}
func split(i, j int, current Direction) {
	switch {
	case current == Up:
		traceBeam(i, j-1, Left)
		traceBeam(i, j+1, Right)
	case current == Down:
		traceBeam(i, j-1, Left)
		traceBeam(i, j+1, Right)
	case current == Left || current == Right:
		traceBeam(i-1, j, Up)
		traceBeam(i+1, j, Down)
	}
}

func bend(i, j int, current Direction, isForward bool) {
	switch {
	case (current == Up && isForward) || (current == Down && !isForward):
		traceBeam(i, j+1, Right)
	case (current == Down && isForward) || (current == Up && !isForward):
		traceBeam(i, j-1, Left)
	case (current == Left && isForward) || (current == Right && !isForward):
		traceBeam(i+1, j, Down)
	case (current == Right && isForward) || (current == Left && !isForward):
		traceBeam(i-1, j, Up)
	}
}

func countEnegized() int {

	sum := 0
	for _, line := range energized {
		for _, b := range line {
			if b == 1 {
				sum++
			}
		}
	}
	return sum
}
