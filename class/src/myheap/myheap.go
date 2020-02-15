package myheap

import "sort"

type Interface interface {
	sort.Interface
	Push(x interface{})
	Pop() interface{}
}

func Init(h Interface) {
	for i := h.Len() - 1; i >= 0; i-- {
		j := (i - 1) / 2 // parent of i
		if h.Less(i, j) {
			h.Swap(i, j)
		}
	}
}

func Fix(h Interface, i int) {
	down(h, i, h.Len())
	up(h, i)
}

func up(h Interface, j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || h.Less(i, j) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

func down(h Interface, i int, n int) {
	for {
		j1 := 2*i + 1          // left child
		if j1 >= n || j1 < 0 { // int overflow?
			break
		}
		j := j1
		j2 := 2*i + 2 // right child
		if j2 < n && h.Less(j2, j1) {
			j = j2
		}
		if h.Less(i, j) {
			break
		}
		h.Swap(i, j)
		i = j
	}
}

func Push(h Interface, v interface{}) {
	h.Push(v)
	up(h, h.Len()-1)
}

func Pop(h Interface) interface{} {
	n := h.Len() - 1
	h.Swap(0, n)
	down(h, 0, n)
	return h.Pop()
}

type IntHeap []int64

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int64))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
