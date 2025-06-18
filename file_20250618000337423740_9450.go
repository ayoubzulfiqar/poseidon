package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync/atomic"

	"github.com/eiannone/keyboard"
)

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {
	var keystrokeCount int64 = 0

	err := keyboard.Open()
	if err != nil {
		fmt.Println("Error opening keyboard:", err)
		return
	}
	defer keyboard.Close()

	clearScreen()
	fmt.Println("Keystroke Counter")
	fmt.Println("Press any key to count, press 'Esc' to exit.")
	fmt.Printf("Keystrokes: %d\n", atomic.LoadInt64(&keystrokeCount))

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Error reading key:", err)
			break
		}

		if key == keyboard.KeyEsc {
			break
		}

		atomic.AddInt64(&keystrokeCount, 1)

		clearScreen()
		fmt.Println("Keystroke Counter")
		fmt.Println("Press any key to count, press 'Esc' to exit.")
		fmt.Printf("Keystrokes: %d\n", atomic.LoadInt64(&keystrokeCount))
	}

	clearScreen()
	fmt.Println("Exiting Keystroke Counter.")
	fmt.Printf("Total Keystrokes: %d\n", atomic.LoadInt64(&keystrokeCount))
}

// Additional implementation at 2025-06-18 00:04:20
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

const keystrokeFile = "keystrokes.txt"

var (
	totalKeystrokes       int
	currentMinuteKeystrokes int
	mu                    sync.Mutex
	exitChan              chan struct{}
	kpmTicker             *time.Ticker
	inputReader           *bufio.Reader
)

func loadCount() {
	data, err := os.ReadFile(keystrokeFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		return
	}
	count, err := strconv.Atoi(string(data))
	if err != nil {
		return
	}
	mu.Lock()
	totalKeystrokes = count
	mu.Unlock()
}

func saveCount() {
	mu.Lock()
	count := totalKeystrokes
	mu.Unlock()
	err := os.WriteFile(keystrokeFile, []byte(strconv.Itoa(count)), 0644)
	if err != nil {
		return
	}
}

func readInput() {
	for {
		r, _, err := inputReader.ReadRune()
		if err != nil {
			select {
			case exitChan <- struct{}{}:
			default:
			}
			return
		}

		if r == 'q' || r == 'Q' {
			select {
			case exitChan <- struct{}{}:
			default:
			}
			return
		}

		mu.Lock()
		totalKeystrokes++
		currentMinuteKeystrokes++
		mu.Unlock()
	}
}

func updateStats() {
	minuteResetTicker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-kpmTicker.C:
			mu.Lock()
			currentKPM := currentMinuteKeystrokes
			total := totalKeystrokes
			mu.Unlock()
			fmt.Printf("\rTotal Keystrokes: %d | Current Minute Keystrokes: %d ", total, currentKPM)
		case <-minuteResetTicker.C:
			mu.Lock()
			currentMinuteKeystrokes = 0
			mu.Unlock()
		case <-exitChan:
			kpmTicker.Stop()
			minuteResetTicker.Stop()
			return
		}
	}
}

func main() {
	exitChan = make(chan struct{})
	inputReader = bufio.NewReader(os.Stdin)
	kpmTicker = time.NewTicker(1 * time.Second)

	loadCount()
	fmt.Printf("Loaded %d previous keystrokes.\n", totalKeystrokes)
	fmt.Println("Type something and press Enter. Press 'q' then Enter to quit.")

	go readInput()
	go updateStats()

	<-exitChan

	fmt.Println("\nExiting keystroke counter.")
	saveCount()
}

// Additional implementation at 2025-06-18 00:05:03
package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"sort"
	"time"

	"golang.org/x/term"
)

const dataFile = "keystrokes.gob"
const displayInterval = 2 * time.Second

type KeystrokeData struct {
	TotalKeystrokes   int
	KeyCounts         map[rune]int
	StartTime         time.Time
	SessionKeystrokes int
}

var currentData KeystrokeData

func loadData() {
	file, err := os.Open(dataFile)
	if err != nil {
		currentData = KeystrokeData{
			KeyCounts: make(map[rune]int),
			StartTime: time.Now(),
		}
		return
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&currentData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding data: %v\n", err)
		currentData = KeystrokeData{
			KeyCounts: make(map[rune]int),
			StartTime: time.Now(),
		}
		return
	}
	currentData.StartTime = time.Now()
	currentData.SessionKeystrokes = 0
}

func saveData() {
	file, err := os.Create(dataFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating data file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(currentData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding data: %v\n", err)
	}
}

func displayStats() {
	fmt.Print("\033[H\033[2J")
	fmt.Println("Keystroke Counter (Press 'q' to quit)")
	fmt.Println("-------------------------------------")
	fmt.Printf("Total Keystrokes (All Time): %d\n", currentData.TotalKeystrokes)

	elapsed := time.Since(currentData.StartTime)
	if elapsed > 0 {
		kpm := float64(currentData.SessionKeystrokes) / elapsed.Minutes()
		fmt.Printf("Keystrokes Per Minute (Session): %.2f\n", kpm)
	} else {
		fmt.Println("Keystrokes Per Minute (Session): 0.00")
	}

	fmt.Println("\nTop 5 Most Pressed Keys:")
	type keyCount struct {
		Rune  rune
		Count int
	}
	var counts []keyCount
	for r, c := range currentData.KeyCounts {
		counts = append(counts, keyCount{Rune: r, Count: c})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	for i, kc := range counts {
		if i >= 5 {
			break
		}
		keyStr := string(kc.Rune)
		if kc.Rune == ' ' {
			keyStr = "[Space]"
		} else if kc.Rune == '\r' || kc.Rune == '\n' {
			keyStr = "[Enter]"
		} else if kc.Rune == '\t' {
			keyStr = "[Tab]"
		} else if kc.Rune < 32 || kc.Rune > 126 {
			keyStr = fmt.Sprintf("[%X]", kc.Rune)
		}
		fmt.Printf("  '%s': %d\n", keyStr, kc.Count)
	}
	fmt.Println("-------------------------------------")
}

func main() {
	loadData()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set raw terminal mode: %v\n", err)
		os.Exit(1)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	inputReader := bufio.NewReader(os.Stdin)
	stopChan := make(chan struct{})
	displayTicker := time.NewTicker(displayInterval)
	defer displayTicker.Stop()

	go func() {
		for {
			select {
			case <-displayTicker.C:
				displayStats()
			case <-stopChan:
				return
			}
		}
	}()

	displayStats()

	for {
		r, _, err := inputReader.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			break
		}

		if r == 'q' || r == 3 {
			break
		}

		currentData.TotalKeystrokes++
		currentData.SessionKeystrokes++
		currentData.KeyCounts[r]++
	}

	close(stopChan)
	time.Sleep(100 * time.Millisecond)

	saveData()
	displayStats()
	fmt.Println("\nExiting keystroke counter.")
}

// Additional implementation at 2025-06-18 00:06:17
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

const dataFileName = "keystrokes_data.json"

// KeystrokeData holds the total keystroke count for persistence
type KeystrokeData struct {
	TotalKeystrokes int `json:"total_keystrokes"`
}

var (
	totalKeystrokes   int
	sessionKeystrokes int
	mu                sync.Mutex // Mutex to protect shared counters
	exitChan          chan struct{}
)

func loadKeystrokes() int {
	data, err := ioutil.ReadFile(dataFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return 0 // File doesn't exist, start from 0
		}
		log.Printf("Error reading keystroke data file: %v", err)
		return 0
	}

	var kd KeystrokeData
	err = json.Unmarshal(data, &kd)
	if err != nil {
		log.Printf("Error unmarshaling keystroke data: %v", err)
		return 0
	}
	return kd.TotalKeystrokes
}

func saveKeystrokes() {
	mu.Lock()
	defer mu.Unlock()

	kd := KeystrokeData{TotalKeystrokes: totalKeystrokes}
	data, err := json.MarshalIndent(kd, "", "  ")
	if err != nil {
		log.Printf("Error marshaling keystroke data: %v", err)
		return
	}

	err = ioutil.WriteFile(dataFileName, data, 0644)
	if err != nil {
		log.Printf("Error writing keystroke data file: %v", err)
	}
}

func draw() {
	mu.Lock()
	defer mu.Unlock()

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Print instructions
	printString(0, 0, termbox.ColorWhite, termbox.ColorDefault, "Keystroke Counter")
	printString(0, 2, termbox.ColorWhite, termbox.ColorDefault, "-----------------")
	printString(0, 4, termbox.ColorGreen, termbox.ColorDefault, fmt.Sprintf("Total Keystrokes: %d", totalKeystrokes))
	printString(0, 5, termbox.ColorCyan, termbox.ColorDefault, fmt.Sprintf("Session Keystrokes: %d", sessionKeystrokes))
	printString(0, 7, termbox.ColorYellow, termbox.ColorDefault, "Press 'r' to reset session count.")
	printString(0, 8, termbox.ColorYellow, termbox.ColorDefault, "Press 'q' or Ctrl+C to quit.")

	termbox.Flush()
}

func printString(x, y int, fg, bg termbox.Attribute, msg string) {
	for i, r := range msg {
		termbox.SetCell(x+i, y, r, fg, bg)
	}
}

func eventLoop() {
	for {
		select {
		case <-exitChan: // Check if exit signal received
			return
		default:
			ev := termbox.PollEvent()
			mu.Lock() // Lock before modifying shared counters
			switch ev.Type {
			case termbox.EventKey:
				if ev.Key == termbox.KeyCtrlC || ev.Ch == 'q' {
					close(exitChan) // Signal main to exit
					mu.Unlock()
					return
				} else if ev.Ch == 'r' {
					sessionKeystrokes = 0
				} else {
					totalKeystrokes++
					sessionKeystrokes++
				}
				draw() // Redraw after key press
			case termbox.EventResize:
				draw() // Redraw on window resize
			case termbox.EventError:
				log.Printf("Termbox error: %v", ev.Err)
				close(exitChan)
				mu.Unlock()
				return
			}
			mu.Unlock() // Unlock after modifying counters and drawing
		}
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatalf("Failed to initialize termbox: %v", err)
	}
	defer termbox.Close()
	defer saveKeystrokes() // Ensure data is saved on exit

	exitChan = make(chan struct{})
	totalKeystrokes = loadKeystrokes()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputAlt)
	termbox.SetOutputMode(termbox.Output256) // Enable more colors if available

	draw() // Initial draw

	go eventLoop() // Start event processing in a goroutine

	<-exitChan // Block main goroutine until exit signal is received
	fmt.Println("\nExiting keystroke counter. Total keystrokes saved.")
}