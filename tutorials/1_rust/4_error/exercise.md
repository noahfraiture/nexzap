## Task: Parse and Validate Numbers

In this exercise, you will use `Option` and `Result` with pattern matching and the `?` operator to parse and validate numbers from strings.

### Instructions

Write a function `parse_positive(s: &str) -> Result<i32, &str>` that parses a string into an `i32` and returns `Ok(n)` if the number is positive, or `Err("Not positive")` if non-positive, or `Err("Invalid number")` if parsing fails. Use the `?` operator for parsing. In `main`, test the function with inputs `"42"`, `"0"`, and `"abc"`.

#### Example:
- Input: `"42"`, `"0"`, `"abc"`
- Output: `Ok(42)`, `Err("Not positive")`, `Err("Invalid number")`

> **Note**: Ensure `fn main` is defined and use `println!` to display results.
