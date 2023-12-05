package main

import (
	"bufio"
	"fmt"
	"strconv"
	"unicode"

	"github.com/mbesida/advent-of-code-2023/common"
)

type VisitedMap = map[[2]int]struct{}
type Matrix = [][]rune

func main() {
	f := common.InputFileHandle("day03")
	defer f.Close()

	matrix := buildMatrix(bufio.NewScanner(f))

	result := partNumbersSum(matrix)
	fmt.Println(result)

}

func isSymbol(s rune) bool {
	task1 := func() bool {
		return (unicode.IsPunct(s) || unicode.IsSymbol(s)) && s != '.'
	}

	task2 := func() bool {
		return s == '*'
	}
	return common.HandleTasks(task1, task2)
}

func buildMatrix(scanner *bufio.Scanner) [][]rune {
	var matrix [][]rune
	for scanner.Scan() {
		fileRow := scanner.Text()
		matrixRow := []rune(fileRow)
		matrix = append(matrix, matrixRow)
	}
	return matrix
}

func partNumbersSum(matrix Matrix) int {
	var sum int
	visitedIndices := make(map[[2]int]struct{})

	find := func(i, j int, adjacentNumbers []int) []int {
		return findNumber(i, j, matrix, visitedIndices, adjacentNumbers)
	}

	for i, row := range matrix {
		for j, s := range row {
			if isSymbol(s) {
				var adjacentNumbers []int
				if i > 0 {
					adjacentNumbers = find(i-1, j, adjacentNumbers)
				}
				if j > 0 {
					adjacentNumbers = find(i, j-1, adjacentNumbers)
				}
				if i < len(matrix)-1 {
					adjacentNumbers = find(i+1, j, adjacentNumbers)
				}
				if j < len(row)-1 {
					adjacentNumbers = find(i, j+1, adjacentNumbers)
				}
				if i > 0 && j > 0 {
					adjacentNumbers = find(i-1, j-1, adjacentNumbers)
				}
				if i < len(matrix)-1 && j < len(row)-1 {
					adjacentNumbers = find(i+1, j+1, adjacentNumbers)
				}
				if i > 0 && j < len(row)-1 {
					adjacentNumbers = find(i-1, j+1, adjacentNumbers)
				}
				if i < len(matrix)-1 && j > 0 {
					adjacentNumbers = find(i+1, j-1, adjacentNumbers)
				}

				t1 := func() int {
					total := 0
					for _, v := range adjacentNumbers {
						total += v
					}
					return total
				}

				t2 := func() int {
					total := 0
					if len(adjacentNumbers) == 2 {
						total = adjacentNumbers[0] * adjacentNumbers[1]
					}
					return total
				}
				sum += common.HandleTasks(t1, t2)
			}
		}
	}
	return sum
}

func findNumber(i, j int, matrix Matrix, visited VisitedMap, adjacentNumbers []int) []int {
	_, ok := visited[[2]int{i, j}]
	if unicode.IsDigit(matrix[i][j]) && !ok {
		result := parseNumber(i, j, matrix, visited)
		adjacentNumbers = append(adjacentNumbers, result)
	}
	return adjacentNumbers
}

func parseNumber(i, j int, matrix Matrix, visited VisitedMap) int {
	visited[[2]int{i, j}] = struct{}{}
	row := matrix[i]
	data := []rune{row[j]}
	for k := j - 1; k >= 0; k-- {
		if !unicode.IsDigit(row[k]) {
			break
		} else {
			visited[[2]int{i, k}] = struct{}{}
			data = append([]rune{row[k]}, data...)
		}
	}
	for k := j + 1; k < len(row); k++ {
		if !unicode.IsDigit(row[k]) {
			break
		} else {
			visited[[2]int{i, k}] = struct{}{}
			data = append(data, row[k])
		}
	}

	number, _ := strconv.Atoi(string(data))
	return number
}
