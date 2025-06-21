import psutil
import time
import logging
import os
LOG_FILE='system_resource_usage.log'
logging.basicConfig(filename=LOG_FILE,level=logging.INFO,format='%(asctime)s - %(levelname)s - %(message)s')
def get_current_stats():
    cpu_percent=psutil.cpu_percent(interval=None)
    mem_info=psutil.virtual_memory()
    disk_info=psutil.disk_usage('/')
    net_io=psutil.net_io_counters()
    return {'cpu_percent':cpu_percent,'memory_percent':mem_info.percent,'disk_percent':disk_info.percent,'net_bytes_sent':net_io.bytes_sent,'net_bytes_recv':net_io.bytes_recv}
def start_resource_logger(interval_seconds=5):
    previous_net_io=psutil.net_io_counters()
    last_log_time=time.time()
    try:
        while True:
            current_time=time.time()
            elapsed_time=current_time-last_log_time
            current_stats=get_current_stats()
            current_net_io=psutil.net_io_counters()
            bytes_sent_delta=current_net_io.bytes_sent-previous_net_io.bytes_sent
            bytes_recv_delta=current_net_io.bytes_recv-previous_net_io.bytes_recv
            net_sent_kbps=(bytes_sent_delta/1024)/elapsed_time if elapsed_time>0 else 0
            net_recv_kbps=(bytes_recv_delta/1024)/elapsed_time if elapsed_time>0 else 0
            log_message=(f"CPU: {current_stats['cpu_percent']:.2f}% | Memory: {current_stats['memory_percent']:.2f}% | Disk: {current_stats['disk_percent']:.2f}% | Net Sent: {net_sent_kbps:.2f} KB/s | Net Recv: {net_recv_kbps:.2f} KB/s")
            logging.info(log_message)
            previous_net_io=current_net_io
            last_log_time=current_time
            time.sleep(interval_seconds)
    except KeyboardInterrupt:
        logging.info("System resource logger stopped by user.")
    except Exception as e:
        logging.error(f"An error occurred: {e}")
if __name__=="__main__":
    start_resource_logger(interval_seconds=5)

# Additional implementation at 2025-06-20 22:33:17
import psutil
import time
import csv
import os
import sys

class SystemResourceLogger:
    def __init__(self, log_file="system_resource_log.csv", interval_seconds=5, monitor_pid=None):
        if not self._check_psutil():
            print("Error: psutil library not found. Please install it using 'pip install psutil'.")
            sys.exit(1)

        self.log_file = log_file
        self.interval = interval_seconds
        self.monitor_pid = monitor_pid
        self._running = False
        self._csv_header = [
            "Timestamp",
            "CPU_Percent",
            "Memory_Total_GB",
            "Memory_Used_GB",
            "Memory_Percent",
            "Disk_Total_GB",
            "Disk_Used_GB",
            "Disk_Percent",
            "Net_Sent_MB",
            "Net_Recv_MB"
        ]
        if self.monitor_pid:
            self._csv_header.extend([
                "Process_PID",
                "Process_Name",
                "Process_CPU_Percent",
                "Process_Memory_RSS_MB",
                "Process_Memory_Percent",
                "Process_Threads"
            ])
        self._initialize_log_file()
        # Store initial network counters for delta calculation
        self._initial_net_io = psutil.net_io_counters()

    def _check_psutil(self):
        try:
            import psutil
            return True
        except ImportError:
            return False

    def _initialize_log_file(self):
        file_exists = os.path.exists(self.log_file)
        with open(self.log_file, 'a', newline='') as f:
            writer = csv.writer(f)
            # Write header only if file is new or empty
            if not file_exists or os.stat(self.log_file).st_size == 0:
                writer.writerow(self._csv_header)

    def _get_system_metrics(self):
        # CPU usage since last call to cpu_percent()
        cpu_percent = psutil.cpu_percent(interval=None) 
        
        # Memory usage
        mem = psutil.virtual_memory()
        mem_total_gb = round(mem.total / (1024**3), 2)
        mem_used_gb = round(mem.used / (1024**3), 2)
        mem_percent = mem.percent

        # Disk usage (for the root partition or default)
        disk = psutil.disk_usage('/') 
        disk_total_gb = round(disk.total / (1024**3), 2)
        disk_used_gb = round(disk.used / (1024**3), 2)
        disk_percent = disk.percent

        # Network usage (delta since last measurement)
        current_net_io = psutil.net_io_counters()
        net_sent_bytes_delta = current_net_io.bytes_sent - self._initial_net_io.bytes_sent
        net_recv_bytes_delta = current_net_io.bytes_recv - self._initial_net_io.bytes_recv
        
        # Update initial counters for the next interval's calculation
        self._initial_net_io = current_net_io 

        net_sent_mb = round(net_sent_bytes_delta / (1024**2), 2)
        net_recv_mb = round(net_recv_bytes_delta / (1024**2), 2)

        return {
            "CPU_Percent": cpu_percent,
            "Memory_Total_GB": mem_total_gb,
            "Memory_Used_GB": mem_used_gb,
            "Memory_Percent": mem_percent,
            "Disk_Total_GB": disk_total_gb,
            "Disk_Used_GB": disk_used_gb,
            "Disk_Percent": disk_percent,
            "Net_Sent_MB": net_sent_mb,
            "Net_Recv_MB": net_recv_mb
        }

    def _get_process_metrics(self):
        if not self.monitor_pid:
            return {}
        
        try:
            process = psutil.Process(self.monitor_pid)
            # Use oneshot() for efficiency when retrieving multiple process attributes
            with process.oneshot(): 
                proc_name = process.name()
                proc_cpu_percent = process.cpu_percent(interval=None) # Non-blocking
                proc_mem_info = process.memory_info()
                proc_mem_rss_mb = round(proc_mem_info.rss / (1024**2), 2)
                proc_mem_percent = process.memory_percent()
                proc_threads = process.num_threads()

            return {
                "Process_PID": self.monitor_pid,
                "Process_Name": proc_name,
                "Process_CPU_Percent": proc_cpu_percent,
                "Process_Memory_RSS_MB": proc_mem_rss_mb,
                "Process_Memory_Percent": proc_mem_percent,
                "Process_Threads": proc_threads
            }
        except psutil.NoSuchProcess:
            print(f"Warning: Process with PID {self.monitor_pid} not found or has terminated. Stopping monitoring for this PID.")
            self.monitor_pid = None # Stop trying to monitor this PID
            return { # Return N/A for process metrics if process is gone
                "Process_PID": "N/A",
                "Process_Name": "N/A",
                "Process_CPU_Percent": "N/A",
                "Process_Memory_RSS_MB": "N/A",
                "Process_Memory_Percent": "N/A",
                "Process_Threads": "N/A"
            }
        except Exception as e:
            print(f"Error getting process metrics for PID {self.monitor_pid}: {e}")
            return { # Return Error for process metrics if other issues occur
                "Process_PID": "Error",
                "Process_Name": "Error",
                "Process_CPU_Percent": "Error",
                "Process_Memory_RSS_MB": "Error",
                "Process_Memory_Percent": "Error",
                "Process_Threads": "Error"
            }

    def _log_data(self, data):
        with open(self.log_file, 'a', newline='') as f:
            writer = csv.writer(f)
            # Ensure all header fields are present, fill missing with empty string
            row = [data.get(header, '') for header in self._csv_header]
            writer.writerow(row)

    def start_logging(self):
        self._running = True
        print(f"Starting system resource logging to {self.log_file} every {self.interval} seconds.")
        if self.monitor_pid:
            print(f"Monitoring process with PID: {self.monitor_pid}")
        print("Press Ctrl+C to stop.")

        # Prime psutil.cpu_percent() for accurate first reading
        psutil.cpu_percent(interval=0.1) 

        try:
            while self._running:
                timestamp = time.strftime("%Y-%m-%d %H:%M:%S")
                
                system_metrics = self._get_system_metrics()
                log_entry = {"Timestamp": timestamp, **system_metrics}

                if self.monitor_pid is not None: # Only try if a PID is set
                    process_metrics = self._get_process_metrics()
                    log_entry.update(process_metrics)
                
                self._log_data(log_entry)
                time.sleep(self.interval)
        except KeyboardInterrupt:
            print("\nLogging stopped by user.")
        except Exception as e:
            print(f"An unexpected error occurred: {e}")
        finally:
            self.stop_logging()

    def stop_logging(self):
        self._running = False
        print("Logger gracefully stopped.")

if __name__ == "__main__":
    # Example 1: Log system-wide resources only
    # logger_system_only = SystemResourceLogger(log_file="system_only_log.csv", interval_seconds=2)
    # logger_system_only.start_logging()

    # Example 2: Log system-wide resources and a specific process
    # Get the current script's PID for demonstration purposes
    current_script_pid = os.getpid() 
    print(f"Attempting to monitor current script's PID: {current_script_pid}")
    
    logger_with_process = SystemResourceLogger(
        log_file="system_and_process_log.csv",
        interval_seconds=3,
        monitor_pid=current_script_pid
    )
    logger_with_process.start_logging()

    # To monitor another process, replace current_script_pid with its PID:
    # Example: logger_with_another_process = SystemResourceLogger(
    #     log_file="system_and_external_process_log.csv", 
    #     interval_seconds=5, 
    #     monitor_pid=1234 # Replace 1234 with an actual PID you want to monitor
    # )
    # logger_with_another_process.start_logging()

# Additional implementation at 2025-06-20 22:34:33
import psutil
import time
import csv
import os
import sys

class SystemResourceLogger:
    def __init__(self, log_file="resource_usage.csv", interval_seconds=1, duration_seconds=None):
        self.log_file = log_file
        self.interval = interval_seconds
        self.duration = duration_seconds
        self.running = False
        self.start_time = None
        self.csv_writer = None
        self.csv_file = None

        self._prev_net_io = None
        self._prev_disk_io = None

    def _get_resource_usage(self):
        """Collects current system resource usage."""
        cpu_percent = psutil.cpu_percent(interval=None)
        mem_info = psutil.virtual_memory()
        disk_io = psutil.disk_io_counters()
        net_io = psutil.net_io_counters()

        disk_read_bytes_per_sec = 0.0
        disk_write_bytes_per_sec = 0.0
        net_bytes_sent_per_sec = 0.0
        net_bytes_recv_per_sec = 0.0

        if self._prev_disk_io:
            disk_read_bytes_per_sec = (disk_io.read_bytes - self._prev_disk_io.read_bytes) / self.interval
            disk_write_bytes_per_sec = (disk_io.write_bytes - self._prev_disk_io.write_bytes) / self.interval
        self._prev_disk_io = disk_io

        if self._prev_net_io:
            net_bytes_sent_per_sec = (net_io.bytes_sent - self._prev_net_io.bytes_sent) / self.interval
            net_bytes_recv_per_sec = (net_io.bytes_recv - self._prev_net_io.bytes_recv) / self.interval
        self._prev_net_io = net_io

        return {
            "timestamp": time.time(),
            "cpu_percent": cpu_percent,
            "mem_percent": mem_info.percent,
            "mem_total_mb": round(mem_info.total / (1024 * 1024), 2),
            "mem_used_mb": round(mem_info.used / (1024 * 1024), 2),
            "mem_free_mb": round(mem_info.free / (1024 * 1024), 2),
            "disk_read_bytes_per_sec": round(disk_read_bytes_per_sec, 2),
            "disk_write_bytes_per_sec": round(disk_write_bytes_per_sec, 2),
            "net_bytes_sent_per_sec": round(net_bytes_sent_per_sec, 2),
            "net_bytes_recv_per_sec": round(net_bytes_recv_per_sec, 2),
        }

    def start_logging(self):
        """Starts logging system resource usage to a CSV file."""
        self.running = True
        self.start_time = time.time()

        file_exists = os.path.exists(self.log_file)

        try:
            self.csv_file = open(self.log_file, 'a', newline='')
            fieldnames = [
                "timestamp", "cpu_percent", "mem_percent", "mem_total_mb",
                "mem_used_mb", "mem_free_mb", "disk_read_bytes_per_sec",
                "disk_write_bytes_per_sec", "net_bytes_sent_per_sec",
                "net_bytes_recv_per_sec"
            ]
            self.csv_writer = csv.DictWriter(self.csv_file, fieldnames=fieldnames)

            if not file_exists or os.stat(self.log_file).st_size == 0:
                self.csv_writer.writeheader()
                self.csv_file.flush()

            print(f"Logging system resources to '{self.log_file}' every {self.interval} seconds...")

            self._prev_disk_io = psutil.disk_io_counters()
            self._prev_net_io = psutil.net_io_counters()
            time.sleep(self.interval)

            while self.running:
                if self.duration and (time.time() - self.start_time) > self.duration:
                    print(f"Logging duration of {self.duration} seconds reached. Stopping.")
                    self.stop_logging()
                    break

                data = self._get_resource_usage()
                self.csv_writer.writerow(data)
                self.csv_file.flush()

                time.sleep(self.interval)

        except KeyboardInterrupt:
            print("\nLogging interrupted by user (Ctrl+C). Stopping.")
        except Exception as e:
            print(f"An error occurred during logging: {e}")
        finally:
            self.stop_logging()

    def stop_logging(self):
        """Stops the logging process and closes the file."""
        if self.running:
            self.running = False
            if self.csv_file:
                self.csv_file.close()
                print(f"Logging stopped. Data saved to '{self.log_file}'.")
            self.csv_file = None
            self.csv_writer = None

if __name__ == "__main__":
    try:
        import psutil
    except ImportError:
        print("Error: psutil library not found.")
        print("Please install it using: pip install psutil")
        sys.exit(1)

    LOG_FILE_NAME = "system_metrics.csv"
    LOG_INTERVAL_SECONDS = 2
    LOG_DURATION_SECONDS = 30 # Set to None for indefinite logging until Ctrl+C

    logger = SystemResourceLogger(
        log_file=LOG_FILE_NAME,
        interval_seconds=LOG_INTERVAL_SECONDS,
        duration_seconds=LOG_DURATION_SECONDS
    )

    try:
        logger.start_logging()
    except Exception as e:
        print(f"An unhandled error occurred during logger execution: {e}")
    finally:
        logger.stop_logging()