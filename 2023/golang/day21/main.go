package main

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/mbesida/advent-of-code-2023/common"
	"github.com/rdleal/go-priorityq/kpq"
)

type Garden [][]int

var garden Garden
var n, m int

func main() {
	f := common.InputFileHandle("day21")
	defer f.Close()

	bytes, _ := io.ReadAll(f)
	start := parseGarden(string(bytes))
	t1 := func() uint64 {
		return countReachability(64, map[common.Point]struct{}{start: struct{}{}}, 0)
	}
	t2 := func() uint64 {
		return task2(start)
	}
	fmt.Println(common.HandleTasks(t1, t2))

}

func printGarden(points []common.Point) {
	for i, row := range garden {
		for j, v := range row {
			if slices.Contains(points, common.Point{i, j}) {
				fmt.Print("0")
				continue
			}
			if v == 1 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func parseGarden(data string) common.Point {
	splitted := strings.Split(data, "\n")
	var start common.Point
	garden = make(Garden, len(splitted))
	for i, line := range splitted {
		garden[i] = make([]int, len(line))
		for j, r := range line {
			if r == '.' {
				garden[i][j] = 0
			} else if r == '#' {
				garden[i][j] = 1
			} else if r == 'S' {
				start = common.Point{i, j}
				garden[i][j] = 0
			}
		}
	}
	n = len(garden)
	m = len(garden[0])
	return start
}

func countReachability(target int, starts map[common.Point]struct{}, steps int) uint64 {
	k := uint64(len(starts))
	if steps == target {
		// printGarden(maps.Keys(starts))
		return k
	}
	nextStarts := make(map[common.Point]struct{})
	for p := range starts {
		i, j := cycle(p.I, n), cycle(p.J, m)
		if garden[i][cycle(j+1, m)] == 0 {
			nextStarts[common.Point{p.I, p.J + 1}] = struct{}{}
		}
		if garden[i][cycle(j-1, m)] == 0 {
			nextStarts[common.Point{p.I, p.J - 1}] = struct{}{}
		}
		if garden[cycle(i+1, n)][j] == 0 {
			nextStarts[common.Point{p.I + 1, p.J}] = struct{}{}
		}
		if garden[cycle(i-1, n)][j] == 0 {
			nextStarts[common.Point{p.I - 1, p.J}] = struct{}{}
		}
	}
	return countReachability(target, nextStarts, steps+1)
}

func cycle(index int, size int) int {
	mod := index % size
	if mod == 0 {
		return 0
	}
	if index < 0 {
		return size + mod
	}
	return mod
}

func calculateDistances(start common.Point) map[common.Point]int {
	distances := make(map[common.Point]int)
	pq := kpq.NewKeyedPriorityQueue[common.Point, int](func(a, b int) bool {
		return a < b
	})
	pq.Push(start, 0)

	for !pq.IsEmpty() {
		p, dist, _ := pq.Pop()
		if _, ok := distances[p]; ok {
			continue
		}
		distances[p] = dist
		for _, diff := range []common.Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
			neigbour := common.Point{p.I + diff.I, p.J + diff.J}
			if neigbour.I >= 0 && neigbour.I < n && neigbour.J >= 0 && neigbour.J < m {
				if garden[neigbour.I][neigbour.J] == 0 {
					pq.Push(neigbour, dist+1)
				}
			}
		}
	}
	return distances
}

func task2(start common.Point) uint64 {
	distances := calculateDistances(start)

	var targetSteps uint64 = 26501365
	// n = 131
	// 26501365 % n = 65
	// (26501365 - 65) / 131 = 202300
	k := (targetSteps - uint64(n/2)) / uint64(n)

	var odd, oddEdges, even, evenEdges uint64
	for _, v := range distances {
		if v%2 == 0 && v > n/2 {
			even++
			evenEdges++
		} else if v%2 == 0 {
			even++
		} else if v%2 == 1 && v > n/2 {
			odd++
			oddEdges++
		} else {
			odd++
		}
	}

	totalOdd := odd * (k + 1) * (k + 1)
	totalEven := even * k * k
	totalOddEdges := oddEdges * (k + 1)
	totalEvenEdges := evenEdges * k

	return totalOdd + totalEven - totalOddEdges + totalEvenEdges
}
