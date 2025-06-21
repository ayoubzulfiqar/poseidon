package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func printTree(path string, indent string, isLast bool) {
	fmt.Print(indent)

	if isLast {
		fmt.Print("└── ")
	} else {
		fmt.Print("├── ")
	}

	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(info.Name())

	if !info.IsDir() {
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("%s    Error reading directory: %v\n", indent, err)
		return
	}

	var filteredEntries []fs.DirEntry
	for _, entry := range entries {
		filteredEntries = append(filteredEntries, entry)
	}

	newIndent := indent
	if isLast {
		newIndent += "    "
	} else {
		newIndent += "│   "
	}

	for i, entry := range filteredEntries {
		entryPath := filepath.Join(path, entry.Name())
		printTree(entryPath, newIndent, i == len(filteredEntries)-1)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: go run main.go <directory_path>")
		os.Exit(1)
	}

	rootPath := args[0]

	info, err := os.Stat(rootPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Error: Directory '%s' does not exist.", rootPath)
		}
		log.Fatalf("Error accessing '%s': %v", rootPath, err)
	}
	if !info.IsDir() {
		log.Fatalf("Error: '%s' is not a directory.", rootPath)
	}

	fmt.Println(rootPath)

	entries, err := os.ReadDir(rootPath)
	if err != nil {
		log.Fatalf("Error reading root directory '%s': %v", rootPath, err)
	}

	var filteredEntries []fs.DirEntry
	for _, entry := range entries {
		filteredEntries = append(filteredEntries, entry)
	}

	for i, entry := range filteredEntries {
		entryPath := filepath.Join(rootPath, entry.Name())
		printTree(entryPath, "", i == len(filteredEntries)-1)
	}
}

// Additional implementation at 2025-06-20 22:47:47
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	ColorReset  = "\033[0m"
	ColorDir    = "\033[1;34m" // Blue bold
	ColorFile   = "\033[0m"    // Default
	ColorSymlink = "\033[0;36m" // Cyan
	ColorInfo   = "\033[0;37m" // Grey
)

type Config struct {
	StartPath    string
	MaxDepth     int // -1 for unlimited, 0 for current dir only, 1 for current + children, etc.
	IncludePatterns []string
	ExcludePatterns []string
	ShowSize     bool
	ShowTime     bool
	UseColors    bool
}

func formatSize(s int64) string {
	const (
		_ = iota
		KB = 1 << (10 * iota)
		MB
		GB
		TB
	)
	switch {
	case s >= TB:
		return fmt.Sprintf("%.1fT", float64(s)/float64(TB))
	case s >= GB:
		return fmt.Sprintf("%.1fG", float64(s)/float64(GB))
	case s >= MB:
		return fmt.Sprintf("%.1fM", float64(s)/float64(MB))
	case s >= KB:
		return fmt.Sprintf("%.1fK", float64(s)/float64(KB))
	default:
		return fmt.Sprintf("%dB", s)

// Additional implementation at 2025-06-20 22:48:58
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// TreeOptions holds configuration for generating the directory tree.
type TreeOptions struct {
	MaxDepth    int    // Maximum depth to traverse (0 for no limit)
	ShowSize    bool   // Whether to display file sizes
	IncludeGlob string // Glob pattern for files/directories to include
	ExcludeGlob string // Glob pattern for files/directories to exclude
}

// formatSize converts bytes into a human-readable string (e.g., 1.2M, 500B).
func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2fG", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2fM", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2fK", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

// printTree recursively prints the directory tree structure.
func printTree(path string, prefix string, depth int, options TreeOptions) error {
	if options.MaxDepth != 0 && depth > options.MaxDepth {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("cannot read directory %s: %w", path, err)
	}

	// Sort entries: directories first, then alphabetically by name.
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir() // Directories come before files
		}
		return entries[i].Name() < entries[j].Name()
	})

	for i, entry := range entries {
		name := entry.Name()
		fullPath := filepath.Join(path, name)

		// Apply include filter
		if options.IncludeGlob != "" {
			matched, err := filepath.Match(options.IncludeGlob, name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Invalid include glob pattern '%s': %v\n", options.IncludeGlob, err)
				continue // Skip on pattern error
			}
			if !matched {
				continue
			}
		}

		// Apply exclude filter
		if options.ExcludeGlob != "" {
			matched, err := filepath.Match(options.ExcludeGlob, name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Invalid exclude glob pattern '%s': %v\n", options.ExcludeGlob, err)
				continue // Skip on pattern error
			}
			if matched {
				continue
			}
		}

		isLast := (i == len(entries)-1)
		var currentPrefix string
		if isLast {
			currentPrefix = prefix + "└── "
		} else {
			currentPrefix = prefix + "├── "
		}

		info, err := entry.Info()
		if err != nil {
			// If info cannot be retrieved, print name with an error message and continue.
			fmt.Printf("%s%s [error getting info: %v]\n", currentPrefix, name, err)
			if entry.IsDir() {
				var nextPrefix string
				if isLast {
					nextPrefix = prefix + "    "
				} else {
					nextPrefix = prefix + "│   "
				}
				// Attempt to recurse even if info failed for the current entry.
				_ = printTree(fullPath, nextPrefix, depth+1, options)
			}
			continue
		}

		var sizeStr string
		if options.ShowSize && !info.IsDir() {
			sizeStr = fmt.Sprintf(" (%s)", formatSize(info.Size()))
		}

		fmt.Printf("%s%s%s\n", currentPrefix, name, sizeStr)

		if info.IsDir() {
			var nextPrefix string
			if isLast {
				nextPrefix = prefix + "    "
			} else {
				nextPrefix = prefix + "│   "
			}
			err := printTree(fullPath, nextPrefix, depth+1, options)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error traversing %s: %v\n", fullPath, err)
			}
		}
	}
	return nil
}

func main() {
	args := os.Args[1:]
	path := "." // Default to current directory

	options := TreeOptions{
		MaxDepth:    0, // 0 means no limit
		ShowSize:    false,
		IncludeGlob: "",
		ExcludeGlob: "",
	}

	// Parse command-line arguments
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-depth":
			if i+1 < len(args) {
				depth, err := strconv.Atoi(args[i+1])
				if err != nil || depth < 0 {
					fmt.Fprintf(os.Stderr, "Error: Invalid depth '%s'. Must be a non-negative integer.\n", args[i+1])
					os.Exit(1)
				}
				options.MaxDepth = depth
				i++ // Consume the next argument
			} else {
				fmt.Fprintf(os.Stderr, "Error: -depth requires a value.\n")
				os.Exit(1)
			}
		case "-size":
			options.ShowSize = true
		case "-include":
			if i+1 < len(args) {
				options.IncludeGlob = args[i+1]
				i++
			} else {
				fmt.Fprintf(os.Stderr, "Error: -include requires a pattern.\n")
				os.Exit(1)
			}
		case "-exclude":
			if i+1 < len(args) {
				options.ExcludeGlob = args[i+1]
				i++
			} else {
				fmt.Fprintf(os.Stderr, "Error: -exclude requires a pattern.\n")
				os.Exit(1)
			}
		default:
			if !strings.HasPrefix(arg, "-") { // Treat as path if not an option
				path = arg
			} else {
				fmt.Fprintf(os.Stderr, "Error: Unknown argument '%s'.\n", arg)
				os.Exit(1)
			}
		}
	}

	// Get and print the root directory name
	rootInfo, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", rootInfo.Name())

	// Start the recursive tree printing
	err = printTree(path, "", 0, options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating tree: %v\n", err)
		os.Exit(1)
	}
}