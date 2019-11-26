package mmr

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestReadSlices(t *testing.T) {
	slices := [][]byte{
		{0, 1, 2},
		[]byte("hello world"),
	}
	var flat []byte
	for _, s := range slices {
		flat = append(flat, s...)
	}
	r, err := ioutil.ReadAll(ReadSlices(slices))
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(flat, r) {
		t.Errorf("expected '%v', got '%v'", flat, r)
	}
}
