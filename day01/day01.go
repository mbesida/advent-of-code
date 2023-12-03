package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode"

	"github.com/s0rg/trie"
)

func main() {
	var choice string

	if len(os.Args) >= 2 {
		choice = os.Args[1]
	}

	var result int
	var err error

	if choice == "2" {
		fmt.Println("solving 2nd task of day 1")
		result, err = commonLogic(digits2)
	} else {
		fmt.Println("solving 1st task of day 1")
		result, err = commonLogic(digits1)
	}

	if err != nil {
		log.Fatalf("somewthing bad happened: %v\n", err)
	}

	fmt.Println(result)
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

	var result int
	var err error

	if digits == [2]int{0, 0} {
		err = fmt.Errorf("no digit found on line %s", s)
	} else if digits[1] == 0 {
		result = digits[0]*10 + digits[0]
	} else {
		result = digits[0]*10 + digits[1]
	}

	return result, err
}
