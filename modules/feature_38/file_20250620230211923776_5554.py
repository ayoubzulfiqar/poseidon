import psutil
import time

def kill_processes_by_name(process_name):
    killed_count = 0
    for proc in psutil.process_iter(['pid', 'name']):
        try:
            if proc.info['name'] == process_name:
                proc.terminate()
                try:
                    proc.wait(timeout=3)
                    killed_count += 1
                except psutil.TimeoutExpired:
                    proc.kill()
                    proc.wait(timeout=3)
                    killed_count += 1
                except psutil.NoSuchProcess:
                    killed_count += 1
        except (psutil.NoSuchProcess, psutil.AccessDenied, psutil.ZombieProcess):
            pass
        except Exception:
            pass
    return killed_count

if __name__ == "__main__":
    process_to_kill = "non_existent_process_for_test.exe"
    killed_count = kill_processes_by_name(process_to_kill)
    if killed_count > 0:
        print(f"Terminated {killed_count} instance(s) of '{process_to_kill}'.")
    else:
        print(f"No processes named '{process_to_kill}' found or terminated.")

# Additional implementation at 2025-06-20 23:03:29
import psutil
import argparse
import sys

def kill_processes_by_name(process_names, force=False, list_only=False):
    """
    Kills processes by name.

    Args:
        process_names (list): A list of process names to kill.
        force (bool): If True, skips confirmation prompt.
        list_only (bool): If True, only lists matching processes without killing.
    """
    found_processes = []
    target_names_lower = [name.lower() for name in process_names]

    print(f"Searching for processes: {', '.join(process_names)}")

    for proc in psutil.process_iter(['pid', 'name', 'cmdline']):
        try:
            proc_name = proc.info['name']
            if proc_name and proc_name.lower() in target_names_lower:
                found_processes.append(proc)
        except (psutil.NoSuchProcess, psutil.AccessDenied, psutil.ZombieProcess):
            # Process might have terminated or become inaccessible during iteration
            continue

    if not found_processes:
        print(f"No processes found matching: {', '.join(process_names)}")
        return

    print("\nFound the following processes:")
    for proc in found_processes:
        try:
            cmdline = ' '.join(proc.info['cmdline']) if proc.info['cmdline'] else proc.info['name']
            print(f"  PID: {proc.info['pid']}, Name: {proc.info['name']}, Command: {cmdline}")
        except (psutil.NoSuchProcess, psutil.AccessDenied, psutil.ZombieProcess):
            print(f"  PID: {proc.info['pid']} (Details unavailable, process might have exited)")
            continue

    if list_only:
        print("\nList-only mode: No processes were killed.")
        return

    if not force:
        confirmation = input("\nDo you want to kill these processes? (yes/no): ").lower()
        if confirmation != 'yes':
            print("Operation cancelled.")
            return

    print("\nAttempting to kill processes...")
    killed_count = 0
    for proc in found_processes:
        try:
            proc_name = proc.info['name']
            proc_pid = proc.info['pid']
            print(f"  Killing process: PID {proc_pid}, Name: {proc_name}...", end=" ")
            proc.kill()
            print("SUCCESS")
            killed_count += 1
        except psutil.NoSuchProcess:
            print("FAILED (Process already terminated)")
        except psutil.AccessDenied:
            print("FAILED (Permission denied. Try running as administrator/root.)")
        except Exception as e:
            print(f"FAILED (An unexpected error occurred: {e})")

    print(f"\nFinished. Successfully killed {killed_count} out of {len(found_processes)} matching processes.")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="A script to kill processes by name, with options for multiple names, confirmation, and dry run."
    )
    parser.add_argument(
        '-n', '--name',
        nargs='+',
        required=True,
        help='One or more process names to kill (e.g., "chrome.exe" "notepad.exe"). Case-insensitive.'
    )
    parser.add_argument(
        '-f', '--force',
        action='store_true',
        help='Skip confirmation prompt and kill processes immediately.'
    )
    parser.add_argument(
        '-l', '--list',
        action='store_true',
        help='Only list matching processes without killing them (dry run).'
    )

    args = parser.parse_args()

    if not args.name:
        print("Error: Please provide at least one process name using -n or --name.")
        parser.print_help()
        sys.exit(1)

    kill_processes_by_name(args.name, args.force, args.list)