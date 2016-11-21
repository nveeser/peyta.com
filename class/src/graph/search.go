package graph

type edgeset []*Edge

type searchIndex struct {
	edgesByID map[ID]edgeset
	seen      map[ID]bool
}

func newSearchIndex(edges []*Edge) *searchIndex {
	edgesByID := make(map[ID]edgeset)
	for _, e := range edges {
		edgesByID[e.Left] = append(edgesByID[e.Left], e)
		edgesByID[e.Right] = append(edgesByID[e.Right], e)
	}
	return &searchIndex{
		edgesByID: edgesByID,
		seen:      make(map[ID]bool),
	}
}

type visitPair func(to, from ID)

type visitNode func(id ID)

func (i *searchIndex) reset() {
	i.seen = make(map[ID]bool)
}

func (i *searchIndex) bfs(start ID, visit visitPair) {
	q := &queue{}
	from := make(map[ID]ID)

	q.push(start)
	i.seen[start] = true

	for {
		id, ok := q.pop()
		if !ok {
			break
		}

		//log.Printf("Queue: %v", q.l)
		for _, e := range i.edgesByID[id] {
			peerID, _ := e.Peer(id)
			if !i.seen[peerID] {
				q.push(peerID)
				i.seen[peerID] = true

				from[peerID] = id
			}
		}
		visit(id, from[id])
	}
}

func (i *searchIndex) dfs(id ID, d Direction, visit visitNode) {
	if i.seen[id] {
		return
	}
	i.seen[id] = true
	for _, e := range i.edgesByID[id] {
		if id, dir := e.Peer(id); dir == d {
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
