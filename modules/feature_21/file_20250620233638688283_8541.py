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
            return '%.2f %s/s' % (value, s)
    return '%.2f B/s' % n

def monitor_bandwidth():
    sample_interval = 1.0 # seconds between samples
    
    old_net_io = psutil.net_io_counters(pernic=True)
    last_time = time.time()

    try:
        while True:
            time.sleep(sample_interval)
            new_net_io = psutil.net_io_counters(pernic=True)
            current_time = time.time()
            time_elapsed = current_time - last_time

            os.system('cls' if os.name == 'nt' else 'clear') 
            print(f"Bandwidth Monitor (Interval: {sample_interval:.1f}s)\n")

            for interface, new_stats in new_net_io.items():
                if interface in old_net_io:
                    old_stats = old_net_io[interface]

                    bytes_sent_diff = new_stats.bytes_sent - old_stats.bytes_sent
                    bytes_recv_diff = new_stats.bytes_recv - old_stats.bytes_recv

                    sent_speed = bytes_sent_diff / time_elapsed
                    recv_speed = bytes_recv_diff / time_elapsed

                    print(f"Interface: {interface}")
                    print(f"  Sent: {bytes_to_human(sent_speed)}")
                    print(f"  Recv: {bytes_to_human(recv_speed)}")
                    print("-" * 30)

            old_net_io = new_net_io
            last_time = current_time

    except KeyboardInterrupt:
        print("\nBandwidth monitor stopped.")
    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    monitor_bandwidth()

# Additional implementation at 2025-06-20 23:37:04
import psutil
import time
import os
import sys

def bytes_to_human_readable(n_bytes):
    """Convert bytes to human-readable format (B, KB, MB, GB)."""
    symbols = ('B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB')
    prefix = {}
    for i, s in enumerate(symbols):
        prefix[s] = 1 << (i * 10)
    for s in reversed(symbols):
        if n_bytes >= prefix[s]:
            value = float(n_bytes) / prefix[s]
            return f"{value:.2f} {s}"
    return f"{n_bytes:.2f} B"

def clear_screen():
    """Clears the console screen."""
    os.system('cls' if os.name == 'nt' else 'clear')

def get_interface_stats(interface_name=None):
    """
    Get network interface statistics.
    Returns a dictionary of {interface_name: {bytes_sent: int, bytes_recv: int}}
    or None if the interface is not found.
    """
    net_io_counters = psutil.net_io_counters(pernic=True)
    if interface_name:
        if interface_name in net_io_counters:
            return {interface_name: net_io_counters[interface_name]}
        else:
            return None # Interface not found
    else:
        return net_io_counters # Return all interfaces

def monitor_bandwidth(interface_name=None, interval=1.0):
    """
    Monitors bandwidth usage for a specified network interface or all interfaces.
    Displays current speed, total data transferred, and peak speeds.
    """
    initial_stats = get_interface_stats(interface_name)
    if initial_stats is None:
        print(f"Error: Interface '{interface_name}' not found.")
        return

    # Store initial total bytes for session total calculation
    session_start_bytes_sent = {iface: stats.bytes_sent for iface, stats in initial_stats.items()}
    session_start_bytes_recv = {iface: stats.bytes_recv for iface, stats in initial_stats.items()}

    # Store previous stats for speed calculation
    previous_stats = initial_stats

    # Store peak speeds
    peak_upload_speed = {iface: 0 for iface in initial_stats.keys()}
    peak_download_speed = {iface: 0 for iface in initial_stats.keys()}

    print("Starting bandwidth monitor. Press Ctrl+C to stop.")
    time.sleep(interval) # Wait for the first interval before calculating

    try:
        while True:
            current_stats = get_interface_stats(interface_name)
            if current_stats is None:
                print(f"Error: Interface '{interface_name}' disappeared. Exiting.")
                break

            clear_screen()
            print("--- Network Bandwidth Monitor ---")
            print(f"Monitoring Interval: {interval:.1f} seconds")
            print("-" * 35)

            for iface, current_io in current_stats.items():
                if iface not in previous_stats:
                    # New interface appeared, initialize its stats
                    previous_stats[iface] = current_io
                    session_start_bytes_sent[iface] = current_io.bytes_sent
                    session_start_bytes_recv[iface] = current_io.bytes_recv
                    peak_upload_speed[iface] = 0
                    peak_download_speed[iface] = 0
                    continue # Skip calculation for this new interface in this cycle

                prev_io = previous_stats[iface]

                # Calculate current speeds
                bytes_sent_delta = current_io.bytes_sent - prev_io.bytes_sent
                bytes_recv_delta = current_io.bytes_recv - prev_io.bytes_recv

                upload_speed = bytes_sent_delta / interval
                download_speed = bytes_recv_delta / interval

                # Update peak speeds
                peak_upload_speed[iface] = max(peak_upload_speed[iface], upload_speed)
                peak_download_speed[iface] = max(peak_download_speed[iface], download_speed)

                # Calculate total data transferred during this session
                total_sent_session = current_io.bytes_sent - session_start_bytes_sent.get(iface, 0)
                total_recv_session = current_io.bytes_recv - session_start_bytes_recv.get(iface, 0)

                print(f"Interface: {iface}")
                print(f"  Upload:   {bytes_to_human_readable(upload_speed)}/s")
                print(f"  Download: {bytes_to_human_readable(download_speed)}/s")
                print(f"  Total Sent (Session): {bytes_to_human_readable(total_sent_session)}")
                print(f"  Total Recv (Session): {bytes_to_human_readable(total_recv_session)}")
                print(f"  Peak Up (Session):    {bytes_to_human_readable(peak_upload_speed[iface])}/s")
                print(f"  Peak Down (Session):  {bytes_to_human_readable(peak_download_speed[iface])}/s")
                print("-" * 35)

                previous_stats[iface] = current_io # Update previous stats for next iteration

            time.sleep(interval)

    except KeyboardInterrupt:
        print("\nBandwidth monitoring stopped.")
    except Exception as e:
        print(f"\nAn error occurred: {e}")

if __name__ == "__main__":
    # Get available network interfaces
    interfaces = psutil.net_io_counters(pernic=True)
    if not interfaces:
        print("No network interfaces found.")
        sys.exit(1)

    interface_names = sorted(interfaces.keys())

    print("Available Network Interfaces:")
    for i, name in enumerate(interface_names):
        print(f"  {i+1}. {name}")
    print("  0. Monitor ALL interfaces")

    selected_interface = None
    while True:
        try:
            choice = input("Enter the number of the interface to monitor (or 0 for all): ")
            choice_int = int(choice)
            if choice_int == 0:
                selected_interface = None # Monitor all
                break
            elif 1 <= choice_int <= len(interface_names):
                selected_interface = interface_names[choice_int - 1]
                break
            else:
                print("Invalid choice. Please enter a number from the list.")
        except ValueError:
            print("Invalid input. Please enter a number.")

    update_interval = 1.0
    while True:
        try:
            interval_str = input(f"Enter update interval in seconds (default: {update_interval:.1f}): ")
            if not interval_str:
                break # Use default
            new_interval = float(interval_str)
            if new_interval > 0:
                update_interval = new_interval
                break
            else:
                print("Interval must be a positive number.")
        except ValueError:
            print("Invalid input. Please enter a number.")

    monitor_bandwidth(selected_interface, update_interval)

# Additional implementation at 2025-06-20 23:38:06
import psutil
import time
import os

class BandwidthMonitor:
    def __init__(self):
        self.last_io_counters = {}
        self.last_time = time.time()
        self.total_sent_session = 0
        self.total_recv_session = 0

    def _bytes_to_human_readable(self, bytes_val):
        if bytes_val is None:
            return "N/A"
        units = ['B', 'KB', 'MB', 'GB', 'TB']
        i = 0
        while bytes_val >= 1024 and i < len(units) - 1:
            bytes_val /= 1024.0
            i += 1
        return f"{bytes_val:.2f} {units[i]}"

    def get_available_interfaces(self):
        return list(psutil.net_io_counters(pernic=True).keys())

    def get_bandwidth_and_totals(self, interface_name):
        current_io = psutil.net_io_counters(pernic=True)
        current_time = time.time()

        if interface_name not in current_io:
            return None, None, None, None

        if interface_name not in self.last_io_counters:
            self.last_io_counters[interface_name] = current_io[interface_name]
            self.last_time = current_time
            return 0, 0, 0, 0

        last_io = self.last_io_counters[interface_name]
        current_io_data = current_io[interface_name]

        delta_time = current_time - self.last_time
        if delta_time == 0:
            return 0, 0, self.total_sent_session, self.total_recv_session

        delta_bytes_sent = current_io_data.bytes_sent - last_io.bytes_sent
        delta_bytes_recv = current_io_data.bytes_recv - last_io.bytes_recv

        self.total_sent_session += delta_bytes_sent
        self.total_recv_session += delta_bytes_recv

        rate_sent = delta_bytes_sent / delta_time
        rate_recv = delta_bytes_recv / delta_time

        self.last_io_counters[interface_name] = current_io_data
        self.last_time = current_time

        return rate_sent, rate_recv, self.total_sent_session, self.total_recv_session

    def run(self, interval=1):
        interfaces = self.get_available_interfaces()

        if not interfaces:
            print("No network interfaces found.")
            return

        print("Available network interfaces:")
        for i, iface in enumerate(interfaces):
            print(f"  {i+1}. {iface}")
        print("  0. Monitor all interfaces")

        while True:
            try:
                choice = input("Enter the number of the interface to monitor (or 0 for all): ")
                choice_idx = int(choice)
                if choice_idx == 0:
                    selected_interfaces = interfaces
                    break
                elif 1 <= choice_idx <= len(interfaces):
                    selected_interfaces = [interfaces[choice_idx - 1]]
                    break
                else:
                    print("Invalid choice. Please try again.")
            except ValueError:
                print("Invalid input. Please enter a number.")

        try:
            while True:
                os.system('cls' if os.name == 'nt' else 'clear')
                print(f"--- Network Bandwidth Monitor (Interval: {interval}s) ---")
                print(f"Monitoring: {', '.join(selected_interfaces)}")
                print("-" * 50)

                for iface in selected_interfaces:
                    rate_sent, rate_recv, total_sent, total_recv = self.get_bandwidth_and_totals(iface)

                    if rate_sent is None:
                        print(f"Interface '{iface}' not found or disconnected.")
                        continue

                    print(f"Interface: {iface}")
                    print(f"  Sent: {self._bytes_to_human_readable(rate_sent)}/s")
                    print(f"  Recv: {self._bytes_to_human_readable(rate_recv)}/s")
                    print(f"  Total Sent (Session): {self._bytes_to_human_readable(total_sent)}")
                    print(f"  Total Recv (Session): {self._bytes_to_human_readable(total_recv)}")
                    print("-" * 50)

                time.sleep(interval)

        except KeyboardInterrupt:
            print("\nBandwidth monitor stopped.")

if __name__ == "__main__":
    monitor = BandwidthMonitor()
    monitor.run(interval=1)