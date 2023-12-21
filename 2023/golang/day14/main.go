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
	North Direction = iota
	West
	South
	East
)

func main() {
	f := common.InputFileHandle("day14")
	defer f.Close()

	rawData, _ := io.ReadAll(f)
	platform := parse(string(rawData))

	t1 := func() int {
		tilted := tiltePlatform(platform, North)
		load := totalLoad(tilted)
		return load
	}
	t2 := func() int {
		lastHash := hash(platform)
		seen := make(map[string]bool)
		seen[lastHash] = true
		sl := []string{lastHash}
		var counter int
		for {
			counter++
			platform = doCycle(platform)
			lastHash = hash(platform)
			if seen[lastHash] {
				break
			}
			seen[lastHash] = true
			sl = append(sl, lastHash)
		}

		firstIndex := slices.Index(sl, lastHash)
		length := counter - firstIndex
		hash := sl[firstIndex+(1000000000-firstIndex)%length]
		platform = parse(hash)

		return totalLoad(platform)
	}

	fmt.Println(common.HandleTasks(t1, t2))
}

func parse(str string) [][]rune {
	var platform [][]rune
	for _, row := range strings.Split(str, "\n") {
		platform = append(platform, []rune(row))
	}
	return platform
}

func printPlatform(platform [][]rune) {
	for _, row := range platform {
		fmt.Println(string(row))
	}
	fmt.Println()
}

func tiltePlatform(platform [][]rune, direction Direction) [][]rune {
	var tiltedPlatform [][]rune

	if slices.Contains([]Direction{North, South}, direction) {
		for j := range platform[0] {
			column := make([]rune, len(platform))
			dotCount := 0
			if direction == North {
				for i := 0; i < len(platform); i++ {
					northSouthTilt(column, platform[i][j], i, &dotCount, direction)
				}
			} else {
				for i := len(platform) - 1; i >= 0; i-- {
					northSouthTilt(column, platform[i][j], i, &dotCount, direction)
				}
			}
			tiltedPlatform = append(tiltedPlatform, column)
		}
		tiltedPlatform = common.TransposeMatrix(tiltedPlatform)
	} else {
		// West and East tilt is the same as North and south but on transposed platform matrix,
		// then result must be transposed back
		transposedPlatform := common.TransposeMatrix(platform)
		var tiltedTransposed [][]rune
		if direction == West {
			tiltedTransposed = tiltePlatform(transposedPlatform, North)
		} else {
			tiltedTransposed = tiltePlatform(transposedPlatform, South)
		}
		tiltedPlatform = common.TransposeMatrix(tiltedTransposed)
	}
	return tiltedPlatform
}

func northSouthTilt(column []rune, r rune, index int, dotCount *int, d Direction) {
	switch {
	case (r == 'O' && *dotCount == 0) || r == '#':
		column[index] = r
		*dotCount = 0
	case r == 'O' && *dotCount != 0:
		column[index] = '.'
		if d == North {
			column[index-*dotCount] = 'O'
		} else {
			column[index+*dotCount] = 'O'
		}
	case r == '.':
		column[index] = '.'
		*dotCount++
	}
}

func totalLoad(tiltedPlatform [][]rune) int {
	total := 0
	length := len(tiltedPlatform)
	for i, row := range tiltedPlatform {
		for _, r := range row {
			if r == 'O' {
				total += (length - i)
			}
		}
	}
	return total
}

func doCycle(platform [][]rune) [][]rune {
	result := platform
	for _, d := range []Direction{North, West, South, East} {
		result = tiltePlatform(result, d)
	}
	return result
}

func hash(platform [][]rune) string {
	var sb strings.Builder
	for i, row := range platform {
		sb.WriteString(string(row))
		if i != len(platform)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
