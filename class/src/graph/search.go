package graph

import "log"

type searchIndex struct {
	nodes map[ID]*snode
}

type snode struct {
	*node
	seen bool
}

type direction bool

const out direction = true
const in direction = false

func (n *snode) peer(e *Edge) (ID, direction) {
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
	nodes := make(map[ID]*snode)

	memo := func(id ID) *snode {
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

type visitPair func(to, from ID)

type visitNode func(id ID)

func (i *searchIndex) reset() {
	for _, n := range i.nodes {
		n.seen = false
	}
}

func (i *searchIndex) bfs(start ID, visit visitPair) {
	q := &queue{}
	from := make(map[ID]ID)

	q.push(start)
	i.nodes[start].seen = true

	for {
		id, ok := q.pop()
		if !ok {
			break
		}
		n := i.nodes[id]

		//log.Printf("Queue: %v", q.l)
		for _, e := range n.edges {
			peerID, _ := n.peer(e)
			peer := i.nodes[peerID]
			if !peer.seen {
				q.push(peerID)
				peer.seen = true

				from[peerID] = n.id
			}
		}
		visit(id, from[id])
	}
}

func (i *searchIndex) dfs(id ID, d direction, visit visitNode) {
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
	l []ID
}

func (q *queue) push(e ID) {
	q.l = append(q.l, e)
}

func (q *queue) pop() (ID, bool) {
	if len(q.l) == 0 {
		return 0, false
	}
	var n ID
	n, q.l = q.l[0], q.l[1:]
	return n, true
}
