package mmr

// This file tests basic functions of the MMR data structure
// independent of any specific implementation of MMRs.

import (
	"testing"
)

// {pos, height}
var heightTable = [][]int{
	{0, 0},
	{1, 0},
	{2, 1},
	{3, 0},
	{4, 0},
	{5, 1},
	{6, 2},
	{7, 0},
	{8, 0},
	{9, 1},
	{10, 0},
	{11, 0},
	{12, 1},
	{13, 2},
	{14, 3},
}

func TestHeight(t *testing.T) {
	for _, c := range heightTable {
		pos, out := c[0], c[1]
		result := height(pos)
		if result != out {
			t.Errorf("height(%d): %d (expected %d)", pos, result, out)
		}
	}
}

func TestPeaks(t *testing.T) {
	// {size, peaks...}
	table := [][]int{
		// binary
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
	for _, c := range table {
		in, out := c[0], c[1:]
		if result := peaks(in); !intSliceEqual(out, result) {
			t.Errorf("peaks(%d): expected '%v', got '%v'", in, out, result)
			continue
		}
	}
}

func TestLeftChild(t *testing.T) {
	// pos, height, first child
	table := [][]int{
		// binary
		{2, 1, 0},
		{5, 1, 3},
		{6, 2, 2},
	}
	for _, vals := range table {
		pos, h, fc := vals[0], vals[1], vals[2]
		if out := leftChild(pos, h); out != fc {
			t.Errorf("leftChild(%d, %d) is %d, expected %d", pos, h, out, fc)
		}
	}
}

func TestChildren(t *testing.T) {
	// pos, height, children...
	table := [][]int{
		// leaves
		{0, 0},
		{1, 0},

		// binary
		{2, 1, 0, 1},
		{5, 1, 3, 4},
		{9, 1, 7, 8},
		{12, 1, 10, 11},
		{6, 2, 2, 5},
		{13, 2, 9, 12},
		{14, 3, 6, 13},
	}
	for _, vals := range table {
		pos, h, expected := vals[0], vals[1], vals[2:]
		if out := children(pos, h); !intSliceEqual(expected, out) {
			t.Errorf("children(%d, %d): expected '%v', got '%v'", pos, h, expected, out)
		}
	}

}
