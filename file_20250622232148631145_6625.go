package main

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

type POINT struct {
	X, Y int32
}

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getCursorPosProc = user32.NewProc("GetCursorPos")
)

func getCursorPos() (POINT, error) {
	var pt POINT
	ret, _, err := getCursorPosProc.Call(uintptr(unsafe.Pointer(&pt)))
	if ret == 0 {
		return POINT{}, fmt.Errorf("GetCursorPos failed: %v", err)
	}
	return pt, nil
}

func main() {
	fmt.Println("Mouse movement tracker started. Press Ctrl+C to exit.")

	var lastPos POINT
	var err error

	lastPos, err = getCursorPos()
	if err != nil {
		fmt.Printf("Error getting initial cursor position: %v\n", err)
		return
	}
	fmt.Printf("Initial position: X=%d, Y=%d\n", lastPos.X, lastPos.Y)

	for {
		currentPos, err := getCursorPos()
		if err != nil {
			fmt.Printf("Error getting cursor position: %v\n", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if currentPos.X != lastPos.X || currentPos.Y != lastPos.Y {
			fmt.Printf("Mouse moved to: X=%d, Y=%d\n", currentPos.X, currentPos.Y)
			lastPos = currentPos
		}

		time.Sleep(50 * time.Millisecond)
	}
}

// Additional implementation at 2025-06-22 23:22:27
