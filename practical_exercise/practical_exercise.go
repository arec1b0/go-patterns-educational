package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Task represents work to be done
type Task struct {
	ID     int
	Action string
	Data   string
}

// Result represents processed work
type Result struct {
	TaskID    int
	Processed string
	Error     error
}

// TaskProcessor interface for strategy pattern
type TaskProcessor interface {
	Process(Task) (string, error)
}

// ConcreteProcessors implementing different strategies
type UppercaseProcessor struct{}
func (p *UppercaseProcessor) Process(t Task) (string, error) {
	if t.Data == "" {
		return "", errors.New("empty data")
	}
	return fmt.Sprintf("UPPERCASE: %s", t.Data), nil
}

type ReverseProcessor struct{}
func (p *ReverseProcessor) Process(t Task) (string, error) {
	if t.Data == "" {
		return "", errors.New("empty data")
	}
	runes := []rune(t.Data)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

// ProcessorFactory - Factory Pattern
func GetProcessor(action string) (TaskProcessor, error) {
	switch action {
	case "uppercase":
		return &UppercaseProcessor{}, nil
	case "reverse":
		return &ReverseProcessor{}, nil
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

// TaskWorkerPool - Worker Pool Pattern
type TaskWorkerPool struct {
	wg      sync.WaitGroup
	tasks   chan Task
	results chan Result
	workers int
}

func NewTaskWorkerPool(workers int) *TaskWorkerPool {
	return &TaskWorkerPool{
		tasks:   make(chan Task),
		results: make(chan Result),
		workers: workers,
	}
}

func (p *TaskWorkerPool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			for task := range p.tasks {
				fmt.Printf("Worker %d processing task %d (%s)\n", workerID, task.ID, task.Action)
				
				processor, err := GetProcessor(task.Action)
				if err != nil {
					p.results <- Result{TaskID: task.ID, Error: err}
					continue
				}
				
				processed, err := processor.Process(task)
				p.results <- Result{
					TaskID:    task.ID,
					Processed: processed,
					Error:     err,
				}
			}
		}(i + 1)
	}
}

func (p *TaskWorkerPool) Submit(task Task) {
	p.tasks <- task
}

func (p *TaskWorkerPool) Results() <-chan Result {
	return p.results
}

func (p *TaskWorkerPool) Stop() {
	close(p.tasks)
	p.wg.Wait()
	close(p.results)
}

func main() {
	// Create and start the worker pool
	pool := NewTaskWorkerPool(3)
	pool.Start()
	
	// Submit tasks
	tasks := []Task{
		{ID: 1, Action: "uppercase", Data: "hello world"},
		{ID: 2, Action: "reverse", Data: "Go is awesome"},
		{ID: 3, Action: "uppercase", Data: "pattern examples"},
		{ID: 4, Action: "unknown", Data: "this will error"},
		{ID: 5, Action: "reverse", Data: ""},
	}
	
	go func() {
		for _, task := range tasks {
			pool.Submit(task)
		}
	}()
	
	// Wait a bit to let processing complete
	time.Sleep(100 * time.Millisecond)
	
	// Collect results (in a real app, you'd handle this differently)
	resultsCount := 0
	for resultsCount < len(tasks) {
		result := <-pool.Results()
		if result.Error != nil {
			fmt.Printf("Task %d failed: %v\n", result.TaskID, result.Error)
		} else {
			fmt.Printf("Task %d result: %s\n", result.TaskID, result.Processed)
		}
		resultsCount++
	}
	
	pool.Stop()
}
