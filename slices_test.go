package mmr

import "testing"

type testcase struct {
	a        []int
	b        []int
	expected bool
}

func TestIntSliceEqual(t *testing.T) {
	cases := []testcase{
		{[]int{1, 1, 2, 3}, []int{1, 1, 2, 3}, true},
		{[]int{1, 1, 2, 3}, []int{1, 1, 2}, false},
		{[]int{1, 1, 2, 3}, []int{}, false},
		{[]int{1, 1, 2, 3}, nil, false},
		{nil, nil, true},
	}

	for _, c := range cases {
		o := intSliceEqual(c.a, c.b)
		if o != c.expected {
			t.Errorf("intSliceEqual(%v, %v): %v, expected %v", c.a, c.b, o, c.expected)
		}
	}
}
