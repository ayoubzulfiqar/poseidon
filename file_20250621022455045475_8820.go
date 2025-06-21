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
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func calculateStandardDeviation(data []float64, mean float64) float64 {
	if len(data) == 0 {
		return 0.0
	}
	sumSqDiff := 0.0
	for _, v := range data {
		diff := v - mean
		sumSqDiff += diff * diff
	}
	// Using population standard deviation
	variance := sumSqDiff / float64(len(data))
	return math.Sqrt(variance)
}

// findOutliers identifies statistical outliers in a dataset using the Z-score method.
// A data point is considered an outlier if its absolute difference from the mean
// is greater than a specified threshold multiplied by the standard deviation.
// It returns two slices: the first contains the outliers, the second contains the non-outliers.
func findOutliers(data []float64, threshold float64) ([]float64, []float64) {
	if len(data) == 0 {
		return nil, nil
	}

	mean := calculateMean(data)
	stdDev := calculateStandardDeviation(data, mean)

	var outliers []float64
	var nonOutliers []float64

	// If standard deviation is 0 (e.g., all data points are the same),
	// no point can be an outlier by this method unless threshold is negative.
	// We treat all as non-outliers in this case.
	if stdDev == 0 {
		return nil, data
	}

	for _, val := range data {
		if math.Abs(val-mean) > threshold*stdDev {
			outliers = append(outliers, val)
		} else {
			nonOutliers = append(nonOutliers, val)
		}
	}
	return outliers, nonOutliers
}

func main() {
	// Sample datasets to test the outlier detection
	dataset1 := []float64{10, 12, 11, 13, 100, 14, 15, 12, 11, 13}
	dataset2 := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	dataset3 := []float64{100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 1000}
	dataset4 := []float64{} // Empty dataset
	dataset5 := []float64{50} // Single element dataset
	dataset6 := []float64{5, 5, 5, 5, 5} // All elements are the same

	// Define the outlier threshold (e.g., 2.0 for 2 standard deviations)
	outlierThreshold := 2.0

	fmt.Println("--- Analyzing Dataset 1 ---")
	fmt.Printf("Data: %v\n", dataset1)
	mean1 := calculateMean(dataset1)
	stdDev1 := calculateStandardDeviation(dataset1, mean1)
	fmt.Printf("Mean: %.2f\n", mean1)
	fmt.Printf("Standard Deviation: %.2f\n", stdDev1)
	outliers1, nonOutliers1 := findOutliers(dataset1, outlierThreshold)
	fmt.Printf("Outliers (threshold %.1f std dev): %v\n", outlierThreshold, outliers1)
	fmt.Printf("Non-Outliers: %v\n", nonOutliers1)
	fmt.Println()

	fmt.Println("--- Analyzing Dataset 2 ---")
	fmt.Printf("Data: %v\n", dataset2)
	mean2 := calculateMean(dataset2)
	stdDev2 := calculateStandardDeviation(dataset2, mean2)
	fmt.Printf("Mean: %.2f\n", mean2)
	fmt.Printf("Standard Deviation: %.2f\n", stdDev2)
	outliers2, nonOutliers2 := findOutliers(dataset2, outlierThreshold)
	fmt.Printf("Outliers (threshold %.1f std dev): %v\n", outlierThreshold, outliers2)
	fmt.Printf("Non-Outliers: %v\n", nonOutliers2)
	fmt.Println()

	fmt.Println("--- Analyzing Dataset 3 ---")
	fmt.Printf("Data: %v\n", dataset3)
	mean3 := calculateMean(dataset3)
	stdDev3 := calculateStandardDeviation(dataset3, mean3)
	fmt.Printf("Mean: %.2f\n", mean3)
	fmt.Printf("Standard Deviation: %.2f\n", stdDev3)
	outliers3, nonOutliers3 := findOutliers(dataset3, outlierThreshold)
	fmt.Printf("Outliers (threshold %.1f std dev): %v\n", outlierThreshold, outliers3)
	fmt.Printf("Non-Outliers: %v\n", nonOutliers3)
	fmt.Println()

	fmt.Println("--- Analyzing Dataset 4 (Empty) ---")
	fmt.Printf("Data: %v\n", dataset4)
	outliers4, nonOutliers4 := findOutliers(dataset4, outlierThreshold)
	fmt.Printf("Outliers (threshold %.1f std dev): %v\n", outlierThreshold, outliers4)
	fmt.Printf("Non-Outliers: %v\n", nonOutliers4)
	fmt.Println()

	fmt.Println("--- Analyzing Dataset 5 (Single Element) ---")
	fmt.Printf("Data: %v\n", dataset5)
	mean5 := calculateMean(dataset5)
	stdDev5 := calculateStandardDeviation(dataset5, mean5)
	fmt.Printf("Mean: %.2f\n", mean5)
	fmt.Printf("Standard Deviation: %.2f\n", stdDev5)
	outliers5, nonOutliers5 := findOutliers(dataset5, outlierThreshold)
	fmt.Printf("Outliers (threshold %.1f std dev): %v\n", outlierThreshold, outliers5)
	fmt.Printf("Non-Outliers: %v\n", nonOutliers5)
	fmt.Println()

	fmt.Println("--- Analyzing Dataset 6 (All Same Elements) ---")
	fmt.Printf("Data: %v\n", dataset6)
	mean6 := calculateMean(dataset6)
	stdDev6 := calculateStandardDeviation(dataset6, mean6)
	fmt.Printf("Mean: %.2f\n", mean6)
	fmt.Printf("Standard Deviation: %.2f\n", stdDev6)
	outliers6, nonOutliers6 := findOutliers(dataset6, outlierThreshold)
	fmt.Printf("Outliers (threshold %.1f std dev): %v\n", outlierThreshold, outliers6)
	fmt.Printf("Non-Outliers: %v\n", nonOutliers6)
	fmt.Println()
}

// Additional implementation at 2025-06-21 02:26:04
package main

import (
	"fmt"
	"math"
	"sort"
)

// calculateMedian calculates the median of a sorted slice of floats.
func calculateMedian(data []float64) float64 {
	n := len(data)
	if n == 0 {
		return 0.0
	}
	if n%2 == 1 {
		return data[n/2]
	}
	return (data[n/2-1] + data[n/2]) / 2.0
}

// calculateQuartile calculates the quartile (Q1 or Q3) of a sorted slice of floats.
// `q` should be 0.25 for Q1 or 0.75 for Q3.
func calculateQuartile(data []float64, q float64) float64 {
	n := len(data)
	if n == 0 {
		return 0.0
	}
	// Using the "inclusive" method (Type 6 or R-7)
	index := q * float64(n-1)
	lowerIndex := int(math.Floor(index))
	upperIndex := int(math.Ceil(index))

	if lowerIndex == upperIndex {
		return data[lowerIndex]
	}

	// Linear interpolation
	return data[lowerIndex]*(1-(index-float64(lowerIndex))) + data[upperIndex]*(index-float64(lowerIndex))
}

// findOutliersIQR finds outliers using the IQR method (Tukey's fences).
// `k` is the multiplier for the IQR, typically 1.5 for outliers.
func findOutliersIQR(data []float64, k float64) ([]float64, []float64) {
	if len(data) == 0 {
		return nil, nil
	}

	sortedData := make([]float64, len(data))
	copy(sortedData, data)
	sort.Float64s(sortedData)

	q1 := calculateQuartile(sortedData, 0.25)
	q3 := calculateQuartile(sortedData, 0.75)
	iqr := q3 - q1

	lowerBound := q1 - k*iqr
	upperBound := q3 + k*iqr

	var lowerOutliers []float64
	var upperOutliers []float64

	for _, val := range data {
		if val < lowerBound {
			lowerOutliers = append(lowerOutliers, val)
		} else if val > upperBound {
			upperOutliers = append(upperOutliers, val)
		}
	}
	return lowerOutliers, upperOutliers
}

// calculateMean calculates the arithmetic mean of a slice of floats.
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

// calculateStandardDeviation calculates the sample standard deviation of a slice of floats.
func calculateStandardDeviation(data []float64) float64 {
	n := len(data)
	if n < 2 { // Sample standard deviation requires at least 2 data points
		return 0.0
	}
	mean := calculateMean(data)
	sumSqDiff := 0.0
	for _, val := range data {
		diff := val - mean
		sumSqDiff += diff * diff
	}
	// Use n-1 for sample standard deviation
	return math.Sqrt(sumSqDiff / float64(n-1))
}

// findOutliersZScore finds outliers using the Z-score method.
// `threshold` is the absolute Z-score value above which a point is considered an outlier.
func findOutliersZScore(data []float64, threshold float64) []float64 {
	n := len(data)
	if n < 2 { // Need at least 2 points to calculate meaningful std dev
		return nil
	}

	mean := calculateMean(data)
	stdDev := calculateStandardDeviation(data)

	if stdDev == 0 { // All data points are identical, no outliers by Z-score
		return nil
	}

	var outliers []float64
	for _, val := range data {
		zScore := (val - mean) / stdDev
		if math.Abs(zScore) > threshold {
			outliers = append(outliers, val)
		}
	}
	return outliers
}

func main() {
	// Sample dataset
	data := []float64{10, 12, 12, 13, 12, 11, 14, 13, 15, 10, 10, 100, 10, 12, 11, 13, 11, 10, 10, 10, -50}

	fmt.Println("Original Data:", data)

	// --- Outlier Detection using IQR Method (Tukey's Fences) ---
	// k is typically 1.5 for "outliers" or 3.0 for "extreme outliers"
	kFactor := 1.5
	lowerOutliersIQR, upperOutliersIQR := findOutliersIQR(data, kFactor)

	fmt.Println("\n--- IQR Method (Tukey's Fences, k =", kFactor, ") ---")
	if len(lowerOutliersIQR) > 0 {
		fmt.Println("Lower Outliers:", lowerOutliersIQR)
	} else {
		fmt.Println("No Lower Outliers found.")
	}
	if len(upperOutliersIQR) > 0 {
		fmt.Println("Upper Outliers:", upperOutliersIQR)
	} else {
		fmt.Println("No Upper Outliers found.")
	}
	if len(lowerOutliersIQR) == 0 && len(upperOutliersIQR) == 0 {
		fmt.Println("No outliers found using IQR method.")
	}

	// --- Outlier Detection using Z-score Method ---
	// Common thresholds: 2.0, 2.5, 3.0
	zScoreThreshold := 2.5
	outliersZScore := findOutliersZScore(data, zScoreThreshold)

	fmt.Println("\n--- Z-score Method (Threshold =", zScoreThreshold, ") ---")
	if len(outliersZScore) > 0 {
		fmt.Println("Outliers:", outliersZScore)
	} else {
		fmt.Println("No outliers found using Z-score method.")
	}
}

// Additional implementation at 2025-06-21 02:26:51
package main

import (
	"fmt"
	"math"
	"sort"
)

// StatisticsReport holds all calculated statistical measures and identified outliers.
type StatisticsReport struct {
	Data              []float64 // The original data used for calculation
	Mean              float64
	Median            float64
	Q1                float64 // First Quartile (25th percentile)
	Q3                float64 // Third Quartile (75th percentile)
	IQR               float64 // Interquartile Range (Q3 - Q1)
	StandardDeviation float64
	LowerOutliers     []float64 // Values below Q1 - 1.5 * IQR
	UpperOutliers     []float64 // Values above Q3 + 1.5 * IQR
	ZScores           []float64 // Z-score for each corresponding data point in 'Data'
}

// CalculateStatistics analyzes a slice of float64 data and returns a StatisticsReport.
// It uses the IQR method for outlier detection and also calculates mean, median,
// standard deviation, and Z-scores.
func CalculateStatistics(data []float64) (*StatisticsReport, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("data slice must contain at least 2 elements for meaningful statistics")
	}

	report := &StatisticsReport{
		Data: make([]float64, len(data)),
	}
	copy(report.Data, data) // Store a copy of original data

	// Sort a copy of the data for median, quartiles, and IQR
	sortedData := make([]float64, len(data))
	copy(sortedData, data)
	sort.Float64s(sortedData)

	// Calculate Mean
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	report.Mean = sum / float64(len(data))

	// Calculate Median
	mid := len(sortedData) / 2
	if len(sortedData)%2 == 0 {
		report.Median = (sortedData[mid-1] + sortedData[mid]) / 2.0
	} else {
		report.Median = sortedData[mid]
	}

	// Calculate Q1 and Q3
	report.Q1 = percentile(sortedData, 0.25)
	report.Q3 = percentile(sortedData, 0.75)
	report.IQR = report.Q3 - report.Q1

	// Calculate Standard Deviation (Population Standard Deviation)
	varianceSum := 0.0
	for _, v := range data {
		varianceSum += math.Pow(v-report.Mean, 2)
	}
	report.StandardDeviation = math.Sqrt(varianceSum / float64(len(data)))

	// Identify Outliers using IQR method
	lowerBound := report.Q1 - 1.5*report.IQR
	upperBound := report.Q3 + 1.5*report.IQR

	for _, v := range data {
		if v < lowerBound {
			report.LowerOutliers = append(report.LowerOutliers, v)
		} else if v > upperBound {
			report.UpperOutliers = append(report.UpperOutliers, v)
		}
	}

	// Calculate Z-Scores for each data point
	report.ZScores = make([]float64, len(data))
	if report.StandardDeviation > 0 { // Avoid division by zero
		for i, v := range data {
			report.ZScores[i] = (v - report.Mean) / report.StandardDeviation
		}
	} else {
		// If std dev is 0, all values are the same, Z-score is 0 for all
		for i := range data {
			report.ZScores[i] = 0.0
		}
	}

	return report, nil
}

// percentile calculates the value at a given percentile (0.0 to 1.0) for sorted data.
// Uses linear interpolation between closest ranks.
func percentile(sortedData []float64, p float64) float64 {
	if len(sortedData) == 0 {
		return 0.0
	}
	if p < 0.0 || p > 1.0 {
		// This case should ideally be handled by the caller or return an error,
		// but for internal use where p is fixed (0.25, 0.75), panic is acceptable.
		panic("percentile must be between 0.0 and 1.0")
	}

	// N = number of data points
	// P = percentile (e.g., 0.25 for 25th percentile)
	// L = (N-1) * P
	// If L is an integer, the percentile is the L-th value (0-indexed).
	// If L is not an integer, interpolate between floor(L) and ceil(L) values.

	index := p * float64(len(sortedData)-1)
	lowerIndex := int(math.Floor(index))
	upperIndex := int(math.Ceil(index))

	if lowerIndex == upperIndex {
		return sortedData[lowerIndex]
	}

	// Linear interpolation
	lowerValue := sortedData[lowerIndex]
	upperValue := sortedData[upperIndex]
	interpolationFactor := index - float64(lowerIndex)

	return lowerValue + (upperValue-lowerValue)*interpolationFactor
}

func main() {
	// Example data set with clear outliers
	data := []float64{
		10.5, 12.1, 12.0, 13.5, 12.2, 11.8, 14.0, 13.1, 15.0, 10.9,
		100.0, // Obvious upper outlier
		12.3, 10.1, 10.2, 12.5, 11.9, 10.3, 12.4, 11.7, 12.6, 11.5,
		-50.0, // Obvious lower outlier
		12.7, 12.8, 12.9, 12.0, 14.1, 11.6, 12.0, 12.1, 15.1, 12.2,
	}

	report, err := CalculateStatistics(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("--- Statistical Outlier Report ---")
	fmt.Printf("Original Data Points: %v\n", report.Data)
	fmt.Printf("Number of Data Points: %d\n", len(report.Data))
	fmt.Printf("Mean: %.2f\n", report.Mean)
	fmt.Printf("Median: %.2f\n", report.Median)
	fmt.Printf("Q1 (25th Percentile): %.2f\n", report.Q1)
	fmt.Printf("Q3 (75th Percentile): %.2f\n", report.Q3)
	fmt.Printf("IQR (Interquartile Range): %.2f\n", report.IQR)
	fmt.Printf("Standard Deviation (Population): %.2f\n", report.StandardDeviation)

	fmt.Println("\n--- Outliers (IQR Method) ---")
	if len(report.LowerOutliers) > 0 {
		fmt.Printf("Lower Outliers (< %.2f): %v\n", report.Q1-1.5*report.IQR, report.LowerOutliers)
	} else {
		fmt.Println("No Lower Outliers found.")
	}
	if len(report.UpperOutliers) > 0 {
		fmt.Printf("Upper Outliers (> %.2f): %v\n", report.Q3+1.5*report.IQR, report.UpperOutliers)
	} else {
		fmt.Println("No Upper Outliers found.")
	}
	if len(report.LowerOutliers) == 0 && len(report.UpperOutliers) == 0 {
		fmt.Println("No outliers found using the IQR method.")
	}

	fmt.Println("\n--- Z-Scores for Data Points ---")
	// Print Z-scores alongside their original values
	for i, val := range report.Data {
		fmt.Printf("Data[%d]: %.2f, Z-Score: %.2f\n", i, val, report.ZScores[i])
	}

	// Example with no outliers
	fmt.Println("\n--- Example with no outliers ---")
	dataNoOutliers := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	reportNoOutliers, err := CalculateStatistics(dataNoOutliers)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Original Data Points: %v\n", reportNoOutliers.Data)
	fmt.Printf("Mean: %.2f, Median: %.2f, IQR: %.2f\n", reportNoOutliers.Mean, reportNoOutliers.Median, reportNoOutliers.IQR)
	if len(reportNoOutliers.Lower

// Additional implementation at 2025-06-21 02:27:57
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// DataPoint represents a single data point with its original value and index.
type DataPoint struct {
	Value float64
	Index int
}

// Statistics holds various calculated statistics for a dataset.
type Statistics struct {
	Mean   float64
	Median float64
	Q1     float64 // First Quartile
	Q3     float64 // Third Quartile
	IQR    float64 // Interquartile Range
	StdDev float64 // Standard Deviation (sample)
	Min    float64
	Max    float64
}

// calculateMedian calculates the median of a sorted slice of floats.
func calculateMedian(data []float64) float64 {
	n := len(data)
	if n == 0 {
		return 0.0
	}
	if n%2 == 1 {
		return data[n/2]
	}
	return (data[n/2-1] + data[n/2]) / 2.0
}

// calculateStatistics computes various statistics for a given dataset.
// It requires the data to be sorted for Q1, Q3, and Median.
func calculateStatistics(data []float64) Statistics {
	n := len(data)
	if n == 0 {
		return Statistics{}
	}

	// Create a copy and sort for quartile calculations
	sortedData := make([]float64, n)
	copy(sortedData, data)
	sort.Float64s(sortedData)

	var stats Statistics
	stats.Min = sortedData[0]
	stats.Max = sortedData[n-1]

	// Mean
	sum := 0.0
	for _, val := range sortedData {
		sum += val
	}
	stats.Mean = sum / float64(n)

	// Median (Q2)
	stats.Median = calculateMedian(sortedData)

	// Q1 and Q3
	if n%2 == 1 {
		// Odd number of elements, median is excluded from halves
		stats.Q1 = calculateMedian(sortedData[:n/2])
		stats.Q3 = calculateMedian(sortedData[n/2+1:])
	} else {
		// Even number of elements, median is average of two middle, so halves include them
		stats.Q1 = calculateMedian(sortedData[:n/2])
		stats.Q3 = calculateMedian(sortedData[n/2:])
	}

	stats.IQR = stats.Q3 - stats.Q1

	// Standard Deviation (sample standard deviation)
	sumSqDiff := 0.0
	for _, val := range sortedData {
		diff := val - stats.Mean
		sumSqDiff += diff * diff
	}
	if n > 1 {
		stats.StdDev = math.Sqrt(sumSqDiff / float64(n-1))
	} else {
		stats.StdDev = 0.0
	}

	return stats
}

// findOutliersIQR finds outliers using the Interquartile Range (IQR) method.
// k is the multiplier for the IQR (commonly 1.5).
func findOutliersIQR(data []DataPoint, k float64) ([]DataPoint, []DataPoint, Statistics) {
	if len(data) == 0 {
		return nil, nil, Statistics{}
	}

	// Extract just the values for statistical calculations
	values := make([]float64, len(data))
	for i, dp := range data {
		values[i] = dp.Value
	}

	stats := calculateStatistics(values)

	lowerBound := stats.Q1 - k*stats.IQR
	upperBound := stats.Q3 + k*stats.IQR

	var outliers []DataPoint
	var nonOutliers []DataPoint

	for _, dp := range data {
		if dp.Value < lowerBound || dp.Value > upperBound {
			outliers = append(outliers, dp)
		} else {
			nonOutliers = append(nonOutliers, dp)
		}
	}

	return outliers, nonOutliers, stats
}

// readDataFromFile reads numerical data from a CSV file.
// It expects a single column of numbers.
func readDataFromFile(filePath string) ([]DataPoint, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// If your CSV has a header, you might want to skip the first record:
	// _, err = reader.Read() // Skip header
	// if err != nil && err != io.EOF {
	//     return nil, fmt.Errorf("failed to read header: %w", err)
	// }

	var data []DataPoint
	recordCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %w", err)
		}

		if len(record) == 0 || (len(record) == 1 && strings.TrimSpace(record[0]) == "") {
			continue // Skip empty lines
		}

		// Assuming data is in the first column
		valStr := strings.TrimSpace(record[0])
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			// Log error but continue if possible, or return error
			fmt.Fprintf(os.Stderr, "Warning: Could not parse '%s' as float on line %d. Skipping.\n", valStr, recordCount+1)
			continue
		}
		data = append(data, DataPoint{Value: val, Index: recordCount})
		recordCount++
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no valid numerical data found in file %s", filePath)
	}

	return data, nil
}

// printResults displays the statistical summary and identified outliers.
func printResults(originalData []DataPoint, outliers []DataPoint, nonOutliers []DataPoint, stats Statistics, k float64) {
	fmt.Println("--- Dataset Summary ---")
	fmt.Printf("Total Data Points: %d\n", len(originalData))
	fmt.Printf("Mean: %.4f\n", stats.Mean)
	fmt.Printf("Median (Q2): %.4f\n", stats.Median)
	fmt.Printf("Q1 (25th Percentile): %.4f\n", stats.Q1)
	fmt.Printf("Q3 (75th Percentile): %.4f\n", stats.Q3)
	fmt.Printf("IQR (Q3 - Q1): %.4f\n", stats.IQR)
	fmt.Printf("Standard Deviation: %.4f\n", stats.StdDev)
	fmt.Printf("Min Value: %.4f\n", stats.Min)
	fmt.Printf("Max Value: %.4f\n", stats.Max)
	fmt.Println("-----------------------")

	fmt.Printf("\n--- Outlier Detection (IQR Method with k=%.1f) ---\n", k)
	lowerBound := stats.Q1 - k*stats.IQR
	upperBound := stats.Q3 + k*stats.IQR
	fmt.Printf("Lower Bound (Q1 - %.1f*IQR): %.4f\n", k, lowerBound)
	fmt.Printf("Upper Bound (Q3 + %.1f*IQR): %.4f\n", k, upperBound)

	fmt.Printf("\nNumber of Outliers Found: %d\n", len(outliers))
	if len(outliers) > 0 {
		fmt.Println("Outliers:")
		for _, dp := range outliers {
			fmt.Printf("  Value: %.4f (Original Index: %d)\n", dp.Value, dp.Index)
		}
	} else {
		fmt.Println("No outliers detected.")
	}

	fmt.Printf("\nNumber of Non-Outliers: %d\n", len(nonOutliers))
	fmt.Println("-------------------------------------------------")
}

func main() {
	// Example 1: Hardcoded data
	fmt.Println("--- Running with Hardcoded Data ---")
	hardcodedValues := []float64{
		10, 12, 12, 13, 12, 11, 14, 13, 15, 10,
		10, 12, 12, 13, 12, 11, 14, 13, 15, 10,
		10, 12, 12, 13, 12, 11, 14, 13, 15, 10,
		100, // Obvious outlier
		-50, // Obvious outlier
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	}

	var hardcodedDataPoints []DataPoint
	for i, val := range hardcodedValues {
		hardcodedDataPoints = append(hardcodedDataPoints, DataPoint{Value: val, Index: i})
	}

	// Default IQR multiplier
	iqrK := 1.5
	outliers, nonOutliers, stats := findOutliersIQR(hardcodedDataPoints, iqrK)
	printResults(hardcodedDataPoints, outliers, nonOutliers, stats, iqrK)

	fmt.Println("\n--- Running with Custom IQR Multiplier (k=3.0) ---")
	iqrKCustom := 3.0
	outliersCustom, nonOutliersCustom, statsCustom := findOutliersIQR(hardcodedDataPoints, iqrKCustom)
	printResults(hardcodedDataPoints, outliersCustom, nonOutliersCustom, statsCustom, iqrKCustom)


	// Example 2: Reading data from a file
	// To test this, create a file named "data.csv" in the same directory
	// with numbers, one per line, e.g.:
	// 10.5
	// 12.1
	// 100.0
	// 11.0
	// 9.8
	// -20.0
	// 13.5
	fmt.Println("\n--- Attempting to Read Data from data.csv ---")
	filePath := "data.csv"
	fileDataPoints, err := readDataFromFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from file '%s': %v\n", filePath, err)
		fmt.Println("Please create a 'data.csv' file with numerical data (one number per line or in the first column) to test file input.")
	} else {
		fmt.Printf("Successfully read %d data points from %s.\n", len(fileDataPoints), filePath)
		// Use the default IQR multiplier for file data
		fileOutliers, fileNonOutliers, fileStats := findOutliersIQR(fileDataPoints, 1.5)
		printResults(fileDataPoints, fileOutliers, fileNonOutliers, fileStats, 1.5)
	}

	// Additional functionality: User input for k value
	fmt.Println("\n--- Enter a custom IQR multiplier (e.g., 1.5, 2.0, 3.0) or press Enter for default (1.5): ---")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	customK := 1.5 // Default
	if input != "" {
		parsedK, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Println("Invalid input for k. Using default 1.5.")
		} else {
			customK = parsedK
		}
	}

	// Determine which dataset to use for user-defined K
	dataToUse := hardcodedDataPoints
	if len(fileDataPoints) > 0 {
		dataToUse = fileDataPoints
	}

	if len(dataToUse) > 0 {
		fmt.Printf("\n--- Running with User-Defined IQR Multiplier (k=%.1f) on Current Data ---\n", customK)
		userOutliers, userNonOutliers, userStats := findOutliersIQR(dataToUse, customK)
		printResults(dataToUse, userOutliers, userNonOutliers, userStats, customK)
	} else {
		fmt.Println("No data available to apply custom K value.")
	}
}