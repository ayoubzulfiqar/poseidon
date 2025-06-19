// server.go
package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

const (
	serverPort = ":8080"
	uploadDir  = "received_files"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Handling new connection from %s", conn.RemoteAddr())

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Error creating upload directory: %v", err)
		return
	}

	var filenameLen uint32
	err := binary.Read(conn, binary.LittleEndian, &filenameLen)
	if err != nil {
		if err != io.EOF {
			log.Printf("Error reading filename length from %s: %v", conn.RemoteAddr(), err)
		}
		return
	}

	filenameBytes := make([]byte, filenameLen)
	_, err = io.ReadFull(conn, filenameBytes)
	if err != nil {
		log.Printf("Error reading filename from %s: %v", conn.RemoteAddr(), err)
		return
	}
	filename := string(filenameBytes)
	safeFilename := filepath.Base(filename) 

	filePath := filepath.Join(uploadDir, safeFilename)
	log.Printf("Receiving file: %s to %s", filename, filePath)

	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file %s: %v", filePath, err)
		return
	}
	defer file.Close()

	bytesReceived, err := io.Copy(file, conn)
	if err != nil {
		log.Printf("Error receiving file content for %s: %v", filename, err)
		return
	}

	log.Printf("Successfully received %s (%d bytes) from %s", filename, bytesReceived, conn.RemoteAddr())

	ackMsg := fmt.Sprintf("File '%s' received successfully (%d bytes)", filename, bytesReceived)
	_, err = conn.Write([]byte(ackMsg))
	if err != nil {
		log.Printf("Error sending acknowledgment to %s: %v", conn.RemoteAddr(), err)
	}
}

func main() {
	log.Println("Starting file transfer server...")
	listener, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer listener.Close()
	log.Printf("Server listening on %s", serverPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

// client.go
package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

const (
	serverAddress = "localhost:8080"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: client <filepath>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	filename := filepath.Base(filePath)

	log.Printf("Attempting to connect to server at %s", serverAddress)
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()
	log.Printf("Connected to server %s", serverAddress)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", filePath, err)
	}
	defer file.Close()

	filenameBytes := []byte(filename)
	filenameLen := uint32(len(filenameBytes))
	err = binary.Write(conn, binary.LittleEndian, filenameLen)
	if err != nil {
		log.Fatalf("Error sending filename length: %v", err)
	}

	_, err = conn.Write(filenameBytes)
	if err != nil {
		log.Fatalf("Error sending filename: %v", err)
	}
	log.Printf("Sent filename: %s", filename)

	bytesSent, err := io.Copy(conn, file)
	if err != nil {
		log.Fatalf("Error sending file content: %v", err)
	}
	log.Printf("Sent %d bytes of file content for %s", bytesSent, filename)

	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Println("Server acknowledgment timed out.")
		} else if err == io.EOF {
			log.Println("Server closed connection without sending acknowledgment.")
		} else {
			log.Fatalf("Error reading server acknowledgment: %v", err)
		}
	} else {
		log.Printf("Server acknowledgment: %s", string(buffer[:n]))
	}

	log.Println("File transfer complete.")
}

// Additional implementation at 2025-06-19 00:26:13
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

const (
	serverBufferSize  = 4096
	receiveDir        = "received_files"
	connectionTimeout = 10 * time.Second // Timeout for initial metadata read
)

var (
	serverPort = flag.Int("port", 8080, "Port to listen on")
	serverWg   sync.WaitGroup
)

func main() {
	flag.Parse()

	if err := os.MkdirAll(receiveDir, 0755); err != nil {
		log.Fatalf("Failed to create receive directory %s: %v", receiveDir, err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *serverPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()
	log.Printf("Server listening on :%d", *serverPort)

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdownChan
		log.Println("Shutting down server...")
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Op == "accept" && opErr.Err.Error() == "use of closed network connection" {
				log.Println("Listener closed, server stopping accept loop.")
				break
			}
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		serverWg.Add(1)
		go handleConnection(conn)
	}

	serverWg.Wait()
	log.Println("Server gracefully shut down.")
}

func handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		serverWg.Done()
		log.Printf("Connection from %s closed.", conn.RemoteAddr())
	}()

	log.Printf("Handling connection from %s", conn.RemoteAddr())

	conn.SetReadDeadline(time.Now().Add(connectionTimeout))

	var filenameLen uint32
	err := binary.Read(conn, binary.LittleEndian, &filenameLen)
	if err != nil {
		log.Printf("[%s] Failed to read filename length: %v", conn.RemoteAddr(), err)
		return
	}

	filenameBytes := make([]byte, filenameLen)
	_, err = io.ReadFull(conn, filenameBytes)
	if err != nil {
		log.Printf("[%s] Failed to read filename: %v", conn.RemoteAddr(), err)
		return
	}
	filename := string(filenameBytes)
	log.Printf("[%s] Receiving file: %s", conn.RemoteAddr(), filename)

	var fileSize int64
	err = binary.Read(conn, binary.LittleEndian, &fileSize)
	if err != nil {
		log.Printf("[%s] Failed to read file size: %v", conn.RemoteAddr(), err)
		return
	}
	log.Printf("[%s] File size: %d bytes", conn.RemoteAddr(), fileSize)

	conn.SetReadDeadline(time.Time{})

	safeFilename := filepath.Base(filename)
	filePath := filepath.Join(receiveDir, safeFilename)

	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("[%s] Failed to create file %s: %v", conn.RemoteAddr(), filePath, err)
		return
	}
	defer file.Close()

	var receivedBytes int64
	buffer := make([]byte, serverBufferSize)
	startTime := time.Now()
	lastProgress := -1

	for receivedBytes < fileSize {
		readAmount := fileSize - receivedBytes
		if readAmount > serverBufferSize {
			readAmount = serverBufferSize
		}

		n, err := conn.Read(buffer[:readAmount])
		if err != nil {
			if err != io.EOF {
				log.Printf("[%s] Error reading from connection: %v", conn.RemoteAddr(), err)
			}
			break
		}

		if n == 0 {
			log.Printf("[%s] Client closed connection prematurely.", conn.RemoteAddr())
			break
		}

		_, err = file.Write(buffer[:n])
		if err != nil {
			log.Printf("[%s] Error writing to file: %v", conn.RemoteAddr(), err)
			return
		}
		receivedBytes += int64(n)

		currentProgress := int(float64(receivedBytes) / float64(fileSize) * 100)
		if currentProgress > lastProgress && currentProgress%5 == 0 {
			log.Printf("[%s] Received: %d/%d bytes (%d%%) for %s", conn.RemoteAddr(), receivedBytes, fileSize, currentProgress, safeFilename)
			lastProgress = currentProgress
		}
	}

	if receivedBytes == fileSize {
		duration := time.Since(startTime)
		log.Printf("[%s] Successfully received file %s (%d bytes) in %s", conn.RemoteAddr(), safeFilename, fileSize, duration)
	} else {
		log.Printf("[%s] File transfer incomplete for %s. Expected %d, received %d.", conn.RemoteAddr(), safeFilename, fileSize, receivedBytes)
		os.Remove(filePath)
	}
}
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

const (
	clientBufferSize = 4096
	dialTimeout      = 5 * time.Second
)

var (
	serverAddress = flag.String("server", "localhost:8080", "Server address (host:port)")
	clientFilePath = flag.String("file", "", "Path to the file to send")
)

func main() {
	flag.Parse()

	if *clientFilePath == "" {
		log.Fatal("Please specify a file to send using -file flag.")
	}

	fileInfo, err := os.Stat(*clientFilePath)
	if os.IsNotExist(err) {
		log.Fatalf("File not found: %s", *clientFilePath)
	}
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}
	if fileInfo.IsDir() {
		log.Fatalf("Cannot send a directory: %s", *clientFilePath)
	}

	conn, err := net.DialTimeout("tcp", *serverAddress, dialTimeout)
	if err != nil {
		log.Fatalf("Failed to connect to server %s: %v", *serverAddress, err)
	}
	defer conn.Close()
	log.Printf("Connected to server %s", *serverAddress)

	file, err := os.Open(*clientFilePath)
	if err != nil {
		log.Fatalf("Failed to open file %s: %v", *clientFilePath, err)
	}
	defer file.Close()

	filename := filepath.Base(*clientFilePath)
	filenameLen := uint32(len(filename))
	fileSize := fileInfo.Size()

	err = binary.Write(conn, binary.LittleEndian, filenameLen)
	if err != nil {
		log.Fatalf("Failed to write filename length: %v", err)
	}

	_, err = conn.Write([]byte(filename))
	if err != nil {
		log.Fatalf("Failed to write filename: %v", err)
	}

	err = binary.Write(conn, binary.LittleEndian, fileSize)
	if err != nil {
		log.Fatalf("Failed to write file size: %v", err)
	}

	log.Printf("Sending file %s (%d bytes) to %s", filename, fileSize, *serverAddress)

	var sentBytes int64
	buffer := make([]byte, clientBufferSize)
	startTime := time.Now()
	lastProgress := -1

	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error reading from file: %v", err)
		}

		if n == 0 {
			break
		}

		bytesWritten, writeErr := conn.Write(buffer[:n])
		if writeErr != nil {
			log.Fatalf("Error writing to connection: %v", writeErr)
		}
		sentBytes += int64(bytesWritten)

		currentProgress := int(float64(sentBytes) / float64(fileSize) * 100)
		if currentProgress > lastProgress && currentProgress%5 == 0 {
			log.Printf("Sent: %d/%d bytes (%d%%)", sentBytes, fileSize, currentProgress)
			lastProgress = currentProgress
		}
	}

	duration := time.Since(startTime)
	if sentBytes == fileSize {
		log.Printf("Successfully sent file %s (%d bytes) in %s", filename, fileSize, duration)
	} else {
		log.Printf("File transfer incomplete. Expected %d, sent %d.", fileSize, sentBytes)
	}
}

// Additional implementation at 2025-06-19 00:27:16
package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	SERVER_ADDR = "localhost:8080"
	BUFFER_SIZE = 4096 // Buffer size for file transfer
	UPLOAD_DIR  = "uploads"
)

// Helper function to write a length-prefixed string
func writeString(conn net.Conn, s string) error {
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(s)))
	if _, err := conn.Write(lenBytes); err != nil {
		return fmt.Errorf("failed to write string length: %w", err)
	}
	if _, err := conn.Write([]byte(s)); err != nil {
		return fmt.Errorf("failed to write string content: %w", err)
	}
	return nil
}

// Helper function to read a length-prefixed string
func readString(conn net.Conn) (string, error) {
	lenBytes := make([]byte, 4)
	if _, err := io.ReadFull(conn, lenBytes); err != nil {
		return "", fmt.Errorf("failed to read string length: %w", err)
	}
	length := binary.BigEndian.Uint32(lenBytes)
	if length == 0 {
		return "", nil
	}
	strBytes := make([]byte, length)
	if _, err := io.ReadFull(conn, strBytes); err != nil {
		return "", fmt.Errorf("failed to read string content: %w", err)
	}
	return string(strBytes), nil
}

// Helper function to write an int64
func writeInt64(conn net.Conn, i int64) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	if _, err := conn.Write(b); err != nil {
		return fmt.Errorf("failed to write int64: %w", err)
	}
	return nil
}

// Helper function to read an int64
func readInt64(conn net.Conn) (int64, error) {
	b := make([]byte, 8)
	if _, err := io.ReadFull(conn, b); err != nil {
		return 0, fmt.Errorf("failed to read int64: %w", err)
	}
	return int64(binary.BigEndian.Uint64(b)), nil
}

// Server-side logic
func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Handling connection from %s\n", conn.RemoteAddr())

	command, err := readString(conn)
	if err != nil {
		fmt.Printf("Error reading command: %v\n", err)
		return
	}

	switch command {
	case "UPLOAD":
		handleUpload(conn)
	case "DOWNLOAD":
		handleDownload(conn)
	case "LIST":
		handleList(conn)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		writeString(conn, "ERROR")
		writeString(conn, "Unknown command")
	}
}

func handleUpload(conn net.Conn) {
	filename, err := readString(conn)
	if err != nil {
		fmt.Printf("Error reading upload filename: %v\n", err)
		writeString(conn, "ERROR")
		writeString(conn, fmt.Sprintf("Failed to read filename: %v", err))
		return
	}

	fileSize, err := readInt64(conn)
	if err != nil {
		fmt.Printf("Error reading upload file size for %s: %v\n", filename, err)
		writeString(conn, "ERROR")
		writeString(conn, fmt.Sprintf("Failed to read file size: %v", err))
		return
	}

	// Sanitize filename to prevent directory traversal
	safeFilename := filepath.Base(filename)
	filePath := filepath.Join(UPLOAD_DIR, safeFilename)

	fmt.Printf("Receiving file: %s (%d bytes) to %s\n", safeFilename, fileSize, filePath)

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filePath, err)
		writeString(conn, "ERROR")
		writeString(conn, fmt.Sprintf("Failed to create file on server: %v", err))
		return
	}
	defer file.Close()

	var receivedBytes int64
	buffer := make([]byte, BUFFER_SIZE)
	for receivedBytes < fileSize {
		bytesToRead := int64(BUFFER_SIZE)
		if fileSize-receivedBytes < bytesToRead {
			bytesToRead = fileSize - receivedBytes
		}

		n, err := io.ReadFull(conn, buffer[:bytesToRead])
		if err != nil {
			fmt.Printf("Error reading file content for %s: %v\n", safeFilename, err)
			writeString(conn, "ERROR")
			writeString(conn, fmt.Sprintf("Failed to read file content: %v", err))
			return
		}

		if _, err := file.Write(buffer[:n]); err != nil {
			fmt.Printf("Error writing file content for %s: %v\n", safeFilename, err)
			writeString(conn, "ERROR")
			writeString(conn, fmt.Sprintf("Failed to write file content: %v", err))
			return
		}
		receivedBytes += int64(n)
	}

	fmt.Printf("Successfully received %s\n", safeFilename)
	writeString(conn, "OK")
	writeString(conn, "File uploaded successfully")
}

func handleDownload(conn net.Conn) {
	filename, err := readString(conn)
	if err != nil {
		fmt.Printf("Error reading download filename: %v\n", err)
		writeString(conn, "ERROR")
		writeString(conn, fmt.Sprintf("Failed to read filename: %v", err))
		return
	}

	safeFilename := filepath.Base(filename)
	filePath := filepath.Join(UPLOAD_DIR, safeFilename)

	fmt.Printf("Client requested download: %s from %s\n", safeFilename, filePath)

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Printf("File not found: %s\n", filePath)
		writeString(conn, "ERROR")
		writeString(conn, "File not found")
		return
	}
	if err != nil {
		fmt.Printf("Error stating file %s: %v\n", filePath, err)
		writeString(conn, "ERROR")
		writeString(conn, fmt.Sprintf("Failed to access file: %v", err))
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		writeString(conn, "ERROR")
		writeString(conn, fmt.Sprintf("Failed to open file: %v", err))
		return
	}
	defer file.Close()

	writeString(conn, "OK")
	writeInt64(conn, fileInfo.Size())

	buffer := make([]byte, BUFFER_SIZE)
	sentBytes, err := io.CopyBuffer(conn, file, buffer)
	if err != nil {
		fmt.Printf("Error sending file content for %s: %v\n", safeFilename, err)
		// No need to send ERROR status here, as we already sent OK and file size.
		// The client will detect the broken pipe.
		return
	}
	fmt.Printf("Successfully sent %s (%d bytes)\n", safeFilename, sentBytes)
}

func handleList(conn net.Conn) {
	fmt.Println("Client requested file list")

	files, err := os.ReadDir(UPLOAD_DIR)
	if err != nil {
		fmt.Printf("Error reading upload directory: %v\n", err)
		writeString(conn, "ERROR")
		writeString(conn, fmt.Sprintf("Failed to list files: %v", err))
		return
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	writeString(conn, "OK")
	writeInt64(conn, int64(len(fileNames))) // Send count of files

	for _, name := range fileNames {
		writeString(conn, name) // Send each filename
	}