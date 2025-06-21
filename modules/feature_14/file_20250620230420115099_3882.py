import os
import zipfile
import sys

def extract_zips_recursively(root_folder):
    if not os.path.isdir(root_folder):
        print(f"Error: '{root_folder}' is not a valid directory.")
        return

    print(f"Starting recursive ZIP extraction in: {root_folder}")

    for dirpath, dirnames, filenames in os.walk(root_folder):
        for filename in filenames:
            if filename.lower().endswith('.zip'):
                zip_filepath = os.path.join(dirpath, filename)
                
                extract_to_dir = os.path.join(dirpath, os.path.splitext(filename)[0])
                
                try:
                    os.makedirs(extract_to_dir, exist_ok=True)
                    with zipfile.ZipFile(zip_filepath, 'r') as zip_ref:
                        print(f"Extracting '{zip_filepath}' to '{extract_to_dir}'...")
                        zip_ref.extractall(extract_to_dir)
                        print(f"Successfully extracted '{filename}'.")
                except zipfile.BadZipFile:
                    print(f"Error: '{zip_filepath}' is a bad or corrupted ZIP file.")
                except Exception as e:
                    print(f"An unexpected error occurred while extracting '{zip_filepath}': {e}")
    
    print("ZIP extraction process completed.")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python script_name.py <folder_path>")
        sys.exit(1)
    
    target_folder = sys.argv[1]
    extract_zips_recursively(target_folder)

# Additional implementation at 2025-06-20 23:05:28
import os
import zipfile
import sys

def extract_all_zips_recursively(source_directory):
    """
    Recursively finds and extracts all ZIP files within a given directory
    and its subdirectories. Each ZIP file is extracted into a new folder
    with the same name as the ZIP file (without the .zip extension)
    in the same directory where the ZIP file is located.

    Args:
        source_directory (str): The path to the directory to start scanning from.
    """
    if not os.path.isdir(source_directory):
        print(f"Error: '{source_directory}' is not a valid directory.")
        return

    print(f"Starting recursive ZIP extraction in: {source_directory}")
    found_zips = 0
    extracted_count = 0
    skipped_count = 0
    error_count = 0

    for root, _, files in os.walk(source_directory):
        for file in files:
            if file.lower().endswith(".zip"):
                zip_path = os.path.join(root, file)
                found_zips += 1
                
                # Determine extraction destination
                # Extract to a folder named after the zip file (without .zip)
                # in the same directory as the zip file.
                zip_name_without_ext = os.path.splitext(file)[0]
                extract_destination = os.path.join(root, zip_name_without_ext)

                print(f"\nProcessing: {zip_path}")
                print(f"  Destination: {extract_destination}")

                try:
                    # Create destination directory if it doesn't exist
                    os.makedirs(extract_destination, exist_ok=True)

                    with zipfile.ZipFile(zip_path, 'r') as zip_ref:
                        # Test the integrity of the archive before extracting
                        try:
                            zip_ref.testzip() 
                            print("  Archive integrity check passed.")
                        except zipfile.BadZipFile:
                            print(f"  Skipping '{zip_path}': Appears to be a bad or corrupted ZIP file (or password protected).")
                            skipped_count += 1
                            continue
                        except Exception as e:
                            print(f"  Error testing integrity of '{zip_path}': {e}")
                            error_count += 1
                            continue

                        print("  Extracting contents...")
                        zip_ref.extractall(extract_destination)
                        print(f"  Successfully extracted '{zip_path}' to '{extract_destination}'.")
                        extracted_count += 1

                except zipfile.BadZipFile:
                    print(f"  Error: '{zip_path}' is not a valid ZIP file or is corrupted.")
                    error_count += 1
                except zipfile.LargeZipFile:
                    print(f"  Error: '{zip_path}' requires ZIP64 support, which might not be fully supported or the file is too large.")
                    error_count += 1
                except PermissionError:
                    print(f"  Error: Permission denied when trying to extract '{zip_path}' to '{extract_destination}'.")
                    error_count += 1
                except FileNotFoundError:
                    print(f"  Error: ZIP file '{zip_path}' not found (might have been moved or deleted during scan).")
                    error_count += 1
                except Exception as e:
                    print(f"  An unexpected error occurred with '{zip_path}': {e}")
                    error_count += 1

    print("\n--- Extraction Summary ---")
    print(f"Total ZIP files found: {found_zips}")
    print(f"Successfully extracted: {extracted_count}")
    print(f"Skipped (e.g., corrupted/password protected): {skipped_count}")
    print(f"Errors encountered: {error_count}")
    print("Extraction process complete.")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python your_script_name.py <directory_path>")
        print("Example: python your_script_name.py C:\\MyZips")
        print("Example: python your_script_name.py /home/user/archives")
    else:
        target_directory = sys.argv[1]
        extract_all_zips_recursively(target_directory)