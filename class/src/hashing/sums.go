package hashing

import "log"

type bucket []int

type table []bucket

func newTable(size int) table {
	var bits uint32
	for size > 0 {
		size = size >> 1
		bits++
	}
	return table(make([]bucket, 1<<(bits+2)))
}

func (t table) hash(key int64) int {
	key = ^key + (key << 18) // key = (key << 18) - key - 1;
	key = key ^ (key >> 31)
	key = key * 21 // key = (key + (key << 2)) + (key << 4);
	key = key ^ (key >> 11)
	key = key + (key << 6)
	key = key ^ (key >> 22)
	return int(key) & (len(t) - 1)
}

func (t table) set(v int64, id int) (first bool) {
	h := t.hash(v)
	t[h] = append(t[h], id)
	return len(t[h]) == 1
}

func (t table) get(v int64) (bucket, bool) {
	h := t.hash(v)
	if t[h] == nil {
		return nil, false
	}
	return t[h], true
}

func CountBuckets(nums []int64) (used int, total int) {
	table := newTable(len(nums))
	log.Printf("Table: %d", len(table))

	var buckets int
	for i, n := range nums {
		first := table.set(n, i)
		if first {
			buckets++
		}
	}
	return buckets, len(table)
}

func SpecialSums(nums []int64, start, end int64) int {
	h := make(map[int64]bucket)
	for i, n := range nums {
		key := n / 1000000
		h[key] = append(h[key], i)
	}

	result := make(map[int]bool)
	var total int
	for k, b1 := range h {
		if b2, ok := h[-k]; ok {
			for _, j := range b1 {
				for _, k := range b2 {
					t := nums[j] + nums[k]
					if t > start && t < end {
						//log.Printf("    Pair (%d, %d) => %d", nums[j], nums[k], t)
						result[int(t)] = true
						if j < k {
							total++
						}
					}
				}
			}
		}
	}
	return len(result)
}
