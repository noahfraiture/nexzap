## Task: Calculate Sum with Goroutines and Channels

In this exercise, you will learn how to use goroutines and channels for concurrent programming in Go.

### Instructions

Write a function `CalculateSum(numbers []int) int` that will:
1. Take a slice of integers as input.
2. Split the work of summing the numbers into two goroutines if the slice has more than one element.
3. Use a channel to communicate partial sums from the goroutines and return the total sum.

#### Steps:
- Declare a channel for partial sums.
- Split the slice and use goroutines for summing each part (if length > 1).
- Return the combined sum from the channel results.

#### Note on Slices:
A slice in Go is a flexible view of an array. Use `len()` to get its length and `[:]` notation to create sub-slices (e.g., `numbers[:2]` for the first two elements, `numbers[2:]` for the rest).

#### Example:
- `CalculateSum([]int{1, 2, 3, 4})` should return `10`.
- `CalculateSum([]int{5})` should return `5`.
- `CalculateSum([]int{}`) should return `0`.

> **Note**: Do not forget to declare the package with `package main` at the top of the file.
