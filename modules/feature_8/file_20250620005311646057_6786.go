package main

import (
	"fmt"
	"os"
	"unicode/utf8"
)

func detectEncoding(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}

	if len(data) >= 2 {
		if data[0] == 0xFF && data[1] == 0xFE {
			return "UTF-16"
		}
		if data[0] == 0xFE && data[1] == 0xFF {
			return "UTF-16"
		}
	}

	if len(data) >= 3 {
		if data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
			return "UTF-8"
		}
	}

	if utf8.Valid(data) {
		isASCII := true
		for _, b := range data {
			if b >= 0x80 {
				isASCII = false
				break
			}
		}
		if isASCII {
			return "ASCII"
		}
		return "UTF-8"
	}

	return "Unknown"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filepath>")
		return
	}

	filePath := os.Args[1]
	encoding := detectEncoding(filePath)
	fmt.Println(encoding)
}

// Additional implementation at 2025-06-20 00:53:40
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

func detectEncoding(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "Unknown", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	bomBytes := make([]byte, 4)
	n, err := io.ReadFull(file, bomBytes)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "Unknown", fmt.Errorf("failed to read initial bytes: %w", err)
	}
	initialBytes := bomBytes[:n]

	if len(initialBytes) >= 2 {
		if initialBytes[0] == 0xFF && initialBytes[1] == 0xFE {
			return "UTF-16LE", nil
		}
		if initialBytes[0] == 0xFE && initialBytes[1] == 0xFF {
			return "UTF-16BE", nil
		}
	}

	if len(initialBytes) >= 3 {
		if initialBytes[0] == 0xEF && initialBytes[1] == 0xBB && initialBytes[2] == 0xBF {
			return "UTF-8", nil
		}
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "Unknown", fmt.Errorf("failed to seek file to start: %w", err)
	}

	reader := bufio.NewReader(file)
	var allASCII = true
	var data bytes.Buffer

	buffer := make([]byte, 4096)
	for {
		n, readErr := reader.Read(buffer)
		if n > 0 {
			chunk := buffer[:n]
			data.Write(chunk)

			for _, b := range chunk {
				if b >= 0x80 {
					allASCII = false
					break
				}
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return "Unknown", fmt.Errorf("failed to read file content: %w", readErr)
		}
	}

	fullContent := data.Bytes()

	if allASCII {
		return "ASCII", nil
	}

	if utf8.Valid(fullContent) {
		return "UTF-8", nil
	}

	return "Unknown", nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file_path>")
		return
	}

	filePath := os.Args[1]
	encoding, err := detectEncoding(filePath)
	if err != nil {
		fmt.Printf("Error detecting encoding for %s: %v\n", filePath, err)
		return
	}

	fmt.Printf("File: %s, Encoding: %s\n", filePath, encoding)
}