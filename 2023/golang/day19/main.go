package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

type Part struct {
	x, m, a, s int
}

type Workflow struct {
	name  string
	rules []func(p Part) string
}

func main() {
	f := common.InputFileHandle("day19")
	defer f.Close()
	bytes, _ := io.ReadAll(f)
	stringData := string(bytes)
	splitted := strings.Split(stringData, "\n\n")
	workflowData, partData := splitted[0], splitted[1]
	workflows := parseWorkFlows(workflowData)
	parts := parseParts(partData)

	for p := range parts {
		currentW := workflows["in"]
	outer:
		for {
			for _, rule := range currentW.rules {
				next := rule(p)
				if next == "" {
					continue
				}
				if next == "A" {
					parts[p] = true
					break outer
				}
				if next != "R" {
					currentW = workflows[next]
					break
				}
				break outer
			}
		}
	}

	sum := 0
	for p, v := range parts {
		if v {
			sum += p.x + p.m + p.a + p.s
		}
	}
	fmt.Println(sum)

}

func parseWorkFlows(data string) map[string]Workflow {
	res := make(map[string]Workflow)
	for _, line := range strings.Split(data, "\n") {
		i := strings.Index(line, "{")
		name := line[:i]
		w := Workflow{name, nil}
		rulesString := line[i+1 : len(line)-1]
		rules := strings.Split(rulesString, ",")
		for _, rule := range rules {
			temp := strings.Split(rule, ":")
			if len(temp) == 1 {
				w.rules = append(w.rules, func(_ Part) string {
					return temp[0]
				})
				continue
			}
			condition, nextName := temp[0], temp[1]
			w.rules = append(w.rules, func(p Part) string {
				e := string(condition[0])
				condString := string(condition[1])
				valueStr := condition[2:]
				value, _ := strconv.Atoi(valueStr)
				result := ""
				switch [2]string{e, condString} {
				case [2]string{"x", "<"}:
					if p.x < value {
						result = nextName
					}
				case [2]string{"x", ">"}:
					if p.x > value {
						result = nextName
					}
				case [2]string{"m", "<"}:
					if p.m < value {
						result = nextName
					}
				case [2]string{"m", ">"}:
					if p.m > value {
						result = nextName
					}
				case [2]string{"a", "<"}:
					if p.a < value {
						result = nextName
					}
				case [2]string{"a", ">"}:
					if p.a > value {
						result = nextName
					}
				case [2]string{"s", "<"}:
					if p.s < value {
						result = nextName
					}
				case [2]string{"s", ">"}:
					if p.s > value {
						result = nextName
					}
				}
				return result
			})
		}
		res[name] = w

	}
	return res
}

func parseParts(data string) map[Part]bool {
	res := make(map[Part]bool)
	for _, line := range strings.Split(data, "\n") {
		elements := strings.Split(line[1:len(line)-1], ",")
		p := Part{}
		for _, v := range elements {
			partData := strings.Split(v, "=")
			value, _ := strconv.Atoi(partData[1])
			switch partData[0] {
			case "x":
				p.x = value
			case "m":
				p.m = value
			case "a":
				p.a = value
			case "s":
				p.s = value
			}
		}
		res[p] = false
	}
	return res
}
