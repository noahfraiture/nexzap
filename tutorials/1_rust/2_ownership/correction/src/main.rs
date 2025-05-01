mod submission;

#[cfg(test)]
mod tests {
    use super::submission::append_and_count;
    use std::string::String;

    #[test]
    fn test_append_and_count_basic() {
        let mut input_string = String::from("Hello");
        let suffix = " World";
        let new_length = append_and_count(&mut input_string, suffix);

        assert_eq!(new_length, 11); // "Hello World" has length 11
        assert_eq!(input_string, "Hello World"); // Ensure the string was modified
    }

    #[test]
    fn test_append_and_count_empty_string() {
        let mut input_string = String::new(); // Empty string
        let suffix = "Rust";
        let new_length = append_and_count(&mut input_string, suffix);

        assert_eq!(new_length, 4); // "Rust" has length 4
        assert_eq!(input_string, "Rust"); // Ensure the string was modified
    }

    #[test]
    fn test_append_and_count_empty_suffix() {
        let mut input_string = String::from("Hello");
        let suffix = ""; // Empty suffix
        let new_length = append_and_count(&mut input_string, suffix);

        assert_eq!(new_length, 5); // "Hello" still has length 5
        assert_eq!(input_string, "Hello"); // No change
    }
}
