import psutil
import time
import os

def bytes_to_human(n):
    symbols = ('B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB')
    prefix = {}
    for i, s in enumerate(symbols):
        prefix[s] = 1 << (i * 10)
    for s in reversed(symbols):
        if n >= prefix[s]:
            value = float(n) / prefix[s]
            return f"{value:.2f} {s}"
    return f"{n:.2f} B"

def monitor_bandwidth(interval=1):
    old_net_io = psutil.net_io_counters(pernic=True)
    
    try:
        while True:
            time.sleep(interval)
            new_net_io = psutil.net_io_counters(pernic=True)
            
            os.system('cls' if os.name == 'nt' else 'clear')
            print(f"{'Interface':<15} {'Download':<15} {'Upload':<15}")
            print("-" * 45)
            
            for interface, new_stats in new_net_io.items():
                if interface in old_net_io:
                    old_stats = old_net_io[interface]
                    
                    bytes_recv_diff = new_stats.bytes_recv - old_stats.bytes_recv
                    bytes_sent_diff = new_stats.bytes_sent - old_stats.bytes_sent
                    
                    download_speed = bytes_recv_diff / interval
                    upload_speed = bytes_sent_diff / interval
                    
                    print(f"{interface:<15} {bytes_to_human(download_speed) + '/s':<15} {bytes_to_human(upload_speed) + '/s':<15}")
                else:
                    print(f"{interface:<15} {'N/A':<15} {'N/A':<15}")
            
            old_net_io = new_net_io
            
    except KeyboardInterrupt:
        print("\nMonitoring stopped.")

if __name__ == "__main__":
    monitor_bandwidth()

# Additional implementation at 2025-06-18 00:18:00
import psutil
import time
import os

def bytes_to_human_readable(n_bytes):
    if n_bytes < 1024:
        return f"{n_bytes:.2f} B"
    elif n_bytes < 1024 * 1024:
        return f"{n_bytes / 1024:.2f} KB"
    elif n_bytes < 1024 * 1024 * 1024:
        return f"{n_bytes / (1024 * 1024)::.2f} MB"
    else:
        return f"{n_bytes / (1024 * 1024 * 1024):.2f} GB"

def bytes_per_second_to_human_readable(n_bytes_per_sec):
    if n_bytes_per_sec < 1024:
        return f"{n_bytes_per_sec:.2f} B/s"
    elif n_bytes_per_sec < 1024 * 1024:
        return f"{n_bytes_per_sec / 1024:.2f} KB/s"
    elif n_bytes_per_sec < 1024 * 1024 * 1024:
        return f"{n_bytes_per_sec / (1024 * 1024):.2f} MB/s"
    else:
        return f"{n_bytes_per_sec / (1024 * 1024 * 1024):.2f} GB/s"

def clear_screen():
    os.system('cls' if os.name == 'nt' else 'clear')

def get_network_interfaces():
    return psutil.net_io_counters(pernic=True)

def main():
    clear_screen()
    print("Available Network Interfaces:")
    interfaces = get_network_interfaces()
    
    if not interfaces:
        print("No network interfaces found or psutil could not retrieve them.")
        return

    interface_names = list(interfaces.keys())
    for i, name in enumerate(interface_names):
        print(f"{i + 1}. {name}")

    selected_interface = None
    while selected_interface is None:
        try:
            choice = int(input("Enter the number of the interface to monitor: "))
            if 1 <= choice <= len(interface_names):
                selected_interface = interface_names[choice - 1]
            else:
                print("Invalid choice. Please enter a number within the range.")
        except ValueError:
            print("Invalid input. Please enter a number.")

    print(f"\nMonitoring interface: {selected_interface}")
    print("Press Ctrl+C to stop.")

    initial_stats = psutil.net_io_counters(pernic=True)[selected_interface]
    last_bytes_sent = initial_stats.bytes_sent
    last_bytes_recv = initial_stats.bytes_recv
    last_time = time.time()

    total_bytes_sent_session = 0
    total_bytes_recv_session = 0
    peak_upload_speed = 0
    peak_download_speed = 0

    clear_screen()

    try:
        while True:
            current_stats = psutil.net_io_counters(pernic=True)[selected_interface]
            current_time = time.time()

            time_diff = current_time - last_time
            if time_diff == 0:
                time.sleep(0.1)
                continue

            bytes_sent_diff = current_stats.bytes_sent - last_bytes_sent
            bytes_recv_diff = current_stats.bytes_recv - last_bytes_recv

            current_upload_speed = bytes_sent_diff / time_diff
            current_download_speed = bytes_recv_diff / time_diff

            total_bytes_sent_session += bytes_sent_diff
            total_bytes_recv_session += bytes_recv_diff

            peak_upload_speed = max(peak_upload_speed, current_upload_speed)
            peak_download_speed = max(peak_download_speed, current_download_speed)

            last_bytes_sent = current_stats.bytes_sent
            last_bytes_recv = current_stats.bytes_recv
            last_time = current_time

            clear_screen()
            print(f"Monitoring Interface: {selected_interface}")
            print("-" * 40)
            print(f"Current Upload:   {bytes_per_second_to_human_readable(current_upload_speed)}")
            print(f"Current Download: {bytes_per_second_to_human_readable(current_download_speed)}")
            print("-" * 40)
            print(f"Total Upload (Session):   {bytes_to_human_readable(total_bytes_sent_session)}")
            print(f"Total Download (Session): {bytes_to_human_readable(total_bytes_recv_session)}")
            print("-" * 40)
            print(f"Peak Upload Speed:   {bytes_per_second_to_human_readable(peak_upload_speed)}")
            print(f"Peak Download Speed: {bytes_per_second_to_human_readable(peak_download_speed)}")
            print("-" * 40)
            print("Press Ctrl+C to stop.")

            time.sleep(1)

    except KeyboardInterrupt:
        print("\nMonitoring stopped.")
    except Exception as e:
        print(f"\nAn error occurred: {e}")

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-18 00:19:04
import psutil
import time
import collections
import sys

def bytes_to_human_readable(bytes_val):
    """Converts a byte value into a human-readable string (e.g., KB, MB, GB)."""
    for unit in ['B', 'KB', 'MB', 'GB', 'TB']:
        if bytes_val < 1024.0:
            return f"{bytes_val:.2f} {unit}"
        bytes_val /= 1024.0
    return f"{bytes_val:.2f} PB"

def monitor_bandwidth(interface_name=None, interval=1.0, threshold_mbps=0):
    """
    Monitors network bandwidth for a specified interface or all interfaces,
    displaying current speeds, total data transferred, rolling averages,
    and triggering alerts if a threshold is exceeded.

    Args:
        interface_name (str, optional): The name of the network interface to monitor.
                                        If None, monitors all active interfaces.
        interval (float): The refresh interval in seconds.
        threshold_mbps (float): Bandwidth threshold in Mbps. An alert is printed
                                if download or upload speed exceeds this value.
                                Set to 0 for no threshold alerts.
    """

    print("Initializing bandwidth monitor...")
    print(f"Refresh interval: {interval} seconds")
    if threshold_mbps > 0:
        print(f"Bandwidth alert threshold: {threshold_mbps:.2f} Mbps")
    else:
        print("No bandwidth threshold alerts configured.")

    # Get initial network stats for all interfaces
    initial_net_io = psutil.net_io_counters(pernic=True)
    
    # Store previous stats for calculating deltas
    previous_net_io = {iface: {'bytes_sent': stats.bytes_sent, 'bytes_recv': stats.bytes_recv}
                       for iface, stats in initial_net_io.items()}

    # Store total data transferred since the monitor started
    total_sent = collections.defaultdict(int)
    total_recv = collections.defaultdict(int)

    # Store a history of speeds for rolling average calculation
    history_length = 10 # Number of past samples to consider for average
    download_history = collections.defaultdict(lambda: collections.deque(maxlen=history_length))
    upload_history = collections.defaultdict(lambda: collections.deque(maxlen=history_length))

    try:
        while True:
            time.sleep(interval)
            current_net_io = psutil.net_io_counters(pernic=True)
            
            print("\n" + "="*50)
            print(f"Network Bandwidth Monitor - {time.strftime('%Y-%m-%d %H:%M:%S')}")
            print("="*50)

            # Determine which interfaces to monitor
            monitored_interfaces = [interface_name] if interface_name else sorted(current_net_io.keys())

            for iface in monitored_interfaces:
                # Skip if interface is not currently active or was not present initially
                if iface not in current_net_io:
                    print(f"Interface '{iface}' not found or inactive. Skipping.")
                    continue
                
                # If a new interface appeared, initialize its previous stats
                if iface not in previous_net_io:
                    previous_net_io[iface] = {'bytes_sent': current_net_io[iface].bytes_sent,
                                              'bytes_recv': current_net_io[iface].bytes_recv}
                    # Skip calculation for this first cycle as no delta can be computed
                    continue 

                current_sent = current_net_io[iface].bytes_sent
                current_recv = current_net_io[iface].bytes_recv

                prev_sent = previous_net_io[iface]['bytes_sent']
                prev_recv = previous_net_io[iface]['bytes_recv']

                # Calculate bytes transferred during the interval
                delta_sent = current_sent - prev_sent
                delta_recv = current_recv - prev_recv

                # Calculate speeds in bytes per second
                upload_speed_bps = delta_sent / interval
                download_speed_bps = delta_recv / interval

                # Update total data transferred since start
                total_sent[iface] += delta_sent
                total_recv[iface] += delta_recv

                # Add current speeds to history for rolling average
                download_history[iface].append(download_speed_bps)
                upload_history[iface].append(upload_speed_bps)

                # Calculate rolling averages
                avg_download_bps = sum(download_history[iface]) / len(download_history[iface]) if download_history[iface] else 0
                avg_upload_bps = sum(upload_history[iface]) / len(upload_history[iface]) if upload_history[iface] else 0

                print(f"Interface: {iface}")
                print(f"  Download Speed: {bytes_to_human_readable(download_speed_bps)}/s ({download_speed_bps * 8 / 1_000_000:.2f} Mbps)")
                print(f"  Upload Speed:   {bytes_to_human_readable(upload_speed_bps)}/s ({upload_speed_bps * 8 / 1_000_000:.2f} Mbps)")
                print(f"  Avg. Download:  {bytes_to_human_readable(avg_download_bps)}/s (over last {len(download_history[iface])} samples)")
                print(f"  Avg. Upload:    {bytes_to_human_readable(avg_upload_bps)}/s (over last {len(upload_history[iface])} samples)")
                print(f"  Total Download: {bytes_to_human_readable(total_recv[iface])}")
                print(f"  Total Upload:   {bytes_to_human_readable(total_sent[iface])}")

                # Check for bandwidth threshold alerts
                if threshold_mbps > 0:
                    download_mbps = download_speed_bps * 8 / 1_000_000
                    upload_mbps = upload_speed_bps * 8 / 1_000_000
                    if download_mbps >= threshold_mbps:
                        print(f"  !!! ALERT: High Download Usage ({download_mbps:.2f} Mbps) on {iface} !!!")
                    if upload_mbps >= threshold_mbps:
                        print(f"  !!! ALERT: High Upload Usage ({upload_mbps:.2f} Mbps) on {iface} !!!")
                print("-" * 40)

                # Update previous stats for the next monitoring cycle
                previous_net_io[iface]['bytes_sent'] = current_sent
                previous_net_io[iface]['bytes_recv'] = current_recv

    except KeyboardInterrupt:
        print("\nMonitoring stopped by user.")
    except Exception as e:
        print(f"\nAn unexpected error occurred: {e}")

def list_interfaces():
    """Prints a list of available network interfaces."""
    print("Available Network Interfaces:")
    interfaces = psutil.net_io_counters(pernic=True)
    if not interfaces:
        print("No network interfaces found.")
        return
    for iface in sorted(interfaces.keys()):
        print(f"- {iface}")
    print("\n")

if __name__ == "__main__":
    # Default configuration values
    default_interval = 1.0  # seconds
    default_threshold_mbps = 10.0 # Mbps (0 for no alerts)

    list_interfaces()

    # User input for interface selection
    selected_interface = input("Enter interface name to monitor (leave blank to monitor all): ").strip()
    if not selected_interface:
        selected_interface = None
        print("Monitoring all active interfaces.")
    else:
        # Validate if the entered interface exists
        if selected_interface not in psutil.net_io_counters(pernic=True):
            print(f"Error: Interface '{selected_interface}' not found. Monitoring all active interfaces instead.")
            selected_interface = None

    # User input for refresh interval
    try:
        user_interval_str = input(f"Enter refresh interval in seconds (default: {default_interval}): ").strip()
        user_interval = float(user_interval_str) if user_interval_str else default_interval
        if user_interval <= 0:
            raise ValueError("Interval must be a positive number.")
    except ValueError:
        print(f"Invalid interval. Using default: {default_interval} seconds.")
        user_interval = default_interval

    # User input for bandwidth threshold
    try:
        user_threshold_str = input(f"Enter bandwidth alert threshold in Mbps (0 for no alert, default: {default_threshold_mbps}): ").strip()
        user_threshold = float(user_threshold_str) if user_threshold_str else default_threshold_mbps
        if user_threshold < 0:
            raise ValueError("Threshold cannot be negative.")
    except ValueError:
        print(f"Invalid threshold. Using default: {default_threshold_mbps} Mbps.")
        user_threshold = default_threshold_mbps

    # Start the bandwidth monitoring
    monitor_bandwidth(selected_interface, user_interval, user_threshold)

# Additional implementation at 2025-06-18 00:19:57
import psutil
import time
import os
import sys

class BandwidthMonitor:
    def __init__(self, interfaces=None, update_interval=1.0, history_size=10):
        self.interfaces = interfaces
        self.update_interval = max(0.1, float(update_interval))
        self.history_size = max(1, int(history_size))
        self.last_counters = {}
        self.interface_data = {}
        self._initialize_interfaces()

    def _initialize_interfaces(self):
        all_interfaces = psutil.net_io_counters(pernic=True)
        if not self.interfaces:
            self.interfaces = sorted(all_interfaces.keys())
        else:
            valid_interfaces = []
            for iface in self.interfaces:
                if iface in all_interfaces:
                    valid_interfaces.append(iface)
            self.interfaces = valid_interfaces

        if not self.interfaces:
            raise ValueError("No valid network interfaces found or specified to monitor.")

        for iface in self.interfaces:
            self.last_counters[iface] = all_interfaces[iface]
            self.interface_data[iface] = {
                'upload_history': [],
                'download_history': [],
                'current_upload': 0,
                'current_download': 0,
                'peak_upload': 0,
                'peak_download': 0,
                'total_sent': 0,
                'total_received': 0,
                'initial_sent': all_interfaces[iface].bytes_sent,
                'initial_received': all_interfaces[iface].bytes_recv
            }

    def _clear_screen(self):
        os.system('cls' if os.name == 'nt' else 'clear')

    def _bytes_to_human_readable(self, n_bytes, unit_suffix="/s"):
        symbols = ('B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB')
        prefix = {}
        for i, s in enumerate(symbols):
            prefix[s] = 1 << (i * 10)
        for s in reversed(symbols):
            if n_bytes >= prefix[s]:
                value = float(n_bytes) / prefix[s]
                return f"{value:.2f} {s}{unit_suffix}"
        return f"{n_bytes:.2f} B{unit_suffix}"

    def _update_stats(self):
        current_counters = psutil.net_io_counters(pernic=True)
        for iface in self.interfaces:
            if iface not in current_counters:
                continue

            prev_sent = self.last_counters[iface].bytes_sent
            prev_received = self.last_counters[iface].bytes_recv
            
            current_sent = current_counters[iface].bytes_sent
            current_received = current_counters[iface].bytes_recv

            upload_speed = (current_sent - prev_sent) / self.update_interval
            download_speed