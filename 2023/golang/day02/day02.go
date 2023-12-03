package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type GameData struct {
	GameId  int
	IsValid bool
}

const (
	REDS   = 12
	GREENS = 13
	BlUES  = 14
)

func main() {
	file, err := os.Open("day02/input")
	if err != nil {
		log.Fatalf("no input file available")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	sum := 0
	for scanner.Scan() {
		row := scanner.Text()
		data, err := parseRow(row)
		if err != nil {
			log.Fatalln(err)
		}
		if data.IsValid {
			sum += data.GameId
		}
	}

	fmt.Println(sum)
}

func parseRow(row string) (*GameData, error) {
	gameSplitted := strings.Split(row, ":")
	if len(gameSplitted) != 2 {
		return nil, invalidRawFormat(row)
	}

	re := regexp.MustCompile(`Game (\d+)`)
	gameIdStr := re.FindStringSubmatch(gameSplitted[0])[1]
	gameId, err := strconv.Atoi(gameIdStr)
	if err != nil {
		return nil, err
	}

	attemptsStr := gameSplitted[1]
	attempts := strings.Split(attemptsStr, ";")

	for _, attempt := range attempts {
		parts := strings.Split(attempt, ",")

		for _, part := range parts {
			rawData := strings.Split(strings.TrimSpace(part), " ")
			if len(rawData) != 2 {
				return nil, invalidRawFormat(row)
			}
			value, err := strconv.Atoi(rawData[0])
			if err != nil {
				return nil, invalidRawFormat(row)
			}
			var limit int

			switch rawData[1] {
			case "red":
				limit = REDS
			case "green":
				limit = GREENS
			case "blue":
				limit = BlUES
			default:
				return nil, invalidRawFormat(row)
			}

			if value > limit {
				return &GameData{gameId, false}, nil
			}
		}
	}

	return &GameData{gameId, true}, nil

}

func invalidRawFormat(line string) error {
	return fmt.Errorf("invalid row format %s", line)
}
