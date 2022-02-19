package mmr

//  Algos: https://github.com/mimblewimble/grin/blob/master/core/src/core/pmmr/pmmr.rs

import (
	"math/bits"
)

func pow2(x int) int {
	if x == 0 {
		return 1
	} else {
		return 2 << (x - 1)
	}
}
func log2(x int) int { return bits.Len(uint(x)) - 1 }

// peaksAndHeights returns the indexes of peaks and their heights for
// an an MMR of size n.
//
// peaksAndHeights runs in log2(n) time.
func peaksAndHeights(n int) (peaks, heights []int) {
	// Source (for peaks calculation):
	// https://github.com/mimblewimble/grin/blob/78220febeda94595159ece675e77e26986a3c11d/core/src/core/pmmr/pmmr.rs#L402

	pos := 0
	for n-pos > 0 {
		peakHeight := log2(n - pos + 1)
		peakSize := pow2(peakHeight) - 1
		peaks = append(peaks, pos+peakSize-1)
		heights = append(heights, peakHeight-1)
		pos += peakSize
	}
	return peaks, heights
}
