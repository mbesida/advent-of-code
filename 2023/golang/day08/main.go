package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

var instructions string

type Direction struct {
	left  string
	right string
}

type Map = map[string]Direction

func main() {
	f := common.InputFileHandle("day08")
	defer f.Close()

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	instructions = scanner.Text()
	scanner.Scan()
	network := buildMap(scanner)

	res := common.HandleTasks(func() int { return travelFirst(network) }, func() int { return travelSecond(network) })
	fmt.Println(res)
}

func travelFirst(network Map) int {
	destination := "ZZZ"
	current := "AAA"
	counter := 0
	n := len(instructions)
	for {

		if current == destination {
			break
		}
		directions := network[current]
		chosenDirection := rune(instructions[counter%n])
		if chosenDirection == 'L' {
			current = directions.left
		} else {
			current = directions.right
		}
		counter++
	}
	return counter
}

func travelSecond(network Map) int {
	n := len(instructions)
	var current []string
	for k := range network {
		if strings.HasSuffix(k, "A") {
			current = append(current, k)
		}
	}
	stepsToEnd := make([]int, len(current))

	counter := 0
	for {
		if endCondition(stepsToEnd) {
			break
		}
		for i, v := range current {
			if stepsToEnd[i] != 0 {
				continue
			}
			directions := network[v]
			chosenDirection := rune(instructions[counter%n])
			if chosenDirection == 'L' {
				current[i] = directions.left
			} else {
				current[i] = directions.right
			}
			if strings.HasSuffix(current[i], "Z") {
				stepsToEnd[i] = counter + 1
			}
		}
		counter++
	}

	res := stepsToEnd[0]
	for i := 1; i < len(stepsToEnd); i++ {
		res = lcd(res, stepsToEnd[i])
	}
	return res
}

func buildMap(scanner *bufio.Scanner) Map {
	var result Map = make(map[string]Direction)
	for scanner.Scan() {
		line := scanner.Text()
		data := strings.Split(line, "=")
		if len(data) != 2 {
			log.Fatalf("incorrect line format %s", line)
		}
		locationId := strings.TrimSpace(data[0])
		leftRight := strings.Split(data[1], ",")
		left := strings.Trim(leftRight[0], " (")
		right := strings.Trim(leftRight[1], " )")
		result[locationId] = Direction{left, right}
	}
	return result
}

func endCondition(stepsToEnd []int) bool {
	for _, v := range stepsToEnd {
		if v == 0 {
			return false
		}
	}
	return true
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcd(a, b int) int {
	return (a * b) / gcd(a, b)
}
