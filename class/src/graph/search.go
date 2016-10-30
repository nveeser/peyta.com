package graph

import "log"

type searchIndex struct {
	nodes map[uint64]*snode
}

type snode struct {
	*node
	seen bool
}

type direction bool

const out direction = true
const in direction = false

func (n *snode) peer(e *Edge) (uint64, direction) {
	switch {
	case n.id == e.Left:
		return e.Right, out
	case n.id == e.Right:
		return e.Left, in
	default:
		log.Fatalf("bfs node has edge it is not in: %d -> %s", n.id, e)
	}
	return 0, in
}

func newIndex(edges []*Edge) *searchIndex {
	nodes := make(map[uint64]*snode)

	memo := func(id uint64) *snode {
		if _, ok := nodes[id]; !ok {
			nodes[id] = &snode{node: &node{id: id}}
		}
		return nodes[id]
	}

	for _, e := range edges {
		left := memo(e.Left)
		left.edges = append(left.edges, e)

		right := memo(e.Right)
		right.edges = append(right.edges, e)
	}
	return &searchIndex{nodes}
}

type visitPair func(n, from uint64)

type visitNode func(id uint64)

func (i *searchIndex) reset() {
	for _, n := range i.nodes {
		n.seen = false
	}
}

func (i *searchIndex) bfs(id uint64, visit visitPair) {
	q := &queue{}
	q.push(id)
	for {
		id, ok := q.pop()
		if !ok {
			break
		}
		//log.Printf("Queue: %v", q.l)
		n := i.nodes[id]
		for _, e := range n.edges {
			peerID, _ := n.peer(e)
			peer := i.nodes[peerID]
			if !peer.seen {
				visit(peer.id, n.id)
				q.push(peerID)
				peer.seen = true
			}
		}
	}
}

func (i *searchIndex) dfs(id uint64, d direction, visit visitNode) {
	node := i.nodes[id]
	if node.seen {
		return
	}
	node.seen = true
	for _, e := range node.edges {
		if id, dir := node.peer(e); dir == d {
			i.dfs(id, d, visit)
		}
	}
	visit(id)
}

type queue struct {
	l []uint64
}

func (q *queue) push(e uint64) {
	q.l = append(q.l, e)
}

func (q *queue) pop() (uint64, bool) {
	if len(q.l) == 0 {
		return 0, false
	}
	var n uint64
	n, q.l = q.l[0], q.l[1:]
	return n, true
}
