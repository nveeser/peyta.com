package graph2

import "log"

type Edge struct {
	Left, Right uint64
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

type node struct {
	id    uint64
	links []uint64
}

type bfs struct {
	*node
	seen     bool
	distance int
}

func BFS(edges []*Edge, first uint64) {
	nodes := make(map[uint64]*bfs)

	memo := func(id uint64) *bfs {
		if _, ok := nodes[id]; !ok {
			nodes[id] = &bfs{node: &node{id: id}}
		}
		return nodes[id]
	}

	for _, e := range edges {
		left := memo(e.Left)
		left.links = append(left.links, e.Right)

		right := memo(e.Right)
		right.links = append(right.links, e.Left)
	}

	q := &queue{}
	current := nodes[first]

	for current != nil {
		log.Printf("Visit: %d (%d)", current.id, current.distance)
		log.Printf("Queue: %v", q.l)
		for _, id := range current.node.links {
			peer := nodes[id]
			if !peer.seen {
				peer.distance = current.distance + 1
				q.push(id)
			}
		}
		log.Printf("Queue: %v", q.l)
		current.seen = true
		current = nil

		if next, ok := q.pop(); ok {
			current = nodes[next]
		}
	}
}
