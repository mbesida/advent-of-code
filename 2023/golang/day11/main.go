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
	parseImage(rawImage)

	expandImage()
	assignUniqueNumbers()

	res := calculateSum()
	fmt.Println(res)
}

func parseImage(rawImage string) {
	splitted := strings.Split(rawImage, "\n")
	for _, s := range splitted {
		row := make([]rune, len(s))
		for j, r := range s {
			row[j] = r
		}
		image = append(image, row)
	}
}

func expandImage() {
	var expandedImage Image

	for i := 0; i < len(image); i++ {
		expandedImage = append(expandedImage, image[i])
		if isBlank(image[i]) {
			expandedImage = append(expandedImage, image[i])
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
				expandColumnImage[i] = make([]rune, len(expandedImage[0])+1)
			}
			for k, row := range expandedImage {
				for j, r := range row {
					if j < i {
						expandColumnImage[k][j] = r
					} else {
						expandColumnImage[k][j+1] = r
					}
				}
				expandColumnImage[k][i] = column[k]
			}
			expandedImage = expandColumnImage
			i++
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
