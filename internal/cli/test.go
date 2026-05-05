package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func testCommand(args []string) error {
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	pattern := fs.String("pattern", "", "test file pattern (e.g., test_*.uad)")
	dir := fs.String("dir", ".", "directory to search for tests")
	exclude := fs.String("exclude", "", "exclude pattern")
	workers := fs.Int("parallel", 4, "number of parallel workers")
	verbose := fs.Bool("v", false, "verbose output")
	formatFlag := fs.String("format", "default", "output format: default, table, json, tap")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Default patterns if none specified
	patterns := []string{"test_*.uad", "*_test.uad"}
	if *pattern != "" {
		patterns = []string{*pattern}
	}
	if fs.NArg() > 0 {
		// If positional argument provided, use it as pattern
		patterns = []string{fs.Arg(0)}
	}

	// Discover tests
	runner := &TestRunner{
		Patterns:    patterns,
		Directory:   *dir,
		ExcludePattern: *exclude,
		Parallel:    *workers,
		Verbose:     *verbose || globalFlags.Verbose,
		Format:      *formatFlag,
	}

	files, err := runner.DiscoverTests()
	if err != nil {
		return fmt.Errorf("test discovery failed: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No test files found")
		return nil
	}

	if runner.Verbose {
		fmt.Printf("Discovered %d test file(s):\n", len(files))
		for _, f := range files {
			fmt.Printf("  - %s\n", f)
		}
		fmt.Println()
	}

	// Run tests
	results, err := runner.RunTests(files)
	if err != nil {
		return err
	}

	// Display results
	formatter := NewFormatter(runner.Format, globalFlags.NoColor)
	if err := formatter.FormatResults(results); err != nil {
		return err
	}

	// Exit with error if tests failed
	if results.Failed > 0 {
		os.Exit(1)
	}

	return nil
}

// TestRunner manages test discovery and execution
type TestRunner struct {
	Patterns       []string
	Directory      string
	ExcludePattern string
	Parallel       int
	Verbose        bool
	Format         string
}

// DiscoverTests finds all test files matching the patterns
func (r *TestRunner) DiscoverTests() ([]string, error) {
	var testFiles []string

	err := filepath.Walk(r.Directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file matches any pattern
		basename := filepath.Base(path)
		for _, pattern := range r.Patterns {
			matched, err := filepath.Match(pattern, basename)
			if err != nil {
				return fmt.Errorf("invalid pattern %s: %w", pattern, err)
			}

			if matched {
				// Check exclusion
				if r.ExcludePattern != "" {
					excluded, err := filepath.Match(r.ExcludePattern, basename)
					if err != nil {
						return fmt.Errorf("invalid exclude pattern %s: %w", r.ExcludePattern, err)
					}
					if excluded {
						continue
					}
				}

				testFiles = append(testFiles, path)
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return testFiles, nil
}

// RunTests executes all test files in parallel
func (r *TestRunner) RunTests(files []string) (*TestResults, error) {
	if r.Verbose {
		fmt.Printf("Running %d test(s) with %d worker(s)...\n\n", len(files), r.Parallel)
	}

	// Run in parallel
	results := RunParallel(files, r.Parallel)

	// Aggregate results
	testResults := &TestResults{
		Total:   len(results),
		Results: results,
	}

	for _, res := range results {
		if res.Success {
			testResults.Passed++
		} else {
			testResults.Failed++
		}
		testResults.Duration += res.Duration
	}

	return testResults, nil
}

// matchesAnyPattern checks if a filename matches any of the patterns
func matchesAnyPattern(filename string, patterns []string) bool {
	basename := filepath.Base(filename)
	for _, pattern := range patterns {
		// Simple glob matching
		matched, _ := filepath.Match(pattern, basename)
		if matched {
			return true
		}
		// Also check if pattern is a substring (for more flexible matching)
		if strings.Contains(basename, strings.Trim(pattern, "*")) {
			return true
		}
	}
	return false
}


