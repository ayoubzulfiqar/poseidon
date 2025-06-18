package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func main() {
	decoder := charmap.Windows1252.NewDecoder()
	reader := transform.NewReader(os.Stdin, decoder)

	_, err := io.Copy(os.Stdout, reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting: %v\n", err)
		os.Exit(1)
	}
}

// Additional implementation at 2025-06-18 00:55:41
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func main() {
	inputFilePath := flag.String("input", "", "Path to the input Windows-1252 file")
	outputFilePath := flag.String("output", "", "Path to the output UTF-8 file")

	flag.Parse()

	if *inputFilePath == "" {
		fmt.Fprintf(os.Stderr, "Error: Input file path is required.\n")
		flag.Usage()
		os.Exit(1)
	}
	if *outputFilePath == "" {
		fmt.Fprintf(os.Stderr, "Error: Output file path is required.\n")
		flag.Usage()
		os.Exit(1)
	}

	inputFile, err := os.Open(*inputFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file %q: %v\n", *inputFilePath, err)
		os.Exit(1)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(*outputFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file %q: %v\n", *outputFilePath, err)
		os.Exit(1)
	}
	defer outputFile.Close()

	decoder := charmap.Windows1252.NewDecoder()

	reader := transform.NewReader(inputFile, decoder)

	bytesCopied, err := io.Copy(outputFile, reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting and writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %q to UTF-8 in %q. Copied %d bytes.\n", *inputFilePath, *outputFilePath, bytesCopied)
}

// Additional implementation at 2025-06-18 00:56:35
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func convertWindows1252ToUTF8(r io.Reader, w io.Writer) error {
	decoder := charmap.Windows1252.NewDecoder()
	reader := transform.NewReader(r, decoder)

	_, err := io.Copy(w, reader)
	if err != nil {
		return fmt.Errorf("failed to copy transformed data: %w", err)
	}
	return nil
}

func main() {
	inputFilePath := flag.String("in", "", "Input file path (default: stdin)")
	outputFilePath := flag.String("out", "", "Output file path (default: stdout)")
	flag.Parse()

	var inputReader io.Reader
	var outputWriter io.Writer
	var err error

	if *inputFilePath == "" {
		inputReader = os.Stdin
		fmt.Fprintln(os.Stderr, "Reading from stdin...")
	} else {
		file, err := os.Open(*inputFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening input file %s: %v\n", *inputFilePath, err)
			os.Exit(1)
		}
		defer file.Close()
		inputReader = file
		fmt.Fprintf(os.Stderr, "Reading from file: %s\n", *inputFilePath)
	}

	if *outputFilePath == "" {
		outputWriter = os.Stdout
		fmt.Fprintln(os.Stderr, "Writing to stdout...")
	} else {
		file, err := os.Create(*outputFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file %s: %v\n", *outputFilePath, err)
			os.Exit(1)
		}
		defer file.Close()
		outputWriter = file
		fmt.Fprintf(os.Stderr, "Writing to file: %s\n", *outputFilePath)
	}

	err = convertWindows1252ToUTF8(inputReader, outputWriter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Conversion failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Conversion complete.")
}

// Additional implementation at 2025-06-18 00:57:37
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func convertWin1252ToUTF8(r io.Reader, w io.Writer) error {
	decoder := charmap.Windows1252.NewDecoder()
	reader := transform.NewReader(r, decoder)

	_, err := io.Copy(w, reader)
	if err != nil {
		return fmt.Errorf("error during conversion: %w", err)
	}
	return nil
}

func main() {
	inputFilePath := flag.String("input", "", "Path to the input file (Windows-1252 encoded). If empty, reads from stdin.")
	outputFilePath := flag.String("output", "", "Path to the output file (UTF-8 encoded). If empty, writes to stdout.")
	flag.Parse()

	var inputReader io.Reader
	var outputWriter io.Writer
	var err error

	if *inputFilePath != "" {
		inputFile, err := os.Open(*inputFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening input file %s: %v\n", *inputFilePath, err)
			os.Exit(1)
		}
		defer inputFile.Close()
		inputReader = inputFile
	} else {
		inputReader = os.Stdin
	}

	if *outputFilePath != "" {
		outputFile, err := os.Create(*outputFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file %s: %v\n", *outputFilePath, err)
			os.Exit(1)
		}
		defer outputFile.Close()
		outputWriter = outputFile
	} else {
		outputWriter = os.Stdout
	}

	if err := convertWin1252ToUTF8(inputReader, outputWriter); err != nil {
		fmt.Fprintf(os.Stderr, "Conversion failed: %v\n", err)
		os.Exit(1)
	}
}