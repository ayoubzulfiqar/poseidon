import os
import argparse

def count_lines_in_file(filepath):
    line_count = 0
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            for line in f:
                stripped_line = line.strip()
                if not stripped_line:
                    continue
                if stripped_line.startswith('#'):
                    continue
                line_count += 1
    except Exception:
        pass
    return line_count

def main():
    parser = argparse.ArgumentParser(
        description="Counts lines of code in Python files."
    )
    parser.add_argument(
        'path',
        type=str,
        help="Path to a Python file or a directory containing Python files."
    )

    args = parser.parse_args()
    target_path = args.path

    total_lines = 0
    processed_files = 0

    if os.path.isfile(target_path):
        if target_path.endswith('.py'):
            lines = count_lines_in_file(target_path)
            total_lines += lines
            processed_files += 1
        else:
            print("Error: Target is not a Python file.")
            return
    elif os.path.isdir(target_path):
        for root, _, files in os.walk(target_path):
            for file in files:
                if file.endswith('.py'):
                    filepath = os.path.join(root, file)
                    lines = count_lines_in_file(filepath)
                    total_lines += lines
                    processed_files += 1
    else:
        print("Error: Path does not exist or is not a file/directory.")
        return

    print(f"Total Python files processed: {processed_files}")
    print(f"Total lines of code: {total_lines}")

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-18 02:05:33
import os
import argparse

def _is_blank_line(line):
    """Checks if a line is blank."""
    return not line.strip()

def _is_comment_line(line):
    """Checks if a line is a comment."""
    stripped_line = line.strip()
    return stripped_line.startswith('#')

def count_loc_in_file(filepath):
    """
    Counts lines of code, blank lines, and comment lines in a single Python file.

    Returns a tuple: (total_lines, code_lines, blank_lines, comment_lines)
    """
    total_lines = 0
    code_lines = 0
    blank_lines = 0
    comment_lines = 0

    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            for line in f:
                total_lines += 1
                if _is_blank_line(line):
                    blank_lines += 1
                elif _is_comment_line(line):
                    comment_lines += 1
                else:
                    code_lines += 1
    except Exception as e:
        print(f"Error reading {filepath}: {e}")
        return 0, 0, 0, 0 # Return zeros on error

    return total_lines, code_lines, blank_lines, comment_lines

def count_loc_in_directory(directory, exclude_dirs=None):
    """
    Counts lines of code in all .py files within a directory and its subdirectories.

    Args:
        directory (str): The path to the directory to scan.
        exclude_dirs (list): A list of directory names to exclude from the scan.

    Returns:
        A dictionary with file-specific counts and a total summary.
        Example:
        {
            'file_counts': {
                'path/to/file1.py': {'total': N, 'code': N, 'blank': N, 'comment': N},
                ...
            },
            'summary': {
                'total_files': N,
                'total_lines': N,
                'total_code_lines': N,
                'total_blank_lines': N,
                'total_comment_lines': N
            }
        }
    """
    if exclude_dirs is None:
        exclude_dirs = []

    file_counts = {}
    total_files = 0
    total_lines = 0
    total_code_lines = 0
    total_blank_lines = 0
    total_comment_lines = 0

    for root, dirs, files in os.walk(directory):
        # Modify dirs in-place to exclude unwanted directories
        dirs[:] = [d for d in dirs if d not in exclude_dirs]

        for file in files:
            if file.endswith('.py'):
                filepath = os.path.join(root, file)
                total, code, blank, comment = count_loc_in_file(filepath)

                file_counts[filepath] = {
                    'total': total,
                    'code': code,
                    'blank': blank,
                    'comment': comment
                }

                total_files += 1
                total_lines += total
                total_code_lines += code
                total_blank_lines += blank
                total_comment_lines += comment

    summary = {
        'total_files': total_files,
        'total_lines': total_lines,
        'total_code_lines': total_code_lines,
        'total_blank_lines': total_blank_lines,
        'total_comment_lines': total_comment_lines
    }

    return {'file_counts': file_counts, 'summary': summary}

def main():
    parser = argparse.ArgumentParser(
        description="Count lines of code in Python files.",
        formatter_class=argparse.RawTextHelpFormatter
    )
    parser.add_argument(
        'path',
        nargs='?',
        default='.',
        help="Path to a Python file or a directory. Defaults to current directory."
    )
    parser.add_argument(
        '-v', '--verbose',
        action='store_true',
        help="Show detailed line counts for each file."
    )
    parser.add_argument(
        '-e', '--exclude',
        nargs='*',
        default=['.git', '__pycache__', 'venv', 'env', 'node_modules', '.vscode'],
        help="""Directory names to exclude from the scan.
Defaults: .git __pycache__ venv env node_modules .vscode
To exclude none, use --exclude (without arguments)."""
    )

    args = parser.parse_args()

    target_path = args.path
    exclude_dirs = args.exclude if args.exclude is not None else []

    if os.path.isfile(target_path):
        if not target_path.endswith('.py'):
            print(f"Error: '{target_path}' is not a Python file.")
            return

        total, code, blank, comment = count_loc_in_file(target_path)
        print(f"File: {target_path}")
        print(f"  Total Lines: {total}")
        print(f"  Code Lines: {code}")
        print(f"  Blank Lines: {blank}")
        print(f"  Comment Lines: {comment}")

    elif os.path.isdir(target_path):
        print(f"Scanning directory: {target_path}")
        print(f"Excluding directories: {', '.join(exclude_dirs) if exclude_dirs else 'None'}\n")

        results = count_loc_in_directory(target_path, exclude_dirs)
        file_counts = results['file_counts']
        summary = results['summary']

        if args.verbose:
            if not file_counts:
                print("No Python files found.")
            else:
                for filepath, counts in file_counts.items():
                    print(f"  File: {filepath}")
                    print(f"    Total: {counts['total']}, Code: {counts['code']}, Blank: {counts['blank']}, Comment: {counts['comment']}")
                print("-" * 50)

        print("--- Summary ---")
        print(f"Total Python Files Scanned: {summary['total_files']}")
        print(f"Total Lines: {summary['total_lines']}")
        print(f"Total Code Lines: {summary['total_code_lines']}")
        print(f"Total Blank Lines: {summary['total_blank_lines']}")
        print(f"Total Comment Lines: {summary['total_comment_lines']}")

    else:
        print(f"Error: Path '{target_path}' does not exist or is not a file/directory.")

if __name__ == '__main__':
    main()

# Additional implementation at 2025-06-18 02:06:00
import os
import argparse

def count_lines_in_file(filepath, ignore_blank, ignore_comments):
    loc_count = 0
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            for line in f:
                stripped_line = line.strip()
                
                if not stripped_line:
                    if ignore_blank:
                        continue
                
                if stripped_line.startswith('#'):
                    if ignore_comments:
                        continue
                
                loc_count += 1
    except Exception:
        return 0
    return loc_count

def is_excluded(path, exclude_paths):
    abs_path = os.path.abspath(path)
    for excluded in exclude_paths:
        abs_excluded = os.path.abspath(excluded)
        if abs_path == abs_excluded or abs_path.startswith(abs_excluded + os.sep):
            return True
    return False

def count_loc(target_path, recursive, ignore_blank, ignore_comments, exclude_paths):
    total_files = 0
    total_loc = 0
    results = []

    if not os.path.exists(target_path):
        return 0, 0, []

    if os.path.isfile(target_path):
        if target_path.endswith('.py') and not is_excluded(target_path, exclude_paths):
            loc = count_lines_in_file(target_path, ignore_blank, ignore_comments)
            results.append((target_path, loc))
            total_files += 1
            total_loc += loc
    elif os.path.isdir(target_path):
        if is_excluded(target_path, exclude_paths):
            return 0, 0, []

        if recursive:
            for root, dirs, files in os.walk(target_path):
                dirs[:] = [d for d in dirs if not is_excluded(os.path.join(root, d), exclude_paths)]

                for file in files:
                    if file.endswith('.py'):
                        filepath = os.path.join(root, file)
                        if not is_excluded(filepath, exclude_paths):
                            loc = count_lines_in_file(filepath, ignore_blank, ignore_comments)
                            results.append((filepath, loc))
                            total_files += 1
                            total_loc += loc
        else:
            for entry in os.listdir(target_path):
                filepath = os.path.join(target_path, entry)
                if os.path.isfile(filepath) and filepath.endswith('.py') and not is_excluded(filepath, exclude_paths):
                    loc = count_lines_in_file(filepath, ignore_blank, ignore_comments)
                    results.append((filepath, loc))
                    total_files += 1
                    total_loc += loc
    
    return total_files, total_loc, results

def main():
    parser = argparse.ArgumentParser(
        description="A command-line tool to count lines of code in Python files.",
        formatter_class=argparse.RawTextHelpFormatter
    )
    parser.add_argument(
        'path',
        nargs='?',
        default='.',
        help="The file or directory to scan. Defaults to current directory ('.')."
    )
    parser.add_argument(
        '-r', '--recursive',
        action='store_true',
        help="Recursively scan subdirectories for Python files."
    )
    parser.add_argument(
        '-b', '--ignore-blank',
        action='store_true',
        help="Ignore blank lines when counting lines of code."
    )
    parser.add_argument(
        '-c', '--ignore-comments',
        action='store_true',
        help="Ignore lines starting with '#' (single-line comments) when counting."
    )
    parser.add_argument(
        '-e', '--exclude',
        nargs='*',
        default=[],
        help="Space-separated list of files or directories to exclude from scanning.\n"
             "Example: --exclude my_script.py tests/ temp/"
    )

    args = parser.parse_args()

    absolute_exclude_paths = [os.path.abspath(p) for p in args.exclude]

    total_files, total_loc, results = count_loc(
        args.path,
        args.recursive,
        args.ignore_blank,
        args.ignore_comments,
        absolute_exclude_paths
    )

    if results:
        print("--- LOC Details ---")
        for filepath, loc in sorted(results):
            print(f"{filepath}: {loc} LOC")
        print("\n--- Summary ---")
        print(f"Total Python Files: {total_files}")
        print(f"Total Lines of Code: {total_loc}")

if __name__ == '__main__':
    main()