package graph

import (
	"math/rand"
)

type contractor struct {
	nodes map[ID]*cNode
	edges []*Edge
}

type cNode struct {
	combined   []*Edge
	replacedBy ID
}

func (c *contractor) resolve(from ID) ID {
	to := c.nodes[from].replacedBy
	if to == from {
		return from
	}
	return c.resolve(to)
}

func (c *contractor) random() *Edge {
	i := rand.Int31() % int32(len(c.edges))
	return c.edges[i]
}

func (c *contractor) erase(e *Edge) bool {
	for i, x := range c.edges {
		if x == e {
			c.edges = append(c.edges[:i], c.edges[i+1:]...)
			return true
		}
	}
	return false
}

func index(edges []*Edge) *contractor {
	c := &contractor{
		nodes: make(map[ID]*cNode),
		edges: append([]*Edge(nil), edges...),
	}

	memo := func(id ID) *cNode {
		if _, ok := c.nodes[id]; !ok {
			c.nodes[id] = &cNode{replacedBy: id}
		}
		return c.nodes[id]
	}

	for _, e := range edges {
		n1 := memo(e.Left)
		n1.combined = append(n1.combined, e)

		n2 := memo(e.Right)
		n2.combined = append(n2.combined, e)
	}
	return c
}

func contract(edges []*Edge) []*Edge {
	c := index(edges)

	for i := len(c.nodes); i > 2; i-- {
		drop := c.random()

		winID := c.resolve(drop.Left)
		loseID := c.resolve(drop.Right)

		win := c.nodes[winID]
		lose := c.nodes[loseID]

		//log.Printf("Contract Edge: %v %d <= %d", drop, winID, loseID)
		//log.Printf("Edges: %v", c.edgeset.edges)
		c.erase(drop)
		c.nodes[loseID].replacedBy = winID

		var combined []*Edge
		combined = append(win.combined, lose.combined...)
		win.combined = nil

		for _, e := range combined {
			if c.resolve(e.Left) != c.resolve(e.Right) {
				//log.Printf("Keep: %s", e)
				win.combined = append(win.combined, e)
			} else {
				c.erase(e)
				//log.Printf("Loop: %s", e)
			}
		}

		//log.Printf("  Deleting Node:[%d]", loseID)
	}
	return c.edges
}
