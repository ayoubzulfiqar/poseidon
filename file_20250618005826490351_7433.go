package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: <program> <filepath> <new_timestamp_RFC3339>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	timestampStr := os.Args[2]

	newTime, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		log.Fatalf("Error parsing timestamp: %v", err)
	}

	err = os.Chtimes(filePath, newTime, newTime)
	if err != nil {
		log.Fatalf("Failed to update timestamps for %s: %v", filePath, err)
	}

	fmt.Printf("Successfully updated timestamps for %s\n", filePath)
}

// Additional implementation at 2025-06-18 00:59:39
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func parseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func updateFileTimes(path string, newAtime, newMtime time.Time, dryRun, verbose bool) error {
	if newAtime.IsZero() {
		newAtime = time.Now()
	}
	if newMtime.IsZero() {
		newMtime = time.Now()
	}

	if verbose {
		fmt.Printf("Processing %s: atime=%s, mtime=%s\n", path, newAtime.Format(time.RFC3339), newMtime.Format(time.RFC3339))
	}

	if dryRun {
		return nil
	}

	err := os.Chtimes(path, newAtime, newMtime)
	if err != nil {
		return fmt.Errorf("failed to update times for %s: %w", path, err)
	}
	return nil
}

func main() {
	var (
		modTimeStr string
		accTimeStr string
		sourceFile string
		recursive  bool
		dryRun     bool
		verbose    bool
		help       bool
	)

	flag.StringVar(&modTimeStr, "m", "", "Set modification time (RFC3339 format, e.g., 2006-01-02T15:04:05Z)")
	flag.StringVar(&accTimeStr, "a", "", "Set access time (RFC3339 format, e.g., 2006-01-02T15:04:05Z)")
	flag.StringVar(&sourceFile, "s", "", "Copy modification and access times from this source file")
	flag.BoolVar(&recursive, "r", false, "Process directories recursively")
	flag.BoolVar(&dryRun, "d", false, "Dry run: show what would be done without making changes")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&help, "h", false, "Show help message")

	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	targetPaths := flag.Args()
	if len(targetPaths) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No target paths specified.")
		flag.Usage()
		os.Exit(1)
	}

	var newMtime time.Time
	var newAtime time.Time
	var err error

	if sourceFile != "" {
		if modTimeStr != "" || accTimeStr != "" {
			fmt.Fprintln(os.Stderr, "Error: Cannot use -s with -m or -a.")
			flag.Usage()
			os.Exit(1)
		}
		srcInfo, err := os.Stat(sourceFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not stat source file %s: %v\n", sourceFile, err)
			os.Exit(1)
		}
		newMtime = srcInfo.ModTime()
		// os.Stat does not provide access time portably.
		// For simplicity and portability, when copying from source,
		// we'll set access time to source's modification time.
		newAtime = srcInfo.ModTime()
		if verbose {
			fmt.Printf("Using times from source file %s: atime=%s, mtime=%s\n", sourceFile, newAtime.Format(time.RFC3339), newMtime.Format(time.RFC3339))
		}
	} else {
		if modTimeStr != "" {
			newMtime, err = parseTime(modTimeStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Invalid modification time format: %v\n", err)
				os.Exit(1)
			}
		}
		if accTimeStr != "" {
			newAtime, err = parseTime(accTimeStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Invalid access time format: %v\n", err)
				os.Exit(1)
			}
		}
	}

	for _, targetPath := range targetPaths {
		info, err := os.Stat(targetPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not stat target path %s: %v\n", targetPath, err)
			continue
		}

		if info.IsDir() {
			if recursive {
				err := filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if d.IsDir() {
						return nil // Skip directories themselves for time update when walking
					}
					return updateFileTimes(path, newAtime, newMtime, dryRun, verbose)
				})
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error walking directory %s: %v\n", targetPath, err)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Skipping directory %s. Use -r for recursive processing.\n", targetPath)
			}
		} else {
			err := updateFileTimes(targetPath, newAtime, newMtime, dryRun, verbose)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

// Additional implementation at 2025-06-18 01:00:45
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
	"syscall"

	"golang.org/x/sys/windows"
)

var (
	targetPath      = flag.String("path", "", "Target file or directory path to update timestamps.")
	setToNow        = flag.Bool("set-to-now", false, "Set timestamps to the current time.")
	sourceFilePath  = flag.String("set-to-file", "", "Set timestamps to match another file's timestamps.")
	recursive       = flag.Bool("recursive", false, "Process directories recursively.")
	setCreationTime = flag.Bool("set-creation-time", false, "Attempt to set creation time (Windows only).")
)

func main() {
	flag.Parse()

	if *targetPath == "" {
		fmt.Println("Error: -path is required.")
		flag.Usage()
		os.Exit(1)
	}

	var atime, mtime time.Time
	var creationTime time.Time

	if *setToNow {
		now := time.Now()
		atime = now
		mtime = now
		creationTime = now
	} else if *sourceFilePath != "" {
		srcInfo, err := os.Stat(*sourceFilePath)
		if err != nil {
			fmt.Printf("Error: Could not get info for source file %s: %v\n", *sourceFilePath, err)
			os.Exit(1)
		}
		mtime = srcInfo.ModTime()
		atime = mtime // Fallback for atime, as os.Stat doesn't universally provide access time

		if *setCreationTime {
			if sysStat, ok := srcInfo.Sys().(*syscall.Win32FileAttributeData); ok {
				creationTime = time.Unix(0, sysStat.CreationTime.Nanoseconds())
			} else {
				fmt.Println("Warning: Could not get creation time from source file on this OS. Using modification time as fallback.")
				creationTime = mtime
			}
		}

	} else {
		fmt.Println("Error: One of -set-to-now or -set-to-file is required.")
		flag.Usage()
		os.Exit(1)
	}

	updateFileTimestamps := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return err
		}

		if info.IsDir() && !*recursive {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}

		fmt.Printf("Updating timestamps for: %s\n", path)

		err = os.Chtimes(path, atime, mtime)
		if err != nil {
			fmt.Printf("Error updating access/modification times for %s: %v\n", path, err)
			return err
		}

		if *setCreationTime {
			if err := setWindowsCreationTime(path, creationTime, atime, mtime); err != nil {
				fmt.Printf("Error updating creation time for %s: %v\n", path, err)
			}
		}
		return nil
	}

	targetInfo, err := os.Stat(*targetPath)
	if err != nil {
		fmt.Printf("Error: Could not get info for target path %s: %v\n", *targetPath, err)
		os.Exit(1)
	}

	if targetInfo.IsDir() {
		if *recursive {
			fmt.Printf("Recursively updating timestamps in directory: %s\n", *targetPath)
			err = filepath.Walk(*targetPath, updateFileTimestamps)
			if err != nil {
				fmt.Printf("Error during recursive walk: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Target is a directory. Updating its own timestamps: %s\n", *targetPath)
			err = os.Chtimes(*targetPath, atime, mtime)
			if err != nil {
				fmt.Printf("Error updating access/modification times for directory %s: %v\n", *targetPath, err)
				os.Exit(1)
			}
			if *setCreationTime {
				if err := setWindowsCreationTime(*targetPath, creationTime, atime, mtime); err != nil {
					fmt.Printf("Error updating creation time for directory %s: %v\n", *targetPath, err)
				}
			}
		}
	} else {
		err = updateFileTimestamps(*targetPath, targetInfo, nil)
		if err != nil {
			os.Exit(1)
		}
	}

	fmt.Println("Timestamp update complete.")
}

func setWindowsCreationTime(path string, creationTime, accessTime, modificationTime time.Time) error {
	handle, err := syscall.CreateFile(
		syscall.StringToUTF16Ptr(path),
		syscall.FILE_WRITE_ATTRIBUTES,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
	if err != nil {
		return fmt.Errorf("failed to open file handle for %s: %w", path, err)
	}
	defer syscall.CloseHandle(handle)

	cTime := windows.NsecToFiletime(creationTime.UnixNano())
	aTime := windows.NsecToFiletime(accessTime.UnixNano())
	mTime := windows.NsecToFiletime(modificationTime.UnixNano())

	err = syscall.SetFileTime(handle, &cTime, &aTime, &mTime)
	if err != nil {
		return fmt.Errorf("failed to set file times for %s: %w", path, err)
	}
	return nil
}