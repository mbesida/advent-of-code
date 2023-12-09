package common

import (
	"fmt"
	"log"
	"os"
)

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

func HandleValue[T any](value1 T, value2 T) T {
	choice := getChoice()

	if choice == "2" {
		return value1
	} else {
		return value2
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
