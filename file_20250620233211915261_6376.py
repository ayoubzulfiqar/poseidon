import socket
import threading
import select

PROXY_HOST = '127.0.0.1'
PROXY_PORT = 8080
BUFFER_SIZE = 4096

def handle_client(client_socket):
    try:
        request = client_socket.recv(BUFFER_SIZE)
        if not request:
            return

        first_line = request.split(b'\n')[0]
        method = first_line.split(b' ')[0]

        if method == b'CONNECT':
            handle_https(client_socket, request)
        else:
            handle_http(client_socket, request)

    except Exception:
        pass
    finally:
        client_socket.close()

def handle_http(client_socket, request):
    target_socket = None
    try:
        headers = request.split(b'\r\n')
        host_header = None
        for header in headers:
            if header.lower().startswith(b'host:'):
                host_header = header.split(b':', 1)[1].strip()
                break

        if not host_header:
            client_socket.sendall(b"HTTP/1.1 400 Bad Request\r\n\r\n")
            return

        host_port = host_header.decode('utf-8').split(':')
        target_host = host_port[0]
        target_port = int(host_port[1]) if len(host_port) > 1 else 80

        target_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        target_socket.connect((target_host, target_port))

        target_socket.sendall(request)

        while True:
            response = target_socket.recv(BUFFER_SIZE)
            if not response:
                break
            client_socket.sendall(response)

    except Exception:
        try:
            client_socket.sendall(b"HTTP/1.1 500 Internal Server Error\r\n\r\n")
        except:
            pass
    finally:
        if target_socket:
            try:
                target_socket.shutdown(socket.SHUT_RDWR)
                target_socket.close()
            except:
                pass

def handle_https(client_socket, request):
    target_socket = None
    try:
        first_line = request.split(b'\n')[0]
        parts = first_line.split(b' ')
        if len(parts) < 2:
            client_socket.sendall(b"HTTP/1.1 400 Bad Request\r\n\r\n")
            return
        
        host_port_str = parts[1].decode('utf-8')
        host_port = host_port_str.split(':')
        target_host = host_port[0]
        target_port = int(host_port[1]) if len(host_port) > 1 else 443

        target_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        target_socket.connect((target_host, target_port))

        client_socket.sendall(b"HTTP/1.1 200 Connection established\r\n\r\n")

        sockets = [client_socket, target_socket]
        while True:
            readable, _, _ = select.select(sockets, [], [], 60)
            if not readable:
                break

            for sock in readable:
                if sock is client_socket:
                    data = client_socket.recv(BUFFER_SIZE)
                    if not data:
                        return
                    target_socket.sendall(data)
                elif sock is target_socket:
                    data = target_socket.recv(BUFFER_SIZE)
                    if not data:
                        return
                    client_socket.sendall(data)

    except Exception:
        try:
            

# Additional implementation at 2025-06-20 23:33:20
import http.server
import http.client
import socketserver
import urllib.parse
import sys

PROXY_PORT = 8080
TARGET_HOST = "httpbin.org"
TARGET_PORT = 80
BLOCKED_DOMAINS = ["example.com", "badsite.net"]

class ProxyHandler(http.server.BaseHTTPRequestHandler):
    protocol_version = "HTTP/1.1"

    def _send_error(self, code, message):
        self.send_response(code)
        self.send_header("Content-type", "text/html")
        self.end_headers()
        self.wfile.write(f"<html><body><h1>{code} {message}</h1></body></html>".encode('utf-8'))

    def do_CONNECT(self):
        self.send_response(200, "Connection Established")
        self.end_headers()

    def do_GET(self):
        self._handle_request("GET")

    def do_POST(self):
        self._handle_request("POST")

    def do_PUT(self):
        self._handle_request("PUT")

    def do_DELETE(self):
        self._handle_request("DELETE")

    def do_HEAD(self):
        self._handle_request("HEAD")

    def do_OPTIONS(self):
        self._handle_request("OPTIONS")

    def _handle_request(self, method):
        parsed_url = urllib.parse.urlparse(self.path)
        
        for blocked_domain in BLOCKED_DOMAINS:
            if blocked_domain in parsed_url.netloc or blocked_domain in TARGET_HOST:
                print(f"Blocked request to {self.path} (domain: {blocked_domain})")
                self._send_error(403, "Forbidden - This domain is blocked by the proxy.")
                return

        target_path = parsed_url.path
        if parsed_url.query:
            target_path += "?" + parsed_url.query

        if parsed_url.netloc:
            target_host = parsed_url.netloc.split(':')[0]
            target_port = int(parsed_url.netloc.split(':')[1]) if ':' in parsed_url.netloc else TARGET_PORT
        else:
            target_host = TARGET_HOST
            target_port = TARGET_PORT

        try:
            conn = http.client.HTTPConnection(target_host, target_port)
            
            content_length = int(self.headers.get('Content-Length', 0))
            request_body = self.rfile.read(content_length) if content_length > 0 else None

            headers_for_target = {}
            for header, value in self.headers.items():
                if header.lower() not in ['proxy-connection', 'connection', 'keep-alive', 'transfer-encoding', 'te', 'trailer', 'proxy-authorization', 'proxy-authenticate', 'content-encoding', 'content-length']:
                    headers_for_target[header] = value
            
            print(f"\n--- Incoming Request ---")
            print(f"Method: {method}")
            print(f"Path: {self.path}")
            print(f"Target: {target_host}:{target_port}{target_path}")
            print(f"Headers: {self.headers}")
            if request_body:
                print(f"Body (first 100 bytes): {request_body[:100]}...")

            conn.request(method, target_path, body=request_body, headers=headers_for_target)
            response = conn.getresponse()

            print(f"\n--- Outgoing Response ---")
            print(f"Status: {response.status} {response.reason}")
            print(f"Headers: {response.getheaders()}")

            self.send_response(response.status, response.reason)
            
            self.send_header("X-Proxy-By", "MyLocalPythonProxy")
            
            for header, value in response.getheaders():
                if header.lower() not in ['transfer-encoding', 'connection', 'content-encoding', 'content-length']:
                    self.send_header(header, value)
            
            response_body = response.read()
            self.send_header("Content-Length", str(len(response_body)))
            self.end_headers()
            self.wfile.write(response_body)

            conn.close()

        except http.client.HTTPException as e:
            print(f"HTTP Error: {e}")
            self._send_error(502, f"Bad Gateway - HTTP Error: {e}")
        except ConnectionRefusedError:
            print(f"Connection refused to {target_host}:{target_port}")
            self._send_error(502, "Bad Gateway - Connection Refused to Target")
        except Exception as e:
            print(f"An unexpected error occurred: {e}")
            self._send_error(500, f"Internal Proxy Error: {e}")

class ThreadingHTTPServer(socketserver.ThreadingMixIn, http.server.HTTPServer):
    daemon_threads = True

if __name__ == "__main__":
    print(f"Starting proxy server on port {PROXY_PORT}")
    print(f"Forwarding requests to {TARGET_HOST}:{TARGET_PORT}")
    print(f"Blocked domains: {BLOCKED_DOMAINS}")
    print(f"To use, configure your browser or application to use http://localhost:{PROXY_PORT} as a proxy.")
    print(f"Press Ctrl+C to stop the server.")

    try:
        with ThreadingHTTPServer(("", PROXY_PORT), ProxyHandler) as httpd:
            httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nProxy server stopped.")
    except Exception as e:
        print(f"Server error: {e}")
        sys.exit(1)

# Additional implementation at 2025-06-20 23:34:23


# Additional implementation at 2025-06-20 23:35:42
import socket
import threading
import sys
import re
import datetime
import select

LISTEN_HOST = '127.0.0.1'
LISTEN_PORT = 8888
BUFFER_SIZE = 4096

BLOCKED_DOMAINS = [
    'example.com',
    'badsite.net',
    'malicious.org',
    'blocked-domain.com'
]

def log_message(level, message, client_address=None, target_address=None):
    timestamp = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    log_entry = f"[{timestamp}] [{level.upper()}]"
    if client_address:
        log_entry += f" Client:{client_address[0]}:{client_address[1]}"
    if target_address:
        log_entry += f" Target:{target_address[0]}:{target_address[1]}"
    log_entry += f" {message}"
    print(log_entry)

class ProxyThread(threading.Thread):
    def __init__(self, client_socket, client_address):
        super().__init__()
        self.client_socket = client_socket
        self.client_address = client_address
        self.target_host = None
        self.target_port = None
        self.target_socket = None

    def run(self):
        try:
            first_data = self.client_socket.recv(BUFFER_SIZE)
            if not first_data:
                return

            header_str = first_data.decode('latin-1')
            
            if header_str.startswith('CONNECT'):
                match = re.search(r'CONNECT\s+([^:]+):(\d+)', header_str)
                if match:
                    self.target_host = match.group(1)
                    self.target_port = int(match.group(2))
                    log_message("info", f"HTTPS CONNECT request for {self.target_host}:{self.target_port}", self.client_address)
                    self._handle_https(first_data)
                else:
                    log_message("warning", "Malformed CONNECT request", self.client_address)
                    self.client_socket.sendall(b"HTTP/1.0 400 Bad Request\r\n\r\n")
            else:
                match = re.search(r'Host:\s*([^:\r\n]+)(?::(\d+))?', header_str, re.IGNORECASE)
                if match:
                    self.target_host = match.group(1)
                    self.target_port = int(match.group(2)) if match.group(2) else 80
                    log_message("info", f"HTTP request for {self.target_host}:{self.target_port}", self.client_address)
                    self._handle_http(first_data)
                else:
                    log_message("warning", "Malformed HTTP request (No Host header)", self.client_address)
                    self.client_socket.sendall(b"HTTP/1.0 400 Bad Request\r\n\r\n")

        except socket.timeout:
            log_message("error", "Socket timeout during initial receive", self.client_address)
        except Exception as e:
            log_message("error", f"Error in proxy thread: {e}", self.client_address)
        finally:
            if self.client_socket:
                self.client_socket.close()
            if self.target_socket:
                self.target_socket.close()

    def _handle_https(self, initial_data):
        if self.target_host in BLOCKED_DOMAINS:
            log_message("blocked", f"Blocked HTTPS connection to {self.target_host}", self.client_address)
            self.client_socket.sendall(b"HTTP/1.0 403 Forbidden\r\n\r\n")
            return

        try:
            self.target_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.target_socket.connect((self.target_host, self.target_port))
            log_message("info", f"Connected to target {self.target_host}:{self.target_port}", self.client_address, (self.target_host, self.target_port))
            
            self.client_socket.sendall(b"HTTP/1.0 200 Connection established\r\n\r\n")
            
            self._forward_data()

        except socket.gaierror:
            log_message("error", f"Could not resolve host: {self.target_host}", self.client_address)
            self.client_socket.sendall(b"HTTP/1.0 502 Bad Gateway\r\n\r\n")
        except ConnectionRefusedError:
            log_message("error", f"Connection refused by target: {self.target_host}:{self.target_port}", self.client_address)
            self.client_socket.sendall(b"HTTP/1.0 502 Bad Gateway\r\n\r\n")
        except Exception as e:
            log_message("error", f"Error handling HTTPS: {e}", self.client_address)
            self.client_socket.sendall(b"HTTP/1.0 500 Internal Proxy Error\r\n\r\n")

    def _handle_http(self, initial_data):
        if self.target_host in BLOCKED_DOMAINS:
            log_message("blocked", f"Blocked HTTP connection to {self.target_host}", self.client_address)
            self.client_socket.sendall(b"HTTP/1.0 403 Forbidden\r\n\r\n")
            return

        try:
            self.target_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.target_socket.connect((self.target_host, self.target_port))
            log_message("info", f"Connected to target {self.target_host}:{self.target_port}", self.client_address, (self.target_host, self.target_port))
            
            self.target_socket.sendall(initial_data)
            
            self._forward_data()

        except socket.gaierror:
            log_message("error", f"Could not resolve host: {self.target_host}", self.client_address)
            self.client_socket.sendall(b"HTTP/1.0 502 Bad Gateway\r\n\r\n")
        except ConnectionRefusedError:
            log_message("error", f"Connection refused by target: {self.target_host}:{self.target_port}", self.client_address)
            self.client_socket.sendall(b"HTTP/1.0 502 Bad Gateway\r\n\r\n")
        except Exception as e:
            log_message("error", f"Error handling HTTP: {e}", self.client_address)
            self.client_socket.sendall(b"HTTP/1.0 500 Internal Proxy Error\r\n\r\n")

    def _forward_data(self):
        sockets = [self.client_socket, self.target_socket]
        while True:
            ready_sockets, _, _ = select.select(sockets, [], [], 60)
            
            if not ready_sockets:
                log_message("timeout", "No data activity, closing connection.", self.client_address, (self.target_host, self.target_port))
                break

            for sock in ready_sockets:
                if sock == self.client_socket:
                    source = self.client_socket
                    destination = self.target_socket
                else:
                    source = self.target_socket
                    destination = self.client_socket

                try:
                    data = source.recv(BUFFER_SIZE)
                    if not data:
                        log_message("info", "Connection closed by peer.", self.client_address, (self.target_host, self.target_port))
                        return
                    
                    destination.sendall(data)
                except socket.error as e:
                    log_message("error", f"Socket error during data forwarding: {e}", self.client_address, (self.target_host, self.target_port))
                    return
                except Exception as e:
                    log_message("error", f"Unexpected error during data forwarding: {e}", self.client_address, (self.target_host, self.target_port))
                    return

def main():
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    
    try:
        server_socket.bind((LISTEN_HOST, LISTEN_PORT))
        server_socket.listen(5)
        log_message("info", f"Proxy server listening on {LISTEN_HOST}:{LISTEN_PORT}")

        while True:
            client_socket, client_address = server_socket.accept()
            log_message("info", f"Accepted connection from {client_address[0]}:{client_address[1]}")
            
            proxy_thread = ProxyThread(client_socket, client_address)
            proxy_thread.daemon = True
            proxy_thread.start()

    except KeyboardInterrupt:
        log_message("info", "Proxy server shutting down.")
    except Exception as e:
        log_message("critical", f"Server error: {e}")
    finally:
        if server_socket:
            server_socket.close()

if __name__ == "__main__":
    main()