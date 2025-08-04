package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extractZip(zipPath string) error {
	destDir := filepath.Join(filepath.Dir(zipPath), strings.TrimSuffix(filepath.Base(zipPath), filepath.Ext(zipPath)))

	fmt.Printf("Extracting %s to %s\n", zipPath, destDir)

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file %s: %w", zipPath, err)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path in zip: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, f.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", fpath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory for %s: %w", fpath, err)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", fpath, err)
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("failed to open file in zip %s: %w", f.Name, err)
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			return fmt.Errorf("failed to copy content for %s: %w", f.Name, err)
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <folder_path>")
		os.Exit(1)
	}

	rootPath := os.Args[1]

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".zip") {
			if err := extractZip(path); err != nil {
				fmt.Fprintf(os.Stderr, "Error extracting %s: %v\n", path, err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking the path %q: %v\n", rootPath, err)
		os.Exit(1)
	}

	fmt.Println("Extraction process completed.")
}

// Additional implementation at 2025-08-04 08:48:40
package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extractAllZips(root string) error {
	fmt.Printf("Starting recursive ZIP extraction from: %s\n", root)

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(strings.ToLower(info.Name()), ".zip") {
			fmt.Printf("Found ZIP file: %s\n", path)

			zipDir := filepath.Dir(path)
			zipFileNameWithoutExt := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			destDir := filepath.Join(zipDir, zipFileNameWithoutExt)

			fmt.Printf("Extracting %s to %s\n", path, destDir)
			if err := unzipFile(path, destDir); err != nil {
				fmt.Printf("Error extracting %s: %v\n", path, err)
			} else {
				fmt.Printf("Successfully extracted %s\n", path)
			}
		}
		return nil
	})
}

func unzipFile(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip file %q: %w", src, err)
	}
	defer r.Close()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory %q: %w", dest, err)
	}

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path in zip: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, f.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %q: %w", fpath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory for %q: %w", fpath, err)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create output file %q: %w", fpath, err)
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("failed to open file in zip %q: %w", f.Name, err)
		}

		_, err = io.Copy(outFile, rc)

		rc.Close()
		outFile.Close()

		if err != nil {
			return fmt.Errorf("failed to copy data for %q: %w", f.Name, err)
		}
	}
	return nil
}

func main() {
	sourceDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Using current directory as source: %s\n", sourceDir)
	fmt.Println("Please ensure this directory or its subdirectories contain .zip files.")
	fmt.Println("Each .zip file will be extracted into a new folder next to it.")
	fmt.Println("For example, 'myarchive.zip' will extract to 'myarchive/'")
	fmt.Println("------------------------------------------------------------")

	if err := extractAllZips(sourceDir); err != nil {
		fmt.Printf("Overall extraction process completed with errors: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("------------------------------------------------------------")
	fmt.Println("Recursive ZIP extraction process finished.")
}

// Additional implementation at 2025-08-04 08:49:55
package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// extractZip extracts a single ZIP file to a specified destination directory.
// It creates the destination directory if it doesn't exist and handles file permissions.
// It also includes a basic check to prevent ZipSlip vulnerabilities.
func extractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file %s: %w", zipPath, err)
	}
	defer r.Close()

	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", destDir, err)
	}

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)

		// Basic ZipSlip prevention: ensure the extracted file path is within the destination directory.
		// Clean the path to resolve any ".." components.
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path in zip (ZipSlip detected): %s", fpath)
		}

		if f.FileInfo().IsDir() {
			// Create directory
			if err := os.MkdirAll(fpath, f.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", fpath, err)
			}
			continue
		}

		// Ensure parent directories exist for the file
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory for %s: %w", fpath, err)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", fpath, err)
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close() // Close the output file even if opening zip entry fails
			return fmt.Errorf("failed to open file in zip %s: %w", f.Name, err)
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()      // Close the reader for the zip entry
		outFile.Close() // Close the output file

		if err != nil {
			return fmt.Errorf("failed to copy data for %s: %w", fpath, err)
		}

		// Set file permissions explicitly
		if err := os.Chmod(fpath, f.Mode()); err != nil {
			return fmt.Errorf("failed to set permissions for %s: %w", fpath, err)
		}
	}
	return nil
}

// extractAllZipsRecursively walks through the given root directory and extracts all found ZIP files.
// Each ZIP file is extracted into a new folder named after the ZIP file (without extension)
// in the same directory as the ZIP file.
func extractAllZipsRecursively(rootPath string) error {
	fmt.Printf("Starting recursive ZIP extraction in: %s\n", rootPath)

	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// If there's an error accessing a path (e.g., permission denied), print it and continue or stop.
			// Returning the error here will stop the entire walk.
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return err // Propagate the error to stop walking
		}

		// Skip directories, we only care about files for extraction
		if info.IsDir() {
			return nil
		}

		// Check if the file is a ZIP file by its extension (case-insensitive)
		if strings.HasSuffix(strings.ToLower(info.Name()), ".zip") {
			fmt.Printf("Found ZIP file: %s\n", path)

			zipFileName := filepath.Base(path)
			zipFileDir := filepath.Dir(path)
			// Determine the extraction folder name by removing the .zip extension
			extractionFolderName := strings.TrimSuffix(zipFileName, filepath.Ext(zipFileName))
			finalDestDir := filepath.Join(zipFileDir, extractionFolderName)

			fmt.Printf("Extracting %s to %s...\n", path, finalDestDir)
			if err := extractZip(path, finalDestDir); err != nil {
				fmt.Printf("Failed to extract %s: %v\n", path, err)
				// Do NOT return err here; we want to continue processing other ZIP files
				// even if one fails.
			} else {
				fmt.Printf("Successfully extracted %s\n", path)
			}
		}
		return nil // Continue walking to the next file/directory
	})
}

func main() {
	// Define the root directory to scan for ZIP files.
	// You can change this to any specific path, e.g., "/path/to/your/zips".
	// For demonstration, we'll use the current working directory.
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		os.Exit(1)
	}

	// Verify that the specified root directory exists.
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		fmt.Printf("Error: Root directory '%s' does not exist.\n", rootDir)
		os.Exit(1)
	}

	// Start the recursive ZIP extraction process.
	if err := extractAllZipsRecursively(rootDir); err != nil {
		fmt.Printf("Recursive ZIP extraction completed with errors: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Recursive ZIP extraction completed successfully.")
	}
}

// Additional implementation at 2025-08-04 08:51:17
package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extractZip(zipFilePath string, destParentDir string) error {
	zipFileName := filepath.Base(zipFilePath)
	zipFileBaseName := strings.TrimSuffix(zipFileName, filepath.Ext(zipFileName))
	extractionDir := filepath.Join(destParentDir, zipFileBaseName)

	fmt.Printf("Extracting %s to %s...\n", zipFilePath, extractionDir)

	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file %s: %w", zipFilePath, err)
	}
	defer r.Close()

	if err := os.MkdirAll(extractionDir, 0755); err != nil {
		return fmt.Errorf("failed to create extraction directory %s: %w", extractionDir, err)
	}

	for _, f := range r.File {
		fpath := filepath.Join(extractionDir, f.Name)

		// Check for ZipSlip vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(extractionDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path in zip: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, f.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", fpath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory for file %s: %w", fpath, err)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", fpath, err)
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("failed to open file in zip %s: %w", f.Name, err)
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return fmt.Errorf("failed to copy data for file %s: %w", f.Name, err)
		}
	}
	fmt.Printf("Successfully extracted %s.\n", zipFilePath)
	return nil
}

func processDirectory(root string) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(strings.ToLower(info.Name()), ".zip") {
			destParentDir := filepath.Dir(path)
			if err := extractZip(path, destParentDir); err != nil {
				fmt.Fprintf(os.Stderr, "Error extracting ZIP %q: %v\n", path, err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking the path %q: %v\n", root, err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory_to_scan>")
		os.Exit(1)
	}

	sourceDir := os.Args[1]

	info, err := os.Stat(sourceDir)
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Directory '%s' does not exist.\n", sourceDir)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: '%s' is not a directory.\n", sourceDir)
		os.Exit(1)
	}

	fmt.Printf("Starting recursive ZIP extraction in: %s\n", sourceDir)
	processDirectory(sourceDir)
	fmt.Println("Finished ZIP extraction process.")
}