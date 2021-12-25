package main

import (
	"fmt"
	"os"
	"strings"
)

type Dir string

const (
	RIGHT Dir = ">"
	DOWN      = "v"
)

type World struct {
	data map[int]map[int]Dir

	width, height int
}

type Coordinate struct {
	x, y int
}

func (w *World) Add(c Coordinate, d Dir) {
	if _, ok := w.data[c.y]; !ok {
		w.data[c.y] = make(map[int]Dir)
	}
	w.data[c.y][c.x] = d
}

func (w World) String() string {
	str := ""
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			d, ok := w.Get(Coordinate{x, y})
			if ok {
				str += fmt.Sprintf("%s", d)
			} else {
				str += "."
			}
		}
		str += "\n"
	}
	return str
}

func (w *World) Get(c Coordinate) (Dir, bool) {
	if _, ok := w.data[c.y]; !ok {
		return "", false
	}
	v, ok := w.data[c.y][c.x]
	return v, ok
}

func (w *World) AllPlaces(d Dir) <-chan Coordinate {
	out := make(chan Coordinate, 5)
	go func() {
		for y, row := range w.data {
			for x, v := range row {
				if v == d {
					out <- Coordinate{x, y}
				}
			}
		}
		close(out)
	}()
	return out
}

func (w *World) CanMove(c Coordinate) (Coordinate, bool, error) {
	d, ok := w.Get(c)
	if !ok {
		return Coordinate{}, false, fmt.Errorf("could not get %+v", c)
	}
	switch d {
	case RIGHT:
		newX := c.x + 1
		if newX >= w.width {
			newX = 0
		}
		if _, populated := w.Get(Coordinate{newX, c.y}); !populated {
			return Coordinate{newX, c.y}, true, nil
		}
	case DOWN:
		newY := c.y + 1
		if newY >= w.height {
			newY = 0
		}
		if _, populated := w.Get(Coordinate{c.x, newY}); !populated {
			return Coordinate{c.x, newY}, true, nil
		}
	default:
		return Coordinate{}, false, fmt.Errorf("unable to determine move ability for %+v", c)
	}
	return Coordinate{}, false, nil
}

func main() {
	lines := strings.Split(input, "\n")
	height, width := len(lines), len(lines[0])
	w := World{
		data:   make(map[int]map[int]Dir),
		height: height,
		width:  width,
	}
	for y, l := range lines {
		for x, c := range strings.Split(l, "") {
			switch c {
			case ">":
				w.Add(Coordinate{x, y}, RIGHT)
			case "v":
				w.Add(Coordinate{x, y}, DOWN)
			}
		}
	}

	fmt.Printf("World:\n%s", w)

	s := 0
	for {
		moves := 0
		tempWorld := World{
			data:   make(map[int]map[int]Dir),
			height: height,
			width:  width,
		}
		newWorld := World{
			data:   make(map[int]map[int]Dir),
			height: height,
			width:  width,
		}
		for c := range w.AllPlaces(RIGHT) {
			newCoordinate, canMove, err := w.CanMove(c)
			if err != nil {
				logFatal("could not check CanMove(%+v): %v", c, err)
			}
			if canMove {
				moves++
				tempWorld.Add(newCoordinate, RIGHT)
				newWorld.Add(newCoordinate, RIGHT)
			} else {
				tempWorld.Add(c, RIGHT)
				newWorld.Add(c, RIGHT)
			}
		}
		// Copy all the downs over and make w the new world.
		for c := range w.AllPlaces(DOWN) {
			tempWorld.Add(c, DOWN)
		}
		for c := range tempWorld.AllPlaces(DOWN) {
			newCoordinate, canMove, err := tempWorld.CanMove(c)
			if err != nil {
				logFatal("could not check CanMove(%+v): %v", c, err)
			}
			if canMove {
				moves++
				newWorld.Add(newCoordinate, DOWN)
			} else {
				newWorld.Add(c, DOWN)
			}
		}
		w = newWorld
		s++
		fmt.Printf("Iteration %d:\n%s", s, newWorld)
		if moves == 0 {
			break
		}
	}
}

func logFatal(s string, args ...interface{}) {
	fmt.Printf(s, args...)
	os.Exit(1)
}

var tiny = `...>...
.......
......>
v.....>
......>
.......
..vvv..`

var sample = `v...>>.vv>
.vv>>.vv..
>>.>v>...v
>>v>>.>.v.
v>v.vv.v..
>.>>..v...
.vv..>.>v.
v.v..>>v.v
....v..v.>`

var input = `>v.v>vv.>.>.vv.>v>...>.vvvvv.>>v.v.v....v>v...>.>vvvv.v..v.>..v.v>>>>..>v>>......vvvv.vvvv.>...v.>.v>..v.>v...vv>v>.>vv.>>>..>v>v>.....>>..
............>..v...vv>....>..v>..>..>.vv.......v>>..vv>>..>>..v..>>.v....>vv...>vvv>vv>>v>.v>...vvvvv....v>.v..>>....>..v.vv.v>...>...v...v
...v..v>...v.vv>>.v>.vv>.vv.>..v>>.v.vvv....>v.v>...>.v>...>v.>v....v.v.>.>v>>..>.......>...>>>>v.vv>v..v...v.>..>v.....>.>vvv.>...>>>v..>>
..>........>.v..v.....>vv>..vvv..>vv......>.....>..vv.v>>v.>>....v.v.v>v.>v...>.vvv>....v>>..>.v>.v>.>>vv.>..>..vv.v.v>v>>vv.>...v...>v>...
vv.v...>.>>.>v..v..>....v.>>v.v...v.....>v.>.vv>v...v.v>v>.>>vv.v.v..vvv>.vv>.v>.>..v.v.>..vv.....>v>>....>.v>v..v..>....>>>...>v.v>.v..vv.
.>....>>..v>v.>>vv>...>>..v.>>.........v..v.v.v>..>..v>.>.>v...vv.>...>v.v.v>vv.v..>>.v>v>.>.....>v..>>..>v>..>>.....vvv.>.>..v.vv..>v.>.>.
.vv..vv>...v>..>...vv...vv..>v..v>...>.>>vvvv..v..v.v>.>>..vvv>v>.v>.v.vv.v.>.v>....v>.v..v.v.>v.....>...>v...>v.v>>......>>>.vvv.vv.>v...>
>.vv>...>v...v.v..>>.>v....vvvv>>v..v.vv>v>..>.>v.>.v>.vv...v>vv>v.>vv...v>...>.>vv.vv..v.>....>..>...>>>.....>v...v.>>v>.v.>...>....>..v.>
>>vv>....v.>..v.>...>..>v.>..>.>..v.v.....vvv>>.>.vv>.>..>...v>.....v>vv>.>>>.>.>.>.v>>>...>>v>v>v>.>>>.>..v>v.>.vv>...vvv>.v..v...v.v.>>.>
.v......>vv...>v>..>..>v>>vv..>..v.v...v..>v>..v.v>.>v>.vv.v....v>v..>.......vv>.v.>>v........>.....>vv.>..v>..>>v....>v>.v...v>v..>vv...>v
vv>v..........>...>.vv>.......>>..>>.v...vv>.>vvvv>>>.v>..v.v>v>>.v..v.vvvvv..v.>>.....vv...vv........v...vv.>>v.>...v.>...>.v>...v>>v>>.v.
..v.vv..v.....>>>>.vvv>v>>....v.v.>.......>vv>..v>v.>....>.v>v.v.>>vvvv.>>..>>>v.........>.>..v.v>v>.>vv.>..>>>vv....>>>.>vv.v.vv..........
v.....v>>>v>..v>.>.vvv..>.....>v....v>..vvv.>.v.>vv.>..>....v..v>>v.>v..>...vvvv>....v>v>v>.....>......>>.v..>v.v....v>.v>...v.vv>>.>v.vv>.
v>.....v>v..v...>.>.vvv>v>..v>.vv...>>v..vvv>......v>.v>vv..v>>..>...>.vv>....v.v......>.>>v>v..>..v..v....v.>vv>....v>>.>v.>v..>.>v..>v..v
.v>.>>v..v>......v>v>.>....v>.vv.>.v>...v.v..v.vv.>v...vv>.>>>vv.>vv....v...>.v.>..v.v.>.v>>>...>>v.v.>>>.>.>vv.>.>..>v>.>>.v.v>..>>.vvv.>.
vvv>.v.vv.....v>.v..vv...v.v.>.v.v>...v.>vvv.>v...>>v.v...vv>..v.>..v>>>..>v.>v.>.v>.v>..>>v>>.v..v.>v>.>.v>>v.>v....v.>v.>vv.>>...v..v..v.
v.v>v>....v.v.>vv.>>..v>.>>......>v.>.vv...v.vv>.v.v...v.v>....v..v>v>vv..>....>..>v.>..v.....>v.>v..v>.v..>v.>vv.v.v.>...>....vv.v.>>.>v..
..>>>vvv..>>>>v.>>.>..v.v>>>.>v.vv.......vv.>>...>v..v.>>>>v>..vv>v>...v>>.>.....v..v...v.>.>vv....>v.vv...v...>v..>v..>..v>.>.>.>>>>..>.v>
vvv>v.v..v>v..vvvvv....vvvvv..>..v>....vv..>.v...v..vvvv.>.>.vv>.v>v>>v>.>..v...vv...v.>.v>......v>>.v.>......>>.>>v>v.v....v>vv.v>.v.>>.>.
vv..>.>vv.v>.vvv..v>.>>.>v..>v.v>.v.v...v...>v>>...>v...v...vvv...vv..>>v>.v>>....>vv..>.v>v..vv..........v.>.>.>v>.v..>.>>.....v>.vv.>>>vv
..>v....>>....>......v>>v.v..v.v.v...>>.....>vvv........>>v..vv>>.v.v...v>..>vvvv.>>v>>>.>.vv...>vv...v>v.v...v..v..>..>vv.v.v>v>>vv.v.....
...>.>v.v.vv..vvv...v>..v>v>>.>>.>...v...v..>>>.>.v.>v...v...v..>...>...v..v...vv>v.vv>vvv>...v..v>>.vvv.>.>>vv...>vv.vv>>.>vvv.........>.>
.>...>..vvv>>..v>.>.v..>.>v>..v>>......vv...v...vv>vv>..v.vv.>.>v.v..v>>.>..>>v.v.vvvv>>....>..>>...v>..>v>vv.v.v>>v>>.>vv..>.v..>>.>....>.
>.v...>v.v.>.>>.v..>>v.v.v>>.>v..>v.v.>v>..vv..v.>v..>....v.>....>v.>.>v>v>>..vv..v>.v.v>v.>v>>.>>v>>....>.v.>..>..>vvvv...vv>v.>.v...>>..v
.>.>..>.>>.vvv>v.v..>v.v..vv..v.v.>v.>.....>vv..>..>v.>>>v.>v>vvvv.>.vv..vvv>>>>>v..>v>v>.>>.v>v.v..>.v.>..>..v>v.>.......vv.>>..>.>....v>v
..vv...vv..v>>vv.v..>..>.v.....>v.v..v.v.>>v>>v>......>>>.....vv.>>v>>>......>>.vv>v>>>v.>..vv....>.v>>.>.>.>>.v..>.....v.v....>>>.>>v.>>>>
.>.v>.vvv>.>..>.vv>>vv>>v.>v>.>v.v.>.v.v..v..v.>.>.>>..>..>.v.>vv.vvv.v>>>.>..>>v>....>..vv.>..>>v>v..>>..>...>.>>v.>..v>.....>>>...v.v.>v>
v.v...>>.>v.v>v...>.>..>>...>..vv>.v.>vv.>v>.>>>.>vv.>>>>....v..>..>.>>.....>.......>.v.v.vv...>..v>>>.vvv>>>v.....>>.vv..v>...>v.>>....v..
..>v>>.v.>.>>>v>.>v.v>..>>>...v>>>..>...vvv.....>>>.vv>>v....v.v>>.>...v....vv..v.vv.v>.>v>.v.v>>.v....>vv.v>v>>vvv>..v.v...v>v>>....>.>..>
.>v..>.>..vv...>v.v..v...>>.>.v>v..v.v...v>v>.>v>.>.vv>>...v.>>...>v.v...v...>..>vv.vv>..v>>v.v.>.v..v...vvv>.v>>>..>>v>..>.>..v>v..v...vvv
v>>...>>..>....v>.v.>>v..>v..vv>..>.>.>>>v>v>vv.>....v..>>..v>.>v.>.vv>>>>v...>v>>....>v.....vv>v..>..>.....>..>.>v>v>..v..>>.>.>.>...vvv>>
...>v.v>.v.>>...vv...>>..v>...v.>v.v>>.v>.v.>>>v>.....vv..>.>..v>..vvvv.>v.>>>.>.>v...v.v..>.>.v.vv....>..>.v.>>.v>..v..>..>vvv.>vv>v>>....
>....v...v.vvv..v..>vvv>>.>>.>.v.vv..v.....v>v..>..>.>.......v>...v.v..>....>v.v..v.v.>..v.v......>.v.v.>..>....v.>.>.....>...>.v>.v.v...v.
v..v>.>.vv...v>.>>....vv.>v>v..>v>v.>v>v>..v.v>v>>.>.vv..>v>v.>>..>>..v..>vv.v>.v..vv....>.v..v>>v.vv..v>>...vv>.>...>>>.>>.v..>.>v.v.>...>
.>>.>v>.>v>>vv..>>>v.>..v....>>..v.v..>v..v..>.vv....v.>...>...>.>.>...>.v>.....>v>...v.>v..>>v.v........v>...>.>.>.v.v>.vv.v..>v>.v>v>vvv>
.vv.>>v>>.>...>....v>v.>.>>....>.v..>>>.v....v..vv.v>.....vvv.>.v..vvvv..vv.>v>.v..vv.>..v...>.>>.>...v.>..v>...>.v....>.v...v>v>..>v..>>..
vv..v.>...v....>v>.>.v>>v>..>vv.v>....>.v..v..vv.......>.v.v.>...v>>...>vv>>..v.>...>....>>...vvv.v.v>>.>..v.v....>....>v.v.v>v>.....v....v
..>v.>>.v..>vv.v..>>......v.......>..v>.>.v>v...vv..v.v...>.>>...v.vv>>>..>.v>.>>v>.v.>>vv>.>>.>...>..v>.>>v>>>>>.vvvv.>.vv......>vv.v>....
.>vv.>..vv>>>>v..v>v.....>vv.>>.v.>.v.>>>v.v..v>>...>vv.v>v.>v>v.v.>v>...>>...v...v>.....v...v..>..>...v.....>.v..>.v>......v.v>...>.vv.>>v
..v>>....v>..v....>>.vv>v..vv.v>>..v...v...>.>>..v>....v..v.v.v.v>v.>...vv>>..vv....v>>>..v.v.v....v.v......v>>vv.>v....v.v..>>.v..vvv..>vv
..>vv>>>>vv>.vv....>..>..>..>>...>..v...>>v......v>...>.vv>>.v.v>....>>..>.v>v>.>v.v>.v.v.v.>v>v>..vv......>v.>>vv...v>v..v>v>...vvvvv>>>v.
>.>>>..vvv>.>v.>v.>.>v>.....v>v..vvv.v>..v>.v.....>.>..>....v.>.v>...v>..>..>>v.>vv.v..>vv>.>v..v>.v>vv.v>..>>.>>..>.vv...v...v.>..vv.v...v
v>>.>>..v..vv>.vvv...>>..>v.v..>...>..v>..>..>>..v>.>>.....>.>vv>....>..v.vv.v..v>.vv.v.v>.>...v>>...v>.v.>.>..vv..>>>>.v>.v>.vv..>>.......
...>>..>...>.v>>v..>.....>.>....>vv>v>vv>>.vvv.>.....>>.>>>vvv..v..>v>.>.v>>v.>vv>>....v>>..>>v>>.vv>.v.>v.vv>>..v........>v.v..>.....v>.v.
v...>>>v.v.v>v>.>.v....>>>..vv>v.....vv..v.>.v.v>....v.v>>.>..v.v...>.v>>..>..v....>.v>...v.v>.v..>......vvv.>vv..v....v..>v.>>v>v..v..>...
.>>v....>..v>>v.vvvvv>...>..vv.>..>.>.>v.v.v>..vv>..>vvv>.vv.>v>vv.........v........v.>.....v>...>v>.vv.v...v...vv.v.>.v..>.....v.....vv>.v
.vv..v>>v..>.v.>>v>v>....>.v.>..>>vv>.v.v.>>..vv.v>>..v>>v..v.v.v.>>.>>..vv....>.>>..>.....>v.>>>.v.v.v>v>v.>.>v>>v......>v>...>v.v.v..>...
.>>...>..vv>....>...>>.>>..v>v.....>>.v.v..>vv>.v>>.v.>>..>v>>.v..>.>vv>v>v...>..>.v>..v>v>.>v>..>>.......v..vv>v>..>v....vvv...>..>..>>.v.
>...>>>v>>.>.>>...>...>.>v.>.v..>.v....>.v>vvv.>....v....>....v.>.>..>..>.>>.>vv.vv>v>....>vv.>.v>.>..v>v.....v.vv....v..>>.>vv>..>.>....v>
>.>v...vv...v..>.vv>v>.v.v.v...v....v.vv.v>.v..>..vv.....>>.>.>>.>v.....v>......vvvvv>>vv..>...>....>..>..v.......v>.v>...vv>.....>.>>..>.>
v.v....v>.v.>>v...>vv.>v.>v>....vv>.v..>>.v>v...v>>.v>.v>...vv..vv.vv>v.>>...vv>.v..vv>.v>.>.>.v>..>.v.v>v......>.vv.v...>>>...>>.v.>v.>.v.
>>>v..>..>....>v...v.vv....>....v>.>v>v....v>>..vv>v.v..>..vvv>>v>..>.....>>v..v.v>v...>.v>v.v>.v......v...vvv.>.>.>>>v.v...>....>>>v......
v>v>.v>v.>>>..>>>v......>v>vvv...>...>>.v..>>>>....v>v.>..vv.>v.>>v....v.....>>....>vv.v...v.v>..>>>.>...v>.v..v.>.>>>....>>v.>..v..v....>.
....>>vv...v.v...>.>..vvv>v.>.v......v.....>..>v>....>>v.v.v>v>.>>v...>>.>>>>>.>.v.v>v......>..v>.>.>.>.....>..>>..>...>v.v>.>>.>...>v.v..>
..>v.>vv>..vvv...>.v...v.>vv.>.v>.>v>>v>..vv>.v.>>>.>>vv>v>>>..vv>.>v>>.....>>>.v..vvv.v.>v....v>.vv...vv>v>...v...>vvvvv>...>...v>..>vv>>v
.>vv..>.>v.....vv>..>..>...v...>>.v...>.>>.>vv.v....v>>.....>v.>>v>>..v..>>>>vv>>v.v.....v.>.>>>.>v.v>v.v..>..>...vv..>>v>.v.vvv.>.>>.vvv..
>...>vvv>vv...v>v>>..........v.....v>>.......>...v>>.>>v>.....>.....>v>..>.>..v..v.>>v>.v.v>v.>v..>>>>.>vv>..v>>..v.v..>v>.vvv..vv....vvv>.
..v.....>.....>.v.vv...>>.v>>.>>vv.>v...>.v..v>.v...v.v.vv....v>.>vv>>...v.....v..vv.vv..>......>v>>....>.v.>vv>vv>>.v.>v.v..>.>>...v..v>.v
.v...v..>.>v..>>.>>..>vv>......v>v.>..>.>....>...v>v..>.vv..v>.vv......v>>.........vvvvv>vvv.v>v.>>v>..v..>>...>>v....v...>>v..vv>...>....v
.vv>>>.v.vvvvv....v>.vv.>...>>..>v.v....v.vv..>..v....>..v>.>vv>.>.>>>>...>.>.>vv.v>>v.v.>.v.>...>vv....>.v.vvvv.v.>>v>>.v..>.v>...>.>vv...
v>v>>..>..>.>>...v>vvvvv..>vv.vv..>.....>vv.>.>v>.>>>.>>>v..vv..vv....>>v.>v.....>v>v>...v.v..v>>...>>>v.>v..>..vv..>.>>>v...>...vvv....>v.
vvv>v>v..v...>v.vvv.>..>.v.>>>..v....>v...vvv.v>.vv.>>vvv>>...>..v>>vv...>.>.>.v.>.>...vv>>..v>vv..vv....>v.>vv>...>..v>.>>...v..>>>..v.v.v
.....>>..>..>>>>.>>>v.>>.v...>v>>>v.>>..>v>>>v>>.>..vv>v.>>v.v.>>>.....>.>v..>v.>v>>.>.>>.>.vv.>>v...v>v..>v>v......>..>...v..>.vv....>..vv
v>v>v>v>v..>....>v>v.>>>>...v.....v.>.>.>v.>.v...>.>.....>v..>...>..>...v..v>.....v>.>vv>>..>.vv.v...>>v.v.>vv.>vv.>v>..vv>v>.v..v.>.vv.v..
v>......>v..v...v....>v..v..vv.....vv....v.vv..vv>v>>.>v....v.v>v..v>.>.v...>..>...>.vv.v..>....>>v>.v....>.>v>...>>>>.....v.....v.>v....>v
>>..vv....vv.v.v..vv.>.>>>v.>.>>vv...>>>>.vv>>...v.v.v>.v..v>..>.vv..>vvv>.vv.vv.v...>>.....>.v>..v.v>..>..v..>v>vv>.vvv.v.v>.vv>...>v...>.
>v.vv>>.>.v.v.v>.vv>v....>v>.>.vv..>>.>.v.v..v.>>vv>>...vv.>.v.vv.>v.>...v>v>..>.v...v...>.>v>v...v>v.....v.v....v>.v>.v..v>..>..>.>>>>.>v.
.vv.>v.v..v>.....>>>.vv...>v...vvv>.v..>..>v>>..>v.vvv>.>>v.>.vv.......>.v>..v.v.>.......v...>vvv..v.>>v..>..v.>...>.>vv.>.>.>.>.>>.>>>...>
.>v.>..>....>.v.>.>>>>v>...>.....vvv>v..v>..>v..v>.>>..>>vv>>.v>vv....v>>..v...v.v...>v..>v...v.v.v>.....>.>>...v>v.>...v.>>.v.....vvv.vv>.
....v....v>v>.v>.>.>v..>vvv>.>>..>>.>.v.>..v......vv>..v>>>.....>.v.v.>v>....v>.>.v>.>v.>v.>.>.v>>.>v.>>.>...>v..>v...>..v>...v...v....>>vv
vv>..>...>>.vvv..v.v.vv.vvv..>v>v>v...>...>>>>.....>.>>>v..v..vvv...v>....vv>>..v...vv>>vv.>>..v>v.....>.v...>..v.>v.v>>v..>.v>.>v.....>.v.
.>..>.......>..vv.>v>.>.>.>vvvv...>>.>>..>>.v..>.vv>.....v>....>......>>v.v.vv..vv.>.>v.>.....>>..>.v>.v.v..v.>..v...>....>.>>.>>>v.>>.vv..
.v.....>>v>v>...>v.>...>..v..v>.v......v>>...v.>vv.>..v....vvv..vv.v.>..>>>v>>..>.v>....>v.>>.>.>v..v>.v.>v..>....v.>..v..vvvvv.v...>vv.>v.
v.>>>>>.v>.>vv>v.v.v.>>....v.>v..v>v>..v...v.....>.v>..v..v>.>v>.v.>v.v..vvvvvv.v..v>.vv>..v.v>v>...>..>v.v...v>v...v.>v.>.v..vv.v>v.vv>.v.
..v>..v...vv.>......v>>v>.>..v..>....>>...>>..v...v.v..vv>.v..v..>.vvv..v.......v..v>v.vv>.v>.>v>.>v..>...>..>>.v....vvv>v.>...>.>v>>v...v.
>.vv..>v>.>..>v.v.......v..>.>v>>vv>>>>v....>...v>>..>v..vvv>..>>..>vv.>v..vv.v>.....v..>...v.>.>vv...>v...>v>>v...>.>v..vv....v.v.>v..v>.v
....v>.v..v..>....>>v..vv>vv>.v.v..v>.>..>.vv.>>>>vv....>.v>..v.>.vvvv.v>v.v>.......v.v....>..vv..>...>>>v.>>..v..>vv.vv>v....>.>.....v>v..
..vv...vv...v>.vvv.v..>..>..v...>....v.v.>...>v.v.vv.>...>..>..vvv...>>.>v..>..v.>.>......v.v>....>..>.>..>>.>...>>>>......v>.>..>v...>vvv.
v>v......>.>...v.>...vv.>vv.v.>..>v>v.>v.>.v>v..vvv...v>>...v.v>v>>>.v.>......>v...>..v..>..v.v.>v...vv.>.v>...>vv....>.....>v.>>..>>.>>.>v
>v.v.>..v.>v.vvv...>v>......v.>v>v>vvv..vv>.v>>v>>.....vvvv>..>v>.v.v...v.>....>.v>.v.>.....>v.v..v....v.>.>v..v.....>>...v....>>..>>>..vv.
..v..>.v>>.v.vv>.>.v.>>..>v>>v>>>.>.>.vv...v..v..v....v..>....vvv.v..v>.v>vv.vv>.>>>>.>..vv.v.....v...v..>.vvv.>.vv..>....v.>v>vv.>.vv..>.>
.v....v>v..>..v>v>..>>....v.>..v>.v.>....v.....v>>>>v..v>.>.v.v.v>.v..>..v....>.>..>vvv..>.v.>..>>vv....>.v.vv.v.v..v...>vv.v...>....vvv>>.
.>v>.v.v.v>>.>..v.....v>>>>.v>>v.>.....v.>>>>.v..vv.>..>.....vv>>v.>vvv.>v..>>....v.>v.vv..vv..v.....>.v>.vvv..v.>v>....v.>>......>>.>v..v>
>...v.>v.>>...>v..v...v....v>vv>>>..>>v....>.>..v..vv>.v.vvv..>....>...v.v>.>.>v.>v>>...vv..v..vv>..>vvvv..v>...vv...>..vv....v..v..>.>...v
.>>..>...v>..v.>.v>..>..>...v...v>v>.>v...v..v...>>.v.v>vv>v.v..>...v>......v....v..>.vvv.>>.v.v.vvv>>>..vv..v..vv.v>.vv...>.v...>v>v..v.v.
v...v.....vv...v.v.>>.v.vv...v.>>......v.>v...v...>..v..>>..>....v...>v>...v>v.vvv.v...>v>v>..v....v........>v...v>...>vv.>>>>.v.>vv.v....v
..>.v..v>....>...>..v>........>.>.>v>...>v>vv>>.>...>>v>..v.>.v.>.v>>>.v.vvv........>>.>>vv..>v>......v.>vv.v.v>v.>>.>..v.>>>>v>>..vvv>.>v>
.>>v..v..vv..vv..v.>v>v.>.vv>.>.>v...>..vvv>>v.....>.>.vv>v>..vv.vv>v>>....>.>>.....vv>v>.v.v.>...>...>.>>.>.>.v.v..v.>v.>........v.>>.v...
.v......v>v>>.v.v>v...v...>>v..vv..>..v.......v>v..>...v>..>v.........>>.....>v.>>..>...>>..v..vv.>...>.>v....>v..>>>..>..>vvv.......>..>>.
v.....v.....>..>v.>.vv.>>.>v..v>>vv.>v>>..>....v>>..>v.>>v..>.v...>....v..v.>.>v.vv....>..v......>>.v..v....vvvv.>..>.vv.>v>..>>....>v>v.>>
..>..v.>.vv>vv.>..v.v..>v>v>>>....vv.....>.>.vvv.v>>.v......>>v>v>vv.>..v.vv.>vv>..vv..v>>.>.v>.vv>v.>>...v>vv.vvv>>...v>.>>.vvvvvv>.v>v...
>v>>v>.vv>>.v>v..v>..vv...>.vv.vvv>..vv.v>.v..v.>..vvv...>.v.v.v.....vv...>.>..>>.>.vv.....>v.>.v.>..>....>.v>.vv.>v..>...>>>...v.>.>.>>v..
>.>.v...>>.v..>.v.>.v..>..>v.>.v.v..vvv.v>.>.>..v.>...v..>..>v.>..v>>.v>...v...v>.>>..v>v..>.....v>v.v.>.>>.vv>v...v.vv.v..v..>.>.>>v..>..>
...>.v.vv.>v......>>v.v>..>.>>..>..>v..v.>v...v..>....>..v...>>.v>.v.>v.v.v...v....>.v>.......v.>.>>v>.vv>.v.v>>.v.v....vv>.>.vv....vv.v.>>
>.v.v.>v.>.>>vvv.>...>vv.>...>..>>...vvv>.>vv.v...v.>v..vv....v...>.v>.>.vv>..>.vv>.v>..vv..>....>.v.>.v>.v.>>.>>.v.>v.vv.v.>>.vvvv>v..v.>v
.>.>..vv>....vv>..v.v....>>>>.vv.vvv.v...>>>>v>.v>vv>...>.>.>>.v..v.v...>...v...v....v..>vv.>...vv.>......>.>..v..>>v>>....v...>>>>.>..>.>.
>v..vv....>..>>>v.>v.v.v.>.......>>>v.>v.v.>v..v.>>.v.v.v.v>.>vv..>..>.>.>>>..v...>vvv>>v.>>..v..>.>>.....>v.v..v>>v.v.>.v.v.>v.>.v.>...>v>
>>..>>v.>.v>>v>>.vvv...>.>..v>.v.>>.vv.........vv.>.vv.....>>v..>>>...>...>..v..v.>>v.v..v>..>>v.>..>>...vv>.vv...v>>v>vvv>v....>..vv.>.v.>
v>......>>vv>.>v.vvv..>>v.vv>>>>...>.>..v..v..>..v.>>...>>.>v>.>..v>...>v..v.vv.v..>>.>.>..>..>>.>>vv..>>..>..>.>......>vv.>>.v.....vvvv..>
>>...vv.>v.v>>.>...>.v>......>>>.>v>>.>..v.v...v..>v>v.v.v.v>...>...v>.v>v.>.vvv...v>.....v>.v>..vvvvv.....>v.v..v.>>.......>v>....>.>.>.v>
>>.>.v>>>..v.>.>.>.>>>.vv>vv>v.>..v...>>v...>.>.>>.>.....>.v...vv..v.vv.v>..v>>>vv.>>vv..>.>.>v.v....v...>..>>v..v.>..v..>..>.v..v>v...>...
.>>..v....>vvv.vv>v>>.v...>vv..v.>>v.>>..v.>.vv.>...>.>.v..>..>.>.>..vvv.v..>vv.>>>>...>v>>.>.>vv>.>..>>v.>..v.>>..>v>>...v..v...>>..>..>.>
.v.>v....>.......>.>>>>>...vv...>.>>v.>.>v....>....v.vv.>vv.v...v>..>v.v>.v>v>>>.v.v.>.>.v>>.v>.v.>v>>>.v>vv.>.>v.vv>...v..v..v>....>v>v>..
..v.>vv........v>v.vvvv.v..vv.v.vv..>.v..>v.v>>>...>..>.v..>v>>.....>>>.v.>vv>v.>.>v.v>..vvv...v.v>.v.>>.>.v>.v..v>v....>v>>v.>.v.v..v.>.v.
v.vv.>>....>v.v..v..>..>v...>>.......>v>..v>...v.>>..>.v..>..>v..vvv>v>..>v...v>>.>vv..v>>.v......>.v>....vv...>.>.vv..v..v.v.>...>v....>v>
.v.>>..>>v.v.>.v.>..>>v..v....v>vvv.>v...>.>.v.v.>.v..vv.vv.v.>v.v>>.v>.......>>vv...v>v>>.>v..>.v>..vv.>.>vv.v>v......>.v>.>.v.>>v...>>vv>
...>.>>.vv>v>>v.vvv>v.>..v>>..v>vv..v>..>vv.v...>>v.>vvv>>...v....vv.v....vvv.v.v.>..v.>v.>.vv.>v>.v.vvv..v..vv.>v.vv.>.>vvvv.v.>.v.....>..
v.....>>.>.>..v>...>..>>vv.v.vv..>vv.v>v>.>v>>vv...>.v>>.....>>v>>vv...>>v.>.v>.v>>vv.>>>v>......>.....>>v>>vvv>v.>v>....v....vv>...>.vv.vv
.v...>vv..>..vv.>v...>.>..v.>v.....>v>v>.>.v>.....v>v.>v.>>v.>v>>v...>.>.v..>..>.....vv>.>vvv>>..>.>.v...v.v...v.v>>>v>.>v...vvv>.v......>.
...v...v.>>...>..v...>..v>.....>v>vvv..v.v.vv>.vv..v...>....>..>v..v.>vv....v..v>vvv.....v>>.vv.>....vv>v.v.v....>..>>..>v>>.>vvv.>....>...
.v>v.>.v...v...>v.>.....>v..v>.>...........v.v......v...v.vv>v>vv>..>.v>>v>>..v..vv.vv.>.vv>>..vv>...v..v>>>>v>>>v>.>>..v.>.v>..>.>.vv>v.>v
>>>.vv.>.>vvv..v.>.vv...>.vv.v..v.v...vv.v>.>.v.v.>>>.>...v..>..>...>.vv.>vvv..>..v.v>vv....>.>v.>v....>>v>>v..>v.>>.v...>.vv>vv...>>v>v>v.
>..v.v>.v>..>.v.>>.......>v>>..v...vv....v>v>>....vvvv>.vv.>.v.....>>.>>.>.v.>>.v...v..>v>>.v....vv..>>.>...>>.vv.v.>vvv.vv>..vv....>>>.v>.
.v.>.v>.>>>v.>.vvv..>....v>..v.>...v>.....>..>>>v>>.>>>...v>>...v>v>.>vv.>.v>.v.v>.>..v.>.v...v..>>.v>..>v>v....v>.>.....>..v...v.>..v.>.>.
.>.v>.v>>....>.v...v>.>...>v.>>v..>v.v>v>vvv>>>.>>>vv.>..v...>>>..v>.v...v...>v>........v....>..>>.>v.v.>.....v>..>vv.v.>....v>....>.>.>>..
.>..v.vv>.v...v...>>v...>>..v>>..v>..v>>.>..>>.vvvv>v>vvv>v>>.>.>....vvvvv>v..>>.......>v>...>v.vv.vv..>..v....v..v.>v>>v.>.v.vv.>.>.>>v>.v
...v...>..>.v>...>>vv>.>v.vv..v...>.>v>v>.>v..v.>....v>.>>.>>vv.>>.>.>vv.v.v..>.v>v........>v>v>..>......>vv>>>.vv....v>>>>..>>>.>v..>...>v
>v.v.>v.v..v.>.v>.>...>>vvv.v>.v>.>>>.v>....>>>>.>>v...v>.>..v.>.>v.v.>>>.....v.>...vv>..vvvv>.>.v>.v>.>vvvv>>..>>vv.>..v..>v..v...>.vv....
vv>v.v.>.>..>.v.v...>>.v>>.v..>>v>>>v>>.v....>....>.>v.....vv>.....>>>..>.v>>>.v.v.vv.v>>.>.>.v.>>..v.>>..v>..>..v.>.>vv>..>v.>v>vv..>.vv>v
>.v.vv.>>>.v>....v>>v...v.....>..v.vvv..>.vv.>>.>.>.vvv>v....>>..vv..>...v.>v.>v.v>>>>...v..v..v>....>.>.....>.v....v.v.vv.v>...v..>v>.vvvv
...v>>.....vv.......v.v......v.vv.v....v.>>.>v...>.>>vv>..>>.....>>v.>.v>..v.>vv>>..>.v....vv.........>...>.v.>.v...v>vv...vv.....>.v>..>..
.v......v>>>.vv>>v>v.>v.>.v..>.>....>>>..>>...>..v.>.vv.>...>>>...vv...>>.>..v>.v..>>.>.>.v....>.vv>v.>v..v.>...vv.>>...v..v.v.>....v...v..
>vv....v>.>..>v....v.>.>v>v.>....vvvv..v>>>..>vv.>..v.>.>.>..>>....v>>>vv>.v..v>v..v..>.vv..v.v.vv.vvvvv>v..v.v>..>..v>v.>..v...v.vvv..vv..
.....>.>v>.vv..v..>v.v.v...v>.v.v....>>v.vv>>v..>>.....vvv..v.....>.v.>.v...>..v.v...>>>>..>v>>>>....>v>>.......vv>v.>>>...>..>>.v.v..>>.>.
..v.>v..vv.>....>.v>vv>>.vv.>>>..v.>v..>..>>v..vv>..v...>>...v.vv..v.....v.....vv>.v..v.>>..v..v....v>vvv..>>v..>v>>.>.v>..vv.v.>>..>v..>>>
..v.v......>v...>>>v>vv...v...v.>...>>..v.....>..>..v.v>v.>..........v.v.>v>..v..v.v.......v.>.v..v>>v..>..v.v..v>v.vvvv..>..vvv..>.....v.v
....>>.>>...>.>v>>.>>.>vv.vv..v...>..v.....>>v>>..........>..>>>>.v...>....v..>vv>>>>.>v.>>....>>>>>...>.v.....vv....>.>>>v.....v..v>vv>v.v
v...v>.v.v>>v.......v>.......>......>.v>>v>v...vv.v>....v.v.v..>..v>.v....v>.>.v....>>>..v.v..>v>.>v.vv>vv.vv>.>v..vv>v..>v.vv....>.>.v>...
>v.>.>>.v>v>........>>vv>vv.v.>..>.vv..v..>v..v..vv>.v>.v.v.>>..v>vv>>>>v.....>v>..v>..>...v..>....v>>.v.v...>v..vv.v.>v..>>v.......v.>vv.>
.v.>.>v.vv.v>.>>>...v>v.>.v...>...>.>>>.vvv>v..>.>..v>>>v...>...vv.v.>v...v.v.......>..vvv....>..>..>..v>...vvvvv..vvv.>>.>>>>...vv.>.>>.v>
..v.>vv>.>>v>>.v>v>.v>>>...>>.>v..>v.>v.>...v.>.vv.v.>v>>.>....vv>>>v.v..v.v>v>v.>>>v...>>>.>v.v..v>>>v>>...vv.vv>>v>v.v>.vv>.v.>.v>..v.v.>
>..vvvvv.>.v>.vv>...v>v.....>v..v>.>.>..v.v.>v..vv...vv..v....>v>>>v>..>..v>v>.>v.>..>>>..>>>v>>.v..vv.v.>..vv>.>>>..vv.>>..v.....>vv>v..v.
v.v>.>..>>.vvv....>...>v...vv>>.vv.>>.vv.vv..v.>.....>..>>>...v..>....>..>>..v.v.....v>.>>>.v>.>.>>>.v.....vv>>.>>v..v.>.>>...>.v.>v>...vv.
>>>v>>.......vv.v.>...v.>>.v>.>.>v>>.>....v.v>..>.v.>.v>v.v...>....>>vv>vv>>vvvv>v..>vv>v.v....>v>..v>>>.vvv.v....v.>.>>.>vv>.vv>v.v>>.>>vv
..>v>vvv.v.v>.>v..>v.vv..>.vvv.>.v>>..>..>>>>>.>.>.v..>..vv.v.v..vv>>vv..>v...>..>vv...>...>>>v....>>.v.>>v>>.vvv..>v.v.>..>..>......v>..>.
vv.v>.>.>v>vvv.>..v>vv....v.>.v..v>..v.....>.v..v.v...>...>.v.>.v..>>..v.v.vv...>v.>v....vv.v.>.vvvv..v.v>....>.vv>..vv.>.>>..vv.>.vvv.>.>.
..v.v.>>v>..>>..>>..v.>.>...v...v...v.>....v>>.>...vvvv>v...v.>>.v>..v.v..>v.v>>.v>...v.>v>>...v.v>>v>v.vv>>.>>..>v>>>>>..v.>.....>.vv.>>..`
