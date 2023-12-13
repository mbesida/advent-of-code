package main

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

type Image [][]rune

var image Image

var galaxies map[int]common.Point = make(map[int]common.Point)

func main() {
	f := common.InputFileHandle("day11")
	defer f.Close()
	bytes, _ := io.ReadAll(f)
	rawImage := string(bytes)

	n := common.HandleValue(1, 1_000_000)
	fmt.Println(n)

	parseImage(rawImage)
	expandImage(1)
	assignUniqueNumbers()
	first := calculateSum()

	parseImage(rawImage)
	expandImage(2)
	assignUniqueNumbers()
	// second := calculateSum()
	// diff := second - first

	var total uint64 = uint64(first)
	// for i := 1; i < n; i++ {
	// 	fmt.Println(total)
	// 	total = total + uint64(diff)
	// }

	fmt.Println(total)

}

func parseImage(rawImage string) {
	image = nil
	splitted := strings.Split(rawImage, "\n")
	for _, s := range splitted {
		row := make([]rune, len(s))
		for j, r := range s {
			row[j] = r
		}
		image = append(image, row)
	}
}

func expandImage(n int) {
	var expandedImage Image
	magnitude := n - 1

	for i := 0; i < len(image); i++ {
		expandedImage = append(expandedImage, image[i])
		if isBlank(image[i]) {
			for k := 0; k < magnitude; k++ {
				expandedImage = append(expandedImage, image[i])
			}
		}
	}

	for i := 0; i < len(expandedImage[0]); i++ {
		column := make([]rune, len(expandedImage))
		for j, row := range expandedImage {
			column[j] = row[i]
		}
		if isBlank(column) {
			expandColumnImage := make([][]rune, len(expandedImage))
			for i := range expandColumnImage {
				expandColumnImage[i] = make([]rune, len(expandedImage[0])+magnitude)
			}
			for k, row := range expandedImage {
				for j, r := range row {
					if j < i {
						expandColumnImage[k][j] = r
					} else {
						expandColumnImage[k][j+magnitude] = r
					}
				}
				for x := 0; x < magnitude; x++ {
					expandColumnImage[k][i+x] = column[k]
				}
			}
			expandedImage = expandColumnImage
			i += magnitude
		}
	}

	image = expandedImage
}

func assignUniqueNumbers() {
	k := 1
	for i, row := range image {
		for j, e := range row {
			if e == '#' {
				galaxies[k] = common.Point{I: i, J: j}
				k++
			}
		}
	}
}

func isBlank(row []rune) bool {
	return !strings.Contains(string(row), "#")
}

func distance(p1, p2 common.Point) int {
	return int(math.Abs(float64(p1.I)-float64(p2.I)) + math.Abs(float64(p1.J)-float64(p2.J)))
}

func calculateSum() int {
	keys := make([]int, 0, len(galaxies))
	for i := range galaxies {
		keys = append(keys, i)
	}

	sum := 0
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			sum += distance(galaxies[keys[i]], galaxies[keys[j]])
		}
	}
	return sum
}
