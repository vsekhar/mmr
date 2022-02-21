package mmr

import "fmt"

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
	} else {
		p = append(p, pathEntry{op: PUSHPOS, pos: l.Pos})
		p = inclusionImpl(pos, r, p)
	}
	p = append(p, pathEntry{op: POP})
	return p
}
