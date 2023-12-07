package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

type HandType int

const (
	HighCard HandType = iota
	OnePair
	TwoPair
	ThreeKind
	FullHouse
	FourKind
	FiveKind
)

var LettersOrder = map[rune]int{
	'A': 13,
	'K': 12,
	'Q': 11,
	'J': 10,
	'T': 9,
	'9': 8,
	'8': 7,
	'7': 6,
	'6': 5,
	'5': 4,
	'4': 3,
	'3': 2,
	'2': 1,
}

type Hand struct {
	label string
	tpe   HandType
}

type Record struct {
	hand Hand
	bid  int
}

func NewHand(label string) Hand {
	letters := make(map[rune]int)

	for _, r := range label {
		letters[r]++
	}

	var tpe HandType
	switch {
	case isFiveKind(letters):
		tpe = FiveKind
	case isFourKind(letters):
		tpe = FourKind
	case isFullHouse(letters):
		tpe = FullHouse
	case isThreeKind(letters):
		tpe = ThreeKind
	case isTwoPair(letters):
		tpe = TwoPair
	case isOnePair(letters):
		tpe = OnePair
	default:
		tpe = HighCard
	}

	return Hand{label, tpe}

}

func main() {
	f := common.InputFileHandle("day07")
	defer f.Close()

	records := parseRecords(f)

	slices.SortFunc(records, compareHands)

	sum := 0
	for i, r := range records {
		sum += (i + 1) * r.bid
	}

	fmt.Println(sum)
}

func compareHands(h1, h2 Record) int {
	if h1.hand.tpe < h2.hand.tpe {
		return -1
	}
	if h1.hand.tpe > h2.hand.tpe {
		return 1
	}
	letters1 := strings.Split(h1.hand.label, "")
	letters2 := strings.Split(h2.hand.label, "")
	for i := 0; i < len(letters1); i++ {
		c := compareLetters(rune(letters1[i][0]), rune(letters2[i][0]))
		if c != 0 {
			return c
		}
	}
	return 0
}

func compareLetters(l1, l2 rune) int {
	if LettersOrder[l1] < LettersOrder[l2] {
		return -1
	}
	if LettersOrder[l1] > LettersOrder[l2] {
		return 1
	}
	return 0
}

func parseRecords(f *os.File) []Record {
	scanner := bufio.NewScanner(f)
	var records []Record
	for scanner.Scan() {
		data := strings.Fields(scanner.Text())
		bid, _ := strconv.Atoi(data[1])
		records = append(records, Record{NewHand(data[0]), bid})
	}
	return records
}

func isFiveKind(letters map[rune]int) bool {
	return len(letters) == 1
}
func isFourKind(letters map[rune]int) bool {
	keys := make([]rune, 0, len(letters))
	for k := range letters {
		keys = append(keys, k)
	}

	return len(letters) == 2 && math.Abs(float64(letters[keys[0]])-float64(letters[keys[1]])) == 3
}
func isFullHouse(letters map[rune]int) bool {
	keys := make([]rune, 0, len(letters))
	for k := range letters {
		keys = append(keys, k)
	}
	return len(letters) == 2 && math.Abs(float64(letters[keys[0]])-float64(letters[keys[1]])) == 1
}
func isThreeKind(letters map[rune]int) bool {
	size := len(letters)
	for _, v := range letters {
		if v == 3 {
			return size == 3
		}
	}
	return false
}
func isTwoPair(letters map[rune]int) bool {
	values := make([]int, 0, len(letters))
	for _, v := range letters {
		values = append(values, v)
	}

	slices.Sort(values)

	return slices.Equal(values, []int{1, 2, 2})
}
func isOnePair(letters map[rune]int) bool {
	values := make([]int, 0, len(letters))
	for _, v := range letters {
		values = append(values, v)
	}

	slices.Sort(values)

	return slices.Equal(values, []int{1, 1, 1, 2})
}
