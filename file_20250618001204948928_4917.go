package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func main() {
	timestamp := time.Now().Format("20060102150405")
	tempDir := os.TempDir()
	screenshotFileName := fmt.Sprintf("screenshot_%s.png", timestamp)
	screenshotPath := tempDir + string(os.PathSeparator) + screenshotFileName

	var screenshotCmd string
	var screenshotArgs []string

	switch runtime.GOOS {
	case "linux":
		screenshotCmd = "scrot"
		screenshotArgs = []string{"-o", screenshotPath}
	case "darwin":
		screenshotCmd = "screencapture"
		screenshotArgs = []string{screenshotPath}
	case "windows":
		screenshotCmd = "powershell.exe"
		psCommand := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms; `+
			`$screen = [System.Windows.Forms.Screen]::PrimaryScreen; `+
			`$bmp = New-Object System.Drawing.Bitmap($screen.Bounds.Width, $screen.Bounds.Height); `+
			`$graphics = [System.Drawing.Graphics]::FromImage($bmp); `+
			`$graphics.CopyFromScreen($screen.Bounds.Location, [System.Drawing.Point]::Empty, $screen.Bounds.Size); `+
			`$bmp.Save(\"%s\", [System.Drawing.Imaging.ImageFormat]::Png)`, screenshotPath)
		screenshotArgs = []string{"-Command", psCommand}
	default:
		fmt.Println("Unsupported operating system for screenshot.")
		os.Exit(1)
	}

	defer func() {
		if err := os.Remove(screenshotPath); err != nil {
			fmt.Printf("Error removing temporary screenshot file %s: %v\n", screenshotPath, err)
		}
	}()

	cmd := exec.Command(screenshotCmd, screenshotArgs...)
	cmd.Stderr = os.Stderr

	fmt.Printf("Taking screenshot using %s...\n", screenshotCmd)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error taking screenshot: %v\n", err)
		fmt.Println("Please ensure the necessary screenshot tool is installed and in your PATH:")
		switch runtime.GOOS {
		case "linux":
			fmt.Println("  - On Linux, try 'sudo apt-get install scrot' or 'sudo pacman -S scrot'")
		case "darwin":
			fmt.Println("  - 'screencapture' is built-in on macOS.")
		case "windows":
			fmt.Println("  - 'powershell.exe' is built-in on Windows.")
		}
		os.Exit(1)
	}

	if _, err := os.Stat(screenshotPath); os.IsNotExist(err) {
		fmt.Printf("Screenshot file was not created at %s. Check permissions or command output.\n", screenshotPath)
		os.Exit(1)
	}

	ocrCmd := exec.Command("tesseract", screenshotPath, "-")
	ocrOutput, err := ocrCmd.Output()
	if err != nil {
		fmt.Printf("Error running OCR: %v\n", err)
		fmt.Println("Please ensure 'tesseract' is installed and in your PATH.")
		fmt.Println("  - On Linux: 'sudo apt-get install tesseract-ocr'")
		fmt.Println("  - On macOS: 'brew install tesseract'")
		fmt.Println("  - On Windows: Download from https://tesseract-ocr.github.io/tessdoc/Installation.html")
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Tesseract stderr: %s\n", exitErr.Stderr)
		}
		os.Exit(1)
	}

	fmt.Println("\n--- OCR Result ---")
	fmt.Println(strings.TrimSpace(string(ocrOutput)))
	fmt.Println("------------------")
}

// Additional implementation at 2025-06-18 00:13:16
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	langFlag      string
	outputFile    string
	clipboardFlag bool
	areaFlag      bool
	keepTempFlag  bool
)

func init() {
	flag.StringVar(&langFlag, "lang", "eng", "Language for OCR (e.g., eng, deu, fra). Requires Tesseract language packs.")
	flag.StringVar(&outputFile, "output", "", "Path to save the OCR text. If empty, prints to console.")
	flag.BoolVar(&clipboardFlag, "clipboard", false, "Copy OCR text to the system clipboard.")
	flag.BoolVar(&areaFlag, "area", false, "Interactively select an area for the screenshot (not supported on all OS/tools).")
	flag.BoolVar(&keepTempFlag, "keep-temp", false, "Do not delete temporary screenshot files.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  A command-line tool to capture a screenshot and perform OCR.\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nNote: Requires 'tesseract' and OS-specific screenshot tools (e.g., 'gnome-screenshot', 'screencapture', 'powershell').\n")
		fmt.Fprintf(os.Stderr, "      For clipboard functionality, 'xclip' (Linux) or 'pbcopy' (macOS) might be needed.\n")
	}
}

func main() {
	flag.Parse()

	tempFile, err := ioutil.TempFile("", "screenshot-*.png")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	tempScreenshotPath := tempFile.Name()
	tempFile.Close()

	if !keepTempFlag {
		defer func() {
			if err := os.Remove(tempScreenshotPath); err != nil {
				log.Printf("Warning: Failed to remove temporary screenshot file %s: %v", tempScreenshotPath, err)
			}
			if outputFile == "" {
				ocrTempTxtPath := strings.TrimSuffix(tempScreenshotPath, ".png") + ".txt"
				if _, err := os.Stat(ocrTempTxtPath); err == nil {
					if err := os.Remove(ocrTempTxtPath); err != nil {
						log.Printf("Warning: Failed to remove temporary OCR output file %s: %v", ocrTempTxtPath, err)
					}
				}
			}
		}()
	}

	log.Printf("Capturing screenshot to: %s", tempScreenshotPath)
	if err := captureScreenshot(tempScreenshotPath, areaFlag); err != nil {
		log.Fatalf("Failed to capture screenshot: %v", err)
	}
	log.Println("Screenshot captured.")

	log.Printf("Performing OCR on %s (language: %s)...", tempScreenshotPath, langFlag)
	ocrText, err := performOCR(tempScreenshotPath, langFlag)
	if err != nil {
		log.Fatalf("Failed to perform OCR: %v", err)
	}
	log.Println("OCR complete.")

	if outputFile != "" {
		if err := ioutil.WriteFile(outputFile, []byte(ocrText), 0644); err != nil {
			log.Fatalf("Failed to write OCR text to file %s: %v", outputFile, err)
		}
		log.Printf("OCR text saved to %s", outputFile)
	} else {
		fmt.Println("\n--- OCR Result ---")
		fmt.Println(ocrText)
		fmt.Println("------------------")
	}

	if clipboardFlag {
		if err := copyToClipboard(ocrText); err != nil {
			log.Printf("Warning: Failed to copy OCR text to clipboard: %v", err)
		} else {
			log.Println("OCR text copied to clipboard.")
		}
	}
}

func captureScreenshot(outputPath string, selectArea bool) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		args := []string{"-x"}
		if selectArea {
			args = append(args, "-i")
		}
		args = append(args, outputPath)
		cmd = exec.Command("screencapture", args...)
	case "linux":
		if selectArea {
			_, errGnome := exec.LookPath("gnome-screenshot")
			_, errScrot := exec.LookPath("scrot")

			if errGnome == nil {
				cmd = exec.Command("gnome-screenshot", "-a", "-f", outputPath)
			} else if errScrot == nil {
				log.Println("Using 'scrot' for area selection. Please select an area.")
				cmd = exec.Command("scrot", "-s", outputPath)
			} else {
				return fmt.Errorf("neither 'gnome-screenshot' nor 'scrot' found for area selection. Try installing one of them or run without --area.")
			}
		} else {
			_, errGnome := exec.LookPath("gnome-screenshot")
			_, errScrot := exec.LookPath("scrot")

			if errGnome == nil {
				cmd = exec.Command("gnome-screenshot", "-f", outputPath)
			} else if errScrot == nil {
				cmd = exec.Command("scrot", outputPath)
			} else {
				return fmt.Errorf("neither 'gnome-screenshot' nor 'scrot' found. Please install one of them.")
			}
		}
	case "windows":
		if selectArea {
			log.Println("Area selection is not directly supported on Windows using built-in command-line tools. Capturing full screen.")
		}
		psScript := fmt.Sprintf(`
			Add-Type -AssemblyName System.Drawing
			Add-Type -AssemblyName System.Windows.Forms
			$screen = [System.Windows.Forms.Screen]::PrimaryScreen
			$bounds = $screen.Bounds
			$bmp = New-Object System.Drawing.Bitmap($bounds.Width, $bounds.Height)
			$graphics = [System.Drawing.Graphics]::FromImage($bmp)
			$graphics.CopyFromScreen($bounds.Left, $bounds.Top, 0, 0, $bounds.Size)
			$bmp.Save('%s', [System.Drawing.Imaging.ImageFormat]::Png)
			$graphics.Dispose()
			$bmp.Dispose()
		`, outputPath)
		cmd = exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psScript)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %v, stderr: %s", err, stderr.String())
	}

	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return fmt.Errorf("screenshot file not found or accessible: %v", err)
	}
	if fileInfo.Size() == 0 {
		return fmt.Errorf("screenshot file is empty, capture might have failed silently")
	}

	return nil
}

func performOCR(imagePath, lang string) (string, error) {
	outputBase := strings.TrimSuffix(imagePath, ".png")
	cmd := exec.Command("tesseract", imagePath, outputBase, "-l", lang)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("tesseract command failed: %v, stderr: %s", err, stderr.String())
	}

	ocrTextPath := outputBase + ".txt"
	textBytes, err := ioutil.ReadFile(ocrTextPath)
	if err != nil {
		return "", fmt.Errorf("failed to read OCR output file %s: %v", ocrTextPath, err)
	}

	return strings.TrimSpace(string(textBytes)), nil
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		_, errXclip := exec.LookPath("xclip")
		_, errXsel := exec.LookPath("xsel")

		if errXclip == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if errXsel == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("neither 'xclip' nor 'xsel' found. Please install one of them for clipboard support.")
		}
	case "windows":
		cmd = exec.Command("cmd", "/C", "clip")
	default:
		return fmt.Errorf("clipboard functionality not supported on %s", runtime.GOOS)
	}

	cmd.Stdin = strings.NewReader(text)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clipboard command failed: %v, stderr: %s", err, stderr.String())
	}
	return nil
}

// Additional implementation at 2025-06-18 00:14:35
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func main() {
	outputFile := flag.String("o", "", "Output OCR text to a file instead of stdout")
	copyToClipboard := flag.Bool("c", false, "Copy OCR text to clipboard")
	ocrLang := flag.String("l", "eng", "OCR language (e.g., eng, deu, fra)")
	selectArea := flag.Bool("s", false, "Select a screen area for the screenshot (macOS/Linux X11 only)")
	flag.Parse()

	tempScreenshotFile, err := ioutil.TempFile("", "screenshot-*.png")
	if err != nil {
		log.Fatalf("Failed to create temporary screenshot file: %v", err)
	}
	tempScreenshotPath := tempScreenshotFile.Name()
	tempScreenshotFile.Close()
	defer os.Remove(tempScreenshotPath)

	if err := takeScreenshot(tempScreenshotPath, *selectArea); err != nil {
		log.Fatalf("Failed to take screenshot: %v", err)
	}

	tempOCRBaseName := strings.TrimSuffix(tempScreenshotPath, ".png")
	tempOCRFile := tempOCRBaseName + ".txt"
	defer os.Remove(tempOCRFile)

	ocrText, err := performOCR(tempScreenshotPath, tempOCRBaseName, *ocrLang)
	if err != nil {
		log.Fatalf("Failed to perform OCR: %v", err)
	}

	if *outputFile != "" {
		if err := ioutil.WriteFile(*outputFile, []byte(ocrText), 0644); err != nil {
			log.Fatalf("Failed to write OCR text to file %s: %v", *outputFile, err)
		}
		fmt.Printf("OCR text saved to %s\n", *outputFile)
	} else {
		fmt.Println(ocrText)
	}

	if *copyToClipboard {
		if err := copyToClipboardFunc(ocrText); err != nil {
			log.Fatalf("Failed to copy OCR text to clipboard: %v", err)
		}
		fmt.Println("OCR text copied to clipboard.")
	}
}

func takeScreenshot(outputPath string, selectArea bool) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		args := []string{}
		if selectArea {
			args = append(args, "-i", "-s")
		}
		args = append(args, outputPath)
		cmd = exec.Command("screencapture", args...)
	case "linux":
		if selectArea {
			_, err := exec.LookPath("gnome-screenshot")
			if err == nil {
				cmd = exec.Command("gnome-screenshot", "-a", "-f", outputPath)
			} else {
				_, err := exec.LookPath("scrot")
				if err == nil {
					cmd = exec.Command("scrot", "-s", outputPath)
				} else {
					return fmt.Errorf("no suitable screenshot tool found for selection (gnome-screenshot, scrot). For Wayland, selection requires 'grim' and 'slurp' which is not universally supported by this tool.")
				}
			}
		} else {
			_, err := exec.LookPath("gnome-screenshot")
			if err == nil {
				cmd = exec.Command("gnome-screenshot", "-f", outputPath)
			} else {
				_, err := exec.LookPath("scrot")
				if err == nil {
					cmd = exec.Command("scrot", outputPath)
				} else {
					_, err := exec.LookPath("grim")
					if err == nil {
						cmd = exec.Command("grim", outputPath)
					} else {
						return fmt.Errorf("no suitable screenshot tool found (gnome-screenshot, scrot, grim)")
					}
				}
			}
		}
	case "windows":
		return fmt.Errorf("automatic screenshot on Windows is not directly supported by this tool via system commands. Please take a screenshot manually and save it to %s, then re-run.", outputPath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if cmd == nil {
		return fmt.Errorf("screenshot command not initialized for %s", runtime.GOOS)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("screenshot command failed: %v\nStderr: %s", err, stderr.String())
	}

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		if _, err := os.Stat(outputPath); err == nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}

	return fmt.Errorf("screenshot file %s not found after command execution", outputPath)
}

func performOCR(imagePath, outputBaseName, lang string) (string, error) {
	_, err := exec.LookPath("tesseract")
	if err != nil {
		return "", fmt.Errorf("tesseract not found in PATH. Please install it (e.g., 'brew install tesseract' on macOS, 'sudo apt-get install tesseract-ocr' on Linux).")
	}

	args := []string{imagePath, outputBaseName}
	if lang != "" {
		args = append(args, "-l", lang)
	}

	cmd := exec.Command("tesseract", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("tesseract command failed: %v\nStderr: %s", err, stderr.String())
	}

	ocrTextPath := outputBaseName + ".txt"
	content, err := ioutil.ReadFile(ocrTextPath)
	if err != nil {
		return "", fmt.Errorf("failed to read OCR output file %s: %v", ocrTextPath, err)
	}

	return strings.TrimSpace(string(content)), nil
}

func copyToClipboardFunc(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		_, err := exec.LookPath("xclip")
		if err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else {
			_, err := exec.LookPath("wl-copy")
			if err == nil {
				cmd = exec.Command("wl-copy")
			} else {
				return fmt.Errorf("no suitable clipboard tool found (xclip, wl-copy). Please install one.")
			}
		}
	case "windows":
		cmd = exec.Command("cmd", "/C", "clip")
	default:
		return fmt.Errorf("unsupported operating system for clipboard: %s", runtime.GOOS)
	}

	if cmd == nil {
		return fmt.Errorf("clipboard command not initialized for %s", runtime.GOOS)
	}

	cmd.Stdin = strings.NewReader(text)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("clipboard command failed: %v\nStderr: %s", err, stderr.String())
	}
	return nil
}

// Additional implementation at 2025-06-18 00:15:56
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func takeScreenshot(outputPath string, selectArea bool) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		if selectArea {
			cmd = exec.Command("screencapture", "-i", outputPath)
		} else {
			cmd = exec.Command("screencapture", outputPath)
		}
	case "linux":
		if _, err := exec.LookPath("gnome-screenshot"); err == nil {
			if selectArea {
				cmd = exec.Command("gnome-screenshot", "-a", "-f", outputPath)
			} else {
				cmd = exec.Command("gnome-screenshot", "-f", outputPath)
			}
		} else if _, err := exec.LookPath("scrot"); err == nil {
			if selectArea {
				cmd = exec.Command("scrot", "-s", outputPath)
			} else {
				cmd = exec.Command("scrot", outputPath)
			}
		} else if _, err := exec.LookPath("maim"); err == nil {
			if selectArea {
				cmd = exec.Command("maim", "-s", outputPath)
			} else {
				cmd = exec.Command("maim", outputPath)
			}
		} else {
			return fmt.Errorf("no screenshot tool found. Please install gnome-screenshot, scrot, or maim")
		}
	case "windows":
		psScript := fmt.Sprintf(`
			Add-Type -AssemblyName System.Windows.Forms
			Add-Type -AssemblyName System.Drawing
			$screen = [System.Windows.Forms.Screen]::PrimaryScreen.Bounds
			$bmp = New-Object System.Drawing.Bitmap($screen.Width, $screen.Height)
			$graphics = [System.Drawing.Graphics]::FromImage($bmp)
			$graphics.CopyFromScreen($screen.Left, $screen.Top, 0, 0, $bmp.Size)
			$bmp.Save('%s', [System.Drawing.Imaging.ImageFormat]::Png)
		`, outputPath)
		cmd = exec.Command("powershell.exe", "-NoProfile", "-Command", psScript)
		if selectArea {
			fmt.Println("Warning: Interactive selection (-s) is not supported on Windows with built-in tools. Taking full screenshot.")
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to take screenshot: %w, stderr: %s", err, stderr.String())
	}
	return nil
}

func performOCR(imagePath, lang string) (string, error) {
	tesseractPath, err := exec.LookPath("tesseract")
	if err != nil {
		return "", fmt.Errorf("tesseract not found. Please install Tesseract OCR and ensure it's in your PATH")
	}

	args := []string{imagePath, "stdout"}
	if lang != "" {
		args = append(args, "-l", lang)
	}

	cmd := exec.Command(tesseractPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("tesseract OCR failed: %w, stderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard")
		} else {
			return fmt.Errorf("no clipboard tool found. Please install xclip or xsel")
		}
	case "windows":
		cmd = exec.Command("cmd", "/c", "clip")
	default:
		return fmt.Errorf("unsupported operating system for clipboard: %s", runtime.GOOS)
	}

	cmd.Stdin = strings.NewReader(text)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w, stderr: %s", err, stderr.String())
	}
	return nil
}

func main() {
	outputFile := ""
	ocrLang := "eng"
	copyClipboard := false
	deleteTemp := false
	selectArea := false

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-o":
			if i+1 < len(args) {
				outputFile = args[i+1]
				i++
			} else {
				log.Fatalf("Error: -o requires an output file path.")
			}
		case "-l":
			if i+1 < len(args) {
				ocrLang = args[i+1]
				i++
			} else {
				log.Fatalf("Error: -l requires a language code.")
			}
		case "-c":
			copyClipboard = true
		case "-d":
			deleteTemp = true
		case "-s":
			selectArea = true
		case "-h", "--help":
			fmt.Println("Usage: go-ocr-screenshot [options]")
			fmt.Println("Options:")
			fmt.Println("  -o <file>    Save OCR text to specified file.")
			fmt.Println("  -l <lang>    Specify OCR language (e.g., 'eng', 'spa'). Default is 'eng'.")
			fmt.Println("  -c           Copy OCR text to clipboard.")
			fmt.Println("  -d           Delete the temporary screenshot file after OCR.")
			fmt.Println("  -s           Select a specific area of the screen for screenshot (interactive).")
			fmt.Println("  -h, --help   Show this help message.")
			os.Exit(0)
		default:
			log.Fatalf("Unknown argument: %s. Use -h for help.", args[i])
		}
	}

	tempFile, err := ioutil.TempFile(os.TempDir(), "screenshot-*.png")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFilePath := tempFile.Name()
	tempFile.Close()

	fmt.Printf("Taking screenshot to: %s\n", tempFilePath)
	err = takeScreenshot(tempFilePath, selectArea)
	if err != nil {
		log.Fatalf("Error taking screenshot: %v", err)
	}
	fmt.Println("Screenshot taken.")

	time.Sleep(500 * time.Millisecond)

	fileInfo, err := os.Stat(tempFilePath)
	if err != nil {
		log.Fatalf("Error accessing screenshot file: %v", err)
	}
	if fileInfo.Size() == 0 {
		log.Fatalf("Screenshot file is empty. This might indicate a problem with the screenshot tool or user cancellation.")
	}

	fmt.Printf("Performing OCR on %s with language '%s'...\n", tempFilePath, ocrLang)
	ocrText, err := performOCR(tempFilePath, ocrLang)
	if err != nil {
		log.Fatalf("Error performing OCR: %v", err)
	}
	fmt.Println("OCR complete.")

	fmt.Println("\n--- OCR Result ---")
	fmt.Println(ocrText)
	fmt.Println("------------------")

	if outputFile != "" {
		err = ioutil.WriteFile(outputFile, []byte(ocrText), 0644)
		if err != nil {
			log.Printf("Warning: Failed to write OCR text to %s: %v", outputFile, err)
		} else {
			fmt.Printf("OCR text saved to %s\n", outputFile)
		}
	}

	if copyClipboard {
		err = copyToClipboard(ocrText)
		if err != nil {
			log.Printf("Warning: Failed to copy OCR text to clipboard: %v", err)
		} else {
			fmt.Println("OCR text copied to clipboard.")
		}
	}

	if deleteTemp {
		err = os.Remove(tempFilePath)
		if err != nil {
			log.Printf("Warning: Failed to delete temporary screenshot file %s: %v", tempFilePath, err)
		} else {
			fmt.Printf("Temporary screenshot file %s deleted.\n", tempFilePath)
		}
	}
}