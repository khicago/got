package contraver

import (
	"sync/atomic"
)

// Option holds configuration options for concurrent execution.
type Option struct {
	Concurrency        int // Number of concurrent tasks
	WaitAtLeastDoneNum int // Number of tasks to complete before returning
}

// OptionFunc is a function that configures an Option.
type OptionFunc func(*Option)

// WithConcurrency sets the concurrency level.
func WithConcurrency(c int) OptionFunc {
	return func(o *Option) {
		o.Concurrency = c
	}
}

// WithWaitAtLeastDoneNum sets the minimum number of tasks to complete before returning.
func WithWaitAtLeastDoneNum(n int) OptionFunc {
	return func(o *Option) {
		o.WaitAtLeastDoneNum = n
	}
}

// RunConcurrent runs function f on each element of elements concurrently.
// concurrency specifies the maximum number of concurrent tasks.
//
// Example:
//
//	RunConcurrent([]int{1, 2, 3, 4, 5}, func(n int) {
//	    fmt.Println(n)
//	    time.Sleep(time.Second)
//	}, 2)
//
// In this example, RunConcurrent will run 5 tasks concurrently, with at most 2 tasks running at the same time.
// RunConcurrent will return immediately after all tasks are started. The tasks will continue to run after RunConcurrent returns.
func RunConcurrent[T any](elements []T, f func(T), concurrency int) {
	// If concurrency is less than or equal to 0, set it to the length of elements
	if concurrency <= 0 {
		concurrency = len(elements)
	}

	// Create a semaphore channel to limit concurrency
	// The semaphore channel is used to control the number of concurrently running goroutines.
	// By sending an empty struct to the channel before starting a goroutine and receiving from
	// the channel after the goroutine completes, we ensure that no more than 'concurrency' number
	// of goroutines are running at the same time. The channel's capacity is set to 'concurrency',
	// which means it can hold up to 'concurrency' number of empty structs. If the channel is full,
	// any attempt to send to it will block until a slot becomes available, effectively limiting
	// the number of concurrent goroutines.
	sem := make(chan struct{}, concurrency)

	// Define the function to run for each element
	run := func(i int, ele T) {
		// Acquire a semaphore slot
		sem <- struct{}{}
		// Ensure the semaphore slot is released after the function completes
		defer func() { <-sem }()

		// Execute the provided function f on the element
		f(ele)
	}

	// Launch a goroutine for each element
	for i := range elements {
		go run(i, elements[i])
	}
}

// TraverseAndWait traverses the elements and calls function f on each element asynchronously.
// It returns only after at least waitAtLeastDoneNum tasks are done. Leftover tasks will continue to run.
// The default concurrency is len(elements).
// The default waitAtLeastDoneNum is len(elements).
//
// Example:
//
//	TraverseAndWait([]int{1, 2, 3, 4, 5}, func(n int) {
//	    fmt.Println(n)
//	    time.Sleep(time.Second)
//	}, WithWaitAtLeastDoneNum(3), WithConcurrency(2))
//
// In this example, TraverseAndWait will return after at least 3 tasks are done,
// with at most 2 tasks running concurrently. The remaining 2 tasks will continue to run after TraverseAndWait returns.
func TraverseAndWait[T any](elements []T, f func(T), opts ...OptionFunc) {
	// Initialize options with default values
	options := Option{
		WaitAtLeastDoneNum: len(elements), // Default to waiting for all tasks to complete
		Concurrency:        len(elements), // Default to full concurrency
	}
	for _, opt := range opts {
		opt(&options)
	}

	// Ensure waitAtLeastDoneNum is within valid range
	if options.WaitAtLeastDoneNum < 0 {
		options.WaitAtLeastDoneNum = 0
	}
	if options.WaitAtLeastDoneNum > len(elements) {
		options.WaitAtLeastDoneNum = len(elements)
	}

	// Channel to signal when enough tasks are done
	done := make(chan struct{})
	// Counter for completed tasks
	numDone := int32(0)

	// Run tasks concurrently
	RunConcurrent(elements, func(e T) {
		f(e)
		// Increment the counter and signal if enough tasks are done
		if atomic.AddInt32(&numDone, 1) >= int32(options.WaitAtLeastDoneNum) {
			select {
			case done <- struct{}{}:
			default:
			}
		}
	}, options.Concurrency)

	// Wait until the required number of tasks are done
	if options.WaitAtLeastDoneNum > 0 {
		<-done
	}
}
