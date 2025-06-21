package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Cell represents a single cell in the spreadsheet.
type Cell struct {
	Input string  // The raw string input (e.g., "10", "=A1+B2")
	Value float64 // The calculated numeric value
	Err   error   // Any error during calculation
}

// Spreadsheet holds the grid of cells.
type Spreadsheet struct {
	cells map[string]*Cell
	mu    sync.RWMutex // Protects access to cells
}

// NewSpreadsheet creates a new empty spreadsheet.
func NewSpreadsheet() *Spreadsheet {
	return &Spreadsheet{
		cells: make(map[string]*Cell),
	}
}

// SetCell sets the input for a given cell address and triggers recalculation.
func (s *Spreadsheet) SetCell(address, input string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cell, ok := s.cells[address]
	if !ok {
		cell = &Cell{}
		s.cells[address] = cell
	}
	cell.Input = input
	cell.Err = nil // Clear previous error

	// Recalculate all cells. This is a naive approach for simplicity.
	// A more efficient approach would involve a dependency graph and topological sort.
	s.recalculateAll()
}

// GetValue returns the calculated value of a cell.
func (s *Spreadsheet) GetValue(address string) (float64, error) {
	s.mu.RLock()
	cell, ok := s.cells[address]
	s.mu.RUnlock()

	if !ok {
		return 0, fmt.Errorf("cell %s not found", address)
	}
	return cell.Value, cell.Err
}

// recalculateAll re-evaluates all cells in the spreadsheet.
// This function clears all calculated values and errors, then re-evaluates each cell.
func (s *Spreadsheet) recalculateAll() {
	// Clear all calculated values and errors before re-evaluation.
	for _, cell := range s.cells {
		cell.Value = 0
		cell.Err = nil
	}

	// Evaluate each cell. The evaluation function will handle dependencies and cycles.
	for addr := range s.cells {
		s.evaluateCell(addr, make(map[string]bool)) // Pass a new visited set for each top-level evaluation
	}
}

// evaluateCell calculates the value of a cell.
// It uses a 'visited' map to detect circular dependencies.
func (s *Spreadsheet) evaluateCell(address string, visited map[string]bool) {
	s.mu.RLock()
	cell, ok := s.cells[address]
	s.mu.RUnlock()

	if !ok {
		return // Cell doesn't exist, nothing to evaluate
	}

	// If already visited in the current evaluation path, it's a circular dependency.
	if visited[address] {
		cell.Err = fmt.Errorf("circular dependency detected involving %s", address)
		return
	}

	visited[address] = true
	defer delete(visited, address) // Remove from visited when done with this path

	input := cell.Input
	if strings.HasPrefix(input, "=") {
		// It's a formula
		formula := input[1:]
		value, err := s.evaluateFormula(formula, visited)
		if err != nil {
			cell.Err = err
		} else {
			cell.Value = value
		}
	} else {
		// It's a literal number or empty
		val, err := strconv.ParseFloat(input, 64)
		if err != nil {
			cell.Err = fmt.Errorf("invalid number format: %s", input)
			cell.Value = 0 // Set to 0 on error
		} else {
			cell.Value = val
		}
	}
}

// evaluateFormula parses and evaluates a formula string.
// It supports basic binary operations (+, -, *, /) and respects operator precedence.
// It does not support parentheses or functions.
func (s *Spreadsheet) evaluateFormula(formula string, visited map[string]bool) (float64, error) {
	formula = strings.TrimSpace(formula)

	// Precedence map for operators
	precedence := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}

	// Find the operator with the lowest precedence that exists in the formula.
	// If multiple operators of the same lowest precedence exist, pick the leftmost.
	bestOp := ""
	bestOpIdx := -1
	minPrecedence := 999 // Arbitrarily high value

	for i, r := range formula {
		op := string(r)
		if p, ok := precedence[op]; ok {
			// If current operator has lower precedence than minPrecedence found so far,
			// or same precedence but is to the left of the current best operator.
			if p < minPrecedence || (p == minPrecedence && i < bestOpIdx) {
				minPrecedence = p
				bestOp = op
				bestOpIdx = i
			}
		}
	}

	if bestOpIdx != -1 {
		// Found an operator, recursively evaluate left and right parts
		leftStr := strings.TrimSpace(formula[:bestOpIdx])
		rightStr := strings.TrimSpace(formula[bestOpIdx+1:])

		leftVal, err := s.getOperandValue(leftStr, visited)
		if err != nil {
			return 0, err
		}
		rightVal, err := s.getOperandValue(rightStr, visited)
		if err != nil {
			return 0, err
		}

		switch bestOp {
		case "+":
			return leftVal + rightVal, nil
		case "-":
			return leftVal - rightVal, nil
		case "*":
			return leftVal * rightVal, nil
		case "/":
			if rightVal == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return leftVal / rightVal, nil
		}
	}

	// If no operator found, it must be a single operand

// Additional implementation at 2025-06-21 04:41:10
package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Cell struct {
	Formula    string
	Value      interface{} // float64 or error
	Dependencies map[string]struct{} // Cells this cell depends on
	Dependents   map[string]struct{} // Cells that depend on this cell
	IsEvaluated bool // For topological sort and cycle detection
}

type Spreadsheet struct {
	Cells map[string]*Cell
}

func NewSpreadsheet() *Spreadsheet {
	return &Spreadsheet{
		Cells: make(map[string]*Cell),
	}
}

func (s *Spreadsheet) SetCell(name, formula string) error {
	if _, ok := s.Cells[name]; !ok {
		s.Cells[name] = &Cell{}
	}

	oldDependencies := s.Cells[name].Dependencies
	if oldDependencies == nil {
		oldDependencies = make(map[string]struct{})
	}

	newDependencies := make(map[string]struct{})
	s.Cells[name].Formula = formula
	s.Cells[name].Dependencies = newDependencies

	for dep := range oldDependencies {
		if cell, ok := s.Cells[dep]; ok {
			delete(cell.Dependents, name)
		}
	}

	re := regexp.MustCompile(`[A-Z]+[0-9]+`)
	matches := re.FindAllString(formula, -1)
	for _, depName := range matches {
		newDependencies[depName] = struct{}{}
		if _, ok := s.Cells[depName]; !ok {
			s.Cells[depName] = &Cell{}
		}
		if s.Cells[depName].Dependents == nil {
			s.Cells[depName].Dependents = make(map[string]struct{})
		}
		s.Cells[depName].Dependents[name] = struct{}{}
	}

	for _, cell := range s.Cells {
		cell.IsEvaluated = false
	}

	return nil
}

func (s *Spreadsheet) EvaluateAll() error {
	for _, cell := range s.Cells {
		cell.IsEvaluated = false
	}

	for name := range s.Cells {
		if err := s.evaluateCell(name, make(map[string]struct{})); err != nil {
			return err
		}
	}
	return nil
}

func (s *Spreadsheet) evaluateCell(name string, visited map[string]struct{}) error {
	cell, ok := s.Cells[name]
	if !ok {
		return fmt.Errorf("cell %s does not exist", name)
	}

	if cell.IsEvaluated {
		return nil
	}

	if _, inVisited := visited[name]; inVisited {
		return fmt.Errorf("circular dependency detected involving cell %s", name)
	}
	visited[name] = struct{}{}

	for depName := range cell.Dependencies {
		if err := s.evaluateCell(depName, visited); err != nil {
			return err
		}
	}

	val, err := s.evaluateFormula(cell.Formula)
	if err != nil {
		cell.Value = err
	} else {
		cell.Value = val
	}
	cell.IsEvaluated = true

	delete(visited, name)
	return nil
}

func (s *Spreadsheet) evaluateFormula(formula string) (float64, error) {
	formula = strings.TrimSpace(formula)

	if val, err := strconv.ParseFloat(formula, 64); err == nil {
		return val, nil
	}

	reCellRef := regexp.MustCompile(`^[A-Z]+[0-9]+$`)
	if reCellRef.MatchString(formula) {
		cell, ok := s.Cells[formula]
		if !ok || cell.Value == nil {
			return 0, fmt.Errorf("cell %s not found or not evaluated", formula)
		}
		if err, isErr := cell.Value.(error); isErr {
			return 0, err
		}
		if val, isFloat := cell.Value.(float64); isFloat {
			return val, nil
		}
		return 0, fmt.Errorf("invalid value type in cell %s", formula)
	}

	// Handle operators with precedence: +,- before *,/
	// This simple recursive descent parser works by finding the lowest precedence operator
	// (addition/subtraction) first. If found, it splits the expression and recursively
	// evaluates the parts. If not, it moves to the next precedence level (multiplication/division).
	// This implicitly handles precedence correctly for simple expressions without parentheses.
	for i := len(formula) - 1; i >= 0; i-- {
		char := formula[i]
		if char == '+' || char == '-' {
			leftStr := formula[:i]
			rightStr := formula[i+1:]

			leftVal, err := s.evaluateFormula(leftStr)
			if err != nil {
				return 0, err
			}
			rightVal, err := s.evaluateFormula(rightStr)
			if err != nil {
				return 0, err
			}

			if char == '+' {
				return leftVal + rightVal, nil
			}
			return leftVal - rightVal, nil
		}
	}

	for i := len(formula) - 1; i >= 0; i-- {
		char := formula[i]
		if char == '*' || char == '/' {
			leftStr := formula[:i]
			rightStr := formula[i+1:]

			leftVal, err := s.evaluateFormula(leftStr)
			if err != nil {
				return 0, err
			}
			rightVal, err := s.evaluateFormula(rightStr)
			if err != nil {
				return 0, err
			}

			if char == '*' {
				return leftVal * rightVal, nil
			}
			if rightVal == 0 {
				return 0, fmt.Errorf("division by zero in formula %s", formula)
			}
			return leftVal / rightVal, nil
		}
	}

	return 0, fmt.Errorf("invalid formula or unsupported syntax: %s", formula)
}

func (s *Spreadsheet) GetCell(name string) (interface{}, error) {
	cell, ok := s.Cells[name]
	if !ok {
		return nil, fmt.Errorf("cell %s does not exist", name)
	}
	if !cell.IsEvaluated {
		return nil, fmt.Errorf("cell %s not yet evaluated", name)
	}
	return cell.Value, nil
}

func (s *Spreadsheet) PrintSpreadsheet() {
	maxCol := 0
	maxRow := 0
	cellNames := make([]string, 0, len(s.Cells))
	for name := range s.Cells {
		cellNames = append(cellNames, name)
		re := regexp.MustCompile(`([A-Z]+)([0-9]+)`)
		matches := re.FindStringSubmatch(name)
		if len(matches) == 3 {
			colStr := matches[1]
			rowStr := matches[2]
			col := 0
			for _, r := range col

// Additional implementation at 2025-06-21 04:42:27


// Additional implementation at 2025-06-21 04:43:37
package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Cell struct {
	Value    float64
	Formula  string
	IsFormula bool
	Error    string
}

type Spreadsheet struct {
	cells map[string]*Cell
}

func NewSpreadsheet() *Spreadsheet {
	return &Spreadsheet{
		cells: make(map[string]*Cell),
	}
}

func (s *Spreadsheet) SetCell(coord, input string) {
	cell := &Cell{
		Formula: input,
	}

	if val, err := strconv.ParseFloat(input, 64); err == nil {
		cell.Value = val
		cell.IsFormula = false
	} else {
		cell.IsFormula = true
	}

	s.cells[coord] = cell
	s.Recalculate()
}

func (s *Spreadsheet) GetCell(coord string) *Cell {
	return s.cells[coord]
}

func (s *Spreadsheet) Recalculate() {
	formulaCoords := []string{}