package myheap

//
// IntHeap Max / Min
//
type minHeap struct {
	IntHeap
}

func (h minHeap) Less(i, j int) bool { return h.IntHeap[i] < h.IntHeap[j] }

type maxHeap struct {
	IntHeap
}

func (h maxHeap) Less(i, j int) bool { return h.IntHeap[i] > h.IntHeap[j] }

type Median struct {
	smalls maxHeap
	larges minHeap
}

func (m *Median) Add(n int64) {
	switch {
	case m.smalls.Len() > 0 && n < m.smalls.IntHeap[0]:
		Push(&m.smalls, n)
		// log.Printf("Add to Smalls: %v", m.smalls)

	case m.larges.Len() > 0 && n > m.larges.IntHeap[0]:
		Push(&m.larges, n)
		// log.Printf("Add to Larges: %v", m.larges)

	case m.smalls.Len() <= m.larges.Len():
		Push(&m.smalls, n)
		//log.Printf("Add to Smalls: %v", m.smalls)

	default:
		Push(&m.larges, n)
		//log.Printf("Add to Larges: %v", m.larges)
	}

	switch m.smalls.Len() - m.larges.Len() {
	case 2:
		//log.Printf("Rebalance Right: %v %v", m.smalls, m.larges)
		v := Pop(&m.smalls).(int64)
		Push(&m.larges, v)

	case -2:
		//log.Printf("Rebalance Left: %v %v", m.smalls, m.larges)
		v := Pop(&m.larges).(int64)
		Push(&m.smalls, v)
	}
}

func (m *Median) Value() int64 {
	switch {
	case m.smalls.Len() >= m.larges.Len():
		return m.smalls.IntHeap[0]

	case m.larges.Len() > 0:
		return m.larges.IntHeap[0]

	default:
		return 0
	}
}

func (m *Median) Size() int {
	return m.larges.Len() + m.smalls.Len()
}
