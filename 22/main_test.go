package main

import (
	// "fmt"
	"testing"
)

func TestRangeMap(t *testing.T) {
	m := NewRangeMap()

	m.Add(Spread{0, 10}, true)

	// want := &RangeMap{
	// 	ranges: []Spread{Spread{0, 10}},
	// 	values: []bool{true},
	// }

	m.Add(Spread{11, 20}, false)

	m.Add(Spread{-10, -1}, false)

	m.Add(Spread{-6, -4}, true)

	m.Add(Spread{5, 15}, true)

	t.Logf("%+v\n", m)

	t.Error("LL")
}
