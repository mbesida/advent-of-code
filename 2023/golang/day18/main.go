package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"regexp"
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

func main() {
	f := common.InputFileHandle("day18")
	defer f.Close()

	fileRows := parseFile(f)
	points := buildPoints(fileRows)
	makePositiveCoordinates(points)
	points = reorganizePoints(points)
	area := calcArea(points)

	fmt.Println(area)
}

func calcArea(points []common.Point) int {
	n := len(points)
	is := make([]int, n)
	js := make([]int, n)
	for k := 0; k < n; k++ {
		is[k] = points[k].I
		js[k] = points[k].J
	}

	// shoelace area
	doubleArea := 0
	for i := 0; i < n-1; i++ {
		doubleArea += is[i] * js[i+1]
		doubleArea -= is[i+1] * js[i]
	}

	perimeter := 0
	for i := 1; i < n; i++ {
		prev := points[i-1]
		p := points[i]
		if prev.I == p.I {
			perimeter += int(math.Abs(float64(p.J) - float64(prev.J)))
		} else {
			perimeter += int(math.Abs(float64(p.I) - float64(prev.I)))
		}
	}

	// Picks's theorem
	countInsidePoints := (doubleArea + 2 - perimeter) / 2

	return countInsidePoints + perimeter
}

func reorganizePoints(points []common.Point) []common.Point {
	reorganizedPoints := points
	first := points[0]
	second := points[1]
	last := points[len(points)-1]
	if first.I == second.I {
		if first.J < second.J {
			tail := reorganizedPoints[1 : len(reorganizedPoints)-1]
			slices.Reverse(tail)
			reorganizedPoints = append([]common.Point{first}, tail...)
			reorganizedPoints = append(reorganizedPoints, last)
		}
	}
	return reorganizedPoints
}

func buildPoints(fileRows []FileRow) []common.Point {
	var points []common.Point
	i, j := 0, 0
	points = append(points, common.Point{i, j})
	for k := 0; k < len(fileRows); k++ {
		diff := fileRows[k].n
		switch fileRows[k].dir {
		case Left:
			j -= diff
		case Right:
			j += diff
		case Up:
			i -= diff
		case Down:
			i += diff
		}
		points = append(points, common.Point{i, j})
	}
	if points[0] != points[len(points)-1] {
		panic("last point not equals to 1st point")
	}
	return points
}

func makePositiveCoordinates(points []common.Point) {
	minPointI := slices.MinFunc(points, func(a, b common.Point) int {
		if a.I < b.I {
			return -1
		}
		return 1
	})
	minPointJ := slices.MinFunc(points, func(a, b common.Point) int {
		if a.J < b.J {
			return -1
		}
		return 1
	})
	var deltaY, deltaX int
	if minPointJ.J < 0 {
		deltaX = -minPointJ.J
	}
	if minPointI.I < 0 {
		deltaY = -minPointI.I
	}
	for i := 0; i < len(points); i++ {
		current := points[i]
		points[i] = common.Point{current.I + deltaY, current.J + deltaX}
	}
}

func parseFile(r io.Reader) []FileRow {
	var res []FileRow
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		data := strings.Fields(line)
		n, _ := strconv.Atoi(data[1])
		re := regexp.MustCompile(`\(#([0-9a-f]{6})\)`)
		colorData := re.FindStringSubmatch(data[2])
		if len(colorData) < 2 {
			log.Fatalf("incorrect format of color %s", line)
		}
		hexString := colorData[1]
		number, _ := strconv.ParseInt(hexString[:len(hexString)-1], 16, 64)
		var hexDir string
		switch hexString[len(hexString)-1] {
		case '2':
			hexDir = "L"
		case '0':
			hexDir = "R"
		case '1':
			hexDir = "D"
		case '3':
			hexDir = "U"
		}

		direction := common.HandleValue(data[0], hexDir)
		x := common.HandleValue(n, int(number))

		switch direction {
		case "L":
			res = append(res, FileRow{Left, x})
		case "R":
			res = append(res, FileRow{Right, x})
		case "D":
			res = append(res, FileRow{Down, x})
		case "U":
			res = append(res, FileRow{Up, x})
		}
	}
	return res
}
