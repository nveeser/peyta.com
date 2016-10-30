package graph

import (
	"fmt"
	"log"
)

// ----------------------------------------
// ById
// ----------------------------------------

type ById []uint64

func (a ById) Len() int           { return len(a) }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ById) Less(i, j int) bool { return a[i] < a[j] }

// ----------------------------------------
// Rows => Edges
// ----------------------------------------

type Row []uint64

func NewEdges(rows []Row) []*Edge {
	var edges []*Edge

	for _, row := range rows {
		a := row[0]
		for _, b := range row[1:] {
			if a < b {
				e := &Edge{Left: a, Right: b}
				edges = append(edges, e)
			}
		}
	}
	return edges
}

// ----------------------------------------
// Edge
// ----------------------------------------

type ID uint64

type Edge struct {
	Left, Right uint64
}

func (e *Edge) String() string {
	return fmt.Sprintf("(%d, %d)", e.Left, e.Right)
}

type node struct {
	id    uint64
	edges []*Edge
}

func TopoSort(edges []*Edge) map[uint64]int {
	g := newIndex(edges)

	current := len(g.nodes)
	r := make(map[uint64]int)

	visit := func(id uint64) {
		log.Printf("  Label: %d = %d", id, current)
		r[id] = current
		current--
	}

	for id, _ := range g.nodes {
		log.Printf("Loop %d", id)
		g.dfs(id, out, visit)
	}

	if current < 0 {
		log.Fatalf("g.current got too small: %d", current)
	}

	return r
}

func FindBFS(edges []*Edge, first uint64) []uint64 {
	var found []uint64
	distance := make(map[uint64]int)

	visit := func(to, from uint64) {
		//log.Printf("Visit: %d (%d)", current.id, current.distance)
		distance[to] = distance[from] + 1
		found = append(found, to)
	}

	newIndex(edges).bfs(first, visit)
	return found
}

func DFSLoop(edges []*Edge) map[uint64]int {
	r := make(map[uint64]int)
	t := 0
	visit := func(to uint64) {
		t++
		r[to] = t
	}
	i := newIndex(edges)
	for id, _ := range i.nodes {
		i.dfs(id, out, visit)
	}
	return r
}

// func (g *Graph) Dump() {
// 	log.Printf("Dump Graph [%p]", g)
// 	for _, n := range g.nodes {
// 		log.Printf("node[%d] %d edges", n.id, len(n.edges))
// 		var links []string
// 		for _, e := range n.edges {
// 			links = append(links, e.String())
// 		}
// 		sort.Strings(links)
// 		log.Printf("   Edges: %s", strings.Join(links, " "))
// 	}
// }
