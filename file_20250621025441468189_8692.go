package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <filepath> <chunk_size_bytes>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	chunkSizeStr := os.Args[2]

	chunkSize, err := strconv.ParseInt(chunkSizeStr, 10, 64)
	if err != nil {
		fmt.Printf("Error parsing chunk size: %v\n", err)
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	buffer := make([]byte, chunkSize)
	chunkNum := 0

	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				if n > 0 {
					fmt.Printf("Processed chunk %d (last partial): %d bytes\n", chunkNum, n)
				}
				break
			}
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		// In a real application, you would process buffer[:n] here.
		// For demonstration, we just print the chunk details.
		fmt.Printf("Processed chunk %d: %d bytes\n", chunkNum, n)

		chunkNum++
	}

	fmt.Println("File chunking complete.")
}

// Additional implementation at 2025-06-21 02:55:45
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// chunkFile reads a large file in chunks and writes each chunk to a separate file.
// It includes progress reporting, configurable chunk size, and a specified output directory.
func chunkFile(inputFilePath string, chunkSize int64, outputDirPath string) error {
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("failed to open input file %s: %w", inputFilePath, err)
	}
	defer inputFile.Close()

	fileInfo, err := inputFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %w", inputFilePath, err)
	}
	fileSize := fileInfo.Size()

	if err := os.MkdirAll(outputDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", outputDirPath, err)
	}

	buffer := make([]byte, chunkSize)
	var bytesReadTotal int64
	chunkNum := 0
	startTime := time.Now()

	fmt.Printf("Starting chunking of '%s' (Size: %s) into chunks of %s...\n",
		inputFilePath, formatBytes(fileSize), formatBytes(chunkSize))
	fmt.Printf("Output directory: %s\n", outputDirPath)

	for {
		n, err := inputFile.Read(buffer)
		if n == 0 {
			if err == io.EOF {
				break // End of file
			}
			return fmt.Errorf("error reading from input file: %w", err)
		}

		chunkFileName := fmt.Sprintf("%s.chunk%05d", filepath.Base(inputFilePath), chunkNum)
		outputFilePath := filepath.Join(outputDirPath, chunkFileName)

		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			return fmt.Errorf("failed to create chunk file %s: %w", outputFilePath, err)
		}

		_, writeErr := outputFile.Write(buffer[:n])
		closeErr := outputFile.Close()

		if writeErr != nil {
			return fmt.Errorf("failed to write to chunk file %s: %w", outputFilePath, writeErr)
		}
		if closeErr != nil {
			return fmt.Errorf("failed to close chunk file %s: %w", outputFilePath, closeErr)
		}

		bytesReadTotal += int64(n)
		chunkNum++

		// Progress reporting
		progress := float64(bytesReadTotal) / float64(fileSize) * 100
		elapsed := time.Since(startTime)
		fmt.Printf("\rProcessed: %s / %s (%.2f%%) - Chunks: %d - Elapsed: %s",
			formatBytes(bytesReadTotal), formatBytes(fileSize), progress, chunkNum, elapsed.Round(time.Second))
	}
	fmt.Println("\nChunking complete!")
	return nil
}

// formatBytes converts bytes to a human-readable format.
func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// parseSize parses a human-readable size string (e.g., "10MB", "512KB") into bytes.
func parseSize(sizeStr string) (int64, error) {
	if len(sizeStr) == 0 {
		return 0, fmt.Errorf("empty size string")
	}

	lastChar := sizeStr[len(sizeStr)-1]
	multiplier := int64(1)
	valueStr := sizeStr

	switch lastChar {
	case 'B', 'b': // Explicit byte suffix, e.g., "100B"
		if len(sizeStr) > 1 && (sizeStr[len(sizeStr)-2] == 'K' || sizeStr[len(sizeStr)-2] == 'M' ||
			sizeStr[len(sizeStr)-2] == 'G' || sizeStr[len(sizeStr)-2] == 'T') {
			// This handles cases like "1KB", "1MB" where 'B' is part of the unit.
			// The logic below will handle these by checking the second-to-last char.
			// For simple 'B' suffix, it's just a number.
		} else {
			valueStr = sizeStr[:len(sizeStr)-1]
		}
	case 'K', 'k':
		multiplier = 1024
		valueStr = sizeStr[:len(sizeStr)-1]
	case 'M', 'm':
		multiplier = 1024 * 1024
		valueStr = sizeStr[:len(sizeStr)-1]
	case 'G', 'g':
		multiplier = 1024 * 1024 * 1024
		valueStr = sizeStr[:len(sizeStr)-1]
	case 'T', 't':
		multiplier = 1024 * 1024 * 1024 * 1024
		valueStr = sizeStr[:len(sizeStr)-1]
	}

	val, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size format: %w", err)
	}
	return val * multiplier, nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s <input_file> <chunk_size> <output_directory>\n", os.Args[0])
		fmt.Println("Example: go run main.go largefile.bin 100MB ./chunks")
		fmt.Println("Chunk size can be specified with B, KB, MB, GB, TB suffixes (e.g., 10MB, 512KB, 100B).")
		os.Exit(1)
	}

	inputFilePath := os.Args[1]
	chunkSizeStr := os.Args[2]
	outputDirPath := os.Args[3]

	chunkSize, err := parseSize(chunkSizeStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing chunk size: %v\n", err)
		os.Exit(1)
	}

	if chunkSize <= 0 {
		fmt.Fprintf(os.Stderr, "Chunk size must be a positive value.\n")
		os.Exit(1)
	}

	if err := chunkFile(inputFilePath, chunkSize, outputDirPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error chunking file: %v\n", err)
		os.Exit(1)
	}
}

// Additional implementation at 2025-06-21 02:56:55
package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Chunk represents a segment of the file.
type Chunk struct {
	ID     int
	Offset int64
	Data   []byte
	Size   int // Actual size of the data in bytes
}

// processChunk simulates processing a chunk of data.
// In a real application, this would contain the actual logic
// for what needs to be done with each chunk (e.g., parsing, compressing, uploading).
func processChunk(chunk Chunk) {
	// Simulate work
	time.Sleep(time.Millisecond * 50) // Adjust as needed for realistic simulation
}

// worker reads chunks from the input channel, processes them, and sends progress updates.
func worker(id int, chunks <-chan Chunk, progress chan<- int64, wg *sync.WaitGroup) {
	defer wg.Done()
	for chunk := range chunks {
		processChunk(chunk)
		progress <- int64(chunk.Size)
	}
}

// byteCountToHumanReadable converts a byte count to a human-readable string.
func byteCountToHumanReadable(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// monitorProgress receives progress updates and prints the current status.
func monitorProgress(totalSize int64, progress <-chan int64, done chan<- struct{}) {
	var processedBytes int64
	startTime := time.Now()
	ticker := time.NewTicker(time.Second) // Update every second
	defer ticker.Stop()

	// Initial print for 0%
	if totalSize > 0 {
		fmt.Printf("\rProcessing: %.2f%% (%s / %s) Elapsed: %s",
			0.0,
			byteCountToHumanReadable(0),
			byteCountToHumanReadable(totalSize),
			time.Duration(0))
	} else {
		fmt.Printf("\rProcessing: %s bytes processed. Elapsed: %s",
			byteCountToHumanReadable(0),
			time.Duration(0))
	}

	for {
		select {
		case bytes, ok := <-progress:
			if !ok { // Channel closed
				// Drain any remaining updates if any were sent just before close
				for bytes := range progress {
					processedBytes += bytes
				}
				// Final update
				if totalSize > 0 {
					percentage := float64(processedBytes) / float64(totalSize) * 100
					fmt.Printf("\rProcessing: %.2f%% (%s / %s) Elapsed: %s\n",
						percentage,
						byteCountToHumanReadable(processedBytes),
						byteCountToHumanReadable(totalSize),
						time.Since(startTime).Round(time.Second))
				} else {
					fmt.Printf("\rProcessing: %s bytes processed. Elapsed: %s\n",
						byteCountToHumanReadable(processedBytes),
						time.Since(startTime).Round(time.Second))
				}
				close(done)
				return
			}
			processedBytes += bytes
		case <-ticker.C:
			if totalSize > 0 {
				percentage := float64(processedBytes) / float64(totalSize) * 100
				elapsed := time.Since(startTime)
				fmt.Printf("\rProcessing: %.2f%% (%s / %s) Elapsed: %s",
					percentage,
					byteCountToHumanReadable(processedBytes),
					byteCountToHumanReadable(totalSize),
					elapsed.Round(time.Second))
			} else {
				fmt.Printf("\rProcessing: %s bytes processed. Elapsed: %s",
					byteCountToHumanReadable(processedBytes),
					time.Since(startTime).Round(time.Second))
			}
		}
	}
}

// processFileInChunks reads a large file in chunks, processes them concurrently, and reports progress.
func processFileInChunks(filePath string, chunkSize int, numWorkers int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	totalFileSize := fileInfo.Size()

	if totalFileSize == 0 {
		fmt.Println("File is empty, nothing to process.")
		return nil
	}

	fmt.Printf("Starting to process file: %s (Size: %s) with chunk size: %s and %d workers.\n",
		filePath, byteCountToHumanReadable(totalFileSize), byteCountToHumanReadable(int64(chunkSize)), numWorkers)

	// Channels for chunk data and progress updates
	chunksChan := make(chan Chunk, numWorkers*2) // Buffered channel for chunks
	progressChan := make(chan int64, numWorkers) // Buffered channel for progress updates
	doneMonitoring := make(chan struct{})        // Signal to stop progress monitor

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, chunksChan, progressChan, &wg)
	}

	// Start progress monitor
	go monitorProgress(totalFileSize, progressChan, doneMonitoring)

	// Read file and send chunks to the channel
	var chunkID int
	var currentOffset int64
	buffer := make([]byte, chunkSize)

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			// Create a copy of the buffer slice for the chunk to avoid data races
			// as the buffer will be reused.
			chunkData := make([]byte, n)
			copy(chunkData, buffer[:n])

			chunksChan <- Chunk{
				ID:     chunkID,
				Offset: currentOffset,
				Data:   chunkData,
				Size:   n,
			}
			chunkID++
			currentOffset += int64(n)
		}

		if err != nil {
			if err == io.EOF {
				break // End of file
			}
			return fmt.Errorf("error reading file: %w", err)
		}
	}

	close(chunksChan)   // Signal to workers that no more chunks will be sent
	wg.Wait()           // Wait for all workers to finish processing
	close(progressChan) // Signal to progress monitor that no more progress updates will be sent

	// Wait for the monitor to finish its final print and exit
	<-doneMonitoring

	fmt.Println("File processing complete.")
	return nil
}

func main() {
	// --- Configuration ---
	inputFilePath := "large_file.bin" // Replace with your large file path
	chunkSize := 1024 * 1024          // 1 MB chunk size
	numWorkers := 4                   // Number of concurrent workers

	// --- Create a dummy large file for testing if it doesn't exist ---
	if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
		fmt.Printf("Creating a dummy large file: %s...\n", inputFilePath)
		dummyFileSize := int64(100 * 1024 * 1024) // 100 MB
		dummyFile, err := os.Create(inputFilePath)
		if err != nil {
			fmt.Printf("Error creating dummy file: %v\n", err)
			return
		}
		defer dummyFile.Close()

		// Write some dummy data
		dummyData := make([]byte, 1024*1024) // 1MB of zeros
		for i := 0; i < len(dummyData); i++ {
			dummyData[i] = byte(i % 256) // Fill with some varying data
		}

		var written int64
		for written < dummyFileSize {
			writeSize := int64(len(dummyData))
			if written+writeSize > dummyFileSize {
				writeSize = dummyFileSize - written
			}
			_, err := dummyFile.Write(dummyData[:writeSize])
			if err != nil {
				fmt.Printf("Error writing to dummy file: %v\n", err)
				return
			}
			written += writeSize
			fmt.Printf("\rWritten: %s / %s", byteCountToHumanReadable(written), byteCountToHumanReadable(dummyFileSize))
		}
		fmt.Println("\nDummy file created.")
	}

	// --- Start processing ---
	err := processFileInChunks(inputFilePath, chunkSize, numWorkers)
	if err != nil {
		fmt.Printf("Error during file processing: %v\n", err)
	}
}

// Additional implementation at 2025-06-21 02:57:41
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Default values for chunking parameters.
const (
	defaultChunkSize = 10 * 1024 * 1024 // 10 MB
	defaultOutputDir = "chunks"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <input_file> [chunk_size_bytes] [output_directory]\n", os.Args[0])
		fmt.Printf("Defaults: chunk_size=%s, output_directory=%s\n", byteCountSI(defaultChunkSize), defaultOutputDir)
		os.Exit(1)
	}

	inputFilePath := os.Args[1]
	chunkSize := defaultChunkSize
	outputDir := defaultOutputDir

	if len(os.Args) >= 3 {
		parsedChunkSize, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil || parsedChunkSize <= 0 {
			log.Printf("Warning: Invalid chunk size '%s'. Using default %s. Error: %v\n", os.Args[2], byteCountSI(defaultChunkSize), err)
		} else {
			chunkSize = int(parsedChunkSize)
		}
	}

	if len(os.Args) >= 4 {
		outputDir = os.Args[3]
	}

	fmt.Printf("Starting file chunking process:\n")
	fmt.Printf("  Input File: %s\n", inputFilePath)
	fmt.Printf("  Chunk Size: %s\n", byteCountSI(chunkSize))
	fmt.Printf("  Output Directory: %s\n", outputDir)

	err := chunkFile(inputFilePath, chunkSize, outputDir)
	if err != nil {
		log.Fatalf("Error during file chunking: %v\n", err)
	}

	fmt.Println("File chunking completed successfully.")
}

// chunkFile reads the input file in specified chunk sizes and writes each chunk
// to a new file in the output directory. It also provides progress updates and MD5 hashes for chunks.
func chunkFile(inputFilePath string, chunkSize int, outputDir string) error {
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("failed to open input file '%s': %w", inputFilePath, err)
	}
	defer inputFile.Close()

	fileInfo, err := inputFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info for '%s': %w", inputFilePath, err)
	}
	fileSize := fileInfo.Size()

	err = os.MkdirAll(outputDir, 0755) // Create output directory with read/write/execute permissions for owner, read/execute for others.
	if err != nil {
		return fmt.Errorf("failed to create output directory '%s': %w", outputDir, err)
	}

	buffer := make([]byte, chunkSize)
	var bytesReadTotal int64
	chunkIndex := 0
	startTime := time.Now()

	fmt.Println("Processing chunks...")

	for {
		n, err := inputFile.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading input file: %w", err)
		}
		if n == 0 { // End of file reached
			break
		}

		chunkData := buffer[:n]
		chunkFileName := fmt.Sprintf("%s_chunk_%05d", filepath.Base(inputFilePath), chunkIndex) // e.g., mylargefile.txt_chunk_00000
		outputFilePath := filepath.Join(outputDir, chunkFileName)

		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			return fmt.Errorf("failed to create chunk file '%s': %w", outputFilePath, err)
		}

		_, err = outputFile.Write(chunkData)
		if err != nil {
			outputFile.Close() // Ensure file is closed on error
			return fmt.Errorf("failed to write to chunk file '%s': %w", outputFilePath, err)
		}

		md5Hash := md5.Sum(chunkData)
		md5String := hex.EncodeToString(md5Hash[:])

		fmt.Printf("  Wrote %s (MD5: %s) to %s\n", byteCountSI(int64(n)), md5String, outputFilePath)

		outputFile.Close() // Close the chunk file after writing

		bytesReadTotal += int64(n)
		chunkIndex++

		// Progress reporting
		if fileSize > 0 { // Avoid division by zero for empty files
			progress := float64(bytesReadTotal) / float64(fileSize) * 100
			elapsed := time.Since(startTime)
			fmt.Printf("  Progress: %.2f%% (%s / %s) Elapsed: %s\n", progress, byteCountSI(bytesReadTotal), byteCountSI(fileSize), elapsed.Round(time.Second))
		} else {
			fmt.Printf("  Processed: %s\n", byteCountSI(bytesReadTotal))
		}
	}

	return nil
}

// byteCountSI converts a byte count to a human-readable string using SI prefixes (e.g., 10 MB).
func byteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}