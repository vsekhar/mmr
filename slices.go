package mmr

import "io"

type sliceReader struct {
	s  [][]byte
	si int64 // of current slice
	i  int64 // within current slice
}

func (s *sliceReader) Read(b []byte) (n int, err error) {
	for n < len(b) && s.si < int64(len(s.s)) {
		n += copy(b, s.s[s.si][s.i:])
		b = b[n:]
		s.i += int64(n)
		if s.i >= int64(len(s.s[s.si])) {
			s.si++
			s.i = 0
		}
	}
	if n == 0 {
		return 0, io.EOF
	}
	return n, nil
}

// readSlices returns an io.Reader that reads from a set of slices sequentially.
func readSlices(s [][]byte) io.Reader {
	return &sliceReader{s: s}
}

func intSliceEqual(a, b []int) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil && len(b) == 0:
		return true
	case len(a) == 0 && b == nil:
		return true
	case a == nil || b == nil:
		return false
	case len(a) != len(b):
		return false
	default:
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
}
