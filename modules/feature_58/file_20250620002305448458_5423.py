import socket

HOST = '127.0.0.1'
PORT = 65432

def echo_server():
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server_socket:
        server_socket.bind((HOST, PORT))
        server_socket.listen()
        print(f"Server listening on {HOST}:{PORT}")
        while True:
            conn, addr = server_socket.accept()
            with conn:
                print(f"Connected by {addr}")
                while True:
                    data = conn.recv(1024)
                    if not data:
                        break
                    print(f"Received from {addr}: {data.decode()}")
                    conn.sendall(data)
                    print(f"Echoed to {addr}: {data.decode()}")
                print(f"Connection with {addr} closed")

if __name__ == "__main__":
    echo_server()

# Additional implementation at 2025-06-20 00:23:37
import socket
import threading
import sys

HOST = '127.0.0.1'
PORT = 12345
BUFFER_SIZE = 1024

def handle_client(client_socket, client_address):
    print(f"[*] Accepted connection from {client_address}")
    try:
        while True:
            data = client_socket.recv(BUFFER_SIZE)
            if not data:
                print(f"[*] Client {client_address} disconnected.")
                break
            decoded_data = data.decode('utf-8').strip()
            print(f"[*] Received from {client_address}: '{decoded_data}'")
            response_message = f"ECHO: {decoded_data}"
            client_socket.sendall(response_message.encode('utf-8'))
            print(f"[*] Sent to {client_address}: '{response_message}'")
    except ConnectionResetError:
        print(f"[*] Client {client_address} forcibly closed the connection.")
    except Exception as e:
        print(f"[*] Error handling client {client_address}: {e}")
    finally:
        client_socket.close()
        print(f"[*] Connection with {client_address} closed.")

def start_server():
    server_socket = None
    try:
        server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        server_socket.bind((HOST, PORT))
        server_socket.listen(5)
        print(f"[*] Listening on {HOST}:{PORT}")
        print("Press Ctrl+C to stop the server.")
        while True:
            client_socket, client_address = server_socket.accept()
            client_handler = threading.Thread(target=handle_client, args=(client_socket, client_address))
            client_handler.daemon = True
            client_handler.start()
    except KeyboardInterrupt:
        print("\n[*] Server is shutting down...")
    except Exception as e:
        print(f"[*] An error occurred: {e}")
    finally:
        if server_socket:
            server_socket.close()
            print("[*] Server socket closed.")
        sys.exit(0)

if __name__ == "__main__":
    start_server()

# Additional implementation at 2025-06-20 00:24:39
import socketserver
import threading
import time
import datetime

class ThreadedEchoHandler(socketserver.BaseRequestHandler):
    def handle(self):
        client_address = self.client_address
        print(f"[{datetime.datetime.now()}] Connected: {client_address}")
        try:
            while True:
                # Receive data from the client
                data = self.request.recv(1024).strip()
                if not data:
                    # Client disconnected or sent empty data
                    print(f"[{datetime.datetime.now()}] Disconnected: {client_address} (No data)")
                    break

                message = data.decode('utf-8')
                print(f"[{datetime.datetime.now()}] Received from {client_address}: '{message}'")

                # Process additional functionality: commands
                if message.upper() == 'QUIT':
                    response = "Goodbye! Connection closing."
                    self.request.sendall(response.encode('utf-8'))
                    print(f"[{datetime.datetime.now()}] Client {client_address} requested QUIT.")
                    break
                elif message.upper() == 'INFO':
                    server_info = f"Server Time: {datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n" \
                                  f"Server Address: {self.server.server_address[0]}:{self.server.server_address[1]}\n" \
                                  f"Thread ID: {threading.get_ident()}"
                    self.request.sendall(server_info.encode('utf-8'))
                else:
                    # Echo functionality with timestamp
                    echo_response = f"[{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] [ECHO] {message}"
                    self.request.sendall(echo_response.encode('utf-8'))
        except ConnectionResetError:
            print(f"[{datetime.datetime.now()}] Connection reset by client: {client_address}")
        except Exception as e:
            print(f"[{datetime.datetime.now()}] Error with {client_address}: {e}")
        finally:
            # Ensure the connection is closed
            self.request.close()
            print(f"[{datetime.datetime.now()}] Connection closed for {client_address}")

class ThreadedEchoServer(socketserver.ThreadingMixIn, socketserver.TCPServer):
    # Enable daemon threads so they exit when the main program exits
    daemon_threads = True
    # Allow the server address to be reused immediately after shutdown
    allow_reuse_address = True

    def __init__(self, server_address, RequestHandlerClass):
        super().__init__(server_address, RequestHandlerClass)
        print(f"[{datetime.datetime.now()}] Server listening on {server_address[0]}:{server_address[1]}...")

def run_server():
    HOST, PORT = "localhost", 9999
    server = ThreadedEchoServer((HOST, PORT), ThreadedEchoHandler)

    # Start a thread with the server -- that thread will then start one more thread for each request
    server_thread = threading.Thread(target=server.serve_forever)
    # Exit the server thread when the main thread terminates
    server_thread.daemon = True
    server_thread.start()
    print(f"[{datetime.datetime.now()}] Server loop running in thread: {server_thread.name}")

    # Keep the main thread alive to prevent daemon threads from exiting
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print(f"[{datetime.datetime.now()}] Shutting down server...")
        server.shutdown()
        server.server_close()
        print(f"[{datetime.datetime.now()}] Server shut down.")

if __name__ == "__main__":
    run_server()