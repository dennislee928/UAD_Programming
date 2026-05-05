package cli

import (
	"sync"
)

// RunParallel executes multiple UAD files in parallel using a worker pool
func RunParallel(files []string, workers int) []*ExecutionResult {
	if workers <= 0 {
		workers = 1
	}
	if workers > len(files) {
		workers = len(files)
	}

	// Create channels
	jobs := make(chan string, len(files))
	results := make(chan *ExecutionResult, len(files))

	// Start workers
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// Send jobs
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var allResults []*ExecutionResult
	for result := range results {
		allResults = append(allResults, result)
	}

	return allResults
}

// worker is a worker function that processes files from the jobs channel
func worker(jobs <-chan string, results chan<- *ExecutionResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for file := range jobs {
		result := executeFile(file)
		results <- result
	}
}


