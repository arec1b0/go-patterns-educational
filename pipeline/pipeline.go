package main

import (
	"fmt"
)

// generate produces numbers and sends them to the output channel
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// square receives numbers, squares them, and sends to output channel
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// sum receives numbers and returns their sum
func sum(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		total := 0
		for n := range in {
			total += n
		}
		out <- total
	}()
	return out
}

func main() {
	// Set up the pipeline
	c1 := generate(1, 2, 3, 4, 5)
	c2 := square(c1)
	c3 := sum(c2)
	
	// Get the final result
	fmt.Println("Sum of squares:", <-c3)
}
