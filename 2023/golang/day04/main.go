package main

import (
	"bufio"
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

const cardRegexp = `Card *(\d+): +((\d{1,2} *)+)\| +((\d{1,2} *)+)`

type Card struct {
	Id     int
	Points int
}

func main() {
	f := common.InputFileHandle("day04")
	defer f.Close()
	scanner := bufio.NewScanner(f)
	cards := make(map[int][]Card)

	for scanner.Scan() {
		card, _ := parseCard(scanner.Text())
		cards[card.Id] = append(cards[card.Id], *card)
	}

	totcalCards := len(cards)

	t1 := func() int {
		sum := 0
		for _, cardSlice := range cards {
			sum += cardSlice[0].Points
		}
		return sum
	}

	t2 := func() int {
		keys := make([]int, 0, len(cards))
		for k := range cards {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, id := range keys {
			for _, card := range cards[id] {
				points := card.Points
				for i := id + 1; i <= id+points; i++ {
					if i <= totcalCards {
						c := cards[i]
						cards[i] = append(cards[i], c[0])
					}
				}

			}
		}
		value := 0
		for _, cardSlice := range cards {
			value += len(cardSlice)
		}
		return value
	}
	fmt.Println(common.HandleTasks(t1, t2))
}

func parseCard(s string) (*Card, error) {
	re := regexp.MustCompile(cardRegexp)
	allMatches := re.FindStringSubmatch(s)

	if len(allMatches) != 6 {
		return nil, fmt.Errorf("invalid card format")
	}

	cardId, _ := strconv.Atoi(allMatches[1])
	winningNumberStrings := strings.Fields(allMatches[2])
	myNumberStrings := strings.Fields(allMatches[4])
	winningNumbers := make(map[int]bool)
	for _, n := range winningNumberStrings {
		value, _ := strconv.Atoi(n)
		winningNumbers[value] = true
	}
	myNumbers := make(map[int]bool)
	for _, n := range myNumberStrings {
		value, _ := strconv.Atoi(n)
		myNumbers[value] = true
	}

	points := 0
	for k := range myNumbers {
		if winningNumbers[k] {
			points += 1
		}
	}
	t1 := func() int {
		if points != 0 {
			return int(math.Pow(2, float64(points-1)))
		} else {
			return points
		}

	}
	t2 := func() int {
		return points
	}

	return &Card{cardId, common.HandleTasks(t1, t2)}, nil
}
