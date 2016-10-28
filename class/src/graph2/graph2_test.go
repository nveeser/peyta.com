package graph2

import "testing"

func TestSimple(t *testing.T) {

	edges := []*Edge{
		&Edge{Left: 1, Right: 2},
		&Edge{Left: 1, Right: 20},
		&Edge{Left: 1, Right: 21},
		&Edge{Left: 1, Right: 22},
		&Edge{Left: 2, Right: 20},
		&Edge{Left: 2, Right: 21},
		&Edge{Left: 2, Right: 22},
	}
	BFS(edges, 1)
	t.Errorf("Something went wrong")
}
