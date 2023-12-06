package main

import (
	"bufio"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

type Range struct {
	start int64
	end   int64
}

func (r Range) inRange(n int64) bool {
	return n >= r.start && n < r.end
}

type Mapping struct {
	destination int64
	source      int64
	rangeLength int64
}

var (
	seeds                 []int64
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

	result := calculate()

	fmt.Println(result)
}

func calculate() int64 {
	locations := make([]int64, len(seeds))

	dictionaries := []map[Range]Range{
		makeRanges(seedToSoil),
		makeRanges(soilToFertilizer),
		makeRanges(fertilizerToWater),
		makeRanges(waterToLight),
		makeRanges(lightToTemperature),
		makeRanges(temperatureToHumidity),
		makeRanges(humidityToLocation),
	}

	for i, seed := range seeds {
		path := make([]int64, len(dictionaries))
		for j, dict := range dictionaries {
			if j == 0 {
				path[j] = findPathInMapping(seed, dict)
			} else {
				path[j] = findPathInMapping(path[j-1], dict)
			}
		}
		locations[i] = path[len(path)-1]
	}

	return slices.Min(locations)
}

func findPathInMapping(seed int64, mapping map[Range]Range) int64 {
	result := int64(-1)
	for k, v := range mapping {
		if k.inRange(seed) {
			length := seed - k.start
			result = v.start + length
			return result
		}
	}
	if result == -1 {
		result = seed
	}
	return result
}

func parseData(scanner *bufio.Scanner) {
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "seeds:"):
			values := strings.Fields(line)[1:]
			for _, s := range values {
				v, _ := strconv.ParseInt(s, 10, 64)
				seeds = append(seeds, v)
			}
			scanner.Scan()
			continue
		case strings.HasPrefix(line, "seed-to-soil"):
			seedToSoil = parseMap(scanner)
		case strings.HasPrefix(line, "soil-to-fertilizer"):
			soilToFertilizer = parseMap(scanner)
		case strings.HasPrefix(line, "fertilizer-to-water"):
			fertilizerToWater = parseMap(scanner)
		case strings.HasPrefix(line, "water-to-light"):
			waterToLight = parseMap(scanner)
		case strings.HasPrefix(line, "light-to-temperature"):
			lightToTemperature = parseMap(scanner)
		case strings.HasPrefix(line, "temperature-to-humidity"):
			temperatureToHumidity = parseMap(scanner)
		case strings.HasPrefix(line, "humidity-to-location"):
			humidityToLocation = parseMap(scanner)
		default:
			log.Fatalf("incorrect file format %s", line)
		}
	}
}

func parseMap(s *bufio.Scanner) []Mapping {
	var mapping []Mapping
	for s.Scan() {
		row := s.Text()
		if row == "" {
			break
		}
		data := strings.Fields(row)
		if len(data) != 3 {
			log.Fatal("incorrect format of a map")
		}
		dest, _ := strconv.ParseInt(data[0], 10, 64)
		source, _ := strconv.ParseInt(data[1], 10, 64)
		length, _ := strconv.ParseInt(data[2], 10, 64)
		mapping = append(mapping, Mapping{dest, source, length})
	}
	return mapping
}

func makeRanges(mapping []Mapping) map[Range]Range {
	result := make(map[Range]Range)
	for _, m := range mapping {
		k := Range{m.source, m.source + m.rangeLength}
		result[k] = Range{m.destination, m.destination + m.rangeLength}
	}
	return result
}
