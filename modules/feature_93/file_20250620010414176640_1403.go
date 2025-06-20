package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage: go run main.go <input_dir> <output_dir> <width> <height>")
		os.Exit(1)
	}

	inputDir := os.Args[1]
	outputDir := os.Args[2]
	widthStr := os.Args[3]
	heightStr := os.Args[4]

	width, err := strconv.Atoi(widthStr)
	if err != nil {
		fmt.Printf("Invalid width: %v\n", err)
		os.Exit(1)
	}
	height, err := strconv.Atoi(heightStr)
	if err != nil {
		fmt.Printf("Invalid height: %v\n", err)
		os.Exit(1)
	}

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory %s: %v\n", outputDir, err)
		os.Exit(1)
	}

	files, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Printf("Error reading input directory %s: %v\n", inputDir, err)
		os.Exit(1)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		inputPath := filepath.Join(inputDir, file.Name())
		ext := strings.ToLower(filepath.Ext(file.Name()))
		baseName := strings.TrimSuffix(file.Name(), ext)
		outputPath := filepath.Join(outputDir, baseName+"_resized"+ext)

		fmt.Printf("Processing %s...\n", inputPath)

		reader, err := os.Open(inputPath)
		if err != nil {
			fmt.Printf("Error opening %s: %v\n", inputPath, err)
			continue
		}

		img, format, err := image.Decode(reader)
		reader.Close()
		if err != nil {
			fmt.Printf("Error decoding %s: %v\n", inputPath, err)
			continue
		}

		resizedImg := image.NewRGBA(image.Rect(0, 0, width, height))
		draw.NearestNeighbor.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Src, nil)

		outputFile, err := os.Create(outputPath)
		if err != nil {
			fmt.Printf("Error creating output file %s: %v\n", outputPath, err)
			continue
		}

		switch format {
		case "jpeg":
			err = jpeg.Encode(outputFile, resizedImg, &jpeg.Options{Quality: 90})
		case "png":
			err = png.Encode(outputFile, resizedImg)
		default:
			fmt.Printf("Unsupported format '%s' for %s. Attempting to save as PNG.\n", format, inputPath)
			outputPath = filepath.Join(outputDir, baseName+"_resized.png")
			outputFile.Close()
			outputFile, err = os.Create(outputPath)
			if err != nil {
				fmt.Printf("Error creating PNG output file %s: %v\n", outputPath, err)
				continue
			}
			err = png.Encode(outputFile, resizedImg)
		}
		outputFile.Close()

		if err != nil {
			fmt.Printf("Error encoding %s: %v\n", outputPath, err)
			continue
		}

		fmt.Printf("Resized %s to %dx%d and saved to %s\n", file.Name(), width, height, outputPath)
	}
	fmt.Println("Bulk resizing complete.")
}

// Additional implementation at 2025-06-20 01:05:09
package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/image/draw"
)

func resizeImage(inputPath, outputPath string, width, height int, quality int, format string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open image %s: %w", inputPath, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image %s: %w", inputPath, err)
	}

	originalBounds := img.Bounds()
	originalWidth := originalBounds.Dx()
	originalHeight := originalBounds.Dy()

	if width == 0 && height == 0 {
		return fmt.Errorf("either width or height must be specified for resizing")
	}

	if width == 0 {
		width = int(float64(originalWidth) * (float64(height) / float64(originalHeight)))
	} else if height == 0 {
		height = int(float64(originalHeight) * (float64(width) / float64(originalWidth)))
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Src, nil)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
	}
	defer outFile.Close()

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(outFile, dst, &jpeg.Options{Quality: quality})
	case "png":
		err = png.Encode(outFile, dst)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to encode image %s: %w", outputPath, err)
	}

	return nil
}

func processDirectory(inputDir, outputDir string, width, height, quality int, format string, workers int) {
	var wg sync.WaitGroup
	jobs := make(chan string)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for inputPath := range jobs {
				relPath, err := filepath.Rel(inputDir, inputPath)
				if err != nil {
					log.Printf("Error getting relative path for %s: %v", inputPath, err)
					continue
				}
				outputPath := filepath.Join(outputDir, relPath)

				outputFileDir := filepath.Dir(outputPath)
				if err := os.MkdirAll(outputFileDir, 0755); err != nil {
					log.Printf("Error creating output directory %s: %v", outputFileDir, err)
					continue
				}

				ext := filepath.Ext(outputPath)
				base := strings.TrimSuffix(outputPath, ext)
				outputPath = fmt.Sprintf("%s.%s", base, format)

				log.Printf("Resizing %s to %s...", inputPath, outputPath)
				if err := resizeImage(inputPath, outputPath, width, height, quality, format); err != nil {
					log.Printf("Failed to resize %s: %v", inputPath, err)
				}
			}
		}()
	}

	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v", path, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
			jobs <- path
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking input directory %s: %v", inputDir, err)
	}

	close(jobs)
	wg.Wait()
	log.Println("All images processed.")
}

func main() {
	inputDir := flag.String("input", "", "Input directory containing images.")
	outputDir := flag.String("output", "resized_images", "Output directory for resized images.")
	width := flag.Int("width", 0, "Target width (pixels). If 0, aspect ratio is maintained based on height.")
	height := flag.Int("height", 0, "Target height (pixels). If 0, aspect ratio is maintained based on width.")
	quality := flag.Int("quality", 90, "JPEG quality (1-100). Only applicable for JPEG output.")
	format := flag.String("format", "jpeg", "Output format (jpeg or png).")
	workers := flag.Int("workers", 4, "Number of concurrent workers (goroutines).")

	flag.Parse()

	if *inputDir == "" {
		log.Fatal("Input directory must be specified using -input flag.")
	}
	if *width == 0 && *height == 0 {
		log.Fatal("Either -width or -height must be specified.")
	}

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory %s: %v", *outputDir, err)
	}

	processDirectory(*inputDir, *outputDir, *width, *height, *quality, *format, *workers)
}