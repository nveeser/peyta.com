package graph

import "testing"

func TestSimple(t *testing.T) {
	g := New()
	g.Add(1, 2)
	g.Add(1, 20)
	g.Add(1, 21)
	g.Add(1, 22)
	g.Add(2, 1)
	g.Add(2, 20)
	g.Add(2, 21)
	g.Add(2, 22)

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

	min := MinCut(rows)

	if len(min) != 2 {
		t.Errorf("got %d wanted 2", len(min))
		t.Logf("min: %v", min)
	}
}

func TestCases2(t *testing.T) {
	rows := []Row{
		Row{1, 2, 3, 4},
		Row{2, 1, 3, 4},
		Row{3, 1, 2, 4},
		Row{4, 1, 2, 3, 5},
		Row{5, 4, 6, 7, 8},
		Row{6, 5, 7, 8},
		Row{7, 5, 6, 8},
		Row{8, 5, 6, 7},
	}

	min := MinCut(rows)

	if len(min) != 1 {
		t.Errorf("got %d wanted 2", len(min))
		t.Logf("min: %v", min)
	}
}
