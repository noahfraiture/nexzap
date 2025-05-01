mod submission;

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_string_length() {
        assert_eq!(submission::string_length("hello"), 5, "Should return 5 for 'hello'");
        assert_eq!(submission::string_length("rust"), 4, "Should return 4 for 'rust'");
        assert_eq!(submission::string_length(""), 0, "Should return 0 for empty string");
        assert_eq!(submission::string_length("aaa"), 3, "Should return 3 for 'aaa'");
    }
}
