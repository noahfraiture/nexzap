# Concurrency in Go

Go’s concurrency is a dream With **goroutines** and **channels**, Go handles parallel tasks—like web requests or data crunching—smoothly. Let’s dive in.

## Goroutines: Lightweight Threads

A **goroutine** is a lightweight thread managed by the Go runtime, not the OS. You can launch thousands without sweating RAM. Start one with `go`.

**Example**:

```go
func worker(id int) {
    println("Worker", id, "done")
}

for range 1000 {
    go worker(i)
}
```

This spawns 1000 goroutines, each printing a message, showing how Go scales effortlessly.

## Channels: Coordinating Work

**Channels** enable safe data sharing between goroutines, preventing race conditions. Create with `make(chan Type)`, send with `ch<-`, and receive with `<-ch`.

**Example**:

```go
ch := make(chan int)
sum := 0
for i := range 3 {
    go func(n int) {
        ch <- n * n
    }(i)
}
for range 3 {
    sum += <-ch
}
println("Sum of squares:", sum)
```

This calculates squares concurrently, with channels collecting results safely. Getting the value from the channel is a blocking operation, so our main thread will wait on the line `sum += <-ch` until data have been pushed in the channel ch.
It is possible to make buffered channel `make(chan int, 3)` that will block a send only when the channel is full.

## Why Use Them?

Goroutines and channels make concurrency intuitive, ideal for scalable systems like APIs. They’re like a synchronized dance—graceful with practice. There's more possibility, like range loop on channel, wait group for goroutine or even traditionnal mutex.
