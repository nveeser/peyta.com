package graph

import (
	"fmt"
	"log"
)

// ----------------------------------------
// ById
// ----------------------------------------

type ByID []ID

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i] < a[j] }

// ----------------------------------------
// Rows => Edges
// ----------------------------------------

type Row []ID

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
	Left, Right ID
}

func (e *Edge) String() string {
	return fmt.Sprintf("(%d, %d)", e.Left, e.Right)
}

type node struct {
	id    ID
	edges []*Edge
}

func TopoSort(edges []*Edge) map[ID]int {
	g := newIndex(edges)

	current := len(g.nodes)
	r := make(map[ID]int)

	visit := func(id ID) {
		//log.Printf("  Label: %d = %d", id, current)
		r[id] = current
		current--
	}

	for id, _ := range g.nodes {
		//log.Printf("Loop %d", id)
		g.dfs(id, out, visit)
	}

	if current < 0 {
		log.Fatalf("g.current got too small: %d", current)
	}

	return r
}

func Distance(edges []*Edge, first ID) map[ID]int {
	distance := make(map[ID]int)

	newIndex(edges).bfs(first, func(to, from ID) {
		//log.Printf("Visit: %d (%d)", current.id, current.distance)
		switch {
		case to == first:
			distance[to] = 0
		default:
			distance[to] = distance[from] + 1
		}
	})

	return distance
}

func WalkBFS(edges []*Edge, first ID) []ID {
	var found []ID
	newIndex(edges).bfs(first, func(to, from ID) {
		//log.Printf("Visit: %d (%d)", current.id, current.distance)
		found = append(found, to)
	})
	return found
}

func Kosaraju(edges []*Edge) map[ID]ID {
	r := make(map[ID]ID)
	t := 0
	visit := func(id ID) {
		t++
		r[ID(t)] = id
	}
	index := newIndex(edges)
	for id, _ := range index.nodes {
		index.dfs(id, in, visit)
	}

	index.reset()

	leaders := make(map[ID]ID)

	for i := t; i > 0; i-- {
		leader := r[ID(i)]
		log.Printf("leader: %d", leader)
		index.dfs(leader, out, func(id ID) {
			leaders[id] = leader
		})
	}
	return leaders
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
