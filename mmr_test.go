package mmr

import (
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

func BenchmarkBranching(b *testing.B) {
	bs := []int{2, 3, 4, 8, 10, 16, 32, 64, 100, 128, 256, 512, 1024, 2048, 8192}
	// create Array, load with values
	for _, branching := range bs {
		// create MMR with specified branching
		b.Run(fmt.Sprintf("b=%d", branching), func(b *testing.B) {
			// extend MMR
		})
	}
}
