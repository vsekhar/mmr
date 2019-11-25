package mmr

//  Algos: https://github.com/mimblewimble/grin/blob/master/core/src/core/pmmr/pmmr.rs

import (
	"math"
	"math/bits"
)

// intPow computes x to the power of y exclusively using integers. If y is negative,
// intPow panics.
func intPow(x, y int) int {
	switch {
	case y < 0:
		panic("intPow cannot raise an integer to a negative power")
	case y == 0:
		return 1
	case x == 2:
		return x << (y - 1)
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

// intLog returns the largest integer smaller than or equal to log_b(n).
func intLog(x, b int) int {
	switch {
	case x < 1:
		panic("n must be greater than 0")
	case b < 1:
		panic("b must be greater than zero")
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

func rightSibling(pos, h, b int) int {
	return pos + (intPow(b, h+1)-1)/(b-1)
}

func leftSibling(pos, h, b int) int {
	return pos - (intPow(b, h+1)-1)/(b-1)
}

func parent(pos, h, b int) int {
	panic("not implemented")
}

// peaks returns the index of peaks in an MMR of size n and branching factor b.
//
// Source: https://github.com/mimblewimble/grin/blob/78220febeda94595159ece675e77e26986a3c11d/core/src/core/pmmr/pmmr.rs#L402
func peaks(n, b int) []int {
	switch {
	case n < 1:
		panic("size of MMR must be greater than 0")
	case b < 2:
		panic("branching factor must be at least 2")
	case b == 2:
		var peaks []int
		p := 0 // partition (advances as we bag peaks)
		for n-p > 0 {
			nextPeak := intPow(b, (intLog(n-p+1, b))) - 1
			peaks = append(peaks, p+nextPeak-1)
			p += nextPeak
		}
		return peaks
	default:
		// Possibly helpful:
		// https://ece.uwaterloo.ca/~dwharder/aads/Lecture_materials/5.04.N-ary_trees.pdf
		panic("branching factors other than 2 are not implemented")
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
	case b == 2:
		var upos = uint(pos)
		// bit-shifting fast-path
		// peakSize := uint(intPow(b, intLog(n, b)+1)) - 1
		// optimised for b=2
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
		panic("branching factors other than 2 are not implemented")
	}
}
