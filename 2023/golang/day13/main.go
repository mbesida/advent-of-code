package main

import (
	"fmt"
	"io"
	"math/bits"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

func main() {
	f := common.InputFileHandle("day13")
	defer f.Close()

	bytes, _ := io.ReadAll(f)
	stringifiedData := string(bytes)

	notes := parseNotes(stringifiedData)

	sum := 0
	for _, n := range notes {
		x := calcReflection1(n)
		// fmt.Println(x)
		sum += x
	}
	fmt.Println(sum)
}

func parseNotes(data string) []Note {
	groups := strings.Split(data, "\n\n")
	n := len(groups)

	rawData := make([][][]rune, n)
	for i := 0; i < n; i++ {
		rows := strings.Split(groups[i], "\n")
		rawData[i] = makeNoteRunes(rows)
	}

	notes := make([]Note, len(rawData))
	for i, matrix := range rawData {
		notes[i] = NewNote(matrix)
	}

	return notes
}

func makeNoteRunes(data []string) [][]rune {
	result := make([][]rune, len(data))
	for i, line := range data {
		runes := []rune(line)
		for j, r := range runes {
			if r == '#' {
				runes[j] = '1'
			} else {
				runes[j] = '0'
			}
		}
		result[i] = runes
	}
	return result
}

func calcReflection1(note Note) int {
	rowReflection, _ := findReflection(note.rows)
	columnReflection, _ := findReflection(note.columns)

	if rowReflection != nil {
		return rowReflection.calc(true)
	}
	if columnReflection != nil {
		return columnReflection.calc(false)
	}
	return 0
}

func calcReflection2(note Note) int {
	rowReflection, potenatialRows := findReflection(note.rows)
	columnReflection, potenatialColumns := findReflection(note.columns)

	if rowReflection != nil {
		// 1
		if len(potenatialRows) > 1 {
			for _, index := range potenatialRows {
				if rowReflection.index != index {

					found := true
					behind, forward := index-1, index+2
					for behind >= 0 && forward < len(note.rows) {
						if note.rows[behind] != note.rows[forward] {
							xor := note.rows[behind] ^ note.rows[forward]
							if bits.OnesCount64(uint64(xor)) != 1 {
								found = false
								break
							}
						}
						behind--
						forward++
					}
					if found {

					}
				}
			}
		}
		return rowReflection.calc(true)
	}
	if columnReflection != nil {
		return columnReflection.calc(false)
	}
	return 0
}

func findReflection(data []int64) (*Reflection, []int) {
	positions := potentialReflectionPoints(data)

	for _, i := range positions {
		if checkPosition(i, data) {
			return NewReflection(i, len(data)), positions
		}
	}

	return nil, positions
}

func checkPosition(i int, data []int64) bool {
	behind, forward := i-1, i+2
	for behind >= 0 && forward < len(data) {
		if data[behind] != data[forward] {
			return false
		}
		behind--
		forward++
	}
	return true
}
