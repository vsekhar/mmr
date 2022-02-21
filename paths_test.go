package mmr

import (
	"testing"
)

func TestInclusionProof(t *testing.T) {
	// {pos, n}
	cases := [][]int{
		{9, 19},
		{10, 19},

		// single peaks
		{0, 1},
		{1, 3},
		{3, 7},
		{5, 15},

		// two peaks at the same level
		{0, 2},
		{1, 2},
		{7, 9},
		{5, 14},
		{8, 14},

		// Stress test
		{4109, 583987},
		{1<<34 - 7, 1<<37 - 499},
	}
	for i, c := range cases {
		pos, n := c[0], c[1]
		t.Logf("Case %d (pos: %d, n: %d)", i, pos, n)
		p := Inclusion(pos, n)
		t.Logf("  Path: %+v", p)
		run(t, p, pos, n)
	}
}
