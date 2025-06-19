package main

import (
	"fmt"
)

func solveSudoku(board *[9][9]int) bool {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if board[r][c] == 0 {
				for num := 1; num <= 9; num++ {
					if isValid(board, r, c, num) {
						board[r][c] = num
						if solveSudoku(board) {
							return true
						}
						board[r][c] = 0 // Backtrack
					}
				}
				return false // No number works for this cell
			}
		}
	}
	return true // All cells filled
}

func isValid(board *[9][9]int, row, col, num int) bool {
	// Check row
	for c := 0; c < 9; c++ {
		if board[row][c] == num {
			return false
		}
	}

	// Check column
	for r := 0; r < 9; r++ {
		if board[r][col] == num {
			return false
		}
	}

	// Check 3x3 box
	startRow := row - row%3
	startCol := col - col%3
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if board[startRow+r][startCol+c] == num {
				return false
			}
		}
	}

	return true
}

func printBoard(board *[9][9]int) {
	for r := 0; r < 9; r++ {
		if r%3 == 0 && r != 0 {
			fmt.Println("---------------------")
		}
		for c := 0; c < 9; c++ {
			if c%3 == 0 && c != 0 {
				fmt.Print("| ")
			}
			fmt.Printf("%d ", board[r][c])
		}
		fmt.Println()
	}
}

func main() {
	puzzle := [9][9]int{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	}

	fmt.Println("Sudoku Puzzle:")
	printBoard(&puzzle)
	fmt.Println("\nSolving...")

	if solveSudoku(&puzzle) {
		fmt.Println("\nSudoku Solved:")
		printBoard(&puzzle)
	} else {
		fmt.Println("\nNo solution exists for the given Sudoku puzzle.")
	}
}

// Additional implementation at 2025-06-19 00:28:39
package main

import (
	"fmt"
)

type Board [9][9]int

func (b Board) print() {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			fmt.Printf("%d ", b[r][c])
			if (c+1)%3 == 0 && c != 8 {
				fmt.Print("| ")
			}
		}
		fmt.Println()
		if (r+1)%3 == 0 && r != 8 {
			fmt.Println("---------------------")
		}
	}
}

func (b Board) findEmpty() (int, int, bool) {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if b[r][c] == 0 {
				return r, c, true
			}
		}
	}
	return 0, 0, false
}

func (b Board) isValid(num, row, col int) bool {
	for c := 0; c < 9; c++ {
		if b[row][c] == num {
			return false
		}
	}

	for r := 0; r < 9; r++ {
		if b[r][col] == num {
			return false
		}
	}

	startRow := (row / 3) * 3
	startCol := (col / 3) * 3
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if b[startRow+r][startCol+c] == num {
				return false
			}
		}
	}

	return true
}

func (b Board) isInitialPuzzleValid() bool {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if b[r][c] != 0 {
				val := b[r][c]
				b[r][c] = 0
				if !b.isValid(val, r, c) {
					return false
				}
				b[r][c] = val
			}
		}
	}
	return true
}

func (b *Board) solve() bool {
	row, col, found := b.findEmpty()
	if !found {
		return true
	}

	for num := 1; num <= 9; num++ {
		if b.isValid(num, row, col) {
			b.Board[row][col] = num
			if b.solve() {
				return true
			}
			b.Board[row][col] = 0
		}
	}
	return false
}

func main() {
	puzzle := Board{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	}

	fmt.Println("Initial Puzzle:")
	puzzle.print()
	fmt.Println()

	if !puzzle.isInitialPuzzleValid() {
		fmt.Println("Error: The initial puzzle is invalid (contains conflicts).")
		return
	}

	if puzzle.solve() {
		fmt.Println("Solved Sudoku:")
		puzzle.print()
	} else {
		fmt.Println("No solution exists for the given Sudoku puzzle.")
	}
}