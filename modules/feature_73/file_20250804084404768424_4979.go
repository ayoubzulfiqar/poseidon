package main

import "fmt"

type Point struct {
	R, C int
}

type Maze struct {
	grid  [][]int
	rows  int
	cols  int
	start Point
	end   Point
}

func solveMazeRecursive(maze *Maze, current Point, path *[]Point, visited [][]bool) bool {
	if current.R < 0 || current.R >= maze.rows || current.C < 0 || current.C >= maze.cols {
		return false
	}
	if maze.grid[current.R][current.C] == 1 {
		return false
	}
	if visited[current.R][current.C] {
		return false
	}

	visited[current.R][current.C] = true
	*path = append(*path, current)

	if current.R == maze.end.R && current.C == maze.end.C {
		return true
	}

	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}

	for i := 0; i < 4; i++ {
		nextR, nextC := current.R+dr[i], current.C+dc[i]
		if solveMazeRecursive(maze, Point{nextR, nextC}, path, visited) {
			return true
		}
	}

	*path = (*path)[:len(*path)-1]

	return false
}

func SolveMaze(maze *Maze) []Point {
	visited := make([][]bool, maze.rows)
	for i := range visited {
		visited[i] = make([]bool, maze.cols)
	}

	var path []Point

	if solveMazeRecursive(maze, maze.start, &path, visited) {
		return path
	}
	return nil
}

func main() {
	grid := [][]int{
		{2, 0, 1, 0},
		{0, 0, 1, 0},
		{1, 0, 0, 0},
		{0, 0, 1, 3},
	}

	rows := len(grid)
	cols := len(grid[0])

	var start, end Point
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == 2 {
				start = Point{r, c}
			} else if grid[r][c] == 3 {
				end = Point{r, c}
			}
		}
	}

	maze := &Maze{
		grid:  grid,
		rows:  rows,
		cols:  cols,
		start: start,
		end:   end,
	}

	path := SolveMaze(maze)

	if path != nil {
		fmt.Println("Path found:")
		for _, p := range path {
			fmt.Printf("(%d, %d) -> ", p.R, p.C)
		}
		fmt.Println("END")
	} else {
		fmt.Println("No path found.")
	}

	fmt.Println("\n--- Second Maze ---")
	grid2 := [][]int{
		{2, 1, 0},
		{0, 1, 0},
		{0, 1, 3},
	}
	rows2 := len(grid2)
	cols2 := len(grid2[0])
	var start2, end2 Point
	for r := 0; r < rows2; r++ {
		for c := 0; c < cols2; c++ {
			if grid2[r][c] == 2 {
				start2 = Point{r, c}
			} else if grid2[r][c] == 3 {
				end2 = Point{r, c}
			}
		}
	}
	maze2 := &Maze{
		grid:  grid2,
		rows:  rows2,
		cols:  cols2,
		start: start2,
		end:   end2,
	}
	path2 := SolveMaze(maze2)
	if path2 != nil {
		fmt.Println("Path found:")
		for _, p := range path2 {
			fmt.Printf("(%d, %d) -> ", p.R, p.C)
		}
		fmt.Println("END")
	} else {
		fmt.Println("No path found.")
	}
}

// Additional implementation at 2025-08-04 08:44:45
package main

import (
	"fmt"
)

const (
	Wall         rune = '#'
	Path         rune = '.'
	Start        rune = 'S'
	End          rune = 'E'
	SolutionPath rune = '*'
)

type Point struct {
	Row, Col int
}

type Maze [][]rune

func printMaze(maze Maze) {
	for _, row := range maze {
		for _, char := range row {
			fmt.Printf("%c ", char)
		}
		fmt.Println()
	}
}

func findStartAndEnd(maze Maze) (Point, Point, bool) {
	var start, end Point
	foundStart, foundEnd := false, false

	for r := 0; r < len(maze); r++ {
		for c := 0; c < len(maze[0]); c++ {
			if maze[r][c] == Start {
				start = Point{r, c}
				foundStart = true
			} else if maze[r][c] == End {
				end = Point{r, c}
				foundEnd = true
			}
		}
	}
	return start, end, foundStart && foundEnd
}

func solveMazeDFS(maze Maze, current Point, end Point, visited [][]bool) bool {
	// Base cases:
	// 1. Out of bounds
	if current.Row < 0 || current.Row >= len(maze) ||
		current.Col < 0 || current.Col >= len(maze[0]) {
		return false
	}

	// 2. Hit a wall
	if maze[current.Row][current.Col] == Wall {
		return false
	}

	// 3. Already visited in the current path attempt (to prevent cycles and redundant exploration)
	if visited[current.Row][current.Col] {
		return false
	}

	// 4. Found the end
	if current == end {
		return true
	}

	// Mark current cell as visited for this path attempt
	visited[current.Row][current.Col] = true

	// Define possible moves (Up, Down, Left, Right)
	moves := []Point{
		{current.Row - 1, current.Col}, // Up
		{current.Row + 1, current.Col}, // Down
		{current.Row, current.Col - 1}, // Left
		{current.Row, current.Col + 1}, // Right
	}

	// Recursively try each move
	for _, move := range moves {
		if solveMazeDFS(maze, move, end, visited) {
			// If a path is found from this move, mark the current cell as part of the solution
			// unless it's the start or end point.
			if maze[current.Row][current.Col] == Path {
				maze[current.Row][current.Col] = SolutionPath
			}
			return true
		}
	}

	// Backtrack: If no path found from this cell, unmark it as visited
	// This allows other potential paths to explore this cell if they approach it differently.
	visited[current.Row][current.Col] = false
	return false
}

func main() {
	sampleMaze := Maze{
		{'S', '.', '.', '#', '.', '.', '.'},
		{'#', '.', '#', '#', '.', '#', '.'},
		{'.', '.', '.', '.', '.', '#', '.'},
		{'.', '#', '#', '#', '.', '#', '.'},
		{'.', '.', '.', '.', '.', '.', 'E'},
	}

	fmt.Println("Original Maze:")
	printMaze(sampleMaze)
	fmt.Println()

	start, end, found := findStartAndEnd(sampleMaze)
	if !found {
		fmt.Println("Error: Start or End point not found in the maze.")
		return
	}

	rows := len(sampleMaze)
	cols := len(sampleMaze[0])
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	if solveMazeDFS(sampleMaze, start, end, visited) {
		fmt.Println("Maze Solved! Path marked with '*':")
		printMaze(sampleMaze)
	} else {
		fmt.Println("No path found in the maze.")
	}

	fmt.Println("\nAnother Maze Example:")
	sampleMaze2 := Maze{
		{'S', '#', '.'},
		{'.', '#', '.'},
		{'.', '.', 'E'},
	}
	fmt.Println("Original Maze 2:")
	printMaze(sampleMaze2)
	fmt.Println()

	start2, end2, found2 := findStartAndEnd(sampleMaze2)
	if !found2 {
		fmt.Println("Error: Start or End point not found in the maze 2.")
		return
	}

	rows2 := len(sampleMaze2)
	cols2 := len(sampleMaze2[0])
	visited2 := make([][]bool, rows2)
	for i := range visited2 {
		visited2[i] = make([]bool, cols2)
	}

	if solveMazeDFS(sampleMaze2, start2, end2, visited2) {
		fmt.Println("Maze 2 Solved! Path marked with '*':")
		printMaze(sampleMaze2)
	} else {
		fmt.Println("No path found in maze 2.")
	}

	fmt.Println("\nUnsolvable Maze Example:")
	sampleMaze3 := Maze{
		{'S', '#', '.'},
		{'.', '#', '.'},
		{'.', '#', 'E'},
	}
	fmt.Println("Original Maze 3:")
	printMaze(sampleMaze3)
	fmt.Println()

	start3, end3, found3 := findStartAndEnd(sampleMaze3)
	if !found3 {
		fmt.Println("Error: Start or End point not found in the maze 3.")
		return
	}

	rows3 := len(sampleMaze3)
	cols3 := len(sampleMaze3[0])
	visited3 := make([][]bool, rows3)
	for i := range visited3 {
		visited3[i] = make([]bool, cols3)
	}

	if solveMazeDFS(sampleMaze3, start3, end3, visited3) {
		fmt.Println("Maze 3 Solved! Path marked with '*':")
		printMaze(sampleMaze3)
	} else {
		fmt.Println("No path found in maze 3.")
	}
}

// Additional implementation at 2025-08-04 08:45:47
package main

import (
	"fmt"
)

const (
	Wall    rune = '#'
	Path    rune = '.'
	Start   rune = 'S'
	End     rune = 'E'
	Visited rune = 'V'
)

type Point struct {
	R, C int
}

func printMaze(maze [][]rune) {
	for _, row := range maze {
		for _, char := range row {
			fmt.Printf("%c ", char)
		}
		fmt.Println()
	}
	fmt.Println()
}

func solveMaze(maze [][]rune, currR, currC int, endR, endC int) bool {
	if currR < 0 || currR >= len(maze) || currC < 0 || currC >= len(maze[0]) {
		return false
	}

	if maze[currR][currC] == Wall || maze[currR][currC] == Visited {
		return false
	}

	if currR == endR && currC == endC {
		return true
	}

	originalChar := maze[currR][currC]
	maze[currR][currC] = Visited

	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}

	for i := 0; i < 4; i++ {
		newR, newC := currR+dr[i], currC+dc[i]
		if solveMaze(maze, newR, newC, endR, endC) {
			return true
		}
	}

	maze[currR][currC] = originalChar
	return false
}

func main() {
	maze1 := [][]rune{
		{'S', '.', '.', '#', '.', '.', '.'},
		{'.', '#', '.', '#', '.', '#', '.'},
		{'.', '#', '.', '.', '.', '#', '.'},
		{'.', '.', '#', '#', '.', '.', '.'},
		{'#', '.', '.', '.', '#', '.', 'E'},
	}

	var startPoint1, endPoint1 Point

	for r, row := range maze1 {
		for c, char := range row {
			if char == Start {
				startPoint1 = Point{R: r, C: c}
			} else if char == End {
				endPoint1 = Point{R: r, C: c}
			}
		}
	}

	fmt.Println("Initial Maze 1:")
	printMaze(maze1)

	found1 := solveMaze(maze1, startPoint1.R, startPoint1.C, endPoint1.R, endPoint1.C)

	fmt.Println("Maze 1 After Solving:")
	printMaze(maze1)

	if found1 {
		fmt.Println("Path found in Maze 1!")
	} else {
		fmt.Println("No path found in Maze 1.")
	}

	fmt.Println("-----------------------------------")

	maze2 := [][]rune{
		{'S', '.', '.', '#', '.', '.', '.'},
		{'.', '#', '.', '#', '.', '#', '.'},
		{'.', '#', '.', '.', '.', '#', '.'},
		{'.', '.', '#', '#', '.', '#', '.'},
		{'#', '.', '.', '.', '#', '.', 'E'},
	}

	var startPoint2, endPoint2 Point

	for r, row := range maze2 {
		for c, char := range row {
			if char == Start {
				startPoint2 = Point{R: r, C: c}
			} else if char == End {
				endPoint2 = Point{R: r, C: c}
			}
		}
	}

	fmt.Println("Initial Maze 2 (No Path):")
	printMaze(maze2)

	found2 := solveMaze(maze2, startPoint2.R, startPoint2.C, endPoint2.R, endPoint2.C)

	fmt.Println("Maze 2 After Solving (No Path):")
	printMaze(maze2)

	if found2 {
		fmt.Println("Path found in Maze 2!")
	} else {
		fmt.Println("No path found in Maze 2.")
	}
}

// Additional implementation at 2025-08-04 08:47:14
package main

import (
	"fmt"
)

// Point represents a coordinate in the maze
type Point struct {
	Row int
	Col int
}

// Maze constants for different cell types
const (
	Wall    rune = '#' // Impassable wall
	Path    rune = '.' // Open path
	Start   rune = 'S' // Starting point
	End     rune = 'E' // Ending point
	Visited rune = 'V' // Temporarily marks cells visited during DFS exploration
	Solved  rune = 'P' // Marks cells that are part of the final found path
)

// solveMazeDFS recursively finds a path from the current point to the end point
// using Depth-First Search. It modifies the maze in place to mark visited cells
// and builds the path slice.
func solveMazeDFS(maze [][]rune, current, end Point, path *[]Point) bool {
	// Base cases for recursion termination:

	// 1. Current point is out of maze bounds
	if current.Row < 0 || current.Row >= len(maze) || current.Col < 0 || current.Col >= len(maze[0]) {
		return false
	}
	// 2. Current point is a wall
	if maze[current.Row][current.Col] == Wall {
		return false
	}
	// 3. Current point has already been visited in the current path (to avoid cycles and redundant work)
	if maze[current.Row][current.Col] == Visited {
		return false
	}
	// 4. Current point is the end point - path found!
	if current == end {
		*path = append(*path, current) // Add the end point to the path
		return true
	}

	// Mark the current cell as visited. Store its original character
	// to restore it if this path doesn't lead to the solution (backtracking).
	originalChar := maze[current.Row][current.Col]
	if originalChar != Start { // Do not overwrite the 'S' character with 'V'
		maze[current.Row][current.Col] = Visited
	}

	// Add the current cell to the potential path.
	*path = append(*path, current)

	// Define possible moves (up, down, left, right)
	dr := []int{-1, 1, 0, 0} // Delta row for neighbors
	dc := []int{0, 0, -1, 1} // Delta column for neighbors

	// Explore each neighbor
	for i := 0; i < 4; i++ {
		next := Point{current.Row + dr[i], current.Col + dc[i]}
		if solveMazeDFS(maze, next, end, path) {
			return true // If a path is found through this neighbor, propagate success
		}
	}

	// Backtrack: If no path was found from this cell (all neighbors failed),
	// remove it from the current path and restore its original state in the maze.
	*path = (*path)[:len(*path)-1] // Remove the last element (current point) from the path
	if originalChar != Start {     // Do not restore 'S'
		maze[current.Row][current.Col] = originalChar // Restore the original character (e.g., '.')
	}
	return false // No path found from this cell
}

// printMaze prints the current state of the maze to the console.
func printMaze(maze [][]rune) {
	for _, row := range maze {
		for _, char := range row {
			fmt.Printf("%c ", char)
		}
		fmt.Println()
	}
}

func main() {
	// --- Test Case 1: Maze with a solvable path ---
	mazeData1 := [][]rune{
		{'#', '#', '#', '#', '#', '#', '#', '#', '#', '#'},
		{'#', 'S', '.', '.', '#', '.', '.', '.', '.', '#'},
		{'#', '#', '#', '.', '#', '.', '#', '#', '.', '#'},
		{'#', '.', '.', '.', '.', '.', '#', '.', '.', '#'},
		{'#', '.', '#', '#', '#', '#', '#', '.', '#', '#'},
		{'#', '.', '.', '.', '.', '.', '.', '.', '.', '#'},
		{'#', '#', '#', '#', '#', '#', '#', '#', 'E', '#'},
	}

	// Create a deep copy of the maze for the solver to modify,
	// keeping the original maze definition clean.
	mazeCopy1 := make([][]rune, len(mazeData1))
	for i := range mazeData1 {
		mazeCopy1[i] = make([]rune, len(mazeData1[i]))
		copy(mazeCopy1[i], mazeData1[i])
	}

	var startPoint1, endPoint1 Point
	foundStart1, foundEnd1 := false, false

	// Find the start and end points in the maze
	for r := 0; r < len(mazeData1); r++ {
		for c := 0; c < len(mazeData1[0]); c++ {
			if mazeData1[r][c] == Start {
				startPoint1 = Point{r, c}
				foundStart1 = true
			} else if mazeData1[r][c] == End {
				endPoint1 = Point{r, c}
				foundEnd1 = true
			}
		}
	}

	if !foundStart1 || !foundEnd1 {
		fmt.Println("Error: Start or End point not found in Test Case 1 maze.")
		return
	}

	fmt.Println("--- Test Case 1: Solvable Maze ---")
	fmt.Println("Original Maze:")
	printMaze(mazeData1)
	fmt.Println()

	var path1 []Point
	if solveMazeDFS(mazeCopy1, startPoint1, endPoint1, &path1) {
		fmt.Println("Path Found!")
		// Mark the found path on the copied maze
		for _, p := range path1 {
			// Do not overwrite Start or End markers
			if mazeCopy1[p.Row][p.Col] != Start && mazeCopy1[p.Row][p.Col] != End {
				mazeCopy1[p.Row][p.Col] = Solved
			}
		}
		fmt.Println("Solved Maze:")
		printMaze(mazeCopy1)
	} else {
		fmt.Println("No path found for Test Case 1.")
	}

	fmt.Println("\n-----------------------------------\n")

	// --- Test Case 2: Maze with no solvable path ---
	mazeData2 := [][]rune{
		{'S', '.', '#'},
		{'.', '.', '#'},
		{'#', 'E', '#'},
	}

	mazeCopy2 := make([][]rune, len(mazeData2))
	for i := range mazeData2 {
		mazeCopy2[i] = make([]rune, len(mazeData2[i]))
		copy(mazeCopy2[i], mazeData2[i])
	}

	var startPoint2, endPoint2 Point
	foundStart2, foundEnd2 := false, false

	for r := 0; r < len(mazeData2); r++ {
		for c := 0; c < len(mazeData2[0]); c++ {
			if mazeData2[r][c] == Start {
				startPoint2 = Point{r, c}
				foundStart2 = true
			} else if mazeData2[r][c] == End {
				endPoint2 = Point{r, c}
				foundEnd2 = true
			}
		}
	}

	if !foundStart2 || !foundEnd2 {
		fmt.Println("Error: Start or End point not found in Test Case 2 maze.")
		return
	}

	fmt.Println("--- Test Case 2: Unsolvable Maze ---")
	fmt.Println("Original Maze:")
	printMaze(mazeData2)
	fmt.Println()

	var path2 []Point
	if solveMazeDFS(mazeCopy2, startPoint2, endPoint2, &path2) {
		fmt.Println("Path Found!")
		for _, p := range path2 {
			if mazeCopy2[p.Row][p.Col] != Start && mazeCopy2[p.Row][p.Col] != End {
				mazeCopy2[p.Row][p.Col] = Solved
			}
		}
		fmt.Println("Solved Maze:")
		printMaze(mazeCopy2)
	} else {
		fmt.Println("No path found for Test Case 2.")
	}
}