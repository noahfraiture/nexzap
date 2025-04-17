# Error Handling in Go

Go’s error handling is explicit, using `error` types and `if err != nil` checks, skipping exceptions. With `defer` for cleanup and rare `panic` for crashes, it’s clear but wordy. Let’s see it with a struct and mutex.

## Errors, Defer, and Mutex

Go returns `error` values to handle issues. Use `defer` for cleanup, like unlocking a mutex, and `panic` for unrecoverable failures.

**Example**:

```go
package main

import (
    "errors"
    "fmt"
    "sync"
)

type Counter struct {
    sync.Mutex
    Value int
}

func (c *Counter) Increment() (int, error) {
    c.Lock()
    defer c.Unlock()
    if c.Value >= 5 {
        return 0, errors.New("counter overflow")
    }
    c.Value++
    return c.Value, nil
}

func main() {
    counter := Counter{Value: 4}
    result, err := counter.Increment()
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Counter:", result)
    }
}
```

Run this to see “Counter: 5”. Try with `Value: 5` to get “Error: counter overflow”. The `Counter` uses a mutex, with `defer` ensuring unlock, and errors if `Value` hits 5.

## Other Tools

- `fmt.Errorf`: Dynamic error messages.
- `panic`: For crashes (use rarely).
- `errors.Is`: Check wrapped errors.

## Why Use It?

Go’s explicit errors and `defer` make concurrent code reliable, perfect for servers. It’s clear but repetitive `if err != nil` checks can slow down small scripts.
