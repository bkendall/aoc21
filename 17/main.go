package main

import (
	"fmt"
	"math"
	"os"
)

type bounds struct {
	startX int
	endX   int
	startY int
	endY   int
}

func main() {

	target := input

	// minX := 0
	// for {
	// 	dist := maxXForVelocity(minX)
	// 	if dist >= target.startX {
	// 		break
	// 	}
	// 	minX++
	// }
	// maxX := minX
	// for {
	// 	dist := maxXForVelocity(maxX)
	// 	if dist > target.endX {
	// 		break
	// 	}
	// 	maxX++
	// }

	// fmt.Printf("X between %d and %d\n", minX, maxX)

	// targetHits := 0
	// maxY := absInt(target.startY)
	// maxVX, maxVY := 0, 0
	// for x := minX; x <= maxX; x++ {
	// 	for y := 0; y <= maxY; y++ {
	// 		fmt.Printf("Checking %d, %d\n", x, y)
	// 		if hitsTarget(x, y, target) {
	// 			targetHits++
	// 			if maxVY < y {
	// 				maxVX, maxVY = x, maxInt(y, maxVY)
	// 			}
	// 		}
	// 	}
	// }

	// fmt.Printf("Target hits: %d\n", targetHits)

	// fmt.Printf("Max X, Y: %d, %d\n", maxVX, maxVY)

	// vy := maxVY
	// y := 0
	// for {
	// 	if vy == 0 {
	// 		break
	// 	}
	// 	y += vy
	// 	vy--
	// }

	// fmt.Printf("Height: %d\n", y)

	targetHits := 0
	for x := 0; x <= target.endX; x++ {
		for y := target.startY; y <= absInt(target.startY); y++ {
			if hitsTarget(x, y, target) {
				fmt.Printf("Hit: %d, %d\n", x, y)
				targetHits++
			}
		}
	}
	fmt.Printf("Target hits: %d\n", targetHits)
}

func maxXForVelocity(x int) int {
	pos := 0
	for x > 0 {
		pos += x
		x--
	}
	return pos
}

func hitsTarget(vx, vy int, target bounds) bool {
	x, y := 0, 0
	for {
		if inTarget(x, y, target) {
			return true
		}
		if y < target.startY {
			return false
		}
		x += vx
		y += vy
		vy--
		vx = maxInt(0, vx-1)
	}
}

func inTarget(x, y int, target bounds) bool {
	if x >= target.startX && x <= target.endX {
		if y >= target.startY && y <= target.endY {
			return true
		}
	}
	return false
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

// target area: x=20..30, y=-10..-5
var sample = bounds{
	startX: 20,
	endX:   30,
	startY: -10,
	endY:   -5,
}

// target area: x=25..67, y=-260..-200
var input = bounds{
	startX: 25,
	endX:   67,
	startY: -260,
	endY:   -200,
}
