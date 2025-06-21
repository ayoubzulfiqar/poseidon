package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		return
	}

	downloadsPath := filepath.Join(homeDir, "Downloads")
	fmt.Printf("Starting cleanup of downloads folder: %s\n", downloadsPath)

	// Define the age threshold for files to be deleted (e.g., 30 days)
	// Files older than this duration will be considered for deletion.
	ageThreshold := 30 * 24 * time.Hour // 30 days

	err = cleanupDownloads(downloadsPath, ageThreshold)
	if err != nil {
		fmt.Printf("Cleanup completed with errors: %v\n", err)
	} else {
		fmt.Println("Downloads cleanup completed successfully.")
	}
}

func cleanupDownloads(folderPath string, threshold time.Duration) error {
	now := time.Now()
	cutoffTime := now.Add(-threshold)

	return filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing %s: %v\n", path, err)
			return nil // Continue walking even if there's an error with one item
		}

		if info.IsDir() {
			// Skip directories, we only want to delete files
			return nil
		}

		// Check if the file's modification time is older than the cutoff
		if info.ModTime().Before(cutoffTime) {
			fmt.Printf("Deleting old file: %s (Modified: %s)\n", path, info.ModTime().Format("2006-01-02 15:04:05"))
			if err := os.Remove(path); err != nil {
				fmt.Printf("Error deleting %s: %v\n", path, err)
				// Do not return error here, try to continue with other files
			}
		}
		return nil
	})
}

// Additional implementation at 2025-06-21 04:28:08
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

var categories = map[string][]string{
	"Documents":   {".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".rtf", ".odt", ".ods", ".odp"},
	"Images":      {".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp", ".heic"},
	"Videos":      {".mp4", ".mov", ".avi", ".mkv", ".flv", ".wmv", ".webm"},
	"Audio":       {".mp3", ".wav", ".aac", ".flac", ".ogg", ".wma"},
	"Archives":    {".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz"},
	"Executables": {".exe", ".dmg", ".app", ".msi", ".deb", ".rpm", ".pkg"},
	"Code":        {".go", ".py", ".js", ".html", ".css", ".java", ".c", ".cpp", ".h", ".sh", ".json", ".xml", ".yml", ".yaml"},
	"Others":      {},
}

func getCategory(ext string) string {
	ext = strings.ToLower(ext)
	for category, extensions := range categories {
		for _, e := range extensions {
			if e == ext {
				return category
			}
		}
	}
	return "Others"
}

func main() {
	downloadsPath := flag.String("path", "", "Path to the downloads folder (e.g., /Users/youruser/Downloads)")
	dryRun := flag.Bool("dry-run", false, "If true, simulate cleanup without making changes")
	ageDays := flag.Int("age", 30, "Process files older than this many days")
	deleteOldUncategorized := flag.Bool("delete-uncategorized", false, "If true, delete files older than 'age' that were not moved to a category (use with caution)")
	flag.Parse()

	if *downloadsPath == "" {
		fmt.Println("Error: Downloads path not specified. Use --path flag.")
		flag.Usage()
		return
	}

	fmt.Printf("Starting downloads cleanup for: %s\n", *downloadsPath)
	fmt.Printf("Dry Run: %t\n", *dryRun)
	fmt.Printf("Processing files older than: %d days\n", *ageDays)
	if *deleteOldUncategorized {
		fmt.Println("Warning: Old uncategorized files will be DELETED from the base path.")
	} else {
		fmt.Println("Old uncategorized files will be moved to 'Others' or ignored if already there.")
	}
	fmt.Println("--------------------------------------------------")

	cleanupDownloads(*downloadsPath, *dryRun, *ageDays, *deleteOldUncategorized)

	fmt.Println("--------------------------------------------------")
	fmt.Println("Cleanup process finished.")
}

func cleanupDownloads(basePath string, dryRun bool, ageDays int, deleteOldUncategorized bool) {
	minModTime := time.Now().AddDate(0, 0, -ageDays)

	for category := range categories {
		categoryPath := filepath.Join(basePath, category)
		if _, err := os.Stat(categoryPath); os.IsNotExist(err) {
			if dryRun {
				fmt.Printf("[DRY RUN] Would create directory: %s\n", categoryPath)
			} else {
				err := os.Mkdir(categoryPath, 0755)
				if err != nil {
					fmt.Printf("Error creating directory %s: %v\n", categoryPath, err)
					return
				}
				fmt.Printf("Created directory: %s\n", categoryPath)
			}
		}
	}

	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil
		}

		if path == basePath {
			return nil
		}
		for category := range categories {
			if path == filepath.Join(basePath, category) {
				return filepath.SkipDir
			}
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			fmt.Printf("Error getting file info for %s: %v\n", path, err)
			return nil
		}

		if info.ModTime().After(minModTime) {
			return nil
		}

		ext := filepath.Ext(d.Name())
		category := getCategory(ext)
		targetDir := filepath.Join(basePath, category)
		newPath := filepath.Join(targetDir, d.Name())

		if filepath.Dir(path) == targetDir {
			return nil
		}

		if dryRun {
			fmt.Printf("[DRY RUN] Would move: %s -> %s\n", path, newPath)
		} else {
			err := os.Rename(path, newPath)
			if err != nil {
				if os.IsExist(err) {
					timestamp := time.Now().Format("_20060102_150405")
					newPathWithTimestamp := filepath.Join(targetDir, strings.TrimSuffix(d.Name(), ext)+timestamp+ext)
					fmt.Printf("Warning: File %s already exists in target. Trying to move to %s\n", newPath, newPathWithTimestamp)
					err = os.Rename(path, newPathWithTimestamp)
					if err != nil {
						fmt.Printf("Error moving %s to %s: %v\n", path, newPathWithTimestamp, err)
					} else {
						fmt.Printf("Moved: %s -> %s\n", path, newPathWithTimestamp)
					}
				} else {
					fmt.Printf("Error moving %s to %s: %v\n", path, newPath, err)
				}
			} else {
				fmt.Printf("Moved: %s -> %s\n", path, newPath)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory for moving: %v\n", err)
	}

	if deleteOldUncategorized {
		fmt.Println("\nChecking for old uncategorized files to delete...")
		err = filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("Error accessing path %s: %v\n", path, err)
				return nil
			}

			if path == basePath {
				return nil
			}
			for category := range categories {
				if path == filepath.Join(basePath, category) {
					return filepath.SkipDir
				}
			}

			if d.IsDir() {
				return nil
			}

			info, err := d.Info()
			if err != nil {
				fmt.Printf("Error getting file info for %s: %v\n", path, err)
				return nil
			}

			if info.ModTime().After(minModTime) {
				return nil
			}

			ext := filepath.Ext(d.Name())
			category := getCategory(ext)

			if category == "Others" && filepath.Dir(path) == basePath {
				if dryRun {
					fmt.Printf("[DRY RUN] Would delete old uncategorized file: %s\n", path)
				} else {
					err := os.Remove(path)
					if err != nil {
						fmt.Printf("Error deleting %s: %v\n", path, err)
					} else {
						fmt.Printf("Deleted old uncategorized file: %s\n", path)
					}
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error walking for deletion: %v\n", err)
		}
	}
}

// Additional implementation at 2025-06-21 04:28:46
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Config struct {
	DownloadsPath string
	MaxAgeDays    int
	DryRun        bool
	Mappings      map[string]string
}

type Stats struct {
	DeletedFiles int
	MovedFiles   int
	DeletedDirs  int
}

func main() {
	var config Config

	flag.StringVar(&config.DownloadsPath, "path", "", "Path to the downloads folder (default: user's downloads)")
	flag.IntVar(&config.MaxAgeDays, "age", 30, "Delete files older than this many days")
	flag.BoolVar(&config.DryRun, "dry-run", false, "Simulate cleanup without making any changes")
	flag.Parse()

	if config.DownloadsPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting user home directory: %v\n", err)
			os.Exit(1)
		}
		config.DownloadsPath = filepath.Join(homeDir, "Downloads")
	}

	config.Mappings = map[string]string{
		".zip":  "Archives",
		".rar":  "Archives",
		".7z":   "Archives",
		".tar":  "Archives",
		".gz":   "Archives",
		".bz2":  "Archives",
		".xz":   "Archives",
		".pdf":  "Documents",
		".doc":  "Documents",
		".docx": "Documents",
		".xls":  "Documents",
		".xlsx": "Documents",
		".ppt":  "Documents",
		".pptx": "Documents",
		".txt":  "Documents",
		".rtf":  "Documents",
		".exe":  "Executables",
		".msi":  "Executables",
		".dmg":  "Executables",
		".app":  "Executables",
		".iso":  "DiskImages",
		".img":  "DiskImages",
		".mp3":  "Audio",
		".wav":  "Audio",
		".flac": "Audio",
		".ogg":  "Audio",
		".mp4":  "Video",
		".mov":  "Video",
		".avi":  "Video",
		".mkv":  "Video",
		".jpg":  "Images",
		".jpeg": "Images",
		".png":  "Images",
		".gif":  "Images",
		".bmp":  "Images",
		".tiff": "Images",
		".webp": "Images",
	}

	fmt.Printf("Starting downloads cleanup for: %s\n", config.DownloadsPath)
	fmt.Printf("Files older than %d days will be deleted.\n", config.MaxAgeDays)
	if config.DryRun {
		fmt.Println("Running in DRY-RUN mode. No changes will be made.")
	}

	stats, err := cleanupDownloads(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cleanup failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n--- Cleanup Summary ---")
	fmt.Printf("Files deleted: %d\n", stats.DeletedFiles)
	fmt.Printf("Files moved: %d\n", stats.MovedFiles)
	fmt.Printf("Empty directories deleted: %d\n", stats.DeletedDirs)
	fmt.Println("-----------------------")
}

func cleanupDownloads(cfg Config) (Stats, error) {
	var stats Stats
	maxAgeDuration := time.Duration(cfg.MaxAgeDays) * 24 * time.Hour
	var directories []string

	err := filepath.Walk(cfg.DownloadsPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", path, err)
			return nil
		}

		if info.IsDir() {
			if path != cfg.DownloadsPath {
				directories = append(directories, path)
			}
			return nil
		}

		age := time.Since(info.ModTime())
		if age > maxAgeDuration {
			fmt.Printf("  [DELETE] %s (age: %s)\n", path, age.Round(time.Hour))
			if !cfg.DryRun {
				if err := os.Remove(path); err != nil {
					fmt.Fprintf(os.Stderr, "    Error deleting %s: %v\n", path, err)
				} else {
					stats.DeletedFiles++
				}
			} else {
				stats.DeletedFiles++
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		if targetSubDir, ok := cfg.Mappings[ext]; ok {
			destDir := filepath.Join(cfg.DownloadsPath, targetSubDir)
			destPath := filepath.Join(destDir, info.Name())

			if _, err := os.Stat(destDir); os.IsNotExist(err) {
				fmt.Printf("  [MKDIR] Creating directory: %s\n", destDir)
				if !cfg.DryRun {
					if err := os.MkdirAll(destDir, 0755); err != nil {
						fmt.Fprintf(os.Stderr, "    Error creating directory %s: %v\n", destDir, err)
						return nil
					}
				}
			}

			if _, err := os.Stat(destPath); err == nil {
				fmt.Printf("  [SKIP] Destination file already exists: %s. Skipping move for %s\n", destPath, path)
				return nil
			} else if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "    Error checking destination %s: %v\n", destPath, err)
				return nil
			}

			fmt.Printf("  [MOVE] %s to %s\n", path, destPath)
			if !cfg.DryRun {
				if err := os.Rename(path, destPath); err != nil {
					fmt.Fprintf(os.Stderr, "    Error moving %s: %v\n", path, err)
				} else {
					stats.MovedFiles++
				}
			} else {
				stats.MovedFiles++
			}
		}
		return nil
	})

	if err != nil {
		return stats, fmt.Errorf("error walking downloads directory: %w", err)
	}

	sort.Slice(directories, func(i, j int) bool {
		return len(directories[i]) > len(directories[j])
	})

	for _, dir := range directories {
		entries, err := os.ReadDir(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Error reading directory %s: %v\n", dir, err)
			continue
		}
		if len(entries) == 0 {
			fmt.Printf("  [DELETE EMPTY DIR] %s\n", dir)
			if !cfg.DryRun {
				if err := os.Remove(dir); err != nil {
					fmt.Fprintf(os.Stderr, "    Error deleting empty directory %s: %v\n", dir, err)
				} else {
					stats.DeletedDirs++
				}
			} else {
				stats.DeletedDirs++
			}
		}
	}

	return stats, nil
}

// Additional implementation at 2025-06-21 04:29:40
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Configuration variables
var (
	downloadsPath string
	ageDays       int
	dryRun        bool
	logFile       string
	logger        *log.Logger
)

func init() {
	defaultDownloadsPath := ""
	homeDir, err := os.UserHomeDir()
	if err == nil {
		defaultDownloadsPath = filepath.Join(homeDir, "Downloads")
	} else {
		// Fallback if home directory can't be determined
		defaultDownloadsPath = "./downloads" // A sensible default for testing or specific use cases
	}

	flag.StringVar(&downloadsPath, "path", defaultDownloadsPath, "Path to the downloads folder")
	flag.IntVar(&ageDays, "age", 30, "Files older than this many days will be considered for deletion")
	flag.BoolVar(&dryRun, "dry-run", false, "If true, no files will be deleted, only reported")
	flag.StringVar(&logFile, "log-file", "", "Path to a log file (optional). If empty, logs to console.")
	flag.Parse()

	// Initialize logger
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file %s: %v", logFile, err)
		}
		logger = log.New(f, "", log.LstdFlags)
	} else {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}
}

func main() {
	logger.Printf("Starting downloads cleanup for: %s", downloadsPath)
	logger.Printf("Files older than %d days will be considered.", ageDays)
	if dryRun {
		logger.Println("DRY RUN mode: No files will be deleted.")
	}

	thresholdTime := time.Now().AddDate(0, 0, -ageDays)

	err := cleanupDownloads(downloadsPath, thresholdTime)
	if err != nil {
		logger.Fatalf("Cleanup failed: %v", err)
	}

	logger.Println("Downloads cleanup completed.")
}

// cleanupDownloads recursively traverses the directory and cleans up files.
func cleanupDownloads(dirPath string, thresholdTime time.Time) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		// If the directory doesn't exist or is not accessible, log and return error
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var subDirs []string

	for _, file := range files {
		filePath := filepath.Join(dirPath, file.Name())

		if file.IsDir() {
			subDirs = append(subDirs, filePath)
			continue // Process subdirectories after files in current dir
		}

		// Skip symbolic links to avoid infinite loops or deleting files outside the target path
		if file.Mode()&os.ModeSymlink != 0 {
			logger.Printf("Skipping symbolic link: %s", filePath)
			continue
		}

		if !shouldDeleteFile(file, thresholdTime) {
			continue
		}

		logger.Printf("Candidate for deletion: %s (Modified: %s)", filePath, file.ModTime().Format("2006-01-02"))

		if !dryRun {
			err := os.Remove(filePath)
			if err != nil {
				logger.Printf("Error deleting %s: %v", filePath, err)
			} else {
				logger.Printf("Deleted: %s", filePath)
			}
		}
	}

	// Recursively clean subdirectories
	for _, subDir := range subDirs {
		err := cleanupDownloads(subDir, thresholdTime)
		if err != nil {
			logger.Printf("Error cleaning subdirectory %s: %v", subDir, err)
		}
	}

	// After cleaning files and subdirectories, check if the current directory is empty
	// This ensures we delete empty directories bottom-up.
	// Only attempt to remove if it's not the root downloads path itself.
	if dirPath != downloadsPath {
		err = removeEmptyDirectory(dirPath)
		if err != nil {
			logger.Printf("Error removing empty directory %s: %v", dirPath, err)
		}
	}

	return nil
}

// shouldDeleteFile determines if a file should be deleted based on its age and extension.
func shouldDeleteFile(file os.FileInfo, thresholdTime time.Time) bool {
	// Files that are too new should not be deleted
	if file.ModTime().After(thresholdTime) {
		return false
	}

	ext := strings.ToLower(filepath.Ext(file.Name()))

	// Policy 1: List of extensions to always keep (e.g., important documents, media, source code).
	// These files will NOT be deleted, regardless of age.
	keepExtensions := map[string]bool{
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
		".ppt": true, ".pptx": true, ".odt": true, ".ods": true, ".odp": true,
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true,
		".mp3": true, ".wav": true, ".flac": true, ".ogg": true,
		".mp4": true, ".avi": true, ".mkv": true, ".mov": true, ".wmv": true,
		".txt": true, ".md": true, ".log": true, ".csv": true, ".json": true, ".xml": true,
		".go": true, ".py": true, ".java": true, ".c": true, ".cpp": true, ".h": true, ".js": true, ".html": true, ".css": true, // Source code/web files
	}

	if keepExtensions[ext] {
		logger.Printf("Keeping file (by extension policy): %s (Modified: %s)", file.Name(), file.ModTime().Format("2006-01-02"))
		return false
	}

	// Policy 2: List of extensions to specifically target for deletion if old.
	// These are typically installers, archives, temporary files, or partial downloads.
	deleteExtensions := map[string]bool{
		".exe": true, ".msi": true, ".dmg": true, ".deb": true, ".rpm": true, ".pkg": true,
		".zip": true, ".tar": true, ".gz": true, ".tgz": true, ".rar": true, ".7z": true, ".bz2": true, ".xz": true,
		".iso": true, ".img": true, ".vhd": true, ".vmdk": true, // Disk images
		".tmp": true, ".temp": true, ".bak": true, ".old": true,
		".crdownload": true, // Chrome partial downloads
		".part":       true, // Firefox partial downloads
		".torrent":    true, // Torrent files
	}

	if deleteExtensions[ext] {
		return true // Delete these if they are older than threshold
	}

	// Policy 3: Default behavior for files not covered by Policy 1 or 2.
	// If a file's extension is not explicitly in `keepExtensions` and not in `deleteExtensions`,
	// it will be deleted if it's older than the threshold. This makes the script
	// quite aggressive for "other" old files.
	// If you want to be more conservative (i.e., only delete files in deleteExtensions),
	// change this `return true` to `return false`.
	return true
}

// removeEmptyDirectory checks if a directory is empty and removes it.
// It does not remove the root downloadsPath.
func removeEmptyDirectory(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s to check if empty: %w", dirPath, err)
	}

	if len(files) == 0 {
		logger.Printf("Candidate for empty directory deletion: %s", dirPath)
		if !dryRun {
			err := os.Remove(dirPath)
			if err != nil {
				return fmt.Errorf("error deleting empty directory %s: %w", dirPath, err)
			}
			logger.Printf("Deleted empty directory: %s", dirPath)
		}
	}
	return nil
}