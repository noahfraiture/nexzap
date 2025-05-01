pub fn divide_length(input: Option<&str>, divisor: i32) -> Result<i32, String> {
    if divisor == 0 {
        return Err("Division by zero".to_string());
    }
    let length = input.unwrap_or("").len() as i32;
    Ok(length / divisor)
}
