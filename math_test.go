package mmr

// This file tests basic functions of the MMR data structure
// independent of any specific implementation of MMRs.

import (
	"testing"
)

func TestHeightBinary(t *testing.T) {
	// {size, branching factor, height}
	table := [][]int{
		{0, 2, 0},
		{1, 2, 0},
		{2, 2, 1},
		{3, 2, 0},
		{4, 2, 0},
		{5, 2, 1},
		{6, 2, 2},
		{7, 2, 0},
		{8, 2, 0},
		{9, 2, 1},
		{10, 2, 0},
		{11, 2, 0},
		{12, 2, 1},
		{13, 2, 2},
		{14, 2, 3},
		// TODO: implement and test other branching factors
	}
	for _, c := range table {
		in, b, out := c[0], c[1], c[2]
		result := height(in, b)
		if result != out {
			t.Errorf("height(%d, %d): %d (expected %d)", in, b, result, out)
		}
	}
}

func TestPeaks(t *testing.T) {
	// {size, peaks...}
	table := [][]int{
		{1, 0},
		{2, 0, 1},
		{3, 2},
		{4, 2, 3},
		{5, 2, 3, 4},
		{6, 2, 5},
		{7, 6},
		{8, 6, 7},
		{9, 6, 7, 8},
	}
	b := 2
	for _, c := range table {
		in := c[0]
		out := c[1:]
		result := peaks(in, b)
		if len(out) != len(result) {
			t.Errorf("peaks(%d, %d): expected '%v', got '%v'", in, b, out, result)
			continue
		}
		for i, o := range out {
			if o != result[i] {
				t.Errorf("peaks @ %d: %d (expected %d)", in, result[i], o)
			}
		}
	}
}

func TestFloorLog(t *testing.T) {
	// {value, base, out}
	table := [][]int{
		// base 2
		{1, 2, 0},
		{2, 2, 1},
		{3, 2, 1},
		{4, 2, 2},
		{5, 2, 2},
		{6, 2, 2},
		{7, 2, 2},
		{8, 2, 3},
		{9, 2, 3},
		{10, 2, 3},
		{11, 2, 3},
		{12, 2, 3},
		{13, 2, 3},
		{14, 2, 3},
		{15, 2, 3},
		{16, 2, 4},

		// base 10
		{7, 10, 0},
		{22, 10, 1},
		{89, 10, 1},
		{420, 10, 2},
		{1427, 10, 3},

		// other: base 7
		{5, 7, 0},
		{9, 7, 1},
		{32, 7, 1},
		{76, 7, 2},
	}
	for _, c := range table {
		in, base, out := c[0], c[1], c[2]
		result := intLog(in, base)
		if result != out {
			t.Errorf("floorLog(%d): %d (expected %d)", in, result, out)
		}
	}
}

func TestIntPow(t *testing.T) {
	// {x, y, out}
	table := [][]int{
		{1, 2, 1},
		{2, 2, 4},
		{3, 2, 9},
		{4, 2, 16},
		{5, 2, 25},

		{4, 3, 64},
		{5, 3, 125},
		{6, 3, 216},
	}
	for _, c := range table {
		x, y, out := c[0], c[1], c[2]
		result := intPow(x, y)
		if result != out {
			t.Errorf("intPow(%d, %d): %d (expected %d)", x, y, result, out)
		}
	}
}

func TestSiblings(t *testing.T) {
	// left, right, height, branching factor
	table := [][]int{
		{0, 1, 0, 2},
		{3, 4, 0, 2},
		{7, 8, 0, 2},
		{10, 11, 0, 2},
	}
	for _, vals := range table {
		left, right, h, b := vals[0], vals[1], vals[2], vals[3]
		if out := rightSibling(left, h, b); out != right {
			t.Errorf("rightSibling(%d, %d, %d) is %d, expected %d", left, h, b, out, right)
		}
		if out := leftSibling(right, h, b); out != left {
			t.Errorf("leftSibling(%d, %d, %d) is %d, expected %d", right, h, b, out, left)
		}
	}
}
