package mmr

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
	if n == 0 || b == 0 {
		panic("n and b must be greater than zero")
	}
	fn, fb := float64(n), float64(b)
	var l float64
	switch b {
	case 2:
		l = math.Log2(fn)
	case 10:
		l = math.Log10(fn)
	default:
		l = math.Log(fn) / math.Log(fb)
	}
	return int(math.Floor(l))
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
	var pos = uint64(n)
	if pos == 0 {
		return 0
	}
	if b == 2 {
		// bit-shifting fast-path
		var peakSize uint64 = math.MaxUint64 >> bits.LeadingZeros64(pos)
		var bitmap uint64
		for peakSize != 0 {
			bitmap <<= 1
			if pos >= peakSize {
				pos -= peakSize
				bitmap |= 1
			}
			peakSize >>= 1
		}
		return int(pos)
	}
	panic("branching factors other than 2 are not implemented")
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
	for i := 0; i < y; i++ {
		ret *= x
	}
	return ret
}
