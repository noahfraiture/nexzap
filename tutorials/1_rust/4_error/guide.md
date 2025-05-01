
# Rust Error Management Cheat Sheet

Rust’s error handling is explicit, using `Option` and `Result` to make code safe but sometimes verbose. This short tutorial covers the basics with quick explanations and tiny examples. Great for robust apps, less so for rapid scripting.

---

## `Option<T>` & `Result<T, E>`: Handling Absence and Errors

**What**: 
- `Option<T>`: `Some(T)` (value) or `None` (no value). For optional data.
- `Result<T, E>`: `Ok(T)` (success) or `Err(E)` (error). For operations that might fail.
**Why**: Safer than `null` or unchecked exceptions, but can feel wordy.

```rust
fn process(s: Option<&str>, n: i32) -> Result<usize, &str> {
    if n == 0 { Err("Zero!") } else { Ok(s.unwrap_or("").len() / n) }
}
println!("{:?}", process(Some("hi"), 2)); // Ok(1)
println!("{:?}", process(None, 0)); // Err("Zero!")
```

---

## `match`: Pattern Matching

**What**: A pattern-matching tool for handling `Option`/`Result` and more, ensuring all cases are covered.
**Why**: Precise control, but can be verbose.

```rust
fn check(n: Option<i32>) -> &str {
    match n {
        Some(x) if x > 0 => "Positive",
        _ => "Other",
    }
}
println!("{}", check(Some(5))); // Positive
```

---

## Methods: `unwrap`, `unwrap_or`, `map`

**What**: 
- `unwrap()`: Gets value or panics. Avoid!
- `unwrap_or(default)`: Gets value or default. Safe.
- `map(f)`: Transforms `Some`/`Ok`.
- There's many more methods on Option/Result
**Why**: Simplifies code, but `unwrap` is risky.

```rust
let n: Option<i32> = Some(2);
println!("{}", n.unwrap_or(0)); // 2
println!("{}", n.map(|x| x * 2).unwrap_or(0)); // 4
```

---

## `?` Operator: Error Propagation

**What**: Unwraps `Some`/`Ok` or returns `None`/`Err` early.
**Why**: Cuts boilerplate, but needs matching return types.

```rust
fn parse(s: &str) -> Result<i32, &str> {
    s.parse::<i32>().map_err(|_| "Invalid")?
}
println!("{:?}", parse("42")); // Ok(42)
```

---

## Tips
- **Pros**: Catches errors early, enforces safety.
- **Cons**: Verbose for quick scripts.
- **Best Practices**: Use `?` for propagation, skip `unwrap`, keep errors clear.

---

## Quick Exercise
1. Return first char of `Option<&str>` as `Option<char>`.
2. Parse string to `i32` with `Result` using `?`.
3. Double `Option<i32>` with `map`, default to 0.

Tame those errors and have fun!
