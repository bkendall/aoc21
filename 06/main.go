package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type LanternFish struct {
	Counter int
}

func NewLanternFish() *LanternFish {
	l := &LanternFish{Counter: 8}
	return l
}

func (l *LanternFish) AddDay() *LanternFish {
	l.Counter--
	if l.Counter < 0 {
		l.Counter = 6
		return NewLanternFish()
	}
	return nil
}

func main() {
	fishStrings := strings.Split(input, ",")
	var lanternFishBuckets = map[int]int{}
	for _, f := range fishStrings {
		lanternFishBuckets[strToInt(f)] += 1
	}

	for d := 0; d < 256; d++ {
		newFish := map[int]int{}
		for age, count := range lanternFishBuckets {
			newAge := age - 1
			if newAge < 0 {
				newAge = 6
				newFish[8] += count
			}
			newFish[newAge] += count
		}
		lanternFishBuckets = newFish
	}

	sum := 0
	for _, count := range lanternFishBuckets {
		sum += count
	}
	fmt.Printf("%d Fish\n", sum)
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

func strToInt(str string) int {
	v, err := strconv.Atoi(str)
	if err != nil {
		logFatal("Failed to convert %#q to int: %v", str, err)
	}
	return v
}

const sample = `3,4,3,1,2`

const input = `5,1,1,3,1,1,5,1,2,1,5,2,5,1,1,1,4,1,1,5,1,1,4,1,1,1,3,5,1,1,1,1,1,1,1,1,1,4,4,4,1,1,1,1,1,4,1,1,1,1,1,5,1,1,1,4,1,1,1,1,1,3,1,1,4,1,4,1,1,2,3,1,1,1,1,4,1,2,2,1,1,1,1,1,1,3,1,1,1,1,1,2,1,1,1,1,1,1,1,4,4,1,4,2,1,1,1,1,1,4,3,1,1,1,1,2,1,1,1,2,1,1,3,1,1,1,2,1,1,1,3,1,3,1,1,1,1,1,1,1,1,1,3,1,1,1,1,3,1,1,1,1,1,1,2,1,1,2,3,1,2,1,1,4,1,1,5,3,1,1,1,2,4,1,1,2,4,2,1,1,1,1,1,1,1,2,1,1,1,1,1,1,1,1,4,3,1,2,1,2,1,5,1,2,1,1,5,1,1,1,1,1,1,2,2,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,4,1,1,1,1,1,3,1,1,5,1,1,1,1,5,1,4,1,1,1,4,1,3,4,1,4,1,1,1,1,1,1,1,1,1,3,5,1,3,1,1,1,1,4,1,5,3,1,1,1,1,1,5,1,1,1,2,2`
