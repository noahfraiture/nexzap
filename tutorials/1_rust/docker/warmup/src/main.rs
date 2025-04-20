mod submission;

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_count_char() {
        assert_eq!(submission::count_char("hello", 'l'), 2, "Should count 2 'l's in 'hello'");
        assert_eq!(submission::count_char("rust", 'x'), 0, "Should count 0 'x's in 'rust'");
        assert_eq!(submission::count_char("", 'a'), 0, "Should count 0 in empty string");
        assert_eq!(submission::count_char("aaa", 'a'), 3, "Should count 3 'a's in 'aaa'");
    }
}
