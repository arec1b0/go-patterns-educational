package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// ErrUserNotFound is a sentinel error that can be checked with errors.Is
var ErrUserNotFound = errors.New("user not found")

// DatabaseError is a custom error type that wraps underlying errors
type DatabaseError struct {
	Operation string // The operation that caused the error
	Err       error  // The underlying error
}

// Error implements the error interface
func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s: %v", e.Operation, e.Err)
}

// Unwrap allows errors.Is and errors.As to work with DatabaseError
func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// getUserFromDB simulates a database user lookup
func getUserFromDB(id int) (string, error) {
	// Simulate a database query that might return sql.ErrNoRows
	if id < 0 {
		// Wrap a standard library error with our custom error type
		return "", &DatabaseError{
			Operation: "getUserFromDB",
			Err:       sql.ErrNoRows,
		}
	}
	
	// Simulate a "not found" case
	if id == 0 {
		// Use %w to wrap the sentinel error
		return "", fmt.Errorf("getting user: %w", ErrUserNotFound)
	}
	
	// Simulate a successful query
	return fmt.Sprintf("User%d", id), nil
}

// processUser demonstrates error handling patterns
func processUser(id int) error {
	username, err := getUserFromDB(id)
	if err != nil {
		return fmt.Errorf("processUser: %w", err)
	}
	
	fmt.Printf("Processing user: %s\n", username)
	// Process the user...
	
	return nil
}

func main() {
	// Example 1: Success case
	username, err := getUserFromDB(1)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Println("Found user:", username)
	
	// Example 2: User not found (wrapped sentinel error)
	_, err = getUserFromDB(0)
	if err != nil {
		// Check for specific error using errors.Is
		if errors.Is(err, ErrUserNotFound) {
			fmt.Println("Handling user not found error")
		}
		fmt.Printf("Error occurred: %v\n", err)
	}
	
	// Example 3: Database error (custom error type)
	_, err = getUserFromDB(-1)
	if err != nil {
		// Type assertion to get more details using errors.As
		var dbErr *DatabaseError
		if errors.As(err, &dbErr) {
			fmt.Printf("Database operation '%s' failed: %v\n", 
				dbErr.Operation, dbErr.Err)
		}
		fmt.Printf("Error occurred: %v\n", err)
	}
	
	// Example 4: Error wrapping in a function
	err = processUser(-1)
	if err != nil {
		fmt.Printf("Process user failed: %v\n", err)
		// The original error is still accessible
		var dbErr *DatabaseError
		if errors.As(err, &dbErr) {
			fmt.Printf("Original error: %v\n", dbErr.Err)
		}
	}
}
