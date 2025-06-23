package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run merge_csv.go <output_file.csv> <input_file1.csv> [input_file2.csv ...]")
		os.Exit(1)
	}

	outputFilePath := os.Args[1]
	inputFilesPaths := os.Args[2:]

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("Error creating output file %s: %v\n", outputFilePath, err)
		os.Exit(1)
	}
	defer outputFile.Close()

	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()

	headerWritten := false

	for _, inputFilePath := range inputFilesPaths {
		inputFile, err := os.Open(inputFilePath)
		if err != nil {
			fmt.Printf("Error opening input file %s: %v\n", inputFilePath, err)
			continue
		}

		csvReader := csv.NewReader(inputFile)

		if !headerWritten {
			header, err := csvReader.Read()
			if err != nil {
				inputFile.Close()
				fmt.Printf("Error reading header from %s: %v\n", inputFilePath, err)
				os.Exit(1)
			}
			if err := csvWriter.Write(header); err != nil {
				inputFile.Close()
				fmt.Printf("Error writing header to output file: %v\n", err)
				os.Exit(1)
			}
			headerWritten = true
		} else {
			_, err := csvReader.Read()
			if err != nil && err != io.EOF {
				fmt.Printf("Warning: Could not skip header in %s: %v\n", inputFilePath, err)
			}
		}

		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("Error reading record from %s: %v\n", inputFilePath, err)
				continue
			}
			if err := csvWriter.Write(record); err != nil {
				inputFile.Close()
				fmt.Printf("Error writing record to output file: %v\n", err)
				os.Exit(1)
			}
		}
		inputFile.Close()
	}

	fmt.Printf("Successfully merged CSV files into %s\n", outputFilePath)
}

// Additional implementation at 2025-06-23 02:03:34
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func mergeCSVFiles(inputPaths []string, outputPath string) error {
	if len(inputPaths) == 0 {
		return fmt.Errorf("no input files provided")
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	firstFile := true

	for _, inputPath := range inputPaths {
		inputFile, err := os.Open(inputPath)
		if err != nil {
			return fmt.Errorf("failed to open input file %s: %w", inputPath, err)
		}
		// Defer closing the input file until the current loop iteration finishes
		defer inputFile.Close()

		reader := csv.NewReader(inputFile)

		if firstFile {
			header, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					fmt.Fprintf(os.Stderr, "Warning: First input file %s is empty or has no header.\n", inputPath)
				} else {
					return fmt.Errorf("failed to read header from %s: %w", inputPath, err)
				}
			} else {
				if err := writer.Write(header); err != nil {
					return fmt.Errorf("failed to write header to output file: %w", err)
				}
			}
			firstFile = false
		} else {
			// For subsequent files, read and discard the header row
			_, err := reader.Read()
			if err != nil && err != io.EOF {
				return fmt.Errorf("failed to read (and skip) header from %s: %w", inputPath, err)
			}
		}

		// Read and write remaining records
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("failed to read record from %s: %w", inputPath, err)
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("failed to write record to output file: %w", err)
			}
		}
	}

	return nil
}

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input_file1.csv> [input_file2.csv...] <output_file.csv>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	outputPath := args[len(args)-1]
	inputPaths := args[:len(args)-1]

	fmt.Printf("Merging %d files into %s...\n", len(inputPaths), outputPath)

	if err := mergeCSVFiles(inputPaths, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error merging CSV files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("CSV files merged successfully.")
}

// Additional implementation at 2025-06-23 02:04:32
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func mergeCSVFiles(inputPaths []string, outputPath string, includeSourceColumn bool) error {
	if len(inputPaths) == 0 {
		return fmt.Errorf("no input files specified")
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
	}
	defer outputFile.Close()

	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()

	headerWritten := false

	for _, inputPath := range inputPaths {
		inputFile, err := os.Open(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to open input file %s: %v. Skipping.\n", inputPath, err)
			continue
		}
		// Defer closing the input file for the current iteration
		defer func(f *os.File) {
			if f != nil {
				f.Close()
			}
		}(inputFile)

		csvReader := csv.NewReader(inputFile)

		sourceFileName := filepath.Base(inputPath)

		if !headerWritten {
			header, err := csvReader.Read()
			if err != nil {
				return fmt.Errorf("failed to read header from %s: %w", inputPath, err)
			}
			if includeSourceColumn {
				header = append(header, "SourceFile")
			}
			if err := csvWriter.Write(header); err != nil {
				return fmt.Errorf("failed to write header to output: %w", err)
			}
			headerWritten = true
		} else {
			// For subsequent files, read and discard their header
			_, err := csvReader.Read()
			if err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "Warning: Failed to read header from subsequent file %s: %v. Assuming no header or empty file.\n", inputPath, err)
			}
		}

		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to read record from %s: %v. Skipping record.\n", inputPath, err)
				continue
			}

			if includeSourceColumn {
				record = append(record, sourceFileName)
			}

			if err := csvWriter.Write(record); err != nil {
				return fmt.Errorf("failed to write record to output: %w", err)
			}
		}
	}

	return nil
}

func main() {
	var outputFilePath string
	var inputFiles string
	var includeSourceColumn bool

	flag.StringVar(&outputFilePath, "output", "merged.csv", "Path to the output merged CSV file")
	flag.StringVar(&outputFilePath, "o", "merged.csv", "Path to the output merged CSV file (shorthand)")
	flag.StringVar(&inputFiles, "input", "", "Comma-separated list of input CSV files to merge")
	flag.StringVar(&inputFiles, "i", "", "Comma-separated list of input CSV files to merge (shorthand)")
	flag.BoolVar(&includeSourceColumn, "source-column", false, "Add a 'SourceFile' column indicating the original file for each row")
	flag.BoolVar(&includeSourceColumn, "s", false, "Add a 'SourceFile' column indicating the original file for each row (shorthand)")

	flag.Parse()

	if inputFiles == "" {
		fmt.Println("Error: No input files specified. Use -input or -i flag.")
		flag.Usage()
		os.Exit(1)
	}

	inputPaths := strings.Split(inputFiles, ",")
	var cleanInputPaths []string
	for _, path := range inputPaths {
		trimmedPath := strings.TrimSpace(path)
		if trimmedPath != "" {
			cleanInputPaths = append(cleanInputPaths, trimmedPath)
		}
	}
	inputPaths = cleanInputPaths

	if len(inputPaths) == 0 {
		fmt.Println("Error: No valid input files found after parsing. Check your -input argument.")
		flag.Usage()
		os.Exit(1)
	}

	err := mergeCSVFiles(inputPaths, outputFilePath, includeSourceColumn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error merging CSV files: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully merged %d CSV files into %s\n", len(inputPaths), outputFilePath)
}

// Additional implementation at 2025-06-23 02:05:51
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type stringArrayFlag []string

func (s *stringArrayFlag) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringArrayFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func processInputFile(inputFile string, csvWriter *csv.Writer, headerWritten *bool) error {
	inputFileHandle, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("error opening input file %s: %w", inputFile, err)
	}
	defer inputFileHandle.Close()

	csvReader := csv.NewReader(inputFileHandle)

	if !*headerWritten {
		header, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("input file %s is empty, skipping", inputFile)
			}
			return fmt.Errorf("error reading header from %s: %w", inputFile, err)
		}
		if err := csvWriter.Write(header); err != nil {
			return fmt.Errorf("error writing header to output file: %w", err)
		}
		*headerWritten = true
	} else {
		_, err := csvReader.Read()
		if err != nil && err != io.EOF {
			return fmt.Errorf("warning: could not skip header in %s: %w", inputFile, err)
		}
	}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading record from %s: %w", inputFile, err)
		}
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("error writing record to output file: %w", err)
		}
	}
	return nil
}

func main() {
	var inputFiles stringArrayFlag
	var outputFile string

	flag.Var(&inputFiles, "in", "Input CSV file(s). Can be specified multiple times, e.g., -in file1.csv -in file2.csv")
	flag.StringVar(&outputFile, "out", "merged.csv", "Output CSV file name.")
	flag.Parse()

	if len(inputFiles) == 0 {
		fmt.Println("Error: No input files specified. Use -in flag.")
		flag.Usage()
		os.Exit(1)
	}

	if outputFile == "" {
		fmt.Println("Error: Output file not specified. Use -out flag.")
		flag.Usage()
		os.Exit(1)
	}

	outputFileHandle, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file %s: %v\n", outputFile, err)
		os.Exit(1)
	}
	defer outputFileHandle.Close()

	csvWriter := csv.NewWriter(outputFileHandle)
	defer csvWriter.Flush()

	headerWritten := false

	for _, inputFile := range inputFiles {
		err := processInputFile(inputFile, csvWriter, &headerWritten)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		csvWriter.Flush()
	}

	if !headerWritten {
		fmt.Println("Warning: No data was written to the output file, possibly all input files were empty or had errors.")
	} else {
		fmt.Printf("Successfully merged CSV files into %s\n", outputFile)
	}
}