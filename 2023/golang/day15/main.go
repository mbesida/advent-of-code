package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

type Lens struct {
	label string
	focal int
}

var boxes map[int][]Lens = make(map[int][]Lens)

func main() {
	f := common.InputFileHandle("day15")
	defer f.Close()

	bytes, _ := io.ReadAll(f)
	line := string(bytes)

	t1 := func() int {
		sum := 0
		for _, step := range steps(line) {
			sum += holidayHash(step)
		}
		return sum
	}

	t2 := func() int {
		for _, step := range steps(line) {
			task2(step)
		}
		sum := 0
		for k, lenses := range boxes {
			sum += calc(k, lenses)
		}

		return sum
	}

	fmt.Println(common.HandleTasks(t1, t2))
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

func task2(step string) {
	if strings.ContainsRune(step, '-') {
		data := strings.Split(step, "-")
		lensLabel := data[0]
		hash := holidayHash(lensLabel)
		box := boxes[hash]
		for i, lense := range box {
			if lense.label == lensLabel {
				box = append(box[:i], box[i+1:]...)
				boxes[hash] = box
				return
			}
		}
	} else if strings.Contains(step, "=") {
		data := strings.Split(step, "=")
		focal, _ := strconv.Atoi(data[1])
		lens := Lens{data[0], focal}
		hash := holidayHash(lens.label)
		box := boxes[hash]
		for i, l := range box {
			if l.label == lens.label {
				box[i] = lens
				return
			}
		}
		boxes[hash] = append(box, lens)
	}
}

func calc(boxHash int, lenses []Lens) int {
	res := 0
	for i, l := range lenses {
		res += (boxHash + 1) * (i + 1) * l.focal
	}
	return res
}
