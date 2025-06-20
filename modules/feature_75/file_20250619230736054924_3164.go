package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const dataFile = "todos.json"

type Todo struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

var todos []Todo
var nextID int

func loadTodos() error {
	data, err := ioutil.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			todos = []Todo{}
			nextID = 1
			return nil
		}
		return fmt.Errorf("error reading file: %w", err)
	}

	err = json.Unmarshal(data, &todos)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	nextID = 1
	for _, todo := range todos {
		if todo.ID >= nextID {
			nextID = todo.ID + 1
		}
	}
	return nil
}

func saveTodos() error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	err = ioutil.WriteFile(dataFile, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}

func addTodo(task string) {
	newTodo := Todo{
		ID:        nextID,
		Task:      task,
		Completed: false,
	}
	todos = append(todos, newTodo)
	nextID++
	fmt.Printf("Added todo: \"%s\" (ID: %d)\n", task, newTodo.ID)
	saveTodos()
}

func listTodos() {
	if len(todos) == 0 {
		fmt.Println("No todos yet!")
		return
	}
	fmt.Println("--- Your Todos ---")
	for _, todo := range todos {
		status := "[ ]"
		if todo.Completed {
			status = "[x]"
		}
		fmt.Printf("%s %d. %s\n", status, todo.ID, todo.Task)
	}
	fmt.Println("------------------")
}

func completeTodo(id int) {
	found := false
	for i := range todos {
		if todos[i].ID == id {
			if todos[i].Completed {
				fmt.Printf("Todo %d is already completed.\n", id)
				return
			}
			todos[i].Completed = true
			found = true
			fmt.Printf("Completed todo: \"%s\" (ID: %d)\n", todos[i].Task, id)
			break
		}
	}
	if !found {
		fmt.Printf("Todo with ID %d not found.\n", id)
	}
	saveTodos()
}

func deleteTodo(id int) {
	originalLen := len(todos)
	newTodos := []Todo{}
	for _, todo := range todos {
		if todo.ID != id {
			newTodos = append(newTodos, todo)
		}
	}
	if len(newTodos) == originalLen {
		fmt.Printf("Todo with ID %d not found.\n", id)
	} else {
		todos = newTodos
		fmt.Printf("Deleted todo with ID: %d\n", id)
	}
	saveTodos()
}

func printHelp() {
	fmt.Println("\n--- Todo List Commands ---")
	fmt.Println("add <task>    - Add a new todo item")
	fmt.Println("list          - List all todo items")
	fmt.Println("complete <id> - Mark a todo item as completed")
	fmt.Println("delete <id>   - Delete a todo item")
	fmt.Println("help          - Show this help message")
	fmt.Println("exit          - Exit the application")
	fmt.Println("--------------------------")
}

func main() {
	err := loadTodos()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load todos: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to your Go Todo List!")
	printHelp()

	for {
		fmt.Print("\nEnter command > ")
		reader.Scan()
		input := strings.TrimSpace(reader.Text())
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		command := strings.ToLower(parts[0])

		switch command {
		case "add":
			if len(parts) < 2 {
				fmt.Println("Usage: add <task>")
				continue
			}
			task := strings.Join(parts[1:], " ")
			addTodo(task)
		case "list":
			listTodos()
		case "complete":
			if len(parts) != 2 {
				fmt.Println("Usage: complete <id>")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid ID. Please enter a number.")
				continue
			}
			completeTodo(id)
		case "delete":
			if len(parts) != 2 {
				fmt.Println("Usage: delete <id>")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid ID. Please enter a number.")
				continue
			}
			deleteTodo(id)
		case "help":
			printHelp()
		case "exit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

// Additional implementation at 2025-06-19 23:08:12
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	dataFile = "todos.json"
)

type Todo struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Priority    string    `json:"priority"` // "High", "Medium", "Low"
	DueDate     time.Time `json:"dueDate"`
	CreatedAt   time.Time `json:"createdAt"`
}

var todos []Todo
var nextID int

func main() {
	loadTodos()
	if len(todos) > 0 {
		maxID := 0
		for _, todo := range todos {
			if todo.ID > maxID {
				maxID = todo.ID
			}
		}
		nextID = maxID + 1
	} else {
		nextID = 1
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		fmt.Println("--- Go Todo List ---")
		fmt.Println("1. Add Todo")
		fmt.Println("2. List Todos")
		fmt.Println("3. Mark Todo as Complete")
		fmt.Println("4. Delete Todo")
		fmt.Println("5. Edit Todo")
		fmt.Println("6. Filter Todos")
		fmt.Println("7. Sort Todos")
		fmt.Println("8. Exit")
		fmt.Print("Choose an option: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			addTodo(reader)
		case "2":
			listTodos(todos)
			fmt.Print("\nPress Enter to continue...")
			reader.ReadString('\n')
		case "3":
			markComplete(reader)
		case "4":
			deleteTodo(reader)
		case "5":
			editTodo(reader)
		case "6":
			filterTodosMenu(reader)
		case "7":
			sortTodosMenu(reader)
		case "8":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
			time.Sleep(1 * time.Second)
		}
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func addTodo(reader *bufio.Reader) {
	fmt.Print("Enter todo description: ")
	desc, _ := reader.ReadString('\n')
	desc = strings.TrimSpace(desc)
	if desc == "" {
		fmt.Println("Description cannot be empty.")
		time.Sleep(1 * time.Second)
		return
	}

	priority := getPriorityInput(reader)
	dueDate := getDueDateInput(reader)

	newTodo := Todo{
		ID:          nextID,
		Description: desc,
		Completed:   false,
		Priority:    priority,
		DueDate:     dueDate,
		CreatedAt:   time.Now(),
	}
	todos = append(todos, newTodo)
	nextID++
	saveTodos()
	fmt.Println("Todo added successfully!")
	time.Sleep(1 * time.Second)
}

func getPriorityInput(reader *bufio.Reader) string {
	for {
		fmt.Print("Enter priority (High, Medium, Low, or leave empty for Medium): ")
		p, _ := reader.ReadString('\n')
		p = strings.TrimSpace(p)
		p = strings.ToLower(p)

		switch p {
		case "high":
			return "High"
		case "medium", "":
			return "Medium"
		case "low":
			return "Low"
		default:
			fmt.Println("Invalid priority. Please enter High, Medium, or Low.")
		}
	}
}

func getDueDateInput(reader *bufio.Reader) time.Time {
	for {
		fmt.Print("Enter due date (YYYY-MM-DD, or leave empty for no due date): ")
		dateStr, _ := reader.ReadString('\n')
		dateStr = strings.TrimSpace(dateStr)

		if dateStr == "" {
			return time.Time{}
		}

		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			fmt.Println("Invalid date format. Please use YYYY-MM-DD.")
			continue
		}
		return t
	}
}

func listTodos(todoList []Todo) {
	if len(todoList) == 0 {
		fmt.Println("No todos to display.")
		return
	}

	fmt.Println("\n--- Your Todos ---")
	for _, todo := range todoList {
		status := "Pending"
		if todo.Completed {
			status = "Completed"
		}
		dueDateStr := "N/A"
		if !todo.DueDate.IsZero() {
			dueDateStr = todo.DueDate.Format("2006-01-02")
		}
		fmt.Printf("ID: %d | Desc: %s | Status: %s | Priority: %s | Due: %s | Created: %s\n",
			todo.ID, todo.Description, status, todo.Priority, dueDateStr, todo.CreatedAt.Format("2006-01-02"))
	}
}

func markComplete(reader *bufio.Reader) {
	listTodos(todos)
	if len(todos) == 0 {
		fmt.Print("\nPress Enter to continue...")
		reader.ReadString('\n')
		return
	}

	fmt.Print("Enter the ID of the todo to mark as complete: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid ID. Please enter a number.")
		time.Sleep(1 * time.Second)
		return
	}

	found := false
	for i := range todos {
		if todos[i].ID == id {
			todos[i].Completed = true
			found = true
			break
		}
	}

	if found {
		saveTodos()
		fmt.Println("Todo marked as complete!")
	} else {
		fmt.Println("Todo not found.")
	}
	time.Sleep(1 * time.Second)
}

func deleteTodo(reader *bufio.Reader) {
	listTodos(todos)
	if len(todos) == 0 {
		fmt.Print("\nPress Enter to continue...")
		reader.ReadString('\n')
		return
	}

	fmt.Print("Enter the ID of the todo to delete: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid ID. Please enter a number.")
		time.Sleep(1 * time.Second)
		return
	}

	originalLen := len(todos)
	newTodos := []Todo{}
	for _, todo := range todos {
		if todo.ID != id {
			newTodos = append(newTodos, todo)
		}
	}
	todos = newTodos

	if len(todos) < originalLen {
		saveTodos()
		fmt.Println("Todo deleted successfully!")
	} else {
		fmt.Println("Todo not found.")
	}
	time.Sleep(1 * time.Second)
}

func editTodo(reader *bufio.Reader) {
	listTodos(todos)
	if len(todos) == 0 {
		fmt.Print("\nPress Enter to continue...")
		reader.ReadString('\n')
		return
	}

	fmt.Print("Enter the ID of the todo to edit: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid ID. Please enter a number.")
		time.Sleep(1 * time.Second)
		return
	}

	foundIndex := -1
	for i, todo := range todos {
		if todo.ID == id {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		fmt.Println("Todo not found.")
		time.Sleep(1 * time.Second)
		return
	}

	fmt.Printf("Current Description: %s\n", todos[foundIndex].Description)
	fmt.Print("Enter new description (leave empty to keep current): ")
	newDesc, _ := reader.ReadString('\n')
	newDesc = strings.TrimSpace(newDesc)
	if newDesc != "" {
		todos[foundIndex].Description = newDesc
	}

	fmt.Printf("Current Priority: %s\n", todos[foundIndex].Priority)
	fmt.Print("Enter new priority (High, Medium, Low, or leave empty to keep current): ")
	newPriority := getPriorityInputForEdit(reader)
	if newPriority != "" {
		todos[foundIndex].Priority = newPriority
	}

	dueDateStr := "N/A"
	if !todos[foundIndex].DueDate.IsZero() {
		dueDateStr = todos[foundIndex].DueDate.Format("2006-01-02")
	}
	fmt.Printf("Current Due Date: %s\n", dueDateStr)
	fmt.Print("Enter new due date (YYYY-MM-DD, 'clear' to remove, or leave empty to keep current): ")
	newDueDateStr, _ := reader.ReadString('\n')
	newDueDateStr = strings.TrimSpace(newDueDateStr)

	if newDueDateStr != "" {
		if strings.ToLower(newDueDateStr) == "clear" {
			todos[foundIndex].DueDate = time.Time{}
		} else {
			t, err := time.Parse("2006-01-02", newDueDateStr)
			if err != nil {
				fmt.Println("Invalid date format. Keeping current due date.")
			} else {
				todos[foundIndex].DueDate = t
			}
		}
	}

	saveTodos()
	fmt.Println("Todo updated successfully!")
	time.Sleep(1 * time.Second)
}

func getPriorityInputForEdit(reader *bufio.Reader) string {
	for {
		p, _ := reader.ReadString('\n')
		p = strings.TrimSpace(p)
		p = strings.ToLower(p)

		switch p {
		case "high":
			return "High"
		case "medium":
			return "Medium"
		case "low":
			return "Low"
		case "":
			return ""
		default:
			fmt.Println("Invalid priority. Please enter High, Medium, Low, or leave empty.")
		}
	}
}

func filterTodosMenu(reader *bufio.Reader) {
	for {
		clearScreen()
		fmt.Println("--- Filter Todos ---")
		fmt.Println("1. Filter by Status (Pending/Completed/All)")
		fmt.Println("2. Filter by Priority (High/Medium/Low/All)")
		fmt.Println("3. Back to Main Menu")
		fmt.Print("Choose an option: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			filterByStatus(reader)
			fmt.Print("\nPress Enter to continue...")
			reader.ReadString('\n')
		case "2":
			filterByPriority(reader)
			fmt.Print("\nPress Enter to continue...")
			reader.ReadString('\n')
		case "3":
			return
		default:
			fmt.Println("Invalid option. Please try again.")
			time.Sleep(1 * time.Second)

// Additional implementation at 2025-06-19 23:09:44
