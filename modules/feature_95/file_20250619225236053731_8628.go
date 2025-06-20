package main

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <folder_to_monitor>", os.Args[0])
	}
	folderToMonitor := os.Args[1]

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Printf("Change detected: Operation=%s, Path=%s", event.Op.String(), event.Name)
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Printf("New file created: %s", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Monitoring error:", err)
			}
		}
	}()

	err = watcher.Add(folderToMonitor)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Monitoring folder: %s. Press Ctrl+C to stop.", folderToMonitor)
	<-done
}

// Additional implementation at 2025-06-19 22:53:19
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

// FolderMonitor monitors a specified folder for file system events.
type FolderMonitor struct {
	folderPath string
	watcher    *fsnotify.Watcher
	logger     *log.Logger
	logFile    *os.File
	done       chan bool
	wg         sync.WaitGroup
}

// NewFolderMonitor creates a new FolderMonitor instance.
func NewFolderMonitor(folderPath, logFilePath string) (*FolderMonitor, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	// Ensure the folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("folder does not exist: %s", folderPath)
	}
	if err != nil {
		return nil, fmt.Errorf("error checking folder: %w", err)
	}

	// Set up logging to a file and console
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		watcher.Close()
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create a multi-writer for logging to both file and console
	writers := []io.Writer{os.Stdout}
	if logFile != nil {
		writers = append(writers, logFile)
	}
	multiWriter := io.MultiWriter(writers...)
	logger := log.New(multiWriter, "[MONITOR] ", log.Ldate|log.Ltime|log.Lshortfile)

	return &FolderMonitor{
		folderPath: folderPath,
		watcher:    watcher,
		logger:     logger,
		logFile:    logFile,
		done:       make(chan bool),
	}, nil
}

// Start begins monitoring the folder.
func (fm *FolderMonitor) Start() error {
	fm.logger.Printf("Starting folder monitor for: %s", fm.folderPath)

	// Add the folder to the watcher
	err := fm.watcher.Add(fm.folderPath)
	if err != nil {
		return fmt.Errorf("failed to add folder to watcher: %w", err)
	}

	fm.wg.Add(1)
	go fm.run()

	return nil
}

// run is the main event loop for the monitor.
func (fm *FolderMonitor) run() {
	defer fm.wg.Done()
	for {
		select {
		case event, ok := <-fm.watcher.Events:
			if !ok {
				fm.logger.Println("Watcher events channel closed.")
				return
			}
			fm.logEvent(event)
		case err, ok := <-fm.watcher.Errors:
			if !ok {
				fm.logger.Println("Watcher errors channel closed.")
				return
			}
			fm.logger.Printf("Watcher error: %v", err)
		case <-fm.done:
			fm.logger.Println("Stopping folder monitor event loop.")
			return
		}
	}
}

// logEvent logs the file system event.
func (fm *FolderMonitor) logEvent(event fsnotify.Event) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] File: %s, Operation: %s", timestamp, event.Name, event.Op.String())

	// Additional details based on operation type
	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		logMessage += " (New file/directory created)"
	case event.Op&fsnotify.Write == fsnotify.Write:
		logMessage += " (File content modified)"
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		logMessage += " (File/directory removed)"
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		logMessage += " (File/directory renamed/moved)"
	case event.Op&fsnotify.Chmod == fsnotify.Chmod:
		logMessage += " (File permissions changed)"
	}

	fm.logger.Println(logMessage)
}

// Stop gracefully stops the folder monitor.
func (fm *FolderMonitor) Stop() {
	fm.logger.Println("Stopping folder monitor...")
	close(fm.done)
	fm.watcher.Close()
	fm.wg.Wait() // Wait for the run goroutine to finish
	if fm.logFile != nil {
		fm.logFile.Close()
	}
	fm.logger.Println("Folder monitor stopped.")
}

func main() {
	// Configuration
	folderToMonitor := "./watched_folder" // Change this to your desired folder
	logFilePath := "./monitor.log"       // Log file path

	// Create the folder if it doesn't exist for demonstration purposes
	if _, err := os.Stat(folderToMonitor); os.IsNotExist(err) {
		err = os.MkdirAll(folderToMonitor, 0755)
		if err != nil {
			log.Fatalf("Failed to create folder %s: %v", folderToMonitor, err)
		}
		log.Printf("Created folder: %s", folderToMonitor)
	}

	monitor, err := NewFolderMonitor(folderToMonitor, logFilePath)
	if err != nil {
		log.Fatalf("Failed to initialize folder monitor: %v", err)
	}
	defer monitor.Stop() // Ensure stop is called on exit

	err = monitor.Start()
	if err != nil {
		log.Fatalf("Failed to start folder monitor: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-sigChan
	monitor.logger.Printf("Received signal: %v. Shutting down...", sig)
}

// Additional implementation at 2025-06-19 22:54:23
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	watchDir string
	logFile  string
	logger   *log.Logger
)

func init() {
	flag.StringVar(&watchDir, "dir", ".", "Directory to monitor for changes")
	flag.StringVar(&logFile, "log", "folder_monitor.log", "Path to the log file")
	flag.Parse()

	absDir, err := filepath.Abs(watchDir)
	if err != nil {
		log.Fatalf("Error getting absolute path for directory %s: %v", watchDir, err)
	}
	watchDir = absDir

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file %s: %v", logFile, err)
	}
	logger = log.New(file, "", 0)

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf("Monitoring directory: %s", watchDir)
	logger.Printf("Monitoring directory: %s", watchDir)
	log.Printf("Logging events to file: %s", logFile)
	logger.Printf("Logging events to file: %s", logFile)
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logEvent(event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
				logger.Printf("Watcher error: %v", err)
			case sig := <-sigChan:
				log.Printf("Received signal %v, shutting down...", sig)
				logger.Printf("Received signal %v, shutting down...", sig)
				done <- true
				return
			}
		}
	}()

	err = watcher.Add(watchDir)
	if err != nil {
		log.Fatalf("Failed to add directory %s to watcher: %v", watchDir, err)
	}

	log.Printf("Started monitoring %s. Press Ctrl+C to stop.", watchDir)
	logger.Printf("Started monitoring %s. Press Ctrl+C to stop.", watchDir)

	<-done
	log.Println("Monitor stopped.")
	logger.Println("Monitor stopped.")
}

func logEvent(event fsnotify.Event) {
	eventType := ""
	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		eventType = "CREATE"
	case event.Op&fsnotify.Write == fsnotify.Write:
		eventType = "WRITE"
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		eventType = "REMOVE"
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		eventType = "RENAME"
	case event.Op&fsnotify.Chmod == fsnotify.Chmod:
		eventType = "CHMOD"
	default:
		eventType = fmt.Sprintf("UNKNOWN_OP(%s)", event.Op.String())
	}

	logMessage := fmt.Sprintf("[%s] %s: %s", time.Now().Format("2006-01-02 15:04:05"), eventType, event.Name)

	log.Println(logMessage)
	logger.Println(logMessage)
}