package main

import (
	"bufio"
	"fmt"
	"math"
	"regexp"
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
	sum := 0
	for scanner.Scan() {
		card, _ := parseCard(scanner.Text())
		sum += card.Points
	}
	fmt.Println(sum)
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
	if points != 0 {
		points = int(math.Pow(2, float64(points-1)))
	}

	return &Card{cardId, points}, nil
}
