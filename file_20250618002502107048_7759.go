package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Matrix [][]float64

func readMatrix() (Matrix, error) {
	fmt.Print("Enter number of rows: ")
	var rows int
	_, err := fmt.Scanln(&rows)
	if err != nil {
		return nil, fmt.Errorf("invalid row input: %w", err)
	}
	if rows <= 0 {
		return nil, fmt.Errorf("number of rows must be positive")
	}

	fmt.Print("Enter number of columns: ")
	var cols int
	_, err = fmt.Scanln(&cols)
	if err != nil {
		return nil, fmt.Errorf("invalid column input: %w", err)
	}
	if cols <= 0 {
		return nil, fmt.Errorf("number of columns must be positive")
	}

	matrix := make(Matrix, rows)
	fmt.Printf("Enter %d rows of %d space-separated numbers:\n", rows, cols)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]float64, cols)
		var rowInput string
		fmt.Printf("Row %d: ", i+1)
		_, err := fmt.Scanln(&rowInput)
		if err != nil {
			return nil, fmt.Errorf("error reading row %d: %w", i+1, err)
		}

		parts := strings.Fields(rowInput)
		if len(parts) != cols {
			return nil, fmt.Errorf("expected %d numbers for row %d, got %d", cols, i+1, len(parts))
		}

		for j, part := range parts {
			val, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid number '%s' in row %d, column %d: %w", part, i+1, j+1, err)
			}
			matrix[i][j] = val
		}
	}
	return matrix, nil
}

func printMatrix(m Matrix) {
	if m == nil || len(m) == 0 || (len(m) > 0 && len(m[0]) == 0) {
		fmt.Println("Empty matrix.")
		return
	}
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			fmt.Printf("%8.2f ", m[i][j])
		}
		fmt.Println()
	}
}

func addMatrices(m1, m2 Matrix) (Matrix, error) {
	if len(m1) == 0 || len(m2) == 0 || len(m1[0]) == 0 || len(m2[0]) == 0 {
		return nil, fmt.Errorf("matrices cannot be empty")
	}
	if len(m1) != len(m2) || len(m1[0]) != len(m2[0]) {
		return nil, fmt.Errorf("matrices must have the same dimensions for addition")
	}

	rows := len(m1)
	cols := len(m1[0])
	result := make(Matrix, rows)
	for i := 0; i < rows; i++ {
		result[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			result[i][j] = m1[i][j] + m2[i][j]
		}
	}
	return result, nil
}

func subtractMatrices(m1, m2 Matrix) (Matrix, error) {
	if len(m1) == 0 || len(m2) == 0 || len(m1[0]) == 0 || len(m2[0]) == 0 {
		return nil, fmt.Errorf("matrices cannot be empty")
	}
	if len(m1) != len(m2) || len(m1[0]) != len(m2[0]) {
		return nil, fmt.Errorf("matrices must have the same dimensions for subtraction")
	}

	rows := len(m1)
	cols := len(m1[0])
	result := make(Matrix, rows)
	for i := 0; i < rows; i++ {
		result[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			result[i][j] = m1[i][j] - m2[i][j]
		}
	}
	return result, nil
}

func multiplyMatrices(m1, m2 Matrix) (Matrix, error) {
	if len(m1) == 0 || len(m1[0]) == 0 || len(m2) == 0 || len(m2[0]) == 0 {
		return nil, fmt.Errorf("matrices cannot be empty")
	}
	if len(m1[0]) != len(m2) {
		return nil, fmt.Errorf("dimensions mismatch for multiplication: columns of first matrix (%d) must equal rows of second matrix (%d)", len(m1[0]), len(m2))
	}

	rows1 := len(m1)
	cols1 := len(m1[0])
	cols2 := len(m2[0])

	result := make(Matrix, rows1)
	for i := 0; i < rows1; i++ {
		result[i] = make([]float64, cols2)
		for j := 0; j < cols2; j++ {
			sum := 0.0
			for k := 0; k < cols1; k++ {
				sum += m1[i][k] * m2[k][j]
			}
			result[i][j] = sum
		}
	}
	return result, nil
}

func scalarMultiplyMatrix(m Matrix, scalar float64) (Matrix, error) {
	if m == nil || len(m) == 0 || len(m[0]) == 0 {
		return nil, fmt.Errorf("matrix cannot be empty for scalar multiplication")
	}
	rows := len(m)
	cols := len(m[0])
	result := make(Matrix, rows)
	for i := 0; i < rows; i++ {
		result[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			result[i][j] = m[i][j] * scalar
		}
	}
	return result, nil
}

func transposeMatrix(m Matrix) (Matrix, error) {
	if m == nil || len(m) == 0 || len(m[0]) == 0 {
		return nil, fmt.Errorf("matrix cannot be empty for transposition")
	}
	rows := len(m)
	cols := len(m[0])

	result := make(Matrix, cols)
	for i := 0; i < cols; i++ {
		result[i] = make([]float64, rows)
		for j := 0; j < rows; j++ {
			result[i][j] = m[j][i]
		}
	}
	return result, nil
}

func main() {
	for {
		fmt.Println("\nMatrix Operations Calculator")
		fmt.Println("1. Add Matrices")
		fmt.Println("2. Subtract Matrices")
		fmt.Println("3. Multiply Matrices")
		fmt.Println("4. Scalar Multiply Matrix")
		fmt.Println("5. Transpose Matrix")
		fmt.Println("6. Exit")
		fmt.Print("Choose an operation: ")

		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		var m1, m2, result Matrix
		var scalar float64

		switch choice {
		case 1:
			fmt.Println("\nEnter Matrix 1:")
			m1, err = readMatrix()
			if err != nil {
				fmt.Printf("Error reading Matrix 1: %v\n", err)
				continue
			}
			fmt.Println("\nEnter Matrix 2:")
			m2, err = readMatrix()
			if err != nil {
				fmt.Printf("Error reading Matrix 2: %v\n", err)
				continue
			}
			result, err = addMatrices(m1, m2)
			if err != nil {
				fmt.Printf("Error performing addition: %v\n", err)
				continue
			}
			fmt.Println("\nResult of Addition:")
			printMatrix(result)
		case 2:
			fmt.Println("\nEnter Matrix 1:")
			m1, err = readMatrix()
			if err != nil {
				fmt.Printf("Error reading Matrix 1: %v\n", err)
				continue
			}
			fmt.Println("\nEnter Matrix 2:")
			m2, err = readMatrix()
			if err != nil {
				fmt.Printf("Error reading Matrix 2: %v\n", err)
				continue
			}
			result, err = subtractMatrices(m1, m2)
			if err != nil {
				fmt.Printf("Error performing subtraction: %v\n", err)
				continue
			}
			fmt.Println("\nResult of Subtraction:")
			printMatrix(result)
		case 3:
			fmt.Println("\nEnter Matrix 1:")
			m1, err = readMatrix()
			if err != nil {
				fmt.Printf("Error reading Matrix 1: %v\n", err)
				continue
			}
			fmt.Println("\nEnter Matrix 2:")
			m2, err = readMatrix()
			if err != nil {
				fmt.Printf("Error reading Matrix 2: %v\n", err)
				continue
			}
			result, err = multiplyMatrices(m1, m2)
			if err != nil {
				fmt.Printf("Error performing multiplication: %v\n", err)
				continue
			}
			fmt.Println("\nResult of Multiplication:")
			printMatrix(result)
		case 4:
			fmt.Println("\nEnter Matrix:")
			m1, err = readMatrix()
			if err != nil {
				fmt.Printf("Error reading Matrix: %v\n", err)
				continue
			}
			fmt.Print("Enter scalar value: ")
			_, err = fmt.Scanln(&scalar)
			if err != nil {
				fmt.Println("Invalid scalar input. Please enter a number.")
				continue
			}
			result, err = scalarMultiplyMatrix(m1, scalar)
			if err != nil {
				fmt.Printf("Error performing scalar multiplication: %v\n", err)
				continue
			}
			fmt.Println("\nResult of Scalar Multiplication:")
			printMatrix(result)
		case 5:
			fmt.Println("\nEnter Matrix:")
			m1, err = readMatrix()
			if err != nil {
				fmt.Printf("Error reading Matrix: %v\n", err)
				continue
			}
			result, err = transposeMatrix(m1)
			if err != nil {
				fmt.Printf("Error performing transpose: %v\n", err)
				continue
			}
			fmt.Println("\nResult of Transposition:")
			printMatrix(result)
		case 6:
			fmt.Println("Exiting calculator. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 6.")
		}
	}
}

// Additional implementation at 2025-06-18 00:25:51
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Matrix represents a 2D matrix of float64 values.
type Matrix struct {
	rows int
	cols int
	data [][]float64
}

// NewMatrix creates a new Matrix with specified dimensions, initialized to zeros.
func NewMatrix(rows, cols int) *Matrix {
	if rows <= 0 || cols <= 0 {
		return nil
	}
	data := make([][]float64, rows)
	for i := range data {
		data[i] = make([]float64, cols)
	}
	return &Matrix{rows: rows, cols: cols, data: data}
}

// NewMatrixFromSlice creates a new Matrix from a 2D slice.
// It performs a deep copy of the data.
func NewMatrixFromSlice(data [][]float64) *Matrix {
	if len(data) == 0 || len(data[0]) == 0 {
		return nil
	}
	rows := len(data)
	cols := len(data[0])
	// Ensure all rows have the same number of columns
	for _, row := range data {
		if len(row) != cols {
			return nil // Irregular matrix
		}
	}

	newData := make([][]float64, rows)
	for i := range data {
		newData[i] = make([]float64, cols)
		copy(newData[i], data[i])
	}
	return &Matrix{rows: rows, cols: cols, data: newData}
}

// Print displays the matrix to the console.
func (m *Matrix) Print() {
	if m == nil {
		fmt.Println("Matrix is nil.")
		return
	}
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			fmt.Printf("%10.4f ", m.data[i][j])
		}
		fmt.Println()
	}
}

// Add performs matrix addition.
func (m *Matrix) Add(other *Matrix) (*Matrix, error) {
	if m.rows != other.rows || m.cols != other.cols {
		return nil, fmt.Errorf("matrix dimensions must match for addition: %dx%d vs %dx%d", m.rows, m.cols, other.rows, other.cols)
	}
	result := NewMatrix(m.rows, m.cols)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[i][j] = m.data[i][j] + other.data[i][j]
		}
	}
	return result, nil
}

// Subtract performs matrix subtraction.
func (m *Matrix) Subtract(other *Matrix) (*Matrix, error) {
	if m.rows != other.rows || m.cols != other.cols {
		return nil, fmt.Errorf("matrix dimensions must match for subtraction: %dx%d vs %dx%d", m.rows, m.cols, other.rows, other.cols)
	}
	result := NewMatrix(m.rows, m.cols)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[i][j] = m.data[i][j] - other.data[i][j]
		}
	}
	return result, nil
}

// Multiply performs matrix multiplication (m * other).
func (m *Matrix) Multiply(other *Matrix) (*Matrix, error) {
	if m.cols != other.rows {
		return nil, fmt.Errorf("incompatible dimensions for multiplication: %dx%d and %dx%d", m.rows, m.cols, other.rows, other.cols)
	}
	result := NewMatrix(m.rows, other.cols)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < other.cols; j++ {
			sum := 0.0
			for k := 0; k < m.cols; k++ {
				sum += m.data[i][k] * other.data[k][j]
			}
			result.data[i][j] = sum
		}
	}
	return result, nil
}

// ScalarMultiply performs scalar multiplication.
func (m *Matrix) ScalarMultiply(scalar float64) *Matrix {
	result := NewMatrix(m.rows, m.cols)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[i][j] = m.data[i][j] * scalar
		}
	}
	return result
}

// Transpose returns the transpose of the matrix.
func (m *Matrix) Transpose() *Matrix {
	result := NewMatrix(m.cols, m.rows)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[j][i] = m.data[i][j]
		}
	}
	return result
}

// Determinant calculates the determinant of a square matrix using cofactor expansion.
func (m *Matrix) Determinant() (float64, error) {
	if m.rows != m.cols {
		return 0, fmt.Errorf("matrix must be square to calculate determinant")
	}

	n := m.rows
	if n == 1 {
		return m.data[0][0], nil
	}
	if n == 2 {
		return m.data[0][0]*m.data[1][1] - m.data[0][1]*m.data[1][0], nil
	}

	det := 0.0
	for j := 0; j < n; j++ {
		subMatrix := m.createSubMatrix(0, j)
		subDet, err := subMatrix.Determinant()
		if err != nil {
			return 0, err
		}
		term := m.data[0][j] * subDet
		if j%2 == 1 { // If column index is odd, subtract the term
			det -= term
		} else { // If column index is even, add the term
			det += term
		}
	}
	return det, nil
}

// createSubMatrix creates a sub-matrix by removing the specified row and column.
func (m *Matrix) createSubMatrix(rowToRemove, colToRemove int) *Matrix {
	subRows := m.rows - 1
	subCols := m.cols - 1
	subMatrix := NewMatrix(subRows, subCols)

	currSubRow := 0
	for i := 0; i < m.rows; i++ {
		if i == rowToRemove {
			continue
		}
		currSubCol := 0
		for j := 0; j < m.cols; j++ {
			if j == colToRemove {
				continue
			}
			subMatrix.data[currSubRow][currSubCol] = m.data[i][j]
			currSubCol++
		}
		currSubRow++
	}
	return subMatrix
}

// Adjugate calculates the adjugate matrix.
func (m *Matrix) Adjugate() (*Matrix, error) {
	if m.rows != m.cols {
		return nil, fmt.Errorf("matrix must be square to calculate adjugate")
	}
	n := m.rows
	if n == 1 {
		return NewMatrixFromSlice([][]float64{{1.0}}), nil // Adjugate of [a] is [1]
	}

	cofactorMatrix := NewMatrix(n, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			subMatrix := m.createSubMatrix(i, j)
			subDet, err := subMatrix.Determinant()
			if err != nil {
				return nil, err
			}
			cofactor := subDet
			if (i+j)%2 == 1 { // Apply sign based on position
				cofactor *= -1
			}
			cofactorMatrix.data[i][j] = cofactor
		}
	}
	return cofactorMatrix.Transpose(), nil // Adjugate is transpose of cofactor matrix
}

// Inverse calculates the inverse of a square matrix.
func (m *Matrix) Inverse() (*Matrix, error) {
	if m.rows != m.cols {
		return nil, fmt.Errorf("matrix must be square to calculate inverse")
	}

	det, err := m.Determinant()
	if err != nil {
		return nil, err
	}

	// Use a small epsilon for floating point comparison to check for singularity
	if math.Abs(det) < 1e-9 {
		return nil, fmt.Errorf("matrix is singular (determinant is zero or very close to zero), cannot calculate inverse")
	}

	adj, err := m.Adjugate()
	if err != nil {
		return nil, err
	}

	invDet := 1.0 / det
	return adj.ScalarMultiply(invDet), nil
}

// IdentityMatrix creates an identity matrix of size n x n.
func IdentityMatrix(n int) *Matrix {
	if n <= 0 {
		return nil
	}
	m := NewMatrix(n, n)
	for i := 0; i < n; i++ {
		m.data[i][i] = 1.0
	}
	return m
}

// Helper function to read a float from console.
func readFloat(reader *bufio.Reader, prompt string) (float64, error) {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	val, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %v", err)
	}
	return val, nil
}

// Helper function to read an integer from console.
func readInt(reader *bufio.Reader, prompt string) (int, error) {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	val, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("invalid integer: %v", err)
	}
	return val, nil
}

// readMatrix reads matrix dimensions and elements from the user.
func readMatrix(reader *bufio.Reader, matrixName string) (*Matrix, error) {
	fmt.Printf("Enter dimensions for %s matrix:\n", matrixName)
	rows, err := readInt(reader, "Enter number of rows: ")
	if err != nil {
		return nil, err
	}
	cols, err := readInt(reader, "Enter number of columns: ")
	if err != nil {
		return nil, err
	}

	m := NewMatrix(rows, cols)
	if m == nil {
		return nil, fmt.Errorf("failed to create matrix with dimensions %dx%d", rows, cols)
	}

	fmt.Printf("Enter elements for %

// Additional implementation at 2025-06-18 00:26:59
package main

import (
	"errors"
	"fmt"
	"math"
)

// Matrix represents a 2D matrix of float64 values.
type Matrix struct {
	rows, cols int
	data       [][]float64
}

// NewMatrix creates a new matrix with the given dimensions, initialized to zeros.
func NewMatrix(rows, cols int) (*Matrix, error) {
	if rows <= 0 || cols <= 0 {
		return nil, errors.New("matrix dimensions must be positive")
	}
	data := make([][]float64, rows)
	for i := range data {
		data[i] = make([]float64, cols)
	}
	return &Matrix{rows: rows, cols: cols, data: data}, nil
}

// NewMatrixFromData creates a new matrix from a 2D slice of float64.
// It performs a deep copy of the input data.
func NewMatrixFromData(data [][]float64) (*Matrix, error) {
	if len(data) == 0 || len(data[0]) == 0 {
		return nil, errors.New("input data cannot be empty or have empty rows")
	}
	rows := len(data)
	cols := len(data[0])
	for _, row := range data {
		if len(row) != cols {
			return nil, errors.New("all rows in input data must have the same number of columns")
		}
	}
	newData := make([][]float64, rows)
	for i := range data {
		newData[i] = make([]float64, cols)
		copy(newData[i], data[i]) // Deep copy
	}
	return &Matrix{rows: rows, cols: cols, data: newData}, nil
}

// Print prints the matrix to the console with formatted output.
func (m *Matrix) Print() {
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			fmt.Printf("%10.4f ", m.data[i][j])
		}
		fmt.Println()
	}
}

// Add performs matrix addition (m + other).
// Returns a new matrix or an error if dimensions do not match.
func (m *Matrix) Add(other *Matrix) (*Matrix, error) {
	if m.rows != other.rows || m.cols != other.cols {
		return nil, errors.New("matrices must have the same dimensions for addition")
	}
	result, _ := NewMatrix(m.rows, m.cols) // Error already checked by dimensions
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[i][j] = m.data[i][j] + other.data[i][j]
		}
	}
	return result, nil
}

// Subtract performs matrix subtraction (m - other).
// Returns a new matrix or an error if dimensions do not match.
func (m *Matrix) Subtract(other *Matrix) (*Matrix, error) {
	if m.rows != other.rows || m.cols != other.cols {
		return nil, errors.New("matrices must have the same dimensions for subtraction")
	}
	result, _ := NewMatrix(m.rows, m.cols) // Error already checked by dimensions
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[i][j] = m.data[i][j] - other.data[i][j]
		}
	}
	return result, nil
}

// ScalarMultiply performs scalar multiplication (m * scalar).
// Returns a new matrix.
func (m *Matrix) ScalarMultiply(scalar float64) *Matrix {
	result, _ := NewMatrix(m.rows, m.cols) // Dimensions are valid from m
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[i][j] = m.data[i][j] * scalar
		}
	}
	return result
}

// Multiply performs matrix multiplication (m * other).
// Returns a new matrix or an error if matrices are not conformable for multiplication.
func (m *Matrix) Multiply(other *Matrix) (*Matrix, error) {
	if m.cols != other.rows {
		return nil, errors.New("number of columns in the first matrix must equal number of rows in the second matrix for multiplication")
	}
	result, _ := NewMatrix(m.rows, other.cols) // Dimensions are valid
	for i := 0; i < m.rows; i++ {
		for j := 0; j < other.cols; j++ {
			sum := 0.0
			for k := 0; k < m.cols; k++ { // m.cols == other.rows
				sum += m.data[i][k] * other.data[k][j]
			}
			result.data[i][j] = sum
		}
	}
	return result, nil
}

// Transpose returns the transpose of the matrix.
// Returns a new matrix.
func (m *Matrix) Transpose() *Matrix {
	result, _ := NewMatrix(m.cols, m.rows) // Dimensions are valid
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[j][i] = m.data[i][j]
		}
	}
	return result
}

// IsSquare checks if the matrix is a square matrix.
func (m *Matrix) IsSquare() bool {
	return m.rows == m.cols
}

// Cofactor returns the minor matrix for the element at (row, col).
// This is used internally for determinant and inverse calculations.
func (m *Matrix) Cofactor(row, col int) (*Matrix, error) {
	if !m.IsSquare() {
		return nil, errors.New("cofactor can only be calculated for square matrices")
	}
	if row < 0 || row >= m.rows || col < 0 || col >= m.cols {
		return nil, errors.New("row or column index out of bounds for cofactor")
	}

	n := m.rows
	minor, _ := NewMatrix(n-1, n-1)
	minorRow, minorCol := 0, 0
	for i := 0; i < n; i++ {
		if i == row {
			continue
		}
		minorCol = 0
		for j := 0; j < n; j++ {
			if j == col {
				continue
			}
			minor.data[minorRow][minorCol] = m.data[i][j]
			minorCol++
		}
		minorRow++
	}
	return minor, nil
}

// Determinant calculates the determinant of a square matrix using Laplace expansion.
// Returns the determinant value or an error if the matrix is not square.
// Note: This method is computationally expensive for large matrices (O(n!)).
func (m *Matrix) Determinant() (float64, error) {
	if !m.IsSquare() {
		return 0, errors.New("determinant can only be calculated for square matrices")
	}

	n := m.rows
	if n == 1 {
		return m.data[0][0], nil
	}
	if n == 2 {
		return m.data[0][0]*m.data[1][1] - m.data[0][1]*m.data[1][0], nil
	}

	det := 0.0
	for j := 0; j < n; j++ {
		minor, err := m.Cofactor(0, j)
		if err != nil {
			return 0, err // Should not happen if IsSquare() is true
		}
		minorDet, err := minor.Determinant()
		if err != nil {
			return 0, err
		}
		term := m.data[0][j] * minorDet
		if j%2 == 1 { // If column index is odd, subtract
			det -= term
		} else { // If column index is even, add
			det += term
		}
	}
	return det, nil
}

// Adjoint calculates the adjoint of a square matrix.
// Returns a new matrix or an error if the matrix is not square.
func (m *Matrix) Adjoint() (*Matrix, error) {
	if !m.IsSquare() {
		return nil, errors.New("adjoint can only be calculated for square matrices")
	}
	n := m.rows
	adj, _ := NewMatrix(n, n) // Dimensions are valid

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			cofactorMinor, err := m.Cofactor(i, j)
			if err != nil {
				return nil, err // Should not happen
			}
			det, err := cofactorMinor.Determinant()
			if err != nil {
				return nil, err
			}
			// Apply sign based on (i+j)
			if (i+j)%2 == 1 {
				det = -det
			}
			adj.data[j][i] = det // Transpose for adjoint (adj[j][i] = cofactor[i][j])
		}
	}
	return adj, nil
}

// Inverse calculates the inverse of an invertible square matrix.
// Returns a new matrix or an error if the matrix is not square or is singular.
func (m *Matrix) Inverse() (*Matrix, error) {
	if !m.IsSquare() {
		return nil, errors.New("inverse can only be calculated for square matrices")
	}

	det, err := m.Determinant()
	if err != nil {
		return nil, err
	}
	// Check for near-zero determinant (singular matrix)
	if math.Abs(det) < 1e-9 { // Using a small epsilon for floating point comparison
		return nil, errors.New("matrix is singular, inverse does not exist")
	}

	adj, err := m.Adjoint()
	if err != nil {
		return nil, err
	}

	// Inverse = (1/det) * Adjoint
	inverse := adj.ScalarMultiply(1 / det)
	return inverse, nil
}

// IdentityMatrix creates an identity matrix of given size n (n x n).
// Returns a new matrix or an error if size is not positive.
func IdentityMatrix(n int) (*Matrix, error) {
	if n <= 0 {
		return nil, errors.New("size for identity matrix must be positive")
	}
	m, _ := NewMatrix(n, n) // Dimensions are valid
	for i := 0; i < n; i++ {
		m.data[i][i] = 1.0
	}
	return m, nil
}

// ZeroMatrix creates a zero matrix of given dimensions.
// Returns a new matrix or an error if dimensions are not positive.
func ZeroMatrix(rows, cols int) (*Matrix, error) {
	return NewMatrix(rows, cols) // NewMatrix already initializes with zeros
}

// ReadMatrixFromConsole prompts the user to enter matrix dimensions and elements.
// Returns a new matrix or an error if input is invalid.
func ReadMatrixFromConsole() (*Matrix, error) {
	var rows, cols int
	fmt.Print("Enter number of rows: ")
	_, err := fmt.Scan(&rows)
	if err != nil {
		return nil, fmt.Errorf("invalid input for rows: %w", err)
	}
	fmt.Print("Enter number of columns: ")
	_, err = fmt.Scan(&cols)
	if err != nil {
		return nil, fmt.Errorf("invalid input for columns: %w", err)
	}

	m, err := NewMatrix(rows, cols)
	if err != nil {
		return nil, err // Propagate error from NewMatrix
	}

	fmt.Println("Enter matrix elements row by row:")
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			fmt.Printf("Enter element [%d][%d]: ", i, j)
			_, err := fmt.Scan(&m.data[i][j])
			if err != nil {
				return nil, fmt.Errorf("invalid input for element [%d][%d]: %w", i, j, err)
			}
		}
	}
	return m, nil
}

// main function provides a menu-driven interface for the matrix calculator.
func main() {
	fmt.Println("--- Go Matrix Operations Calculator ---")

	for {
		fmt.Println("\nChoose an operation:")
		fmt.Println("1. Add Matrices")
		fmt.Println("2. Subtract Matrices")
		fmt.Println("3. Scalar Multiply Matrix")
		fmt.Println("4. Multiply Matrices")
		fmt.Println("5. Transpose Matrix")
		fmt.Println("6. Calculate Determinant")
		fmt.Println("7. Calculate Inverse")
		fmt.Println("8. Create Identity Matrix")
		fmt.Println("9. Create Zero Matrix")
		fmt.Println("0. Exit")

		var choice int
		fmt.Print("Enter your choice: ")
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			// Clear the invalid input from the buffer to prevent infinite loop
			var discard string
			fmt.Scanln(&discard

// Additional implementation at 2025-06-18 00:27:59
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Matrix represents a 2D matrix.
type Matrix struct {
	Rows int
	Cols int
	Data [][]float64
}

// NewMatrix creates a new matrix with specified dimensions, initialized to zeros.
func NewMatrix(rows, cols int) (Matrix, error) {
	if rows <= 0 || cols <= 0 {
		return Matrix{}, fmt.Errorf("matrix dimensions must be positive")
	}
	data := make([][]float64, rows)
	for i := range data {
		data[i] = make([]float64, cols)
	}
	return Matrix{Rows: rows, Cols: cols, Data: data}, nil
}

// IdentityMatrix creates an identity matrix of size n.
func IdentityMatrix(n int) (Matrix, error) {
	if n <= 0 {
		return Matrix{}, fmt.Errorf("identity matrix size must be positive")
	}
	mat, err := NewMatrix(n, n)
	if err != nil {
		return Matrix{}, err
	}
	for i := 0; i < n; i++ {
		mat.Data[i][i] = 1.0
	}
	return mat, nil
}

// ZeroMatrix creates a zero matrix of specified dimensions.
func ZeroMatrix(rows, cols int) (Matrix, error) {
	return NewMatrix(rows, cols) // NewMatrix already initializes to zeros
}

// Add adds two matrices.
func (m Matrix) Add(other Matrix) (Matrix, error) {
	if m.Rows != other.Rows || m.Cols != other.Cols {
		return Matrix{}, fmt.Errorf("matrices must have the same dimensions for addition")
	}
	result, _ := NewMatrix(m.Rows, m.Cols)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			result.Data[i][j] = m.Data[i][j] + other.Data[i][j]
		}
	}
	return result, nil
}

// Subtract subtracts one matrix from another.
func (m Matrix) Subtract(other Matrix) (Matrix, error) {
	if m.Rows != other.Rows || m.Cols != other.Cols {
		return Matrix{}, fmt.Errorf("matrices must have the same dimensions for subtraction")
	}
	result, _ := NewMatrix(m.Rows, m.Cols)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			result.Data[i][j] = m.Data[i][j] - other.Data[i][j]
		}
	}
	return result, nil
}

// ScalarMultiply multiplies a matrix by a scalar.
func (m Matrix) ScalarMultiply(scalar float64) Matrix {
	result, _ := NewMatrix(m.Rows, m.Cols)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			result.Data[i][j] = m.Data[i][j] * scalar
		}
	}
	return result
}

// Multiply multiplies two matrices.
func (m Matrix) Multiply(other Matrix) (Matrix, error) {
	if m.Cols != other.Rows {
		return Matrix{}, fmt.Errorf("number of columns in the first matrix must equal number of rows in the second matrix for multiplication")
	}
	result, _ := NewMatrix(m.Rows, other.Cols)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < other.Cols; j++ {
			sum := 0.0
			for k := 0; k < m.Cols; k++ {
				sum += m.Data[i][k] * other.Data[k][j]
			}
			result.Data[i][j] = sum
		}
	}
	return result, nil
}

// Transpose returns the transpose of the matrix.
func (m Matrix) Transpose() Matrix {
	result, _ := NewMatrix(m.Cols, m.Rows)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			result.Data[j][i] = m.Data[i][j]
		}
	}
	return result
}

// Trace returns the trace of a square matrix (sum of diagonal elements).
func (m Matrix) Trace() (float64, error) {
	if m.Rows != m.Cols {
		return 0, fmt.Errorf("matrix must be square to calculate trace")
	}
	trace := 0.0
	for i := 0; i < m.Rows; i++ {
		trace += m.Data[i][i]
	}
	return trace, nil
}

// SubMatrix returns a submatrix by removing a specified row and column.
func (m Matrix) SubMatrix(excludeRow, excludeCol int) (Matrix, error) {
	if m.Rows <= 1 || m.Cols <= 1 {
		return Matrix{}, fmt.Errorf("cannot create submatrix from a 1x1 or smaller matrix")
	}
	if excludeRow < 0 || excludeRow >= m.Rows || excludeCol < 0 || excludeCol >= m.Cols {
		return Matrix{}, fmt.Errorf("invalid row or column to exclude for submatrix")
	}

	result, _ := NewMatrix(m.Rows-1, m.Cols-1)
	currRow := 0
	for i := 0; i < m.Rows; i++ {
		if i == excludeRow {
			continue
		}
		currCol := 0
		for j := 0; j < m.Cols; j++ {
			if j == excludeCol {
				continue
			}
			result.Data[currRow][currCol] = m.Data[i][j]
			currCol++
		}
		currRow++
	}
	return result, nil
}

// Determinant calculates the determinant of a square matrix.
func (m Matrix) Determinant() (float64, error) {
	if m.Rows != m.Cols {
		return 0, fmt.Errorf("matrix must be square to calculate determinant")
	}

	n := m.Rows

	if n == 1 {
		return m.Data[0][0], nil
	}
	if n == 2 {
		return m.Data[0][0]*m.Data[1][1] - m.Data[0][1]*m.Data[1][0], nil
	}

	det := 0.0
	for j := 0; j < n; j++ {
		sub, err := m.SubMatrix(0, j)
		if err != nil {
			return 0, err // Should not happen for n > 1
		}
		subDet, err := sub.Determinant()
		if err != nil {
			return 0, err
		}
		term := m.Data[0][j] * subDet
		if j%2 == 1 { // If column index is odd, subtract
			det -= term
		} else { // If column index is even, add
			det += term
		}
	}
	return det, nil
}

// CofactorMatrix calculates the cofactor matrix.
func (m Matrix) CofactorMatrix() (Matrix, error) {
	if m.Rows != m.Cols {
		return Matrix{}, fmt.Errorf("matrix must be square to calculate cofactor matrix")
	}

	n := m.Rows
	if n == 1 {
		return Matrix{Rows: 1, Cols: 1, Data: [][]float64{{1.0}}}, nil // Cofactor of a 1x1 matrix is 1
	}

	cofactorMat, _ := NewMatrix(n, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			sub, err := m.SubMatrix(i, j)
			if err != nil {
				return Matrix{}, err // Should not happen
			}
			subDet, err := sub.Determinant()
			if err != nil {
				return Matrix{}, err
			}
			cofactor := subDet
			if (i+j)%2 == 1 {
				cofactor = -cofactor
			}
			cofactorMat.Data[i][j] = cofactor
		}
	}
	return cofactorMat, nil
}

// Adjugate calculates the adjugate (or adjoint) matrix.
func (m Matrix) Adjugate() (Matrix, error) {
	cofactorMat, err := m.CofactorMatrix()
	if err != nil {
		return Matrix{}, err
	}
	return cofactorMat.Transpose(), nil
}

// Inverse calculates the inverse of a square matrix.
func (m Matrix) Inverse() (Matrix, error) {
	if m.Rows != m.Cols {
		return Matrix{}, fmt.Errorf("matrix must be square to calculate inverse")
	}

	det, err := m.Determinant()
	if err != nil {
		return Matrix{}, err
	}

	const epsilon = 1e-9 // Tolerance for checking against zero
	if math.Abs(det) < epsilon {
		return Matrix{}, fmt.Errorf("matrix is singular (determinant is zero), cannot calculate inverse")
	}

	adj, err := m.Adjugate()
	if err != nil {
		return Matrix{}, err
	}

	return adj.ScalarMultiply(1.0 / det), nil
}

// PrintMatrix prints the matrix to the console.
func PrintMatrix(m Matrix) {
	if m.Rows == 0 || m.Cols == 0 {
		fmt.Println("[Empty Matrix]")
		return
	}
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			fmt.Printf("%10.4f ", m.Data[i][j])
		}
		fmt.Println()
	}
}

// readInt reads an integer from stdin.
func readInt(prompt string) (int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	val, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("invalid integer input: %w", err)
	}
	return val, nil
}

// readFloat64 reads a float64 from stdin.
func readFloat64(prompt string) (float64, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	val, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float input: %w", err)
	}
	return val, nil
}

// ReadMatrixFromInput prompts the user to enter matrix dimensions and elements.
func ReadMatrixFromInput(name string) (Matrix, error) {
	fmt.Printf("Enter dimensions for Matrix %s:\n", name)
	rows, err := readInt("Enter number of rows: ")
	if err != nil {
		return Matrix{}, err
	}
	cols, err := readInt("Enter number of columns: ")
	if err != nil {
		return Matrix{}, err
	}

	mat, err := NewMatrix(rows, cols)
	if err != nil {
		return Matrix{}, err
	}

	fmt.Printf("Enter elements for Matrix %s (%d rows, %d columns):\n", name, rows, cols)
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < rows; i++ {
		fmt.Printf("Enter elements for row %d (space-separated): ", i+1)
		line, _ := reader.ReadString('\n')
		parts := strings.Fields(line)
		if len(parts) != cols {
			return Matrix{}, fmt.Errorf("expected %d elements for row %d, got %d", cols, i+1, len(parts))
		}
		for j := 0; j < cols; j++ {
			val, err := strconv.ParseFloat(parts[j], 64)
			if err != nil {
				return Matrix{}, fmt.Errorf("invalid float element '%s' in row %d, column %d: %w", parts[j], i+1, j+1, err)
			}
			mat.Data[i][j] = val
		}
	}
	return mat, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n--- Matrix Operations Calculator ---")
		fmt.Println("1. Add Matrices")
		fmt.Println("2. Subtract Matrices")
		fmt.Println("3. Scalar Multiply Matrix")
		fmt.Println("4. Multiply Matrices")
		fmt.Println("5. Transpose Matrix")
		fmt.Println("6. Calculate Determinant")
		fmt.Println("7. Calculate Trace")
		fmt.Println("8. Calculate Cofactor Matrix")
		fmt.Println("9. Calculate Adjugate Matrix")
		fmt.Println("10. Calculate Inverse Matrix")
		fmt.Println("11. Create Identity Matrix")
		fmt.Println("12. Create Zero Matrix")
		fmt.Println("13. Exit")
		fmt.Print("Enter your choice: ")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Invalid choice. Please enter a number.")
			continue
		}

		var matA, matB Matrix
		var scalar float64
		var result Matrix
		var det float64
		var trace float64
		var size int

		switch choice {
		case 1, 2, 4: // Operations requiring two matrices
			matA, err = ReadMatrixFromInput("A")
			if err != nil {
				fmt.Println("Error reading Matrix A:", err)
				continue
			}
			matB, err = ReadMatrixFromInput("B")
			if err != nil {
				fmt.Println("Error reading Matrix B:", err)
				continue
			}
			fmt.Println("\nMatrix A:")
			PrintMatrix(matA)
			fmt.Println("\nMatrix B:")
			PrintMatrix(matB)

			switch choice {
			case 1:
				result, err = matA.Add(matB)
				if err == nil {
					fmt.Println("\nResult of A + B:")
					PrintMatrix(result)
				}
			case 2:
				result, err = matA.Subtract(matB)
				if err == nil {
					fmt.Println("\nResult of A - B:")
					PrintMatrix(result)
				}
			case 4:
				result, err = matA.Multiply(matB)
				if err == nil {
					fmt.Println("\nResult of A * B:")
					PrintMatrix(result)
				}
			}
		case 3, 5, 6, 7, 8, 9, 10: // Operations requiring one matrix
			matA, err = ReadMatrixFromInput("A")
			if err != nil {
				fmt.Println("Error reading Matrix A:", err)
				continue
			}
			fmt.Println("\nMatrix A:")
			PrintMatrix(matA)

			switch choice {
			case 3:
				scalar, err = readFloat64("Enter scalar value: ")
				if err != nil {
					fmt.Println("Error reading scalar:", err)
					continue
				}
				result = matA.ScalarMultiply(scalar)
				fmt.Println("\nResult of Scalar Multiplication:")
				PrintMatrix(result)
			case 5:
				result = matA.Transpose()
				fmt.Println("\nTranspose of Matrix A:")
				PrintMatrix(result)
			case 6:
				det, err = matA.Determinant()
				if err == nil {
					fmt.Printf("\nDeterminant of Matrix A: %.4f\n", det)
				}
			case 7:
				trace, err = matA.Trace()
				if err == nil {
					fmt.Printf("\nTrace of Matrix A: %.4f\n", trace)
				}
			case 8:
				result, err = matA.CofactorMatrix()
				if err == nil {
					fmt.Println("\nCofactor Matrix of A:")
					PrintMatrix(result)
				}
			case 9:
				result, err = matA.Adjugate()
				if err == nil {
					fmt.Println("\nAdjugate Matrix of A:")
					PrintMatrix(result)
				}
			case 10:
				result, err = matA.Inverse()
				if err == nil {
					fmt.Println("\nInverse Matrix of A:")
					PrintMatrix(result)
				}
			}
		case 11: // Create Identity Matrix
			size, err = readInt("Enter size for Identity Matrix (n x n): ")
			if err != nil {
				fmt.Println("Error reading size:", err)
				continue
			}
			result, err = IdentityMatrix(size)
			if err == nil {
				fmt.Println("\nIdentity Matrix:")
				PrintMatrix(result)
			}
		case 12: // Create Zero Matrix
			rows, err := readInt("Enter number of rows for Zero Matrix: ")
			if err != nil {
				fmt.Println("Error reading rows:", err)
				continue
			}
			cols, err := readInt("Enter number of columns for Zero Matrix: ")
			if err != nil {
				fmt.Println("Error reading columns:", err)
				continue
			}
			result, err = ZeroMatrix(rows, cols)
			if err == nil {
				fmt.Println("\nZero Matrix:")
				PrintMatrix(result)
			}
		case 13:
			fmt.Println("Exiting calculator. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please select a valid option (1-13).")
		}

		if err != nil {
			fmt.Println("Operation failed:", err)
		}
	}
}