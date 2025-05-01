## Task: Manage String Ownership and Borrowing

In this exercise, you will practice Rust’s ownership and borrowing rules by manipulating a `String` without triggering compile-time errors.

### Instructions

Write a function `append_and_count(s: &mut String, suffix: &str) -> usize` that appends `suffix` to `s` and returns the new length of `s`.Ensure the original `String` remains usable after the function call.

#### Example:
- Input: `String::from("Hello")`, `" World"`
- Output: `Modified string: Hello World, Length: 11`

> **Note**: Use borrowing to avoid moving ownership.
