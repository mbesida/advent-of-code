package main

import (
	"strconv"

	"github.com/mbesida/advent-of-code-2023/common"
)

type Reflection struct {
	index int
	pairs [][2]int
}

func NewReflection(index, length int) *Reflection {
	if index >= length-1 {
		panic("position can't be last in note")
	}
	var pairs [][2]int
	pairs = append(pairs, [2]int{index, index + 1})
	behind, forward := index-1, index+2
	for behind >= 0 && forward < length {
		pairs = append(pairs, [2]int{behind, forward})
		behind--
		forward++
	}
	return &Reflection{index, pairs}
}

func (r *Reflection) calc(isRow bool) int {
	if isRow {
		return 100 * (r.index + 1)
	}
	return r.index + 1
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
