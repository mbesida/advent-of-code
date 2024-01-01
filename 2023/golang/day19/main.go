package main

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

type Part struct {
	x, m, a, s int
}

type Rule struct {
	element, condition string
	value              int
}

type Condition struct {
	rule *Rule
	next string
}
type Workflow struct {
	name  string
	rules []func(p Part) string
	conds []Condition
}

type StateItem struct {
	lower, upper             int
	lowerStrict, upperStrict bool
}

func NewState() StateItem {
	return StateItem{1, 4000, false, false}
}

func (s StateItem) adjust() (int, int) {
	a, b := s.lower, s.upper
	if s.lowerStrict {
		a++
	}
	if s.upperStrict {
		b--
	}
	return a, b
}

func (s StateItem) isNotValid() bool {
	res := s.diff()
	return res <= 0
}
func (s StateItem) diff() int {
	lower, upper := s.adjust()
	return upper - lower + 1
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

	res := common.HandleTasks(
		func() uint64 { return part1(parts, workflows) },
		func() uint64 { return part2(workflows) },
	)

	fmt.Println(res)
}

func part1(parts map[Part]bool, workflows map[string]Workflow) uint64 {
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

	var sum uint64 = 0
	for p, v := range parts {
		if v {
			sum += uint64(p.x + p.m + p.a + p.s)
		}
	}
	return sum
}

func part2(workflows map[string]Workflow) uint64 {
	return traverse(workflows)
}

func parseWorkFlows(data string) map[string]Workflow {
	res := make(map[string]Workflow)
	for _, line := range strings.Split(data, "\n") {
		i := strings.Index(line, "{")
		name := line[:i]
		w := Workflow{name, nil, nil}
		rulesString := line[i+1 : len(line)-1]
		rules := strings.Split(rulesString, ",")
		for _, rule := range rules {
			temp := strings.Split(rule, ":")
			if len(temp) == 1 {
				w.rules = append(w.rules, func(_ Part) string {
					return temp[0]
				})
				w.conds = append(w.conds, Condition{nil, temp[0]})
				continue
			}
			condition, nextName := temp[0], temp[1]
			e := string(condition[0])
			condString := string(condition[1])
			valueStr := condition[2:]
			value, _ := strconv.Atoi(valueStr)
			w.conds = append(w.conds, Condition{&Rule{e, condString, value}, nextName})
			w.rules = append(w.rules, func(p Part) string {
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

type Foo struct {
	name  string
	state [4]StateItem
}

func traverse(workflows map[string]Workflow) uint64 {
	var sum uint64 = 0
	var agg [][4]StateItem

	queue := []Foo{{"in", [4]StateItem{NewState(), NewState(), NewState(), NewState()}}}
	for len(queue) != 0 {
		current := queue[0]
		queue = queue[1:]
		if current.name == "A" {
			var combinations uint64 = 1
			for _, s := range current.state {
				combinations *= uint64(s.diff())
			}
			agg = append(agg, current.state)
			sum += combinations
			continue
		}
		if current.name == "R" {
			continue
		}

		w := workflows[current.name]
		left, right := current.state, current.state
		for _, cond := range w.conds {
			if cond.rule != nil {
				switch [2]string{cond.rule.element, cond.rule.condition} {
				case [2]string{"x", "<"}:
					left = updateState(0, right, cond.rule.value, true, true)
					right = updateState(0, right, cond.rule.value, true, false)

				case [2]string{"x", ">"}:
					left = updateState(0, right, cond.rule.value, false, true)
					right = updateState(0, right, cond.rule.value, false, false)

				case [2]string{"m", "<"}:
					left = updateState(1, right, cond.rule.value, true, true)
					right = updateState(1, right, cond.rule.value, true, false)

				case [2]string{"m", ">"}:
					left = updateState(1, right, cond.rule.value, false, true)
					right = updateState(1, right, cond.rule.value, false, false)

				case [2]string{"a", "<"}:
					left = updateState(2, right, cond.rule.value, true, true)
					right = updateState(2, right, cond.rule.value, true, false)

				case [2]string{"a", ">"}:
					left = updateState(2, right, cond.rule.value, false, true)
					right = updateState(2, right, cond.rule.value, false, false)

				case [2]string{"s", "<"}:
					left = updateState(3, right, cond.rule.value, true, true)
					right = updateState(3, right, cond.rule.value, true, false)

				case [2]string{"s", ">"}:
					left = updateState(3, right, cond.rule.value, false, true)
					right = updateState(3, right, cond.rule.value, false, false)
				}

				if !slices.ContainsFunc(left[:], func(si StateItem) bool { return si.isNotValid() }) {
					queue = append(queue, Foo{cond.next, left})
				} else {
					fmt.Println("----")
				}
			} else {
				queue = append(queue, Foo{cond.next, right})
			}

		}
	}
	for _, v := range agg {
		fmt.Println(v)
	}
	fmt.Println(len(agg))

	return sum
}

func updateState(i int, states [4]StateItem, value int, upperBound bool, notOpposite bool) [4]StateItem {
	newStates := states
	if upperBound {
		if notOpposite {
			if states[i].upper > value {
				newStates[i] = StateItem{states[i].lower, value, states[i].upperStrict, true}
			}
		} else {
			if states[i].lower <= value {
				newStates[i] = StateItem{value, states[i].upper, false, states[i].upperStrict}
			}
		}
	} else {
		if notOpposite {
			if states[i].lower < value {
				newStates[i] = StateItem{value, states[i].upper, true, states[i].upperStrict}
			}
		} else {
			if states[i].upper >= value {
				newStates[i] = StateItem{states[i].lower, value, states[i].lowerStrict, false}
			}
		}
	}
	return newStates
}