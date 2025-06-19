package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"syscall"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	gdi32            = syscall.NewLazyDLL("gdi32.dll")
	getDC            = user32.NewProc("GetDC")
	releaseDC        = user32.NewProc("ReleaseDC")
	getSystemMetrics = user32.NewProc("GetSystemMetrics")
	createCompatibleDC = gdi32.NewProc("CreateCompatibleDC")
	createCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	selectObject     = gdi32.NewProc("SelectObject")
	bitBlt           = gdi32.NewProc("BitBlt")
	deleteObject     = gdi32.NewProc("DeleteObject")
	deleteDC         = gdi32.NewProc("DeleteDC")
	getDIBits        = gdi32.NewProc("GetDIBits")
	getObject        = gdi32.NewProc("GetObjectW")
)

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
	SRCCOPY     = 0x00CC0020
	DIB_RGB_COLORS = 0
	BI_RGB = 0
)

type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

type BITMAP struct {
	BmType       int32
	BmWidth      int32
	BmHeight     int32
	BmWidthBytes int32
	BmPlanes     uint16
	BmBitsPixel  uint16
	BmBits       uintptr
}

func captureScreen(filename string) error {
	hDC, _, _ := getDC.Call(0)
	if hDC == 0 {
		return fmt.Errorf("GetDC failed")
	}
	defer releaseDC.Call(0, hDC)

	width, _, _ := getSystemMetrics.Call(SM_CXSCREEN)
	height, _, _ := getSystemMetrics.Call(SM_CYSCREEN)

	hMemDC, _, _ := createCompatibleDC.Call(hDC)
	if hMemDC == 0 {
		return fmt.Errorf("CreateCompatibleDC failed")
	}
	defer deleteDC.Call(hMemDC)

	hBitmap, _, _ := createCompatibleBitmap.Call(hDC, width, height)
	if hBitmap == 0 {
		return fmt.Errorf("CreateCompatibleBitmap failed")
	}
	defer deleteObject.Call(hBitmap)

	oldBitmap, _, _ := selectObject.Call(hMemDC, hBitmap)
	if oldBitmap == 0 {
		return fmt.Errorf("SelectObject failed")
	}
	defer selectObject.Call(hMemDC, oldBitmap)

	ret, _, _ := bitBlt.Call(hMemDC, 0, 0, width, height, hDC, 0, 0, SRCCOPY)
	if ret == 0 {
		return fmt.Errorf("BitBlt failed")
	}

	var bm BITMAP
	ret, _, _ = getObject.Call(hBitmap, uintptr(unsafe.Sizeof(bm)), uintptr(unsafe.Pointer(&bm)))
	if ret == 0 {
		return fmt.Errorf("GetObject failed")
	}

	bmiHeader := BITMAPINFOHEADER{
		BiSize:        uint32(unsafe.Sizeof(BITMAPINFOHEADER{})),
		BiWidth:       int32(width),
		BiHeight:      -int32(height),
		BiPlanes:      1,
		BiBitCount:    32,
		BiCompression: BI_RGB,
	}

	imageSize := int(width) * int(height) * 4
	bits := make([]byte, imageSize)

	ret, _, _ = getDIBits.Call(
		hMemDC,
		hBitmap,
		0,
		uintptr(height),
		uintptr(unsafe.Pointer(&bits[0])),
		uintptr(unsafe.Pointer(&bmiHeader)),
		DIB_RGB_COLORS,
	)
	if ret == 0 {
		return fmt.Errorf("GetDIBits failed")
	}

	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	for i := 0; i < len(bits); i += 4 {
		img.Pix[i+0] = bits[i+2]
		img.Pix[i+1] = bits[i+1]
		img.Pix[i+2] = bits[i+0]
		img.Pix[i+3] = 0xFF
	}

	outFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	return png.Encode(outFile, img)
}

func main() {
	fmt.Println("Capturing screen...")
	err := captureScreen("screenshot.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error capturing screen: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Screenshot saved to screenshot.png")
}

// Additional implementation at 2025-06-18 23:52:11
package main

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"github.com/go-vgo/robotgo"
)

const (
	captureDelaySeconds = 3       // Delay before capturing the screen
	outputDir           = "screenshots" // Directory to save screenshots
)

func main() {
	fmt.Printf("Screen capture will begin in %d seconds...\n", captureDelaySeconds)
	time.Sleep(captureDelaySeconds * time.Second)

	// Ensure the output directory exists
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory %s: %v\n", outputDir, err)
		return
	}

	// Generate a unique filename based on the current timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("screenshot_%s.png", timestamp)
	fullPath := filepath.Join(outputDir, filename)

	fmt.Printf("Capturing screen and saving to %s...\n", fullPath)

	// Capture the entire screen
	img, err := robotgo.CaptureImg()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error capturing screen: %v\n", err)
		return
	}

	// Create the output file
	file, err := os.Create(fullPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file %s: %v\n", fullPath, err)
		return
	}
	defer file.Close() // Ensure the file is closed

	// Encode the image to PNG format and write to the file
	err = png.Encode(file, img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding image to PNG: %v\n", err)
		return
	}

	fmt.Printf("Screenshot saved successfully to %s\n", fullPath)
}

// Additional implementation at 2025-06-18 23:52:48
