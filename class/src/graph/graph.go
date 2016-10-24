package graph

import (
	"fmt"
	"log"
	"math/rand"
)

// ----------------------------------------
// MinCut
// ----------------------------------------

func key(left, right uint64) string {
	if left < right {
		return fmt.Sprintf("%d%d", left, right)
	}
	return fmt.Sprintf("%d%d", right, left)
}

type Row []uint64

func MakeEdges(rows []Row) []*Edge {
	var edges []*Edge

	exists := make(map[string]bool)
	for _, row := range rows {
		a := row[0]
		for _, b := range row[1:] {
			k := key(a, b)
			if _, ok := exists[k]; !ok {
				e := &Edge{Left: a, Right: b}
				edges = append(edges, e)
				exists[k] = true
			}
		}
	}
	return edges
}

func MinCut(count int, edges []*Edge) []*Edge {
	var result []*Edge

	for i := 0; i < count; i++ {
		g := New()
		for _, e := range edges {
			g.Add(e)
		}
		//g.Dump()
		g.Contract()
		log.Printf("[%d/%d] Size: %d", i, count, g.EdgeLen())

		if result == nil || len(result) > g.EdgeLen() {
			result = g.edgeset.edges
		}
	}
	return result
}

// ----------------------------------------
// Graph
// ----------------------------------------

func New() *Graph {
	return &Graph{
		nodes:   make(map[uint64]*Node),
		edgeset: &edgeSet{},
	}
}

type Graph struct {
	nodes   map[uint64]*Node
	edgeset *edgeSet
}

func (g *Graph) EdgeLen() int {
	return g.edgeset.Len()
}

func (g *Graph) Add(e *Edge) {
	n1 := g.node(e.Left)
	e.leftn = n1
	n2 := g.node(e.Right)
	e.rightn = n2

	n1.edges.add(e)
	n2.edges.add(e)

	g.edgeset.add(e)
}

func (g *Graph) node(id uint64) *Node {
	n, ok := g.nodes[id]
	if !ok {
		n = &Node{id: id, edges: &edgeSet{}}
		g.nodes[id] = n
	}
	return n
}

func (g *Graph) Dump() {
	log.Printf("Dump Graph [%p]", g)
	for _, n := range g.nodes {
		log.Printf("Node[%d] %d edges", n.id, n.edges.Len())
		for _, e := range n.edges.edges {
			log.Printf("   Edge: %s", e)
		}
	}
}

func (g *Graph) Contract() {
	for len(g.nodes) > 2 {
		//log.Printf("------ Nodes: %d --------", len(g.nodes))

		drop := g.edgeset.random()

		win := drop.leftn
		lose := drop.rightn
		//log.Printf("Contract Edge: %v %d => %d", drop, lose.id, win.id)

		g.edgeset.erase(drop)

		for _, e := range lose.edges.edges {
			e.update(lose, win)
			if e.isLoop() {
				win.edges.erase(e)
				g.edgeset.erase(e)
			} else {
				win.edges.add(e)
			}
		}

		//log.Printf("  Deleting Node:[%d]", lose.id)
		lose.edges = nil
		delete(g.nodes, lose.id)

		// for _, e := range win.edges.edges {
		// 	log.Printf("   Edge: %s", e)
		// }
	}

}

// ----------------------------------------
// Edge
// ----------------------------------------

type Edge struct {
	Left, Right   uint64
	leftn, rightn *Node
}

func (e *Edge) String() string {
	if e.leftn.id < e.rightn.id {
		return fmt.Sprintf("(%d, %d)", e.leftn.id, e.rightn.id)
	}
	return fmt.Sprintf("(%d, %d)", e.rightn.id, e.leftn.id)
}

func (e Edge) isLoop() bool {
	return e.leftn == e.rightn
}

func (e *Edge) update(before, after *Node) {
	pre := e.String()
	switch before {
	case e.leftn:
		e.leftn = after
	case e.rightn:
		e.rightn = after
	default:
		log.Fatalf("update(): %s does not point to node: %d", e, before.id)
	}
	post := e.String()
	_ = fmt.Sprintf("   updated edge %s => %s", pre, post)
	//log.Print(m)
}

// ----------------------------------------
// Node
// ----------------------------------------

type Node struct {
	id    uint64
	edges *edgeSet
}

// ----------------------------------------
// edgeSet
// ----------------------------------------

type edgeSet struct {
	edges []*Edge
}

func (s *edgeSet) add(e *Edge) {
	s.edges = append(s.edges, e)
}

func (s *edgeSet) random() *Edge {
	i := rand.Int31() % int32(len(s.edges))
	return s.edges[i]
}

func (s *edgeSet) erase(e *Edge) bool {
	for i, x := range s.edges {
		if x == e {
			s.edges = append(s.edges[:i], s.edges[i+1:]...)
			return true
		}
	}
	return false
}

func (s *edgeSet) Len() int {
	return len(s.edges)
}

// func (s *edgeSet) Swap(i, j int) {
// 	s.edges[i], s.edges[j] = s.edges[j], s.edges[i]
// }
// func (s *edgeSet) Less(i, j int) bool {
// 	if s.edges[i].left.id != s.edges[j].left.id {
// 		return s.edges[i].left.id < s.edges[j].left.id
// 	}
// 	if s.edges[i].right.id != s.edges[j].right.id {
// 		return s.edges[i].right.id < s.edges[j].right.id
// 	}
// 	return false
// }
