import os
import sys
import tempfile
import subprocess
import platform

def main():
    screenshot_path = None
    ocr_output_base = None
    try:
        with tempfile.NamedTemporaryFile(suffix=".png", delete=False) as temp_screenshot_file:
            screenshot_path = temp_screenshot_file.name

        current_os = platform.system()
        screenshot_command = []

        if current_os == "Linux":
            screenshot_command = ["scrot", screenshot_path]
        elif current_os == "Darwin":
            screenshot_command = ["screencapture", screenshot_path]
        elif current_os == "Windows":
            screenshot_command = ["scrot", screenshot_path]
        else:
            sys.exit(1)

        subprocess.run(screenshot_command, check=True)

        ocr_output_base = tempfile.NamedTemporaryFile(delete=False).name
        ocr_command = ["tesseract", screenshot_path, ocr_output_base, "-l", "eng"]
        subprocess.run(ocr_command, check=True)

        with open(ocr_output_base + ".txt", "r") as f:
            sys.stdout.write(f.read())

    except Exception:
        sys.exit(1)
    finally:
        if screenshot_path and os.path.exists(screenshot_path):
            os.remove(screenshot_path)
        if ocr_output_base:
            if os.path.exists(ocr_output_base + ".txt"):
                os.remove(ocr_output_base + ".txt")
            if os.path.exists(ocr_output_base):
                os.remove(ocr_output_base)

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-18 02:17:42
import os
import subprocess
import tempfile
import argparse
import platform
import sys
from PIL import Image

try:
    import pytesseract
except ImportError:
    print("Error: pytesseract not found. Please install it using 'pip install pytesseract'.", file=sys.stderr)
    print("Also, ensure Tesseract OCR is installed and in your PATH.", file=sys.stderr)
    sys.exit(1)

try:
    import pyperclip
except ImportError:
    print("Warning: pyperclip not found. Clipboard functionality will be disabled.", file=sys.stderr)
    pyperclip = None

if platform.system() == "Windows":
    try:
        import pyautogui
    except ImportError:
        print("Error: pyautogui not found. Please install it using 'pip install pyautogui'.", file=sys.stderr)
        print("pyautogui is required for screenshot functionality on Windows.", file=sys.stderr)
        sys.exit(1)

def take_screenshot(output_path: str, region: bool = False):
    system = platform.system()

    if system == "Linux":
        try:
            if region:
                print("Please select a region on your screen...")
                subprocess.run(["scrot", "-s", output_path], check=True)
            else:
                subprocess.run(["scrot", output_path], check=True)
        except FileNotFoundError:
            print("Error: 'scrot' command not found.", file=sys.stderr)
            print("Please install scrot (e.g., 'sudo apt-get install scrot' on Debian/Ubuntu).", file=sys.stderr)
            sys.exit(1)
        except subprocess.CalledProcessError:
            print("Error: Failed to take screenshot with scrot. Ensure X server is running and scrot is configured.", file=sys.stderr)
            sys.exit(1)
    elif system == "Darwin":
        try:
            if region:
                print("Please select a region on your screen...")
                subprocess.run(["screencapture", "-i", output_path], check=True)
            else:
                subprocess.run(["screencapture", output_path], check=True)
        except FileNotFoundError:
            print("Error: 'screencapture' command not found.", file=sys.stderr)
            print("This command should be available on macOS by default.", file=sys.stderr)
            sys.exit(1)
        except subprocess.CalledProcessError:
            print("Error: Failed to take screenshot with screencapture.", file=sys.stderr)
            sys.exit(1)
    elif system == "Windows":
        if region:
            print("Warning: Interactive region selection is not natively supported for Windows via command-line system tools.")
            print("Taking a full screen screenshot instead.")
        try:
            pyautogui.screenshot(output_path)
        except Exception as e:
            print(f"Error: Failed to take screenshot with pyautogui: {e}", file=sys.stderr)
            sys.exit(1)
    else:
        print(f"Error: Unsupported operating system: {system}", file=sys.stderr)
        sys.exit(1)

def perform_ocr(image_path: str, lang: str = 'eng') -> str:
    try:
        text = pytesseract.image_to_string(Image.open(image_path), lang=lang)
        return text
    except pytesseract.TesseractNotFoundError:
        print("Error: Tesseract OCR engine not found.", file=sys.stderr)
        print("Please install Tesseract OCR and ensure it's in your system's PATH.", file=sys.stderr)
        print("Download from: https://tesseract-ocr.github.io/tessdoc/Installation.html", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"Error during OCR processing: {e}", file=sys.stderr)
        sys.exit(1)

def main():
    parser = argparse.ArgumentParser(
        description="Command-line screenshot OCR tool.",
        formatter_class=argparse.RawTextHelpFormatter
    )
    parser.add_argument(
        "-r", "--region",
        action="store_true",
        help="Select a specific region of the screen for OCR.\n"
             "(Interactive selection available on Linux/macOS. Full screen on Windows.)"
    )
    parser.add_argument(
        "-l", "--lang",
        type=str,
        default="eng",
        help="Language for OCR (e.g., 'eng', 'fra', 'deu').\n"
             "Requires corresponding Tesseract language packs to be installed."
    )
    parser.add_argument(
        "-c", "--clipboard",
        action="store_true",
        help="Copy the OCR result to the system clipboard."
    )
    parser.add_argument(
        "-o", "--output",
        type=str,
        help="Save the OCR result to a specified text file."
    )

    args = parser.parse_args()

    temp_screenshot_path = None
    try:
        with tempfile.NamedTemporaryFile(suffix=".png", delete=False) as temp_file:
            temp_screenshot_path = temp_file.name

        take_screenshot(temp_screenshot_path, args.region)
        ocr_text = perform_ocr(temp_screenshot_path, args.lang)

        print("\n--- OCR Result ---")
        print(ocr_text.strip())
        print("------------------\n")

        if args.clipboard:
            if pyperclip:
                try:
                    pyperclip.copy(ocr_text)
                    print("OCR result copied to clipboard.")
                except pyperclip.PyperclipException as e:
                    print(f"Warning: Could not copy to clipboard: {e}", file=sys.stderr)
                    print("Ensure you have a clipboard tool installed (e.g., xclip/xsel on Linux).", file=sys.stderr)
            else:
                print("Warning: pyperclip is not installed. Cannot copy to clipboard.", file=sys.stderr)

        if args.output:
            try:
                with open(args.output, "w", encoding="utf-8") as f:
                    f.write(ocr_text)
                print(f"OCR result saved to '{args.output}'.")
            except IOError as e:
                print(f"Error: Could not save output to file '{args.output}': {e}", file=sys.stderr)

    finally:
        if temp_screenshot_path and os.path.exists(temp_screenshot_path):
            os.remove(temp_screenshot_path)

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-18 02:18:27
import argparse
import os
import sys
import time
import subprocess
import tempfile
from PIL import Image, ImageGrab
import pytesseract

def get_clipboard_content():
    """Retrieves text content from the clipboard."""
    try:
        import pyperclip
        return pyperclip.paste()
    except ImportError:
        print("Warning: pyperclip not found. Cannot access clipboard.", file=sys.stderr)
        return None
    except Exception as e:
        print(f"Error accessing clipboard: {e}", file=sys.stderr)
        return None

def copy_to_clipboard(text):
    """Copies text content to the clipboard."""
    try:
        import pyperclip
        pyperclip.copy(text)
        print("OCR text copied to clipboard.")
    except ImportError:
        print("Warning: pyperclip not found. Cannot copy to clipboard.", file=sys.stderr)
    except Exception as e:
        print(f"Error copying to clipboard: {e}", file=sys.stderr)

def capture_screenshot(full_screen=False, region=None, delay=0):
    """
    Captures a screenshot using platform-specific tools or Pillow.
    Returns the path to the captured image file.
    """
    if delay > 0:
        print(f"Waiting {delay} seconds before taking screenshot...")
        time.sleep(delay)

    temp_file_path = None
    try:
        # Create a temporary file for the screenshot
        fd, temp_file_path = tempfile.mkstemp(suffix=".png")
        os.close(fd) # Close the file descriptor immediately

        if sys.platform == "darwin":  # macOS
            if region:
                # Interactive region selection
                print("Select a region on your screen...")
                subprocess.run(["screencapture", "-i", "-s", temp_file_path], check=True)
            else:
                # Full screen
                subprocess.run(["screencapture", temp_file_path], check=True)
        elif sys.platform.startswith("linux"):  # Linux
            if region:
                # Interactive region selection
                print("Select a region on your screen...")
                # gnome-screenshot -a for area, -f for file
                subprocess.run(["gnome-screenshot", "-a", "-f", temp_file_path], check=True)
            else:
                # Full screen
                subprocess.run(["gnome-screenshot", "-f", temp_file_path], check=True)
        elif sys.platform == "win32":  # Windows
            if region:
                # For Windows, interactive region selection via command-line is not straightforward.
                # Prompt user for coordinates or use Pillow's grab with bbox.
                print("Interactive region selection is not directly supported via command-line on Windows.")
                print("Please provide coordinates for the region (x1, y1, x2, y2).")
                try:
                    coords_str = input("Enter coordinates (e.g., 0,0,800,600): ")
                    x1, y1, x2, y2 = map(int, coords_str.split(','))
                    screenshot = ImageGrab.grab(bbox=(x1, y1, x2, y2))
                    screenshot.save(temp_file_path)
                except ValueError:
                    print("Invalid coordinates. Aborting screenshot.", file=sys.stderr)
                    return None
                except Exception as e:
                    print(f"Error capturing region: {e}", file=sys.stderr)
                    return None
            else:
                # Full screen using Pillow
                screenshot = ImageGrab.grab()
                screenshot.save(temp_file_path)
        else:
            print(f"Unsupported operating system: {sys.platform}", file=sys.stderr)
            return None

        if not os.path.exists(temp_file_path) or os.path.getsize(temp_file_path) == 0:
            print("Screenshot capture failed or resulted in an empty file.", file=sys.stderr)
            return None

        return temp_file_path

    except FileNotFoundError as e:
        print(f"Error: Screenshot tool not found. Please ensure it's installed and in your PATH.", file=sys.stderr)
        if sys.platform == "darwin":
            print("For macOS, 'screencapture' is built-in.", file=sys.stderr)
        elif sys.platform.startswith("linux"):
            print("For Linux, try installing 'gnome-screenshot' (e.g., 'sudo apt install gnome-screenshot').", file=sys.stderr)
        return None
    except subprocess.CalledProcessError as e:
        print(f"Error during screenshot capture: {e}", file=sys.stderr)
        print(f"Stderr: {e.stderr.decode() if e.stderr else 'N/A'}", file=sys.stderr)
        return None
    except Exception as e:
        print(f"An unexpected error occurred during screenshot capture: {e}", file=sys.stderr)
        return None

def perform_ocr(image_path, lang='eng'):
    """Performs OCR on the given image file."""
    try:
        image = Image.open(image_path)
        text = pytesseract.image_to_string(image, lang=lang)
        return text
    except FileNotFoundError:
        print("Error: Tesseract is not installed or not in your PATH.", file=sys.stderr)
        print("Please install Tesseract OCR engine (e.g., 'sudo apt install tesseract-ocr').", file=sys.stderr)
        print("If installed, ensure its executable is in your system's PATH or set pytesseract.pytesseract.tesseract_cmd.", file=sys.stderr)
        return None
    except Exception as e:
        print(f"Error during OCR processing: {e}", file=sys.stderr)
        return None

def main():
    parser = argparse.ArgumentParser(
        description="Command-line screenshot OCR tool.",
        formatter_class=argparse.RawTextHelpFormatter
    )
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument(
        "--full",
        action="store_true",
        help="Capture a full-screen screenshot."
    )
    group.add_argument(
        "--region",
        action="store_true",
        help="Capture a user-selected region screenshot.\n"
             "  - macOS/Linux: Interactive selection.\n"
             "  - Windows: Prompts for coordinates (x1,y1,x2,y2)."
    )
    parser.add_argument(
        "--output-file",
        "-o",
        type=str,
        help="Save the OCR'd text to a specified file."
    )
    parser.add_argument(
        "--clipboard",
        "-c",
        action="store_true",
        help="Copy the OCR'd text to the system clipboard."
    )
    parser.add_argument(
        "--lang",
        "-l",
        type=str,
        default="eng",
        help="Specify the OCR language (e.g., 'eng', 'deu', 'fra').\n"
             "Requires corresponding Tesseract language packs."
    )
    parser.add_argument(
        "--delay",
        "-d",
        type=int,
        default=0,
        help="Delay in seconds before taking the screenshot."
    )

    args = parser.parse_args()

    screenshot_path = None
    try:
        screenshot_path = capture_screenshot(
            full_screen=args.full,
            region=args.region,
            delay=args.delay
        )

        if not screenshot_path:
            print("Screenshot capture failed. Exiting.", file=sys.stderr)
            sys.exit(1)

        print(f"Screenshot captured: {screenshot_path}")

        ocr_text = perform_ocr(screenshot_path, lang=args.lang)

        if ocr_text is None:
            print("OCR failed. Exiting.", file=sys.stderr)
            sys.exit(1)

        if not ocr_text.strip():
            print("No text found in the screenshot.")
        else:
            print("\n--- OCR Result ---")
            print(ocr_text.strip())
            print("------------------")

            if args.output_file:
                try:
                    with open(args.output_file, "w", encoding="utf-8") as f:
                        f.write(ocr_text.strip())
                    print(f"OCR text saved to {args.output_file}")
                except IOError as e:
                    print(f"Error saving to file {args.output_file}: {e}", file=sys.stderr)

            if args.clipboard:
                copy_to_clipboard(ocr_text.strip())

    finally:
        # Clean up the temporary screenshot file
        if screenshot_path and os.path.exists(screenshot_path):
            try:
                os.remove(screenshot_path)
                print(f"Cleaned up temporary file: {screenshot_path}")
            except OSError as e:
                print(f"Error removing temporary file {screenshot_path}: {e}", file=sys.stderr)

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-18 02:19:53
import os
import subprocess
import argparse
import tempfile
import platform
from PIL import ImageGrab

def capture_screenshot(output_path):
    system = platform.system()

    if system == "Darwin":
        try:
            subprocess.run(["screencapture", "-x", output_path], check=True, capture_output=True)
        except subprocess.CalledProcessError as e:
            raise RuntimeError(f"macOS screenshot failed: {e.stderr.decode()}")
        except FileNotFoundError:
            raise RuntimeError("screencapture command not found. Is macOS installed?")
    elif system == "Linux":
        try:
            subprocess.run(["scrot", output_path], check=True, capture_output=True)
        except (subprocess.CalledProcessError, FileNotFoundError):
            try:
                subprocess.run(["gnome-screenshot", "-f", output_path], check=True, capture_output=True)
            except subprocess.CalledProcessError as e:
                raise RuntimeError(f"Linux screenshot failed (gnome-screenshot): {e.stderr.decode()}")
            except FileNotFoundError:
                raise RuntimeError("Neither 'scrot' nor 'gnome-screenshot' found. Please install one of them.")
    elif system == "Windows":
        try:
            screenshot = ImageGrab.grab()
            screenshot.save(output_path)
        except Exception as e:
            raise RuntimeError(f"Windows screenshot failed (Pillow ImageGrab): {e}")
    else:
        raise NotImplementedError(f"Screenshotting not implemented for {system}")

def perform_ocr(image_path, lang='eng'):
    temp_output_base = tempfile.NamedTemporaryFile(delete=False, suffix=".txt").name
    temp_output_base_no_ext = temp_output_base[:-4]

    try:
        subprocess.run(
            ["tesseract", image_path, temp_output_base_no_ext, "-l", lang],
            check=True,
            capture_output=True
        )
        
        with open(temp_output_base, 'r', encoding='utf-8') as f:
            ocr_text = f.read()
        return ocr_text
    except subprocess.CalledProcessError as e:
        raise RuntimeError(f"Tesseract OCR failed: {e.stderr.decode()}")
    except FileNotFoundError:
        raise RuntimeError("Tesseract command not found. Please install Tesseract OCR.")
    finally:
        if os.path.exists(temp_output_base):
            os.remove(temp_output_base)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Command-line screenshot OCR tool using system utilities."
    )
    parser.add_argument(
        "-l", "--lang",
        default="eng",
        help="Language for OCR (e.g., 'eng', 'spa', 'fra'). Requires Tesseract language packs."
    )
    parser.add_argument(
        "-o", "--output",
        help="Path to save the OCR text. If not specified, text is printed to stdout."
    )
    parser.add_argument(
        "--keep-temp",
        action="store_true",
        help="Keep the temporary screenshot file for debugging."
    )

    args = parser.parse_args()

    temp_screenshot_path = None
    try:
        with tempfile.NamedTemporaryFile(suffix=".png", delete=False) as temp_file:
            temp_screenshot_path = temp_file.name

        print(f"Capturing screenshot to {temp_screenshot_path}...")
        capture_screenshot(temp_screenshot_path)
        print("Screenshot captured. Performing OCR...")

        ocr_result = perform_ocr(temp_screenshot_path, args.lang)

        if args.output:
            with open(args.output, 'w', encoding='utf-8') as f:
                f.write(ocr_result)
            print(f"OCR text saved to {args.output}")
        else:
            print("\n--- OCR Result ---")
            print(ocr_result)
            print("------------------")

    except RuntimeError as e:
        print(f"Error: {e}")
    except NotImplementedError as e:
        print(f"Error: {e}")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
    finally:
        if temp_screenshot_path and os.path.exists(temp_screenshot_path) and not args.keep_temp:
            os.remove(temp_screenshot_path)
            print(f"Temporary screenshot file {temp_screenshot_path} removed.")
        elif args.keep_temp:
            print(f"Temporary screenshot file kept at {temp_screenshot_path}.")