package graph

import (
	"fmt"
	"log"
)

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

func (e *Edge) Peer(id ID) (ID, Direction) {
	switch {
	case id == e.Left:
		return e.Right, right
	case id == e.Right:
		return e.Left, left
	default:
		log.Fatalf("id %d not on edge: %s", id, e)
	}
	return 0, right
}

// ----------------------------------------
// Edge
// ----------------------------------------

type Direction bool

const right Direction = true
const left Direction = false

func (d Direction) String() string {
	if d == right {
		return "right"
	}
	return "left"
}

// ----------------------------------------
// Ascending
// ----------------------------------------

// Ascending is for sorting a slice of IDs ascending.
type Ascending []ID

func (a Ascending) Len() int           { return len(a) }
func (a Ascending) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Ascending) Less(i, j int) bool { return a[i] < a[j] }

// ----------------------------------------
// Rows => Edges
// ----------------------------------------

type VertexRow []ID

func NewEdges(rows []VertexRow) []*Edge {
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
// Functions
// ----------------------------------------

func TopoSort(edges []*Edge) map[ID]int {
	g := newIndex(edges)

	current := len(g.edgesByID)
	r := make(map[ID]int)

	visit := func(id ID) {
		//log.Printf("  Label: %d = %d", id, current)
		r[id] = current
		current--
	}

	for id, _ := range g.edgesByID {
		//log.Printf("Loop %d", id)
		g.dfs(id, right, visit)
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

type Group []ID

func Kosaraju(edges []*Edge) []Group {
	r := make(map[ID]ID)
	t := 0
	visit := func(id ID) {
		t++
		r[ID(t)] = id
	}
	index := newIndex(edges)
	for id, _ := range index.edgesByID {
		index.dfs(id, left, visit)
	}

	index.reset()

	leaders := make(map[ID]ID)
	group := make(map[ID]Group)

	for i := t; i > 0; i-- {
		leader := r[ID(i)]
		//log.Printf("leader: %d", leader)
		index.dfs(leader, right, func(id ID) {
			leaders[id] = leader
			g := group[leader]
			group[leader] = append(g, id)
		})
	}

	var result []Group
	for _, g := range group {
		result = append(result, g)
	}
	return result
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
