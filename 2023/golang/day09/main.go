package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

func main() {
	f := common.InputFileHandle("day09")
	defer f.Close()

	scanner := bufio.NewScanner(f)

	sum := 0
	for scanner.Scan() {
		numbers := strings.Fields(scanner.Text())
		values := make([]int, len(numbers))
		for i, v := range numbers {
			x, _ := strconv.Atoi(v)
			values[i] = x
		}
		p := common.HandleTasks(func() int { return predictRowValue(values) }, func() int { return historyRowValue(values) })
		sum += p
	}

	fmt.Println("result is", sum)
}

func predictRowValue(values []int) int {
	currentValues := values
	var lastValues []int
	for {
		lastValues = append(lastValues, currentValues[len(currentValues)-1])
		res := buildNextValues(currentValues)
		currentValues = res.nextValues
		if res.stopCondition {
			break
		}
	}

	agg := 0
	for _, v := range lastValues {
		agg += v
	}

	return agg
}

func historyRowValue(values []int) int {

	currentValues := values
	var firstValues []int
	for {
		firstValues = append(firstValues, currentValues[0])
		res := buildNextValues(currentValues)
		currentValues = res.nextValues
		if res.stopCondition {
			break
		}
	}

	agg := 0
	for i := len(firstValues) - 1; i >= 0; i-- {
		agg = firstValues[i] - agg

	}
	return agg
}

func buildNextValues(currentValues []int) struct {
	nextValues    []int
	stopCondition bool
} {
	nextValues := make([]int, len(currentValues)-1)
	stopCondition := true
	for i := 0; i < len(nextValues); i++ {
		nextValues[i] = currentValues[i+1] - currentValues[i]
		if nextValues[i] != 0 && stopCondition {
			stopCondition = false
		}
	}

	return struct {
		nextValues    []int
		stopCondition bool
	}{nextValues, stopCondition}
}
