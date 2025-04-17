# Error Handling in Go

Go handles errors explicitly using return values, encouraging developers to check and manage errors directly rather than relying on exceptions. Let's explore error handling with the `errors` package and `error` type, along with its pros and cons.

## Working with Errors in Go

In Go, functions that can fail often return an `error` as the last return value. You check it with an `if` statement. The `errors` package helps create custom errors, and you can wrap errors with additional context using `fmt.Errorf`.

**Example**:

package main

import (
    "errors"
    "fmt"
    "os"
)

func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero is not allowed")
    }
    return a / b, nil
}

func main() {
    // Example of handling a custom error
    result, err := Divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)

    // Example of handling a system error
    content, err := os.ReadFile("example.txt")
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }
    fmt.Println("File content:", string(content))
}

Run this to see output like:
```
Error: division by zero is not allowed
```
Or, if the file operation fails:
```
Error reading file: open example.txt: no such file or directory
```

## Why Use Explicit Error Handling?

Go’s approach to errors forces you to handle potential failures explicitly, reducing unhandled exceptions and improving code reliability. It’s simple but requires more verbose code compared to try-catch mechanisms in other languages.
