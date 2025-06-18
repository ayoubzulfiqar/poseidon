import psutil
import time
import datetime
import os

def get_system_resource_usage():
    cpu_percent = psutil.cpu_percent()
    virtual_memory = psutil.virtual_memory()
    disk_usage = psutil.disk_usage('/')
    net_io = psutil.net_io_counters()

    return {
        "timestamp": datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
        "cpu_percent": cpu_percent,
        "memory_total_gb": round(virtual_memory.total / (1024**3), 2),
        "memory_used_gb": round(virtual_memory.used / (1024**3), 2),
        "memory_percent": virtual_memory.percent,
        "disk_total_gb": round(disk_usage.total / (1024**3), 2),
        "disk_used_gb": round(disk_usage.used / (1024**3), 2),
        "disk_percent": disk_usage.percent,
        "net_bytes_sent_mb": round(net_io.bytes_sent / (1024**2), 2),
        "net_bytes_recv_mb": round(net_io.bytes_recv / (1024**2), 2)
    }

def log_resource_usage(log_file_path, interval_seconds=5, duration_seconds=60):
    print(f"Starting resource usage logger. Logging to '{log_file_path}' every {interval_seconds} seconds.")
    if duration_seconds > 0:
        print(f"Logging will run for approximately {duration_seconds} seconds.")
    else:
        print("Logging will run continuously until manually stopped (Ctrl+C).")

    start_time = time.time()
    
    if not os.path.exists(log_file_path) or os.path.getsize(log_file_path) == 0:
        header = "Timestamp,CPU_Percent,Mem_Total_GB,Mem_Used_GB,Mem_Percent,Disk_Total_GB,Disk_Used_GB,Disk_Percent,Net_Sent_MB,Net_Recv_MB\n"
        with open(log_file_path, 'w') as f:
            f.write(header)

    try:
        psutil.cpu_percent() 

        while True:
            current_time = time.time()
            if duration_seconds > 0 and (current_time - start_time) > duration_seconds:
                print("Logging duration completed.")
                break

            usage_data = get_system_resource_usage()
            
            log_entry = (
                f"{usage_data['timestamp']},"
                f"{usage_data['cpu_percent']},"
                f"{usage_data['memory_total_gb']},"
                f"{usage_data['memory_used_gb']},"
                f"{usage_data['memory_percent']},"
                f"{usage_data['disk_total_gb']},"
                f"{usage_data['disk_used_gb']},"
                f"{usage_data['disk_percent']},"
                f"{usage_data['net_bytes_sent_mb']},"
                f"{usage_data['net_bytes_recv_mb']}\n"
            )

            with open(log_file_path, 'a') as f:
                f.write(log_entry)
            
            time.sleep(interval_seconds)

    except KeyboardInterrupt:
        print("\nResource usage logger stopped by user.")
    except Exception as e:
        print(f"\nAn error occurred: {e}")

if __name__ == "__main__":
    LOG_FILE = "system_resource_log.csv"
    LOG_INTERVAL_SECONDS = 5
    LOG_DURATION_SECONDS = 60 

    log_resource_usage(LOG_FILE, LOG_INTERVAL_SECONDS, LOG_DURATION_SECONDS)

# Additional implementation at 2025-06-18 00:11:17
import psutil
import time
import datetime
import os
import signal

# Configuration
LOG_FILE = "system_resource_usage.log"
LOG_INTERVAL_SECONDS = 5
# On Windows, use a path like 'C:\\'
MONITORED_DISK_PATH = "/"

# Global flag for graceful shutdown
running = True

def signal_handler(signum, frame):
    """Handles SIGINT (Ctrl+C) and SIGTERM for graceful shutdown."""
    global running
    print("\nCtrl+C detected. Shutting down gracefully...")
    running = False

def get_resource_usage():
    """
    Collects system-wide CPU, Memory, Disk, and Network usage.
    Returns a dictionary of metrics.
    """
    timestamp = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    # CPU Usage
    # psutil.cpu_percent(interval=None) calculates the percentage since the last call
    # or since boot if it's the first call.
    cpu_percent = psutil.cpu_percent(interval=None)

    # Memory Usage
    virtual_memory = psutil.virtual_memory()
    mem_total_gb = virtual_memory.total / (1024**3)
    mem_used_gb = virtual_memory.used / (1024**3)
    mem_percent = virtual_memory.percent

    # Disk Usage for the specified path
    disk_total_gb = disk_used_gb = disk_percent = "N/A"
    try:
        disk_usage = psutil.disk_usage(MONITORED_DISK_PATH)
        disk_total_gb = disk_usage.total / (1024**3)
        disk_used_gb = disk_usage.used / (1024**3)
        disk_percent = disk_usage.percent
    except Exception as e:
        print(f"Warning: Could not get disk usage for '{MONITORED_DISK_PATH}': {e}")

    # Network I/O (cumulative since boot)
    net_io = psutil.net_io_counters()
    bytes_sent_mb = net_io.bytes_sent / (1024**2)
    bytes_recv_mb = net_io.bytes_recv / (1024**2)

    return {
        "timestamp": timestamp,
        "cpu_percent": cpu_percent,
        "mem_total_gb": mem_total_gb,
        "mem_used_gb": mem_used_gb,
        "mem_percent": mem_percent,
        "disk_path": MONITORED_DISK_PATH,
        "disk_total_gb": disk_total_gb,
        "disk_used_gb": disk_used_gb,
        "disk_percent": disk_percent,
        "net_bytes_sent_mb": bytes_sent_mb,
        "net_bytes_recv_mb": bytes_recv_mb
    }

def log_resource_usage(data):
    """
    Formats and logs the collected resource usage data to a file and prints to console.
    """
    # Format disk usage values if they are numbers
    formatted_disk_total = f"{data['disk_total_gb']:.2f}" if isinstance(data['disk_total_gb'], float) else data['disk_total_gb']
    formatted_disk_used = f"{data['disk_used_gb']:.2f}" if isinstance(data['disk_used_gb'], float) else data['disk_used_gb']
    formatted_disk_percent = f"{data['disk_percent']:.2f}" if isinstance(data['disk_percent'], float) else data['disk_percent']

    log_entry = (
        f"Timestamp: {data['timestamp']} | "
        f"CPU: {data['cpu_percent']:.2f}% | "
        f"Memory: {data['mem_used_gb']:.2f}GB/{data['mem_total_gb']:.2f}GB ({data['mem_percent']:.2f}%) | "
        f"Disk ({data['disk_path']}): {formatted_disk_used}GB/{formatted_disk_total}GB ({formatted_disk_percent}%) | "
        f"Net I/O: Sent {data['net_bytes_sent_mb']:.2f}MB, Recv {data['net_bytes_recv_mb']:.2f}MB"
    )
    try:
        with open(LOG_FILE, "a") as f:
            f.write(log_entry + "\n")
        print(f"Logged: {log_entry}")
    except IOError as e:
        print(f"Error writing to log file '{LOG_FILE}': {e}")

def main():
    """Main function to set up and run the resource logger."""
    # Register signal handlers for graceful shutdown
    signal.signal(signal.SIGINT, signal_handler)  # Handles Ctrl+C
    signal.signal(signal.SIGTERM, signal_handler) # Handles kill command

    print(f"Starting system resource usage logger.")
    print(f"Logs will be saved to '{LOG_FILE}' every {LOG_INTERVAL_SECONDS} seconds.")
    print("Press Ctrl+C to stop.")

    # Write a header to the log file if it's new or empty
    if not os.path.exists(LOG_FILE) or os.path.getsize(LOG_FILE) == 0:
        try:
            with open(LOG_FILE, "w") as f:
                f.write("--- System Resource Usage Log ---\n")
                f.write(f"Log started at: {datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
                f.write(f"Logging interval: {LOG_INTERVAL_SECONDS} seconds\n")
                f.write(f"Monitored disk path: {MONITORED_DISK_PATH}\n")
                f.write("---------------------------------\n")
        except IOError as e:
            print(f"Error creating/writing header to log file '{LOG_FILE}': {e}")
            return # Exit if we can't even write the header

    # Main logging loop
    while running:
        try:
            usage_data = get_resource_usage()
            log_resource_usage(usage_data)
        except Exception as e:
            print(f"An unexpected error occurred during resource collection or logging: {e}")
        time.sleep(LOG_INTERVAL_SECONDS)

    print("Logger stopped.")

if __name__ == "__main__":
    main()