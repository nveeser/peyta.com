package graph

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strings"
)

// ----------------------------------------
// Edge
// ----------------------------------------

type Edge struct {
	Left, Right uint64
}

func (e *Edge) String() string {
	return fmt.Sprintf("(%d, %d)", e.Left, e.Right)
}

// ----------------------------------------
// Node
// ----------------------------------------

type Node struct {
	id    uint64
	edges []*Edge
}

// ----------------------------------------
// MinCut
// ----------------------------------------

type Row []uint64

func MakeEdges(rows []Row) []*Edge {
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

func MinCut(edges []*Edge, count int) []*Edge {
	var result []*Edge
	for i := 0; i < count; i++ {
		g := New(edges)
		//log.Printf("g : %+v", g)
		g.Dump()
		min := contract(g)
		log.Printf("[%d/%d] Size: %d  -- %v", i, count, len(min), min)

		if result == nil || len(result) > len(min) {
			result = min
		}
	}
	return result
}

// ----------------------------------------
// Graph
// ----------------------------------------

func New(edges []*Edge) *Graph {
	g := &Graph{nodes: make(map[uint64]*Node)}

	get := func(id uint64) *Node {
		if _, ok := g.nodes[id]; !ok {
			g.nodes[id] = &Node{id: id}
		}
		return g.nodes[id]
	}

	for _, e := range edges {
		n1 := get(e.Left)
		n1.edges = append(n1.edges, e)

		n2 := get(e.Right)
		n2.edges = append(n2.edges, e)

		g.edges = append(g.edges, e)
	}
	return g
}

type Graph struct {
	nodes map[uint64]*Node
	edges []*Edge
}

func (g *Graph) EdgeLen() int {
	return len(g.edges)
}

func (g *Graph) Dump() {
	log.Printf("Dump Graph [%p]", g)
	for _, n := range g.nodes {
		log.Printf("Node[%d] %d edges", n.id, len(n.edges))
		var links []string
		for _, e := range n.edges {
			links = append(links, e.String())
		}
		sort.Strings(links)
		log.Printf("   Edges: %s", strings.Join(links, " "))
	}
}

type contractor struct {
	nodes   map[uint64]*cNode
	edgeset *edgeSet
}

func (c *contractor) resolve(from uint64) uint64 {
	to := c.nodes[from].replacedBy
	if to == from {
		return from
	}
	return c.resolve(to)
}

type cNode struct {
	combined   []*Edge
	replacedBy uint64
}

func contract(g *Graph) []*Edge {
	c := &contractor{
		nodes:   make(map[uint64]*cNode),
		edgeset: &edgeSet{edges: g.edges},
	}
	for k, n := range g.nodes {
		c.nodes[k] = &cNode{
			combined:   n.edges,
			replacedBy: k,
		}
	}

	for i := len(c.nodes); i > 2; i-- {
		drop := c.edgeset.random()

		winID := c.resolve(drop.Left)
		loseID := c.resolve(drop.Right)

		win := c.nodes[winID]
		lose := c.nodes[loseID]

		log.Printf("Contract Edge: %v %d <= %d", drop, winID, loseID)

		c.edgeset.erase(drop)
		c.nodes[loseID].replacedBy = winID

		var combined []*Edge
		combined = append(win.combined, lose.combined...)
		win.combined = nil

		for _, e := range combined {
			if c.resolve(e.Left) != c.resolve(e.Right) {
				//log.Printf("[] Keep: %s", e)
				win.combined = append(win.combined, e)
			} else {
				c.edgeset.erase(e)
				//log.Printf("[] Loop: %s", e)
			}
		}

		//log.Printf("  Deleting Node:[%d]", loseID)
	}
	return c.edgeset.edges
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
