package graph

import "testing"

func TestSimple(t *testing.T) {
	g := New()
	g.Add(&Edge{Left: 1, Right: 2})
	g.Add(&Edge{Left: 1, Right: 20})
	g.Add(&Edge{Left: 1, Right: 21})
	g.Add(&Edge{Left: 1, Right: 22})
	g.Add(&Edge{Left: 2, Right: 1})
	g.Add(&Edge{Left: 2, Right: 20})
	g.Add(&Edge{Left: 2, Right: 21})
	g.Add(&Edge{Left: 2, Right: 22})

	if len(g.nodes) != 5 {
		t.Errorf("got %d nodes wanted 5", len(g.nodes))
	}
	if g.edgeset.Len() != 8 {
		t.Errorf("got %d edges wanted 8", g.edgeset.Len())
	}

	g.Contract()

	if len(g.nodes) != 2 {
		t.Errorf("got %d nodes wanted 2", len(g.nodes))
	}
	if g.edgeset.Len() != 2 {
		t.Errorf("got %d edges wanted 2", g.edgeset.Len())
	}
}

func TestCases(t *testing.T) {
	rows := [][]uint64{
		[]uint64{1, 2, 3, 4, 7},
		[]uint64{2, 1, 3, 4},
		[]uint64{3, 1, 2, 4},
		[]uint64{4, 1, 2, 3, 5},
		[]uint64{5, 4, 6, 7, 8},
		[]uint64{6, 5, 7, 8},
		[]uint64{7, 1, 5, 6, 8},
		[]uint64{8, 5, 6, 7},
	}

	edges := MakeEdges(rows)
	min := MinCut(edges)

	if len(min) != 2 {
		t.Errorf("got %d wanted 2", len(min))
		t.Logf("min: %v", min)
	}
}

func TestCases2(t *testing.T) {
	rows := [][]uint64{
		[]uint64{1, 2, 3, 4},
		[]uint64{2, 1, 3, 4},
		[]uint64{3, 1, 2, 4},
		[]uint64{4, 1, 2, 3, 5},
		[]uint64{5, 4, 6, 7, 8},
		[]uint64{6, 5, 7, 8},
		[]uint64{7, 5, 6, 8},
		[]uint64{8, 5, 6, 7},
	}

	edges := MakeEdges(rows)
	min := MinCut(edges)

	if len(min) != 1 {
		t.Errorf("got %d wanted 2", len(min))
		t.Logf("min: %v", min)
	}
}
