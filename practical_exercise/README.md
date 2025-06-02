# Practical Exercise: Strategy, Factory, Worker Pool Patterns

This lesson combines multiple Go design patterns in a practical task-processing scenario.

## Patterns Demonstrated
- **Strategy Pattern:** Different processors for different actions.
- **Factory Pattern:** Processor selection based on action string.
- **Worker Pool Pattern:** Multiple workers process tasks concurrently.

## How it works
- Tasks are submitted to a worker pool.
- Each worker uses the appropriate processor to handle the task.
- Results and errors are collected and printed.

## Run it
```sh
go run practical_exercise.go
```

## Output Example
```
Worker 1 processing task 1 (uppercase)
Worker 2 processing task 2 (reverse)
...
Task 4 failed: unknown action: unknown
Task 5 failed: empty data
Task 3 result: UPPERCASE: pattern examples
...
```

## Pattern Use Case
Flexible, extensible processing of heterogeneous tasks with error handling.
