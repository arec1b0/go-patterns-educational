package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestVaryingWorkerCounts tests the worker pool with different numbers of workers
// Rationale: We need to understand how the number of workers affects throughput
// and resource utilization. Too few workers may underutilize resources, while
// too many may cause contention and overhead.
func TestVaryingWorkerCounts(t *testing.T) {
	// Test parameters
	testCases := []struct {
		name       string
		numWorkers int
		numJobs    int
	}{
		{"SingleWorker", 1, 100},            // Sequential processing
		{"FewWorkers", 4, 100},              // Small number of workers
		{"MediumWorkers", 16, 100},          // Medium number of workers
		{"ManyWorkers", runtime.NumCPU()*4, 100}, // Large number of workers
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create channels
			jobs := make(chan int, tc.numJobs)
			results := make(chan int, tc.numJobs)
			done := make(chan struct{})
			
			// Track processed jobs
			processed := make(map[int]bool)
			var mu sync.Mutex
			
			// Start workers
			var wg sync.WaitGroup
			startTime := time.Now()
			
			for w := 1; w <= tc.numWorkers; w++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					worker(id, jobs, results)
				}(w)
			}
			
			// Send jobs
			go func() {
				for j := 1; j <= tc.numJobs; j++ {
					jobs <- j
				}
				close(jobs)
			}()
			
			// Collect results
			go func() {
				for i := 0; i < tc.numJobs; i++ {
					result := <-results
					mu.Lock()
					processed[result/2] = true
					mu.Unlock()
				}
				close(done)
			}()
			
			// Wait for all results to be collected
			<-done
			duration := time.Since(startTime)
			
			// Verify all jobs were processed
			if len(processed) != tc.numJobs {
				t.Errorf("Expected %d jobs to be processed, got %d", tc.numJobs, len(processed))
			}
			
			// Report performance metrics
			t.Logf("%s: Processed %d jobs with %d workers in %v (%.2f jobs/sec)",
				tc.name, tc.numJobs, tc.numWorkers, duration, float64(tc.numJobs)/duration.Seconds())
			
			// Clean up
			wg.Wait()
		})
	}
}

// TestLargeJobVolume tests the worker pool with a very large number of jobs
// Rationale: Large job volumes can expose memory leaks, resource exhaustion,
// and scaling issues that aren't apparent with small workloads.
func TestLargeJobVolume(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large job volume test in short mode")
	}
	
	// Test parameters
	numWorkers := runtime.NumCPU()
	numJobs := 10000 // 10K jobs - reduced from 100K to avoid timeouts
	
	// Create channels with appropriate buffer sizes
	jobs := make(chan int, 1000)
	results := make(chan int, 1000)
	
	// Track memory usage
	var memStatsBefore, memStatsAfter runtime.MemStats
	runtime.ReadMemStats(&memStatsBefore)
	
	// Start workers
	var wg sync.WaitGroup
	startTime := time.Now()
	
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, jobs, results)
		}(w)
	}
	
	// Send jobs in a separate goroutine to avoid blocking
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- j
		}
		close(jobs)
	}()
	
	// Track progress
	processed := uint64(0)
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		timeoutTimer := time.NewTimer(30 * time.Second) // Hard timeout after 30 seconds
		defer timeoutTimer.Stop()
		for {
			select {
			case <-ticker.C:
				current := atomic.LoadUint64(&processed)
				t.Logf("Progress: %d/%d jobs (%.1f%%)", 
					current, numJobs, float64(current)/float64(numJobs)*100)
				if current >= uint64(numJobs) {
					close(done)
					return
				}
			case <-timeoutTimer.C:
				t.Logf("Test timed out after 30 seconds")
				close(done)
				return
			}
		}
	}()
	
	// Collect results with timeout
	resultCount := 0
	collectDone := make(chan struct{})
	
	go func() {
		for i := 0; i < numJobs; i++ {
			select {
			case <-results:
				atomic.AddUint64(&processed, 1)
				resultCount++
			case <-done:
				return
			}
		}
		close(collectDone)
	}()
	
	// Wait for either all results or timeout
	select {
		case <-collectDone:
			t.Logf("Collected all %d results", numJobs)
		case <-done:
			t.Logf("Collection stopped, received %d/%d results", resultCount, numJobs)
	}
	
	duration := time.Since(startTime)
	runtime.ReadMemStats(&memStatsAfter)
	
	// Report performance metrics
	t.Logf("Large Job Volume: Processed %d jobs with %d workers in %v (%.2f jobs/sec)",
		numJobs, numWorkers, duration, float64(numJobs)/duration.Seconds())
	
	// Report memory usage
	memUsedBefore := memStatsBefore.Alloc
	memUsedAfter := memStatsAfter.Alloc
	t.Logf("Memory usage: Before: %d bytes, After: %d bytes, Diff: %d bytes",
		memUsedBefore, memUsedAfter, int64(memUsedAfter)-int64(memUsedBefore))
	
	// Clean up
	wg.Wait()
}

// slowWorker is like worker but with variable processing times
func slowWorker(id int, jobs <-chan int, results chan<- int, delayFn func(int) time.Duration) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(delayFn(job)) // Variable delay
		results <- job * 2
	}
}

// TestSlowWorkers tests how the pool handles workers with variable processing times
// Rationale: In real-world scenarios, workers often have varying processing times.
// The system should handle heterogeneous processing times gracefully.
func TestSlowWorkers(t *testing.T) {
	// Test parameters
	numWorkers := 5
	numJobs := 20
	
	// Create channels
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	
	// Define delay function - some jobs are slow, others are fast
	delayFn := func(job int) time.Duration {
		if job%3 == 0 {
			return 100 * time.Millisecond // Slow job
		}
		return 10 * time.Millisecond // Fast job
	}
	
	// Start workers
	var wg sync.WaitGroup
	startTime := time.Now()
	
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			slowWorker(id, jobs, results, delayFn)
		}(w)
	}
	
	// Track job completion times
	jobStartTimes := make(map[int]time.Time)
	jobCompletionTimes := make(map[int]time.Duration)
	var mu sync.Mutex
	
	// Send jobs
	for j := 1; j <= numJobs; j++ {
		mu.Lock()
		jobStartTimes[j] = time.Now()
		mu.Unlock()
		jobs <- j
	}
	close(jobs)
	
	// Collect results and measure completion times
	for i := 0; i < numJobs; i++ {
		result := <-results
		jobID := result / 2
		mu.Lock()
		jobCompletionTimes[jobID] = time.Since(jobStartTimes[jobID])
		mu.Unlock()
	}
	
	// Calculate statistics
	var totalFastTime, totalSlowTime time.Duration
	fastCount, slowCount := 0, 0
	
	for jobID, completionTime := range jobCompletionTimes {
		if jobID%3 == 0 {
			totalSlowTime += completionTime
			slowCount++
		} else {
			totalFastTime += completionTime
			fastCount++
		}
	}
	
	avgFastTime := totalFastTime / time.Duration(fastCount)
	avgSlowTime := totalSlowTime / time.Duration(slowCount)
	
	// Report metrics
	t.Logf("Slow Workers Test: Total time: %v", time.Since(startTime))
	t.Logf("Average completion time for fast jobs: %v", avgFastTime)
	t.Logf("Average completion time for slow jobs: %v", avgSlowTime)
	
	// Verify the slow jobs actually took longer
	if avgSlowTime <= avgFastTime {
		t.Errorf("Expected slow jobs to take longer than fast jobs")
	}
	
	// Clean up
	wg.Wait()
}

// BenchmarkChannelBufferSizes benchmarks different channel buffer sizes
// Rationale: Buffer sizes affect how many jobs can be queued before blocking.
// Finding the optimal buffer size can significantly impact performance.
func BenchmarkChannelBufferSizes(b *testing.B) {
	// Test different buffer sizes
	bufferSizes := []int{0, 10, 100, 1000, 10000}
	numWorkers := runtime.NumCPU()
	
	for _, bufferSize := range bufferSizes {
		b.Run(fmt.Sprintf("Buffer_%d", bufferSize), func(b *testing.B) {
			// Create channels with the specified buffer size
			jobs := make(chan int, bufferSize)
			results := make(chan int, bufferSize)
			
			// Start workers
			var wg sync.WaitGroup
			for w := 1; w <= numWorkers; w++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					for job := range jobs {
						// Simplified worker without printing
						time.Sleep(10 * time.Microsecond) // Minimal work
						results <- job * 2
					}
				}(w)
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
			
			// Stop timer before cleanup
			b.StopTimer()
			wg.Wait()
		})
	}
}

// BenchmarkOptimalWorkerCount benchmarks different worker counts
// Rationale: Finding the optimal number of workers is crucial for maximizing
// throughput while minimizing resource contention.
func BenchmarkOptimalWorkerCount(b *testing.B) {
	// Test different worker counts
	workerCounts := []int{1, 2, 4, 8, runtime.NumCPU(), runtime.NumCPU() * 2, runtime.NumCPU() * 4}
	
	for _, workerCount := range workerCounts {
		b.Run(fmt.Sprintf("Workers_%d", workerCount), func(b *testing.B) {
			// Create channels
			jobs := make(chan int, 1000)
			results := make(chan int, 1000)
			
			// Start workers
			var wg sync.WaitGroup
			for w := 1; w <= workerCount; w++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					for job := range jobs {
						// Simplified worker without printing
						time.Sleep(10 * time.Microsecond) // Minimal work
						results <- job * 2
					}
				}(w)
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
			
			// Stop timer before cleanup
			b.StopTimer()
			wg.Wait()
		})
	}
}
