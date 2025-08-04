package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <directory> <search_string>")
		os.Exit(1)
	}

	searchDir := os.Args[1]
	searchString := os.Args[2]

	err := filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing path %q: %v\n", path, err)
			return nil // Continue walking even if there's an error with one path
		}
		if !info.IsDir() {
			if containsString(path, searchString) {
				fmt.Println(path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking the directory %q: %v\n", searchDir, err)
		os.Exit(1)
	}
}

func containsString(filePath, searchString string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false // Silently skip files that cannot be opened
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), searchString) {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		return false // Silently skip files that cause scanning errors
	}
	return false
}

// Additional implementation at 2025-08-04 08:04:21
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	startPath := flag.String("path", ".", "Starting directory for the search")
	searchString := flag.String("str", "", "String to search for (required)")
	caseSensitive := flag.Bool("case", false, "Perform case-sensitive search")
	includeExtensions := flag.String("exts", "", "Comma-separated list of file extensions to include (e.g., 'go,txt')")
	excludeDirectories := flag.String("exclude-dirs", "", "Comma-separated list of directory names to exclude (e.g., 'node_modules,.git')")
	maxWorkers := flag.Int("workers", 5, "Maximum number of concurrent file processing workers")

	flag.Parse()

	if *searchString == "" {
		fmt.Println("Error: Search string (-str) is required.")
		flag.Usage()
		os.Exit(1)
	}

	var allowedExts map[string]struct{}
	if *includeExtensions != "" {
		allowedExts = make(map[string]struct{})
		for _, ext := range strings.Split(*includeExtensions, ",") {
			allowedExts[strings.TrimSpace(ext)] = struct{}{}
		}
	}

	var excludedDirs map[string]struct{}
	if *excludeDirectories != "" {
		excludedDirs = make(map[string]struct{})
		for _, dir := range strings.Split(*excludeDirectories, ",") {
			excludedDirs[strings.TrimSpace(dir)] = struct{}{}
		}
	}

	var fileProcessingWg sync.WaitGroup
	filePaths := make(chan string)
	results := make(chan string) // Channel for results to print them sequentially

	// Start file processing workers
	for i := 0; i < *maxWorkers; i++ {
		fileProcessingWg.Add(1)
		go func() {
			defer fileProcessingWg.Done()
			for filePath := range filePaths {
				processFile(filePath, *searchString, *caseSensitive, results)
			}
		}()
	}

	// Goroutine to close results channel when all file processing workers are done
	go func() {
		fileProcessingWg.Wait() // Wait for all file processing goroutines to finish
		close(results)          // Then close the results channel
	}()

	// Start a single goroutine to print results
	var printerWg sync.WaitGroup
	printerWg.Add(1)
	go func() {
		defer printerWg.Done()
		for res := range results {
			fmt.Print(res)
		}
	}()

	// Walk the directory tree and send file paths to the filePaths channel
	err := filepath.WalkDir(*startPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing path %s: %v\n", path, err)
			return nil // Continue walking
		}

		if d.IsDir() {
			// Check if directory should be excluded
			dirName := filepath.Base(path)
			if _, ok := excludedDirs[dirName]; ok {
				return filepath.SkipDir // Skip this directory and its contents
			}
			return nil
		}

		// It's a file
		if allowedExts != nil {
			ext := strings.TrimPrefix(filepath.Ext(path), ".")
			if _, ok := allowedExts[ext]; !ok {
				return nil // Skip file if extension not allowed
			}
		}

		filePaths <- path // Send file path to workers
		return nil
	})

	close(filePaths) // No more files to send to workers

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory: %v\n", err)
	}

	printerWg.Wait() // Wait for the results printer to finish before main exits
}

func processFile(filePath, searchString string, caseSensitive bool, results chan<- string) {
	file, err := os.Open(filePath)
	if err != nil {
		results <- fmt.Sprintf("Error opening file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	foundInFile := false
	var fileResults strings.Builder

	// Prepare search strings for case-insensitive comparison
	searchStrForCompare := searchString
	if !caseSensitive {
		searchStrForCompare = strings.ToLower(searchString)
	}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		lineToCheck := line
		if !caseSensitive {
			lineToCheck = strings.ToLower(line)
		}

		if strings.Contains(lineToCheck, searchStrForCompare) {
			if !foundInFile {
				fileResults.WriteString(fmt.Sprintf("\nFile: %s\n", filePath))
				foundInFile = true
			}
			fileResults.WriteString(fmt.Sprintf("  Line %d: %s\n", lineNum, line))
		}
	}

	if err := scanner.Err(); err != nil {
		results <- fmt.Sprintf("Error reading file %s: %v\n", filePath, err)
	}

	if foundInFile {
		results <- fileResults.String()
	}
}

// Additional implementation at 2025-08-04 08:05:03
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// searchFile reads a file line by line and checks for the search term.
// It sends the file path to the results channel if a match is found.
func searchFile(filePath string, searchTerm string, compiledRegex *regexp.Regexp, ignoreCase bool, isRegex bool, results chan<- string) {
	file, err := os.Open(filePath)
	if err != nil {
		// Silently skip unreadable files to avoid excessive error output for common permissions issues.
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		found := false

		if isRegex {
			if compiledRegex.MatchString(line) {
				found = true
			}
		} else {
			if ignoreCase {
				if strings.Contains(strings.ToLower(line), strings.ToLower(searchTerm)) {
					found = true
				}
			} else {
				if strings.Contains(line, searchTerm) {
					found = true
				}
			}
		}

		if found {
			results <- filePath
			return // Found in this file, no need to scan further
		}
	}

	if err := scanner.Err(); err != nil {
		// Error reading file, but continue processing other files.
		// fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filePath, err)
	}
}

func main() {
	var rootPath string
	var searchTerm string
	var isRegex bool
	var ignoreCase bool
	var excludeDirsStr string
	var workers int

	flag.StringVar(&rootPath, "path", ".", "Root directory to start searching")
	flag.StringVar(&searchTerm, "str", "", "The string or regex to search for")
	flag.BoolVar(&isRegex, "regex", false, "Treat search string as a regular expression")
	flag.BoolVar(&ignoreCase, "ignore-case", false, "Perform case-insensitive search")
	flag.StringVar(&excludeDirsStr, "exclude-dir", ".git,node_modules,vendor", "Comma-separated list of directory names to exclude")
	flag.IntVar(&workers, "workers", runtime.NumCPU(), "Number of concurrent workers")

	flag.Parse()

	if searchTerm == "" {
		fmt.Println("Error: Search term cannot be empty. Use -str flag.")
		flag.Usage()
		os.Exit(1)
	}

	excludeDirs := strings.Split(excludeDirsStr, ",")
	for i, dir := range excludeDirs {
		excludeDirs[i] = strings.TrimSpace(dir)
	}

	var compiledRegex *regexp.Regexp
	if isRegex {
		var err error
		if ignoreCase {
			// (?i) flag makes the regex case-insensitive
			compiledRegex, err = regexp.Compile("(?i)" + searchTerm)
		} else {
			compiledRegex, err = regexp.Compile(searchTerm)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error compiling regex: %v\n", err)
			os.Exit(1)
		}
	}

	filesToProcess := make(chan string, 100) // Buffer for files to be processed
	results := make(chan string)             // Channel for found files
	var wg sync.WaitGroup                    // WaitGroup for worker goroutines

	// Start result collector goroutine
	go func() {
		foundFiles := make(map[string]struct{}) // Use a map to store unique results
		for filePath := range results {
			if _, exists := foundFiles[filePath]; !exists {
				fmt.Println(filePath)
				foundFiles[filePath] = struct{}{}
			}
		}
	}()

	// Start worker goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range filesToProcess {
				searchFile(filePath, searchTerm, compiledRegex, ignoreCase, isRegex, results)
			}
		}()
	}

	// Start directory walking goroutine
	go func() {
		err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				// Error accessing path, but continue walking.
				// fmt.Fprintf(os.Stderr, "Error accessing path %s: %v\n", path, err)
				return nil
			}

			if d.IsDir() {
				dirName := filepath.Base(path)
				for _, excluded := range excludeDirs {
					if dirName == excluded {
						return filepath.SkipDir // Skip this directory
					}
				}
			} else {
				if d.Type().IsRegular() {
					filesToProcess <- path
				}
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking directory %s: %v\n", rootPath, err)
		}
		close(filesToProcess) // No more files will be sent to workers
	}()

	// Wait for all worker goroutines to finish processing files.
	wg.Wait()
	close(results) // All workers are done, no more results will be sent to the collector.
}

// Additional implementation at 2025-08-04 08:06:30
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	searchDir       string
	searchString    string
	caseInsensitive bool
	includeExts     string
	printLines      bool
	numWorkers      int
)

var (
	wg       sync.WaitGroup
	fileChan chan string
)

func init() {
	flag.StringVar(&searchDir, "dir", ".", "Directory to start searching from.")
	flag.StringVar(&searchString, "str", "", "String to search for.")
	flag.BoolVar(&caseInsensitive, "i", false, "Perform case-insensitive search.")
	flag.StringVar(&includeExts, "ext", "", "Comma-separated list of file extensions to include (e.g., 'go,txt'). If empty, all files are checked.")
	flag.BoolVar(&printLines, "lines", false, "Print line numbers where the string is found.")
	flag.IntVar(&numWorkers, "workers", 4, "Number of concurrent workers to process files.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  A script to find files containing a specific string with additional functionality.\n\n")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if searchString == "" {
		fmt.Println("Error: Search string cannot be empty. Use -str flag.")
		flag.Usage()
		os.Exit(1)
	}

	if numWorkers < 1 {
		fmt.Println("Error: Number of workers must be at least 1.")
		os.Exit(1)
	}

	if caseInsensitive {
		searchString = strings.ToLower(searchString)
	}

	var extensions map[string]struct{}
	if includeExts != "" {
		extensions = make(map[string]struct{})
		for _, ext := range strings.Split(includeExts, ",") {
			extensions[strings.ToLower(strings.TrimSpace(ext))] = struct{}{}
		}
	}

	fileChan = make(chan string, numWorkers*2)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(extensions)
	}

	err := filepath.WalkDir(searchDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return nil
		}
		if !d.IsDir() {
			fileChan <- path
		}
		return nil
	})

	close(fileChan)
	wg.Wait()

	if err != nil {
		fmt.Printf("Error walking directory %q: %v\n", searchDir, err)
		os.Exit(1)
	}
}

func worker(extensions map[string]struct{}) {
	defer wg.Done()
	for filePath := range fileChan {
		processFile(filePath, extensions)
	}
}

func processFile(filePath string, extensions map[string]struct{}) {
	if len(extensions) > 0 {
		ext := strings.ToLower(filepath.Ext(filePath))
		if ext != "" && ext[0] == '.' {
			ext = ext[1:]
		}
		if _, ok := extensions[ext]; !ok {
			return
		}
	}

	searchFile(filePath)
}

func searchFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %q: %v\n", filePath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	foundInFile := false
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		lineForSearch := line
		if caseInsensitive {
			lineForSearch = strings.ToLower(line)
		}

		if strings.Contains(lineForSearch, searchString) {
			if !foundInFile {
				fmt.Printf("Found in: %s\n", filePath)
				foundInFile = true
			}
			if printLines {
				fmt.Printf("  Line %d: %s\n", lineNumber, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %q: %v\n", filePath, err)
	}
}