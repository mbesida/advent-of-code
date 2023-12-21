package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

func main() {
	f := common.InputFileHandle("day15")
	defer f.Close()

	bytes, _ := io.ReadAll(f)
	line := string(bytes)

	sum := 0
	for _, step := range steps(line) {
		sum += holidayHash(step)
	}
	fmt.Println(sum)
}

func steps(line string) []string {
	return strings.Split(strings.TrimSpace(line), ",")
}

func holidayHash(str string) int {
	current := 0
	for _, r := range str {
		current += int(r)
		current *= 17
		current %= 256
	}
	return current
}
