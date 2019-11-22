// Package mmr is an implementation of Merkle Mountain Ranges.
//
// See:
//  https://github.com/mimblewimble/grin/blob/master/doc/mmr.md
//  https://github.com/opentimestamps/opentimestamps-server/blob/master/doc/merkle-mountain-range.md
//  Algos: https://github.com/mimblewimble/grin/blob/master/core/src/core/pmmr/pmmr.rs
package mmr

import (
	"hash"
	"math"
	"math/bits"
)

// Interface provides methods specific to querying an MMR.
//
// An MMR is created by passing an underlying Array to NewMMR. New values can be added
// to the MMR by adding them to the underlying array.
//
// Modifying or deleting a value from the underlying Array will likely corrupt an MMR
// built on top of it, and methods of the MMR may subsequently panic.
//
// The methods of Interface are not goroutine safe.
type Interface interface {
	// GetIndex returns the index at which the value with hash h can be found. If the
	// hash cannot be found, ok is false.
	GetIndex(h []byte) (i int, ok bool)

	// Digest returns a hash that summarizes the MMR in its current state.
	Digest() []byte

	// Prove provides a set of hashes that prove the inclusion of the value at index i
	// or nil if such a proof cannot be constructed.
	//
	// prefix comprises the hashes of the children of record i, if any. suffix comprises
	// the hashes of nodes from record i to a peak as well as all subsequent peaks.
	// digest is the hash of the tree overall (equivalent to calling Digest).
	//
	// The caller can verify the proof by writing the data from prefix into a hash function,
	// followed by the data at record i, followed by suffix. The result should equal the
	// digest.
	//
	Proof(i int) (prefix []byte, suffix []byte, digest []byte)
}

// An Array is any user-provided ordered container with the accessors needed to build it
// into an MMR.
//
// To add a value to an MMR, add it to its underlying Array.
//
// Modifying or deleting a value from the underlying Array will likely corrupt an MMR
// built on top of it, and methods of the MMR may subsequently panic.
//
// Methods of Array need not be goroutine safe since the methods of Interface are not
// goroutine safe.
type Array interface {
	// Len returns the length of the Array.
	Len() int

	// HashAt hashes the bytes in prefix (if any), and then the value at index i.
	HashAt(prefix []byte, i int) []byte
}

type mmr struct {
	array     Array
	hashes    [][]byte
	indexes   map[string]int // need to convert []byte hashes to string
	hasher    hash.Hash
	branching int
}

// New returns a new MMR constructed from Array a using Hash h with branching factor b.
// If a is nil or b is less than two, New panics.
func New(a Array, h hash.Hash, b int) Interface {
	if a == nil {
		panic("array required")
	}
	if b < 2 {
		panic("branching factor must be greater than 2")
	}
	ret := &mmr{
		array:     a,
		hashes:    make([][]byte, a.Len()),
		indexes:   make(map[string]int),
		hasher:    h,
		branching: b,
	}
	ret.extend()
	return ret
}

func (m *mmr) extend() int {
	if m.array.Len() < len(m.hashes) {
		panic("array length decreased; MMR does not support deletion")
	}
	s := m.Len()
	i := s
	for ; i < m.array.Len(); i++ {
		h := height(i, m.branching)
		prefix := []byte{}
		if h > 0 {
		}
		hash := m.array.HashAt(prefix, i)
		m.hashes = append(m.hashes, hash)
		m.indexes[string(hash)] = i
	}
	return i - s
}

func (m *mmr) Len() int {
	return m.array.Len()
}

func (m *mmr) GetIndex(hash []byte) (i int, ok bool) {
	m.extend()
	i, ok = m.indexes[string(hash)]
	return
}

func (m *mmr) Digest() []byte {
	m.extend()
	ps := peaks(m.Len(), m.branching)
	ret := make([]byte, 0)
	h := []byte{}
	for _, p := range ps {
		h = m.array.HashAt(h, p)
	}
	// TODO: bag peaks
	return ret
}

// childPrefix returns the concatenation of the hashes of the children of a node
// at pos, or an empty slice if the node has no children.
func (m *mmr) childPrefix(pos int) []byte {
	h := height(pos, m.branching)
	if h > 0 {
		var prefix []byte
		for i := firstChild(pos, h, m.branching); i < pos; i++ {
			prefix = append(prefix, m.hashes[i]...)
		}
		return prefix
	}
	return []byte{}
}

func firstChild(pos, h, b int) int {
	return pos - intPow(b, h)
}

func parent(pos, h, b int) int {
	panic("not implemented")
}

func (m *mmr) Proof(i int) (prefix []byte, suffix []byte, digest []byte) {
	m.extend()
	prefix = m.childPrefix(i)
	_ = height(i, m.branching)
	// figure out firstChild
	// parent := firstChild + intPow(m.branching, h+1))

	digest = m.Digest()
	return nil, nil, nil
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
