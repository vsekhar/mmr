package mmr

func leftChild(pos, h int) int  { return pos - pow2(h) }
func rightChild(pos, h int) int { return pos - 1 }

type Iterator struct {
	n       int
	peaks   []int // indexes of peaks for an MMR of size n
	heights []int // heights of peaks for an MMR of size n
}

// Next returns the index and height (both counting from zero) of the
// next node in the iterator's sequence.
//
// Next runs in constant time.
func (i *Iterator) Next() *Node {
	r := new(Node)
	r.Pos = i.n
	np := len(i.peaks)
	if len(i.heights) >= 2 && i.heights[np-1] == i.heights[np-2] {
		r.Left, r.Right = i.peaks[np-2], i.peaks[np-1]
		i.peaks = i.peaks[:np-2]
		i.peaks = append(i.peaks, r.Pos)
		r.Height = i.heights[np-1] + 1
		i.heights = i.heights[:np-2]
		i.heights = append(i.heights, r.Height)
	} else {
		i.peaks = append(i.peaks, r.Pos)
		r.Height = 0
		i.heights = append(i.heights, r.Height)
	}
	// Will the next iteration create a new peak, then that is this node's parent
	if len(i.heights) >= 2 && i.heights[len(i.heights)-1] == i.heights[len(i.heights)-2] {
		// Next node will be this node's parent
		r.Parent = r.Pos + 1
	} else {
		r.Parent = r.Pos + pow2(r.Height+1) // inverse of leftChild
	}
	i.n++

	if r.Height > 0 {
		r.HasChildren = true
		r.Left = leftChild(r.Pos, r.Height)
		r.Right = rightChild(r.Pos, r.Height)
	}
	return r
}

// Begin returns an iterator that is initialized to the beginning of
// an MMR.
//
// The first call to Next on the returned iterator will return the
// node at position 0.
//
// Begin is equivalent to calling IterJustBefore(0).
func Begin() *Iterator { return new(Iterator) }

// IterAt returns an iterator at a point just before pos.
//
// The first call to Next on the returned iterator will return the
// node at position pos.
func IterJustBefore(pos int) *Iterator {
	peaks, heights := peaksAndHeights(pos)
	return &Iterator{
		n:       pos,
		peaks:   peaks,
		heights: heights,
	}
}
