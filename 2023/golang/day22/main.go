package main

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
	"golang.org/x/exp/maps"
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

func (b Brick) supports(nextLevel []Brick) []Brick {
	var supportedByB []Brick
	for _, nb := range nextLevel {
		if b.holds(nb) {
			supportedByB = append(supportedByB, nb)
		}
	}
	return supportedByB
}

func (b Brick) supportsExclusive(currentLevel, nextLevel []Brick) []Brick {
	var allSupportedByB []Brick
	for _, nb := range nextLevel {
		if b.holds(nb) {
			allSupportedByB = append(allSupportedByB, nb)
		}
	}
	supportedOnlyByB := allSupportedByB
	for _, other := range currentLevel {
		if b != other {
			supportedOnlyByB = slices.DeleteFunc(supportedOnlyByB, func(a Brick) bool { return other.holds(a) })
		}
	}
	return supportedOnlyByB
}

func main() {
	f := common.InputFileHandle("day22")
	defer f.Close()

	bricks := parseData(f)

	res := process(bricks)
	fmt.Println(res)
}

func sortByStart(a, b Brick) int {
	if a.start.z < b.start.z {
		return -1
	}
	return 1
}
func sortByEnd(a, b Brick) int {
	if a.end.z < b.end.z {
		return -1
	}
	return 1
}

func process(bricks []Brick) int {
	fallen := fall(bricks)
	slices.SortFunc(fallen, sortByStart)
	startMap := make(map[int][]Brick)
	endMap := make(map[int][]Brick)
	for _, b := range fallen {
		startMap[b.start.z] = append(startMap[b.start.z], b)
		endMap[b.end.z] = append(endMap[b.end.z], b)
	}
	t1 := func() int {
		return task1(fallen, startMap, endMap)
	}
	t2 := func() int {
		return task2(fallen, startMap, endMap)
	}

	return common.HandleTasks(t1, t2)
}

func task1(bricks []Brick, startMap, endMap map[int][]Brick) int {
	countRemovable := 0
	for _, b := range bricks {
		currentLevel := endMap[b.end.z]
		nextLevel := startMap[b.end.z+1]

		if len(b.supportsExclusive(currentLevel, nextLevel)) == 0 {
			countRemovable++
		}
	}

	return countRemovable
}

// slow, but does a job
func task2(bricks []Brick, startMap, endMap map[int][]Brick) int {
	supportsOnly := make(map[Brick][]Brick)
	for _, b := range bricks {
		currentLevel := endMap[b.end.z]
		nextLevel := startMap[b.end.z+1]
		supprtedOnly := b.supportsExclusive(currentLevel, nextLevel)
		if len(supprtedOnly) != 0 {
			supportsOnly[b] = supprtedOnly
		}
	}
	counter := 0
	for supporter, suportees := range supportsOnly {
		currentBricks := make(map[Brick]struct{})
		currentBricks[supporter] = struct{}{}
		currentSupportees := suportees
		for _, sp := range suportees {
			currentBricks[sp] = struct{}{}
		}

		for {
			var newCurrentSupportees []Brick
			for _, sp := range currentSupportees {
				prevLevel := endMap[sp.start.z-1]
				nextLevel := startMap[sp.end.z+1]
				var skip bool
				for _, plb := range prevLevel {
					if plb.holds(sp) && !slices.Contains(maps.Keys(currentBricks), plb) {
						skip = true
						break
					}
				}
				if skip {
					continue
				}
				currentBricks[sp] = struct{}{}
				newCurrentSupportees = append(newCurrentSupportees, sp.supports(nextLevel)...)
			}
			if len(newCurrentSupportees) == 0 {
				break
			}
			currentSupportees = newCurrentSupportees
		}
		counter += len(currentBricks) - 1
	}
	return counter
}

func fall(bricks []Brick) []Brick {
	slices.SortFunc(bricks, sortByStart) // sort by start.z

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

	for i := 0; i < len(bricks); i++ {
		b := bricks[i]
		currentZ := b.start.z
		for !checkOccupation(grid, b, currentZ) {
			currentZ--
		}
		currentZ++
		diff := b.start.z - currentZ
		endZ := b.end.z - diff
		fallen[i] = Brick{Coords{b.start.x, b.start.y, currentZ}, Coords{b.end.x, b.end.y, endZ}}
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
