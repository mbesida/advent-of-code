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

	res := travel(network)
	fmt.Println(res)
}

func travel(network Map) int {
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
