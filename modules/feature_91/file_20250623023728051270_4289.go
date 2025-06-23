package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

type ResourceStats struct {
	Timestamp         time.Time
	NumGoroutines     int
	NumCPU            int
	MemAlloc          uint64
	MemTotalAlloc     uint64
	MemSys            uint64
	MemHeapSys        uint64
	MemHeapInuse      uint64
	MemStackInuse     uint64
	MemFrees          uint64
	MemLookups        uint64
	MemGCNum          uint32
	MemGCPauseTotalNs uint64
}

func GetResourceStats() ResourceStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return ResourceStats{
		Timestamp:         time.Now(),
		NumGoroutines:     runtime.NumGoroutine(),
		NumCPU:            runtime.NumCPU(),
		MemAlloc:          m.Alloc,
		MemTotalAlloc:     m.TotalAlloc,
		MemSys:            m.Sys,
		MemHeapSys:        m.HeapSys,
		MemHeapInuse:      m.HeapInuse,
		MemStackInuse:     m.StackInuse,
		MemFrees:          m.Frees,
		MemLookups:        m.Lookups,
		MemGCNum:          m.NumGC,
		MemGCPauseTotalNs: m.PauseTotalNs,
	}
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats := GetResourceStats()

		log.Printf(
			"[%s] Goroutines: %d, CPUs: %d, MemAlloc: %s, MemSys: %s, HeapInuse: %s, GC Cycles: %d, GC Pause Total: %s",
			stats.Timestamp.Format("2006-01-02 15:04:05.000"),
			stats.NumGoroutines,
			stats.NumCPU,
			byteCountToHuman(stats.MemAlloc),
			byteCountToHuman(stats.MemSys),
			byteCountToHuman(stats.MemHeapInuse),
			stats.MemGCNum,
			time.Duration(stats.MemGCPauseTotalNs).String(),
		)
	}
}

func byteCountToHuman(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// Additional implementation at 2025-06-23 02:38:32
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// SystemMetrics holds collected system resource data.
type SystemMetrics struct {
	Timestamp    time.Time            `json:"timestamp"`
	CPUPercent   []float64            `json:"cpu_percent_per_core"` // Per-core usage percentage
	MemTotal     uint64               `json:"mem_total_bytes"`
	MemUsed      uint64               `json:"mem_used_bytes"`
	MemFree      uint64               `json:"mem_free_bytes"`
	MemAvailable uint64               `json:"mem_available_bytes"`
	DiskUsage    *disk.UsageStat      `json:"disk_usage,omitempty"` // Overall disk usage for root partition
	NetIO        []net.IOCountersStat `json:"net_io_counters"`      // Per-interface network I/O
}

// ResourceLogger manages the system resource logging process.
type ResourceLogger struct {
	LogFilePath string
	Interval    time.Duration
	LogFile     *os.File
	Logger      *log.Logger // Logger for writing to the file
}

// NewResourceLogger creates and initializes a new ResourceLogger.
// It opens the specified log file for appending and sets up the internal logger.
func NewResourceLogger(logFilePath string, interval time.Duration) (*ResourceLogger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %q: %w", logFilePath, err)
	}

	return &ResourceLogger{
		LogFilePath: logFilePath,
		Interval:    interval,
		LogFile:     file,
		Logger:      log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds), // Include microseconds for precise timestamps
	}, nil
}

// Close closes the underlying log file.
func (rl *ResourceLogger) Close() error {
	if rl.LogFile != nil {
		return rl.LogFile.Close()
	}
	return nil
}

// collectMetrics gathers system resource usage using gopsutil.
func (rl *ResourceLogger) collectMetrics() (*SystemMetrics, error) {
	metrics := &SystemMetrics{
		Timestamp: time.Now(),
	}

	// CPU usage (per-core)
	cpuPercents, err := cpu.Percent(0, true) // 0 for non-blocking, true for per-CPU
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU percent: %w", err)
	}
	metrics.CPUPercent = cpuPercents

	// Memory usage
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual memory info: %w", err)
	}
	metrics.MemTotal = vmem.Total
	metrics.MemUsed = vmem.Used
	metrics.MemFree = vmem.Free
	metrics.MemAvailable = vmem.Available

	// Disk usage (for the root partition or first available)
	partitions, err := disk.Partitions(false) // Only physical partitions
	if err == nil && len(partitions) > 0 {
		var targetMountpoint string
		for _, p := range partitions {
			if p.Mountpoint == "/" { // Prioritize root partition
				targetMountpoint = p.Mountpoint
				break
			}
		}
		if targetMountpoint == "" { // Fallback to the first partition if root not found
			targetMountpoint = partitions[0].Mountpoint
		}

		diskUsage, err := disk.Usage(targetMountpoint)
		if err != nil {
			// Log a warning but don't fail metric collection
			log.Printf("Warning: Failed to get disk usage for %q: %v", targetMountpoint, err)
		} else {
			metrics.DiskUsage = diskUsage
		}
	} else if err != nil {
		log.Printf("Warning: Failed to get disk partitions: %v", err)
	}

	// Network I/O (per-interface)
	netIO, err := net.IOCounters(true) // true for per-interface
	if err != nil {
		return nil, fmt.Errorf("failed to get network IO counters: %w", err)
	}
	metrics.NetIO = netIO

	return metrics, nil
}

// logMetrics writes the collected metrics to the log file in JSON Lines format.
func (rl *ResourceLogger) logMetrics(metrics *SystemMetrics) {
	data, err := json.Marshal(metrics)
	if err != nil {
		rl.Logger.Printf("Error marshaling metrics to JSON: %v", err)
		return
	}
	rl.Logger.Println(string(data)) // Writes the JSON string followed by a newline
}

// StartLogging begins the continuous logging process.
// It runs in a loop, collecting and logging metrics at the specified interval,
// until the provided context is cancelled.
func (rl *ResourceLogger) StartLogging(ctx context.Context) {
	ticker := time.NewTicker(rl.Interval)
	defer ticker.Stop() // Ensure the ticker is stopped when the function exits

	log.Printf("Resource logging started. Logging to %q every %s.", rl.LogFilePath, rl.Interval)

	for {
		select {
		case <-ticker.C: // Tick event: time to collect and log metrics
			metrics, err := rl.collectMetrics()
			if err != nil {
				log.Printf("Error collecting metrics: %v", err) // Log to stderr
				rl.Logger.Printf("Error collecting metrics: %v", err) // Also log to file
				continue
			}
			rl.logMetrics(metrics)
		case <-ctx.Done(): // Context cancelled: time to shut down
			log.Println("Resource logging stopped by context cancellation.")
			return
		}
	}
}

func main() {
	// Configuration for the logger
	logFilePath := "system_resource_log.jsonl" // Output file in JSON Lines format
	logInterval := 5 * time.Second             // Log every 5 seconds

	// Initialize the resource logger
	logger, err := NewResourceLogger(logFilePath, logInterval)
	if err != nil {
		log.Fatalf("Failed to initialize resource logger: %v", err)
	}
	// Ensure the log file is closed when main exits
	defer func() {
		if err := logger.Close(); err != nil {
			log.Printf("Error closing log file: %v", err)
		}
		log.Println("Log file closed.")
	}()

	// Setup graceful shutdown using context and OS signals
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	// Register to receive SIGINT (Ctrl+C) and SIGTERM signals
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine to listen for OS signals and cancel the context
	go func() {
		sig := <-sigChan // Block until a signal is received
		log.Printf("Received signal: %s. Initiating graceful shutdown...", sig)
		cancel() // Cancel the context, which will stop the logger's loop
	}()

	// Start the resource logger. This call blocks until the context is cancelled.
	logger.StartLogging(ctx)

	log.Println("Application gracefully shut down.")
}