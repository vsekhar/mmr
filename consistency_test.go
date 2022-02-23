package mmr

import (
	"fmt"
	"testing"
)

func TestConsistency(t *testing.T) {
	// {n1, n2}
	cases := [][]int{
		{1, 2},
		{1, 3},
		{7, 14},

		// Table of 19:
		//
		//           14
		//      6          13
		//   2    5     9      12      17
		// 0  1  3  4  7  8  10  11  15  16  18
		{1, 19},
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
		{19, 19},

		// Stress test
		{4109, 583987},
		{1<<34 - 7, 1<<37 - 499},
	}
	for _, c := range cases {
		n1, n2 := c[0], c[1]
		t.Run(fmt.Sprintf("case {%d, %d}", n1, n2), func(t *testing.T) {
			m1, m2 := New(n1), New(n2)
			path := append(m1.Digest(), m1.To(m2)...)
			path = path.eval()
			d2 := m2.Digest()
			if !d2.Equals(path) {
				t.Errorf("digest mismatch\n  digest: %+v\n  path: %+v", d2, path)
			}
		})
	}
}
