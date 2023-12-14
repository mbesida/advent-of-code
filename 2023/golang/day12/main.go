package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

type RowData struct {
	numbers []int
	pattern string
}

func main() {
	f := common.InputFileHandle("day12")
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var data []RowData
	for scanner.Scan() {
		data = append(data, parseLine(scanner.Text()))
	}

	sum := 0
	for _, rd := range data {
		sum += arrangements(rd)
	}

	fmt.Println(sum)
}

func parseLine(line string) RowData {
	splitted := strings.Fields(line)
	if len(splitted) != 2 {
		log.Fatalf("incorrect file format on line %s", line)
	}
	numbersData := strings.Split(splitted[1], ",")
	numbers := make([]int, len(numbersData))
	for i, v := range numbersData {
		numbers[i], _ = strconv.Atoi(v)
	}
	return RowData{numbers, splitted[0]}
}

func isValid(rd RowData) bool {
	partitions := strings.FieldsFunc(rd.pattern, func(r rune) bool { return r == '.' })

	if len(partitions) != len(rd.numbers) {
		return false
	}

	for i, p := range partitions {
		if len(p) != rd.numbers[i] {
			return false
		}
	}

	return true
}

func generateCases(toFill int, questionIndexes []int, pattern string) []string {
	var res []string
	numbers := generateNumbers(len(questionIndexes), toFill)
	countOfQuestionMarks := len(questionIndexes)
	for _, n := range numbers {
		binaryRepr := bitMap(countOfQuestionMarks, n)
		patternBytes := []byte(pattern)
		for i := 0; i < len(binaryRepr); i++ {
			if binaryRepr[i] == '1' {
				patternBytes[questionIndexes[i]] = '#'
			}
		}
		res = append(res, strings.ReplaceAll(string(patternBytes), "?", "."))
	}

	return res
}

func bitMap(questionsNumber, n int) string {
	pattern := fmt.Sprintf("%%0%db", questionsNumber)
	return fmt.Sprintf(pattern, n)
}

func countOnes(n int) int {
	count := 0
	for n > 0 {
		n &= (n - 1)
		count++
	}
	return count
}

func generateNumbers(n int, k int) []int {
	var res []int
	for i := 1; i < (1 << n); i++ {
		if countOnes(i) == k {
			res = append(res, i)
		}
	}
	return res
}

func questionIndexes(s string) []int {
	var res []int
	for i := 0; i < len(s); i++ {
		if s[i] == '?' {
			res = append(res, i)
		}
	}
	return res
}

func arrangements(rd RowData) int {
	numbersSum := 0
	for _, v := range rd.numbers {
		numbersSum += v
	}

	toFill := numbersSum - strings.Count(rd.pattern, "#")
	if toFill == 0 {
		return 1
	}

	cases := generateCases(toFill, questionIndexes(rd.pattern), rd.pattern)

	totalValid := 0
	for _, c := range cases {
		if isValid(RowData{rd.numbers, c}) {
			totalValid++
		}
	}

	return totalValid
}
