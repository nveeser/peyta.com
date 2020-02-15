package inversion

func Invert(a []uint64) ([]uint64, int) {
	switch len(a) {
	case 1:
		return a, 0
	case 2:
		if a[0] > a[1] {
			a[0], a[1] = a[1], a[0]
			return a, 1
		}
		return a, 0
	}

	n := len(a) / 2
	left, linv := Invert(a[0:n])
	right, rinv := Invert(a[n:])

	var sorted []uint64
	inversions := linv + rinv

	//log.Printf("Merge: %v and %v", left, right)

	i := 0
	j := 0
	total := len(left) + len(right)
	for i+j < total {
		//lvalue := ""
		// if i < len(left) {
		// 	lvalue = fmt.Sprintf("%d", left[i])
		// }
		//rvalue := ""
		// if j < len(right) {
		// 	rvalue = fmt.Sprintf("%d", right[j])
		// }
		//log.Printf("Index left[%d] = %s right[%d] = %s", i, lvalue, j, rvalue)

		switch {

		case j == len(right):
			sorted = append(sorted, left[i])
			i++

		case i == len(left):
			sorted = append(sorted, right[j])
			j++

		case left[i] < right[j]:
			sorted = append(sorted, left[i])
			i++

		default:
			sorted = append(sorted, right[j])
			j++
			//log.Printf("Inversions: %d", len(left)-i+1)
			inversions += len(left) - i
		}
	}
	//log.Printf("Sorted: %v (%d)", sorted, inversions)
	return sorted, inversions
}
