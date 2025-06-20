package main

import (
	"fmt"
	"math"
)

func calculateMean(data []float64) float64 {
	if len(data) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, val := range data {
		sum += val
	}
	return sum / float64(len(data))
}

func detectSkewness(data []float64) (float64, string) {
	n := float64(len(data))

	if n < 2 {
		return 0.0, "Not enough data points (requires at least 2) to calculate meaningful skewness."
	}

	mean := calculateMean(data)

	sumOfSquares := 0.0
	sumOfCubes := 0.0

	for _, val := range data {
		diff := val - mean
		sumOfSquares += diff * diff
		sumOfCubes += math.Pow(diff, 3)
	}

	stdDevPop := math.Sqrt(sumOfSquares / n)

	if stdDevPop == 0 {
		return 0.0, "No variance in data (all values are the same), skewness is undefined."
	}

	skewness := (sumOfCubes / n) / math.Pow(stdDevPop, 3)

	var interpretation string
	absSkewness := math.Abs(skewness)

	if absSkewness < 0.2 {
		interpretation = "Data is approximately symmetrical."
	} else if skewness > 0 {
		if absSkewness < 0.5 {
			interpretation = "Data is moderately positively (right) skewed."
		} else {
			interpretation = "Data is highly positively (right) skewed."
		}
	} else {
		if absSkewness < 0.5 {
			interpretation = "Data is moderately negatively (left) skewed."
		} else {
			interpretation = "Data is highly negatively (left) skewed."
		}
	}

	return skewness, interpretation
}

func main() {
	data1 := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	skew1, interp1 := detectSkewness(data1)
	fmt.Printf("Dataset: %v\nSkewness: %.4f\nInterpretation: %s\n\n", data1, skew1, interp1)

	data2 := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 100}
	skew2, interp2 := detectSkewness(data2)
	fmt.Printf("Dataset: %v\nSkewness: %.4f\nInterpretation: %s\n\n", data2, skew2, interp2)

	data3 := []float64{1, 10, 20, 30, 40, 50, 60, 70, 80, 90}
	skew3, interp3 := detectSkewness(data3)
	fmt.Printf("Dataset: %v\nSkewness: %.4f\nInterpretation: %s\n\n", data3, skew3, interp3)

	data4 := []float64{5, 5, 5, 5, 5}
	skew4, interp4 := detectSkewness(data4)
	fmt.Printf("Dataset: %v\nSkewness: %.4f\nInterpretation: %s\n\n", data4, skew4, interp4)

	data5 := []float64{1, 2}
	skew5, interp5 := detectSkewness(data5)
	fmt.Printf("Dataset: %v\nSkewness: %.4f\nInterpretation: %s\n\n", data5, skew5, interp5)

	data6 := []float64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	skew6, interp6 := detectSkewness(data6)
	fmt.Printf("Dataset: %v\nSkewness: %.4f\nInterpretation: %s\n\n", data6, skew6, interp6)

	data7 := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100, 200, 300}
	skew7, interp7 := detectSkewness(data7)
	fmt.Printf("Dataset: %v\nSkewness: %.4f\nInterpretation: %s\n\n", data7, skew7, interp7)

	data8 := []float64{100, 90, 80, 70, 60, 50, 40, 30, 20, 10, 5, 2, 1}
	skew8, interp8 := detectSkewness(data8)
	fmt.Printf("Dataset: %v\nSkewness: %.4f\nInterpretation: %s\n\n", data8, skew8, interp8)
}

// Additional implementation at 2025-06-19 22:37:35
package main

import (
	"fmt"
	"math"
	"sort"
)

// calculateMean computes the arithmetic mean of a slice of float64.
func calculateMean(data []float64) float64 {
	if len(data) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

// calculateMedian computes the median of a slice of float64.
func calculateMedian(data []float64) float64 {
	if len(data) == 0 {
		return 0.0
	}
	sortedData := make([]float64, len(data))
	copy(sortedData, data)
	sort.Float64s(sortedData)

	mid := len(sortedData) / 2
	if len(sortedData)%2 == 0 {
		return (sortedData[mid-1] + sortedData[mid]) / 2.0
	}
	return sortedData[mid]
}

// calculateStandardDeviation computes the sample standard deviation of a slice of float64.
func calculateStandardDeviation(data []float64, mean float64) float64 {
	if len(data) < 2 { // Need at least 2 points for sample std dev
		return 0.0
	}
	sumSqDiff := 0.0
	for _, v := range data {
		diff := v - mean
		sumSqDiff += diff * diff
	}
	return math.Sqrt(sumSqDiff / float64(len(data)-1)) // Sample standard deviation
}

// calculateSkewnessCoefficient computes the Fisher-Pearson coefficient of skewness (third standardized moment).
// Skewness = E[((X - mu) / sigma)^3]
func calculateSkewnessCoefficient(data []float64, mean float64, stdDev float64) float64 {
	if len(data) < 3 || stdDev == 0 { // Need at least 3 points for meaningful skewness, and non-zero std dev
		return 0.0
	}
	sumCubedDiff := 0.0
	for _, v := range data {
		normalizedDiff := (v - mean) / stdDev
		sumCubedDiff += math.Pow(normalizedDiff, 3)
	}
	// This is the population-like third standardized moment.
	// For sample skewness, a more robust formula exists, but this is commonly used.
	return sumCubedDiff / float64(len(data))
}

// calculateFrequencyMap computes the frequency of each unique string in a slice.
func calculateFrequencyMap(data []string) map[string]int {
	freqMap := make(map[string]int)
	for _, v := range data {
		freqMap[v]++
	}
	return freqMap
}

// calculateGiniImpurity computes the Gini impurity for categorical data.
// Gini = 1 - sum(p_i^2)
func calculateGiniImpurity(freqMap map[string]int, total int) float64 {
	if total == 0 {
		return 0.0
	}
	sumSqProportions := 0.0
	for _, count := range freqMap {
		p := float64(count) / float64(total)
		sumSqProportions += p * p
	}
	return 1.0 - sumSqProportions
}

// calculateEntropy computes the entropy for categorical data.
// Entropy = -sum(p_i * log2(p_i))
func calculateEntropy(freqMap map[string]int, total int) float64 {
	if total == 0 {
		return 0.0
	}
	entropy := 0.0
	for _, count := range freqMap {
		p := float64(count) / float64(total)
		if p > 0 { // Avoid log(0)
			entropy -= p * math.Log2(p)
		}
	}
	return entropy
}

// AnalyzeNumericalData calculates and reports skewness for numerical data.
// It returns a boolean indicating if the data is considered skewed based on the threshold,
// along with mean, median, standard deviation, and skewness coefficient.
func AnalyzeNumericalData(data []float64, skewnessThreshold float64) (bool, float64, float64, float64, float64) {
	fmt.Println("\n--- Numerical Data Skewness Analysis ---")
	if len(data) == 0 {
		fmt.Println("No numerical data provided.")
		return false, 0, 0, 0, 0
	}

	mean := calculateMean(data)
	median := calculateMedian(data)
	stdDev := calculateStandardDeviation(data, mean)
	skewness := calculateSkewnessCoefficient(data, mean, stdDev)

	fmt.Printf("Data Points: %d\n", len(data))
	fmt.Printf("Mean: %.4f\n", mean)
	fmt.Printf("Median: %.4f\n", median)
	fmt.Printf("Standard Deviation: %.4f\n", stdDev)
	fmt.Printf("Skewness Coefficient (Fisher-Pearson): %.4f\n", skewness)

	isSkewed := false
	if math.Abs(skewness) > skewnessThreshold {
		isSkewed = true
		fmt.Printf("Conclusion: Data is likely SKEWED (absolute skewness %.4f > threshold %.4f).\n", math.Abs(skewness), skewnessThreshold)
		if skewness > 0 {
			fmt.Println("  (Positive skew: Tail on the right, Mean > Median)")
		} else {
			fmt.Println("  (Negative skew: Tail on the left, Mean < Median)")
		}
	} else {
		fmt.Printf("Conclusion: Data is NOT significantly skewed (absolute skewness %.4f <= threshold %.4f).\n", math.Abs(skewness), skewnessThreshold)
	}

	return isSkewed, mean, median, stdDev, skewness
}

// AnalyzeCategoricalData calculates and reports imbalance for categorical data.
// It returns a boolean indicating if the data is considered imbalanced based on the thresholds,
// along with the frequency map, Gini impurity, and entropy.
func AnalyzeCategoricalData(data []string, giniThreshold float64, entropyThreshold float64) (bool, map[string]int, float64, float64) {
	fmt.Println("\n--- Categorical Data Imbalance Analysis ---")
	if len(data) == 0 {
		fmt.Println("No categorical data provided.")
		return false, nil, 0, 0
	}

	total := len(data)
	freqMap := calculateFrequencyMap(data)
	gini := calculateGiniImpurity(freqMap, total)
	entropy := calculateEntropy(freqMap, total)

	fmt.Printf("Total Data Points: %d\n", total)
	fmt.Println("Frequency Distribution:")
	for category, count := range freqMap {
		fmt.Printf("  - %s: %d (%.2f%%)\n", category, count, float64(count)/float64(total)*100)
	}
	fmt.Printf("Gini Impurity: %.4f\n", gini)
	fmt.Printf("Entropy: %.4f\n", entropy)

	isImbalanced := false
	// Lower Gini/Entropy means more pure/less diverse (more imbalanced towards one category)
	// Higher Gini/Entropy means more mixed/diverse (less imbalanced)
	// So, we check if Gini/Entropy is *below* a certain threshold for imbalance.
	if gini < giniThreshold || entropy < entropyThreshold {
		isImbalanced = true
		fmt.Printf("Conclusion: Data is likely IMBALANCED (Gini %.4f < threshold %.4f OR Entropy %.4f < threshold %.4f).\n", gini, giniThreshold, entropy, entropyThreshold)
	} else {
		fmt.Printf("Conclusion: Data is NOT significantly imbalanced (Gini %.4f >= threshold %.4f AND Entropy %.4f >= threshold %.4f).\n", gini, giniThreshold, entropy, entropyThreshold)
	}

	return isImbalanced, freqMap, gini, entropy
}

func main() {
	// Example 1: Numerical Data (Right-Skewed)
	numericalData1 := []float64{1, 2, 2, 3, 3, 3, 4, 4, 5, 10, 15, 20}
	skewed1, _, _, _, _ := AnalyzeNumericalData(numericalData1, 0.5) // Skewness threshold: |skewness| > 0.5

	// Example 2: Numerical Data (Symmetric/Normal-like)
	numericalData2 := []float64{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	skewed2, _, _, _, _ := AnalyzeNumericalData(numericalData2, 0.5)

	// Example 3: Numerical Data (Left-Skewed)
	numericalData3 := []float64{1, 5, 10, 15, 15, 16, 17, 17, 18, 18, 19, 20}
	skewed3, _, _, _, _ := AnalyzeNumericalData(numericalData3, 0.5)

	// Example 4: Categorical Data (Imbalanced)
	categoricalData1 := []string{"Male", "Male", "Male", "Female", "Male", "Male", "Male", "Female"}
	imbalanced1, _, _, _ := AnalyzeCategoricalData(categoricalData1, 0.4, 0.7) // Gini/Entropy thresholds for imbalance

	// Example 5: Categorical Data (Balanced)
	categoricalData2 := []string{"Red", "Blue", "Green", "Red", "Blue", "Green", "Red", "Blue", "Green"}
	imbalanced2, _, _, _ := AnalyzeCategoricalData(categoricalData2, 0.4, 0.7)

	// Example 6: Empty Data
	AnalyzeNumericalData([]float64{}, 0.5)
	AnalyzeCategoricalData([]string{}, 0.4, 0.7)

	fmt.Println("\n--- Overall Skewness/Imbalance Summary ---")
	fmt.Printf("Numerical Data 1 (Right-Skewed Example): Skewed = %t\n", skewed1)
	fmt.Printf("Numerical Data 2 (Symmetric Example): Skewed = %t\n", skewed2)
	fmt.Printf("Numerical Data 3 (Left-Skewed Example): Skewed = %t\n", skewed3)
	fmt.Printf("Categorical Data 1 (Imbalanced Example): Imbalanced = %t\n", imbalanced1)
	fmt.Printf("Categorical Data 2 (Balanced Example): Imbalanced = %t\n", imbalanced2)
}

// Additional implementation at 2025-06-19 22:38:40
package main

import (
	"fmt"
	"math"
	"sort"
)

// SkewnessResult holds the calculated statistical measures and skewness classification.
type SkewnessResult struct {
	Mean             float64
	Median           float64
	StandardDeviation float64
	SkewnessCoefficient float64
	Classification   string
	Message          string
}

// calculateMean computes the arithmetic mean of a slice of float64.
func calculateMean(data []float64) float64 {
	if len(data) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, val := range data {
		sum += val
	}
	return sum / float64(len(data))
}

// calculateMedian computes the median of a slice of float64.
// It sorts the data, so a copy should be passed if the original order needs to be preserved.
func calculateMedian(data []float64) float64 {
	if len(data) == 0 {
		return 0.0
	}
	sortedData := make([]float64, len(data))
	copy(sortedData, data)
	sort.Float64s(sortedData)

	n := len(sortedData)
	if n%2 == 1 {
		return sortedData[n/2]
	}
	return (sortedData[n/2-1] + sortedData[n/2]) / 2.0
}

// calculateStandardDeviation computes the population standard deviation of a slice of float64.
func calculateStandardDeviation(data []float64, mean float64) float64 {
	if len(data) == 0 {
		return 0.0
	}
	sumOfSquares := 0.0
	for _, val := range data {
		diff := val - mean
		sumOfSquares += diff * diff
	}
	return math.Sqrt(sumOfSquares / float64(len(data)))
}

// calculateSkewnessCoefficient computes the Fisher-Pearson coefficient of skewness.
// It returns the coefficient and a boolean indicating if the calculation was successful (false if std dev is zero).
func calculateSkewnessCoefficient(data []float64, mean float64, stdDev float64) (float64, bool) {
	if len(data) < 3 { // Skewness is typically meaningful for N >= 3
		return 0.0, false
	}
	if stdDev == 0.0 { // All values are the same, skewness is undefined or 0 (perfectly symmetric)
		return 0.0, false
	}

	sumOfCubes := 0.0
	for _, val := range data {
		diff := val - mean
		sumOfCubes += diff * diff * diff
	}
	thirdMoment := sumOfCubes / float64(len(data))
	return thirdMoment / math.Pow(stdDev, 3), true
}

// classifySkewness categorizes the skewness coefficient into descriptive terms.
func classifySkewness(skewness float64) string {
	absSkewness := math.Abs(skewness)
	if absSkewness < 0.1 {
		return "Symmetric"
	} else if absSkewness < 0.5 {
		if skewness > 0 {
			return "Slightly Right (Positive) Skewed"
		}
		return "Slightly Left (Negative) Skewed"
	} else if absSkewness < 1.0 {
		if skewness > 0 {
			return "Moderately Right (Positive) Skewed"
		}
		return "Moderately Left (Negative) Skewed"
	} else {
		if skewness > 0 {
			return "Highly Right (Positive) Skewed"
		}
		return "Highly Left (Negative) Skewed"
	}
}

// DetectSkewness analyzes a slice of float64 data for skewness and returns a SkewnessResult.
func DetectSkewness(data []float64) SkewnessResult {
	result := SkewnessResult{}

	if len(data) == 0 {
		result.Message = "Error: No data provided."
		return result
	}
	if len(data) == 1 {
		result.Mean = data[0]
		result.Median = data[0]
		result.StandardDeviation = 0.0
		result.SkewnessCoefficient = 0.0
		result.Classification = "Symmetric"
		result.Message = "Warning: Skewness is not well-defined for single data point."
		return result
	}

	result.Mean = calculateMean(data)
	result.Median = calculateMedian(data)
	result.StandardDeviation = calculateStandardDeviation(data, result.Mean)

	skewCoeff, ok := calculateSkewnessCoefficient(data, result.Mean, result.StandardDeviation)
	if !ok {
		result.SkewnessCoefficient = 0.0 // Default to 0 if std dev is 0 or not enough data
		result.Classification = "Symmetric"
		if result.StandardDeviation == 0.0 {
			result.Message = "All data points are identical. Distribution is perfectly symmetric."
		} else {
			result.Message = "Not enough data points to reliably calculate skewness coefficient (N < 3)."
		}
	} else {
		result.SkewnessCoefficient = skewCoeff
		result.Classification = classifySkewness(skewCoeff)
		result.Message = "Skewness analysis complete."
	}

	return result
}

func main() {
	// Example 1: Symmetric data (normal-like)
	data1 := []float64{10, 20, 30, 40, 50, 60, 70, 80, 90}
	fmt.Println("--- Data Set 1 (Symmetric) ---")
	res1 := DetectSkewness(data1)
	fmt.Printf("Data Points: %v\n", data1)
	fmt.Printf("Mean: %.2f\n", res1.Mean)
	fmt.Printf("Median: %.2f\n", res1.Median)
	fmt.Printf("Standard Deviation: %.2f\n", res1.StandardDeviation)
	fmt.Printf("Skewness Coefficient: %.4f\n", res1.SkewnessCoefficient)
	fmt.Printf("Classification: %s\n", res1.Classification)
	fmt.Printf("Message: %s\n\n", res1.Message)

	// Example 2: Right-skewed data (positive skew)
	data2 := []float64{10, 20, 30, 40, 50, 60, 70, 80, 100, 200, 300}
	fmt.Println("--- Data Set 2 (Right-Skewed) ---")
	res2 := DetectSkewness(data2)
	fmt.Printf("Data Points: %v\n", data2)
	fmt.Printf("Mean: %.2f\n", res2.Mean)
	fmt.Printf("Median: %.2f\n", res2.Median)
	fmt.Printf("Standard Deviation: %.2f\n", res2.StandardDeviation)
	fmt.Printf("Skewness Coefficient: %.4f\n", res2.SkewnessCoefficient)
	fmt.Printf("Classification: %s\n", res2.Classification)
	fmt.Printf("Message: %s\n\n", res2.Message)

	// Example 3: Left-skewed data (negative skew)
	data3 := []float64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 10, 5, 2}
	fmt.Println("--- Data Set 3 (Left-Skewed) ---")
	res3 := DetectSkewness(data3)
	fmt.Printf("Data Points: %v\n", data3)
	fmt.Printf("Mean: %.2f\n", res3.Mean)
	fmt.Printf("Median: %.2f\n", res3.Median)
	fmt.Printf("Standard Deviation: %.2f\n", res3.StandardDeviation)
	fmt.Printf("Skewness Coefficient: %.4f\n", res3.SkewnessCoefficient)
	fmt.Printf("Classification: %s\n", res3.Classification)
	fmt.Printf("Message: %s\n\n", res3.Message)

	// Example 4: Data with all identical values
	data4 := []float64{5, 5, 5, 5, 5}
	fmt.Println("--- Data Set 4 (Identical Values) ---")
	res4 := DetectSkewness(data4)
	fmt.Printf("Data Points: %v\n", data4)
	fmt.Printf("Mean: %.2f\n", res4.Mean)
	fmt.Printf("Median: %.2f\n", res4.Median)
	fmt.Printf("Standard Deviation: %.2f\n", res4.StandardDeviation)
	fmt.Printf("Skewness Coefficient: %.4f\n", res4.SkewnessCoefficient)
	fmt.Printf("Classification: %s\n", res4.Classification)
	fmt.Printf("Message: %s\n\n", res4.Message)

	// Example 5: Empty data set
	data5 := []float64{}
	fmt.Println("--- Data Set 5 (Empty) ---")
	res5 := DetectSkewness(data5)
	fmt.Printf("Data Points: %v\n", data5)
	fmt.Printf("Mean: %.2f\n", res5.Mean)
	fmt.Printf("Median: %.2f\n", res5.Median)
	fmt.Printf("Standard Deviation: %.2f\n", res5.StandardDeviation)
	fmt.Printf("Skewness Coefficient: %.4f\n", res5.SkewnessCoefficient)
	fmt.Printf("Classification: %s\n", res5.Classification)
	fmt.Printf("Message: %s\n\n", res5.Message)

	// Example 6: Single data point
	data6 := []float64{42.0}
	fmt.Println("--- Data Set 6 (Single Point) ---")
	res6 := DetectSkewness(data6)
	fmt.Printf("Data Points: %v\n", data6)
	fmt.Printf("Mean: %.2f\n", res6.Mean)
	fmt.Printf("Median: %.2f\n", res6.Median)
	fmt.Printf("Standard Deviation: %.2f\n", res6.StandardDeviation)
	fmt.Printf("Skewness Coefficient: %.4f\n", res6.SkewnessCoefficient)
	fmt.Printf("Classification: %s\n", res6.Classification)
	fmt.Printf("Message: %s\n\n", res6.Message)

	// Example 7: Two data points (not enough for skewness coefficient)
	data7 := []float64{10, 20}
	fmt.Println("--- Data Set 7 (Two Points) ---")
	res7 := DetectSkewness(data7)
	fmt.Printf("Data Points: %v\n", data7)
	fmt.Printf("Mean: %.2f\n", res7.Mean)
	fmt.Printf("Median: %.2f\n", res7.Median)
	fmt.Printf("Standard Deviation: %.2f\n", res7.StandardDeviation)
	fmt.Printf("Skewness Coefficient: %.4f\n", res7.SkewnessCoefficient)
	fmt.Printf("Classification: %s\n", res7.Classification)
	fmt.Printf("Message: %s\n\n", res7.Message)
}