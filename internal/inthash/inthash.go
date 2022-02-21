package inthash

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
)

type IntHash struct {
	h   hash.Hash64
	buf []byte
}

func New() *IntHash {
	return &IntHash{h: fnv.New64(), buf: make([]byte, 8)}
}

func (i *IntHash) Write(x int) {
	if x < 0 {
		panic("inthash: write only accepts non-negative values")
	}
	i.h.Reset()
	i.h.Write(i.buf)
	binary.LittleEndian.PutUint64(i.buf, uint64(x))
	i.h.Write(i.buf)
	i.buf = i.h.Sum(i.buf[:0])
}

func (i *IntHash) Sum() int {
	return int(binary.LittleEndian.Uint64(i.buf))
}
