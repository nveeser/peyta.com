package graph

import (
	"fmt"
	"log"
	"math/rand"
)

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

func (g *Graph) Add(a, b uint64) {
	n1 := g.node(a)
	n2 := g.node(b)

	e := &Edge{n1, n2}
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

// ----------------------------------------
// Edge
// ----------------------------------------

type Edge struct {
	left, right *Node
}

func (e *Edge) String() string {
	if e.left.id < e.right.id {
		return fmt.Sprintf("(%d, %d)", e.left.id, e.right.id)
	}
	return fmt.Sprintf("(%d, %d)", e.right.id, e.left.id)
}

func (e Edge) isLoop() bool {
	return e.left == e.right
}

func (e *Edge) update(before, after *Node) {
	pre := e.String()
	switch before {
	case e.left:
		e.left = after
	case e.right:
		e.right = after
	default:
		log.Fatalf("update(): %s does not point to node: %d", e, before.id)
	}
	post := e.String()
	_ = fmt.Sprintf("   updated edge %s => %s", pre, post)
	//log.Print(m)
}

func (e *Edge) opposite(n *Node) *Node {
	switch n {
	case e.left:
		return e.right
	case e.right:
		return e.left
	default:
		log.Fatalf("opposite(): %s does not connect to node[%d]", e, n.id)
	}
	return nil
}

// ----------------------------------------
// Node
// ----------------------------------------

type Node struct {
	id    uint64
	edges *edgeSet
}

func (n *Node) String() string {
	return fmt.Sprintf("%d  %v", n.id, n.Peers())
}

func (n *Node) Peers() []uint64 {
	var r []uint64
	for _, e := range n.edges.edges {
		peer := e.opposite(n)
		r = append(r, peer.id)
	}
	return r
}

// ----------------------------------------
// Contract
// ----------------------------------------

type Row []uint64

func key(left, right uint64) string {
	if left < right {
		return fmt.Sprintf("%d%d", left, right)
	}
	return fmt.Sprintf("%d%d", right, left)
}

func MinCut(rows []Row) []*Edge {
	var result []*Edge

	for i := 0; i < len(rows); i++ {
		g := New()
		exists := make(map[string]bool)
		for _, row := range rows {
			a := row[0]
			for _, b := range row[1:] {
				k := key(a, b)
				if _, ok := exists[k]; !ok {
					g.Add(a, b)
					exists[k] = true
				}
			}
		}
		//g.Dump()
		g.Contract()
		log.Printf("Size: %d", g.EdgeLen())

		if result == nil || len(result) > g.EdgeLen() {
			result = g.edgeset.edges
		}
	}
	return result
}

// ----------------------------------------
// Contract
// ----------------------------------------

func (g *Graph) Contract() {
	for len(g.nodes) > 2 {
		//log.Printf("------ Nodes: %d --------", len(g.nodes))

		drop := g.edgeset.random()

		win := drop.left
		lose := drop.right
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
func (s *edgeSet) Swap(i, j int) {
	s.edges[i], s.edges[j] = s.edges[j], s.edges[i]
}
func (s *edgeSet) Less(i, j int) bool {
	if s.edges[i].left.id != s.edges[j].left.id {
		return s.edges[i].left.id < s.edges[j].left.id
	}
	if s.edges[i].right.id != s.edges[j].right.id {
		return s.edges[i].right.id < s.edges[j].right.id
	}
	return false
}
