package main

import (
	"fmt"
	"math"
	"sort"
)

// Polynomial represents a polynomial equation.
// The coefficients are stored in ascending order of power.
// e.g., for ax^2 + bx + c, coeffs would be [c, b, a]
type Polynomial []float64

// Degree returns the degree of the polynomial.
func (p Polynomial) Degree() int {
	for i := len(p) - 1; i >= 0; i-- {
		if math.Abs(p[i]) > 1e-9 { // Check for non-zero coefficient
			return i
		}
	}
	return 0 // All coefficients are zero, consider it degree 0
}

// Evaluate calculates the value of the polynomial at a given x.
func (p Polynomial) Evaluate(x float64) float64 {
	result := 0.0
	for i, coeff := range p {
		result += coeff * math.Pow(x, float64(i))
	}
	return result
}

// Derivative returns the derivative of the polynomial.
func (p Polynomial) Derivative() Polynomial {
	degree := p.Degree()
	if degree == 0 {
		return Polynomial{0.0} // Derivative of a constant is 0
	}
	deriv := make(Polynomial, degree)
	for i := 1; i <= degree; i++ {
		deriv[i-1] = p[i] * float64(i)
	}
	return deriv
}

// Solve finds the real roots of the polynomial.
// It dispatches to specific solvers based on the degree.
// For higher degrees, it uses a numerical method (Newton-Raphson)
// with multiple initial guesses to find multiple real roots.
func (p Polynomial) Solve() []float64 {
	degree := p.Degree()

	// Handle special cases for zero polynomial
	if degree == 0 && math.Abs(p[0]) < 1e-9 { // P(x) = 0
		// Infinite solutions, but we can't return infinite numbers.
		// For practical purposes, return an empty slice as no specific 'x' is a root.
		return []float64{}
	}
	if degree == 0 && math.Abs(p[0]) > 1e-9 { // P(x) = C (C != 0)
		return []float64{} // No solutions
	}

	switch degree {
	case 1:
		return p.solveLinear()
	case 2:
		return p.solveQuadratic()
	default:
		// For degree >= 3, use a numerical method like Newton-Raphson.
		// We'll try multiple initial guesses to find multiple roots.
		return p.solveNumerical()
	}
}

// solveLinear solves a linear equation ax + b = 0.
// Assumes degree is 1.
func (p Polynomial) solveLinear() []float64 {
	// p[0] is b, p[1] is a
	a := p[1]
	b := p[0]

	if math.Abs(a) < 1e-9 { // Should not happen if degree is truly 1, but for safety
		if math.Abs(b) < 1e-9 {
			return []float64{} // 0 = 0, infinite solutions, but no specific root
		}
		return []float64{} // b = 0 (b!=0), no solution
	}
	return []float64{-b / a}
}

// solveQuadratic solves a quadratic equation ax^2 + bx + c = 0.
// Assumes degree is 2.
func (p Polynomial) solveQuadratic() []float64 {
	// p[0] is c, p[1] is b, p[2] is a
	c := p[0]
	b := p[1]
	a := p[2]

	discriminant := b*b - 4*a*c

	if discriminant < -1e-9 { // Negative discriminant (complex roots)
		return []float64{} // Return no real roots
	} else if discriminant < 1e-9 { // Discriminant near zero (one real root)
		return []float64{-b / (2 * a)}
	} else { // Positive discriminant (two distinct real roots)
		sqrtDisc := math.Sqrt(discriminant)
		x1 := (-b + sqrtDisc) / (2 * a)
		x2 := (-b - sqrtDisc) / (2 * a)
		return []float64{x1, x2}
	}
}

// solveNumerical uses Newton-Raphson to find real roots.
// It tries multiple initial guesses to find distinct roots.
func (p Polynomial) solveNumerical() []float64 {
	roots := make(map[float64]struct{}) // Use a map to store unique roots
	deriv := p.Derivative()

	// Define a range for initial guesses.
	// This is a heuristic. A more robust approach might involve
	// analyzing the polynomial's behavior (e.g., Sturm's theorem).
	minGuess := -10.0
	maxGuess := 10.0
	numGuesses := 100 // Number of initial guesses

	for i := 0; i < numGuesses; i++ {
		initialGuess := minGuess + (maxGuess-minGuess)*float64(i)/float64(numGuesses-1)
		root := newtonRaphson(p, deriv, initialGuess, 1e-7, 1000) // tolerance, maxIterations

		if !math.IsNaN(root) { // Check if a root was found
			// Add root to map, rounding to avoid floating point precision issues
			roundedRoot := math.Round(root*1e6) / 1e6 // Round to 6 decimal places
			roots[roundedRoot] = struct{}{}
		}
	}

	// Convert map keys to a slice
	var result []float64
	for r := range roots {
		result = append(result, r)
	}
	sort.Float64s(result) // Sort roots for consistent output
	return result
}

// newtonRaphson performs the Newton-Raphson method to find a single root.
func newtonRaphson(p, pPrime Polynomial, initialGuess, tolerance float64, maxIterations int) float64 {
	x := initialGuess
	for i := 0; i < maxIterations; i++ {
		fx := p.Evaluate(x)
		fPrimeX := pPrime.Evaluate(x)

		if math.Abs(fx) < tolerance { // Root found
			return x
		}
		if math.Abs(fPrimeX) < 1e-9 { // Derivative is zero, method fails or hits a local extremum
			return math.NaN() // Indicate failure to converge or a problematic point
		}
		x = x - fx/fPrimeX
	}
	return math.NaN() // Did not converge within maxIterations
}

func main() {
	// Test cases
	// P(x) = 0 (constant zero polynomial)
	p0 := Polynomial{0}
	fmt.Printf("Solving %v: Roots = %v\n", p0, p0.Solve())

	// P(x) = 5 (constant non-zero polynomial)
	p1 := Polynomial{5}
	fmt.Printf("Solving %v: Roots = %v\n", p1, p1.Solve())

	// P(x) = 2x - 4 = 0  => x = 2
	p2 := Polynomial{-4, 2}
	fmt.Printf("Solving %v: Roots = %v\n", p2, p2.Solve())

	// P(x) = x^2 - 4 = 0 => x = 2, x = -2
	p3 := Polynomial{-4, 0, 1}
	fmt.Printf("Solving %v: Roots = %v\n", p3, p3.Solve())

	// P(x) = x^2 + 2x + 1 = 0 => x = -1 (double root)
	p4 := Polynomial{1, 2, 1}
	fmt.Printf("Solving %v: Roots = %v\n", p4, p4.Solve())

	// P(x) = x^2 + 1 = 0 (no real roots)
	p5 := Polynomial{1, 0, 1}
	fmt.Printf("Solving %v: Roots = %v\n", p5, p5.Solve())

	// P(x) = x^3 - 6x^2 + 11x - 6 = 0 => (x-1)(x-2)(x-3) = 0 => x = 1, 2, 3
	p6 := Polynomial{-6, 11, -6, 1}
	fmt.Printf("Solving %v: Roots = %v\n", p6, p6.Solve())

	// P(x) = x^4 - 10x^2 + 9 = 0 => (x^2-1)(x^2-9) = 0 => x = 1, -1, 3, -3
	p7 := Polynomial{9, 0, -10, 0, 1}
	fmt.Printf("Solving %v: Roots = %v\n", p7, p7.Solve())

	// P(x) = x^3 - x - 1 = 0 (one real root, approx 1.3247)
	p8 := Polynomial{-1, -1, 0, 1}
	fmt.Printf("Solving %v: Roots = %v\n", p8, p8.Solve())

	// P(x) = x^5 - 2x^4 - 4x^3 + 8x^2 + 3x - 6 = 0
	// Roots: -sqrt(3), -1, 1, sqrt(3), 2 (approx -1.732, -1, 1, 1.732, 2)
	p9 := Polynomial{-6, 3, 8, -4, -2, 1}
	fmt.Printf("Solving %v: Roots = %v\n", p9, p9.Solve())

	// P(x) = x^2 + 0x + 0 = 0, should be degree 2, but effectively x^2=0
	p10 := Polynomial{0, 0, 1}
	fmt.Printf("Solving %v: Roots = %v\n", p10, p10.Solve())

	// P(x) = 0x + 0 = 0, should be degree 0
	p11 := Polynomial{0, 0}
	fmt.Printf("Solving %v: Roots = %v\n", p11, p11.Solve())

	// P(x) = 0x + 5 = 0, should be degree 0
	p12 := Polynomial{5, 0}
	fmt.Printf("Solving %v: Roots = %v\n", p12, p12.Solve())
}

// Additional implementation at 2025-06-20 01:35:21
package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// Polynomial represents a polynomial equation.
// Coefficients are stored such that coeffs[i] is the coefficient of x^i.
// Example: 3x^2 + 2x - 1 would be []float64{-1, 2, 3}
type Polynomial []float64

// Degree returns the degree of the polynomial.
func (p Polynomial) Degree() int {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] != 0 {
			return i
		}
	}
	return 0 // Constant polynomial (e.g., P(x) = 0 or P(x) = 5)
}

// Evaluate calculates the value of the polynomial at a given x.
func (p Polynomial) Evaluate(x float64) float64 {
	result := 0.0
	for i, coeff := range p {
		result += coeff * math.Pow(x, float64(i))
	}
	return result
}

// Derivative returns the derivative of the polynomial.
func (p Polynomial) Derivative() Polynomial {
	if len(p) <= 1 {
		return Polynomial{0} // Derivative of constant is 0
	}
	deriv := make(Polynomial, len(p)-1)
	for i := 1; i < len(p); i++ {
		deriv[i-1] = p[i] * float64(i)
	}
	return deriv
}

// Solve finds the real roots of the polynomial.
// It handles linear and quadratic equations directly.
// For higher degrees, it uses a numerical method (Newton-Raphson) to find one or more real roots.
// It returns a slice of roots.
func (p Polynomial) Solve() []float64 {
	degree := p.Degree()

	// Handle special cases for degree 0
	if degree == 0 {
		if len(p) > 0 && p[0] == 0 {
			// P(x) = 0, infinite solutions (any x is a root).
			// We return nil to indicate no specific roots, handled by caller.
			return nil
		}
		return nil // P(x) = C (non-zero constant), no roots
	}

	if degree == 1 {
		// Linear equation: ax + b = 0 => x = -b/a
		// p[0] is b, p[1] is a
		if p[1] == 0 { // Should not happen if degree is 1, but as a safeguard
			return nil
		}
		return []float64{-p[0] / p[1]}
	}

	if degree == 2 {
		// Quadratic equation: ax^2 + bx + c = 0
		// p[0] is c, p[1] is b, p[2] is a
		a, b, c := p[2], p[1], p[0]
		discriminant := b*b - 4*a*c
		if discriminant < 0 {
			return nil // No real roots
		} else if discriminant == 0 {
			return []float64{-b / (2 * a)} // One real root (double root)
		} else {
			sqrtDiscriminant := math.Sqrt(discriminant)
			x1 := (-b + sqrtDiscriminant) / (2 * a)
			x2 := (-b - sqrtDiscriminant) / (2 * a)
			return []float64{x1, x2} // Two distinct real roots
		}
	}

	// For higher degrees, use Newton-Raphson to find real roots.
	// This method might not find all roots and requires good initial guesses.
	// We try a few initial guesses to increase chances of finding multiple roots.
	var roots []float64
	// Initial guesses cover a range to find different roots if they exist.
	initialGuesses := []float64{-10.0, -5.0, -2.0, -1.0, 0.0, 1.0, 2.0, 5.0, 10.0}
	for _, guess := range initialGuesses {
		root := newtonRaphson(p, guess, 1e-9, 100) // Tolerance 1e-9, max 100 iterations
		if !math.IsNaN(root) {
			// Check if this root is already found (due to multiple guesses converging to same root)
			found := false
			for _, r := range roots {
				if math.Abs(r-root) < 1e-7 { // Check for approximate equality
					found = true
					break
				}
			}
			if !found {
				roots = append(roots, root)
			}
		}
	}

	// Sort roots for consistent output
	sort.Float64s(roots)
	return roots
}

// newtonRaphson implements the Newton-Raphson method to find a root of a polynomial.
// p: the polynomial
// initialGuess: starting point for the iteration
// tolerance: desired accuracy for f(x) to be considered zero
// maxIterations: maximum number of iterations
func newtonRaphson(p Polynomial, initialGuess, tolerance float64, maxIterations int) float64 {
	x := initialGuess
	pPrime := p.Derivative()

	for i := 0; i < maxIterations; i++ {
		fx := p.Evaluate(x)
		fPrimeX := pPrime.Evaluate(x)

		if math.Abs(fx) < tolerance {
			return x // Found a root within tolerance
		}
		if fPrimeX == 0 {
			// Derivative is zero, Newton-Raphson fails.
			// This can happen at local extrema or inflection points.
			return math.NaN()
		}

		x = x - fx/fPrimeX
	}
	return math.NaN() // No convergence within maxIterations
}

// parseEquation parses a string equation like "2x^2 + 3x - 5 = 0" into a Polynomial.
// It assumes the equation is in the form P(x) = 0.
// This parser is simplified and may not handle all valid polynomial string formats (e.g., no parentheses).
func parseEquation(eq string) (Polynomial, error) {
	eq = strings.ReplaceAll(eq, " ", "")
	parts := strings.Split(eq, "=")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid equation format: must contain exactly one '='")
	}
	if parts[1] != "0" {
		return nil, fmt.Errorf("unsupported equation format: right side must be '0'")
	}

	termsStr := parts[0]
	// Replace "-" with "+-" to easily split by "+"
	termsStr = strings.ReplaceAll(termsStr, "-", "+-")
	rawTerms := strings.Split(termsStr, "+")

	coeffs := make(map[int]float64) // Map degree to coefficient
	maxDegree := 0

	for _, term := range rawTerms {
		if term == "" {
			continue
		}

		coeff := 1.0
		degree := 0
		isNegative := false

		if strings.HasPrefix(term, "-") {
			isNegative = true
			term = term[1:]
		}

		partsX := strings.Split(term, "x")
		if len(partsX) == 1 { // Constant term or just 'x' (if 'x' is not split)
			if term == "x" {
				coeff = 1.0
				degree = 1
			} else if term == "" {
				continue
			} else {
				val, err := strconv.ParseFloat(term, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid number in term '%s': %w", term, err)
				}
				coeff = val
				degree = 0
			}
		} else if len(partsX) == 2 { // Term with 'x'
			coeffStr := partsX[0]
			degreeStr := partsX[1]

			if coeffStr == "" { // e.g., "x" or "-x"
				coeff = 1.0
			} else {
				val, err := strconv.ParseFloat(coeffStr, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid coefficient in term '%s': %w", term, err)
				}
				coeff = val
			}

			if degreeStr == "" { // e.g., "3x"
				degree = 1
			} else if strings.HasPrefix(degreeStr, "^") { // e.g., "x^2"
				degVal, err := strconv.Atoi(degreeStr[1:])
				if err != nil {
					return nil, fmt.Errorf("invalid degree in term '%s': %w", term, err)
				}
				degree = degVal
			} else {
				return nil, fmt.Errorf("invalid term format '%s'", term)
			}
		} else {
			return nil, fmt.Errorf("invalid term format '%s'", term)
		}

		if isNegative {
			coeff = -coeff
		}

		coeffs[degree] += coeff
		if degree > maxDegree {
			maxDegree = degree
		}
	}

	// Convert map to slice
	polySlice := make(Polynomial, maxDegree+1)
	for deg, val := range coeffs {
		if deg >= len(polySlice) {
			// This case should ideally not be hit if maxDegree is tracked correctly,
			// but serves as a safeguard for unexpected parsing results.
			newPolySlice := make(Polynomial, deg+1)
			copy(newPolySlice, polySlice)
			polySlice = newPolySlice
		}
		polySlice[deg] = val
	}

	return polySlice, nil
}

func main() {
	equations := []string{
		"5 = 0",                  // Degree 0, no solution
		"2x + 4 = 0",             // Linear, x = -2
		"3x - 9 = 0",             // Linear, x = 3
		"x^2 - 4 = 0",            // Quadratic, x = 2, x = -2
		"x^2 + 2x + 1 = 0",       // Quadratic, x = -1 (double root)
		"x^2 + 1 = 0",            // Quadratic, no real roots
		"2x^2 + 3x - 5 = 0",      // Quadratic, x = 1, x = -2.5
		"x^3 - 6x^2 + 11x - 6 = 0", // Cubic, roots: 1, 2, 3
		"x^4 - 10x^2 + 9 = 0",    // Quartic, roots: 1, -1, 3, -3
		"x^5 - 2x^4 - 4x^3 + 8x^2 + 3x - 6 = 0", // Quintic, roots approx: -1.732, -1, 1.414, 2, 1.732
		"x = 0",                  // Linear, x = 0
		"0 = 0",                  // Degree 0, P(x)=0, infinite solutions
		"-x^2 + 4x - 4 = 0",      // Quadratic, x = 2 (double root)
		"x^3 - x = 0",            // Cubic, roots: -1, 0, 1
	}

	for _, eqStr := range equations {
		fmt.Printf("Solving: %s\n", eqStr)
		poly, err := parseEquation(eqStr)
		if err != nil {
			fmt.Printf("  Error parsing equation: %v\n", err)
			continue
		}

		roots := poly.Solve()

		if poly.Degree() == 0 && len(poly) > 0 && poly[0] == 0 {
			fmt.Println("  Solution: All real numbers (0 = 0)")
		} else if len(roots) == 0 {
			fmt.Println("  No real roots found.")
		} else {
			fmt.Printf("  Roots: %v\n", roots)
		}
		fmt.Println("------------------------------------")
	}
}

// Additional implementation at 2025-06-20 01:36:27
package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

// Polynomial represents a polynomial equation.
// The slice stores coefficients where coeffs[i] is the coefficient of x^i.
// For example, {1, 2, 3} represents 3x^2 + 2x + 1.
type Polynomial []float64

// Degree returns the degree of the polynomial.
func (p Polynomial) Degree() int {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] != 0 {
			return i
		}
	}
	return 0 // Zero polynomial or constant 0
}

// Evaluate evaluates the polynomial at a given x value.
func (p Polynomial) Evaluate(x float64) float64 {
	result := 0.0
	for i, coeff := range p {
		result += coeff * math.Pow(x, float64(i))
	}
	return result
}

// Derivative computes the derivative of the polynomial.
func (p Polynomial) Derivative() Polynomial {
	degree := p.Degree()
	if degree == 0 {
		return Polynomial{0} // Derivative of a constant is 0
	}
	deriv := make(Polynomial, degree)
	for i := 1; i <= degree; i++ {
		deriv[i-1] = p[i] * float64(i)
	}
	return deriv
}

// Solve finds the roots of the polynomial.
// It handles linear and quadratic equations analytically.
// For higher degrees, it attempts to find one or more real roots using Newton-Raphson
// with multiple initial guesses. This method is not guaranteed to find all roots
// (especially complex ones for degrees > 2) and might only find real roots if they converge.
func (p Polynomial) Solve() []complex128 {
	degree := p.Degree()

	// Handle zero polynomial (all coefficients are zero)
	if degree == 0 && len(p) > 0 && p[0] == 0 {
		// All coefficients are zero, infinite solutions.
		// Returning an empty slice signifies no specific finite roots.
		return []complex128{}
	}

	// Handle constant non-zero polynomial (e.g., P(x) = 5)
	if degree == 0 {
		return []complex128{} // No roots for P(x) = C where C != 0
	}

	// Linear equation: ax + b = 0
	if degree == 1 {
		b := p[0] // constant term
		a := p[1] // coefficient of x
		if a == 0 {
			return []complex128{} // Should not happen if degree is truly 1
		}
		root := -b / a
		return []complex128{complex(root, 0)}
	}

	// Quadratic equation: ax^2 + bx + c = 0
	if degree == 2 {
		c := p[0] // constant term
		b := p[1] // coefficient of x
		a := p[2] // coefficient of x^2
		if a == 0 {
			// Degenerates to linear or constant
			return Polynomial{b, c}.Solve()
		}

		discriminant := b*b - 4*a*c
		if discriminant >= 0 {
			// Real roots
			root1 := (-b + math.Sqrt(discriminant)) / (2 * a)
			root2 := (-b - math.Sqrt(discriminant)) / (2 * a)
			return []complex128{complex(root1, 0), complex(root2, 0)}
		} else {
			// Complex roots
			realPart := -b / (2 * a)
			imagPart := math.Sqrt(math.Abs(discriminant)) / (2 * a)
			root1 := complex(realPart, imagPart)
			root2 := complex(realPart, -imagPart)
			return []complex128{root1, root2}
		}
	}

	// For higher degrees (degree > 2), use Newton-Raphson to find real roots.
	maxIterations := 1000
	tolerance := 1e-9

	deriv := p.Derivative()
	if deriv.Degree() == 0 && deriv[0] == 0 {
		// Derivative is zero polynomial, implies original was constant.
		// This case should ideally not be reached if degree > 2.
		return []complex128{}
	}

	initialGuesses := []float64{0.0, 1.0, -1.0, 0.5, -0.5, 10.0, -10.0}
	var foundRoots []complex128

	for _, guess := range initialGuesses {
		x0 := guess
		for i := 0; i < maxIterations; i++ {
			fx0 := p.Evaluate(x0)
			fPrimeX0 := deriv.Evaluate(x0)

			if math.Abs(fx0) < tolerance {
				// Found a root, check if it's new
				isNewRoot := true
				for _, r := range foundRoots {
					if cmplx.Abs(r - complex(x0, 0)) < tolerance {
						isNewRoot = false
						break
					}
				}
				if isNewRoot {
					foundRoots = append(foundRoots, complex(x0, 0))
				}
				break // Converged for this guess
			}

			if math.Abs(fPrimeX0) < tolerance {
				// Derivative is too close to zero, Newton-Raphson might diverge.
				break
			}

			x0 = x0 - fx0/fPrimeX0
		}
	}

	return foundRoots
}

func main() {
	p1 := Polynomial{4, 2} // 2x + 4 = 0
	fmt.Printf("Polynomial: 2x + 4\nRoots: %v\n", p1.Solve())

	p2 := Polynomial{-4, 0, 1} // x^2 - 4 = 0
	fmt.Printf("\nPolynomial: x^2 - 4\nRoots: %v\n", p2.Solve())

	p3 := Polynomial{1, 0, 1} // x^2 + 1 = 0
	fmt.Printf("\nPolynomial: x^2 + 1\nRoots: %v\n", p3.Solve())

	p4 := Polynomial{-6, 11, -6, 1} // x^3 - 6x^2 + 11x - 6 = 0
	fmt.Printf("\nPolynomial: x^3 - 6x^2 + 11x - 6\nRoots (Newton-Raphson): %v\n", p4.Solve())

	p5 := Polynomial{5} // 5 = 0
	fmt.Printf("\nPolynomial: 5\nRoots: %v\n", p5.Solve())

	p6 := Polynomial{0} // 0 = 0
	fmt.Printf("\nPolynomial: 0\nRoots: %v\n", p6.Solve())

	p7 := Polynomial{1, -3, 0, 0, 0, 1} // x^5 - 3x + 1 = 0
	fmt.Printf("\nPolynomial: x^5 - 3x + 1\nRoots (Newton-Raphson): %v\n", p7.Solve())

	p8 := Polynomial{0, 0, 0, 1, 2} // 2x^4 + x^3 = 0
	fmt.Printf("\nPolynomial: 2x^4 + x^3\nRoots (Newton-Raphson): %v\n", p8.Solve())
}

// Additional implementation at 2025-06-20 01:37:38
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// evaluatePolynomial evaluates a polynomial at a given x.
// coeffs are ordered from highest degree to constant term (e.g., [a, b, c] for ax^2 + bx + c).
func evaluatePolynomial(coeffs []float64, x float64) float64 {
	result := 0.0
	degree := len(coeffs) - 1
	for i, coeff := range coeffs {
		result += coeff * math.Pow(x, float64(degree-i))
	}
	return result
}

// derivePolynomial returns the coefficients of the derivative of a polynomial.
// coeffs are ordered from highest degree to constant term.
func derivePolynomial(coeffs []float64) []float64 {
	if len(coeffs) <= 1 { // Constant or zero polynomial
		return []float64{0.0}
	}
	derivedCoeffs := make([]float64, len(coeffs)-1)
	degree := len(coeffs) - 1
	for i, coeff := range coeffs[:degree] { // Iterate up to the second to last coefficient
		derivedCoeffs[i] = coeff * float64(degree-i)
	}
	return derivedCoeffs
}

// newtonRaphson finds a root of a polynomial using the Newton-Raphson method.
// It returns one real root if found within maxIterations and tolerance.
func newtonRaphson(coeffs []float64, initialGuess float64, maxIterations int, tolerance float64) (float64, bool) {
	x := initialGuess
	derivedCoeffs := derivePolynomial(coeffs)

	for i := 0; i < maxIterations; i++ {
		fx := evaluatePolynomial(coeffs, x)
		fPrimeX := evaluatePolynomial(derivedCoeffs, x)

		if math.Abs(fx) < tolerance {
			return x, true // Found a root
		}
		if math.Abs(fPrimeX) < 1e-10 { // Avoid division by zero or very small derivative
			return x, false // Derivative is too small, likely a local extremum or multiple root
		}

		x = x - fx/fPrimeX
	}
	return x, false // Did not converge
}

// solvePolynomial solves a polynomial equation based on its degree.
// coeffs are ordered from highest degree to constant term.
func solvePolynomial(coeffs []float64) {
	// Determine the true degree by removing leading zero coefficients
	trueCoeffs := []float64{}
	foundNonZero := false
	for _, c := range coeffs {
		if c != 0 {
			foundNonZero = true
		}
		if foundNonZero {
			trueCoeffs = append(trueCoeffs, c)
		}
	}

	if len(trueCoeffs) == 0 { // All coefficients were zero
		fmt.Println("Equation is 0 = 0. Infinitely many solutions.")
		return
	}

	degree := len(trueCoeffs) - 1

	switch degree {
	case 0: // c = 0
		if trueCoeffs[0] == 0 { // This case should be covered by len(trueCoeffs) == 0
			fmt.Println("Equation is 0 = 0. Infinitely many solutions.")
		} else {
			fmt.Println("No solution (constant non-zero).")
		}
	case 1: // ax + b = 0
		a := trueCoeffs[0]
		b := trueCoeffs[1]
		// 'a' cannot be zero here due to trueCoeffs logic
		x := -b / a
		fmt.Printf("Solution: x = %.6f\n", x)
	case 2: // ax^2 + bx + c = 0
		a := trueCoeffs[0]
		b := trueCoeffs[1]
		c := trueCoeffs[2]

		delta := b*b - 4*a*c
		if delta >= 0 {
			x1 := (-b + math.Sqrt(delta)) / (2 * a)
			x2 := (-b - math.Sqrt(delta)) / (2 * a)
			if delta == 0 {
				fmt.Printf("Solution: x = %.6f (double root)\n", x1)
			} else {
				fmt.Printf("Solutions: x1 = %.6f, x2 = %.6f\n", x1, x2)
			}
		} else {
			realPart := -b / (2 * a)
			imagPart := math.Sqrt(math.Abs(delta)) / (2 * a)
			fmt.Printf("Solutions: x1 = %.6f + %.6fi, x2 = %.6f - %.6fi\n", realPart, imagPart, realPart, imagPart)
		}
	default: // Degree 3 or higher, use numerical method for one real root
		fmt.Printf("Solving degree %d polynomial numerically (finding one real root)...\n", degree)
		// Try a few initial guesses to find a root
		initialGuesses := []float64{0.0, 1.0, -1.0, 0.5, -0.5, 10.0, -10.0, 100.0, -100.0}
		found := false
		for _, guess := range initialGuesses {
			root, converged := newtonRaphson(trueCoeffs, guess, 100, 1e-7) // Increased iterations, tightened tolerance
			if converged {
				fmt.Printf("Found a real root: x = %.6f (using initial guess %.1f)\n", root, guess)
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Could not find a real root using Newton-Raphson with common initial guesses.")
			fmt.Println("Consider trying different initial guesses or a different numerical method for all roots.")
		}
	}
}

// getCoefficients prompts the user for coefficients and parses them.
// It returns the coefficients, a boolean indicating if the user wants to exit, and an error.
func getCoefficients() ([]float64, bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter coefficients (e.g., '1 2 3' for x^2 + 2x + 3) or 'exit' to quit: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if strings.ToLower(input) == "exit" {
		return nil, true, nil
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil, false, fmt.Errorf("no coefficients entered")
	}

	coeffs := make([]float64, len(parts))
	for i, part := range parts {
		val, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil, false, fmt.Errorf("invalid coefficient '%s': %w", part, err)
		}
		coeffs[i] = val
	}
	return coeffs, false, nil
}

func main() {
	fmt.Println("Polynomial Equation Solver")
	fmt.Println("This solver finds exact solutions for linear and quadratic equations.")
	fmt.Println("For higher degree polynomials, it attempts to find one real root using Newton-Raphson.")
	fmt.Println("Coefficients should be entered from highest degree to constant term.")
	fmt.Println("Example: For x^2 + 2x + 3 = 0, enter '1 2 3'")
	fmt.Println("Type 'exit' to quit.")
	fmt.Println("--------------------------------------------------")

	for {
		coeffs, isExit, err := getCoefficients()
		if isExit {
			fmt.Println("Exiting solver.")
			break
		}
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		solvePolynomial(coeffs)
		fmt.Println("--------------------------------------------------")
	}
}