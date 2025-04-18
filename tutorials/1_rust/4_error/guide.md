# Rust Error Management Cheat Sheet

## `Option<T>` & `Result<T, E>`
- **Option**: `Some(T)` (value) or `None` (no value). Use for optional values.
- **Result**: `Ok(T)` (success) or `Err(E)` (failure). Use for operations that may fail.

### Basic Usage & Pattern Matching
```rust
fn process(opt: Option<i32>, res: Result<i32, &str>) -> (i32, i32) {
    let opt_val = match opt {
        Some(v) => v,          // Extract value
        None => 0,             // Default
    };
    let res_val = match res {
        Ok(v) => v,            // Success
        Err(e) => {
            println!("Error: {}", e);
            0                  // Default
        }
    };
    (opt_val, res_val)
}

let x = Some(5);
let y: Result<i32, &str> = Ok(10);
let z: Result<i32, &str> = Err("Failed");
println!("{}", x.unwrap_or(0));     // 5
println!("{}", y.unwrap_or(0));     // 10
println!("{}", z.unwrap_or(0));     // 0
println!("{:?}", x.map(|n| n * 2)); // Some(10)
```

### Common Methods
- `unwrap()`: Get `T` or panic (`None`/`Err`).
- `unwrap_or(default)`: Get `T` or `default`.
- `map(f)`: Apply `f` to `Some`/`Ok`, skip `None`/`Err`.

## `?` Operator
- **Purpose**: Propagates `None`/`Err` early, unwraps `Some`/`Ok` to `T`.
- **Use**: In functions returning `Option`/`Result`.

### Example
```rust
fn parse_and_double(s: &str) -> Result<i32, &str> {
    let num = s.parse::<i32>().map_err(|_| "Invalid")?;
    Ok(num * 2)
}

fn first_positive(nums: &[i32]) -> Option<i32> {
    let first = nums.get(0)?;
    if *first > 0 { Some(*first) } else { None }
}

fn main() {
    println!("{:?}", parse_and_double("42"));   // Ok(84)
    println!("{:?}", parse_and_double("abc"));  // Err("Invalid")
    println!("{:?}", first_positive(&[5, 10])); // Some(5)
    println!("{:?}", first_positive(&[]));      // None
}
```

## Tips
- Use `match` for explicit `Option`/`Result` handling.
- Prefer `unwrap_or` or `?` over `unwrap` to avoid panics.
- Chain `map`, `unwrap_or` for concise code.
