package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// WorkerPoolMetrics collects performance metrics for the worker pool
type WorkerPoolMetrics struct {
	StartTime      time.Time
	EndTime        time.Time
	JobsProcessed  uint64
	ErrorCount     uint64
	ProcessingTime map[int]time.Duration // Job ID -> processing time
	WorkerStats    map[int]*WorkerStat   // Worker ID -> stats
	mu             sync.Mutex
}

// WorkerStat tracks statistics for an individual worker
type WorkerStat struct {
	JobsProcessed uint64
	TotalTime     time.Duration
	MaxTime       time.Duration
	MinTime       time.Duration
	IdleTime      time.Duration
}

// NewWorkerPoolMetrics creates a new metrics collector
func NewWorkerPoolMetrics() *WorkerPoolMetrics {
	return &WorkerPoolMetrics{
		StartTime:      time.Now(),
		ProcessingTime: make(map[int]time.Duration),
		WorkerStats:    make(map[int]*WorkerStat),
	}
}

// RecordJobStart records the start of job processing
func (m *WorkerPoolMetrics) RecordJobStart(workerID, jobID int) time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Initialize worker stats if needed
	if _, exists := m.WorkerStats[workerID]; !exists {
		m.WorkerStats[workerID] = &WorkerStat{
			MinTime: time.Hour, // Initialize to a high value
		}
	}
	
	return time.Now()
}

// RecordJobCompletion records the completion of a job
func (m *WorkerPoolMetrics) RecordJobCompletion(workerID, jobID int, startTime time.Time) {
	processingTime := time.Since(startTime)
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Record job processing time
	m.ProcessingTime[jobID] = processingTime
	
	// Update worker stats
	if stat, exists := m.WorkerStats[workerID]; exists {
		stat.JobsProcessed++
		stat.TotalTime += processingTime
		
		if processingTime > stat.MaxTime {
			stat.MaxTime = processingTime
		}
		
		if processingTime < stat.MinTime {
			stat.MinTime = processingTime
		}
	}
	
	// Update total jobs processed
	atomic.AddUint64(&m.JobsProcessed, 1)
}

// RecordError increments the error count
func (m *WorkerPoolMetrics) RecordError() {
	atomic.AddUint64(&m.ErrorCount, 1)
}

// RecordWorkerIdle records idle time for a worker
func (m *WorkerPoolMetrics) RecordWorkerIdle(workerID int, idleTime time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if stat, exists := m.WorkerStats[workerID]; exists {
		stat.IdleTime += idleTime
	}
}

// Complete marks the metrics collection as complete
func (m *WorkerPoolMetrics) Complete() {
	m.EndTime = time.Now()
}

// Summary returns a summary of the metrics
func (m *WorkerPoolMetrics) Summary() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	totalDuration := m.EndTime.Sub(m.StartTime)
	throughput := float64(m.JobsProcessed) / totalDuration.Seconds()
	
	// Calculate average processing time
	var totalProcessingTime time.Duration
	for _, duration := range m.ProcessingTime {
		totalProcessingTime += duration
	}
	avgProcessingTime := totalProcessingTime / time.Duration(len(m.ProcessingTime))
	
	// Calculate worker utilization
	workerUtilization := make(map[int]float64)
	for workerID, stat := range m.WorkerStats {
		if totalDuration > 0 {
			workerUtilization[workerID] = float64(stat.TotalTime) / float64(totalDuration)
		}
	}
	
	// Build summary string
	summary := fmt.Sprintf("Worker Pool Metrics Summary:\n")
	summary += fmt.Sprintf("- Total Duration: %v\n", totalDuration)
	summary += fmt.Sprintf("- Jobs Processed: %d\n", m.JobsProcessed)
	summary += fmt.Sprintf("- Errors: %d\n", m.ErrorCount)
	summary += fmt.Sprintf("- Throughput: %.2f jobs/sec\n", throughput)
	summary += fmt.Sprintf("- Avg Processing Time: %v\n", avgProcessingTime)
	
	summary += "\nWorker Statistics:\n"
	for workerID, stat := range m.WorkerStats {
		avgTime := time.Duration(0)
		if stat.JobsProcessed > 0 {
			avgTime = stat.TotalTime / time.Duration(stat.JobsProcessed)
		}
		
		summary += fmt.Sprintf("Worker %d:\n", workerID)
		summary += fmt.Sprintf("  - Jobs Processed: %d\n", stat.JobsProcessed)
		summary += fmt.Sprintf("  - Avg Time: %v\n", avgTime)
		summary += fmt.Sprintf("  - Min Time: %v\n", stat.MinTime)
		summary += fmt.Sprintf("  - Max Time: %v\n", stat.MaxTime)
		summary += fmt.Sprintf("  - Utilization: %.2f%%\n", workerUtilization[workerID]*100)
	}
	
	return summary
}

// CSV exports the metrics in CSV format
func (m *WorkerPoolMetrics) CSV() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	csv := "job_id,processing_time_ms,worker_id\n"
	
	// Map worker IDs to jobs
	jobToWorker := make(map[int]int)
	// In a real implementation, this would be populated during job processing
	// We're using a placeholder value for this example
	for jobID := range m.ProcessingTime {
		jobToWorker[jobID] = 1 // Default to worker 1 for simplicity
	}
	
	// Generate CSV rows
	for jobID, duration := range m.ProcessingTime {
		workerID := jobToWorker[jobID]
		csv += fmt.Sprintf("%d,%d,%d\n", 
			jobID, 
			duration.Milliseconds(),
			workerID)
	}
	
	return csv
}
