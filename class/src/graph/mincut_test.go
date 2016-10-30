package graph

import "testing"

func TestMinCut1(t *testing.T) {
	rows := []Row{
		[]ID{1, 2, 3, 4, 7},
		[]ID{2, 1, 3, 4},
		[]ID{3, 1, 2, 4},
		[]ID{4, 1, 2, 3, 5},
		[]ID{5, 4, 6, 7, 8},
		[]ID{6, 5, 7, 8},
		[]ID{7, 1, 5, 6, 8},
		[]ID{8, 5, 6, 7},
	}

	edges := NewEdges(rows)
	min := MinCut(edges, 10)

	if len(min) != 2 {
		t.Errorf("got %d wanted 2", len(min))
		t.Logf("min: %v", min)
	}
}

func TestMinCut2(t *testing.T) {
	rows := []Row{
		[]ID{1, 2, 3, 4},
		[]ID{2, 1, 3, 4},
		[]ID{3, 1, 2, 4},
		[]ID{4, 1, 2, 3, 5},
		[]ID{5, 4, 6, 7, 8},
		[]ID{6, 5, 7, 8},
		[]ID{7, 5, 6, 8},
		[]ID{8, 5, 6, 7},
	}

	edges := NewEdges(rows)
	min := MinCut(edges, len(rows))

	if len(min) != 1 {
		t.Errorf("got %d wanted 2", len(min))
		t.Logf("min: %v", min)
	}
}
