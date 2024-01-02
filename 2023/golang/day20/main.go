package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
	"golang.org/x/exp/maps"
)

const WarmupCount = 1000

var modules map[string]Module = make(map[string]Module)

type Module interface {
	Name() string
	Destinations() []string
	Receive(signal int, from string) Module
	NextSignal() int //-1 means do not send signal, 0 - low, 1 - high
	SetDestinations(dest []string)
}

func send(m Module) (int, int, []Module) {
	low, high := 0, 0
	signal := m.NextSignal()
	var nextModules []Module
	if signal != -1 {
		if signal == 0 {
			low += len(m.Destinations())
		} else {
			high += len(m.Destinations())
		}
		for _, d := range m.Destinations() {
			dest, ok := modules[d]
			if ok {
				newModule := dest.Receive(signal, m.Name())
				nextModules = append(nextModules, newModule)
			}
		}
	}
	return low, high, nextModules
}

type BroadCast struct {
	name         string
	destinations []string
}

func (b *BroadCast) Name() string                   { return b.name }
func (b *BroadCast) Receive(_ int, _ string) Module { return nil }
func (b *BroadCast) NextSignal() int                { return 0 }
func (b *BroadCast) Destinations() []string         { return b.destinations }
func (b *BroadCast) SetDestinations(dest []string) {
	b.destinations = dest
}

type FlipFlop struct {
	name         string
	isOn         bool
	nextToSend   int
	destinations []string
}

func NewFlipFlop(name string) *FlipFlop {
	return &FlipFlop{name, false, -1, nil}
}

func (b *FlipFlop) Name() string { return b.name }
func (b *FlipFlop) Receive(signal int, _ string) Module {
	m := NewFlipFlop(b.name)
	m.destinations = b.destinations
	m.isOn = b.isOn
	m.nextToSend = b.nextToSend
	m.name = b.name
	if signal != 1 {
		if b.isOn {
			m.nextToSend = 0
		} else {
			m.nextToSend = 1
		}
		m.isOn = !b.isOn
	} else {
		m.nextToSend = -1
	}
	return m
}
func (b *FlipFlop) NextSignal() int        { return b.nextToSend }
func (b *FlipFlop) Destinations() []string { return b.destinations }
func (b *FlipFlop) SetDestinations(dest []string) {
	b.destinations = dest
}

type Conjunction struct {
	name         string
	inputs       map[string]int
	destinations []string
}

func NewConjunction(name string) *Conjunction {
	return &Conjunction{name, make(map[string]int), nil}
}

func (b *Conjunction) Name() string { return b.name }

func (b *Conjunction) Receive(signal int, from string) Module {
	m := NewConjunction(b.name)
	m.name = b.name
	m.inputs = b.inputs
	m.destinations = b.destinations
	if _, ok := b.inputs[from]; ok {
		m.inputs[from] = signal
	} else {
		log.Fatalf("incorrect state: signal from unkown module %s", from)
	}
	return m
}
func (b *Conjunction) NextSignal() int {
	for _, v := range b.inputs {
		if v == 0 {
			return 1
		}
	}
	return 0
}
func (b *Conjunction) Destinations() []string { return b.destinations }
func (b *Conjunction) SetDestinations(dest []string) {
	b.destinations = dest
}
func (b *Conjunction) addInput(input string) {
	b.inputs[input] = 0
}

type Test struct{ name string }

func (b *Test) Name() string                   { return b.name }
func (b *Test) Receive(_ int, _ string) Module { return nil }
func (b *Test) NextSignal() int                { return -1 }
func (b *Test) Destinations() []string         { return nil }
func (b *Test) SetDestinations(dest []string)  {}

func main() {
	f := common.InputFileHandle("day20")
	defer f.Close()
	bytes, _ := io.ReadAll(f)
	parseInput(string(bytes))

	res := warmup(modules)
	fmt.Println(res)
}

func warmup(modules map[string]Module) uint64 {
	broadcaster := modules["broadcaster"]
	lowCounter, highCounter := 0, 0
	for i := 0; i < WarmupCount; i++ {
		lowCounter++
		nextModules := []Module{broadcaster}
		for nextModules != nil {
			var newNextModules []Module
			for _, m := range nextModules {
				low, high, next := send(modules[m.Name()])
				if low == 0 && high == 0 {
					continue
				}
				newNextModules = append(newNextModules, next...)
				lowCounter += low
				highCounter += high
			}
			nextModules = newNextModules
			for _, m := range nextModules {
				modules[m.Name()] = m
			}
		}
	}

	fmt.Println(lowCounter, highCounter)
	return uint64(lowCounter) * uint64(highCounter)
}

func parseInput(s string) {
	patterns := make(map[string][]string)
	for _, line := range strings.Split(s, "\n") {
		data := strings.Split(line, " -> ")
		input := data[0]
		output := data[1]
		patterns[input] = strings.Split(output, ", ")
	}

	for _, key := range maps.Keys(patterns) {
		var name string
		if key == "broadcaster" {
			name = key
			modules[key] = &BroadCast{key, nil}
		}
		name = key[1:]
		if key[0] == '%' {
			modules[name] = NewFlipFlop(name)
			patterns[name] = patterns[key]
		}
		if key[0] == '&' {
			modules[name] = NewConjunction(name)
			patterns[name] = patterns[key]
		}
	}

	maps.DeleteFunc(patterns, func(key string, _ []string) bool {
		return key[0] == '%' || key[0] == '&'
	})

	for key, outputs := range patterns {
		setOutputs(key, modules, outputs)
	}
	for key, outputs := range patterns {
		addInput(key, modules, outputs)
	}
}

func setOutputs(key string, modules map[string]Module, outputs []string) {
	var outModules []string
	for _, out := range outputs {
		m, ok := modules[out]
		if ok {
			outModules = append(outModules, m.Name())
		} else {
			outModules = append(outModules, out)
		}
	}
	modules[key].SetDestinations(outModules)
}

func addInput(key string, modules map[string]Module, outputs []string) {
	for _, out := range outputs {
		conj, ok := modules[out].(*Conjunction)
		if ok {
			conj.addInput(key)
		}
	}
}
