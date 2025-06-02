package main

import "fmt"

// Logger is the interface that all loggers must implement
type Logger interface {
	Log(message string)
}

// ConsoleLogger logs messages to the console
type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(message string) {
	fmt.Printf("[Console] %s\n", message)
}

// FileLogger simulates logging to a file
type FileLogger struct {
	filename string
}

func (l *FileLogger) Log(message string) {
	// In a real implementation, this would write to a file
	fmt.Printf("[File:%s] %s\n", l.filename, message)
}

// NetworkLogger simulates logging to a network service
type NetworkLogger struct {
	host string
	port int
}

func (l *NetworkLogger) Log(message string) {
	// In a real implementation, this would send the log over the network
	fmt.Printf("[Network:%s:%d] %s\n", l.host, l.port, message)
}

// LoggerType represents the type of logger to create
type LoggerType string

const (
	Console LoggerType = "console"
	File    LoggerType = "file"
	Network LoggerType = "network"
)

// LoggerConfig holds configuration for creating loggers
type LoggerConfig struct {
	Type     LoggerType
	Filename string // Used for FileLogger
	Host     string // Used for NetworkLogger
	Port     int    // Used for NetworkLogger
}

// NewLogger is the factory function that creates the appropriate logger
func NewLogger(config LoggerConfig) (Logger, error) {
	switch config.Type {
	case Console:
		return &ConsoleLogger{}, nil
	case File:
		if config.Filename == "" {
			return nil, fmt.Errorf("filename is required for file logger")
		}
		return &FileLogger{filename: config.Filename}, nil
	case Network:
		if config.Host == "" || config.Port == 0 {
			return nil, fmt.Errorf("host and port are required for network logger")
		}
		return &NetworkLogger{host: config.Host, port: config.Port}, nil
	default:
		return nil, fmt.Errorf("unknown logger type: %s", config.Type)
	}
}

func main() {
	// Create different types of loggers using the factory
	consoleLogger, _ := NewLogger(LoggerConfig{Type: Console})
	fileLogger, _ := NewLogger(LoggerConfig{
		Type:     File,
		Filename: "app.log",
	})
	networkLogger, _ := NewLogger(LoggerConfig{
		Type: Network,
		Host: "localhost",
		Port: 514,
	})

	// Use the loggers
	consoleLogger.Log("This is a console log message")
	fileLogger.Log("This is a file log message")
	networkLogger.Log("This is a network log message")

	// Example of error handling
	_, err := NewLogger(LoggerConfig{Type: File})
	if err != nil {
		fmt.Printf("Error creating file logger: %v\n", err)
	}
}
