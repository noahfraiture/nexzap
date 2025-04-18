pub fn count_char(s: &str, c: char) -> usize {
    s.chars().filter(|i| *i == c).count()
}
