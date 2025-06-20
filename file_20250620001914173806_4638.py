import pandas as pd
import os

def merge_csv_files(input_files, output_file):
    """
    Merges multiple CSV files into a single CSV file.

    Args:
        input_files (list): A list of paths to the input CSV files.
        output_file (str): The path for the merged output CSV file.
    """
    all_dataframes = []

    for file in input_files:
        df = pd.read_csv(file)
        all_dataframes.append(df)

    merged_df = pd.concat(all_dataframes, ignore_index=True)
    merged_df.to_csv(output_file, index=False)

if __name__ == "__main__":
    # --- Create dummy CSV files for demonstration ---
    # These files will be created in the same directory as the script.
    # In a real scenario, these files would already exist.
    dummy_file_1 = "sales_q1.csv"
    dummy_file_2 = "sales_q2.csv"
    dummy_file_3 = "sales_q3.csv"
    
    with open(dummy_file_1, "w") as f:
        f.write("Date,Product,Revenue\n")
        f.write("2023-01-01,Laptop,1200\n")
        f.write("2023-01-05,Mouse,25\n")
        f.write("2023-02-10,Keyboard,75\n")

    with open(dummy_file_2, "w") as f:
        f.write("Date,Product,Revenue\n")
        f.write("2023-04-01,Monitor,300\n")
        f.write("2023-05-15,Webcam,50\n")
        f.write("2023-06-20,Headphones,100\n")

    with open(dummy_file_3, "w") as f:
        f.write("Date,Product,Revenue\n")
        f.write("2023-07-01,Printer,250\n")
        f.write("2023-08-05,Router,80\n")
        f.write("2023-09-10,Speaker,60\n")
    # --- End of dummy file creation ---

    # Define the list of input CSV files
    input_csv_files = [dummy_file_1, dummy_file_2, dummy_file_3]

    # Define the name for the output merged CSV file
    output_merged_csv = "all_sales_data.csv"

    # Call the function to merge the CSV files
    merge_csv_files(input_csv_files, output_merged_csv)

    # --- Optional: Clean up dummy files after merging ---
    # Uncomment the following lines if you want to remove the temporary dummy files
    # os.remove(dummy_file_1)
    # os.remove(dummy_file_2)
    # os.remove(dummy_file_3)
    # os.remove(output_merged_csv) # Uncomment to also remove the output file
    # --- End of optional cleanup ---

# Additional implementation at 2025-06-20 00:20:11
import pandas as pd
import os
import glob

def get_files_to_merge():
    files = []
    while True:
        choice = input("Do you want to (1) list files or (2) merge all CSVs in a directory? (1/2): ").strip()
        if choice == '1':
            while True:
                file_path = input("Enter CSV file path (or 'done' to finish): ").strip()
                if file_path.lower() == 'done':
                    break
                if not file_path.lower().endswith('.csv'):
                    print(f"'{file_path}' is not a CSV file. Please enter a valid CSV path.")
                    continue
                if not os.path.exists(file_path):
                    print(f"File not found: '{file_path}'. Please check the path.")
                    continue
                files.append(file_path)
            if not files:
                print("No files selected. Please select at least one CSV file.")
                continue
            break
        elif choice == '2':
            directory = input("Enter the directory path containing CSV files: ").strip()
            if not os.path.isdir(directory):
                print(f"Directory not found: '{directory}'. Please check the path.")
                continue
            csv_files_in_dir = glob.glob(os.path.join(directory, '*.csv'))
            if not csv_files_in_dir:
                print(f"No CSV files found in directory: '{directory}'.")
                continue
            files.extend(csv_files_in_dir)
            break
        else:
            print("Invalid choice. Please enter '1' or '2'.")
    return files

def get_output_filename():
    while True:
        output_name = input("Enter the desired output CSV filename (e.g., merged_data.csv): ").strip()
        if not output_name:
            print("Output filename cannot be empty.")
            continue
        if not output_name.lower().endswith('.csv'):
            output_name += '.csv'
        return output_name

def get_drop_duplicates_option():
    while True:
        choice = input("Do you want to drop duplicate rows in the merged file? (yes/no): ").strip().lower()
        if choice in ['yes', 'y']:
            return True
        elif choice in ['no', 'n']:
            return False
        else:
            print("Invalid choice. Please enter 'yes' or 'no'.")

def merge_csv_files():
    print("--- CSV File Merger ---")

    input_files = get_files_to_merge()
    if not input_files:
        print("No CSV files selected for merging. Exiting.")
        return

    output_file = get_output_filename()
    drop_duplicates = get_drop_duplicates_option()

    all_dataframes = []
    print("\nReading CSV files...")
    for i, file_path in enumerate(input_files):
        try:
            df = pd.read_csv(file_path)
            all_dataframes.append(df)
            print(f"  Successfully read: {os.path.basename(file_path)}")
        except FileNotFoundError:
            print(f"  Error: File not found - {file_path}")
        except pd.errors.EmptyDataError:
            print(f"  Warning: File is empty - {file_path}. Skipping.")
        except Exception as e:
            print(f"  Error reading {file_path}: {e}")

    if not all_dataframes:
        print("No valid data to merge. Exiting.")
        return

    print("\nMerging dataframes...")
    try:
        merged_df = pd.concat(all_dataframes, ignore_index=True)
        print("Dataframes merged successfully.")

        if drop_duplicates:
            initial_rows = len(merged_df)
            merged_df.drop_duplicates(inplace=True)
            rows_after_dedup = len(merged_df)
            print(f"Dropped {initial_rows - rows_after_dedup} duplicate rows.")

        print(f"\nSaving merged data to '{output_file}'...")
        merged_df.to_csv(output_file, index=False)
        print(f"Successfully saved merged data to '{output_file}'")
        print(f"Total rows in merged file: {len(merged_df)}")

    except Exception as e:
        print(f"An error occurred during merging or saving: {e}")

if __name__ == "__main__":
    merge_csv_files()