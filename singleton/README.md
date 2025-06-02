# Singleton Pattern (Go)

This example demonstrates the thread-safe singleton pattern in Go using `sync.Once`.

## How it works
- Uses `sync.Once` to ensure thread-safe initialization
- Lazy initialization (only created when first needed)
- Global point of access to the instance
- Prevents multiple instances in concurrent scenarios

## Run it
```sh
go run singleton.go
```

## Expected Output
```
Initialized config singleton
Are configs the same instance? true
Config 1: &{APIKey:test-api-key-12345 DatabaseURL:postgres://localhost:5432/mydb}
Config 2: &{APIKey:test-api-key-12345 DatabaseURL:postgres://localhost:5432/mydb}
After modification:
Config 1: &{APIKey:new-api-key-67890 DatabaseURL:postgres://localhost:5432/mydb}
Config 2: &{APIKey:new-api-key-67890 DatabaseURL:postgres://localhost:5432/mydb}
```

## Pattern Use Case
- Configuration management
- Database connection pools
- Logging systems
- Caching mechanisms

## Key Concepts
- Thread-safe initialization
- Global access point
- Single instance guarantee
- Lazy initialization
