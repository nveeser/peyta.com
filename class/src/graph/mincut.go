package graph

// ----------------------------------------
// MinCut
// ----------------------------------------

func MinCut(edges []*Edge, count int) []*Edge {
	var result []*Edge
	for i := 0; i < count; i++ {
		min := contract(edges)
		if result == nil || len(result) > len(min) {
			result = min
		}
	}
	return result
}
