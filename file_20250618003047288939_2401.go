package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func main() {
	fmt.Println("Simple Go Voice Recorder")
	fmt.Println("Press Enter to start recording.")
	fmt.Println("Press Enter again to stop recording.")

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	outputDir := "recordings"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", outputDir, err)
		return
	}

	timestamp := time.Now().Format("20060102_150405")
	outputFile := filepath.Join(outputDir, fmt.Sprintf("recording_%s.wav", timestamp))

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("rec", outputFile, "rate", "44100", "channels", "1", "enc", "signed-integer", "bits", "16")
	case "windows":
		cmd = exec.Command("sox", "-d", outputFile, "rate", "44100", "channels", "1", "enc", "signed-integer", "bits", "16")
	default:
		fmt.Println("Unsupported operating system for direct recording command.")
		fmt.Println("Please ensure 'rec' (SoX) is installed and in your PATH.")
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Starting recording to %s...\n", outputFile)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error starting recording: %v\n", err)
		fmt.Println("Please ensure 'rec' (SoX) is installed and in your system's PATH.")
		fmt.Println("For Linux/macOS: `brew install sox` or `sudo apt-get install sox`")
		fmt.Println("For Windows: Download from sox.sourceforge.net and add to PATH.")
		return
	}

	fmt.Println("Recording... Press Enter to stop.")
	_, _ = reader.ReadString('\n')

	err = cmd.Process.Kill()
	if err != nil {
		fmt.Printf("Error stopping recording process: %v\n", err)
	} else {
		fmt.Println("Recording stopped.")
	}

	_ = cmd.Wait()

	fmt.Printf("Recording saved to: %s\n", outputFile)
	fmt.Println("Exiting.")
}

// Additional implementation at 2025-06-18 00:31:22
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	recordingsDir = "recordings"
	fileExtension = ".wav" // Simulating WAV files
)

// AudioInput represents a source of audio data (e.g., microphone).
type AudioInput interface {
	Start() error
	Read(p []byte) (n int, err error)
	Stop() error
}

// AudioOutput represents a destination for audio data (e.g., speakers).
type AudioOutput interface {
	Start() error
	Write(p []byte) (n int, err error)
	Stop() error
}

// DummyAudioDevice simulates an audio device by generating/consuming dummy data.
// In a real application, this would interact with OS audio APIs (e.g., PortAudio, WASAPI, CoreAudio).
type DummyAudioDevice struct {
	mu      sync.Mutex
	running bool
}

// NewDummyAudioDevice creates a new dummy audio device.
func NewDummyAudioDevice() *DummyAudioDevice {
	return &DummyAudioDevice{}
}

// Start initializes the dummy device.
func (d *DummyAudioDevice) Start() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.running {
		return fmt.Errorf("dummy device already running")
	}
	d.running = true
	fmt.Println("[Dummy Audio] Device started.")
	return nil
}

// Read simulates reading audio data from a microphone.
// It generates dummy byte data.
func (d *DummyAudioDevice) Read(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.running {
		return 0, io.EOF // Signal that the device is not running
	}
	// Simulate reading a small amount of data
	for i := range p {
		p[i] = byte(time.Now().Nanosecond() % 256) // Dummy data
	}
	time.Sleep(10 * time.Millisecond) // Simulate real-time delay
	return len(p), nil
}

// Write simulates writing audio data to speakers.
// It just "consumes" the data.
func (d *DummyAudioDevice) Write(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.running {
		return 0, fmt.Errorf("dummy device not running")
	}
	// In a real scenario, this would send data to audio output buffer.
	time.Sleep(10 * time.Millisecond) // Simulate real-time delay
	return len(p), nil
}

// Stop de-initializes the dummy device.
func (d *DummyAudioDevice) Stop() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.running {
		return fmt.Errorf("dummy device not running")
	}
	d.running = false
	fmt.Println("[Dummy Audio] Device stopped.")
	return nil
}

// Recorder manages the audio recording process.
type Recorder struct {
	audioInput AudioInput
	recording  bool
	outputFile *os.File
	stopChan   chan struct{}
	wg         sync.WaitGroup
	mu         sync.Mutex // Protects recording state
}

// NewRecorder creates a new Recorder instance.
func NewRecorder(input AudioInput) *Recorder {
	return &Recorder{
		audioInput: input,
	}
}

// StartRecording begins the recording process.
// It creates a new file and starts reading from the audio input.
func (r *Recorder) StartRecording(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.recording {
		return fmt.Errorf("already recording")
	}

	if err := os.MkdirAll(recordingsDir, 0755); err != nil {
		return fmt.Errorf("failed to create recordings directory: %w", err)
	}

	filePath := filepath.Join(recordingsDir, name+fileExtension)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create recording file: %w", err)
	}

	if err := r.audioInput.Start(); err != nil {
		file.Close()
		os.Remove(filePath) // Clean up partially created file
		return fmt.Errorf("failed to start audio input: %w", err)
	}

	r.outputFile = file
	r.recording = true
	r.stopChan = make(chan struct{})
	r.wg.Add(1)

	fmt.Printf("Recording started: %s\n", name)

	go func() {
		defer r.wg.Done()
		defer r.outputFile.Close()
		buffer := make([]byte, 4096) // Buffer for reading audio data
		for {
			select {
			case <-r.stopChan:
				fmt.Println("Recording stopped.")
				return
			default:
				n, err := r.audioInput.Read(buffer)
				if err != nil {
					if err != io.EOF {
						log.Printf("Error reading audio input: %v", err)
					}
					// If audio input signals EOF (e.g., device stopped), stop recording.
					log.Println("Audio input stream ended, stopping recording goroutine.")
					return
				}
				if n > 0 {
					if _, err := r.outputFile.Write(buffer[:n]); err != nil {
						log.Printf("Error writing to recording file: %v", err)
						return
					}
				}
			}
		}
	}()
	return nil
}

// StopRecording stops the current recording.
func (r *Recorder) StopRecording() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.recording {
		return fmt.Errorf("not currently recording")
	}

	close(r.stopChan) // Signal the recording goroutine to stop
	r.wg.Wait()       // Wait for the goroutine to finish

	if err := r.audioInput.Stop(); err != nil {
		log.Printf("Error stopping audio input: %v", err)
	}

	r.recording = false
	r.outputFile = nil // Clear file handle
	return nil
}

// Player manages the audio playback process.
type Player struct {
	audioOutput AudioOutput
	playing     bool
	stopChan    chan struct{}
	wg          sync.WaitGroup
	mu          sync.Mutex // Protects playing state
}

// NewPlayer creates a new Player instance.
func NewPlayer(output AudioOutput) *Player {
	return &Player{
		audioOutput: output,
	}
}

// PlayRecording plays a recorded file.
func (p *Player) PlayRecording(name string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.playing {
		return fmt.Errorf("already playing")
	}

	filePath := filepath.Join(recordingsDir, name+fileExtension)
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open recording file: %w", err)
	}
	// file.Close() is deferred in the goroutine, not here.

	if err := p.audioOutput.Start(); err != nil {
		file.Close() // Close file if audio output fails to start
		return fmt.Errorf("failed to start audio output: %w", err)
	}

	p.playing = true
	p.stopChan = make(chan struct{})
	p.wg.Add(1)

	fmt.Printf("Playing: %s\n", name)

	go func() {
		defer p.wg.Done()
		defer file.Close() // Ensure file is closed when goroutine exits
		buffer := make([]byte, 4096) // Buffer for playing audio data
		for {
			select {
			case <-p.stopChan:
				fmt.Println("Playback stopped.")
				p.StopPlayback() // Update player state
				return
			default:
				n, err := file.Read(buffer)
				if err != nil {
					if err != io.EOF {
						log.Printf("Error reading recording file: %v", err)
					}
					fmt.Println("Playback finished.")
					p.StopPlayback() // Update player state and stop audio output
					return
				}
				if n > 0 {
					if _, err := p.audioOutput.Write(buffer[:n]); err != nil {
						log.Printf("Error writing to audio output: %v", err)
						p.StopPlayback() // Stop if there's an output error
						return
					}
				}
			}
		}
	}()
	return nil
}

// StopPlayback stops the current playback.
func (p *Player) StopPlayback() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.playing {
		return fmt.Errorf("not currently playing")
	}

	// Signal the playback goroutine to stop, if it's not already done.
	select {
	case <-p.stopChan:
		// Already closed
	default:
		close(p.stopChan)
	}

	p.wg.Wait() // Wait for the goroutine to finish

	if err := p.audioOutput.Stop(); err != nil {
		log.Printf("Error stopping audio output: %v", err)
	}

	p.playing = false
	return nil
}

// RecordingManager handles listing and deleting recordings.
type RecordingManager struct{}

// NewRecordingManager creates a new RecordingManager instance.
func NewRecordingManager() *RecordingManager {
	return &RecordingManager{}
}

// ListRecordings returns a list of available recording names.
func (rm *RecordingManager) ListRecordings() ([]string, error) {
	files, err := ioutil.ReadDir(recordingsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil // Directory doesn't exist, no recordings
		}
		return nil, fmt.Errorf("failed to read recordings directory: %w", err)
	}

	var names []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), fileExtension) {
			name := strings.TrimSuffix(file.Name(), fileExtension)
			names = append(names, name)
		}
	}
	return names, nil
}

// DeleteRecording deletes a specific recording file.
func (rm *RecordingManager) DeleteRecording(name string) error {
	filePath := filepath.Join(recordingsDir, name+fileExtension)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("recording '%s' not found", name)
	}
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete recording '%s': %w", name, err)
	}
	fmt.Printf("Recording '%s' deleted.\n", name)
	return nil
}

// main function for the CLI application.
func main() {
	// Initialize dummy audio devices
	dummyAudioInput := NewDummyAudioDevice()
	dummyAudioOutput := NewDummyAudioDevice()

	recorder := NewRecorder(dummyAudioInput)
	player := NewPlayer(dummyAudioOutput)
	manager := NewRecordingManager()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Go Voice Recorder (Dummy Audio)")
	fmt.Println("Type 'help' for commands.")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		command := strings.ToLower(parts[0])

		switch command {
		case "record":
			if len(parts) < 2 {
				fmt.Println("Usage: record <name>")
				continue
			}
			name := parts[1]
			if err := recorder.StartRecording(name); err != nil {
				fmt.Printf("Error starting recording: %v\n", err)
			}
		case "stop":
			if err := recorder.StopRecording(); err != nil {
				fmt.Printf("Error stopping recording: %v\n", err)
			}
		case "play":
			if len(parts