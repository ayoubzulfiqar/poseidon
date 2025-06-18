package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}

	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		fmt.Printf("%d\t%s\n", lineNumber, scanner.Text())
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file %s: %v", filePath, err)
	}
}

// Additional implementation at 2025-06-17 23:35:54
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func addLineNumbers(reader io.Reader, writer io.Writer, skipEmpty bool) error {
	scanner := bufio.NewScanner(reader)
	lineNum := 1
	bufWriter := bufio.NewWriter(writer)
	defer bufWriter.Flush()

	for scanner.Scan() {
		line := scanner.Text()
		if skipEmpty && len(line) == 0 {
			_, err := bufWriter.WriteString(line + "\n")
			if err != nil {
				return fmt.Errorf("failed to write empty line: %w", err)
			}
			continue
		}

		formattedLine := fmt.Sprintf("%6d  %s\n", lineNum, line)
		_, err := bufWriter.WriteString(formattedLine)
		if err != nil {
			return fmt.Errorf("failed to write line %d: %w", lineNum, err)
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}
	return nil
}

func main() {
	var inputFile string
	var outputFile string
	var skipEmpty bool

	flag.StringVar(&inputFile, "in", "", "Input file path (default: stdin)")
	flag.StringVar(&outputFile, "out", "", "Output file path (default: stdout)")
	flag.BoolVar(&skipEmpty, "s", false, "Skip numbering and incrementing for empty lines")
	flag.BoolVar(&skipEmpty, "skip-empty", false, "Skip numbering and incrementing for empty lines (long form)")

	flag.Parse()

	var reader io.Reader = os.Stdin
	var writer io.Writer = os.Stdout
	var err error

	if inputFile != "" {
		file, err := os.Open(inputFile)
		if err != nil {
			log.Fatalf("Error opening input file %s: %v", inputFile, err)
		}
		defer file.Close()
		reader = file
	}

	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			log.Fatalf("Error creating output file %s: %v", outputFile, err)
		}
		defer func(f *os.File) {
			if cerr := f.Close(); cerr != nil {
				log.Printf("Error closing output file %s: %v", f.Name(), cerr)
			}
		}(file)
		writer = file
	}

	if err := addLineNumbers(reader, writer, skipEmpty); err != nil {
		log.Fatalf("Error processing file: %v", err)
	}
}