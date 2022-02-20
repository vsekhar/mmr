package mmr

// Algos: https://github.com/mimblewimble/grin/blob/master/core/src/core/pmmr/pmmr.rs
// Finding peaks: https://github.com/mimblewimble/grin/blob/78220febeda94595159ece675e77e26986a3c11d/core/src/core/pmmr/pmmr.rs#L402

import (
	"math/bits"
)

func pow2(x int) int { return 2 << (x - 1) }
func log2(x int) int { return bits.Len(uint(x)) - 1 }

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

// Node records the position of a node, its height, its parents and its
// children (if any).
//
// The zero Node is invalid.
type Node struct {
	Pos         int
	Height      int
	Parent      int
	Left, Right int // set only if Height > 0
}

func (n Node) HasChildren() bool { return n.Height > 0 }

// Iterator supports walking the MMR data structure in amortized constant time.
//
// The zero iterator is a valid iterator at a position just before the first Node
// of an MMR.
type Iterator struct {
	n       int
	peaks   []int // indexes of peaks for an MMR of size n
	heights []int // heights of peaks for an MMR of size n
}

// Next returns the index and height (both counting from zero) of the
// next node in the iterator's sequence.
//
// Next runs in amortized constant time.
func (i *Iterator) Next() Node {
	var r Node
	r.Pos = i.n
	np := len(i.peaks)

	// Adding a non-leaf?
	if len(i.heights) >= 2 && i.heights[np-1] == i.heights[np-2] {
		r.Left, r.Right = i.peaks[np-2], i.peaks[np-1]
		i.peaks = i.peaks[:np-2]
		r.Height = i.heights[np-1] + 1
		i.heights = i.heights[:np-2]
	}

	// Record parent
	if len(i.heights) >= 1 && r.Height == i.heights[len(i.heights)-1] {
		r.Parent = r.Pos + 1 // from right child (next node)
	} else {
		r.Parent = r.Pos + pow2(r.Height+1) // from left child (inverse of left child below)
	}

	// Record children
	if r.Height > 0 {
		r.Left = r.Pos - pow2(r.Height)
		r.Right = r.Pos - 1
	}

	i.peaks = append(i.peaks, r.Pos)
	i.heights = append(i.heights, r.Height)
	i.n++

	return r
}

// Begin returns an iterator that is initialized to the beginning of
// an MMR.
//
// The first call to Next on the returned iterator will return the
// node at position 0.
//
// Begin is equivalent to calling IterJustBefore(0).
//
// Begin runs in constant time.
func Begin() *Iterator {
	return new(Iterator)
}

// IterAt returns an iterator at a point just before pos.
//
// The first call to Next on the returned iterator will return the
// node at position pos.
//
// IterJustBefore runs in log2(pos) time.
func IterJustBefore(pos int) *Iterator {
	peaks, heights := peaksAndHeights(pos)
	return &Iterator{
		n:       pos,
		peaks:   peaks,
		heights: heights,
	}
}
