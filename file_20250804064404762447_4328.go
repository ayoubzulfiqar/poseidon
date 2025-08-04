package main

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {
	fmt.Println("Starting mouse movement tracker. Press Ctrl+C to exit.")
	for {
		x, y := robotgo.GetMousePos()
		fmt.Printf("Mouse Position: X=%d, Y=%d\n", x, y)
		time.Sleep(100 * time.Millisecond) // Update every 100ms
	}
}

// Additional implementation at 2025-08-04 06:44:52
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-vgo/robotgo"
)

// MouseTracker holds the state for tracking mouse movements.
type MouseTracker struct {
	lastX, lastY int
	lastUpdateTime time.Time
	totalDistance float64
	startTime time.Time
	logFile *os.File
	csvWriter *csv.Writer
	pollingInterval time.Duration
	stopChan chan struct{}
}

// NewMouseTracker creates and initializes a new MouseTracker.
func NewMouseTracker(logFileName string, interval time.Duration) (*MouseTracker, error) {
	file, err := os.Create(logFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	writer := csv.NewWriter(file)
	// Write CSV header
	if err := writer.Write([]string{"Timestamp", "X", "Y", "InstantaneousSpeedPxPerSec", "TotalDistancePx", "ElapsedTimeSec", "AverageSpeedPxPerSec"}); err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}
	writer.Flush()

	x, y := robotgo.GetMousePos()
	currentTime := time.Now()

	return &MouseTracker{
		lastX: x,
		lastY: y,
		lastUpdateTime: currentTime,
		totalDistance: 0.0,
		startTime: currentTime,
		logFile: file,
		csvWriter: writer,
		pollingInterval: interval,
		stopChan: make(chan struct{}),
	}, nil
}

// StartTracking begins polling the mouse position and logging data.
func (mt *MouseTracker) StartTracking() {
	log.Println("Mouse tracking started. Press Ctrl+C to stop.")

	ticker := time.NewTicker(mt.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mt.updateAndLog()
		case <-mt.stopChan:
			log.Println("Stopping mouse tracking.")
			return
		}
	}
}

// StopTracking signals the tracker to stop.
func (mt *MouseTracker) StopTracking() {
	close(mt.stopChan)
}

// Close cleans up resources, like closing the log file.
func (mt *MouseTracker) Close() {
	if mt.logFile != nil {
		mt.csvWriter.Flush()
		if err := mt.csvWriter.Error(); err != nil {
			log.Printf("Error flushing CSV writer: %v", err)
		}
		if err := mt.logFile.Close(); err != nil {
			log.Printf("Error closing log file: %v", err)
		}
	}
	log.Println("Tracker resources closed.")
}

// updateAndLog gets the current mouse position, calculates distance/speed, and logs it.
func (mt *MouseTracker) updateAndLog() {
	currentX, currentY := robotgo.GetMousePos()
	currentTime := time.Now()

	timeDelta := currentTime.Sub(mt.lastUpdateTime).Seconds()
	if timeDelta == 0 { // Avoid division by zero if updates are too fast or system clock resolution is low
		timeDelta = 0.001 // Small non-zero value
	}

	dx := float64(currentX - mt.lastX)
	dy := float64(currentY - mt.lastY)
	distanceMoved := math.Sqrt(dx*dx + dy*dy)

	mt.totalDistance += distanceMoved

	instantaneousSpeed := distanceMoved / timeDelta // Pixels per second

	mt.lastX = currentX
	mt.lastY = currentY
	mt.lastUpdateTime = currentTime // Update last update time

	totalElapsedTime := time.Since(mt.startTime).Seconds()
	averageSpeed := mt.totalDistance / totalElapsedTime
	if totalElapsedTime == 0 { // Avoid division by zero for average speed at very start
		averageSpeed = 0
	}

	// Log to CSV
	record := []string{
		currentTime.Format(time.RFC3339Nano),
		strconv.Itoa(currentX),
		strconv.Itoa(currentY),
		fmt.Sprintf("%.2f", instantaneousSpeed),
		fmt.Sprintf("%.2f", mt.totalDistance),
		fmt.Sprintf("%.2f", totalElapsedTime),
		fmt.Sprintf("%.2f", averageSpeed),
	}
	if err := mt.csvWriter.Write(record); err != nil {
		log.Printf("Error writing record to CSV: %v", err)
	}
	mt.csvWriter.Flush()
}

func main() {
	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logFileName := "mouse_movements.csv"
	pollingInterval := 100 * time.Millisecond // Check mouse position every 100ms

	tracker, err := NewMouseTracker(logFileName, pollingInterval)
	if err != nil {
		log.Fatalf("Failed to create mouse tracker: %v", err)
	}
	defer tracker.Close() // Ensure resources are closed on exit

	go tracker.StartTracking() // Run tracking in a goroutine

	// Wait for termination signal
	<-sigChan
	tracker.StopTracking() // Signal the tracker to stop
	// Give the goroutine a moment to process the stop signal and flush
	time.Sleep(200 * time.Millisecond)
	log.Println("Application terminated.")
}