package main

import (
	"bufio"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
	"golang.org/x/exp/maps"
)

type Range struct {
	start int64
	end   int64
}

func NewRange(start, length int64) Range {
	return Range{start, start + length}
}

func (r Range) inRange(n int64) bool {
	return n >= r.start && n < r.end
}
func (r Range) insideRange(other Range) bool {
	return other.start <= r.start && other.end >= r.end
}
func (r Range) interleaveLeft(other Range) bool {
	return r.start < other.start && other.start <= r.end && r.end < other.end
}

type Mapping struct {
	destination int64
	source      int64
	rangeLength int64
}

var (
	seeds                 []int64
	seedRanges            []Range
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

	dictionaries := []map[Range]Range{
		makeRanges(seedToSoil),
		makeRanges(soilToFertilizer),
		makeRanges(fertilizerToWater),
		makeRanges(waterToLight),
		makeRanges(lightToTemperature),
		makeRanges(temperatureToHumidity),
		makeRanges(humidityToLocation),
	}

	t1 := func() int64 {
		locations := make([]int64, len(seeds))
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

	t2 := func() int64 {
		currentRanges := seedRanges
		for _, dict := range dictionaries {
			currentRanges = trasnformRanges(currentRanges, dict)
			for i, r := range currentRanges {
				for k, v := range dict {
					if r.insideRange(k) {
						x := r.start - k.start
						length := r.end - r.start
						currentRanges[i] = NewRange(v.start+x, length)
						break
					}
				}
			}
		}
		minLocation := currentRanges[0].start
		for _, cr := range currentRanges {
			if cr.start < minLocation {
				minLocation = cr.start
			}
		}
		return minLocation
	}

	return common.HandleTasks(t1, t2)

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
			if len(seeds)%2 != 0 {
				log.Fatal("seeds should be in ranges")
			}
			for i := 0; i < len(seeds); i += 2 {
				seedRanges = append(seedRanges, NewRange(seeds[i], seeds[i+1]))
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
		k := NewRange(m.source, m.rangeLength)
		result[k] = NewRange(m.destination, m.rangeLength)
	}
	return result
}

func trasnformRanges(ranges []Range, dict map[Range]Range) []Range {
	var newRanges []Range

	keys := maps.Keys(dict)
	slices.SortFunc(keys, func(k1, k2 Range) int {
		if k1.start < k2.start {
			return -1
		}
		return 0
	})

	for _, r := range ranges {
		var tempRanges []Range

		for _, s := range keys {
			if r.interleaveLeft(s) {
				tempRanges = append(tempRanges, Range{r.start, s.start})
				tempRanges = append(tempRanges, Range{s.start, r.end})
			} else if s.interleaveLeft(r) {
				tempRanges = append(tempRanges, Range{r.start, s.end})
				tempRanges = append(tempRanges, Range{s.end, r.end})
			} else if r.insideRange(s) {
				tempRanges = append(tempRanges, r)
			} else if s.insideRange(r) {
				if r.start != s.start {
					tempRanges = append(tempRanges, Range{r.start, s.start})
					tempRanges = append(tempRanges, Range{s.start, s.end})
					if r.end != s.end {
						tempRanges = append(tempRanges, Range{s.end, r.end})
					}
				} else {
					tempRanges = append(tempRanges, Range{r.start, s.end})
					tempRanges = append(tempRanges, Range{s.end, r.end})
				}
			}
			if len(tempRanges) != 0 {
				break
			}

		}
		if len(tempRanges) == 0 {
			newRanges = append(newRanges, r)
		} else {
			newRanges = append(newRanges, tempRanges...)
		}
	}
	return newRanges
}
