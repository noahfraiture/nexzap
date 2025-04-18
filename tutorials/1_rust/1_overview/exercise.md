## Task: Build a Simple CLI Counter

In this exercise, you will create a simple command-line program that uses Rust’s type system and `cargo` to count occurrences of a character in a string, demonstrating Rust’s safety and tooling.

### Instructions

Write a function `count_char(input: &str, target: char) -> usize` that takes a string slice and a character, then returns the number of times the character appears in the string. In `main`, prompt the user for a string and a character, then print the count. Use `std::io` for input and handle potential input errors gracefully.

#### Example:
- Input: `hello`, `l`
- Output: `The character 'l' appears 2 times in "hello".`

> **Note**: Use `cargo init` to set up the project, ensure `fn main` is defined, and import necessary modules.
