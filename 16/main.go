package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var versionSum uint64

func main() {
	versionSum = 0
	fmt.Printf("%d\n", uint64(math.MaxUint64))
	fmt.Printf("%f\n", float64(math.MaxFloat64))

	s := hexToStr(input)

	fmt.Printf("I: %s\n", s)

	values, rem := parsePacket(s)

	fmt.Printf("rem: %s\n", rem)
	fmt.Printf("value: %d\n", values)

	// fmt.Printf("versionSum: %d\n", versionSum)
}

func hexToStr(str string) string {
	s := ""
	for _, c := range strings.Split(str, "") {
		inst, err := strconv.ParseUint(c, 16, 64)
		if err != nil {
			logFatal("could not parse input: %v\n", err)
		}
		s += fmt.Sprintf("%04b", inst)
	}
	return s
}

func parsePacket(s string) (uint64, string) {
	// fmt.Printf("parsePacket: %s; after %d\n", s, breakAfter)
	values := []uint64{}

	vString := s[:3]
	s = s[3:]
	// fmt.Printf("Version %s, %s", version, s)

	version, err := strconv.ParseUint(vString, 2, 64)
	if err != nil {
		logFatal("could not parse %q: %v", vString, err)
	}
	// fmt.Printf("Version: %d\n", version)
	versionSum += version

	iString := s[:3]
	s = s[3:]
	instruction, err := strconv.ParseUint(iString, 2, 64)
	if err != nil {
		logFatal("could not parse %q: %v", iString, err)
	}
	// fmt.Printf("instruction: %d\n", instruction)

	switch instruction {
	case 4:
		l, rem := parseLiteral(s)
		s = rem
		fmt.Printf("LITERAL: %d\n", l)
		return l, s
	case 0:
		fmt.Printf("SUM...\n")
		values, s = parseOperator(s)
		sum := uint64(0)
		for _, v := range values {
			sum += v
		}
		fmt.Printf("SUM: %d\n", sum)
		return sum, s
	case 1:
		fmt.Printf("PRODUCT...\n")
		values, s = parseOperator(s)
		product := uint64(1)
		fmt.Printf("Product of: %v\n", values)
		for _, v := range values {
			product = product * v
		}
		fmt.Printf("PRODUCT: %d\n", product)
		return product, s
	case 2:
		fmt.Printf("MIN...\n")
		values, s = parseOperator(s)
		var min uint64
		min = math.MaxUint64
		for _, v := range values {
			if v < min {
				min = v
			}
		}
		return min, s
	case 3:
		fmt.Printf("MAX...\n")
		values, s = parseOperator(s)
		var max uint64
		max = 0
		for _, v := range values {
			if v > max {
				max = v
			}
		}
		return max, s
	case 5:
		fmt.Printf(">...\n")
		values, s = parseOperator(s)
		if len(values) != 2 {
			logFatal("got more greater thans than expected: %+v\n", values)
		}
		if values[0] > values[1] {
			return 1, s
		}
		return 0, s
	case 6:
		fmt.Printf("<...\n")
		values, s = parseOperator(s)
		if len(values) != 2 {
			logFatal("got more less thans than expected: %+v\n", values)
		}
		if values[0] < values[1] {
			return 1, s
		}
		return 0, s
	case 7:
		fmt.Printf("EQUALS...\n")
		values, s = parseOperator(s)
		if len(values) != 2 {
			logFatal("got more equals than expected: %+v\n", values)
		}
		fmt.Printf("eq?: %d, %d\n", values[0], values[1])
		if values[0] == values[1] {
			return 1, s
		}
		return 0, s
	default:
		logFatal("not implemented %d\n", instruction)
	}
	return 0, ""
}

func parseLiteral(s string) (uint64, string) {
	bts := ""
	for {
		switch string(s[0]) {
		case "1":
			bts += s[1:5]
			s = s[5:]
		case "0":
			bts += s[1:5]
			s = s[5:]
			literal, err := strconv.ParseUint(bts, 2, 64)
			if err != nil {
				logFatal("could not parse %q: %v", bts, err)
			}
			return literal, s
		}
	}
}

func parseOperator(s string) ([]uint64, string) {
	ltString := s[:1]
	s = s[1:]

	switch ltString {
	case "0":
		lString := s[:15]
		s = s[15:]
		length, err := strconv.ParseUint(lString, 2, 64)
		if err != nil {
			logFatal("failed to parse length %q: %v", lString, err)
		}
		fmt.Printf("length: %d\n", length)
		subPackets := s[:length]
		s = s[length:]
		values := []uint64{}
		var value uint64
		for len(subPackets) > 0 {
			value, subPackets = parsePacket(subPackets)
			values = append(values, value)
		}
		return values, s
	case "1":
		cString := s[:11]
		s = s[11:]
		count, err := strconv.ParseUint(cString, 2, 64)
		if err != nil {
			logFatal("failed to parse count %q: %v", cString, err)
		}
		fmt.Printf("count: %d\n", count)
		values := []uint64{}
		var value uint64
		for c := uint64(0); c < count; c++ {
			value, s = parsePacket(s)
			values = append(values, value)
		}
		return values, s
	default:
		logFatal("parseOperator invalid length type %q\n", ltString)
	}

	return nil, s
}

func logFatal(s string, args ...interface{}) {
	fmt.Printf(s, args...)
	os.Exit(1)
}

const sample = `9C0141080250320F1802104A08`

const input = `820D4A801EE00720190CA005201682A00498014C04BBB01186C040A200EC66006900C44802BA280104021B30070A4016980044C800B84B5F13BFF007081800FE97FDF830401BF4A6E239A009CCE22E53DC9429C170013A8C01E87D102399803F1120B4632004261045183F303E4017DE002F3292CB04DE86E6E7E54100366A5490698023400ABCC59E262CFD31DDD1E8C0228D938872A472E471FC80082950220096E55EF0012882529182D180293139E3AC9A00A080391563B4121007223C4A8B3279B2AA80450DE4B72A9248864EAB1802940095CDE0FA4DAA5E76C4E30EBE18021401B88002170BA0A43000043E27462829318F83B00593225F10267FAEDD2E56B0323005E55EE6830C013B00464592458E52D1DF3F97720110258DAC0161007A084228B0200DC568FB14D40129F33968891005FBC00E7CAEDD25B12E692A7409003B392EA3497716ED2CFF39FC42B8E593CC015B00525754B7DFA67699296DD018802839E35956397449D66997F2013C3803760004262C4288B40008747E8E114672564E5002256F6CC3D7726006125A6593A671A48043DC00A4A6A5B9EAC1F352DCF560A9385BEED29A8311802B37BE635F54F004A5C1A5C1C40279FDD7B7BC4126ED8A4A368994B530833D7A439AA1E9009D4200C4178FF0880010E8431F62C880370F63E44B9D1E200ADAC01091029FC7CB26BD25710052384097004677679159C02D9C9465C7B92CFACD91227F7CD678D12C2A402C24BF37E9DE15A36E8026200F4668AF170401A8BD05A242009692BFC708A4BDCFCC8A4AC3931EAEBB3D314C35900477A0094F36CF354EE0CCC01B985A932D993D87E2017CE5AB6A84C96C265FA750BA4E6A52521C300467033401595D8BCC2818029C00AA4A4FBE6F8CB31CAE7D1CDDAE2E9006FD600AC9ED666A6293FAFF699FC168001FE9DC5BE3B2A6B3EED060`

// 781436321536 is too high.
// 436797739167464
