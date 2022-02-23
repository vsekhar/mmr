package mmr_test

import (
	"testing"

	"github.com/vsekhar/mmr"
	"github.com/vsekhar/mmr/internal/testdata"
)

func TestFirstNode(t *testing.T) {
	zero := mmr.MMR{}
	if n1 := mmr.Advance(&zero); n1 != testdata.FirstNode {
		t.Errorf("zero MMR does not start with first node, got %+v", n1)
	}
	new0 := mmr.New(0)
	if n3 := mmr.Advance(&new0); n3 != testdata.FirstNode {
		t.Errorf("New(0) does not start with first node, got node at pos %d", n3.Pos)
	}
}
func TestAdvance(t *testing.T) {
	m := mmr.MMR{}
	for i, c := range testdata.Sequence {
		node := mmr.Advance(&m)
		if c != node {
			t.Errorf("pos %d: expected {%v}, got {%v}", i, c, node)
		}
	}
}

func TestAt(t *testing.T) {
	for i, n := range testdata.Sequence {
		n2 := mmr.At(n.Pos)
		if n != n2 {
			t.Errorf("pos %d: expected %+v, got %+v", i, n, n2)
		}
	}
}
