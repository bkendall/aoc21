package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type space struct {
	energy   int
	flashing bool
}

type World struct {
	m       [][]space
	flashes int
}

type coordinate struct {
	x int
	y int
}

func (w *World) Push(arr []space) {
	w.m = append(w.m, arr)
}

func (w *World) maxX() int {
	return len(w.m[0]) - 1
}

func (w *World) maxY() int {
	return len(w.m) - 1
}

func (w *World) Step() int {
	flashStack := []coordinate{}
	flashesThisStep := 0

	for y := 0; y < len(w.m); y++ {
		for x := 0; x < len(w.m[y]); x++ {
			w.m[y][x].energy++
			if w.m[y][x].energy == 10 {
				w.flashes++
				flashesThisStep++
				w.m[y][x].energy = 0
				w.m[y][x].flashing = true
				flashStack = append(flashStack, coordinate{x, y})
			}
		}
	}

	for {
		var flashed coordinate
		if len(flashStack) == 0 {
			break
		}
		flashed, flashStack = flashStack[0], flashStack[1:]
		for _, c := range w.Surrounding(flashed) {
			if w.m[c.y][c.x].flashing {
				continue
			}
			if w.m[c.y][c.x].energy == 9 {
				w.flashes++
				flashesThisStep++
				w.m[c.y][c.x].energy = 0
				w.m[c.y][c.x].flashing = true
				flashStack = append(flashStack, c)
			} else {
				w.m[c.y][c.x].energy++
			}
		}
	}

	for y := 0; y < len(w.m); y++ {
		for x := 0; x < len(w.m[y]); x++ {
			w.m[y][x].flashing = false
		}
	}

	return flashesThisStep
}

func (w *World) Surrounding(c coordinate) []coordinate {
	arr := []coordinate{}
	x, y := c.x, c.y
	// W.
	if x > 0 {
		arr = append(arr, coordinate{x - 1, y})
		// NW.
		if y > 0 {
			arr = append(arr, coordinate{x - 1, y - 1})
		}
	}
	// N.
	if y > 0 {
		arr = append(arr, coordinate{x, y - 1})
		// NE.
		if x < w.maxX() {
			arr = append(arr, coordinate{x + 1, y - 1})
		}
	}
	// E.
	if x < w.maxX() {
		arr = append(arr, coordinate{x + 1, y})
		// SE.
		if y < w.maxY() {
			arr = append(arr, coordinate{x + 1, y + 1})
		}
	}
	// S.
	if y < w.maxY() {
		arr = append(arr, coordinate{x, y + 1})
		// SW.
		if x > 0 {
			arr = append(arr, coordinate{x - 1, y + 1})
		}
	}
	return arr
}

func main() {
	lines := strings.Split(input, "\n")

	world := World{}

	for _, l := range lines {
		line := []space{}
		for _, c := range strings.Split(l, "") {
			line = append(line, space{energy: strToInt(c)})
		}
		world.Push(line)
	}

	s := 0
	for {
		f := world.Step()
		fmt.Printf("After step %d: %d\n", s+1, world.flashes)
		if f == 100 {
			fmt.Printf("Flashed everyone at step %d\n", s+1)
			break
		}
		s++
	}
}

func logFatal(s string, args ...interface{}) {
	fmt.Printf(s, args...)
	os.Exit(1)
}

func strToInt(str string) int {
	v, err := strconv.Atoi(str)
	if err != nil {
		logFatal("Failed to convert %#q to int: %v", str, err)
	}
	return v
}

var sample = `5483143223
2745854711
5264556173
6141336146
6357385478
4167524645
2176841721
6882881134
4846848554
5283751526`

var input = `4836484555
4663841772
3512484556
1481547572
7741183422
8683222882
4215244233
1544712171
5725855786
1717382281`
