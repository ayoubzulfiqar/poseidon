package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/gordonklaus/portaudio"
)

const (
	sampleRate      = 44100
	channels        = 1
	bitsPerSample   = 32
	framesPerBuffer = 512
	outputFilename  = "recording.wav"
)

var (
	recordedSamples = &bytes.Buffer{}
)

func writeWAVHeader(w io.Writer, audioDataSize int) error {
	binary.Write(w, binary.BigEndian, []byte("RIFF"))
	totalFileSize := uint32(36 + audioDataSize)
	binary.Write(w, binary.LittleEndian, totalFileSize)
	binary.Write(w, binary.BigEndian, []byte("WAVE"))

	binary.Write(w, binary.BigEndian, []byte("fmt "))
	binary.Write(w, binary.LittleEndian, uint32(16))
	binary.Write(w, binary.LittleEndian, uint16(3))
	binary.Write(w, binary.LittleEndian, uint16(channels))
	binary.Write(w, binary.LittleEndian, uint32(sampleRate))
	byteRate := uint32(sampleRate * channels * bitsPerSample / 8)
	binary.Write(w, binary.LittleEndian, byteRate)
	blockAlign := uint16(channels * bitsPerSample / 8)
	binary.Write(w, binary.LittleEndian, blockAlign)
	binary.Write(w, binary.LittleEndian, uint16(bitsPerSample))

	binary.Write(w, binary.BigEndian, []byte("data"))
	binary.Write(w, binary.LittleEndian, uint32(audioDataSize))

	return nil
}

func main() {
	fmt.Println("Initializing PortAudio...")
	err := portaudio.Initialize()
	if err != nil {
		fmt.Printf("Error initializing PortAudio: %v\n", err)
		return
	}
	defer portaudio.Terminate()

	fmt.Println("Opening default input stream...")
	stream, err := portaudio.OpenDefaultStream(channels, 0, sampleRate, framesPerBuffer, func(in []float32) {
		err := binary.Write(recordedSamples, binary.LittleEndian, in)
		if err != nil {
			fmt.Printf("Error writing samples to buffer: %v\n", err)
		}
	})
	if err != nil {
		fmt.Printf("Error opening stream: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Println("Starting stream...")
	err = stream.Start()
	if err != nil {
		fmt.Printf("Error starting stream: %v\n", err)
		return
	}
	defer stream.Stop()

	fmt.Println("Recording... Press Ctrl+C to stop.")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	fmt.Println("Stopping recording...")

	err = stream.Stop()
	if err != nil {
		fmt.Printf("Error stopping stream: %v\n", err)
	}

	fmt.Println("Saving recorded audio to", outputFilename)

	file, err := os.Create(outputFilename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	audioDataSize := recordedSamples.Len()
	err = writeWAVHeader(file, audioDataSize)
	if err != nil {
		fmt.Printf("Error writing WAV header: %v\n", err)
		return
	}

	_, err = recordedSamples.WriteTo(file)
	if err != nil {
		fmt.Printf("Error writing audio data to file: %v\n", err)
		return
	}

	fmt.Println("Recording saved successfully.")
}

// Additional implementation at 2025-06-18 02:09:27
package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/go-audio/wav"
	"github.com/gordonklaus/portaudio"
)

const (
	sampleRate       = 44100 // CD quality
	channels         = 1     // Mono
	framesPerBuffer  = 512   // Buffer size for audio processing
	maxRecordSeconds = 10    // Maximum recording duration in seconds
)

func main() {
	fmt.Println("Go Voice Recorder")
	fmt.Println("-----------------")
	fmt.Println("Press Enter to start recording.")
	fmt.Println("Press Enter again to stop recording.")
	fmt.Println(fmt.Sprintf("Recording will automatically stop after %d seconds.", maxRecordSeconds))

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n') // Wait for first Enter to initiate

	err := portaudio.Initialize()
	if err != nil {
		fmt.Printf("Error initializing PortAudio: %v\n", err)
		return
	}
	defer portaudio.Terminate()

	var recordedSamples []int16 // Buffer to store recorded samples

	// Channel to signal when recording should stop (either by user or timeout)
	stopRecording := make(chan struct{})

	// Stream setup: 1 input channel, 0 output channels, sample rate, buffer size, callback function
	stream, err := portaudio.OpenDefaultStream(channels, 0, sampleRate, framesPerBuffer, func(in []int16) {
		// This callback is executed when new audio data is available from the microphone
		recordedSamples = append(recordedSamples, in...)
	})
	if err != nil {
		fmt.Printf("Error opening audio stream: %v\n", err)
		return
	}
	defer stream.Close() // Ensure stream is closed when main exits

	fmt.Println("Recording started...")

	// Start recording in a separate goroutine to allow concurrent user input
	go func() {
		err := stream.Start()
		if err != nil {
			fmt.Printf("Error starting stream: %v\n", err)
			close(stopRecording) // Signal to stop if stream fails to start
			return
		}

		// Wait for either a manual stop signal or the maximum duration timeout
		select {
		case <-stopRecording:
			// Manual stop triggered by user input
		case <-time.After(time.Duration(maxRecordSeconds) * time.Second):
			fmt.Println("\nMaximum recording duration reached.")
			select {
			case <-stopRecording: // Check if already closed
			default:
				close(stopRecording) // Ensure stop signal is sent if timeout occurs first
			}
		}

		err = stream.Stop()
		if err != nil {
			fmt.Printf("Error stopping stream: %v\n", err)
		}
		fmt.Println("Recording stopped.")
	}()

	// Wait for second Enter from the user to manually stop recording
	_, _ = reader.ReadString('\n')
	// If stopRecording is already closed by timeout, this is a no-op, which is fine.
	select {
	case <-stopRecording: // Check if already closed by timeout
	default:
		close(stopRecording) // Signal the goroutine to stop
	}

	// Give a moment for the stream.Stop() to complete and the goroutine to finish its cleanup
	time.Sleep(500 * time.Millisecond)

	if len(recordedSamples) == 0 {
		fmt.Println("No audio recorded.")
		return
	}

	// Prompt user for filename to save the recording
	fmt.Print("Enter filename to save (e.g., my_recording.wav): ")
	filename, _ := reader.ReadString('\n')
	filename = trimNewline(filename) // Clean up newline characters
	if filename == "" {
		// Generate a default filename if none is provided
		filename = fmt.Sprintf("recording_%s.wav", time.Now().Format("20060102_150405"))
		fmt.Printf("No filename entered, saving as %s\n", filename)
	}

	// Create and open the WAV file for writing
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close() // Ensure the file is closed

	// Initialize WAV encoder: file writer, sample rate, bit depth (16), number of channels, format (1 for PCM)
	enc := wav.NewEncoder(file, sampleRate, 16, channels, 1)
	defer enc.Close() // Ensure the encoder flushes and closes the WAV header

	// Write all recorded samples to the WAV encoder
	for _, sample := range recordedSamples {
		// Write each 16-bit sample using little-endian byte order
		err = enc.Write(binary.LittleEndian, sample)
		if err != nil {
			fmt.Printf("Error writing sample to WAV: %v\n", err)
			return
		}
	}

	fmt.Printf("Recording successfully saved to %s\n", filename)
}

// trimNewline removes trailing newline and carriage return characters from a string.
func trimNewline(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '\r' { // For Windows CRLF
		s = s[:len(s)-1]
	}
	return s
}

// Additional implementation at 2025-06-18 02:10:51
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	recordingsDir := "recordings"
	if err := os.MkdirAll(recordingsDir, 0755); err != nil {
		fmt.Printf("Error creating recordings directory: %v\n", err)
		return
	}

	fmt.Println("Press Enter to start recording...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	timestamp := time.Now().Format("20060102_150405")
	outputFile := filepath.Join(recordingsDir, fmt.Sprintf("recording_%s.wav", timestamp))

	cmd := exec.Command("ffmpeg",
		"-f", "alsa",
		"-i", "default",
		"-acodec", "pcm_s16le",
		"-ar", "44100",
		"-ac", "1",
		outputFile,
	)

	fmt.Printf("Recording started to %s...\n", outputFile)
	fmt.Println("Press Enter to stop recording.")

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error starting recording: %v\n", err)
		return
	}

	bufio.NewReader(os.Stdin).ReadBytes('\n')

	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		fmt.Printf("Error sending interrupt signal to ffmpeg: %v. Attempting to kill.\n", err)
		if err := cmd.Process.Kill(); err != nil {
			fmt.Printf("Error killing ffmpeg process: %v\n", err)
		}
	}

	err = cmd.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Recording stopped. ffmpeg exited with error: %v\n", exitErr)
		} else {
			fmt.Printf("Error waiting for recording process: %v\n", err)
		}
	} else {
		fmt.Println("Recording stopped successfully.")
	}

	fmt.Printf("Recording saved to: %s\n", outputFile)
}