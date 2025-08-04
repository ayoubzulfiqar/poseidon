package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var asciiArtFont = map[rune][]string{
	' ': {
		"     ",
		"     ",
		"     ",
		"     ",
		"     ",
	},
	'A': {
		" AAA ",
		"A   A",
		"AAAAA",
		"A   A",
		"A   A",
	},
	'B': {
		"BBBB ",
		"B   B",
		"BBBB ",
		"B   B",
		"BBBB ",
	},
	'C': {
		" CCCCC",
		"C     ",
		"C     ",
		"C     ",
		" CCCCC",
	},
	'D': {
		"DDDD ",
		"D   D",
		"D   D",
		"D   D",
		"DDDD ",
	},
	'E': {
		"EEEEE",
		"E    ",
		"EEEEE",
		"E    ",
		"EEEEE",
	},
	'F': {
		"FFFFF",
		"F    ",
		"FFFF ",
		"F    ",
		"F    ",
	},
	'G': {
		" GGG ",
		"G    ",
		"G GGG",
		"G   G",
		" GGG ",
	},
	'H': {
		"H   H",
		"H   H",
		"HHHHH",
		"H   H",
		"H   H",
	},
	'I': {
		"IIIII",
		"  I  ",
		"  I  ",
		"  I  ",
		"IIIII",
	},
	'J': {
		"JJJJJ",
		"    J",
		"    J",
		"J   J",
		" JJJ ",
	},
	'K': {
		"K   K",
		"K  K ",
		"KKK  ",
		"K  K ",
		"K   K",
	},
	'L': {
		"L    ",
		"L    ",
		"L    ",
		"L    ",
		"LLLLL",
	},
	'M': {
		"M   M",
		"MM MM",
		"M M M",
		"M   M",
		"M   M",
	},
	'N': {
		"N   N",
		"NN  N",
		"N N N",
		"N  NN",
		"N   N",
	},
	'O': {
		" OOO ",
		"O   O",
		"O   O",

// Additional implementation at 2025-08-04 06:15:01
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"  // Register GIF format
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"os"
)

const defaultCharset = " .:-=+*#%@"

func main() {
	imagePath := flag.String("image", "", "Path to the image file")
	outputWidth := flag.Int("width", 100, "Desired output width in characters")
	invertColors := flag.Bool("invert", false, "Invert ASCII character mapping (dark pixels become light chars)")
	charset := flag.String("charset", defaultCharset, "Custom character set from darkest to lightest")

	flag.Parse()

	if *imagePath == "" {
		fmt.Println("Usage: go run main.go -image <path_to_image> [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(*imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening image file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding image: %v\n", err)
		os.Exit(1)
	}

	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	if *outputWidth <= 0 {
		fmt.Fprintf(os.Stderr, "Output width must be greater than 0.\n")
		os.Exit(1)
	}

	// Adjust height for console character aspect ratio (typically 2:1 height:width)
	aspectRatioCorrection := 2.0
	newHeight := int(float64(originalHeight) / float64(originalWidth) * float64(*outputWidth) / aspectRatioCorrection)
	if newHeight == 0 {
		newHeight = 1
	}

	resizedImg := image.NewGray(image.Rect(0, 0, *outputWidth, newHeight))
	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Src, nil)

	charSetLen := len(*charset)
	if charSetLen == 0 {
		fmt.Fprintf(os.Stderr, "Character set cannot be empty.\n")
		os.Exit(1)
	}

	for y := 0; y < newHeight; y++ {
		for x := 0; x < *outputWidth; x++ {
			grayColor := resizedImg.At(x, y).(color.Gray)
			brightness := grayColor.Y

			var charIndex int
			if *invertColors {
				charIndex = int(float64(255-brightness) / 255.0 * float64(charSetLen-1))
			} else {
				charIndex = int(float64(brightness) / 255.0 * float64(charSetLen-1))
			}

			if charIndex < 0 {
				charIndex = 0
			} else if charIndex >= charSetLen {
				charIndex = charSetLen - 1
			}

			fmt.Print(string((*charset)[charIndex]))
		}
		fmt.Println()
	}
}

// Additional implementation at 2025-08-04 06:15:57
package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif" // Register GIF format
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"log"
	"os"
	"flag"
	"image/draw"
)

// Default character set, from light to dark (for dark console background)
const defaultChars = " .:-=+*#%@"

// toGrayscale converts a color.Color to an 8-bit grayscale value (0-255).
// It uses the NTSC standard for luminance calculation.
func toGrayscale(c color.Color) uint8 {
	r, g, b, _ := c.RGBA()
	// RGBA values are 16-bit, so divide by 257 to get 8-bit equivalent.
	// Y = 0.299R + 0.587G + 0.114B
	return uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 257)
}

// resizeImage resizes an image to a new width, maintaining aspect ratio.
// It uses a simple nearest-neighbor scaling algorithm.
func resizeImage(img image.Image, newWidth int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if newWidth <= 0 || newWidth == width {
		return img // No resizing needed or invalid width
	}

	// Calculate new height maintaining aspect ratio
	newHeight := (height * newWidth) / width
	if newHeight == 0 { // Ensure height is at least 1
		newHeight = 1
	}

	// Create a new RGBA image with the desired dimensions
	resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Iterate over the new image pixels and map them back to original image
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Map new coordinates back to original image coordinates
			origX := (x * width) / newWidth
			origY := (y * height) / newHeight
			resized.Set(x, y, img.At(origX, origY))
		}
	}
	return resized
}

func main() {
	imagePath := flag.String("image", "", "Path to the image file (JPG, PNG, GIF)")
	outputWidth := flag.Int("width", 100, "Desired output width in characters (0 for original width)")
	charSet := flag.String("chars", defaultChars, "Custom character set from light to dark (e.g., ' .:-=+*#%@')")

	flag.Parse()

	if *imagePath == "" {
		fmt.Println("Usage: go run main.go -image <path_to_image> [-width <int>] [-chars <string>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(*imagePath)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}

	// Resize image if a specific width is provided and different from original
	if *outputWidth > 0 && *outputWidth != img.Bounds().Dx() {
		img = resizeImage(img, *outputWidth)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Convert the image to RGBA format for consistent pixel access.
	// This handles various input image types (e.g., Paletted, YCbCr).
	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, img, bounds.Min, draw.Src)

	charLen := len(*charSet)
	if charLen == 0 {
		log.Fatal("Character set cannot be empty.")
	}

	// Iterate over each pixel, convert to grayscale, and map to a character
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixelColor := rgbaImg.At(x, y)
			grayValue := toGrayscale(pixelColor)

			// Map grayscale value (0-255) to an index in the character set.
			// The division by 256.0 ensures a float result before scaling.
			charIndex := int(float64(grayValue) / 256.0 * float64(charLen))
			
			// Ensure charIndex does not exceed the bounds of the character set
			if charIndex >= charLen {
				charIndex = charLen - 1
			}
			
			fmt.Print(string((*charSet)[charIndex]))
		}
		fmt.Println() // New line after each row of characters
	}
}

// Additional implementation at 2025-08-04 06:16:34
package main

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif" // Register GIF format
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"os"
	"strings"
	"flag"
)

// Default character set for brightness mapping (from light to dark)
const defaultCharSet = " .:-=+*#%@"

// loadImage loads an image from the given file path.
func loadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open image file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %w", err)
	}
	return img, nil
}

// resizeImage resizes the image to the target width, maintaining aspect ratio.
// It also adjusts the height to account for the typical aspect ratio of console characters (taller than wide).
func resizeImage(img image.Image, targetWidth int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width == 0 || height == 0 {
		return nil // Avoid division by zero or zero-sized image
	}

	// Adjust aspect ratio for console characters (typically characters are about twice as tall as wide)
	// A factor of 0.55 to 0.6 is common for monospace fonts.
	aspectRatio := float64(height) / float64(width)
	targetHeight := int(float64(targetWidth) * aspectRatio * 0.55) 

	if targetWidth == 0 || targetHeight == 0 {
		return nil // Avoid creating zero-sized image
	}

	resizedImg := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	draw.ApproxBiLinear.Scale(resizedImg, resizedImg.Bounds(), img, bounds, draw.Src, nil)
	return resizedImg
}

// pixelToBrightness converts an RGBA color to a grayscale brightness value (0-255).
// It uses the luminosity method for a more perceptually accurate brightness.
func pixelToBrightness(r, g, b, a uint32) uint8 {
	// Convert 16-bit color components to 8-bit
	r8 := uint8(r >> 8)
	g8 := uint8(g >> 8)
	b8 := uint8(b >> 8)
	a8 := uint8(a >> 8)

	// If alpha is 0, it's fully transparent, consider it black (or darkest)
	if a8 == 0 {
		return 0 
	}

	// Luminosity method: 0.299*R + 0.587*G + 0.114*B
	brightness := float64(r8)*0.299 + float64(g8)*0.587 + float64(b8)*0.114
	return uint8(brightness)
}

// mapBrightnessToChar maps a brightness value (0-255) to a character from the character set.
// If invert is true, it maps dark pixels to light characters and vice-versa.
func mapBrightnessToChar(brightness uint8, charSet string, invert bool) string {
	idx := int(float64(brightness) / 256.0 * float64(len(charSet)))
	
	// Ensure index is within bounds
	if idx >= len(charSet) {
		idx = len(charSet) - 1
	}
	if idx < 0 {
		idx = 0
	}

	if invert {
		idx = len(charSet) - 1 - idx
	}
	return string(charSet[idx])
}

// generateAsciiArt converts a resized image into a 2D slice of ASCII characters.
func generateAsciiArt(img image.Image, charSet string, invert bool) [][]string {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	asciiArt := make([][]string, height)
	for y := 0; y < height; y++ {
		asciiArt[y] = make([]string, width)
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			brightness := pixelToBrightness(r, g, b, a)
			asciiArt[y][x] = mapBrightnessToChar(brightness, charSet, invert)
		}
	}
	return asciiArt
}

// printAsciiArt prints the ASCII art to the console.
// If color is true, it uses ANSI escape codes to print characters with their corresponding pixel color.
func printAsciiArt(asciiArt [][]string, sourceImg image.Image, color bool) {
	asciiHeight := len(asciiArt)
	if asciiHeight == 0 {
		return
	}
	asciiWidth := len(asciiArt[0])

	for y := 0; y < asciiHeight; y++ {
		for x := 0; x < asciiWidth; x++ {
			char := asciiArt[y][x]
			if color {
				// Get the color from the source image (which is the resized image)
				r, g, b, _ := sourceImg.At(x, y).RGBA()
				// Print character with 24-bit ANSI color (foreground)
				fmt.Printf("\033[38;2;%d;%d;%dm%s\033[0m", r>>8, g>>8, b>>8, char)
			} else {
				fmt.Print(char)
			}
		}
		fmt.Println() // Newline after each row
	}
}

// saveAsciiArt saves the generated ASCII art to a text file.
func saveAsciiArt(asciiArt [][]string, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, row := range asciiArt {
		_, err := writer.WriteString(strings.Join(row, "") + "\n")
		if err != nil {
			return fmt.Errorf("could not write to file: %w", err)
		}
	}
	return writer.Flush() // Ensure all buffered data is written to file
}

func main() {
	// Define command-line flags
	imagePath := flag.String("path", "", "Path to the image file (required)")
	outputWidth := flag.Int("width", 100, "Target width of the ASCII art in characters")
	charSet := flag.String("chars", defaultCharSet, "Custom character set for brightness mapping (from light to dark)")
	invertColors := flag.Bool("invert", false, "Invert brightness mapping (dark pixels become light characters)")
	enableColor := flag.Bool("color", false, "Enable color output (requires ANSI compatible terminal)")
	outputPath := flag.String("output", "", "Optional path to save ASCII art to a text file instead of printing to console")

	flag.Parse()

	// Validate required flags
	if *imagePath == "" {
		fmt.Println("Error: Image path is required. Use -path flag.")
		flag.Usage() // Print usage instructions
		os.Exit(1)
	}

	// Load the image
	img, err := loadImage(*imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading image: %v\n", err)
		os.Exit(1)
	}

	// Resize the image
	resizedImg := resizeImage(img, *outputWidth)
	if resizedImg == nil {
		fmt.Fprintf(os.Stderr, "Error: Resized image is nil. Check image dimensions or target width.\n")
		os.Exit(1)
	}

	// Generate ASCII art from the resized image
	asciiArt := generateAsciiArt(resizedImg, *charSet, *invertColors)

	// Output the ASCII art
	if *outputPath != "" {
		err := saveAsciiArt(asciiArt, *outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving ASCII art to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("ASCII art successfully saved to %s\n", *outputPath)
	} else {
		printAsciiArt(asciiArt, resizedImg, *enableColor)
	}
}