package graph

import "testing"

// func TestSimple(t *testing.T) {
// 	input := []*Edge{
// 		&Edge{Left: 1, Right: 2},
// 		&Edge{Left: 1, Right: 20},
// 		&Edge{Left: 1, Right: 21},
// 		&Edge{Left: 1, Right: 22},
// 		&Edge{Left: 2, Right: 1},
// 		&Edge{Left: 2, Right: 20},
// 		&Edge{Left: 2, Right: 21},
// 		&Edge{Left: 2, Right: 22},
// 	}
// 	g := New(input)

// 	if len(g.nodes) != 5 {
// 		t.Errorf("got %d nodes wanted 5", len(g.nodes))
// 	}
// 	if len(g.edges) != 8 {
// 		t.Errorf("got %d edges wanted 8", len(g.edges))
// 	}

// 	minEdges := contract(g)

// 	if len(g.nodes) != 5 {
// 		t.Errorf("got %d nodes wanted 5", len(g.nodes))
// 	}
// 	if len(g.edges) != 8 {
// 		t.Errorf("got %d edges wanted 8", len(g.edges))
// 	}

// 	if len(minEdges) != 2 {
// 		t.Errorf("got %d edges wanted 2", len(minEdges))
// 	}
// }

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

	edges := MakeEdges(rows)
	min := MinCut(edges, 10)

	if len(min) != 2 {
		t.Errorf("got %d wanted 2", len(min))
		t.Logf("min: %v", min)
	}
}

// func TestCases2(t *testing.T) {
// 	rows := []Row{
// 		[]uint64{1, 2, 3, 4},
// 		[]uint64{2, 1, 3, 4},
// 		[]uint64{3, 1, 2, 4},
// 		[]uint64{4, 1, 2, 3, 5},
// 		[]uint64{5, 4, 6, 7, 8},
// 		[]uint64{6, 5, 7, 8},
// 		[]uint64{7, 5, 6, 8},
// 		[]uint64{8, 5, 6, 7},
// 	}

// 	edges := MakeEdges(rows)
// 	min := MinCut(edges, len(rows))

// 	if len(min) != 1 {
// 		t.Errorf("got %d wanted 2", len(min))
// 		t.Logf("min: %v", min)
// 	}
// }
