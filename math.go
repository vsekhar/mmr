package mmr

//  Algos: https://github.com/mimblewimble/grin/blob/master/core/src/core/pmmr/pmmr.rs

import (
	"fmt"
	"math"
	"math/bits"
)

// intPow computes x to the power of y exclusively using integers. If y is negative,
// intPow panics.
//
// This function runs in O(y) time in the general case and O(1) time if x is a power
// of 2.
func intPow(x, y int) int {
	switch {
	case y < 0:
		panic("intPow cannot raise an integer to a negative power")
	case y == 0:
		return 1
	case x == 2:
		return x << (y - 1)
	case x == 10:
		return int(math.Pow10(y))
	case x&(x-1) == 0:
		return x << ((bits.Len(uint(x)) - 1) * (y - 1))
	default:
		ret := 1
		for i := 0; i < y; i++ {
			ret *= x
		}
		return ret
	}
}

// intLog returns the largest integer smaller than or equal to log_b(x).
//
// This function runs in O(1) time.
func intLog(x, b int) int {
	switch {
	case x < 1:
		panic(fmt.Sprintf("x must be greater than 0, x=%d", x))
	case b < 1:
		panic(fmt.Sprintf("b must be greater than 0, b=%d", b))
	case b == 2:
		return bits.Len(uint(x)) - 1
	case b == 10:
		return int(math.Log10(float64(x)))
	case b&(b-1) == 0:
		if x&(x-1) != 0 {
			return bits.Len(uint(x)) / bits.Len(uint(b))
		}
		return int(math.Log2(float64(x))/float64(bits.Len(uint(b)))) + 1
	default:
		return int(math.Log(float64(x)) / math.Log(float64(b)))
	}
}

func firstChild(pos, h, b int) int {
	return pos - intPow(b, h)
}

// children returns the indexes of the children of the node at index pos and height h
// in an MMR with branching factor b. If the node at pos is a leaf, children returns nil.
func children(pos, h, b int) []int {
	switch {
	case h < 0:
		panic("height cannot be negative")
	case h == 0:
		return nil
	default:
		delta := siblingDelta(h-1, b)
		children := make([]int, 0, b)
		for c := firstChild(pos, h, b); c < pos; c += delta {
			children = append(children, c)
		}
		return children
	}
}

func rightSibling(pos, h, b int) int {
	return pos + siblingDelta(h, b)
}

func leftSibling(pos, h, b int) int {
	return pos - siblingDelta(h, b)
}

func siblingDelta(h, b int) int {
	return (intPow(b, h+1) - 1) / (b - 1)
}

// leftEdgePos returns the zero-based index at the left-most edge at zero-based
// height h of the MMR with branching factor b.
func leftEdgePos(h, b int) int {
	return siblingDelta(h, b) - 1
}

func parent(pos, h, b int) int {
	panic("not implemented")
}

// peaks returns the index of peaks in an MMR of size n and branching factor b.
//
// Source: https://github.com/mimblewimble/grin/blob/78220febeda94595159ece675e77e26986a3c11d/core/src/core/pmmr/pmmr.rs#L402
func peaks(n, b int) (peaks []int) {
	switch {
	case n < 1:
		panic("size of MMR must be greater than 0")
	case b < 2:
		panic("branching factor must be at least 2")
	case n < b:
		for i := 0; i < n; i++ {
			peaks = append(peaks, i)
		}
		return peaks
	case b == 2:
		p := 0 // partition (advances as we bag peaks)
		for n-p > 0 {
			nextPeak := intPow(b, (intLog(n-p+1, b))) - 1
			peaks = append(peaks, p+nextPeak-1)
			p += nextPeak
		}
		return peaks
	default:
		pos := n - 1
		p := 0
		for {
			cpos := pos - p
			h := intLog((cpos*(b-1)/b)+1, b) // height IF pos was on left edge
			i := leftEdgePos(h, b)           // pos IF pos was on left edge
			peaks = append(peaks, p+i)
			if i == cpos {
				return peaks // pos is indeed on left edge
			}
			s := (intPow(b, h+1) - 1) / (b - 1) // size of perfect tree to left
			p += s
		}
	}
}

// height returns the height (counting from 0) of the node at index n in the MMR with
// branching factor b.
func height(pos, b int) int {
	switch {
	case pos < 0:
		panic("index cannot be negative")
	case b < 2:
		panic("branching factor must be at least 2")
	case pos == 0:
		return 0
	case pos < b:
		return 0
	case b == 2: // fast path
		var upos = uint(pos)
		const allOnes = (1 << bits.UintSize) - 1
		var peakSize uint = allOnes >> bits.LeadingZeros(upos)
		for peakSize != 0 {
			if upos >= peakSize {
				upos -= peakSize
			}
			peakSize >>= 1
		}
		return int(upos)
	default:
		// Subtract away perfect trees to the left of pos until pos is on the
		// left edge. Nodes on left edge d at (zero-based) height h are sums:
		//  d(0)=0
		//  d(h)=d(h-1) + b^h
		//  d(h)=(b*(1-b^h))/(1-b)
		//  h = log_b((d*(1-b)/b) + 1)
		for pos > 0 {
			h := intLog((pos*(b-1)/b)+1, b) // height IF pos was on left edge
			i := leftEdgePos(h, b)          // pos IF pos was on left edge
			if i == pos {
				return h // pos is indeed on left edge, done
			}
			s := (intPow(b, h+1) - 1) / (b - 1) // size of perfect tree to left
			pos -= s
		}
		return 0
	}
}

func siblings(pos, h, b int) (pre, post []int) {
	panic("not implemented")
}

type pathEntry struct {
	pre  []int
	post []int
}

// path returns a sequence of nodes whose hashes comprise a proof of the node
// at pos.
func path(pos, b int) (p []pathEntry) {
	h := height(pos, b)
	cs := children(pos, h, b)
	if cs != nil {
		cpe := pathEntry{pre: make([]int, 0, len(cs))}
		for _, c := range cs {
			cpe.pre = append(cpe.pre, c)
		}
		p = append(p, cpe)
	}
	return nil
}
