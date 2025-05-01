## Task: Filter and Transform Names

In this exercise, you’ll practice Rust’s iterator methods to process a vector of names, filtering and transforming them based on specific conditions.

### Instructions

Write a function `process_names(names: Vec<String>) -> Vec<String>` that:
- Takes a `Vec<String>` containing names.
- Uses iterators to:
  - Filter out names shorter than 4 characters.
  - Transform the remaining names to uppercase.
  - Collect the results into a new `Vec<String>`.
- Use at least `filter` and `map` iterator methods, and chain them for conciseness.

#### Example:
- Input: `vec!["Jo".to_string(), "Anna".to_string(), "Robert".to_string(), "Li".to_string()]`
- Output: `vec!["ANNA".to_string(), "ROBERT".to_string()]`

> **Note**: Use `iter()` to avoid taking ownership of the vector’s elements, and ensure the solution is efficient by chaining operations.

---

This exercise is small but requires thinking about:
- Iterator method chaining with `filter` and `map`.
- Working with `String` and string methods (e.g., `len()`, `to_uppercase()`).
- Using `iter()` for borrowing to preserve the input vector.
- Collecting results into the correct output type.

Give it a try, and if you want, share your solution or ask for feedback!
