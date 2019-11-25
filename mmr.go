// Package mmr is an implementation of Merkle Mountain Ranges.
//
// See:
//  https://github.com/mimblewimble/grin/blob/master/doc/mmr.md
//  https://github.com/opentimestamps/opentimestamps-server/blob/master/doc/merkle-mountain-range.md
package mmr

import (
	"golang.org/x/crypto/sha3"
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

	// Digest returns a hash that summarizes the MMR as it was when it had size n.
	// Digest panics if n is greater than the length of the underlying Array.
	Digest(n int) []byte

	// Prove provides a set of hashes that prove the inclusion of the value at index i
	// or nils if such a proof cannot be constructed.
	//
	// Verifing the proof comprises the following:
	//
	//  1) hash each array in sequence[0]
	//  2) hash the timestamp of record i
	//  3) hash the data of record i
	//  4) hash the salt of record i
	//  5) hash each array in sequence[1], and store this hash
	//  6) hash each array in sequence[2]
	//  7) hash the value stored in step 5
	//  8) hash each array in sequence [3]
	//
	// The result should equal digest.
	Proof(i int) (sequence [][]byte, digest []byte)
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
	HashAt(i int) []byte
}

type mmr struct {
	array     Array
	hashes    [][]byte
	indexes   map[string]int // need to convert []byte hashes to string
	hasher    sha3.ShakeHash
	branching int
}

// New returns a new MMR constructed from Array a using Hash h with branching factor b.
// If a is nil or b is less than two, New panics.
func New(a Array, h sha3.ShakeHash, b int) Interface {
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
		// prefix := []byte{}
		if h > 0 {
			// TODO
		}
		hash := m.array.HashAt(i)
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

func (m *mmr) Digest(n int) []byte {
	m.extend()
	// ps := peaks(m.Len(), m.branching)
	ret := make([]byte, 0)
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

func (m *mmr) Proof(i int) (sequence [][]byte, digest []byte) {
	m.extend()
	return nil, nil
}
