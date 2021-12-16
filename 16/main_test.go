package main

import (
	"testing"
)

func TestParsePacket(t *testing.T) {
	if got, _ := parsePacket(hexToStr("D2FE28")); got != uint64(2021) {
		t.Errorf("Unexpected LITERAL operation; got %d, want %d", got, 2021)
	}
	if got, _ := parsePacket(hexToStr("C200B40A82")); got != uint64(3) {
		t.Errorf("Unexpected SUM operation; got %d, want %d", got, 3)
	}

	if got, _ := parsePacket(hexToStr("04005AC33890")); got != uint64(54) {
		t.Errorf("Unexpected PRODUCT operation; got %d, want %d", got, 3)
	}

	if got, _ := parsePacket(hexToStr("880086C3E88112")); got != uint64(7) {
		t.Errorf("Unexpected MIN operation; got %d, want %d", got, 7)
	}

	if got, _ := parsePacket(hexToStr("CE00C43D881120")); got != uint64(9) {
		t.Errorf("Unexpected MAX operation; got %d, want %d", got, 9)
	}

	if got, _ := parsePacket(hexToStr("D8005AC2A8F0")); got != uint64(1) {
		t.Errorf("Unexpected GREATER THAN operation; got %d, want %d", got, 1)
	}

	if got, _ := parsePacket(hexToStr("F600BC2D8F")); got != uint64(0) {
		t.Errorf("Unexpected NOT GREATER THAN operation; got %d, want %d", got, 0)
	}

	if got, _ := parsePacket(hexToStr("9C005AC2F8F0")); got != uint64(0) {
		t.Errorf("Unexpected EQUAL TO operation; got %d, want %d", got, 0)
	}

	if got, _ := parsePacket(hexToStr("9C0141080250320F1802104A08")); got != uint64(1) {
		t.Errorf("Unexpected FULL THING operation; got %d, want %d", got, 1)
	}
}

func TestVersionSum(t *testing.T) {
	// t.Skip()
	versionSum = 0
	parsePacket(hexToStr("8A004A801A8002F478"))
	if got, want := versionSum, uint64(16); got != want {
		t.Errorf("Unexpected version sum; got %d, want %d", got, want)
	}

	versionSum = 0
	parsePacket(hexToStr("620080001611562C8802118E34"))
	if got, want := versionSum, uint64(12); got != want {
		t.Errorf("Unexpected version sum; got %d, want %d", got, want)
	}

	versionSum = 0
	parsePacket(hexToStr("C0015000016115A2E0802F182340"))
	if got, want := versionSum, uint64(23); got != want {
		t.Errorf("Unexpected version sum; got %d, want %d", got, want)
	}

	versionSum = 0
	parsePacket(hexToStr("A0016C880162017C3686B18A3D4780"))
	if got, want := versionSum, uint64(31); got != want {
		t.Errorf("Unexpected version sum; got %d, want %d", got, want)
	}

	versionSum = 0
	parsePacket(hexToStr(input))
	if got, want := versionSum, uint64(974); got != want {
		t.Errorf("Unexpected version sum; got %d, want %d", got, want)
	}
}

// 100111 0 000000001010000 010000 1 00000000010 010100 00001 100100 00011 110001 100000000010000100000100101000001000
// VVVTTT I LLLLLLLLLLLLLLL VVVTTT I LLLLLLLLLLL VVVTTT	AAAAA VVVTTT AAAAA VVVTTT
// equals num packets: 80bts   sum count: 2      literal    1 literal    3
