package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func watchCommand(args []string) error {
	fs := flag.NewFlagSet("watch", flag.ExitOnError)
	pattern := fs.String("pattern", "", "file pattern to watch (e.g., *.uad)")
	dir := fs.String("dir", ".", "directory to watch")
	debounce := fs.Duration("debounce", 500*time.Millisecond, "debounce duration")
	clear := fs.Bool("clear", true, "clear screen between runs")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Default patterns
	patterns := []string{"*.uad"}
	if *pattern != "" {
		patterns = []string{*pattern}
	}
	if fs.NArg() > 0 {
		patterns = []string{fs.Arg(0)}
	}

	watcher := &Watcher{
		Patterns: patterns,
		Directory: *dir,
		Debounce: *debounce,
		Clear:    *clear,
	}

	fmt.Printf("Watching for changes in %s (pattern: %v)\n", *dir, patterns)
	fmt.Println("Press Ctrl+C to stop...")
	fmt.Println()

	return watcher.Watch()
}

// Watcher watches files for changes and triggers re-execution
type Watcher struct {
	Patterns  []string
	Directory string
	Debounce  time.Duration
	Clear     bool
	lastRun   time.Time
}

// Watch starts watching for file changes
func (w *Watcher) Watch() error {
	// Simple polling-based file watcher (no external dependencies)
	fileStates := make(map[string]time.Time)

	// Initialize file states
	files, err := w.discoverFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		info, err := os.Stat(file)
		if err == nil {
			fileStates[file] = info.ModTime()
		}
	}

	// Initial run
	if len(files) > 0 {
		w.runFiles(files)
	}

	// Watch loop
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		// Rediscover files (in case new files were added)
		files, err := w.discoverFiles()
		if err != nil {
			continue
		}

		// Check for changes
		changed := false
		for _, file := range files {
			info, err := os.Stat(file)
			if err != nil {
				continue
			}

			modTime := info.ModTime()
			if lastMod, exists := fileStates[file]; !exists || modTime.After(lastMod) {
				fileStates[file] = modTime
				changed = true
			}
		}

		// If changed and debounce period passed, run
		if changed && time.Since(w.lastRun) > w.Debounce {
			w.runFiles(files)
		}
	}

	return nil
}

// discoverFiles finds all files matching the patterns
func (w *Watcher) discoverFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(w.Directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		basename := filepath.Base(path)
		for _, pattern := range w.Patterns {
			matched, err := filepath.Match(pattern, basename)
			if err != nil {
				return err
			}
			if matched {
				files = append(files, path)
				break
			}
		}

		return nil
	})

	return files, err
}

// runFiles executes all discovered files
func (w *Watcher) runFiles(files []string) {
	w.lastRun = time.Now()

	if w.Clear {
		clearScreen()
	}

	fmt.Printf("=== Running %d file(s) at %s ===\n\n",
		len(files), w.lastRun.Format("15:04:05"))

	// Run files sequentially in watch mode
	passed := 0
	failed := 0

	for _, file := range files {
		result := executeFile(file)
		if result.Success {
			passed++
			fmt.Printf("✓ %s (%dms)\n", filepath.Base(file), result.Duration.Milliseconds())
		} else {
			failed++
			fmt.Printf("✗ %s (%dms)\n", filepath.Base(file), result.Duration.Milliseconds())
			if result.Error != nil {
				fmt.Printf("  Error: %v\n", result.Error)
			}
		}
	}

	fmt.Printf("\n%d passed, %d failed\n\n", passed, failed)
	fmt.Println("Watching for changes...")
}

// clearScreen clears the terminal screen
func clearScreen() {
	// ANSI escape code to clear screen and move cursor to top-left
	fmt.Print("\033[2J\033[H")
}


