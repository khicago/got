# ProRetry

ProRetry 提供了多种退避算法和可配置的重试机制。

## 特性

- 多种退避算法：常数、线性、指数和斐波那契
- 灵活的重试机制
- 可配置的初始重试间隔和退避策略
- 支持自定义错误处理和可重试错误判定
- 清晰的错误处理，包括重试次数和最后一次错误 

## 使用方法

### 基本用法

```go
import (
    "github.com/khicago/got/util/proretry"
)

func main() {
    err := proretry.Run(
        func() error {
            // 你的函数逻辑
            return nil
        },
        3, // 最大重试次数
    )

    if err != nil {
        // 处理错误
    }
}
```

### 自定义退避策略

```go
err := proretry.Run(
    yourFunction,
    5,
    proretry.WithBackoff(proretry.ExponentialBackoff(100 * time.Millisecond)),
)
```

### 自定义可重试错误

```go
err := proretry.Run(
    yourFunction,
    3,
    proretry.WithRetryableErrs(ErrTemporary, ErrTimeout),
)
```

### 使用自定义错误判定函数

```go
err := proretry.Run(
    yourFunction,
    3,
    proretry.WithRetryableErrFunc(func(err error) bool {
        return err != nil && err.Error() == "temporary error"
    }),
)
```

## API 参考

### 退避算法

- `ConstantBackoff(interval time.Duration) Backoff`
- `LinearBackoff(initInterval time.Duration) Backoff`
- `ExponentialBackoff(initInterval time.Duration) Backoff`
- `FibonacciBackoff(initInterval time.Duration) Backoff`

### 重试选项

- `WithInitInterval(interval time.Duration) RetryOption`
- `WithRetryableErrs(errs ...error) RetryOption`
- `WithRetryableErrFunc(f func(error) bool) RetryOption`
- `WithBackoff(backoff Backoff) RetryOption`

### 主要函数

- `Run(fn RetryFunc, maxRetries int, opts ...RetryOption) error`
 