package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 8 {
		fmt.Println("Usage: go run main.go <filename> <year> <month> <day> <hour> <minute> <second>")
		fmt.Println("Example: go run main.go my_file.txt 2023 10 27 15 30 00")
		os.Exit(1)
	}

	filename := os.Args[1]
	year, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Invalid year: %v", err)
	}
	monthInt, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf("Invalid month: %v", err)
	}
	day, err := strconv.Atoi(os.Args[4])
	if err != nil {
		log.Fatalf("Invalid day: %v", err)
	}
	hour, err := strconv.Atoi(os.Args[5])
	if err != nil {
		log.Fatalf("Invalid hour: %v", err)
	}
	minute, err := strconv.Atoi(os.Args[6])
	if err != nil {
		log.Fatalf("Invalid minute: %v", err)
	}
	second, err := strconv.Atoi(os.Args[7])
	if err != nil {
		log.Fatalf("Invalid second: %v", err)
	}

	month := time.Month(monthInt)

	newTime := time.Date(year, month, day, hour, minute, second, 0, time.Local)

	// os.Chtimes changes the access and modification times of the named file.
	// The "creation time" (birth time) is not directly modifiable via os.Chtimes
	// and typically requires platform-specific system calls, which are not
	// part of the standard os package for portable timestamp manipulation.
	// This program updates the modification and access times, which are
	// commonly what is referred to when discussing file timestamp updates.
	err = os.Chtimes(filename, newTime, newTime)
	if err != nil {
		log.Fatalf("Failed to update timestamps for %s: %v", filename, err)
	}

	fmt.Printf("Successfully updated modification and access timestamps for %s to %s\n", filename, newTime.Format("2006-01-02 15:04:05"))
}

// Additional implementation at 2025-06-20 23:48:18
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// parseTime parses a string into a time.Time.
// It supports "now" for the current time, and various common layouts.
// It also supports relative times like "now-1h", "now+2d", etc.
func parseTime(s string) (time.Time, error) {
	if s == "now" {
		return time.Now(), nil
	}

	// Check for relative time (e.g., "now-1h", "now+24h")
	if len(s) > 3 && s[:3] == "now" {
		durationStr := s[3:]
		d, err := time.ParseDuration(durationStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid duration format: %w", err)
		}
		return time.Now().Add(d), nil
	}

	// Try common layouts
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00", // With timezone
		"2006-01-02T15:04:05",       // Without timezone
		"2006-01-02 15:04:05",
		"2006-01-02",
		"15:04:05", // Assume today's date
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			// If only time is provided, combine with today's date
			if layout == "15:04:05" {
				now := time.Now()
				return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), now.Location()), nil
			}
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported time format: %s", s)
}

// updateFileTimestamps updates the access and modification times of a file.
// If a time is zero, it means that timestamp should not be changed.
// Note: Go's standard library (os.Chtimes) does not provide a portable way to change
// a file's creation (birth) time. This function updates modification and access times.
func updateFileTimestamps(path string, newModTime, newAccessTime time.Time, dryRun bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %w", path, err)
	}

	// os.Stat().ModTime() returns the modification time.
	// Access time is not directly available from os.FileInfo on all systems.
	// When os.Chtimes is called, it requires both access and modification times.
	// If a new time is not provided, we use the file's current modification time
	// as a fallback for the access time, or the actual modification time for mod time.
	currentModTime := info.ModTime()
	currentAccessTime := info.ModTime() // Fallback for access time if not explicitly set

	var setModTime time.Time
	if !newModTime.IsZero() {
		setModTime = newModTime
	} else {
		setModTime = currentModTime
	}

	var setAccessTime time.Time
	if !newAccessTime.IsZero() {
		setAccessTime = newAccessTime
	} else {
		setAccessTime = currentAccessTime
	}

	if dryRun {
		fmt.Printf("Dry run: Would update %s: mod_time=%s, access_time=%s\n", path, setModTime.Format(time.RFC3339), setAccessTime.Format(time.RFC3339))
		return nil
	}

	fmt.Printf("Updating %s: mod_time=%s, access_time=%s\n", path, setModTime.Format(time.RFC3339), setAccessTime.Format(time.RFC3339))
	err = os.Chtimes(path, setAccessTime, setModTime)
	if err != nil {
		return fmt.Errorf("failed to change times for %s: %w", path, err)
	}
	return nil
}

func main() {
	var (
		modTimeStr    string
		accessTimeStr string
		recursive     bool
		dryRun        bool
		help          bool
	)

	flag.StringVar(&modTimeStr, "m", "", "Set modification time (e.g., 'now', 'now-1h', '2023-01-01 10:00:00'). If empty, modification time is not changed.")
	flag.StringVar(&accessTimeStr, "a", "", "Set access time (e.g., 'now', 'now-1h', '2023-01-01 10:00:00'). If empty, access time is not changed.")
	flag.BoolVar(&recursive, "r", false, "Recursively process directories.")
	flag.BoolVar(&dryRun, "dry-run", false, "Perform a dry run without making any changes.")
	flag.BoolVar(&help, "h", false, "Show help message.")

	flag.Parse()

	if help {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file1> [file2 ...]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No files or directories specified.")
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file1> [file2 ...]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	var newModTime time.Time
	var newAccessTime time.Time
	var err error

	if modTimeStr != "" {
		newModTime, err = parseTime(modTimeStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing modification time: %v\n", err)
			os.Exit(1)
		}
	}

	if accessTimeStr != "" {
		newAccessTime, err = parseTime(accessTimeStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing access time: %v\n", err)
			os.Exit(1)
		}
	}

	if newModTime.IsZero() && newAccessTime.IsZero() {
		fmt.Fprintln(os.Stderr, "Error: At least one of -m or -a must be specified.")
		os.Exit(1)
	}

	for _, targetPath := range flag.Args() {
		fileInfo, err := os.Stat(targetPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", targetPath, err)
			continue
		}

		if fileInfo.IsDir() {
			if !recursive {
				fmt.Fprintf(os.Stderr, "Skipping directory %s (use -r for recursive processing)\n", targetPath)
				continue
			}

			err = filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error walking %s: %v\n", path, err)
					return nil // Continue walking even if one file errors
				}
				// Skip directories themselves, only process regular files and symlinks
				if d.IsDir() {
					return nil
				}
				if d.Type().IsRegular() || d.Type()&fs.ModeSymlink != 0 {
					if err := updateFileTimestamps(path, newModTime, newAccessTime, dryRun); err != nil {
						fmt.Fprintf(os.Stderr, "Error updating %s: %v\n", path, err)
					}
				}
				return nil
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error during recursive walk for %s: %v\n", targetPath, err)
			}
		} else {
			if err := updateFileTimestamps(targetPath, newModTime, newAccessTime, dryRun); err != nil {
				fmt.Fprintf(os.Stderr, "Error updating %s: %v\n", targetPath, err)
			}
		}
	}
}

// Additional implementation at 2025-06-20 23:49:14
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func parseTime(timeStr string) (time.Time, error) {
	if timeStr == "now" {
		return time.Now(), nil
	}
	layouts := []string{
		time.RFC3339,                 // "2006-01-02T15:04:05Z07:00"
		"2006-01-02 15:04:05",        // "YYYY-MM-DD HH:MM:SS"
		"2006-01-02T15:04:05",        // "YYYY-MM-DDTHH:MM:SS"
		"2006-01-02",                 // "YYYY-MM-DD"
		"2006/01/02 15:04:05",        // "YYYY/MM/DD HH:MM:SS"
		"2006/01/02",                 // "YYYY/MM/DD"
	}
	for _, layout := range layouts {
		t, err := time.Parse(layout, timeStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported time format: %s. Use 'now', RFC3339, or common date/datetime formats (e.g., '2023-01-01T12:00:00Z', '2023-01-01 12:00:00', '2023-01-01')", timeStr)
}

func updateFileTimestamp(
	filePath string,
	targetMTime *time.Time,
	targetATime *time.Time,
	dryRun bool,
	onlyOlder bool,
	onlyNewer bool,
) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %w", filePath, err)
	}

	currentMTime := fileInfo.ModTime()
	// os.FileInfo does not directly expose access time.
	// For Chtimes, if targetATime is nil, we'll pass currentMTime as a fallback
	// or rely on the OS to update atime based on access.
	// To truly preserve atime, platform-specific syscalls would be needed to read it.
	// For simplicity and portability, we use currentMTime as a proxy if atime is not specified.
	currentATime := currentMTime

	var newATime time.Time
	if targetATime != nil {
		newATime = *targetATime
	} else {
		newATime = currentATime // Use currentMTime as a proxy for currentATime if not specified
	}

	var newMTime time.Time
	if targetMTime != nil {
		newMTime = *targetMTime
	} else {
		newMTime = currentMTime // Keep current modification time if not specified
	}

	// Apply conditional updates based on modification time
	if onlyOlder && !newMTime.Before(currentMTime) {
		fmt.Printf("Skipping %s: new modification time is not older than current, but --only-older is set.\n", filePath)
		return nil
	}
	if onlyNewer && !newMTime.After(currentMTime) {
		fmt.Printf("Skipping %s: new modification time is not newer than current, but --only-newer is set.\n", filePath)
		return nil
	}

	if dryRun {
		fmt.Printf("Dry run: Would update %s (Mod: %s -> %s, Acc: %s -> %s)\n",
			filePath, currentMTime.Format(time.RFC3339), newMTime.Format(time.RFC3339),
			currentATime.Format(time.RFC3339), newATime.Format(time.RFC3339))
		return nil
	}

	err = os.Chtimes(filePath, newATime, newMTime)
	if err != nil {
		return fmt.Errorf("failed to change times for %s: %w", filePath, err)
	}

	fmt.Printf("Updated %s (Mod: %s -> %s, Acc: %s -> %s)\n",
		filePath, currentMTime.Format(time.RFC3339), newMTime.Format(time.RFC3339),
		currentATime.Format(time.RFC3339), newATime.Format(time.RFC3339))
	return nil
}

func main() {
	var (
		path        string
		recursive   bool
		mtimeStr    string
		atimeStr    string
		extsStr     string
		dryRun      bool
		onlyOlder   bool
		onlyNewer   bool
	)

	flag.StringVar(&path, "path", ".", "File or directory path to update timestamps. Defaults to current directory.")
	flag.BoolVar(&recursive, "recursive", false, "Recursively process files in subdirectories.")
	flag.StringVar(&mtimeStr, "mtime", "", "Set modification time. Use 'now' or a specific time (e.g., '2023-01-01T12:00:00Z', '2023-01-01', '2023-01-01 12:00:00').")
	flag.StringVar(&atimeStr, "atime", "", "Set access time. Same formats as -mtime.")
	flag.StringVar(&extsStr, "ext", "", "Comma-separated list of file extensions to process (e.g., '.txt,.log'). If empty, all files are processed.")
	flag.BoolVar(&dryRun, "dryrun", false, "Perform a dry run without actually changing timestamps.")
	flag.BoolVar(&onlyOlder, "only-older", false, "Only update if the current modification time is older than the target modification time.")
	flag.BoolVar(&onlyNewer, "only-newer", false, "Only update if the current modification time is newer than the target modification time.")
	flag.Parse()

	if mtimeStr == "" && atimeStr == "" {
		fmt.Println("Error: At least one of -mtime or -atime must be specified.")
		flag.Usage()
		os.Exit(1)
	}

	if onlyOlder && onlyNewer {
		fmt.Println("Error: Cannot use both -only-older and -only-newer simultaneously.")
		flag.Usage()
		os.Exit(1)
	}

	var targetMTime *time.Time
	if mtimeStr != "" {
		t, err := parseTime(mtimeStr)
		if err != nil {
			fmt.Printf("Error parsing -mtime: %v\n", err)
			os.Exit(1)
		}
		targetMTime = &t
	}

	var targetATime *time.Time
	if atimeStr != "" {
		t, err := parseTime(atimeStr)
		if err != nil {
			fmt.Printf("Error parsing -atime: %v\n", err)
			os.Exit(1)
		}
		targetATime = &t
	}

	var allowedExts map[string]struct{}
	if extsStr != "" {
		allowedExts = make(map[string]struct{})
		for _, ext := range strings.Split(extsStr, ",") {
			e := strings.TrimSpace(ext)
			if e != "" {
				if !strings.HasPrefix(e, ".") {
					e = "." + e
				}
				allowedExts[strings.ToLower(e)] = struct{}{}
			}
		}
	}

	processFile := func(filePath string, info fs.FileInfo) error {
		if info.IsDir() {
			return nil // Skip directories
		}

		if allowedExts != nil {
			fileExt := strings.ToLower(filepath.Ext(filePath))
			if _, ok := allowedExts[fileExt]; !ok {
				return nil // Skip if extension not allowed
			}
		}

		err := updateFileTimestamp(filePath, targetMTime, targetATime, dryRun, onlyOlder, onlyNewer)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", filePath, err)
		}
		return nil
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error accessing path %s: %v\n", path, err)
		os.Exit(1)
	}

	if fileInfo.IsDir() && recursive {
		err = filepath.Walk(path, func(p string, info fs