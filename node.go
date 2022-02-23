package mmr

// Node records the position of a node, its height, its parents and its
// children (if any).
//
// The zero Node is invalid.
type Node struct {
	Pos    int
	Height int
}

func (n Node) HasChildren() bool { return n.Height > 0 }

func (n Node) child(childPos int) Node {
	if !(n.Height > 0) {
		panic("leaf has no child")
	}
	r := Node{
		Pos:    childPos,
		Height: n.Height - 1,
	}
	return r
}

// LeftChild returns the Node corresponding to the left child of n.
//
// If n is a leaf, LeftChild panics.
//
// LeftChild runs in constant time.
func (n Node) LeftChild() Node {
	return n.child(leftChild(n.Pos, n.Height))
}

// RightChild returns the Node corresponding to the right child of n
//
// If n is a leaf, RightChild panics.
//
// RightChild runs in constant time.
func (n Node) RightChild() Node {
	return n.child(rightChild(n.Pos, n.Height))
}

// At returns the Node at position pos.
//
// At runs in log2(pos) time.
func At(pos int) Node {
	m := New(pos)
	return Advance(&m)
}
