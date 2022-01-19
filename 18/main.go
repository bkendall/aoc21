package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Pair struct {
	leftValue  int
	rightValue int
	leftPair   *Pair
	rightPair  *Pair

	parent *Pair
}

var beginNumber = regexp.MustCompile(`^\d+,.+`)

var endNumber = regexp.MustCompile(`\d+$`)

func newPair(s string) *Pair {
	p := &Pair{}
	// fmt.Printf("Parsing %q\n", s)
	if beginNumber.MatchString(s) {
		arr := strings.SplitN(s, ",", 2)
		p.leftValue = strToInt(arr[0])
		// fmt.Printf("Found left value: %d\n", p.leftValue)
		s = arr[1]
	} else {
		counter := 1
		i := 1
		for counter > 0 {
			switch string(s[i]) {
			case "[":
				counter++
			case "]":
				counter--
			}
			i++
		}
		subStr := s[1 : i-1]
		s = s[i:]
		p.leftPair = newPair(subStr)
		p.leftPair.parent = p
	}
	s = strings.TrimLeft(s, ",")
	// fmt.Printf("Whats this: %q\n", s)
	if endNumber.MatchString(s) {
		arr := strings.Split(s, ",")
		p.rightValue = strToInt(arr[len(arr)-1])
		// fmt.Printf("Found right value: %d\n", p.rightValue)
	} else if len(s) > 0 {
		counter := 0
		i := 0
		for i == 0 || counter > 0 {
			switch string(s[i]) {
			case "[":
				counter++
			case "]":
				counter--
			}
			i++
		}
		subStr := s[1 : i-1]
		p.rightPair = newPair(subStr)
		p.rightPair.parent = p
	}
	return p
}

func (p Pair) String() string {
	s := "["
	if p.leftPair != nil {
		s += p.leftPair.String()
	} else {
		s += fmt.Sprintf("%d", p.leftValue)
	}
	s += ","
	if p.rightPair != nil {
		s += p.rightPair.String()
	} else {
		s += fmt.Sprintf("%d", p.rightValue)
	}
	s += "]"
	return s
}

func (p *Pair) resolve(depth int) *Pair {
	fmt.Printf("resolve depth %d\n", depth)
	if depth >= 3 {
		if p.leftPair != nil {
			fmt.Printf("too deep!\n")
			ptr := p.parent
			for ptr != nil {
				if ptr.leftPair == nil {
					ptr.leftValue += p.leftValue
					break
				}
				ptr = ptr.parent
			}
		}
	} else {
		p.leftPair.resolve(depth + 1)
	}
	return p
}

func (p *Pair) Add(n *Pair) *Pair {
	new := &Pair{
		leftPair:  p,
		rightPair: n,
	}
	return new.resolve(0)
}

func main() {
	lines := strings.Split(sample, "\n")
	params := make([]*Pair, len(lines))
	for i, l := range lines {
		params[i] = newPair(l[1 : len(l)-1])
		// fmt.Printf("%s\n", params[i])
	}

	sum := params[0]
	params = params[1:]
	for i, p := range params {
		sum = sum.Add(p)
		fmt.Printf("%d: %s\n", i, sum)
	}
}

func strToInt(str string) int {
	v, err := strconv.Atoi(str)
	if err != nil {
		logFatal("Failed to convert %#q to int: %v", str, err)
	}
	return v
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

var sample = `[[[[4,3],4],4],[7,[[8,4],9]]]
[1,1]`

var fullSample = `[[[0,[5,8]],[[1,7],[9,6]]],[[4,[1,2]],[[1,4],2]]]
[[[5,[2,8]],4],[5,[[9,9],0]]]
[6,[[[6,2],[5,6]],[[7,6],[4,7]]]]
[[[6,[0,7]],[0,9]],[4,[9,[9,0]]]]
[[[7,[6,4]],[3,[1,3]]],[[[5,5],1],9]]
[[6,[[7,3],[3,2]]],[[[3,8],[5,7]],4]]
[[[[5,4],[7,7]],8],[[8,3],8]]
[[9,3],[[9,9],[6,[4,9]]]]
[[2,[[7,7],7]],[[5,8],[[9,3],[0,2]]]]
[[[[5,2],5],[8,[3,7]]],[[5,[7,5]],[4,4]]]`

var input = ``
