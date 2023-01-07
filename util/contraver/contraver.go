package contraver

import (
	"sync/atomic"
)

type Option struct {
	Concurrency        int
	WaitAtLeastDoneNum int
}

type OptionFunc func(*Option)

// WithConcurrency - 设置并发度
func WithConcurrency(c int) OptionFunc {
	return func(o *Option) {
		o.Concurrency = c
	}
}

// WithWaitAtLeastDoneNum - 设置完成多少个任务就退出
func WithWaitAtLeastDoneNum(n int) OptionFunc {
	return func(o *Option) {
		o.WaitAtLeastDoneNum = n
	}
}

// RunConcurrent runs function f on each element of elements concurrently.
// concurrency - 并发度
//
// example:
//
//	RunConcurrent([]int{1, 2, 3, 4, 5}, func(n int) {
//	    fmt.Println(n)
//	    time.Sleep(time.Second)
//	}, 2)
//
// in this example, RunConcurrent will run 5 tasks concurrently, and at
// most 2 tasks are running concurrently. RunConcurrent will return
// immediately after all tasks are started. after RunConcurrent returns,
// the 5 tasks will continue to run.
func RunConcurrent[T any](elements []T, f func(t T), concurrency int) {
	if concurrency <= 0 {
		concurrency = len(elements)
	}

	// 遍历每个元素并发执行函数
	sem := make(chan struct{}, concurrency)

	run := func(i int, ele T) {
		// 这个信号量需要在这里声明，而不是在外面声明
		// 这样主函数才能直接退出，而不用等待所有协程都完成
		sem <- struct{}{}
		defer func() { <-sem }()

		f(ele)
	}

	// option 2:
	// 一开始就创建所有协程，但这样可能占用过多资源
	for i := range elements {

		go run(i, elements[i])
	}

}

// TraverseAndWait traverses the elements and calls function f on each element asynchronously.
// It returns only after at least waitAtLeastDoneNum tasks are done. leftover tasks will continue to run.
// default concurrency is len(elements)
// default waitAtLeastDoneNum is len(elements)
//
// example:
//
//	TraverseAndWait([]int{1, 2, 3, 4, 5}, func(n int) {
//	    fmt.Println(n)
//	    time.Sleep(time.Second)
//	}, WithWaitAtLeastDoneNum(3), WithConcurrency(2))
//
// in this example, TraverseAndWait will return after at least 3 tasks are done,
// and at most 2 tasks are running concurrently. after TraverseAndWait returns,
// the remaining 2 tasks will continue to run.
func TraverseAndWait[T any](elements []T, f func(t T), opts ...OptionFunc) {
	options := Option{
		WaitAtLeastDoneNum: len(elements), // 默认所有都完成才返回
		Concurrency:        len(elements), // 默认所有都并发
	}
	for _, opt := range opts {
		opt(&options)
	}

	// Defensively ensure that waitNum is not less than 0
	if options.WaitAtLeastDoneNum < 0 {
		options.WaitAtLeastDoneNum = 0
	}

	// Defensively ensure that waitNum is not more than the length of elements
	if options.WaitAtLeastDoneNum > len(elements) {
		options.WaitAtLeastDoneNum = len(elements)
	}

	// Set up the done channel and counter variables
	done := make(chan struct{})
	numDone := int32(0)

	RunConcurrent(elements, func(e T) {
		f(e)
		if atomic.AddInt32(&numDone, 1) >= int32(options.WaitAtLeastDoneNum) {
			select {
			case done <- struct{}{}:
			default:
			}
		}
	}, options.Concurrency)

	// Wait until enough tasks are done
	if options.WaitAtLeastDoneNum > 0 {
		<-done
	}
}
