package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

type wins struct {
	one uint64
	two uint64
}

var cache = map[string]wins{}

const winningScore = 21

func playGame(g game, rolls []int) wins {
	playerOneIfTrue := g.p1Turn
	if len(rolls) == 0 {
		if res, ok := cache[g.Hash()]; ok {
			return res
		}
	}
	if len(rolls) < 3 {
		results := wins{}
		for i := 1; i <= 3; i++ {
			rs := rolls[:]
			rs = append(rs, i)
			r := playGame(g, rs)
			results.one += r.one
			results.two += r.two
		}
		return results
	}
	// fmt.Printf("Running step...\n")
	dist := 0
	for _, v := range rolls {
		dist += v
	}
	if playerOneIfTrue {
		g.onePos += dist
		if g.onePos > 10 {
			g.onePos = g.onePos % 10
			if g.onePos == 0 {
				g.onePos = 10
			}
		}
		g.oneScore += g.onePos
		if g.oneScore >= winningScore {
			return wins{one: 1}
		} else {
			g.p1Turn = !g.p1Turn
			res := playGame(g, []int{})
			cache[g.Hash()] = res
			return res
		}
	} else {
		g.twoPos += dist
		if g.twoPos > 10 {
			g.twoPos = g.twoPos % 10
			if g.twoPos == 0 {
				g.twoPos = 10
			}
		}
		g.twoScore += g.twoPos
		if g.twoScore >= winningScore {
			return wins{two: 1}
		} else {
			g.p1Turn = !g.p1Turn
			res := playGame(g, []int{})
			cache[g.Hash()] = res
			return playGame(g, []int{})
		}
	}
}

func main() {
	game := input
	game.p1Turn = true

	res := playGame(game, []int{})

	fmt.Printf("Results: %+v\n", res)

	// fmt.Printf("P1: %+v\nP2: %+v\n", oneLastThreeRolls, twoLastThreeRolls)
	// fmt.Printf("Last roll: %d\n", lastRoll)
	// fmt.Printf("Game over. P1: %d; P2: %d\n", game.oneScore, game.twoScore)
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
	n, err := strconv.ParseInt(str, 2, 64)
	if err != nil {
		logFatal("could not parse %q: %v\n", str, err)
	}
	return int(n)
}

type game struct {
	onePos   int
	oneScore int
	twoPos   int
	twoScore int

	p1Turn bool
}

func (g game) Hash() string {
	return fmt.Sprintf("%d;%d;%d;%d;%t", g.onePos, g.oneScore, g.twoPos, g.twoScore, g.p1Turn)
}

var sample = game{
	onePos: 4,
	twoPos: 8,
}

var input = game{
	onePos: 7,
	twoPos: 10,
}
