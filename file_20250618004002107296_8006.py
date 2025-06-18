import os
import time

def update_file_timestamps(file_path, new_access_time=None, new_modification_time=None):
    if not os.path.exists(file_path):
        return

    if new_access_time is None and new_modification_time is None:
        os.utime(file_path, None)
    else:
        stat_info = os.stat(file_path)
        
        atime_to_set = new_access_time if new_access_time is not None else stat_info.st_atime
        mtime_to_set = new_modification_time if new_modification_time is not None else stat_info.st_mtime
        
        os.utime(file_path, (atime_to_set, mtime_to_set))

# Additional implementation at 2025-06-18 00:41:18
import os
import argparse
import datetime
import time
import sys

try:
    import win32api
    import win32con
    import pywintypes
    _HAS_WIN32 = True
except ImportError:
    _HAS_WIN32 = False

def _set_windows_creation_time(filepath, new_time_dt):
    if not _HAS_WIN32:
        return False

    try:
        ft = pywintypes.Time(new_time_dt)

        handle = win32api.CreateFile(
            filepath,
            win32con.GENERIC_WRITE,
            win32con.FILE_SHARE_READ | win32con.FILE_SHARE_WRITE | win32con.FILE_SHARE_DELETE,
            None,
            win32con.OPEN_EXISTING,
            win32con.FILE_ATTRIBUTE_NORMAL,
            None
        )

        # Get current times to preserve access and modification if not explicitly set
        # This is a more robust way than setting all three to `ft`
        # However, for this specific function, we assume `new_time_dt` is the desired creation time,
        # and `os.utime` will handle access/modification.
        # So, we only set creation time, leaving access/modification as they are.
        # But SetFileTime requires all three. So we read current ones.
        creation_time, access_time, modification_time = win32api.GetFileTime(handle)
        win32api.SetFileTime(handle, ft, access_time, modification_time)
        win32api.CloseHandle(handle)
        return True
    except Exception:
        return False

def update_file_timestamps(
    filepath,
    target_time=None,
    copy_from_file=None,
    no_access_time=False,
    no_modification_time=False,
    no_creation_time=False,
    preview_mode=False
):
    if not os.path.exists(filepath):
        return False

    atime_val = None
    mtime_val = None
    ctime_dt_val = None # For Windows creation time (datetime object)

    if copy_from_file:
        if not os.path.exists(copy_from_file):
            return False
        stat_info = os.stat(copy_from_file)
        atime_val = stat_info.st_atime
        mtime_val = stat_info.st_mtime
        if sys.platform == "win32" and _HAS_WIN32:
            ctime_dt_val = datetime.datetime.fromtimestamp(stat_info.st_ctime)
        elif hasattr(stat_info, 'st_birthtime'):
            ctime_dt_val = datetime.datetime.fromtimestamp(stat_info.st_birthtime)
        else:
            ctime_dt_val = None

    elif target_time:
        atime_val = target_time.timestamp()
        mtime_val = target_time.timestamp()
        ctime_dt_val = target_time

    else:
        current_time_ts = time.time()
        atime_val = current_time_ts
        mtime_val = current_time_ts
        ctime_dt_val = datetime.datetime.fromtimestamp(current_time_ts)

    original_stat = os.stat(filepath)
    original_atime = original_stat.st_atime
    original_mtime = original_stat.st_mtime

    final_atime = original_atime if no_access_time else atime_val
    final_mtime = original_mtime if no_modification_time else mtime_val

    if not preview_mode:
        try:
            os.utime(filepath, (final_atime, final_mtime))

            if sys.platform == "win32" and _HAS_WIN32 and not no_creation_time and ctime_dt_val:
                _set_windows_creation_time(filepath, ctime_dt_val)

            return True
        except Exception:
            return False
    else:
        return True

def parse_datetime_string(dt_str):
    formats = [
        "%Y-%m-%d %H:%M:%S",
        "%Y-%m-%d %H:%M",
        "%Y-%m-%d",
        "%Y/%m/%d %H:%M:%S",
        "%Y/%m/%d %H:%M",
        "%Y/%m/%d"
    ]
    for fmt in formats:
        try:
            return datetime.datetime.strptime(dt_str, fmt)
        except ValueError:
            continue
    raise ValueError("Invalid date/time format. Expected YYYY-MM-DD HH:MM:SS or similar.")

def main():
    parser = argparse.ArgumentParser(add_help=False)
    parser.add_argument('path', type=str, help='File or directory path to update.')
    parser.add_argument('-r', '--recursive', action='store_true', help='Process directories recursively.')
    parser.add_argument('-s', '--set-time', type=str, help='Set all timestamps to a specific time (YYYY-MM-DD HH:MM:SS).')
    parser.add_argument('-c', '--copy-from', type=str, help='Copy timestamps from another file.')
    parser.add_argument('--no-access-time', action='store_true', help='Do not change access time.')
    parser.add_argument('--no-modification-time', action='store_true', help='Do not change modification time.')
    parser.add_argument('--no-creation-time', action='store_true', help='Do not change creation time (Windows only).')
    parser.add_argument('-p', '--preview', action='store_true', help='Show what would be changed without making actual changes.')
    parser.add_argument('-h', '--help', action='help', default=argparse.SUPPRESS, help='Show this help message and exit.')

    args = parser.parse_args()

    target_dt = None
    if args.set_time:
        try:
            target_dt = parse_datetime_string(args.set_time)
        except ValueError as e:
            sys.stderr.write(f"Error: {e}\n")
            sys.exit(1)

    if args.copy_from and not os.path.exists(args.copy_from):
        sys.stderr.write(f"Error: Source file for copying '{args.copy_from}' does not exist.\n")
        sys.exit(1)

    if args.copy_from and args.set_time:
        sys.stderr.write("Error: Cannot use both --set-time and --copy-from simultaneously.\n")
        sys.exit(1)

    if not os.path.exists(args.path):
        sys.stderr.write(f"Error: Target path '{args.path}' does not exist.\n")
        sys.exit(1)

    if os.path.isfile(args.path):
        if update_file_timestamps(
            args.path,
            target_time=target_dt,
            copy_from_file=args.copy_from,
            no_access_time=args.no_access_time,
            no_modification_time=args.no_modification_time,
            no_creation_time=args.no_creation_time,
            preview_mode=args.preview
        ):
            sys.stdout.write(f"Processed: {args.path} (Preview: {args.preview})\n")
        else:
            sys.stderr.write(f"Failed to process: {args.path}\n")
    elif os.path.isdir(args.path):
        if not args.recursive:
            sys.stderr.write(f"Error: '{args.path}' is a directory. Use -r or --recursive to process its contents.\n")
            sys.exit(1)

        for root, _, files in os.walk(args.path):
            for filename in files:
                filepath = os.path.join(root, filename)
                if update_file_timestamps(
                    filepath,
                    target_time=target_dt,
                    copy_from_file=args.copy_from,
                    no_access_time=args.no_access_time,
                    no_modification_time=args.no_modification_time,
                    no_creation_time=args.no_creation_time,
                    preview_mode=args.preview
                ):
                    sys.stdout.write(f"Processed: {filepath} (Preview: {args.preview})\n")
                else:
                    sys.stderr.write(f"Failed to process: {filepath}\n")

if __name__ == '__main__':
    main()