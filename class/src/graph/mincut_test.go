package graph

import "testing"

func TestMinCut1(t *testing.T) {
	rows := []Row{
		[]uint64{1, 2, 3, 4, 7},
		[]uint64{2, 1, 3, 4},
		[]uint64{3, 1, 2, 4},
		[]uint64{4, 1, 2, 3, 5},
		[]uint64{5, 4, 6, 7, 8},
		[]uint64{6, 5, 7, 8},
		[]uint64{7, 1, 5, 6, 8},
		[]uint64{8, 5, 6, 7},
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
		[]uint64{1, 2, 3, 4},
		[]uint64{2, 1, 3, 4},
		[]uint64{3, 1, 2, 4},
		[]uint64{4, 1, 2, 3, 5},
		[]uint64{5, 4, 6, 7, 8},
		[]uint64{6, 5, 7, 8},
		[]uint64{7, 5, 6, 8},
		[]uint64{8, 5, 6, 7},
	}

	edges := NewEdges(rows)
	min := MinCut(edges, len(rows))

	if len(min) != 1 {
		t.Errorf("got %d wanted 2", len(min))
		t.Logf("min: %v", min)
	}
}
