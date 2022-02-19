package mmr

// This file tests basic functions of the MMR data structure
// independent of any specific implementation of MMRs.

import (
	"testing"

	"github.com/vsekhar/mmr/internal/bruteforce"
)

func TestPeaksAndHeights(t *testing.T) {
	sizes := []int{
		1, 3, 5, 13, 16, 32, 50, 900, 20000,
	}
sizes_loop:
	for i, s := range sizes {
		m := bruteforce.New(s)
		mathPeaks, mathHeights := peaksAndHeights(s)
		if len(m.Peaks) != len(mathPeaks) {
			t.Errorf("case %d: expected peaks (%v), got (%v)", i, m.Peaks, mathPeaks)
			continue
		}
		for i, p := range m.Peaks {
			if p.Index != mathPeaks[i] {
				t.Errorf("case %d: expected peaks (%v), got (%v)", i, m.Peaks, mathPeaks)
				continue sizes_loop
			}
			if p.Height != mathHeights[i] {
				t.Errorf("case %d: expected peak at %d to have height %d, got %d", i, p.Index, mathHeights[i], p.Height)
			}
		}
	}
}

func TestPath(t *testing.T) {
	// TODO
}
