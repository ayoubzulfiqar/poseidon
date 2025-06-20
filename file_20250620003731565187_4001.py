import os
import sys

def hide_file(filepath):
    abs_filepath = os.path.abspath(filepath)
    if not os.path.exists(abs_filepath):
        print(f"Error: File not found at '{abs_filepath}'")
        return False

    try:
        if sys.platform.startswith('win'):
            command = f'attrib +h "{abs_filepath}"'
            result = os.system(command)
            if result == 0:
                print(f"Successfully hid '{abs_filepath}'.")
                return True
            else:
                print(f"Failed to hide '{abs_filepath}'. Command exited with code {result}.")
                return False
        else:
            dirname = os.path.dirname(abs_filepath)
            basename = os.path.basename(abs_filepath)
            if not basename.startswith('.'):
                new_filepath = os.path.join(dirname, '.' + basename)
                os.rename(abs_filepath, new_filepath)
                print(f"Successfully hid '{abs_filepath}' as '{new_filepath}'.")
                return True
            else:
                print(f"'{abs_filepath}' is already hidden.")
                return True
    except OSError as e:
        print(f"Error hiding '{abs_filepath}': {e}")
        return False
    except Exception as e:
        print(f"An unexpected error occurred while hiding '{abs_filepath}': {e}")
        return False

def unhide_file(filepath):
    abs_filepath = os.path.abspath(filepath)
    if not os.path.exists(abs_filepath):
        print(f"Error: File not found at '{abs_filepath}'")
        return False

    try:
        if sys.platform.startswith('win'):
            command = f'attrib -h "{abs_filepath}"'
            result = os.system(command)
            if result == 0:
                print(f"Successfully unhid '{abs_filepath}'.")
                return True
            else:
                print(f"Failed to unhide '{abs_filepath}'. Command exited with code {result}.")
                return False
        else:
            dirname = os.path.dirname(abs_filepath)
            basename = os.path.basename(abs_filepath)
            if basename.startswith('.'):
                new_filepath = os.path.join(dirname, basename[1:])
                os.rename(abs_filepath, new_filepath)
                print(f"Successfully unhid '{abs_filepath}' as '{new_filepath}'.")
                return True
            else:
                print(f"'{abs_filepath}' is not hidden.")
                return True
    except OSError as e:
        print(f"Error unhiding '{abs_filepath}': {e}")
        return False
    except Exception as e:
        print(f"An unexpected error occurred while unhiding '{abs_filepath}': {e}")
        return False

if __name__ == "__main__":
    while True:
        print("\n--- File Hider/Unhider ---")
        print("1. Hide a file")
        print("2. Unhide a file")
        print("3. Exit")
        choice = input("Enter your choice (1/2/3): ")

        if choice == '1':
            file_path = input("Enter the path of the file to hide: ")
            hide_file(file_path)
        elif choice == '2':
            file_path = input("Enter the path of the file to unhide: ")
            unhide_file(file_path)
        elif choice == '3':
            print("Exiting program.")
            break
        else:
            print("Invalid choice. Please enter 1, 2, or 3.")

# Additional implementation at 2025-06-20 00:38:25
import os
import sys
import subprocess

def _is_windows():
    """Checks if the current operating system is Windows."""
    return sys.platform.startswith('win')

def _hide_file_windows(filepath):
    """Hides a file on Windows using the 'attrib +h' command."""
    try:
        # Use subprocess.run for better control and error handling
        # 'attrib +h' sets the hidden attribute
        subprocess.run(['attrib', '+h', filepath], check=True, capture_output=True, text=True)
        print(f"Successfully hid '{filepath}'.")
    except FileNotFoundError:
        print(f"Error: 'attrib' command not found. This should not happen on Windows.")
    except subprocess.CalledProcessError as e:
        print(f"Error hiding '{filepath}': {e.stderr.strip()}")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")

def _unhide_file_windows(filepath):
    """Unhides a file on Windows using the 'attrib -h' command."""
    try:
        # 'attrib -h' removes the hidden attribute
        subprocess.run(['attrib', '-h', filepath], check=True, capture_output=True, text=True)
        print(f"Successfully unhid '{filepath}'.")
    except FileNotFoundError:
        print(f"Error: 'attrib' command not found. This should not happen on Windows.")
    except subprocess.CalledProcessError as e:
        print(f"Error unhiding '{filepath}': {e.stderr.strip()}")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")

def _hide_file_unix(filepath):
    """Hides a file on Unix-like systems by renaming it with a dot prefix."""
    abs_filepath = os.path.abspath(filepath)
    dirname, basename = os.path.split(abs_filepath)

    if basename.startswith('.'):
        print(f"'{filepath}' is already hidden.")
        return

    hidden_filepath = os.path.join(dirname, '.' + basename)

    try:
        os.rename(abs_filepath, hidden_filepath)
        print(f"Successfully hid '{filepath}' as '{os.path.basename(hidden_filepath)}'.")
    except FileNotFoundError:
        print(f"Error: File '{filepath}' not found.")
    except OSError as e:
        print(f"Error hiding '{filepath}': {e}")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")

def _unhide_file_unix(filepath):
    """Unhides a file on Unix-like systems by removing its dot prefix."""
    abs_filepath = os.path.abspath(filepath)
    dirname, basename = os.path.split(abs_filepath)

    # Case 1: User provided the hidden path directly (e.g., .myfile.txt)
    if basename.startswith('.'):
        unhidden_basename = basename[1:]
        unhidden_filepath = os.path.join(dirname, unhidden_basename)
        
        if not os.path.exists(abs_filepath):
            print(f"Error: Hidden file '{filepath}' not found.")
            return

        try:
            os.rename(abs_filepath, unhidden_filepath)
            print(f"Successfully unhid '{filepath}' to '{unhidden_basename}'.")
        except OSError as e:
            print(f"Error unhiding '{filepath}': {e}")
        except Exception as e:
            print(f"An unexpected error occurred: {e}")
    # Case 2: User provided the unhidden path (e.g., myfile.txt), but it might be hidden
    else:
        hidden_version_path = os.path.join(dirname, '.' + basename)
        
        if os.path.exists(abs_filepath):
            print(f"'{filepath}' is already unhidden.")
            return
        elif os.path.exists(hidden_version_path):
            try:
                os.rename(hidden_version_path, abs_filepath)
                print(f"Successfully unhid '{os.path.basename(hidden_version_path)}' to '{basename}'.")
            except OSError as e:
                print(f"Error unhiding '{filepath}': {e}")
            except Exception as e:
                print(f"An unexpected error occurred: {e}")
        else:
            print(f"Error: File '{filepath}' or its hidden version ('.{basename}') not found.")


def hide_file(filepath):
    """Hides a file based on the operating system."""
    if _is_windows():
        _hide_file_windows(filepath)
    else:
        _hide_file_unix(filepath)

def unhide_file(filepath):
    """Unhides a file based on the operating system."""
    if _is_windows():
        _unhide_file_windows(filepath)
    else:
        _unhide_file_unix(filepath)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python script.py <hide|unhide> <filepath>")
        sys.exit(1)

    action = sys.argv[1].lower()
    filepath = sys.argv[2]

    if action == "hide":
        hide_file(filepath)
    elif action == "unhide":
        unhide_file(filepath)
    else:
        print("Invalid action. Please use 'hide' or 'unhide'.")
        sys.exit(1)

# Additional implementation at 2025-06-20 00:39:17
import os
import sys
import subprocess

def hide_file_windows(filepath):
    try:
        subprocess.run(['attrib', '+h', filepath], check=True, creationflags=subprocess.CREATE_NO_WINDOW)
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        return False

def unhide_file_windows(filepath):
    try:
        subprocess.run(['attrib', '-h', filepath], check=True, creationflags=subprocess.CREATE_NO_WINDOW)
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        return False

def is_hidden_windows(filepath):
    try:
        result = subprocess.run(['attrib', filepath], capture_output=True, text=True, check=True, creationflags=subprocess.CREATE_NO_WINDOW)
        return ' H ' in result.stdout
    except (subprocess.CalledProcessError, FileNotFoundError):
        return False

def hide_file_unix(filepath):
    directory = os.path.dirname(filepath)
    filename = os.path.basename(filepath)
    if filename.startswith('.'):
        return True
    new_filepath = os.path.join(directory, '.' + filename)
    try:
        os.rename(filepath, new_filepath)
        return True
    except OSError:
        return False

def unhide_file_unix(filepath):
    directory = os.path.dirname(filepath)
    filename = os.path.basename(filepath)
    if not filename.startswith('.'):
        return True
    new_filepath = os.path.join(directory, filename[1:])
    try:
        os.rename(filepath, new_filepath)
        return True
    except OSError:
        return False

def is_hidden_unix(filepath):
    return os.path.basename(filepath).startswith('.')

def list_files_in_directory(directory):
    if not os.path.isdir(directory):
        return None
    files_info = []
    is_hidden_func = is_hidden_windows if sys.platform == 'win32' else is_hidden_unix
    
    for item in os.listdir(directory):
        full_path = os.path.join(directory, item)
        if os.path.isfile(full_path):
            files_info.append((item, is_hidden_func(full_path)))
    return files_info

def main():
    if len(sys.argv) < 3:
        print("Usage: python script.py <action> <path>")
        print("Actions:")
        print("  hide <filepath>")
        print("  unhide <filepath>")
        print("  check <filepath>")
        print("  list <directory_path>")
        sys.exit(1)

    action = sys.argv[1].lower()
    path = sys.argv[2]

    if sys.platform == 'win32':
        hide_func = hide_file_windows
        unhide_func = unhide_file_windows
        is_hidden_func = is_hidden_windows
    else:
        hide_func = hide_file_unix
        unhide_func = unhide_file_unix
        is_hidden_func = is_hidden_unix

    if action == 'hide':
        if not os.path.exists(path):
            print(f"Error: File or directory not found at '{path}'")
            sys.exit(1)
        if hide_func(path):
            print(f"Successfully hid '{path}'.")
        else:
            print(f"Failed to hide '{path}'.")
    elif action == 'unhide':
        if not os.path.exists(path):
            print(f"Error: File or directory not found at '{path}'")
            sys.exit(1)
        if unhide_func(path):
            print(f"Successfully unhid '{path}'.")
        else:
            print(f"Failed to unhide '{path}'.")
    elif action == 'check':
        if not os.path.exists(path):
            print(f"Error: File or directory not found at '{path}'")
            sys.exit(1)
        if is_hidden_func(path):
            print(f"'{path}' is hidden.")
        else:
            print(f"'{path}' is not hidden.")
    elif action == 'list':
        files_info = list_files_in_directory(path)
        if files_info is None:
            print(f"Error: Directory not found or not a directory at '{path}'.")
            sys.exit(1)
        if not files_info:
            print(f"No files found in '{path}'.")
        else:
            print(f"Files in '{path}':")
            for filename, hidden_status in files_info:
                status_str = " (Hidden)" if hidden_status else ""
                print(f"  - {filename}{status_str}")
    else:
        print(f"Error: Unknown action '{action}'.")
        sys.exit(1)

if __name__ == '__main__':
    main()

# Additional implementation at 2025-06-20 00:40:39
import os
import sys
import subprocess

def _hide_file_windows(filepath):
    try:
        subprocess.run(['attrib', '+h', filepath], check=True, capture_output=True, text=True)
        return True
    except subprocess.CalledProcessError:
        return False
    except FileNotFoundError:
        return False
    except Exception:
        return False

def _unhide_file_windows(filepath):
    try:
        subprocess.run(['attrib', '-h', filepath], check=True, capture_output=True, text=True)
        return True
    except subprocess.CalledProcessError:
        return False
    except FileNotFoundError:
        return False
    except Exception:
        return False

def _hide_file_unix(filepath):
    dirname = os.path.dirname(filepath)
    basename = os.path.basename(filepath)

    if basename.startswith('.'):
        return True

    new_filepath = os.path.join(dirname, '.' + basename)
    try:
        os.rename(filepath, new_filepath)
        return True
    except OSError:
        return False
    except Exception:
        return False

def _unhide_file_unix(filepath):
    dirname = os.path.dirname(filepath)
    basename = os.path.basename(filepath)

    if not basename.startswith('.'):
        return True

    new_basename = basename[1:]
    new_filepath = os.path.join(dirname, new_basename)
    try:
        os.rename(filepath, new_filepath)
        return True
    except OSError:
        return False
    except Exception:
        return False

def hide_file(filepath):
    if sys.platform.startswith('win'):
        return _hide_file_windows(filepath)
    else:
        return _hide_file_unix(filepath)

def unhide_file(filepath):
    if sys.platform.startswith('win'):
        return _unhide_file_windows(filepath)
    else:
        return _unhide_file_unix(filepath)

def main():
    if len(sys.argv) < 3:
        sys.exit(1)

    action = sys.argv[1].lower()
    filepath = sys.argv[2]

    if action == 'hide':
        if not hide_file(filepath):
            sys.exit(1)
    elif action == 'unhide':
        if not unhide_file(filepath):
            sys.exit(1)
    else:
        sys.exit(1)

if __name__ == '__main__':
    main()