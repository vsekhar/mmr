package mmr

import "fmt"

// Consistency produces two paths that can be used to show that a larger
// MMR is consistent with (i.e. includes) a smaller MMR.
//
// p1 produces the digest of the smaller MMR of size n1.
//
// The combined path `p = append(p1, p2...)` produces the digest of the
// larger MMR of size n2.
//
// If n2 > n1, Consistency panics.
//
// Consistency runs in log2(n2) time.
func Consistency(n1, n2 int) (p1, p2 path) {
	if n1 > n2 {
		panic(fmt.Sprintf("consistency cannot be shown from MMR of size %d to MMR of size %d", n1, n2))
	}
	// p1 pushes peaks of MMR @ n1, producing its digest
	peaks1, _ := peaksAndHeights(n1)
	for _, p := range peaks1 {
		p1 = append(p1, pathEntry{op: PUSHPOS, pos: p})
	}

	// p2 can be appended to p1 to produce digest of MMR @ n2
	//
	// Proving inclusion of n1 in n2 is appending to p1 inclusion proof
	// of last peak of n1 in the appropriate peak of n2?
	lastPeak := peaks1[len(peaks1)-1]
	peaks2, _ := peaksAndHeights(n2)
	found := false
	for _, p := range peaks2 {
		if !found {
			if p < lastPeak {
				continue
			}
			found = true
			node := At(p)
			p2 = inclusionImpl(lastPeak, node, p2)
		} else {
			p2 = append(p2, pathEntry{op: PUSHPOS, pos: p})
		}
	}

	// drop instructions in p2 up to and including the PUSHVAL instruction
	// (these will be covered in p1)
	p2Start := 0
	for ; p2Start < len(p2); p2Start++ {
		if p2[p2Start].op == PUSHVAL {
			p2Start++
			break
		}
	}
	p2 = p2[p2Start:]

	return p1, p2
}
