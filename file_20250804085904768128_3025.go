package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Cell struct {
	RawInput  string
	Value     float64
	IsFormula bool
	Evaluated bool
}

type Spreadsheet struct {
	cells map[string]*Cell
}

func NewSpreadsheet() *Spreadsheet {
	return &Spreadsheet{
		cells: make(map[string]*Cell),
	}
}

func (s *Spreadsheet) SetCell(address, input string) {
	cell := &Cell{
		RawInput:  input,
		IsFormula: strings.HasPrefix(input, "="),
		Evaluated: false,
	}
	s.cells[address] = cell
	for _, c := range s.cells {
		c.Evaluated = false
	}
}

func (s *Spreadsheet) GetCellValue(address string) (float64, error) {
	return s.getCellValueRecursive(address, make(map[string]bool))
}

func (s *Spreadsheet) getCellValueRecursive(address string, visited map[string]bool) (float64, error) {
	cell, ok := s.cells[address]
	if !ok {
		return 0, fmt.Errorf("cell %s not found", address)
	}

	if visited[address] {
		return 0, fmt.Errorf("circular dependency detected involving cell %s", address)
	}

	if cell.Evaluated {
		return cell.Value, nil
	}

	visited[address] = true
	defer delete(visited, address)

	if cell.IsFormula {
		formula := strings.TrimPrefix(cell.RawInput, "=")
		val, err := s.evaluateFormula(formula, visited)
		if err != nil {
			return 0, err
		}
		cell.Value = val
		cell.Evaluated = true
		return val, nil
	} else {
		val, err := strconv.ParseFloat(cell.RawInput, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number format for cell %s: %w", address, err)
		}
		cell.Value = val
		cell.Evaluated = true
		return val, nil
	}
}

var cellRefRegex = regexp.MustCompile(`^[A-Z]+[0-9]+$`)

func (s *Spreadsheet) evaluateFormula(formula string, visited map[string]bool) (float64, error) {
	tokenPattern := regexp.MustCompile(`([A-Z]+[0-9]+|\d+\.?\d*|[+\-*/()])`)
	matches := tokenPattern.FindAllString(formula, -1)

	tokens := []string{}
	for _, m := range matches {
		if m != "" {
			tokens = append(tokens, m)
		}
	}

	parser := &formulaParser{
		tokens:  tokens,
		pos:     0,
		visited: visited,
		ss:      s,
	}

	val, err := parser.parseExpression()
	if err != nil {
		return 0, err
	}
	if parser.pos != len(parser.tokens) {
		return 0, fmt.Errorf("unexpected tokens at end of formula: %v", parser.tokens[parser.pos:])
	}
	return val, nil
}

type formulaParser struct {
	tokens  []string
	pos     int
	visited map[string]bool
	ss      *Spreadsheet
}

func (p *formulaParser) peek() string {
	if p.pos >= len(p.tokens) {
		return ""
	}
	return p.tokens[p.pos]
}

func (p *formulaParser) consume() string {
	if p.pos >= len(p.tokens) {
		return ""
	}
	token := p.tokens[p.pos]
	p.pos++
	return token
}

func (p *formulaParser) parseExpression() (float64, error) {
	left, err := p.parseTerm()
	if err != nil {
		return 0, err
	}

	for p.peek() == "+" || p.peek() == "-" {
		op := p.consume()
		right, err

// Additional implementation at 2025-08-04 09:00:36
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Cell represents a single cell in the spreadsheet.
type Cell struct {
	Value     float64 // The calculated numeric value of the cell
	Formula   string  // The raw formula string (e.g., "=A1+B2" or "123")
	IsFormula bool    // True if the cell contains a formula, false if it's a literal number
}

// Spreadsheet holds all cells.
type Spreadsheet struct {
	Cells map[string]*Cell // Key: "A1", "B2", etc.
}

// NewSpreadsheet creates and initializes a new Spreadsheet.
func NewSpreadsheet() *Spreadsheet {
	return &Spreadsheet{
		Cells: make(map[string]*Cell),
	}
}

// SetCell sets the content of a cell. It parses the input string to determine if it's a number

// Additional implementation at 2025-08-04 09:02:21
import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Cell represents a single cell in the spreadsheet.
type Cell struct {
	RawInput  string  // The original string input (e.g., "10", "=A1+B2")
	Value     float64 // The calculated numeric value of the cell
	IsFormula bool    // True if the cell contains a formula
}

// Spreadsheet manages the collection of cells and their dependencies.
type Spreadsheet struct {
	Cells map[string]*Cell // Stores cells, keyed by their coordinate (e.g., "A1", "B2")

	// Dependencies maps a cell coordinate to a set of coordinates it depends on.
	// Example: Dependencies["A3"] = {"A1":{}, "A2":{}} if A3 = A1 + A2
	Dependencies map[string]map[string]struct{}

	// Dependents maps a cell coordinate to a set of coordinates that depend on it.
	// This is the reverse of Dependencies, used for efficient recalculation.
	// Example: Dependents["A1"] = {"A3":{}} if A3 = A1 + A2
	Dependents map[string]map[string]struct{}
}

// NewSpreadsheet creates and initializes a new empty spreadsheet.
func NewSpreadsheet() *Spreadsheet {
	return &Spreadsheet{
		Cells:        make(map[string]*Cell),
		Dependencies: make(map[string]map[string]struct{}),
		Dependents:   make(map[string]map[string]struct{}),
	}
}

// SetCell sets the value or formula for a given cell coordinate.
// It also updates the dependency graph and triggers a recalculation of the entire sheet.
func (s *Spreadsheet) SetCell(coord, input string) error {
	// Clear old dependencies for this cell before setting new ones
	if oldDeps, ok := s.Dependencies[coord]; ok {
		for depCoord := range oldDeps {
			delete(s.Dependents[depCoord], coord)
			if len(s.Dependents[depCoord]) == 0 {
				delete(s.Dependents, depCoord)
			}
		}
	}
	delete(s.Dependencies, coord) // Remove the cell's own dependency entry

	cell := &Cell{RawInput: input}
	s.Cells[coord] = cell

	if strings.HasPrefix(input, "=") {
		cell.IsFormula = true
		formula := input[1:] // Remove the leading '='

		// Find all cell references in the formula (e.g., A1, B2, C10)
		re := regexp.MustCompile(`[A-Z]+[0-9]+`)
		matches := re.FindAllString(formula, -1)

		s.Dependencies[coord] = make(map[string]struct{})
		for _, depCoord := range matches {
			s.Dependencies[coord][depCoord] = struct{}{} // Add to current cell's dependencies

			// Add current cell to the dependents list of the cell it depends on
			if s.Dependents[depCoord] == nil {
				s.Dependents[depCoord] = make(map[string]struct{})
			}
			s.Dependents[depCoord][coord] = struct{}{}
		}
	} else {
		// Attempt to parse as a number
		val, err := strconv.ParseFloat(input, 64)
		if err != nil {
			// If not a number, treat as 0 for calculation purposes (or could be NaN/error)
			cell.Value = 0.0
			cell.IsFormula = false
		} else {
			cell.Value = val
			cell.IsFormula = false
		}
	}

	// After updating a cell, recalculate the entire sheet to ensure consistency
	return s.Recalculate()
}

// GetCell retrieves the calculated value of a cell.
func (s *Spreadsheet) GetCell(coord string) (float64, error) {
	cell, exists := s.Cells[coord]
	if !exists {
		return 0, fmt.Errorf("cell %s does not exist", coord)
	}
	if math.IsNaN(cell.Value) {
		return 0, fmt.Errorf("cell %s contains a circular reference or evaluation error", coord)
	}
	return cell.Value, nil
}

// Recalculate performs a topological sort of the dependency graph and evaluates all formulas.
// It detects and marks circular dependencies.
func (s *Spreadsheet) Recalculate() error {
	// 1. Build the graph for topological sort (Kahn's algorithm)
	// Adjacency list: cell -> cells it directly depends on (for graph traversal)
	graph := make(map[string][]string)
	// In-degrees: number of incoming edges for each cell (number of dependencies)
	inDegree := make(map[string]int)

	// Initialize in-degrees for all cells that are part of a formula or have dependents
	for coord, cell := range s.Cells {
		if cell.IsFormula {
			inDegree[coord] = 0 // Will be updated by actual dependencies
		}
	}
	// Also include cells that are dependencies but might not be explicitly set yet
	for coord := range s.Dependents {
		if _, ok := s.Cells[coord]; !ok {
			s.Cells[coord] = &Cell{RawInput: "", Value: 0.0, IsFormula: false} // Create dummy cell
		}
	}

	// Populate graph and in-degrees based on Dependencies map
	for coord, deps := range s.Dependencies {
		for depCoord := range deps {
			// An edge exists from depCoord to coord (depCoord -> coord)
			// because coord depends on depCoord.
			graph[depCoord] = append(graph[depCoord], coord)
			inDegree[coord]++
		}
	}

	// 2. Initialize queue with cells that have no incoming edges (no dependencies or non-formulas)
	queue := []string{}
	for coord, cell := range s.Cells {
		if !cell.IsFormula { // Non-formula cells are "source" nodes
			queue = append(queue, coord)
		} else if inDegree[coord] == 0 { // Formula cells with no dependencies
			queue = append(queue, coord)
		}
	}

	// 3. Perform topological sort and evaluate cells
	evaluatedCount := 0
	for len(queue) > 0 {
		currentCoord := queue[0]
		queue = queue[1:]
		evaluatedCount++

		cell := s.Cells[currentCoord]
		if cell == nil {
			// This case should ideally not happen if all dependencies are initialized as dummy cells.
			continue
		}

		if cell.IsFormula {
			val, err := s.evaluateFormula(currentCoord, cell.RawInput[1:])
			if err != nil {
				// Mark cell value as NaN to indicate an error (e.g., circular ref, div by zero)
				cell.Value = math.NaN()
				fmt.Printf("Error evaluating %s: %v\n", currentCoord, err)
			} else {
				cell.Value = val
			}
		}

		// Decrement in-degree of neighbors (cells that depend on currentCoord)
		for _, neighborCoord := range graph[currentCoord] {
			inDegree[neighborCoord]--
			if inDegree[neighborCoord] == 0 {
				queue = append(queue, neighborCoord)
			}
		}
	}

	// 4. Check for cycles
	// If not all formula cells were evaluated, it means there's a cycle.
	// Non-formula cells are always evaluated first.
	for coord, cell := range s.Cells {
		if cell.IsFormula && inDegree[coord] > 0 {
			// This cell is part of a cycle, mark its value as NaN
			cell.Value = math.NaN()
			fmt.Printf("Circular reference detected involving cell: %s\n", coord)
		}
	}

	return nil
}

// Lexer for parsing formula expressions.
type Lexer struct {
	input string
	pos   int
}

// Token types
const (
	TOKEN_NUMBER = iota
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_MULTIPLY
	TOKEN_DIVIDE
	TOKEN_LPAREN
	TOKEN_RPAREN
	TOKEN_CELLREF
	TOKEN_EOF // End of file/input
)

// Token represents a lexical token from the formula string.
type Token struct {
	Type  int
	Value string
}

// NewLexer creates a new lexer for the given input string.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NextToken returns the next token from the input string.
func (l *Lexer) NextToken() Token {
	// Skip whitespace
	for l.pos < len(l.input) && unicode.IsSpace(rune(l.input[l.pos])) {
		l.pos++
	}

	if l.pos >= len(l.input) {
		return Token{Type: TOKEN_EOF}
	}

	char := rune(l.input[l.pos])

	switch char {
	case '+':
		l.pos++
		return Token{Type: TOKEN_PLUS, Value: "+"}
	case '-':
		l.pos++
		return Token{Type: TOKEN_MINUS, Value: "-"}
	case '*':
		l.pos++
		return Token{Type: TOKEN_MULTIPLY, Value: "*"}
	case '/':
		l.pos++
		return Token{Type: TOKEN_DIVIDE, Value: "/"}
	case '(':
		l.pos++
		return Token{Type: TOKEN_LPAREN, Value: "("}
	case ')':
		l.pos++
		return Token{Type: TOKEN