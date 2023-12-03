package day01

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"unicode"

	"github.com/s0rg/trie"
)

func Solve1() (int, error) {
	return commonLogic(digits1)
}

func Solve2() (int, error) {
	return commonLogic(digits2)
}

func commonLogic(handleDigits func(string) [2]int) (int, error) {
	file, err := os.Open("day01/input1")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	sum := 0
	for scanner.Scan() {
		number, err := calibrationValue(scanner.Text(), handleDigits)
		if err != nil {
			return 0, err
		}
		sum += number
	}

	return sum, nil
}

func digits1(s string) [2]int {
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
	return digits
}

func digits2(s string) [2]int {
	digits := [2]int{0, 0}
	keys := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	table := make(map[string]int, len(keys))
	numbers := trie.New[int]()

	for i, k := range keys {
		numbers.Add(k, i+1)
		table[k] = i + 1
	}

	var current string
	for _, ch := range s {
		var value int

		if unicode.IsDigit(ch) {
			value = int(ch) - '0'
			current = ""
		} else {
			current += string(ch)

			keys, exists := numbers.Suggest(current)

			if !exists {
				current = current[1:]
				continue
			}

			if len(keys) != 1 {
				continue
			}

			if v, ok := table[current]; ok {
				value = v
				current = string(ch)
			} else {
				continue
			}
		}

		if digits[0] == 0 {
			digits[0] = value
		} else {
			digits[1] = value
		}
	}

	return digits
}

func calibrationValue(s string, handleDigits func(string) [2]int) (int, error) {
	digits := handleDigits(s)

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
