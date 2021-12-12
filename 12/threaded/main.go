package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var run func(string, []string, map[string]int, func(string) int)

// var maxFn func(string) int

func main() {
	inputStrings := strings.Split(input, "\n")
	smallCaves := map[string]bool{}
	g := map[string]map[string]bool{}
	for _, l := range inputStrings {
		parts := strings.Split(l, "-")
		start, end := parts[0], parts[1]
		if _, ok := g[start]; !ok {
			g[start] = map[string]bool{}
		}
		if _, ok := g[end]; !ok {
			g[end] = map[string]bool{}
		}
		g[start][end] = true
		g[end][start] = true

		if strings.ToLower(start) == start {
			if start != "start" && start != "end" {
				smallCaves[start] = true
			}
		}
		if strings.ToLower(end) == end {
			if end != "start" && end != "end" {
				smallCaves[end] = true
			}
		}
	}

	// fmt.Printf("Map: %+v\n", g)
	fmt.Printf("SmallCaves: %+v\n", smallCaves)

	// maxFn = maxVisits
	successfulPaths := map[string]bool{}
	lock := sync.RWMutex{}

	group := sync.WaitGroup{}

	run = func(curr string, path []string, log map[string]int, maxFn func(string) int) {
		// fmt.Printf("run: %q\n\t%v\n\t%v\n", curr, path, log)

		// If we've reached the end, save the path and bail.
		if curr == "end" {
			path = append(path, curr)
			// successfulPaths = append(successfulPaths, copyArr(path))
			k := strings.Join(path, ",")
			lock.Lock()
			successfulPaths[k] = true
			lock.Unlock()
			return
		}

		// We need to mark the current place as visited.
		log[curr]++
		// And tack it onto the current path.
		path = append(path, curr)

		for next := range g[curr] {
			max, currentVisits := maxFn(next), log[next]
			// fmt.Printf("Next? %v, %d, %d\n", next, max, currentVisits)
			if currentVisits < max {
				run(next, copyArr(path), copyMap(log), maxFn)
			}
		}
	}
	for doubleSmallVisit := range smallCaves {
		doubleSmallVisit := doubleSmallVisit
		specialMaxFn := func(s string) int {
			if s == doubleSmallVisit {
				return 2
			}
			return maxVisits(s)
		}
		group.Add(1)
		go func() {
			run("start", []string{}, map[string]int{}, specialMaxFn)
			group.Done()
		}()
	}

	group.Wait()

	paths := []string{}
	for p := range successfulPaths {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	// for _, p := range paths {
	// 	fmt.Printf("%s\n", p)
	// }
	fmt.Printf("Total paths: %d\n", len(paths))
}

func maxVisits(s string) int {
	if strings.ToLower(s) == s {
		return 1
	}
	return math.MaxInt
}

func copyMap(m map[string]int) map[string]int {
	n := map[string]int{}
	for k, v := range m {
		n[k] = v
	}
	return n
}

func copyArr(a []string) []string {
	n := []string{}
	n = append(n, a...)
	return n
}

func sortString(s string) string {
	arr := strings.Split(s, "")
	sort.Strings(arr)
	return strings.Join(arr, "")
}

func dropFromArr(arr []string, s string) []string {
	n := []string{}
	for _, st := range arr {
		if st != s {
			n = append(n, st)
		}
	}
	return n
}

func containsAll(str string, chars string) bool {
	for _, c := range chars {
		if !strings.Contains(str, string(c)) {
			return false
		}
	}
	return true
}

func strForNum(in map[string]int, v int) (string, bool) {
	for k, val := range in {
		if val == v {
			return k, true
		}
	}
	return "", false
}

func logFatal(s string, args ...interface{}) {
	fmt.Printf(s, args...)
	os.Exit(1)
}

func minInt(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

func maxInt(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func absInt(a int) int {
	return int(math.Abs(float64(a)))
}

func strToInt(str string) int {
	v, err := strconv.Atoi(str)
	if err != nil {
		logFatal("Failed to convert %#q to int: %v", str, err)
	}
	return v
}

const sample = `start-A
start-b
A-c
A-b
b-d
A-end
b-end`

const largeSample = `dc-end
HN-start
start-kj
dc-start
dc-HN
LN-dc
HN-end
kj-sa
kj-HN
kj-dc`

const largerSample = `fs-end
he-DX
fs-he
start-DX
pj-DX
end-zg
zg-sl
zg-pj
pj-he
RW-he
fs-DX
pj-RW
zg-RW
start-pj
he-WI
zg-he
pj-fs
start-RW`

const input = `vn-DD
qm-DD
MV-xy
end-xy
KG-end
end-kw
qm-xy
start-vn
MV-vn
vn-ko
lj-KG
DD-xy
lj-kh
lj-MV
ko-MV
kw-qm
qm-MV
lj-kw
VH-lj
ko-qm
ko-start
MV-start
DD-ko`
