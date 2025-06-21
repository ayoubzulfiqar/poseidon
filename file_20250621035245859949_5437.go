package main

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
)

func main() {
	initialContent, err := clipboard.ReadAll()
	if err != nil {
		fmt.Printf("Warning: Could not read initial clipboard content: %v\n", err)
		initialContent = ""
	}
	lastContent := initialContent

	fmt.Println("Monitoring clipboard changes...")
	fmt.Println("Press Ctrl+C to exit.")
	fmt.Printf("Current clipboard content (on start): %s\n", lastContent)

	for {
		currentContent, err := clipboard.ReadAll()
		if err != nil {
			fmt.Printf("Error reading clipboard: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if currentContent != lastContent {
			fmt.Printf("Clipboard changed: %s\n", currentContent)
			lastContent = currentContent
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// Additional implementation at 2025-06-21 03:53:35
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"golang.design/x/clipboard"
)

const (
	logFilePath     = "clipboard_monitor.log"
	historyCapacity = 10 // Store last 10 unique clipboard entries
	debouncePeriod  = 500 * time.Millisecond // Ignore identical content copied within this period
)

var (
	logger         *log.Logger
	lastContent    []byte
	lastChangeTime time.Time
	history        []string
)

func init() {
	// Initialize clipboard
	err := clipboard.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize clipboard: %v\n", err)
		os.Exit(1)
	}

	// Set up logging to a file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
		os.Exit(1)
	}
	logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	fmt.Printf("Clipboard monitor started. Logging to %s\n", logFilePath)

	history = make([]string, 0, historyCapacity)
}

func main() {
	// Watch for clipboard changes of type Text
	clipboardChan := clipboard.Watch(context.Background(), clipboard.TypeText)

	for content := range clipboardChan { // This loop blocks until new content is available
		// Convert content to string for comparison and logging
		currentContent := string(content)

		// Debounce identical content
		if currentContent == string(lastContent) && time.Since(lastChangeTime) < debouncePeriod {
			// Content is the same and within debounce period, ignore
			continue
		}

		// Update last content and time
		lastContent = content
		lastChangeTime = time.Now()

		// Log the change
		logMessage := fmt.Sprintf("Clipboard changed: %s", currentContent)
		logger.Println(logMessage)
		fmt.Println(logMessage) // Also print to console

		// Update history
		updateHistory(currentContent)

		// Optionally print current history
		printHistory()
	}
}

// updateHistory adds a new unique entry to the history, maintaining capacity.
func updateHistory(newEntry string) {
	// Check if the new entry is already the latest in history
	if len(history) > 0 && history[len(history)-1] == newEntry {
		return // Already the latest, no need to add again
	}

	// Remove duplicates from history if newEntry already exists
	for i := 0; i < len(history); i++ {
		if history[i] == newEntry {
			history = append(history[:i], history[i+1:]...)
			break
		}
	}

	// Add new entry to the end
	history = append(history, newEntry)

	// Trim history if it exceeds capacity
	if len(history) > historyCapacity {
		history = history[len(history)-historyCapacity:]
	}
}

// printHistory prints the current in-memory history to console.
func printHistory() {
	fmt.Println("\n--- Clipboard History ---")
	if len(history) == 0 {
		fmt.Println("No history yet.")
		return
	}
	for i, entry := range history {
		fmt.Printf("%d: %s\n", i+1, entry)
	}
	fmt.Println("-------------------------\n")
}