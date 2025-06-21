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
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func findDuplicates(rootPath string) (map[string][]string, error) {
	fileHashes := make(map[string][]string)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing path %q: %v\n", path, err)
			return nil
		}

		if !info.IsDir() && info.Mode().IsRegular() {
			hash, err := calculateFileHash(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error calculating hash for %q: %v\n", path, err)
				return nil
			}
			fileHashes[hash] = append(fileHashes[hash], path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the path %q: %w", rootPath, err)
	}

	return fileHashes, nil
}

func main() {
	rootPath := "."

	if len(os.Args) > 1 {
		rootPath = os.Args[1]
	}

	fmt.Printf("Scanning for duplicate files in: %s\n", rootPath)

	duplicateGroups, err := findDuplicates(rootPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	foundDuplicates := false
	for hash, paths := range duplicateGroups {
		if len(paths) > 1 {
			foundDuplicates = true
			fmt.Printf("\nDuplicate files (Hash: %s):\n", hash)
			for _, p := range paths {
				fmt.Printf("  - %s\n", p)
			}
		}
	}

	if !foundDuplicates {
		fmt.Println("\nNo duplicate files found.")
	}
}

// Additional implementation at 2025-06-21 00:47:36
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	bufferSize = 65536 // 64KB buffer for hashing
	numWorkers = 4     // Number of concurrent hash calculation workers
)

func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.CopyBuffer(hash, file, make([]byte, bufferSize)); err != nil {
		return "", fmt.Errorf("failed to hash file %s: %w", filePath, err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func findDuplicates(rootPath string) (map[string][]string, error) {
	fileHashes := make(map[string][]string)
	filesToProcess := make(chan string)
	results := make(chan struct {
		hash string
		path string
		err  error
	})
	var wg sync.WaitGroup
	var collectorWg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range filesToProcess {
				hash, err := calculateFileHash(filePath)
				results <- struct {
					hash string
					path string
					err  error
				}{hash, filePath, err}
			}
		}()
	}

	go func() {
		err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error accessing path %q: %v\n", path, err)
				return nil
			}
			if !info.IsDir() {
				filesToProcess <- path
			}
			return nil
		})
		close(filesToProcess)
		if err != nil {
			log.Printf("Error walking the path %q: %v\n", rootPath, err)
		}
	}()

	collectorWg.Add(1)
	go func() {
		defer collectorWg.Done()
		for res := range results {
			if res.err != nil {
				log.Printf("Error processing file %s: %v\n", res.path, res.err)
				continue
			}
			fileHashes[res.hash] = append(fileHashes[res.hash], res.path)
		}
	}()

	wg.Wait()
	close(results)
	collectorWg.Wait()

	duplicates := make(map[string][]string)
	for hash, paths := range fileHashes {
		if len(paths) > 1 {
			duplicates[hash] = paths
		}
	}

	return duplicates, nil
}

func main() {
	rootPath := flag.String("path", ".", "Root directory to scan for duplicate files")
	flag.Parse()

	if *rootPath == "" {
		log.Fatal("Please provide a root directory using the -path flag.")
	}

	absPath, err := filepath.Abs(*rootPath)
	if err != nil {
		log.Fatalf("Failed to resolve absolute path for %q: %v", *rootPath, err)
	}

	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		log.Fatalf("Path does not exist: %q", absPath)
	}
	if err != nil {
		log.Fatalf("Failed to stat path %q: %v", absPath, err)
	}
	if !info.IsDir() {
		log.Fatalf("Path is not a directory: %q", absPath)
	}

	fmt.Printf("Scanning for duplicate files in: %s\n", absPath)
	duplicateFiles, err := findDuplicates(absPath)
	if err != nil {
		log.Fatalf("Error finding duplicates: %v", err)
	}

	if len(duplicateFiles) == 0 {
		fmt.Println("No duplicate files found.")
		return
	}

	fmt.Println("\n--- Duplicate Files Found ---")
	for hash, paths := range duplicateFiles {
		fmt.Printf("Hash: %s\n", hash)
		for _, path := range paths {
			fmt.Printf("  - %s\n", path)
		}
		fmt.Println()
	}
}

// Additional implementation at 2025-06-21 00:48:04
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// fileInfo stores basic information about a file
type fileInfo struct {
	Path string
	Size int64
}

// hashResult stores the hash and path of a file
type hashResult struct {
	Path string
	Hash string
	Err  error
}

// calculateMD5Hash computes the MD5 hash of a file
func calculateMD5Hash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to read file %s for hashing: %w", filePath, err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// findDuplicateFiles walks the given root directory and identifies duplicate files based on content hash.
// It uses file size as a preliminary filter and processes hashing concurrently.
func findDuplicateFiles(root string) (map[string][]string, error) {
	filesBySize := make(map[int64][]fileInfo)
	filesByHash := make(map[string][]string)
	var mu sync.Mutex // Mutex for protecting filesByHash map writes

	log.Printf("Scanning directory: %s", root)
	startScan := time.Now()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v", path, err)
			return nil // Continue walking even if one path has an error
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				log.Printf("Error getting file info for %s: %v", path, err)
				return nil
			}
			mu.Lock()
			filesBySize[info.Size()] = append(filesBySize[info.Size()], fileInfo{Path: path, Size: info.Size()})
			mu.Unlock()
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", root, err)
	}

	log.Printf("Scan completed in %v. Found %d unique file sizes.", time.Since(startScan), len(filesBySize))

	var wg sync.WaitGroup
	hashChan := make(chan fileInfo)
	resultsChan := make(chan hashResult)

	// Start worker goroutines for hashing
	numWorkers := 4 // Adjust based on CPU cores/IO capacity
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fi := range hashChan {
				hash, err := calculateMD5Hash(fi.Path)
				resultsChan <- hashResult{Path: fi.Path, Hash: hash, Err: err}
			}
		}()
	}

	// Send files to hash workers
	go func() {
		for _, files := range filesBySize {
			if len(files) > 1 { // Only hash files that have potential duplicates (same size)
				for _, fi := range files {
					hashChan <- fi
				}
			}
		}
		close(hashChan)
	}()

	// Collect results from hash workers
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	log.Println("Starting file hashing...")
	startHash := time.Now()
	hashedFilesCount := 0

	for res := range resultsChan {
		if res.Err != nil {
			log.Printf("Error hashing file %s: %v", res.Path, res.Err)
			continue
		}
		mu.Lock()
		filesByHash[res.Hash] = append(filesByHash[res.Hash], res.Path)
		mu.Unlock()
		hashedFilesCount++
	}

	log.Printf("Hashing completed for %d files in %v.", hashedFilesCount, time.Since(startHash))

	// Filter out unique files (those with only one path per hash)
	duplicates := make(map[string][]string)
	for hash, paths := range filesByHash {
		if len(paths) > 1 {
			duplicates[hash] = paths
		}
	}

	return duplicates, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory_path>")
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

	fmt.Printf("Finding duplicate files in: %s\n", rootDir)
	startTime := time.Now()

	duplicates, err := findDuplicateFiles(rootDir)
	if err != nil {
		log.Fatalf("Failed to find duplicates: %v", err)
	}

	totalDuplicatesFound := 0
	if len(duplicates) == 0 {
		fmt.Println("\nNo duplicate files found.")
	} else {
		fmt.Println("\n--- Duplicate Files Found ---")
		for hash, paths := range duplicates {
			fmt.Printf("Hash: %s\n", hash)
			for _, path := range paths {
				fmt.Printf("  - %s\n", path)
				totalDuplicatesFound++
			}
			fmt.Println()
		}
	}

	fmt.Printf("Scan finished in %v\n", time.Since(startTime))
	fmt.Printf("Total duplicate groups: %d\n", len(duplicates))
	fmt.Printf("Total duplicate files (including originals in groups): %d\n", totalDuplicatesFound)
}

// Additional implementation at 2025-06-21 00:49:01
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	bufferSize = 65536 // 64KB buffer for reading files
	numWorkers = 4     // Number of concurrent hash calculation workers
)

type fileInfo struct {
	Path string
	Hash string
	Err  error
}

func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	hasher := sha256.New()
	buf := make([]byte, bufferSize)

	if _, err := io.CopyBuffer(hasher, file, buf); err != nil {
		return "", fmt.Errorf("failed to hash file %s: %w", filePath, err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func findDuplicates(rootPath string) (map[string][]string, error) {
	duplicateFiles := make(map[string][]string)
	filePaths := make(chan string)
	results := make(chan fileInfo)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range filePaths {
				hash, err := calculateFileHash(filePath)
				results <- fileInfo{Path: filePath, Hash: hash, Err: err}
			}
		}()
	}

	go func() {
		err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error accessing path %s: %v", path, err)
				return nil
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			if info.Size() == 0 {
				return nil
			}

			filePaths <- path
			return nil
		})
		close(filePaths)
		if err != nil {
			log.Printf("Error walking directory %s: %v", rootPath, err)
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if res.Err != nil {
			log.Printf("Error processing file %s: %v", res.Path, res.Err)
			continue
		}
		duplicateFiles[res.Hash] = append(duplicateFiles[res.Hash], res.Path)
	}

	return duplicateFiles, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory_to_scan>")
		os.Exit(1)
	}

	rootDirectory := os.Args[1]

	fmt.Printf("Scanning directory: %s\n", rootDirectory)

	duplicates, err := findDuplicates(rootDirectory)
	if err != nil {
		log.Fatalf("Error finding duplicates: %v", err)
	}

	foundDuplicates := false
	for hash, paths := range duplicates {
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
	}
}