package graph

import (
	"reflect"
	"sort"
	"testing"
)

func TestFindBFS(t *testing.T) {
	edges := []*Edge{
		{1, 2},
		{1, 3},
		{4, 5},
		{4, 6},
		{6, 7},
	}

	cases := []struct {
		start    uint64
		expected []uint64
	}{
		{1, []uint64{1, 2, 3}},
		{2, []uint64{1, 2, 3}},
		{3, []uint64{1, 2, 3}},
		{4, []uint64{4, 5, 6, 7}},
		{5, []uint64{4, 5, 6, 7}},
		{6, []uint64{4, 5, 6, 7}},
		{7, []uint64{4, 5, 6, 7}},
	}

	for _, tc := range cases {
		found := FindBFS(edges, tc.start)

		sort.Sort(ById(found))
		if !reflect.DeepEqual(found, tc.expected) {
			t.Errorf("got %v wanted %v", found, tc.expected)
		}
	}
}

func TestTopoSort(t *testing.T) {
	input := []*Edge{
		{1, 2},
		{2, 4},
		{1, 3},
		{3, 4},
	}

	r := TopoSort(input)
	if r[1] != 1 {
		t.Errorf("got %d wanted 1", r[1])
		t.Logf("Map: %v", r)
	}
	if r[4] != 4 {
		t.Errorf("got %d wanted 4", r[4])
		t.Logf("Map: %v", r)
	}
}

// func TestLoopCount(t *testing.T) {
// 	input := []*Edge{
// 		{9, 6},
// 		{6, 3},
// 		{3, 9},
// 		{6, 8},
// 		{8, 2},
// 		{2, 5},
// 		{5, 8},
// 		{7, 9},
// 		{7, 4},
// 		{4, 1},
// 		{1, 7},
// 	}

// 	expected := map[uint64]int{
// 		1: 7, 2: 3, 3: 1, 4: 8, 5: 2, 6: 5, 7: 9, 8: 4, 9: 6,
// 	}

// 	got := DFSLoop(input)

// 	for k, v := range expected {
// 		if got[k] != v {
// 			t.Errorf("got r[%d] = %d wanted %d", k, got[k], v)
// 		}
// 	}
// 	for k, v := range got {
// 		t.Logf("got[%d] = %d", k, v)
// 	}
// }
