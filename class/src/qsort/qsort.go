package qsort

import (
	"log"
	"strings"
)

// SetupPivot is the function which finds the next pivot element and
// moves it to the the first slot.
type SetupPivot func(v []uint64) int

func FirstElement(v []uint64) int {
	return 0
}

func LastElement(v []uint64) int {
	p := len(v) - 1
	return p
}

func Median3(v []uint64) int {
	l := len(v) - 1
	m := l / 2
	//log.Printf(" Median: v[%d]=%d v[%d]=%d v[%d]=%d", 0, v[0], m, v[m], l, v[l])
	switch {
	case v[0] < v[l]:
		switch {
		case v[m] < v[0]:
			return 0
		case v[l] < v[m]:
			return l
		default:
			return m
		}
	default: // v[l] < v[0]
		switch {
		case v[m] < v[l]:
			return l
		case v[0] < v[m]:
			return 0
		default:
			return m
		}
	}
	return 0
}

func Sort(v []uint64, p SetupPivot) int {
	s := sorter{p}
	return s.sort(v, 1)
}

type sorter struct {
	findPivot SetupPivot
}

func (s sorter) log(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (s sorter) sort(v []uint64, depth int) int {
	f := strings.Repeat("  ", depth)

	s.log("%sInput: %v", f, v)
	if len(v) <= 1 {
		return 0
	}

	//log.Printf("%s MedianPivot: v[%d]=%d", f, v[p], p)
	Desc(v).Swap(0, s.findPivot(v))

	p := 0

	count := int(0)
	//s.log("%sPivot: %v", f, v)
	for j := 1; j < len(v); j++ {
		s.log("%s Compare I = %d P = %d V = %v", f, j, p, v)
		if v[j] < v[0] {
			p++
			if p != j {
				s.log("%s Swap v[%d] = %d <-> v[%d] = %d", f, p, v[p], j, v[j])
				v[p], v[j] = v[j], v[p]
				s.log("%s   Done p=%d V: %v", f, p, v)
			}
		}
	}
	v[0], v[p] = v[p], v[0]
	s.log("%s   Recurse V: %v", f, v)
	count += s.sort(v[0:p], depth+1)
	count += s.sort(v[p+1:], depth+1)
	count += (len(v) - 1)
	return count
}

type Desc []uint64

func (s Desc) Len() int           { return len(s) }
func (s Desc) Less(i, j int) bool { return s[i] < s[j] }
func (s Desc) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
