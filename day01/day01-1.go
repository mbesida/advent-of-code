package day01

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func Solve1() (int, error) {
	file, err := os.Open("day01/input1")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	sum := 0
	for scanner.Scan() {
		number, err := calibrationValue(scanner.Text())
		if err != nil {
			return 0, err
		}
		sum += number
	}

	return sum, nil
}

func calibrationValue(s string) (int, error) {
	digits := [2]int{0, 0}

	for _, ch := range s {
		if unicode.IsDigit(ch) {
			value := int(ch) - '0'
			if digits[0] == 0 {
				digits[0] = value
			} else {
				digits[1] = value
			}
		}
	}

	var resultString string = "0"
	var err error

	if digits == [2]int{0, 0} {
		err = fmt.Errorf("no digit found on line %s", s)
	} else if digits[1] == 0 {
		resultString = fmt.Sprintf("%[1]d%[1]d", digits[0])
	} else {
		resultString = fmt.Sprintf("%d%d", digits[0], digits[1])
	}

	if err != nil {
		return 0, err
	} else {
		return strconv.Atoi(resultString)
	}
}
