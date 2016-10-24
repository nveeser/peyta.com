package qsort

import (
	"fmt"
	"sort"
	"testing"
)

var forum100 = []uint64{
	57, 97, 17, 31, 54, 98, 87, 27, 89,
	81, 18, 70, 3, 34, 63, 100, 46, 30,
	99, 10, 33, 65, 96, 38, 48, 80, 95,
	6, 16, 19, 56, 61, 1, 47, 12, 73, 49,
	41, 37, 40, 59, 67, 93, 26, 75, 44,
	58, 66, 8, 55, 94, 74, 83, 7, 15, 86,
	42, 50, 5, 22, 90, 13, 69, 53, 43, 24,
	92, 51, 23, 39, 78, 85, 4, 25, 52, 36,
	60, 68, 9, 64, 79, 14, 45, 2, 77, 84,
	11, 71, 35, 72, 28, 76, 82, 88, 32, 21,
	20, 91, 62, 29,
}

type testCase struct {
	name           string
	input          []uint64
	expectedFirst  int
	expectedLast   int
	expectedMedian int
}

func TestSort(t *testing.T) {
	cases := []testCase{
		testCase{
			name:           "Sorted",
			input:          []uint64{1, 2, 3, 4, 5},
			expectedFirst:  10,
			expectedLast:   10,
			expectedMedian: 6,
		},
		testCase{
			name:           "FiveUnsorted",
			input:          []uint64{2, 1, 3, 4, 5},
			expectedFirst:  7,
			expectedLast:   10,
			expectedMedian: 6,
		},
		testCase{
			name:           "Reverse",
			input:          []uint64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			expectedFirst:  10,
			expectedLast:   10,
			expectedMedian: 19,
		},
		testCase{
			name:           "ForumCase10",
			input:          []uint64{3, 9, 8, 4, 6, 10, 2, 5, 7, 1},
			expectedFirst:  25,
			expectedLast:   29,
			expectedMedian: 21,
		},
		testCase{
			name:           "ForumCase100",
			input:          forum100,
			expectedFirst:  615,
			expectedLast:   587,
			expectedMedian: 518,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runPivot := func(p SetupPivot, expectedCount int) {
				expected := make([]uint64, len(tc.input))
				copy(expected, tc.input)
				sort.Sort(Desc(expected))

				got := make([]uint64, len(tc.input))
				copy(got, tc.input)
				count := Sort(got, p)

				for i := range got {
					if got[i] != expected[i] {
						t.Errorf("element[%d]: got %v wanted %v", i, got[i], expected[i])
					}
				}
				if count != expectedCount {
					t.Errorf("count wrong: got %d wanted %d", count, expectedCount)
				}
			}
			t.Run("First", func(t *testing.T) { runPivot(FirstElement, tc.expectedFirst) })
			t.Run("Last", func(t *testing.T) { runPivot(LastElement, tc.expectedLast) })
			t.Run("Median", func(t *testing.T) { runPivot(Median3, tc.expectedMedian) })
		})
	}
}

func TestMedian(t *testing.T) {
	cases := []struct {
		input    []uint64
		expected int
	}{
		{
			input:    []uint64{2, 4},
			expected: 0,
		},
		{
			input:    []uint64{8, 2, 4, 5, 7, 1},
			expected: 2,
		},
		{
			input:    []uint64{5, 2, 1, 4, 3},
			expected: 4,
		},
		{
			input:    []uint64{3, 2, 5, 4, 1},
			expected: 0,
		},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			m := Median3(tc.input)
			if m != tc.expected {
				t.Fatalf("Median is wrong got %d wanted %d", m, tc.expected)
			}
		})
	}
}
