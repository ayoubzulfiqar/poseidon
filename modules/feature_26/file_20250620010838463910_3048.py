import os
import zipfile
import sys

def extract_all_zips_recursively(root_directory):
    if not os.path.isdir(root_directory):
        print(f"Error: Directory '{root_directory}' not found.")
        return

    print(f"Starting ZIP extraction in: {root_directory}")
    extracted_count = 0
    failed_count = 0

    for dirpath, dirnames, filenames in os.walk(root_directory):
        for filename in filenames:
            if filename.lower().endswith('.zip'):
                zip_file_path = os.path.join(dirpath, filename)
                
                extraction_target_dir = os.path.join(dirpath, os.path.splitext(filename)[0])

                try:
                    os.makedirs(extraction_target_dir, exist_ok=True)

                    with zipfile.ZipFile(zip_file_path, 'r') as zip_ref:
                        print(f"Extracting '{zip_file_path}' to '{extraction_target_dir}'...")
                        zip_ref.extractall(extraction_target_dir)
                        extracted_count += 1
                except zipfile.BadZipFile:
                    print(f"Error: '{zip_file_path}' is corrupted. Skipping.")
                    failed_count += 1
                except Exception as e:
                    print(f"Error extracting '{zip_file_path}': {e}")
                    failed_count += 1
    
    print("\n--- Summary ---")
    print(f"Processed: {extracted_count + failed_count} ZIPs")
    print(f"Extracted: {extracted_count}")
    print(f"Failed: {failed_count}")

if __name__ == "__main__":
    if len(sys.argv) > 1:
        target_directory = sys.argv[1]
    else:
        target_directory = os.getcwd()

    extract_all_zips_recursively(target_directory)

# Additional implementation at 2025-06-20 01:09:43
import os
import zipfile
import sys

def extract_all_zips_recursively(root_directory):
    """
    Extracts all ZIP files found recursively within the specified root directory.

    Each ZIP file will be extracted into a new subdirectory located in the same
    directory as the ZIP file. The new subdirectory will be named after the
    ZIP file (without extension) with "_extracted" appended.

    Args:
        root_directory (str): The path to the directory to start scanning from.
    """
    if not os.path.isdir(root_directory):
        print(f"Error: Directory not found: {root_directory}")
        return

    print(f"Scanning '{root_directory}' for ZIP files...")
    extracted_count = 0
    error_count = 0

    for dirpath, dirnames, filenames in os.walk(root_directory):
        for filename in filenames:
            if filename.lower().endswith('.zip'):
                zip_path = os.path.join(dirpath, filename)
                
                # Create a destination directory for the extracted contents
                # e.g., 'my_archive.zip' -> 'my_archive_extracted/'
                base_name = os.path.splitext(filename)[0]
                extract_to_dir = os.path.join(dirpath, base_name + "_extracted")

                try:
                    # Ensure the destination directory exists
                    os.makedirs(extract_to_dir, exist_ok=True)
                    
                    with zipfile.ZipFile(zip_path, 'r') as zip_ref:
                        print(f"  Extracting '{zip_path}' to '{extract_to_dir}'...")
                        zip_ref.extractall(extract_to_dir)
                        extracted_count += 1
                except zipfile.BadZipFile:
                    print(f"  Error: '{zip_path}' is a bad or corrupted ZIP file.")
                    error_count += 1
                except Exception as e:
                    print(f"  An unexpected error occurred while extracting '{zip_path}': {e}")
                    error_count += 1
    
    print(f"\nExtraction complete.")
    print(f"Successfully extracted {extracted_count} ZIP files.")
    if error_count > 0:
        print(f"Encountered {error_count} errors during extraction.")

if __name__ == "__main__":
    # Default directory to scan.
    # You can change this to a specific path, or pass it as a command-line argument.
    # Example: python your_script_name.py /path/to/your/folder
    
    if len(sys.argv) > 1:
        target_directory = sys.argv[1]
    else:
        # Use the current working directory as default if no argument is provided
        target_directory = os.getcwd() 
        print(f"No directory specified. Using current directory: {target_directory}")

    extract_all_zips_recursively(target_directory)

# Additional implementation at 2025-06-20 01:10:55
import os
import zipfile
import sys

def extract_all_zips_recursively(root_dir):
    if not os.path.isdir(root_dir):
        print(f"Error: Directory not found - {root_dir}")
        return

    print(f"Starting recursive ZIP extraction in: {root_dir}")

    for dirpath, dirnames, filenames in os.walk(root_dir):
        for filename in filenames:
            if filename.lower().endswith(".zip"):
                zip_filepath = os.path.join(dirpath, filename)
                
                extraction_dir_name = os.path.splitext(filename)[0]
                extraction_path = os.path.join(dirpath, extraction_dir_name)

                if os.path.exists(extraction_path):
                    if os.path.isdir(extraction_path) and os.listdir(extraction_path):
                        print(f"Skipping '{zip_filepath}': Extraction directory '{extraction_path}' already exists and is not empty.")
                        continue
                    elif os.path.isdir(extraction_path) and not os.listdir(extraction_path):
                        print(f"Warning: Empty extraction directory '{extraction_path}' found. Using it for '{zip_filepath}'.")
                    else:
                        print(f"Error: Cannot create extraction directory '{extraction_path}'. A file with the same name already exists.")
                        continue
                
                try:
                    print(f"Extracting '{zip_filepath}' to '{extraction_path}'...")
                    os.makedirs(extraction_path, exist_ok=True)
                    with zipfile.ZipFile(zip_filepath, 'r') as zip_ref:
                        zip_ref.extractall(extraction_path)
                    print(f"Successfully extracted: {zip_filepath}")
                except zipfile.BadZipFile:
                    print(f"Error: '{zip_filepath}' is not a valid ZIP file or is corrupted.")
                except Exception as e:
                    print(f"An unexpected error occurred while extracting '{zip_filepath}': {e}")
                print("-" * 50)

if __name__ == "__main__":
    if len(sys.argv) > 1:
        target_directory = sys.argv[1]
    else:
        target_directory = os.getcwd()

    extract_all_zips_recursively(target_directory)