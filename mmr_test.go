package mmr

import (
	"crypto/rand"
	"fmt"
	"testing"

	"golang.org/x/crypto/sha3"
)

type testArray [][]byte

func (t *testArray) Len() int {
	return len(*t)
}

func (t *testArray) HashAt(i int) []byte {
	h := sha3.Sum512((*t)[i])
	return h[:]
}

func TestMMR(t *testing.T) {
	array := &testArray{
		{0, 1, 2},
		[]byte("hello world"),
	}
	mmr := New(array, 2)
	for i := range *array {
		h := array.HashAt(i)
		if n, ok := mmr.GetIndex(h); !ok || n != i {
			t.Errorf("indexes don't match: %d from array, %d, %t from MMR", i, n, ok)
		}
	}
}

type byteArray [][]byte

func (b *byteArray) Len() int {
	return len(*b)
}

func (b *byteArray) HashAt(i int) []byte {
	r := sha3.Sum512((*b)[i])
	return r[:]
}

func BenchmarkBranching(b *testing.B) {
	bs := []int{
		2,
		3,
		4,
		8,
		10,
		13, // prime
		16,
		23, // prime
		32,
		47, // prime
		64,
		100,
		101, // prime
		128,
		256,
		257, // prime
		509, // prime
		512,
		1021, // prime
		1024,
		2048,
	}
	const arraySize = 2 << 11  // 4096
	const elementSize = 2 << 5 // 64
	array := new(byteArray)
	*array = make([][]byte, arraySize)
	for i := 0; i < arraySize; i++ {
		var buf [elementSize]byte
		n, err := rand.Read(buf[:])
		if n != elementSize {
			b.Fatal("short read")
		}
		if err != nil {
			b.Error(err)
		}
		(*array)[i] = buf[:]
	}
	// initial hashing
	for _, branching := range bs {
		b.Run(fmt.Sprintf("creation, b=%d", branching), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := New(array, branching)
				if len(m.(*mmr).hashes) != arraySize {
					b.Errorf("MMR hashes length %d, expected %d", len(m.(*mmr).hashes), arraySize)
				}
			}
		})
	}

	// proving
	for _, branching := range bs {
		m := New(array, branching)
		b.Run(fmt.Sprintf("proving, b=%d", branching), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				proofs := bs // use branching factor array as pos array for proofs
				for _, p := range proofs {
					_, _ = m.Proof(p)
				}
			}
		})
	}
}
