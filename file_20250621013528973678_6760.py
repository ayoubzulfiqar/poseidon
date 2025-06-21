import subprocess
import os
import tempfile
import platform
import sys

def main():
    temp_dir = tempfile.gettempdir()
    screenshot_filename = f"screenshot_ocr_{os.getpid()}_{os.urandom(4).hex()}.png"
    screenshot_path = os.path.join(temp_dir, screenshot_filename)

    system = platform.system()
    screenshot_command = []

    if system == "Linux":
        screenshot_command = ["scrot", "-z", screenshot_path]
    elif system == "Darwin":
        screenshot_command = ["screencapture", screenshot_path]

    if not screenshot_command:
        sys.stderr.write("Error: Unsupported operating system or no suitable screenshot tool found.\n")
        sys.stderr.write("This script supports 'scrot' on Linux and 'screencapture' on macOS.\n")
        sys.stderr.write("Please ensure 'tesseract' is also installed and in your PATH.\n")
        sys.exit(1)

    try:
        screenshot_process = subprocess.run(screenshot_command, check=True, capture_output=True)
        if screenshot_process.returncode != 0:
            sys.stderr.write(f"Screenshot command failed: {screenshot_process.stderr.decode().strip()}\n")
            sys.exit(1)

        ocr_command = ["tesseract", screenshot_path, "stdout"]
        ocr_process = subprocess.run(ocr_command, check=True, capture_output=True, text=True)

        if ocr_process.returncode != 0:
            sys.stderr.write(f"OCR command failed: {ocr_process.stderr.strip()}\n")
            sys.exit(1)

        sys.stdout.write(ocr_process.stdout.strip())
        sys.stdout.write("\n")

    except FileNotFoundError as e:
        sys.stderr.write(f"Error: Required command '{e.filename}' not found. Please ensure it is installed and in your system's PATH.\n")
        sys.stderr.write("For screenshot: 'scrot' (Linux) or 'screencapture' (macOS).\n")
        sys.stderr.write("For OCR: 'tesseract'.\n")
        sys.exit(1)
    except subprocess.CalledProcessError as e:
        sys.stderr.write(f"An error occurred during command execution: {e}\n")
        sys.stderr.write(f"Command: {' '.join(e.cmd)}\n")
        sys.stderr.write(f"Stderr: {e.stderr.decode().strip()}\n")
        sys.exit(1)
    except Exception as e:
        sys.stderr.write(f"An unexpected error occurred: {e}\n")
        sys.exit(1)
    finally:
        if os.path.exists(screenshot_path):
            os.remove(screenshot_path)

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-21 01:36:29
import subprocess
import argparse
import os
import sys
import tempfile
import platform

try:
    from PIL import ImageGrab
except ImportError:
    ImageGrab = None

try:
    import pyperclip
    PYPERCLIP_AVAILABLE = True
except ImportError:
    pyperclip = None
    PYPERCLIP_AVAILABLE = False

def check_command_exists(cmd):
    """Checks if a command exists in the system's PATH."""
    try:
        # Use shutil.which for a more robust check if available, otherwise fall back to system commands
        import shutil
        if shutil.which(cmd):
            return True
    except ImportError:
        pass # shutil not available, fall back to subprocess

    # Fallback for systems without shutil or if shutil.which fails
    try:
        if platform.system() == "Windows":
            subprocess.run(["where", cmd], check=True, capture_output=True, text=True)
        else:
            subprocess.run(["which", cmd], check=True, capture_output=True, text=True)
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        return False

def take_screenshot(output_path):
    """
    Captures a screenshot using system-specific tools.
    For Linux/macOS, it attempts an interactive selection.
    For Windows, it captures the entire screen.
    """
    system = platform.system()
    try:
        if system == "Linux":
            if not check_command_exists("scrot"):
                raise RuntimeError("scrot not found. Please install it (e.g., sudo apt-get install scrot).")
            # -s: interactive selection, -o: output file, -z: silent
            subprocess.run(["scrot", "-s", "-o", output_path, "-z"], check=True, capture_output=True)
            # scrot returns 0 even if selection is cancelled, so check if file was created
            if not os.path.exists(output_path) or os.path.getsize(output_path) == 0:
                raise RuntimeError("Screenshot cancelled or failed to capture.")
        elif system == "Darwin": # macOS
            if not check_command_exists("screencapture"):
                raise RuntimeError("screencapture not found. This is a macOS built-in tool, something is wrong.")
            # -i: interactive selection, -s: selection mode (not full screen)
            subprocess.run(["screencapture", "-i", "-s", output_path], check=True, capture_output=True)
        elif system == "Windows":
            if ImageGrab is None:
                raise RuntimeError("Pillow (PIL) is required for screenshots on Windows. Please install it (pip install Pillow).")
            # On Windows, ImageGrab captures the entire screen. Interactive selection is not
            # easily available via built-in command-line tools without complex PowerShell or GUI interaction.
            img = ImageGrab.grab()
            img.save(output_path)
        else:
            raise RuntimeError(f"Unsupported operating system: {system}")
    except subprocess.CalledProcessError as e:
        raise RuntimeError(f"Screenshot command failed: {e.stderr.decode().strip()}")
    except FileNotFoundError:
        raise RuntimeError(f"Screenshot command not found for {system}. Please ensure it's installed and in your PATH.")

def perform_ocr(image_path, lang="eng"):
    """
    Performs OCR on the given image using Tesseract.
    """
    if not check_command_exists("tesseract"):
        raise RuntimeError("Tesseract-OCR not found. Please install it and ensure it's in your PATH.")

    try:
        # Tesseract command: tesseract <image_path> stdout -l <lang>
        result = subprocess.run(["tesseract", image_path, "stdout", "-l", lang],
                                check=True, capture_output=True, text=True, encoding='utf-8')
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        # Tesseract often outputs errors to stderr even on success, or if language data is missing.
        # If stdout is empty, it's likely a real failure.
        if not e.stdout.strip():
            raise RuntimeError(f"OCR failed. Tesseract error: {e.stderr.strip()}")
        return e.stdout.strip() # Return partial output if available
    except FileNotFoundError:
        raise RuntimeError("Tesseract command not found. Please ensure it's installed and in your PATH.")

def copy_to_clipboard(text):
    """Copies text to the system clipboard."""
    if PYPERCLIP_AVAILABLE:
        try:
            pyperclip.copy(text)
            return True
        except pyperclip.PyperclipException as e:
            sys.stderr.write(f"Warning: Could not copy to clipboard: {e}\n")
            return False
    else:
        sys.stderr.write("Warning: pyperclip not installed. Cannot copy to clipboard. (pip install pyperclip)\n")
        return False

def open_image(image_path):
    """Opens the image using the default system viewer."""
    system = platform.system()
    try:
        if system == "Windows":
            os.startfile(image_path)
        elif system == "Darwin": # macOS
            subprocess.run(["open", image_path], check=True)
        elif system == "Linux":
            subprocess.run(["xdg-open", image_path], check=True)
        else:
            sys.stderr.write(f"Warning: Cannot open image on unsupported OS: {system}\n")
    except FileNotFoundError:
        sys.stderr.write(f"Warning: Default image viewer command not found for {system}.\n")
    except subprocess.CalledProcessError as e:
        sys.stderr.write(f"Warning: Failed to open image: {e}\n")

def main():
    parser = argparse.ArgumentParser(
        description="Command-line screenshot OCR tool using system utilities."
    )
    parser.add_argument(
        "-o", "--output",
        help="Path to save the OCR'd text to a file."
    )
    parser.add_argument(
        "-l", "--lang",
        default="eng",
        help="Language for OCR (e.g., 'eng', 'fra', 'spa'). Requires Tesseract language packs."
    )
    parser.add_argument(
        "-c", "--clipboard",
        action="store_true",
        help="Copy the OCR'd text to the system clipboard."
    )
    parser.add_argument(
        "-v", "--view",
        action="store_true",
        help="Open the captured screenshot image after OCR."
    )

    args = parser.parse_args()

    temp_image_file = None
    try:
        # Create a temporary file for the screenshot
        # NamedTemporaryFile ensures unique name and handles cleanup on close/delete
        with tempfile.NamedTemporaryFile(suffix=".png", delete=False) as tmp:
            temp_image_file = tmp.name

        sys.stdout.write("Taking screenshot...\n")
        take_screenshot(temp_image_file)
        sys.stdout.write(f"Screenshot saved to {temp_image_file}\n")

        sys.stdout.write("Performing OCR...\n")
        ocr_text = perform_ocr(temp_image_file, args.lang)

        sys.stdout.write("\n--- OCR Result ---\n")
        sys.stdout.write(ocr_text)
        sys.stdout.write("\n------------------\n")

        if args.output:
            try:
                with open(args.output, "w", encoding="utf-8") as f:
                    f.write(ocr_text)
                sys.stdout.write(f"OCR text saved to {args.output}\n")
            except IOError as e:
                sys.stderr.write(f"Error: Could not write to output file {args.output}: {e}\n")

        if args.clipboard:
            if copy_to_clipboard(ocr_text):
                sys.stdout.write("OCR text copied to clipboard.\n")

        if args.view:
            open_image(temp_image_file)

    except RuntimeError as e:
        sys.stderr.write(f"Error: {e}\n")
        sys.exit(1)
    except Exception as e:
        sys.stderr.write(f"An unexpected error occurred: {e}\n")
        sys.exit(1)
    finally:
        if temp_image_file and os.path.exists(temp_image_file):
            os.remove(temp_image_file)
            sys.stdout.write(f"Cleaned up temporary screenshot file: {temp_image_file}\n")

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-21 01:37:26
import subprocess
import platform
import os
import argparse
import tempfile
from PIL import Image

try:
    import pytesseract
except ImportError:
    # This error will be caught by check_dependencies, but good to have a direct message
    pass

def get_screenshot_command(output_path, interactive_selection=False):
    """
    Returns the appropriate system command for taking a screenshot.
    """
    system = platform.system()
    if system == "Linux":
        if interactive_selection:
            # -s: interactive selection, -z: don't open image viewer
            return ["scrot", "-s", "-z", output_path]
        else:
            return ["scrot", "-z", output_path]
    elif system == "Darwin": # macOS
        if interactive_selection:
            # -s: interactive selection
            return ["screencapture", "-s", output_path]
        else:
            return ["screencapture", output_path]
    elif system == "Windows":
        # Windows does not have a simple built-in command-line tool for interactive selection.
        # This PowerShell command captures the full primary screen.
        # Interactive selection is not supported via command-line system tools on Windows.
        powershell_script = f"""
        Add-Type -AssemblyName System.Drawing
        Add-Type -AssemblyName System.Windows.Forms
        $bounds = [System.Windows.Forms.Screen]::PrimaryScreen.Bounds
        $bmp = New-Object System.Drawing.Bitmap($bounds.Width, $bounds.Height)
        $graphics = [System.Drawing.Graphics]::FromImage($bmp)
        $graphics.CopyFromScreen($bounds.Location, [System.Drawing.Point]::Empty, $bounds.Size)
        $bmp.Save('{output_path}', [System.Drawing.Imaging.ImageFormat]::Png)
        """
        return ["powershell", "-Command", powershell_script]
    else:
        raise OSError(f"Unsupported operating system: {system}")

def copy_to_clipboard(text):
    """
    Copies the given text to the system clipboard using OS-specific commands.
    """
    system = platform.system()
    try:
        if system == "Linux":
            p = subprocess.Popen(['xclip', '-selection', 'clipboard'], stdin=subprocess.PIPE)
            p.communicate(input=text.encode('utf-8'))
        elif system == "Darwin": # macOS
            p = subprocess.Popen(['pbcopy'], stdin=subprocess.PIPE)
            p.communicate(input=text.encode('utf-8'))
        elif system == "Windows":
            p = subprocess.Popen(['clip'], stdin=subprocess.PIPE)
            p.communicate(input=text.encode('utf-8'))
        else:
            print(f"Clipboard copy not supported on {system}.")
            return False
        return True
    except FileNotFoundError:
        print(f"Error: Clipboard utility not found for {system}. Please ensure it's installed and in your PATH.")
        if system == "Linux":
            print("Try: sudo apt-get install xclip (Debian/Ubuntu) or sudo yum install xclip (Fedora/RHEL)")
        return False
    except Exception as e:
        print(f"An error occurred while copying to clipboard: {e}")
        return False

def check_dependencies():
    """
    Checks if required system tools and Python libraries are installed.
    """
    system = platform.system()
    missing_deps = []

    # Check for pytesseract Python library
    try:
        import pytesseract
    except ImportError:
        missing_deps.append("Python library 'pytesseract' (pip install pytesseract)")

    # Check for Tesseract OCR executable
    try:
        subprocess.run(["tesseract", "--version"], capture_output=True, check=True)
    except (subprocess.CalledProcessError, FileNotFoundError):
        missing_deps.append("Tesseract OCR engine (https://tesseract-ocr.github.io/tessdoc/Installation.html)")

    # Check for screenshot tool
    if system == "Linux":
        try:
            subprocess.run(["scrot", "--version"], capture_output=True, check=True)
        except (subprocess.CalledProcessError, FileNotFoundError):
            missing_deps.append("scrot (sudo apt-get install scrot or equivalent)")
    elif system == "Darwin": # macOS
        # screencapture is built-in
        pass
    elif system == "Windows":
        # PowerShell is built-in
        pass

    # Check for clipboard tool
    if system == "Linux":
        try:
            subprocess.run(["xclip", "-version"], capture_output=True, check=True)
        except (subprocess.CalledProcessError, FileNotFoundError):
            missing_deps.append("xclip (sudo apt-get install xclip or equivalent)")
    elif system == "Darwin": # macOS
        # pbcopy is built-in
        pass
    elif system == "Windows":
        # clip is built-in
        pass

    if missing_deps:
        print("Error: The following required dependencies are missing or not in your PATH:")
        for dep in missing_deps:
            print(f"- {dep}")
        print("\nPlease install them and ensure they are accessible from your command line.")
        return False
    return True

def main():
    if not check_dependencies():
        exit(1)

    parser = argparse.ArgumentParser(
        description="Capture a screenshot, perform OCR, and output the text."
    )
    parser.add_argument(
        "-c", "--clipboard", action="store_true",
        help="Copy the OCR'd text to the system clipboard."
    )
    parser.add_argument(
        "-o", "--output-file", type=str,
        help="Save the OCR'd text to a specified file."
    )
    parser.add_argument(
        "-l", "--lang", type=str, default="eng",
        help="Language(s) for OCR (e.g., 'eng', 'spa', 'eng+spa'). Requires corresponding Tesseract language packs."
    )
    parser.add_argument(
        "-a", "--area", action="store_true",
        help="Select a specific area of the screen for the screenshot (interactive)."
             "Note: On Windows, this will still capture the full screen as interactive selection is not supported via command-line system tools."
    )

    args = parser.parse_args()

    temp_image_file = None
    try:
        # Create a temporary file for the screenshot
        with tempfile.NamedTemporaryFile(suffix=".png", delete=False) as tmp:
            temp_image_file = tmp.name

        print("Taking screenshot...")
        screenshot_cmd = get_screenshot_command(temp_image_file, args.area)
        
        # Special handling for Windows PowerShell command as it's a single string
        if platform.system() == "Windows":
            result = subprocess.run(screenshot_cmd, shell=True, check=True, capture_output=True)
        else:
            result = subprocess.run(screenshot_cmd, check=True, capture_output=True)

        if result.returncode != 0:
            print(f"Error taking screenshot: {result.stderr.decode().strip()}")
            return

        # Check if the screenshot file was actually created and is not empty
        if not os.path.exists(temp_image_file) or os.path.getsize(temp_image_file) == 0:
            print("Error: Screenshot file was not created or is empty. This might happen if you cancelled the interactive selection.")
            if platform.system() == "Linux" and args.area:
                print("If using 'scrot -s', ensure you select an area.")
            return

        print(f"Screenshot saved to {temp_image_file}. Performing OCR...")

        # Perform OCR using pytesseract
        try:
            image = Image.open(temp_image_file)
            ocr_text = pytesseract.image_to_string(image, lang=args.lang)
        except pytesseract.TesseractNotFoundError:
            print("Error: Tesseract OCR engine not found. Please install it and ensure it's in your PATH.")
            return
        except Exception as e:
            print(f"Error during OCR processing: {e}")
            return

        if not ocr_text.strip():
            print("No text found in the screenshot.")
            return

        print("\n--- OCR Result ---")
        print(ocr_text.strip())
        print("------------------")

        if args.clipboard:
            if copy_to_clipboard(ocr_text):
                print("OCR text copied to clipboard.")
            else:
                print("Failed to copy text to clipboard.")

        if args.output_file:
            try:
                with open(args.output_file, "w", encoding="utf-8") as f:
                    f.write(ocr_text.strip())
                print(f"OCR text saved to {args.output_file}")
            except IOError as e:
                print(f"Error saving to file {args.output_file}: {e}")

    except FileNotFoundError as e:
        print(f"Error: Required command not found. Please ensure it's installed and in your PATH: {e}")
        if "scrot" in str(e) and platform.system() == "Linux":
            print("Try: sudo apt-get install scrot (Debian/Ubuntu) or sudo yum install scrot (Fedora/RHEL)")
        elif "xclip" in str(e) and platform.system() == "Linux":
            print("Try: sudo apt-get install xclip (Debian/Ubuntu) or sudo yum install xclip (Fedora/RHEL)")
        elif "tesseract" in str(e):
            print("Please install Tesseract OCR: https://tesseract-ocr.github.io/tessdoc/Installation.html")
    except subprocess.CalledProcessError as e:
        print(f"Error executing command: {e}")
        print(f"Command: {e.cmd}")
        print(f"Stderr: {e.stderr.decode().strip()}")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
    finally:
        if temp_image_file and os.path.exists(temp_image_file):
            os.remove(temp_image_file)

if __name__ == "__main__":
    main()