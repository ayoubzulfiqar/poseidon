package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Sample(data []interface{}, k int) []interface{} {
	n := len(data)

	if n == 0 || k <= 0 {
		return []interface{}{}
	}

	if k >= n {
		shuffledData := make([]interface{}, n)
		copy(shuffledData, data)
		rand.Shuffle(n, func(i, j int) {
			shuffledData[i], shuffledData[j] = shuffledData[j], shuffledData[i]
		})
		return shuffledData
	}

	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}

	sampled := make([]interface{}, k)
	for i := 0; i < k; i++ {
		r := i + rand.Intn(n-i)
		indices[i], indices[r] = indices[r], indices[i]
		sampled[i] = data[indices[i]]
	}

	return sampled
}

func main() {
	items := []interface{}{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape", "honeydew", "kiwi", "lemon"}

	fmt.Println("Original items:", items)

	sample1 := Sample(items, 3)
	fmt.Println("Sample 1 (k=3):", sample1)

	sample2 := Sample(items, 5)
	fmt.Println("Sample 2 (k=5):", sample2)

	sample3 := Sample(items, len(items))
	fmt.Println("Sample 3 (k=len(items)):", sample3)

	sample4 := Sample(items, 15)
	fmt.Println("Sample 4 (k=15):", sample4)

	sample5 := Sample(items, 0)
	fmt.Println("Sample 5 (k=0):", sample5)

	numbers := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println("\nOriginal numbers:", numbers)
	sampleNumbers := Sample(numbers, 4)
	fmt.Println("Sample numbers (k=4):", sampleNumbers)

	emptyItems := []interface{}{}
	fmt.Println("\nOriginal empty items:", emptyItems)
	sampleEmpty := Sample(emptyItems, 2)
	fmt.Println("Sample empty (k=2):", sampleEmpty)

	singleItem := []interface{}{"only one"}
	fmt.Println("\nOriginal single item:", singleItem)
	sampleSingle1 := Sample(singleItem, 1)
	fmt.Println("Sample single (k=1):", sampleSingle1)
	sampleSingle2 := Sample(singleItem, 5)
	fmt.Println("Sample single (k=5):", sampleSingle2)
}

// Additional implementation at 2025-08-04 07:16:45
package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// Sampler holds a random number generator for reproducible sampling.
type Sampler struct {
	rng *rand.Rand
}

// NewSampler creates a new Sampler instance.
// If seed is 0, it uses the current time as a seed, otherwise it uses the provided seed.
func NewSampler(seed int64) *Sampler {
	var source rand.Source
	if seed == 0 {
		source = rand.NewSource(time.Now().UnixNano())
	} else {
		source = rand.NewSource(seed)
	}
	return &Sampler{
		rng: rand.New(source),
	}
}

// SampleOne returns a single random element from the given slice.
// Returns nil if the slice is empty.
func (s *Sampler) SampleOne(data []interface{}) interface{} {
	if len(data) == 0 {
		return nil
	}
	return data[s.rng.Intn(len(data))]
}

// SampleK returns k unique random elements from the given slice without replacement.
// It uses the Fisher-Yates shuffle algorithm to efficiently select k elements.
// Returns an error if k is negative or greater than the number of available elements.
func (s *Sampler) SampleK(data []interface{}, k int) ([]interface{}, error) {
	if k < 0 {
		return nil, errors.New("k cannot be negative")
	}
	if k == 0 {
		return []interface{}{}, nil
	}
	if k > len(data) {
		return nil, fmt.Errorf("k (%d) cannot be greater than the number of elements (%d)", k, len(data))
	}

	// Create a copy to avoid modifying the original slice
	shuffledData := make([]interface{}, len(data))
	copy(shuffledData, data)

	// Perform a partial Fisher-Yates shuffle for the first k elements
	for i := 0; i < k; i++ {
		// Pick a random index from the remaining unshuffled part (from i to end)
		j := s.rng.Intn(len(shuffledData)-i) + i
		// Swap the current element with the randomly chosen element
		shuffledData[i], shuffledData[j] = shuffledData[j], shuffledData[i]
	}

	return shuffledData[:k], nil
}

// WeightedItem represents an item with an associated weight.
type WeightedItem struct {
	Value  interface{}
	Weight float64
}

// SampleWeighted returns a single random element based on its weight.
// The probability of an item being selected is proportional to its weight.
// Returns an error if the slice is empty or if all weights sum to zero or less.
func (s *Sampler) SampleWeighted(items []WeightedItem) (interface{}, error) {
	if len(items) == 0 {
		return nil, errors.New("cannot sample from an empty slice of weighted items")
	}

	totalWeight := 0.0
	for _, item := range items {
		totalWeight += item.Weight
	}

	if totalWeight <= 0 {
		return nil, errors.New("total weight must be greater than zero")
	}

	// Generate a random number between 0 (inclusive) and totalWeight (exclusive)
	r := s.rng.Float64() * totalWeight

	// Iterate through items, accumulating weights until r falls within an item's range
	currentWeight := 0.0
	for _, item := range items {
		currentWeight += item.Weight
		if r < currentWeight { // Use < for consistency with Float64() which is [0, 1)
			return item.Value, nil
		}
	}

	// This part should theoretically not be reached if totalWeight > 0 and r is within [0, totalWeight).
	// It's a fallback for extreme floating point precision edge cases, returning the last item.
	return items[len(items)-1].Value, nil
}

func main() {
	// Initialize Sampler with a fixed seed for reproducible results.
	// Use 0 for a truly random seed based on current time.
	sampler := NewSampler(42) // Using a fixed seed for demonstration

	fmt.Println("--- SampleOne ---")
	data := []interface{}{"apple", "banana", "cherry", "date", "elderberry"}
	for i := 0; i < 5; i++ {
		fmt.Printf("Random fruit: %v\n", sampler.SampleOne(data))
	}
	fmt.Println()

	fmt.Println("--- SampleK (without replacement) ---")
	numbers := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sampleSize := 3
	sampledNumbers, err := sampler.SampleK(numbers, sampleSize)
	if err != nil {
		fmt.Printf("Error sampling K: %v\n", err)
	} else {
		fmt.Printf("Sampled %d unique numbers: %v\n", sampleSize, sampledNumbers)
	}

	// Test error case: k > len(data)
	sampleSize = 12
	sampledNumbers, err = sampler.SampleK(numbers, sampleSize)
	if err != nil {
		fmt.Printf("Error sampling K (expected): %v\n", err)
	}
	// Test error case: k < 0
	sampleSize = -1
	sampledNumbers, err = sampler.SampleK(numbers, sampleSize)
	if err != nil {
		fmt.Printf("Error sampling K (expected): %v\n", err)
	}
	fmt.Println()

	fmt.Println("--- SampleWeighted ---")
	weightedItems := []WeightedItem{
		{Value: "common", Weight: 70},
		{Value: "uncommon", Weight: 20},
		{Value: "rare", Weight: 10},
		{Value: "legendary", Weight: 0.5},
	}

	// Collect results to observe distribution
	results := make(map[interface{}]int)
	numSamples := 100000
	for i := 0; i < numSamples; i++ {
		item, err := sampler.SampleWeighted(weightedItems)
		if err != nil {
			fmt.Printf("Error sampling weighted: %v\n", err)
			break
		}
		results[item]++
	}

	fmt.Printf("Weighted sample distribution over %d samples:\n", numSamples)
	// Sort keys for consistent output, assuming keys are string-representable
	var keys []interface{}
	for k := range results {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprintf("%v", keys[i]) < fmt.Sprintf("%v", keys[j])
	})

	for _, k := range keys {
		fmt.Printf("  %v: %d (%.2f%%)\n", k, results[k], float64(results[k])/float64(numSamples)*100)
	}

	// Test error case for weighted sampling: empty slice
	emptyWeightedItems := []WeightedItem{}
	_, err = sampler.SampleWeighted(emptyWeightedItems)
	if err != nil {
		fmt.Printf("Error sampling weighted (expected empty): %v\n", err)
	}

	// Test error case for weighted sampling: zero or negative total weight
	zeroWeightItems := []WeightedItem{
		{Value: "zero", Weight: 0},
		{Value: "negative", Weight: -5},
	}
	_, err = sampler.SampleWeighted(zeroWeightItems)
	if err != nil {
		fmt.Printf("Error sampling weighted (expected zero/negative total weight): %v\n", err)
	}
}

// Additional implementation at 2025-08-04 07:17:22


// Additional implementation at 2025-08-04 07:18:56
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Sampler holds the random number generator for reproducible sampling.
type Sampler struct {
	rng *rand.Rand
}

// NewSampler creates a new Sampler with an optional seed.
// If seed is 0, it uses the current time in nanoseconds as the seed,
// providing a different sequence of random numbers each time the program runs.
// If a non-zero seed is provided, the sequence of random numbers will be
// deterministic for that seed, useful for reproducible results.
func NewSampler(seed int64) *Sampler {
	var source rand.Source
	if seed == 0 {
		source = rand.NewSource(time.Now().UnixNano())
	} else {
		source = rand.NewSource(seed)
	}
	return &Sampler{
		rng: rand.New(source),
	}
}

// SampleWithoutReplacement takes a slice of interface{} and returns a new slice
// containing 'count' randomly selected unique elements from the input slice.
// It returns an error if count is negative or greater than the length of the input slice.
func (s *Sampler) SampleWithoutReplacement(data []interface{}, count int) ([]interface{}, error) {
	if count < 0 {
		return nil, fmt.Errorf("sample count cannot be negative: %d", count)
	}
	if count > len(data) {
		return nil, fmt.Errorf("sample count (%d) cannot be greater than data length (%d) for sampling without replacement", count, len(data))
	}

	if count == 0 {
		return []interface{}{}, nil
	}
	if count == len(data) {
		// If sampling all elements, just shuffle and return a copy
		shuffled := make([]interface{}, len(data))
		copy(shuffled, data)
		s.rng.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		return shuffled, nil
	}

	// Perform a partial Fisher-Yates shuffle to select 'count' unique elements.
	// Create a copy to avoid modifying the original slice.
	temp := make([]interface{}, len(data))
	copy(temp, data)

	result := make([]interface{}, count)
	for i := 0; i < count; i++ {
		// Pick a random index from the remaining unsampled elements (from i to len(temp)-1)
		idx := s.rng.Intn(len(temp)-i) + i
		result[i] = temp[idx]
		// Swap the picked element with the element at the current position 'i'.
		// This effectively moves the picked element to the "sampled" part of the slice
		// and ensures it won't be picked again.
		temp[idx], temp[i] = temp[i], temp[idx]
	}
	return result, nil
}

// SampleWithReplacement takes a slice of interface{} and returns a new slice
// containing 'count' randomly selected elements from the input slice.
// Elements can be selected multiple times.
// It returns an error if count is negative or if the input data slice is empty.
func (s *Sampler) SampleWithReplacement(data []interface{}, count int) ([]interface{}, error) {
	if count < 0 {
		return nil, fmt.Errorf("sample count cannot be negative: %d", count)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("cannot sample from an empty data slice")
	}

	result := make([]interface{}, count)
	for i := 0; i < count; i++ {
		idx := s.rng.Intn(len(data))
		result[i] = data[idx]
	}
	return result, nil
}

// SampleByPercentageWithoutReplacement samples a percentage of elements without replacement.
// The percentage should be between 0.0 and 100.0.
// The number of elements to sample is calculated as (percentage / 100.0) * len(data).
// The result is rounded down to the nearest integer.
func (s *Sampler) SampleByPercentageWithoutReplacement(data []interface{}, percentage float64) ([]interface{}, error) {
	if percentage < 0.0 || percentage > 100.0 {
		return nil, fmt.Errorf("percentage must be between 0.0 and 100.0, got %.2f", percentage)
	}
	if len(data) == 0 {
		return []interface{}{}, nil // Sampling from empty data results in empty sample
	}

	count := int(float64(len(data)) * percentage / 100.0)
	return s.SampleWithoutReplacement(data, count)
}

// SampleByPercentageWithReplacement samples a percentage of elements with replacement.
// The percentage should be between 0.0 and 100.0.
// The number of elements to sample is calculated as (percentage / 100.0) * len(data).
// The result is rounded down to the nearest integer.
// Note: For sampling with replacement, the effective count can be greater than len(data)
// if the percentage implies it (e.g., 200% of 10 items means 20 items).
func (s *Sampler) SampleByPercentageWithReplacement(data []interface{}, percentage float64) ([]interface{}, error) {
	if percentage < 0.0 || percentage > 100.0 { // Percentage can be > 100 for replacement, but for simplicity, let's keep it 0-100 for now.
		// If the intent is to allow >100% for replacement, remove this check or adjust.
		// For typical "percentage of data" use, 0-100 is standard.
		return nil, fmt.Errorf("percentage must be between 0.0 and 100.0, got %.2f", percentage)
	}
	if len(data) == 0 {
		return []interface{}{}, nil // Sampling from empty data results in empty sample
	}

	count := int(float64(len(data)) * percentage / 100.0)
	return s.SampleWithReplacement(data, count)
}

func main() {
	// Example data slices
	items := []interface{}{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape", "honeydew", "kiwi", "lemon"}
	numbers := []interface{}{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150}
	emptyData := []interface{}{}

	fmt.Println("Original Items:", items)
	fmt.Println("Original Numbers:", numbers)
	fmt.Println("--------------------------------------------------")

	// Create a sampler with a fixed seed for reproducible results
	sampler := NewSampler(42) // Using seed 42

	// --- Demonstrate SampleWithoutReplacement ---
	fmt.Println("--- Sampling Without Replacement (Fixed Count) ---")
	sampleSize := 3
	sampledItems, err := sampler.SampleWithoutReplacement(items, sampleSize)
	if err != nil {
		fmt.Printf("Error sampling items: %v\n", err)
	} else {
		fmt.Printf("Sampled %d items: %v\n", sampleSize, sampledItems)
	}

	// Test edge case: count > len(data) for without replacement
	_, err = sampler.SampleWithoutReplacement(items, 20)
	if err != nil {
		fmt.Printf("Expected error for count > len(data): %v\n", err)
	}

	// Test sampling 0 elements
	sampledZero, err := sampler.SampleWithoutReplacement(items, 0)
	if err != nil {
		fmt.Printf("Error sampling zero items: %v\n", err)
	} else {
		fmt.Printf("Sampled 0 items: %v (length %d)\n", sampledZero, len(sampledZero))
	}

	// Test sampling from empty data
	sampledFromEmpty, err := sampler.SampleWithoutReplacement(emptyData, 0)
	if err != nil {
		fmt.Printf("Error sampling from empty data: %v\n", err)
	} else {
		fmt.Printf("Sampled 0 items from empty data: %v (length %d)\n", sampledFromEmpty, len(sampledFromEmpty))
	}
	_, err = sampler.SampleWithoutReplacement(emptyData, 1)
	if err != nil {
		fmt.Printf("Expected error sampling >0 from empty data: %v\n", err)
	}

	fmt.Println("--------------------------------------------------")

	// --- Demonstrate SampleByPercentageWithoutReplacement ---
	fmt.Println("--- Sampling Without Replacement (By Percentage) ---")
	samplePercentage := 30.0 // 30% of 10 items = 3 items
	sampledNumbersByPercent, err := sampler.SampleByPercentageWithoutReplacement(numbers, samplePercentage)
	if err != nil {
		fmt.Printf("Error sampling numbers by percentage: %v\n", err)
	} else {
		fmt.Printf("Sampled %.1f%% of numbers: %v (length %d)\n", samplePercentage, sampledNumbersByPercent, len(sampledNumbersByPercent))
	}

	samplePercentage = 5.0 // 5% of 15 items = 0.75 -> 0 items
	sampledNumbersBySmallPercent, err := sampler.SampleByPercentageWithoutReplacement(numbers, samplePercentage)
	if err != nil {
		fmt.Printf("Error sampling numbers by small percentage: %v\n", err)
	} else {
		fmt.Printf("Sampled %.1f%% of numbers: %v (length %d)\n", samplePercentage, sampledNumbersBySmallPercent, len(sampledNumbersBySmallPercent))
	}

	// Test percentage out of range
	_, err = sampler.SampleByPercentageWithoutReplacement(items, 110.0)
	if err != nil {
		fmt.Printf("Expected error for percentage > 100: %v\n", err)
	}

	fmt.Println("--------------------------------------------------")

	// --- Demonstrate SampleWithReplacement ---
	fmt.Println("--- Sampling With Replacement (Fixed Count) ---")
	sampleSize = 5
	sampledItemsWithReplacement, err := sampler.SampleWithReplacement(items, sampleSize)
	if err != nil {
		fmt.Printf("Error sampling items with replacement: %v\n", err)
	} else {
		fmt.Printf("Sampled %d items (with replacement): %v\n", sampleSize, sampledItemsWithReplacement)
	}

	// Test sampling more than original length with replacement
	sampleSize = 15
	sampledItemsMore, err := sampler.SampleWithReplacement(items, sampleSize)
	if err != nil {
		fmt.Printf("Error sampling more items with replacement: %v\n", err)
	} else {
		fmt.Printf("Sampled %d items (with replacement, > original length): %v\n", sampleSize, sampledItemsMore)
	}

	// Test sampling from empty data
	_, err = sampler.SampleWithReplacement(emptyData, 5)
	if err != nil {
		fmt.Printf("Expected error sampling from empty data: %v\n", err)
	}

	fmt.Println("--------------------------------------------------")

	// --- Demonstrate SampleByPercentageWithReplacement ---
	fmt.Println("--- Sampling With Replacement (By Percentage) ---")
	samplePercentage = 50.0 // 50% of 15 numbers = 7 items
	sampledNumbersWithReplacementByPercent, err := sampler.SampleByPercentageWithReplacement(numbers, samplePercentage)
	if err != nil {
		fmt.Printf("Error sampling numbers with replacement by percentage: %v\n", err)
	} else {
		fmt.Printf("Sampled %.1f%% of numbers (with replacement): %v (length %d)\n", samplePercentage, sampledNumbersWithReplacementByPercent, len(sampledNumbersWithReplacementByPercent))
	}

	fmt.Println("--------------------------------------------------")

	// --- Demonstrating different seeds ---
	fmt.Println("--- Demonstrating different seeds ---")
	sampler2 := NewSampler(0) // Uses time.Now().UnixNano(), likely different each run
	sampledItems2, err := sampler2.SampleWithoutReplacement(items, 3)
	if err != nil {
		fmt.Printf("Error sampling items (sampler2): %v\n", err)
	} else {
		fmt.Printf("Sampled 3 items (sampler2, different seed): %v\n", sampledItems2)
	}

	sampler3 := NewSampler(42) // Same seed as the first sampler
	sampledItems3, err := sampler3.SampleWithoutReplacement(items, 3)
	if err != nil {
		fmt.Printf("Error sampling items (sampler3): %v\n", err)
	} else {
		fmt.Printf("Sampled 3 items (sampler3, same seed as first): %v\n", sampledItems3)
	}
}