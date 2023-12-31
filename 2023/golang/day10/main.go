package main

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

type Direction int

const (
	undefined Direction = iota
	north
	east
	south
	west
)

type Map = [][]rune
type ParsedData struct {
	startingPoint common.Point
	network       Map
}

type Path []common.Point

func main() {
	f := common.InputFileHandle("day10")
	defer f.Close()

	data, _ := io.ReadAll(f)
	dataString := string(data)

	parsedData := parseData(dataString)

	t1 := func() int {
		return findFarthest(parsedData)
	}
	t2 := func() int {
		path := buildPath(parsedData)
		return enclosedTiles(parsedData, path)
	}
	res := common.HandleTasks(t1, t2)

	fmt.Println(res)
}

func parseData(data string) ParsedData {
	rows := strings.Split(data, "\n")

	var network Map
	initX := -1
	initY := -1
	for k, row := range rows {
		n := len(row)
		runeRow := make([]rune, n)
		for i := 0; i < n; i++ {
			r := rune(row[i])
			if initX == -1 && initY == -1 && r == 'S' {
				initX = k
				initY = i
			}
			runeRow[i] = r
		}
		network = append(network, runeRow)
	}

	return ParsedData{common.Point{initX, initY}, network}
}

func buildPath(data ParsedData) Path {
	point := data.startingPoint
	var path Path
	path = append(path, point)
	path = recurse(point.I, point.J, data.network, undefined, path)

	return path
}

func findFarthest(data ParsedData) int {
	path := buildPath(data)
	return len(path) / 2
}

func enclosedTiles(data ParsedData, path Path) int {
	count := 0
	for i, row := range data.network {
		for j := range row {
			if !slices.Contains(path, common.Point{i, j}) {
				intersections := 0
				for k := j + 1; k < len(row); k++ {
					p := common.Point{i, k}
					if slices.Contains(path, p) && strings.ContainsRune("|F7", row[k]) {
						intersections++
					}
				}
				if intersections%2 == 1 {
					count++
				}
			}
		}
	}
	return count
}

func recurse(i, j int, grid Map, previous Direction, path Path) Path {
	nextI, nextJ := 0, 0
	switch {
	case northWakable(i, j, grid) && previous != north:
		nextI, nextJ = i-1, j
		previous = south
	case eastWakable(i, j, grid) && previous != east:
		nextI, nextJ = i, j+1
		previous = west
	case southWakable(i, j, grid) && previous != south:
		nextI, nextJ = i+1, j
		previous = north
	case westWakable(i, j, grid) && previous != west:
		nextI, nextJ = i, j-1
		previous = east
	default:
		nextI, nextJ = -1, -1
	}
	if nextI == -1 && nextJ == -1 {
		return []common.Point{}
	}

	pointToAdd := common.Point{I: nextI, J: nextJ}
	if grid[nextI][nextJ] == 'S' {
		first := path[0]
		second := path[1]

		if grid[second.I][second.J] == grid[pointToAdd.I][pointToAdd.J] {
			grid[first.I][first.J] = grid[pointToAdd.I][pointToAdd.J]
		} else if previous == south && second.J > j || previous == east && second.I > i {
			grid[first.I][first.J] = 'F'
		} else if previous == west && second.I > i || previous == south && second.J < j {
			grid[first.I][first.J] = '7'
		} else if previous == north && second.J < j || previous == west && second.I < i {
			grid[first.I][first.J] = 'J'
		} else if previous == north && second.J > j || previous == east && second.I < i {
			grid[first.I][first.J] = 'L'
		}

		return path
	}

	path = append(path, pointToAdd)

	return recurse(nextI, nextJ, grid, previous, path)
}

func northWakable(i, j int, grid Map) bool {
	return (i-1 >= 0) && strings.ContainsRune("|LJS", grid[i][j]) && strings.ContainsRune("|F7S", grid[i-1][j])
}
func eastWakable(i, j int, grid Map) bool {
	return (j+1 < len(grid[i])) && strings.ContainsRune("-LFS", grid[i][j]) && strings.ContainsRune("-J7S", grid[i][j+1])
}
func southWakable(i, j int, grid Map) bool {
	return (i+1 < len(grid)) && strings.ContainsRune("|7FS", grid[i][j]) && strings.ContainsRune("|LJS", grid[i+1][j])
}
func westWakable(i, j int, grid Map) bool {
	return (j-1 >= 0) && strings.ContainsRune("-J7S", grid[i][j]) && strings.ContainsRune("-FLS", grid[i][j-1])
}
