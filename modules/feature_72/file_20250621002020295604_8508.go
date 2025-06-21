package main

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func main() {
	targetFile := "watched_file.txt"

	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		file, err := os.Create(targetFile)
		if err != nil {
			log.Fatalf("Error creating file %s: %v", targetFile, err)
		}
		file.Close()
	} else if err != nil {
		log.Fatalf("Error checking file %s: %v", targetFile, err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Printf("Event: %s - %s", event.Name, event.Op.String())
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	err = watcher.Add(targetFile)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}

// Additional implementation at 2025-06-21 00:21:19
import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Event struct {
	Path string
	Op   fsnotify.Op
	Time time.Time
}

type Watcher struct {
	watcher          *fsnotify.Watcher
	events           chan Event
	errors           chan error
	done             chan struct{}
	debounceDuration time.Duration
	debounceTimers   map[string]*time.Timer
	mu               sync.Mutex
	watchedDirs      map[string]bool
	ctx              context.Context
	cancel           context.CancelFunc
	wg               sync.WaitGroup
}

func NewWatcher(ctx context.Context, debounceDuration time.Duration) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fsnotify watcher: %w", err)
	}

	childCtx, cancel := context.WithCancel(ctx)

	w := &Watcher{
		watcher:          fsWatcher,
		events:           make(chan Event),
		errors:           make(chan error),
		done:             make(chan struct{}),
		debounceDuration: debounceDuration,
		debounceTimers:   make(map[string]*time.Timer),
		watchedDirs:      make(map[string]bool),
		ctx:              childCtx,
		cancel:           cancel,
	}

	w.wg.Add(1)
	go w.run()

	return w, nil
}

func (w *Watcher) Add(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Path not found, skipping: %s", path)
			return nil
		}
		return fmt.Errorf("failed to stat path %s: %w", path, err)
	}

	if info.IsDir() {
		return filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Printf("Error walking path %s: %v", p, err)
				return err
			}
			if d.IsDir() {
				w.mu.Lock()
				if _, ok := w.watchedDirs[p]; !ok {
					if err := w.watcher.Add(p); err != nil {
						log.Printf("Failed to add directory %s to watcher: %v", p, err)
						w.mu.Unlock()
						return fmt.Errorf("failed to add directory %s: %w", p, err)
					}
					w.watchedDirs[p] = true
					log.Printf("Watching directory: %s", p)
				}
				w.mu.Unlock()
			}
			return nil
		})
	} else {
		dir := filepath.Dir(path)
		w.mu.Lock()
		if _, ok := w.watchedDirs[dir]; !ok {
			if err := w.watcher.Add(dir); err != nil {
				log.Printf("Failed to add directory %s for file %s to watcher: %v", dir, path, err)
				w.mu.Unlock()
				return fmt.Errorf("failed to add directory %s for file %s: %w", dir, path, err)
			}
			w.watchedDirs[dir] = true
			log.Printf("Watching directory for file: %s (dir: %s)", path, dir)
		}
		w.mu.Unlock()
	}
	return nil
}

func (w *Watcher) Remove(path string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, ok := w.watchedDirs[path]; ok {
		if err := w.watcher.Remove(path); err != nil {
			return fmt.Errorf("failed to remove path %s from watcher: %w", path, err)
		}
		delete(w.watchedDirs, path)
		log.Printf("Stopped watching: %s", path)
	} else {
		log.Printf("Path %s was not explicitly watched, skipping remove.", path)
	}
	return nil
}

func (w *Watcher) Events() <-chan Event {
	return w.events
}

func (w *Watcher) Errors() <-chan error {
	return w.errors
}

func (w *Watcher) Close() error {
	w.cancel()
	w.wg.Wait()

	w.mu.Lock()
	for _, timer := range w.debounceTimers {
		timer.Stop()
	}
	w.debounceTimers = make(map[string]*time.Timer)
	w.mu.Unlock()

	close(w.events)
	close(w.errors)

	return w.watcher.Close()
}

func (w *Watcher) run() {
	defer w.wg.Done()
	defer close(w.done)

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			w.handleFsNotifyEvent(event)

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			w.errors <- fmt.Errorf("fsnotify error: %w", err)

		case <-w.ctx.Done():
			log.Println("Watcher context cancelled, shutting down.")
			return
		}
	}
}

func (w *Watcher) handleFsNotifyEvent(event fsnotify.Event) {
	if event.Op&fsnotify.Create == fsnotify.Create {
		info, err := os.Stat(event.Name)
		if err == nil && info.IsDir() {
			log.Printf("New directory created: %s, adding recursively.", event.Name)
			if err := w.Add(event.Name); err != nil {
				w.errors <- fmt.Errorf("failed to add new directory %s: %w", event.Name, err)
			}
		}
	}

	if event.Op&fsnotify.Remove == fsnotify.Remove {
		w.mu.Lock()
		_, wasWatchedDir := w.watchedDirs[event.Name]
		w.mu.Unlock()

		if wasWatchedDir {
			log.Printf("Watched directory removed: %s, cleaning up.", event.Name)
			w.mu.Lock()
			delete(w.watchedDirs, event.Name)
			w.mu.Unlock()
		}
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if timer, ok := w.debounceTimers[event.Name]; ok {
		timer.Stop()
	}

	w.debounceTimers[event.Name] = time.AfterFunc(w.debounceDuration, func() {
		select {
		case w.events <- Event{Path: event.Name, Op: event.Op, Time: time.Now()}:
		case <-w.ctx.Done():
			log.Printf("Context cancelled, dropping debounced event for %s", event.Name)
		}

		w.mu.Lock()
		delete(w.debounceTimers, event.Name)
		w.mu.Unlock()
	})
}

func main() {
	testDir, err := os.MkdirTemp("", "watcher_test_")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(testDir)

	log.Printf("Watching directory: %s", testDir)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watcher, err := NewWatcher(ctx, 500*time.Millisecond)
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}
	defer func() {
		log.Println("Closing watcher...")
		if err := watcher.Close(); err != nil {
			log.Printf("Error closing watcher: %v", err)
		}
		log.Println("Watcher closed.")
	}()

	if err := watcher.Add(testDir); err != nil {
		log.Fatalf("Failed to add directory %s to watcher: %v", testDir, err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events():
				log.Printf("EVENT: %s - %s", event.Path, event.Op)
			case err := <-watcher.Errors():
				log.Printf("ERROR: %v", err)
			case <-ctx.Done():
				log.Println("Event consumer shutting down.")
				return
			}
		}
	}()

	log.Println("Simulating file changes in 3 seconds...")
	time.Sleep(3 * time.Second)

	filePath1 := filepath.Join(testDir, "test_file_1.txt")
	log.Printf("Creating file: %s", filePath1)
	if err := os.WriteFile(filePath1, []byte("hello"), 0644); err != nil {
		log.Printf("Error creating file: %v", err)
	}
	time.Sleep(100 * time.Millisecond)

	log.Printf("Modifying file multiple times (should debounce): %s", filePath1)
	for i := 0; i < 5; i++ {
		if err := os.WriteFile(filePath1, []byte(fmt.Sprintf("hello %d", i)), 0644); err != nil {
			log.Printf("Error modifying file: %v", err)
		}
		time.Sleep(50 * time.Millisecond)
	}
	time.Sleep(1 * time.Second)

	subDir := filepath.Join(testDir, "sub_dir")
	log.Printf("Creating subdirectory: %s", subDir)
	if err := os.Mkdir(subDir, 0755); err != nil {
		log.Printf("Error creating sub dir: %v", err)
	}
	time.Sleep(100 * time.Millisecond)

	filePath2 := filepath.Join(subDir, "test_file_2.txt")
	log.Printf("Creating file in subdirectory: %s", filePath2)
	if err := os.WriteFile(filePath2, []byte("nested hello"), 0644);