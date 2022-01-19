package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

var countArr = []int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9}

func joinInts(arr []int) int {
	s := ""
	for _, i := range arr {
		s = fmt.Sprintf("%s%d", s, i)
	}
	return strToInt(s)
}

func next() int {
	counter := joinInts(countArr)

	for i := len(countArr) - 1; i >= 0; i-- {
		if countArr[i] == 1 {
			countArr[i] = 9
		} else {
			countArr[i] -= 1
			break
		}
	}

	return counter
}

type ALU struct {
	registers    map[string]int
	instructions []Instruction
	input        func() int
}

func (a ALU) Run() {
	for _, i := range a.instructions {
		switch i.instruction {
		case "inp":
			a.registers[i.one] = a.input()
			// fmt.Printf("Read in %d\n", a.registers[i.one])
		case "mul":
			var two int
			if isRegister(i.two) {
				two = a.registers[i.two]
			} else {
				two = strToInt(i.two)
			}
			a.registers[i.one] *= two
		case "eql":
			var two int
			if isRegister(i.two) {
				two = a.registers[i.two]
			} else {
				two = strToInt(i.two)
			}
			if a.registers[i.one] == two {
				a.registers[i.one] = 1
			} else {
				a.registers[i.one] = 0
			}
		case "add":
			var two int
			if isRegister(i.two) {
				two = a.registers[i.two]
			} else {
				two = strToInt(i.two)
			}
			a.registers[i.one] += two
		case "div":
			var two int
			if isRegister(i.two) {
				two = a.registers[i.two]
			} else {
				two = strToInt(i.two)
			}
			res := a.registers[i.one] / two
			a.registers[i.one] = res
		case "mod":
			var two int
			if isRegister(i.two) {
				two = a.registers[i.two]
			} else {
				two = strToInt(i.two)
			}
			res := a.registers[i.one] % two
			a.registers[i.one] = res
		default:
			logFatal("Unregistered instruction %q\n", i.instruction)
		}
	}
}

func (a ALU) String() string {
	s := ""
	for k, v := range a.registers {
		s += fmt.Sprintf("\t%s: %d\n", k, v)
	}
	return s
}

func isRegister(s string) bool {
	return s == "w" || s == "x" || s == "y" || s == "z"
}

type Instruction struct {
	instruction string
	one, two    string
}

type Result struct {
	modelNumber int
	zRegister   int
}

func worker(id int, insts []Instruction, jobs <-chan int, results chan<- Result) {
	for modelNumber := range jobs {
		n := strings.Split(fmt.Sprintf("%d", modelNumber), "")
		input := func() int {
			r := n[0]
			n = n[1:]
			// fmt.Printf("r: %s\nn: %v\n", r, n)
			return strToInt(r)
		}
		a := ALU{
			registers:    make(map[string]int),
			instructions: insts,
			input:        input,
		}
		a.Run()

		z := -1
		if v, ok := a.registers["z"]; ok {
			z = v
		}
		results <- Result{modelNumber: modelNumber, zRegister: z}
	}
}

func main() {
	insts := []Instruction{}
	for _, l := range strings.Split(input, "\n") {
		// fmt.Printf("Parsing %q\n", l)
		arr := strings.Split(l, " ")
		i := Instruction{}
		i.instruction = arr[0]
		i.one = arr[1]
		if len(arr) >= 3 {
			i.two = arr[2]
		}
		insts = append(insts, i)
	}

	const numWorkers = 32
	jobs := make(chan int, numWorkers)
	results := make(chan Result, numWorkers)
	wg := sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		i := i
		wg.Add(1)
		go func() {
			worker(i, insts, jobs, results)
			wg.Done()
		}()
	}

	rwg := sync.WaitGroup{}
	done := false
	rwg.Add(1)
	go func() {
		defer rwg.Done()
		c := 0
		for result := range results {
			c++
			if c%1e6 == 0 {
				fmt.Printf("Found a model number: %+v\n", result)
			}
			if result.zRegister == 0 {
				done = true
			}
		}
	}()

	for !done {
		mn := next()
		jobs <- mn
		// done = true
	}
	close(jobs)
	wg.Wait()

	close(results)
	rwg.Wait()
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
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		logFatal("could not parse %q: %v\n", str, err)
	}
	return int(n)
}

var tiny = `inp z
inp x
mul z 3
eql z x`

var sample = `inp w
add z w
mod z 2
div w 2
add y w
mod y 2
div w 2
add x w
mod x 2
div w 2
mod w 2`

var input = `inp w
mul x 0
add x z
mod x 26
div z 1
add x 14
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 14
mul y x
add z y
inp w
mul x 0
add x z
mod x 26
div z 1
add x 14
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 2
mul y x
add z y
inp w
mul x 0
add x z
mod x 26
div z 1
add x 14
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 1
mul y x
add z y
inp w
mul x 0
add x z
mod x 26
div z 1
add x 12
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 13
mul y x
add z y
inp w
mul x 0
add x z
mod x 26
div z 1
add x 15
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 5
mul y x
add z y
inp w
mul x 0
add x z
mod x 26
div z 26
add x -12
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 5
mul y x
add z y
inp w
mul x 0
add x z
mod x 26
div z 26
add x -12
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 5
mul y x
add z y
inp w
mul x 0
add x z
mod x 26
div z 1
add x 12
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 9
mul y x
add z y
inp w
mul x 0
add x z
mod x 26
div z 26
add x -7
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 3
mul y x
add z y
inp w
mul x 0
add x z
mod x 26
div z 1
add x 13
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 13
mul y x
add z y
inp w
mul x 0
add x z
mod x 26
div z 26
add x -8
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 2
mul y x
add z y
inp w
mul x 0
add x z
mod x 26
div z 26
add x -5
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 1
mul y x
add z y
inp w
mul x 0
add x z
mod x 26
div z 26
add x -10
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 11
mul y x
add z y
inp w
mul x 0
add x z
mod x 26
div z 26
add x -7
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 8
mul y x
add z y`
