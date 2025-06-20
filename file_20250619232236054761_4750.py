import os

def find_files_containing_string(root_directory, search_string):
    found_files = []
    for dirpath, dirnames, filenames in os.walk(root_directory):
        for filename in filenames:
            filepath = os.path.join(dirpath, filename)
            try:
                with open(filepath, 'r', encoding='utf-8', errors='ignore') as f:
                    content = f.read()
                    if search_string in content:
                        found_files.append(filepath)
            except Exception:
                pass
    return found_files

directory_to_search = os.getcwd()
search_string_to_find = "your_string_here"

files_containing_string = find_files_containing_string(directory_to_search, search_string_to_find)

for file_path in files_containing_string:
    print(file_path)

# Additional implementation at 2025-06-19 23:23:51
import os
import argparse

def _process_single_file(file_path, search_string_processed, case_sensitive, allowed_file_types):
    """
    Helper function to process a single file for the search string.
    Prints matching lines directly.
    """
    if allowed_file_types:
        _, ext = os.path.splitext(file_path)
        if ext.lower() not in allowed_file_types:
            return

    try:
        with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
            for line_num, line in enumerate(f, 1):
                line_to_check = line if case_sensitive else line.lower()
                if search_string_processed in line_to_check:
                    print(f"{file_path}:{line_num}:{line.strip()}")
    except IOError:
        pass
    except Exception:
        pass

def main():
    parser = argparse.ArgumentParser(
        description="Finds files containing a specific string with additional functionality."
    )
    parser.add_argument(
        "directory",
        type=str,
        help="The starting directory to search."
    )
    parser.add_argument(
        "search_string",
        type=str,
        help="The string to search for."
    )
    parser.add_argument(
        "-r", "--recursive",
        action="store_true",
        help="Search subdirectories recursively."
    )
    parser.add_argument(
        "-i", "--ignore-case",
        action="store_true",
        help="Perform a case-insensitive search."
    )
    parser.add_argument(
        "-t", "--file-types",
        nargs="+",
        help="Filter by file extensions (e.g., .txt .py). Include the dot."
    )

    args = parser.parse_args()

    if not os.path.isdir(args.directory):
        return

    search_string_processed = args.search_string if not args.ignore_case else args.search_string.lower()

    allowed_file_types = {ft.lower() for ft in args.file_types} if args.file_types else None

    if args.recursive:
        for root, _, files in os.walk(args.directory):
            for filename in files:
                file_path = os.path.join(root, filename)
                _process_single_file(file_path, search_string_processed, not args.ignore_case, allowed_file_types)
    else:
        for filename in os.listdir(args.directory):
            file_path = os.path.join(args.directory, filename)
            if os.path.isfile(file_path):
                _process_single_file(file_path, search_string_processed, not args.ignore_case, allowed_file_types)

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-19 23:25:03
import os
import argparse

def find_string_in_files(start_dir, search_string, recursive=False, ignore_case=False, file_types=None):
    if not os.path.isdir(start_dir):
        print(f"Error: Directory '{start_dir}' not found.")
        return

    search_string_processed = search_string.lower() if ignore_case else search_string

    if file_types:
        file_types_processed = [f".{ft.lstrip('.').lower()}" for ft in file_types]
    else:
        file_types_processed = None

    for root, _, files in os.walk(start_dir):
        if not recursive and root != start_dir:
            continue

        for filename in files:
            file_path = os.path.join(root, filename)

            if file_types_processed:
                file_ext = os.path.splitext(filename)[1].lower()
                if file_ext not in file_types_processed:
                    continue

            try:
                with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
                    for line_num, line in enumerate(f, 1):
                        line_processed = line.lower() if ignore_case else line
                        if search_string_processed in line_processed:
                            print(f"Found in: {file_path} (Line {line_num}): {line.strip()}")
            except IOError as e:
                print(f"Error reading file {file_path}: {e}")

def main():
    parser = argparse.ArgumentParser(
        description="Finds files containing a specific string with additional functionality."
    )
    parser.add_argument(
        "directory",
        type=str,
        help="The directory to start searching from."
    )
    parser.add_argument(
        "search_string",
        type=str,
        help="The string to search for."
    )
    parser.add_argument(
        "-r", "--recursive",
        action="store_true",
        help="Search subdirectories recursively."
    )
    parser.add_argument(
        "-i", "--ignore-case",
        action="store_true",
        help="Perform a case-insensitive search."
    )
    parser.add_argument(
        "-f", "--file-types",
        type=str,
        help="Comma-separated list of file extensions to include (e.g., 'txt,py,md')."
    )

    args = parser.parse_args()

    file_types_list = None
    if args.file_types:
        file_types_list = [ft.strip() for ft in args.file_types.split(',')]

    find_string_in_files(
        start_dir=args.directory,
        search_string=args.search_string,
        recursive=args.recursive,
        ignore_case=args.ignore_case,
        file_types=file_types_list
    )

if __name__ == "__main__":
    main()