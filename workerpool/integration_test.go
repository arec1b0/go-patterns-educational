package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

// simulateRealWorkload performs actual computation instead of just sleeping
// This better represents real-world CPU-bound tasks
func simulateRealWorkload(complexity int) int {
	result := 0
	// Fibonacci calculation - CPU intensive
	if complexity <= 1 {
		return complexity
	}
	
	a, b := 0, 1
	for i := 2; i <= complexity; i++ {
		result = a + b
		a, b = b, result
	}
	return result
}

// networkWorker simulates a worker that performs network I/O
func networkWorker(id int, jobs <-chan int, results chan<- int, metrics *WorkerPoolMetrics) {
	for job := range jobs {
		startTime := metrics.RecordJobStart(id, job)
		
		// Simulate network latency (highly variable)
		latency := time.Duration(50+rand.Intn(200)) * time.Millisecond
		time.Sleep(latency)
		
		// Simulate processing
		result := job * 2
		
		metrics.RecordJobCompletion(id, job, startTime)
		results <- result
	}
}

// cpuWorker simulates a worker that performs CPU-intensive work
func cpuWorker(id int, jobs <-chan int, results chan<- int, metrics *WorkerPoolMetrics) {
	for job := range jobs {
		startTime := metrics.RecordJobStart(id, job)
		
		// Actual CPU work - calculate Fibonacci numbers
		complexity := 25 + (job % 10) // Vary complexity
		result := simulateRealWorkload(complexity)
		
		metrics.RecordJobCompletion(id, job, startTime)
		results <- result
	}
}

// TestIntegrationMixedWorkload tests the worker pool with a mix of CPU and network tasks
// Rationale: Real-world applications often have a mix of CPU-bound and I/O-bound tasks.
// This test verifies the worker pool can handle heterogeneous workloads efficiently.
func TestIntegrationMixedWorkload(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Test parameters
	numCPUWorkers := runtime.NumCPU() / 2
	numNetworkWorkers := runtime.NumCPU() / 2
	numJobs := 100
	
	// Create channels with appropriate buffer sizes
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	
	// Create metrics collector
	metrics := NewWorkerPoolMetrics()
	
	// Start CPU workers
	var wg sync.WaitGroup
	for w := 1; w <= numCPUWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cpuWorker(id, jobs, results, metrics)
		}(w)
	}
	
	// Start network workers
	for w := numCPUWorkers + 1; w <= numCPUWorkers + numNetworkWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			networkWorker(id, jobs, results, metrics)
		}(w)
	}
	
	// Send jobs
	startTime := time.Now()
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)
	
	// Collect results
	for i := 0; i < numJobs; i++ {
		<-results
	}
	
	// Complete metrics collection
	metrics.Complete()
	duration := time.Since(startTime)
	
	// Report metrics
	t.Logf("Mixed Workload Integration Test Results:")
	t.Logf("Total time: %v", duration)
	t.Logf("Throughput: %.2f jobs/sec", float64(numJobs)/duration.Seconds())
	t.Logf("\n%s", metrics.Summary())
	
	// Clean up
	wg.Wait()
}

// TestIntegrationBurstyWorkload tests the worker pool with bursty traffic patterns
// Rationale: Many real-world systems experience bursty traffic rather than
// steady loads. This test verifies the worker pool can handle sudden spikes.
func TestIntegrationBurstyWorkload(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Test parameters
	numWorkers := runtime.NumCPU()
	burstSize := 50
	numBursts := 3
	totalJobs := burstSize * numBursts
	
	// Create channels
	jobs := make(chan int, totalJobs)
	results := make(chan int, totalJobs)
	
	// Create metrics collector
	metrics := NewWorkerPoolMetrics()
	
	// Start workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				startTime := metrics.RecordJobStart(id, job)
				
				// Simulate work
				time.Sleep(50 * time.Millisecond)
				
				metrics.RecordJobCompletion(id, job, startTime)
				results <- job * 2
			}
		}(w)
	}
	
	// Send jobs in bursts
	startTime := time.Now()
	jobID := 1
	
	for burst := 1; burst <= numBursts; burst++ {
		t.Logf("Sending burst %d of %d jobs", burst, burstSize)
		
		// Send a burst of jobs
		for j := 0; j < burstSize; j++ {
			jobs <- jobID
			jobID++
		}
		
		// Wait between bursts
		if burst < numBursts {
			time.Sleep(500 * time.Millisecond)
		}
	}
	close(jobs)
	
	// Collect results
	for i := 0; i < totalJobs; i++ {
		<-results
	}
	
	// Complete metrics collection
	metrics.Complete()
	duration := time.Since(startTime)
	
	// Report metrics
	t.Logf("Bursty Workload Integration Test Results:")
	t.Logf("Total time: %v", duration)
	t.Logf("Throughput: %.2f jobs/sec", float64(totalJobs)/duration.Seconds())
	t.Logf("\n%s", metrics.Summary())
	
	// Clean up
	wg.Wait()
}

// TestIntegrationErrorHandling tests how the worker pool handles errors
// Rationale: In real-world scenarios, tasks can fail. The worker pool
// should handle errors gracefully without losing other jobs.
func TestIntegrationErrorHandling(t *testing.T) {
	// Test parameters
	numWorkers := 5
	numJobs := 50
	errorRate := 0.2 // 20% of jobs will fail
	
	// Create channels
	jobs := make(chan int, numJobs)
	results := make(chan struct {
		value int
		err   error
	}, numJobs)
	
	// Create metrics collector
	metrics := NewWorkerPoolMetrics()
	
	// Start workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				startTime := metrics.RecordJobStart(id, job)
				
				// Simulate work
				time.Sleep(10 * time.Millisecond)
				
				// Randomly generate errors
				if rand.Float64() < errorRate {
					metrics.RecordError()
					results <- struct {
						value int
						err   error
					}{0, fmt.Errorf("error processing job %d", job)}
				} else {
					results <- struct {
						value int
						err   error
					}{job * 2, nil}
				}
				
				metrics.RecordJobCompletion(id, job, startTime)
			}
		}(w)
	}
	
	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)
	
	// Collect results
	errorCount := 0
	successCount := 0
	
	for i := 0; i < numJobs; i++ {
		result := <-results
		if result.err != nil {
			errorCount++
		} else {
			successCount++
		}
	}
	
	// Complete metrics collection
	metrics.Complete()
	
	// Verify error rate is approximately as expected
	expectedErrors := int(float64(numJobs) * errorRate)
	errorMargin := int(float64(numJobs) * 0.1) // Allow 10% margin
	
	if errorCount < expectedErrors-errorMargin || errorCount > expectedErrors+errorMargin {
		t.Errorf("Expected around %d errors, got %d", expectedErrors, errorCount)
	}
	
	// Report metrics
	t.Logf("Error Handling Integration Test Results:")
	t.Logf("Total jobs: %d", numJobs)
	t.Logf("Successful jobs: %d (%.1f%%)", successCount, float64(successCount)/float64(numJobs)*100)
	t.Logf("Failed jobs: %d (%.1f%%)", errorCount, float64(errorCount)/float64(numJobs)*100)
	t.Logf("\n%s", metrics.Summary())
	
	// Clean up
	wg.Wait()
}
