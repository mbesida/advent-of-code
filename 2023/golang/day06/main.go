package main

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

type Race struct {
	t        int
	distance int
}

func main() {
	f := common.InputFileHandle("day06")
	defer f.Close()
	dataBytes, _ := io.ReadAll(f)
	dataString := string(dataBytes)

	races, err := parseRaces(dataString)
	if err != nil {
		log.Fatalf("error happened %s", err)
	}

	res := calculate(races)

	fmt.Println(res)

}

func calculate(races []Race) int {
	acc := 1
	for _, r := range races {
		counter := 0
		for i := 1; i < r.t; i++ {
			if i*(r.t-i) > r.distance {
				counter++
			}
		}
		acc *= counter

	}
	return acc
}

func parseRaces(s string) ([]Race, error) {
	e := fmt.Errorf("incorrect file format")

	lines := strings.Split(s, "\n")
	if len(lines) != 2 {
		return nil, e
	}
	timeData, foundT := strings.CutPrefix(lines[0], "Time:")
	distanceData, foundD := strings.CutPrefix(lines[1], "Distance:")

	if !foundT || !foundD {
		return nil, e
	}

	timeSlice := strings.Fields(timeData)
	distanceSlice := strings.Fields(distanceData)

	if len(timeSlice) != len(distanceSlice) {
		return nil, e
	}

	numberOfRaces := len(timeSlice)
	result := make([]Race, numberOfRaces)
	for i := 0; i < numberOfRaces; i++ {
		t, _ := strconv.Atoi(timeSlice[i])
		d, _ := strconv.Atoi(distanceSlice[i])
		result[i] = Race{t, d}
	}
	return result, nil
}
