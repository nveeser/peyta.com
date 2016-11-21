package graph

import "testing"

func TestSimpleGraphs(t *testing.T) {
	cases := []struct {
		name     string
		input    []*SizeEdge
		expected map[ID]Size
	}{
		{
			"Basic",
			[]*SizeEdge{
				NewSizeEdge(1, 3, 3),
				NewSizeEdge(1, 2, 1),
				NewSizeEdge(2, 3, 1),
				NewSizeEdge(2, 4, 4),
				NewSizeEdge(3, 4, 6),
				NewSizeEdge(3, 4, 6),
			},
			map[ID]Size{1: 0, 2: 1, 3: 2, 4: 5},
		},
		{
			"BackTrack",
			[]*SizeEdge{
				NewSizeEdge(1, 2, 1),
				NewSizeEdge(2, 3, 1),
				NewSizeEdge(3, 4, 1),
				NewSizeEdge(1, 5, 6),
				NewSizeEdge(4, 5, 2),
			},
			map[ID]Size{1: 0, 2: 1, 3: 2, 4: 3, 5: 5},
		},
		{
			"Disconnected",
			[]*SizeEdge{
				NewSizeEdge(1, 2, 1),
				NewSizeEdge(2, 3, 1),
				NewSizeEdge(3, 4, 1),
				NewSizeEdge(1, 5, 6),
				NewSizeEdge(4, 5, 2),

				NewSizeEdge(6, 7, 2),
				NewSizeEdge(7, 8, 3),
				NewSizeEdge(8, 6, 4),
			},
			map[ID]Size{1: 0, 2: 1, 3: 2, 4: 3, 5: 5, 6: Inf, 7: Inf, 8: Inf},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, e := range tc.input {
				t.Logf("Edge: %s", e)
			}
			got := FindDistances(tc.input, 1)

			for k, v := range tc.expected {
				if got[k] != Size(v) {
					t.Errorf("start: %d got r[%d] = %d wanted %d", 1, k, got[k], v)
				}
			}
			for k, v := range got {
				t.Logf("got[%d] = %d", k, v)
			}
		})
	}
}

func TestLines(t *testing.T) {
	input := []string{
		"1 2,3 3,2",
		"2 4,4",
		"3 2,1 4,2 5,3",
		"4 5,2 6,1",
		"5 6,2",
	}

	expected := map[ID]int{1: 0, 2: 3, 3: 2, 4: 4, 5: 5, 6: 5}

	var edges []*SizeEdge
	for _, l := range input {
		e, err := ParseNodePaths(l, " ")
		if err != nil {
			t.Logf("Line: %s", l)
			t.Fatalf("got error parsing the line: %s", err)
		}
		edges = append(edges, e...)
	}

	for _, e := range edges {
		t.Logf("Edge: %s", e)
	}

	got := FindDistances(edges, 1)
	for k, v := range expected {
		if got[k] != Size(v) {
			t.Errorf("start: %d got r[%d] = %d wanted %d", 1, k, got[k], v)
		}
	}
	for k, v := range got {
		t.Logf("got[%d] = %d", k, v)
	}

}
