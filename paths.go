package mmr

import (
	"fmt"
	"strings"
)

//go:generate stringer -type=OpCode

type OpCode int

const (
	UNSPECIFIED OpCode = iota
	PUSHPOS
	PUSHINPUT
	POP
)

// PathEntry is a single operation in a Path.
type PathEntry struct {
	op  OpCode
	pos int
}

// Path represents a path through an MMR that is used to prove inclusion and consistency.
type Path []PathEntry

func (p *Path) push(e PathEntry) { *p = append(*p, e) }
func (p *Path) pop() PathEntry   { var v PathEntry; v, *p = (*p)[len(*p)-1], (*p)[:len(*p)-1]; return v }

// String returns a string describing a Path.
func (p Path) String() string {
	s := make([]string, 0, len(p))
	for i, e := range p {
		s = append(s, fmt.Sprintf("%d:%s(%d)", i, e.op, e.pos))
	}
	return strings.Join(s, ", ")
}

func (p Path) Equals(p2 Path) bool {
	if len(p) != len(p2) {
		return false
	}
	for i := range p {
		if p[i] != p2[i] {
			return false
		}
	}
	return true
}

// evalWithValue executes and condenses Path p into the minimum Path corresponding
// to the digest of an MMR.
//
// The result of the evalWithValue method of Path can be compared to the result
// of the MMR Digest function to determine if the path constitutes a valid proof up to
// the digest of a given MMR.
//
// If pos is non-negative, it is a position value that must be used
// exactly once by p (e.g. for inclusion proofs).
//
// evalWithValue runs in log2(n)^2 where n is the size of the originating MMR.
func (p Path) evalWithValue(pos int) Path {
	if pos < 0 {
		panic("bad position value for inclusionDigest (use consistencyDigest?)")
	}
	return p.evalImpl(pos)
}

func (p Path) eval() Path {
	return p.evalImpl(-1)
}

func (p Path) evalImpl(pos int) Path {
	var s Path
	valUse := 0
	for _, o := range p {
		switch o.op {
		case PUSHPOS:
			s.push(PathEntry{op: PUSHPOS, pos: o.pos})
		case PUSHINPUT:
			if pos < 0 {
				panic("path expects a value, none provided")
			}
			s.push(PathEntry{op: PUSHPOS, pos: pos})
			valUse++
		case POP:
			s.pop()
			s.pop()
			s.push(PathEntry{op: PUSHPOS, pos: o.pos})
		default:
			panic(fmt.Sprintf("path: bad instruction in condense (%s)", o.op))
		}
	}
	if pos >= 0 && valUse != 1 {
		panic(fmt.Sprintf("pos used: %d times, expected 1", valUse))
	}
	return s
}
