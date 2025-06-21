package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <program_name> <filename>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	wordCounts := make(map[string]int)
	scanner := bufio.NewScanner(file)

	wordRegex := regexp.MustCompile(`[^a-zA-Z0-9']+`)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ToLower(line)

		words := wordRegex.Split(line, -1)

		for _, word := range words {
			word = strings.TrimSpace(word)
			if word != "" {
				wordCounts[word]++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	type wordFrequency struct {
		Word  string
		Count int
	}

	var frequencies []wordFrequency
	for word, count := range wordCounts {
		frequencies = append(frequencies, wordFrequency{Word: word, Count: count})
	}

	sort.Slice(frequencies, func(i, j int) bool {
		if frequencies[i].Count != frequencies[j].Count {
			return frequencies[i].Count > frequencies[j].Count
		}
		return frequencies[i].Word < frequencies[j].Word
	})

	for _, wf := range frequencies {
		fmt.Printf("%s: %d\n", wf.Word, wf.Count)
	}
}

// Additional implementation at 2025-06-20 23:51:14
package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

// WordFreq represents a word and its frequency.
type WordFreq struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

func main() {
	filePath := flag.String("file", "", "Path to the text file (required)")
	topN := flag.Int("top", 0, "Display only the top N most frequent words (0 for all)")
	caseSensitive := flag.Bool("case-sensitive", false, "Perform case-sensitive word counting")
	outputFormat := flag.String("format", "plain", "Output format: plain, csv, or json")

	flag.Parse()

	if *filePath == "" {
		fmt.Println("Error: --file is required.")
		flag.Usage()
		os.Exit(1)
	}

	wordCounts, err := countWordFrequencies(*filePath, *caseSensitive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing file: %v\n", err)
		os.Exit(1)
	}

	sortedWords := sortWordFrequencies(wordCounts)

	if *topN > 0 && *topN < len(sortedWords) {
		sortedWords = sortedWords[:*topN]
	}

	err = printResults(sortedWords, *outputFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error printing results: %v\n", err)
		os.Exit(1)
	}
}

// countWordFrequencies reads a file and counts word occurrences.
func countWordFrequencies(filePath string, caseSensitive bool) (map[string]int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	wordCounts := make(map[string]int)
	// Regex to find words: sequences of letters, numbers, and apostrophes.
	// This helps in stripping punctuation from words like "hello," or "it's".
	wordRegex := regexp.MustCompile(`[a-zA-Z0-9']+`)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords) // Splits the input by whitespace

	for scanner.Scan() {
		token := scanner.Text()
		// Extract actual words from the token, handling attached punctuation.
		matches := wordRegex.FindAllString(token, -1)
		for _, word := range matches {
			if !caseSensitive {
				word = strings.ToLower(word)
			}
			wordCounts[word]++
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return wordCounts, nil
}

// sortWordFrequencies converts the map to a slice of WordFreq and sorts it.
func sortWordFrequencies(wordCounts map[string]int) []WordFreq {
	var sortedWords []WordFreq
	for word, count := range wordCounts {
		sortedWords = append(sortedWords, WordFreq{Word: word, Count: count})
	}

	sort.Slice(sortedWords, func(i, j int) bool {
		if sortedWords[i].Count != sortedWords[j].Count {
			return sortedWords[i].Count > sortedWords[j].Count // Sort by count descending
		}
		return sortedWords[i].Word < sortedWords[j].Word // Then by word alphabetically ascending
	})

	return sortedWords
}

// printResults outputs the word frequencies based on the specified format.
func printResults(words []WordFreq, format string) error {
	switch format {
	case "plain":
		for _, wf := range words {
			fmt.Printf("%-20s %d\n", wf.Word, wf.Count)
		}
	case "csv":
		writer := csv.NewWriter(os.Stdout)
		writer.Write([]string{"Word", "Count"}) // CSV header
		for _, wf := range words {
			err := writer.Write([]string{wf.Word, fmt.Sprintf("%d", wf.Count)})
			if err != nil {
				return fmt.Errorf("error writing CSV row: %w", err)
			}
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			return fmt.Errorf("error flushing CSV writer: %w", err)
		}
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ") // Pretty print JSON
		if err := encoder.Encode(words); err != nil {
			return fmt.Errorf("error encoding JSON: %w", err)
		}
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
	return nil
}