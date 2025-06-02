package main

import (
	"fmt"
	"sync"
)

// Config holds application configuration
type Config struct {
	APIKey      string
	DatabaseURL string
}

var (
	instance *Config
	once     sync.Once
)

// GetConfig returns the singleton instance of Config
func GetConfig() *Config {
	once.Do(func() {
		// This initialization happens exactly once
		instance = &Config{
			APIKey:      "test-api-key-12345",
			DatabaseURL: "postgres://localhost:5432/mydb",
		}
		fmt.Println("Initialized config singleton")
	})
	return instance
}

func main() {
	// Get configuration multiple times
	config1 := GetConfig()
	config2 := GetConfig()
	
	// Both variables point to the same instance
	fmt.Println("Are configs the same instance?", config1 == config2) // true
	
	// Show config values
	fmt.Printf("Config 1: %+v\n", config1)
	fmt.Printf("Config 2: %+v\n", config2)
	
	// Modifications affect all references
	config1.APIKey = "new-api-key-67890"
	fmt.Println("After modification:")
	fmt.Printf("Config 1: %+v\n", config1)
	fmt.Printf("Config 2: %+v\n", config2)
}
