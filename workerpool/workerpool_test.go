package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestWorkerPool tests the basic functionality of the worker pool
func TestWorkerPool(t *testing.T) {
	// Test parameters
	numWorkers := 3
	numJobs := 10
	
	// Create channels
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	done := make(chan struct{})
	
	// Track processed jobs
	var wg sync.WaitGroup
	processed := make(map[int]bool)
	var mu sync.Mutex
	
	// Start workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, jobs, results)
		}(w)
	}
	
	// Send jobs
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- j
		}
		close(jobs)
	}()
	
	// Collect results in a separate goroutine
	go func() {
		for i := 0; i < numJobs; i++ {
			result := <-results
			mu.Lock()
			if processed[result/2] {
				t.Errorf("Job %d was processed multiple times", result/2)
			}
			processed[result/2] = true
			mu.Unlock()
		}
		close(done)
	}()
	
	// Wait for all results to be collected
	<-done
	
	// Verify all jobs were processed
	if len(processed) != numJobs {
		t.Errorf("Expected %d jobs to be processed, got %d", numJobs, len(processed))
	}
	
	// Verify all job numbers are present
	for i := 1; i <= numJobs; i++ {
		if !processed[i] {
			t.Errorf("Job %d was not processed", i)
		}
	}
	
	// Clean up - workers will exit when jobs channel is closed
	wg.Wait()
}

// TestWorkerPoolWithTimeout tests worker pool with timeout
func TestWorkerPoolWithTimeout(t *testing.T) {
	numWorkers := 2
	jobs := make(chan int, 5)
	results := make(chan int, 5)
	
	// Start workers
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}
	
	// Send a job that takes time
	jobs <- 42
	close(jobs)
	
	// Wait for result with timeout
	select {
	case result := <-results:
		if result != 84 { // 42 * 2
			t.Errorf("Expected result 84, got %d", result)
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Timeout waiting for worker result")
	}
}

// BenchmarkWorkerPool benchmarks the worker pool
func BenchmarkWorkerPool(b *testing.B) {
	numWorkers := 10
	jobs := make(chan int, b.N)
	results := make(chan int, b.N)
	
	// Start workers
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}
	
	// Reset timer after setup
	b.ResetTimer()
	
	// Send jobs
	for i := 0; i < b.N; i++ {
		jobs <- i
	}
	close(jobs)
	
	// Collect results
	for i := 0; i < b.N; i++ {
		<-results
	}
}

// Example_workerPool demonstrates how to use the worker pool
func Example_workerPool() {
	// Create channels
	jobs := make(chan int, 3)
	results := make(chan int, 3)
	
	// Start worker
	go worker(1, jobs, results)
	
	// Send jobs
	jobs <- 5
	jobs <- 10
	close(jobs)
	
	// Collect results - use a slice to collect all results first
	var outputs []string
	outputs = append(outputs, fmt.Sprintln(<-results))
	outputs = append(outputs, fmt.Sprint(<-results))
	
	// Print outputs in a deterministic order
	for _, output := range outputs {
		fmt.Print(output)
	}
	
	// Unordered output:
	// Worker 1 processing job 5
	// Worker 1 processing job 10
	// 10
	// 20
}
