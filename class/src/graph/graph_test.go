package graph

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
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
	cases := []struct {
		input    []*Edge
		expected map[ID]int
	}{
		{
			input:    []*Edge{{1, 2}, {2, 4}, {1, 3}, {3, 4}},
			expected: map[ID]int{1: 1, 4: 4},
		},
		// {
		// 	input:    []*Edge{{1, 4}, {2, 8}, {3, 6}, {4, 7}, {5, 2}, {6, 9}, {7, 1}, {8, 5}, {8, 6}, {9, 7}, {9, 3}},
		// 	expected: map[ID]int{1: 1, 4: 4},
		// },
		// {
		// 	input:    []*Edge{{6, 9}, {3, 6}, {9, 3}, {8, 6}, {2, 8}, {5, 2}, {8, 5}, {9, 7}, {4, 7}, {1, 4}, {7, 1}},
		// 	expected: map[ID]int{1: 1, 4: 4},
		// },
	}

	for _, tc := range cases {
		r := TopoSort(tc.input)
		t.Logf("TopoSort Result: %s", dumpMap(r))
		for k, v := range tc.expected {
			if r[k] != v {
				t.Errorf("got %d wanted %d", r[k], v)
			}
		}
	}
}

func dumpMap(m map[ID]int) string {
	var keys []ID
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Sort(Ascending(keys))

	var out []string
	for _, k := range keys {
		out = append(out, fmt.Sprintf("%d:%d", k, m[k]))
	}
	return "map{" + strings.Join(out, " ") + "}"
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
		{3, 6, 9},
		{2, 5, 8},
		{1, 4, 7},
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

func TestLargestGroups(t *testing.T) {
	cases := []struct {
		input    []*Edge
		expected []int
	}{
		{
			input:    []*Edge{{1, 4}, {2, 8}, {3, 6}, {4, 7}, {5, 2}, {6, 9}, {7, 1}, {8, 5}, {8, 6}, {9, 7}, {9, 3}},
			expected: []int{3, 3, 3, 0, 0},
		},
		{
			input:    []*Edge{{6, 9}, {3, 6}, {9, 3}, {8, 6}, {2, 8}, {5, 2}, {8, 5}, {9, 7}, {4, 7}, {1, 4}, {7, 1}},
			expected: []int{3, 3, 3, 0, 0},
		},
		{
			input:    []*Edge{{1, 2}, {2, 3}, {3, 1}, {3, 4}, {5, 4}, {6, 4}, {8, 6}, {6, 7}, {7, 8}, {4, 3}, {4, 6}},
			expected: []int{7, 1, 0, 0, 0},
		},
		{
			input:    []*Edge{{1, 2}, {2, 3}, {2, 4}, {2, 5}, {3, 6}, {4, 5}, {4, 7}, {5, 2}, {5, 6}, {5, 7}, {6, 3}, {6, 8}, {7, 8}, {7, 10}, {8, 7}, {9, 7}, {10, 9}, {10, 11}, {11, 12}, {12, 10}},
			expected: []int{6, 3, 2, 1, 0},
		},
	}

	for i, tc := range cases {
		groups := Kosaraju(tc.input)

		got := LargestGroups(groups, 5)

		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("%d got %v wanted %v", i, got, tc.expected)
		}
	}
}
