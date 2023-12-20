package main

import (
	"fmt"
	"io"
	"math/bits"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

func main() {
	f := common.InputFileHandle("day13")
	defer f.Close()

	bytes, _ := io.ReadAll(f)
	stringifiedData := string(bytes)

	notes := parseNotes(stringifiedData)

	sum := 0
	for _, n := range notes {
		x := calcReflection(n)
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
		notes[i] = NewNote(matrix)
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

func calcReflection(note Note) int {
	potentialRows := potentialReflectionPoints(note.rows)
	potentialColumns := potentialReflectionPoints(note.columns)
	rowReflection := findReflection(note.rows, potentialRows)
	columnReflection := findReflection(note.columns, potentialColumns)

	t1 := func() int {
		if rowReflection != -1 {
			return calculate(rowReflection, true)
		}
		if columnReflection != -1 {
			return calculate(columnReflection, false)
		}
		return 0
	}
	t2 := func() int {
		res := findAlternative(note.rows, rowReflection, potentialRows)
		if res != -1 {
			return calculate(res, true)
		}
		res = findAlternative(note.columns, columnReflection, potentialColumns)
		if res != -1 {
			return calculate(res, false)
		}
		return 0
	}
	return common.HandleTasks(t1, t2)

}

func findAlternative(data []int64, pos int, potential []int) int {
	dataLen := len(data)
	if pos == -1 || len(potential) > 1 {
		for _, index := range potential {
			if pos != index {
				if checkPosition(index, data, false) {
					return index
				}
			}
		}
	}

	for i := 0; i < dataLen-1; i++ {
		if isOneStepDistance(data[i], data[i+1]) {
			if checkPosition(i, data, true) {
				return i
			}
		}
	}

	return -1
}

func findReflection(data []int64, potentialPositions []int) int {
	for _, i := range potentialPositions {
		if checkPosition(i, data, true) {
			return i
		}
	}
	return -1
}

// directComparisonOnly means that if we sould check only equality of elements
func checkPosition(i int, data []int64, directComparisonOnly bool) bool {
	behind, forward := i-1, i+2
	for behind >= 0 && forward < len(data) {
		if directComparisonOnly {
			if data[behind] != data[forward] {
				return false
			}
		} else {
			if data[behind] != data[forward] && !isOneStepDistance(data[behind], data[forward]) {
				return false
			}
		}
		behind--
		forward++
	}
	return true
}

func isOneStepDistance(x, y int64) bool {
	xor := x ^ y
	return bits.OnesCount64(uint64(xor)) == 1
}

type Note struct {
	rows, columns []int64
}

func NewNote(matrix [][]rune) Note {
	transposedMatrix := common.TransposeMatrix(matrix)
	rows := convertToNumbers(matrix)
	columns := convertToNumbers(transposedMatrix)
	return Note{rows, columns}
}

func convertToNumbers(data [][]rune) []int64 {
	numbers := make([]int64, len(data))
	for i, row := range data {
		number, _ := strconv.ParseInt(string(row), 2, 64)
		numbers[i] = number
	}
	return numbers
}

func potentialReflectionPoints(data []int64) []int {
	var points []int
	for i := 0; i < len(data)-1; i++ {
		if data[i] == data[i+1] {
			points = append(points, i)
		}
	}
	return points
}

func calculate(index int, isRow bool) int {
	if isRow {
		return 100 * (index + 1)
	}
	return index + 1
}
