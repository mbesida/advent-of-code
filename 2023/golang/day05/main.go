package main

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

type Mapping struct {
	destination int64
	source      int64
	rangeLength int64
}

var (
	seeds                 []int
	seedToSoil            []Mapping
	soilToFertilizer      []Mapping
	fertilizerToWater     []Mapping
	waterToLight          []Mapping
	lightToTemperature    []Mapping
	temperatureToHumidity []Mapping
	humidityToLocation    []Mapping
)

func main() {
	f := common.InputFileHandle("day05")
	defer f.Close()
	scanner := bufio.NewScanner(f)
	parseData(scanner)
}

func parseData(scanner *bufio.Scanner) {
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "seeds:") {
			values := strings.Fields(line)[1:]
			for _, s := range values {
				v, _ := strconv.Atoi(s)
				seeds = append(seeds, v)
			}
		}
		if strings.HasPrefix(line, "seed-to-soil map:") {
			for scanner.Scan() {
				data := strings.Fields(scanner.Text())
				dest, _ := strconv.Atoi(data)
				seedToSoil = append(seedToSoil, Ma)
				for i, s := range data {
					strconv.Atoi(s)
					seedToSoil = append(seedToSoil, Ma)
				}
			}
		}
		if inSeedToSoil {

		}
	}
}
