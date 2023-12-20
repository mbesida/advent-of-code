package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

func main() {
	f := common.InputFileHandle("day14")
	defer f.Close()

	rawData, _ := io.ReadAll(f)
	platform := parse(string(rawData))
	tilted := tiltePlatform(platform)

	load := totalLoad(tilted)
	fmt.Println(load)
}

func parse(str string) [][]rune {
	var platform [][]rune
	for _, row := range strings.Split(str, "\n") {
		platform = append(platform, []rune(row))
	}
	return platform
}

func tiltePlatform(platform [][]rune) [][]rune {
	var tiltedPlatform [][]rune
	for j := range platform[0] {
		column := make([]rune, len(platform))
		dotCount := 0
		for i, row := range platform {
			switch {
			case (row[j] == 'O' && dotCount == 0) || row[j] == '#':
				column[i] = row[j]
				dotCount = 0
			case row[j] == 'O' && dotCount != 0:
				column[i] = '.'
				column[i-dotCount] = 'O'
			case row[j] == '.':
				column[i] = '.'
				dotCount++
			}
		}
		tiltedPlatform = append(tiltedPlatform, column)
	}
	return common.TransposeMatrix(tiltedPlatform)
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
