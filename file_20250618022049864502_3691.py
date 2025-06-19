import socket
import struct
import time

def get_ntp_time(host="pool.ntp.org"):
    NTP_PORT = 123
    NTP_DELTA = 2208988800

    ntp_request = bytearray(48)
    ntp_request[0] = 0x1B

    try:
        with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as client_socket:
            client_socket.settimeout(5)
            client_socket.sendto(ntp_request, (host, NTP_PORT))
            response_data, _ = client_socket.recvfrom(48)

            if len(response_data) != 48:
                raise ValueError("NTP response has incorrect length.")

            ntp_integer_part, ntp_fractional_part = struct.unpack('!II', response_data[40:48])

            unix_timestamp = ntp_integer_part - NTP_DELTA

            return unix_timestamp
    except socket.timeout:
        raise TimeoutError(f"NTP request to {host} timed out.")
    except socket.gaierror:
        raise ConnectionError(f"Could not resolve NTP host: {host}")
    except ConnectionRefusedError:
        raise ConnectionRefusedError(f"Connection refused by NTP server: {host}")
    except Exception as e:
        raise RuntimeError(f"An error occurred: {e}")

if __name__ == '__main__':
    try:
        current_unix_time = get_ntp_time()
        print(f"Current Unix time from NTP: {current_unix_time}")
        print(f"Current UTC time from NTP: {time.gmtime(current_unix_time)}")
        print(f"Current local time from NTP: {time.ctime(current_unix_time)}")
    except Exception as e:
        print(f"Failed to get NTP time: {e}")

# Additional implementation at 2025-06-18 02:22:06
import ntplib
import datetime
import sys

def get_ntp_time(server_address='pool.ntp.org'):
    """
    Connects to an NTP server, retrieves time information, and displays it.
    Includes error handling for network issues and NTP specific errors.
    """
    client = ntplib.NTPClient()
    try:
        # Request NTP time from the specified server
        response = client.request(server_address, version=3)

        # Extract and display core time information
        print(f"NTP Server: {server_address}")
        print(f"Reference ID: {response.ref_id}")
        print(f"Stratum: {response.stratum}")
        print(f"Leap Indicator: {response.leap}")
        print(f"Precision: {response.precision}")
        print(f"Root Delay: {response.root_delay:.6f} seconds")
        print(f"Root Dispersion: {response.root_dispersion:.6f} seconds")

        # Calculate and display synchronized time
        # tx_time is the time when the server sent the response
        synchronized_time_utc = datetime.datetime.fromtimestamp(response.tx_time, datetime.timezone.utc)
        print(f"Server Transmit Time (UTC): {synchronized_time_utc.strftime('%Y-%m-%d %H:%M:%S.%f')[:-3]}")

        # Display network-related metrics
        print(f"Offset (local vs server): {response.offset:.6f} seconds")
        print(f"Delay (round trip): {response.delay:.6f} seconds")

        # Additional functionality: Display local time and difference
        local_time_now = datetime.datetime.now(datetime.timezone.utc)
        print(f"Local System Time (UTC): {local_time_now.strftime('%Y-%m-%d %H:%M:%S.%f')[:-3]}")

        # The offset is (local_time - server_time), so a positive offset means local is ahead
        # A more direct difference can be calculated from the synchronized time
        # This is the difference between the local system's current time and the time received from NTP
        time_difference = local_time_now.timestamp() - response.tx_time
        print(f"Difference (Local - Server): {time_difference:.6f} seconds")

    except ntplib.NTPException as e:
        print(f"NTP Error: Could not synchronize with {server_address}. {e}", file=sys.stderr)
        print("Possible reasons: Server unreachable, firewall blocking NTP (UDP port 123), or invalid server address.", file=sys.stderr)
    except Exception as e:
        print(f"An unexpected error occurred: {e}", file=sys.stderr)
        print("Please check your network connection or the server address.", file=sys.stderr)

if __name__ == "__main__":
    # Allow specifying a server address as a command-line argument
    if len(sys.argv) > 1:
        target_server = sys.argv[1]
        print(f"Attempting to get NTP time from: {target_server}")
        get_ntp_time(target_server)
    else:
        print("No server specified. Using default 'pool.ntp.org'.")
        print("Usage: python your_script_name.py [ntp_server_address]")
        get_ntp_time()

# Additional implementation at 2025-06-18 02:22:37
import ntplib
from datetime import datetime

class NTPClient:
    def __init__(self, server='pool.ntp.org'):
        self.client = ntplib.NTPClient()
        self.server = server

    def get_time_info(self):
        try:
            response = self.client.request(self.server, version=3)
            print(f"NTP Server: {self.server}")
            print(f"Reference ID: {ntplib.ref_id_to_text(response.ref_id)}")
            print(f"Stratum: {response.stratum}")
            print(f"Offset: {response.offset:.4f} seconds")
            print(f"Delay: {response.delay:.4f} seconds")
            print(f"NTP Time (UTC): {datetime.fromtimestamp(response.tx_time)}")
            print(f"Local Time (UTC): {datetime.utcfromtimestamp(response.orig_time)}")
            print(f"Root Delay: {response.root_delay:.4f} seconds")
            print(f"Root Dispersion: {response.root_dispersion:.4f} seconds")
            print(f"Leap Indicator: {response.leap}")
            print(f"Precision: {response.precision}")
            return response
        except ntplib.NTPException as e:
            print(f"NTP Error: {e}")
            return None
        except Exception as e:
            print(f"An unexpected error occurred: {e}")
            return None

    def get_current_ntp_time(self):
        try:
            response = self.client.request(self.server, version=3)
            return datetime.fromtimestamp(response.tx_time)
        except ntplib.NTPException as e:
            print(f"NTP Error getting current time: {e}")
            return None
        except Exception as e:
            print(f"An unexpected error occurred getting current time: {e}")
            return None

    def get_offset_and_delay(self):
        try:
            response = self.client.request(self.server, version=3)
            return response.offset, response.delay
        except ntplib.NTPException as e:
            print(f"NTP Error getting offset/delay: {e}")
            return None, None
        except Exception as e:
            print(f"An unexpected error occurred getting offset/delay: {e}")
            return None, None

if __name__ == "__main__":
    ntp_client = NTPClient(server='pool.ntp.org')

    print("--- Detailed NTP Time Information ---")
    ntp_client.get_time_info()

    print("\n--- Current NTP Time ---")
    current_ntp_time = ntp_client.get_current_ntp_time()
    if current_ntp_time:
        print(f"Current NTP Time (UTC): {current_ntp_time}")
        print(f"Current Local System Time (UTC): {datetime.utcnow()}")

    print("\n--- Offset and Delay ---")
    offset, delay = ntp_client.get_offset_and_delay()
    if offset is not None and delay is not None:
        print(f"Offset from server: {offset:.4f} seconds")
        print(f"Network Delay: {delay:.4f} seconds")

# Additional implementation at 2025-06-18 02:22:55
import socket
import struct
import datetime
import time

NTP_SERVER = "pool.ntp.org"
NTP_PORT = 123
NTP_EPOCH_OFFSET = 2208988800

def get_ntp_time_extended(server=NTP_SERVER, port=NTP_PORT, timeout=5):
    try:
        client = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        client.settimeout(timeout)

        request_packet = b'\x1b' + b'\x00' * 47

        t1_client = time.time()

        client.sendto(request_packet, (server, port))
        response_packet, addr = client.recvfrom(1024)

        t4_client = time.time()

        def ntp_to_float(ntp_timestamp_64bit):
            seconds = ntp_timestamp_64bit >> 32
            fraction = ntp_timestamp_64bit & 0xFFFFFFFF
            return float(seconds) + float(fraction) / (2**32)

        t_originate_full = struct.unpack('!Q', response_packet[24:32])[0]
        t_receive_full = struct.unpack('!Q', response_packet[32:40])[0]
        t_transmit_full = struct.unpack('!Q', response_packet[40:48])[0]

        t1_server_echo_float = ntp_to_float(t_originate_full) - NTP_EPOCH_OFFSET
        t2_server_float = ntp_to_float(t_receive_full) - NTP_EPOCH_OFFSET
        t3_server_float = ntp_to_float(t_transmit_full) - NTP_EPOCH_OFFSET

        ntp_time_unix_epoch = t3_server_float
        ntp_datetime = datetime.datetime.fromtimestamp(ntp_time_unix_epoch, datetime.timezone.utc)

        delay = (t4_client - t1_client) - (t3_server_float - t2_server_float)
        offset = ((t2_server_float - t1_client) + (t3_server_float - t4_client)) / 2.0

        return ntp_datetime, addr[0], offset, delay

    except socket.timeout:
        return None, None, None, None
    except socket.error:
        return None, None, None, None
    except Exception:
        return None, None, None, None
    finally:
        if 'client' in locals() and client:
            client.close()