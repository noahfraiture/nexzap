# What is Go?

Go is a compiled programming language that produces a single binary, making deployment a breeze. It comes with the `go` tool to build projects, manage packages, or scaffold new ones. Go tries to be a new C++ like a lot of langauges.

## Key Features
- **Statically Typed with Type Inference**: Variables need a type at compile time, but Go’s compiler can infer them. Strict type checking catches errors before runtime.
- **Minimalist Syntax**: With ~25 keywords, Go avoids complex abstractions like classes, using structs and interfaces instead.
- **Concurrency with Goroutines**: Go’s goroutines—lightweight threads managed by the Go runtime—make concurrent programming simple. Channels help goroutines communicate safely.
- **Standard Library**: A robust set of packages for HTTP, file handling, and JSON parsing lets you build apps without reaching for external libraries.
- **Explicit Design**: No built-in `map`, `filter`, or `reduce`. You’ll write your own loops, keeping code clear but sometimes verbose.

## Pros
- **Simplicity**: Easy to learn and read, with a small syntax that keeps codebases consistent.
- **Compiled to Binary**: Go compiles to a single, efficient binary for fast execution and easy deployment across platforms.
- **Concurrency**: Goroutines are lightweight and efficient, making parallel tasks (like web servers) straightforward to implement.
- **Tooling**: The `go` tool automates formatting, testing, and dependency management for a smooth workflow.
- **Clear Error Handling**: Explicit error checks (e.g., `if err != nil`) make code predictable and robust.

## Cons
- **Limited Abstraction**: Features like interfaces are simple but can feel restrictive compared to languages with deeper type systems (e.g., Rust).
- **Verbose Error Handling**: Checking errors explicitly adds lines of code, which can feel repetitive.
- **No Built-In Conveniences**: Lack of functions like `map` or `filter` means writing more manual loops than in languages like JavaScript.
- **Basic Standard Library Scope**: While powerful, the standard library lacks some specialized tools (e.g., advanced data structures), requiring custom implementations.

## Why Learn Go?
Go is a pragmatic choice for building web servers, CLI tools, or distributed systems where simplicity and efficiency matter. It’s not a do-everything language like a bloated IDE, but more like a trusty hammer—perfect for certain nails, less so for others. If you value clarity and control, Go’s a solid pick.

## How to start ?

1. Download Go
2. Create a package with `go mod init main`
3. Create your first file of package `main` with a `func main() {}` function
