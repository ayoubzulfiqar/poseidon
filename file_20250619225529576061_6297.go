package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func printDirTree(path string, prefix string, isLast bool) {
	info, err := os.Lstat(path)
	if err != nil {
		log.Printf("Error accessing %s: %v", path, err)
		return
	}

	connector := "├── "
	if isLast {
		connector = "└── "
	}

	fmt.Printf("%s%s%s\n", prefix, connector, info.Name())

	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			log.Printf("Error reading directory %s: %v", path, err)
			return
		}

		sort.Slice(entries, func(i, j int) bool {
			aIsDir := entries[i].IsDir()
			bIsDir := entries[j].IsDir()

			if aIsDir && !bIsDir {
				return true
			}
			if !aIsDir && bIsDir {
				return false
			}
			return entries[i].Name() < entries[j].Name()
		})

		newPrefix := prefix
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}

		for i, entry := range entries {
			childPath := filepath.Join(path, entry.Name())
			childIsLast := (i == len(entries)-1)
			printDirTree(childPath, newPrefix, childIsLast)
		}
	}
}

func main() {
	startPath := "."

	if len(os.Args) > 1 {
		startPath = os.Args[1]
	}

	absPath, err := filepath.Abs(startPath)
	if err != nil {
		log.Fatalf("Error getting absolute path for %s: %v", startPath, err)
	}

	rootInfo, err := os.Lstat(absPath)
	if err != nil {
		log.Fatalf("Error accessing root path %s: %v", absPath, err)
	}
	fmt.Printf("%s\n", rootInfo.Name())

	if rootInfo.IsDir() {
		entries, err := os.ReadDir(absPath)
		if err != nil {
			log.Fatalf("Error reading root directory %s: %v", absPath, err)
		}

		sort.Slice(entries, func(i, j int) bool {
			aIsDir := entries[i].IsDir()
			bIsDir := entries[j].IsDir()

			if aIsDir && !bIsDir {
				return true
			}
			if !aIsDir && bIsDir {
				return false
			}
			return entries[i].Name() < entries[j].Name()
		})

		for i, entry := range entries {
			childPath := filepath.Join(absPath, entry.Name())
			childIsLast := (i == len(entries)-1)
			printDirTree(childPath, "", childIsLast)
		}
	}
}

// Additional implementation at 2025-06-19 22:56:49
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorBlue   = "\033[34m" // Directories
	colorGreen  = "\033[32m" // Files
	colorYellow = "\033[33m" // Sizes
	colorRed    = "\033[31m" // Errors
)

var (
	rootPath         string
	maxDepth         int
	excludePatterns  string
	showSizes        bool
	compiledExcludes []*regexp.Regexp
)

func init() {
	flag.StringVar(&rootPath, "path", ".", "The starting directory path.")
	flag.IntVar(&maxDepth, "depth", 0, "Maximum depth to traverse (0 for unlimited).")
	flag.StringVar(&excludePatterns, "exclude", ".git,node_modules,vendor,.DS_Store,Thumbs.db,*.log,*.tmp", "Comma-separated regex patterns to exclude files/directories.")
	flag.BoolVar(&showSizes, "sizes", false, "Show file sizes.")
	flag.Parse()

	if excludePatterns != "" {
		patterns := strings.Split(excludePatterns, ",")
		for _, p := range patterns {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			// Convert simple glob-like patterns to regex
			if strings.Contains(p, "*") || strings.Contains(p, "?") {
				p = regexp.QuoteMeta(p) // Escape special regex characters
				p = strings.ReplaceAll(p, "\\*", ".*")
				p = strings.ReplaceAll(p, "\\?", ".")
				p = "^" + p + "$" // Match whole name
			} else {
				// For exact matches, ensure it matches the whole name
				p = "^" + regexp.QuoteMeta(p) + "$"
			}

			re, err := regexp.Compile(p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Invalid regex pattern '%s': %v\n", p, err)
				continue
			}
			compiledExcludes = append(compiledExcludes, re)
		}
	}
}

func main() {
	fmt.Println(colorBlue + rootPath + colorReset)

	if err := printTree(rootPath, "", 1); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating tree: %v\n", err)
		os.Exit(1)
	}
}

func printTree(path string, prefix string, currentDepth int) error {
	if maxDepth > 0 && currentDepth > maxDepth {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("%s%s [Error: %v]%s\n", prefix, colorRed+"<unreadable directory>"+colorReset, err, colorReset)
		return nil
	}

	var filteredEntries []fs.DirEntry
	for _, entry := range entries {
		name := entry.Name()
		shouldExclude := false
		for _, re := range compiledExcludes {
			if re.MatchString(name) {
				shouldExclude = true
				break
			}
		}
		if !shouldExclude {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	sort.Slice(filteredEntries, func(i, j int) bool {
		iIsDir := filteredEntries[i].IsDir()
		jIsDir := filteredEntries[j].IsDir()

		if iIsDir && !jIsDir {
			return true
		}
		if !iIsDir && jIsDir {
			return false
		}
		return filteredEntries[i].Name() < filteredEntries[j].Name()
	})

	for i, entry := range filteredEntries {
		isLast := (i == len(filteredEntries)-1)
		entryPrefix := "├── "
		nextPrefix := "│   "
		if isLast {
			entryPrefix = "└── "
			nextPrefix = "    "
		}

		fullPath := filepath.Join(path, entry.Name())
		info, err := entry.Info()
		if err != nil {
			fmt.Printf("%s%s%s [Error: %v]%s\n", prefix, entryPrefix, colorRed+entry.Name()+colorReset, err, colorReset)
			continue
		}

		if info.IsDir() {
			fmt.Printf("%s%s%s%s%s\n", prefix, entryPrefix, colorBlue, entry.Name(), colorReset)
			if err := printTree(fullPath, prefix+nextPrefix, currentDepth+1); err != nil {
				return err
			}
		} else {
			sizeStr := ""
			if showSizes {
				sizeStr = formatSize(info.Size())
			}
			fmt.Printf("%s%s%s%s%s%s\n", prefix, entryPrefix, colorGreen, entry.Name(), colorReset, sizeStr)
		}
	}
	return nil
}

func formatSize(size int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case size >= GB:
		return fmt.Sprintf(" (%s%.1fGB%s)", colorYellow, float64(size)/GB, colorReset)
	case size >= MB:
		return fmt.Sprintf(" (%s%.1fMB%s)", colorYellow, float64(size)/MB, colorReset)
	case size >= KB:
		return fmt.Sprintf(" (%s%.1fKB%s)", colorYellow, float64(size)/KB, colorReset)
	default:
		return fmt.Sprintf(" (%s%dB%s)", colorYellow, size, colorReset)
	}
}