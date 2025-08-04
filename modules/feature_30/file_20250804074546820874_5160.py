import os
import sys
import subprocess

def hide_file(filepath):
    if sys.platform == "win32":
        try:
            subprocess.run(f'attrib +h "{filepath}"', shell=True, check=True, capture_output=True)
            return True
        except subprocess.CalledProcessError:
            return False
    else:
        dirname = os.path.dirname(filepath)
        basename = os.path.basename(filepath)
        if basename.startswith('.'):
            return True
        hidden_filepath = os.path.join(dirname, '.' + basename)
        try:
            os.rename(filepath, hidden_filepath)
            return True
        except OSError:
            return False

def unhide_file(filepath):
    if sys.platform == "win32":
        try:
            subprocess.run(f'attrib -h "{filepath}"', shell=True, check=True, capture_output=True)
            return True
        except subprocess.CalledProcessError:
            return False
    else:
        dirname = os.path.dirname(filepath)
        basename = os.path.basename(filepath)
        if not basename.startswith('.'):
            return True
        unhidden_basename = basename[1:]
        unhidden_filepath = os.path.join(dirname, unhidden_basename)
        try:
            os.rename(filepath, unhidden_filepath)
            return True
        except OSError:
            return False

if __name__ == "__main__":
    test_file_name = "example_file.txt"
    
    with open(test_file_name, "w") as f:
        f.write("This is a test file for hiding and unhiding.")

    if hide_file(test_file_name):
        path_after_hide = test_file_name
        if sys.platform != "win32":
            path_after_hide = os.path.join(os.path.dirname(test_file_name), '.' + os.path.basename(test_file_name))
        
        unhide_file(path_after_hide)
        
        if os.path.exists(test_file_name):
            os.remove(test_file_name)
        elif sys.platform != "win32" and os.path.exists(path_after_hide):
            os.remove(path_after_hide)
    else:
        if os.path.exists(test_file_name):
            os.remove(test_file_name)

# Additional implementation at 2025-08-04 07:46:42
import os
import sys
import subprocess

def hide_file_windows(filepath):
    """Hides a file or directory on Windows using attrib +h."""
    try:
        # Use subprocess.run for better control and security than os.system
        # check=True will raise a CalledProcessError if the command returns a non-zero exit code
        subprocess.run(['attrib', '+h', filepath], check=True, capture_output=True, text=True)
        print(f"Successfully hidden: {filepath}")
    except subprocess.CalledProcessError as e:
        print(f"Error hiding '{filepath}': {e.stderr.strip()}")
    except FileNotFoundError:
        print(f"Error: 'attrib' command not found. Ensure it's in your system PATH.")
    except Exception as e:
        print(f"An unexpected error occurred hiding '{filepath}': {e}")

def unhide_file_windows(filepath):
    """Unhides a file or directory on Windows using attrib -h."""
    try:
        subprocess.run(['attrib', '-h', filepath], check=True, capture_output=True, text=True)
        print(f"Successfully unhidden: {filepath}")
    except subprocess.CalledProcessError as e:
        print(f"Error unhiding '{filepath}': {e.stderr.strip()}")
    except FileNotFoundError:
        print(f"Error: 'attrib' command not found. Ensure it's in your system PATH.")
    except Exception as e:
        print(f"An unexpected error occurred unhiding '{filepath}': {e}")

def hide_file_unix(filepath):
    """Hides a file or directory on Unix-like systems by adding a dot prefix."""
    dirname, basename = os.path.split(filepath)
    if not basename.startswith('.'):
        new_filepath = os.path.join(dirname, '.' + basename)
        try:
            os.rename(filepath, new_filepath)
            print(f"Successfully hidden: {filepath} -> {new_filepath}")
        except FileNotFoundError:
            print(f"Error: File not found '{filepath}'")
        except PermissionError:
            print(f"Error: Permission denied to rename '{filepath}'")
        except OSError as e:
            print(f"Error hiding '{filepath}': {e}")
    else:
        print(f"Already hidden (starts with '.'): {filepath}")

def unhide_file_unix(filepath):
    """Unhides a file or directory on Unix-like systems by removing a dot prefix."""
    dirname, basename = os.path.split(filepath)
    # Ensure it's a hidden file we can unhide (not '.' or '..')
    if basename.startswith('.') and len(basename) > 1:
        new_filepath = os.path.join(dirname, basename[1:])
        try:
            os.rename(filepath, new_filepath)
            print(f"Successfully unhidden: {filepath} -> {new_filepath}")
        except FileNotFoundError:
            print(f"Error: File not found '{filepath}'")
        except PermissionError:
            print(f"Error: Permission denied to rename '{filepath}'")
        except OSError as e:
            print(f"Error unhiding '{filepath}': {e}")
    else:
        print(f"Not a hidden file (does not start with '.') or invalid name: {filepath}")

def main():
    if len(sys.argv) < 3:
        print("Usage: python hide_unhide.py <hide|unhide> <file_or_directory1> [file_or_directory2 ...]")
        sys.exit(1)

    action = sys.argv[1].lower()
    files_to_process = sys.argv[2:]

    if action not in ['hide', 'unhide']:
        print("Invalid action. Please use 'hide' or 'unhide'.")
        sys.exit(1)

    is_windows = os.name == 'nt'

    for filepath in files_to_process:
        # Resolve absolute path for robustness
        abs_filepath = os.path.abspath(filepath)

        # Check if the path exists before attempting to modify it
        if not os.path.exists(abs_filepath):
            print(f"Warning: Path does not exist - {abs_filepath}")
            continue

        if action == 'hide':
            if is_windows:
                hide_file_windows(abs_filepath)
            else: # Assume posix for non-windows (Linux, macOS, BSD)
                hide_file_unix(abs_filepath)
        elif action == 'unhide':
            if is_windows:
                unhide_file_windows(abs_filepath)
            else: # Assume posix for non-windows
                unhide_file_unix(abs_filepath)

if __name__ == '__main__':
    main()

# Additional implementation at 2025-08-04 07:48:01
import os
import sys
import subprocess

# Determine the operating system
is_windows = sys.platform.startswith('win')

def is_file_hidden(filepath):
    """
    Checks if a file is hidden.
    On Windows, checks the 'H' attribute using 'attrib'.
    On Unix-like, checks if the filename starts with a dot.
    """
    if not os.path.exists(filepath):
        print(f"Error: File not found at '{filepath}'")
        return False

    if is_windows:
        try:
            # Use subprocess to run 'attrib' command
            # capture_output=True captures stdout and stderr
            # text=True decodes output as text
            # check=True raises CalledProcessError for non-zero exit codes
            # creationflags=subprocess.CREATE_NO_WINDOW prevents a console window from popping up
            result = subprocess.run(
                ['attrib', filepath],
                capture_output=True,
                text=True,
                check=True,
                creationflags=subprocess.CREATE_NO_WINDOW
            )
            # The output for a hidden file typically contains 'H' in the attribute string.
            # Example: "A H        C:\path\to\file.txt"
            # We split the output by the uppercase basename to isolate the attribute string.
            attribute_string = result.stdout.upper().split(os.path.basename(filepath).upper())[0]
            return 'H' in attribute_string
        except subprocess.CalledProcessError as e:
            print(f"Error checking hidden status for '{filepath}': {e}")
            print(f"Stderr: {e.stderr.strip()}")
            return False
        except FileNotFoundError:
            print("Error: 'attrib' command not found. Make sure it's in your system PATH.")
            return False
        except Exception as e:
            print(f"An unexpected error occurred while checking hidden status: {e}")
            return False
    else:  # Unix-like system
        # A file is considered hidden if its name starts with a dot '.'
        return os.path.basename(filepath).startswith('.')

def hide_file(filepath):
    """
    Hides a file.
    On Windows, uses 'attrib +h'.
    On Unix-like, renames the file to start with a dot.
    """
    if not os.path.exists(filepath):
        print(f"Error: File not found at '{filepath}'")
        return False

    if is_file_hidden(filepath):
        print(f"'{filepath}' is already hidden.")
        return True

    try:
        if is_windows:
            subprocess.run(
                ['attrib', '+h', filepath],
                check=True,
                creationflags=subprocess.CREATE_NO_WINDOW
            )
            print(f"Successfully hid '{filepath}'.")
            return True
        else:  # Unix-like system
            dirname = os.path.dirname(filepath)
            basename = os.path.basename(filepath)
            new_filepath = os.path.join(dirname, '.' + basename)
            os.rename(filepath, new_filepath)
            print(f"Successfully hid '{filepath}' as '{os.path.basename(new_filepath)}'.")
            return True
    except FileNotFoundError:
        print(f"Error: File not found or command not accessible for '{filepath}'.")
        return False
    except PermissionError:
        print(f"Error: Permission denied to hide '{filepath}'. Run as administrator/sudo.")
        return False
    except subprocess.CalledProcessError as e:
        print(f"Error hiding '{filepath}': {e}")
        print(f"Stderr: {e.stderr.strip()}")
        return False
    except OSError as e:
        print(f"OS error hiding '{filepath}': {e}")
        return False
    except Exception as e:
        print(f"An unexpected error occurred while hiding '{filepath}': {e}")
        return False

def unhide_file(filepath):
    """
    Unhides a file.
    On Windows, uses 'attrib -h'.
    On Unix-like, renames the file to remove the leading dot.
    """
    if not os.path.exists(filepath):
        print(f"Error: File not found at '{filepath}'")
        return False

    if not is_file_hidden(filepath):
        print(f"'{filepath}' is not hidden (or not hidden by this method).")
        return True

    try:
        if is_windows:
            subprocess.run(
                ['attrib', '-h', filepath],
                check=True,
                creationflags=subprocess.CREATE_NO_WINDOW
            )
            print(f"Successfully unhid '{filepath}'.")
            return True
        else:  # Unix-like system
            dirname = os.path.dirname(filepath)
            basename = os.path.basename(filepath)
            if not basename.startswith('.'):
                print(f"'{filepath}' is not hidden (does not start with '.').")
                return True
            new_filepath = os.path.join(dirname, basename[1:])
            os.rename(filepath, new_filepath)
            print(f"Successfully unhid '{filepath}' as '{os.path.basename(new_filepath)}'.")
            return True
    except FileNotFoundError:
        print(f"Error: File not found or command not accessible for '{filepath}'.")
        return False
    except PermissionError:
        print(f"Error: Permission denied to unhide '{filepath}'. Run as administrator/sudo.")
        return False
    except subprocess.CalledProcessError as e:
        print(f"Error unhiding '{filepath}': {e}")
        print(f"Stderr: {e.stderr.strip()}")
        return False
    except OSError as e:
        print(f"OS error unhiding '{filepath}': {e}")
        return False
    except Exception as e:
        print(f"An unexpected error occurred while unhiding '{filepath}': {e}")
        return False

def list_hidden_files_in_directory(directory_path):
    """
    Lists all hidden files in the specified directory.
    """
    if not os.path.isdir(directory_path):
        print(f"Error: Directory not found at '{directory_path}'")
        return []

    hidden_files = []
    try:
        for item_name in os.listdir(directory_path):
            full_path = os.path.join(directory_path, item_name)
            if os.path.isfile(full_path): # Only check files, not directories
                if is_file_hidden(full_path):
                    hidden_files.append(full_path)
        return hidden_files
    except PermissionError:
        print(f"Error: Permission denied to list files in '{directory_path}'.")
        return []
    except Exception as e:
        print(f"An unexpected error occurred while listing files: {e}")
        return []

def get_valid_filepath(prompt):
    """Helper function to get a file path from user and validate it."""
    while True:
        filepath = input(prompt).strip()
        if os.path.exists(filepath):
            if os.path.isfile(filepath):
                return filepath
            else:
                print("Error: The path points to a directory, not a file. Please enter a file path.")
        else:
            print("Error: File does not exist. Please enter a valid file path.")

def get_valid_directory_path(prompt):
    """Helper function to get a directory path from user and validate it."""
    while True:
        dirpath = input(prompt).strip()
        if os.path.isdir(dirpath):
            return dirpath
        else:
            print("Error: Directory does not exist. Please enter a valid directory path.")

def main():
    """Main function to run the file hiding/unhiding program."""
    while True:
        print("\n--- File Hider/Unhider ---")
        print("1. Hide a file")
        print("2. Unhide a file")
        print("3. Check if a file is hidden")
        print("4. List hidden files in a directory")
        print("5. Exit")

        choice = input("Enter your choice (1-5): ").strip()

        if choice == '1':
            file_to_hide = get_valid_filepath("Enter the full path of the file to hide: ")
            if file_to_hide:
                hide_file(file_to_hide)
        elif choice == '2':
            file_to_unhide = get_valid_filepath("Enter the full path of the file to unhide: ")
            if file_to_unhide:
                unhide_file(file_to_unhide)
        elif choice == '3':
            file_to_check = get_valid_filepath("Enter the full path of the file to check: ")
            if file_to_check:
                if is_file_hidden(file_to_check):
                    print(f"'{file_to_check}' IS hidden.")
                else:
                    print(f"'{file_to_check}' IS NOT hidden.")
        elif choice == '4':
            dir_to_list = get_valid_directory_path("Enter the full path of the directory to list hidden files: ")
            if dir_to_list:
                hidden_files = list_hidden_files_in_directory(dir_to_list)
                if hidden_files:
                    print(f"\nHidden files in '{dir_to_list}':")
                    for f in hidden_files:
                        print(f"- {f}")
                else:
                    print(f"No hidden files found in '{dir_to_list}'.")
        elif choice == '5':
            print("Exiting program. Goodbye!")
            break
        else:
            print("Invalid choice. Please enter a number between 1 and 5.")

if __name__ == "__main__":
    main()

# Additional implementation at 2025-08-04 07:49:04
import os
import sys
import subprocess
import argparse

def hide_file(filepath):
    if sys.platform == "win32":
        try:
            subprocess.run(["attrib", "+h", filepath], check=True, creationflags=subprocess.CREATE_NO_WINDOW)
        except subprocess.CalledProcessError:
            pass
    else:
        directory, filename = os.path.split(filepath)
        if not filename.startswith("."):
            new_filepath = os.path.join(directory, "." + filename)
            try:
                os.rename(filepath, new_filepath)
            except OSError:
                pass

def unhide_file(filepath):
    if sys.platform == "win32":
        try:
            subprocess.run(["attrib", "-h", filepath], check=True, creationflags=subprocess.CREATE_NO_WINDOW)
        except subprocess.CalledProcessError:
            pass
    else:
        directory, filename = os.path.split(filepath)
        if filename.startswith("."):
            new_filepath = os.path.join(directory, filename[1:])
            try:
                os.rename(filepath, new_filepath)
            except OSError:
                pass

def is_file_hidden(filepath):
    if not os.path.exists(filepath):
        return False

    if sys.platform == "win32":
        try:
            result = subprocess.run(["attrib", filepath], capture_output=True, text=True, check=True, creationflags=subprocess.CREATE_NO_WINDOW)
            return " H " in result.stdout.upper()
        except (subprocess.CalledProcessError, FileNotFoundError):
            return False
    else:
        filename = os.path.basename(filepath)
        return filename.startswith(".")

def list_files(directory=".", include_hidden=False):
    files = []
    for entry in os.listdir(directory):
        full_path = os.path.join(directory, entry)
        if os.path.isfile(full_path):
            if include_hidden:
                files.append(full_path)
            else:
                if not is_file_hidden(full_path):
                    files.append(full_path)
    return files

def list_hidden_files(directory="."):
    hidden_files = []
    for entry in os.listdir(directory):
        full_path = os.path.join(directory, entry)
        if os.path.isfile(full_path) and is_file_hidden(full_path):
            hidden_files.append(full_path)
    return hidden_files

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("action", choices=["hide", "unhide", "list", "listhidden"], help="hide|unhide|list|listhidden")
    parser.add_argument("path", nargs="?", default=".", help="file or directory path")

    args = parser.parse_args()

    if args.action == "hide":
        if not args.path or not os.path.exists(args.path) or os.path.isdir(args.path):
            print("Error: File path required.")
            return
        hide_file(args.path)
        if is_file_hidden(args.path):
            print(f"Hidden: {args.path}")
        else:
            print(f"Failed to hide: {args.path}")

    elif args.action == "unhide":
        if not args.path or not os.path.exists(args.path) or os.path.isdir(args.path):
            print("Error: File path required.")
            return
        unhide_file(args.path)
        if not is_file_hidden(args.path):
            print(f"Unhidden: {args.path}")
        else:
            print(f"Failed to unhide: {args.path}")

    elif args.action == "list":
        if os.path.isfile(args.path):
            print("Error: Directory path required.")
            return
        print(f"Non-hidden files in '{args.path}':")
        files = list_files(args.path, include_hidden=False)
        if not files:
            print("  None.")
        for f in files:
            print(f"  {f}")

    elif args.action == "listhidden":
        if os.path.isfile(args.path):
            print("Error: Directory path required.")
            return
        print(f"Hidden files in '{args.path}':")
        hidden_files = list_hidden_files(args.path)
        if not hidden_files:
            print("  None.")
        for f in hidden_files:
            print(f"  {f}")

if __name__ == "__main__":
    main()