package main

import (
	"log"
	"net/http"
)

func main() {
	const staticDir = "./static"

	fs := http.FileServer(http.Dir(staticDir))

	http.Handle("/", fs)

	log.Println("Serving static files from", staticDir, "on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Additional implementation at 2025-06-21 04:12:38
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type customFileServer struct {
	fs      http.Handler
	rootDir string
}

func NewCustomFileServer(dir string) *customFileServer {
	return &customFileServer{
		fs:      http.FileServer(http.Dir(dir)),
		rootDir: dir,
	}
}

func (cfs *customFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	cleanPath := filepath.Clean(r.URL.Path)
	if strings.HasPrefix(cleanPath, "/") {
		cleanPath = cleanPath[1:]
	}
	fullPath := filepath.Join(cfs.rootDir, cleanPath)

	info, err := os.Stat(fullPath)

	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "<html><body><h1>404 Not Found</h1><p>The requested resource '%s' was not found on this server.</p></body></html>", r.URL.Path)
		return
	}

	if err != nil {
		log.Printf("Error checking file %s: %v", fullPath, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if info.IsDir() && !strings.HasSuffix(r.URL.Path, "/") {
		http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
		return
	}

	cfs.fs.ServeHTTP(w, r)
}

func main() {
	if err := os.MkdirAll("static/subdir", 0755); err != nil {
		log.Fatalf("Failed to create static directory: %v", err)
	}
	if err := os.WriteFile("static/index.html", []byte("<html><body><h1>Welcome!</h1><p>This is the main index page.</p></body></html>"), 0644); err != nil {
		log.Fatalf("Failed to write index.html: %v", err)
	}
	if err := os.WriteFile("static/about.html", []byte("<html><body><h1>About Us</h1><p>We are a simple Go server.</p></body></html>"), 0644); err != nil {
		log.Fatalf("Failed to write about.html: %v", err)
	}
	if err := os.WriteFile("static/subdir/test.txt", []byte("This is a test file in a subdirectory."), 0644); err != nil {
		log.Fatalf("Failed to write test.txt: %v", err)
	}

	fileServer := NewCustomFileServer("static")

	mux := http.NewServeMux()
	mux.Handle("/", fileServer)

	log.Println("Server starting on :8080")
	log.Println("Access http://localhost:8080/")
	log.Println("Access http://localhost:8080/about.html")
	log.Println("Access http://localhost:8080/subdir/test.txt")
	log.Println("Access http://localhost:8080/nonexistent.html (should show custom 404)")
	log.Println("Access http://localhost:8080/subdir/ (should list directory contents)")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Additional implementation at 2025-06-21 04:13:43
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type staticFileServer struct {
	fs      http.Handler
	rootDir string
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (s *staticFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s %s from %s", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path, r.RemoteAddr)

	lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

	s.fs.ServeHTTP(lrw, r)

	if lrw.statusCode == http.StatusNotFound {
		s.serveCustom404(w, r)
	}
}

func (s *staticFileServer) serveCustom404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.New("404").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>404 Not Found</title>
			<style>
				body { font-family: sans-serif; text-align: center; margin-top: 50px; background-color: #f4f4f4; color: #333; }
				.container { max-width: 600px; margin: auto; padding: 20px; background: #fff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
				h1 { color: #e74c3c; margin-bottom: 10px; }
				p { font-size: 1.1em; line-height: 1.6; }
				code { background-color: #eee; padding: 2px 5px; border-radius: 3px; }
				a { color: #3498db; text-decoration: none; }
				a:hover { text-decoration: underline; }
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Oops! 404 Not Found</h1>
				<p>The requested URL <code>{{.URL}}</code> was not found on this server.</p>
				<p>It looks like you've stumbled upon a page that doesn't exist.</p>
				<p>Perhaps you can go back to the <a href="/">homepage</a> or check the address for typos.</p>
			</div>
		</body>
		</html>
	`)
	if err != nil {
		log.Printf("Error parsing 404 template: %v", err)
		fmt.Fprintf(w, "404 Not Found: %s", r.URL.Path)
		return
	}

	data := struct {
		URL string
	}{
		URL: r.URL.Path,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing 404 template: %v", err)
		fmt.Fprintf(w, "404 Not Found: %s", r.URL.Path)
	}
}

func main() {
	staticDir := "./static"
	port := ":8080"

	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Printf("Creating static directory: %s", staticDir)
		if err := os.MkdirAll(staticDir, 0755); err != nil {
			log.Fatalf("Failed to create static directory: %v", err)
		}
		dummyIndexPath := filepath.Join(staticDir, "index.html")
		if err := os.WriteFile(dummyIndexPath, []byte("<!DOCTYPE html><html><head><title>Go Static Server</title><style>body{font-family:sans-serif;text-align:center;margin-top:50px;}h1{color:#2ecc71;}</style></head><body><h1>Hello from Go Static Server!</h1><p>This is the <code>index.html</code> file.</p><p>Try navigating to <a href=\"/test.txt\">/test.txt</a> or a non-existent page like <a href=\"/nonexistent\">/nonexistent</a>.</p></body></html>"), 0644); err != nil {
			log.Printf("Failed to create dummy index.html: %v", err)
		}
		dummyFilePath := filepath.Join(staticDir, "test.txt")
		if err := os.WriteFile(dummyFilePath, []byte("This is a test file served by the Go static server.\n\nHello, world!"), 0644); err != nil {
			log.Printf("Failed to create dummy test.txt: %v", err)
		}
	}

	fileServer := http.FileServer(http.Dir(staticDir))

	customHandler := &staticFileServer{
		fs:      http.StripPrefix("/", fileServer),
		rootDir: staticDir,
	}

	http.Handle("/", customHandler)

	log.Printf("Starting static file server on port %s, serving from directory: %s", port, staticDir)
	log.Printf("Access it at http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// Additional implementation at 2025-06-21 04:14:27
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type staticFileServer struct {
	rootDir      string
	notFoundPath string
	fileServer   http.Handler
}

func NewStaticFileServer(rootDir, notFoundPath string) *staticFileServer {
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		log.Fatalf("Error getting absolute path for root directory %s: %v", rootDir, err)
	}

	fs := http.FileServer(http.Dir(absRootDir))

	return &staticFileServer{
		rootDir:      absRootDir,
		notFoundPath: notFoundPath,
		fileServer:   fs,
	}
}

func (s *staticFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)

	// Construct the full file system path.
	// filepath.Clean on r.URL.Path handles ".." and ensures a canonical path.
	// filepath.Join ensures the path is within s.rootDir.
	requestedFilePath := filepath.Join(s.rootDir, filepath.Clean(r.URL.Path))

	// Check if the file exists and is accessible.
	info, err := os.Stat(requestedFilePath)

	if os.IsNotExist(err) {
		// File does not exist.
		s.serveCustom404(w, r)
	} else if err != nil {
		// Other error (e.g., permission denied, or path is not valid).
		log.Printf("Error accessing file %s: %v", requestedFilePath, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		// File or directory exists.
		// Add Cache-Control header for static assets (optional, but good practice)
		if !info.IsDir() { // Only for files, not directories
			w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
		}

		// Serve the file using the underlying http.FileServer.
		// http.FileServer handles redirects for trailing slashes on directories,
		// serving index.html, and content-type negotiation.
		s.fileServer.ServeHTTP(w, r)
	}

	log.Printf("Served %s %s in %v", r.Method, r.URL.Path, time.Since(start))
}

func (s *staticFileServer) serveCustom404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	notFoundFullPath := filepath.Join(s.rootDir, s.notFoundPath)
	
	// Check if the custom 404 page exists
	_, err := os.Stat(notFoundFullPath)
	if os.IsNotExist(err) {
		// Fallback to a simple text 404 if custom page not found
		http.Error(w, "404 Not Found", http.StatusNotFound)
		log.Printf("Custom 404 page not found at %s. Serving default 404.", notFoundFullPath)
		return
	}

	// Serve the custom 404 page
	http.ServeFile(w, r, notFoundFullPath)
	log.Printf("Served custom 404 for %s", r.URL.Path)
}

func main() {
	const port = ":8080"
	const staticDir = "public"
	const notFoundPage = "404.html"

	// Create the static directory if it doesn't exist
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Printf("Creating static directory: %s", staticDir)
		if err := os.Mkdir(staticDir, 0755); err != nil {
			log.Fatalf("Failed to create static directory: %v", err)
		}
	}

	// Create a dummy index.html if it doesn't exist
	indexPath := filepath.Join(staticDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		log.Printf("Creating dummy index.html at %s", indexPath)
		err := os.WriteFile(indexPath, []byte("<h1>Welcome to the Go Static Server!</h1><p>Try navigating to /nonexistent.html to see the custom 404.</p>"), 0644)
		if err != nil {
			log.Fatalf("Failed to create dummy index.html: %v", err)
		}
	}

	// Create a dummy 404.html if it doesn't exist
	notFoundFullPath := filepath.Join(staticDir, notFoundPage)
	if _, err := os.Stat(notFoundFullPath); os.IsNotExist(err) {
		log.Printf("Creating dummy 404.html at %s", notFoundFullPath)
		err := os.WriteFile(notFoundFullPath, []byte("<!DOCTYPE html><html><head><title>404 Not Found</title></head><body><h1>Oops! Page Not Found</h1><p>The page you requested could not be found.</p><p><a href=\"/\">Go to Home</a></p></body></html>"), 0644)
		if err != nil {
			log.Fatalf("Failed to create dummy 404.html: %v", err)
		}
	}

	// Create the custom static file server handler
	handler := NewStaticFileServer(staticDir, notFoundPage)

	log.Printf("Starting server on %s, serving files from %s", port, staticDir)
	log.Printf("Custom 404 page: %s", notFoundPage)

	// Register the handler for all paths
	http.Handle("/", handler)

	// Start the HTTP server
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}