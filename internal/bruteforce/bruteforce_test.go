package bruteforce

import (
	"testing"

	"github.com/vsekhar/mmr/internal/testdata"
)

func TestMMR(t *testing.T) {
	n := 50
	m := New(n)
	if m.Len() != n {
		t.Errorf("expected length of %d, got length of %d", n, m.Len())
	}

	itr := m.At(0)
	for i, node := range testdata.Sequence {
		bruteNode := m.At(i)
		if itr != bruteNode {
			t.Errorf("next iterator out of sync")
		}
		itr = itr.Next()
		if bruteNode.Index != i {
			t.Errorf("expected index %d, got %d", i, bruteNode.Index)
		}
		if bruteNode.Height != node.Height {
			t.Errorf("expected height %d, got %d", node.Height, bruteNode.Height)
		}
		if node.HasChildren() {
			left, right := node.LeftChild().Pos, node.RightChild().Pos
			if bruteNode.Left == nil {
				t.Errorf("expected left child with index %d, got nil", left)
			} else if bruteNode.Left.Index != left {
				t.Errorf("expected left index %d, got %d", left, bruteNode.Left.Index)
			}
			if bruteNode.Right == nil {
				t.Errorf("expected right child with index %d, got nil", right)
			} else if bruteNode.Right.Index != right {
				t.Errorf("expected right index %d, got %d", right, bruteNode.Right.Index)
			}
		} else {
			if bruteNode.Left != nil {
				t.Errorf("expected no left child, got left child with index %d", bruteNode.Left.Index)
			}
			if bruteNode.Left != nil {
				t.Errorf("expected no right child, got right child with index %d", bruteNode.Right.Index)
			}
		}
	}
}

func TestPeaks(t *testing.T) {
	// {size, peaks...}
	cases := [][]int{
		{1, 0},
		{2, 0, 1},
		{8, 6, 7},
		{14, 6, 13},
		{15, 14},
		{16, 14, 15},
	}
	for i, c := range cases {
		size, peaks := c[0], c[1:]
		m := New(size)
		if len(m.Peaks) != len(peaks) {
			t.Errorf("case %d: expected %d peaks, got %d", i, peaks, len(m.Peaks))
		}
		for j, p := range peaks {
			if p != m.Peaks[j].Index {
				t.Errorf("case %d, peak %d: expected index %d, got %d", i, j, p, m.Peaks[j].Index)
			}
		}
	}
}
