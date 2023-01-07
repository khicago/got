package contraver

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func generateElements(n int) []int {
	elements := make([]int, n)
	for i := 0; i < n; i++ {
		elements[i] = i
	}
	return elements
}

func TestRunConcurrent(t *testing.T) {
	cases := []struct {
		name        string
		concurrency int
		elemsSize   int
	}{
		// Case 1: 测试在元素数为0的情况下，函数是否能正确处理并不引发panic。
		{"Case 1", 1, 0},
		// Case 2: 单并发测试，测试在只有一个元素和一个并发时，函数能否正常工作。
		{"Case 2", 1, 1},
		// Case 3: 单并发多元素测试，测试在并发数小于元素数时，函数能否正常工作。
		{"Case 3", 1, 10},
		// Case 4: 并发数小于元素数的情况，验证是否能正确处理所有元素并且最大并发执行数不超过设置的并发度。
		{"Case 4", 5, 10},
		// Case 5: 并发数等于元素数的情况，验证是否能正确处理所有元素并且最大并发执行数不超过设置的并发度。
		{"Case 5", 10, 10},
		// Case 6: 并发数大于元素数的情况，验证函数是否仍然能正常运行并处理所有的元素。
		{"Case 6", 20, 10},
		// Case 7: 并发数为0的情况，期望函数能自动将并发度设置为与元素数相等，并且最大并发执行数不超过元素数。
		{"Case 7", 0, 10},
		// Case 8: 并发数为负数的情况，期望函数能自动将并发度设置为与元素数相等，并且最大并发执行数不超过元素数。
		{"Case 8", -5, 10},
	}

	for _, test := range cases {
		// 使用 t.Run() 运行子测试，test.name 就是当前测试用例的名字
		t.Run(test.name, func(t *testing.T) {
			// 最大并发执行数量。这个值用于验证函数 RunConcurrent 中的并发控制是否有效，
			// 即并发数是否被限制在给定的并发度以内。
			var maxConcurrentExecutions int32
			// 当前并发执行数量。这个值主要用于计算上面的最大并发执行数量。
			var currentConcurrentExecutions int32
			// 已处理元素数量。这个值用于验证函数 RunConcurrent 是否处理了所有的元素，
			// 当所有的并发任务执行完毕后，该值应该和输入数组的长度相等。
			var processedElements int32

			// 并发执行
			runFunc := func(t int) {
				// 每次执行，已处理元素数量加一
				atomic.AddInt32(&processedElements, 1)
				// 当前并发执行数量加一
				atomic.AddInt32(&currentConcurrentExecutions, 1)
				// 如果当前并发执行数量大于最大并发执行数量
				// 那么就更新最大并发执行数量
				if val := atomic.LoadInt32(&currentConcurrentExecutions); val > atomic.LoadInt32(&maxConcurrentExecutions) {
					atomic.StoreInt32(&maxConcurrentExecutions, val)
				}
				time.Sleep(100 * time.Millisecond) // 用来模拟处理时间
				// 完成后，当前并发执行数量减一
				atomic.AddInt32(&currentConcurrentExecutions, -1)
			}

			// 调用待测试的函数 RunConcurrent
			RunConcurrent(generateElements(test.elemsSize), runFunc, test.concurrency)

			time.Sleep(1100 * time.Millisecond) // 等过一段时间，让所有并发执行的函数有时间完成

			// 验证全部的元素是否被处理过
			assert.Equal(t, int32(test.elemsSize), processedElements, "All elements should be processed")
			// 如果设置的并发度小于等于0或者大于元素数量，
			// 实际并发度就会变为元素数量
			if test.concurrency <= 0 || test.concurrency > test.elemsSize {
				test.concurrency = test.elemsSize
			}
			// 验证最大并发执行数量是否符合预期
			assert.Equal(t, int32(test.concurrency), maxConcurrentExecutions, "Max concurrent executions should be equal concurrency")
		})
	}
}

func TestTraverseAndWait(t *testing.T) {
	elements := generateElements(5)
	testCases := []struct {
		waitNum     int
		concurrency int
		want        int
	}{
		// Case 1: 当等待任务数为 5，即所有任务， 并发度为 2 时，应该所有的任务都被执行，因此返回 5
		{5, 2, 5},
		// Case 2: 当等待任务数为 3， 并发度为 2 时，至少有 3个任务被执行
		{3, 2, 3},
		// Case 3: 当等待任务数为 5，即所有任务， 并发度为 2 时，应该所有的任务都被执行，因此返回 5
		{5, 2, 5},
		// Case 4: 当等待任务数为 3， 并发度为 5 时，至少有 3个任务被执行
		{3, 5, 3},
		// Case 5: 当等待任务数为 5，即所有任务， 并发度为 10 时，因为等待任务数等于所有任务数，因此返回 5
		{5, 10, 5},
	}

	for _, tc := range testCases {
		// 最大并发执行数量。这个值用于验证函数 RunConcurrent 中的并发控制是否有效，
		// 即并发数是否被限制在给定的并发度以内。
		var maxConcurrentExecutions int32
		// 当前并发执行数量。这个值主要用于计算上面的最大并发执行数量。
		var currentConcurrentExecutions int32
		// 已处理元素数量。这个值用于验证函数 RunConcurrent 是否处理了所有的元素，
		// 当所有的并发任务执行完毕后，该值应该和输入数组的长度相等。
		var processedElements int32
		f := func(i int) {
			// 模拟处理元素
			if v := atomic.AddInt32(&currentConcurrentExecutions, 1); v > atomic.LoadInt32(&maxConcurrentExecutions) {
				atomic.StoreInt32(&maxConcurrentExecutions, v)
			}
			time.Sleep(time.Duration(i) * time.Millisecond)
			atomic.AddInt32(&processedElements, 1)
			atomic.AddInt32(&currentConcurrentExecutions, -1)
		}

		opts := []OptionFunc{WithWaitAtLeastDoneNum(tc.waitNum), WithConcurrency(tc.concurrency)}
		TraverseAndWait(elements, f, opts...)

		// 验证处理完毕的元素数量
		got := atomic.LoadInt32(&processedElements)
		if got != int32(tc.want) {
			t.Errorf("for waitNum=%v and concurrency=%v, unexpected doneCount: got %v, want %v", tc.waitNum, tc.concurrency, got, tc.want)
		}
	}
}

func BenchmarkRunConcurrent(b *testing.B) {
	elements := generateElements(1000) // 生成大量元素用于测试
	runFunc := func(t int) {}          // 被并发执行的函数，因为我们主要关注内存使用和通用性能，所以这里并不做特殊处理

	// 使用 ReportAllocs 方法记录内存分配
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		RunConcurrent(elements, runFunc, 10)
	}
}
