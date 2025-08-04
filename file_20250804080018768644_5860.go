package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	openPorts := make(chan int, 1024)
	results := make([]int, 0)

	fmt.Println("Scanning ports on localhost (1-1024)...")

	for i := 1; i <= 1024; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			address := fmt.Sprintf("localhost:%d", port)
			conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
			if err != nil {
				return
			}
			conn.Close()
			openPorts <- port
		}(i)
	}

	go func() {
		wg.Wait()
		close(openPorts)
	}()

	for port := range openPorts {
		results = append(results, port)
	}

	sort.Ints(results)

	if len(results) == 0 {
		fmt.Println("No open ports found in the range 1-1024.")
	} else {
		fmt.Println("Open ports found:")
		for _, port := range results {
			fmt.Printf("Port %d is open\n", port)
		}
	}
}

// Additional implementation at 2025-08-04 08:01:22
package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

const (
	targetHost       = "127.0.0.1"
	startPort        = 1
	endPort          = 1024
	timeout          = 500 * time.Millisecond
	concurrencyLimit = 100
)

func main() {
	fmt.Printf("Scanning %s from port %d to %d...\n", targetHost, startPort, endPort)

	var openPorts []int
	var mu sync.Mutex
	var wg sync.WaitGroup

	guard := make(chan struct{}, concurrencyLimit)

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		guard <- struct{}{}

		go func(p int) {
			defer wg.Done()
			defer func() { <-guard }()

			address := fmt.Sprintf("%s:%d", targetHost, p)
			conn, err := net.DialTimeout("tcp", address, timeout)
			if err != nil {
				return
			}
			conn.Close()
			mu.Lock()
			openPorts = append(openPorts, p)
			mu.Unlock()
		}(port)
	}

	wg.Wait()
	close(guard)

	sort.Ints(openPorts)

	if len(openPorts) > 0 {
		fmt.Println("\nOpen ports found:")
		for _, p := range openPorts {
			fmt.Printf("  Port %d is open\n", p)
		}
	} else {
		fmt.Println("\nNo open ports found in the specified range.")
	}
}

// Additional implementation at 2025-08-04 08:01:50
package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	defaultStartPort  = 1
	defaultEndPort    = 65535
	defaultWorkers    = 100
	connectionTimeout = 500 * time.Millisecond // Timeout for TCP connection attempts
	bannerReadTimeout = 200 * time.Millisecond // Timeout for reading service banner
	bannerReadSize    = 1024                   // Max bytes to read for banner
)

// PortResult holds the outcome of scanning a single port.
type PortResult struct {
	Port    int
	IsOpen  bool
	Service string
	Error   error
}

// getServiceBanner attempts to read a small banner from an open TCP connection.
// It sets a read deadline and tries to clean up the banner for display.
func getServiceBanner(conn net.Conn) string {
	conn.SetReadDeadline(time.Now().Add(bannerReadTimeout))
	buffer := make([]byte, bannerReadSize)
	n, err := conn.Read(buffer)
	if err != nil {
		// Check for timeout errors specifically
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return "timeout reading banner"
		}
		return fmt.Sprintf("error reading banner: %v", err)
	}

	// Basic cleanup: replace non-printable ASCII with spaces and limit length
	banner := string(buffer[:n])
	cleanedBanner := ""
	for _, r := range banner {
		if r >= 32 && r <= 126 { // Printable ASCII characters
			cleanedBanner += string(r)
		} else if r == '\n' || r == '\r' || r == '\t' { // Common whitespace
			cleanedBanner += " "
		}
	}
	if len(cleanedBanner) > 80 { // Truncate long banners
		cleanedBanner = cleanedBanner[:80] + "..."
	}
	return cleanedBanner
}

// worker is a goroutine that continuously pulls ports from the jobs channel,
// scans them, and sends the results to the results channel.
func worker(target string, jobs <-chan int, results chan<- PortResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for port := range jobs {
		addr := fmt.Sprintf("%s:%d", target, port)
		conn, err := net.DialTimeout("tcp", addr, connectionTimeout)
		if err != nil {
			results <- PortResult{Port: port, IsOpen: false, Error: err}
			continue
		}
		defer conn.Close()

		serviceBanner := getServiceBanner(conn)
		results <- PortResult{Port: port, IsOpen: true, Service: serviceBanner}
	}
}

func main() {
	target := "127.0.0.1" // Fixed to localhost as per instructions
	startPort := defaultStartPort
	endPort := defaultEndPort
	numWorkers := defaultWorkers

	// Parse command-line arguments for port range and worker count
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-s", "--start":
			if i+1 < len(args) {
				p, err := strconv.Atoi(args[i+1])
				if err == nil && p >= 1 && p <= 65535 {
					startPort = p
				} else {
					fmt.Printf("Warning: Invalid start port '%s'. Using default %d.\n", args[i+1], defaultStartPort)
				}
				i++
			}
		case "-e", "--end":
			if i+1 < len(args) {
				p, err := strconv.Atoi(args[i+1])
				if err == nil && p >= 1 && p <= 65535 {
					endPort = p
				} else {
					fmt.Printf("Warning: Invalid end port '%s'. Using default %d.\n", args[i+1], defaultEndPort)
				}
				i++
			}
		case "-w", "--workers":
			if i+1 < len(args) {
				w, err := strconv.Atoi(args[i+1])
				if err == nil && w > 0 {
					numWorkers = w
				} else {
					fmt.Printf("Warning: Invalid worker count '%s'. Using default %d.\n", args[i+1], defaultWorkers)
				}
				i++
			}
		case "-h", "--help":
			fmt.Println("Usage: go run scanner.go [OPTIONS]")
			fmt.Println("Options:")
			fmt.Println("  -s, --start <port>   Start port (default: 1, min: 1, max: 65535)")
			fmt.Println("  -e, --end <port>     End port (default: 65535, min: 1, max: 65535)")
			fmt.Println("  -w, --workers <num>  Number of concurrent workers (default: 100, min: 1)")
			fmt.Println("  -h, --help           Show this help message")
			os.Exit(0)
		default:
			fmt.Printf("Unknown argument: %s. Use -h for help.\n", args[i])
			os.Exit(1)
		}
	}

	if startPort > endPort {
		fmt.Println("Error: Start port cannot be greater than end port.")
		os.Exit(1)
	}

	fmt.Printf("Scanning %s from port %d to %d with %d workers...\n", target, startPort, endPort, numWorkers)

	jobs := make(chan int, numWorkers) // Buffered channel for ports to scan
	results := make(chan PortResult, (endPort-startPort+1)) // Buffered channel for results
	var wg sync.WaitGroup // WaitGroup to wait for all workers to complete

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(target, jobs, results, &wg)
	}

	// Populate the jobs channel with all ports to be scanned
	for p := startPort; p <= endPort; p++ {
		jobs <- p
	}
	close(jobs) // Close the jobs channel to signal workers no more jobs are coming

	// Wait for all workers to finish processing their jobs
	wg.Wait()
	close(results) // Close the results channel after all workers have sent their results

	// Collect and print results
	foundOpenPorts := false
	fmt.Println("\n--- Scan Results ---")
	for res := range results {
		if res.IsOpen {
			foundOpenPorts = true
			if res.Service != "" {
				fmt.Printf("Port %d is OPEN (Service: %s)\n", res.Port, res.Service)
			} else {
				fmt.Printf("Port %d is OPEN\n", res.Port)
			}
		}
	}

	if !foundOpenPorts {
		fmt.Println("No open ports found in the specified range.")
	}
	fmt.Println("--- Scan Complete ---")
}

// Additional implementation at 2025-08-04 08:02:40
package main

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	targetHost       = "127.0.0.1"
	startPort        = 1
	endPort          = 1024 // Common ports range for quick localhost scan
	concurrencyLimit = 100  // Number of concurrent goroutines
	timeoutSeconds   = 1    // Timeout for each connection attempt in seconds
)

// PortScanResult holds the result for a single port scan
type PortScanResult struct {
	Port   int
	IsOpen bool
	Error  error // Stores the error if the port is not open
}

// scanPort attempts to connect to a specific port and reports the result
func scanPort(port int, results chan<- PortScanResult, wg *sync.WaitGroup) {
	defer wg.Done()

	address := net.JoinHostPort(targetHost, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeoutSeconds)*time.Second)

	if err != nil {
		results <- PortScanResult{Port: port, IsOpen: false, Error: err}
		return
	}
	defer conn.Close()

	results <- PortScanResult{Port: port, IsOpen: true, Error: nil}
}

func main() {
	fmt.Printf("Scanning ports %d-%d on %s...\n", startPort, endPort, targetHost)

	var wg sync.WaitGroup
	// Buffered channel to hold all results, size equals the total number of ports to scan
	results := make(chan PortScanResult, endPort-startPort+1)

	// Semaphore to limit the number of concurrent goroutines
	sem := make(chan struct{}, concurrencyLimit)

	for i := startPort; i <= endPort; i++ {
		sem <- struct{}{} // Acquire a slot in the semaphore (blocks if limit reached)
		wg.Add(1)
		go func(port int) {
			defer func() { <-sem }() // Release the slot when the goroutine finishes
			scanPort(port, results, &wg)
		}(i)
	}

	wg.Wait()    // Wait for all goroutines to finish their execution
	close(results) // Close the results channel as no more results will be sent

	var openPorts []int
	var closedPorts int
	var filteredPorts int // Ports that timed out or connection refused

	for res := range results {
		if res.IsOpen {
			openPorts = append(openPorts, res.Port)
		} else {
			// Differentiate between truly closed (connection refused) and filtered (timeout)
			if netErr, ok := res.Error.(net.Error); ok && netErr.Timeout() {
				filteredPorts++
			} else if res.Error != nil && (res.Error.Error() == "connection refused" || res.Error.Error() == "connect: connection refused") {
				closedPorts++
			} else {
				// Other errors, treat as filtered for simplicity in this context
				filteredPorts++
			}
		}
	}

	sort.Ints(openPorts) // Sort open ports for cleaner, ordered output

	fmt.Println("\n--- Scan Results ---")
	if len(openPorts) > 0 {
		fmt.Println("Open Ports:")
		for _, p := range openPorts {
			fmt.Printf("  Port %d is OPEN\n", p)
		}
	} else {
		fmt.Println("No open ports found in the specified range.")
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total Ports Scanned: %d\n", endPort-startPort+1)
	fmt.Printf("  Open Ports: %d\n", len(openPorts))
	fmt.Printf("  Closed Ports: %d\n", closedPorts)
	fmt.Printf("  Filtered/Other Errors: %d\n", filteredPorts)
	fmt.Println("Scan complete.")
}