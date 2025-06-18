import psutil
import sys
import time

def kill_processes_by_name(process_name):
    killed_count = 0
    for proc in psutil.process_iter(['pid', 'name']):
        try:
            if proc.info['name'] == process_name:
                print(f"Attempting to terminate process {proc.info['name']} (PID: {proc.info['pid']})...")
                proc.terminate()
                try:
                    proc.wait(timeout=3)
                    print(f"Process {proc.info['name']} (PID: {proc.info['pid']}) terminated successfully.")
                    killed_count += 1
                except psutil.TimeoutExpired:
                    print(f"Process {proc.info['name']} (PID: {proc.info['pid']}) did not terminate gracefully. Forcibly killing...")
                    proc.kill()
                    proc.wait(timeout=1)
                    print(f"Process {proc.info['name']} (PID: {proc.info['pid']}) forcibly killed.")
                    killed_count += 1
        except (psutil.NoSuchProcess, psutil.AccessDenied, psutil.ZombieProcess):
            continue
        except Exception as e:
            print(f"Error processing PID {proc.info.get('pid', 'N/A')}: {e}")
            continue
    if killed_count == 0:
        print(f"No processes found with the name '{process_name}'.")
    else:
        print(f"Successfully killed {killed_count} processes with the name '{process_name}'.")

if __name__ == "__main__":
    if len(sys.argv) > 1:
        target_process_name = sys.argv[1]
    else:
        target_process_name = "non_existent_process_12345.exe"
    kill_processes_by_name(target_process_name)

# Additional implementation at 2025-06-18 00:21:46
import psutil
import sys

def kill_processes_by_name(process_name):
    """
    Kills all processes matching the given name.
    Handles various exceptions during process iteration and termination.
    """
    killed_count = 0
    print(f"Attempting to kill processes named: '{process_name}'")

    # Iterate over all running processes
    for proc in psutil.process_iter(['pid', 'name']):
        try:
            # Check if the process name matches the target name
            if proc.info['name'] == process_name:
                pid = proc.info['pid']
                current_process_name = proc.info['name']
                print(f"Found process '{current_process_name}' with PID {pid}. Attempting to kill...")
                
                # Terminate the process
                proc.kill() 
                print(f"Successfully killed process '{current_process_name}' with PID {pid}.")
                killed_count += 1
        except psutil.NoSuchProcess:
            # Process might have terminated between iteration and access
            print(f"Process with PID {proc.info['pid']} no longer exists (already terminated).")
        except psutil.AccessDenied:
            # Permission issues, usually requires elevated privileges
            print(f"Permission denied to kill process '{proc.info['name']}' with PID {proc.info['pid']}. Try running as administrator/root.")
        except psutil.ZombieProcess:
            # Process is a zombie, cannot be killed directly
            print(f"Process '{proc.info['name']}' with PID {proc.info['pid']} is a zombie process and cannot be killed directly.")
        except Exception as e:
            # Catch any other unexpected errors
            print(f"An unexpected error occurred while processing PID {proc.info['pid']}: {e}")

    if killed_count == 0:
        print(f"No processes found with the name '{process_name}'.")
    else:
        print(f"Successfully killed {killed_count} process(es) named '{process_name}'.")

if __name__ == "__main__":
    # Check if a process name is provided as a command-line argument
    if len(sys.argv) > 1:
        target_process_name = sys.argv[1]
    else:
        # If not, prompt the user for input
        target_process_name = input("Please enter the exact name of the process to kill: ").strip()

    if not target_process_name:
        print("No process name provided. Exiting.")
    else:
        kill_processes_by_name(target_process_name)

# Additional implementation at 2025-06-18 00:22:17
import psutil
import sys
import os

def kill_processes_by_name(target_name, confirm=True, list_only=False, case_sensitive=False):
    if not target_name:
        print("Error: Process name cannot be empty.")
        return

    matching_processes = []
    current_pid = os.getpid()

    print(f"Searching for processes matching '{target_name}'...")

    for proc in psutil.process_iter(['pid', 'name', 'cmdline']):
        try:
            proc_info = proc.info
            proc_pid = proc_info['pid']

            if proc_pid == current_pid:
                continue

            proc_name = proc_info['name']
            proc_cmdline = proc_info['cmdline']

            compare_target = target_name if case_sensitive else target_name.lower()
            compare_proc_name = proc_name if case_sensitive else proc_name.lower()
            compare_cmdline = ' '.join(proc_cmdline) if proc_cmdline else ''
            if not case_sensitive:
                compare_cmdline = compare_cmdline.lower()

            if compare_target in compare_proc_name or (proc_cmdline and compare_target in compare_cmdline):
                matching_processes.append(proc)

        except (psutil.NoSuchProcess, psutil.AccessDenied, psutil.ZombieProcess):
            continue

    if not matching_processes:
        print(f"No processes found matching '{target_name}'.")
        return

    print("\nFound the following processes:")
    for proc in matching_processes:
        try:
            print(f"  PID: {proc.pid}, Name: {proc.name()}, Cmdline: {' '.join(proc.cmdline())}")
        except (psutil.NoSuchProcess, psutil.AccessDenied):
            print(f"  PID: {proc.pid}, (Process info unavailable)")
            continue

    if list_only:
        print("\nList-only mode: No processes were terminated.")
        return

    if confirm:
        user_input = input(f"\nDo you want to terminate {len(matching_processes)} process(es)? (yes/no): ").lower()
        if user_input != 'yes':
            print("Operation cancelled by user.")
            return

    print("\nAttempting to terminate processes...")
    killed_count = 0
    failed_count = 0

    for proc in matching_processes:
        proc_display_name = "Unknown"
        try:
            proc_display_name = proc.name()
            proc.terminate()
            proc.wait(timeout=3)
            print(f"Successfully terminated process: PID {proc.pid}, Name: {proc_display_name}")
            killed_count += 1
        except psutil.NoSuchProcess:
            print(f"Process PID {proc.pid} ({proc_display_name}) already terminated.")
            killed_count += 1
        except psutil.AccessDenied:
            print(f"Access denied: Cannot terminate process PID {proc.pid}, Name: {proc_display_name}. (Try running as administrator/root)")
            failed_count += 1
        except psutil.TimeoutExpired:
            print(f"Process PID {proc.pid}, Name: {proc_display_name} did not terminate within timeout. Trying kill...")
            try:
                proc.kill()
                proc.wait(timeout=3)
                print(f"Successfully force-killed process: PID {proc.pid}, Name: {proc_display_name}")
                killed_count += 1
            except (psutil.NoSuchProcess, psutil.AccessDenied):
                print(f"Failed to force-kill process PID {proc.pid}, Name: {proc_display_name}.")
                failed_count += 1
            except Exception as e:
                print(f"An unexpected error occurred while force-killing PID {proc.pid}: {e}")
                failed_count += 1
        except Exception as e:
            print(f"An unexpected error occurred while terminating PID {proc.pid}: {e}")
            failed_count += 1

    print(f"\n--- Summary ---")
    print(f"Total processes found: {len(matching_processes)}")
    print(f"Processes successfully terminated: {killed_count}")
    print(f"Processes failed to terminate: {failed_count}")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python script_name.py <process_name> [--no-confirm] [--list-only] [--case-sensitive]")
        sys.exit(1)

    process_name_arg = sys.argv[1]
    no_confirm_arg = "--no-confirm" in sys.argv
    list_only_arg = "--list-only" in sys.argv
    case_sensitive_arg = "--case-sensitive" in sys.argv

    kill_processes_by_name(
        target_name=process_name_arg,
        confirm=not no_confirm_arg,
        list_only=list_only_arg,
        case_sensitive=case_sensitive_arg
    )