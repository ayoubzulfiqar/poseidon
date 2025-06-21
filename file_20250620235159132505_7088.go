package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func handleProxy(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s from %s", r.Method, r.URL.String(), r.RemoteAddr)

	if r.Method == http.MethodConnect {
		handleHTTPS(w, r)
	} else {
		handleHTTP(w, r)
	}
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	if !r.URL.IsAbs() {
		http.Error(w, "This is a proxy server. Please use it as a proxy.", http.StatusBadRequest)
		return
	}

	proxyReq, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating proxy request: %v", err), http.StatusInternalServerError)
		log.Printf("Error creating proxy request for %s: %v", r.URL.String(), err)
		return
	}

	for name, values := range r.Header {
		if !strings.EqualFold(name, "Proxy-Connection") &&
			!strings.EqualFold(name, "Connection") &&
			!strings.EqualFold(name, "Keep-Alive") &&
			!strings.EqualFold(name, "Te") &&
			!strings.EqualFold(name, "Trailers") &&
			!strings.EqualFold(name, "Transfer-Encoding") &&
			!strings.EqualFold(name, "Upgrade") {
			for _, value := range values {
				proxyReq.Header.Add(name, value)
			}
		}
	}

	if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if prior := proxyReq.Header.Get("X-Forwarded-For"); prior != "" {
			clientIP = prior + ", " + clientIP
		}
		proxyReq.Header.Set("X-Forwarded-For", clientIP)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending request to target: %v", err), http.StatusBadGateway)
		log.Printf("Error sending request to target %s: %v", r.URL.String(), err)
		return
	}
	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error copying response body for %s: %v", r.URL.String(), err)
	}
	log.Printf("Proxied %s %s with status %d", r.Method, r.URL.String(), resp.StatusCode)
}

func handleHTTPS(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.DialTimeout("tcp", r.URL.Host, 10*time.Second)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error dialing target host for HTTPS: %v", err), http.StatusServiceUnavailable)
		log.Printf("Error dialing target host %s for HTTPS: %v", r.URL.Host, err)
		return
	}
	defer destConn.Close()

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		log.Print("Hijacking not supported by http.ResponseWriter")
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error hijacking connection: %v", err), http.StatusInternalServerError)
		log.Printf("Error hijacking connection: %v", err)
		return
	}
	defer clientConn.Close()

	log.Printf("Establishing HTTPS tunnel to %s", r.URL.Host)

	go func() {
		if _, err := io.Copy(destConn, clientConn); err != nil {
			log.Printf("Error copying from client to destination for %s: %v", r.URL.Host, err)
		}
		destConn.Close()
	}()
	if _, err := io.Copy(clientConn, destConn); err != nil {
		log.Printf("Error copying from destination to client for %s: %v", r.URL.Host, err)
	}
	clientConn.Close()
	log.Printf("HTTPS tunnel to %s closed", r.URL.Host)
}

func main() {
	port := "8080"
	addr := fmt.Sprintf(":%s", port)

	log.Printf("Starting proxy server on localhost%s", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(handleProxy),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not start server: %v", err)
	}
}

// Additional implementation at 2025-06-20 23:53:00
package main

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	proxyPort        = ":8080"          // Default proxy port
	readWriteTimeout = 10 * time.Second // Timeout for read/write operations on connections
	shutdownTimeout  = 5 * time.Second  // Timeout for graceful server shutdown
)

// handleHTTP processes standard HTTP requests.
func handleHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("HTTP Request: %s %s from %s", req.Method, req.URL.String(), req.RemoteAddr)

	// Create a new request to forward.
	// Use req.Context() to propagate cancellation signals (e.g., client disconnects).
	proxyReq, err := http.NewRequestWithContext(req.Context(), req.Method, req.URL.String(), req.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		log.Printf("Error creating proxy request: %v", err)
		return
	}

	// Copy headers from the original request to the proxy request.
	// Exclude hop-by-hop headers that are handled by the transport.
	for name, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// Set a User-Agent if not present, for better compatibility with some servers.
	if proxyReq.Header.Get("User-Agent") == "" {
		proxyReq.Header.Set("User-Agent", "Go-Proxy/1.0")
	}

	// Execute the proxy request using a custom HTTP client with a timeout.
	client := &http.Client{
		Timeout: readWriteTimeout,
	}
	resp, err := client.Do(proxyReq)
	if err != nil {
		// Check if the error is due to context cancellation (e.g., client disconnected).
		if req.Context().Err() != nil {
			log.Printf("Client disconnected during HTTP request to %s: %v", req.URL.String(), req.Context().Err())
			return // Do not send error response if client already gone.
		}
		http.Error(w, "Error forwarding request", http.StatusBadGateway)
		log.Printf("Error forwarding request to %s: %v", req.URL.String(), err)
		return
	}
	defer resp.Body.Close()

	log.Printf("HTTP Response: %s %s Status: %s", req.Method, req.URL.String(), resp.Status)

	// Copy headers from the proxy response to the client response.
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set the status code from the proxy response.
	w.WriteHeader(resp.StatusCode)

	// Copy the response body from the proxy response to the client response.
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v", err)
	}
}

// handleConnect processes HTTPS CONNECT requests, establishing a TCP tunnel.
func handleConnect(w http.ResponseWriter, req *http.Request) {
	log.Printf("HTTPS CONNECT Request: %s from %s", req.URL.Host, req.RemoteAddr)

	// Establish a direct TCP connection to the target host.
	targetConn, err := net.DialTimeout("tcp", req.URL.Host, readWriteTimeout)
	if err != nil {
		http.Error(w, "Error connecting to target", http.StatusServiceUnavailable)
		log.Printf("Error connecting to target %s: %v", req.URL.Host, err)
		return
	}
	defer targetConn.Close() // Ensure target connection is closed.

	// Hijack the client connection from the HTTP server to get the raw TCP connection.
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		log.Printf("Hijacking not supported for client %s", req.RemoteAddr)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, "Error hijacking client connection", http.StatusInternalServerError)
		log.Printf("Error hijacking client connection for %s: %v", req.RemoteAddr, err)
		return
	}
	defer clientConn.Close() // Ensure client connection is closed.

	// Send 200 OK to the client to indicate that the tunnel is established.
	_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		log.Printf("Error writing 200 OK to client %s: %v", req.RemoteAddr, err)
		return
	}

	// Start bidirectional data copying between client and target.
	// Use a channel to wait for both copy operations to complete or error.
	done := make(chan struct{}, 2)

	// Goroutine to copy data from client to target.
	go func() {
		_, err := io.Copy(targetConn, clientConn)
		if err != nil && err != io.EOF {
			log.Printf("Error copying from client to target for %s: %v", req.RemoteAddr, err)
		}
		done <- struct{}{}
	}()

	// Goroutine to copy data from target to client.
	go func() {
		_, err := io.Copy(clientConn, targetConn)
		if err != nil && err != io.EOF {
			log.Printf("Error copying from target to client for %s: %v", req.RemoteAddr, err)
		}
		done <- struct{}{}
	}()

	// Wait for both copy operations to finish before handleConnect returns.
	<-done
	<-done
}

func main() {
	// Create a custom HTTP server multiplexer.
	mux := http.NewServeMux()

	// Register a handler for all incoming requests.
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodConnect {
			handleConnect(w, req)
		} else {
			handleHTTP(w, req)
		}
	})

	// Configure the HTTP server.
	server := &http.Server{
		Addr:         proxyPort,
		Handler:      mux,
		ReadTimeout:  readWriteTimeout,
		WriteTimeout: readWriteTimeout,
		IdleTimeout:  30 * time.Second, // Keep-alive connections.
	}

	// Setup graceful shutdown:
	// Create a context that is cancelled when an interrupt signal is received.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop() // Release resources associated with the context.

	// Start the server in a goroutine so it doesn't block the main thread.
	go func() {
		log.Printf("Proxy server starting on %s", proxyPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v", proxyPort, err)
		}
	}()

	// Wait for the context to be cancelled (i.e., an interrupt signal is received).
	<-ctx.Done()
	log.Println("Shutting down proxy server...")

	// Create a new context with a timeout for the server shutdown operation.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt to gracefully shut down the server.
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Proxy server gracefully stopped.")
}

// Additional implementation at 2025-06-20 23:54:22
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ProxyServer represents our localhost proxy server.
type ProxyServer struct {
	listenAddr string
	client     *http.Client // http.Client for making outgoing requests.
}

// NewProxyServer creates a new ProxyServer instance configured to listen on the given address.
func NewProxyServer(addr string) *ProxyServer {
	return &ProxyServer{
		listenAddr: addr,
		client: &http.Client{
			Timeout: 30 * time.Second, // Set a timeout for outgoing HTTP requests.
		},
	}
}

// Start begins listening for incoming client connections and handles them.
func (p *ProxyServer) Start() {
	listener, err := net.Listen("tcp", p.listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", p.listenAddr, err)
	}
	log.Printf("Proxy server listening on %s", p.listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go p.handleClient(conn) // Handle each client connection in a new goroutine.
	}
}

// handleClient processes an incoming client connection, reading the request and dispatching to HTTP or HTTPS handlers.
func (p *ProxyServer) handleClient(clientConn net.Conn) {
	defer clientConn.Close()

	clientBuf := bufio.NewReader(clientConn)
	req, err := http.ReadRequest(clientBuf)
	if err != nil {
		if err != io.EOF { // EOF is expected if the client closes the connection cleanly.
			log.Printf("Error reading client request from %s: %v", clientConn.RemoteAddr(), err)
		}
		return
	}

	log.Printf("[%s] %s %s", clientConn.RemoteAddr(), req.Method, req.URL.String())

	if req.Method == http.MethodConnect {
		p.handleHTTPS(clientConn, req)
	} else {
		p.handleHTTP(clientConn, req)
	}
}

// handleHTTP proxies standard HTTP requests.
func (p *ProxyServer) handleHTTP(clientConn net.Conn, req *http.Request) {
	// Reconstruct the target URL. For HTTP proxy requests, req.URL might be relative (e.g., "/path")
	// and the actual host is in the Host header.
	targetURL := req.URL
	if !targetURL.IsAbs() {
		// Construct an absolute URL using the Host header.
		scheme := "http" // Assume HTTP for non-CONNECT requests.
		// If the proxy itself is accessed via HTTPS (e.g., behind a TLS terminator),
		// this might need to be adjusted based on X-Forwarded-Proto or similar.
		// For a simple localhost proxy, HTTP is usually sufficient.
		targetURL = &url.URL{
			Scheme:   scheme,
			Host:     req.Host,
			Path:     req.URL.Path,
			RawQuery: req.URL.RawQuery,
			Fragment: req.URL.Fragment,
		}
	}

	// Create a new request to send to the target server.
	proxyReq, err := http.NewRequest(req.Method, targetURL.String(), req.Body)
	if err != nil {
		log.Printf("Error creating proxy request for %s: %v", targetURL.String(), err)
		http.Error(clientConn, "502 Bad Gateway", http.StatusBadGateway)
		return
	}

	// Copy headers from the original request to the proxy request, excluding hop-by-hop headers.
	for name, values := range req.Header {
		if !isHopByHop(name) {
			for _, value := range values {
				proxyReq.Header.Add(name, value)
			}
		}
	}

	// Add/modify custom headers for demonstration purposes.
	// This is an example of extending functionality: adding an X-Forwarded-For header.
	if clientIP, _, err := net.SplitHostPort(clientConn.RemoteAddr().String()); err == nil {
		proxyReq.Header.Set("X-Forwarded-For", clientIP)
	}
	proxyReq.Header.Set("X-Proxy-By", "GoLocalhostProxy")
	proxyReq.Header.Set("Via", "1.1 GoLocalhostProxy")

	// Make the request to the target server.
	resp, err := p.client.Do(proxyReq)
	if err != nil {
		log.Printf("Error making request to target %s: %v", targetURL.String(), err)
		http.Error(clientConn, "502 Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Write the status line back to the client.
	statusLine := fmt.Sprintf("HTTP/%d.%d %s\r\n", resp.ProtoMajor, resp.ProtoMinor, resp.Status)
	_, err = clientConn.Write([]byte(statusLine))
	if err != nil {
		log.Printf("Error writing status line to client: %v", err)
		return
	}

	// Copy headers from the target response to the client response, excluding hop-by-hop headers.
	for name, values := range resp.Header {
		if !isHopByHop(name) {
			for _, value := range values {
				_, err = clientConn.Write([]byte(fmt.Sprintf("%s: %s\r\n", name, value)))
				if err != nil {
					log.Printf("Error writing header %s to client: %v", name, err)
					return
				}
			}
		}
	}
	_, err = clientConn.Write([]byte("\r\n")) // End of headers.
	if err != nil {
		log.Printf("Error writing end of headers to client: %v", err)
		return
	}

	// Copy the response body from the target to the client.
	_, err = io.Copy(clientConn, resp.Body)
	if err != nil {
		log.Printf("Error copying response body to client: %v", err)
	}
}

// handleHTTPS handles CONNECT requests for establishing an SSL/TLS tunnel.
func (p *ProxyServer) handleHTTPS(clientConn net.Conn, req *http.Request) {
	// The target address for CONNECT requests is in req.URL.Host (e.g., "www.google.com:443").
	targetAddr := req.URL.Host
	if !strings.Contains(targetAddr, ":") {
		targetAddr += ":443" // Default to HTTPS port 443 if not specified.
	}

	log.Printf("Establishing HTTPS tunnel to %s for %s", targetAddr, clientConn.RemoteAddr())

	// Establish a direct TCP connection to the target server.
	targetConn, err := net.DialTimeout("tcp", targetAddr, 10*time.Second)
	if err != nil {
		log.Printf("Error dialing target %s: %v", targetAddr, err)
		http.Error(clientConn, "502 Bad Gateway", http.StatusBadGateway)
		return
	}
	defer targetConn.Close()

	// Inform the client that the tunnel is established.
	_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	if err != nil {
		log.Printf("Error writing 200 OK to client for CONNECT: %v", err)
		return
	}

	// Start bi-directional copying between client and target connections.
	// This allows the client and target to communicate directly over the established tunnel.
	done := make(chan struct{}, 2) // Use a channel to wait for both copy operations to complete.

	go func() {
		_, err := io.Copy(targetConn, clientConn)
		if err != nil && err != io.EOF { // Ignore EOF, which means one side closed the connection.
			log.Printf("Error copying from client to target for %s: %v", targetAddr, err)
		}
		done <- struct{}{}
	}()

	go func() {
		_, err := io.Copy(clientConn, targetConn)
		if err != nil && err != io.EOF {
			log.Printf("Error copying from target to client for %s: %v", targetAddr, err)
		}
		done <- struct{}{}
	}()

	// Wait for both copy operations to finish before closing connections.
	<-done
	<-done
}

// isHopByHop checks if a header is a hop-by-hop header that should not be forwarded by a proxy.
func isHopByHop(header string) bool {
	header = strings.ToLower(header)
	switch header {
	case "connection", "keep-alive", "proxy-authenticate", "proxy-authorization",
		"te", "trailer", "transfer-encoding", "upgrade":
		return true
	default:
		return false
	}
}

func main() {
	proxy := NewProxyServer(":8080") // The proxy will listen on localhost port 8080.
	proxy.Start()
}

// Additional implementation at 2025-06-20 23:55:02
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	proxyPort = ":8080"
	proxyName = "GoProxy/1.0"
)

func main() {
	log.Printf("Starting %s on port %s", proxyName, proxyPort)
	listener, err := net.Listen("tcp", proxyPort)
	if err != nil {
		log.Fatalf("Failed to start proxy: %v", err)
	}
	defer listener.Close()

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept client connection: %v", err)
			continue
		}
		go handleClient(clientConn)
	}
}

func handleClient(clientConn net.Conn) {
	defer clientConn.Close()

	clientReader := bufio.NewReader(clientConn)

	// Read the first line of the request to determine method and target
	firstLine, err := clientReader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading first line from client %s: %v", clientConn.RemoteAddr(), err)
		return
	}

	parts := strings.Split(strings.TrimSpace(firstLine), " ")
	if len(parts) < 3 {
		log.Printf("Invalid HTTP request line from client %s: %s", clientConn.RemoteAddr(), firstLine)
		return
	}

	method := parts[0]
	requestURL := parts[1]
	httpVersion := parts[2]

	log.Printf("[%s] %s %s", clientConn.RemoteAddr(), method, requestURL)

	if method == http.MethodConnect {
		handleHTTPS(clientConn, clientReader, requestURL)
	} else {
		handleHTTP(clientConn, clientReader, method, requestURL, httpVersion)
	}
}

func handleHTTPS(clientConn net.Conn, clientReader *bufio.Reader, requestURL string) {
	// For CONNECT, the requestURL is host:port
	targetAddr := requestURL

	// Discard remaining headers for CONNECT request
	for {
		line, err := clientReader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading CONNECT headers from client %s: %v", clientConn.RemoteAddr(), err)
			return
		}
		if strings.TrimSpace(line) == "" {
			break // End of headers
		}
	}

	targetConn, err := net.DialTimeout("tcp", targetAddr, 5*time.Second)
	if err != nil {
		log.Printf("Failed to connect to target %s for HTTPS: %v", targetAddr, err)
		clientConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
		return
	}
	defer targetConn.Close()

	// Inform client that connection is established
	clientConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))

	// Bidirectional data transfer
	go io.Copy(targetConn, clientConn)
	io.Copy(clientConn, targetConn)
}

func handleHTTP(clientConn net.Conn, clientReader *bufio.Reader, method, requestURL, httpVersion string) {
	// Reconstruct the first line
	fullRequest := []byte(fmt.Sprintf("%s %s %s\r\n", method, requestURL, httpVersion))

	// Read remaining headers from client
	var hostHeader string
	headers := make(map[string]string)
	for {
		line, err := clientReader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading HTTP headers from client %s: %v", clientConn.RemoteAddr(), err)
			return
		}
		fullRequest = append(fullRequest, []byte(line)...)
		if strings.TrimSpace(line) == "" {
			break // End of headers
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[strings.ToLower(key)] = value
			if strings.ToLower(key) == "host" {
				hostHeader = value
			}
		}
	}

	// Determine target address
	var targetAddr string
	if hostHeader != "" {
		targetAddr = hostHeader
		if !strings.Contains(targetAddr, ":") {
			targetAddr += ":80" // Default HTTP port
		}
	} else {
		// Fallback if Host header is missing, try to parse from requestURL
		parsedURL, err := url.Parse(requestURL)
		if err == nil && parsedURL.Host != "" {
			targetAddr = parsedURL.Host
			if !strings.Contains(targetAddr, ":") {
				targetAddr += ":80"
			}
		} else {
			log.Printf("Could not determine target host for %s %s", method, requestURL)
			clientConn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
			return
		}
	}

	targetConn, err := net.DialTimeout("tcp", targetAddr, 5*time.Second)
	if err != nil {
		log.Printf("Failed to connect to target %s for HTTP: %v", targetAddr, err)
		clientConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
		return
	}
	defer targetConn.Close()

	// Modify request: Add X-Proxy-By header
	// Replace the first CRLF after the request line with the new header and CRLF
	modifiedRequest := strings.Replace(string(fullRequest), "\r\n", fmt.Sprintf("\r\nX-Proxy-By: %s\r\n", proxyName), 1)
	// Some clients send Proxy-Connection, which should be changed to Connection for the upstream server
	modifiedRequest = strings.Replace(modifiedRequest, "Proxy-Connection:", "Connection:", 1)

	// Write modified request to target
	_, err = targetConn.Write([]byte(modifiedRequest))
	if err != nil {
		log.Printf("Error writing request to target %s: %v", targetAddr, err)
		return
	}

	// Read response from target and forward to client
	targetReader := bufio.NewReader(targetConn)
	responseFirstLine, err := targetReader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading response first line from target %s: %v", targetAddr, err)
		return
	}

	// Log response status
	responseParts := strings.SplitN(strings.TrimSpace(responseFirstLine), " ", 3)
	statusCode := "N/A"
	if len(responseParts) >= 2 {
		statusCode = responseParts[1]
	}
	log.Printf("[%s] %s %s -> %s %s", clientConn.RemoteAddr(), method, requestURL, statusCode, targetAddr)

	// Write response first line to client
	_, err = clientConn.Write([]byte(responseFirstLine))
	if err != nil {
		log.Printf("Error writing response first line to client %s: %v", clientConn.RemoteAddr(), err)
		return
	}

	// Copy remaining response headers and body
	_, err = io.Copy(clientConn, targetReader)
	if err != nil {
		log.Printf("Error copying response body from target to client: %v", err)
	}
}