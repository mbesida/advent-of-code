package main

import (
	"fmt"
	"log"

	"github.com/mbesida/advent-of-code-2023/day01"
)

func main() {
	res, err := day01.Solve1()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
