package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
)

const (
	shortCodeLength = 6
	storageFile     = "urls.json"
	port            = ":8080"
)

type URLStore struct {
	urls     map[string]string
	mu       sync.RWMutex
	filename string
}

func NewURLStore(filename string) *URLStore {
	store := &URLStore{
		urls:     make(map[string]string),
		filename: filename,
	}
	if err := store.Load(); err != nil {
		log.Printf("Failed to load URLs from %s: %v. Starting with empty store.", filename, err)
	}
	return store
}

func (s *URLStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := ioutil.ReadFile(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error reading file %s: %w", s.filename, err)
	}

	if len(data) == 0 {
		s.urls = make(map[string]string)
		return nil
	}

	if err := json.Unmarshal(data, &s.urls); err != nil {
		return fmt.Errorf("error unmarshalling JSON from %s: %w", s.filename, err)
	}
	return nil
}

func (s *URLStore) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.MarshalIndent(s.urls, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling URLs to JSON: %w", err)
	}

	if err := ioutil.WriteFile(s.filename, data, 0644); err != nil {
		return fmt.Errorf("error writing to file %s: %w", s.filename, err)
	}
	return nil
}

func (s *URLStore) Get(shortCode string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	longURL, ok := s.urls[shortCode]
	return longURL, ok
}

func (s *URLStore) Add(longURL string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var shortCode string
	for {
		shortCode = generateShortCode(shortCodeLength)
		if _, exists := s.urls[shortCode]; !exists {
			s.urls[shortCode] = longURL
			break
		}
	}

	if err := s.Save(); err != nil {
		log.Printf("Warning: Failed to save URLs after adding new entry: %v", err)
	}
	return shortCode, nil
}

func generateShortCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Error generating random bytes: %v", err)
	}

	for i := 0; i < length; i++ {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

func shortenHandler(store *URLStore, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	shortCode, err := store.Add(longURL)
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		log.Printf("Error adding URL: %v", err)
		return
	}

	shortURL := fmt.Sprintf("http://localhost%s/%s", port, shortCode)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s", shortURL)
}

func redirectHandler(store *URLStore, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortCode := path.Base(r.URL.Path)
	if shortCode == "" || shortCode == "/" {
		http.Error(w, "Welcome to the URL Shortener! Use POST /shorten with form data 'url=<your_long_url>' to shorten a URL.", http.StatusOK)
		return
	}

	longURL, ok := store.Get(shortCode)
	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func main() {
	store := NewURLStore(storageFile)

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		shortenHandler(store, w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		redirectHandler(store, w, r)
	})

	log.Printf("URL Shortener service starting on port %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Additional implementation at 2025-06-18 00:51:52
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

const (
	storageFileName = "urls.json"
	base62Charset   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	serverPort      = ":8080"
)

type Storage struct {
	URLs   map[string]string `json:"urls"`
	NextID int64             `json:"next_id"`
}

var (
	data     Storage
	mu       sync.RWMutex
	filePath string
)

func init() {
	data = Storage{
		URLs:   make(map[string]string),
		NextID: 1,
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	filePath = filepath.Join(dir, storageFileName)

	loadURLs()
}

func loadURLs() {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Open(filePath)
	if os.IsNotExist(err) {
		log.Printf("Storage file %s not found, starting with empty data.", filePath)
		return
	}
	if err != nil {
		log.Fatalf("Failed to open storage file: %v", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read storage file: %v", err)
	}

	if len(bytes) == 0 {
		log.Printf("Storage file %s is empty, starting with empty data.", filePath)
		return
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		log.Fatalf("Failed to unmarshal URL data from file: %v", err)
	}
	log.Printf("Loaded %d URLs from %s. Next ID: %d", len(data.URLs), filePath, data.NextID)
}

func saveURLs() error {
	mu.RLock()
	defer mu.RUnlock()

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal URL data: %w", err)
	}

	tmpFilePath := filePath + ".tmp"
	if err := ioutil.WriteFile(tmpFilePath, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write temporary storage file: %w", err)
	}

	if err := os.Rename(tmpFilePath, filePath); err != nil {
		return fmt.Errorf("failed to rename temporary file to storage file: %w", err)
	}

	log.Printf("Saved %d URLs to %s. Next ID: %d", len(data.URLs), filePath, data.NextID)
	return nil
}

func generateShortID(id int64) string {
	if id == 0 {
		return string(base62Charset[0])
	}

	shortURL := []byte{}
	for id > 0 {
		remainder := id % int64(len(base62Charset))
		shortURL = append([]byte{base62Charset[remainder]}, shortURL...)
		id /= int64(len(base62Charset))
	}
	return string(shortURL)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	if !isValidURL(longURL) {
		http.Error(w, "Invalid URL format. Must start with http:// or https://", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for short, long := range data.URLs {
		if long == longURL {
			fmt.Fprintf(w, "http://localhost%s/%s\n", serverPort, short)
			return
		}
	}

	shortID := generateShortID(data.NextID)
	data.URLs[shortID] = longURL
	data.NextID++

	if err := saveURLs(); err != nil {
		log.Printf("Error saving URLs: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "http://localhost%s/%s\n", serverPort, shortID)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortID := r.URL.Path[1:]

	if shortID == "" {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Welcome to the Go URL Shortener!\n\n")
		fmt.Fprintf(w, "To shorten a URL, send a POST request to http://localhost%s/shorten with 'url' form parameter.\n", serverPort)
		fmt.Fprintf(w, "Example: curl -X POST -d \"url=https://example.com\" http://localhost%s/shorten\n", serverPort)
		return
	}

	mu.RLock()
	longURL, found := data.URLs[shortID]
	mu.RUnlock()

	if !found {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func isValidURL(url string) bool {
	return len(url) > 7 && (url[:7] == "http://" || url[:8] == "https://")
}

func main() {
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)

	log.Printf("URL Shortener service starting on port %s", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}

// Additional implementation at 2025-06-18 00:53:06
package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

const (
	defaultShortCodeLength = 7
	storageFileName        = "urls.json"
)

// URLMapping represents a single URL mapping entry.
type URLMapping struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// URLShortenerService manages URL mappings and persistence.
type URLShortenerService struct {
	mu       sync.RWMutex
	mappings map[string]URLMapping // shortCode -> URLMapping
	filePath string
}

// NewURLShortenerService creates and initializes a new URLShortenerService.
func NewURLShortenerService(dataDir string) (*URLShortenerService, error) {
	service := &URLShortenerService{
		mappings: make(map[string]URLMapping),
		filePath: filepath.Join(dataDir, storageFileName),
	}

	if err := service.loadMappings(); err != nil {
		return nil, fmt.Errorf("failed to load mappings: %w", err)
	}
	return service, nil
}

// loadMappings loads URL mappings from the storage file.
func (s *URLShortenerService) loadMappings() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Storage file %s not found, starting with empty mappings.", s.filePath)
			return nil
		}
		return fmt.Errorf("error reading storage file: %w", err)
	}

	var loadedMappings []URLMapping
	if len(data) > 0 {
		if err := json.Unmarshal(data, &loadedMappings); err != nil {
			return fmt.Errorf("error unmarshalling mappings: %w", err)
		}
	}

	for _, mapping := range loadedMappings {
		s.mappings[mapping.ShortCode] = mapping
	}
	log.Printf("Loaded %d mappings from %s", len(s.mappings), s.filePath)
	return nil
}

// saveMappings saves current URL mappings to the storage file.
func (s *URLShortenerService) saveMappings() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var mappingsToSave []URLMapping
	for _, mapping := range s.mappings {
		mappingsToSave = append(mappingsToSave, mapping)
	}

	data, err := json.MarshalIndent(mappingsToSave, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling mappings: %w", err)
	}

	if err := ioutil.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("error writing storage file: %w", err)
	}
	return nil
}

// generateShortCode generates a unique short code.
func (s *URLShortenerService) generateShortCode() (string, error) {
	for {
		// Calculate bytes needed to guarantee defaultShortCodeLength base64 characters
		// Base64 encodes 3 bytes into 4 characters. So

// Additional implementation at 2025-06-18 00:54:02
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	defaultShortURLLength = 6
	charset               = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	storageFileName       = "urls.json"
)

// URLShortener manages the mapping of short URLs to long URLs and persists them to a file.
type URLShortener struct {
	mu       sync.RWMutex
	urlMap   map[string]string
	filePath string
}

// NewURLShortener creates a new URLShortener instance and loads existing mappings from a file.
func NewURLShortener(dir string) (*URLShortener, error) {
	us := &URLShortener{
		urlMap:   make(map[string]string),
		filePath: filepath.Join(dir, storageFileName),
	}

	if err := us.load(); err != nil {
		// If file doesn't exist, it's not an error, just start with an empty map
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load URL mappings: %w", err)
		}
	}
	return us, nil
}

// load reads the URL mappings from the storage file.
func (us *URLShortener) load() error {
	us.mu.Lock()
	defer us.mu.Unlock()

	data, err := ioutil.ReadFile(us.filePath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		// File is empty, initialize with an empty map
		us.urlMap = make(map[string]string)
		return nil
	}

	return json.Unmarshal(data, &us.urlMap)
}

// save writes the current URL mappings to the storage file.
func (us *URLShortener) save() error {
	us.mu.RLock()
	defer us.mu.RUnlock()

	data, err := json.MarshalIndent(us.urlMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal URL mappings: %w", err)
	}

	// Write to a temporary file first to ensure atomicity
	tmpFilePath := us.filePath + ".tmp"
	if err := ioutil.WriteFile(tmpFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Rename the temporary file to the actual file
	if err := os.Rename(tmpFilePath, us.filePath); err != nil {
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}
	return nil
}

// generateShortCode generates a random short code of a specified length.
func generateShortCode(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// ShortenURL creates a short URL for a given long URL.
// If customShortURL is provided, it attempts to use that; otherwise, it generates a random one.
// Returns the short URL or an error if the custom URL is taken or generation fails.
func (us *URLShortener) ShortenURL(longURL, customShortURL string) (string, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	if customShortURL != "" {
		if _, exists := us.urlMap[customShortURL]; exists {
			return "", fmt.Errorf("custom short URL '%s' is already taken", customShortURL)
		}
		us.urlMap[customShortURL] = longURL
		go us.save() // Save asynchronously
		return customShortURL, nil
	}

	// Generate a random short URL
	for i := 0; i < 10; i++ { // Try a few times to avoid collisions
		shortCode := generateShortCode(defaultShortURLLength)
		if _, exists := us.urlMap[shortCode]; !exists {
			us.urlMap[shortCode] = longURL
			go us.save() // Save asynchronously
			return shortCode, nil
		}
	}
	return "", fmt.Errorf("failed to generate a unique short URL after multiple attempts")
}

// GetLongURL retrieves the long URL associated with a short URL.
func (us *URLShortener) GetLongURL(shortURL string) (string, bool) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	longURL, ok := us.urlMap[shortURL]
	return longURL, ok
}

// ShortenRequest represents the JSON request body for shortening a URL.
type ShortenRequest struct {
	LongURL        string `json:"long_url"`
	CustomShortURL string `json:"custom_short_url,omitempty"`
}

// ShortenResponse represents the JSON response body for a successful URL shortening.
type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

// ErrorResponse represents the JSON response body for an error.
type ErrorResponse struct {
	Error string `json:"error"`
}

// handleShorten is the HTTP handler for creating short URLs.
// It expects a POST request with a JSON body containing "long_url" and optionally "custom_short_url".
func (us *URLShortener) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.LongURL == "" {
		http.Error(w, "Long URL is required", http.StatusBadRequest)
		return
	}

	shortURL, err := us.ShortenURL(req.LongURL, req.CustomShortURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to shorten URL: %v", err), http.StatusInternalServerError)
		return
	}

	resp := ShortenResponse{ShortURL: shortURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleRedirect is the HTTP handler for redirecting short URLs to their long counterparts.
// It expects a GET request with the short code as part of the path.
func (us *URLShortener) handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[1:] // Remove leading '/'

	if shortCode == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	longURL, ok := us.GetLongURL(shortCode)
	if !ok {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound) // Use StatusFound (302) for temporary redirect
}

func main() {
	// Create a directory for storage if it doesn't exist
	storageDir := "./data"
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	shortener, err := NewURLShortener(storageDir)
	if err != nil {
		log.Fatalf("Failed to initialize URL shortener: %v", err)
	}

	// Register handlers
	http.HandleFunc("/shorten", shortener.handleShorten)
	http.HandleFunc("/", shortener.handleRedirect) // Catch-all for short URLs

	port := ":8080"
	log.Printf("URL Shortener service starting on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}