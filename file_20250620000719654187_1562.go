package main

import (
	"fmt"
	"math"
)

func solveLinearEquations(matrix [][]float64) ([]float64, error) {
	n := len(matrix)
	if n == 0 {
		return nil, fmt.Errorf("empty matrix")
	}
	m := len(matrix[0])
	if m != n+1 {
		return nil, fmt.Errorf("invalid augmented matrix dimensions: expected N rows and N+1 columns, got %d rows and %d columns", n, m)
	}

	// Gaussian Elimination (Forward Elimination)
	for k := 0; k < n; k++ {
		// Find pivot row
		iMax := k
		for i := k + 1; i < n; i++ {
			if math.Abs(matrix[i][k]) > math.Abs(matrix[iMax][k]) {
				iMax = i
			}
		}

		// Swap rows if necessary
		if iMax != k {
			matrix[k], matrix[iMax] = matrix[iMax], matrix[k]
		}

		// Check for singular matrix (or no unique solution)
		if math.Abs(matrix[k][k]) < 1e-9 { // Use a small epsilon for float comparison
			return nil, fmt.Errorf("matrix is singular or has no unique solution")
		}

		// Eliminate elements below the pivot
		for i := k + 1; i < n; i++ {
			factor := matrix[i][k] / matrix[k][k]
			for j := k; j < m; j++ {
				matrix[i][j] -= factor * matrix[k][j]
			}
		}
	}

	// Back Substitution
	solutions := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		sum := 0.0
		for j := i + 1; j < n; j++ {
			sum += matrix[i][j] * solutions[j]
		}
		if math.Abs(matrix[i][i]) < 1e-9 {
			return nil, fmt.Errorf("matrix is singular or has no unique solution during back substitution")
		}
		solutions[i] = (matrix[i][n] - sum) / matrix[i][i]
	}

	return solutions, nil
}

func main() {
	// Example 1: Unique solution
	// 2x + y - z = 8
	// -3x - y + 2z = -11
	// -2x + y + 2z = -3
	matrix1 := [][]float64{
		{2, 1, -1, 8},
		{-3, -1, 2, -11},
		{-2, 1, 2, -3},
	}
	fmt.Println("Solving System 1:")
	solutions1, err1 := solveLinearEquations(matrix1)
	if err1 != nil {
		fmt.Println("Error:", err1)
	} else {
		fmt.Println("Solutions:", solutions1)
	}
	fmt.Println()

	// Example 2: Another unique solution
	// x + 2y + 3z = 6
	// 2x + y - z = 1
	// 3x + y + 2z = 5
	matrix2 := [][]float64{
		{1, 2, 3, 6},
		{2, 1, -1, 1},
		{3, 1, 2, 5},
	}
	fmt.Println("Solving System 2:")
	solutions2, err2 := solveLinearEquations(matrix2)
	if err2 != nil {
		fmt.Println("Error:", err2)
	} else {
		fmt.Println("Solutions:", solutions2)
	}
	fmt.Println()

	// Example 3: Singular matrix (no unique solution, inconsistent system)
	// x + y = 2
	// x + y = 3
	matrix3 := [][]float64{
		{1, 1, 2},
		{1, 1, 3},
	}
	fmt.Println("Solving System 3 (Singular/No Solution):")
	solutions3, err3 := solveLinearEquations(matrix3)
	if err3 != nil {
		fmt.Println("Error:", err3)
	} else {
		fmt.Println("Solutions:", solutions3)
	}
	fmt.Println()

	// Example 4: Singular matrix (infinite solutions, dependent system)
	// x + y = 2
	// 2x + 2y = 4
	matrix4 := [][]float64{
		{1, 1, 2},
		{2, 2, 4},
	}
	fmt.Println("Solving System 4 (Singular/Infinite Solutions):")
	solutions4, err4 := solveLinearEquations(matrix4)
	if err4 != nil {
		fmt.Println("Error:", err4)
	} else {
		fmt.Println("Solutions:", solutions4)
	}
	fmt.Println()

	// Example 5: Empty matrix
	matrix5 := [][]float64{}
	fmt.Println("Solving System 5 (Empty Matrix):")
	solutions5, err5 := solveLinearEquations(matrix5)
	if err5 != nil {
		fmt.Println("Error:", err5)
	} else {
		fmt.Println("Solutions:", solutions5)
	}
	fmt.Println()

	// Example 6: Invalid dimensions (2 equations, 2 variables, but 3 columns)
	matrix6 := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
	}
	fmt.Println("Solving System 6 (Invalid Dimensions):")
	solutions6, err6 := solveLinearEquations(matrix6)
	if err6 != nil {
		fmt.Println("Error:", err6)
	} else {
		fmt.Println("Solutions:", solutions6)
	}
	fmt.Println()
}

// Additional implementation at 2025-06-20 00:08:31
package main

import (
	"fmt"
	"math"
)

const epsilon = 1e-9

func solveLinearEquations(matrix [][]float64) ([]float64, error) {
	n := len(matrix)
	if n == 0 {
		return nil, fmt.Errorf("empty matrix provided")
	}
	if len(matrix[0]) != n+1 {
		return nil, fmt.Errorf("invalid matrix dimensions: expected %dx%d, got %dx%d", n, n+1, n, len(matrix[0]))
	}

	augMatrix := make([][]float64, n)
	for i := range matrix {
		augMatrix[i] = make([]float64, n+1)
		copy(augMatrix[i], matrix[i])
	}

	for k := 0; k < n; k++ {
		iMax := k
		for i := k + 1; i < n; i++ {
			if math.Abs(augMatrix[i][k]) > math.Abs(augMatrix[iMax][k]) {
				iMax = i
			}
		}

		if math.Abs(augMatrix[iMax][k]) < epsilon {
			return nil, fmt.Errorf("matrix is singular or has no unique solution")
		}

		augMatrix[k], augMatrix[iMax] = augMatrix[iMax], augMatrix[k]

		for i := k + 1; i < n; i++ {
			factor := augMatrix[i][k] / augMatrix[k][k]
			for j := k; j <= n; j++ {
				augMatrix[i][j] -= factor * augMatrix[k][j]
			}
		}
	}

	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		sum := 0.0
		for j := i + 1; j < n; j++ {
			sum += augMatrix[i][j] * x[j]
		}

		if math.Abs(augMatrix[i][i]) < epsilon {
			return nil, fmt.Errorf("division by zero during back-substitution (matrix is singular)")
		}
		x[i] = (augMatrix[i][n] - sum) / augMatrix[i][i]
	}

	return x, nil
}

func main() {
	matrix1 := [][]float64{
		{1, 1, 1, 6},
		{0, 2, 5, -4},
		{2, 5, -1, 27},
	}
	fmt.Println("Solving system 1:")
	solution1, err1 := solveLinearEquations(matrix1)
	if err1 != nil {
		fmt.Printf("Error: %v\n", err1)
	} else {
		fmt.Printf("Solution: x = %v\n", solution1)
	}
	fmt.Println("--------------------")

	matrix2 := [][]float64{
		{2, 1, 7},
		{1, -3, -7},
	}
	fmt.Println("Solving system 2:")
	solution2, err2 := solveLinearEquations(matrix2)
	if err2 != nil {
		fmt.Printf("Error: %v\n", err2)
	} else {
		fmt.Printf("Solution: x = %v\n", solution2)
	}
	fmt.Println("--------------------")

	matrix3 := [][]float64{
		{1, 1, 2},
		{2, 2, 4},
	}
	fmt.Println("Solving system 3 (singular):")
	solution3, err3 := solveLinearEquations(matrix3)
	if err3 != nil {
		fmt.Printf("Error: %v\n", err3)
	} else {
		fmt.Printf("Solution: x = %v\n", solution3)
	}
	fmt.Println("--------------------")

	matrix4 := [][]float64{
		{1, 1, 2},
		{1.000000001, 1, 2.000000001},
	}
	fmt.Println("Solving system 4 (ill-conditioned):")
	solution4, err4 := solveLinearEquations(matrix4)
	if err4 != nil {
		fmt.Printf("Error: %v\n", err4)
	} else {
		fmt.Printf("Solution: x = %v\n", solution4)
	}
	fmt.Println("--------------------")

	matrix5 := [][]float64{}
	fmt.Println("Solving system 5 (empty):")
	solution5, err5 := solveLinearEquations(matrix5)
	if err5 != nil {
		fmt.Printf("Error: %v\n", err5)
	} else {
		fmt.Printf("Solution: x = %v\n", solution5)
	}
	fmt.Println("--------------------")

	matrix7 := [][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
	}
	fmt.Println("Solving system 7 (invalid dimensions):")
	solution7, err7 := solveLinearEquations(matrix7)
	if err7 != nil {
		fmt.Printf("Error: %v\n", err7)
	} else {
		fmt.Printf("Solution: x = %v\n", solution7)
	}
	fmt.Println("--------------------")

	matrix8 := [][]float64{
		{3, 9},
	}
	fmt.Println("Solving system 8 (1x1):")
	solution8, err8 := solveLinearEquations(matrix8)
	if err8 != nil {
		fmt.Printf("Error: %v\n", err8)
	} else {
		fmt.Printf("Solution: x = %v\n", solution8)
	}
	fmt.Println("--------------------")
}

// Additional implementation at 2025-06-20 00:09:37
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// solveLinearSystem solves a system of linear equations Ax = b using Gaussian elimination with partial pivoting.
// It returns the solution vector x or an error if the matrix is singular or dimensions are incompatible.
func solveLinearSystem(A [][]float64, b []float64) ([]float64, error) {
	n := len(A)
	if n == 0 {
		return nil, fmt.Errorf("empty matrix A")
	}
	if len(A[0]) != n {
		return nil, fmt.Errorf("matrix A must be square (n x n)")
	}
	if len(b) != n {
		return nil, fmt.Errorf("vector b must have n elements")
	}

	// Create an augmented matrix [A|b] for easier manipulation during elimination.
	augmentedMatrix := make([][]float64, n)
	for i := 0; i < n; i++ {
		augmentedMatrix[i] = make([]float64, n+1)
		copy(augmentedMatrix[i], A[i])
		augmentedMatrix[i][n] = b[i]
	}

	// Forward Elimination (to Upper Triangular Form)
	for k := 0; k < n; k++ {
		// Partial Pivoting: Find the row with the largest absolute value in the current column k
		// to ensure numerical stability and avoid division by zero.
		pivotRow := k
		for i := k + 1; i < n; i++ {
			if math.Abs(augmentedMatrix[i][k]) > math.Abs(augmentedMatrix[pivotRow][k]) {
				pivotRow = i
			}
		}

		// Swap the current row k with the pivotRow.
		augmentedMatrix[k], augmentedMatrix[pivotRow] = augmentedMatrix[pivotRow], augmentedMatrix[k]

		// Check if the pivot element is effectively zero after pivoting.
		// If it is, the matrix is singular or ill-conditioned, and a unique solution does not exist.
		if math.Abs(augmentedMatrix[k][k]) < 1e-9 { // Use a small epsilon for float comparison
			return nil, fmt.Errorf("matrix is singular or ill-conditioned, no unique solution exists")
		}

		// Eliminate elements below the pivot in the current column k.
		for i := k + 1; i < n; i++ {
			factor := augmentedMatrix[i][k] / augmentedMatrix[k][k]
			for j := k; j < n+1; j++ {
				augmentedMatrix[i][j] -= factor * augmentedMatrix[k][j]
			}
		}
	}

	// Back Substitution to find the solution vector x.
	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		sum := 0.0
		for j := i + 1; j < n; j++ {
			sum += augmentedMatrix[i][j] * x[j]
		}
		// Double check for division by zero, though it should ideally be caught during forward elimination.
		if math.Abs(augmentedMatrix[i][i]) < 1e-9 {
			return nil, fmt.Errorf("division by zero during back substitution (matrix likely singular)")
		}
		x[i] = (augmentedMatrix[i][n] - sum) / augmentedMatrix[i][i]
	}

	return x, nil
}

// readMatrix reads an N x N matrix from standard input, prompting the user for each row.
func readMatrix(scanner *bufio.Scanner, n int) ([][]float64, error) {
	matrix := make([][]float64, n)
	fmt.Printf("Enter the %d rows of the matrix A, each with %d space-separated numbers:\n", n, n)
	for i := 0; i < n; i++ {
		fmt.Printf("Row %d: ", i+1)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return nil, fmt.Errorf("error reading row %d: %w", i+1, err)
			}
			return nil, fmt.Errorf("failed to read row %d (EOF or empty input)", i+1)
		}
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) != n {
			return nil, fmt.Errorf("expected %d numbers in row %d, got %d", n, i+1, len(parts))
		}

		row := make([]float64, n)
		for j, s := range parts {
			val, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid number '%s' in row %d, column %d: %w", s, i+1, j+1, err)
			}
			row[j] = val
		}
		matrix[i] = row
	}
	return matrix, nil
}

// readVector reads an N-element vector from standard input, prompting the user.
func readVector(scanner *bufio.Scanner, n int, name string) ([]float64, error) {
	vector := make([]float64, n)
	fmt.Printf("Enter the %d space-separated numbers for vector %s:\n", n, name)
	fmt.Printf("Vector %s: ", name)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading vector %s: %w", name, err)
		}
		return nil, fmt.Errorf("failed to read vector %s (EOF or empty input)", name)
	}
	line := strings.TrimSpace(scanner.Text())
	parts := strings.Fields(line)
	if len(parts) != n {
		return nil, fmt.Errorf("expected %d numbers for vector %s, got %d", n, name, len(parts))
	}

	for i, s := range parts {
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number '%s' in vector %s at position %d: %w", s, name, i+1, err)
		}
		vector[i] = val
	}
	return vector, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("--- Linear Equation Solver (Ax = b) ---")
	fmt.Print("Enter the size of the system (N): ")
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading system size: %v\n", err)
		} else {
			fmt.Println("Error reading system size (EOF or empty input).")
		}
		return
	}
	nStr := strings.TrimSpace(scanner.Text())
	n, err := strconv.Atoi(nStr)
	if err != nil || n <= 0 {
		fmt.Println("Invalid system size. Please enter a positive integer.")
		return
	}

	A, err := readMatrix(scanner, n)
	if err != nil {
		fmt.Printf("Error reading matrix A: %v\n", err)
		return
	}

	b, err := readVector(scanner, n, "b")
	if err != nil {
		fmt.Printf("Error reading vector b: %v\n", err)
		return
	}

	fmt.Println("\n--- System Entered ---")
	fmt.Println("Matrix A:")
	for i := 0; i < n; i++ {
		fmt.Printf("  %v\n", A[i])
	}
	fmt.Printf("Vector b: %v\n", b)

	x, err := solveLinearSystem(A, b)
	if err != nil {
		fmt.Printf("\nError solving system: %v\n", err)
		return
	}

	fmt.Println("\n--- Solution ---")
	fmt.Println("Solution vector x:")
	for i, val := range x {
		fmt.Printf("  x[%d] = %.6f\n", i, val)
	}
}

// Additional implementation at 2025-06-20 00:10:40
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type SolutionError string

func (e SolutionError) Error() string {
	return string(e)
}

const (
	ErrNoSolution        SolutionError = "no unique solution: inconsistent system"
	ErrInfiniteSolutions SolutionError = "infinite solutions: dependent system"
	ErrInvalidInput      SolutionError = "invalid input format"
	ErrDimensionMismatch SolutionError = "dimension mismatch"
)

func solveLinearEquations(augmentedMatrix [][]float64) ([]float64, error) {
	n := len(augmentedMatrix)
	if n == 0 {
		return nil, ErrDimensionMismatch
	}
	m := len(augmentedMatrix[0])
	if m != n+1 {
		return nil, ErrDimensionMismatch
	}

	for i := 0; i < n; i++ {
		maxRow := i
		for k := i + 1; k < n; k++ {
			if math.Abs(augmentedMatrix[k][i]) > math.Abs(augmentedMatrix[maxRow][i]) {
				maxRow = k
			}
		}
		augmentedMatrix[i], augmentedMatrix[maxRow] = augmentedMatrix[maxRow], augmentedMatrix[i]

		if math.Abs(augmentedMatrix[i][i]) < 1e-9 {
			if math.Abs(augmentedMatrix[i][n]) < 1e-9 {
				return nil, ErrInfiniteSolutions
			} else {
				return nil, ErrNoSolution
			}
		}

		for k := i + 1; k < n; k++ {
			factor := augmentedMatrix[k][i] / augmentedMatrix[i][i]
			for j := i; j < m; j++ {
				augmentedMatrix[k][j] -= factor * augmentedMatrix[i][j]
			}
		}
	}

	for i := n - 1; i >= 0; i-- {
		isCoeffZero := true
		for j := 0; j < n; j++ {
			if math.Abs(augmentedMatrix[i][j]) > 1e-9 {
				isCoeffZero = false
				break
			}
		}
		if isCoeffZero {
			if math.Abs(augmentedMatrix[i][n]) > 1e-9 {
				return nil, ErrNoSolution
			} else {
				return nil, ErrInfiniteSolutions
			}
		}
	}

	solution := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		sum := 0.0
		for j := i + 1; j < n; j++ {
			sum += augmentedMatrix[i][j] * solution[j]
		}
		if math.Abs(augmentedMatrix[i][i]) < 1e-9 {
			return nil, ErrInfiniteSolutions
		}
		solution[i] = (augmentedMatrix[i][n] - sum) / augmentedMatrix[i][i]
	}

	return solution, nil
}

func readMatrixInput() ([][]float64, error) {
	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the number of equations (N): ")
	reader.Scan()
	nStr := strings.TrimSpace(reader.Text())
	n, err := strconv.Atoi(nStr)
	if err != nil || n <= 0 {
		return nil, ErrInvalidInput
	}

	fmt.Printf("Enter the augmented matrix (N rows, N+1 columns).\n")
	fmt.Printf("Each row should have N coefficients followed by the constant term, separated by spaces.\n")

	augmentedMatrix := make([][]float64, n)
	for i := 0; i < n; i++ {
		fmt.Printf("Enter row %d: ", i+1)
		reader.Scan()
		line := strings.TrimSpace(reader.Text())
		parts := strings.Fields(line)

		if len(parts) != n+1 {
			return nil, fmt.Errorf("%w: row %d has %d elements, expected %d", ErrInvalidInput, i+1, len(parts), n+1)
		}

		row := make([]float64, n+1)
		for j, part := range parts {
			val, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return nil, fmt.Errorf("%w: invalid number '%s' in row %d, column %d", ErrInvalidInput, part, i+1, j+1)
			}
			row[j] = val
		}
		augmentedMatrix[i] = row
	}
	return augmentedMatrix, nil
}

func main() {
	augmentedMatrix, err := readMatrixInput()
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	solution, err := solveLinearEquations(augmentedMatrix)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Solution:")
	for i, val := range solution {
		fmt.Printf("x%d = %.6f\n", i+1, val)
	}
}