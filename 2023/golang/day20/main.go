package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
)

type Item struct {
	from, to string
	signal   bool
}

type StatefulModule interface {
	Name() string
	Update(signal bool, from string) int //-1 means do not send signal, 0 - low, 1 - high
}

type FlipFlop struct {
	name string
	isOn bool
}

func (f *FlipFlop) Name() string { return f.name }
func (f *FlipFlop) Update(signal bool, _ string) int {
	if signal {
		return -1
	}
	var pulse int
	if f.isOn {
		pulse = 0
	} else {
		pulse = 1
	}
	f.isOn = !f.isOn
	return pulse
}

type Conjunction struct {
	name   string
	inputs map[string]bool
}

func (c *Conjunction) Name() string { return c.name }
func (c *Conjunction) Update(signal bool, from string) int {
	_, ok := c.inputs[from]
	if !ok {
		log.Fatalf("incorrect state: there is no %s in inputs map of conjunction %s", from, c.name)
	}
	c.inputs[from] = signal

	for _, receivedSignal := range c.inputs {
		if !receivedSignal {
			return 1
		}
	}

	return 0
}

var configuration map[string][]string = make(map[string][]string)
var state map[string]StatefulModule = make(map[string]StatefulModule)

func process(from, to string, signal bool) int {
	sm, ok := state[to]
	if ok {
		return sm.Update(signal, from)
	}
	if to == "broadcaster" {
		return 0
	}
	return -1
}

func main() {
	f := common.InputFileHandle("day20")
	defer f.Close()
	bytes, _ := io.ReadAll(f)
	readInitialState(string(bytes))
	res := doWarmup()
	fmt.Println(res)
}

func doWarmup() uint64 {
	var low, high int
	warmupCount := 1000
	for i := 0; i < warmupCount; i++ {
		l, h := propagateSignals(0, 0, []Item{{"button", "broadcaster", false}})
		low += l
		high += h
	}
	fmt.Println(low, high)
	return uint64(low) * uint64(high)
}

func propagateSignals(low, high int, items []Item) (int, int) {
	if len(items) == 0 {
		return low, high
	}
	var newItems []Item
	for _, ni := range items {
		if ni.signal {
			high++
		} else {
			low++
		}

		nextSignal := process(ni.from, ni.to, ni.signal)
		if nextSignal != -1 {
			outputs := configuration[ni.to]
			next := nextSignal != 0
			for _, out := range outputs {
				newItems = append(newItems, Item{ni.to, out, next})
			}
		}
	}
	return propagateSignals(low, high, newItems)
}

func readInitialState(str string) {
	for _, line := range strings.Split(str, "\n") {
		data := strings.Split(line, " -> ")
		input := data[0]
		output := data[1]
		outData := strings.Split(output, ", ")
		if input[0] == '%' {
			configuration[input[1:]] = outData
			state[input[1:]] = &FlipFlop{input[1:], false}
		} else if input[0] == '&' {
			configuration[input[1:]] = outData
			state[input[1:]] = &Conjunction{input[1:], make(map[string]bool)}
		} else if input == "broadcaster" {
			configuration[input] = outData
		}
	}

	for k, outs := range configuration {
		for _, out := range outs {
			if sm, ok := state[out]; ok {
				if c, ok := sm.(*Conjunction); ok {
					c.inputs[k] = false
				}
			}
		}
	}
}
