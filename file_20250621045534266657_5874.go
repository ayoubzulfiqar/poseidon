package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading: %v\n", err)
			os.Exit(1)
		}

		if b == '\r' {
			nextByte, err := reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					writer.WriteByte('\r')
					break
				}
				fmt.Fprintf(os.Stderr, "Error reading after \\r: %v\n", err)
				os.Exit(1)
			}

			if nextByte == '\n' {
				writer.WriteByte('\n')
			} else {
				writer.WriteByte('\r')
				writer.WriteByte(nextByte)
			}
		} else {
			writer.WriteByte(b)
		}
	}
}

// Additional implementation at 2025-06-21 04:56:20
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func convertFile(filePath string, createBackup, verbose bool) error {
	if verbose {
		fmt.Printf("Processing %s...\n", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	originalContent := string(content)
	convertedContent := strings.ReplaceAll(originalContent, "\r\n", "\n")

	if originalContent == convertedContent {
		if verbose {
			fmt.Printf("No CRLF line endings found in %s. Skipping.\n", filePath)
		}
		return nil // No changes needed
	}

	if createBackup {
		backupPath := filePath + ".bak"
		if verbose {
			fmt.Printf("Creating backup %s...\n", backupPath)
		}
		err = os.WriteFile(backupPath, content, 0644)
		if err != nil {
			return fmt.Errorf("error creating backup for %s at %s: %w", filePath, backupPath, err)
		}
	}

	if verbose {
		fmt.Printf("Writing converted content to %s...\n", filePath)
	}
	err = os.WriteFile(filePath, []byte(convertedContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing converted content to %s: %w", filePath, err)
	}

	if verbose {
		fmt.Printf("Successfully converted %s.\n", filePath)
	}
	return nil
}

func main() {
	createBackup := flag.Bool("backup", false, "Create a .bak file before modifying the original file.")
	verbose := flag.Bool("v", false, "Enable verbose output.")
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Usage: go run your_script.go [options] <file_or_directory1> [file_or_directory2...]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	for _, filePath := range args {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", filePath, err)
			continue
		}

		if fileInfo.IsDir() {
			if *verbose {
				fmt.Printf("Walking directory %s...\n", filePath)
			}
			err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error walking path %s: %v\n", path, err)
					return nil // Continue walking other paths, don't stop the entire walk
				}
				if !info.IsDir() {
					// Only process regular files
					if err := convertFile(path, *createBackup, *verbose); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to process file %s: %v\n", path, err)
					}
				}
				return nil
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error during directory walk for %s: %v\n", filePath, err)
			}
		} else {
			// It's a regular file
			if err := convertFile(filePath, *createBackup, *verbose); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to process file %s: %v\n", filePath, err)
			}
		}
	}
}

// Additional implementation at 2025-06-21 04:57:31
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"flag"
)

// processFile converts CRLF to LF for a single file.
// It can write to an output file or modify in-place with an optional backup.
func processFile(inputPath, outputPath string, inPlace, createBackup, verbose bool) error {
	if verbose {
		fmt.Printf("Processing file: %s\n", inputPath)
	}

	content, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", inputPath, err)
	}

	// Convert CRLF to LF
	convertedContent := strings.ReplaceAll(string(content), "\r\n", "\n")

	// If no conversion happened, no need to write
	if string(content) == convertedContent {
		if verbose {
			fmt.Printf("No CRLF found or no change needed for %s\n", inputPath)
		}
		return nil
	}

	targetPath := outputPath
	if inPlace {
		targetPath = inputPath
		if createBackup {
			backupPath := inputPath + ".bak"
			if verbose {
				fmt.Printf("Creating backup: %s -> %s\n", inputPath, backupPath)
			}
			if err := ioutil.WriteFile(backupPath, content, 0644); err != nil {
				return fmt.Errorf("failed to create backup file %s: %w", backupPath, err)
			}
		}
	} else if targetPath == "" {
		// This case should be caught by flag validation in main, but as a safeguard.
		return fmt.Errorf("output path must be specified when not using in-place mode")
	}

	if verbose {
		fmt.Printf("Writing converted content to: %s\n", targetPath)
	}
	if err := ioutil.WriteFile(targetPath, []byte(convertedContent), 0644); err != nil {
		return fmt.Errorf("failed to write converted content to %s: %w", targetPath, err)
	}

	return nil
}

// processDirectory recursively processes files in a directory.
func processDirectory(inputDir, outputDir string, inPlace, createBackup, recursive, verbose bool, includeExts map[string]bool) error {
	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %q: %w", path, err)
		}

		if info.IsDir() {
			if path == inputDir { // Don't skip the root directory itself
				return nil
			}
			if !recursive {
				if verbose {
					fmt.Printf("Skipping directory (not recursive): %s\n", path)
				}
				return filepath.SkipDir // Skip subdirectories if not recursive
			}
			// If recursive, continue walking
			return nil
		}

		// Check file extension
		ext := strings.ToLower(filepath.Ext(path))
		if len(includeExts) > 0 && !includeExts[ext] {
			if verbose {
				fmt.Printf("Skipping file (unmatched extension): %s\n", path)
			}
			return nil // Skip file if extension not in whitelist
		}

		// Determine output path for directory mode
		var currentOutputPath string
		if !inPlace {
			relPath, err := filepath.Rel(inputDir, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path for %s: %w", path, err)
			}
			currentOutputPath = filepath.Join(outputDir, relPath)
			// Ensure output directory exists for the current file
			outputFileDir := filepath.Dir(currentOutputPath)
			if err := os.MkdirAll(outputFileDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory %s: %w", outputFileDir, err)
			}
		}

		if err := processFile(path, currentOutputPath, inPlace, createBackup, verbose); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", path, err)
			// Don't stop the whole walk for one file error, just log it.
		}
		return nil
	})
}

func main() {
	var inputPath string
	var outputPath string
	var inPlace bool
	var createBackup bool
	var recursive bool
	var verbose bool
	var extensions string

	flag.StringVar(&inputPath, "input", "", "Input file or directory path.")
	flag.StringVar(&outputPath, "output", "", "Output file or directory path. Required if not using --in-place.")
	flag.BoolVar(&inPlace, "in-place", false, "Modify files in-place.")
	flag.BoolVar(&createBackup, "backup", false, "Create a .bak file before modifying in-place. Only applicable with --in-place.")
	flag.BoolVar(&recursive, "recursive", false, "Process files in subdirectories when input is a directory.")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output.")
	flag.StringVar(&extensions, "ext", "", "Comma-separated list of file extensions to process (e.g., .txt,.go). If empty, all files are processed.")

	flag.Parse()

	if inputPath == "" {
		fmt.Fprintf(os.Stderr, "Error: --input is required.\n")
		flag.Usage()
		os.Exit(1)
	}

	if !inPlace && outputPath == "" {
		fmt.Fprintf(os.Stderr, "Error: --output is required when not using --in-place.\n")
		flag.Usage()
		os.Exit(1)
	}

	if createBackup && !inPlace {
		fmt.Fprintf(os.Stderr, "Error: --backup can only be used with --in-place.\n")
		flag.Usage()
		os.Exit(1)
	}

	includeExts := make(map[string]bool)
	if extensions != "" {
		for _, ext := range strings.Split(extensions, ",") {
			ext = strings.TrimSpace(ext)
			if ext != "" {
				if !strings.HasPrefix(ext, ".") {
					ext = "." + ext // Ensure extension starts with a dot
				}
				includeExts[strings.ToLower(ext)] = true
			}
		}
	}

	info, err := os.Stat(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid input path %s: %v\n", inputPath, err)
		os.Exit(1)
	}

	if info.IsDir() {
		if !inPlace && outputPath != "" {
			// Ensure output directory exists if it's a new one for directory processing
			if err := os.MkdirAll(outputPath, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to create output directory %s: %v\n", outputPath, err)
				os.Exit(1)
			}
		}
		if err := processDirectory(inputPath, outputPath, inPlace, createBackup, recursive, verbose, includeExts); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing directory %s: %v\n", inputPath, err)
			os.Exit(1)
		}
	} else {
		// Input is a file
		if outputPath != "" && !inPlace {
			// If output is specified and not in-place, check if output path is a directory or a file
			outputInfo, err := os.Stat(outputPath)
			if err == nil && outputInfo.IsDir() {
				// If output is an existing directory, then the output file should be within it, with the same name as input
				outputPath = filepath.Join(outputPath, filepath.Base(inputPath))
			} else if os.IsNotExist(err) {
				// If output path doesn't exist, assume it's a file path. Ensure its parent directory exists.
				outputDir := filepath.Dir(outputPath)
				if outputDir != "." { // Don't try to create current directory
					if err := os.MkdirAll(outputDir, 0755); err != nil {
						fmt.Fprintf(os.Stderr, "Error: Failed to create output directory %s: %v\n", outputDir, err)
						os.Exit(1)
					}
				}
			} else if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Invalid output path %s: %v\n", outputPath, err)
				os.Exit(1)
			}
		}
		if err := processFile(inputPath, outputPath, inPlace, createBackup, verbose); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing file %s: %v\n", inputPath, err)
			os.Exit(1)
		}
	}

	if verbose {
		fmt.Println("Conversion complete.")
	}
}

// Additional implementation at 2025-06-21 04:58:07
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	inputPath = flag.String("input", "", "Path to the file or directory to process.")
	dryRun    = flag.Bool("dry-run", false, "Perform a dry run without making any changes.")
	overwrite = flag.Bool("overwrite", false, "Overwrite the original file(s). If false, creates new file(s) with '.lf' suffix.")
	backup    = flag.Bool("backup", false, "Create a '.bak' file before overwriting (only effective with -overwrite).")
	recursive = flag.Bool("recursive", false, "Process directories recursively.")
	verbose   = flag.Bool("verbose", false, "Print detailed information about processed files.")
)

func main() {
	flag.Parse()

	if *inputPath == "" {
		fmt.Println("Error: -input flag is required.")
		flag.Usage()
		os.Exit(1)
	}

	if *backup && !*overwrite {
		fmt.Println("Warning: -backup flag has no effect without -overwrite.")
	}

	info, err := os.Stat(*inputPath)
	if err != nil {
		fmt.Printf("Error: Could not access path '%s': %v\n", *inputPath, err)
		os.Exit(1)
	}

	if info.IsDir() {
		if !*recursive {
			fmt.Printf("Error: '%s' is a directory. Use -recursive to process directories.\n", *inputPath)
			os.Exit(1)
		}
		err = filepath.Walk(*inputPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Error walking path %s: %v\n", path, err)
				return err
			}
			if !info.IsDir() {
				return processFile(path)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error during directory traversal: %v\n", err)
			os.Exit(1)
		}
	} else {
		err = processFile(*inputPath)
		if err != nil {
			fmt.Printf("Error processing file '%s': %v\n", *inputPath, err)
			os.Exit(1)
		}
	}

	if *dryRun {
		fmt.Println("\nDry run complete. No changes were made.")
	} else {
		fmt.Println("\nProcessing complete.")
	}
}

func processFile(filePath string) error {
	if *verbose {
		fmt.Printf("Processing %s...\n", filePath)
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file '%s': %w", filePath, err)
	}

	if !bytes.Contains(content, []byte{'\r', '\n'}) {
		if *verbose {
			fmt.Printf("  '%s' already uses LF line endings or has no CRLF. Skipping.\n", filePath)
		}
		return nil
	}

	modifiedContent := bytes.ReplaceAll(content, []byte{'\r', '\n'}, []byte{'\n'})

	if *dryRun {
		fmt.Printf("  Would convert CRLF to LF in '%s'.\n", filePath)
		return nil
	}

	outputPath := filePath
	if !*overwrite {
		outputPath = filePath + ".lf"
	}

	if *backup && *overwrite {
		backupPath := filePath + ".bak"
		if *verbose {
			fmt.Printf("  Creating backup '%s'...\n", backupPath)
		}
		err = os.Rename(filePath, backupPath)
		if err != nil {
			return fmt.Errorf("failed to create backup for '%s': %w", filePath, err)
		}
	}

	err = ioutil.WriteFile(outputPath, modifiedContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file '%s': %w", outputPath, err)
	}

	if *overwrite {
		fmt.Printf("  Converted CRLF to LF in '%s'. Original file overwritten.\n", filePath)
	} else {
		fmt.Printf("  Converted CRLF to LF. New file created: '%s'.\n", outputPath)
	}

	return nil
}