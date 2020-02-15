package myheap

import "testing"

func TestMinHeap(t *testing.T) {
	input := []int64{5, 3, 9, 2, 10, 1}

	h := &IntHeap{}
	for _, n := range input {
		t.Logf("Add %d", n)
		Push(h, n)
		t.Logf("Heap: %v", h)
	}

	t.Logf("Array: %v", h)

	expected := []int64{1, 2, 3, 5, 9, 10}
	for i, want := range expected {
		if h.Len() == 0 {
			t.Fatalf("got empty heap")
		}
		t.Logf("Heap %v", h)
		v := Pop(h).(int64)
		if v != want {
			t.Errorf("Pop(%d) got %d wanted %d", i, v, want)
		}
	}
}

func TestInitHeap(t *testing.T) {
	input := []int64{5, 3, 9, 2, 10, 1}

	h := IntHeap(input)
	Init(&h)

	expected := []int64{1, 2, 3, 5, 9, 10}
	for i, want := range expected {
		if h.Len() == 0 {
			t.Fatalf("got empty heap")
		}
		t.Logf("Heap %v", h)
		v := Pop(&h).(int64)
		if v != want {
			t.Errorf("Pop(%d) got %d wanted %d", i, v, want)
		}
	}
}
