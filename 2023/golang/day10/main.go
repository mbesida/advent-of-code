package main

import (
	"fmt"
	"io"
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
	x       int
	y       int
	network Map
}

func main() {
	f := common.InputFileHandle("day10")
	defer f.Close()

	data, _ := io.ReadAll(f)
	dataString := string(data)

	parsedData := parseData(dataString)
	res := findFarthest(parsedData)

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

	return ParsedData{initX, initY, network}
}

func findFarthest(data ParsedData) int {
	return recurse(data.x, data.y, 0, data.network, undefined) / 2
}

func recurse(i, j, agg int, grid Map, previous Direction) int {
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
		return -1
	}
	if grid[nextI][nextJ] == 'S' {
		return agg + 1
	}
	return recurse(nextI, nextJ, agg+1, grid, previous)
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
