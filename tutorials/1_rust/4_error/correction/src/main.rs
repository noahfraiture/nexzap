mod submission;

#[cfg(test)]
mod tests {
    use super::submission::divide_length;

    #[test]
    fn test_divide_length_basic() {
        let result = divide_length(Some("hello"), 2);
        assert_eq!(result, Ok(2)); // "hello".len() = 5, 5 / 2 = 2
    }

    #[test]
    fn test_divide_length_none_input() {
        let result = divide_length(None, 3);
        assert_eq!(result, Ok(0)); // None => length = 0, 0 / 3 = 0
    }

    #[test]
    fn test_divide_length_zero_divisor() {
        let result = divide_length(Some("hi"), 0);
        assert_eq!(result, Err("Division by zero".to_string())); // Divisor = 0 => error
    }

    #[test]
    fn test_divide_length_empty_string() {
        let result = divide_length(Some(""), 4);
        assert_eq!(result, Ok(0)); // "".len() = 0, 0 / 4 = 0
    }
}
