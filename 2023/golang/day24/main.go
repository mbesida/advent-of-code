package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

const (
	RangeMin = 200000000000000
	RangeMax = 400000000000000
)

type Hailstone struct {
	x, y, z    int64
	dx, dy, dz int64
}

// x1 + t1*v1x = x2 + t2*v2x
// y1 + t1*v1y = y2 + t2*v2y
// solving equations for t1 and t2
// if t1 and t2 are within [0, infinity) then rays intersect
func (h Hailstone) intersectsWithinRange(another Hailstone) bool {
	denominator := float64(h.dx)*float64(another.dy) - float64(h.dy)*float64(another.dx)

	if denominator == 0 { // rays are parallel
		return false
	}

	t1 := (float64(another.dy)*float64(another.x-h.x) - float64(another.dx)*float64(another.y-h.y)) / denominator
	t2 := (float64(h.dy)*float64(another.x-h.x) - float64(h.dx)*float64(another.y-h.y)) / denominator

	if t1 >= 0 && t2 >= 0 {
		xCoord := float64(h.x) + t1*float64(h.dx)
		yCoord := float64(h.y) + t1*float64(h.dy)
		return xCoord >= RangeMin && xCoord <= RangeMax && yCoord >= RangeMin && yCoord <= RangeMax
	}

	return false
}

func main() {
	f := common.InputFileHandle("day24")
	defer f.Close()
	hailstones := parseHailstones(f)

	counter := 0
	for i := 0; i < len(hailstones)-1; i++ {
		for j := i + 1; j < len(hailstones); j++ {
			if hailstones[i].intersectsWithinRange(hailstones[j]) {
				counter++
			}
		}
	}
	fmt.Println(counter)
}

func parseHailstones(reader io.Reader) []Hailstone {
	var res []Hailstone
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), "@")
		positionData := strings.Split(data[0], ",")
		velocityData := strings.Split(data[1], ",")
		res = append(res, Hailstone{
			x:  parseInteger(positionData[0]),
			y:  parseInteger(positionData[1]),
			z:  parseInteger(positionData[2]),
			dx: parseInteger(velocityData[0]),
			dy: parseInteger(velocityData[1]),
			dz: parseInteger(velocityData[2]),
		})

	}
	return res
}

func parseInteger(s string) int64 {
	res, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	return res
}
