package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
)

var code39CharMap = map[rune]string{
	'0': "000110100",
	'1': "100100001",
	'2': "001100001",
	'3': "101100000",
	'4': "000110001",
	'5': "100110000",
	'6': "001110000",
	'7': "000100101",
	'8': "100100100",
	'9': "001100100",
	'A': "100001001",
	'B': "001001001",
	'C': "101001000",
	'D': "000011001",
	'E': "100011000",
	'F': "001011000",
	'G': "000001101",
	'H': "100001100",
	'I': "001001100",
	'J': "000011100",
	'K': "100000011",
	'L': "001000011",
	'M': "101000010",
	'N': "000010011",
	'O': "100010010",
	'P': "001010010",
	'Q': "000000111",
	'R': "100000110",
	'S': "001000110",
	'T': "000010110",
	'U': "110000001",
	'V': "011000001",
	'W': "111000000",
	'X': "010010001",
	'Y': "110010000",
	'Z': "011010000",
	'-': "010000101",
	'.': "110000100",
	' ': "011000100",
	'$': "010101000",
	'/': "010100010",
	'+': "010001010",
	'%': "000101010",
	'*': "011010100",
}

func encodeCode39(data string) ([]int, error) {
	var pattern []int

	startPattern, ok := code39CharMap['*']
	if !ok {
		return nil, fmt.Errorf("internal error: start character '*' pattern not found")
	}
	for _, r := range startPattern {
		if r == '0' {
			pattern = append(pattern, 0)
		} else {
			pattern = append(pattern, 1)
		}
	}

	pattern = append(pattern, 0) // Inter-character gap (narrow space)

	for _, char := range strings.ToUpper(data) {
		charPattern, ok := code39CharMap[char]
		if !ok {
			return nil, fmt.Errorf("invalid character for Code 39: %c", char)
		}
		for _, r := range charPattern {
			if r == '0' {
				pattern = append(pattern, 0)
			} else {
				pattern = append(pattern, 1)
			}
		}
		pattern = append(pattern, 0) // Inter-character gap (narrow space)
	}

	stopPattern, ok := code39CharMap['*']
	if !ok {
		return nil, fmt.Errorf("internal error: stop character '*' pattern not found")
	}
	for _, r := range stopPattern {
		if r == '0' {
			pattern = append(pattern, 0)
		} else {
			pattern = append(pattern, 1)
		}
	}

	return pattern, nil
}

func generateBarcodeImage(pattern []int, filename string) error {
	const (
		narrowWidth = 2
		wideWidth   = 6
		barHeight   = 100
		paddingX    = 20
		paddingY    = 20
	)

	totalBarcodeWidth := 0
	for _, p := range pattern {
		if p == 0 {
			totalBarcodeWidth += narrowWidth
		} else {
			totalBarcodeWidth += wideWidth
		}
	}

	imgWidth := totalBarcodeWidth + 2*paddingX
	imgHeight := barHeight + 2*paddingY

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			img.Set(x, y, color.White)
		}
	}

	currentX := paddingX
	for i, p := range pattern {
		var width int
		if p == 0 {
			width = narrowWidth
		} else {
			width = wideWidth
		}

		isBar := false
		idxInBlock := i % 10
		if idxInBlock == 0 || idxInBlock == 2 || idxInBlock == 4 || idxInBlock == 6 || idxInBlock == 8 {
			isBar = true
		}

		if isBar {
			for y := paddingY; y < paddingY+barHeight; y++ {
				for x := currentX; x < currentX+width; x++ {
					img.Set(x, y, color.Black)
				}
			}
		}
		currentX += width
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <data_to_encode> [output_filename.png]")
		os.Exit(1)
	}

	data := os.Args[1]
	outputFilename := "barcode.png"
	if len(os.Args) > 2 {
		outputFilename = os.Args[2]
	}

	pattern, err := encodeCode39(data)
	if err != nil {
		fmt.Printf("Error encoding data: %v\n", err)
		os.Exit(1)
	}

	err = generateBarcodeImage(pattern, outputFilename)
	if err != nil {
		fmt.Printf("Error generating barcode image: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Barcode for \"%s\" generated successfully as %s\n", data, outputFilename)
}

// Additional implementation at 2025-06-17 23:33:56
