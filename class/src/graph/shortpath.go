package graph

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Size int

func (s Size) String() string {
	if s == Inf {
		return "INF"
	}
	return strconv.Itoa(int(s))
}

const Inf Size = 1000000

func NewSizeEdge(left, right ID, size Size) *SizeEdge {
	return &SizeEdge{&Edge{left, right}, size}
}

type SizeEdge struct {
	*Edge
	Size Size
}

func (s SizeEdge) String() string {
	return fmt.Sprintf("%s/%s", s.Edge, s.Size)
}

type pathNode struct {
	edges    []*SizeEdge
	distance Size
	seen     bool
	path     []ID
}

func newSizeIndex(edges []*SizeEdge) *pathIndex {
	nodes := make(map[ID]*pathNode)

	node := func(id ID) *pathNode {
		n, ok := nodes[id]
		if !ok {
			n = &pathNode{distance: Inf}
			nodes[id] = n
		}
		return n
	}

	for _, e := range edges {
		n := node(e.Left)
		n.edges = append(n.edges, e)

		n = node(e.Right)
		n.edges = append(n.edges, e)
	}

	return &pathIndex{
		unseen: nodes,
		seen:   make(map[ID]*pathNode),
	}
}

type pathIndex struct {
	unseen map[ID]*pathNode
	seen   map[ID]*pathNode
}

func (i *pathIndex) removeClosest() (ID, *pathNode) {
	best := Inf
	var r ID
	var rn *pathNode
	for id, n := range i.unseen {
		if n.distance <= best {
			best = n.distance
			r = id
		}
	}
	rn = i.unseen[r]
	delete(i.unseen, r)
	i.seen[r] = rn
	return r, rn
}

func FindDistances(edges []*SizeEdge, start ID) map[ID]Size {
	index := newSizeIndex(edges)

	if _, ok := index.unseen[start]; !ok {
		return nil
	}

	index.unseen[start].distance = 0

	for len(index.unseen) > 0 {
		id, n := index.removeClosest()
		//log.Printf("Closest: %d / %d", id, n.distance)
		for _, e := range n.edges {
			peerID, _ := e.Peer(id)
			if peer, ok := index.unseen[peerID]; ok {
				//log.Printf("   Eval: %s", e)

				if n.distance+e.Size < peer.distance {
					peer.distance = n.distance + e.Size
					peer.path = append(n.path, peerID)
					//log.Printf("    distance[%d] = %d", peerID, peer.distance)
				}
			}
		}
	}

	out := make(map[ID]Size)
	for id, n := range index.seen {
		out[id] = n.distance
	}
	return out
}

func ParseNodePaths(line string, delim string) ([]*SizeEdge, error) {
	x := strings.Split(strings.TrimSpace(line), delim)

	left, err := strconv.ParseUint(x[0], 10, 64)
	if err != nil {
		return nil, err
	}
	var r []*SizeEdge

	for _, spec := range x[1:] {
		p := strings.Split(strings.TrimSpace(spec), ",")
		right, err := strconv.ParseUint(p[0], 10, 64)
		if err != nil {
			return nil, err
		}
		size, err := strconv.ParseInt(p[1], 10, 32)
		if err != nil {
			return nil, err
		}
		e := &SizeEdge{&Edge{ID(left), ID(right)}, Size(size)}
		r = append(r, e)
	}

	return r, nil
}

type rankable interface {
	rank() int
	String() string
}

type minheap struct {
	s []rankable
}

func (h *minheap) size() int { return len(h.s) }

func (h *minheap) insert(value rankable) {
	k := len(h.s)
	h.s = append(h.s, value)

	for k != 0 {
		parent := (k - 1) / 2
		if h.inOrder(parent, k) {
			break
		}
		h.swap(parent, k)
		k = parent
	}
}

func (h *minheap) extract() (rankable, bool) {
	if len(h.s) == 0 {
		return nil, false
	}
	value := h.s[0]

	log.Printf("Extract: %s", value)

	max := len(h.s) - 1
	h.swap(0, max)
	h.s = h.s[0:max]

	k := 0
	for x := 0; x < 10; x++ {
		l := 2*k + 1
		r := 2*k + 2

		log.Printf("K: %d Left: s[%d] Right: s[%d] max: %d", k, l, r, len(h.s))

		largest := k
		if l < len(h.s) && !h.inOrder(largest, l) {
			largest = l
		}
		if r < len(h.s) && !h.inOrder(largest, r) {
			largest = r
		}
		if largest == k {
			break
		}
		h.swap(k, largest)
		k = largest
	}
	return value, true
}

func (h *minheap) swap(i, j int) {
	h.s[i], h.s[j] = h.s[j], h.s[i]
}

func (h *minheap) inOrder(parent, child int) bool {
	return h.s[parent].rank() < h.s[child].rank()
}
