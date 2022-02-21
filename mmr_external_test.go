package mmr_test

import (
	"testing"

	"github.com/vsekhar/mmr"
	"github.com/vsekhar/mmr/internal/testdata"
)

func TestFirstNode(t *testing.T) {
	zero := &mmr.Iterator{}
	if n1 := zero.Next(); n1 != testdata.FirstNode {
		t.Errorf("zero iterator does not start with first node, got %+v", n1)
	}
	if n2 := mmr.Begin().Next(); n2 != testdata.FirstNode {
		t.Errorf("Begin iterator does not start with first node, got %+v", n2)
	}
	if n3 := mmr.IterJustBefore(0).Next(); n3 != testdata.FirstNode {
		t.Errorf("IterJustBefore(0) does not start with first node, got ")
	}
}
func TestIterator(t *testing.T) {
	itr := mmr.Begin()
	for i, c := range testdata.Sequence {
		node := itr.Next()
		if c != node {
			t.Errorf("pos %d: expected {%v}, got {%v}", i, c, node)
		}
	}
}

func TestNodeAt(t *testing.T) {
	for i, n := range testdata.Sequence {
		n2 := mmr.At(n.Pos)
		if n != n2 {
			t.Errorf("pos %d: expected %+v, got %+v", i, n, n2)
		}
	}
}
