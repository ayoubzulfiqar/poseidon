import socket
import threading
import sys

PROXY_HOST = '127.0.0.1'
PROXY_PORT = 8888
BUFFER_SIZE = 4096

class ProxyThread(threading.Thread):
    def __init__(self, client_socket, client_address):
        threading.Thread.__init__(self)
        self.client_socket = client_socket
        self.client_address = client_address
        self.target_socket = None

    def run(self):
        try:
            first_chunk = self.client_socket.recv(BUFFER_SIZE)
            if not first_chunk:
                return

            first_line_str = first_chunk.decode('latin-1').split('\n')[0]
            parts = first_line_str.split(' ')

            method = parts[0]
            target_address = None
            target_port = None

            if method == 'CONNECT':
                host_port = parts[1]
                if ':' in host_port:
                    target_address, target_port_str = host_port.split(':')
                    target_port = int(target_port_str)
                else:
                    target_address = host_port
                    target_port = 443
                
                self.target_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
                self.target_socket.connect((target_address, target_port))

                self.client_socket.sendall(b"HTTP/1.1 200 Connection established\r\n\r\n")

                self._relay_data(self.client_socket, self.target_socket)

            else:
                full_request_headers = first_chunk.decode('latin-1')
                
                host_header = None
                for line in full_request_headers.split('\n'):
                    if line.lower().startswith('host:'):
                        host_header = line.split(':')[1].strip()
                        break
                
                if host_header:
                    target_address = host_header
                    target_port = 80
                    
                    self.target_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
                    self.target_socket.connect((target_address, target_port))

                    self.target_socket.sendall(first_chunk)

                    self._relay_data(self.client_socket, self.target_socket)
                else:
                    self.client_socket.sendall(b"HTTP/1.1 400 Bad Request\r\n\r\n")

        except Exception:
            pass
        finally:
            if self.client_socket:
                self.client_socket.close()
            if self.target_socket:
                self.target_socket.close()

    def _relay_data(self, sock1, sock2):
        while True:
            try:
                data_from_1 = sock1.recv(BUFFER_SIZE)
                if not data_from_1:
                    break
                sock2.sendall(data_from_1)

                data_from_2 = sock2.recv(BUFFER_SIZE)
                if not data_from_2:
                    break
                sock1.sendall(data_from_2)
            except socket.error:
                break
            except Exception:
                break

def main():
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    server_socket.bind((PROXY_HOST, PROXY_PORT))
    server_socket.listen(5)

    while True:
        try:
            client_socket, client_address = server_socket.accept()
            proxy_thread = ProxyThread(client_socket, client_address)
            proxy_thread.daemon = True
            proxy_thread.start()
        except KeyboardInterrupt:
            break
        except Exception:
            pass
    server_socket.close()

if __name__ == '__main__':
    main()

# Additional implementation at 2025-06-19 23:27:32
import socket
import threading
import sys
import datetime

class ProxyServer:
    def __init__(self, host='127.0.0.1', port=8888, buffer_size=4096):
        self.proxy_host = host
        self.proxy_port = port
        self.buffer_size = buffer_size
        self.server_socket = None
        self.blocked_domains = {"example.com", "badsite.net"} # Example blocked domains
        self.custom_response_header = "X-Proxy-By: PythonProxy/1.0"

    def _log(self, message):
        timestamp = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        print(f"[{timestamp}] {message}")

    def start(self):
        try:
            self.server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            self.server_socket.bind((self.proxy_host, self.proxy_port))
            self.server_socket.listen(5)
            self._log(f"Proxy server listening on {self.proxy_host}:{self.proxy_port}")

            while True:
                client_socket, client_address = self.server_socket.accept()
                self._log(f"Accepted connection from {client_address[0]}:{client_address[1]}")
                client_handler = threading.Thread(target=self.handle_client, args=(client_socket, client_address))
                client_handler.daemon = True
                client_handler.start()

        except KeyboardInterrupt:
            self._log("Proxy server shutting down.")
        except Exception as e:
            self._log(f"Server error: {e}")
        finally:
            if self.server_socket:
                self.server_socket.close()

    def handle_client(self, client_socket, client_address):
        try:
            first_line = client_socket.recv(self.buffer_size).decode('latin-1')
            if not first_line:
                return

            lines = first_line.split('\r\n')
            request_line = lines[0]
            method, path, http_version = self._parse_request_line(request_line)

            if method == "CONNECT":
                self.handle_https(client_socket, request_line)
            else:
                # Reconstruct the full initial request bytes
                full_request_bytes = first_line.encode('latin-1')
                self.handle_http(client_socket, full_request_bytes, request_line)

        except socket.error as e:
            self._log(f"Socket error handling client {client_address}: {e}")
        except Exception as e:
            self._log(f"Error handling client {client_address}: {e}")
        finally:
            client_socket.close()

    def handle_http(self, client_socket, initial_request_bytes, request_line):
        try:
            headers_str = initial_request_bytes.decode('latin-1').split('\r\n\r\n', 1)[0]
            host, port = self._get_host_port_from_headers(headers_str)

            if not host:
                client_socket.sendall(b"HTTP/1.1 400 Bad Request\r\n\r\n")
                self._log(f"Bad request: No Host header found in {request_line}")
                return

            if host in self.blocked_domains:
                client_socket.sendall(b"HTTP/1.1 403 Forbidden\r\nContent-Type: text/plain\r\n\r\nAccess to this domain is blocked by the proxy.\r\n")
                self._log(f"Blocked access to {host} for HTTP request: {request_line}")
                return

            self._log(f"HTTP Request: {request_line} -> {host}:{port}")

            target_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            target_socket.connect((host, port))
            target_socket.sendall(initial_request_bytes)

            response_buffer = b""
            while True:
                data = target_socket.recv(self.buffer_size)
                if not data:
                    break
                response_buffer += data

                # Simple check for end of headers (double CRLF)
                if b'\r\n\r\n' in response_buffer:
                    break
            
            # Modify response headers
            modified_response_bytes = self._add_custom_header(response_buffer)
            
            client_socket.sendall(modified_response_bytes)

            # Continue forwarding the rest of the body if not fully received
            while True:
                data = target_socket.recv(self.buffer_size)
                if not data:
                    break
                client_socket.sendall(data)

            self._log(f"HTTP Response forwarded for {request_line}")

        except socket.timeout:
            self._log(f"Timeout connecting to {host}:{port}")
            client_socket.sendall(b"HTTP/1.1 504 Gateway Timeout\r\n\r\n")
        except socket.error as e:
            self._log(f"Socket error in HTTP handling for {request_line}: {e}")
            client_socket.sendall(b"HTTP/1.1 502 Bad Gateway\r\n\r\n")
        except Exception as e:
            self._log(f"Error in HTTP handling for {request_line}: {e}")
            client_socket.sendall(b"HTTP/1.1 500 Internal Proxy Error\r\n\r\n")
        finally:
            if 'target_socket' in locals() and target_socket:
                target_socket.close()

    def handle_https(self, client_socket, request_line):
        try:
            host, port = self._get_host_port_from_connect(request_line)

            if not host:
                client_socket.sendall(b"HTTP/1.1 400 Bad Request\r\n\r\n")
                self._log(f"Bad CONNECT request: {request_line}")
                return

            if host in self.blocked_domains:
                client_socket.sendall(b"HTTP/1.1 403 Forbidden\r\nContent-Type: text/plain\r\n\r\nAccess to this domain is blocked by the proxy.\r\n")
                self._log(f"Blocked access to {host} for HTTPS request: {request_line}")
                return

            self._log(f"HTTPS CONNECT: {request_line} -> {host}:{port}")

            target_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            target_socket.connect((host, port))

            client_socket.sendall(b"HTTP/1.1 200 Connection established\r\n\r\n")
            self._log(f"Connection established for {host}:{port}")

            # Bidirectional forwarding
            forward_client_to_target = threading.Thread(target=self._forward_data, args=(client_socket, target_socket, f"Client->Target {host}:{port}"))
            forward_target_to_client = threading.Thread(target=self._forward_data, args=(target_socket, client_socket, f"Target->Client {host}:{port}"))

            forward_client_to_target.daemon = True
            forward_target_to_client.daemon = True

            forward_client_to_target.start()
            forward_target_to_client.start()

            # Wait for both threads to finish (or for an error to occur)
            forward_client_to_target.join()
            forward_target_to_client.join()

            self._log(f"HTTPS connection closed for {host}:{port}")

        except socket.timeout:
            self._log(f"Timeout connecting to {host}:{port}")
            client_socket.sendall(b"HTTP/1.1 504 Gateway Timeout\r\n\r\n")
        except socket.error as e:
            self._log(f"Socket error in HTTPS handling for {request_line}: {e}")
            client_socket.sendall(b"HTTP/1.1 502 Bad Gateway\r\n\r\n")
        except Exception as e:
            self._log(f"Error in HTTPS handling for {request_line}: {e}")
            client_socket.sendall(b"HTTP/1.1 500 Internal Proxy Error\r\n\r\n")
        finally:
            if 'target_socket' in locals() and target_socket:
                target_socket.close()

    def _forward_data(self, source_socket, destination_socket, log_prefix=""):
        try:
            while True:
                data = source_socket.recv(self.buffer_size)
                if not data:
                    break
                destination_socket.sendall(data)
        except socket.error as e:
            self._log(f"Forwarding error ({log_prefix}): {e}")
        except Exception as e:
            self._log(f"Unexpected error during forwarding ({log_prefix}): {e}")
        finally:
            # It's important not to close the sockets here directly
            # as they are managed by the main handle_client/handle_https logic.
            # The connection will naturally close when one side breaks.
            pass

    def _parse_request_line(self, request_line):
        parts = request_line.split(' ', 2)
        if len(parts) == 3:
            return parts[0], parts[1], parts[2]
        return None, None, None

    def _get_host_port_from_headers(self, headers_str):
        host = None
        port = 80 # Default for HTTP
        for line in headers_str.split('\r\n'):
            if line.lower().startswith("host:"):
                host_port_str = line[len("host:"):].strip()
                if ':' in host_port_str:
                    host, port_str = host_port_str.split(':')
                    try:
                        port = int(port_str)
                    except ValueError:
                        pass # Keep default port if invalid
                else:
                    host = host_port_str
                break
        return host, port

    def _get_host_port_from_connect(self, request_line):
        # Example: CONNECT www.google.com:443 HTTP/1.1
        parts = request_line.split(' ')
        if len(parts) >= 2:
            host_port_str = parts[1]
            if ':' in host_port_str:
                host, port_str = host_port_str.split(':')
                try:
                    port = int(port_str)
                    return host, port
                except ValueError:
                    pass
        return None, None # Invalid format

    def _add_custom_header(self, response_bytes):
        # Find the end of the headers (first \r\n\r\n)
        header_end_index = response_bytes.find(b'\r\n\r\n')
        if header_end_index == -1:
            return response_bytes # No headers found or malformed

        headers_part = response_bytes[:header_end_index + 2] # +2 to include the first \r\n of the double CRLF
        body_part = response_bytes[header_end_index + 4:] # +4 to skip the double CRLF

        # Convert headers to string to easily insert
        headers_str = headers_part.decode('latin-1')

        # Find the first CRLF after the status line to insert the header
        # This assumes the status line is the first line.
        first_crlf_index = headers_str.find('\r\n')
        if first_crlf_index == -1:
            return response_bytes # Malformed headers

        # Insert the custom header after the status line
        modified_headers_str = (
            headers_str[:first_crlf_index + 2] +
            self.custom_response_header + "\r\n" +
            headers_str[first_crlf_index + 2:]
        )

        return modified_headers_str.encode('latin-1') + body_part

if __name__ == '__main__':
    proxy = ProxyServer()
    proxy.start()

# Additional implementation at 2025-06-19 23:28:13
import http.server
import socketserver
import requests
import sys

PROXY_PORT = 8888
TARGET_HOST = "httpbin.org"
TARGET_PORT = 80

BLOCKED_PATHS = ["/deny", "/status/403"]

class ProxyHandler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        self._handle_request()

    def do_POST(self):
        self._handle_request()

    def do_PUT(self):
        self._handle_request()

    def do_DELETE(self):
        self._handle_request()

    def do_HEAD(self):
        self._handle_request()

    def do_OPTIONS(self):
        self._handle_request()

    def do_PATCH(self):
        self._handle_request()

    def _handle_request(self):
        print(f"--- Incoming Request ---")
        print(f"Client: {self.client_address[0]}:{self.client_address[1]}")
        print(f"Method: {self.command}")
        print(f"Path: {self.path}")
        print(f"Headers:")
        for header, value in self.headers.items():
            print(f"  {header}: {value}")

        for blocked_path in BLOCKED_PATHS:
            if self.path.startswith(blocked_path):
                print(f"--- Request Blocked: {self.path} ---")
                self.send_error(403, "Forbidden: This path is blocked by the proxy.")
                return

        target_url = f"http://{TARGET_HOST}:{TARGET_PORT}{self.path}"

        req_headers = {}
        EXCLUDE_REQUEST_HEADERS = ['connection', 'keep-alive', 'proxy-authenticate', 'proxy-authorization', 'te', 'trailers', 'transfer-encoding', 'upgrade', 'host']
        for header, value in self.headers.items():
            if header.lower() not in EXCLUDE_REQUEST_HEADERS:
                req_headers[header] = value
        
        content_length = int(self.headers.get('Content-Length', 0))
        request_body = self.rfile.read(content_length) if content_length > 0 else None

        try:
            response = requests.request(
                method=self.command,
                url=target_url,
                headers=req_headers,
                data=request_body,
                allow_redirects=False,
                stream=True
            )

            print(f"--- Outgoing Response ---")
            print(f"Status: {response.status_code}")
            print(f"Headers:")
            for header, value in response.headers.items():
                print(f"  {header}: {value}")

            self.send_response(response.status_code)

            self.send_header("X-Proxy-Modified", "True")
            self.send_header("X-Proxy-Info", "Python Localhost Proxy")

            EXCLUDE_RESPONSE_HEADERS = ['connection', 'keep-alive', 'proxy-authenticate', 'proxy-authorization', 'te', 'trailers', 'transfer-encoding', 'upgrade', 'content-length']
            for header, value in response.headers.items():
                if header.lower() not in EXCLUDE_RESPONSE_HEADERS:
                    self.send_header(header, value)
            
            self.end_headers()

            for chunk in response.iter_content(chunk_size=8192):
                self.wfile.write(chunk)
            
            print(f"--- Request Handled: {self.path} ---")

        except requests.exceptions.RequestException as e:
            print(f"--- Proxy Error: Request to target failed: {e} ---")
            self.send_error(502, f"Bad Gateway: Could not connect to target server. Error: {e}")
        except Exception as e:
            print(f"--- Unexpected Proxy Error: {e} ---")
            self.send_error(500, f"Internal Proxy Error: {e}")

class ThreadingHTTPServer(socketserver.ThreadingMixIn, http.server.HTTPServer):
    pass

if __name__ == "__main__":
    print(f"Starting proxy server on port {PROXY_PORT}")
    print(f"Forwarding requests to {TARGET_HOST}:{TARGET_PORT}")
    print(f"Blocked paths: {BLOCKED_PATHS}")

    try:
        with ThreadingHTTPServer(("", PROXY_PORT), ProxyHandler) as httpd:
            httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nProxy server stopped.")
    except Exception as e:
        print(f"An unhandled error occurred: {e}", file=sys.stderr)

# Additional implementation at 2025-06-19 23:29:21
import socket
import threading
import select
import sys
import datetime
from urllib.parse import urlparse

PROXY_HOST = '127.0.0.1'
PROXY_PORT = 8080
BUFFER_SIZE = 4096
TIMEOUT = 60

BLOCKED_DOMAINS = [
    'example.com',
    'badsite.net'
]
CUSTOM_RESPONSE_HEADER_NAME = 'X-Proxy-Served-By'
CUSTOM_RESPONSE_HEADER_VALUE = 'MyPythonProxy'

def log_message(level, message):
    timestamp = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
    print(f"[{timestamp}] [{level}] {message}")

class ProxyHandler(threading.Thread):
    def __init__(self, client_socket, client_address):
        super().__init__()
        self.client_socket = client_socket
        self.client_address = client_address
        self.target_socket = None
        self.is_https = False

    def run(self):
        try:
            self.handle_request()
        except Exception as e:
            log_message("ERROR", f"Error handling client {self.client_address}: {e}")
        finally:
            self.close_connections()

    def handle_request(self):
        first_line_bytes = self.client_socket.recv(BUFFER_SIZE)
        if not first_line_bytes:
            return

        first_line = first_line_bytes.decode('latin-1', errors='ignore')
        
        first_line_end = first_line.find('\r\n')
        if first_line_end == -1:
            log_message("WARNING", f"Incomplete first line from {self.client_address}")
            return

        request_line = first_line[:first_line_end]
        log_message("INFO", f"Request from {self.client_address}: {request_line}")

        parts = request_line.split(' ')
        if len(parts) < 3:
            log_message("WARNING", f"Malformed request line from {self.client_address}: {request_line}")
            return

        method, url, http_version = parts[0], parts[1], parts[2]

        if method == 'CONNECT':
            self.is_https = True
            host_port = url.split(':')
            target_host = host_port[0]
            target_port = int(host_port[1]) if len(host_port) > 1 else 443
        else:
            parsed_url = urlparse(url)
            target_host = parsed_url.hostname
            target_port = parsed_url.port if parsed_url.port else (80 if parsed_url.scheme == 'http' else 443)
            
            path = parsed_url.path
            if parsed_url.query:
                path += '?' + parsed_url.query
            if parsed_url.fragment:
                path += '#' + parsed_url.fragment
            
            request_line = f"{method} {path} {http_version}"

        if not target_host:
            log_message("WARNING", f"Could not determine target host from URL: {url}")
            return

        if any(blocked_domain in target_host for blocked_domain in BLOCKED_DOMAINS):
            log_message("BLOCKED", f"Blocking request to {target_host} from {self.client_address}")
            self.client_socket.sendall(b"HTTP/1.1 403 Forbidden\r\nContent-Length: 19\r\n\r\nDomain Blocked By Proxy")
            return

        try:
            self.target_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.target_socket.settimeout(TIMEOUT)
            self.target_socket.connect((target_host, target_port))
            log_message("INFO", f"Connected to target: {target_host}:{target_port}")
        except Exception as e:
            log_message("ERROR", f"Could not connect to target {target_host}:{target_port}: {e}")
            self.client_socket.sendall(b"HTTP/1.1 502 Bad Gateway\r\n\r\n")
            return

        if self.is_https:
            self.client_socket.sendall(b"HTTP/1.1 200 Connection Established\r\n\r\n")
            self.tunnel_data()
        else:
            full_request = first_line_bytes.replace(first_line_bytes[:first_line_end], request_line.encode('latin-1'), 1)
            self.target_socket.sendall(full_request)
            self.forward_http_response()

    def tunnel_data(self):
        sockets = [self.client_socket, self.target_socket]
        while True:
            readable, _, _ = select.select(sockets, [], [], TIMEOUT)
            if not readable:
                log_message("WARNING", f"Timeout during data tunneling for {self.client_address}")
                break

            for sock in readable:
                try:
                    data = sock.recv(BUFFER_SIZE)
                    if not data:
                        log_message("INFO", f"Connection closed by {sock.getpeername()}")
                        return

                    if sock is self.client_socket:
                        self.target_socket.sendall(data)
                    else:
                        self.client_socket.sendall(data)
                except socket.timeout:
                    log_message("WARNING", f"Socket timeout during tunneling for {self.client_address}")
                    return
                except Exception as e:
                    log_message("ERROR", f"Error during data tunneling for {self.client_address}: {e}")
                    return

    def forward_http_response(self):
        response_data = b""
        try:
            while True:
                chunk = self.target_socket.recv(BUFFER_SIZE)
                if not chunk:
                    break
                response_data += chunk
                if b'\r\n\r\n' in response_data:
                    break
            
            if not response_data:
                log_message("WARNING", f"No response received from target for {self.client_address}")
                return

            headers_end_index = response_data.find(b'\r\n\r\n')
            if headers_end_index == -1:
                log_message("WARNING", f"Incomplete response headers from target for {self.client_address}")
                self.client_socket.sendall(response_data)
                return

            headers_raw = response_data[:headers_end_index + 4]
            body_data = response_data[headers_end_index + 4:]

            modified_headers_raw = self.add_custom_header(headers_raw)
            
            self.client_socket.sendall(modified_headers_raw)
            self.client_socket.sendall(body_data)

            while True:
                data = self.target_socket.recv(BUFFER_SIZE)
                if not data:
                    break
                self.client_socket.sendall(data)

            log_message("INFO", f"HTTP response forwarded to {self.client_address}")

        except Exception as e:
            log_message("ERROR", f"Error forwarding HTTP response for {self.client_address}: {e}")

    def add_custom_header(self, headers_bytes):
        headers_str = headers_bytes.decode('latin-1', errors='ignore')
        lines = headers_str.split('\r\n')
        
        header_end_idx = -1
        for i, line in enumerate(lines):
            if not line.strip():
                header_end_idx = i
                break
        
        if header_end_idx != -1:
            lines.insert(header_end_idx, f"{CUSTOM_RESPONSE_HEADER_NAME}: {CUSTOM_RESPONSE_HEADER_VALUE}")
        else:
            lines.append(f"{CUSTOM_RESPONSE_HEADER_NAME}: {CUSTOM_RESPONSE_HEADER_VALUE}")
            lines.append("")

        return '\r\n'.join(lines).encode('latin-1')

    def close_connections(self):
        if self.client_socket:
            try:
                self.client_socket.shutdown(socket.SHUT_RDWR)
                self.client_socket.close()
            except OSError as e:
                log_message("WARNING", f"Error closing client socket {self.client_address}: {e}")
            self.client_socket = None
        if self.target_socket:
            try:
                self.target_socket.shutdown(socket.SHUT_RDWR)
                self.target_socket.close()
            except OSError as e:
                log_message("WARNING", f"Error closing target socket: {e}")
            self.target_socket = None
        log_message("INFO", f"Connections closed for {self.client_address}")

def start_proxy_server():
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    
    try:
        server_socket.bind((PROXY_HOST, PROXY_PORT))
        server_socket.listen(5)
        log_message("INFO", f"Proxy server listening on {PROXY_HOST}:{PROXY_PORT}")
    except Exception as e:
        log_message("CRITICAL", f"Failed to start proxy server: {e}")
        sys.exit(1)

    while True:
        try:
            client_socket, client_address = server_socket.accept()
            log_message("INFO", f"Accepted connection from {client_address}")
            handler = ProxyHandler(client_socket, client_address)
            handler.daemon = True
            handler.start()
        except KeyboardInterrupt:
            log_message("INFO", "Proxy server shutting down.")
            break
        except Exception as e:
            log_message("ERROR", f"Error accepting client connection: {e}")

    server_socket.close()

if __name__ == "__main__":
    start_proxy_server()