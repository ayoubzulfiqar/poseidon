package main

import (
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/kbinani/screenshot"
)

func main() {
	n := screenshot.NumActiveDisplays()

	if n == 0 {
		fmt.Println("No active displays found.")
		return
	}

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.Capture(bounds.Min.X, bounds.Min.Y, bounds.Dx(), bounds.Dy())
		if err != nil {
			fmt.Printf("Failed to capture display %d: %v\n", i, err)
			continue
		}

		fileName := fmt.Sprintf("screenshot-display-%d-%s.png", i, time.Now().Format("20060102-150405"))
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf("Failed to create file %s: %v\n", fileName, err)
			continue
		}
		defer file.Close()

		err = png.Encode(file, img)
		if err != nil {
			fmt.Printf("Failed to encode image to PNG for %s: %v\n", fileName, err)
			continue
		}

		fmt.Printf("Screenshot saved to %s\n", fileName)
	}
}

// Additional implementation at 2025-08-04 09:04:16
package main

import (
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/vova616/go-screenshot/screenshot"
)

func main() {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("screenshot_%s.png", timestamp)

	fmt.Printf("Capturing screen and saving to %s...\n", filename)

	img, err := screenshot.CaptureScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error capturing screen: %v\n", err)
		os.Exit(1)
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file %s: %v\n", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding image to PNG: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Screenshot saved successfully to %s\n", filename)
}

// Additional implementation at 2025-08-04 09:05:20
package main

import (
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/kbinani/screenshot"
)

func main() {
	// Additional functionality: Add a delay before capturing
	delaySeconds := 3
	fmt.Printf("Capturing screen in %d seconds...\n", delaySeconds)
	time.Sleep(time.Duration(delaySeconds) * time.Second)

	// Additional functionality: Allow specifying output filename via command-line argument
	outputFileName := "screenshot.png"
	if len(os.Args) > 1 {
		outputFileName = os.Args[1]
		if _, err := os.Stat(outputFileName); err == nil {
			fmt.Printf("Warning: File '%s' already exists and will be overwritten.\n", outputFileName)
		}
	}

	n := screenshot.NumDisplay()
	if n == 0 {
		fmt.Println("No displays found.")
		os.Exit(1)
	}

	// Capture the primary display (display 0)
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.Capture(bounds.Min.X, bounds.Min.Y, bounds.Dx(), bounds.Dy())
	if err != nil {
		fmt.Printf("Error capturing screen: %v\n", err)
		os.Exit(1)
	}

	file, err := os.Create(outputFileName)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", outputFileName, err)
		os.Exit(1)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		fmt.Printf("Error encoding image to PNG: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Screenshot saved to %s\n", outputFileName)
}