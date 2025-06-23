package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Complex struct {
	Real float64
	Imag float64
}

func (c Complex) String() string {
	if c.Imag == 0 {
		return fmt.Sprintf("%.4f", c.Real)
	}
	if c.Real == 0 {
		return fmt.Sprintf("%.4fi", c.Imag)
	}
	if c.Imag < 0 {
		return fmt.Sprintf("%.4f - %.4fi", c.Real, -c.Imag)
	}
	return fmt.Sprintf("%.4f + %.4fi", c.Real, c.Imag)
}

func (c Complex) Add(other Complex) Complex {
	return Complex{Real: c.Real + other.Real, Imag: c.Imag + other.Imag}
}

func (c Complex) Subtract(other Complex) Complex {
	return Complex{Real: c.Real - other.Real, Imag: c.Imag - other.Imag}
}

func (c Complex) Multiply(other Complex) Complex {
	realPart := c.Real*other.Real - c.Imag*other.Imag
	imagPart := c.Real*other.Imag + c.Imag*other.Real
	return Complex{Real: realPart, Imag: imagPart}
}

func (c Complex) Divide(other Complex) (Complex, error) {
	denominator := other.Real*other.Real + other.Imag*other.Imag
	if denominator == 0 {
		return Complex{}, fmt.Errorf("division by zero complex number")
	}
	realPart := (c.Real*other.Real + c.Imag*other.Imag) / denominator
	imagPart := (c.Imag*other.Real - c.Real*other.Imag) / denominator
	return Complex{Real: realPart, Imag: imagPart}, nil
}

func (c Complex) Magnitude() float64 {
	return math.Sqrt(c.Real*c.Real + c.Imag*c.Imag)
}

func (c Complex) Conjugate() Complex {
	return Complex{Real: c.Real, Imag: -c.Imag}
}

func readFloat(prompt string) (float64, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	val, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %v", err)
	}
	return val, nil
}

func readComplex(prompt string) (Complex, error) {
	fmt.Println(prompt)
	realPart, err := readFloat("  Enter real part: ")
	if err != nil {
		return Complex{}, err
	}
	imagPart, err := readFloat("  Enter imaginary part: ")
	if err != nil {
		return Complex{}, err
	}
	return Complex{Real: realPart, Imag: imagPart}, nil
}

func readOperation() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter operation (+, -, *, /, mag, conj, exit): ")
	op, _ := reader.ReadString('\n')
	op = strings.TrimSpace(strings.ToLower(op))
	switch op {
	case "+", "-", "*", "/", "mag", "conj", "exit":
		return op, nil
	default:
		return "", fmt.Errorf("invalid operation: %s", op)
	}
}

func main() {
	fmt.Println("Welcome to the Go Complex Number Calculator!")

	for {
		op, err := readOperation()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		if op == "exit" {
			fmt.Println("Exiting calculator. Goodbye!")
			break
		}

		var c1, c2 Complex
		var err1, err2 error

		if op == "mag" || op == "conj" {
			c1, err1 = readComplex("Enter the complex number:")
			if err1 != nil {
				fmt.Println("Error:", err1)
				continue
			}
		} else {
			c1, err1 = readComplex("Enter the first complex number:")
			if err1 != nil {
				fmt.Println("Error:", err1)
				continue
			}
			c2, err2 = readComplex("Enter the second complex number:")
			if err2 != nil {
				fmt.Println("Error:", err2)
				continue
			}
		}

		var result Complex
		var resErr error
		var floatResult float64

		switch op {
		case "+":
			result = c1.Add(c2)
			fmt.Printf("Result: %s + %s = %s\n", c1, c2, result)
		case "-":
			result = c1.Subtract(c2)
			fmt.Printf("Result: %s - %s = %s\n", c1, c2, result)
		case "*":
			result = c1.Multiply(c2)
			fmt.Printf("Result: %s * %s = %s\n", c1, c2, result)
		case "/":
			result, resErr = c1.Divide(c2)
			if resErr != nil {
				fmt.Println("Error:", resErr)
			} else {
				fmt.Printf("Result: %s / %s = %s\n", c1, c2, result)
			}
		case "mag":
			floatResult = c1.Magnitude()
			fmt.Printf("Magnitude of %s = %.4f\n", c1, floatResult)
		case "conj":
			result = c1.Conjugate()
			fmt.Printf("Conjugate of %s = %s\n", c1, result)
		default:
			fmt.Println("Unknown operation. Please try again.")
		}
		fmt.Println("------------------------------------")
	}
}

// Additional implementation at 2025-06-22 23:29:37
package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Complex represents a complex number.
type Complex struct {
	Real float64
	Imag float64
}

// NewComplex creates a new Complex number.
func NewComplex(real, imag float64) Complex {
	return Complex{Real: real, Imag: imag}
}

// Add returns the sum of two complex numbers.
func (c Complex) Add(other Complex) Complex {
	return NewComplex(c.Real+other.Real, c.Imag+other.Imag)
}

// Subtract returns the difference of two complex numbers.
func (c Complex) Subtract(other Complex) Complex {
	return NewComplex(c.Real-other.Real, c.Imag-other.Imag)
}

// Multiply returns the product of two complex numbers.
// (a + bi)(c + di) = (ac - bd) + (ad + bc)i
func (c Complex) Multiply(other Complex) Complex {
	realPart := c.Real*other.Real - c.Imag*other.Imag
	imagPart := c.Real*other.Imag + c.Imag*other.Real
	return NewComplex(realPart, imagPart)
}

// Divide returns the quotient of two complex numbers.
// (a + bi) / (c + di) = ((ac + bd) + (bc - ad)i) / (c^2 + d^2)
func (c Complex) Divide(other Complex) (Complex, error) {
	denominator := other.Real*other.Real + other.Imag*other.Imag
	if denominator == 0 {
		return Complex{}, fmt.Errorf("division by zero complex number")
	}
	realPart := (c.Real*other.Real + c.Imag*other.Imag) / denominator
	imagPart := (c.Imag*other.Real - c.Real*other.Imag) / denominator
	return NewComplex(realPart, imagPart), nil
}

// Abs returns the magnitude (absolute value) of the complex number.
// |a + bi| = sqrt(a^2 + b^2)
func (c Complex) Abs() float64 {
	return math.Sqrt(c.Real*c.Real + c.Imag*c.Imag)
}

// Conjugate returns the complex conjugate of the number.
// Conjugate of (a + bi) is (a - bi)
func (c Complex) Conjugate() Complex {
	return NewComplex(c.Real, -c.Imag)
}

// String returns a string representation of the complex number.
func (c Complex) String() string {
	if c.Imag == 0 {
		return fmt.Sprintf("%.4f", c.Real)
	}
	if c.Real == 0 {
		return fmt.Sprintf("%.4fi", c.Imag)
	}
	if c.Imag < 0 {
		return fmt.Sprintf("%.4f - %.4fi", c.Real, math.Abs(c.Imag))
	}
	return fmt.Sprintf("%.4f + %.4fi", c.Real, c.Imag)
}

// ParseComplex parses a string into a Complex number.
// Supports formats like "3", "4i", "3+4i", "3-4i", "-2-5i", "i", "-i".
func ParseComplex(s string) (Complex, error) {
	s = strings.ReplaceAll(s, " ", "") // Remove spaces

	// Handle pure imaginary numbers like "i" or "-i"
	if s == "i" {
		return NewComplex(0, 1), nil
	}
	if s == "-i" {
		return NewComplex(0, -1), nil
	}

	// Check for 'i' to determine if it's a complex number
	iIndex := strings.LastIndex(s, "i")

	if iIndex == -1 { // No 'i', assume it's a pure real number
		realPart, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return Complex{}, fmt.Errorf("invalid real number format: %w", err)
		}
		return NewComplex(realPart, 0), nil
	}

	// It contains 'i', so it's either pure imaginary or complex
	var realStr, imagStr string
	var sign int = 1 // 1 for +, -1 for -

	// Find the sign separating real and imaginary parts
	plusIndex := strings.LastIndex(s, "+")
	minusIndex := strings.LastIndex(s, "-") // This needs careful handling for negative real parts

	// Determine the split point for real and imaginary parts
	splitIndex := -1
	if plusIndex != -1 && plusIndex < iIndex { // Ensure '+' is before 'i'
		splitIndex = plusIndex
	}
	// If there's a minus, and it's not the very first character (negative real part)
	// and it's before 'i'
	if minusIndex != -1 && minusIndex > 0 && minusIndex < iIndex {
		if splitIndex == -1 || minusIndex > splitIndex { // Take the last operator before 'i'
			splitIndex = minusIndex
		}
	}

	if splitIndex != -1 { // Found a separator, it's a+bi or a-bi
		realStr = s[:splitIndex]
		imagStr = s[splitIndex : iIndex]
	} else { // No separator, it's just bi or -bi
		imagStr = s[:iIndex]
	}

	// Parse real part
	realPart := 0.0
	if realStr != "" {
		r, err := strconv.ParseFloat(realStr, 64)
		if err != nil {
			return Complex{}, fmt.Errorf("invalid real part '%s': %w", realStr, err)
		}
		realPart = r
	}

	// Parse imaginary part
	imagPart := 0.0
	if imagStr == "" || imagStr == "+" { // "i" or "a+i"
		imagPart = 1.0
	} else if imagStr == "-" { // "-i" or "a-i"
		imagPart = -1.0
	} else {
		// Remove the sign if it's already handled by the split
		if imagStr[0] == '+' || imagStr[0] == '-' {
			imagStr = imagStr[1:]
		}
		im, err := strconv.ParseFloat(imagStr, 64)
		if err != nil {
			return Complex{}, fmt.Errorf("invalid imaginary part '%s': %w", imagStr, err)
		}
		imagPart = im
	}

	return NewComplex(realPart, imagPart), nil
}

func main() {
	fmt.Println("--- Complex Number Calculator ---")

	// Example 1: Basic Operations
	c1 := NewComplex(3, 4)
	c2 := NewComplex(1, -2)
	fmt.Printf("c1 = %s\n", c1)
	fmt.Printf("c2 = %s\n", c2)

	sum := c1.Add(c2)
	fmt.Printf("c1 + c2 = %s\n", sum)

	diff := c1.Subtract(c2)
	fmt.Printf("c1 - c2 = %s\n", diff)

	prod := c1.Multiply(c2)
	fmt.Printf("c1 * c2 = %s\n", prod)

	quotient, err := c1.Divide(c2)
	if err != nil {
		fmt.Printf("Error dividing c1 by c2: %v\n", err)
	} else {
		fmt.Printf("c1 / c2 = %s\n", quotient)
	}

	fmt.Printf("Magnitude of c1 (|c1|) = %.4f\n", c1.Abs())
	fmt.Printf("Conjugate of c1 (c1*) = %s\n", c1.Conjugate())

	fmt.Println("\n--- Division by Zero Test ---")
	c3 := NewComplex(0, 0)
	_, err = c1.Divide(c3)
	if err != nil {
		fmt.Printf("Attempted c1 / (0+0i): %v\n", err)
	}

	fmt.Println("\n--- Parsing Examples ---")
	testStrings := []string{
		"3",
		"4i",
		"3+4i",
		"3-4i",
		"-2-5i",
		"i",
		"-i",
		"1.23",
		"-0.5+1.7i",
		"7.0i",
		" -1.0 - 2.0i ",
	}

	for _, s := range testStrings {
		parsed, err := ParseComplex(s)
		if err != nil {
			fmt.Printf("Failed to parse \"%s\": %v\n", s, err)
		} else {
			fmt.Printf("Parsed \"%s\" as %s\n", s, parsed)
		}
	}

	fmt.Println("\n--- Combined Operations with Parsing ---")
	s1 := "2+3i"
	s2 := "1-i"

	p1, err1 := ParseComplex(s1)
	p2, err2 := ParseComplex(s2)

	if err1 != nil || err2 != nil {
		fmt.Printf("Error parsing: %v, %v\n", err1, err2)
		return
	}

	fmt.Printf("(%s + %s) * %s = %s\n", p1, p2, p1.Conjugate(), (p1.Add(p2)).Multiply(p1.Conjugate()))

	s3 := "10"
	s4 := "2i"
	p3, err3 := ParseComplex(s3)
	p4, err4 := ParseComplex(s4)
	if err3 != nil || err4 != nil {
		fmt.Printf("Error parsing: %v, %v\n", err3, err4)
		return
	}
	divResult, err := p3.Divide(p4)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("%s / %s = %s\n", p3, p4, divResult)
	}
}

// Additional implementation at 2025-06-22 23:30:28
package main

import (
	"fmt"
	"math"
)

// Complex represents a complex number z = real + imag*i
type Complex struct {
	real float64
	imag float64
}

// NewComplex creates a new Complex number
func NewComplex(r, i float64) Complex {
	return Complex{real: r, imag: i}
}

// String returns the string representation of the complex number
// Implements the fmt.Stringer interface.
func (c Complex) String() string {
	if c.imag == 0 {
		return fmt.Sprintf("%.4f", c.real)
	}
	if c.real == 0 {
		return fmt.Sprintf("%.4fi", c.imag)
	}
	if c.imag < 0 {
		return fmt.Sprintf("%.4f - %.4fi", c.real, math.Abs(c.imag))
	}
	return fmt.Sprintf("%.4f + %.4fi", c.real, c.imag)
}

// Add returns the sum of two complex numbers (c + other)
func (c Complex) Add(other Complex) Complex {
	return NewComplex(c.real+other.real, c.imag+other.imag)
}

// Subtract returns the difference of two complex numbers (c - other)
func (c Complex) Subtract(other Complex) Complex {
	return NewComplex(c.real-other.real, c.imag-other.imag)
}

// Multiply returns the product of two complex numbers (c * other)
// Formula: (a + bi)(c + di) = (ac - bd) + (ad + bc)i
func (c Complex) Multiply(other Complex) Complex {
	return NewComplex(c.real*other.real-c.imag*other.imag,
		c.real*other.imag+c.imag*other.real)
}

// Divide returns the quotient of two complex numbers (c / other)
// Formula: (a + bi) / (c + di) = [(ac + bd) + (bc - ad)i] / (c^2 + d^2)
// Handles division by zero by relying on IEEE 754 floating-point behavior (Inf/NaN).
func (c Complex) Divide(other Complex) Complex {
	denominator := other.real*other.real + other.imag*other.imag
	realPart := (c.real*other.real + c.imag*other.imag) / denominator
	imagPart := (c.imag*other.real - c.real*other.imag) / denominator
	return NewComplex(realPart, imagPart)
}

// Conjugate returns the complex conjugate of the number
// conjugate(a + bi) = a - bi
func (c Complex) Conjugate() Complex {
	return NewComplex(c.real, -c.imag)
}

// Magnitude returns the magnitude (absolute value or modulus) of the complex number
// |z| = sqrt(real^2 + imag^2)
func (c Complex) Magnitude() float64 {
	return math.Sqrt(c.real*c.real + c.imag*c.imag)
}

// Phase returns the phase (argument) of the complex number in radians
// arg(z) = atan2(imag, real)
func (c Complex) Phase() float64 {
	return math.Atan2(c.imag, c.real)
}

// Equals checks if two complex numbers are approximately equal within a given tolerance.
// Useful for comparing floating-point numbers.
func (c Complex) Equals(other Complex, tolerance float64) bool {
	return math.Abs(c.real-other.real) < tolerance &&
		math.Abs(c.imag-other.imag) < tolerance
}

// FromPolar creates a complex number from its polar coordinates (magnitude and phase in radians)
func FromPolar(magnitude, phase float64) Complex {
	return NewComplex(magnitude*math.Cos(phase), magnitude*math.Sin(phase))
}

// PowerInt raises the complex number to an integer power n.
// Uses De Moivre's theorem: (r(cos(theta) + i sin(theta)))^n = r^n(cos(n*theta) + i sin(n*theta))
func (c Complex) PowerInt(n int) Complex {
	if n == 0 {
		return NewComplex(1, 0) // z^0 = 1
	}
	if c.real == 0 && c.imag == 0 {
		return NewComplex(0, 0) // 0^n = 0 for n > 0
	}

	r := c.Magnitude()
	theta := c.Phase()

	rPowN := math.Pow(r, float64(n))
	nTheta := float64(n) * theta

	return NewComplex(rPowN*math.Cos(nTheta), rPowN*math.Sin(nTheta))
}

func main() {
	// --- Demonstration of Complex Number Calculator ---

	// Initialize complex numbers
	z1 := NewComplex(3, 4)
	z2 := NewComplex(1, -2)
	z3 := NewComplex(0, 0) // Zero complex number

	fmt.Println("--- Initial Complex Numbers ---")
	fmt.Println("z1 =", z1)
	fmt.Println("z2 =", z2)
	fmt.Println("z3 =", z3)
	fmt.Println("")

	// --- Basic Arithmetic Operations ---
	fmt.Println("--- Arithmetic Operations ---")
	sum := z1.Add(z2)
	fmt.Println("z1 + z2 =", sum)

	diff := z1.Subtract(z2)
	fmt.Println("z1 - z2 =", diff)

	prod := z1.Multiply(z2)
	fmt.Println("z1 * z2 =", prod)

	quot := z1.Divide(z2)
	fmt.Println("z1 / z2 =", quot)

	// Division by zero cases
	divByZeroNum := NewComplex(5, 0)
	divByZeroDen := NewComplex(0, 0)
	fmt.Println("5 / 0 =", divByZeroNum.Divide(divByZeroDen)) // Should result in Inf + Inf i
	fmt.Println("0 / 0 =", z3.Divide(divByZeroDen))           // Should result in NaN + NaN i
	fmt.Println("")

	// --- Complex Number Properties ---
	fmt.Println("--- Complex Number Properties ---")
	fmt.Println("Conjugate of z1 =", z1.Conjugate())
	fmt.Println("Magnitude of z1 =", z1.Magnitude())
	fmt.Println("Phase of z1 (radians) =", z1.Phase())
	fmt.Println("Phase of z1 (degrees) =", z1.Phase()*180/math.Pi)
	fmt.Println("")

	// --- Conversion and Comparison ---
	fmt.Println("--- Conversion and Comparison ---")
	// Reconstruct z1 from its polar coordinates
	magZ1 := z1.Magnitude()
	phaseZ1 := z1.Phase()
	z1FromPolar := FromPolar(magZ1, phaseZ1)
	fmt.Println("z1 from polar (mag=%.4f, phase=%.4f rad) = %s", magZ1, phaseZ1, z1FromPolar)
	// Check for approximate equality due to floating point precision
	fmt.Println("z1 equals z1FromPolar (tolerance 1e-9)?", z1.Equals(z1FromPolar, 1e-9))
	fmt.Println("")

	// --- Advanced Operation: Integer Power ---
	fmt.Println("--- Advanced Operation: Integer Power ---")
	z4 := NewComplex(1, 1) // 1 + i
	fmt.Println("z4 =", z4)
	fmt.Println("z4^0 =", z4.PowerInt(0))   // Should be 1 + 0i
	fmt.Println("z4^1 =", z4.PowerInt(1))   // Should be 1 + 1i
	fmt.Println("z4^2 =", z4.PowerInt(2))   // (1+i)^2 = 1 + 2i - 1 = 2i
	fmt.Println("z4^3 =", z4.PowerInt(3))   // (1+i)^3 = 2i(1+i) = 2i - 2 = -2 + 2i
	fmt.Println("z4^4 =", z4.PowerInt(4))   // (1+i)^4 = (-2+2i)(1+i) = -2 - 2i + 2i - 2 = -4
	fmt.Println("z4^-1 =", z4.PowerInt(-1)) // 1/(1+i) = (1-i)/2 = 0.5 - 0.5i

	z5 := NewComplex(0, 2) // 2i
	fmt.Println("z5 =", z5)
	fmt.Println("z5^2 =", z5.PowerInt(2)) // (2i)^2 = 4i^2 = -4
	fmt.Println("z5^3 =", z5.PowerInt(3)) // (2i)^3 = -4 * 2i = -8i
	fmt.Println("")
}

// Additional implementation at 2025-06-22 23:31:06
