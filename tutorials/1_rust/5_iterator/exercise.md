## Task: Filter and Transform a Number List

In this exercise, you will use Rust’s iterators to filter and transform a vector of numbers, leveraging lazy evaluation and common iterator methods.

### Instructions

Write a function `process_numbers(nums: &[i32]) -> Vec<i32>` that takes a slice of integers, filters out non-positive numbers, doubles the remaining numbers, and collects the result into a `Vec<i32>`. Use `filter` and `map` in a chained iterator expression. In `main`, test with `[1, -2, 3, 0, 4]`.

#### Example:

- Input: `[1, -2, 3, 0, 4]`
- Output: `[2, 6, 8]`

> **Note**: Ensure `fn main` is defined and print the result.
