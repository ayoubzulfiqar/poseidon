import os
import sys

def _print_tree_recursive(current_path, indent, is_last_sibling):
    prefix = "└── " if is_last_sibling else "├── "
    print(indent + prefix + os.path.basename(current_path))

    child_indent_base = indent + ("    " if is_last_sibling else "│   ")

    try:
        contents = sorted(os.listdir(current_path))
    except PermissionError:
        print(child_indent_base + "    [Permission Denied]")
        return
    except NotADirectoryError:
        return

    dirs = [d for d in contents if os.path.isdir(os.path.join(current_path, d))]
    files = [f for f in contents if os.path.isfile(os.path.join(current_path, f))]

    all_items = dirs + files

    for i, item in enumerate(all_items):
        item_path = os.path.join(current_path, item)
        is_last_item_in_list = (i == len(all_items) - 1)

        if os.path.isdir(item_path):
            _print_tree_recursive(item_path, child_indent_base, is_last_item_in_list)
        else:
            file_prefix = "└── " if is_last_item_in_list else "├── "
            print(child_indent_base + file_prefix + os.path.basename(item_path))

def generate_directory_tree(start_path):
    if not os.path.exists(start_path):
        print(f"Error: Path '{start_path}' does not exist.")
        return
    if not os.path.isdir(start_path):
        print(f"Error: Path '{start_path}' is not a directory.")
        return

    print(os.path.basename(os.path.abspath(start_path)))

    try:
        contents = sorted(os.listdir(start_path))
    except PermissionError:
        print("    [Permission Denied]")
        return

    dirs = [d for d in contents if os.path.isdir(os.path.join(start_path, d))]
    files = [f for f in contents if os.path.isfile(os.path.join(start_path, f))]

    all_items = dirs + files

    for i, item in enumerate(all_items):
        item_path = os.path.join(start_path, item)
        is_last_item = (i == len(all_items) - 1)
        if os.path.isdir(item_path):
            _print_tree_recursive(item_path, "│   ", is_last_item)
        else:
            prefix = "└── " if is_last_item else "├── "
            print("│   " + prefix + os.path.basename(item_path))

if __name__ == "__main__":
    target_path = "."

    if len(sys.argv) > 1:
        target_path = sys.argv[1]

    generate_directory_tree(target_path)

# Additional implementation at 2025-06-21 03:58:06
import os
import fnmatch
import shutil

def _should_exclude(name, exclude_list, exclude_patterns):
    if exclude_list and name in exclude_list:
        return True
    if exclude_patterns:
        for pattern in exclude_patterns:
            if fnmatch.fnmatch(name, pattern):
                return True
    return False

def _print_tree(current_path, indent, current_depth, max_depth,
                exclude_dirs, exclude_files,
                exclude_dir_patterns, exclude_file_patterns):

    if max_depth is not None and current_depth > max_depth:
        return

    try:
        entries = sorted(os.listdir(current_path))
    except PermissionError:
        print(f"{indent}├── [Permission Denied] {os.path.basename(current_path)}")
        return
    except FileNotFoundError:
        print(f"{indent}├── [Not Found] {os.path.basename(current_path)}")
        return

    dirs = []
    files = []

    for entry in entries:
        full_path = os.path.join(current_path, entry)
        if os.path.isdir(full_path):
            if not _should_exclude(entry, exclude_dirs, exclude_dir_patterns):
                dirs.append(entry)
        else:
            if not _should_exclude(entry, exclude_files, exclude_file_patterns):
                files.append(entry)

    all_entries = dirs + files

    for i, entry in enumerate(all_entries):
        is_last = (i == len(all_entries) - 1)
        prefix = "└── " if is_last else "├── "
        new_indent = indent + ("    " if is_last else "│   ")

        full_path = os.path.join(current_path, entry)

        if os.path.isdir(full_path):
            print(f"{indent}{prefix}{entry}/")
            _print_tree(full_path, new_indent, current_depth + 1, max_depth,
                        exclude_dirs, exclude_files, exclude_dir_patterns, exclude_file_patterns)
        else:
            print(f"{indent}{prefix}{entry}")

def generate_tree(start_path, max_depth=None, exclude_dirs=None, exclude_files=None,
                  exclude_dir_patterns=None, exclude_file_patterns=None):
    if not os.path.exists(start_path):
        print(f"Error: Path '{start_path}' does not exist.")
        return
    if not os.path.isdir(start_path):
        print(f"Error: Path '{start_path}' is not a directory.")
        return

    print(f"{os.path.basename(os.path.abspath(start_path))}/")
    _print_tree(start_path, "", 0, max_depth,
                exclude_dirs, exclude_files, exclude_dir_patterns, exclude_file_patterns)

if __name__ == "__main__":
    test_dir = "test_tree_root"
    
    if os.path.exists(test_dir):
        shutil.rmtree(test_dir)

    os.makedirs(os.path.join(test_dir, "dir1", "subdir_a"), exist_ok=True)
    os.makedirs(os.path.join(test_dir, "dir1", "subdir_b"), exist_ok=True)
    os.makedirs(os.path.join(test_dir, "dir2", "temp_dir"), exist_ok=True)
    os.makedirs(os.path.join(test_dir, ".git"), exist_ok=True)
    os.makedirs(os.path.join(test_dir, "__pycache__"), exist_ok=True)

    with open(os.path.join(test_dir, "file1.txt"), "w") as f: f.write("content")
    with open(os.path.join(test_dir, "file2.log"), "w") as f: f.write("content")
    with open(os.path.join(test_dir, "dir1", "file_a.py"), "w") as f: f.write("content")
    with open(os.path.join(test_dir, "dir1", "file_b.txt"), "w") as f: f.write("content")
    with open(os.path.join(test_dir, "dir1", "subdir_a", "nested_file.txt"), "w") as f: f.write("content")
    with open(os.path.join(test_dir, "dir2", "temp_file.tmp"), "w") as f: f.write("content")
    with open(os.path.join(test_dir, "dir2", "important.txt"), "w") as f: f.write("content")

    print("--- Full Tree (Max Depth 3) ---")
    generate_tree(test_dir, max_depth=3)
    print("\n--- Tree with Exclusions and Depth Limit ---")
    generate_tree(
        test_dir,
        max_depth=2,
        exclude_dirs=['temp_dir', '__pycache__'],
        exclude_files=['file1.txt'],
        exclude_dir_patterns=['.git'],
        exclude_file_patterns=['*.log', '*.tmp']
    )
    print("\n--- Tree with only specific file exclusion ---")
    generate_tree(
        test_dir,
        exclude_files=['file2.log', 'file_a.py']
    )
    print("\n--- Tree with only specific directory exclusion ---")
    generate_tree(
        test_dir,
        exclude_dirs=['dir1']
    )

    shutil.rmtree(test_dir)

# Additional implementation at 2025-06-21 03:59:21
import os
import shutil

def generate_tree(start_path, exclude_dirs=None, exclude_exts=None, max_depth=None):
    """
    Generates a directory tree diagram for the given path with additional functionality.

    Args:
        start_path (str): The root directory from which to generate the tree.
        exclude_dirs (list, optional): A list of directory names to exclude.
                                       E.g., ['.git', '__pycache__']. Defaults to None.
        exclude_exts (list, optional): A list of file extensions (e.g., ['.pyc', '.log']) to exclude.
                                       Leading dot is optional. Defaults to None.
        max_depth (int, optional): The maximum depth to traverse. None means no limit.
                                   Depth 0 is the start_path itself. Depth 1 includes its direct children.
                                   Defaults to None.
    """
    if not os.path.isdir(start_path):
        print(f"Error: '{start_path}' is not a valid directory.")
        return

    if exclude_dirs is None:
        exclude_dirs = []
    if exclude_exts is None:
        exclude_exts = []

    # Normalize exclude_exts to include leading dot