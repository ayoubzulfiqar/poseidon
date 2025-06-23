package main

import (
	"fmt"
	"sort"
)

type Interval struct {
	Start int
	End   int
}

func mergeIntervals(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return []Interval{}
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].Start < intervals[j].Start
	})

	merged := []Interval{intervals[0]}

	for i := 1; i < len(intervals); i++ {
		lastMerged := &merged[len(merged)-1]
		current := intervals[i]

		if current.Start <= lastMerged.End {
			if current.End > lastMerged.End {
				lastMerged.End = current.End
			}
		} else {
			merged = append(merged, current)
		}
	}

	return merged
}

func main() {
	// Example 1: Basic merging
	intervals1 := []Interval{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	fmt.Println("Original intervals 1:", intervals1)
	merged1 := mergeIntervals(intervals1)
	fmt.Println("Merged intervals 1:", merged1)

	// Example 2: No overlap
	intervals2 := []Interval{{1, 2}, {3, 4}, {5, 6}}
	fmt.Println("Original intervals 2:", intervals2)
	merged2 := mergeIntervals(intervals2)
	fmt.Println("Merged intervals 2:", merged2)

	// Example 3: Complete overlap
	intervals3 := []Interval{{1, 10}, {2, 5}, {3, 7}}
	fmt.Println("Original intervals 3:", intervals3)
	merged3 := mergeIntervals(intervals3)
	fmt.Println("Merged intervals 3:", merged3)

	// Example 4: Empty input
	intervals4 := []Interval{}
	fmt.Println("Original intervals 4:", intervals4)
	merged4 := mergeIntervals(intervals4)
	fmt.Println("Merged intervals 4:", merged4)

	// Example 5: Single interval
	intervals5 := []Interval{{10, 20}}
	fmt.Println("Original intervals 5:", intervals5)
	merged5 := mergeIntervals(intervals5)
	fmt.Println("Merged intervals 5:", merged5)

	// Example 6: Intervals with same start but different end
	intervals6 := []Interval{{1, 5}, {1, 3}, {6, 8}}
	fmt.Println("Original intervals 6:", intervals6)
	merged6 := mergeIntervals(intervals6)
	fmt.Println("Merged intervals 6:", merged6)

	// Example 7: Intervals that touch
	intervals7 := []Interval{{1, 3}, {3, 5}, {6, 8}}
	fmt.Println("Original intervals 7:", intervals7)
	merged7 := mergeIntervals(intervals7)
	fmt.Println("Merged intervals 7:", merged7)
}

// Additional implementation at 2025-06-22 22:42:50
