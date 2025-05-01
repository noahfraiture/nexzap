pub fn append_and_count(s: &mut String, suffix: &str) -> usize {
    s.push_str(suffix);
    s.len()
}
