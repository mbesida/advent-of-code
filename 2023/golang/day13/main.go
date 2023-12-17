package main

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

type Note [][]rune

func main() {
	f := common.InputFileHandle("day13")
	defer f.Close()

	bytes, _ := io.ReadAll(f)
	stringifiedData := string(bytes)

	notes := parseNotes(stringifiedData)

	sum := 0
	for _, n := range notes {
		x := calcReflection(n)
		// fmt.Println(x)
		sum += x
	}
	fmt.Println(sum)
}

func parseNotes(data string) []Note {
	groups := strings.Split(data, "\n\n")
	n := len(groups)

	notes := make([]Note, n)
	for i := 0; i < n; i++ {
		rows := strings.Split(groups[i], "\n")
		notes[i] = makeNoteRunes(rows)
	}

	return notes
}

func makeNoteRunes(data []string) [][]rune {
	result := make([][]rune, len(data))
	for i, line := range data {
		result[i] = []rune(line)
	}
	return result
}

func calcReflection(data [][]rune) int {
	transposedData := common.TransposeMatrix(data)

	rows, a := findReflectionLinePosition(data)
	columns, b := findReflectionLinePosition(transposedData)

	var res int
	if a {
		res = 100 * rows
	} else if b {
		res = columns
	}

	return res
}

func findReflectionLinePosition(data [][]rune) (int, bool) {
	var positions []int
	dataLen := len(data)
	for i := 0; i < dataLen-1; i++ {
		if slices.Compare(data[i], data[i+1]) == 0 {
			positions = append(positions, i)
		}
	}

	found := false
	pos := 0
	if len(positions) != 0 {
		for _, i := range positions {
			foundI := true
			behind, forward := i-1, i+2
			for behind >= 0 && forward < dataLen {
				if slices.Compare(data[behind], data[forward]) != 0 {
					foundI = false
					break
				}
				behind--
				forward++
			}
			if foundI {
				found = foundI
				pos = i
				break
			}
		}
	}

	return pos + 1, found
}
