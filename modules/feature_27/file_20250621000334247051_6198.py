import socket
import os
import struct
import sys

HOST = '127.0.0.1'
PORT = 65432
BUFFER_SIZE = 4096

def start_server():
    if not os.path.exists('received_files'):
        os.makedirs('received_files')

    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.bind((HOST, PORT))
        s.listen()
        print(f"Server listening on {HOST}:{PORT}")
        while True:
            conn, addr = s.accept()
            with conn:
                print(f"Connected by {addr}")
                try:
                    filename_len_bytes = conn.recv(4)
                    if not filename_len_bytes:
                        print("Client disconnected before sending filename length.")
                        continue
                    filename_len = struct.unpack('!I', filename_len_bytes)[0]

                    filename = conn.recv(filename_len).decode('utf-8')
                    print(f"Receiving file: {filename}")

                    file_size_bytes = conn.recv(8)
                    if not file_size_bytes:
                        print("Client disconnected before sending file size.")
                        continue
                    file_size = struct.unpack('!Q', file_size_bytes)[0]

                    received_bytes = 0
                    output_filepath = os.path.join('received_files', filename)
                    with open(output_filepath, 'wb') as f:
                        while received_bytes < file_size:
                            data = conn.recv(min(BUFFER_SIZE, file_size - received_bytes))
                            if not data:
                                break
                            f.write(data)
                            received_bytes += len(data)
                    print(f"File '{filename}' received successfully. Size: {received_bytes} bytes")

                except Exception as e:
                    print(f"Error during transfer: {e}")
                finally:
                    conn.close()

def send_file(filepath):
    if not os.path.exists(filepath):
        print(f"Error: File '{filepath}' not found.")
        return

    filename = os.path.basename(filepath)
    file_size = os.path.getsize(filepath)

    try:
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            s.connect((HOST, PORT))
            print(f"Connected to server {HOST}:{PORT}")

            filename_bytes = filename.encode('utf-8')
            filename_len = len(filename_bytes)
            s.sendall(struct.pack('!I', filename_len))

            s.sendall(filename_bytes)

            s.sendall(struct.pack('!Q', file_size))

            sent_bytes = 0
            with open(filepath, 'rb') as f:
                while sent_bytes < file_size:
                    data = f.read(BUFFER_SIZE)
                    if not data:
                        break
                    s.sendall(data)
                    sent_bytes += len(data)
            print(f"File '{filename}' sent successfully. Size: {sent_bytes} bytes")

    except Exception as e:
        print(f"Error during transfer: {e}")

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Usage: python script.py server")
        print("       python script.py client <filepath>")
        sys.exit(1)

    mode = sys.argv[1].lower()

    if mode == 'server':
        start_server()
    elif mode == 'client':
        if len(sys.argv) < 3:
            print("Usage: python script.py client <filepath>")
            sys.exit(1)
        filepath_to_send = sys.argv[2]
        send_file(filepath_to_send)
    else:
        print("Invalid mode. Use 'server' or 'client'.")
        sys.exit(1)

# Additional implementation at 2025-06-21 00:04:22
import socket
import os
import sys
import threading
import json

HOST = '127.0.0.1'
PORT = 65432
FILE_DIR = 'shared_files'
BUFFER_SIZE = 4096

if not os.path.exists(FILE_DIR):
    os.makedirs(FILE_DIR)

def handle_client(conn, addr):
    print(f"Connected by {addr}")
    try:
        while True:
            header_buffer = b''
            while True:
                chunk = conn.recv(1)
                if not chunk:
                    break
                header_buffer += chunk
                try:
                    command_data = json.loads(header_buffer.decode('utf-8'))
                    break
                except json.JSONDecodeError:
                    if len(header_buffer) > BUFFER_SIZE * 2: # Prevent infinite buffer growth
                        raise ValueError("Header too large or malformed")
                    continue
            
            if not header_buffer:
                break

            command = command_data.get('command')
            filename = command_data.get('filename')
            file_size = command_data.get('size')

            if command == 'UPLOAD':
                print(f"Receiving file: {filename} from {addr}")
                filepath = os.path.join(FILE_DIR, filename)
                try:
                    with open(filepath, 'wb') as f:
                        bytes_received = 0
                        while bytes_received < file_size:
                            data = conn.recv(min(BUFFER_SIZE, file_size - bytes_received))
                            if not data:
                                break
                            f.write(data)
                            bytes_received += len(data)
                    print(f"Successfully received {filename}")
                    conn.sendall(b"UPLOAD_SUCCESS")
                except Exception as e:
                    print(f"Error receiving file {filename}: {e}")
                    conn.sendall(f"UPLOAD_ERROR: {e}".encode('utf-8'))

            elif command == 'DOWNLOAD':
                print(f"Client {addr} requested to download: {filename}")
                filepath = os.path.join(FILE_DIR, filename)
                if os.path.exists(filepath) and os.path.isfile(filepath):
                    try:
                        file_size = os.path.getsize(filepath)
                        response_header = json.dumps({'status': 'READY', 'filename': filename, 'size': file_size}).encode('utf-8')
                        conn.sendall(response_header + b'\n')
                        
                        ack = conn.recv(BUFFER_SIZE).decode('utf-8')
                        if ack == "ACK_READY":
                            with open(filepath, 'rb') as f:
                                while True:
                                    bytes_read = f.read(BUFFER_SIZE)
                                    if not bytes_read:
                                        break
                                    conn.sendall(bytes_read)
                            print(f"Successfully sent {filename} to {addr}")
                        else:
                            print(f"Client {addr} did not acknowledge download readiness.")
                    except Exception as e:
                        print(f"Error sending file {filename}: {e}")
                        error_header = json.dumps({'status': 'ERROR', 'message': str(e)}).encode('utf-8')
                        conn.sendall(error_header + b'\n')
                else:
                    print(f"File not found: {filename}")
                    error_header = json.dumps({'status': 'NOT_FOUND', 'filename': filename}).encode('utf-8')
                    conn.sendall(error_header + b'\n')

            elif command == 'LIST':
                print(f"Client {addr} requested file list.")
                files = [f for f in os.listdir(FILE_DIR) if os.path.isfile(os.path.join(FILE_DIR, f))]
                response_header = json.dumps({'status': 'LIST_READY', 'files': files}).encode('utf-8')
                conn.sendall(response_header + b'\n')

            elif command == 'QUIT':
                print(f"Client {addr} disconnected.")
                break
            else:
                print(f"Unknown command from {addr}: {command}")
                conn.sendall(b"UNKNOWN_COMMAND")

    except ConnectionResetError:
        print(f"Client {addr} disconnected unexpectedly.")
    except Exception as e:
        print(f"Error handling client {addr}: {e}")
    finally:
        conn.close()
        print(f"Connection with {addr} closed.")

def start_server():
    print(f"Server starting on {HOST}:{PORT}")
    print(f"Serving files from: {os.path.abspath(FILE_DIR)}")
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    try:
        server_socket.bind((HOST, PORT))
        server_socket.listen(5)
        print("Server listening...")
        while True:
            conn, addr = server_socket.accept()
            client_thread = threading.Thread(target=handle_client, args=(conn, addr))
            client_thread.start()
    except Exception as e:
        print(f"Server error: {e}")
    finally:
        server_socket.close()

def send_command(sock, command, filename=None, file_size=None):
    header_data = {'command': command}
    if filename:
        header_data['filename'] = filename
    if file_size is not None:
        header_data['size'] = file_size
    header_json = json.dumps(header_data)
    sock.sendall(header_json.encode('utf-8'))

def receive_response_header(sock):
    buffer = b''
    while True:
        chunk = sock.recv(1)
        if not chunk:
            return None
        buffer += chunk
        if buffer.endswith(b'\n'):
            break
    try:
        return json.loads(buffer.decode('utf-8').strip())
    except json.JSONDecodeError:
        print(f"Error decoding JSON header: {buffer.decode('utf-8')}")
        return None

def client_upload(sock, filepath):
    if not os.path.exists(filepath) or not os.path.isfile(filepath):
        print(f"Error: File not found locally: {filepath}")
        return

    filename = os.path.basename(filepath)
    file_size = os.path.getsize(filepath)

    print(f"Uploading {filename} ({file_size} bytes)...")
    send_command(sock, 'UPLOAD', filename, file_size)

    try:
        with open(filepath, 'rb') as f:
            bytes_sent = 0
            while bytes_sent < file_size:
                bytes_read = f.read(BUFFER_SIZE)
                if not bytes_read:
                    break
                sock.sendall(bytes_read)
                bytes_sent += len(bytes_read)
        print(f"File {filename} sent.")
        response = sock.recv(BUFFER_SIZE).decode('utf-8')
        if response == "UPLOAD_SUCCESS":
            print(f"Server confirmed upload success for {filename}.")
        else:
            print(f"Server reported an error during upload: {response}")
    except Exception as e:
        print(f"Error during upload: {e}")

def client_download(sock, filename):
    print(f"Requesting download of {filename}...")
    send_command(sock, 'DOWNLOAD', filename)

    response_header = receive_response_header(sock)
    if not response_header:
        print("Failed to receive response header from server.")
        return

    status = response_header.get('status')
    if status == 'READY':
        remote_filename = response_header.get('filename')
        file_size = response_header.get('size')
        print(f"Server ready to send {remote_filename} ({file_size} bytes).")

        local_filepath = os.path.join(FILE_DIR, remote_filename)
        sock.sendall(b"ACK_READY")

        try:
            with open(local_filepath, 'wb') as f:
                bytes_received = 0
                while bytes_received < file_size:
                    data = sock.recv(min(BUFFER_SIZE, file_size - bytes_received))
                    if not data:
                        print("Connection closed prematurely during download.")
                        break
                    f.write(data)
                    bytes_received += len(data)
            print(f"Successfully downloaded {remote_filename} to {local_filepath}")
        except Exception as e:
            print(f"Error during download: {e}")
    elif status == 'NOT_FOUND':
        print(f"Server reported file not found: {filename}")
    elif status == 'ERROR':
        print(f"Server reported an error: {response_header.get('message', 'Unknown error')}")
    else:
        print(f"Unexpected server response status: {status}")

def client_list_files(sock):
    print("Requesting file list from server...")
    send_command(sock, 'LIST')

    response_header = receive_response_header(sock)
    if not response_header:
        print("Failed to receive response header from server.")
        return

    status = response_header.get('status')
    if status == 'LIST_READY':
        files = response_header.get('files', [])
        print("\n--- Files on Server ---")
        if files:
            for f in files:
                print(f"- {f}")
        else:
            print("No files found on server.")
        print("-----------------------")
    else:
        print(f"Unexpected server response status for list: {status}")

def start_client():
    print(f"Client connecting to {HOST}:{PORT}")
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        client_socket.connect((HOST, PORT))
        print("Connected to server.")

        while True:
            print("\n--- File Transfer Client ---")
            print("1. Upload File")
            print("2. Download File")
            print("3. List Files on Server")
            print("4. Quit")
            choice = input("Enter your choice: ")

            if choice == '1':
                filepath = input("Enter path to file to upload: ")
                client_upload(client_socket, filepath)
            elif choice == '2':
                filename = input("Enter filename to download: ")
                client_download(client_socket, filename)
            elif choice == '3':
                client_list_files(client_socket)
            elif choice == '4':
                print("Disconnecting from server.")
                send_command(client_socket, 'QUIT')
                break
            else:
                print("Invalid choice. Please try again.")

    except ConnectionRefusedError:
        print("Connection refused. Make sure the server is running.")
    except Exception as e:
        print(f"Client error: {e}")
    finally:
        client_socket.close()

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python file_transfer_tool.py [server|client]")
        sys.exit(1)

    mode = sys.argv[1].lower()
    if mode == 'server':
        start_server()
    elif mode == 'client':
        start_client()
    else:
        print("Invalid mode. Please choose 'server' or 'client'.")
        sys.exit(1)

# Additional implementation at 2025-06-21 00:05:52
import socket
import threading
import os
import sys

# --- Configuration ---
HOST = '127.0.0.1'
PORT = 12345
SERVER_FILE_DIR = 'server_files'
CLIENT_DOWNLOAD_DIR = 'client_downloads'
BUFFER_SIZE = 4096 # Standard buffer size for network operations

# Ensure directories exist
if not os.path.exists(SERVER_FILE_DIR):
    os.makedirs(SERVER_FILE_DIR)
if not os.path.exists(CLIENT_DOWNLOAD_DIR):
    os.makedirs(CLIENT_DOWNLOAD_DIR)

# --- Helper Functions for Network Communication ---

def _send_prefixed_message(sock, message_bytes):
    """Sends a message prefixed by its 8-byte length."""
    msg_len = str(len(message_bytes)).ljust(8).encode('utf-8')
    sock.sendall(msg_len)
    sock.sendall(message_bytes)

def _recv_prefixed_message(sock):
    """Receives a message prefixed by its 8-byte length."""
    len_bytes = sock.recv(8)
    if not len_bytes:
        return None # Connection closed
    
    try:
        msg_len = int(len_bytes.decode('utf-8').strip())
    except ValueError:
        print("Error: Invalid length header received.")
        return None

    chunks = []
    bytes_recd = 0
    while bytes_recd < msg