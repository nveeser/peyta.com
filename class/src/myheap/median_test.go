package myheap

import (
	"strconv"
	"testing"
)

func TestMedia(t *testing.T) {
	cases := []struct {
		input    []int64
		expected []int64
	}{
		{
			[]int64{1, 2, 3, 4, 5, 6},
			[]int64{1, 1, 2, 2, 3, 3},
		},
		{
			[]int64{10, 30, 20, 31, 32, 33, 40, 1, 2},
			[]int64{10, 10, 20, 20, 30, 30, 31, 30, 30},
		},
	}

	for i, tc := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			m := &Median{}
			for i, n := range tc.input {
				m.Add(n)
				t.Logf("Added %d Median: %d Struct: %v", n, m.Value(), m)
				if m.Value() != tc.expected[i] {
					t.Errorf("got median %d wanted %d", m.Value(), tc.expected[i])
				}
			}
		})
	}
}
