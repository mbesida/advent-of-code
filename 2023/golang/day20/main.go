package main

import (
	"fmt"
	"io"
	"log"
	"slices"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
	"golang.org/x/exp/maps"
)

type Item struct {
	from, to string
	signal   bool
}
type State []Item

type Module interface {
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

func (c *Conjunction) NextSignal() int {
	for _, receivedSignal := range c.inputs {
		if !receivedSignal {
			return 1
		}
	}
	return 0
}
func (c *Conjunction) Name() string { return c.name }
func (c *Conjunction) Update(signal bool, from string) int {
	_, ok := c.inputs[from]
	if !ok {
		log.Fatalf("incorrect state: there is no %s in inputs map of conjunction %s", from, c.name)
	}
	c.inputs[from] = signal

	return c.NextSignal()
}

var configuration map[string][]string = make(map[string][]string)
var state map[string]Module = make(map[string]Module)

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
	t1 := func() uint64 {
		return doWarmup()
	}
	t2 := func() uint64 {
		return countMinPresses()
	}
	fmt.Println(common.HandleTasks(t1, t2))
}

func doWarmup() uint64 {
	var low, high int
	warmupCount := 1000
	for i := 0; i < warmupCount; i++ {
		l, h := propagateSignals(0, 0, []Item{{"button", "broadcaster", false}}, func(_ Item) {})
		low += l
		high += h
	}
	fmt.Println(low, high)
	return uint64(low) * uint64(high)
}

// "lk" "zv" "sp" "xt" - are closest modules to "dg" module which is the only connection to "rx"(our target module)
// probably this problem doesn't have general solution(maybe not works on other input data)
//
// here, idea is that high signal reaches on of these modules with some periodicity
// so we can count on which iteration each module receives high signal and find
// least common multiple of those numbers
// (meaning that on that iteration all these module receive high signal simulteneously)
// This is similar to what we had on day08
func countMinPresses() uint64 {
	counters := map[string]uint64{
		"lk": 0,
		"zv": 0,
		"sp": 0,
		"xt": 0,
	}

	k := 1
	for {
		fn := func(item Item) {
			if slices.Contains(maps.Keys(counters), item.from) && item.signal {
				if counters[item.from] == 0 {
					counters[item.from] = uint64(k)
				}
			}
		}
		propagateSignals(0, 0, []Item{{"button", "broadcaster", false}}, fn)

		if !slices.ContainsFunc(maps.Values(counters), func(i uint64) bool { return i == 0 }) {
			break
		}

		k++
	}

	return lcm(maps.Values(counters))
}

func propagateSignals(low, high int, state State, fn func(i Item)) (int, int) {

	if len(state) == 0 {
		return low, high
	}

	var newState State
	for _, ni := range state {
		if ni.signal {
			high++
		} else {
			low++
		}

		nextSignal := process(ni.from, ni.to, ni.signal)
		if nextSignal != -1 {
			fn(ni) //process for part 2, do nothing for part 1
			outputs := configuration[ni.to]
			next := nextSignal != 0
			for _, out := range outputs {
				newState = append(newState, Item{ni.to, out, next})
			}
		}

	}

	return propagateSignals(low, high, newState, fn)
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

func lcm(nums []uint64) uint64 {
	if len(nums) == 1 {
		return nums[0]
	}
	a := nums[0]
	b := lcm(nums[1:])
	return (a * b) / gcd(a, b)
}

func gcd(a, b uint64) uint64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}
