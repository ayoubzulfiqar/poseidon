package main

import (
	"fmt"
	"math"
	"strconv"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func ToFraction(f float64) string {
	if f == 0 {
		return "0"
	}

	isNegative := f < 0
	absF := math.Abs(f)

	numerator := absF
	denominator := 1.0
	epsilon := 1e-9
	maxDenominator := 1e15

	for math.Abs(numerator-math.Round(numerator)) > epsilon && denominator < maxDenominator {
		numerator *= 10
		denominator *= 10
	}

	numInt := int64(math.Round(numerator))
	denInt := int64(math.Round(denominator))

	commonDivisor := gcd(numInt, denInt)
	numInt /= commonDivisor
	denInt /= commonDivisor

	if denInt == 1 {
		if isNegative {
			return "-" + strconv.FormatInt(numInt, 10)
		}
		return strconv.FormatInt(numInt, 10)
	}

	if isNegative {
		return "-" + strconv.FormatInt(numInt, 10) + "/" + strconv.FormatInt(denInt, 10)
	}
	return strconv.FormatInt(numInt, 10) + "/" + strconv.FormatInt(denInt, 10)
}

func main() {
	fmt.Println("0.5 =", ToFraction(0.5))
	fmt.Println("0.25 =", ToFraction(0.25))
	fmt.Println("0.125 =", ToFraction(0.125))
	fmt.Println("0.333 =", ToFraction(0.333))
	fmt.Println("0.666667 =", ToFraction(0.666667))
	fmt.Println("1.5 =", ToFraction(1.5))
	fmt.Println("3.25 =", ToFraction(3.25))
	fmt.Println("5.0 =", ToFraction(5.0))
	fmt.Println("0.0 =", ToFraction(0.0))
	fmt.Println("-0.5 =", ToFraction(-0.5))
	fmt.Println("-2.75 =", ToFraction(-2.75))
	fmt.Println("0.0001 =", ToFraction(0.0001))
	fmt.Println("0.0000001 =", ToFraction(0.0000001))
	fmt.Println("0.000000001 =", ToFraction(0.000000001))
	fmt.Println("0.0000000000000001 =", ToFraction(0.0000000000000001))
	fmt.Println("0.1 =", ToFraction(0.1))
	fmt.Println("0.9999999999999999 =", ToFraction(0.9999999999999999))
	fmt.Println("2/3 (float) =", ToFraction(2.0/3.0))
	fmt.Println("1/3 (float) =", ToFraction(1.0/3.0))
}

// Additional implementation at 2025-06-20 23:22:09
package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// gcd calculates the Greatest Common Divisor of two integers using the Euclidean algorithm.
func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// simplifyFraction simplifies a fraction given its numerator and denominator.
func simplifyFraction(numerator, denominator int64) (int64, int64) {
	if denominator == 0 {
		return numerator, denominator // Division by zero, handle upstream
	}
	if numerator == 0 {
		return 0, 1 // 0/X is always 0/1
	}

	commonDivisor := gcd(abs(numerator), abs(denominator))
	return numerator / commonDivisor, denominator / commonDivisor
}

// abs returns the absolute value of an int64.
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// decimalToFraction converts a decimal string to its fractional representation.
// It handles terminating decimals and returns the fraction as a string (e.g., "1/2").
// It returns an error for invalid input or if the number is too large/small to process.
func decimalToFraction(decimalStr string) (string, error) {
	// Trim whitespace
	decimalStr = strings.TrimSpace(decimalStr)

	// Handle empty string
	if decimalStr == "" {
		return "", fmt.Errorf("empty input string")
	}

	// Determine sign
	isNegative := false
	if strings.HasPrefix(decimalStr, "-") {
		isNegative = true
		decimalStr = decimalStr[1:]
	} else if strings.HasPrefix(decimalStr, "+") {
		decimalStr = decimalStr[1:]
	}

	// Split into integer and fractional parts
	parts := strings.Split(decimalStr, ".")
	integerPartStr := parts[0]
	fractionalPartStr := ""
	if len(parts) > 1 {
		fractionalPartStr = parts[1]
	}
	if len(parts) > 2 {
		return "", fmt.Errorf("invalid decimal format: multiple decimal points")
	}

	// Convert integer part
	integerPart, err := strconv.ParseInt(integerPartStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid integer part '%s': %w", integerPartStr, err)
	}

	// Calculate numerator and denominator
	var numerator, denominator int64

	if fractionalPartStr == "" {
		// No fractional part, e.g., "5" -> "5/1"
		numerator = integerPart
		denominator = 1
	} else {
		// Calculate power of 10 for denominator
		numDecimalPlaces := len(fractionalPartStr)
		if numDecimalPlaces > 18 { // Limit to avoid int64 overflow for math.Pow10 (max 10^18 fits in int64)
			return "", fmt.Errorf("too many decimal places for precise conversion (max 18 digits after decimal point)")
		}
		powerOf10 := int64(math.Pow10(numDecimalPlaces))

		// Convert fractional part to integer
		fractionalPartInt, err := strconv.ParseInt(fractionalPartStr, 10, 64)
		if err != nil {
			return "", fmt.Errorf("invalid fractional part '%s': %w", fractionalPartStr, err)
		}

		// Combine integer and fractional parts for the numerator
		// Example: 1.25 -> (1 * 100) + 25 = 125
		numerator = integerPart*powerOf10 + fractionalPartInt
		denominator = powerOf10
	}

	// Apply sign to numerator
	if isNegative {
		numerator = -numerator
	}

	// Simplify the fraction
	numerator, denominator = simplifyFraction(numerator, denominator)

	// Handle special case for 0
	if numerator == 0 {
		return "0/1", nil
	}

	return fmt.Sprintf("%d/%d", numerator, denominator), nil
}

// fractionToDecimal converts a fraction string (e.g., "1/2", "-3/4") to a float64.
// It returns an error for invalid input or division by zero.
func fractionToDecimal(fractionStr string) (float64, error) {
	fractionStr = strings.TrimSpace(fractionStr)
	parts := strings.Split(fractionStr, "/")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid fraction format: %s. Expected 'numerator/denominator'", fractionStr)
	}

	numeratorStr := parts[0]
	denominatorStr := parts[1]

	numerator, err := strconv.ParseFloat(numeratorStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid numerator '%s': %w", numeratorStr, err)
	}

	denominator, err := strconv.ParseFloat(denominatorStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid denominator '%s': %w", denominatorStr, err)
	}

	if denominator == 0 {
		return 0, fmt.Errorf("division by zero: denominator is zero")
	}

	return numerator / denominator, nil
}

func main() {
	// Test cases for decimalToFraction
	decimalTests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"0", "0/1", false},
		{"1", "1/1", false},
		{"0.5", "1/2", false},
		{"-0.5", "-1/2", false},
		{"1.25", "5/4", false},
		{"-1.25", "-5/4", false},
		{"0.125", "1/8", false},
		{"0.333", "333/1000", false}, // Note: This is for terminating decimal, not repeating 1/3
		{"10.0", "10/1", false},
		{"-7", "-7/1", false},
		{"0.0", "0/1", false},
		{"0.00", "0/1", false},
		{"123", "123/1", false},
		{"0.0001", "1/10000", false},
		{"-0.0001", "-1/10000", false},
		{"2.75", "11/4", false},
		{"", "", true},
		{"abc", "", true},
		{"1.2.3", "", true},
		{"1.000000000000000001", "", true}, // Too many decimal places (18 is max for int64 power of 10)
		{"9223372036854775807.0", "9223372036854775807/1", false}, // Max int64
		{"-9223372036854775808.0", "-9223372036854775808/1", false}, // Min int64
	}

	fmt.Println("--- Decimal to Fraction Conversion ---")
	for _, tt := range decimalTests {
		fraction, err := decimalToFraction(tt.input)
		if (err != nil) != tt.hasError {
			fmt.Printf("FAIL: decimalToFraction(\"%s\") expected error: %t, got error: %v\n", tt.input, tt.hasError, err)
		} else if !tt.hasError && fraction != tt.expected {
			fmt.Printf("FAIL: decimalToFraction(\"%s\") expected \"%s\", got \"%s\"\n", tt.input, tt.expected, fraction)
		} else if !tt.hasError {
			fmt.Printf("PASS: decimalToFraction(\"%s\") -> \"%s\"\n", tt.input, fraction)
		} else {
			fmt.Printf("PASS: decimalToFraction(\"%s\") -> Error: %v\n", tt.input, err)
		}
	}

	fmt.Println("\n--- Fraction to Decimal Conversion ---")
	fractionTests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"1/2", 0.5, false},
		{"-1/2", -0.5, false},
		{"5/4", 1.25, false},
		{"1/1", 1.0, false},
		{"0/1", 0.0, false},
		{"10/1", 10.0, false},
		{"333/1000", 0.333, false},
		{"1/0", 0.0, true},
		{"abc/def", 0.0, true},
		{"1", 0.0, true}, // Invalid format
		{"", 0.0, true},
	}

	for _, tt := range fractionTests {
		decimal, err := fractionToDecimal(tt.input)
		if (err != nil) != tt.hasError {
			fmt.Printf("FAIL: fractionToDecimal(\"%s\") expected error: %t, got error: %v\n", tt.input, tt.hasError, err)
		} else if !tt.hasError && math.Abs(decimal-tt.expected) > 1e-9 { // Use a small epsilon for float comparison
			fmt.Printf("FAIL: fractionToDecimal(\"%s\") expected %f, got %f\n", tt.input, tt.expected, decimal)
		} else if !tt.hasError {
			fmt.Printf("PASS: fractionToDecimal(\"%s\") -> %f\n", tt.input, decimal)
		} else {
			fmt.Printf("PASS: fractionToDecimal(\"%s\") -> Error: %v\n", tt.input, err)
		}
	}

	fmt.Println("\n--- Combined Example ---")
	testDecimal := "3.14159"
	fmt.Printf("Original Decimal: %s\n", testDecimal)
	fraction, err := decimalToFraction(testDecimal)
	if err != nil {
		fmt.Printf("Error converting to fraction: %v\n", err)
	} else {
		fmt.Printf("Converted to Fraction: %s\n", fraction)
		decimalBack, err := fractionToDecimal(fraction)
		if err != nil {
			fmt.Printf("Error converting back to decimal: %v\n", err)
		} else {
			fmt.Printf("Converted back to Decimal: %f\n", decimalBack)
		}
	}

	testDecimal = "0.666666" // Approximation of 2/3
	fmt.Printf("\nOriginal Decimal: %s\n", testDecimal)
	fraction, err = decimalToFraction(testDecimal)
	if err != nil {
		fmt.Printf("Error converting to fraction: %v\n", err)
	} else {
		fmt.Printf("Converted to Fraction: %s\n", fraction)
		decimalBack, err = fractionToDecimal(fraction)
		if err != nil {
			fmt.Printf("Error converting back to decimal: %v\n", err)
		} else {
			fmt.Printf("Converted back to Decimal: %f\n", decimalBack)
		}
	}
}