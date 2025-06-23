package main

import (
	"fmt"
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Welcome to the Home Page!")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Additional implementation at 2025-06-23 02:17:00
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

var (
	books    = make(map[int]Book)
	nextID   = 1
	booksMux sync.Mutex // Mutex to protect access to the books map
)

func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func writeErrorResponse(w http.ResponseWriter, status int, message string) {
	writeJSONResponse(w, status, map[string]string{"error": message})
}

func getBooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/books" {
		writeErrorResponse(w, http.StatusNotFound, "Endpoint not found")
		return
	}
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	booksMux.Lock()
	defer booksMux.Unlock()

	// Convert map to slice for JSON encoding
	bookList := make([]Book, 0, len(books))
	for _, book := range books {
		bookList = append(bookList, book)
	}
	writeJSONResponse(w, http.StatusOK, bookList)
}

func getBookByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[1] != "books" {
		writeErrorResponse(w, http.StatusNotFound, "Invalid URL path")
		return
	}

	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	booksMux.Lock()
	defer booksMux.Unlock()

	book, ok := books[id]
	if !ok {
		writeErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}
	writeJSONResponse(w, http.StatusOK, book)
}

func createBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/books" {
		writeErrorResponse(w, http.StatusNotFound, "Endpoint not found")
		return
	}
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var newBook Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if newBook.Title == "" || newBook.Author == "" || newBook.Year == 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Title, Author, and Year are required")
		return
	}

	booksMux.Lock()
	defer booksMux.Unlock()

	newBook.ID = nextID
	books[nextID] = newBook
	nextID++

	writeJSONResponse(w, http.StatusCreated, newBook)
}

func updateBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[1] != "books" {
		writeErrorResponse(w, http.StatusNotFound, "Invalid URL path")
		return
	}

	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	var updatedBook Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	booksMux.Lock()
	defer booksMux.Unlock()

	_, ok := books[id]
	if !ok {
		writeErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	// Ensure the ID from the URL is used, not the one from the body
	updatedBook.ID = id
	books[id] = updatedBook
	writeJSONResponse(w, http.StatusOK, updatedBook)
}

func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[1] != "books" {
		writeErrorResponse(w, http.StatusNotFound, "Invalid URL path")
		return
	}

	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	booksMux.Lock()
	defer booksMux.Unlock()

	_, ok := books[id]
	if !ok {
		writeErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	delete(books, id)
	writeJSONResponse(w, http.StatusNoContent, nil) // No content to return for successful deletion
}

func main() {
	// Initialize some dummy data
	booksMux.Lock()
	books[nextID] = Book{ID: nextID, Title: "The Go Programming Language", Author: "Alan A. A. Donovan and Brian W. Kernighan", Year: 2015}
	nextID++
	books[nextID] = Book{ID: nextID, Title: "Clean Code", Author: "Robert C. Martin", Year: 2008}
	nextID++
	booksMux.Unlock()

	mux := http.NewServeMux()

	// Route for /books (GET all, POST new)
	mux.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getBooksHandler(w, r)
		case http.MethodPost:
			createBookHandler(w, r)
		default:
			writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	// Route for /books/{id} (GET by ID, PUT update, DELETE)
	mux.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		// This handler catches /books/ and anything deeper like /books/123
		// We need to differentiate based on the path length
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) == 2 && parts[0] == "books" { // e.g., /books/123
			switch r.Method {
			case http.MethodGet:
				getBookByIDHandler(w, r)
			case http.MethodPut:
				updateBookHandler(w, r)
			case http.MethodDelete:
				deleteBookHandler(w, r)
			default:
				writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		} else {
			writeErrorResponse(w, http.StatusNotFound, "Endpoint not found")
		}
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("Server starting on port 8080...")
	log.Fatal(server.ListenAndServe())
}