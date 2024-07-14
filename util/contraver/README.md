# Contraver

Contraver is a powerful Go library for concurrent task execution, 
offering a flexible, efficient, and easy-to-use API for handling 
various concurrency scenarios.

## Features

- 🚀 **Flexible Concurrency Control**: Easily manage the number of concurrent tasks
- 🔧 **Configurable Execution Options**: Customize concurrency level and wait conditions
- 🧩 **Generic Support**: Works with any data type
- 🛡️ **Graceful Error Handling**: Built-in resource management for safe execution
- 🎯 **Simple and Intuitive API**: Easy to integrate into existing projects

## Installation

Using Go modules, simply import Contraver in your project:

```go
import "github.com/khicago/got/util/contraver"
```

Then run:

```bash
go mod tidy
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "time"
    "github.com/khicago/got/util/contraver"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5}

    contraver.RunConcurrent(numbers, func(n int) {
        fmt.Println(n)
        time.Sleep(time.Second)
    }, 2)

    // Output will be printed concurrently, with at most 2 tasks running simultaneously
}
```

### Advanced Usage with Options

```go
package main

import (
    "fmt"
    "time"
    "github.com/khicago/got/util/contraver"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5}

    contraver.TraverseAndWait(numbers, func(n int) {
        fmt.Println(n)
        time.Sleep(time.Second)
    }, contraver.WithConcurrency(2), contraver.WithWaitAtLeastDoneNum(3))

    fmt.Println("At least 3 tasks are done")
}
```

## API Reference

### RunConcurrent

```go
func RunConcurrent[T any](elements []T, f func(T), concurrency int)
```

`RunConcurrent` executes the given function `f` on each element of `elements` concurrently. The `concurrency` parameter specifies the maximum number of concurrent tasks.

#### How it works:
- Uses a semaphore channel to limit concurrency
- Launches a goroutine for each element
- Returns immediately after all tasks are started
- Tasks continue to run after `RunConcurrent` returns

### TraverseAndWait

```go
func TraverseAndWait[T any](elements []T, f func(T), opts ...OptionFunc)
```

`TraverseAndWait` traverses the elements and calls function `f` on each element asynchronously. It returns only after at least `waitAtLeastDoneNum` tasks are done. Leftover tasks will continue to run.

#### Options:
- `WithConcurrency(n int)`: Sets the maximum number of concurrent tasks
- `WithWaitAtLeastDoneNum(n int)`: Sets the minimum number of tasks to complete before returning

#### How it works:
- Uses `RunConcurrent` internally for task execution
- Employs atomic operations to track completed tasks
- Uses a channel to signal when enough tasks are done
- Supports flexible configuration through option functions

## Principles and Implementation Details

### Concurrency Control
Contraver uses a semaphore pattern to limit concurrency. A channel with a capacity equal to the desired concurrency level acts as a semaphore. Before starting a goroutine, an empty struct is sent to this channel. If the channel is full, the send operation blocks, effectively limiting the number of concurrent goroutines.

### Task Tracking
For `TraverseAndWait`, an atomic counter is used to keep track of completed tasks. This ensures thread-safe incrementing of the counter across multiple goroutines.

### Flexible Configuration
The package uses the functional options pattern to provide a clean and extensible way to configure the behavior of `TraverseAndWait`. This allows for easy addition of new options in the future without breaking existing code.

### Generic Implementation
By using Go's generics, the package can work with slices of any type, providing flexibility while maintaining type safety.

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## License
This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.
