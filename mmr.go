package mmr

// Algos: https://github.com/mimblewimble/grin/blob/master/core/src/core/pmmr/pmmr.rs
// Finding peaks: https://github.com/mimblewimble/grin/blob/78220febeda94595159ece675e77e26986a3c11d/core/src/core/pmmr/pmmr.rs#L402

import (
	"fmt"
	"math/bits"
)

func pow2(x int) int { return 2 << (x - 1) }
func log2(x int) int { return bits.Len(uint(x)) - 1 }

func leftChild(pos, height int) int  { return pos - pow2(height) }
func rightChild(pos, height int) int { return pos - 1 }

// peaksAndHeights returns the indexes of peaks and their heights for
// an an MMR of size n.
//
// peaksAndHeights runs in log2(n) time.
func peaksAndHeights(n int) (peaks, heights []int) {
	pos := 0
	for n-pos > 0 {
		peakHeight := log2(n - pos + 1)
		peakSize := pow2(peakHeight) - 1
		peaks = append(peaks, pos+peakSize-1)
		heights = append(heights, peakHeight-1)
		pos += peakSize
	}
	return peaks, heights
}

// MMR represents a Merkle Mountain Range data structure. MMR does not itself store
// data, but provides the logic necessary to construct MMRs and derive proofs.
//
// The zero MMR is valid and corresponds to an empty MMR.
type MMR struct {
	peaks   []int // indexes of peaks for an MMR of size n
	heights []int // heights of peaks for an MMR of size n
}

func (i MMR) Pe2aks() []int { return i.peaks }
func (i MMR) Len() int {
	if len(i.peaks) == 0 {
		return 0
	}
	return i.peaks[len(i.peaks)-1] + 1
}

// Advance extends MMR m by one Node and returns the new Node.
//
// Advance runs in amortized constant time.
func Advance(m *MMR) Node {
	var r Node
	r.Pos = m.Len()
	np := len(m.peaks)

	// Non-leaf?
	if len(m.heights) >= 2 && m.heights[np-1] == m.heights[np-2] {
		r.Height = m.heights[np-1] + 1
		m.peaks = m.peaks[:np-2]
		m.heights = m.heights[:np-2]
	}
	m.peaks = append(m.peaks, r.Pos)
	m.heights = append(m.heights, r.Height)

	return r
}

// New returns an MMR of size n.
//
// New runs in log2(n) time.
func New(n int) MMR {
	peaks, heights := peaksAndHeights(n)
	return MMR{
		peaks:   peaks,
		heights: heights,
	}
}

// Digest returns a path that computes the digest of the MMR.
func (m MMR) Digest() Path {
	var path Path
	for _, p := range m.peaks {
		path = append(path, PathEntry{op: PUSHPOS, pos: p})
	}
	return path
}

func inclusionImpl(pos int, node Node, p Path) Path {
	if node.Pos == pos {
		p = append(p, PathEntry{op: PUSHINPUT, pos: pos})
		return p
	}
	l, r := node.LeftChild(), node.RightChild()
	if l.Pos >= pos {
		p = inclusionImpl(pos, l, p)
		p = append(p, PathEntry{op: PUSHPOS, pos: r.Pos})
	} else {
		p = append(p, PathEntry{op: PUSHPOS, pos: l.Pos})
		p = inclusionImpl(pos, r, p)
	}
	p = append(p, PathEntry{op: POP, pos: node.Pos})
	return p
}

// Has returns a path that proves the inclusion of the node at position pos
// in the MMR.
func (m MMR) Has(pos int) Path {
	if pos >= m.Len() {
		panic(fmt.Sprintf("pos %d does not exist in MMR of size %d", pos, m.Len()))
	}
	var path Path
	found := false
	for i, p := range m.peaks {
		if !found {
			if p < pos {
				path = append(path, PathEntry{op: PUSHPOS, pos: p})
				continue
			}
			found = true
			node := Node{Pos: p, Height: m.heights[i]}
			path = inclusionImpl(pos, node, path)
		} else {
			path = append(path, PathEntry{op: PUSHPOS, pos: p})
		}
	}
	return path
}

// To returns a path that can be appended to the digest of m1 to prove that
// the larger MMR m2 contains the smaller MMR m1. I.e.
//
//   m2.Digest().Equals(append(m1.Digest(), m1.To(m2)...).Digest())
//
// If m2 is not equal to or larger than m1, To panics.
func (m1 MMR) To(m2 MMR) Path {
	var path Path
	if m2.Len() < m1.Len() {
		panic(fmt.Sprintf("consistency cannot be shown from MMR of size %d to MMR of size %d", m1.Len(), m2.Len()))
	}
	lastPeak := m1.peaks[len(m1.peaks)-1]
	found := false
	for i, p := range m2.peaks {
		if !found {
			if p < lastPeak {
				continue
			}
			found = true
			node := Node{Pos: p, Height: m2.heights[i]}
			path = inclusionImpl(lastPeak, node, path)
		} else {
			path = append(path, PathEntry{op: PUSHPOS, pos: p})
		}
	}

	// drop instructions in p2 up to and including the PUSHVAL instruction
	// (these will come from the prepended digest of m2)
	start := 0
	for ; start < len(path); start++ {
		if path[start].op == PUSHINPUT {
			start++
			break
		}
	}
	path = path[start:]
	return path
}
