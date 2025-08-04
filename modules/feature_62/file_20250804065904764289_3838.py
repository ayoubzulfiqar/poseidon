import os
import subprocess
import sys

def hide_file(filepath):
    try:
        abs_filepath = os.path.abspath(filepath)
        if not os.path.exists(abs_filepath):
            raise FileNotFoundError()

        if sys.platform == "win32":
            subprocess.run(['attrib', '+h', abs_filepath], capture_output=True, text=True, check=True)
        else:
            dirname = os.path.dirname(abs_filepath)
            basename = os.path.basename(abs_filepath)
            
            if not basename.startswith('.'):
                new_filepath = os.path.join(dirname, '.' + basename)
                os.rename(abs_filepath, new_filepath)
    except (FileNotFoundError, PermissionError, subprocess.CalledProcessError, OSError):
        pass
    except Exception:
        pass

def unhide_file(filepath):
    try:
        abs_filepath = os.path.abspath(filepath)
        if not os.path.exists(abs_filepath):
            raise FileNotFoundError()

        if sys.platform == "win32":
            subprocess.run(['attrib', '-h', abs_filepath], capture_output=True, text=True, check=True)
        else:
            dirname = os.path.dirname(abs_filepath)
            basename = os.path.basename(abs_filepath)
            
            if basename.startswith('.'):
                new_filepath = os.path.join(dirname, basename[1:])
                os.rename(abs_filepath, new_filepath)
    except (FileNotFoundError, PermissionError, subprocess.CalledProcessError, OSError):
        pass
    except Exception:
        pass

if __name__ == "__main__":
    test_filename = "test_hidden_file.txt"
    
    try:
        with open(test_filename, "w") as f:
            f.write("This is a test file.")

        hide_file(test_filename)
        unhide_file(test_filename)

    except Exception:
        pass
    finally:
        try:
            if os.path.exists(test_filename):
                os.remove(test_filename)
            
            if sys.platform != "win32":
                dot_prefixed_name = "." + test_filename
                if os.path.exists(dot_prefixed_name):
                    os.remove(dot_prefixed_name)
        except Exception:
            pass

# Additional implementation at 2025-08-04 06:59:37
import os
import sys
import subprocess

def _hide_windows(file_path):
    try:
        subprocess.check_call(['attrib', '+h', file_path], shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        print(f"Successfully hid: {file_path}")
        return True
    except subprocess.CalledProcessError as e:
        print(f"Error hiding {file_path}: {e.stderr.decode().strip()}")
        return False
    except FileNotFoundError:
        print(f"Error: 'attrib' command not found. Ensure it's in your system PATH.")
        return False
    except Exception as e:
        print(f"An unexpected error occurred hiding {file_path}: {e}")
        return False

def _unhide_windows(file_path):
    try:
        subprocess.check_call(['attrib', '-h', file_path], shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        print(f"Successfully unhid: {file_path}")
        return True
    except subprocess.CalledProcessError as e:
        print(f"Error unhiding {file_path}: {e.stderr.decode().strip()}")
        return False
    except FileNotFoundError:
        print(f"Error: 'attrib' command not found. Ensure it's in your system PATH.")
        return False
    except Exception as e:
        print(f"An unexpected error occurred unhiding {file_path}: {e}")
        return False

def _hide_unix(file_path):
    directory = os.path.dirname(file_path)
    basename = os.path.basename(file_path)

    if basename.startswith('.'):
        print(f"File {file_path} is already hidden.")
        return True

    new_path = os.path.join(directory, '.' + basename)
    try:
        os.rename(file_path, new_path)
        print(f"Successfully hid: {file_path} -> {new_path}")
        return True
    except FileNotFoundError:
        print(f"Error: File not found at {file_path}")
        return False
    except OSError as e:
        print(f"Error hiding {file_path}: {e}")
        return False
    except Exception as e:
        print(f"An unexpected error occurred hiding {file_path}: {e}")
        return False

def _unhide_unix(file_path):
    directory = os.path.dirname(file_path)
    basename = os.path.basename(file_path)

    if not basename.startswith('.'):
        print(f"File {file_path} is not hidden by dot prefix method.")
        return True

    new_basename = basename[1:]
    new_path = os.path.join(directory, new_basename)
    try:
        os.rename(file_path, new_path)
        print(f"Successfully unhid: {file_path} -> {new_path}")
        return True
    except FileNotFoundError:
        print(f"Error: File not found at {file_path}")
        return False
    except OSError as e:
        print(f"Error unhiding {file_path}: {e}")
        return False
    except Exception as e:
        print(f"An unexpected error occurred unhiding {file_path}: {e}")
        return False

def hide_file(file_path):
    if sys.platform.startswith('win'):
        return _hide_windows(file_path)
    else:
        return _hide_unix(file_path)

def unhide_file(file_path):
    if sys.platform.startswith('win'):
        return _unhide_windows(file_path)
    else:
        return _unhide_unix(file_path)

def main():
    while True:
        print("\nFile Hider/Unhider")
        print("1. Hide a file")
        print("2. Unhide a file")
        print("3. Exit")

        choice = input("Enter your choice (1/2/3): ").strip()

        if choice == '1':
            file_path = input("Enter the full path of the file to hide: ").strip()
            if not file_path:
                print("File path cannot be empty.")
                continue
            hide_file(file_path)
        elif choice == '2':
            file_path = input("Enter the full path of the file to unhide: ").strip()
            if not file_path:
                print("File path cannot be empty.")
                continue
            unhide_file(file_path)
        elif choice == '3':
            print("Exiting program.")
            break
        else:
            print("Invalid choice. Please enter 1, 2, or 3.")

if __name__ == '__main__':
    main()

# Additional implementation at 2025-08-04 07:00:11
import os
import sys
import subprocess

IS_WINDOWS = sys.platform == 'win32'
if IS_WINDOWS:
    import stat
    FILE_ATTRIBUTE_HIDDEN = stat.FILE_ATTRIBUTE_HIDDEN

def hide_file(filepath):
    if not os.path.exists(filepath):
        print(f"Error: File not found at '{filepath}'")
        return False
    try:
        if IS_WINDOWS:
            subprocess.run(['attrib', '+h', filepath], check=True, capture_output=True, text=True)
            print(f"Successfully hid '{filepath}' (Windows).")
        else:
            directory = os.path.dirname(filepath)
            filename = os.path.basename(filepath)
            if not filename.startswith('.'):
                new_filepath = os.path.join(directory, '.' + filename)
                os.rename(filepath, new_filepath)
                print(f"Successfully hid '{filepath}' as '{new_filepath}' (Unix).")
            else:
                print(f"'{filepath}' is already hidden (starts with '.').")
        return True
    except subprocess.CalledProcessError as e:
        print(f"Error hiding '{filepath}' (Windows attrib failed): {e.stderr.strip()}")
        return False
    except OSError as e:
        print(f"Error hiding '{filepath}': {e}")
        return False

def unhide_file(filepath):
    if not os.path.exists(filepath):
        print(f"Error: File not found at '{filepath}'")
        return False
    try:
        if IS_WINDOWS:
            subprocess.run(['attrib', '-h', filepath], check=True, capture_output=True, text=True)
            print(f"Successfully unhid '{filepath}' (Windows).")
        else:
            directory = os.path.dirname(filepath)
            filename = os.path.basename(filepath)
            if filename.startswith('.'):
                new_filepath = os.path.join(directory, filename[1:])
                os.rename(filepath, new_filepath)
                print(f"Successfully unhid '{filepath}' as '{new_filepath}' (Unix).")
            else:
                print(f"'{filepath}' is not hidden (does not start with '.').")
        return True
    except subprocess.CalledProcessError as e:
        print(f"Error unhiding '{filepath}' (Windows attrib failed): {e.stderr.strip()}")
        return False
    except OSError as e:
        print(f"Error unhiding '{filepath}': {e}")
        return False

def list_files_in_directory(directory):
    if not os.path.isdir(directory):
        print(f"Error: Directory not found at '{directory}'")
        return []
    files = []
    try:
        for item in os.listdir(directory):
            full_path = os.path.join(directory, item)
            if os.path.isfile(full_path):
                files.append(item)
    except OSError as e:
        print(f"Error listing files in '{directory}': {e}")
    return files

def list_hidden_files_in_directory(directory):
    if not os.path.isdir(directory):
        print(f"Error: Directory not found at '{directory}'")
        return []
    hidden_files = []
    try:
        for item in os.listdir(directory):
            full_path = os.path.join(directory, item)
            if os.path.isfile(full_path):
                if IS_WINDOWS:
                    try:
                        attributes = os.stat(full_path).st_file_attributes
                        if attributes & FILE_ATTRIBUTE_HIDDEN:
                            hidden_files.append(item)
                    except OSError:
                        pass
                else:
                    if item.startswith('.'):
                        hidden_files.append(item)
    except OSError as e:
        print(f"Error listing hidden files in '{directory}': {e}")
    return hidden_files

def get_user_file_selection(file_list, prompt_message="Select a file by number:"):
    if not file_list:
        print("No files available to select.")
        return None
    for i, file_name in enumerate(file_list):
        print(f"{i + 1}. {file_name}")
    while True:
        try:
            choice = input(prompt_message + " ")
            index = int(choice) - 1
            if 0 <= index < len(file_list):
                return file_list[index]
            else:
                print("Invalid number. Please try again.")
        except ValueError:
            print("Invalid input. Please enter a number.")

def main_menu():
    while True:
        print("\n--- File Hider/Unhider ---")
        print("1. Hide a file")
        print("2. Unhide a file")
        print("3. List hidden files")
        print("4. Exit")
        choice = input("Enter your choice: ")
        if choice == '1':
            target_dir = input("Enter the directory path (or press Enter for current directory): ")
            if not target_dir:
                target_dir = os.getcwd()
            files = list_files_in_directory(target_dir)
            if files:
                selected_file = get_user_file_selection(files, "Select a file to HIDE:")
                if selected_file:
                    full_path = os.path.join(target_dir, selected_file)
                    hide_file(full_path)
            else:
                print("No files found in the specified directory.")
        elif choice == '2':
            target_dir = input("Enter the directory path (or press Enter for current directory): ")
            if not target_dir:
                target_dir = os.getcwd()
            files_to_unhide = list_hidden_files_in_directory(target_dir)
            if files_to_unhide:
                selected_file = get_user_file_selection(files_to_unhide, "Select a file to UNHIDE:")
                if selected_file:
                    full_path = os.path.join(target_dir, selected_file)
                    unhide_file(full_path)
            else:
                print("No hidden files found in the specified directory.")
        elif choice == '3':
            target_dir = input("Enter the directory path (or press Enter for current directory): ")
            if not target_dir:
                target_dir = os.getcwd()
            hidden_files = list_hidden_files_in_directory(target_dir)
            if hidden_files:
                print(f"\nHidden files in '{target_dir}':")
                for f in hidden_files:
                    print(f"- {f}")
            else:
                print(f"No hidden files found in '{target_dir}'.")
        elif choice == '4':
            print("Exiting program.")
            break
        else:
            print("Invalid choice. Please try again.")

main_menu()

# Additional implementation at 2025-08-04 07:01:14
import os
import sys
import subprocess
import stat

def hide_file(filepath):
    if not os.path.exists(filepath):
        return False
    if sys.platform.startswith('win'):
        try:
            subprocess.run(['attrib', '+h', filepath], check=True, capture_output=True, text=True)
            return True
        except subprocess.CalledProcessError:
            return False
        except Exception:
            return False
    else:
        dirname, basename = os.path.split(filepath)
        if not basename.startswith('.'):
            new_filepath = os.path.join(dirname, '.' + basename)
            try:
                os.rename(filepath, new_filepath)
                return True
            except OSError:
                return False
        return True

def unhide_file(filepath):
    if not os.path.exists(filepath):
        return False
    if sys.platform.startswith('win'):
        try:
            subprocess.run(['attrib', '-h', filepath], check=True, capture_output=True, text=True)
            return True
        except subprocess.CalledProcessError:
            return False
        except Exception:
            return False
    else:
        dirname, basename = os.path.split(filepath)
        if basename.startswith('.'):
            new_basename = basename[1:]
            new_filepath = os.path.join(dirname, new_basename)
            try:
                os.rename(filepath, new_filepath)
                return True
            except OSError:
                return False
        return True

def is_file_hidden(filepath):
    if not os.path.exists(filepath):
        return False
    if sys.platform.startswith('win'):
        try:
            attributes = os.stat(filepath).st_file_attributes
            return (attributes & stat.FILE_ATTRIBUTE_HIDDEN) != 0
        except OSError:
            return False
    else:
        basename = os.path.basename(filepath)
        return basename.startswith('.')

if __name__ == "__main__":
    test_file_name = "test_file_for_hiding.txt"
    hidden_test_file_name_unix = "." + test_file_name

    if sys.platform.startswith('win'):
        if os.path.exists(test_file_name):
            unhide_file(test_file_name)
            os.remove(test_file_name)
    else:
        if os.path.exists(test_file_name):
            os.remove(test_file_name)
        if os.path.exists(hidden_test_file_name_unix):
            os.remove(hidden_test_file_name_unix)

    with open(test_file_name, 'w') as f:
        f.write("This is a test file.")

    print(is_file_hidden(test_file_name))

    hide_file(test_file_name)
    
    if sys.platform.startswith('win'):
        print(is_file_hidden(test_file_name))
        unhide_file(test_file_name)
    else:
        print(is_file_hidden(hidden_test_file_name_unix))
        unhide_file(hidden_test_file_name_unix)

    print(is_file_hidden(test_file_name))

    if os.path.exists(test_file_name):
        os.remove(test_file_name)
    if not sys.platform.startswith('win') and os.path.exists(hidden_test_file_name_unix):
        os.remove(hidden_test_file_name_unix)