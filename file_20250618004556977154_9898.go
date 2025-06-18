package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to hash file %s: %w", filePath, err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func deduplicateDirectory(dirPath string) ([]string, error) {
	fileHashes := make(map[string]string)
	var deletedFiles []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing path %q: %v\n", path, err)
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		hash, err := calculateFileHash(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not hash %s: %v\n", path, err)
			return nil
		}

		if originalPath, found := fileHashes[hash]; found {
			fmt.Printf("Duplicate found: %s (original: %s)\n", path, originalPath)
			if err := os.Remove(path); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting duplicate %s: %v\n", path, err)
			} else {
				deletedFiles = append(deletedFiles, path)
				fmt.Printf("Deleted: %s\n", path)
			}
		} else {
			fileHashes[hash] = path
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %q: %w", dirPath, err)
	}

	return deletedFiles, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory_path>")
		os.Exit(1)
	}

	dirToDeduplicate := os.Args[1]

	info, err := os.Stat(dirToDeduplicate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accessing path %s: %v\n", dirToDeduplicate, err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: %s is not a directory.\n", dirToDeduplicate)
		os.Exit(1)
	}

	fmt.Printf("Starting deduplication in directory: %s\n", dirToDeduplicate)

	deleted, err := deduplicateDirectory(dirToDeduplicate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Deduplication failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nDeduplication complete.\n")
	if len(deleted) > 0 {
		fmt.Printf("Deleted %d duplicate files:\n", len(deleted))
		for _, file := range deleted {
			fmt.Printf("- %s\n", file)
		}
	} else {
		fmt.Println("No duplicate files found or deleted.")
	}
}

// Additional implementation at 2025-06-18 00:46:34
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// calculateFileHash calculates the SHA256 hash of a file.
func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to hash file %s: %w", filePath, err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// findDuplicates walks the given directory and identifies duplicate files based on their SHA256 hash.
// It returns a map where keys are hashes and values are lists of file paths sharing that hash.
func findDuplicates(rootDir string) (map[string][]string, error) {
	fileHashes := make(map[string][]string)

	var wg sync.WaitGroup
	pathsChan := make(chan string, 1000) // Buffered channel for file paths
	resultsChan := make(chan struct {
		hash string
		path string
		err  error
	}, 1000) // Buffered channel for hash results

	numWorkers := runtime.NumCPU()
	if numWorkers == 0 {
		numWorkers = 1
	}

	// Start worker goroutines to calculate hashes
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathsChan {
				hash, err := calculateFileHash(path)
				resultsChan <- struct {
					hash string
					path string
					err  error
				}{hash, path, err}
			}
		}()
	}

	// Goroutine to walk the directory and send paths to pathsChan
	walkErrChan := make(chan error, 1) // Channel to signal walk errors
	go func() {
		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// Log the error but continue walking for individual file/dir access issues
				fmt.Fprintf(os.Stderr, "Warning: Error accessing path %s: %v\n", path, err)
				return nil
			}
			if !info.IsDir() {
				pathsChan <- path
			}
			return nil
		})
		close(pathsChan) // Close pathsChan when walking is done
		walkErrChan <- err // Send the final walk error (nil or actual error)
		close(walkErrChan)
	}()

	// Goroutine to close resultsChan when all workers are done
	go func() {
		wg.Wait() // Wait for all workers to finish processing pathsChan
		close(resultsChan)
	}()

	// Process results from resultsChan sequentially
	for res := range resultsChan {
		if res.err != nil {
			fmt.Fprintf(os.Stderr, "Error processing file %s: %v\n", res.path, res.err)
			continue
		}
		fileHashes[res.hash] = append(fileHashes[res.hash], res.path)
	}

	// Check for initial walk error after all results are processed
	if err := <-walkErrChan; err != nil {
		return nil, fmt.Errorf("directory walk failed: %w", err)
	}

	return fileHashes, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: deduplicator <directory_to_scan>")
		os.Exit(1)
	}

	rootDir := os.Args[1]

	info, err := os.Stat(rootDir)
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Directory '%s' does not exist.\n", rootDir)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accessing directory '%s': %v\n", rootDir, err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: '%s' is not a directory.\n", rootDir)
		os.Exit(1)
	}

	fmt.Printf("Scanning directory: %s\n", rootDir)

	duplicateGroups, err := findDuplicates(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding duplicates: %v\n", err)
		os.Exit(1)
	}

	foundDuplicates := false
	for hash, paths := range duplicateGroups {
		if len(paths) > 1 {
			foundDuplicates = true
			fmt.Printf("\nDuplicate files (Hash: %s):\n", hash)
			for _, path := range paths {
				fmt.Printf("  - %s\n", path)
			}
		}
	}

	if !foundDuplicates {
		fmt.Println("\nNo duplicate files found.")
	} else {
		fmt.Println("\nScan complete. Above are the duplicate groups found.")
	}
}

// Additional implementation at 2025-06-18 00:47:31
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// calculateFileHash computes the SHA256 hash of a file.
func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to hash file %s: %w", filePath, err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// scanDirectory walks the given root directory, calculates hashes for all regular files,
// and returns a map where keys are hashes and values are lists of file paths with that hash.
func scanDirectory(root string) (map[string][]string, error) {
	fileHashes := make(map[string][]string)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Log the error but continue walking
			log.Printf("Error accessing path %s: %v", path, err)
			return nil // Don't stop the walk for individual errors
		}

		if d.Type().IsRegular() { // Only process regular files
			hash, hashErr := calculateFileHash(path)
			if hashErr != nil {
				log.Printf("Error hashing file %s: %v", path, hashErr)
				return nil // Don't stop the walk for hashing errors
			}

			// Add to map
			fileHashes[hash] = append(fileHashes[hash], path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", root, err)
	}

	return fileHashes, nil
}

// reportDuplicates prints out the identified duplicate files.
func reportDuplicates(fileHashes map[string][]string) {
	foundDuplicates := false
	// Collect hashes that have duplicates to sort them for consistent output
	var duplicateHashes []string
	for hash, paths := range fileHashes {
		if len(paths) > 1 {
			foundDuplicates = true
			duplicateHashes = append(duplicateHashes, hash)
		}
	}

	sort.Strings(duplicateHashes) // Sort hashes for consistent output

	for _, hash := range duplicateHashes {
		paths := fileHashes[hash]
		fmt.Printf("Duplicate files (Hash: %s):\n", hash)
		// Sort paths for consistent output
		sort.Strings(paths)
		for _, p := range paths {
			fmt.Printf("  - %s\n", p)
		}
		fmt.Println()
	}

	if !foundDuplicates {
		fmt.Println("No duplicate files found.")
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run dedupe.go <directory_to_scan>")
		os.Exit(1)
	}

	rootDir := os.Args[1]

	// Check if the provided path is a directory
	info, err := os.Stat(rootDir)
	if err != nil {
		log.Fatalf("Error accessing path %s: %v", rootDir, err)
	}
	if !info.IsDir() {
		log.Fatalf("Path %s is not a directory.", rootDir)
	}

	fmt.Printf("Scanning directory: %s\n", rootDir)
	fileHashes, err := scanDirectory(rootDir)
	if err != nil {
		log.Fatalf("Failed to scan directory: %v", err)
	}

	fmt.Println("\n--- Duplicate Files Report ---")
	reportDuplicates(fileHashes)
	fmt.Println("--- Report End ---")
}