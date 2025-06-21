package main

import (
	"fmt"
)

func generatePascalsTriangle(numRows int) [][]int {
	if numRows <= 0 {
		return [][]int{}
	}

	triangle := make([][]int, numRows)

	for i := 0; i < numRows; i++ {
		triangle[i] = make([]int, i+1)
		triangle[i][0] = 1
		if i > 0 {
			triangle[i][i] = 1
		}

		for j := 1; j < i; j++ {
			triangle[i][j] = triangle[i-1][j-1] + triangle[i-1][j]
		}
	}
	return triangle
}

func main() {
	numRows := 7
	pascalsTriangle := generatePascalsTriangle(numRows)

	for _, row := range pascalsTriangle {
		for _, val := range row {
			fmt.Printf("%d ", val)
		}
		fmt.Println()
	}
}

// Additional implementation at 2025-06-21 04:24:50
package main

import "fmt"

func GeneratePascalTriangle(numRows int) [][]int {
	if numRows <= 0 {
		return [][]int{}
	}

	triangle := make([][]int, numRows)

	for i := 0; i < numRows; i++ {
		triangle[i] = make([]int, i+1)
		triangle[i][0] = 1
		triangle[i][i] = 1

		for j := 1; j < i; j++ {
			triangle[i][j] = triangle[i-1][j-1] + triangle[i-1][j]
		}
	}
	return triangle
}

func GetPascalRow(rowIndex int) []int {
	if rowIndex < 0 {
		return []int{}
	}

	row := make([]int, rowIndex+1)
	row[0] = 1

	for k := 1; k <= rowIndex; k++ {
		row[k] = (row[k-1] * (rowIndex - k + 1)) / k
	}
	return row
}

func GetPascalElement(row, col int) int {
	if col < 0 || col > row || row < 0 {
		return 0
	}

	if col == 0 || col == row {
		return 1
	}

	if col > row/2 {
		col = row - col
	}

	res := 1
	for i := 0; i < col; i++ {
		res = res * (row - i) / (i + 1)
	}
	return res
}

func main() {
	fmt.Println("Pascal's Triangle (5 rows):")
	triangle := GeneratePascalTriangle(5)
	for _, row := range triangle {
		fmt.Println(row)
	}

	fmt.Println("\nSpecific Row (Row 4, 0-indexed):")
	row4 := GetPascalRow(4)
	fmt.Println(row4)

	fmt.Println("\nSpecific Element (Row 5, Col 2, 0-indexed):")
	element := GetPascalElement(5, 2)
	fmt.Println(element)

	fmt.Println("\nSpecific Element (Row 0, Col 0):")
	element00 := GetPascalElement(0, 0)
	fmt.Println(element00)

	fmt.Println("\nSpecific Element (Row 3, Col 3):")
	element33 := GetPascalElement(3, 3)
	fmt.Println(element33)

	fmt.Println("\nSpecific Element (Invalid: Row -1, Col 0):")
	invalidElement := GetPascalElement(-1, 0)
	fmt.Println(invalidElement)

	fmt.Println("\nSpecific Row (Invalid: Row -1):")
	invalidRow := GetPascalRow(-1)
	fmt.Println(invalidRow)

	fmt.Println("\nPascal's Triangle (Invalid: 0 rows):")
	emptyTriangle := GeneratePascalTriangle(0)
	fmt.Println(emptyTriangle)
}

// Additional implementation at 2025-06-21 04:25:20
package main

import (
	"fmt"
)

func generatePascalsTriangle(numRows int) [][]int {
	if numRows <= 0 {
		return [][]int{}
	}

	triangle := make([][]int, numRows)

	for i := 0; i < numRows; i++ {
		triangle[i] = make([]int, i+1)
		triangle[i][0] = 1
		triangle[i][i] = 1

		for j := 1; j < i; j++ {
			triangle[i][j] = triangle[i-1][j-1] + triangle[i-1][j]
		}
	}
	return triangle
}

func getRow(rowIndex int) []int {
	if rowIndex < 0 {
		return []int{}
	}

	row := make([]int, rowIndex+1)
	row[0] = 1

	for i := 1; i <= rowIndex; i++ {
		for j := i; j > 0; j-- {
			row[j] += row[j-1]
		}
	}
	return row
}

func getElement(row, col int) int {
	if col < 0 || col > row {
		return 0
	}
	if col == 0 || col == row {
		return 1
	}
	if col > row/2 {
		col = row - col
	}

	res := 1
	for i := 0; i < col; i++ {
		res = res * (row - i) / (i + 1)
	}
	return res
}

func main() {
	fmt.Println("--- Pascal's Triangle (5 rows) ---")
	triangle5 := generatePascalsTriangle(5)
	for _, row := range triangle5 {
		fmt.Println(row)
	}

	fmt.Println("\n--- Pascal's Triangle (7 rows) ---")
	triangle7 := generatePascalsTriangle(7)
	for _, row := range triangle7 {
		fmt.Println(row)
	}

	fmt.Println("\n--- Specific Row (Row 4, 0-indexed) ---")
	row4 := getRow(4)
	fmt.Println(row4)

	fmt.Println("\n--- Specific Row (Row 0, 0-indexed) ---")
	row0 := getRow(0)
	fmt.Println(row0)

	fmt.Println("\n--- Specific Row (Row 6, 0-indexed) ---")
	row6 := getRow(6)
	fmt.Println(row6)

	fmt.Println("\n--- Specific Element (Row 4, Col 2) ---")
	elem4_2 := getElement(4, 2)
	fmt.Println(elem4_2)

	fmt.Println("\n--- Specific Element (Row 6, Col 3) ---")
	elem6_3 := getElement(6, 3)
	fmt.Println(elem6_3)

	fmt.Println("\n--- Specific Element (Row 5, Col 0) ---")
	elem5_0 := getElement(5, 0)
	fmt.Println(elem5_0)

	fmt.Println("\n--- Specific Element (Row 5, Col 5) ---")
	elem5_5 := getElement(5, 5)
	fmt.Println(elem5_5)

	fmt.Println("\n--- Specific Element (Invalid: Row 3, Col 4) ---")
	elem3_4 := getElement(3, 4)
	fmt.Println(elem3_4)
}

// Additional implementation at 2025-06-21 04:26:42
package main

import (
	"fmt"
	"strings"
)

func generatePascalsTriangle(numRows int) [][]int {
	if numRows <= 0 {
		return [][]int{}
	}

	triangle := make([][]int, numRows)

	for i := 0; i < numRows; i++ {
		triangle[i] = make([]int, i+1)
		triangle[i][0] = 1
		if i > 0 {
			triangle[i][i] = 1
		}

		for j := 1; j < i; j++ {
			triangle[i][j] = triangle[i-1][j-1] + triangle[i-1][j]
		}
	}
	return triangle
}

func printPascalsTriangle(triangle [][]int) {
	if len(triangle) == 0 {
		fmt.Println("No triangle to print.")
		return
	}

	maxWidth := 0
	if len(triangle) > 0 {
		lastRow := triangle[len(triangle)-1]
		for _, val := range lastRow {
			width := len(fmt.Sprintf("%d", val))
			if width > maxWidth {
				maxWidth = width
			}
		}
	}

	lastRowWidth := 0
	if len(triangle) > 0 {
		lr := triangle[len(triangle)-1]
		lastRowWidth = (len(lr) * maxWidth) + (len(lr) - 1)
		if len(lr) == 1 {
			lastRowWidth = maxWidth
		}
	}

	for _, row := range triangle {
		rowStr := []string{}
		for _, val := range row {
			rowStr = append(rowStr, fmt.Sprintf("%*d", maxWidth, val))
		}
		
		totalRowWidth := (len(row) * maxWidth) + (len(row) - 1)
		if len(row) == 1 {
			totalRowWidth = maxWidth
		}
		
		padding := (lastRowWidth - totalRowWidth) / 2
		fmt.Printf("%s%s\n", strings.Repeat(" ", padding), strings.Join(rowStr, " "))
	}
}

func getPascalsElement(row, col int) (int, error) {
	if col < 0 || col > row {
		return 0, fmt.Errorf("invalid column %d for row %d", col, row)
	}
	if row < 0 {
		return 0, fmt.Errorf("invalid row %d", row)
	}

	if col > row/2 {
		col = row - col
	}

	res := 1
	for i := 0; i < col; i++ {
		res = res * (row - i) / (i + 1)
	}
	return res, nil
}

func printSpecificRow(rowNum int, triangle [][]int) {
	if rowNum < 0 || rowNum >= len(triangle) {
		fmt.Printf("Row %d is out of bounds for the generated triangle.\n", rowNum)
		return
	}
	fmt.Printf("Row %d: %v\n", rowNum, triangle[rowNum])
}

func main() {
	numRows := 7

	fmt.Println("Generating Pascal's Triangle:")
	pascalsTriangle := generatePascalsTriangle(numRows)
	printPascalsTriangle(pascalsTriangle)

	fmt.Println("\n--- Additional Functionality ---")

	targetRow := 5
	targetCol := 2
	element, err := getPascalsElement(targetRow, targetCol)
	if err != nil {
		fmt.Printf("Error getting element (%d, %d): %v\n", targetRow, targetCol, err)
	} else {
		fmt.Printf("Element at Row %d, Column %d (0-indexed): %d\n", targetRow, targetCol, element)
	}

	rowToPrint := 3
	printSpecificRow(rowToPrint, pascalsTriangle)

	rowToPrint = 0
	printSpecificRow(rowToPrint, pascalsTriangle)

	rowToPrint = 6
	printSpecificRow(rowToPrint, pascalsTriangle)

	_, err = getPascalsElement(3, 4)
	if err != nil {
		fmt.Printf("Error getting element (3, 4): %v\n", err)
	}
}