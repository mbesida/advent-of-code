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

	sum := 0
	for scanner.Scan() {
		rd := parseLine(scanner.Text())
		sum += arrangements(rd)
	}

	fmt.Println(sum)
}

func parseLine(line string) RowData {
	splitted := strings.Fields(line)
	if len(splitted) != 2 {
		log.Fatalf("incorrect file format on line %s", line)
	}
	n := common.HandleValue(1, 5)
	var numbers []string
	for i := 1; i <= n; i++ {
		numbers = append(numbers, splitted[1])
	}

	var patterns []string
	for i := 1; i <= n; i++ {
		patterns = append(patterns, splitted[0])
	}
	numbersData := strings.Split(strings.Join(numbers, ","), ",")
	ints := make([]int, len(numbersData))
	for i, v := range numbersData {
		ints[i], _ = strconv.Atoi(v)
	}
	return RowData{ints, strings.Join(patterns, "?")}
}

type Key struct {
	a, b, c int
	d       bool
}

var memo map[Key]int = make(map[Key]int)

func arrangements(rd RowData) int {
	count := ways(rd.pattern, rd.numbers, 0, false)
	clear(memo)
	return count
}

func ways(s string, numbers []int, current int, consumedHash bool) int {
	key := Key{len(s), len(numbers), current, consumedHash}
	if v, ok := memo[key]; ok {
		return v
	}

	if len(s) == 0 {
		if len(numbers) == 0 && current == 0 {
			return 1
		}
		return 0
	}

	switch s[0] {
	case '.':
		if current != 0 {
			return 0
		}
		return ways(s[1:], numbers, current, false)

	case '#':
		switch {
		case len(numbers) == 0 && current == 0:
			return 0

		case len(numbers) == 0:
			return ways(s[1:], numbers, current-1, true)

		case consumedHash && current == 0:
			return 0

		case current == 0:
			current = numbers[0]
			numbers = numbers[1:]
		}
		return ways(s[1:], numbers, current-1, true)

	default:
		count := ways("."+s[1:], numbers, current, consumedHash)
		count += ways("#"+s[1:], numbers, current, consumedHash)
		memo[key] = count
		return count
	}
}
