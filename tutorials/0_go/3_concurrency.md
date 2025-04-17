# Concurrency in Go

Go’s concurrency is a dream compared to C++’s threading nightmares, where you’re one mutex shy of chaos. With **goroutines** and **channels**, Go handles parallel tasks—like web requests or data crunching—smoothly. Let’s dive in.

## Goroutines: Lightweight Threads

A **goroutine** is a lightweight thread managed by the Go runtime, not the OS. You can launch thousands without sweating RAM. Start one with `go`.

**Example**:

```go
func worker(id int) {
    println("Worker", id, "done")
}

for i := 0; i < 1000; i++ {
    go worker(i)
}
```

This spawns 1000 goroutines, each printing a message, showing how Go scales effortlessly.

## Channels: Coordinating Work

**Channels** enable safe data sharing between goroutines, preventing race conditions. Create with `make(chan Type)`, send with `<-`, and receive with `<-`.

**Example**:

```go
ch := make(chan int)
sum := 0
for i := 1; i <= 3; i++ {
    go func(n int) {
        ch <- n * n
    }(i)
}
for i := 0; i < 3; i++ {
    sum += <-ch
}
println("Sum of squares:", sum)
```

This calculates squares concurrently, with channels collecting results safely. Getting the value from the channel is a blocking operation, so our main thread will wait on the line `sum += <-ch` until data have been pushed in the channel ch.

## Why Use Them?

Goroutines and channels make concurrency intuitive, ideal for scalable systems like APIs. They’re like a synchronized dance—graceful with practice. There's more possibility, like range loop on channel, wait group for goroutine or even traditionnal mutex.
