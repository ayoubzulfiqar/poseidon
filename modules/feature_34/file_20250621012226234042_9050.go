package main

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

const (
	idleTimeout   = 5 * time.Second // Duration after which to lock workstation
	checkInterval = 1 * time.Second // How often to check for idle time
)

// LASTINPUTINFO struct for GetLastInputInfo
type LASTINPUTINFO struct {
	CbSize uint32
	DwTime uint32 // Tick count of the last input event
}

// Windows API calls
var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getLastInputInfo = user32.NewProc("GetLastInputInfo")
	lockWorkStation  = user32.NewProc("LockWorkStation")
)

// getIdleTime retrieves the idle time on Windows.
func getIdleTime() (time.Duration, error) {
	var lastInputInfo LASTINPUTINFO
	lastInputInfo.CbSize = uint32(unsafe.Sizeof(lastInputInfo))

	ret, _, err := getLastInputInfo.Call(uintptr(unsafe.Pointer(&lastInputInfo)))
	if ret == 0 {
		return 0, fmt.Errorf("GetLastInputInfo failed: %v", err)
	}

	currentTickCount := uint32(syscall.GetTickCount())
	idleMilliseconds := currentTickCount - lastInputInfo.DwTime

	return time.Duration(idleMilliseconds) * time.Millisecond, nil
}

// lockWorkstation locks the Windows workstation.
func lockWorkstation() error {
	ret, _, err := lockWorkStation.Call()
	if ret == 0 {
		return fmt.Errorf("LockWorkStation failed: %v", err)
	}
	return nil
}

func main() {
	if runtime.GOOS != "windows" {
		fmt.Println("This program is designed for Windows only due to OS-specific API calls.")
		os.Exit(1)
	}

	fmt.Printf("Workstation locker started. Will lock after %s of idle time.\n", idleTimeout)
	fmt.Printf("Checking every %s...\n", checkInterval)

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		idle, err := getIdleTime()
		if err != nil {
			fmt.Printf("Error getting idle time: %v\n", err)
			continue
		}

		fmt.Printf("Current idle time: %s\n", idle)

		if idle >= idleTimeout {
			fmt.Printf("Idle time %s exceeded %s. Locking workstation...\n", idle, idleTimeout)
			err := lockWorkstation()
			if err != nil {
				fmt.Printf("Error locking workstation: %v\n", err)
			} else {
				fmt.Println("Workstation locked successfully.")
				os.Exit(0)
			}
		}
	}
}

// Additional implementation at 2025-06-21 01:23:12
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"
)

const (
	// Configuration
	idleThresholdMinutes   = 5  // Lock after 5 minutes of idle time
	warningDurationSeconds = 10 // Show warning for 10 seconds before locking
	checkIntervalSeconds   = 5  // Check idle time every 5 seconds

	// Windows API constants and structs
	lastInputInfoSize = 8 // Size of LASTINPUTINFO struct (cbSize field)
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procGetLastInputInfo = user32.NewProc("GetLastInputInfo")
	procLockWorkStation  = user32.NewProc("LockWorkStation")
	procGetTickCount     = kernel32.NewProc("GetTickCount")
)

// LASTINPUTINFO struct for GetLastInputInfo
type LASTINPUTINFO struct {
	CbSize uint32
	DwTime uint32
}

// getTickCount retrieves the number of milliseconds that have elapsed since the system started.
func getTickCount() uint32 {
	ret, _, _ := procGetTickCount.Call()
	return uint32(ret)
}

// getIdleTime retrieves the idle time in milliseconds.
func getIdleTime() (time.Duration, error) {
	var lastInputInfo LASTINPUTINFO
	lastInputInfo.CbSize = lastInputInfoSize // Set the size of the struct

	// Call GetLastInputInfo
	ret, _, err := procGetLastInputInfo.Call(uintptr(unsafe.Pointer(&lastInputInfo)))
	if ret == 0 {
		return 0, fmt.Errorf("GetLastInputInfo failed: %v", err)
	}

	// Calculate idle time
	currentTickCount := getTickCount()
	idleTimeMs := currentTickCount - lastInputInfo.DwTime
	return time.Duration(idleTimeMs) * time.Millisecond, nil
}

// lockWorkStation locks the workstation.
func lockWorkStation() error {
	ret, _, err := procLockWorkStation.Call()
	if ret == 0 {
		return fmt.Errorf("LockWorkStation failed: %v", err)
	}
	return nil
}

func main() {
	fmt.Printf("Workstation Locker started.\n")
	fmt.Printf("Idle threshold: %d minutes\n", idleThresholdMinutes)
	fmt.Printf("Warning duration: %d seconds\n", warningDurationSeconds)
	fmt.Printf("Check interval: %d seconds\n", checkIntervalSeconds)
	fmt.Printf("Press Ctrl+C to exit.\n\n")

	idleThreshold := time.Duration(idleThresholdMinutes) * time.Minute
	warningDuration := time.Duration(warningDurationSeconds) * time.Second
	checkInterval := time.Duration(checkIntervalSeconds) * time.Second

	// Setup signal handling for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nExiting Workstation Locker.")
		os.Exit(0)
	}()

	for {
		idleTime, err := getIdleTime()
		if err != nil {
			fmt.Printf("Error getting idle time: %v\n", err)
			time.Sleep(checkInterval)
			continue
		}

		fmt.Printf("Current idle time: %s\n", idleTime.Round(time.Second))

		if idleTime >= idleThreshold {
			fmt.Printf("Workstation idle for %s. Locking in %s...\n", idleTime.Round(time.Second), warningDuration.Round(time.Second))
			time.Sleep(warningDuration) // Wait for the warning duration

			// Re-check idle time after warning period, in case user moved mouse
			idleTimeAfterWarning, err := getIdleTime()
			if err != nil {
				fmt.Printf("Error re-checking idle time: %v\n", err)
				time.Sleep(checkInterval)
				continue
			}

			if idleTimeAfterWarning >= idleThreshold {
				fmt.Println("Still idle. Locking workstation now!")
				if err := lockWorkStation(); err != nil {
					fmt.Printf("Failed to lock workstation: %v\n", err)
				} else {
					fmt.Println("Workstation locked successfully.")
					// After locking, sleep for a longer period to prevent immediate re-lock attempts
					// or to allow the system to settle after locking.
					time.Sleep(time.Minute) // Sleep for a minute after locking
				}
			} else {
				fmt.Printf("User activity detected during warning period. Not locking.\n")
			}
		}

		time.Sleep(checkInterval)
	}
}