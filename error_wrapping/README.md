# Error Wrapping in Go

This example demonstrates idiomatic error handling in Go, including error wrapping, custom error types, and the use of `errors.Is` and `errors.As`.

## Key Concepts

### Error Wrapping
- Use `fmt.Errorf("context: %w", err)` to wrap errors with additional context
- Preserves the original error for inspection
- Enables error chain traversal with `errors.Is` and `errors.As`

### Sentinel Errors
- Predefined error values that can be checked with `errors.Is`
- Example: `var ErrUserNotFound = errors.New("user not found")`

### Custom Error Types
- Define types that implement the `error` interface
- Can include additional context (e.g., operation name, parameters)
- Should implement `Unwrap() error` to work with error inspection

## Run it
```sh
go run error_wrapping.go
```

## Expected Output
```
Found user: User1
Handling user not found error
Error occurred: getting user: user not found
Database operation 'getUserFromDB' failed: sql: no rows in result set
Error occurred: database error during getUserFromDB: sql: no rows in result set
Process user failed: processUser: database error during getUserFromDB: sql: no rows in result set
Original error: sql: no rows in result set
```

## Pattern Use Case
- When you need to add context to errors
- When you want to distinguish between different error types
- When you need to handle specific error conditions differently
- When building libraries where callers need to inspect errors

## Best Practices
1. Always check errors, never ignore them
2. Add context to errors when propagating them up the call stack
3. Use `errors.Is` for sentinel error checks
4. Use `errors.As` for type assertions on custom error types
5. In libraries, return errors rather than logging them
6. Keep error messages clear and helpful for debugging
