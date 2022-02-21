package mmr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vsekhar/mmr/internal/inthash"
)

//go:generate stringer -type=opcode

type opcode int

const (
	UNSPECIFIED opcode = iota
	PUSHPOS
	PUSHVAL
	PUSHEMPTY
	POP
)

type pathEntry struct {
	op  opcode
	pos int
}

type path []pathEntry

func (p path) String() string {
	s := make([]string, 0, len(p))
	for i, e := range p {
		s = append(s, fmt.Sprintf("%d:%s(%d)", i, e.op, e.pos))
	}
	return strings.Join(s, ", ")
}

// Inclusion returns a path corresponding to the nodes needed to prove inclusion
// of the node at pos in an MMR of size n.
//
// Inclusion runs in log2(n) time.
func Inclusion(pos, n int) path {
	if pos >= n {
		panic(fmt.Sprintf("pos %d does not exist in MMR of size %d", pos, n))
	}
	var r path
	peaks, _ := peaksAndHeights(n)
	found := false
	for _, p := range peaks {
		if !found {
			if p < pos {
				r = append(r, pathEntry{op: PUSHPOS, pos: p})
				continue
			}
			found = true
			node := At(p)
			r = inclusionImpl(pos, node, r)
		} else {
			r = append(r, pathEntry{op: PUSHPOS, pos: p})
		}
	}
	r = append(r, pathEntry{op: PUSHEMPTY})

	for i := 0; i < len(peaks); i++ {
		r = append(r, pathEntry{op: POP})
	}

	return r
}

func inclusionImpl(pos int, node Node, p path) path {
	if node.Pos == pos {
		p = append(p, pathEntry{op: PUSHVAL})
		return p
	}
	l, r := node.LeftChild(), node.RightChild()
	if l.Pos >= pos {
		p = inclusionImpl(pos, l, p)
		p = append(p, pathEntry{op: PUSHPOS, pos: r.Pos})
		p = append(p, pathEntry{op: POP})
		return p
	}
	p = append(p, pathEntry{op: PUSHPOS, pos: l.Pos})
	p = inclusionImpl(pos, r, p)
	p = append(p, pathEntry{op: POP})
	return p
}

// TODO: consistency proof of an MMR of size n with an MMR of size m > n.
func Consistency(n1, n2 int) path {
	panic("unimplemented")
}

type stack []int

func (s *stack) pop() int   { var v int; v, *s = (*s)[len(*s)-1], (*s)[:len(*s)-1]; return v }
func (s *stack) push(x int) { *s = append(*s, x) }

func run(t *testing.T, p path, pos, n int) {
	var s stack
	var peakStack stack
	peaks, _ := peaksAndHeights(n)
	for _, p := range peaks {
		peakStack.push(p)
	}

	// Construct the digest (intHashing position numbers)
	digest := inthash.New()
	for len(peakStack) > 0 {
		digest.Write(peakStack.pop())
	}
	if digest.Sum() == 0 {
		t.Errorf("bad hash") // sanity check
	}

	computedDigest := inthash.New()
	computingDigest := false
	for ip, o := range p {
		switch o.op {
		case PUSHPOS:
			s.push(o.pos)
		case PUSHVAL:
			s.push(pos)
		case PUSHEMPTY:
			s.push(-1)
		case POP:
			n1 := At(s.pop())
			n2 := At(s.pop())
			if !computingDigest {
				if n1.Pos != -1 {
					// Interior node, they should be well-formed
					if n1.Height != n2.Height {
						t.Errorf("%d: different heights", ip)
					}
					if n1.Parent != n2.Parent {
						t.Errorf("%d: different parents", ip)
					}
					s.push(n2.Parent)
				} else {
					// Start computing digest
					if n2.Pos < 0 {
						t.Errorf("%d: two virtual nodes", ip)
					}
					computedDigest.Write(n2.Pos)
					computingDigest = true
					s.push(-1)
				}
			} else {
				// Continue computing digest
				if n1.Pos != -1 {
					t.Errorf("%d: unexpected non-virtual node", ip)
				}
				if n2.Pos < 0 {
					t.Errorf("%d: unexpected virtual node", ip)
				}
				computedDigest.Write(n2.Pos)
				s.push(-1)
			}

		}
	}

	if len(s) != 1 {
		t.Errorf("residual stack entries: expected 1, got %d", len(s))
	}

	if digest.Sum() != computedDigest.Sum() {
		t.Errorf("digest mismatch: got %d, expected %d", computedDigest.Sum(), digest.Sum())
	}

	return
}
