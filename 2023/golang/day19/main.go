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

type State struct {
	lower, upper Part
}

func NewState() State {
	return State{Part{1, 1, 1, 1}, Part{4000, 4000, 4000, 4000}}
}

func (f State) combinations() uint64 {
	var combinations uint64 = 1
	combinations *= uint64(f.upper.x) - uint64(f.lower.x) + 1
	combinations *= uint64(f.upper.m) - uint64(f.lower.m) + 1
	combinations *= uint64(f.upper.a) - uint64(f.lower.a) + 1
	combinations *= uint64(f.upper.s) - uint64(f.lower.s) + 1
	return combinations
}
func (f State) updateX(value int, upper bool) State {
	var new State
	if upper {
		new = State{f.lower, Part{value, f.upper.m, f.upper.a, f.upper.s}}
	} else {
		new = State{Part{value, f.lower.m, f.lower.a, f.lower.s}, f.upper}
	}
	return new
}
func (f State) updateM(value int, upper bool) State {
	var new State
	if upper {
		new = State{f.lower, Part{f.upper.x, value, f.upper.a, f.upper.s}}
	} else {
		new = State{Part{f.lower.x, value, f.lower.a, f.lower.s}, f.upper}
	}
	return new
}
func (f State) updateA(value int, upper bool) State {
	var new State
	if upper {
		new = State{f.lower, Part{f.upper.x, f.upper.m, value, f.upper.s}}
	} else {
		new = State{Part{f.lower.x, f.lower.m, value, f.lower.s}, f.upper}
	}
	return new
}
func (f State) updateS(value int, upper bool) State {
	var new State
	if upper {
		new = State{f.lower, Part{f.upper.x, f.upper.m, f.upper.a, value}}
	} else {
		new = State{Part{f.lower.x, f.lower.m, f.lower.a, value}, f.upper}
	}
	return new
}

type QueueItem struct {
	name  string
	state State
}

func traverse(workflows map[string]Workflow) uint64 {
	var sum uint64 = 0

	queue := []QueueItem{{"in", NewState()}}
	for len(queue) != 0 {
		current := queue[0]
		queue = queue[1:]
		if current.name == "A" {
			sum += current.state.combinations()
			continue
		}
		if current.name == "R" {
			continue
		}

		w := workflows[current.name]
		for _, cond := range w.conds {
			st := current.state
			if cond.rule == nil {
				queue = append(queue, QueueItem{cond.next, st})
				continue
			}
			var head, tail State
			val := cond.rule.value
			switch [2]string{cond.rule.element, cond.rule.condition} {
			case [2]string{"x", "<"}:
				head = updateState(st, st.upper.x, val, true, func() State {
					return st.updateX(val-1, true)
				})
				tail = updateState(st, st.lower.x, val-1, false, func() State {
					return st.updateX(val, false)
				})

			case [2]string{"x", ">"}:
				head = updateState(st, st.lower.x, val, false, func() State {
					return st.updateX(val+1, false)
				})
				tail = updateState(st, st.upper.x, val+1, true, func() State {
					return st.updateX(val, true)
				})

			case [2]string{"m", "<"}:
				head = updateState(st, st.upper.m, val, true, func() State {
					return st.updateM(val-1, true)
				})
				tail = updateState(st, st.lower.m, val-1, false, func() State {
					return st.updateM(val, false)
				})

			case [2]string{"m", ">"}:
				head = updateState(st, st.lower.m, val, false, func() State {
					return st.updateM(val+1, false)
				})
				tail = updateState(st, st.upper.m, val+1, true, func() State {
					return st.updateM(val, true)
				})

			case [2]string{"a", "<"}:
				head = updateState(st, st.upper.a, val, true, func() State {
					return st.updateA(val-1, true)
				})
				tail = updateState(st, st.lower.a, val-1, false, func() State {
					return st.updateA(val, false)
				})

			case [2]string{"a", ">"}:
				head = updateState(st, st.lower.a, val, false, func() State {
					return st.updateA(val+1, false)
				})
				tail = updateState(st, st.upper.a, val+1, true, func() State {
					return st.updateA(val, true)
				})

			case [2]string{"s", "<"}:
				head = updateState(st, st.upper.s, val, true, func() State {
					return st.updateS(val-1, true)
				})
				tail = updateState(st, st.lower.s, val-1, false, func() State {
					return st.updateS(val, false)
				})

			case [2]string{"s", ">"}:
				head = updateState(st, st.lower.s, val, false, func() State {
					return st.updateS(val+1, false)
				})
				tail = updateState(st, st.upper.s, val+1, true, func() State {
					return st.updateS(val, true)
				})
			}
			queue = append(queue, QueueItem{cond.next, head})
			current.state = tail

		}
	}

	return sum
}

func updateState(state State, current, value int, upper bool, update func() State) State {
	if upper && current > value {
		return update()
	}
	if !upper && current < value {
		return update()
	}
	return state
}
