# Rust Iterators

Iterators are just one of Rust’s many features—not a core pillar, but I picked them because they’re super cool and fun to explore!

## Basics

- **Iterator**: Trait for producing a sequence of values via `next()`.
- **Types**: Created from collections (e.g., `Vec`, slices) or custom types.
- **Lazy**: Operations (e.g., `map`, `filter`) don’t execute until consumed (e.g., `collect`, `for`).

### Creating Iterators

```rust
let v = vec![1, 2, 3];
let iter = v.iter();         // Immutable: &i32
let mut_iter = v.iter_mut(); // Mutable: &mut i32
let into_iter = v.into_iter(); // Owned: i32
```

## Common Methods

- **Transform**: `map(f)`, `filter(p)`, `enumerate()`.
- **Consume**: `collect()`, `fold(init, f)`, `for_each(f)`.
- **Combine**: `chain(other)`, `zip(other)`.
- **Check**: `any(p)`, `all(p)`, `find(p)`.

### Examples

```rust
fn main() {
    let v = vec![1, 2, 3, 4];
    
    // Transform and collect
    let doubled: Vec<i32> = v.iter().map(|x| x * 2).collect(); // [2, 4, 6, 8]
    
    // Filter even numbers
    let evens: Vec<&i32> = v.iter().filter(|&&x| x % 2 == 0).collect(); // [2, 4]
    
    // Check if any number > 3
    let has_large = v.iter().any(|&x| x > 3); // true
    
    // Sum with fold
    let sum = v.iter().fold(0, |acc, &x| acc + x); // 10
    
    // Iterate with for
    for (i, &val) in v.iter().enumerate() {
        println!("Index {}: {}", i, val);
    }
}
```

## Tips

- Use `iter()` for references, `into_iter()` for ownership.
- Chain methods for concise operations (e.g., `map().filter().collect()`).
- Avoid unnecessary `collect()`; use `for` or consumers like `any` for efficiency.
