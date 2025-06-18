package main

import (
	"fmt"
	"math"
)

func mean(data []float64) float64 {
	if len(data) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func stdDev(data []float64) (float64, error) {
	n := float64(len(data))
	if n == 0 {
		return 0, fmt.Errorf("cannot calculate standard deviation for an empty slice")
	}

	m := mean(data)
	sumSqDiff := 0.0
	for _, v := range data {
		diff := v - m
		sumSqDiff += diff * diff
	}
	return math.Sqrt(sumSqDiff / n), nil
}

func covariance(x, y []float64) (float64, error) {
	n := float64(len(x))
	if n == 0 {
		return 0, fmt.Errorf("cannot calculate covariance for empty slices")
	}
	if len(x) != len(y) {
		return 0, fmt.Errorf("slices must have the same length for covariance calculation")
	}

	meanX := mean(x)
	meanY := mean(y)

	sumProductDiffs := 0.0
	for i := 0; i < len(x); i++ {
		sumProductDiffs += (x[i] - meanX) * (y[i] - meanY)
	}
	return sumProductDiffs / n, nil
}

func pearsonCorrelation(x, y []float64) (float64, error) {
	if len(x) == 0 || len(y) == 0 {
		return 0, fmt.Errorf("input slices cannot be empty")
	}
	if len(x) != len(y) {
		return 0, fmt.Errorf("input slices must have the same length")
	}

	cov, err := covariance(x, y)
	if err != nil {
		return 0, fmt.Errorf("error calculating covariance: %w", err)
	}

	stdDevX, err := stdDev(x)
	if err != nil {
		return 0, fmt.Errorf("error calculating standard deviation for X: %w", err)
	}

	stdDevY, err := stdDev(y)
	if err != nil {
		return 0, fmt.Errorf("error calculating standard deviation for Y: %w", err)
	}

	denominator := stdDevX * stdDevY
	if denominator == 0 {
		return 0, fmt.Errorf("cannot calculate correlation: standard deviation of one or both variables is zero (all values are the same)")
	}

	return cov / denominator, nil
}

func main() {
	dataX := []float64{10, 20, 30, 40, 50}
	dataY := []float64{15, 25, 35, 45, 55}

	dataA := []float64{1, 2, 3, 4, 5}
	dataB := []float64{5, 4, 3, 2, 1}

	dataP := []float64{1, 2, 3, 4, 5}
	dataQ := []float64{1, 1, 1, 1, 1}

	dataR := []float64{1, 2, 3}
	dataS := []float64{1, 2, 3, 4}

	dataT := []float64{}
	dataU := []float64{}

	fmt.Println("--- Correlation Examples ---")

	corr1, err1 := pearsonCorrelation(dataX, dataY)
	if err1 != nil {
		fmt.Printf("Error calculating correlation for dataX, dataY: %v\n", err1)
	} else {
		fmt.Printf("Correlation between dataX and dataY: %.4f\n", corr1)
	}

	corr2, err2 := pearsonCorrelation(dataA, dataB)
	if err2 != nil {
		fmt.Printf("Error calculating correlation for dataA, dataB: %v\n", err2)
	} else {
		fmt.Printf("Correlation between dataA and dataB: %.4f\n", corr2)
	}

	corr3, err3 := pearsonCorrelation(dataP, dataQ)
	if err3 != nil {
		fmt.Printf("Error calculating correlation for dataP, dataQ: %v\n", err3)
	} else {
		fmt.Printf("Correlation between dataP and dataQ: %.4f\n", corr3)
	}

	corr4, err4 := pearsonCorrelation(dataR, dataS)
	if err4 != nil {
		fmt.Printf("Error calculating correlation for dataR, dataS: %v\n", err4)
	} else {
		fmt.Printf("Correlation between dataR and dataS: %.4f\n", corr4)
	}

	corr5, err5 := pearsonCorrelation(dataT, dataU)
	if err5 != nil {
		fmt.Printf("Error calculating correlation for dataT, dataU: %v\n", err5)
	} else {
		fmt.Printf("Correlation between dataT and dataU: %.4f\n", corr5)
	}

	hoursStudied := []float64{10, 9, 2, 15, 10, 16, 11, 16}
	examScore := []float64{95, 80, 10, 90, 70, 90, 85, 90}

	corr6, err6 := pearsonCorrelation(hoursStudied, examScore)
	if err6 != nil {
		fmt.Printf("Error calculating correlation for hoursStudied, examScore: %v\n", err6)
	} else {
		fmt.Printf("Correlation between Hours Studied and Exam Score: %.4f\n", corr6)
	}
}

// Additional implementation at 2025-06-18 00:01:06
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// calculateMean calculates the arithmetic mean of a slice of float64.
func calculateMean(data []float64) float64 {
	sum := 0.0
	for _, val := range data {
		sum += val
	}
	return sum / float64(len(data))
}

// calculateCorrelation calculates the Pearson correlation coefficient between two data sets.
// It returns the coefficient and an error if the data sets are invalid (e.g., different lengths, empty).
func calculateCorrelation(x, y []float64) (float64, error) {
	if len(x) == 0 || len(y) == 0 {
		return 0, fmt.Errorf("input data sets cannot be empty")
	}
	if len(x) != len(y) {
		return 0, fmt.Errorf("input data sets must have the same length")
	}
	if len(x) < 2 { // Need at least two points for variance
		return 0, fmt.Errorf("at least two data points are required to calculate correlation")
	}

	meanX := calculateMean(x)
	meanY := calculateMean(y)

	numerator := 0.0
	sumSqDevX := 0.0
	sumSqDevY := 0.0

	for i := 0; i < len(x); i++ {
		devX := x[i] - meanX
		devY := y[i] - meanY
		numerator += devX * devY
		sumSqDevX += devX * devX
		sumSqDevY += devY * devY
	}

	denominator := math.Sqrt(sumSqDevX * sumSqDevY)

	if denominator == 0 {
		// This happens if one or both datasets have zero variance (all values are the same).
		// In such cases, Pearson correlation is undefined.
		return 0, fmt.Errorf("cannot calculate correlation: variance of one or both datasets is zero")
	}

	return numerator / denominator, nil
}

// parseNumbers parses a comma-separated string of numbers into a slice of float64.
func parseNumbers(input string) ([]float64, error) {
	if input == "" {
		return nil, fmt.Errorf("input cannot be empty")
	}

	parts := strings.Split(input, ",")
	data := make([]float64, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue // Skip empty parts if there are extra commas
		}
		val, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number '%s' in input: %w", part, err)
		}
		data = append(data, val)
	}
	return data, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Correlation Coefficient Calculator")
	fmt.Println("Enter data sets as comma-separated numbers (e.g., 1,2,3,4,5)")
	fmt.Println("Type 'exit' to quit at any data input prompt.")

	for {
		fmt.Println("\n--- New Calculation ---")

		fmt.Print("Enter data set X: ")
		scanner.Scan()
		rawX := strings.TrimSpace(scanner.Text())
		if strings.ToLower(rawX) == "exit" {
			fmt.Println("Exiting program.")
			return
		}
		xData, err := parseNumbers(rawX)
		if err != nil {
			fmt.Printf("Error parsing X data: %v\n", err)
			continue
		}

		fmt.Print("Enter data set Y: ")
		scanner.Scan()
		rawY := strings.TrimSpace(scanner.Text())
		if strings.ToLower(rawY) == "exit" {
			fmt.Println("Exiting program.")
			return
		}
		yData, err := parseNumbers(rawY)
		if err != nil {
			fmt.Printf("Error parsing Y data: %v\n", err)
			continue
		}

		correlation, err := calculateCorrelation(xData, yData)
		if err != nil {
			fmt.Printf("Error calculating correlation: %v\n", err)
		} else {
			fmt.Printf("Correlation Coefficient (r): %.4f\n", correlation)
		}
	}
}

// Additional implementation at 2025-06-18 00:02:07
package main

import (
	"errors"
	"fmt"
	"math"
)

func mean(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func stdDev(data []float64) float64 {
	if len(data) < 2 {
		return 0.0
	}
	m := mean(data)
	sumSqDiff := 0.0
	for _, v := range data {
		diff := v - m
		sumSqDiff += diff * diff
	}
	return math.Sqrt(sumSqDiff / float64(len(data)-1))
}

func calculatePearsonCorrelation(x, y []float64) (float64, error) {
	if len(x) == 0 || len(y) == 0 {
		return 0, errors.New("input slices cannot be empty")
	}
	if len(x) != len(y) {
		return 0, errors.New("input slices must have the same length")
	}
	if len(x) < 2 {
		return 0, errors.New("at least two data points are required for correlation calculation")
	}

	meanX := mean(x)
	meanY := mean(y)

	numerator := 0.0
	sumSqDiffX := 0.0
	sumSqDiffY := 0.0

	for i := 0; i < len(x); i++ {
		diffX := x[i] - meanX
		diffY := y[i] - meanY
		numerator += diffX * diffY
		sumSqDiffX += diffX * diffX
		sumSqDiffY += diffY * diffY
	}

	denominator := math.Sqrt(sumSqDiffX * sumSqDiffY)

	if denominator == 0 {
		return 0, nil
	}

	return numerator / denominator, nil
}

func main() {
	dataX := []float64{10, 20, 30, 40, 50}
	dataY := []float64{15, 25, 35, 45, 55}

	dataA := []float64{1, 2, 3, 4, 5}
	dataB := []float64{5, 4, 3, 2, 1}

	dataP := []float64{1, 2, 3, 4, 5}
	dataQ := []float64{1, 1, 1, 1, 1}

	fmt.Printf("Dataset X: %v\n", dataX)
	fmt.Printf("Mean X: %.2f\n", mean(dataX))
	fmt.Printf("Standard Deviation X: %.2f\n", stdDev(dataX))

	fmt.Printf("Dataset Y: %v\n", dataY)
	fmt.Printf("Mean Y: %.2f\n", mean(dataY))
	fmt.Printf("Standard Deviation Y: %.2f\n", stdDev(dataY))

	corrXY, errXY := calculatePearsonCorrelation(dataX, dataY)
	if errXY != nil {
		fmt.Printf("Error calculating correlation for X, Y: %v\n", errXY)
	} else {
		fmt.Printf("Pearson Correlation (X, Y): %.4f\n", corrXY)
	}
	fmt.Println("---")

	fmt.Printf("Dataset A: %v\n", dataA)
	fmt.Printf("Mean A: %.2f\n", mean(dataA))
	fmt.Printf("Standard Deviation A: %.2f\n", stdDev(dataA))

	fmt.Printf("Dataset B: %v\n", dataB)
	fmt.Printf("Mean B: %.2f\n", mean(dataB))
	fmt.Printf("Standard Deviation B: %.2f\n", stdDev(dataB))

	corrAB, errAB := calculatePearsonCorrelation(dataA, dataB)
	if errAB != nil {
		fmt.Printf("Error calculating correlation for A, B: %v\n", errAB)
	} else {
		fmt.Printf("Pearson Correlation (A, B): %.4f\n", corrAB)
	}
	fmt.Println("---")

	fmt.Printf("Dataset P: %v\n", dataP)
	fmt.Printf("Mean P: %.2f\n", mean(dataP))
	fmt.Printf("Standard Deviation P: %.2f\n", stdDev(dataP))

	fmt.Printf("Dataset Q: %v\n", dataQ)
	fmt.Printf("Mean Q: %.2f\n", mean(dataQ))
	fmt.Printf("Standard Deviation Q: %.2f\n", stdDev(dataQ))

	corrPQ, errPQ := calculatePearsonCorrelation(dataP, dataQ)
	if errPQ != nil {
		fmt.Printf("Error calculating correlation for P, Q: %v\n", errPQ)
	} else {
		fmt.Printf("Pearson Correlation (P, Q): %.4f\n", corrPQ)
	}
	fmt.Println("---")

	dataUnequalX := []float64{1, 2, 3}
	dataUnequalY := []float64{4, 5}
	fmt.Printf("Dataset UnequalX: %v\n", dataUnequalX)
	fmt.Printf("Dataset UnequalY: %v\n", dataUnequalY)
	corrUnequal, errUnequal := calculatePearsonCorrelation(dataUnequalX, dataUnequalY)
	if errUnequal != nil {
		fmt.Printf("Error calculating correlation for UnequalX, UnequalY: %v\n", errUnequal)
	} else {
		fmt.Printf("Pearson Correlation (UnequalX, UnequalY): %.4f\n", corrUnequal)
	}
	fmt.Println("---")

	dataEmptyX := []float64{}
	dataEmptyY := []float64{1, 2, 3}
	fmt.Printf("Dataset EmptyX: %v\n", dataEmptyX)
	fmt.Printf("Dataset EmptyY: %v\n", dataEmptyY)
	corrEmpty, errEmpty := calculatePearsonCorrelation(dataEmptyX, dataEmptyY)
	if errEmpty != nil {
		fmt.Printf("Error calculating correlation for EmptyX, EmptyY: %v\n", errEmpty)
	} else {
		fmt.Printf("Pearson Correlation (EmptyX, EmptyY): %.4f\n", corrEmpty)
	}
	fmt.Println("---")
}

// Additional implementation at 2025-06-18 00:02:51
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

// readData reads numerical data from a CSV file.
// It expects two columns of numbers.
func readData(filePath string) ([]float64, []float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 2 // Expect exactly two columns

	var xData []float64
	var yData []float64

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("error reading CSV record: %w", err)
		}

		if len(record) != 2 {
			return nil, nil, fmt.Errorf("invalid record format: expected 2 columns, got %d: %v", len(record), record)
		}

		x, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid number in column 1: %s", record[0])
		}
		y, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid number in column 2: %s", record[1])
		}

		xData = append(xData, x)
		yData = append(yData, y)
	}

	if len(xData) == 0 {
		return nil, nil, fmt.Errorf("no data found in the file")
	}
	if len(xData) != len(yData) {
		return nil, nil, fmt.Errorf("mismatched data lengths for X and Y")
	}

	return xData, yData, nil
}

// calculateMean calculates the arithmetic mean of a slice of float64.
func calculateMean(data []float64) float64 {
	sum := 0.0
	for _, val := range data {
		sum += val
	}
	return sum / float64(len(data))
}

// calculateStandardDeviation calculates the sample standard deviation of a slice of float64.
func calculateStandardDeviation(data []float64, mean float64) float64 {
	if len(data) < 2 {
		return 0.0 // Standard deviation is undefined or zero for less than 2 data points
	}
	sumOfSquares := 0.0
	for _, val := range data {
		diff := val - mean
		sumOfSquares += diff * diff
	}
	return math.Sqrt(sumOfSquares / float64(len(data)-1)) // Sample standard deviation
}

// calculateCovariance calculates the sample covariance between two slices of float64.
func calculateCovariance(x, y []float64, meanX, meanY float64) float64 {
	if len(x) != len(y) {
		panic("data slices must have the same length for covariance calculation")
	}
	if len(x) < 2 {
		return 0.0 // Covariance is undefined or zero for less than 2 data points
	}

	sumOfProducts := 0.0
	for i := 0; i < len(x); i++ {
		sumOfProducts += (x[i] - meanX) * (y[i] - meanY)
	}
	return sumOfProducts / float64(len(x)-1) // Sample covariance
}

// calculateCorrelationCoefficient calculates the Pearson correlation coefficient.
func calculateCorrelationCoefficient(x, y []float64) (float64, error) {
	if len(x) != len(y) {
		return 0, fmt.Errorf("data sets must have the same number of elements")
	}
	if len(x) < 2 {
		return 0, fmt.Errorf("at least two data points are required to calculate correlation")
	}

	meanX := calculateMean(x)
	meanY := calculateMean(y)

	stdDevX := calculateStandardDeviation(x, meanX)
	stdDevY := calculateStandardDeviation(y, meanY)

	if stdDevX == 0 || stdDevY == 0 {
		return 0, fmt.Errorf("standard deviation is zero for one or both datasets, correlation is undefined")
	}

	covariance := calculateCovariance(x, y, meanX, meanY)

	correlation := covariance / (stdDevX * stdDevY)
	return correlation, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <data_file.csv>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	xData, yData, err := readData(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading data: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully loaded %d data points from %s\n", len(xData), filePath)

	meanX := calculateMean(xData)
	meanY := calculateMean(yData)
	fmt.Printf("Mean of X: %.4f\n", meanX)
	fmt.Printf("Mean of Y: %.4f\n", meanY)

	stdDevX := calculateStandardDeviation(xData, meanX)
	stdDevY := calculateStandardDeviation(yData, meanY)
	fmt.Printf("Standard Deviation of X: %.4f\n", stdDevX)
	fmt.Printf("Standard Deviation of Y: %.4f\n", stdDevY)

	covariance := calculateCovariance(xData, yData, meanX, meanY)
	fmt.Printf("Covariance (X, Y): %.4f\n", covariance)

	correlation, err := calculateCorrelationCoefficient(xData, yData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calculating correlation: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Pearson Correlation Coefficient: %.4f\n", correlation)
}