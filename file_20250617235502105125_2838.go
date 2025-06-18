package main

import (
	"fmt"
)

func DamerauLevenshtein(s1, s2 string) int {
	len1 := len(s1)
	len2 := len(s2)

	d := make([][]int, len1+1)
	for i := range d {
		d[i] = make([]int, len2+1)
	}

	for i := 0; i <= len1; i++ {
		d[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		d[0][j] = j
	}

	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			deletion := d[i-1][j] + 1
			insertion := d[i][j-1] + 1
			substitution := d[i-1][j-1] + cost

			d[i][j] = min(deletion, min(insertion, substitution))

			if i > 1 && j > 1 && s1[i-1] == s2[j-2] && s1[i-2] == s2[j-1] {
				transposition := d[i-2][j-2] + 1
				d[i][j] = min(d[i][j], transposition)
			}
		}
	}

	return d[len1][len2]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println("kitten vs sitting:", DamerauLevenshtein("kitten", "sitting"))
	fmt.Println("saturday vs sunday:", DamerauLevenshtein("saturday", "sunday"))
	fmt.Println("ab vs ba:", DamerauLevenshtein("ab", "ba"))
	fmt.Println("abc vs acb:", DamerauLevenshtein("abc", "acb"))
	fmt.Println("apple vs apply:", DamerauLevenshtein("apple", "apply"))
	fmt.Println("pale vs ple:", DamerauLevenshtein("pale", "ple"))
	fmt.Println("test vs test:", DamerauLevenshtein("test", "test"))
	fmt.Println(" vs test:", DamerauLevenshtein("", "test"))
	fmt.Println("test vs :", DamerauLevenshtein("test", ""))
	fmt.Println("go vs og:", DamerauLevenshtein("go", "og"))
}

// Additional implementation at 2025-06-17 23:56:06
package main

func minInt(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}
	minVal := nums[0]
	for _, num := range nums {
		if num < minVal {
			minVal = num
		}
	}
	return minVal
}

func DamerauLevenshteinDistance(s1, s2 string) int {
	return DamerauLevenshteinWithCosts(s1, s2, 1, 1, 1, 1)
}

func DamerauLevenshteinWithCosts(s1, s2 string, insertCost, deleteCost, substituteCost, transposeCost int) int {
	len1 := len(s1)
	len2 := len(s2)

	d := make([][]int, len1+1)
	for i := range d {
		d[i] = make([]int, len2+1)
	}

	for i := 0; i <= len1; i++ {
		d[i][0] = i * deleteCost
	}
	for j := 0; j <= len2; j++ {
		d[0][j] = j * insertCost
	}

	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := substituteCost
			if s1[i-1] == s2[j-1] {
				cost = 0
			}

			d[i][j] = minInt(
				d[i-1][j]+deleteCost,
				d[i][j-1]+insertCost,
				d[i-1][j-1]+cost,
			)

			if i > 1 && j > 1 && s1[i-1] == s2[j-2] && s1[i-2] == s2[j-1] {
				d[i][j] = minInt(
					d[i][j],
					d[i-2][j-2]+transposeCost,
				)
			}
		}
	}

	return d[len1][len2]
}