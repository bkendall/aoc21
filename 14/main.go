package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	inputStrings := strings.Split(inputInsertions, "\n")

	polymer := map[string]map[string]int{}
	polymerParts := strings.Split(inputStart, "")
	for i := 0; i < len(polymerParts)-1; i++ {
		if _, ok := polymer[polymerParts[i]]; !ok {
			polymer[polymerParts[i]] = map[string]int{}
		}
		polymer[polymerParts[i]][polymerParts[i+1]]++
	}

	fmt.Printf("Poly: %+v\n", polymer)

	templates := map[string]map[string]string{}
	for _, l := range inputStrings {
		parts := strings.Split(l, " -> ")
		pair, res := parts[0], parts[1]
		pairArr := strings.Split(pair, "")
		if _, ok := templates[pairArr[0]]; !ok {
			templates[pairArr[0]] = map[string]string{}
		}
		templates[pairArr[0]][pairArr[1]] = res
	}

	fmt.Printf("Templates: %+v\n", templates)

	steps := 40
	for i := 0; i < steps; i++ {
		fmt.Printf("Step %d...\n", i+1)
		newPolymer := map[string]map[string]int{}
		for one, m := range polymer {
			if _, ok := newPolymer[one]; !ok {
				newPolymer[one] = map[string]int{}
			}
			for two := range m {
				// fmt.Printf("Looking at %s, %s\n", one, two)
				if _, ok := templates[one]; ok {
					if add := templates[one][two]; add != "" {
						// fmt.Printf("Adding %s, %s\n", one, add)
						newPolymer[one][add] += polymer[one][two]
						if _, ok := newPolymer[add]; !ok {
							newPolymer[add] = map[string]int{}
						}
						newPolymer[add][two] += polymer[one][two]
					} else {
						newPolymer[one][two] += polymer[one][two]
					}
					// fmt.Printf("Polymer: %+v\n", newPolymer)
				}
			}
		}
		polymer = newPolymer
	}

	m := map[string]int{}
	for c, _ := range polymer {
		for cc, count := range polymer[c] {
			m[cc] += count
		}
	}
	m[polymerParts[0]] += 1
	fmt.Printf("M: %+v\n", m)

	maxV, maxC := math.MinInt, ""
	minV, minC := math.MaxInt, ""
	for k, v := range m {
		if v > maxV {
			maxV, maxC = v, k
		}
		if v < minV {
			minV, minC = v, k
		}
	}

	fmt.Printf("Max: %s - %d\n", maxC, maxV)
	fmt.Printf("Min: %s - %d\n", minC, minV)

	fmt.Printf("Diff: %d\n", maxV-minV)
}

func maxVisits(s string) int {
	if strings.ToLower(s) == s {
		return 1
	}
	return math.MaxInt
}

func copyMap(m map[string]map[string]int) map[string]map[string]int {
	n := map[string]map[string]int{}
	for k, sm := range m {
		n[k] = map[string]int{}
		for kk, v := range sm {
			n[k][kk] = v
		}
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

const sampleStart = `NNCB`

const sampleInsertions = `CH -> B
HH -> N
CB -> H
NH -> C
HB -> C
HC -> B
HN -> C
NN -> C
BH -> H
NC -> B
NB -> B
BN -> B
BB -> N
BC -> B
CC -> N
CN -> C`

const inputStart = `KOKHCCHNKKFHBKVVHNPN`

const inputInsertions = `BN -> C
OS -> K
BK -> C
KO -> V
HF -> K
PS -> B
OK -> C
OC -> B
FH -> K
NV -> F
HO -> H
KK -> H
CV -> P
SC -> C
FK -> N
VV -> F
FN -> F
KP -> O
SB -> O
KF -> B
CH -> K
VF -> K
BH -> H
KV -> F
CO -> N
PK -> N
NH -> P
NN -> C
PP -> H
SH -> N
VO -> O
NC -> F
BC -> B
HC -> H
FS -> C
PN -> F
CK -> K
CN -> V
HS -> S
CB -> N
OF -> B
OV -> K
SK -> S
HP -> C
SN -> P
SP -> B
BP -> C
VP -> C
BS -> K
FV -> F
PH -> P
FF -> P
VK -> F
BV -> S
VB -> S
BF -> O
BB -> H
OB -> B
VS -> P
KB -> P
SF -> N
PF -> S
HH -> P
KN -> K
PC -> B
NB -> O
VC -> P
PV -> H
KH -> O
OP -> O
NF -> K
HN -> P
FC -> H
PO -> B
OH -> C
ON -> N
VN -> B
VH -> F
FO -> B
FP -> B
BO -> H
CC -> P
CS -> K
NO -> V
CF -> N
PB -> H
KS -> P
HK -> S
HB -> K
HV -> O
SV -> H
CP -> S
NP -> N
FB -> B
KC -> V
NS -> P
OO -> V
SO -> O
NK -> K
SS -> H`
