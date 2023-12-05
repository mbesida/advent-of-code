package common

import (
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

func getChoice() string {
	var choice string

	if len(os.Args) >= 2 {
		choice = os.Args[1]
	}

	return choice
}
