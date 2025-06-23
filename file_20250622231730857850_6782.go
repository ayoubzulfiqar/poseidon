package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tuotoo/qrcode"
)

func main() {
	content := "https://golang.org" // Default content

	if len(os.Args) > 1 {
		content = os.Args[1]
	}

	// Generate QR code as ASCII art
	// qrcode.M is Medium error correction level
	// 1 is the module size (character width for ASCII output)
	// true includes a quiet zone border
	// true specifies ASCII output
	qrBytes, err := qrcode.Generate(content, qrcode.M, 1, true, true)
	if err != nil {
		log.Fatalf("Failed to generate QR code: %v", err)
	}

	fmt.Print(string(qrBytes))
}

// Additional implementation at 2025-06-22 23:18:13
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/skip2/go-qrcode"
)

func main() {
	content := flag.String("content", "", "The data to encode in the QR code (required)")
	outputFile := flag.String("output", "", "Output file to save the ASCII QR code (default: stdout)")
	blackChar := flag.String("black", "██", "Character(s) for black modules")
	whiteChar := flag.String("white", "  ", "Character(s) for white modules")
	quietZone := flag.Int("quietzone", 4, "Size of the quiet zone (border) around the QR code")
	useColor := flag.Bool("color", false, "Use ANSI colors for output (black on white background)")

	flag.Parse()

	if *content == "" {
		fmt.Println("Error: --content is required.")
		flag.Usage()
		os.Exit(1)
	}

	qr, err := qrcode.New(*content, qrcode.Medium)
	if err != nil {
		log.Fatalf("Failed to generate QR code: %v", err)
	}

	asciiArt := renderQRCodeASCII(qr, *blackChar, *whiteChar, *quietZone, *useColor)

	if *outputFile != "" {
		err = os.WriteFile(*outputFile, []byte(asciiArt), 0644)
		if err != nil {
			log.Fatalf("Failed to write to output file %s: %v", *outputFile, err)
		}
		fmt.Printf("QR code saved to %s\n", *outputFile)
	} else {
		fmt.Print(asciiArt)
	}
}

// renderQRCodeASCII generates the ASCII art representation of a QR code.
func renderQRCodeASCII(qr *qrcode.QRCode, blackChar, whiteChar string, quietZone int, useColor bool) string {
	var sb strings.Builder
	size := qr.Size

	// ANSI color codes for background
	const (
		reset   = "\033[0m"
		bgBlack = "\033[40m"
		bgWhite = "\033[47m"
	)

	// Top quiet zone
	for i := 0; i < quietZone; i++ {
		if useColor {
			sb.WriteString(bgWhite) // Quiet zone is always white
		}
		for j := 0; j < size+2*quietZone; j++ {
			sb.WriteString(whiteChar)
		}
		if useColor {
			sb.WriteString(reset)
		}
		sb.WriteString("\n")
	}

	for y := 0; y < size; y++ {
		// Left quiet zone
		if useColor {
			sb.WriteString(bgWhite)
		}
		for i := 0; i < quietZone; i++ {
			sb.WriteString(whiteChar)
		}
		if useColor {
			sb.WriteString(reset) // Reset after quiet zone, then manage color for QR modules
		}

		// QR code modules
		currentColor := ""
		if useColor {
			// Initialize currentColor based on the first module's color
			if qr.Module(y, 0) {
				currentColor = bgBlack
			} else {
				currentColor = bgWhite
			}
			sb.WriteString(currentColor)
		}

		for x := 0; x < size; x++ {
			isBlack := qr.Module(y, x)
			charToPrint := whiteChar
			targetColor := bgWhite

			if isBlack {
				charToPrint = blackChar
				targetColor = bgBlack
			}

			if useColor && targetColor != currentColor {
				sb.WriteString(reset)
				sb.WriteString(targetColor)
				currentColor = targetColor
			}
			sb.WriteString(charToPrint)
		}

		// Right quiet zone
		if useColor {
			sb.WriteString(reset) // Reset after QR modules
			sb.WriteString(bgWhite)
		}
		for i := 0; i < quietZone; i++ {
			sb.WriteString(whiteChar)
		}
		if useColor {
			sb.WriteString(reset)
		}
		sb.WriteString("\n")
	}

	// Bottom quiet zone
	for i := 0; i < quietZone; i++ {
		if useColor {
			sb.WriteString(bgWhite)
		}
		for j := 0; j < size+2*quietZone; j++ {
			sb.WriteString(whiteChar)
		}
		if useColor {
			sb.WriteString(reset)
		}
		sb.WriteString("\n")
	}

	return sb.String()
}