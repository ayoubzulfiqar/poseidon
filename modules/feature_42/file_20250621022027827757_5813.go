package main

import (
	"fmt"
	"math"
)

// calculateMean calculates the mean of a slice of float64.
func calculateMean(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

// calculateCorrelation calculates the Pearson correlation coefficient between two datasets.
// It returns the correlation coefficient and an error if the inputs are invalid.
func calculateCorrelation(x, y []float64) (float64, error) {
	if len(x) == 0 || len(y) == 0 {
		return 0, fmt.Errorf("input datasets cannot be empty")
	}
	if len(x) != len(y) {
		return 0, fmt.Errorf("input datasets must have the same length")
	}
	if len(x) < 2 {
		return 0, fmt.Errorf("correlation requires at least two data points")
	}

	meanX := calculateMean(x)
	meanY := calculateMean(y)

	sumProdDeviations := 0.0
	sumSqDeviationsX := 0.0
	sumSqDeviationsY := 0.0

	for i := 0; i < len(x); i++ {
		devX := x[i] - meanX
		devY := y[i] - meanY
		sumProdDeviations += devX * devY
		sumSqDeviationsX += devX * devX
		sumSqDeviationsY += devY * devY
	}

	denominator := math.Sqrt(sumSqDeviationsX * sumSqDeviationsY)

	if denominator == 0 {
		return 0.0, nil
	}

	return sumProdDeviations / denominator, nil
}

func main() {
	dataX1 := []float64{1, 2, 3, 4, 5}
	dataY1 := []float64{2, 4, 5, 4, 5}
	corr1, err1 := calculateCorrelation(dataX1, dataY1)
	if err1 != nil {
		fmt.Printf("Error: %v\n", err1)
	} else {
		fmt.Printf("Correlation (dataX1, dataY1): %.4f\n", corr1)
	}

	dataX2 := []float64{1, 2, 3, 4, 5}
	dataY2 := []float64{5, 4, 3, 2, 1}
	corr2, err2 := calculateCorrelation(dataX2, dataY2)
	if err2 != nil {
		fmt.Printf("Error: %v\n", err2)
	} else {
		fmt.Printf("Correlation (dataX2, dataY2): %.4f\n", corr2)
	}

	dataX3 := []float64{1, 2, 3, 4, 5}
	dataY3 := []float64{1, 5, 2, 4, 3}
	corr3, err3 := calculateCorrelation(dataX3, dataY3)
	if err3 != nil {
		fmt.Printf("Error: %v\n", err3)
	} else {
		fmt.Printf("Correlation (dataX3, dataY3): %.4f\n", corr3)
	}

	dataX4 := []float64{10, 20, 30, 40, 50}
	dataY4 := []float64{10, 20, 30, 40, 50}
	corr4, err4 := calculateCorrelation(dataX4, dataY4)
	if err4 != nil {
		fmt.Printf("Error: %v\n", err4)
	} else {
		fmt.Printf("Correlation (dataX4, dataY4): %.4f\n", corr4)
	}

	dataX5 := []float64{1, 1, 1, 1, 1}
	dataY5 := []float64{1, 2, 3, 4, 5}
	corr5, err5 := calculateCorrelation(dataX5, dataY5)
	if err5 != nil {
		fmt.Printf("Error: %v\n", err5)
	} else {
		fmt.Printf("Correlation (dataX5, dataY5): %.4f\n", corr5)
	}

	dataX6 := []float64{1, 2, 3}
	dataY6 := []float64{4, 5, 6, 7}
	corr6, err6 := calculateCorrelation(dataX6, dataY6)
	if err6 != nil {
		fmt.Printf("Error: %v\n", err6)
	} else {
		fmt.Printf("Correlation (dataX6, dataY6): %.4f\n", corr6)
	}

	dataX7 := []float64{1}
	dataY7 := []float64{2}
	corr7, err7 := calculateCorrelation(dataX7, dataY7)
	if err7 != nil {
		fmt.Printf("Error: %v\n", err7)
	} else {
		fmt.Printf("Correlation (dataX7, dataY7): %.4f\n", corr7)
	}

	dataX8 := []float64{}
	dataY8 := []float64{}
	corr8, err8 := calculateCorrelation(dataX8, dataY8)
	if err8 != nil {
		fmt.Printf("Error: %v\n", err8)
	} else {
		fmt.Printf("Correlation (dataX8, dataY8): %.4f\n", corr8)
	}
}

// Additional implementation at 2025-06-21 02:21:13
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func calculateCorrelation(x, y []float64) (float64, error) {
	n := len(x)
	if n != len(y) {
		return 0, fmt.Errorf("datasets must have the same number of elements: len(x)=%d, len(y)=%d", n, len(y))
	}
	if n < 2 {
		return 0, fmt.Errorf("at least two data points are required to calculate correlation, got %d", n)
	}

	var sumX, sumY, sumXY, sumX2, sumY2 float64
	for i := 0; i < n; i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
		sumY2 += y[i] * y[i]
	}

	numerator := float64(n)*sumXY - sumX*sumY
	denominatorX := float64(n)*sumX2 - sumX*sumX
	denominatorY := float64(n)*sumY2 - sumY*sumY

	denominator := math.Sqrt(denominatorX * denominatorY)

	if denominator == 0 {
		return 0, fmt.Errorf("cannot calculate correlation: one or both datasets have no variance (all values are the same)")
	}

	return numerator / denominator, nil
}

func readFloatData(prompt string) ([]float64, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return nil, fmt.Errorf("no data entered")
	}

	parts := strings.Fields(input)
	data := make([]float64, len(parts))
	for i, part := range parts {
		val, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number '%s': %w", part, err)
		}
		data[i] = val
	}
	return data, nil
}

func main() {
	fmt.Println("Pearson Correlation Coefficient Calculator")
	fmt.Println("----------------------------------------")
	fmt.Println("Enter data points separated by spaces (e.g., 1.0 2.5 3.0):")

	xData, err := readFloatData("Enter X values: ")
	if err != nil {
		fmt.Printf("Error reading X data: %v\n", err)
		return
	}

	yData, err := readFloatData("Enter Y values: ")
	if err != nil {
		fmt.Printf("Error reading Y data: %v\n", err)
		return
	}

	correlation, err := calculateCorrelation(xData, yData)
	if err != nil {
		fmt.Printf("Error calculating correlation: %v\n", err)
		return
	}

	fmt.Printf("\nCorrelation Coefficient (r): %.4f\n", correlation)
}

// Additional implementation at 2025-06-21 02:21:41
package main

import (
	"errors"
	"fmt"
	"math"
)

// calculateMean calculates the arithmetic mean of a slice of float64 numbers.
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

// calculateStandardDeviation calculates the sample standard deviation of a slice of float64 numbers.
// It requires the mean of the data.
func calculateStandardDeviation(data []float64, mean float64) float64 {
	if len(data) < 2 { // Standard deviation requires at least 2 data points for sample std dev
		return 0.0
	}
	sumOfSquares := 0.0
	for _, val := range data {
		diff := val - mean
		sumOfSquares += diff * diff
	}
	// Using N-1 for sample standard deviation
	return math.Sqrt(sumOfSquares / float64(len(data)-1))
}

// calculateCovariance calculates the sample covariance between two slices of float64 numbers.
// It requires the means of both datasets.
func calculateCovariance(x, y []float64, meanX, meanY float64) float64 {
	if len(x) != len(y) || len(x) < 2 {
		return 0.0 // Should be caught by calculatePearsonCorrelation's checks
	}
	sumProductOfDeviations := 0.0
	for i := 0; i < len(x); i++ {
		sumProductOfDeviations += (x[i] - meanX) * (y[i] - meanY)
	}
	// Using N-1 for sample covariance
	return sumProductOfDeviations / float64(len(x)-1)
}

// calculatePearsonCorrelation calculates the Pearson correlation coefficient between two datasets.
// It returns the coefficient and an error if the inputs are invalid.
func calculatePearsonCorrelation(x, y []float64) (float64, error) {
	if len(x) == 0 || len(y) == 0 {
		return 0.0, errors.New("input datasets cannot be empty")
	}
	if len(x) != len(y) {
		return 0.0, errors.New("input datasets must have the same number of elements")
	}
	if len(x) < 2 {
		return 0.0, errors.New("at least two data points are required to calculate correlation")
	}

	meanX := calculateMean(x)
	meanY := calculateMean(y)

	stdDevX := calculateStandardDeviation(x, meanX)
	stdDevY := calculateStandardDeviation(y, meanY)

	if stdDevX == 0 || stdDevY == 0 {
		// If standard deviation is zero, all values in the dataset are the same.
		// Correlation is undefined in this case as there's no variance.
		return 0.0, errors.New("standard deviation of one or both datasets is zero (all values are the same)")
	}

	covariance := calculateCovariance(x, y, meanX, meanY)

	correlation := covariance / (stdDevX * stdDevY)

	return correlation, nil
}

// printDataset prints a given dataset with a label for clarity.
func printDataset(label string, data []float64) {
	fmt.Printf("%s: %v\n", label, data)
}

func main() {
	// Example 1: Positive correlation
	dataX1 := []float64{10, 20, 30, 40, 50}
	dataY1 := []float64{15, 25, 35, 45, 55}
	fmt.Println("--- Example 1: Positive Correlation ---")
	printDataset("Data X1", dataX1)
	printDataset("Data Y1", dataY1)
	corr1, err1 := calculatePearsonCorrelation(dataX1, dataY1)
	if err1 != nil {
		fmt.Printf("Error calculating correlation: %v\n", err1)
	} else {
		fmt.Printf("Pearson Correlation (X1, Y1): %.4f\n", corr1)
	}
	fmt.Println()

	// Example 2: Negative correlation
	dataX2 := []float64{10, 20, 30, 40, 50}
	dataY2 := []float64{50, 40, 30, 20, 10}
	fmt.Println("--- Example 2: Negative Correlation ---")
	printDataset("Data X2", dataX2)
	printDataset("Data Y2", dataY2)
	corr2, err2 := calculatePearsonCorrelation(dataX2, dataY2)
	if err2 != nil {
		fmt.Printf("Error calculating correlation: %v\n", err2)
	} else {
		fmt.Printf("Pearson Correlation (X2, Y2): %.4f\n", corr2)
	}
	fmt.Println()

	// Example 3: No correlation (random-ish)
	dataX3 := []float64{1, 2, 3, 4, 5}
	dataY3 := []float64{5, 1, 4, 2, 3}
	fmt.Println("--- Example 3: No Correlation ---")
	printDataset("Data X3", dataX3)
	printDataset("Data Y3", dataY3)
	corr3, err3 := calculatePearsonCorrelation(dataX3, dataY3)
	if err3 != nil {
		fmt.Printf("Error calculating correlation: %v\n", err3)
	} else {
		fmt.Printf("Pearson Correlation (X3, Y3): %.4f\n", corr3)
	}
	fmt.Println()

	// Example 4: Error case - different lengths
	dataX4 := []float64{1, 2, 3}
	dataY4 := []float64{4, 5}
	fmt.Println("--- Example 4: Error - Different Lengths ---")
	printDataset("Data X4", dataX4)
	printDataset("Data Y4", dataY4)
	corr4, err4 := calculatePearsonCorrelation(dataX4, dataY4)
	if err4 != nil {
		fmt.Printf("Error calculating correlation: %v\n", err4)
	} else {
		fmt.Printf("Pearson Correlation (X4, Y4): %.4f\n", corr4)
	}
	fmt.Println()

	// Example 5: Error case - insufficient data points
	dataX5 := []float64{1}
	dataY5 := []float64{2}
	fmt.Println("--- Example 5: Error - Insufficient Data ---")
	printDataset("Data X5", dataX5)
	printDataset("Data Y5", dataY5)
	corr5, err5 := calculatePearsonCorrelation(dataX5, dataY5)
	if err5 != nil {
		fmt.Printf("Error calculating correlation: %v\n", err5)
	} else {
		fmt.Printf("Pearson Correlation (X5, Y5): %.4f\n", corr5)
	}
	fmt.Println()

	// Example 6: Error case - zero standard deviation in one dataset
	dataX6 := []float64{5, 5, 5, 5, 5}
	dataY6 := []float64{1, 2, 3, 4, 5}
	fmt.Println("--- Example 6: Error - Zero Standard Deviation ---")
	printDataset("Data X6", dataX6)
	printDataset("Data Y6", dataY6)
	corr6, err6 := calculatePearsonCorrelation(dataX6, dataY6)
	if err6 != nil {
		fmt.Printf("Error calculating correlation: %v\n", err6)
	} else {
		fmt.Printf("Pearson Correlation (X6, Y6): %.4f\n", corr6)
	}
	fmt.Println()
}