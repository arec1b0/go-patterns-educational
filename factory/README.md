# Factory Pattern (Go)

This example demonstrates the Factory pattern in Go, which provides an interface for creating objects in a superclass but allows subclasses to alter the type of objects that will be created.

## How it works
- Defines a common `Logger` interface that all loggers must implement
- Provides concrete implementations: `ConsoleLogger`, `FileLogger`, and `NetworkLogger`
- Uses a factory function `NewLogger` to create the appropriate logger based on configuration
- Handles validation and error cases

## Run it
```sh
go run factory.go
```

## Expected Output
```
[Console] This is a console log message
[File:app.log] This is a file log message
[Network:localhost:514] This is a network log message
Error creating file logger: filename is required for file logger
```

## Pattern Use Case
- When you don't know ahead of time what class of objects you need
- When a class wants its subclasses to specify the objects it creates
- When classes delegate responsibility to one of several helper subclasses

## Key Concepts
- Encapsulates object creation
- Promotes loose coupling
- Centralizes complex creation logic
- Makes adding new types easier without changing existing code
