# Pipeline Pattern (Go)

This example demonstrates the pipeline pattern using Go channels and goroutines for staged data processing.

## How it works
- **Stage 1 (Generate)**: Produces a stream of numbers
- **Stage 2 (Square)**: Squares each number
- **Stage 3 (Sum)**: Sums all squared numbers

Each stage runs concurrently, connected by channels, allowing for efficient data flow and processing.

## Run it
```sh
go run pipeline.go
```

## Expected Output
```
Sum of squares: 55
```

## Pattern Use Case
- Data processing pipelines
- Stream processing
- Multi-stage computations

## Key Concepts
- Channels for communication between stages
- Goroutines for concurrent execution
- Deferred channel closing
- Type-safe pipeline stages
