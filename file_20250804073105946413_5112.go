package main

import (
	"errors"
	"fmt"
	"math"
)

// Polynomial represents a polynomial as a slice of coefficients.
// The index i corresponds to the coefficient of x^i.
// For example, [c0, c1, c2] represents c2*x^2 + c1*x + c0.
type Polynomial []float64

// Evaluate calculates the value of the polynomial at a given x.
func (p Polynomial) Evaluate(x float64) float64 {
	result := 0.0
	for i, coeff := range p {
		result += coeff * math.Pow(x, float64(i))
	}
	return result
}

// Derivative calculates the derivative of the polynomial.
func (p Polynomial) Derivative() Polynomial {
	if len(p) <= 1 {
		return Polynomial{0.0} // Derivative of a constant or empty polynomial is 0
	}
	deriv := make(Polynomial, len(p)-1)
	for i := 1; i < len(p); i++ {
		deriv[i-1] = float64(i) * p[i]
	}
	return deriv
}

// SolveLinear solves a linear equation (ax + b = 0).
// The polynomial is expected to be of degree 0 or 1.
// p[0] is b, p[1] is a.
func SolveLinear(p Polynomial) ([]float64, error) {
	if len(p) == 0 {
		return nil, errors.New("empty polynomial, no equation")
	}
	if len(p) == 1 { // Constant equation: c = 0
		if p[0] == 0 {
			return nil, errors.New("infinite solutions (0=0)")
		}
		return nil, errors.New("no solution (constant != 0)")
	}
	// Now len(p) >= 2. p[0] is b, p[1] is a.
	if p[1] == 0 { // a = 0
		if p[0] == 0 { // 0x + 0 = 0
			return nil, errors.New("infinite solutions (0x+0=0)")
		}
		return nil, errors.New("no solution (0x+b=0, b!=0)")
	}
	return []float64{-p[0] / p[1]}, nil
}

// SolveQuadratic solves a quadratic equation (ax^2 + bx + c = 0).
// The polynomial is expected to be of degree 0, 1, or 2.
// p[0] is c, p[1] is b, p[2] is a.
func SolveQuadratic(p Polynomial) ([]float64, error) {
	if len(p) < 3 { // Not a quadratic equation (degree < 2)
		return SolveLinear(p) // Delegate to linear solver (or constant)
	}

	a := p[2]
	b := p[1]
	c := p[0]

	if a == 0 { // It's actually a linear equation (or constant)
		return SolveLinear(p[:2]) // Pass [c, b] to linear solver
	}

	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return nil, errors.New("no real roots")
	} else if discriminant == 0 {
		return []float64{-b / (2 * a)}, nil
	} else {
		sqrtDiscriminant := math.Sqrt(discriminant)
		x1 := (-b + sqrtDiscriminant) / (2 * a)
		x2 := (-b - sqrtDiscriminant) / (2 * a)
		return []float64{x1, x2}, nil
	}
}

// SolveNewtonRaphson finds a root of a polynomial using the Newton-Raphson method.
// It requires an initial guess, a tolerance for convergence, and a maximum number of iterations.
// This method finds only one root and requires the derivative of the polynomial.
func SolveNewtonRaphson(p Polynomial, initialGuess float64, tolerance float64, maxIterations int) (float64, error) {
	if len(p) == 0 {
		return 0, errors.New("empty polynomial")
	}
	if len(p) == 1 { // Constant polynomial
		if p[0] == 0 {
			return 0, errors.New("infinite solutions (0=0)")
		}
		return 0, errors.New("no root (constant != 0)")
	}

	deriv := p.Derivative()
	// Check if derivative is effectively zero everywhere (e.g., for a constant polynomial, already handled)
	// Or if the derivative at the initial guess is zero, which would cause division by zero.
	if len(deriv) == 0 || math.Abs(deriv.Evaluate(initialGuess)) < 1e-10 {
		return 0, errors.New("derivative is zero or too close to zero, Newton-Raphson cannot proceed")
	}

	x := initialGuess
	for i := 0; i < maxIterations; i++ {
		fx := p.Evaluate(x)
		fPrimeX := deriv.Evaluate(x)

		if math.Abs(fPrimeX) < 1e-10 { // Avoid division by zero or very small derivative
			return 0, errors.New("derivative too close to zero, Newton-Raphson cannot proceed")
		}

		xNew := x - fx/fPrimeX

		if math.Abs(xNew-x) < tolerance {
			return xNew, nil
		}
		x = xNew
	}
	return 0, errors.New("Newton-Raphson did not converge within max iterations")
}

func main() {
	fmt.Println("Polynomial Equation Solver in Go")

	// Example 1: Linear Equation (2x + 4 = 0)
	// Coefficients: c0=4, c1=2 => [4, 2]
	pLinear := Polynomial{4, 2}
	fmt.Printf("\nSolving linear equation: %.2fx + %.2f = 0\n", pLinear[1], pLinear[0])
	roots, err := SolveLinear(pLinear)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Root(s): %v\n", roots)
	}

	// Example 2: Quadratic Equation (x^2 - 5x + 6 = 0)
	// Coefficients: c0=6, c1=-5, c2=1 => [6, -5, 1]
	pQuadratic := Polynomial{6, -5, 1}
	fmt.Printf("\nSolving quadratic equation: %.2fx^2 + %.2fx + %.2f = 0\n", pQuadratic[2], pQuadratic[1], pQuadratic[0])
	roots, err = SolveQuadratic(pQuadratic)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Root(s): %v\n", roots)
	}

	// Example 3: Quadratic Equation (x^2 + 4 = 0) - No real roots
	// Coefficients: c0=4, c1=0, c2=1 => [4, 0, 1]
	pNoRealRoots := Polynomial{4, 0, 1}
	fmt.Printf("\nSolving quadratic equation: %.2fx^2 + %.2fx + %.2f = 0\n", pNoRealRoots[2], pNoRealRoots[1], pNoRealRoots[0])
	roots, err = SolveQuadratic(pNoRealRoots)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Root(s): %v\n", roots)
	}

	// Example 4: Quadratic Equation (x^2 - 4x + 4 = 0) - One real root
	// Coefficients: c0=4, c1=-4, c2=1 => [4, -4, 1]
	pOneRealRoot := Polynomial{4, -4, 1}
	fmt.Printf("\nSolving quadratic equation: %.2fx^2 + %.2fx + %.2f = 0\n", pOneRealRoot[2], pOneRealRoot[1], pOneRealRoot[0])
	roots, err = SolveQuadratic(pOneRealRoot)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Root(s): %v\n", roots)
	}

	// Example 5: Higher Degree Polynomial (x^3 - 6x^2 + 11x - 6 = 0) using Newton-Raphson
	// Coefficients: c0=-6, c1=11, c2=-6, c3=1 => [-6, 11, -6, 1]
	// Known roots: 1, 2, 3
	pCubic := Polynomial{-6, 11, -6, 1}
	fmt.Printf("\nSolving cubic equation: %.2fx^3 + %.2fx^2 + %.2fx + %.2f = 0 using Newton-Raphson\n", pCubic[3], pCubic[2], pCubic[1], pCubic[0])

	// Try to find root near 0.5 (should converge to 1)
	root, err := SolveNewtonRaphson(pCubic, 0.5, 1e-6, 100)
	if err != nil {
		fmt.Printf("Newton-Raphson Error (initial guess 0.5): %v\n", err)
	} else {
		fmt.Printf("Newton-Raphson Root (initial guess 0.5): %f (P(root)=%f)\n", root, pCubic.Evaluate(root))
	}

	// Try to find root near 2.5 (should converge to 2 or 3 depending on path)
	root, err = SolveNewtonRaphson(pCubic, 2.5, 1e-6, 100)
	if err != nil {
		fmt.Printf("Newton-Raphson Error (initial guess 2.5): %v\n", err)
	} else {
		fmt.Printf("Newton-Raphson Root (initial guess 2.5): %f (P(root)=%f)\n", root, pCubic.Evaluate(root))
	}

	// Example 6: Constant polynomial (5 = 0)
	pConstant := Polynomial{5}
	fmt.Printf("\nSolving constant equation: %.2f = 0\n", pConstant[0])
	_, err = SolveNewtonRaphson(pConstant, 0, 1e-6, 100)
	if err != nil {
		fmt.Printf("Newton-Raphson Error: %v\n", err)
	}

	// Example 7: Zero polynomial (0 = 0)
	pZero := Polynomial{0}
	fmt.Printf("\nSolving zero equation: %.2f = 0\n", pZero[0])
	_, err = SolveNewtonRaphson(pZero, 0, 1e-6, 100)
	if err != nil {
		fmt.Printf("Newton-Raphson Error: %v\n", err)
	}

	// Example 8: Linear equation with a=0 (0x + 5 = 0)
	pLinearZeroA := Polynomial{5, 0}
	fmt.Printf("\nSolving linear equation: %.2fx + %.2f = 0\n", pLinearZeroA[1], pLinearZeroA[0])
	roots, err = SolveLinear(pLinearZeroA)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Root(s): %v\n", roots)
	}

	// Example 9: Linear equation with a=0, b=0 (0x + 0 = 0)
	pLinearZeroAB := Polynomial{0, 0}
	fmt.Printf("\nSolving linear equation: %.2fx + %.2f = 0\n", pLinearZeroAB[1], pLinearZeroAB[0])
	roots, err = SolveLinear(pLinearZeroAB)
	if err != nil {

// Additional implementation at 2025-08-04 07:32:36
package main

import (
	"fmt"
	"math"
	"math/cmplx
	"sort"
	"strings"
)

// Polynomial represents a polynomial equation.
// Coefficients are stored such that coeffs[i] is the coefficient of x^i.
// For example, [1, 2, 3] represents 1 + 2x + 3x^2.
type Polynomial struct {
	Coeffs []float64
}

// NewPolynomial creates a new Polynomial from a slice of coefficients.
// The slice should be ordered from the constant term up to the highest degree.
func NewPolynomial(coeffs ...float64) *Polynomial {
	// Remove trailing zeros to get the true degree
	degree := len(coeffs) - 1
	for degree >= 0 && coeffs[degree] == 0 {
		degree--
	}
	if degree < 0 {
		return &Polynomial{Coeffs: []float64{0}} // Zero polynomial
	}
	return &Polynomial{Coeffs: coeffs[:degree+1]}
}

// Degree returns the degree of the polynomial.
func (p *Polynomial) Degree() int {
	return len(p.Coeffs) - 1
}

// Evaluate calculates the value of the polynomial at a given x.
func (p *Polynomial) Evaluate(x float64) float64 {
	result := 0.0
	for i, coeff := range p.Coeffs {
		result += coeff * math.Pow(x, float64(i))
	}
	return result
}

// Derivative returns a new Polynomial representing the derivative of the current polynomial.
func (p *Polynomial) Derivative() *Polynomial {
	if p.Degree() < 1 {
		return NewPolynomial(0) // Derivative of a constant is 0
	}
	derivCoeffs := make([]float64, p.Degree())
	for i := 1; i <= p.Degree(); i++ {
		derivCoeffs[i-1] = p.Coeffs[i] * float64(i)
	}
	return NewPolynomial(derivCoeffs...)
}

// SolveLinear solves a linear equation (ax + b = 0).
// Returns a slice of roots.
func (p *Polynomial) SolveLinear() []float64 {
	if p.Degree() != 1 {
		return nil // Not a linear equation
	}
	b := p.Coeffs[0]
	a := p.Coeffs[1]

	if a == 0 {
		if b == 0 {
			// 0x + 0 = 0, infinite solutions (or all real numbers)
			// For simplicity, return an empty slice to indicate no unique solution.
			return []float64{}
		}
		// 0x + b = 0 where b != 0, no solution
		return []float64{}
	}
	return []float64{-b / a}
}

// SolveQuadratic solves a quadratic equation (ax^2 + bx + c = 0).
// Returns a slice of complex128 roots.
func (p *Polynomial) SolveQuadratic() []complex128 {
	if p.Degree() != 2 {
		return nil // Not a quadratic equation
	}
	c := p.Coeffs[0]
	b := p.Coeffs[1]
	a := p.Coeffs[2]

	if a == 0 { // Degenerates to a linear equation
		linearPoly := NewPolynomial(c, b)
		realRoots := linearPoly.SolveLinear()
		complexRoots := make([]complex128, len(realRoots))
		for i, r := range realRoots {
			complexRoots[i] = complex(r, 0)
		}
		return complexRoots
	}

	discriminant := b*b - 4*a*c
	if discriminant >= 0 {
		sqrtDiscriminant := math.Sqrt(discriminant)
		x1 := (-b + sqrtDiscriminant) / (2 * a)
		x2 := (-b - sqrtDiscriminant) / (2 * a)
		return []complex128{complex(x1, 0), complex(x2, 0)}
	} else {
		sqrtDiscriminant := cmplx.Sqrt(complex(discriminant, 0))
		x1 := (-complex(b, 0) + sqrtDiscriminant) / (2 * complex(a, 0))
		x2 := (-complex(b, 0) - sqrtDiscriminant) / (2 * complex(a, 0))
		return []complex128{x1, x2}
	}
}

// SolveNewtonRaphson attempts to find a single real root using the Newton-Raphson method.
// It requires an initial guess, a tolerance for convergence, and a maximum number of iterations.
// Returns the found root and an error if convergence fails.
func (p *Polynomial) SolveNewtonRaphson(initialGuess float64, tolerance float64, maxIterations int) (float64, error) {
	if p.Degree() == 0 {
		if p.Coeffs[0] == 0 {
			return 0, fmt.Errorf("zero polynomial, infinite roots")
		}
		return 0, fmt.Errorf("constant non-zero polynomial, no roots")
	}

	x := initialGuess
	deriv := p.Derivative()

	for i := 0; i < maxIterations; i++ {
		fx := p.Evaluate(x)
		fPrimeX := deriv.Evaluate(x)

		if math.Abs(fx) < tolerance {
			return x, nil // Converged
		}

		if math.Abs(fPrimeX) < 1e-10 { // Avoid division by zero or very small derivative
			return x, fmt.Errorf("derivative too close to zero at x = %f, cannot proceed", x)
		}

		x = x - fx/fPrimeX
	}
	return x, fmt.Errorf("Newton-Raphson did not converge within %d iterations", maxIterations)
}

// String returns a string representation of the polynomial.
func (p *Polynomial) String() string {
	var terms []string
	// Iterate from highest degree down to constant term
	for i := p.Degree(); i >= 0; i-- {
		coeff := p.Coeffs[i]
		if coeff == 0 {
			continue // Skip zero coefficients
		}

		absCoeff := math.Abs(coeff)
		var term string

		// Determine the coefficient part of the term
		coeffStr := ""
		if absCoeff != 1 || i == 0 { // If coeff is not 1 (or -1), or it'