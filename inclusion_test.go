package mmr

import (
	"fmt"
	"testing"
)

func TestInclusionProof(t *testing.T) {
	// {pos, n}
	cases := [][]int{
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

		// Table of 19:
		//
		//           14
		//      6          13
		//   2    5     9      12      17
		// 0  1  3  4  7  8  10  11  15  16  18
		{1, 19},
		{2, 19},
		{2, 19},
		{3, 19},
		{4, 19},
		{5, 19},
		{6, 19},
		{7, 19},
		{8, 19},
		{9, 19},
		{10, 19},
		{11, 19},
		{12, 19},
		{13, 19},
		{14, 19},
		{15, 19},
		{16, 19},
		{17, 19},
		{18, 19},

		// Stress test
		{4109, 583987},
		{1<<34 - 7, 1<<37 - 499},
	}
	for _, c := range cases {
		pos, n := c[0], c[1]
		t.Run(fmt.Sprintf("case {%d, %d}", pos, n), func(t *testing.T) {
			// t.Logf("Case %d (pos: %d, n: %d)", i, pos, n)
			p := Inclusion(pos, n)
			// t.Logf("  Path: %+v", p)
			s := run(t, p, pos, n)
			if !s.matchesDigest(n) {
				t.Errorf("digest check failed")
			}
		})
	}
}
