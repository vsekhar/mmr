package mmr

import (
	"testing"

	"github.com/vsekhar/mmr/internal/bruteforce"
)

func TestIterator(t *testing.T) {
	itr := Begin()
	for i, c := range Sequence {
		node := itr.Next()
		if c != *node {
			t.Errorf("pos %d: expected {%v}, got {%v}", i, c, node)
		}
	}
}

func nodeEquals(n *Node, bn *bruteforce.Node) bool {
	if n == nil || bn == nil {
		return false
	}
	if n.Pos != bn.Index {
		return false
	}
	if n.Height != bn.Height {
		return false
	}
	if n.HasChildren {
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
		if !nodeEquals(node, bruteNode) {
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
		if !nodeEquals(node, bruteNode) {
			t.Errorf("pos %d: got (%v), expected (%v)", i, node, bruteNode)
		}
	}
}
