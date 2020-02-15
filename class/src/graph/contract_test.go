package graph

import "testing"

func TestContract(t *testing.T) {
	input := []*Edge{
		&Edge{1, 2},
		&Edge{1, 20},
		&Edge{1, 21},
		&Edge{1, 22},
		&Edge{2, 1},
		&Edge{2, 20},
		&Edge{2, 21},
		&Edge{2, 22},
	}

	minEdges := contract(input)

	if len(minEdges) != 2 {
		t.Errorf("got %d edges wanted 2", len(minEdges))
	}
}
