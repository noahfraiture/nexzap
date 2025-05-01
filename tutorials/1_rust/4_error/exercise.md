## Task: Safe String Length Divider

In this exercise, you’ll practice Rust’s error handling with `Option` and `Result` to safely compute a string’s length divided by a given number, using `String` for errors to keep things simple.

### Instructions

Write a function `divide_length(input: Option<&str>, divisor: i32) -> Result<i32, String>` that:
- Takes an `Option<&str>` (input string) and an `i32` (divisor).
- Returns `Ok(result)` with the string’s length divided by the divisor, using 0 as the length if `input` is `None`.
- Returns `Err("Division by zero".to_string())` if the divisor is 0.
- Use the `?` operator and at least one `Option`/`Result` method (e.g., `unwrap_or`).

#### Example:
- `divide_length(Some("hello"), 2)` → `Ok(2)`  // "hello".len() = 5, 5 / 2 = 2
- `divide_length(None, 3)` → `Ok(0)`  // default length = 0, 0 / 3 = 0
- `divide_length(Some("hi"), 0)` → `Err("Division by zero".to_string())`

> **Note**: Ensure safe handling of `Option` and `Result` without using `unwrap`.
