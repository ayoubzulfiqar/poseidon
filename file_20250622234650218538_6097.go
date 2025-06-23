// server.go
package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"
	UPLOAD_DIR  = "uploads" // Directory to save received files
)

func main() {
	fmt.Println("Server Running...")
	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)

	listener, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	// Create upload directory if it doesn't exist
	if _, err := os.Stat(UPLOAD_DIR); os.IsNotExist(err) {
		os.Mkdir(UPLOAD_DIR, 0755)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr().String())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 1. Read filename length (uint32)
	var filenameLen uint32
	err := binary.Read(conn, binary.BigEndian, &filenameLen)
	if err != nil {
		fmt.Println("Error reading filename length:", err)
		return
	}

	// 2. Read filename
	filenameBytes := make([]byte, filenameLen)
	_, err = io.ReadFull(conn, filenameBytes)
	if err != nil {
		fmt.Println("Error reading filename:", err)
		return
	}
	filename := string(filenameBytes)
	fmt.Printf("Receiving file: %s from %s\n", filename, conn.RemoteAddr().String())

	// Sanitize filename to prevent path traversal
	safeFilename := filepath.Base(filename)
	filePath := filepath.Join(UPLOAD_DIR, safeFilename)

	// 3. Read file content length (uint64)
	var fileContentLen uint64
	err = binary.Read(conn, binary.BigEndian, &fileContentLen)
	if err != nil {
		fmt.Println("Error reading file content length:", err)
		return
	}

	// 4. Create file and write content
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Use io.CopyN to limit the bytes read from the connection
	// This prevents reading past the end of the file content into the next message
	bytesReceived, err := io.CopyN(file, conn, int64(fileContentLen))
	if err != nil {
		fmt.Println("Error receiving file content:", err)
		// Attempt to send an error response back to the client
		conn.Write([]byte("ERROR: " + err.Error()))
		return
	}

	fmt.Printf("Received %d bytes for file %s\n", bytesReceived, safeFilename)

	if bytesReceived != int64(fileContentLen) {
		fmt.Printf("Warning: Expected %d bytes, but received %d for file %s\n", fileContentLen, bytesReceived, safeFilename)
		conn.Write([]byte("ERROR: Incomplete file transfer"))
		return
	}

	fmt.Printf("File %s received successfully and saved to %s\n", safeFilename, filePath)
	conn.Write([]byte("SUCCESS: File received"))
}

// client.go
package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run client.go <filepath_to_send>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	filename := filepath.Base(filePath)

	fmt.Println("Client connecting to " + SERVER_HOST + ":" + SERVER_PORT)
	conn, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err.Error())
		os.Exit(1)
	}

	fileSize := uint64(fileInfo.Size())
	filenameBytes := []byte(filename)
	filenameLen := uint32(len(filenameBytes))

	// 1. Send filename length
	err = binary.Write(conn, binary.BigEndian, filenameLen)
	if err != nil {
		fmt.Println("Error sending filename length:", err)
		return
	}

	// 2. Send filename
	_, err = conn.Write(filenameBytes)
	if err != nil {
		fmt.Println("Error sending filename:", err)
		return
	}

	// 3. Send file content length
	err = binary.Write(conn, binary.BigEndian, fileSize)
	if err != nil {
		fmt.Println("Error sending file size:", err)
		return
	}

	// 4. Send file content
	bytesSent, err := io.Copy(conn, file)
	if err != nil {
		fmt.Println("Error sending file content:", err)
		return
	}

	fmt.Printf("Sent %d bytes for file %s\n", bytesSent, filename)

	if bytesSent != int64(fileSize) {
		fmt.Printf("Warning: Expected to send %d bytes, but sent %d\n", fileSize, bytesSent)
	}

	// Read response from server
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // Set a read deadline
	n, err := conn.Read(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("Server response timeout.")
		} else {
			fmt.Println("Error reading server response:", err)
		}
		return
	}
	fmt.Println("Server response:", string(buffer[:n]))
}

// Additional implementation at 2025-06-22 23:48:05
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	defaultServerPort = "8080"
	bufferSize        = 4096
)

type FileMetadata struct {
	Filename string
	Filesize int64
}

func startServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server: Listen error: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Printf("Server: Listening on :%s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Server: Accept error: %v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Server: Connection from %s\n", conn.RemoteAddr())

	decoder := json.NewDecoder(conn)
	var metadata FileMetadata
	if err := decoder.Decode(&metadata); err != nil {
		fmt.Fprintf(os.Stderr, "Server: Metadata decode error from %s: %v\n", conn.RemoteAddr(), err)
		return
	}

	fmt.Printf("Server: Receiving '%s' (%s) from %s\n", metadata.Filename, byteCountToHumanReadable(metadata.Filesize), conn.RemoteAddr())

	targetDir := "received_files"
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Server: Directory creation error '%s': %v\n", targetDir, err)
		return
	}

	filePath := filepath.Join(targetDir, filepath.Base(metadata.Filename))
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server: File creation error '%s': %v\n", filePath, err)
		return
	}
	defer file.Close()

	receivedBytes := int64(0)
	buffer := make([]byte, bufferSize)
	startTime := time.Now()

	for receivedBytes < metadata.Filesize {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "Server: Read error from %s: %v\n", conn.RemoteAddr(), err)
			}
			break
		}

		_, err = file.Write(buffer[:n])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Server: Write error to file '%s': %v\n", filePath, err)
			return
		}
		receivedBytes += int64(n)

		progress := float64(receivedBytes) / float64(metadata.Filesize) * 100
		fmt.Printf("\rServer: %s: %.2f%% (%s/%s)", metadata.Filename, progress, byteCountToHumanReadable(receivedBytes), byteCountToHumanReadable(metadata.Filesize))
	}
	fmt.Println()

	duration := time.Since(startTime)
	if receivedBytes == metadata.Filesize {
		fmt.Printf("Server: Received '%s' (%s) from %s in %s\n", metadata.Filename, byteCountToHumanReadable(receivedBytes), conn.RemoteAddr(), duration.Round(time.Millisecond))
	} else {
		fmt.Printf("Server: Partial transfer of '%s' from %s: %s received (expected %s) in %s\n", metadata.Filename, conn.RemoteAddr(), byteCountToHumanReadable(receivedBytes), byteCountToHumanReadable(metadata.Filesize), duration.Round(time.Millisecond))
	}
}

func startClient(serverAddr, filePath string) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Client: File info error for '%s': %v\n", filePath, err)
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Client: Connect error to %s: %v\n", serverAddr, err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Printf("Client: Connected to %s\n", serverAddr)

	metadata := FileMetadata{
		Filename: filepath.Base(filePath),
		Filesize: fileInfo.Size(),
	}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(metadata); err != nil {
		fmt.Fprintf(os.Stderr, "Client: Metadata encode error: %v\n", err)
		return
	}
	fmt.Printf("Client: Sending '%s' (%s) to %s\n", metadata.Filename, byteCountToHumanReadable(metadata.Filesize), serverAddr)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Client: File open error '%s': %v\n", filePath, err)
		return
	}
	defer file.Close()

	sentBytes := int64(0)
	buffer := make([]byte, bufferSize)
	startTime := time.Now()

	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "Client: Read error from file '%s': %v\n", filePath, err)
			return
		}

		_, writeErr := conn.Write(buffer[:n])
		if writeErr != nil {
			fmt.Fprintf(os.Stderr, "Client: Network write error: %v\n", writeErr)
			return
		}
		sentBytes += int64(n)

		progress := float64(sentBytes) / float64(metadata.Filesize) * 100
		fmt.Printf("\rClient: %s: %.2f%% (%s/%s)", metadata.Filename, progress, byteCountToHumanReadable(sentBytes), byteCountToHumanReadable(metadata.Filesize))
	}
	fmt.Println()

	duration := time.Since(startTime)
	fmt.Printf("Client: Sent '%s' (%s) to %s in %s\n", metadata.Filename, byteCountToHumanReadable(sentBytes), serverAddr, duration.Round(time.Millisecond))
}

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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  Server: go run main.go server [port]")
		fmt.Println("  Client: go run main.go client [server_addr:port] [file_to_send]")
		os.Exit(1)
	}

	mode := os.Args[1]

	switch mode {
	case "server":
		port := defaultServerPort
		if len(os.Args) > 2 {
			port = os.Args[2]
			if _, err := strconv.Atoi(port); err != nil {
				fmt.Fprintf(os.Stderr, "Invalid port: %s\n", port)
				os.Exit(1)
			}
		}
		startServer(port)
	case "client":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run main.go client [server_addr:port] [file_to_send]")
			os.Exit(1)
		}
		serverAddr := os.Args[2]
		filePath := os.Args[3]
		startClient(serverAddr, filePath)
	default:
		fmt.Println("Invalid mode. Use 'server' or 'client'.")
		os.Exit(1)
	}
}

// Additional implementation at 2025-06-22 23:49:25
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"

	SERVER_FILES_DIR     = "server_files"
	CLIENT_DOWNLOADS_DIR = "client_downloads"

	BUFFER_SIZE = 4096 // 4KB buffer for file transfers
)

// --- Server Logic ---

func startServer() {
	// Create server files directory if it doesn't exist
	if _, err := os.Stat(SERVER_FILES_DIR); os.IsNotExist(err) {
		os.Mkdir(SERVER_FILES_DIR, 0755)
	}

	fmt.Printf("Starting server on %s:%s...\n", SERVER_HOST, SERVER_PORT)
	listener, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server listening. Waiting for connections...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		go handleClientConnection(conn)
	}
}

func handleClientConnection(conn net.Conn) {
	fmt.Printf("Client connected from %s\n", conn.RemoteAddr().String())
	defer func() {
		fmt.Printf("Client %s disconnected.\n", conn.RemoteAddr().String())
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Minute)) // Set a timeout for inactivity
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// Client closed connection
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Printf("Client %s timed out.\n", conn.RemoteAddr().String())
			} else {
				fmt.Println("Error reading from client:", err.Error())
			}
			break
		}

		message = strings.TrimSpace(message)
		parts := strings.Fields(message)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])

		switch command {
		case "LS":
			fmt.Printf("Client %s requested file list.\n", conn.RemoteAddr().String())
			listFiles(writer)
		case "GET":
			if len(parts) < 2 {
				fmt.Fprintf(writer, "ERROR: Missing filename for GET command.\n")
				writer.Flush()
				continue
			}
			filename := parts[1]
			fmt.Printf("Client %s requested file: %s\n", conn.RemoteAddr().String(), filename)
			sendFile(writer, filename)
		case "PUT":
			if len(parts) < 3 {
				fmt.Fprintf(writer, "ERROR: Missing filename or size for PUT command.\n")
				writer.Flush()
				continue
			}
			filename := parts[1]
			sizeStr := parts[2]
			fileSize, err := strconv.ParseInt(sizeStr, 10, 64)
			if err != nil {
				fmt.Fprintf(writer, "ERROR: Invalid file size for PUT command.\n")
				writer.Flush()
				continue
			}
			fmt.Printf("Client %s wants to send file: %s (size: %d bytes)\n", conn.RemoteAddr().String(), filename, fileSize)
			receiveFile(reader, writer, filename, fileSize)
		case "QUIT":
			fmt.Printf("Client %s sent QUIT command.\n", conn.RemoteAddr().String())
			fmt.Fprintf(writer, "BYE\n")
			writer.Flush()
			return // Exit handleClientConnection goroutine
		default:
			fmt.Fprintf(writer, "ERROR: Unknown command '%s'\n", command)
			writer.Flush()
		}
	}
}

func listFiles(writer *bufio.Writer) {
	files, err := ioutil.ReadDir(SERVER_FILES_DIR)
	if err != nil {
		fmt.Fprintf(writer, "ERROR: Could not read server files directory: %s\n", err.Error())
		writer.Flush()
		return
	}

	fmt.Fprintf(writer, "FILE_LIST\n")
	for _, file := range files {
		if !file.IsDir() {
			fmt.Fprintf(writer, "%s (Size: %d bytes)\n", file.Name(), file.Size())
		}
	}
	fmt.Fprintf(writer, "END_LIST\n")
	writer.Flush()
}

func sendFile(writer *bufio.Writer, filename string) {
	filePath := filepath.Join(SERVER_FILES_DIR, filename)
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Fprintf(writer, "FILE_NOT_FOUND: %s\n", filename)
		writer.Flush()
		return
	}
	if err != nil {
		fmt.Fprintf(writer, "ERROR_STAT_FILE: %s\n", err.Error())
		writer.Flush()
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(writer, "ERROR_OPEN_FILE: %s\n", err.Error())
		writer.Flush()
		return
	}
	defer file.Close()

	fmt.Fprintf(writer, "FILE_START %s %d\n", filename, fileInfo.Size())
	writer.Flush() // Ensure header is sent before file data

	bytesSent, err := io.CopyBuffer(writer, file, make([]byte, BUFFER_SIZE))
	if err != nil {
		fmt.Printf("Error sending file data: %s\n", err.Error())
		// Client will detect incomplete transfer.
	} else {
		fmt.Printf("Sent %d bytes for file %s.\n", bytesSent, filename)
	}
	writer.Flush() // Ensure all file data is flushed
	// No explicit FILE_END needed, client will read until EOF or specified size
}

func receiveFile(reader *bufio.Reader, writer *bufio.Writer, filename string, fileSize int64) {
	filePath := filepath.Join(SERVER_FILES_DIR, filename)

	// Check if file already exists and handle it (e.g., overwrite or rename)
	// For simplicity, we'll overwrite.
	outputFile, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintf(writer, "ERROR_CREATE_FILE: %s\n", err.Error())
		writer.Flush()
		return
	}
	defer outputFile.Close()

	fmt.Fprintf(writer, "READY_TO_RECEIVE\n")
	writer.Flush()

	bytesReceived, err := io.CopyN(outputFile, reader, fileSize)
	if err != nil {
		fmt.Printf("Error receiving file data for %s: %s\n", filename, err.Error())
		fmt.Fprintf(writer, "ERROR_RECEIVING_DATA: %s\n", err.Error())
		os.Remove(filePath) // Clean up incomplete file
	} else if bytesReceived != fileSize {
		fmt.Printf("Incomplete file transfer for %s. Expected %d, got %d.\n", filename, fileSize, bytesReceived)
		fmt.Fprintf(writer, "ERROR_INCOMPLETE_TRANSFER: Expected %d, got %d\n", fileSize, bytesReceived)
		os.Remove(filePath) // Clean up incomplete file
	} else {
		fmt.Printf("Successfully received %d bytes for file %s.\n", bytesReceived, filename)
		fmt.Fprintf(writer, "FILE_RECEIVED_OK\n")
	}
	writer.Flush()
}

// --- Client Logic ---

func startClient() {
	// Create client downloads directory if it doesn't exist
	if _, err := os.Stat(CLIENT_DOWNLOADS_DIR); os.IsNotExist(err) {
		os.Mkdir(CLIENT_DOWNLOADS_DIR, 0755)
	}

	fmt.Printf("Connecting to server at %s:%s...\n", SERVER_HOST, SERVER_PORT)
	conn, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connected to server.")

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter commands: ls, get <filename>, put <filename>, quit")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break // EOF or error
		}
		input := scanner.Text()
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])

		switch command {
		case "LS":
			fmt.Fprintf(writer, "LS\n")
			writer.Flush()
			handleLsResponse(reader)
		case "GET":
			if len(parts) < 2 {
				fmt.Println("Usage: get <filename>")
				continue
			}
			filename := parts[1]
			fmt.Fprintf(writer, "GET %s\n", filename)
			writer.Flush()
			handleGetResponse(reader, filename)
		case "PUT":
			if len(parts) < 2 {
				fmt.Println("Usage: put <filename>")
				continue
			}
			filename := parts[1]
			sendClientFile(writer, reader, filename)
		case "QUIT":
			fmt.Fprintf(writer, "QUIT\n")
			writer.Flush()
			response, _ := reader.ReadString('\n')
			fmt.Println(strings.TrimSpace(response))
			return // Exit client
		default:
			fmt.Println("Unknown command. Please use: ls, get <filename>, put <filename>, quit")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}

func handleLsResponse(reader *bufio.Reader) {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading LS response:", err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "FILE_LIST" {
			fmt.Println("--- Server Files ---")
			continue
		}
		if line == "END_LIST" {
			fmt.Println("--------------------")
			return
		}
		if strings.HasPrefix(line, "ERROR:") {
			fmt.Println("Server Error:", line)
			return
		}
		fmt.Println(line)
	}
}

func handleGetResponse(reader *bufio.Reader, requestedFilename string) {
	header, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading GET response header:", err)
		return
	}
	header = strings.TrimSpace(header)
	parts := strings.Fields(header)

	if len(parts) == 0 {
		fmt.Println("Empty response from server.")
		return
	}

	switch parts[0] {
	case "FILE_NOT_FOUND:":
		fmt.Printf("Server: File '%s' not found.\n", requestedFilename)
		return
	case "ERROR_STAT_FILE:", "ERROR_OPEN_FILE:", "ERROR_RECEIVING_DATA:", "ERROR_INCOMPLETE_TRANSFER:":
		fmt.Printf("Server Error: %s\n", header)
		return
	case "FILE_START":
		if len(parts) < 3 {
			fmt.Println("Invalid FILE_START header from server.")
			return
		}
		filename := parts[1]
		sizeStr := parts[2]
		fileSize, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			fmt.Println("Invalid file size in FILE_START header:", err)
			return
		}

		fmt.Printf("Receiving file '%s' (%d bytes)...\n", filename, fileSize)
		filePath := filepath.Join(CLIENT_DOWNLOADS_DIR, filename)
		outputFile, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating local file:", err)
			return
		}
		defer outputFile.Close()

		bytesReceived, err := io.CopyN(outputFile, reader, fileSize)
		if err != nil && err != io.EOF { // io.EOF is expected if the server closes connection after sending all bytes
			fmt.Printf("Error receiving file data for %s: %s\n", filename, err.Error())
			os.Remove(filePath) // Clean up incomplete file
			return
		}

		if bytesReceived != fileSize {
			fmt.Printf("Incomplete file transfer for %s. Expected %d, got %d.\n", filename, fileSize, bytesReceived)
