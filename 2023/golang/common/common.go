package common

import (
	"fmt"
	"log"
	"os"
)

type Point struct {
	I, J int
}

func HandleTasksWithErrors[T any](task1Func func() (T, error), task2Func func() (T, error)) (T, error) {
	choice := getChoice()

	if choice == "2" {
		return task2Func()
	} else {
		return task1Func()
	}
}

func HandleTasks[T any](task1Func func() T, task2Func func() T) T {
	choice := getChoice()

	if choice == "2" {
		return task2Func()
	} else {
		return task1Func()
	}
}

func ExecuteActions(action1, action2 func()) {
	choice := getChoice()

	if choice == "2" {
		action2()
	} else {
		action1()
	}
}

func HandleValue[T any](value1 T, value2 T) T {
	choice := getChoice()

	if choice == "2" {
		return value2
	} else {
		return value1
	}
}

func InputFileHandle(day string) *os.File {
	inputFile := fmt.Sprintf("%s/input", day)
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("can't open input file")
	}
	return file
}

func getChoice() string {
	var choice string

	if len(os.Args) >= 2 {
		choice = os.Args[1]
	}

	return choice
}

func TransposeMatrix[T any](matrix [][]T) [][]T {
	if len(matrix) <= 0 {
		panic("zero sized matrix can't be transposed")
	}
	transposed := make([][]T, len(matrix[0]))
	for i := 0; i < len(matrix[0]); i++ {
		column := make([]T, len(matrix))
		for j, row := range matrix {
			column[j] = row[i]
		}
		transposed[i] = column
	}
	return transposed
}
