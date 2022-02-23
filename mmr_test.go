package mmr

// This file tests basic functions of the MMR data structure
// independent of any specific implementation of MMRs.

import (
	"fmt"
	"testing"

	"github.com/vsekhar/mmr/internal/bruteforce"
)

func TestPeaksAndHeights(t *testing.T) {
	sizes := []int{
		1, 3, 5, 13, 16, 32, 50, 900, 20000,
	}
	for i, s := range sizes {
		t.Run(fmt.Sprintf("size %d", s), func(t *testing.T) {
			m := bruteforce.New(s)
			mathPeaks, mathHeights := peaksAndHeights(s)
			if len(m.Peaks) != len(mathPeaks) {
				t.Errorf("case %d: expected peaks (%v), got (%v)", i, m.Peaks, mathPeaks)
			}
			for i, p := range m.Peaks {
				if p.Index != mathPeaks[i] {
					t.Errorf("case %d: expected peaks (%v), got (%v)", i, m.Peaks, mathPeaks)
				}
				if p.Height != mathHeights[i] {
					t.Errorf("case %d: expected peak at %d to have height %d, got %d", i, p.Index, mathHeights[i], p.Height)
				}
			}
		})
	}
}

func TestIteratorStart(t *testing.T) {
	isBeforeFirstNode := func(i MMR) bool {
		if i.Len() != 0 {
			return false
		}
		if len(i.peaks) != 0 {
			return false
		}
		if len(i.heights) != 0 {
			return false
		}
		return true
	}
	for i, m := range []MMR{
		{}, // zero
		New(0),
	} {
		if !isBeforeFirstNode(m) {
			t.Errorf("iterator %d is not at before-first-node: %+v", i, m)
		}
	}
}

func bruteforceNodeDeepEquals(n Node, bn *bruteforce.Node) bool {
	if bn == nil {
		return false
	}
	if n.Pos != bn.Index {
		return false
	}
	if n.Height != bn.Height {
		return false
	}
	if n.HasChildren() {
		if bn.Left == nil || bn.Left.Index != n.LeftChild().Pos {
			return false
		}
		if bn.Right == nil || bn.Right.Index != n.RightChild().Pos {
			return false
		}
	}
	// TODO: compare bn.Parent with computed parent of n.
	return true
}

func TestAdvanceBruteForce(t *testing.T) {
	n := 500
	bm := bruteforce.New(n * 2) // ensure parents exist
	bruteNode := bm.At(0)
	m := MMR{}
	node := Advance(&m)
	for i := 0; i < n; i++ {
		if !bruteforceNodeDeepEquals(node, bruteNode) {
			t.Errorf("pos %d: mismatch", i)
		}
		node = Advance(&m)
		bruteNode = bruteNode.Next()
	}
}

func TestAt(t *testing.T) {
	n := 50
	m := bruteforce.New(n * 2) // ensure parents exist
	for i := 0; i < n; i++ {
		bruteNode := m.At(i)
		m := New(i)
		node := Advance(&m)
		if !bruteforceNodeDeepEquals(node, bruteNode) {
			t.Errorf("pos %d: got (%v), expected (%v)", i, node, bruteNode)
		}
	}
}
