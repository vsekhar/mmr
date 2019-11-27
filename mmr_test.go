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
	bs := []int{2, 3, 4, 8, 10, 16, 32, 64, 100, 128, 256, 512, 1024, 2048, 8192}
	const arraySize = 2 << 1   // 4096
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
	for _, branching := range bs {
		mmr := New(array, branching)
		// create MMR with specified branching
		b.Run(fmt.Sprintf("b=%d", branching), func(b *testing.B) {
			_, _ = mmr.Proof(arraySize - 1)
		})
	}
}
