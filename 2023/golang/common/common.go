package common

import (
	"os"
)

func HandleTasks[T any](task1Func func() (T, error), task2Func func() (T, error)) (T, error) {
	var choice string

	if len(os.Args) >= 2 {
		choice = os.Args[1]
	}

	if choice == "2" {
		return task2Func()
	} else {
		return task1Func()
	}
}

func ToTask[T any](fromFunc func() T) func() (T, error) {
	return func() (T, error) {
		result := fromFunc()
		return result, nil
	}
}
