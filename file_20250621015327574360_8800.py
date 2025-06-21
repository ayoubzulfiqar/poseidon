import os
import shutil
import time

def sync_folders(source_folder, destination_folder):
    if not os.path.exists(source_folder):
        return

    os.makedirs(destination_folder, exist_ok=True)

    for root, dirs, files in os.walk(source_folder):
        relative_path = os.path.relpath(root, source_folder)
        destination_dir = os.path.join(destination_folder, relative_path)

        os.makedirs(destination_dir, exist_ok=True)

        for file_name in files:
            source_file_path = os.path.join(root, file_name)
            destination_file_path = os.path.join(destination_dir, file_name)

            if not os.path.exists(destination_file_path) or \
               os.path.getmtime(source_file_path) > os.path.getmtime(destination_file_path):
                try:
                    shutil.copy2(source_file_path, destination_file_path)
                except Exception:
                    pass

    for root, dirs, files in os.walk(destination_folder, topdown=False):
        relative_path = os.path.relpath(root, destination_folder)
        source_dir = os.path.join(source_folder, relative_path)

        for file_name in files:
            destination_file_path = os.path.join(root, file_name)
            source_file_path = os.path.join(source_dir, file_name)
            if not os.path.exists(source_file_path):
                try:
                    os.remove(destination_file_path)
                except Exception:
                    pass

        for dir_name in dirs:
            destination_sub_dir_path = os.path.join(root, dir_name)
            source_sub_dir_path = os.path.join(source_dir, dir_name)
            if not os.path.exists(source_sub_dir_path):
                try:
                    shutil.rmtree(destination_sub_dir_path)
                except Exception:
                    pass

if __name__ == "__main__":
    test_source = "source_folder_sync_test"
    test_destination = "destination_folder_sync_test"

    if os.path.exists(test_source):
        shutil.rmtree(test_source)
    if os.path.exists(test_destination):
        shutil.rmtree(test_destination)

    os.makedirs(os.path.join(test_source, "sub_dir1"), exist_ok=True)
    os.makedirs(os.path.join(test_source, "sub_dir2"), exist_ok=True)
    with open(os.path.join(test_source, "file1.txt"), "w") as f:
        f.write("This is file1.")
    with open(os.path.join(test_source, "sub_dir1", "file2.txt"), "w") as f:
        f.write("This is file2.")
    with open(os.path.join(test_source, "sub_dir2", "file3.txt"), "w") as f:
        f.write("This is file3.")
    with open(os.path.join(test_source, "file_to_update.txt"), "w") as f:
        f.write("Initial content.")

    print("--- Initial Sync ---")
    sync_folders(test_source, test_destination)
    print("Initial sync complete. Check 'destination_folder_sync_test'.")

    print("\n--- Test Case 1: Add new file to source ---")
    with open(os.path.join(test_source, "new_file_source.txt"), "w") as f:
        f.write("This is a new file.")
    sync_folders(test_source, test_destination)
    print("New file added and synced. Check 'destination_folder_sync_test/new_file_source.txt'.")

    print("\n--- Test Case 2: Update existing file in source ---")
    time.sleep(0.1)
    with open(os.path.join(test_source, "file_to_update.txt"), "w") as f:
        f.write("Updated content.")
    sync_folders(test_source, test_destination)
    print("File updated and synced. Check 'destination_folder_sync_test/file_to_update.txt'.")

    print("\n--- Test Case 3: Remove file from source ---")
    os.remove(os.path.join(test_source, "file1.txt"))
    sync_folders(test_source, test_destination)
    print("File removed from source and synced. Check 'destination_folder_sync_test/file1.txt' should be gone.")

    print("\n--- Test Case 4: Remove directory from source ---")
    shutil.rmtree(os.path.join(test_source, "sub_dir2"))
    sync_folders(test_source, test_destination)
    print("Directory removed from source and synced. Check 'destination_folder_sync_test/sub_dir2' should be gone.")

    print("\n--- Test Case 5: Add file to destination (should be removed) ---")
    with open(os.path.join(test_destination, "extra_file_dest.txt"), "w") as f:
        f.write("This file should be removed.")
    sync_folders(test_source, test_destination)
    print("Extra file in destination should be removed. Check 'destination_folder_sync_test/extra_file_dest.txt' should be gone.")

    print("\n--- Test Case 6: Add directory to destination (should be removed) ---")
    os.makedirs(os.path.join(test_destination, "extra_dir_dest"), exist_ok=True)
    with open(os.path.join(test_destination, "extra_dir_dest", "extra_file_in_extra_dir.txt"), "w") as f:
        f.write("This file should be removed with its dir.")
    sync_folders(test_source, test_destination)
    print("Extra directory in destination should be removed. Check 'destination_folder_sync_test/extra_dir_dest' should be gone.")

    print("\n--- All test cases complete. Cleaning up. ---")
    if os.path.exists(test_source):
        shutil.rmtree(test_source)
    if os.path.exists(test_destination):
        shutil.rmtree(test_destination)
    print("Cleanup complete.")

# Additional implementation at 2025-06-21 01:54:14
import os
import shutil
import hashlib
import argparse
import logging
import fnmatch

logger = None

def setup_logging(log_file=None):
    global logger
    logger = logging.getLogger("folder_sync")
    logger.setLevel(logging.INFO)

    ch = logging.StreamHandler()
    ch.setLevel(logging.INFO)
    formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
    ch.setFormatter(formatter)
    logger.addHandler(ch)

    if log_file:
        fh = logging.FileHandler(log_file)
        fh.setLevel(logging.INFO)
        fh.setFormatter(formatter)
        logger.addHandler(fh)

def get_file_checksum(filepath, block_size=65536):
    hasher = hashlib.md5()
    try:
        with open(filepath, 'rb') as f:
            buf = f.read(block_size)
            while len(buf) > 0:
                hasher.update(buf)
                buf = f.read(block_size)
        return hasher.hexdigest()
    except OSError as e:
        logger.warning(f"Could not read file for checksum '{filepath}': {e}")
        return None

def get_file_info(filepath, use_checksum):
    try:
        stat = os.stat(filepath)
        info = {
            'mtime': stat.st_mtime,
            'size': stat.st_size,
            'checksum': None
        }
        if use_checksum:
            info['checksum'] = get_file_checksum(filepath)
        return info
    except OSError as e:
        logger.warning(f"Could not get info for '{filepath}': {e}")
        return None

def should_exclude(path, exclude_patterns):
    basename = os.path.basename(path)
    for pattern in exclude_patterns:
        if fnmatch.fnmatch(basename, pattern):
            return True
    return False

def sync_folders(source_dir, dest_dir, dry_run, exclude_patterns, use_checksum):
    logger.info(f"Starting sync from '{source_dir}' to '{dest_dir}' (Dry-run: {dry_run}, Checksum: {use_checksum})")

    if not os.path.exists(source_dir):
        logger.error(f"Source directory '{source_dir}' does not exist.")
        return
    if not os.path.exists(dest_dir):
        logger.info(f"Destination directory '{dest_dir}' does not exist. Creating it.")
        if not dry_run:
            try:
                os.makedirs(dest_dir)
            except OSError as e:
                logger.error(f"Failed to create destination directory '{dest_dir}': {e}")
                return

    source_files = {}
    dest_files = {}
    source_dirs = set()
    dest_dirs = set()

    for root, dirs, files in os.walk(source_dir):
        dirs[:] = [d for d in dirs if not should_exclude(os.path.join(root, d), exclude_patterns)]
        for d in dirs:
            rel_path = os.path.relpath(os.path.join(root, d), source_dir)
            if not should_exclude(rel_path, exclude_patterns):
                source_dirs.add(rel_path)

        for file in files:
            full_path = os.path.join(root, file)
            rel_path = os.path.relpath(full_path, source_dir)
            if not should_exclude(rel_path, exclude_patterns):
                info = get_file_info(full_path, use_checksum)
                if info:
                    source_files[rel_path] = info

    for root, dirs, files in os.walk(dest_dir):
        dirs[:] = [d for d in dirs if not should_exclude(os.path.join(root, d), exclude_patterns)]
        for d in dirs:
            rel_path = os.path.relpath(os.path.join(root, d), dest_dir)
            if not should_exclude(rel_path, exclude_patterns):
                dest_dirs.add(rel_path)

        for file in files:
            full_path = os.path.join(root, file)
            rel_path = os.path.relpath(full_path, dest_dir)
            if not should_exclude(rel_path, exclude_patterns):
                info = get_file_info(full_path, use_checksum)
                if info:
                    dest_files[rel_path] = info

    files_to_copy = []
    files_to_delete = []

    for rel_path, src_info in source_files.items():
        dest_info = dest_files.get(rel_path)
        src_full_path = os.path.join(source_dir, rel_path)
        dest_full_path = os.path.join(dest_dir, rel_path)

        if not dest_info:
            files_to_copy.append((src_full_path, dest_full_path, "New file"))
        else:
            if use_checksum:
                if src_info['checksum'] != dest_info['checksum']:
                    files_to_copy.append((src_full_path, dest_full_path, "Checksum mismatch"))
            else:
                if src_info['mtime'] > dest_info['mtime'] or src_info['size'] != dest_info['size']:
                    files_to_copy.append((src_full_path, dest_full_path, "Modified (mtime/size)"))

    for rel_path in dest_files:
        if rel_path not in source_files:
            files_to_delete.append(os.path.join(dest_dir, rel_path))

    for src, dest, reason in files_to_copy:
        logger.info(f"COPY: {reason}: '{src}' to '{dest}'")
        if not dry_run:
            dest_parent_dir = os.path.dirname(dest)
            if not os.path.exists(dest_parent_dir):
                logger.info(f"Creating directory: '{dest_parent_dir}'")
                try:
                    os.makedirs(dest_parent_dir, exist_ok=True)
                except OSError as e:
                    logger.error(f"Error creating directory '{dest_parent_dir}': {e}")
                    continue
            try:
                shutil.copy2(src, dest)
            except Exception as e:
                logger.error(f"Error copying '{src}' to '{dest}': {e}")

    for dest_file in files_to_delete:
        logger.info(f"DELETE: '{dest_file}'")
        if not dry_run:
            try:
                os.remove(dest_file)
            except Exception as e:
                logger.error(f"Error deleting '{dest_file}': {e}")

    dest_dirs_sorted = sorted(list(dest_dirs), key=lambda x: x.count(os.sep), reverse=True)
    for rel_path in dest_dirs_sorted:
        full_path = os.path.join(dest_dir, rel_path)
        if not os.path.exists(os.path.join(source_dir, rel_path)):
            if os.path.isdir(full_path):
                try:
                    if not os.listdir(full_path):
                        logger.info(f"DELETE EMPTY DIR: '{full_path}'")
                        if not dry_run:
                            os.rmdir(full_path)
                except OSError as e:
                    logger.error(f"Error deleting directory '{full_path}': {e}")

    logger.info("Sync complete.")

def main():
    parser = argparse.ArgumentParser(description="One-way folder synchronization script.")
    parser.add_argument("source", help="Source directory path.")
    parser.add_argument("destination", help="Destination directory path.")
    parser.add_argument("--dry-run", action="store_true",
                        help="Perform a dry run without making any changes.")
    parser.add_argument("--log-file", help="Path to a log file for output.")
    parser.add_argument("--exclude", nargs='*', default=[],
                        help="List of file/folder name patterns to exclude (e.g., '*.log', 'temp_folder'). Uses fnmatch.")
    parser.add_argument("--checksum", action="store_true",
                        help="Use MD5 checksums for file comparison instead of modification time and size.")

    args = parser.parse_args()

    setup_logging(args.log_file)

    sync_folders(args.source, args.destination, args.dry_run, args.exclude, args.checksum)

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-21 01:54:58
import os
import shutil
import logging
import argparse
import sys

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

def setup_logging(log_file=None):
    if log_file:
        file_handler = logging.FileHandler(log_file)
        file_handler.setFormatter(logging.Formatter('%(asctime)s - %(levelname)s - %(message)s'))
        logger.addHandler(file_handler)

def sync_folders(source_path, destination_path, dry_run=False):
    source_path = os.path.abspath(source_path)
    destination_path = os.path.abspath(destination_path)

    if not os.path.exists(source_path):
        logger.error(f"Source path '{source_path}' does not exist.")
        return

    if not os.path.isdir(source_path):
        logger.error(f"Source path '{source_path}' is not a directory.")
        return

    logger.info(f"Starting one-way sync from '{source_path}' to '{destination_path}'")
    if dry_run:
        logger.info("DRY RUN: No changes will be made.")

    logger.info("Phase 1: Copying/Updating files and creating directories...")
    try:
        if not os.path.exists(destination_path):
            if dry_run:
                logger.info(f"DRY RUN: Would create destination directory '{destination_path}'")
            else:
                os.makedirs(destination_path)
                logger.info(f"Created destination directory '{destination_path}'")

        for dirpath, dirnames, filenames in os.walk(source_path):
            relative_path = os.path.relpath(dirpath, source_path)
            dest_dirpath = os.path.join(destination_path, relative_path)

            if not os.path.exists(dest_dirpath):
                if dry_run:
                    logger.info(f"DRY RUN: Would create directory '{dest_dirpath}'")
                else:
                    os.makedirs(dest_dirpath)
                    logger.info(f"Created directory '{dest_dirpath}'")

            for filename in filenames:
                src_file = os.path.join(dirpath, filename)
                dest_file = os.path.join(dest_dirpath, filename)

                should_copy = False
                if not os.path.exists(dest_file):
                    should_copy = True
                    action = "Copying new file"
                else:
                    src_mtime = os.path.getmtime(src_file)
                    dest_mtime = os.path.getmtime(dest_file)
                    if src_mtime > dest_mtime:
                        should_copy = True
                        action = "Updating modified file"
                    elif os.path.getsize(src_file) != os.path.getsize(dest_file):
                        should_copy = True
                        action = "Updating file (size mismatch)"

                if should_copy:
                    if dry_run:
                        logger.info(f"DRY RUN: {action}: '{src_file}' to '{dest_file}'")
                    else:
                        try:
                            shutil.copy2(src_file, dest_file)
                            logger.info(f"{action}: '{src_file}' to '{dest_file}'")
                        except Exception as e:
                            logger.error(f"Error {action.lower()} '{src_file}': {e}")

    except Exception as e:
        logger.error(f"Error during Phase 1 (copy/update): {e}")
        return

    logger.info("Phase 2: Deleting extra files/directories in destination...")
    try:
        source_items = set()
        for root, dirs, files in os.walk(source_path):
            relative_root = os.path.relpath(root, source_path)
            if relative_root == ".":
                relative_root = ""
            for d in dirs:
                source_items.add(os.path.normpath(os.path.join(relative_root, d)))
            for f in files:
                source_items.add(os.path.normpath(os.path.join(relative_root, f)))

        dest_items_to_check = []
        for root, dirs, files in os.walk(destination_path):
            relative_root = os.path.relpath(root, destination_path)
            if relative_root == ".":
                relative_root = ""

            for d in dirs:
                full_dest_path = os.path.join(root, d)
                relative_dest_path = os.path.normpath(os.path.join(relative_root, d))
                if relative_dest_path not in source_items:
                    dest_items_to_check.append((full_dest_path, "dir"))

            for f in files:
                full_dest_path = os.path.join(root, f)
                relative_dest_path = os.path.normpath(os.path.join(relative_root, f))
                if relative_dest_path not in source_items:
                    dest_items_to_check.append((full_dest_path, "file"))

        dest_items_to_check.sort(key=lambda x: len(x[0]), reverse=True)

        for item_path, item_type in dest_items_to_check:
            if not os.path.exists(item_path):
                continue

            if item_type == "file":
                if dry_run:
                    logger.info(f"DRY RUN: Would delete file '{item_path}'")
                else:
                    try:
                        os.remove(item_path)
                        logger.info(f"Deleted file '{item_path}'")
                    except Exception as e:
                        logger.error(f"Error deleting file '{item_path}': {e}")
            elif item_type == "dir":
                if dry_run:
                    logger.info(f"DRY RUN: Would delete directory '{item_path}'")
                else:
                    try:
                        shutil.rmtree(item_path)
                        logger.info(f"Deleted directory '{item_path}'")
                    except Exception as e:
                        logger.error(f"Error deleting directory '{item_path}': {e}")

    except Exception as e:
        logger.error(f"Error during Phase 2 (deletion): {e}")

    logger.info("Synchronization complete.")

def main():
    parser = argparse.ArgumentParser(
        description="One-way mirror synchronization of two folders.",
        formatter_class=argparse.RawTextHelpFormatter
    )
    parser.add_argument("source", help="Path to the source folder.")
    parser.add_argument("destination", help="Path to the destination folder.")
    parser.add_argument("-d", "--dry-run", action="store_true",
                        help="Perform a dry run (simulate changes without making them).")
    parser.add_argument("-l", "--log-file", type=str,
                        help="Path to a log file to output sync details.")

    args = parser.parse_args()

    setup_logging(args.log_file)

    sync_folders(args.source, args.destination, args.dry_run)

if __name__ == "__main__":
    main()