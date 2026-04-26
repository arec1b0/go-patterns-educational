package main

import (
	"strings"
	"testing"
)

func TestSummaryHandlesNoProcessedJobs(t *testing.T) {
	metrics := NewWorkerPoolMetrics()
	metrics.Complete()

	summary := metrics.Summary()

	if !strings.Contains(summary, "- Jobs Processed: 0") {
		t.Fatalf("expected summary to report zero jobs processed, got: %s", summary)
	}

	if !strings.Contains(summary, "- Throughput: 0.00 jobs/sec") {
		t.Fatalf("expected summary to report zero throughput, got: %s", summary)
	}

	if !strings.Contains(summary, "- Avg Processing Time: 0s") {
		t.Fatalf("expected summary to report zero average processing time, got: %s", summary)
	}
}
