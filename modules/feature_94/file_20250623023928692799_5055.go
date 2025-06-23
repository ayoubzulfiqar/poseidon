package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ROWS  = 10
	COLS  = 10
	MINES = 15
)

type Cell struct {
	hasMine       bool
	isRevealed    bool
	isFlagged     bool
	adjacentMines int
}

var board [ROWS][COLS]Cell
var gameOver bool
var gameWon bool
var firstMove bool = true
var revealedCells int

func main() {
	rand.Seed(time.Now().UnixNano())
	initializeBoard()
	reader := bufio.NewReader(os.Stdin)

	for !gameOver && !gameWon {
		printBoard()
		fmt.Print("Enter command (e.g., r 1 2 for reveal, f 3 4 for flag): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Fields(input)

		if len(parts) != 3 {
			fmt.Println("Invalid input format. Use 'r row col' or 'f row col'.")
			continue
		}

		cmd := parts[0]
		row, err1 := strconv.Atoi(parts[1])
		col, err2 := strconv.Atoi(parts[2])

		if err1 != nil || err2 != nil || !isValid(row, col) {
			fmt.Println("Invalid row or column. Please enter numbers within board limits.")
			continue
		}

		if firstMove {
			placeMines(row, col)
			calculateAdjacentMines()
			firstMove = false
		}

		switch cmd {
		case "r":
			revealCell(row, col)
		case "f":
			flagCell(row, col)
		default:
			fmt.Println("Unknown command. Use 'r' for reveal or 'f' for flag.")
		}

		checkWin()
	}

	printBoard()
	if gameWon {
		fmt.Println("Congratulations! You won!")
	} else if gameOver {
		fmt.Println("Game Over! You hit a mine.")
	}
}

func initializeBoard() {
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			board[r][c] = Cell{}
		}
	}
	gameOver = false
	gameWon = false
	revealedCells = 0
}

func placeMines(firstClickR, firstClickC int) {
	minesPlaced := 0
	for minesPlaced < MINES {
		r := rand.Intn(ROWS)
		c := rand.Intn(COLS)

		if (r >= firstClickR-1 && r <= firstClickR+1) && (c >= firstClickC-1 && c <= firstClickC+1) {
			continue
		}

		if !board[r][c].hasMine {
			board[r][c].hasMine = true
			minesPlaced++
		}
	}
}

func calculateAdjacentMines() {
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			if board[r][c].hasMine {
				continue
			}
			count := 0
			for dr := -1; dr <= 1; dr++ {
				for dc := -1; dc <= 1; dc++ {
					if dr == 0 && dc == 0 {
						continue
					}
					nr, nc := r+dr, c+dc
					if isValid(nr, nc) && board[nr][nc].hasMine {
						count++
					}
				}
			}
			board[r][c].adjacentMines = count
		}
	}
}

func printBoard() {
	fmt.Print("   ")
	for c := 0; c < COLS; c++ {
		fmt.Printf("%2d ", c)
	}
	fmt.Println()
	fmt.Print("  ")
	for c := 0; c < COLS; c++ {
		fmt.Print("---")
	}
	fmt.Println()

	for r := 0; r < ROWS; r++ {
		fmt.Printf("%2d|", r)
		for c := 0; c < COLS; c++ {
			cell := board[r][c]
			if gameOver && cell.hasMine {
				fmt.Print(" * ")
			} else if cell.isRevealed {
				if cell.hasMine {
					fmt.Print(" * ")
				} else if cell.adjacentMines == 0 {
					fmt.Print("   ")
				} else {
					fmt.Printf(" %d ", cell.adjacentMines)
				}
			} else if cell.isFlagged {
				fmt.Print(" F ")
			} else {
				fmt.Print(" # ")
			}
		}
		fmt.Println()
	}
}

func revealCell(r, c int) {
	if !isValid(r, c) || board[r][c].isRevealed || board[r][c].isFlagged {
		return
	}

	board[r][c].isRevealed = true
	revealedCells++

	if board[r][c].hasMine {
		gameOver = true
		return
	}

	if board[r][c].adjacentMines == 0 {
		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				if dr == 0 && dc == 0 {
					continue
				}
				revealCell(r+dr, c+dc)
			}
		}
	}
}

func flagCell(r, c int) {
	if !isValid(r, c) || board[r][c].isRevealed {
		return
	}
	board[r][c].isFlagged = !board[r][c].isFlagged
}

func checkWin() {
	totalNonMineCells := (ROWS * COLS) - MINES
	if revealedCells == totalNonMineCells {
		gameWon = true
	}
}

func isValid(r, c int) bool {
	return r >= 0 && r < ROWS && c >= 0 && c < COLS
}

// Additional implementation at 2025-06-23 02:40:38
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Cell struct {
	IsMine        bool
	IsRevealed    bool
	IsFlagged     bool
	AdjacentMines int
}

type Game struct {
	Board         [][]Cell
	Rows          int
	Cols          int
	Mines         int
	RevealedCount int
	StartTime     time.Time
	IsGameOver    bool
	HasWon        bool
}

const (
	SmallRows  = 8
	SmallCols  = 8
	SmallMines = 10

	MediumRows  = 12
	MediumCols  = 16
	MediumMines = 30

	LargeRows  = 16
	LargeCols  = 30
	LargeMines = 99
)

func NewGame(rows, cols, mines int) *Game {
	game := &Game{
		Rows:          rows,
		Cols:          cols,
		Mines:         mines,
		Board:         make([][]Cell, rows),
		RevealedCount: 0,
		IsGameOver:    false,
		HasWon:        false,
	}

	for r := 0; r < rows; r++ {
		game.Board[r] = make([]Cell, cols)
	}

	game.placeMines()
	game.calculateAdjacentMines()

	return game
}

func (g *Game) placeMines() {
	rand.Seed(time.Now().UnixNano())
	minesPlaced := 0
	for minesPlaced < g.Mines {
		r := rand.Intn(g.Rows)
		c := rand.Intn(g.Cols)
		if !g.Board[r][c].IsMine {
			g.Board[r][c].IsMine = true
			minesPlaced++
		}
	}
}

func (g *Game) calculateAdjacentMines() {
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {
			if !g.Board[r][c].IsMine {
				count := 0
				for dr := -1; dr <= 1; dr++ {
					for dc := -1; dc <= 1; dc++ {
						if dr == 0 && dc == 0 {
							continue
						}
						nr, nc := r+dr, c+dc
						if nr >= 0 && nr < g.Rows && nc >= 0 && nc < g.Cols && g.Board[nr][nc].IsMine {
							count++
						}
					}
				}
				g.Board[r][c].AdjacentMines = count
			}
		}
	}
}

func (g *Game) DisplayBoard() {
	fmt.Print("   ")
	for c := 0; c < g.Cols; c++ {
		fmt.Printf("%2d ", c)
	}
	fmt.Println()
	fmt.Print("  +" + strings.Repeat("---", g.Cols) + "+")
	fmt.Println()

	for r := 0; r < g.Rows; r++ {
		fmt.Printf("%2d |", r)
		for c := 0; c < g.Cols; c++ {
			cell := g.Board[r][c]
			if cell.IsRevealed {
				if cell.IsMine {
					fmt.Print(" * ")
				} else if cell.AdjacentMines == 0 {
					fmt.Print("   ")
				} else {
					fmt.Printf(" %d ", cell.AdjacentMines)
				}
			} else if cell.IsFlagged {
				fmt.Print(" F ")
			} else {
				fmt.Print(" # ")
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +" + strings.Repeat("---", g.Cols) + "+")
	fmt.Printf("Mines: %d | Covered: %d | Flags: %d\n", g.Mines, (g.Rows*g.Cols)-g.RevealedCount, g.countFlags())
	if g.StartTime.IsZero() {
		fmt.Println("Time: 0s")
	} else {
		fmt.Printf("Time: %s\n", time.Since(g.StartTime).Round(time.Second))
	}
}

func (g *Game) countFlags() int {
	count := 0
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {
			if g.Board[r][c].IsFlagged {
				count++
			}
		}
	}
	return count
}

func (g *Game) RevealCell(r, c int) bool {
	if r < 0 || r >= g.Rows || c < 0 || c >= g.Cols {
		fmt.Println("Invalid coordinates.")
		return true
	}
	cell := &g.Board[r][c]

	if cell.IsRevealed {
		fmt.Println("Cell already revealed.")
		return true
	}
	if cell.IsFlagged {
		fmt.Println("Cell is flagged. Unflag it first to reveal.")
		return true
	}

	if g.StartTime.IsZero() {
		g.StartTime = time.Now()
	}

	cell.IsRevealed = true
	g.RevealedCount++

	if cell.IsMine {
		g.IsGameOver = true
		g.HasWon = false
		return false
	}

	if cell.AdjacentMines == 0 {
		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				if dr == 0 && dc == 0 {
					continue
				}
				nr, nc := r+dr, c+dc
				if nr >= 0 && nr < g.Rows && nc >= 0 && nc < g.Cols && !g.Board[nr][nc].IsRevealed && !g.Board[nr][nc].IsFlagged {
					g.RevealCell(nr, nc)
				}
			}
		}
	}

	if g.RevealedCount == (g.Rows*g.Cols)-g.Mines {
		g.IsGameOver = true
		g.HasWon = true
		return false
	}

	return true
}

func (g *Game) ToggleFlag(r, c int) {
	if r < 0 || r >= g.Rows || c < 0 || c >= g.Cols {
		fmt.Println("Invalid coordinates.")
		return
	}
	cell := &g.Board[r][c]

	if cell.IsRevealed {
		fmt.Println("Cannot flag a revealed cell.")
		return
	}

	cell.IsFlagged = !cell.IsFlagged
	if cell.IsFlagged {
		fmt.Printf("Cell (%d, %d) flagged.\n", r, c)
	} else {
		fmt.Printf("Cell (%d, %d) unflagged.\n", r, c)
	}
}

func parseInput(input string) (command string, r, c int, err error) {
	parts := strings.Fields(strings.ToLower(input))
	if len(parts) < 3 {
		return "", 0, 0, fmt.Errorf("invalid command format. Use 'reveal R C' or 'flag R C'")
	}

	cmd := parts[0]
	rowStr := parts[1]
	colStr := parts[2]

	r, err = strconv.Atoi(rowStr)
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid row number: %s", rowStr)
	}
	c, err = strconv.Atoi(colStr)
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid column number: %s", colStr)
	}

	if cmd != "reveal" && cmd != "flag" {
		return "", 0, 0, fmt.Errorf("unknown command: %s. Use 'reveal' or 'flag'", cmd)
	}

	return cmd, r, c, nil
}

func gameLoop(game *Game) {
	reader := bufio.NewReader(os.Stdin)

	for !game.IsGameOver {
		game.DisplayBoard()
		fmt.Print("Enter command (e.g., 'reveal 0 0', 'flag 1 2'): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		cmd, r, c, err := parseInput(input)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		switch cmd {
		case "reveal":
			if !game.RevealCell(r, c) {
				break
			}
		case "flag":
			game.ToggleFlag(r, c)
		default:
			fmt.Println("Unknown command. Please use 'reveal' or 'flag'.")
		}
	}

	game.DisplayBoard()
	if game.HasWon {
		fmt.Printf("\nCONGRATULATIONS! You won in %s!\n", time.Since(game.StartTime).Round(time.Second))
	} else {
		fmt.Printf("\nGAME OVER! You hit a mine. Better luck next time!\n")
	}
}

func chooseDifficulty() (rows, cols, mines int) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nChoose difficulty:")
		fmt.Println("1. Small (8x8, 10 mines)")
		fmt.Println("2. Medium (12x16, 30 mines)")
		fmt.Println("3. Large (16x30, 99 mines)")
		fmt.Print("Enter choice (1, 2, or 3): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			return SmallRows, SmallCols, SmallMines
		case "2":
			return MediumRows, MediumCols, MediumMines
		case "3":
			return LargeRows, LargeCols, LargeMines
		default:
			fmt.Println("Invalid choice. Please enter 1, 2, or 3.")
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		rows, cols, mines := chooseDifficulty()
		game := NewGame(rows, cols, mines)
		gameLoop(game)

		fmt.Print("Play again? (y/n): ")
		playAgainInput, _ := reader.ReadString('\n')
		playAgainInput = strings.TrimSpace(strings.ToLower(playAgainInput))
		if playAgainInput != "y" {
			break
		}
	}
	fmt.Println("Thanks for playing Minesweeper!")
}

// Additional implementation at 2025-06-23 02:41:57
package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type CellState int

const (
	Covered CellState = iota
	Revealed
	Flagged
)

type CellContent int

const (
	Mine  CellContent = -1
	Empty CellContent = 0
)

type Cell struct {
	Content CellContent
	State   CellState
}

type Board struct {
	Cells      [][]Cell
	Rows       int
	Cols       int
	Mines      int
	FlagsUsed  int
	GameOver   bool
	GameWon    bool
	FirstClick bool
}

func NewBoard(rows, cols, mines int) *Board {
	board := &Board{
		Rows:       rows,
		Cols:       cols,
		Mines:      mines,
		Cells:      make([][]Cell, rows),
		FirstClick: true,
	}
	for r := 0; r < rows; r++ {
		board.Cells[r] = make([]Cell, cols)
		for c := 0; c < cols; c++ {
			board.Cells[r][c] = Cell{Content: Empty, State: Covered}
		}
	}
	return board
}

func (b *Board) PlaceMines(firstClickR, firstClickC int) {
	rand.Seed(time.Now().UnixNano())
	minesPlaced := 0
	for minesPlaced < b.Mines {
		r := rand.Intn(b.Rows)
		c := rand.Intn(b.Cols)

		isNearFirstClick := false
		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				if r == firstClickR+dr && c == firstClickC+dc {
					isNearFirstClick = true
					break
				}
			}
			if isNearFirstClick {
				break
			}
		}

		if b.Cells[r][c].Content != Mine && !isNearFirstClick {
			b.Cells[r][c].Content = Mine
			minesPlaced++
		}
	}
	b.CalculateNumbers()
}

func (b *Board) CalculateNumbers() {
	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			if b.Cells[r][c].Content == Mine {
				continue
			}

			mineCount := 0
			for dr := -1; dr <= 1; dr++ {
				for dc := -1; dc <= 1; dc++ {
					if dr == 0 && dc == 0 {
						continue
					}
					nr, nc := r+dr, c+dc
					if b.IsValidCoord(nr, nc) && b.Cells[nr][nc].Content == Mine {
						mineCount++
					}
				}
			}
			b.Cells[r][c].Content = CellContent(mineCount)
		}
	}
}

func (b *Board) IsValidCoord(r, c int) bool {
	return r >= 0 && r < b.Rows && c >= 0 && c < b.Cols
}

func (b *Board) PrintBoard() {
	fmt.Print("   ")
	for c := 0; c < b.Cols; c++ {
		fmt.Printf("%2d ", c)
	}
	fmt.Println()

	fmt.Print("  +")
	for c := 0; c < b.Cols; c++ {
		fmt.Print("---")
	}
	fmt.Println("+")

	for r := 0; r < b.Rows; r++ {
		fmt.Printf("%2d|", r)
		for c := 0; c < b.Cols; c++ {
			cell := b.Cells[r][c]
			switch cell.State {
			case Covered:
				fmt.Print(" # ")
			case Flagged:
				fmt.Print(" F ")
			case Revealed:
				if cell.Content == Mine {
					fmt.Print(" * ")
				} else if cell.Content == Empty {
					fmt.Print("   ")
				} else {
					fmt.Printf(" %d ", cell.Content)
				}
			}
		}
		fmt.Println("|")
	}

	fmt.Print("  +")
	for c := 0; c < b.Cols; c++ {
		fmt.Print("---")
	}
	fmt.Println("+")
	fmt.Printf("Mines: %d, Flags Used: %d\n", b.Mines, b.FlagsUsed)
}

func (b *Board) RevealCell(r, c int) {
	if !b.IsValidCoord(r, c) || b.Cells[r][c].State == Revealed || b.Cells[r][c].State == Flagged {
		return
	}

	if b.FirstClick {
		b.PlaceMines(r, c)
		b.FirstClick = false
	}

	b.Cells[r][c].State = Revealed

	if b.Cells[r][c].Content == Mine {
		b.GameOver = true
		fmt.Println("BOOM! Game Over.")
		return
	}

	if b.Cells[r][c].Content == Empty {
		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				if dr == 0 && dc == 0 {
					continue
				}
				b.RevealCell(r+dr, c+dc)
			}
		}
	}
}

func (b *Board) FlagCell(r, c int) {
	if !b.IsValidCoord(r, c) || b.Cells[r][c].State == Revealed {
		return
	}

	if b.Cells[r][c].State == Flagged {
		b.Cells[r][c].State = Covered
		b.FlagsUsed--
	} else {
		if b.FlagsUsed < b.Mines {
			b.Cells[r][c].State = Flagged
			b.FlagsUsed++
		} else {
			fmt.Println("Cannot place more flags than total mines.")
		}
	}
}

func (b *Board) CheckWin() bool {
	revealedNonMines := 0
	totalNonMines := (b.Rows * b.Cols) - b.Mines

	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			cell := b.Cells[r][c]
			if cell.State == Revealed && cell.Content != Mine {
				revealedNonMines++
			}
		}
	}

	if revealedNonMines == totalNonMines {
		b.GameWon = true
		b.GameOver = true
		fmt.Println("Congratulations! You won!")
		return true
	}
	return false
}

func main() {
	var input string

	fmt.Println("Welcome to Go Minesweeper!")

	rows, cols, mines := 10, 10, 15
	fmt.Printf("Enter board dimensions (rows cols mines, e.g., 10 10 15). Press Enter for default %d %d %d: ", rows, cols, mines)
	fmt.Scanln(&input)

	parts := strings.Fields(input)
	if len(parts) == 3 {
		r, errR := strconv.Atoi(parts[0])
		c, errC := strconv.Atoi(parts[1])
		m, errM := strconv.Atoi(parts[2])
		if errR == nil && errC == nil && errM == nil && r > 0 && c > 0 && m > 0 && m < r*c {
			rows, cols, mines = r, c, m
		} else {
			fmt.Println("Invalid input for dimensions. Using default.")
		}
	} else if len(parts) > 0 {
		fmt.Println("Invalid input format for dimensions. Using default.")
	}

	board := NewBoard(rows, cols, mines)

	for !board.GameOver {
		board.PrintBoard()
		fmt.Print("Enter action (r for reveal, f for flag) and coordinates (row col), e.g., r 0 0 or f 2 3: ")
		fmt.Scanln(&input)

		parts = strings.Fields(input)
		if len(parts) != 3 {
			fmt.Println("Invalid input format. Please use 'r row col' or 'f row col'.")
			continue
		}

		action := strings.ToLower(parts[0])
		r, errR := strconv.Atoi(parts[1])
		c, errC := strconv.Atoi(parts[2])

		if errR != nil || errC != nil || !board.IsValidCoord(r, c) {
			fmt.Println("Invalid row or column. Please enter valid numbers within board limits.")
			continue
		}

		switch action {
		case "r":
			board.RevealCell(r, c)
		case "f":
			board.FlagCell(r, c)
		default:
			fmt.Println("Invalid action. Use 'r' for reveal or 'f' for flag.")
		}

		if !board.GameOver {
			board.CheckWin()
		}
	}

	// Game over, print final board revealing all mines if lost
	if !board.GameWon {
		for r := 0; r < board.Rows; r++ {
			for c := 0; c < board.Cols; c++ {
				if board.Cells[r][c].Content == Mine {
					board.Cells[r][c].State = Revealed // Reveal all mines on loss
				}
			}
		}
	}
	board.PrintBoard()

	if board.GameWon {
		fmt.Println("You successfully cleared the minefield!")
	} else {
		fmt.Println("Better luck next time!")
	}
}

// Additional implementation at 2025-06-23 02:42:54
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Cell struct {
	hasMine       bool
	isRevealed    bool
	isFlagged     bool
	adjacentMines int
}

type GameState int

const (
	Playing GameState = iota
	Won
	Lost
)

type Board struct {
	grid          [][]Cell
	rows          int
	cols          int
	mines         int
	revealedCount int
	flagsUsed     int
}

const (
	EasyRows    = 9
	EasyCols    = 9
	EasyMines   = 10
	MediumRows  = 16
	MediumCols  = 16
	MediumMines = 40
	HardRows    = 16
	HardCols    = 30
	HardMines   = 99
)

const (
	UnrevealedChar = "■"
	FlagChar       = "F"
	MineChar       = "M"
	EmptyChar      = " "
)

var reader = bufio.NewReader(os.Stdin)

type Game struct {
	board Board
	state GameState
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func printInstructions() {
	fmt.Println("Welcome to Go Minesweeper!")
	fmt.Println("Instructions:")
	fmt.Println("  To reveal a cell: type 'r <row> <col>' (e.g., 'r 5 3')")
	fmt.Println("  To flag/unflag a cell: type 'f <row> <col>' (e.g., 'f 2 7')")
	fmt.Println("  Rows and columns are 0-indexed.")
	fmt.Println("Good luck!\n")
}

func newGame(rows, cols, mines int) *Game {
	board := Board{
		rows:  rows,
		cols:  cols,
		mines: mines,
		grid:  make([][]Cell, rows),
	}
	for r := 0; r < rows; r++ {
		board.grid[r] = make([]Cell, cols)
	}

	return &Game{
		board: board,
		state: Playing,
	}
}

func (b *Board) placeMines(firstClickR, firstClickC int) {
	rand.Seed(time.Now().UnixNano())
	minesPlaced := 0
	for minesPlaced < b.mines {
		r := rand.Intn(b.rows)
		c := rand.Intn(b.cols)

		// Ensure the cell is not already a mine and not in the first click area (3x3 around it)
		if !b.grid[r][c].hasMine && !((r >= firstClickR-1 && r <= firstClickR+1) && (c >= firstClickC-1 && c <= firstClickC+1)) {
			b.grid[r][c].hasMine = true
			minesPlaced++
		}
	}
}

func (b *Board) calculateAdjacentMines() {
	for r := 0; r < b.rows; r++ {
		for c := 0; c < b.cols; c++ {
			if b.grid[r][c].hasMine {
				continue
			}

			count := 0
			for dr := -1; dr <= 1; dr++ {
				for dc := -1; dc <= 1; dc++ {
					if dr == 0 && dc == 0 {
						continue
					}
					nr, nc := r+dr, c+dc
					if nr >= 0 && nr < b.rows && nc >= 0 && nc < b.cols && b.grid[nr][nc].hasMine {
						count++
					}
				}
			}
			b.grid[r][c].adjacentMines = count
		}
	}
}

func (b *Board) displayBoard() {
	fmt.Print("   ")
	for c := 0; c < b.cols; c++ {
		fmt.Printf("%2d ", c)
	}
	fmt.Println()

	fmt.Print("   ")
	for c := 0; c < b.cols; c++ {
		fmt.Print("---")
	}
	fmt.Println()

	for r := 0; r < b.rows; r++ {
		fmt.Printf("%2d |", r)
		for c := 0; c < b.cols; c++ {
			cell := b.grid[r][c]
			if cell.isRevealed {
				if cell.hasMine {
					fmt.Printf(" %s ", MineChar)
				} else if cell.adjacentMines == 0 {
					fmt.Printf(" %s ", EmptyChar)
				} else {
					fmt.Printf(" %d ", cell.adjacentMines)
				}
			} else if cell.isFlagged {
				fmt.Printf(" %s ", FlagChar)
			} else {
				fmt.Printf(" %s ", UnrevealedChar)
			}
		}
		fmt.Println()
	}
	fmt.Printf("Mines: %d, Flags Used: %d\n", b.mines, b.flagsUsed)
}

func getInput() (int, int, string, error) {
	fmt.Print("Enter command (e.g., 'r 0 0' or 'f 1 2'): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	parts := strings.Fields(input)

	if len(parts) != 3 {
		return 0, 0, "", fmt.Errorf("invalid input format. Use 'r <row> <col>' or 'f <row> <col>'")
	}

	action := strings.ToLower(parts[0])
	if action != "r" && action != "f" {
		return 0, 0, "", fmt.Errorf("invalid action. Use 'r' for reveal or 'f' for flag")
	}

	row, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, "", fmt.Errorf("invalid row number: %v", err)
	}
	col, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, "", fmt.Errorf("invalid column number: %v", err)
	}

	return row, col, action, nil
}

func (g *Game) revealCell(r, c int) bool {
	if r < 0 || r >= g.board.rows || c < 0 || c >= g.board.cols {
		return false
	}

	cell := &g.board.grid[r][c]

	if cell.isRevealed || cell.isFlagged {
		return false
	}

	cell.isRevealed = true
	g.board.revealedCount++

	if cell.hasMine {
		return true
	}

	if cell.adjacentMines == 0 {
		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				if dr == 0 && dc == 0 {
					continue
				}
				g.revealCell(r+dr, c+dc)
			}
		}
	}
	return false
}

func (g *Game) flagCell(r, c int) {
	if r < 0 || r >= g.board.rows || c < 0 || c >= g.board.cols {
		fmt.Println("Invalid coordinates.")
		return
	}

	cell := &g.board.grid[r][c]

	if cell.isRevealed {
		fmt.Println("Cannot flag a revealed cell.")
		return
	}

	cell.isFlagged = !cell.isFlagged
	if cell.isFlagged {
		g.board.flagsUsed++
	} else {
		g.board.flagsUsed--
	}
}

func (g *Game) checkWin() bool {
	nonMineCells := (g.board.rows * g.board.cols) - g.board.mines
	if g.board.revealedCount == nonMineCells {
		g.state = Won
		return true
	}
	return false
}

func main() {
	for {
		clearScreen()
		printInstructions()

		fmt.Println("Choose difficulty:")
		fmt.Println("1. Easy (9x9, 10 mines)")
		fmt.Println("2. Medium (16x16, 40 mines)")
		fmt.Println("3. Hard (16x30, 99 mines)")
		fmt.Print("Enter choice (1-3): ")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choice, _ := strconv.Atoi(choiceStr)

		var rows, cols, mines int
		switch choice {
		case 1:
			rows, cols, mines = EasyRows, EasyCols, EasyMines
		case 2:
			rows, cols, mines = MediumRows, MediumCols, MediumMines
		case 3:
			rows, cols, mines = HardRows, HardCols, HardMines
		default:
			fmt.Println("Invalid choice, defaulting to Easy.")
			rows, cols, mines = EasyRows, EasyCols, EasyMines
		}

		game := newGame(rows, cols, mines)
		firstClick := true

		for game.state == Playing {
			clearScreen()
			game.board.displayBoard()

			row, col, action, err := getInput()
			if err != nil {
				fmt.Println("Error:", err)
				time.Sleep(2 * time.Second)
				continue
			}

			if row < 0 || row >= game.board.rows || col < 0 || col >= game.board.cols {
				fmt.Println("Coordinates out of bounds.")
				time.Sleep(2 * time.Second)
				continue
			}

			if firstClick {
				game.board.placeMines(row, col)
				game.board.calculateAdjacentMines()
				firstClick = false
			}

			if action == "r" {
				if game.board.grid[row][col].isFlagged {
					fmt.Println("Cannot reveal a flagged cell. Unflag it first.")
					time.Sleep(2 * time.Second)
					continue
				}
				if game.revealCell(row, col) {
					game.state = Lost
				}
			} else if action == "f" {
				game.flagCell(row, col)
			}

			if game.state == Playing {
				game.checkWin()
			}
		}

		clearScreen()
		game.board.displayBoard()
		if game.state == Won {
			fmt.Println("\nCongratulations! You won!")
		} else {
			fmt.Println("\nGame Over! You hit a mine!")
			// Reveal all mines on loss
			for r := 0; r < game.board.rows; r++ {
				for c := 0; c < game.board.cols; c++ {
					if game.board.grid[r][c].hasMine {
						game.board.grid[r][c].isRevealed = true
					}
				}
			}
			game.board.displayBoard() // Display board with mines revealed
		}

		fmt.Print("Play again? (y/n): ")
		playAgain, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(playAgain)) != "y" {
			break
		}
	}
	fmt.Println("Thanks for playing!")
}