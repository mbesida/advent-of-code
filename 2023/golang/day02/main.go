package main

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

const (
	REDS   = 12
	GREENS = 13
	BlUES  = 14
)

type GameData struct {
	GameId    int
	IsValid   bool
	MaxReds   int
	MaxGreens int
	MaxBlues  int
}

func NewGameData(id int, isValid bool) *GameData {
	return &GameData{id, isValid, 0, 0, 0}
}

func main() {
	file := common.InputFileHandle("day02")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	sum := 0
	for scanner.Scan() {
		row := scanner.Text()
		data, err := parseRow(row)
		if err != nil {
			log.Fatalln(err)
		}

		t1 := func() int {
			return task1(data)
		}
		t2 := func() int {
			return task2(data)
		}

		sum += common.HandleTasks(t1, t2)
	}

	fmt.Println(sum)
}

func task1(data *GameData) int {
	result := 0

	if data.IsValid {
		result += data.GameId
	}

	return result
}

func task2(data *GameData) int {
	return data.MaxBlues * data.MaxGreens * data.MaxReds
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

	gameData := NewGameData(gameId, true)

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

			switch rawData[1] {
			case "red":
				gameData.MaxReds = max(gameData.MaxReds, value)
			case "green":
				gameData.MaxGreens = max(gameData.MaxGreens, value)
			case "blue":
				gameData.MaxBlues = max(gameData.MaxBlues, value)
			default:
				return nil, invalidRawFormat(row)
			}
		}
	}

	if gameData.MaxBlues > BlUES || gameData.MaxGreens > GREENS || gameData.MaxReds > REDS {
		gameData.IsValid = false
	}

	return gameData, nil

}

func invalidRawFormat(line string) error {
	return fmt.Errorf("invalid row format %s", line)
}
