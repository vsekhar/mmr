package mmr

// This file tests basic functions of the MMR data structure
// independent of any specific implementation of MMRs.

import (
	"testing"
)

func intSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// {pos, branching factor, height}
var heightTable = [][]int{
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

	/*
		// non-binary
		{0, 3, 0},
		{1, 3, 0},
		{2, 3, 0},
		{3, 3, 1},
		{4, 3, 0},
		{5, 3, 0},
		{6, 3, 0},
		{7, 3, 1},
		{8, 3, 0},
		{9, 3, 0},
		{10, 3, 0},
		{11, 3, 1},
		{12, 3, 2},
	*/
}

func TestHeight(t *testing.T) {
	doTestHeight(t)
}

func BenchmarkHeight(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doTestHeight(b)
	}
}

func doTestHeight(t testing.TB) {
	for _, c := range heightTable {
		pos, b, out := c[0], c[1], c[2]
		result := height(pos, b)
		if result != out {
			t.Errorf("height(%d, %d): %d (expected %d)", pos, b, result, out)
		}
	}
}

func TestPeaks(t *testing.T) {
	// {size, branching factor, peaks...}
	table := [][]int{
		// binary
		{1, 2, 0},
		{2, 2, 0, 1},
		{3, 2, 2},
		{4, 2, 2, 3},
		{5, 2, 2, 3, 4},
		{6, 2, 2, 5},
		{7, 2, 6},
		{8, 2, 6, 7},
		{9, 2, 6, 7, 8},

		/*
			// ternary
			{1, 3, 0},
			{2, 3, 0, 1},
			{3, 3, 0, 1, 2},
			{4, 3, 3},
		*/
	}
	for _, c := range table {
		in := c[0]
		b := c[1]
		out := c[2:]
		result := peaks(in, b)
		if len(out) != len(result) || !intSliceEqual(out, result) {
			t.Errorf("peaks(%d, %d): expected '%v', got '%v'", in, b, out, result)
			continue
		}
	}
}

func TestIntLog(t *testing.T) {
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

		// base powers of 2
		{3, 4, 0},
		{4, 4, 1},
		{5, 4, 1},
		{16, 4, 2},
		{63, 4, 2},
		{64, 4, 3},
		{15, 16, 0},
		{16, 16, 1},
		{255, 16, 1},
		{256, 16, 2},

		// base 3
		{1, 3, 0},
		{2, 3, 0},
		{3, 3, 1},
		{7, 3, 1},
		{9, 3, 2},
		{10, 3, 2},
		{80, 3, 3},
		{81, 3, 4},

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
			t.Errorf("intLog(%d, %d): %d (expected %d)", in, base, result, out)
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
		{6, 2, 36},
		{7, 2, 49},
		{8, 2, 64},

		{4, 3, 64},
		{5, 3, 125},
		{6, 3, 216},

		// powers of two
		{4, 5, 1024},
		{16, 4, 65536},
		{256, 2, 65536},
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
	// TODO: invalid siblings?
	// left, right, height, branching factor
	table := [][]int{
		// binary
		{0, 1, 0, 2},
		// 2 @ 1
		{3, 4, 0, 2},
		// 5 @ 1
		// 6 @ 2
		{7, 8, 0, 2},
		// 9 @ 1
		{10, 11, 0, 2},
		// 12 @ 1
		// 13 @ 2
		// 14 @ 3
		{2, 5, 1, 2},
		{9, 12, 1, 2},
		{6, 13, 2, 2},

		// ternary
		{0, 1, 0, 3},
		{1, 2, 0, 3},
		// 3 @ 1
		{4, 5, 0, 3},
		{5, 6, 0, 3},
		// 7 @ 1
		{8, 9, 0, 3},
		{9, 10, 0, 3},
		// 11 @ 1
		{3, 7, 1, 3},
		{7, 11, 1, 3},
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

func TestFirstChild(t *testing.T) {
	// pos, height, branching, first child
	table := [][]int{
		// binary
		{2, 1, 2, 0},
		{5, 1, 2, 3},
		{6, 2, 2, 2},

		// ternary
		{3, 1, 3, 0},
		{7, 1, 3, 4},
		{11, 1, 3, 8},
		{12, 2, 3, 3},
	}
	for _, vals := range table {
		pos, h, b, fc := vals[0], vals[1], vals[2], vals[3]
		if out := firstChild(pos, h, b); out != fc {
			t.Errorf("firstChild(%d, %d, %d) is %d, expected %d", pos, h, b, out, fc)
		}
	}
}
