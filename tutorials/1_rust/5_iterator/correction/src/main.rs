mod submission;

#[cfg(test)]
mod tests {
    use super::submission::process_names;

    #[test]
    fn test_process_names_basic() {
        let input = vec![
            "Jo".to_string(),
            "Anna".to_string(),
            "Robert".to_string(),
            "Li".to_string(),
        ];
        let expected = vec!["ANNA".to_string(), "ROBERT".to_string()];
        let result = process_names(input);
        assert_eq!(result, expected);
    }

    #[test]
    fn test_process_names_empty() {
        let input: Vec<String> = vec![];
        let expected: Vec<String> = vec![];
        let result = process_names(input);
        assert_eq!(result, expected);
    }

    #[test]
    fn test_process_names_all_short() {
        let input = vec!["Jo".to_string(), "Li".to_string(), "Sam".to_string()];
        let expected: Vec<String> = vec![];
        let result = process_names(input);
        assert_eq!(result, expected);
    }

    #[test]
    fn test_process_names_all_valid() {
        let input = vec!["Anna".to_string(), "Bethany".to_string()];
        let expected = vec!["ANNA".to_string(), "BETHANY".to_string()];
        let result = process_names(input);
        assert_eq!(result, expected);
    }
}
