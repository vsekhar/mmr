package mmr

// This file tests basic functions of the MMR data structure
// independent of any specific implementation of MMRs.

import (
	"math"
	"testing"
)

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
	{13, 3, 0},
	{14, 3, 0},
	{15, 3, 0},
	{16, 3, 1},
	{17, 3, 0},
	{18, 3, 0},
	{19, 3, 0},
	{20, 3, 1},
	{21, 3, 0},
	{22, 3, 0},
	{23, 3, 0},
	{24, 3, 1},
	{25, 3, 2},
	{26, 3, 0},
	{27, 3, 0},
	{28, 3, 0},
	{29, 3, 1},
	{30, 3, 0},
	{31, 3, 0},
	{32, 3, 0},
	{33, 3, 1},
	{34, 3, 0},
	{35, 3, 0},
	{36, 3, 0},
	{37, 3, 1},
	{38, 3, 2},
	{39, 3, 3},
	{40, 3, 0},

	{0, 4, 0},
	{1, 4, 0},
	{2, 4, 0},
	{3, 4, 0},
	{4, 4, 1},
	{5, 4, 0},
	{6, 4, 0},
	{7, 4, 0},
	{8, 4, 0},
	{9, 4, 1},

	{0, 5, 0},
	{1, 5, 0},
	{2, 5, 0},
	{3, 5, 0},
	{4, 5, 0},
	{5, 5, 1},
	{6, 5, 0},
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

		// ternary
		{1, 3, 0},
		{2, 3, 0, 1},
		{3, 3, 0, 1, 2},
		{4, 3, 3},
		{5, 3, 3, 4},
		{6, 3, 3, 4, 5},
		{7, 3, 3, 4, 5, 6},
		{8, 3, 3, 7},
		{9, 3, 3, 7, 8},
		{10, 3, 3, 7, 8, 9},
		{11, 3, 3, 7, 8, 9, 10},
		{12, 3, 3, 7, 11},
		{13, 3, 12},
	}
	for _, c := range table {
		in, b, out := c[0], c[1], c[2:]
		if result := peaks(in, b); !intSliceEqual(out, result) {
			t.Errorf("peaks(%d, %d): expected '%v', got '%v'", in, b, out, result)
			continue
		}
	}
}

// {value, base, out}
var intLogTable = [][]int{
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

func TestIntLog(t *testing.T) {
	for _, c := range intLogTable {
		in, base, out := c[0], c[1], c[2]
		result := intLog(in, base)
		if result != out {
			t.Errorf("intLog(%d, %d): %d (expected %d)", in, base, result, out)
		}
		builtin := int(math.Log(float64(in)) / math.Log(float64(base)))
		if result != builtin {
			t.Errorf("intLog(%d, %d): %d (builtin %d) - test case may be wrong", in, base, result, builtin)
		}
	}
}

func BenchmarkIntLog(b *testing.B) {
	b.Run("builtin", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, c := range intLogTable {
				_ = int(math.Pow(float64(c[0]), float64(c[1])))
			}
		}
	})
	b.Run("custom", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, c := range intLogTable {
				_ = intLog(c[0], c[1])
			}
		}
	})
}

// {x, y, out}
var intPowTable = [][]int{
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

func TestIntPow(t *testing.T) {
	for _, c := range intPowTable {
		x, y, out := c[0], c[1], c[2]
		result := intPow(x, y)
		if result != out {
			t.Errorf("intPow(%d, %d): %d (expected %d)", x, y, result, out)
		}
		builtin := int(math.Pow(float64(x), float64(y)))
		if result != builtin {
			t.Errorf("intPow(%d, %d): %d (builtin %d) - test case may be wrong", x, y, result, builtin)
		}
	}
}

func BenchmarkIntPow(b *testing.B) {
	b.Run("builtin", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, c := range intPowTable {
				_ = int(math.Pow(float64(c[0]), float64(c[1])))
			}
		}
	})
	b.Run("custom", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, c := range intPowTable {
				_ = intPow(c[0], c[1])
			}
		}
	})
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

func TestChildren(t *testing.T) {
	// pos, height, branching, children...
	table := [][]int{
		// leaves
		{0, 0, 2},

		// binary
		{2, 1, 2, 0, 1},
		{5, 1, 2, 3, 4},
		{9, 1, 2, 7, 8},
		{12, 1, 2, 10, 11},
		{6, 2, 2, 2, 5},
		{13, 2, 2, 9, 12},
		{14, 3, 2, 6, 13},

		// ternary
		{3, 1, 3, 0, 1, 2},
		{16, 1, 3, 13, 14, 15},
		{25, 2, 3, 16, 20, 24},
		{39, 3, 3, 12, 25, 38},
	}
	for _, vals := range table {
		pos, h, b, expected := vals[0], vals[1], vals[2], vals[3:]
		if out := children(pos, h, b); !intSliceEqual(expected, out) {
			t.Errorf("children(%d, %d, %d): expected '%v', got '%v'", pos, h, b, expected, out)
		}
	}

}
