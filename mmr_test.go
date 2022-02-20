package mmr

// This file tests basic functions of the MMR data structure
// independent of any specific implementation of MMRs.

import (
	"testing"

	"github.com/vsekhar/mmr/internal/bruteforce"
)

func TestPeaksAndHeights(t *testing.T) {
	sizes := []int{
		1, 3, 5, 13, 16, 32, 50, 900, 20000,
	}
sizes_loop:
	for i, s := range sizes {
		m := bruteforce.New(s)
		mathPeaks, mathHeights := peaksAndHeights(s)
		if len(m.Peaks) != len(mathPeaks) {
			t.Errorf("case %d: expected peaks (%v), got (%v)", i, m.Peaks, mathPeaks)
			continue
		}
		for i, p := range m.Peaks {
			if p.Index != mathPeaks[i] {
				t.Errorf("case %d: expected peaks (%v), got (%v)", i, m.Peaks, mathPeaks)
				continue sizes_loop
			}
			if p.Height != mathHeights[i] {
				t.Errorf("case %d: expected peak at %d to have height %d, got %d", i, p.Index, mathHeights[i], p.Height)
			}
		}
	}
}

func TestIteratorStart(t *testing.T) {
	isBeforeFirstNode := func(i *Iterator) bool {
		if i.n != 0 {
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
	for i, itr := range []*Iterator{
		{}, // zero
		Begin(),
		IterJustBefore(0),
	} {
		if !isBeforeFirstNode(itr) {
			t.Errorf("iterator %d is not at before-first-node: %+v", i, *itr)
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
		if bn.Left == nil || bn.Left.Index != n.Left {
			return false
		}
		if bn.Right == nil || bn.Right.Index != n.Right {
			return false
		}
	}
	if bn.Parent != nil && bn.Parent.Index != n.Parent {
		return false
	}
	return true
}

func TestIteratorBruteForce(t *testing.T) {
	n := 500
	m := bruteforce.New(n * 2) // ensure parents exist
	bruteNode := m.At(0)
	itr := Begin()
	node := itr.Next()
	for i := 0; i < n; i++ {
		if !bruteforceNodeDeepEquals(node, bruteNode) {
			t.Errorf("pos %d: mismatch", i)
		}
		node = itr.Next()
		bruteNode = bruteNode.Next()
	}
}

func TestIterAt(t *testing.T) {
	n := 50
	m := bruteforce.New(n * 2) // ensure parents exist
	for i := 0; i < n; i++ {
		bruteNode := m.At(i)
		node := IterJustBefore(i).Next()
		if !bruteforceNodeDeepEquals(node, bruteNode) {
			t.Errorf("pos %d: got (%v), expected (%v)", i, node, bruteNode)
		}
	}
}

func TestPath(t *testing.T) {
	// TODO
}
