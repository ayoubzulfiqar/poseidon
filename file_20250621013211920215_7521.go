package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Point represents a coordinate in the maze grid
type Point struct {
	R, C int
}

// Maze struct holds the maze grid, dimensions, start, and end points
type Maze struct {
	Grid   [][]rune // '#' for wall, ' ' for path, 'S' for start, 'E' for end, 'X' for solution
	Width  int      // Logical width of cells
	Height int      // Logical height of cells
	Start  Point
	End    Point
}

// GenCell is used internally for maze generation
type GenCell struct {
	Visited bool
	Walls   [4]bool // 0:Top, 1:Right, 2:Bottom, 3:Left
}

const (
	WallTop = iota
	WallRight
	WallBottom
	WallLeft
)

// NewMaze initializes a new Maze struct
func NewMaze(width, height int) *Maze {
	m := &Maze{
		Width:  width,
		Height: height,
	}
	m.Grid = make([][]rune, 2*height+1)
	for i := range m.Grid {
		m.Grid[i] = make([]rune, 2*width+1)
		for j := range m.Grid[i] {
			m.Grid[i][j] = '#' // Initialize all as walls
		}
	}
	return m
}

// GenerateMaze creates a maze using Depth-First Search (DFS) algorithm
func (m *Maze) GenerateMaze() {
	genGrid := make([][]GenCell, m.Height)
	for r := range genGrid {
		genGrid[r] = make([]GenCell, m.Width)
		for c := range genGrid[r] {
			genGrid[r][c].Walls = [4]bool{true, true, true, true} // All walls initially present
		}
	}

	// Stack for DFS
	stack := []Point{}

	// Start from a random cell
	startR, startC := rand.Intn(m.Height), rand.Intn(m.Width)
	current := Point{startR, startC}
	genGrid[current.R][current.C].Visited = true
	stack = append(stack, current)

	// Directions: Up, Down, Left, Right
	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}

	for len(stack) > 0 {
		current = stack[len(stack)-1] // Peek
		
		// Find unvisited neighbors
		neighbors := []Point{}
		validDirections := []int{}
		for i := 0; i < 4; i++ {
			nR, nC := current.R+dr[i], current.C+dc[i]
			if nR >= 0 && nR < m.Height && nC >= 0 && nC < m.Width && !genGrid[nR][nC].Visited {
				neighbors = append(neighbors, Point{nR, nC})
				validDirections = append(validDirections, i)
			}
		}

		if len(neighbors) > 0 {
			// Choose a random unvisited neighbor
			idx := rand.Intn(len(neighbors))
			next := neighbors[idx]
			direction := validDirections[idx]

			// Remove wall between current and next cell
			switch direction {
			case WallTop: // next is above current
				genGrid[current.R][current.C].Walls[WallTop] = false
				genGrid[next.R][next.C].Walls[WallBottom] = false
			case WallBottom: // next is below current
				genGrid[current.R][current.C].Walls[WallBottom] = false
				genGrid[next.R][next.C].Walls[WallTop] = false
			case WallLeft: // next is left of current
				genGrid[current.R][current.C].Walls[WallLeft] = false
				genGrid[next.R][next.C].Walls[WallRight] = false
			case WallRight: // next is right of current
				genGrid[current.R][current.C].Walls[WallRight] = false
				genGrid[next.R][next.C].Walls[WallLeft] = false
			}

			genGrid[next.R][next.C].Visited = true
			stack = append(stack, next)
		} else {
			// No unvisited neighbors, backtrack
			stack = stack[:len(stack)-1]
		}
	}

	// Convert genGrid to Maze.Grid for display
	for r := 0; r < m.Height; r++ {
		for c := 0; c < m.Width; c++ {
			// Carve out the cell itself
			m.Grid[2*r+1][2*c+1] = ' '

			// Carve out walls based on genGrid
			if !genGrid[r][c].Walls[WallTop] {
				m.Grid[2*r][2*c+1] = ' '
			}
			if !genGrid[r][c].Walls[WallRight] {
				m.Grid[2*r+1][2*c+2] = ' '
			}
			if !genGrid[r][c].Walls[WallBottom] {
				m.Grid[2*r+2][2*c+1] = ' '
			}
			if !genGrid[r][c].Walls[WallLeft] {
				m.Grid[2*r+1][2*c] = ' '
			}
		}
	}

	// Set start and end points
	m.Start = Point{1, 1} // Top-left corner of the path cells
	m.End = Point{2*m.Height - 1, 2*m.Width - 1} // Bottom-right corner of the path cells
	m.Grid[m.Start.R][m.Start.C] = 'S'
	m.Grid[m.End.R][m.End.C] = 'E'
}

// SolveMaze finds a path from start to end using Breadth-First Search (BFS)
func (m *Maze) SolveMaze() []Point {
	queue := [][]Point{} // Queue of paths
	
	// Visited array to prevent cycles and redundant checks
	visited := make([][]bool, len(m.Grid))
	for i := range visited {
		visited[i] = make([]bool, len(m.Grid[0]))
	}

	// Start point
	startPath := []Point{m.Start}
	queue = append(queue, startPath)
	visited[m.Start.R][m.Start.C] = true

	// Directions: Up, Down, Left, Right
	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}

	for len(queue) > 0 {
		currentPath := queue[0]
		queue = queue[1:] // Dequeue
		current := currentPath[len(currentPath)-1]

		if current == m.End {
			return currentPath // Path found
		}

		for i := 0; i < 4; i++ {
			nR, nC := current.R+dr[i], current.C+dc[i]

			// Check bounds
			if nR < 0 || nR >= len(m.Grid) || nC < 0 || nC >= len(m.Grid[0]) {
				continue
			}

			// Check if it's a path and not visited
			if (m.Grid[nR][nC] == ' ' || m.Grid[nR][nC] == 'E') && !visited[nR][nC] {
				visited[nR][nC] = true
				newPath := append([]Point{}, currentPath...) // Create a new slice to avoid modifying original path
				newPath = append(newPath, Point{nR, nC})
				queue = append(queue, newPath)
			}
		}
	}
	return nil // No path found
}

// PrintMaze prints the maze to the console
func (m *Maze) PrintMaze() {
	for _, row := range m.Grid {
		for _, cell := range row {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed random number generator

	mazeWidth := 25
	mazeHeight := 15

	maze := NewMaze(mazeWidth, mazeHeight)
	maze.GenerateMaze()

	fmt.Println("Generated Maze:")
	maze.PrintMaze()

	solutionPath := maze.SolveMaze()

	if solutionPath != nil {
		fmt.Println("\nSolved Maze:")
		// Mark the solution path on the maze grid
		for i, p := range solutionPath {
			if p != maze.Start && p != maze.End {
				maze.Grid[p.R][p.C] = 'X'
			} else if p == maze.Start && i != 0 { // If start is part of path but not the very first point (already 'S')
				maze.Grid[p.R][p.C] = 'S'
			} else if p == maze.End && i != len(solutionPath)-1 { // If end is part of path but not the very last point (already 'E')
				maze.Grid[p.R][p.C] = 'E'
			}
		}
		maze.PrintMaze()
	} else {
		fmt.Println("\nNo solution found for the maze.")
	}
}

// Additional implementation at 2025-06-21 01:33:23
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Point represents a coordinate in the maze
type Point struct {
	R, C int
}

// Maze represents the maze structure
type Maze struct {
	Width, Height int
	Grid          [][]int // 0: path, 1: wall, 2: start, 3: end, 4: solved path
	Start, End    Point
}

// NewMaze creates a new maze with specified dimensions
func NewMaze(width, height int) *Maze {
	// Ensure odd dimensions for internal cells to be odd and walls even
	// This simplifies the generation logic where cells are (2r+1, 2c+1)
	if width%2 == 0 {
		width++
	}
	if height%2 == 0 {
		height++
	}

	maze := &Maze{
		Width:  width,
		Height: height,
		Grid:   make([][]int, height),
	}

	for r := 0; r < height; r++ {
		maze.Grid[r] = make([]int, width)
		for c := 0; c < width; c++ {
			maze.Grid[r][c] = 1 // Initialize all as walls
		}
	}

	// Set start and end points (top-left and bottom-right corners of the path cells)
	maze.Start = Point{1, 1}
	maze.End = Point{height - 2, width - 2}

	return maze
}

// Generate uses Depth-First Search (DFS) to create the maze
func (m *Maze) Generate() {
	rand.Seed(time.Now().UnixNano())

	stack := []Point{}
	visited := make([][]bool, m.Height)
	for r := range visited {
		visited[r] = make([]bool, m.Width)
	}

	// Start DFS from the start point (physical grid coordinates)
	startCell := m.Start
	stack = append(stack, startCell)
	visited[startCell.R][startCell.C] = true
	m.Grid[startCell.R][startCell.C] = 0 // Carve path

	dr := []int{-2, 2, 0, 0} // Row changes for N, S, E, W (skipping a wall)
	dc := []int{0, 0, -2, 2} // Column changes for N, S, E, W (skipping a wall)

	for len(stack) > 0 {
		curr := stack[len(stack)-1] // Peek
		
		// Find unvisited neighbors
		unvisitedNeighbors := []Point{}
		for i := 0; i < 4; i++ {
			nr, nc := curr.R+dr[i], curr.C+dc[i]
			
			// Check bounds and if it's an unvisited cell (odd coordinates)
			if nr >= 0 && nr < m.Height && nc >= 0 && nc < m.Width &&
			   nr%2 == 1 && nc%2 == 1 && // Ensure it's a cell, not a wall
			   !visited[nr][nc] {
				unvisitedNeighbors = append(unvisitedNeighbors, Point{nr, nc})
			}
		}

		if len(unvisitedNeighbors) > 0 {
			// Choose a random unvisited neighbor
			next := unvisitedNeighbors[rand.Intn(len(unvisitedNeighbors))]

			// Carve path to the next cell
			m.Grid[next.R][next.C] = 0
			
			// Carve path through the wall between current and next
			wallR := curr.R + (next.R-curr.R)/2
			wallC := curr.C + (next.C-curr.C)/2
			m.Grid[wallR][wallC] = 0

			visited[next.R][next.C] = true
			stack = append(stack, next)
		} else {
			// Backtrack
			stack = stack[:len(stack)-1]
		}
	}
	
	// Ensure start and end points are marked after generation
	m.Grid[m.Start.R][m.Start.C] = 2
	m.Grid[m.End.R][m.End.C] = 3
}

// Solve uses Breadth-First Search (BFS) to find a path from start to end
func (m *Maze) Solve() bool {
	queue := []Point{}
	visited := make([][]bool, m.Height)
	for r := range visited {
		visited[r] = make([]bool, m.Width)
	}

	parent := make(map[Point]Point) // To reconstruct the path

	queue = append(queue, m.Start)
	visited[m.Start.R][m.Start.C] = true

	dr := []int{-1, 1, 0, 0} // Row changes for N, S, E, W
	dc := []int{0, 0, -1, 1} // Column changes for N, S, E, W

	found := false
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == m.End {
			found = true
			break
		}

		for i := 0; i < 4; i++ {
			nr, nc := curr.R+dr[i], curr.C+dc[i]

			if nr >= 0 && nr < m.Height && nc >= 0 && nc < m.Width &&
				!visited[nr][nc] && m.Grid[nr][nc] != 1 { // Not a wall
				
				visited[nr][nc] = true
				parent[Point{nr, nc}] = curr
				queue = append(queue, Point{nr, nc})
			}
		}
	}

	if found {
		// Reconstruct path
		curr := m.End
		for curr

// Additional implementation at 2025-06-21 01:34:25
package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Cell struct {
	Row, Col int
	Walls    [4]bool // 0: Top, 1: Right, 2: Bottom, 3: Left
	Visited  bool    // Used during maze generation
	IsPath   bool    // Used to mark the solution path
}

type Maze struct {
	Width, Height int
	Grid          [][]Cell
	Rand          *rand.Rand
}

func NewMaze(width, height int) *Maze {
	grid := make([][]Cell, height)
	for r := 0; r < height; r++ {
		grid[r] = make([]Cell, width)
		for c := 0; c < width; c++ {
			grid[r][c] = Cell{
				Row:   r,
				Col:   c,
				Walls: [4]bool{true, true, true, true},
			}
		}
	}
	return &Maze{
		Width:  width,
		Height: height,
		Grid:   grid,
		Rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (m *Maze) removeWall(c1, c2 *Cell) {
	if c1.Row == c2.Row {
		if c1.Col < c2.Col {
			c1.Walls[1] = false
			c2.Walls[3] = false
		} else {
			c1.Walls[3] = false
			c2.Walls[1] = false
		}
	} else if c1.Col == c2.Col {
		if c1.Row < c2.Row {
			c1.Walls[2] = false
			c2.Walls[0] = false
		} else {
			c1.Walls[0] = false
			c2.Walls[2] = false
		}
	}
}

func (m *Maze) getUnvisitedNeighbors(cell *Cell) []*Cell {
	neighbors := []*Cell{}
	if cell.Row > 0 && !m.Grid[cell.Row-1][cell.Col].Visited {
		neighbors = append(neighbors, &m.Grid[cell.Row-1][cell.Col])
	}
	if cell.Col < m.Width-1 && !m.Grid[cell.Row][cell.Col+1].Visited {
		neighbors = append(neighbors, &m.Grid[cell.Row][cell.Col+1])
	}
	if cell.Row < m.Height-1 && !m.Grid[cell.Row+1][cell.Col].Visited {
		neighbors = append(neighbors, &m.Grid[cell.Row+1][cell.Col])
	}
	if cell.Col > 0 && !m.Grid[cell.Row][cell.Col-1].Visited {
		neighbors = append(neighbors, &m.Grid[cell.Row][cell.Col-1])
	}
	return neighbors
}

func (m *Maze) Generate() {
	stack := []*Cell{}
	startCell := &m.Grid[0][0]
	startCell.Visited = true
	stack = append(stack, startCell)

	for len(stack) > 0 {
		currentCell := stack[len(stack)-1]
		unvisitedNeighbors := m.getUnvisitedNeighbors(currentCell)

		if len(unvisitedNeighbors) > 0 {
			nextCell := unvisitedNeighbors[m.Rand.Intn(len(unvisitedNeighbors))]

			m.removeWall(currentCell, nextCell)
			nextCell.Visited = true
			stack = append(stack, nextCell)
		} else {
			stack = stack[:len(stack)-1]
		}
	}
}

func (m *Maze) Solve(startRow, startCol, endRow, endCol int) bool {
	queue := []*Cell{}
	parentMap := make(map[*Cell]*Cell)
	visited := make(map[*Cell]bool)

	startCell := &m.Grid[startRow][startCol]
	endCell := &m.Grid[endRow][endCol]

	queue = append(queue, startCell)
	visited[startCell] = true

	found := false
	for len(queue) > 0 {
		currentCell := queue[0]
		queue = queue[1:]

		if currentCell == endCell {
			found = true
			break
		}

		moves := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

		for i, move := range moves {
			if !currentCell.Walls[i] {
				nextRow, nextCol := currentCell.Row+move[0], currentCell.Col+move[1]

				if nextRow >= 0 && nextRow < m.Height && nextCol >= 0 && nextCol < m.Width {
					nextCell := &m.Grid[nextRow][nextCol]
					if !visited[nextCell] {
						visited[nextCell] = true
						parentMap[nextCell] = currentCell
						queue = append(queue, nextCell)
					}
				}
			}
		}
	}

	if found {
		pathCell := endCell
		for pathCell != nil && pathCell != startCell {
			pathCell.IsPath = true
			pathCell = parentMap[pathCell]
		}
		startCell.IsPath = true
		return true
	}
	return false
}

func (m *Maze) Print() {
	for c := 0; c < m.Width; c++ {
		fmt.Print("+-")
	}
	fmt.Println("+")

	for r := 0; r < m.Height; r++ {
		fmt.Print("|")
		for c := 0; c < m.Width; c++ {
			cell := m.Grid[r][c]
			if cell.IsPath {
				fmt.Print("o")
			} else {
				fmt.Print(" ")
			}
			if cell.Walls[1] {
				fmt.Print("|")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()

		fmt.Print("+")
		for c := 0; c < m.Width; c++ {
			cell := m.Grid[r][c]
			if cell.Walls[2] {
				fmt.Print("--")
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Println("+")
	}
}

func main() {
	mazeWidth := 20
	mazeHeight := 10

	maze := NewMaze(mazeWidth, mazeHeight)
	maze.Generate()

	fmt.Println("Generated Maze:")
	maze.Print()

	fmt.Println("\nSolving Maze...")
	if maze.Solve(0, 0, mazeWidth-1, mazeHeight-1) {
		fmt.Println("Maze Solved! Path marked with 'o':")
		maze.Print()
	} else {
		fmt.Println("Could not find a solution.")
	}
}