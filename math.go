package mmr

//  Algos: https://github.com/mimblewimble/grin/blob/master/core/src/core/pmmr/pmmr.rs

import (
	"math"
	"math/bits"
)

func firstChild(pos, h, b int) int {
	return pos - intPow(b, h)
}

func rightSibling(pos, h, b int) int {
	return pos + intPow(b, h+1) - 1
}

func leftSibling(pos, h, b int) int {
	return pos - intPow(b, h+1) + 1
}

func parent(pos, h, b int) int {
	panic("not implemented")
}

// peaks returns the index of peaks in an MMR of size n and branching factor b.
//
// Source: https://github.com/mimblewimble/grin/blob/78220febeda94595159ece675e77e26986a3c11d/core/src/core/pmmr/pmmr.rs#L402
func peaks(n, b int) []int {
	if n < 1 {
		panic("size of MMR must be positive")
	}
	if b < 2 {
		panic("branching factor must be at least 2")
	}
	var peaks []int
	p := 0 // partition (advances as we bag peaks)
	for n-p > 0 {
		nextPeak := intPow(b, (intLog(n-p+1, b))) - 1
		peaks = append(peaks, p+nextPeak-1)
		p += nextPeak
	}
	return peaks
}

// intLog returns the largest integer smaller than or equal to log_b(n).
func intLog(n, b int) int {
	if n < 1 || b < 1 {
		panic("n and b must be greater than zero")
	}
	if b == 2 {
		return bits.Len(uint(n)) - 1
	}
	fn, fb := float64(n), float64(b)
	if b&(b-1) == 0 {
		return int(math.Log2(fn) / math.Log2(fb))
	}
	if b == 10 {
		return int(math.Log10(fn))
	}
	return int(math.Log(fn) / math.Log(fb))
}

// height returns the height (counting from 0) of the node at index n in the MMR with
// branching factor b.
func height(n, b int) int {
	if n < 0 {
		panic("index cannot be negative")
	}
	if b < 2 {
		panic("branching factor must be at least 2")
	}
	if n == 0 {
		return 0
	}
	if b == 2 {
		// bit-shifting fast-path
		var un = uint64(n)
		var peakSize uint64 = math.MaxUint64 >> bits.LeadingZeros64(un)
		for peakSize != 0 {
			if un >= peakSize {
				un -= peakSize
			}
			peakSize >>= 1
		}
		return int(un)
	}
	// 6 = 110
	// 7 = 111
	// N = b^{h+1} - 1
	// N + 1 = b^{h+1}
	// log_b(N+1) = h + 1
	// h = log_b(N+1) - 1
	var peakSize = intPow(b, intLog(n+1, b)) - 1
	for peakSize != 0 {
		if n >= peakSize {
			n -= peakSize
		}
		peakSize /= b
	}
	return n
}

// intPow computes x to the power of y exclusively using integers. If y is negative,
// intPow panics.
func intPow(x, y int) int {
	if y < 0 {
		panic("intPow cannot raise an integer to a negative power")
	}
	ret := 1
	if x == 2 {
		return x << (y - 1)
	}
	if x&(x-1) == 0 {
		return x << ((bits.Len(uint(x)) - 1) * (y - 1))
	}
	for i := 0; i < y; i++ {
		ret *= x
	}
	return ret
}
