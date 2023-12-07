package main

import (
	"bufio"
	"cmp"
	"fmt"
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
	'J': common.HandleTasks(func() int { return 10 }, func() int { return 0 }),
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
	label       string
	tpe         HandType
	possibleTpe HandType
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

	frequencies := sortedFrequencies(letters)

	tpe := handType(frequencies)

	t1 := func() Hand {
		return Hand{label, tpe, -1}
	}

	t2 := func() Hand {
		if strings.Contains(label, "J") {
			jCount := letters['J']
			topNonJ := 'A'
			value := 1
			for k, v := range letters {
				if k != 'J' && v >= value {
					topNonJ = k
					value = v
				}
			}
			delete(letters, 'J')
			letters[topNonJ] += jCount
			possibleFrequencies := sortedFrequencies(letters)
			possibleTpe := handType(possibleFrequencies)
			return Hand{label, tpe, possibleTpe}
		}
		return Hand{label, tpe, tpe}
	}

	return common.HandleTasks(t1, t2)
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
	t1 := func() int {
		return cmp.Compare(h1.hand.tpe, h2.hand.tpe)
	}

	t2 := func() int {
		return cmp.Compare(h1.hand.possibleTpe, h2.hand.possibleTpe)
	}

	cmpRes := common.HandleTasks(t1, t2)
	if cmpRes != 0 {
		return cmpRes
	}

	labelFunc := func(h Hand) func() string { return func() string { return h.label } }

	letters1 := strings.Split(labelFunc(h1.hand)(), "")
	letters2 := strings.Split(labelFunc(h2.hand)(), "")
	for i := 0; i < len(letters1); i++ {
		c := compareLetters(rune(letters1[i][0]), rune(letters2[i][0]))
		if c != 0 {
			return c
		}
	}
	return 0
}

func compareLetters(l1, l2 rune) int {
	t1 := func() map[rune]int {
		return LettersOrder
	}

	t2 := func() map[rune]int {
		return LettersOrder
	}

	mapping := common.HandleTasks(t1, t2)

	if mapping[l1] < mapping[l2] {
		return -1
	}
	if mapping[l1] > mapping[l2] {
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

func isFiveKind(sortedFrequencies []int) bool {
	return len(sortedFrequencies) == 1
}
func isFourKind(sortedFrequencies []int) bool {
	return slices.Equal(sortedFrequencies, []int{1, 4})
}
func isFullHouse(sortedFrequencies []int) bool {
	return slices.Equal(sortedFrequencies, []int{2, 3})
}
func isThreeKind(sortedFrequencies []int) bool {
	return slices.Equal(sortedFrequencies, []int{1, 1, 3})
}
func isTwoPair(sortedFrequencies []int) bool {
	return slices.Equal(sortedFrequencies, []int{1, 2, 2})
}
func isOnePair(sortedFrequencies []int) bool {
	return slices.Equal(sortedFrequencies, []int{1, 1, 1, 2})
}

func handType(sortedFrequencies []int) HandType {
	var tpe HandType

	switch {
	case isFiveKind(sortedFrequencies):
		tpe = FiveKind
	case isFourKind(sortedFrequencies):
		tpe = FourKind
	case isFullHouse(sortedFrequencies):
		tpe = FullHouse
	case isThreeKind(sortedFrequencies):
		tpe = ThreeKind
	case isTwoPair(sortedFrequencies):
		tpe = TwoPair
	case isOnePair(sortedFrequencies):
		tpe = OnePair
	default:
		tpe = HighCard
	}
	return tpe
}

func sortedFrequencies(letters map[rune]int) []int {
	frequencies := make([]int, 0, len(letters))
	for _, v := range letters {
		frequencies = append(frequencies, v)
	}
	slices.Sort(frequencies)
	return frequencies
}
