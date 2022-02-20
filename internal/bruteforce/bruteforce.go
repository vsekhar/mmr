// Package bruteforce constructs and fully populates an in-memory MMR
// nodes recording positions, heights and pointers to in-memeory
// children. These MMRs can be used to validate more efficient
// algorithms to traverse MMRs that do not construct full trees.
package bruteforce

import (
	"math/bits"
)

type Node struct {
	Index               int
	Height              int
	Parent, Left, Right *Node
}

// Next returns a pointer to the Next node in the mmr,
// or nil if there is no Next node.
//
// Next runs in amortized constant time.
func (n *Node) Next() *Node {
	// TODO: don't depend on the parent existing, restricts
	// iteration to about half of the nodes in the brute force MMR.
	// Hmm... that would require math, which is what the brute force
	// MMR is meant to avoid. So just create really big brute force MMRs.
	if n.Parent == nil {
		return nil
	}
	if n.Parent.Right == n {
		return n.Parent
	}
	c := n.Parent.Right
	for c.Left != nil {
		c = c.Left
	}
	return c
}

type mmr struct {
	Peaks []*Node
}

// Len returns the number of nodes in the MMR
//
// Len runs in log2(n) time.
func (m mmr) Len() int {
	i := 0
	for _, n := range m.Peaks {
		i += n.Index + 1 - i
	}
	return i
}

// At returns a pointer to the node At index pos.
//
// At runs in log2(n) time.
func (m mmr) At(pos int) *Node {
	for _, n := range m.Peaks {
		if n.Index >= pos {
			return atImpl(n, pos)
		}
	}
	return nil
}

func atImpl(n *Node, pos int) *Node {
	if n.Index == pos {
		return n
	}
	if n.Left.Index >= pos {
		return atImpl(n.Left, pos)
	}
	return atImpl(n.Right, pos)
}

func New(n int) *mmr {
	numPeaks := bits.Len(uint(n)) - 1 // log2
	m := &mmr{Peaks: make([]*Node, 0, numPeaks)}
	for i := 0; i < n; i++ {
		numPeaks := len(m.Peaks)
		if numPeaks >= 2 && m.Peaks[numPeaks-1].Height == m.Peaks[numPeaks-2].Height {
			cur := &Node{
				Index:  i,
				Left:   m.Peaks[numPeaks-2],
				Right:  m.Peaks[numPeaks-1],
				Height: m.Peaks[numPeaks-1].Height + 1,
			}
			cur.Left.Parent = cur
			cur.Right.Parent = cur
			m.Peaks = m.Peaks[:numPeaks-2]
			m.Peaks = append(m.Peaks, cur)
		} else {
			m.Peaks = append(m.Peaks, &Node{Index: i})
		}
	}
	return m
}
