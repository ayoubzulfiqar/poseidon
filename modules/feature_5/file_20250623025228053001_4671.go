package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"  // Register GIF format for image.Decode
	_ "image/jpeg" // Register JPEG format for image.Decode
	_ "image/png"  // Register PNG format for image.Decode
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

func main() {
	inputDir := flag.String("input", "", "Input directory containing images")
	outputDir := flag.String("output", "", "Output directory for resized images")
	width := flag.Int("width", 0, "Target width (0 to maintain aspect ratio)")
	height := flag.Int("height", 0, "Target height (0 to maintain aspect ratio)")
	quality := flag.Int("quality", 90, "JPEG quality (1-100)")

	flag.Parse()

	if *inputDir == "" || *outputDir == "" {
		fmt.Println("Error: Input and output directories must be specified.")
		flag.Usage()
		return
	}

	if *width == 0 && *height == 0 {
		fmt.Println("Error: Either width or height (or both) must be specified.")
		flag.Usage()
		return
	}

	err := os.MkdirAll(*outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory %s: %v\n", *outputDir, err)
		return
	}

	err = resizeImages(*inputDir, *outputDir, *width, *height, *quality)
	if err != nil {
		fmt.Printf("Error during bulk resize: %v\n", err)
	}
}

func resizeImages(inputDir, outputDir string, targetWidth, targetHeight, jpegQuality int) error {
	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil // Skip directories
		}

		relPath, err := filepath.Rel(inputDir, path)
		if err != nil {
			return err
		}

		initialOutputPath := filepath.Join(outputDir, relPath)

		outputFileDir := filepath.Dir(initialOutputPath)
		if err := os.MkdirAll(outputFileDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory %s: %w", outputFileDir, err)
		}

		finalOutputPath, err := resizeImage(path, initialOutputPath, targetWidth, targetHeight, jpegQuality)
		if err != nil {
			fmt.Printf("Warning: Could not resize %s: %v\n", path, err)
		} else {
			fmt.Printf("Resized %s to %s\n", path, finalOutputPath)
		}
		return nil
	})
}

func resizeImage(inputPath, initialOutputPath string, targetWidth, targetHeight, jpegQuality int) (string, error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open input image %s: %w", inputPath, err)
	}
	defer inputFile.Close()

	img, format, err := image.Decode(inputFile)
	if err != nil {
		return "", fmt.Errorf("failed to decode image %s: %w", inputPath, err)
	}

	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	newWidth, newHeight := targetWidth, targetHeight

	if newWidth == 0 && newHeight == 0 {
		return "", fmt.Errorf("target width and height cannot both be zero")
	}

	if newWidth == 0 {
		newWidth = int(float64(originalWidth) * (float64(newHeight) / float64(originalHeight)))
	} else if newHeight == 0 {
		newHeight = int(float64(originalHeight) * (float64(newWidth) / float64(originalWidth)))
	}

	dstImg := image.NewNRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(dstImg, dstImg.Bounds(), img, img.Bounds(), draw.Over, nil)

	outputExt := filepath.Ext(initialOutputPath)
	outputFormat := format

	switch format {
	case "jpeg":
		// Keep as JPEG
	case "png":
		// Keep as PNG
	case "gif":
		// Convert GIF to PNG for simplicity (single frame, transparency)
		outputFormat = "png"
		outputExt = ".png"
	default:
		// Unknown format, default to PNG
		outputFormat = "png"
		outputExt = ".png"
	}

	finalOutputPath := initialOutputPath
	if strings.ToLower(filepath.Ext(initialOutputPath)) != outputExt {
		finalOutputPath = strings.TrimSuffix(initialOutputPath, filepath.Ext(initialOutputPath)) + outputExt
	}

	outputFile, err := os.Create(finalOutputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output image %s: %w", finalOutputPath, err)
	}
	defer outputFile.Close()

	switch outputFormat {
	case "jpeg":
		err = jpeg.Encode(outputFile, dstImg, &jpeg.Options{Quality: jpegQuality})
	case "png":
		err = png.Encode(outputFile, dstImg)
	default:
		err = fmt.Errorf("unsupported output format: %s", outputFormat)
	}

	if err != nil {
		return "", fmt.Errorf("failed to encode image %s: %w", finalOutputPath, err)
	}

	return finalOutputPath, nil
}

// Additional implementation at 2025-06-23 02:53:52
package main

import (
	"fmt"
	"image"
	"image/draw" // Standard library for image drawing and scaling
	"image/jpeg"
	"image/png"
	_ "image/gif"  // Register GIF format decoder
	_ "image/jpeg" // Register JPEG format decoder
	_ "image/png"  // Register PNG format decoder
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ImageResizeOptions holds options for image resizing
type ImageResizeOptions struct {
	InputPath  string
	OutputPath string
	Width      int
	Height     int
	Format     string // "jpeg", "png"
	Quality    int    // For JPEG (1-100)
}

// resizeImage resizes a single image based on options
func resizeImage(options ImageResizeOptions) error {
	file, err := os.Open(options.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open image %s: %w", options.InputPath, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image %s: %w", options.InputPath, err)
	}

	originalBounds := img.Bounds()
	originalWidth := originalBounds.Dx()
	originalHeight := originalBounds.Dy()

	targetWidth := options.Width
	targetHeight := options.Height

	// Calculate target dimensions preserving aspect ratio if only one dimension is provided
	if targetWidth == 0 && targetHeight == 0 {
		return fmt.Errorf("either width or height must be specified for resizing")
	} else if targetWidth == 0 {
		targetWidth = int(float64(originalWidth) * (float64(targetHeight) / float64(originalHeight)))
	} else if targetHeight == 0 {
		targetHeight = int(float64(originalHeight) * (float64(targetWidth) / float64(originalWidth)))
	}

	// Create a new image with the target dimensions
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Use BiLinear for scaling, which is part of the standard library's image/draw package.
	// This provides a reasonable quality balance for standard library usage.
	draw.BiLinear.Scale(dst, dst.Bounds(), img, originalBounds, draw.Over, nil)

	outFile, err := os.Create(options.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", options.OutputPath, err)
	}
	defer outFile.Close()

	switch strings.ToLower(options.Format) {
	case "jpeg":
		err = jpeg.Encode(outFile, dst, &jpeg.Options{Quality: options.Quality})
	case "png":
		err = png.Encode(outFile, dst)
	default:
		return fmt.Errorf("unsupported output format: %s", options.Format)
	}

	if err != nil {
		return fmt.Errorf("failed to encode image to %s: %w", options.OutputPath, err)
	}

	return nil
}

// processImages finds and processes images concurrently
func processImages(inputDir, outputDir string, width, height int, format string, quality, workers int) {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatalf("Failed to create output directory %s: %v", outputDir, err)
		}
	}

	var imagePaths []string
	err := filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
				imagePaths = append(imagePaths, path)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to walk input directory %s: %v", inputDir, err)
	}

	if len(imagePaths) == 0 {
		log.Printf("No supported images found in %s", inputDir)
		return
	}

	log.Printf("Found %d images to process in %s", len(imagePaths), inputDir)
	log.Printf("Processing with %d workers...", workers)

	var wg sync.WaitGroup
	jobs := make(chan string, workers) // Buffered channel to limit concurrent jobs

	// Start worker goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for inputPath := range jobs {
				// Determine output file path and name
				relPath, err := filepath.Rel(inputDir, inputPath)
				if err != nil {
					log.Printf("Error getting relative path for %s: %v", inputPath, err)
					continue
				}

				outputFileName := strings.TrimSuffix(filepath.Base(relPath), filepath.Ext(relPath)) + "." + format
				outputPath := filepath.Join(outputDir, outputFileName)

				options := ImageResizeOptions{
					InputPath:  inputPath,
					OutputPath: outputPath,
					Width:      width,
					Height:     height,
					Format:     format,
					Quality:    quality,
				}

				if err := resizeImage(options); err != nil {
					log.Printf("Error processing %s: %v", inputPath, err)
				} else {
					log.Printf("Successfully processed %s to %s", filepath.Base(inputPath), filepath.Base(outputPath))
				}
			}
		}()
	}

	// Send image paths to the jobs channel
	for _, path := range imagePaths {
		jobs <- path
	}
	close(jobs) // Close the channel to signal workers no more jobs are coming

	wg.Wait() // Wait for all workers to finish
	log.Println("All images processed.")
}

func main() {
	// Parse command line arguments manually for simplicity in a single file
	args := os.Args[1:] // Skip the program name

	inputDir := ""
	outputDir := ""
	width := 0
	height := 0
	format := "jpeg" // Default output format
	quality := 80    // Default JPEG quality
	workers := 4     // Default number of concurrent workers

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-inputDir":
			if i+1 < len(args) {
				inputDir = args[i+1]
				i++
			}
		case "-outputDir":
			if i+1 < len(args) {
				outputDir = args[i+1]
				i++
			}
		case "-width":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &width)
				i++
			}
		case "-height":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &height)
				i++
			}
		case "-format":
			if i+1 < len(args) {
				format = args[i+1]
				i++
			}
		case "-quality":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &quality)
				i++
			}
		case "-workers":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &workers)
				i++
			}
		case "-h", "--help":
			fmt.Println("Usage: go run main.go -inputDir <dir> -outputDir <dir> [-width <px>] [-height <px>] [-format <jpeg|png>] [-quality <1-100>] [-workers <num>]")
			fmt.Println("  -inputDir: Source directory containing images.")
			fmt.Println("  -outputDir: Destination directory for resized images.")
			fmt.Println("  -width: Target width in pixels. If height is 0, aspect ratio is preserved.")
			fmt.Println("  -height: Target height in pixels. If width is 0, aspect ratio is preserved.")
			fmt.Println("  -format: Output image format (jpeg or png). Default is jpeg.")
			fmt.Println("  -quality: JPEG quality (1-100). Default is 80. Only applicable for JPEG format.")
			fmt.Println("  -workers: Number of concurrent workers. Default is 4.")
			os.Exit(0)
		}
	}

	if inputDir == "" || outputDir == "" {
		log.Fatal("Error: -inputDir and -outputDir are required.")
	}
	if width == 0 && height == 0 {
		log.Fatal("Error: Either -width or -height must be specified.")
	}
	if quality < 1 || quality > 100 {
		log.Fatal("Error: JPEG quality must be between 1 and 100.")
	}
	if workers < 1 {
		log.Fatal("Error: Number of workers must be at least 1.")
	}

	processImages(inputDir, outputDir, width, height, format, quality, workers)
}