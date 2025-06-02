# Worker Pool Pattern (Go)

This example demonstrates a simple worker pool using goroutines and channels in Go.

## How it works
- Launches a fixed number of worker goroutines.
- Each worker pulls jobs from a channel, processes them, and sends results to another channel.
- The main function sends jobs and collects results.

## Run it
```sh
# Run the example
go run workerpool.go

# Run tests
go test -v

# Run benchmarks
go test -bench=.

# With race detector
go test -race
```

## Output Example
```
Worker 1 processing job 1
Worker 2 processing job 2
...
Result: 2
Result: 4
...
```

## Testing
The test suite includes:

1. **Unit Tests**
   - Verifies all jobs are processed
   - Ensures no job is processed multiple times
   - Tests timeout handling

2. **Benchmarks**
   - Measures worker pool performance
   - Helps determine optimal worker count

3. **Example**
   - Shows basic usage in documentation
   - Run with `go test -run=Example_workerPool`

## Pattern Use Case
- Parallel processing of independent tasks 
- Web request handling
- File processing
- Any CPU-bound or I/O-bound work that can be parallelized

## Best Practices
- Always close channels when done
- Use buffered channels for better performance
- Consider context for cancellation
- Handle worker panics gracefully
- Monitor worker health in production
