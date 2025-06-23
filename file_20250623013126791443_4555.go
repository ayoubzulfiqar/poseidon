package main

import (
	"fmt"
)

// maxSlidingWindow calculates the maximum value in each sliding window of size k.
func maxSlidingWindow(nums []int, k int) []int {
	if len(nums) == 0 || k == 0 {
		return []int{}
	}
	if k == 1 {
		return nums
	}

	var result []int
	// deque stores indices of elements in decreasing order of their values.
	// The front of the deque always holds the index of the maximum element in the current window.
	var deque []int

	for i, num := range nums {
		// Remove elements from the front of the deque that are out of the current window.
		if len(deque) > 0 && deque[0] <= i-k {
			deque = deque[1:]
		}

		// Remove elements from the back of the deque that are smaller than or equal to the current element.
		// These elements can no longer be the maximum in any future window that includes the current element.
		for len(deque) > 0 && nums[deque[len(deque)-1]] <= num {
			deque = deque[:len(deque)-1]
		}

		// Add the current element's index to the back of the deque.
		deque = append(deque, i)

		// Once the window is fully formed (i.e., i >= k-1), the maximum for the current window
		// is the element at the index stored at the front of the deque.
		if i >= k-1 {
			result = append(result, nums[deque[0]])
		}
	}

	return result
}

func main() {
	// Test cases
	nums1 := []int{1, 3, -1, -3, 5, 3, 6, 7}
	k1 := 3
	fmt.Println(maxSlidingWindow(nums1, k1)) // Expected: [3 3 5 5 6 7]

	nums2 := []int{1}
	k2 := 1
	fmt.Println(maxSlidingWindow(nums2, k2)) // Expected: [1]

	nums3 := []int{1, -1}
	k3 := 1
	fmt.Println(maxSlidingWindow(nums3, k3)) // Expected: [1 -1]

	nums4 := []int{7, 2, 4}
	k4 := 2
	fmt.Println(maxSlidingWindow(nums4, k4)) // Expected: [7 4]

	nums5 := []int{9, 11}
	k5 := 2
	fmt.Println(maxSlidingWindow(nums5, k5)) // Expected: [11]

	nums6 := []int{4, 3, 2, 1}
	k6 := 2
	fmt.Println(maxSlidingWindow(nums6, k6)) // Expected: [4 3 2]

	nums7 := []int{}
	k7 := 0
	fmt.Println(maxSlidingWindow(nums7, k7)) // Expected: []

	nums8 := []int{1, 2, 3, 4, 5}
	k8 := 5
	fmt.Println(maxSlidingWindow(nums8, k8)) // Expected: [5]
}

// Additional implementation at 2025-06-23 01:32:31
package main

import "fmt"

func maxSlidingWindow(nums []int, k int) []int {
	if len(nums) == 0 || k == 0 {
		return []int{}
	}
	if k == 1 {
		return nums
	}

	var result []int
	var deque []int

	for i, num := range nums {
		if len(deque) > 0 && deque[0] < i-k+1 {
			deque = deque[1:]
		}

		for len(deque) > 0 && nums[deque[len(deque)-1]] <= num {
			deque = deque[:len(deque)-1]
		}

		deque = append(deque, i)

		if i >= k-1 {
			result = append(result, nums[deque[0]])
		}
	}

	return result
}

func main() {
	nums1 := []int{1, 3, -1, -3, 5, 3, 6, 7}
	k1 := 3
	fmt.Println(maxSlidingWindow(nums1, k1))

	nums2 := []int{1}
	k2 := 1
	fmt.Println(maxSlidingWindow(nums2, k2))

	nums3 := []int{1, 2, 3, 4, 5}
	k3 := 2
	fmt.Println(maxSlidingWindow(nums3, k3))

	nums4 := []int{9, 11}
	k4 := 2
	fmt.Println(maxSlidingWindow(nums4, k4))

	nums5 := []int{4, -2}
	k5 := 2
	fmt.Println(maxSlidingWindow(nums5, k5))

	nums6 := []int{}
	k6 := 0
	fmt.Println(maxSlidingWindow(nums6, k6))

	nums7 := []int{1, 3, 1, 2, 0, 5}
	k7 := 3
	fmt.Println(maxSlidingWindow(nums7, k7))
}

// Additional implementation at 2025-06-23 01:33:34
package main

import (
	"fmt"
)

func maxSlidingWindow(nums []int, k int) []int {
	if len(nums) == 0 || k == 0 || k > len(nums) {
		return []int{}
	}

	var res []int
	var deque []int // Stores indices of elements in decreasing order of their values

	for i := 0; i < len(nums); i++ {
		// Remove elements from the front of the deque that are out of the current window
		// The element at deque[0] is the index of the maximum element in the current window.
		// If this index is i - k, it means it's about to leave the window.
		if len(deque) > 0 && deque[0] == i-k {
			deque = deque[1:] // Pop front
		}

		// Remove elements from the back of the deque that are smaller than the current element
		// We want to maintain a decreasing order of values in the deque.
		// If nums[i] is greater than or equal to the value at the last index in the deque,
		// then the element at the last index can never be the maximum in any future window
		// (because nums[i] is larger and comes later). So, remove it.
		for len(deque) > 0 && nums[deque[len(deque)-1]] <= nums[i] {
			deque = deque[:len(deque)-1] // Pop back
		}

		// Add the current element's index to the back of the deque
		deque = append(deque, i) // Push back

		// If the window is fully formed (i.e., i >= k - 1), add the maximum to the result
		// The maximum element for the current window is always at deque[0] (its index).
		if i >= k-1 {
			res = append(res, nums[deque[0]])
		}
	}

	return res
}

func main() {
	nums1 := []int{1, 3, -1, -3, 5, 3, 6, 7}
	k1 := 3
	fmt.Println(maxSlidingWindow(nums1, k1))

	nums2 := []int{1}
	k2 := 1
	fmt.Println(maxSlidingWindow(nums2, k2))

	nums3 := []int{1, -1}
	k3 := 1
	fmt.Println(maxSlidingWindow(nums3, k3))

	nums4 := []int{7, 2, 4}
	k4 := 2
	fmt.Println(maxSlidingWindow(nums4, k4))

	nums5 := []int{1, 2, 3, 4, 5}
	k5 := 3
	fmt.Println(maxSlidingWindow(nums5, k5))

	nums6 := []int{5, 4, 3, 2, 1}
	k6 := 3
	fmt.Println(maxSlidingWindow(nums6, k6))

	nums7 := []int{}
	k7 := 0
	fmt.Println(maxSlidingWindow(nums7, k7))

	nums8 := []int{1, 2, 3}
	k8 := 4
	fmt.Println(maxSlidingWindow(nums8, k8))

	nums9 := []int{1, 3, 1, 2, 0, 5}
	k9 := 3
	fmt.Println(maxSlidingWindow(nums9, k9))
}