package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

type Note struct {
	rows, columns []int64
}

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

	rawData := make([][][]rune, n)
	for i := 0; i < n; i++ {
		rows := strings.Split(groups[i], "\n")
		rawData[i] = makeNoteRunes(rows)
	}

	notes := make([]Note, len(rawData))
	for i, matrix := range rawData {
		rows := convertToNumbers(matrix)
		columns := convertToNumbers(common.TransposeMatrix(matrix))
		notes[i] = Note{rows, columns}
	}

	return notes
}

func makeNoteRunes(data []string) [][]rune {
	result := make([][]rune, len(data))
	for i, line := range data {
		runes := []rune(line)
		for j, r := range runes {
			if r == '#' {
				runes[j] = '1'
			} else {
				runes[j] = '0'
			}
		}
		result[i] = runes
	}
	return result
}

func convertToNumbers(data [][]rune) []int64 {
	numbers := make([]int64, len(data))
	for i, row := range data {
		number, _ := strconv.ParseInt(string(row), 2, 64)
		numbers[i] = number
	}
	return numbers
}

func calcReflection(note Note) int {
	rows, foundRows := findReflectionLinePosition(note.rows)
	columns, foundColumns := findReflectionLinePosition(note.columns)

	t1 := func() int {
		var res int
		if foundRows {
			res = 100 * rows
		} else if foundColumns {
			res = columns
		}
		return res
	}

	t2 := func() int {
		// var res int
		// if a {
		// 	for i, row := range data {
		// 		for j, v := range row {

		// 		}
		// 	}
		// 	res = 100 * rows
		// } else if b {
		// 	res = columns
		// }
		return 0
	}

	return common.HandleTasks(t1, t2)
}

func findReflectionLinePosition(data []int64) (int, bool) {
	var positions []int
	dataLen := len(data)
	for i := 0; i < dataLen-1; i++ {
		if data[i] == data[i+1] {
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
				if data[behind] != data[forward] {
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
