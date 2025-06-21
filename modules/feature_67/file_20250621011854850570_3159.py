import os

def rename_files_sequentially(folder_path):
    if not os.path.isdir(folder_path):
        print(f"Error: Folder '{folder_path}' not found or is not a directory.")
        return

    try:
        files = [f for f in os.listdir(folder_path) if os.path.isfile(os.path.join(folder_path, f))]
        files.sort()
    except OSError as e:
        print(f"Error accessing folder '{folder_path}': {e}")
        return

    if not files:
        print(f"No files found in '{folder_path}'.")
        return

    print(f"Attempting to rename {len(files)} files in '{folder_path}'...")
    for index, filename in enumerate(files):
        old_path = os.path.join(folder_path, filename)
        name, ext = os.path.splitext(filename)
        new_filename = f"{index + 1}{ext}"
        new_path = os.path.join(folder_path, new_filename)

        try:
            os.rename(old_path, new_path)
            print(f"Renamed '{filename}' to '{new_filename}'")
        except OSError as e:
            print(f"Failed to rename '{filename}' to '{new_filename}': {e}")
    print("Renaming process finished.")

if __name__ == "__main__":
    target_folder = input("Please enter the full path to the folder: ")
    rename_files_sequentially(target_folder)

# Additional implementation at 2025-06-21 01:19:24
import os
import sys

def rename_files_sequentially():
    print("--- Sequential File Renamer ---")

    while True:
        folder_path = input("Enter the folder path to rename files in (or 'q' to quit): ").strip()
        if folder_path.lower() == 'q':
            sys.exit()
        if not os.path.isdir(folder_path):
            print(f"Error: Folder '{folder_path}' not found or is not a directory. Please try again.")
        else:
            break

    while True:
        start_num_str = input("Enter the starting number for renaming (default is 1): ").strip()
        if not start_num_str:
            start_num = 1
            break
        try:
            start_num = int(start_num_str)
            if start_num < 0:
                print("Starting number cannot be negative. Please enter a positive integer.")
            else:
                break
        except ValueError:
            print("Invalid input. Please enter an integer.")

    target_extension = input("Enter file extension to target (e.g., 'txt', 'jpg', or leave blank for all files): ").strip().lower()
    if target_extension and not target_extension.startswith('.'):
        target_extension = '.' + target_extension

    while True:
        padding_str = input("Enter desired number of digits for padding (e.g., 3 for 001, 002; 0 for no padding): ").strip()
        if not padding_str:
            padding = 0
            break
        try:
            padding = int(padding_str)
            if padding < 0:
                print("Padding cannot be negative. Please enter a non-negative integer.")
            else:
                break
        except ValueError:
            print("Invalid input. Please enter an integer.")

    files_to_rename = []
    try:
        for filename in os.listdir(folder_path):
            file_path = os.path.join(folder_path, filename)
            if os.path.isfile(file_path):
                _, ext = os.path.splitext(filename)
                if not target_extension or ext.lower() == target_extension:
                    files_to_rename.append(filename)
    except OSError as e:
        print(f"Error accessing folder: {e}")
        sys.exit()

    if not files_to_rename:
        print("No files found to rename with the specified criteria.")
        return

    files_to_rename.sort()

    print("\n--- Preview of Renaming ---")
    print(f"Files in '{folder_path}' will be renamed as follows:")
    current_number = start_num
    for old_name in files_to_rename:
        _, original_ext = os.path.splitext(old_name)
        if padding > 0:
            new_name = f"{current_number:0{padding}d}{original_ext}"
        else:
            new_name = f"{current_number}{original_ext}"
        print(f"  '{old_name}' -> '{new_name}'")
        current_number += 1

    confirm = input("\nDo you want to proceed with renaming these files? (yes/no): ").strip().lower()
    if confirm != 'yes':
        print("Renaming cancelled.")
        return

    print("\n--- Renaming Files ---")
    current_number = start_num
    renamed_count = 0
    for old_name in files_to_rename:
        old_file_path = os.path.join(folder_path, old_name)
        _, original_ext = os.path.splitext(old_name)

        if padding > 0:
            new_name = f"{current_number:0{padding}d}{original_ext}"
        else:
            new_name = f"{current_number}{original_ext}"
        new_file_path = os.path.join(folder_path, new_name)

        try:
            if os.path.exists(new_file_path):
                print(f"Warning: '{new_file_path}' already exists. Skipping '{old_name}'.")
            else:
                os.rename(old_file_path, new_file_path)
                print(f"Renamed '{old_name}' to '{new_name}'")
                renamed_count += 1
        except OSError as e:
            print(f"Error renaming '{old_name}': {e}")
        current_number += 1

    print(f"\nRenaming complete. {renamed_count} files renamed.")

if __name__ == "__main__":
    rename_files_sequentially()

# Additional implementation at 2025-06-21 01:20:32
import os

def rename_files_sequentially():
    print("--- File Renamer with Sequential Numbers ---")

    folder_path = input("Enter the folder path where files need to be renamed: ").strip()

    if not os.path.isdir(folder_path):
        print(f"Error: The path '{folder_path}' is not a valid directory.")
        return

    try:
        start_num_str = input("Enter the starting number for renaming (e.g., 1): ").strip()
        start_number = int(start_num_str)
        if start_number < 0:
            print("Error: Starting number cannot be negative.")
            return
    except ValueError:
        print("Error: Invalid starting number. Please enter an integer.")
        return

    target_extension = input("Enter the file extension to target (e.g., .txt, .jpg). Leave blank to target all files: ").strip()
    if target_extension and not target_extension.startswith('.'):
        target_extension = '.' + target_extension

    preview_input = input("Do you want to preview the changes before applying? (yes/no): ").strip().lower()
    preview_only = preview_input == 'yes'

    print(f"\nScanning folder: {folder_path}")

    files_to_rename = []
    for filename in os.listdir(folder_path):
        file_path = os.path.join(folder_path, filename)
        if os.path.isfile(file_path):
            if not target_extension or filename.lower().endswith(target_extension.lower()):
                files_to_rename.append(filename)

    if not files_to_rename:
        print("No files found matching the criteria in the specified folder.")
        return

    files_to_rename.sort()

    max_number = start_number + len(files_to_rename) - 1
    padding_digits = len(str(max_number))

    rename_map = {}
    current_number = start_number
    for old_filename in files_to_rename:
        name, ext = os.path.splitext(old_filename)
        new_filename = f"{current_number:0{padding_digits}d}{ext}"
        if old_filename == new_filename:
            print(f"Skipping '{old_filename}' as its new name would be identical.")
            current_number += 1
            continue
        rename_map[old_filename] = new_filename
        current_number += 1

    print("\n--- Proposed Renaming Plan ---")
    if not rename_map:
        print("No files will be renamed based on the current plan.")
        return

    for old_name, new_name in rename_map.items():
        print(f"'{old_name}' -> '{new_name}'")

    if preview_only:
        print("\nPreview mode enabled. No files were actually renamed.")
        print("Run the script again and choose 'no' for preview to apply changes.")
        return

    confirm = input("\nDo you want to proceed with renaming these files? (yes/no): ").strip().lower()
    if confirm != 'yes':
        print("Renaming cancelled by user.")
        return

    print("\n--- Applying Renaming ---")
    successful_renames = 0
    failed_renames = 0
    for old_filename, new_filename in rename_map.items():
        old_path = os.path.join(folder_path, old_filename)
        new_path = os.path.join(folder_path, new_filename)
        try:
            os.rename(old_path, new_path)
            print(f"Renamed '{old_filename}' to '{new_filename}'")
            successful_renames += 1
        except OSError as e:
            print(f"Error renaming '{old_filename}' to '{new_filename}': {e}")
            failed_renames += 1

    print(f"\nRenaming complete.")
    print(f"Successfully renamed: {successful_renames} files.")
    print(f"Failed to rename: {failed_renames} files.")

if __name__ == "__main__":
    rename_files_sequentially()

# Additional implementation at 2025-06-21 01:21:47
import os
import sys

def rename_files_sequentially():
    print("--- Sequential File Renamer ---")

    while True:
        folder_path = input("Enter the full path to the folder: ").strip()
        if not folder_path:
            print("Folder path cannot be empty. Please try again.")
            continue
        if not os.path.isdir(folder_path):
            print(f"Error: '{folder_path}' is not a valid directory. Please try again.")
        else:
            break

    while True:
        try:
            start_num_str = input("Enter the starting number for the sequence (e.g., 1, 100): ").strip()
            start_number = int(start_num_str)
            if start_number < 0:
                print("Starting number cannot be negative. Please enter a non-negative integer.")
            else:
                break
        except ValueError:
            print("Invalid input. Please enter an integer.")

    while True:
        try:
            padding_str = input("Enter the desired number of digits for padding (e.g., 3 for 001, 010): ").strip()
            padding = int(padding_str)
            if padding < 1:
                print("Padding must be at least 1. Please enter a positive integer.")
            else:
                break
        except ValueError:
            print("Invalid input. Please enter an integer.")

    try:
        files = [f for f in os.listdir(folder_path) if os.path.isfile(os.path.join(folder_path, f))]
        if not files:
            print("No files found in the specified folder. Exiting.")
            return
    except PermissionError:
        print(f"Error: Permission denied to access '{folder_path}'.")
        return
    except OSError as e:
        print(f"Error accessing folder '{folder_path}': {e}")
        return

    files.sort()

    renaming_map = {}
    current_number = start_number

    print("\n--- Preview of Renaming ---")
    print("Original Name -> New Name")
    print("---------------------------")

    for original_filename in files:
        name, ext = os.path.splitext(original_filename)
        new_filename = f"{current_number:0{padding}d}{ext}"
        renaming_map[original_filename] = new_filename
        print(f"{original_filename} -> {new_filename}")
        current_number += 1

    print("---------------------------")

    confirm = input("\nDo you want to proceed with these renames? (yes/no): ").strip().lower()
    if confirm != 'yes':
        print("Renaming cancelled by user.")
        return

    print("\n--- Performing Renaming ---")
    renamed_count = 0
    failed_count = 0

    for original_filename, new_filename in renaming_map.items():
        original_path = os.path.join(folder_path, original_filename)
        new_path = os.path.join(folder_path, new_filename)

        try:
            os.rename(original_path, new_path)
            print(f"Renamed '{original_filename}' to '{new_filename}'")
            renamed_count += 1
        except FileExistsError:
            print(f"Error: Cannot rename '{original_filename}' to '{new_filename}'. A file with the new name already exists.")
            failed_count += 1
        except PermissionError:
            print(f"Error: Permission denied to rename '{original_filename}' to '{new_filename}'.")
            failed_count += 1
        except OSError as e:
            print(f"Error renaming '{original_filename}' to '{new_filename}': {e}")
            failed_count += 1

    print("\n--- Renaming Summary ---")
    print(f"Successfully renamed: {renamed_count} files")
    print(f"Failed to rename: {failed_count} files")
    print("Operation complete.")

if __name__ == "__main__":
    rename_files_sequentially()