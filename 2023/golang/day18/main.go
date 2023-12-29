package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type FileRow struct {
	dir Direction
	n   int
}

type Line struct {
	p1, p2 common.Point
}

func NewLine(p1, p2 common.Point) Line {
	var line Line
	if p1.I == p2.I {
		if p1.J < p2.J {
			line = Line{p1, p2}
		} else {
			line = Line{p2, p1}
		}
	}
	if p1.J == p2.J {
		if p1.I < p2.I {
			line = Line{p1, p2}
		} else {
			line = Line{p2, p1}
		}
	}

	return line
}

func (l Line) isHorizontal() bool {
	return l.p1.I == l.p2.I
}

func (l Line) intersects(other Line) bool {
	if l.isHorizontal() && other.isHorizontal() {
		return false
	}

	doIntersect := func(horizontal, vertical Line) bool {
		hi := horizontal.p1.I
		vj := vertical.p1.J
		return horizontal.p1.J <= vj && horizontal.p2.J >= vj && hi >= vertical.p1.I && hi <= vertical.p2.I
	}

	if l.isHorizontal() {
		return doIntersect(l, other)
	}
	if other.isHorizontal() {
		return doIntersect(other, l)
	}

	return false
}

var terrain [][]int
var lines []Line

func main() {
	f := common.InputFileHandle("day18")
	defer f.Close()

	lines := parseFile(f)
	buildTerrain(lines)

	fmt.Println(countMeters())
}

func buildTerrain(fileRows []FileRow) {
	var i, j, startI, startJ int
	var n, m int
	for k, row := range fileRows {
		if k == 0 {
			switch row.dir {
			case Left:
				n, m = 1, row.n+1
				i, j = 0, 0
				startI, startJ = 0, m-1
			case Up:
				n, m = row.n+1, 1
				i, j = 0, 0
				startI, startJ = n-1, 0
			case Right:
				n, m = 1, row.n+1
				i, j = 0, m-1
				startI, startJ = 0, 0
			case Down:
				n, m = row.n+1, 1
				i, j = n-1, 0
				startI, startJ = 0, 0
			}
			continue
		}

		switch row.dir {
		case Left:
			if j-row.n >= 0 {
				j -= row.n
			} else {
				diff := row.n - j
				m += diff
				j = 0
				startJ += diff
			}
		case Up:
			if i-row.n >= 0 {
				i -= row.n
			} else {
				diff := row.n - i
				n += diff
				i = 0
				startI += diff
			}
		case Right:
			if j+row.n < m {
				j += row.n
			} else {
				diff := j + row.n - m + 1
				m += diff
				j = m - 1
			}
		case Down:
			if i+row.n < n {
				i += row.n
			} else {
				diff := i + row.n - n + 1
				n += diff
				i = n - 1
			}
		}
	}
	terrain = make([][]int, n)
	for i := range terrain {
		terrain[i] = make([]int, m)
	}

	markTerrainBorders(startI, startJ, n, m, fileRows)
	fillTerrain(n, m)
}

func parseFile(r io.Reader) []FileRow {
	var res []FileRow
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		data := strings.Fields(line)
		n, _ := strconv.Atoi(data[1])
		switch data[0] {
		case "L":
			res = append(res, FileRow{Left, n})
		case "R":
			res = append(res, FileRow{Right, n})
		case "D":
			res = append(res, FileRow{Down, n})
		case "U":
			res = append(res, FileRow{Up, n})
		}
	}
	return res
}

func markTerrainBorders(startI, startJ, n, m int, fileRows []FileRow) {
	// build terrain borders
	i, j := startI, startJ
	terrain[i][j] = 1
	for _, row := range fileRows {
		p1 := common.Point{i, j}
		var p2 common.Point
		for k := 0; k < row.n; k++ {
			switch row.dir {
			case Left:
				p2 = common.Point{i, j - 1}
				terrain[i][j-1] = 1
				j--
			case Up:
				p2 = common.Point{i - 1, j}
				terrain[i-1][j] = 1
				i--
			case Right:
				p2 = common.Point{i, j + 1}
				terrain[i][j+1] = 1
				j++
			case Down:
				p2 = common.Point{i + 1, j}
				terrain[i+1][j] = 1
				i++
			}
		}
		lines = append(lines, NewLine(p1, p2))

	}
}

func fillTerrain(n, m int) {
	linesByRow := make(map[int][]Line)
	for _, l := range lines {
		current := linesByRow[l.p1.I]
		current = append(current, l)
		linesByRow[l.p1.I] = current
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if terrain[i][j] == 0 && isInPolygon(i, j, linesByRow) {
				terrain[i][j] = 1
			}
		}
	}
}

// ray cast algo
func isInPolygon(i, j int, linesByRow map[int][]Line) bool {
	l1 := NewLine(common.Point{i, j}, common.Point{i, math.MaxInt})
	intersectCount := 0
	currentLines := linesByRow[i] // ignore lines that start at line i
	for _, l2 := range lines {
		if l2.p1.J < j || slices.Contains(currentLines, l2) {
			continue
		}

		if l1.intersects(l2) {
			intersectCount++
		}
	}
	return intersectCount%2 == 1
}

func countMeters() int {
	sum := 0
	for i := range terrain {
		for j := range terrain[i] {
			if terrain[i][j] == 1 {
				sum++
			}
		}
	}
	return sum
}
