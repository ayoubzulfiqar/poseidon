package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

const (
	SPI_SETDESKWALLPAPER = 0x0014
	SPIF_UPDATEINIFILE   = 0x01
	SPIF_SENDCHANGE      = 0x02
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	systemParametersInfo = user32.NewProc("SystemParametersInfoW")
)

func setWallpaperWindows(imagePath string) error {
	pathPtr, err := syscall.UTF16PtrFromString(imagePath)
	if err != nil {
		return fmt.Errorf("failed to convert path to UTF16 pointer: %w", err)
	}

	ret, _, err := systemParametersInfo.Call(
		uintptr(SPI_SETDESKWALLPAPER),
		uintptr(0),
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(SPIF_UPDATEINIFILE|SPIF_SENDCHANGE),
	)

	if ret == 0 {
		return fmt.Errorf("SystemParametersInfo failed: %w", err)
	}
	return nil
}

func setWallpaperDarwin(imagePath string) error {
	script := fmt.Sprintf(`tell application "Finder" to set desktop picture to POSIX file "%s"`, imagePath)
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set wallpaper via osascript: %w\nOutput: %s", err, output)
	}
	return nil
}

func setWallpaperLinux(imagePath string) error {
	gsettingsCmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://"+imagePath)
	if err := gsettingsCmd.Run(); err == nil {
		exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-options", "zoom").Run()
		return nil
	}

	xfconfCmd := exec.Command("xfconf-query", "-c", "xfce4-desktop", "-p", "/backdrop/screen0/monitor0/workspace0/last-image", "-s", imagePath)
	if err := xfconfCmd.Run(); err == nil {
		return nil
	}
	
	fehCmd := exec.Command("feh", "--bg-fill", imagePath)
	if err := fehCmd.Run(); err == nil {
		return nil
	}

	return fmt.Errorf("failed to set wallpaper on Linux. Tried gsettings, xfconf-query, and feh. Ensure one of these is installed and your desktop environment is supported.")
}

func setWallpaper(imagePath string) error {
	absPath, err := filepath.Abs(imagePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("image file not found: %s", absPath)
	}

	switch runtime.GOOS {
	case "windows":
		return setWallpaperWindows(absPath)
	case "darwin":
		return setWallpaperDarwin(absPath)
	case "linux":
		return setWallpaperLinux(absPath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func listWallpapers(dirPath string) error {
	if dirPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		dirPath = filepath.Join(homeDir, "Pictures")
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			dirPath = homeDir
		}
	}

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	fmt.Printf("Wallpapers in %s:\n", dirPath)
	found := false
	for _, file := range files {
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".bmp" || ext == ".gif" {
				fmt.Println(filepath.Join(dirPath, file.Name()))
				found = true
			}
		}
	}
	if !found {
		fmt.Println("No common image files found.")
	}
	return nil
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  wallpaper set <path/to/image.jpg>")
	fmt.Println("  wallpaper list [directory]")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "set":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: Missing image path for 'set' command.")
			printUsage()
			os.Exit(1)
		}
		imagePath := os.Args[2]
		if err := setWallpaper(imagePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting wallpaper: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Wallpaper set successfully.")
	case "list":
		dirPath := ""
		if len(os.Args) >= 3 {
			dirPath = os.Args[2]
		}
		if err := listWallpapers(dirPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error listing wallpapers: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown command '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}

// Additional implementation at 2025-06-18 00:33:23


// Additional implementation at 2025-06-18 00:34:16
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	defaultWallpaperDir = "./wallpapers" // Default directory for wallpaper images
	defaultInterval     = 5 * time.Minute // Default interval for scheduled changes
)

func setWallpaper(imagePath string) error {
	absPath, err := filepath.Abs(imagePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", imagePath, err)
	}

	fmt.Printf("Attempting to set wallpaper to: %s\n", absPath)

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		script := fmt.Sprintf(`
			$regPath = "HKCU:\Control Panel\Desktop"
			Set-ItemProperty -Path $regPath -Name Wallpaper -Value "%s"
			Rundll32.exe user32.dll,UpdatePerUserSystemParameters 1, True
			`, absPath)
		cmd = exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)
	case "darwin": // macOS
		script := fmt.Sprintf(`tell application "Finder" to set desktop picture to POSIX file "%s"`, absPath)
		cmd = exec.Command("osascript", "-e", script)
	case "linux":
		cmd = exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://"+absPath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command to set wallpaper: %w", err)
	}

	fmt.Println("Wallpaper set successfully (or command executed).")
	return nil
}

func listWallpapers(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	var wallpaperPaths []string
	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".bmp":  true,
		".gif":  true,
		".webp": true,
	}

	for _, file := range files {
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if imageExtensions[ext] {
				wallpaperPaths = append(wallpaperPaths, filepath.Join(dir, file.Name()))
			}
		}
	}
	return wallpaperPaths, nil
}

func getRandomWallpaper(dir string) (string, error) {
	wallpapers, err := listWallpapers(dir)
	if err != nil {
		return "", err
	}
	if len(wallpapers) == 0 {
		return "", fmt.Errorf("no wallpaper images found in %s", dir)
	}

	randomIndex := rand.Intn(len(wallpapers))
	return wallpapers[randomIndex], nil
}

func startScheduler(dir string, interval time.Duration) {
	fmt.Printf("Starting wallpaper scheduler. Changing every %s from %s...\n", interval, dir)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("Scheduler: Time to change wallpaper...")
		wallpaper, err := getRandomWallpaper(dir)
		if err != nil {
			log.Printf("Scheduler error getting random wallpaper: %v", err)
			continue
		}
		if err := setWallpaper(wallpaper); err != nil {
			log.Printf("Scheduler error setting wallpaper %s: %v", wallpaper, err)
		}
	}
}

func printUsage() {
	fmt.Println("Usage: wallpaper-manager [command] [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  set <image_path>          Set a specific image as wallpaper.")
	fmt.Println("  list                      List all available wallpaper images in the default directory.")
	fmt.Println("  random                    Set a random wallpaper from the default directory.")
	fmt.Println("  schedule [interval]       Start a scheduler to change wallpaper periodically.")
	fmt.Println("                            Interval can be like '1h', '30m', '5s'. Default is 5m.")
	fmt.Println("  help                      Display this help message.")
	fmt.Printf("\nDefault wallpaper directory: %s\n", defaultWallpaperDir)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	if _, err := os.Stat(defaultWallpaperDir); os.IsNotExist(err) {
		fmt.Printf("Creating default wallpaper directory: %s\n", defaultWallpaperDir)
		if err := os.MkdirAll(defaultWallpaperDir, 0755); err != nil {
			log.Fatalf("Failed to create wallpaper directory: %v", err)
		}
	}

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "set":
		if len(args) < 1 {
			fmt.Println("Error: Missing image path for 'set' command.")
			printUsage()
			return
		}
		imagePath := args[0]
		if err := setWallpaper(imagePath); err != nil {
			log.Fatalf("Failed to set wallpaper: %v", err)
		}
	case "list":
		wallpapers, err := listWallpapers(defaultWallpaperDir)
		if err != nil {
			log.Fatalf("Failed to list wallpapers: %v", err)
		}
		if len(wallpapers) == 0 {
			fmt.Printf("No wallpaper images found in %s.\n", defaultWallpaperDir)
			return
		}
		fmt.Printf("Wallpapers found in %s:\n", defaultWallpaperDir)
		for _, wp := range wallpapers {
			fmt.Println("-", wp)
		}
	case "random":
		wallpaper, err := getRandomWallpaper(defaultWallpaperDir)
		if err != nil {
			log.Fatalf("Failed to get random wallpaper: %v", err)
		}
		if err := setWallpaper(wallpaper); err != nil {
			log.Fatalf("Failed to set random wallpaper: %v", err)
		}
	case "schedule":
		interval := defaultInterval
		if len(args) > 0 {
			parsedInterval, err := time.ParseDuration(args[0])
			if err != nil {
				fmt.Printf("Error: Invalid interval format '%s'. Using default %s.\n", args[0], defaultInterval)
			} else {
				interval = parsedInterval
			}
		}
		startScheduler(defaultWallpaperDir, interval)
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}