package mmr

import (
	"fmt"
	"strings"
	"testing"
)

//go:generate stringer -type=opcode

type opcode int

const (
	UNSPECIFIED opcode = iota
	PUSHPOS
	PUSHVAL
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

type stack []int

func (s *stack) pop() int   { var v int; v, *s = (*s)[len(*s)-1], (*s)[:len(*s)-1]; return v }
func (s *stack) push(x int) { *s = append(*s, x) }

// matchesDigest returns true if the positions in s (in pop order)
// match the positions of the peaks of an MMR of size n.
//
// matchesDigest runs in log2(n) time.
func (s stack) matchesDigest(n int) bool {
	var peakStack stack
	peaks, _ := peaksAndHeights(n)
	if len(s) != len(peaks) {
		return false
	}
	for _, p := range peaks {
		peakStack.push(p)
	}
	for i := range s {
		if s[i] != peakStack[i] {
			return false
		}
	}
	return true
}

func run(t *testing.T, p path, pos, n int) stack {
	var s stack

	valUse := 0
	for ip, o := range p {
		switch o.op {
		case PUSHPOS:
			s.push(o.pos)
		case PUSHVAL:
			s.push(pos)
			valUse++
		case POP:
			n1, n2 := At(s.pop()), At(s.pop())
			if n1.Height != n2.Height {
				t.Errorf("%d: different heights", ip)
			}
			if n1.Parent != n2.Parent {
				t.Errorf("%d: different parents", ip)
			}
			s.push(n2.Parent)
		}
	}

	if pos >= 0 && valUse != 1 {
		t.Errorf("valUse: %d, expected 1", valUse)
	}

	return s
}
