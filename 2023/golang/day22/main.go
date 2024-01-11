package main

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

type Coords struct {
	x, y, z int
}

type Brick struct {
	start, end Coords
}

func (b Brick) holds(other Brick) bool {
	if b.end.z+1 == other.start.z {
		return max(b.start.x, other.start.x) <= min(b.end.x, other.end.x) &&
			max(b.start.y, other.start.y) <= min(b.end.y, other.end.y)
	}
	return false
}

func main() {
	f := common.InputFileHandle("day22")
	defer f.Close()

	bricks := parseData(f)

	res := process(bricks)
	fmt.Println(res)
}

func sortFunc(a, b Brick) int {
	if a.start.z < b.start.z {
		return -1
	}
	return 1
}

func process(bricks []Brick) int {
	fallen := fall(bricks)
	slices.SortFunc(fallen, sortFunc)
	startMap := make(map[int][]Brick)
	endMap := make(map[int][]Brick)
	for _, b := range fallen {
		startMap[b.start.z] = append(startMap[b.start.z], b)
		endMap[b.end.z] = append(endMap[b.end.z], b)
	}

	var safeToRemove []Brick
	for _, b := range fallen {
		currentLevel := endMap[b.end.z]
		nextLevel := startMap[b.end.z+1]
		var holdedByB []Brick
		for _, nb := range nextLevel {
			if b.holds(nb) {
				holdedByB = append(holdedByB, nb)
			}
		}
		if len(holdedByB) == 0 {
			safeToRemove = append(safeToRemove, b)
			continue
		}
		holdedByOthers := make(map[Brick]int)
		for _, currentLevelB := range currentLevel {
			if currentLevelB != b {
				for _, hb := range holdedByB {
					if currentLevelB.holds(hb) {
						holdedByOthers[hb]++
					}
				}
			}
		}

		if len(holdedByB) == len(holdedByOthers) {
			safeToRemove = append(safeToRemove, b)
		}
	}

	return len(safeToRemove)
}

func fall(bricks []Brick) []Brick {
	slices.SortFunc(bricks, sortFunc) // sort by start.z

	var mX, mY, mZ int
	for _, b := range bricks {
		mX = max(mX, b.end.x)
		mY = max(mY, b.end.y)
		mZ = max(mZ, b.end.z)
	}
	grid := make([][][]int, mX+1)
	for i := 0; i <= mX; i++ {
		grid[i] = make([][]int, mY+1)
		for j := 0; j <= mY; j++ {
			grid[i][j] = make([]int, mZ+1)
			grid[i][j][0] = 1
		}
	}

	// normalize
	minZ := bricks[0].start.z
	for i, b := range bricks {
		startZ := b.start.z - minZ + 1
		endZ := b.end.z - minZ + 1
		bricks[i] = Brick{Coords{b.start.x, b.start.y, startZ}, Coords{b.end.x, b.end.y, endZ}}
	}

	fallen := make([]Brick, len(bricks))
	fallen[0] = bricks[0]
	occupy(grid, bricks[0])

	for i := 1; i < len(bricks); i++ {
		if fallen[i-1].start.z == bricks[i].start.z {
			fallen[i] = bricks[i]
		} else {
			b := bricks[i]
			currentZ := b.start.z
			for !checkOccupation(grid, b, currentZ) {
				currentZ--
			}
			currentZ++
			diff := b.start.z - currentZ
			endZ := b.end.z - diff
			fallen[i] = Brick{Coords{b.start.x, b.start.y, currentZ}, Coords{b.end.x, b.end.y, endZ}}
		}
		occupy(grid, fallen[i])
	}

	return fallen
}

func occupy(grid [][][]int, b Brick) {
	for i := 0; i <= b.end.x-b.start.x; i++ {
		for j := 0; j <= b.end.y-b.start.y; j++ {
			for k := 0; k <= b.end.z-b.start.z; k++ {
				grid[b.start.x+i][b.start.y+j][b.start.z+k] = 1
			}
		}
	}
}

func checkOccupation(grid [][][]int, b Brick, currentZ int) bool {
	for i := b.start.x; i <= b.end.x; i++ {
		for j := b.start.y; j <= b.end.y; j++ {
			if grid[i][j][currentZ] == 1 {
				return true
			}
		}
	}
	return false
}

func parseData(reader io.Reader) []Brick {
	bs, _ := io.ReadAll(reader)
	var bricks []Brick

	for _, row := range bytes.Split(bs, []byte("\n")) {
		bricks = append(bricks, parseBrick(string(row)))
	}

	return bricks
}

func parseBrick(line string) Brick {
	data := strings.Split(line, "~")
	start := parseCoords(data[0])
	end := parseCoords(data[1])
	if start.x > end.x || (start.x == end.x && start.y > end.y) {
		start, end = end, start
	}
	return Brick{start, end}
}

func parseCoords(s string) Coords {
	data := strings.Split(s, ",")
	x, _ := strconv.Atoi(data[0])
	y, _ := strconv.Atoi(data[1])
	z, _ := strconv.Atoi(data[2])
	return Coords{x, y, z}
}

func printBricks(bricks []Brick) {
	for _, b := range bricks {
		fmt.Println(b)
	}
}
