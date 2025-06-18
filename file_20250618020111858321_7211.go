package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gen2brain/beeep"
)

var (
	workDuration       = flag.Int("work", 25, "Duration of a work session in minutes")
	shortBreakDuration = flag.Int("short-break", 5, "Duration of a short break in minutes")
	longBreakDuration  = flag.Int("long-break", 15, "Duration of a long break in minutes")
	cycles             = flag.Int("cycles", 4, "Number of work sessions before a long break")
	silent             = flag.Bool("silent", false, "Disable notifications")
)

func main() {
	flag.Parse()

	fmt.Println("Starting Pomodoro Timer!")
	fmt.Printf("Work: %d min, Short Break: %d min, Long Break: %d min, Cycles: %d\n",
		*workDuration, *shortBreakDuration, *longBreakDuration, *cycles)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nPomodoro timer stopped.")
		os.Exit(0)
	}()

	for i := 1; ; i++ {
		fmt.Printf("\n--- Cycle %d ---\n", i)

		fmt.Println("Starting Work Session...")
		runTimer(*workDuration*time.Minute, "Work Session")
		if !*silent {
			sendNotification("Pomodoro", "Work session finished! Take a break.")
		}

		if i%*cycles == 0 {
			fmt.Println("Starting Long Break...")
			runTimer(*longBreakDuration*time.Minute, "Long Break")
			if !*silent {
				sendNotification("Pomodoro", "Long break finished! Time to work.")
			}
		} else {
			fmt.Println("Starting Short Break...")
			runTimer(*shortBreakDuration*time.Minute, "Short Break")
			if !*silent {
				sendNotification("Pomodoro", "Short break finished! Time to work.")
			}
		}
	}
}

func runTimer(duration time.Duration, sessionType string) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	endTime := time.Now().Add(duration)

	for range ticker.C {
		remaining := endTime.Sub(time.Now())
		if remaining <= 0 {
			printTime(0, sessionType)
			fmt.Println()
			break
		}
		printTime(remaining, sessionType)
	}
}

func printTime(d time.Duration, sessionType string) {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	fmt.Printf("\r%s: %02d:%02d", sessionType, minutes, seconds)
}

func sendNotification(title, message string) {
	err := beeep.Notify(title, message, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending notification: %v\n", err)
	}
}

// Additional implementation at 2025-06-18 02:01:46


// Additional implementation at 2025-06-18 02:02:51
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
)

type PomodoroState int

const (
	StateIdle PomodoroState = iota
	StateWorking
	StateShortBreak
	StateLongBreak
	StatePaused
)

func (s PomodoroState) String() string {
	switch s {
	case StateIdle:
		return "Idle"
	case StateWorking:
		return "Working"
	case StateShortBreak:
		return "Short Break"
	case StateLongBreak:
		return "Long Break"
	case StatePaused:
		return "Paused"
	default:
		return "Unknown"
	}
}

type PomodoroConfig struct {
	WorkDuration             time.Duration
	ShortBreakDuration       time.Duration
	LongBreakDuration        time.Duration
	PomodorosBeforeLongBreak int
}

type PomodoroTimer struct {
	Config PomodoroConfig

	state        PomodoroState
	currentPhase time.Duration
	remaining    time.Duration

	pomodoroCount  int
	totalPomodoros int

	timer      *time.Timer
	stopChan   chan struct{}
	pauseChan  chan struct{}
	resumeChan chan struct{}
	skipChan   chan struct{}

	mu sync.Mutex
}

func NewPomodoroTimer(config PomodoroConfig) *PomodoroTimer {
	if config.WorkDuration == 0 {
		config.WorkDuration = 25 * time.Minute
	}
	if config.ShortBreakDuration == 0 {
		config.ShortBreakDuration = 5 * time.Minute
	}
	if config.LongBreakDuration == 0 {
		config.LongBreakDuration = 15 * time.Minute
	}
	if config.PomodorosBeforeLongBreak == 0 {
		config.PomodorosBeforeLongBreak = 4
	}

	return &PomodoroTimer{
		Config:       config,
		state:        StateIdle,
		stopChan:     make(chan struct{}),
		pauseChan:    make(chan struct{}),
		resumeChan:   make(chan struct{}),
		skipChan:     make(chan struct{}),
		mu:           sync.Mutex{},
	}
}

func (pt *PomodoroTimer) Start() {
	pt.mu.Lock()
	if pt.state != StateIdle {
		pt.mu.Unlock()
		fmt.Println("Timer is already running or paused. Use 'p' to pause, 's' to stop, 'k' to skip.")
		return
	}
	pt.state = StateWorking
	pt.currentPhase = pt.Config.WorkDuration
	pt.remaining = pt.Config.WorkDuration
	pt.mu.Unlock()

	pt.clearScreen()
	fmt.Printf("Starting Pomodoro: %s for %s\n", pt.state, pt.currentPhase)
	pt.notify("Pomodoro Started", fmt.Sprintf("Time to focus for %s!", pt.currentPhase))
	go pt.runTimer()
}

func (pt *PomodoroTimer) runTimer() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		pt.mu.Lock()
		select {
		case <-pt.stopChan:
			pt.state = StateIdle
			pt.remaining = 0
			pt.clearScreen()
			fmt.Println("Pomodoro timer stopped.")
			pt.mu.Unlock()
			return
		case <-pt.pauseChan:
			pt.state = StatePaused
			pt.clearScreen()
			fmt.Printf("Timer Paused. Remaining: %s\n", pt.remaining)
			pt.mu.Unlock()
			select {
			case <-pt.resumeChan:
				pt.mu.Lock()
				if pt.currentPhase == pt.Config.WorkDuration {
					pt.state = StateWorking
				} else if pt.currentPhase == pt.Config.ShortBreakDuration {
					pt.state = StateShortBreak
				} else if pt.currentPhase == pt.Config.LongBreakDuration {
					pt.state = StateLongBreak
				}
				pt.clearScreen()
				fmt.Printf("Timer Resumed. %s: %s\n", pt.state, pt.remaining)
				pt.mu.Unlock()
			case <-pt.stopChan:
				pt.mu.Lock()
				pt.state = StateIdle
				pt.remaining = 0
				pt.clearScreen()
				fmt.Println("Pomodoro timer stopped.")
				pt.mu.Unlock()
				return
			case <-pt.skipChan:
				pt.mu.Unlock()
				pt.nextPhase()
				continue
			}
		case <-pt.skipChan:
			pt.mu.Unlock()
			pt.nextPhase()
			continue
		case <-ticker.C:
			if pt.state == StatePaused {
				pt.mu.Unlock()
				continue
			}

			pt.remaining -= 1 * time.Second
			pt.displayTimer()

			if pt.remaining <= 0 {
				pt.mu.Unlock()
				pt.nextPhase()
				continue
			}
		}
		pt.mu.Unlock()
	}
}

func (pt *PomodoroTimer) nextPhase() {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	switch pt.state {
	case StateWorking:
		pt.pomodoroCount++
		pt.totalPomodoros++
		pt.notify("Pomodoro Completed!", fmt.Sprintf("You've completed %d pomodoros. Take a break!", pt.pomodoroCount))
		if pt.pomodoroCount%pt.Config.PomodorosBeforeLongBreak == 0 {
			pt.state = StateLongBreak
			pt.currentPhase = pt.Config.LongBreakDuration
			pt.remaining = pt.Config.LongBreakDuration
			pt.clearScreen()
			fmt.Printf("Starting Long Break: %s for %s\n", pt.state, pt.currentPhase)
			pt.notify("Long Break Time!", fmt.Sprintf("Enjoy your %s long break.", pt.currentPhase))
		} else {
			pt.state = StateShortBreak
			pt.currentPhase = pt.Config.ShortBreakDuration
			pt.remaining = pt.Config.ShortBreakDuration
			pt.clearScreen()
			fmt.Printf("Starting Short Break: %s for %s\n", pt.state, pt.currentPhase)
			pt.notify("Short Break Time!", fmt.Sprintf("Enjoy your %s short break.", pt.currentPhase))
		}
	case StateShortBreak, StateLongBreak:
		pt.state = StateWorking
		pt.currentPhase = pt.Config.WorkDuration
		pt.remaining = pt.Config.WorkDuration
		pt.clearScreen()
		fmt.Printf("Starting Work Session: %s for %s\n", pt.state, pt.currentPhase)
		pt.notify("Break Over!", fmt.Sprintf("Time to focus again for %s.", pt.currentPhase))
	case StateIdle, StatePaused:
		fmt.Println("Cannot transition from Idle or Paused directly. Start or Resume first.")
		return
	}
	pt.displayTimer()
}

func (pt *PomodoroTimer) Pause() {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	if pt.state != StateWorking && pt.state != StateShortBreak && pt.state != StateLongBreak {
		fmt.Println("Timer is not running to be paused.")
		return
	}
	pt.pauseChan <- struct{}{}
}

func (pt *PomodoroTimer) Resume() {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	if pt.state != StatePaused {
		fmt.Println("Timer is not paused to be resumed.")
		return
	}
	pt.resumeChan <- struct{}{}
}

func (pt *PomodoroTimer) Stop() {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	if pt.state == StateIdle {
		fmt.Println("Timer is not running.")
		return
	}
	pt.stopChan <- struct{}{}
}

func (pt *PomodoroTimer) Skip() {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	if pt.state == StateIdle {
		fmt.Println("Timer is not running to skip a phase.")
		return
	}
	pt.skipChan <- struct{}{}
}

func (pt *PomodoroTimer) displayTimer() {
	pt.clearScreen()
	pt.mu.Lock()
	defer pt.mu.Unlock()
	fmt.Printf("State: %s | Remaining: %s | Pomodoros: %d | Total: %d\n",
		pt.state, pt.remaining.Round(time.Second), pt.pomodoroCount, pt.totalPomodoros)
	fmt.Println("Commands: (s)top, (p)ause, (r)esume, (k)skip, (q)uit")
}

func (pt *PomodoroTimer) clearScreen() {
	cmd := ""
	if runtime.GOOS == "windows" {
		cmd = "cls"
	} else {
		cmd = "clear"
	}
	c := exec.Command(cmd)
	c.Stdout = os.Stdout
	c.Run()
}

func (pt *PomodoroTimer) notify(title, message string) {
	err := beeep.Notify(title, message, "")
	if err != nil {
		fmt.Printf("Error sending notification: %v\n", err)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	config := PomodoroConfig{
		WorkDuration:             25 * time.Minute,
		ShortBreakDuration:       5 * time.Minute,
		LongBreakDuration:        15 * time.Minute,
		PomodorosBeforeLongBreak: 4,
	}

	fmt.Println("Welcome to Go Pomodoro Timer!")
	fmt.Println("Do you want to use custom durations? (y/N)")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "y" {
		fmt.Print("Enter Work Duration (minutes, default 25): ")
		workStr, _ := reader.ReadString('\n')
		if val, err := strconv.Atoi(strings.TrimSpace(workStr)); err == nil && val > 0 {
			config.WorkDuration = time.Duration(val) * time.Minute
		}

		fmt.Print("Enter Short Break Duration (minutes, default 5): ")
		shortBreakStr, _ := reader.ReadString('\n')
		if val, err := strconv.

// Additional implementation at 2025-06-18 02:03:38
