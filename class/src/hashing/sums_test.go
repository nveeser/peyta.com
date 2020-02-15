package hashing

import "testing"

func TestTable(t *testing.T) {
	table := newTable(1000)

	table.set(1340, 10)

	b, ok := table.get(1340)
	if !ok || len(b) != 1 || b[0] != 10 {
		t.Errorf("got %v %s wanted {10}, true", b, ok)
	}
}
