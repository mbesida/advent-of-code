package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
)

type VisitedMap = map[[2]int]struct{}
type Matrix = [][]rune

func main() {
	file, err := os.Open("day03/input")
	if err != nil {
		log.Fatal("can't open input data file")
	}
	defer file.Close()

	matrix := buildMatrix(bufio.NewScanner(file))

	result := partNumbersSum(matrix)
	fmt.Println(result)

}

func isSymbol(s rune) bool {
	return (unicode.IsPunct(s) || unicode.IsSymbol(s)) && s != '.'
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

	find := func(i, j int) int {
		return findNumber(i, j, matrix, visitedIndices)
	}

	for i, row := range matrix {
		for j, s := range row {
			if isSymbol(s) {
				if i > 0 {
					sum += find(i-1, j)
				}
				if j > 0 {
					sum += find(i, j-1)
				}
				if i < len(matrix)-1 {
					sum += find(i+1, j)
				}
				if j < len(row)-1 {
					sum += find(i, j+1)
				}
				if i > 0 && j > 0 {
					sum += find(i-1, j-1)
				}
				if i < len(matrix)-1 && j < len(row)-1 {
					sum += find(i+1, j+1)
				}
				if i > 0 && j < len(row)-1 {
					sum += find(i-1, j+1)
				}
				if i < len(matrix)-1 && j > 0 {
					sum += find(i+1, j-1)
				}
			}
		}
	}
	return sum
}

func findNumber(i, j int, matrix Matrix, visited VisitedMap) int {
	var result int
	_, ok := visited[[2]int{i, j}]
	if unicode.IsDigit(matrix[i][j]) && !ok {
		result = parseNumber(i, j, matrix, visited)
	}
	return result
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
