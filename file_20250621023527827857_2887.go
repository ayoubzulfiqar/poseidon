package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	urlRegex := regexp.MustCompile(`(https?|ftp|file)://[^\s/$.?#].[^\s]*`)

	foundURLs := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := urlRegex.FindAllString(line, -1)
		for _, url := range matches {
			foundURLs[url] = true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	for url := range foundURLs {
		fmt.Println(url)
	}
}

// Additional implementation at 2025-06-21 02:36:11
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filepath>\n", os.Args[0])
		os.Exit(1)
	}

	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Regular expression to find URLs.
	// This regex broadly matches http(s), ftp, and file URLs.
	// It's designed to be reasonably comprehensive without being overly complex.
	re := regexp.MustCompile(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]*[-A-Za-z0-9+&@#/%=~_|]`)

	// Use a map to store unique URLs
	uniqueURLs := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		foundURLs := re.FindAllString(line, -1)
		for _, url := range foundURLs {
			uniqueURLs[url] = true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Extract unique URLs into a slice for sorting (additional functionality: sorted output)
	var sortedURLs []string
	for url := range uniqueURLs {
		sortedURLs = append(sortedURLs, url)
	}
	sort.Strings(sortedURLs) // Sort them alphabetically

	// Print all unique URLs
	for _, url := range sortedURLs {
		fmt.Println(url)
	}
}

// Additional implementation at 2025-06-21 02:36:59
package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filepath>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		os.Exit(1)
	}

	urlRegex := regexp.MustCompile(`(https?|ftp|file):\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

	foundURLs := urlRegex.FindAllString(string(content), -1)

	uniqueURLsMap := make(map[string]bool)
	var uniqueURLsList []string

	for _, url := range foundURLs {
		if _, exists := uniqueURLsMap[url]; !exists {
			uniqueURLsMap[url] = true
			uniqueURLsList = append(uniqueURLsList, url)
		}
	}

	sort.Strings(uniqueURLsList)

	if len(uniqueURLsList) == 0 {
		fmt.Println("No URLs found in the file.")
		return
	}

	fmt.Println("Found URLs:")
	for _, url := range uniqueURLsList {
		fmt.Println(url)
	}
}