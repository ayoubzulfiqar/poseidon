package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run main.go <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file %q: %v\n", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	urlRegex := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		foundURLs := urlRegex.FindAllString(line, -1)

		for _, url := range foundURLs {
			fmt.Println(url)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from file %q: %v\n", filename, err)
		os.Exit(1)
	}
}

// Additional implementation at 2025-06-21 00:02:51
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename> [--scheme <scheme>]")
		os.Exit(1)
	}

	filePath := os.Args[1]

	var schemeFilter string
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--scheme" {
			if i+1 < len(os.Args) {
				schemeFilter = os.Args[i+1]
				i++ // Skip the next argument as it's the scheme value
			} else {
				fmt.Println("Error: --scheme requires a value.")
				os.Exit(1)
			}
		}
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Regular expression to find URLs. This regex is designed to capture common URL patterns
	// including http, https, and ftp schemes, followed by valid URL characters.
	re := regexp.MustCompile(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]*[-A-Za-z0-9+&@#/%=~_|]`)

	allURLs := re.FindAllString(string(content), -1)

	// Store all unique URLs found, regardless of scheme filter
	allUniqueURLs := make(map[string]bool)
	for _, url := range allURLs {
		trimmedURL := strings.TrimSpace(url)
		allUniqueURLs[trimmedURL] = true
	}

	// Prepare the final set of URLs to display based on the scheme filter
	displayURLs := make(map[string]bool)
	if schemeFilter != "" {
		for url := range allUniqueURLs {
			if strings.HasPrefix(url, schemeFilter+"://") {
				displayURLs[url] = true
			}
		}
	} else {
		// If no scheme filter, display all unique URLs
		displayURLs = allUniqueURLs
	}

	fmt.Printf("--- URL Extraction Report for '%s' ---\n", filePath)
	fmt.Printf("Total URLs found (including duplicates): %d\n", len(allURLs))
	fmt.Printf("Total unique URLs found (before any scheme filter): %d\n", len(allUniqueURLs))

	if schemeFilter != "" {
		fmt.Printf("Total unique URLs found (filtered by scheme '%s'): %d\n", schemeFilter, len(displayURLs))
		fmt.Printf("Unique URLs (filtered by scheme '%s'):\n", schemeFilter)
	} else {
		fmt.Printf("Unique URLs:\n")
	}

	if len(displayURLs) == 0 {
		fmt.Println("  No URLs found matching the criteria.")
	} else {
		// Iterate and print the URLs. Order is not guaranteed as it's from a map.
		for url := range displayURLs {
			fmt.Printf("  %s\n", url)
		}
	}
	fmt.Println("---------------------------------------")
}