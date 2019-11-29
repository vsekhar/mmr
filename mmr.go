// Package mmr is an implementation of Merkle Mountain Ranges.
//
// See:
//  https://github.com/mimblewimble/grin/blob/master/doc/mmr.md
//  https://github.com/opentimestamps/opentimestamps-server/blob/master/doc/merkle-mountain-range.md
package mmr

import (
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/sha3"
)

const (
	hashLengthBytes = 64
)

type pathHashes struct {
	pre  [][]byte // record hashes of siblings to the left
	post [][]byte // record hashes of siblings to the right
}

type digest struct {
	n    int
	hash []byte
}

type proof struct {
	// path is a sequence of sibling record hashes used to construct successive parent
	// record hashes all the way up to a node's peak. If path is empty, the node is
	// itself a peak.
	path []pathHashes

	// rightPeaks is a sequence of peak record hashes to the right of the peak reached
	// via path. It is ordered from right to left. If rightPeaks is empty, the node is
	// the rightmost (newest) node which is both a leaf and an peak.
	rightPeaks [][]byte

	// leftPeaks is a sequence of peak record hashes to the left of the peak reached
	// via path. It is ordered from right to left. If leftPeaks is empty, the node is
	// the leftmost (largest) peak.
	leftPeaks [][]byte

	// digest is the digest of the full tree at the time this proof was generated.
	digest digest
}

// TODO: proving digest a is in digest b, by proving that each of the peaks of digest a
// is in digest b. I.e. a digest proof is a []proof and validates if all the constituent
// proofs validate.

// Interface provides methods specific to querying an MMR.
//
// An MMR is created by passing an underlying Array to New. New values can be added
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

	// AuditProof provides a set of hashes that prove the inclusion of the value
	// at index i or nils if such a proof cannot be constructed.
	//
	// Verifing the proof comprises the following:
	//
	//  1) hash each array in sequence[0]
	//  2) hash the timestamp of record i provided when the record was added
	//  3) hash the data of record i
	//  4) hash the salt of record i
	//  5) hash each array in sequence[1], and store this hash
	//  6) hash each array in sequence[2]
	//  7) hash the value stored in step 5
	//  8) hash each array in sequence [3]
	//
	// The result should equal digest.
	AuditProof(i int) (sequence [][]byte, digest []byte)

	// ConsistencyProof TBD
	//
	//  - prove digest to digest
	//  - do we need the digests? just the sizes? one or the other?
	ConsistencyProof(n int)
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

	// HashAt returns a hash of element i. Any hash function producing an array of bytes may
	// be used.
	HashAt(i int) []byte
}

// Timestamping and salting is a function for the Array. I.e. adding to the Array should
// involve hashing those additional values. MMR is only concerned with the overall resulting
// hash of data.

type hashSet struct {
	ofData []byte // from Array.HashAt()
	ofNode []byte // of children in order from first to last (if any) followed by hash ofData.
}

type mmr struct {
	array     Array
	hashes    []hashSet
	index     map[string]int // key is base58 encoded hash of data
	hasher    sha3.ShakeHash
	branching int
}

// New returns a new MMR constructed from Array a using Hash h with branching factor b.
// If a is nil or b is less than two, New panics.
func New(a Array, b int) Interface {
	if a == nil {
		panic("array required")
	}
	if b < 2 {
		panic("branching factor must be greater than 2")
	}
	ret := &mmr{
		array:     a,
		hashes:    make([]hashSet, 0, a.Len()),
		index:     make(map[string]int),
		hasher:    sha3.NewShake256(),
		branching: b,
	}
	return ret
}

func (m *mmr) nodeHash(children []int, dataHash []byte) []byte {
	m.hasher.Reset()
	for _, i := range children {
		m.hasher.Write(m.hashes[i].ofNode)
	}
	m.hasher.Write(dataHash)
	ret := new([hashLengthBytes]byte)[:]
	n, err := m.hasher.Read(ret)
	if n != hashLengthBytes {
		panic("short read from hasher")
	}
	if err != nil {
		panic(err)
	}
	return ret
}

// extend updates the MMR's data structure to cover elements [0, n), if necessary,
// and returns the number of elements that were hashed.
func (m *mmr) extend(n int) int {
	if n < len(m.hashes) {
		return 0
	}
	if m.array.Len() < n {
		n = m.array.Len()
	}

	delta := n - len(m.hashes)
	for i := len(m.hashes); i < n; i++ {
		var hash hashSet
		hash.ofData = m.array.HashAt(i)
		m.index[base58.Encode(hash.ofData)] = i
		cs := children(i, height(i, m.branching), m.branching)
		hash.ofNode = m.nodeHash(cs, hash.ofData)
		m.hashes = append(m.hashes, hash)
	}
	return delta
}

func (m *mmr) Len() int {
	return m.array.Len()
}

func (m *mmr) GetIndex(hash []byte) (i int, ok bool) {
	// already hashed?
	if i, ok = m.index[base58.Encode(hash)]; ok {
		return
	}
	// complete search
	m.extend(m.Len())
	i, ok = m.index[base58.Encode(hash)]
	return
}

func (m *mmr) Digest(n int) []byte {
	m.extend(n)
	// ps := peaks(m.Len(), m.branching)
	ret := make([]byte, 0)
	// TODO: bag peaks
	return ret
}

func (m *mmr) AuditProof(i int) (sequence [][]byte, digest []byte) {
	// extend hashes up to the newest peak needed for the proof
	ps := peaks(i, m.branching)
	maxPeak := 0
	for _, p := range ps {
		if p > maxPeak {
			maxPeak = p
		}
	}
	m.extend(maxPeak + 1)

	// TODO: compute proof

	panic("not implemented")
}

func (m *mmr) ConsistencyProof(n int) {
	panic("not implemented")
}
