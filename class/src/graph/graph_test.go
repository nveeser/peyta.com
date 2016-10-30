package graph

import (
	"reflect"
	"sort"
	"testing"
)

func TestWalkBFS(t *testing.T) {
	edges := []*Edge{
		{1, 2},
		{1, 3},
		{4, 5},
		{4, 6},
		{6, 7},
	}

	cases := []struct {
		start    ID
		expected []ID
	}{
		{1, []ID{1, 2, 3}},
		{2, []ID{1, 2, 3}},
		{3, []ID{1, 2, 3}},
		{4, []ID{4, 5, 6, 7}},
		{5, []ID{4, 5, 6, 7}},
		{6, []ID{4, 5, 6, 7}},
		{7, []ID{4, 5, 6, 7}},
	}

	for _, tc := range cases {
		found := WalkBFS(edges, tc.start)

		sort.Sort(Ascending(found))
		if !reflect.DeepEqual(found, tc.expected) {
			t.Errorf("got %v wanted %v", found, tc.expected)
		}
	}
}

func TestDistance(t *testing.T) {
	edges := []*Edge{
		{1, 2},
		{1, 3},
		{4, 5},
		{4, 6},
		{6, 7},
	}

	cases := []struct {
		start    ID
		expected map[ID]int
	}{
		{1, map[ID]int{1: 0, 2: 1, 3: 1}},
		{2, map[ID]int{1: 1, 2: 0, 3: 2}},
		{3, map[ID]int{1: 1, 2: 2, 3: 0}},
		{4, map[ID]int{4: 0, 5: 1, 6: 1, 7: 2}},
		{5, map[ID]int{4: 1, 5: 0, 6: 2, 7: 3}},
		{6, map[ID]int{4: 1, 5: 2, 6: 0, 7: 1}},
		{7, map[ID]int{4: 2, 5: 3, 6: 1, 7: 0}},
	}

	for _, tc := range cases {
		t.Logf("Distance(%d)", tc.start)
		got := Distance(edges, tc.start)

		for k, v := range tc.expected {
			if got[k] != v {
				t.Errorf("start: %d got r[%d] = %d wanted %d", tc.start, k, got[k], v)
			}
		}
		for k, v := range got {
			t.Logf("got[%d] = %d", k, v)
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

type GroupBySize []Group

func (a GroupBySize) Len() int      { return len(a) }
func (a GroupBySize) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a GroupBySize) Less(i, j int) bool {
	switch {
	case len(a[i]) != len(a[j]):
		return len(a[i]) < len(a[j])
	default:
		return a[i][0] < a[j][0]
	}
}

func TestKosaraju(t *testing.T) {
	input := []*Edge{
		{6, 9},
		{3, 6},
		{9, 3},
		{8, 6},
		{2, 8},
		{5, 2},
		{8, 5},
		{9, 7},
		{4, 7},
		{1, 4},
		{7, 1},
	}

	expected := []Group{
		{1, 4, 7},
		{2, 5, 8},
		{3, 6, 9},
	}

	got := Kosaraju(input)
	for _, group := range got {
		sort.Sort(Ascending(group))
	}
	sort.Sort(GroupBySize(got))

	for k, v := range got {
		t.Logf("got[%d] = %d", k, v)
	}

	for i := range expected {
		sort.Sort(Ascending(got[i]))

		if !reflect.DeepEqual(got[i], expected[i]) {
			t.Errorf("[%d] got %v wanted %v", i, got[i], expected[i])
		}
	}
}
